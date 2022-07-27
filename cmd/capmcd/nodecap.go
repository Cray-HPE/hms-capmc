//
// MIT License
//
// (C) Copyright [2019-2022] Hewlett Packard Enterprise Development LP
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

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"path"
	"strings"

	base "github.com/Cray-HPE/hms-base"
	"github.com/Cray-HPE/hms-capmc/internal/capmc"
	"github.com/Cray-HPE/hms-smd/pkg/sm"
)

// List of seen PowerControl.Name values. This may not be the best long term.
// TODO Revist this in light of new information learned.
const (
	// Proposed Cray nC Accelerator PowerControl.Name
	CrayAccelPCName = "Accelerator0 Power Control"
	// Cray nC Node PowerControl.Name
	CrayNodePCName = "Node Power Control"
	// BMC Node PowerControl.Name (based on DMTF reference implementation)
	GenericNodePCName = "System Power Control"
	// GIGA-BYTE BMC Node PowerControl.Name
	GigaByteNodePCName = "Chassis Power Control"
	// Intel BMC Node PowerControl.Name
	IntelNodePCName = "Server Power Control"
	HPENodePCName   = "HPE Power Control"
	HPEApolloPCName = "PowerLimit Resource for AccPowerService"
)

// newPowerCapNidError creates a new PowerCapNid structure initialized as
// an error response for the get/set/capabilities power cap APIs.
func newPowerCapNidError(nid, ecode int, emsg string) capmc.PowerCapNid {
	return capmc.PowerCapNid{Nid: nid, E: ecode, ErrMsg: emsg}
}

// doPowerCapCapabilities is the HTTP handler for the
// get_power_cap_capabilities API
func (d *CapmcD) doPowerCapCapabilities(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		sendJsonError(w, http.StatusMethodNotAllowed,
			fmt.Sprintf("(%s) Not Allowed", r.Method))
		return
	}

	var args capmc.NidlistRequest
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&args)
	if err != nil {
		sendJsonError(w, http.StatusBadRequest,
			fmt.Sprintf("Bad Request: JSON: %s", err))
		return
	}

	log.Printf("Info: CAPMC Get Power Cap Capabilities - %v", args.Nids)

	// The incoming NID list could be invalid. Do simple validation
	// before contacting Hardware State Manager.
	var niQuery HSMQuery
	if len(args.Nids) > 0 {
		var invalidNIDs []int
		// Duplicate NIDs aren't an error in the Cascade CAPMC API.
		niQuery.NIDs, invalidNIDs = validateNIDs(true, args.Nids)
		if len(invalidNIDs) > 0 {
			sendJsonError(w, http.StatusBadRequest,
				fmt.Sprintf("invalid nids: %v", invalidNIDs))
			return
		}
	}

	// Get requested NidInfos from HSM
	// the GetNidInfo function will properly update the HSMQuery Types array when the NIDs list is empty
	var nidInfoList []*capmc.NidInfo
	nidInfoList, err = d.GetNidInfo(niQuery)
	if err != nil {
		if nErr, ok := err.(*InvalidNIDsError); ok {
			// format a json response to the incorrect nids error and stop processing
			var data capmc.GetPowerCapCapabilitiesResponse
			data.E = 22 // EINVAL
			data.ErrMsg = "Invalid argument"
			for _, nid := range nErr.NIDs {
				data.Nids = append(data.Nids,
					newPowerCapNidError(nid, 22, "Undefined NID"))
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			err = json.NewEncoder(w).Encode(data)
			return
		}

		// unknown error - just bail
		status := http.StatusInternalServerError
		log.Printf("Error: CAPMC Get Power Cap Capabilities: %s\n", err.Error())
		sendJsonError(w, status, err.Error())
		return
	}

	//create an xname to nid lookup to be help build the NID list for each PowerCapGroup later
	xnameNidInfoLookup := make(map[string]*capmc.NidInfo)
	for _, nidInfo := range nidInfoList {
		xnameNidInfoLookup[nidInfo.Xname] = nidInfo
	}

	// Get requested ComponentEndpoints from HSM
	// ComponentEndpoints contain the PowerControl data structures that describe the capabilities
	var ceQuery HSMQuery
	if len(args.Nids) > 0 {
		ceQuery.NIDs = args.Nids
	} else {
		ceQuery.Types = append(ceQuery.Types, "node")
	}
	restrict := getRestrictStr(ceQuery)

	var componentEndpoints []*sm.ComponentEndpoint
	componentEndpoints, err = d.GetComponentEndpoints(restrict)
	if err != nil {
		log.Printf("Error: CAPMC Get Power Cap Capabilities: %s\n", err.Error())
		sendJsonError(w, http.StatusBadRequest,
			fmt.Sprintf("Bad Request: JSON: %s", err))
		return
	}

	//create an xname to componentEndpoint lookup to be used when querying for ComponentEndpoint PowerControl data
	//this will eliminate the need to do multiple linear searches later
	xnameComponentLookup := make(map[string]*sm.ComponentEndpoint)
	for _, componentEndpoint := range componentEndpoints {
		xnameComponentLookup[componentEndpoint.ID] = componentEndpoint
	}

	//Get the full SystemHWInventory
	//At this time, only the “all” and “s0” xname values are accepted by HSM Hardware Inventory queries
	//The HSM plan is to expand this to accept any xname, group, or partition as described in CASMHMS-2285.
	hwInventory, err := d.GetHWInventoryQuery("all")
	if err != nil {
		log.Printf("Error: CAPMC GetHWInventoryQuery failed: %s\n", err.Error())
		sendJsonError(w, http.StatusBadRequest,
			fmt.Sprintf("Bad Request: JSON: %s", err))
		return
	}
	//filter hwInventory.Nodes according to NidlistRequest
	//Per Manny III, the HSM HWInventory Query API currently only supports getting the whole SystemHWInventory
	//so it is necessary to prune the resulting HWInventory list here
	if len(args.Nids) > 0 {
		var prunedNodeArray []*sm.HWInvByLoc
		nodeArray := hwInventory.Nodes
		for _, node := range *nodeArray {
			if _, ok := xnameNidInfoLookup[node.ID]; ok {
				prunedNodeArray = append(prunedNodeArray, node)
			}
		}
		//reset Nodes array with pruned list
		if len(prunedNodeArray) > 0 {
			hwInventory.Nodes = &prunedNodeArray
		}
	}

	var powerCapGroups []capmc.PowerCapGroup
	//get the monikerGroups from the SystemHWInventory
	uniqueMonikerGroups := convertSystemHWInventoryToUniqueMonikerGroups(hwInventory)
	for _, uniqueMonikerGroup := range uniqueMonikerGroups {
		//create the Nids array in this moniker group from the Xnames list
		for _, xname := range uniqueMonikerGroup.Xnames {
			if nidInfo, ok := xnameNidInfoLookup[xname]; ok {
				uniqueMonikerGroup.Nids = append(uniqueMonikerGroup.Nids, int(nidInfo.Nid))
			}
		}
		powerCapGroup, err := buildPowerCapCapabilitiesGroup(uniqueMonikerGroup, xnameComponentLookup)
		if err != nil {
			log.Printf("Error: CAPMC Get Power Cap Capabilities failed: %s\n", err.Error())
			sendJsonError(w, http.StatusBadRequest,
				fmt.Sprintf("Bad Request: JSON: %s", err))
			return
		}
		powerCapGroups = append(powerCapGroups, powerCapGroup)
	}

	var data capmc.GetPowerCapCapabilitiesResponse

	data.E = 0
	data.ErrMsg = ""
	data.Groups = powerCapGroups

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("Error: CAPMC Get Power Cap Capabilities encoding JSON response: %s\n", err)
	}
	return
}

// doPowerCapGet is the HTTP handler for the get_power_cap API
func (d *CapmcD) doPowerCapGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		sendJsonError(w, http.StatusMethodNotAllowed,
			fmt.Sprintf("(%s) Not Allowed", r.Method))
		return
	}

	var (
		args  capmc.NidlistRequest
		query HSMQuery
	)

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

	if len(args.Nids) > 0 {
		// Strip out duplicates. Cascade CAPMC treated duplicates as
		// OK and would return duplicate information. That seems silly.
		// A NID based query is specific enough as CAPMC can't use HSM
		// filtering so it can return state information if necessary.
		query.NIDs, _ = validateNIDs(true, args.Nids)
	} else {
		// Default list: all enabled ready compute nodes
		query.Roles = append(query.Roles, "Compute")
		query.States = append(query.States, "Ready")
		query.Enabled = append(query.Enabled, true)
	}

	log.Printf("Info: CAPMC Get Power Cap - %v", args.Nids)

	var data capmc.PowerCapResponse

	// NOTE None of these are documented in the API guide.
	// * Invalid argument (so no valid "ready" NIDs) [EINVAL]
	// * Undefined NID [EINVAL]
	// * Invalid state, NID is 'empty' [EINVAL]
	// * Invalid state, NID is 'disabled' [EINVAL]

	nodes, err := d.GetNodesByNID(query)
	if err != nil {
		var nidError *InvalidNIDsError

		// 'Invalid' NIDs requested
		if errors.As(err, &nidError) {
			for _, nid := range nidError.NIDs {
				data.Nids = append(data.Nids,
					newPowerCapNidError(nid, 22,
						"Undefined NID"))
			}

			// Continue and see if any of the requested NIDs
			// aren't "Ready" but won't retrieve any caps.
			data.E = 22 // EINVAL
			data.ErrMsg = "Invalid argument"
		} else {
			log.Printf("Error: %s", err)
			sendJsonError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	if len(args.Nids) > 0 {
		// Sort out Nodes that aren't ready. Specifically can't use
		// HSM for this task as this information needs to be passed
		// back to the caller.
		for _, node := range nodes {
			if !node.Enabled || node.State != string(base.StateReady) {
				data.Nids = append(data.Nids,
					newPowerCapNidError(node.Nid, 22,
						"Invalid state, NID is not 'ready'"))
			}
		}

		// The request contains nodes that aren't ready.
		if len(data.Nids) > 0 {
			data.E = 22 // EINVAL
			data.ErrMsg = "Invalid argument"
		}
	}

	// Olympus power control can have more than 1 target URL to gather power
	// information from. Create additional node entries, one for each power
	// URL that is needed.
	var newNodes []*NodeInfo
	for _, node := range nodes {
		if node.RfControlsCnt > 0 {
			for _, cap := range node.PowerCaps {
				if cap.Name == "Node Power Limit" {
					node.RfPowerURL = cap.Path
				} else {
					nni := *node
					nni.RfPowerURL = cap.Path
					newNodes = append(newNodes, &nni)
				}
			}
		}
	}

	if len(newNodes) > 0 {
		nodes = append(nodes, newNodes...)
	}

	// Only get power caps if all the NIDs were 'good'.
	if data.E == 0 {
		var failed int
		cmd := bmcCmd{cmd: bmcCmdGetPowerCap}
		waitNum, waitChan := d.queueBmcCmd(cmd, nodes)
		var controlMap = make(map[int][]capmc.PowerCapControl)
		for i := 0; i < waitNum; i++ {
			result := <-waitChan
			if result.rc != 0 {
				data.Nids = append(data.Nids,
					newPowerCapNidError(result.ni.Nid,
						result.rc,
						"Error getting power cap from NID"))
				log.Printf("Notice: get power cap failed: %s", result.msg)
				failed++
				continue
			}

			var rfPower capmc.Power
			err := json.Unmarshal([]byte(result.msg), &rfPower)
			if err != nil {
				data.Nids = append(data.Nids,
					newPowerCapNidError(result.ni.Nid,
						74, // EBADMSG (Linux)
						fmt.Sprintf("Error decoding Redfish Power data for NID: %s", err)))
				log.Printf("Notice: Unmarshal failed: %s", err)
				failed++
				continue
			}

			// This would be nice to use but not all versions
			// of the schema support PowerControl@odata.count.
			// Looking at you Intel...
			if d.debug {
				log.Printf("Debug: PowerControl Count %d",
					rfPower.PowerCtlCnt)
			}

			if rfPower.Error != nil {
				log.Printf("Notice: %s %s: Invalid license for power capping for NID %d (%s)",
					result.ni.BmcType, result.ni.BmcFQDN,
					result.ni.Nid, result.ni.Hostname)
				data.Nids = append(data.Nids,
					newPowerCapNidError(result.ni.Nid,
						-1, "Invalid license"))
				failed++
				continue
			}

			pctlLen := len(rfPower.PowerCtl)
			hpePctlLen := len(rfPower.ActualPowerLimits) +
				len(rfPower.PowerLimitRanges) +
				len(rfPower.PowerLimits)
			ctlLen := result.ni.RfControlsCnt

			if pctlLen < 1 && hpePctlLen < 1 && ctlLen < 1 {
				log.Printf("Notice: %s %s: No Redfish power control data for NID %d (%s)",
					result.ni.BmcType, result.ni.BmcFQDN,
					result.ni.Nid, result.ni.Hostname)
				data.Nids = append(data.Nids,
					newPowerCapNidError(result.ni.Nid,
						66, // ENODATA (Linux)
						"No Redfish Power data for NID"))
				failed++
				continue
			}

			var val *int
			var controls []capmc.PowerCapControl
			if ctlLen > 0 {
				if rfPower.SetPoint != nil {
					val = rfPower.SetPoint
				} else {
					var unconstrained int
					val = &unconstrained
				}
				controlMap[result.ni.Nid] = append(controlMap[result.ni.Nid],
					capmc.PowerCapControl{Name: rfPower.Name, Val: val})
			} else if hpePctlLen > 0 {
				// Handle Apollo 6500 AccPowerService power cap query
				for _, pl := range rfPower.PowerLimits {
					name := "Node Power Limit"
					if pl.PowerLimitInWatts != nil {
						val = pl.PowerLimitInWatts
					} else {
						var unconstrained int
						// Per Cascade 0 is unconstrained.
						// So if there isn't a control it
						// must be by definition unconstrained.
						val = &unconstrained
					}
					controls = append(controls,
						capmc.PowerCapControl{Name: name, Val: val})
				}
			} else {
				// Handle standard Redfish PowerControl power cap query
				for _, pc := range rfPower.PowerCtl {
					var name string
					if isHpeServer(result.ni) {
						name = "Node Power Limit"
					} else {
						name = pc.Name
					}

					if pc.PowerLimit != nil {
						val = pc.PowerLimit.LimitInWatts
					} else {
						var unconstrained int
						// Per Cascade 0 is unconstrained.
						// So if there isn't a control it
						// must be by definition unconstrained.
						val = &unconstrained
					}
					controls = append(controls,
						capmc.PowerCapControl{Name: name, Val: val})
				}
			}

			// need to find at least one known power cap control
			if len(controls) < 1 && len(controlMap[result.ni.Nid]) < 1 {
				data.Nids = append(data.Nids,
					newPowerCapNidError(result.ni.Nid,
						66, // ENODATA (Linux)
						"No Redfish power cap controls found for NID"))
				failed++
				continue
			}
			if ctlLen < 1 {
				data.Nids = append(data.Nids,
					capmc.PowerCapNid{Nid: result.ni.Nid, Controls: controls})
			}
		}

		for nid, e := range controlMap {
			data.Nids = append(data.Nids,
				capmc.PowerCapNid{Nid: nid, Controls: e})
		}

		if failed > 0 {
			data.E = 52 // EBADE ?
			data.ErrMsg = "Invalid exchange"
		}
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("Error: encoding JSON response: %s\n", err)
	}
}

// doPowerCapGet is the HTTP handler for the get_power_cap API
func (d *CapmcD) doPowerCapSet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		sendJsonError(w, http.StatusMethodNotAllowed,
			fmt.Sprintf("(%s) Not Allowed", r.Method))
		return
	}

	var args capmc.SetPowerCapRequest

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

	if len(args.Nids) == 0 {
		sendJsonError(w, http.StatusBadRequest, "no NIDs")
		return
	}

	var (
		data  capmc.PowerCapResponse
		nodes []*NodeInfo
	)

	// NOTE None of these checks are documented in the API guide.

	nidsMap := make(map[int]int)
	nids := []int{}
	for n, nid := range args.Nids {
		// check for duplicates
		if _, ok := nidsMap[nid.Nid]; !ok {
			nidsMap[nid.Nid] = n
			nids = append(nids, nid.Nid)
		} else {
			// keep track of duplicates
			data.Nids = append(data.Nids,
				newPowerCapNidError(nid.Nid, 22, "Duplicate NID"))
		}
	}

	// NOTE Use the post duplicate check NID list rather than the full
	// input args here as the later could get very ugly (large) in the log.
	log.Printf("Info: CAPMC Set Power Cap - %v", nids)

	var query = HSMQuery{
		NIDs: nids,
	}
	nodes, err = d.GetNodesByNID(query)
	if err != nil {
		var nidError *InvalidNIDsError

		// 'Invalid' NIDs requested
		if errors.As(err, &nidError) {
			for _, nid := range nidError.NIDs {
				data.Nids = append(data.Nids,
					newPowerCapNidError(nid, 22,
						"Undefined NID"))
			}

			// Continue and see if any of the requested
			// NIDs aren't "Ready" but won't set any caps.
			data.E = 22 // EINVAL
			data.ErrMsg = "Invalid argument"
		} else {
			log.Printf("Error: %s", err)
			sendJsonError(w, http.StatusInternalServerError,
				err.Error())
			return
		}
	}

	bmcCmds := make(map[*NodeInfo]bmcCmd)
	var newNodes []*NodeInfo
	for _, node := range nodes {
		if !node.Enabled || node.State != string(base.StateReady) {
			data.Nids = append(data.Nids,
				newPowerCapNidError(node.Nid, 22,
					"Invalid state, NID is not 'ready'"))
			continue
		}

		if node.Role != string(base.RoleCompute) {
			data.Nids = append(data.Nids,
				newPowerCapNidError(node.Nid, 22,
					"Invalid type, not a compute node"))
			continue
		}

		controls := args.Nids[nidsMap[node.Nid]].Controls
		if len(controls) < 1 {
			data.Nids = append(data.Nids,
				newPowerCapNidError(node.Nid, 22,
					"Invalid command, 'control' is not an object"))
			continue
		}

		// Loop through all the controls and generate a powerGen structure
		// that will be used to generate the payload later
		pGen, err := generateControls(node, controls)

		if err != nil {
			data.Nids = append(data.Nids,
				newPowerCapNidError(node.Nid, 22, err.Error()))
			continue
		}

		if len(pGen) == 0 {
			data.Nids = append(data.Nids,
				newPowerCapNidError(node.Nid, 22,
					"No NID supported controls specified"))
			continue
		}

		// Loop through each of the controls that need their own action to
		// generate a json payload to add to the bmcCmds map. The index to the
		// bmcCmds map is the address of a NodeInfo structure.
		var payload []byte
		for n, pg := range pGen {
			payload, err = generatePayload(n, pg)

			if err != nil {
				log.Printf("Error: %s", err)
				continue
			}

			if d.debug {
				log.Printf("Debug: payload=%s", payload)
			}

			newNodes = append(newNodes, n)
			bmcCmds[n] = bmcCmd{
				cmd:     bmcCmdSetPowerCap,
				payload: payload,
			}
		}
	}

	// The request contained invalid NIDs, controls, and/or values.
	if len(data.Nids) > 0 {
		// Punt! The request was invalid.
		// Do Not Pass Go. Do Not Collect $200.
		data.E = 22 // EINVAL
		data.ErrMsg = "Invalid Argument"
	}

	// There were no supported PowerControls
	if len(bmcCmds) <= 0 {
		log.Printf("Info: no supported power capping controls for request")
		data.E = 22 // EINVAL
		data.ErrMsg = "No supported power capping controls"
	}

	// Only set power caps if all the NIDs, controls, and values were 'good'
	if data.E == 0 {
		var failed int
		if len(newNodes) > 0 {
			nodes = newNodes
		}
		waitNum, waitChan := d.queueBmcCmds(bmcCmds, nodes)
		for i := 0; i < waitNum; i++ {
			result := <-waitChan
			if result.rc != 0 {
				data.Nids = append(data.Nids,
					newPowerCapNidError(result.ni.Nid,
						result.rc,
						"Error setting power cap for NID"))
				log.Printf("Notice: set power cap failed: %s", result.msg)
				failed++
				continue
			}
		}

		if failed > 0 {
			data.E = 52 // EBADE ?
			data.ErrMsg = "Invalid exchange"
		}
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("Error encoding JSON response: %s\n", err)
	}
}

func generateControls(node *NodeInfo, controls []capmc.PowerCapControl) (map[*NodeInfo]powerGen, error) {
	powerCtl := make([]capmc.PowerControl, node.RfPwrCtlCnt)
	seen := make(map[string]bool)
	pControls := make(map[*NodeInfo]powerGen)
	nControls := []capmc.RFControl{}

	for _, control := range controls {
		var (
			min, max int = -1, -1
			ok       bool
			pc       PowerCap
		)

		if pc, ok = node.PowerCaps[control.Name]; !ok {
			log.Printf("Notice: skipping undefined control for NID: %s", control.Name)
			continue
		}

		if _, ok = seen[control.Name]; ok {
			errStr := fmt.Sprintf("Duplicate control specified: %s", control.Name)
			return nil, errors.New(errStr)
		}

		seen[control.Name] = true

		min = pc.Min
		max = pc.Max
		// Trick the marshaller to put a 0 into a json payload by using the
		// address of a variable that contains a 0.
		zero := 0

		// The vaule of Zero is used by most vendors as a method of turning off
		// power capping so it is a valid option.
		if (min != -1) && (*control.Val < min) && (*control.Val != 0) {
			errStr := fmt.Sprintf("Control (%s) value (%d) is less than minimum (%d)",
				control.Name, *control.Val, min)
			return nil, errors.New(errStr)
		}
		if (max != -1) && (*control.Val > max) {
			errStr := fmt.Sprintf("Control (%s) value (%d) is greater than maximum (%d)",
				control.Name, *control.Val, max)
			return nil, errors.New(errStr)
		}

		if node.RfControlsCnt > 0 {
			pcPath := node.PowerCaps[control.Name].Path

			ctl := capmc.RFControl{}
			if *control.Val > 0 {
				ctl = capmc.RFControl{
					Oid:         pcPath,
					SetPoint:    control.Val,
					ControlMode: "Automatic",
				}
			} else {
				ctl = capmc.RFControl{
					Oid:         pcPath,
					ControlMode: "Disabled",
				}
			}
			nControls = append(nControls, ctl)
		} else if isHpeApollo6500(node) {
			pControls[node] = powerGen{
				powerLimit: capmc.HpeConfigurePowerLimit{
					PowerLimits: []capmc.HpePowerLimits{
						{
							PowerLimitInWatts: control.Val,
							ZoneNumber:        &zero,
						},
					},
				},
			}
		} else {
			powerCtl[pc.PwrCtlIndex] =
				capmc.PowerControl{
					PowerLimit: &capmc.PowerLimit{
						LimitInWatts: control.Val,
					},
				}
		}
	}

	if len(powerCtl) > 0 && powerCtl[0].PowerLimit != nil {
		pControls[node] = powerGen{powerCtl: powerCtl}
	}

	if len(nControls) > 0 {
		deepCtl := capmc.RFControlsDeep{Members: nControls}
		pControls[node] = powerGen{controls: deepCtl}
		fullPath := nControls[0].Oid
		// Use the control URL as the basis for the deep patch URL. Chop the
		// last 2 components from the path. It does not matter which control
		// is the first element in the array, the URLs are of the same format.
		url := path.Dir(path.Dir(fullPath))
		node.RfPowerURL = url + "/Controls.Deep"
	}

	return pControls, nil
}

type powerGen struct {
	powerCtl   []capmc.PowerControl
	powerLimit capmc.HpeConfigurePowerLimit
	controls   capmc.RFControlsDeep
}

// generatePayload - take a standard Redfish power cap structure or the
// specialized HPE Apollo power cap structure and generate the []byte payload
// needed for the BMC command.
func generatePayload(node *NodeInfo, pGen powerGen) ([]byte, error) {
	var payload []byte
	var err error
	var power interface{} = nil

	if node.RfControlsCnt > 0 {
		if len(pGen.controls.Members) == 0 {
			return nil, errors.New("missing power limit information")
		}
		power = pGen.controls
	} else if isHpeApollo6500(node) {
		if len(pGen.powerLimit.PowerLimits) == 0 {
			return nil, errors.New("missing power limit information")
		}
		// HPE Apollo 6500 power capping structure
		power = capmc.HpeConfigurePowerLimit{PowerLimits: pGen.powerLimit.PowerLimits}
	} else {
		if len(pGen.powerCtl) == 0 {
			return nil, errors.New("missing power control information")
		}
		// Standard Redfish power capping structure
		power = capmc.Power{PowerCtl: pGen.powerCtl}
	}

	payload, err = json.Marshal(power)

	return payload, err
}

//buildPowerCapCapabilitiesGroup - build a PowerCapGroup
func buildPowerCapCapabilitiesGroup(monikerGroup PowerCapCapabilityMonikerGroup, xnameComponentLookup map[string]*sm.ComponentEndpoint) (group capmc.PowerCapGroup, err error) {
	if monikerGroup.Xnames == nil {
		return group, errors.New(InvalidArguments)
	}
	representativeXname := monikerGroup.Xnames[0]

	//retrieve the componentEndpoint for one of the monikerGroup's xnames
	if componentEndpoint, ok := xnameComponentLookup[representativeXname]; ok {
		//get these values from the monikerGroup
		group.Name = monikerGroup.Name
		group.Desc = monikerGroup.Desc

		powerCtlArray := componentEndpoint.RedfishSystemInfo.PowerCtlInfo.PowerCtl
		//TODO: this value requires completion of CASMHMS-2297, so this will need to be defered for now
		group.Static = 0
		var controls []capmc.PowerCapCapabilityControl
		if powerCtlArray != nil && len(powerCtlArray) > 0 {
			//Get the node PowerControl data from the node object in the PowerCtl array
			//general node power characteristics are always defined in the first entry of the powerCtlArray, so examine it separately, first.
			powerCtl0 := powerCtlArray[0]
			if powerCtl0 != nil {
				group.Supply = int(powerCtl0.PowerCapacityWatts) //PowerControl.PowerCapacityWatts
				oem := powerCtl0.OEM
				if oem != nil {
					cray := oem.Cray
					if cray != nil {
						group.Powerup = int(cray.PowerResetWatts) //PowerControl.OEM.Cray.PowerResetWatts
						powerLimit := cray.PowerLimit
						if powerLimit != nil {
							group.HostLimitMax = int(powerLimit.Max) //PowerControl.OEM.Cray.PowerLimit.Max
							group.HostLimitMin = int(powerLimit.Min) //PowerControl.OEM.Cray.PowerLimit.Min
						}
					}
					hpe := oem.HPE
					if hpe != nil {
						powerLimit := hpe.PowerLimit
						group.HostLimitMax = int(powerLimit.Max) //PowerControl.OEM.HPE.PowerLimit.Max
						group.HostLimitMin = int(powerLimit.Min) //PowerControl.OEM.HPE.PowerLimit.Min
					}
				}
			}

			for _, powerCtl := range powerCtlArray {
				var pccControl capmc.PowerCapCapabilityControl = capmc.PowerCapCapabilityControl{Name: string(powerCtl.Name), Desc: string(powerCtl.Name), Max: 0, Min: 0}
				oem := powerCtl.OEM
				if oem != nil {
					cray := oem.Cray
					if cray != nil {
						powerLimit := cray.PowerLimit
						if powerLimit != nil {
							pccControl.Max = int(powerLimit.Max) //PowerControl.OEM.Cray.PowerLimit.Max
							pccControl.Min = int(powerLimit.Min) //PowerControl.OEM.Cray.PowerLimit.Min
						}
					}
					hpe := oem.HPE
					if hpe != nil {
						powerLimit := hpe.PowerLimit
						pccControl.Max = int(powerLimit.Max) //PowerControl.OEM.HPE.PowerLimit.Max
						pccControl.Min = int(powerLimit.Min) //PowerControl.OEM.HPE.PowerLimit.Min
					}
				}
				controls = append(controls, pccControl)
			}
		}

		controlsArray := componentEndpoint.RedfishSystemInfo.Controls
		if controlsArray != nil && len(controlsArray) > 0 {
			for _, controlElem := range controlsArray {
				var pccControl capmc.PowerCapCapabilityControl = capmc.PowerCapCapabilityControl{
					Name: string(controlElem.Control.Name),
					Desc: string(controlElem.Control.Name),
					Max:  0,
					Min:  0,
				}
				max := int(controlElem.Control.SettingRangeMax)
				min := int(controlElem.Control.SettingRangeMin)
				if controlElem.Control.PhysicalContext == "Chassis" {
					group.Supply = int(controlElem.Control.SettingRangeMax)
					group.Powerup = 0 // TODO: FIND A SOURCE FOR THIS
					group.HostLimitMax = max
					group.HostLimitMin = min
				}
				pccControl.Max = max
				pccControl.Min = min

				controls = append(controls, pccControl)
			}
		}
		group.Controls = controls

		//set the nids for this PowerCapGroup
		group.Nids = monikerGroup.Nids
	}

	return group, err
}

//get the set of uniqueMonikerGroups contained in the provided SystemHWInventory
func convertSystemHWInventoryToUniqueMonikerGroups(hwInventory sm.SystemHWInventory) (uniqueMonikerGroups []PowerCapCapabilityMonikerGroup) {
	//iterate across the nodes in hwInventory list of nodes, and prune any node not in the NidlistRequest
	monikerMap := make(map[string]PowerCapCapabilityMonikerGroup)
	nodeArray := hwInventory.Nodes

	for _, node := range *nodeArray {
		var monikerType PowerCapCapabilityMonikerType = PowerCapCapabilityMonikerType{
			Version:          "3_",
			SSD:              "",
			BaseBoardType:    "",
			BaseBoardSubType: "",
			CPUID:            "",
			TDP:              "",
			NumCores:         "",
			MemSizeGiB:       "",
			MemSpeedMHZ:      "",
			Accelerator:      "",
		}
		//pull the processor related field values
		processors := node.Processors
		if (processors != nil) && (len(*processors) > 0) {
			for _, processor := range *processors {
				populatedFRUPointer := processor.PopulatedFRU
				if populatedFRUPointer != nil {
					HMSProcessorFRUInfoPointer := populatedFRUPointer.HMSProcessorFRUInfo
					if HMSProcessorFRUInfoPointer != nil {
						if HMSProcessorFRUInfoPointer.TotalCores != "" {
							monikerType.NumCores = string(HMSProcessorFRUInfoPointer.TotalCores) + "c_"
						}
						if HMSProcessorFRUInfoPointer.ProcessorId.VendorID != "" {
							monikerType.CPUID = strings.Replace(string(HMSProcessorFRUInfoPointer.ProcessorId.VendorID), " ", "", -1) + "_"
						}
						break
					}
				}
			}
		}
		//pull the memory related field values
		dimms := node.Memory
		if (dimms != nil) && (len(*dimms) > 0) {
			for _, dimm := range *dimms {
				populatedFRUPointer := dimm.PopulatedFRU
				if populatedFRUPointer != nil {
					HMSMemoryFRUInfoPointer := populatedFRUPointer.HMSMemoryFRUInfo
					if HMSMemoryFRUInfoPointer != nil {
						if HMSMemoryFRUInfoPointer.OperatingSpeedMhz != "" {
							monikerType.MemSpeedMHZ = string(HMSMemoryFRUInfoPointer.OperatingSpeedMhz) + "MHz_"
							break
						}
					}
				}
			}
		}
		//pull the NodeLocationInfo related field values
		HMSNodeLocationInfoPointer := node.HMSNodeLocationInfo
		if HMSNodeLocationInfoPointer != nil {
			monikerType.MemSizeGiB = string(HMSNodeLocationInfoPointer.MemorySummary.TotalSystemMemoryGiB) + "GiB_"
		}
		nodeAccels := node.NodeAccels
		if (nodeAccels != nil) && (len(*nodeAccels) > 0) {
			for _, nodeAccel := range *nodeAccels {
				if nodeAccel.PopulatedFRU != nil {
					monikerType.Accelerator = string(nodeAccel.PopulatedFRU.FRUID)
				}
			}
		} else {
			monikerType.Accelerator = "NoAccel"
		}
		name := fmt.Sprintf("%s%s%s%s%s%s%s%s%s%s", monikerType.Version, monikerType.SSD, monikerType.BaseBoardType,
			monikerType.BaseBoardSubType, monikerType.CPUID, monikerType.TDP, monikerType.NumCores, monikerType.MemSizeGiB, monikerType.MemSpeedMHZ, monikerType.Accelerator)
		if mGroup, ok := monikerMap[name]; ok {
			// this moniker name is already in the monikerMap, so just append the node.ID to its list of xnames
			mGroup.Xnames = append(mGroup.Xnames, node.ID)
			monikerMap[name] = mGroup
		} else {
			monikerGroup := PowerCapCapabilityMonikerGroup{Name: name, Desc: name, Xnames: []string{node.ID}}
			monikerMap[name] = monikerGroup
		}
	}
	for _, monikerGroup := range monikerMap {
		uniqueMonikerGroups = append(uniqueMonikerGroups, monikerGroup)
	}
	return uniqueMonikerGroups
}
