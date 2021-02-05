// MIT License
//
// (C) Copyright [2019-2021] Hewlett Packard Enterprise Development LP
//
// Permission is hereby granted, free of charge, to any person obtaining a
// copy of this software and associated documentation files (the "Software"),
// to deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included
// in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
// THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
// OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
// ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.
//
// This file contains the component control code for xname referenced
// components for capmcd
//

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"stash.us.cray.com/HMS/hms-capmc/internal/capmc"
)

// These are the HTTP handlers.
// They call the handler function after that with the command filled in.

// doGroupOn handles a power on request on group of components
func (d *CapmcD) doGroupOn(w http.ResponseWriter, r *http.Request) {
	d.doGroupOnOffCtrl(w, r, bmcCmdPowerOn)
}

// doGroupOff handles a power off request on group of components
func (d *CapmcD) doGroupOff(w http.ResponseWriter, r *http.Request) {
	d.doGroupOnOffCtrl(w, r, bmcCmdPowerOff)
}

// doGroupReinit handles a power reinit request on group of components
func (d *CapmcD) doGroupReinit(w http.ResponseWriter, r *http.Request) {
	d.doGroupOnOffCtrl(w, r, bmcCmdPowerRestart)
}

// doGroupStatus handles a get power status request on groups of components
func (d *CapmcD) doGroupStatus(w http.ResponseWriter, r *http.Request) {
	d.doGroupOnOffCtrl(w, r, bmcCmdPowerStatus)
}

// doGroupOnOffCtrl function looks like an HTTP handler API function that's
// registered Go's http server. It's not; it's called by the on/off/reinit/status
// handlers with the command filled in.
func (d *CapmcD) doGroupOnOffCtrl(w http.ResponseWriter, r *http.Request, command string) {

	if d.debug {
		log.Printf("Debug: doGroupOnOffCtrl command = %s\n", command)
	}

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		sendJsonError(w, http.StatusMethodNotAllowed,
			fmt.Sprintf("(%s) Not Allowed", r.Method))
		return
	}

	var args capmc.GroupControl
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&args)
	if err != nil {
		sendJsonError(w, http.StatusBadRequest, fmt.Sprintf("Bad Request: %s", err))
		return
	}

	if args.Groups == nil || len(args.Groups) == 0 {
		sendJsonError(w, http.StatusBadRequest, "Bad Request: Missing required Group list")
		return
	}

	// Only for get_group_status
	fmap, err := capmc.NodeStatusFilterParse(args.Filter)
	if err != nil {
		sendJsonError(w, http.StatusBadRequest, fmt.Sprintf("%s", err))
		return
	}

	if args.Force {
		command, err = d.getForceOption(command)
		if err != nil {
			sendJsonError(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	var (
		nl    []*NodeInfo
		query HSMQuery
	)

	// A role block prevents command from working on the component. The
	// HSM will do the filtering based on a negated role.
	roles, _ := cmdBlockRole(command)
	query.Roles = stringSliceMap(roles, func(s string) string {
		return "!" + s
	})

	query.Groups = args.Groups
	nl, err = d.GetNodesByGroup(query)
	if err != nil {
		var status int

		if _, ok := err.(*InvalidGroupsError); ok {
			status = http.StatusBadRequest
		} else {
			log.Printf("Error: %s", err)
			status = http.StatusInternalServerError
		}

		sendJsonError(w, status, err.Error())
		return
	}

	if nl == nil {
		sendJsonError(w, http.StatusNotFound, "No components in group")
		return
	}

	err = d.checkForDisabledComponents(nl, "xname")

	if err != nil {
		sendJsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	log.Printf("Info: Group power command: %s, groups: %v, reason: %s\n",
		command, args.Groups, args.Reason)

	if command == bmcCmdPowerStatus {
		statusData := d.doCompStatus(nl, command, fmap)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(statusData)
	} else {
		data := d.doCompOnOffCtrl(nl, command)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}

	return
}
