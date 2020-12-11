// Copyright 2017-2020 Hewlett Packard Enterprise Development LP
//
// This file contains the node control code for capmcd Redfish interface
//
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"

	base "stash.us.cray.com/HMS/hms-base"
	"stash.us.cray.com/HMS/hms-capmc/internal/capmc"
)

const unlimited int = -1

// defaultNodeRules can be overriden via config file.
// TODO Determine appropriate defaults. The values below are guesses.
var defaultNodeRules = PowerOpRules{
	MinOffTime: unlimited,
	MaxOffTime: unlimited,
	Off:        OpRule{Latency: 60, MaxReq: unlimited},
	On:         OpRule{Latency: 120, MaxReq: unlimited},
	Reinit:     OpRule{Latency: 180, MaxReq: unlimited},
}

// These are the HTTP handlers.
// They call the handler function after that with the command filled in.

// doNodeOn handles a node on request
func (d *CapmcD) doNodeOn(w http.ResponseWriter, r *http.Request) {
	d.doNodeOnOffCtrl(w, r, bmcCmdPowerOn)
}

// doNodeOff handles a node off request
func (d *CapmcD) doNodeOff(w http.ResponseWriter, r *http.Request) {
	d.doNodeOnOffCtrl(w, r, bmcCmdPowerOff)
}

// doNodeRestart handles a node reinit or restart request
func (d *CapmcD) doNodeRestart(w http.ResponseWriter, r *http.Request) {
	d.doNodeOnOffCtrl(w, r, bmcCmdPowerRestart)
}

// doNodeOnOffCtrl function looks like an HTTP handler API function that's
// registered Go's http server. It's not; it's called by the on/off/reinit
// handlers with the command filled in.
func (d *CapmcD) doNodeOnOffCtrl(w http.ResponseWriter, r *http.Request, command string) {

	if d.debug {
		log.Printf("Debug: doNodeOnOffCtrl command = %s\n", command)
	}

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		sendJsonError(w, http.StatusMethodNotAllowed,
			fmt.Sprintf("(%s) Not Allowed", r.Method))
		return
	}

	var args capmc.NodePowerRequest
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&args)
	if err != nil {
		sendJsonError(w, http.StatusBadRequest, fmt.Sprintf("Bad Request: %s", err))
		return
	}

	if args.Nids == nil || len(args.Nids) == 0 {
		sendJsonError(w, http.StatusBadRequest, "Bad Request: Missing required NID list")
		return
	}

	if args.Force {
		command, err = d.getForceOption(command)
		if err != nil {
			sendJsonError(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	var query HSMQuery

	// A role block prevents command from working on the component. The
	// HSM will do the filtering based on a negated role.
	roles, _ := cmdBlockRole(command)
	query.Roles = stringSliceMap(roles, func(s string) string {
		return "!" + s
	})

	// The incoming NID list could be invalid. Do simple validation
	// before contacting the Hardware State Manager.
	if len(args.Nids) > 0 {
		var bad []int

		query.NIDs, bad = validateNIDs(false, args.Nids)
		if len(bad) > 0 {
			sendJsonError(w, http.StatusBadRequest,
				fmt.Sprintf("invalid/duplicate nids: %v", bad))
			return
		}
	}

	var nl []*NodeInfo

	// TODO The next several statements should be factored out and shared
	// with the node on/off/restart|status calls.
	nl, err = d.GetNodesByNID(query)
	if err != nil {
		var status int

		if _, ok := err.(*InvalidNIDsError); ok {
			status = http.StatusBadRequest
		} else {
			log.Printf("Error: %s", err)
			status = http.StatusInternalServerError
		}

		sendJsonError(w, status, err.Error())
		return
	}

	if nl == nil {
		sendJsonError(w, http.StatusNotFound, "No nodes in machine")
		return
	}

	err = d.checkForDisabledComponents(nl, "nid")

	if err != nil {
		sendJsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	log.Printf("Info: Node power command: %s, nids: %v, reason: %s\n",
		command, args.Nids, args.Reason)

	var data capmc.NodePowerResponse
	data.Nids = make([]*capmc.NodePowerNidErr, 0, 1)

	var targetedXname []string
	for _, v := range nl {
		targetedXname = append(targetedXname, v.Hostname)
	}

	err = d.reserveComponents(targetedXname, command)
	defer d.releaseComponents(targetedXname)

	if err != nil {
		errstr := fmt.Sprintf("Error: Failed to reserve nodes while performing a %s.", command)
		log.Printf(errstr)
		data.ErrResponse.E = 37 // ENOLCK
		data.ErrResponse.ErrMsg = errstr
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
		return
	}

	waitNum, waitChan := d.queueBmcCmd(bmcCmd{cmd: command}, nl)

	var failures int
	for i := 0; i < waitNum; i++ {
		// Wait for each task to complete.
		result := <-waitChan
		// If any BMC/Node op fails then record overall error code for
		// the call. The HTTP code returned will still be 200.
		if result.rc != 0 {
			failures++
			// TODO Look at using the stack instead of heap????
			// The original API doesn't include a nid list
			// identifying failed operations this does.
			nidErr := new(capmc.NodePowerNidErr)
			nidErr.Nid = result.ni.Nid
			nidErr.ErrResponse.E = result.rc
			nidErr.ErrResponse.ErrMsg = result.msg
			data.Nids = append(data.Nids, nidErr)
		}
	}

	if failures > 0 {
		data.ErrResponse.E = -1
		data.ErrResponse.ErrMsg = fmt.Sprintf("Errors encountered with %d/%d NIDs",
			failures, waitNum)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
	return
}

// doNodeStatus handles a node status request
// The original API didn't have an "indeterminate" or "error" states. It's
// likely that with Redfish backend handling we'll encounter errors
// trying to talk to the hardware, so we need an extra state.
func (d *CapmcD) doNodeStatus(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		sendJsonError(w, http.StatusMethodNotAllowed,
			fmt.Sprintf("(%s) Not Allowed", r.Method))
		return
	}

	var args capmc.NodeStatusRequest
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&args)
	if err != nil {
		sendJsonError(w, http.StatusBadRequest,
			fmt.Sprintf("Bad Request: JSON: %s", err))
		return
	}

	filter, err := capmc.NodeStatusFilterParse(args.Filter)
	if err != nil {
		sendJsonError(w, http.StatusBadRequest, fmt.Sprintf("%s", err))
		return
	}

	if args.Source == "" {
		log.Printf("Info: no status source specified using HSM")
		args.Source = "HSM"
	}

	var useRedfish bool
	switch strings.ToLower(args.Source) {
	case "hsm", "hms", "sm", "smd", "software":
		// default
	case "redfish", "hardware":
		useRedfish = true
	default:
		sendJsonError(w, http.StatusBadRequest,
			fmt.Sprintf("unknown status source '%s'", args.Source))
		return
	}

	// Always ignore Nodes with these States
	var query = HSMQuery{
		States: []string{
			"!" + string(base.StateEmpty),
			"!" + string(base.StatePopulated),
			"!" + string(base.StateUnknown),
		},
	}

	// The incoming NID list could be invalid. Do simple validation
	// before contacting the Hardware State Manager.
	if len(args.Nids) > 0 {
		var invalidNIDs []int

		// Duplicate NIDs are ok in the Cascade CAPMC API.
		query.NIDs, invalidNIDs = validateNIDs(true, args.Nids)
		if len(invalidNIDs) > 0 {
			sendJsonError(w, http.StatusBadRequest,
				fmt.Sprintf("invalid nids: %v", invalidNIDs))
			return
		}
	}

	// set for use in log message
	if len(args.Filter) < 1 {
		args.Filter = capmc.FilterShowAll + " (implied)"
	}

	log.Printf("Info: Node power command: status, filter: %s, NIDs: %v\n",
		args.Filter, args.Nids)

	var data capmc.NodeStatusResponse

	if useRedfish {
		var nl []*NodeInfo

		nl, err = d.GetNodesByNID(query)
		if err != nil {
			var status int

			if _, ok := err.(*InvalidNIDsError); ok {
				status = http.StatusBadRequest
			} else {
				log.Printf("Error: %s\n", err)
				status = http.StatusInternalServerError
			}

			sendJsonError(w, status, err.Error())
			return
		}

		if nl == nil {
			sendJsonError(w, http.StatusNotFound, "No nodes in machine")
			return
		}

		cidToNidMap := make(map[string]int)
		for _, node := range nl {
			cidToNidMap[node.Hostname] = node.Nid
		}

		cdata := d.doCompStatus(nl, bmcCmdPowerStatus, filter)

		// The JSON encoder omits empty lists. The Cascade CAPMC API
		// response contains more lists than this, but these are the
		// only ones CAPMC can report using Redfish.
		data.On = make([]int, 0, 1)
		data.Off = make([]int, 0, 1)
		data.Undefined = make([]int, 0, 1)

		// NOTE it is safe to ignore error here as it isn't possible
		// for the component Id (node) list used to fetch power state
		// to not contain all of the possible NIDs.
		data.Off, _ = compIdsToNids(cdata.Off, cidToNidMap)
		data.On, _ = compIdsToNids(cdata.On, cidToNidMap)
		data.Undefined, _ = compIdsToNids(cdata.Undefined, cidToNidMap)

		data.ErrResponse.E = cdata.ErrResponse.E
		// Convert any error message to NID based
		data.ErrResponse.ErrMsg =
			strings.Replace(data.ErrResponse.ErrMsg, "Xnames", "NIDs", 1)
	} else {
		data, err = d.GetNidStatus(query, filter)
		if err != nil {
			var status int

			if _, ok := err.(*InvalidNIDsError); ok {
				status = http.StatusBadRequest
			} else {
				log.Printf("Error: %s\n", err)
				status = http.StatusInternalServerError
			}

			sendJsonError(w, status, err.Error())
			return
		}
	}

	// Sorting of integer lists is mostly for convenience. It's not
	// strictly needed for capmc.
	sort.Ints(data.Disabled)
	sort.Ints(data.Halt)
	sort.Ints(data.Off)
	sort.Ints(data.On)
	sort.Ints(data.Ready)
	sort.Ints(data.Standby)
	sort.Ints(data.Undefined)
	if data.Flags != nil {
		sort.Ints(data.Flags.Alert)
		sort.Ints(data.Flags.Reserved)
		sort.Ints(data.Flags.Warning)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
	return
}

// doNodeRules handles a node rules request
// This never was a very popular API, so these values are ok for now. The
// original API call expected an empty JSON object to be posted.  This
// call accepts GET and is backward compatible with the old POST behavior.
func (d *CapmcD) doNodeRules(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		w.Header().Set("Allow", "GET, POST")
		sendJsonError(w, http.StatusMethodNotAllowed,
			fmt.Sprintf("(%s) Not Allowed", r.Method))
		return
	}

	switch r.Method {
	case http.MethodGet:
		// NOTE Finding no good examples or documentation on what
		// should happen when a GET has a body we'll follow the
		// maxim "be liberal in what you accept, and conservative
		// in what you send".  Log it and move on.
		if len(r.TransferEncoding) > 0 || r.ContentLength > 0 {
			log.Printf("Unexpected: %s %s with body",
				r.Method, r.URL.Path)
		}

		// There are no supported query parameters at this time but
		// following the above "be liberal in what you accept, and
		// conservative in what you send", we'll only error out when
		// the request is malformed. Log it and move on.
		err := r.ParseForm()
		if err != nil {
			sendJsonError(w, http.StatusBadRequest,
				fmt.Sprintf("Bad Request: %s", err))
			return
		}

		if len(r.Form) > 0 {
			log.Printf("Unexpected: %s %s with query %s",
				r.Method, r.URL.Path, r.URL.RawQuery)
		}
	case http.MethodPost:
		// NOTE Finding no good examples or documentation on what
		// should happen when a POST has query parameters we'll
		// follow the maxim "be liberal in what you accept, and
		// conservative in what you send".  Log it and move on.
		m, err := url.ParseQuery(r.URL.RawQuery)
		if len(m) > 0 {
			log.Printf("Unexpected: %s %s with query %s",
				r.Method, r.URL.Path, r.URL.RawQuery)
		}

		// There are no supported body parameters at this time but
		// following the above "be liberal in what you accept, and
		// conservative in what you send", we'll only error out when
		// the request is malformed. Log it and move on.
		var v interface{}
		err = json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			sendJsonError(w, http.StatusBadRequest,
				fmt.Sprintf("Bad Request: JSON: %s", err))
			return
		}
	}

	rules := capmc.GetNodeRulesResponse{
		MaxOffTime:        d.config.NodeRules.MaxOffTime,
		MinOffTime:        d.config.NodeRules.MinOffTime,
		LatencyNodeOff:    d.config.NodeRules.Off.Latency,
		LatencyNodeOn:     d.config.NodeRules.On.Latency,
		LatencyNodeReinit: d.config.NodeRules.Reinit.Latency,
		MaxOffReqCount:    d.config.NodeRules.Off.MaxReq,
		MaxOnReqCount:     d.config.NodeRules.On.MaxReq,
		MaxReinitReqCount: d.config.NodeRules.Reinit.MaxReq,
		ErrResponse:       capmc.ErrResponse{E: 0, ErrMsg: ""},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rules)
}
