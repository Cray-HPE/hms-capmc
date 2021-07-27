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
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	base "github.com/Cray-HPE/hms-base"
	"github.com/Cray-HPE/hms-capmc/internal/capmc"
	compcreds "github.com/Cray-HPE/hms-compcredentials"
	rf "github.com/Cray-HPE/hms-smd/pkg/redfish"
	"github.com/Cray-HPE/hms-smd/pkg/sm"
)

// HSMQuery is used for restricting HSM queries
type HSMQuery struct {
	ComponentIDs []string
	Groups       []string
	NIDs         []int
	Roles        []string
	States       []string
	Types        []string
	Enabled      []bool
}

type InvalidCompIDsError struct {
	err     string
	CompIDs []string
}

func (e *InvalidCompIDsError) Error() string {
	return fmt.Sprintf("%s: %v", e.err, e.CompIDs)
}

type InvalidGroupsError struct {
	err    string
	Groups []string
}

func (e *InvalidGroupsError) Error() string {
	return fmt.Sprintf("%s: %v", e.err, e.Groups)
}

type InvalidNIDsError struct {
	err  string
	NIDs []int
}

func (e *InvalidNIDsError) Error() string {
	return fmt.Sprintf("%s: %v", e.err, e.NIDs)
}

// GetFromHSM makes a request to the Hardware State Manager and unpacks the results.
func (d *CapmcD) GetFromHSM(path, restrict string, v interface{}) error {
	URI := d.hsmURL.String() + path
	if restrict != "" {
		URI += "?" + restrict
	}
	req, err := http.NewRequest(http.MethodGet, URI, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	body, err := d.doRequest(req)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, v)
	if err != nil {
		return err
	}

	return err
}

// GetComponents retrieves Components from the Hardware State Manager.
func (d *CapmcD) GetComponents(restrict string) ([]*base.Component, error) {
	var components base.ComponentArray
	err := d.GetFromHSM("/State/Components", restrict, &components)
	return components.Components, err
}

// GetComponentEndpoints retrieves ComponentEndpoints from the Hardware State
// Manager.
func (d *CapmcD) GetComponentEndpoints(restrict string) ([]*sm.ComponentEndpoint, error) {
	var componentEndpoints sm.ComponentEndpointArray
	err := d.GetFromHSM("/Inventory/ComponentEndpoints", restrict, &componentEndpoints)
	return componentEndpoints.ComponentEndpoints, err
}

// GetRedfishEndpoints retrieves RedfishEndpoints from the
// Hardware State Manager.
func (d *CapmcD) GetRedfishEndpoints(restrict string) ([]*sm.RedfishEndpoint, error) {
	var redfishEndpoints sm.RedfishEndpointArray
	err := d.GetFromHSM("/Inventory/RedfishEndpoints", restrict, &redfishEndpoints)
	return redfishEndpoints.RedfishEndpoints, err
}

// GetGroups retrieves a the requested group info from the Hardware State Manager.
func (d *CapmcD) GetGroups(groups []string) ([]sm.Group, error) {
	var resp []sm.Group
	params := url.Values{}
	for _, group := range groups {
		params.Add("group", group)
	}
	err := d.GetFromHSM("/groups", params.Encode(), &resp)
	return resp, err
}

// GetComponentsQuery retrives specific Components from the Hardware Stage Manager.
func (d *CapmcD) GetComponentsQuery(xname, restrict string) ([]*base.Component, error) {
	var components base.ComponentArray
	URI := fmt.Sprintf("/State/Components/Query/%s", xname)
	err := d.GetFromHSM(URI, restrict, &components)
	return components.Components, err
}

//GetHWInventory retrieves SystemHWInventory data from the
//Hardware State Manager.
func (d *CapmcD) GetHWInventoryQuery(xname string) (sm.SystemHWInventory, error) {
	var hwInventory sm.SystemHWInventory
	var restrict = ""
	err := d.GetFromHSM(fmt.Sprintf("/Inventory/Hardware/Query/%s", xname), restrict, &hwInventory)
	return hwInventory, err
}

func getRestrictStr(query HSMQuery) string {
	params := url.Values{}
	if query.ComponentIDs != nil {
		for _, id := range query.ComponentIDs {
			params.Add("id", id)
		}
	}
	if query.Enabled != nil {
		for _, eb := range query.Enabled {
			params.Add("enabled", strconv.FormatBool(eb))
		}
	}
	if query.NIDs != nil {
		for _, nid := range query.NIDs {
			params.Add("nid", strconv.Itoa(nid))
		}
	}
	if query.Roles != nil {
		for _, role := range query.Roles {
			params.Add("role", role)
		}
	}
	if query.States != nil {
		for _, state := range query.States {
			params.Add("state", state)
		}
	}
	if query.Types != nil {
		for _, ctype := range query.Types {
			params.Add("type", ctype)
		}
	}
	return params.Encode()
}

// convertPowerCtlToPowerCaps converts HSM PowerCtlInfo into a CAPMC
// PowerCap adding to the mapping of CAPMC power capping controls to Redfish
// PowerControls.
func convertPowerCtlsToPowerCaps(pwrCtlInfo rf.PowerCtlInfo) map[string]PowerCap {
	pwrCaps := make(map[string]PowerCap)
	for i, pwrCtl := range pwrCtlInfo.PowerCtl {
		var min, max int

		// no range check; buyer beware (default)
		min, max = -1, -1
		ctlName, ok := defaultPowerCtlToCap[pwrCtl.Name]
		if !ok {
			// HPE Proliant iLO devices do not have a Name in their PowerControl
			// field. Need to check the #odata.id to determine if the system is
			// an HPE with iLO or not.
			if strings.Contains(pwrCtl.Oid, "Chassis/1/Power") {
				pwrCtl.Name = "HPE Power Control"
				ctlName = "node"
			} else {
				log.Printf("Notice: Skipped unknown PowerControl: %s",
					pwrCtl.Name)
				continue
			}
		}

		if pwrCtl.OEM != nil {
			if pwrCtl.OEM.Cray != nil {
				min = pwrCtl.OEM.Cray.PowerLimit.Min
				max = pwrCtl.OEM.Cray.PowerLimit.Max
			}
		}

		// remove any #fragment in the path
		var path string
		c := strings.Index(pwrCtl.Oid, "#")
		if c < 0 {
			path = pwrCtl.Oid
		} else {
			path = pwrCtl.Oid[:c]
		}

		if path == "" {
			// A sane default fallback.
			path = pwrCtlInfo.PowerURL
		}

		pwrCaps[ctlName] = PowerCap{
			Name:        pwrCtl.Name,
			Path:        path,
			Min:         min,
			Max:         max,
			PwrCtlIndex: i,
		}
	}

	// an empty map and a nil map aren't the same thing
	if len(pwrCaps) > 0 {
		return pwrCaps
	} else {
		return nil
	}
}

// GetNodes retrieves information from the Hardware State Manager
// and composits it into the old Hardware Management DB NodeInfo format
// returning a list of NodeInfo.
// TODO Remove/rework as this is a TEMPORARY shim supporting the previous
//      hardware management database NodeInfo format.
func (d *CapmcD) GetNodes(query HSMQuery) ([]*NodeInfo, error) {
	var (
		components         []*base.Component
		componentEndpoints []*sm.ComponentEndpoint
		credentials        map[string]compcreds.CompCredentials
		cepQuery           HSMQuery
	)

	if d.ccs == nil {
		log.Printf("The secure store isn't ready, won't be able to get credentials")
		// The connection to the secure store isn't ready
		return nil, fmt.Errorf("Connection to the secure store isn't ready. Can not get redfish credentials.")
	}

	credentials = make(map[string]compcreds.CompCredentials)

	clist1 := query.ComponentIDs

	// If there are too many ComponentIDs in the list, need to call State
	// Manager in chunks due to limited URL length
	for len(clist1) > 0 || len(components) == 0 {
		if len(clist1) > MaxComponentQuery {
			query.ComponentIDs = clist1[:MaxComponentQuery]
			clist1 = clist1[MaxComponentQuery:]
		} else {
			query.ComponentIDs = clist1
			// This sets the list to empty, breaking us out of the for loop
			clist1 = clist1[len(clist1):]
		}

		// Form the query parameter string used for the HSM queries
		restrict := getRestrictStr(query)

		// Get requested Components from HSM
		newComponents, err := d.GetComponents(restrict)
		if err != nil {
			log.Printf("Error requesting components\n")
			return nil, err
		}
		// No components found. Bail
		if len(newComponents) == 0 && len(components) == 0 {
			log.Printf("No components returned\n")
			return []*NodeInfo{}, nil
		}

		// Add most recent components to main list
		components = append(components, newComponents...)

		// We'll use the list of components that we got back from our component
		// query for this next query only when components requested equals
		// components returned. Otherwise, we just need the list of components
		// for querying the secure store for redfish credentials.
		if len(query.ComponentIDs) != len(newComponents) {
			for _, comp := range newComponents {
				cepQuery.ComponentIDs = append(cepQuery.ComponentIDs, comp.ID)
			}
		} else {
			cepQuery = query
		}

		clist2 := cepQuery.ComponentIDs

		// We may have a new list that is larger than the chunk size depending
		// on the initial query parameters.
		for len(clist2) > 0 || len(componentEndpoints) == 0 {
			if len(clist2) > MaxComponentQuery {
				cepQuery.ComponentIDs = clist2[:MaxComponentQuery]
				clist2 = clist2[MaxComponentQuery:]
			} else {
				cepQuery.ComponentIDs = clist2
				// This sets the list to empty, breaking us out of the for loop
				clist2 = clist2[len(clist2):]
			}

			restrict = getRestrictStr(cepQuery)

			// Get requested ComponentEndpoinots from HSM
			newComponentEndpoints, err := d.GetComponentEndpoints(restrict)
			if err != nil {
				log.Printf("Error requesting component endpoints\n")
				return nil, err
			}

			// Add most recent components to main list
			componentEndpoints = append(componentEndpoints,
				newComponentEndpoints...)

			// Use the secure store.
			// TODO: Consider doing the Vault call in parallel with GetComponentEndpoints()
			newCredentials, err := d.ccs.GetCompCreds(cepQuery.ComponentIDs)
			if err != nil {
				log.Printf("Error requesting credentials")
				return nil, err
			}

			// Merge the new credential map with the main credential map
			for key, val := range newCredentials {
				credentials[key] = val
			}
		}
	}

	// Make an associative array of NodeInfo indexed by ID (xname)
	nodeMap := make(map[string]*NodeInfo)
	for _, component := range components {
		// skip 'Empty' components as they won't have a match in
		// component endpoints
		if component.State == string(base.StateEmpty) {
			continue
		}
		ni := new(NodeInfo)
		ni.Hostname = component.ID
		// TODO Look into fixing base/component.go so NID is
		// an int not a json.number!
		nid, _ := component.NID.Int64()
		ni.Nid = int(nid)
		ni.Role = component.Role
		ni.State = component.State
		ni.Type = component.Type
		// HSM should never return nil but just in case...
		if component.Enabled != nil {
			ni.Enabled = *component.Enabled
		}
		nodeMap[component.ID] = ni
	}

	// Add componentEndpoint and credential information to entries in the
	// nodeMap
	for _, componentEndpoint := range componentEndpoints {
		var (
			ok bool
			ni *NodeInfo
		)

		ni, ok = nodeMap[componentEndpoint.ID]
		if !ok {
			log.Printf("Warning: %s ComponentEndpoint w/o Component",
				componentEndpoint.ID)
			continue
		}

		if componentEndpoint.RedfishChassisInfo != nil &&
			componentEndpoint.RedfishChassisInfo.Actions != nil {
			ni.RfResetTypes = componentEndpoint.RedfishChassisInfo.Actions.ChassisReset.AllowableValues
			ni.RfActionURI = componentEndpoint.RedfishChassisInfo.Actions.ChassisReset.Target
			if componentEndpoint.RedfishChassisInfo.Actions.OEM != nil {
				ni.RfEpoURI = componentEndpoint.RedfishChassisInfo.Actions.OEM.ChassisEmergencyPower.Target
			}
		}

		if componentEndpoint.RedfishSystemInfo != nil &&
			componentEndpoint.RedfishSystemInfo.Actions != nil {
			ni.RfActionURI = componentEndpoint.RedfishSystemInfo.Actions.ComputerSystemReset.Target
			ni.RfResetTypes = componentEndpoint.RedfishSystemInfo.Actions.ComputerSystemReset.AllowableValues
			// Make sure there is no leading hostname. HSM was
			// adding the RfEndpointFQDN for a time.
			ni.RfPowerURL = strings.TrimPrefix(componentEndpoint.RedfishSystemInfo.PowerURL, componentEndpoint.RfEndpointFQDN)
			ni.RfPwrCtlCnt = len(componentEndpoint.RedfishSystemInfo.PowerCtlInfo.PowerCtl)
			ni.PowerCaps = convertPowerCtlsToPowerCaps(componentEndpoint.RedfishSystemInfo.PowerCtlInfo)
		}
		if componentEndpoint.RedfishManagerInfo != nil &&
			componentEndpoint.RedfishManagerInfo.Actions != nil {
			ni.RfActionURI = componentEndpoint.RedfishManagerInfo.Actions.ManagerReset.Target
			ni.RfResetTypes = componentEndpoint.RedfishManagerInfo.Actions.ManagerReset.AllowableValues
			// TODO: Does something need to be done with OEM actions if they are available?
		}
		if componentEndpoint.RedfishOutletInfo != nil &&
			componentEndpoint.RedfishOutletInfo.Actions != nil {
			ni.RfActionURI = componentEndpoint.RedfishOutletInfo.Actions.PowerControl.Target
			ni.RfResetTypes = componentEndpoint.RedfishOutletInfo.Actions.PowerControl.AllowableValues
		}

		ni.Domain = componentEndpoint.Domain
		ni.FQDN = componentEndpoint.FQDN
		ni.BmcFQDN = componentEndpoint.RfEndpointFQDN
		ni.BmcPath = componentEndpoint.OdataID
		ni.RfType = componentEndpoint.RedfishType
		ni.RfSubtype = componentEndpoint.RedfishSubtype

		cred, ok := credentials[componentEndpoint.ID]
		if !ok {
			log.Printf("Warning: No credentials for %s",
				componentEndpoint.ID)
			continue
		}

		ni.BmcUser = cred.Username
		ni.BmcPass = cred.Password
		ni.BmcType = base.GetHMSTypeString(componentEndpoint.RfEndpointID)
	}

	nodeList := make([]*NodeInfo, 0, len(nodeMap))
	for _, ni := range nodeMap {
		nodeList = append(nodeList, ni)
	}

	return nodeList, nil
}

// TODO Remove/rework as this is TEMPORARY (hopefully) since the hardware
// state manager does not have a single API to get the node information
// required.

// GetNodesByNID gets nodes from the hardware state manager (HSM) filtering
// the results by NID.  This combines the results of two HSM API calls and
// Vault data into a unified node information structure.
func (d *CapmcD) GetNodesByNID(query HSMQuery) ([]*NodeInfo, error) {

	// Tell HSM to filter by component type 'node' if the caller wants all
	// nodes. Otherwise, filter by the specified NIDs.
	if len(query.NIDs) == 0 {
		query.Types = append(query.Types, "node")
	}

	nodes, err := d.GetNodes(query)
	if err != nil {
		return nil, err
	}

	// Check for missing NIDs if the node list came back a different
	// size than our requested list of NIDs.
	if len(query.NIDs) != 0 && len(query.NIDs) != len(nodes) {
		var (
			bad []int
			blk int
		)

		// default
		errmsg := "error retrieving nids; all nids not found"
		err = errors.New(errmsg)

		// NIDs maybe missing because of a blocked role
		if len(query.Roles) > 0 {
			query.Roles = stringSliceMap(query.Roles,
				func(s string) string {
					return strings.Trim(s, "!")
				})

			// figure out which NIDs were blacklisted
			brnodes, err := d.GetNodes(query)
			if err != nil {
				return nil, err
			}

			// keep count of blocked (blacklisted) NIDs
			blk = len(brnodes)
		}

		// Filter by NID
		nodesByNID := make(map[int]*NodeInfo)
		for _, node := range nodes {
			nodesByNID[node.Nid] = node
		}

		reqNodes := nodes[:0]
		for _, nid := range query.NIDs {
			node, ok := nodesByNID[nid]
			if !ok {
				bad = append(bad, nid)
				continue
			}

			reqNodes = append(reqNodes, node)
		}

		if len(bad) > 0 {
			if blk == 0 {
				errmsg = "nids not found"
			} else {
				if len(bad) == blk {
					errmsg = "nids role blocked"
				} else {
					errmsg = "nids role blocked/not found"
				}
			}
			err = &InvalidNIDsError{errmsg, bad}
		} else {
			log.Printf("Notice: requested NIDs (%d) != returned nodes (%d) but no missing NIDs found!", len(query.NIDs), len(nodes))
			err = nil
		}

		return reqNodes, err
	}

	return nodes, nil
}

// GetNodesByXname gets nodes from the hardware state manager (HSM) filtering
// the results by Xname.  This combines the results of three HSM API calls into
// a unified node information structure.
func (d *CapmcD) GetNodesByXname(query HSMQuery) ([]*NodeInfo, error) {

	nodes, err := d.GetNodes(query)
	if err != nil {
		return nil, err
	}

	// Check for missing xnames if the node list came back a
	// different size than our requested list of xnames.
	if len(query.ComponentIDs) != 0 && len(query.ComponentIDs) != len(nodes) {
		var (
			bad []string
			blk int
		)

		if len(query.Roles) > 0 {
			query.Roles = stringSliceMap(query.Roles,
				func(s string) string {
					return strings.Trim(s, "!")
				})

			// figure out which nodes were blacklisted
			brnodes, err := d.GetNodes(query)
			if err != nil {
				return nil, err
			}

			// keep count of blocked (blacklisted) nodes
			blk = len(brnodes)
		}

		// Filter by Xname, the Hostname is the Xname.
		nodesMap := make(map[string]*NodeInfo)
		for _, node := range nodes {
			nodesMap[node.Hostname] = node
		}

		reqNodes := nodes[:0]
		for _, xname := range query.ComponentIDs {
			node, ok := nodesMap[xname]
			if !ok {
				bad = append(bad, xname)
			} else {
				reqNodes = append(reqNodes, node)
			}
		}

		if len(bad) > 0 {
			var errmsg string

			if blk == 0 {
				errmsg = "xnames not found"
			} else {
				if len(bad) == blk {
					errmsg = "xnames role blocked"
				} else {
					errmsg = "xnames role blocked/not found"
				}
			}
			err = &InvalidCompIDsError{errmsg, bad}
		}

		return reqNodes, err
	}

	return nodes, nil
}

// GetNodesByGroup gets nodes from the hardware state manager (HSM) filtering
// the results by Group.  This combines the results of three HSM API calls into
// a unified node information structure.
func (d *CapmcD) GetNodesByGroup(query HSMQuery) ([]*NodeInfo, error) {
	var bad []string

	groups, err := d.GetGroups(query.Groups)

	groupMap := make(map[string]sm.Group)
	for _, group := range groups {
		groupMap[group.Label] = group
	}

	for _, group := range query.Groups {
		smGroup, ok := groupMap[group]
		if !ok {
			bad = append(bad, group)
		} else {
			query.ComponentIDs = append(query.ComponentIDs, smGroup.Members.IDs...)
		}
	}

	if len(bad) > 0 {
		err = &InvalidGroupsError{"groups not found", bad}
		return nil, err
	}

	groupNodes, err := d.GetNodes(query)
	if err != nil {
		return nil, err
	}

	return groupNodes, nil
}

// NOTE: This could be done using GetNodesByNID but there really is no
// need do all the work (2 HSM API calls, a Vault access, and simulated
// table join) as all the information required is in Components.

// GetNidInfo retrieves Components of type=node from the Hardware State
// Manager returning NidInfo format.
func (d *CapmcD) GetNidInfo(query HSMQuery) ([]*capmc.NidInfo, error) {

	// Tell HSM to filter by component type 'node' if the caller wants all
	// nodes. Otherwise, filter by the specified NIDs.
	if len(query.NIDs) == 0 {
		query.Types = append(query.Types, "node")
	}
	restrict := getRestrictStr(query)

	nodes, err := d.GetComponents(restrict)
	if err != nil {
		return nil, err
	}

	// Check for missing NIDs if the node list came back a different
	// size than our requested list of NIDs.
	if len(query.NIDs) != 0 && len(query.NIDs) != len(nodes) {
		bad := checkComponentsForMissingNIDs(query.NIDs, nodes)
		if len(bad) > 0 {
			return nil, &InvalidNIDsError{"nids not found", bad}
		}
	}

	reqNodes := make([]*capmc.NidInfo, 0, len(nodes))
	for _, node := range nodes {
		reqNodes = append(reqNodes, newNidInfoFromComponent(node))
	}

	return reqNodes, nil
}

// newNidInfoFromComponent returns a new NidInfo based on the node Component
func newNidInfoFromComponent(node *base.Component) *capmc.NidInfo {
	// TODO When/If the Cascade 'v0' API is supported return a JSON
	//      blob using Cname rather than Xname. The 'v1' API Cascade
	//      CAPMC for Shasta will use Xname.
	// XT/XC CName ~= Shasta XName (Component ID)
	info := new(capmc.NidInfo)
	info.Xname = node.ID

	// FIXME The rationale for base/component.go being a
	// json.number is odd.
	// As such need the ugly two step process below.
	nid, _ := node.NID.Int64()
	info.Nid = int(nid)

	info.Role = node.Role

	return info
}

// NOTE: This could be done using GetNodes but there really is no
// need do all the work (2 HSM API calls, a Vault access, and simulated
// table join) as all the information required is in Components.

// GetComponentStatus retrieves the state of Components from the Hardware
// State Manager returning a XnameStatusResponse.
func (d *CapmcD) GetComponentStatus(query HSMQuery, filter uint) (capmc.XnameStatusResponse, error) {
	var (
		cs      capmc.XnameStatusResponse
		csFlags capmc.XnameStatusFlags
	)

	components, err := d.GetComponents(getRestrictStr(query))
	if err != nil {
		return cs, err
	}

	// Check for missing Component IDs (xnames) if the component list
	// came back different size than our requested list of Component IDs.
	if len(query.ComponentIDs) != 0 && len(query.ComponentIDs) != len(components) {
		bad := checkComponentsForMissingIDs(query.ComponentIDs, components)
		if len(bad) > 0 {
			return cs, &InvalidCompIDsError{"xnames not found", bad}
		}
	}

	for _, component := range components {
		switch component.Flag {
		case string(base.FlagAlert):
			if filter&capmc.FilterShowAlertBit != 0 {
				csFlags.Alert = append(csFlags.Alert, component.ID)
			}
		case string(base.FlagLocked):
			if filter&capmc.FilterShowLockedBit != 0 {
				csFlags.Locked = append(csFlags.Locked, component.ID)
			}
		case string(base.FlagOK), "":
			if filter&capmc.FilterShowOKBit != 0 {
				csFlags.OK = append(csFlags.OK, component.ID)
			}
		case string(base.FlagUnknown):
			if filter&capmc.FilterShowUnknownBit != 0 {
				csFlags.Unknown = append(csFlags.Unknown, component.ID)
			}
		case string(base.FlagWarning):
			if filter&capmc.FilterShowWarningBit != 0 {
				csFlags.Warning = append(csFlags.Warning, component.ID)
			}
		default:
			// This indicates that a new HMSFlag was added and
			// CAPMC is out of sync with HSM
			log.Printf("Error: unknown flag '%s'; skipping\n", component.Flag)
		}

		if component.Enabled != nil && *component.Enabled == false {
			// count as a 'flag'
			if filter&capmc.FilterShowDisabledBit != 0 {
				csFlags.Disabled = append(csFlags.Disabled, component.ID)
			}
		}

		// This is more general than required. The incoming HSM query
		// parameters can limit the categories returned.
		switch component.State {
		case string(base.StateEmpty):
			if filter&capmc.FilterShowEmptyBit != 0 {
				cs.Empty = append(cs.Empty, component.ID)
			}
		case string(base.StateHalt):
			if filter&capmc.FilterShowHaltBit != 0 {
				cs.Halt = append(cs.Halt, component.ID)
			}
		case string(base.StateOff):
			if filter&capmc.FilterShowOffBit != 0 {
				cs.Off = append(cs.Off, component.ID)
			}
		case string(base.StateOn):
			if filter&capmc.FilterShowOnBit != 0 {
				cs.On = append(cs.On, component.ID)
			}
		case string(base.StatePopulated):
			if filter&capmc.FilterShowPopulatedBit != 0 {
				cs.Populated = append(cs.Populated, component.ID)
			}
		case string(base.StateReady):
			if filter&capmc.FilterShowReadyBit != 0 {
				cs.Ready = append(cs.Ready, component.ID)
			}
		case string(base.StateStandby):
			if filter&capmc.FilterShowStandbyBit != 0 {
				cs.Standby = append(cs.Standby, component.ID)
			}
		case string(base.StateUnknown):
			if filter&capmc.FilterShowUnknownBit != 0 {
				cs.Unknown = append(cs.Unknown, component.ID)
			}
		default:
			// This indicates that a new HMSState was
			// added and CAPMC is out of sync with HSM
			log.Printf("Error: unknown state '%s'; mapping to undefined status\n", component.State)
			if filter&capmc.FilterShowUndefinedBit != 0 {
				cs.Undefined = append(cs.Undefined, component.ID)
			}
		}
	}

	flags := len(csFlags.Alert) + len(csFlags.Locked) + len(csFlags.OK) +
		len(csFlags.Unknown) + len(csFlags.Warning)
	if flags > 0 {
		cs.Flags = &csFlags
	}

	return cs, err
}

// NOTE: This could be done using GetNodesByNID but there really is no
// need do all the work (2 HSM API calls, a Vault access, and simulated
// table join) as all the information required is in Components.

// GetNidStatus retrieves the state of Components type=node from the Hardware
// State Manager returning a NodeStatusResponse.
func (d *CapmcD) GetNidStatus(query HSMQuery, filter uint) (capmc.NodeStatusResponse, error) {
	var (
		ns      capmc.NodeStatusResponse
		nsFlags capmc.NodeStatusFlags
	)

	// Tell HSM to filter by component type 'node' if the caller wants all
	// nodes. Otherwise, filter by the specified NIDs.
	if len(query.NIDs) == 0 {
		query.Types = append(query.Types, "node")
	}
	restrict := getRestrictStr(query)

	nodes, err := d.GetComponents(restrict)
	if err != nil {
		return ns, err
	}

	// Check for missing NIDs if the node list came back a different
	// size than our requested list of NIDs.
	if len(query.NIDs) != 0 && len(query.NIDs) != len(nodes) {
		bad := checkComponentsForMissingNIDs(query.NIDs, nodes)
		if len(bad) > 0 {
			return ns, &InvalidNIDsError{"nids not found", bad}
		}
	}

	for _, node := range nodes {
		// Sigh. Only necessary because HSM didn't use int correctly
		// making this a json.number!
		nid, _ := strconv.Atoi(string(node.NID))

		// flags is an unspecified, in the API document,
		// field the Cascade implementation returns when set.
		switch node.Flag {
		case string(base.FlagOK), "":
			// do nothing; OK isn't interesting
		case string(base.FlagWarning):
			if filter&capmc.FilterShowWarningBit != 0 {
				nsFlags.Warning = append(nsFlags.Warning, nid)
			}
		case string(base.FlagAlert):
			if filter&capmc.FilterShowAlertBit != 0 {
				nsFlags.Alert = append(nsFlags.Alert, nid)
			}
		case string(base.FlagLocked):
			// map Locked to Reserved
			if filter&capmc.FilterShowReservedBit != 0 {
				nsFlags.Reserved = append(nsFlags.Reserved, nid)
			}
		case string(base.FlagUnknown):
			log.Printf("Notice: NID %d: flag 'unknown' set; skipping flag\n", nid)
		default:
			// This indicates that a new HMSFlag was added and
			// CAPMC is out of sync with HSM
			log.Printf("Error: unknown flag '%s'; skipping\n", node.Flag)
		}

		// Shasta no longer has Disabled as a State. For the
		// purpose of being Cascade API compatible treat a
		// node that isn't enabled as having a Disabled status.
		// This check is done first since a node can only
		// appear in one status list. Disabled overrides state.
		if node.Enabled != nil && *node.Enabled == false {
			// count as disabled
			if filter&capmc.FilterShowDisabledBit != 0 {
				ns.Disabled = append(ns.Disabled, nid)
			}
			continue
		}

		// This is more general than required. The incoming
		// HSM query parameters can limit the categories
		// returned.
		switch node.State {
		case string(base.StateEmpty):
			if filter&capmc.FilterShowEmptyBit != 0 {
				ns.Empty = append(ns.Empty, nid)
			}
		case string(base.StateHalt):
			if filter&capmc.FilterShowHaltBit != 0 {
				ns.Halt = append(ns.Halt, nid)
			}
		case string(base.StateOff):
			if filter&capmc.FilterShowOffBit != 0 {
				ns.Off = append(ns.Off, nid)
			}
		case string(base.StateOn):
			if filter&capmc.FilterShowOnBit != 0 {
				ns.On = append(ns.On, nid)
			}
		case string(base.StatePopulated):
			if filter&capmc.FilterShowPopulatedBit != 0 {
				ns.Populated = append(ns.Populated, nid)
			}
		case string(base.StateReady):
			if filter&capmc.FilterShowReadyBit != 0 {
				ns.Ready = append(ns.Ready, nid)
			}
		case string(base.StateStandby):
			if filter&capmc.FilterShowStandbyBit != 0 {
				ns.Standby = append(ns.Standby, nid)
			}
		case string(base.StateUnknown):
			if filter&capmc.FilterShowUnknownBit != 0 {
				ns.Unknown = append(ns.Unknown, nid)
			}
		default:
			// This indicates that a new HMSState was
			// added and CAPMC is out of sync with HSM
			log.Printf("Error: unknown state '%s'; mapping to undefined status\n", node.State)
			if filter&capmc.FilterShowUndefinedBit != 0 {
				ns.Undefined = append(ns.Undefined, nid)
			}
		}
	}

	flags := len(nsFlags.Alert) + len(nsFlags.Reserved) + len(nsFlags.Warning)
	if flags > 0 {
		ns.Flags = &nsFlags
	}

	return ns, err
}

func (d *CapmcD) checkForDisabledComponents(nl []*NodeInfo, t string) error {
	var (
		bad []string
		err error
	)

	for _, ni := range nl {
		if ni.Enabled == false {
			if t == "nid" {
				bad = append(bad, strconv.Itoa(ni.Nid))
			} else if t == "xname" {
				bad = append(bad, ni.Hostname)
			}
		}
	}

	if len(bad) > 0 {
		var label string

		if t == "nid" {
			label = "nodes disabled"
		} else if t == "xname" {
			label = "components disabled"
		}

		err = &InvalidCompIDsError{label, bad}
	}

	return err
}

// checkComponentsForMissingIDs returns missing Component IDs not found in
// components
func checkComponentsForMissingIDs(ids []string, components []*base.Component) (missing []string) {
	componentIDs := make(map[string]bool)

	for _, component := range components {
		componentIDs[component.ID] = true
	}

	for _, id := range ids {
		_, ok := componentIDs[id]
		if !ok {
			missing = append(missing, id)
		}
	}

	return missing
}

// checkComponentsForMissingNIDs returns missing NIDs not found in nodes
func checkComponentsForMissingNIDs(nids []int, nodes []*base.Component) (missing []int) {
	componentNIDs := make(map[int]bool)

	for _, node := range nodes {
		nid, _ := node.NID.Int64()
		componentNIDs[int(nid)] = true
	}

	for _, nid := range nids {
		_, ok := componentNIDs[nid]
		if !ok {
			missing = append(missing, nid)
		}
	}

	return missing
}

// validateNIDs validates an array of nids into valid and invalid arrays.
func validateNIDs(dupsOK bool, nids []int) ([]int, []int) {

	seenNIDs := make(map[int]bool)
	validNIDs := nids[:0]
	invalidNIDs := []int{}

	for _, nid := range nids {
		if _, seenNID := seenNIDs[nid]; !seenNID {
			seenNIDs[nid] = true
			if nid < 0 {
				invalidNIDs = append(invalidNIDs, nid)
			} else {
				validNIDs = append(validNIDs, nid)
			}
		} else {
			// keep track of duplicates?
			if !dupsOK {
				invalidNIDs = append(invalidNIDs, nid)
			}
		}
	}

	return validNIDs, invalidNIDs
}

// doRequest sends a HTTP request
func (d *CapmcD) doRequest(req *http.Request) ([]byte, error) {

	//This func is only for HSM access, so use the non-cert transport.
	rsp, err := d.smClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer rsp.Body.Close()
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	if rsp.StatusCode >= http.StatusMultipleChoices {
		// TODO improve message
		log.Printf("Error: status=%s body=%s -> HTTP %s %s",
			rsp.Status, body, req.Method, req.URL.String())
		// TODO need better error handling
		contentType := rsp.Header.Get("Content-Type")
		switch contentType {
		case "text/plain":
			log.Printf("Whoa got a text/plain error response!")
			err = fmt.Errorf("%s", body)
		case "application/json":
			var errRsp capmc.ErrResponse
			log.Printf("Received application/json error response")
			err = json.Unmarshal(body, &errRsp)
			if err != nil {
				log.Printf("Failed to Unmarshal ErrResponse: %s", err)
				err = fmt.Errorf("%s", body)
			} else {
				err = fmt.Errorf("%s", errRsp.ErrMsg)
			}
		case "application/problem+json":
			// unmarshal RFC7807 format here
			var pd base.ProblemDetails
			err = json.Unmarshal(body, &pd)
			if err != nil {
				log.Printf("Info: Failed to Unmarshal RFC8707 data: %s", err)
				err = fmt.Errorf("%s", body)
				break
			}
			// parse out the problem details - not sure what will be included
			// so build with parts that are present
			errMsg := "Hardware State Manager error"
			if pd.Status >= 400 && pd.Title != "" {
				errMsg = fmt.Sprintf("%s: Error: %d %s", errMsg, pd.Status, pd.Title)
			}
			if pd.Detail != "" {
				errMsg = fmt.Sprintf("%s: Details: %s", errMsg, pd.Detail)
			}
			if pd.Instance != "" {
				errMsg = fmt.Sprintf("%s: Instance: %s", errMsg, pd.Instance)
			}
			// NOTE: Type is filled with "about:blank" if there is not an active
			//  URL that will describe the problem when dereferenced.  We may as
			//  well skip the blank type information.
			if pd.Type != "about:blank" {
				errMsg = fmt.Sprintf("%s: Type: %s", errMsg, pd.Type)
			}
			err = fmt.Errorf("%s", errMsg)
		default:
			log.Printf("Unsupported Content-Type: %s in Error Response", contentType)
			err = fmt.Errorf("%s", body)
		}
		return body, err
	}

	return body, nil
}
