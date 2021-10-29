/*
 * MIT License
 *
 * (C) Copyright [2019-2021] Hewlett Packard Enterprise Development LP
 *
 * Permission is hereby granted, free of charge, to any person obtaining a
 * copy of this software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation
 * the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the
 * Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included
 * in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
 * THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
 * OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
 * ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
 * OTHER DEALINGS IN THE SOFTWARE.
 */

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	base "github.com/Cray-HPE/hms-base"
	"github.com/Cray-HPE/hms-capmc/internal/capmc"
)

// getCabinetType determines if the xname provided exists in a Mountain, Hill,
// or River cabinet.
func (d *CapmcD) getCabinetType(xname string) base.HMSClass {
	cType := base.HMSClass("Unknown")

	params := url.Values{}
	params.Add("id", xname)
	components, err := d.GetComponents(params.Encode())

	if len(components) == 1 && err == nil {
		for _, comp := range components {
			cType = base.HMSClass(comp.Class)
		}
	}

	return cType
}

// handleDependentComponents takes a list of xnames and generates a new list
// of xnames that includes or excludes components based on the power action
// being requested.
// Rosetta: The Rosetta must be powered off before the RouterModule or iPDU
//          outlet are turned off to prevent possible hardware damage. The
//          switches are set to auto power on the Rosettas when the RouterModule
//          or iPDU outlet are turned on.
func (d *CapmcD) handleDependentComponents(xnames []string, cmd string) []string {
	var newList []string
	xmap := make(map[string]bool)

	for _, xname := range xnames {
		xmap[xname] = true
		hmsType := base.GetHMSType(xname)
		switch hmsType {
		case base.HSNBoard:
			if cmd == bmcCmdPowerOn || cmd == bmcCmdPowerForceOn {
				cabType := d.getCabinetType(xname)
				switch cabType {
				case base.ClassHill, base.ClassMountain:
					// Remove the Rosetta from the target list
					if _, ok := xmap[base.GetHMSCompParent(xname)]; ok {
						delete(xmap, xname)
					}
				case base.ClassRiver:
					// TODO: CASMHMS-4223
				default:
					log.Printf("Notice: Could not determine the cabinet "+
						"type of %s\n", xname)
				}
			}
		case base.RouterModule:
			if cmd == bmcCmdPowerOff || cmd == bmcCmdPowerForceOff {
				// Add the Rosetta to the target list
				xmap[xname+"e0"] = true
			} else if cmd == bmcCmdPowerOn || cmd == bmcCmdPowerForceOn {
				// Remove the Rosetta from the target list
				if _, ok := xmap[xname+"e0"]; ok {
					delete(xmap, xname+"e0")
				}
			}
		case base.CabinetPDUPowerConnector, base.CabinetPDUOutlet:
			// TODO: CASMHMS-4223
		}
	}

	for xname := range xmap {
		newList = append(newList, xname)
	}

	return newList
}

// GenerateXnameDescendantList takes a list of xnames and generates a new list
// of xnames that includes all of the descendants of the original list of
// xnames. Any xnames in the original list that are not found in HSM are put
// onto a separate list for use later.
func (d *CapmcD) GenerateXnameDescendantList(query HSMQuery) ([]string, error) {
	var (
		newList []string
		badList []string
		squery  HSMQuery = HSMQuery{
			Roles:  query.Roles,
			States: query.States,
			// TODO Maybe these should come from the config file
			Types: []string{
				"chassis",
				"cabinetpdu",
				"routermodule",
				"hsnboard",
				"computemodule",
				"node",
			},
		}
	)

	xmap := make(map[string]bool)
	for _, xname := range query.ComponentIDs {
		_, inMap := xmap[xname]
		if inMap == false {
			hmsType := base.GetHMSType(xname)
			if hmsType == base.Node ||
				hmsType == base.CabinetPDUOutlet ||
				hmsType == base.CabinetPDUPowerConnector ||
				hmsType == base.HSNBoard {
				xmap[xname] = true
			} else {
				// Retrieve a list of component structures
				// including for the xname itself on Mountain.
				components, err := d.GetComponentsQuery(xname,
					getRestrictStr(squery))
				if err != nil {
					return nil, err
				}

				if len(components) > 0 {
					for _, comp := range components {
						xmap[comp.ID] = true
					}
				} else {
					badList = append(badList, xname)
				}
			}
		}
	}

	for xname := range xmap {
		newList = append(newList, xname)
	}

	if len(badList) > 0 {
		return newList,
			&InvalidCompIDsError{"invalid xnames", badList}
	}

	return newList, nil
}

// GenerateXnamePrereqList takes a list of xnames and generates a new list
// of xnames that includes all of the ancestors of the original list of
// xnames. Any xnames in the original list that are not found in HSM are put
// onto a separate list for use later.
func (d *CapmcD) GenerateXnamePrereqList(query HSMQuery) ([]string, error) {
	var (
		newList []string
		badList []string
		squery  HSMQuery = HSMQuery{Roles: query.Roles}
	)

	xmap := make(map[string]bool)
	compList, err := d.GetComponents(getRestrictStr(squery))
	if err != nil {
		return nil, err
	}

	compMap := make(map[string]bool)
	for _, comp := range compList {
		compMap[comp.ID] = true
	}

	for _, xname := range query.ComponentIDs {
		if _, valid := compMap[xname]; valid {
			switch base.GetHMSType(xname) {
			case base.Node:
				// Need to determine which iPDU sockets this node is plugged into
				pduS1 := ""
				pduS2 := ""
				if pduS1 != "" {
					xmap[pduS1] = true
				}
				if pduS2 != "" {
					xmap[pduS2] = true
				}
				xmap[xname] = true
				// Strip the node and BMC field
				xname = base.GetHMSCompParent(xname)
				xname = base.GetHMSCompParent(xname)
				fallthrough
			case base.HSNBoard:
				if _, valid := compMap[xname]; valid {
					xmap[xname] = true
				}
				// Strip the enclosure field
				xname = base.GetHMSCompParent(xname)
			case base.RouterModule, base.ComputeModule:
				if _, valid := compMap[xname]; valid {
					xmap[xname] = true
				}
				// Strip the compute/switch module field
				xname = base.GetHMSCompParent(xname)
				fallthrough
			case base.Chassis:
				if _, valid := compMap[xname]; valid {
					xmap[xname] = true
				}
			default:
				badList = append(badList, xname)
			}
		} else {
			badList = append(badList, xname)
		}
	}

	for xname := range xmap {
		newList = append(newList, xname)
	}

	if len(badList) > 0 {
		return newList,
			&InvalidCompIDsError{"Invalid Component IDs", badList}
	}

	return newList, nil
}

// These are the HTTP handlers.
// They call the handler function after that with the command filled in.

// doXnameOn handles a node on request
func (d *CapmcD) doXnameOn(w http.ResponseWriter, r *http.Request) {
	d.doXnameOnOffCtrl(w, r, bmcCmdPowerOn)
}

// doXnameOff handles a node off request
func (d *CapmcD) doXnameOff(w http.ResponseWriter, r *http.Request) {
	d.doXnameOnOffCtrl(w, r, bmcCmdPowerOff)
}

// doXnameReinit handles a node reinit request
func (d *CapmcD) doXnameReinit(w http.ResponseWriter, r *http.Request) {
	d.doXnameOnOffCtrl(w, r, bmcCmdPowerRestart)
}

// doXnameStatus handles a status request for a component reference via an xname
func (d *CapmcD) doXnameStatus(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		sendJsonError(w, http.StatusMethodNotAllowed,
			fmt.Sprintf("(%s) Not Allowed", r.Method))
		return
	}

	var args capmc.XnameStatusRequest
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&args)
	if err != nil {
		if err == io.EOF {
			sendJsonError(w, http.StatusBadRequest, "no request")
		} else {
			sendJsonError(w, http.StatusBadRequest, err.Error())
		}

		return
	}

	filter, err := capmc.StatusFilterParse(args.Filter)
	if err != nil {
		sendJsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	if args.Source == "" {
		log.Printf("Info: no status source specified using Redfish")
		args.Source = "redfish"
	}

	var useHSM bool
	switch strings.ToLower(args.Source) {
	case "hsm", "hms", "sm", "smd", "software":
		useHSM = true
	case "redfish", "hardware":
		// default
	default:
		sendJsonError(w, http.StatusBadRequest,
			fmt.Sprintf("unknown status source '%s'", args.Source))
		return
	}

	var query HSMQuery

	if len(args.Xnames) > 0 {
		var bad []string

		query.ComponentIDs, bad = base.ValidateCompIDs(args.Xnames, false)
		if len(bad) > 0 {
			sendJsonError(w, http.StatusBadRequest,
				fmt.Sprintf("invalid/duplicate xnames: %v", bad))
			return
		}
	}

	// Only query hardware that CAPMC can control power on. If in the future
	// a new piece of hardware is available for power control, it will need
	// to be added/updated in the configmap for On. If we cannot find the power
	// controls for On, fallback to returning everything.
	pc, ok := d.config.PowerControls["On"]
	if ok {
		vt := pc.CompSeq

		for _, t := range vt {
			query.Types = append(query.Types, t)
		}
	}

	// By default show only enabled components. Asking for disabled
	// will only show disabled components that match the other flags.
	// The different States are OR'd together. The Enabled flag is
	// AND'd with the States.
	if len(args.Filter) > 0 {
		if filter&capmc.FilterShowHaltBit != 0 {
			query.States = append(query.States, "Halt")
		}
		// When using Redfish we always want the actual power states, not the
		// ones from HSM.
		if filter&capmc.FilterShowOffBit != 0 || useHSM == false {
			query.States = append(query.States, "Off")
		}
		if filter&capmc.FilterShowOnBit != 0 || useHSM == false {
			query.States = append(query.States, "On")
		}
		if filter&capmc.FilterShowReadyBit != 0 {
			query.States = append(query.States, "Ready")
		}
		if filter&capmc.FilterShowStandbyBit != 0 {
			query.States = append(query.States, "Standby")
		}
		// Need to set the enabled bit if show_disabled wasn't used
		if filter&capmc.FilterShowDisabledBit == 0 {
			query.Enabled = append(query.Enabled, true)
		} else {
			// Only query on false if we are not using show_all
			if filter != capmc.FilterShowAllBit {
				query.Enabled = append(query.Enabled, false)
			}
		}
	} else {
		// The default, set it for use in log messages
		args.Filter = capmc.FilterShowAll + " (implied)"
	}

	var compIDStr string
	if len(args.Xnames) > 0 {
		compIDStr = fmt.Sprintf("%v", args.Xnames)
	} else {
		compIDStr = "[all]"
	}

	log.Printf("Info: Xname power command: status, filter: %s, xnames: %v\n",
		args.Filter, compIDStr)

	var data capmc.XnameStatusResponse

	if useHSM {
		data, err = d.GetComponentStatus(query, filter)
		if err != nil {
			var status int

			if _, ok := err.(*InvalidCompIDsError); ok {
				status = http.StatusBadRequest
			} else {
				log.Printf("Error: %s\n", err)
				status = http.StatusInternalServerError
			}

			sendJsonError(w, status, err.Error())
			return
		}
	} else {
		nl, err := d.GetNodesByXname(query)

		if err != nil {
			var status int

			if _, ok := err.(*InvalidCompIDsError); ok {
				status = http.StatusBadRequest
			} else {
				log.Printf("Error: %s", err)
				status = http.StatusInternalServerError
			}

			sendJsonError(w, status, err.Error())
			return
		}

		if nl == nil {
			sendJsonError(w, http.StatusNotFound,
				"No matching components found")
			return
		}

		// Since the hardware power state is either On or Off, be sure at least
		// one of those filters is set so the query returns meaningful data.
		if filter&capmc.FilterShowOffBit == 0 &&
			filter&capmc.FilterShowOnBit == 0 {
			filter |= capmc.FilterShowOffBit
			filter |= capmc.FilterShowOnBit
		}
		data = d.doCompStatus(nl, bmcCmdPowerStatus, filter)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
	return
}

// doXnameOnOffCtrl function looks like an HTTP handler API function that's
// registered Go's http server. It's not; it's called by the on/off/reinit
// handlers with the command filled in.
func (d *CapmcD) doXnameOnOffCtrl(w http.ResponseWriter, r *http.Request, command string) {

	if d.debug {
		log.Printf("Debug: doNodeOnOffCtrl command = %s\n", command)
	}

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		sendJsonError(w, http.StatusMethodNotAllowed,
			fmt.Sprintf("(%s) Not Allowed", r.Method))
		return
	}

	var args capmc.XnameControl
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&args)
	if err != nil {
		sendJsonError(w, http.StatusBadRequest, fmt.Sprintf("Bad Request: %s", err))
		return
	}

	if args.Recurse && args.Prereq {
		sendJsonError(w, http.StatusBadRequest,
			"Bad Request: recursive and prereq options are mutually exclusive")
		return
	}

	if args.Xnames == nil {
		sendJsonError(w, http.StatusBadRequest, "Bad Request: Missing required xnames parameter")
		return
	}

	if len(args.Xnames) == 0 {
		sendJsonError(w, http.StatusBadRequest, "Bad Request: Required xnames list is empty")
		return
	}

	if args.Force {
		command, err = d.getForceOption(command)
		if err != nil {
			sendJsonError(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	var eData capmc.XnameControlResponse
	handleErr := func(err error) bool {
		var (
			compIDError *InvalidCompIDsError
			status      int
		)

		if errors.As(err, &compIDError) {
			eData.Xnames = append(eData.Xnames,
				MakeXnameErrors(compIDError)...)

			if args.Continue {
				log.Printf("Notice: ignoring bad component ids %s\n",
					compIDError.Error())
			} else {
				status = http.StatusBadRequest
				eData.ErrResponse = capmc.ErrResponseEINVAL
			}
		} else {
			log.Printf("Error: %s\n", err)
			status = http.StatusInternalServerError
			eData.ErrResponse.E = status
			eData.ErrResponse.ErrMsg = err.Error()
			eData.Xnames = eData.Xnames[:0]
		}

		if status != 0 {
			SendResponseJSON(w, status, eData)

			return false
		}

		return true
	}

	duplicatesOK := args.Force || args.Continue
	xnames, badXnames := base.ValidateCompIDs(args.Xnames, duplicatesOK)
	if len(badXnames) > 0 {
		continueCmd := handleErr(&InvalidCompIDsError{
			err:     "invalid/duplicate xnames",
			CompIDs: badXnames,
		})
		if !continueCmd {
			return
		}
	}

	// Some components need special cases to prevent errors and failures
	xnames = d.handleDependentComponents(xnames, command)

	// Check the list for the special cases of s0 and all. If we find one,
	// simply use an empty list. This will indicate that we want everything.
	for _, v := range xnames {
		if v == "s0" || v == "all" {
			xnames = xnames[:0]
			break
		}
	}

	var query HSMQuery

	// A role block prevents command from working on the component. The
	// HSM will do the filtering based on a negated role.
	roles, _ := d.cmdBlockRole(command)
	query.Roles = stringSliceMap(roles, func(s string) string {
		return "!" + s
	})

	if len(xnames) > 0 {
		squery := HSMQuery{
			ComponentIDs: xnames,
			Roles:        query.Roles,
			States:       []string{"!Empty"},
		}

		switch {
		case args.Recurse:
			query.ComponentIDs, err = d.GenerateXnameDescendantList(squery)
		case args.Prereq:
			query.ComponentIDs, err = d.GenerateXnamePrereqList(squery)
		default:
			query.ComponentIDs = xnames
			err = nil
		}

		if err != nil {
			continueCmd := handleErr(err)
			if !continueCmd {
				return
			}
		}
	}

	if args.Continue {
		query.Enabled = []bool{true}
		query.States = append(query.States, "!Empty")
	}

	var nl []*NodeInfo
	nl, err = d.GetNodesByXname(query)
	if err != nil {
		continueCmd := handleErr(err)
		if !continueCmd {
			return
		}
	}

	if nl == nil || len(nl) == 0 {
		sendJsonError(w, http.StatusNotFound, "No nodes found to operate on")
		return
	}

	err = d.checkForDisabledComponents(nl, "xname")

	if err != nil {
		sendJsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	log.Printf("Info: Xname power command: %s, xnames: %v, reason: %s\n",
		command, xnames, args.Reason)

	data := d.doCompOnOffCtrl(nl, command)

	// add accumulated ignored errors if any
	if len(eData.Xnames) > 0 {
		data.Xnames = append(data.Xnames, eData.Xnames...)
		if len(data.ErrResponse.ErrMsg) > 0 {
			data.ErrResponse.ErrMsg += "; "
		}
		data.ErrResponse.ErrMsg += fmt.Sprintf("Errors encountered with %d/%d Xnames for %s",
			len(data.Xnames), len(args.Xnames), command)
	}

	SendResponseJSON(w, http.StatusOK, data)

	return
}

// MakeXnameErrors returns a structure of E, ErrMsg, and xnames array
func MakeXnameErrors(err *InvalidCompIDsError) []*capmc.XnameControlErr {
	// massage the message
	msg := strings.ReplaceAll(err.err, "xnames", "xname")

	var xnames []*capmc.XnameControlErr

	for _, xname := range err.CompIDs {
		xnames = append(xnames, capmc.MakeXnameError(xname, 22, msg))
	}

	return xnames
}
