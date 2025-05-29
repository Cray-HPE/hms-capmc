/*
 * MIT License
 *
 * (C) Copyright [2019-2022,2025] Hewlett Packard Enterprise Development LP
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

// This file is the core of the capmcd application, and includes the main entry point

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"path"
	"regexp"
	"strings"
	"time"

	base "github.com/Cray-HPE/hms-base/v2"
	"github.com/Cray-HPE/hms-capmc/internal/capmc"
	rf "github.com/Cray-HPE/hms-smd/v2/pkg/redfish"
)

// warRedFishPowerOID takes a NodeInfo and constructs what the HSM normally
// returns in the ComponentEndpoint PowerURL. This is a work around that is
// only required for a HSM that doesn't use the oid to determine the path
// it stores in PowerURL.
func warRedfishPowerOID(ni *NodeInfo) string {
	log.Printf("Notice: HSM ComponentEndpoint empty RFPowerURL; constructing path")
	base := path.Base(ni.BmcPath)

	// base works for everything (so far) but Intel
	matched, _ := regexp.MatchString(`QSBP.*`, base)
	if matched {
		log.Printf("Info: Using Intel Redfish Power OID workaround")
		base = "RackMount/Baseboard"
	}

	oid := path.Join(path.Dir(path.Dir(ni.BmcPath)),
		"Chassis", base, "Power")

	return oid
}

// decodeBmcResponse decodes a response from a Redfish endpoint into a
// standard message string.
func (d *CapmcD) decodeBmcResponse(ni *NodeInfo, resp *http.Response) (msg string) {
	log.Printf("Notice: HTTP %s %s %s", resp.Status, resp.Request.Method,
		resp.Request.URL.String())
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error: reading response body: %s\n", err)
		msg = "Internal Server Error"
		return
	}
	defer resp.Body.Close()

	contentType, _, _ := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	switch contentType {
	case "application/json":
		var rferr rf.RedfishError

		err := json.Unmarshal(body, &rferr)
		if err == nil {
			msg = fmt.Sprintf("Redfish Error Message: %s",
				rferr.Error.Message)

			if len(rferr.Error.ExtendedInfo) > 0 {
				msg += " ExtendedInfo:"
			}
			for _, einfo := range rferr.Error.ExtendedInfo {
				msg += fmt.Sprintf(" Message: %s",
					einfo.Message)
				msg += fmt.Sprintf(" Resolution: %s",
					einfo.Resolution)
			}

			break
		} else {
			// Early Dell Redfish implementations use a
			// "non-standard" error response. Try that form next.
			var dellerr struct {
				Error struct {
					Code    string `json:"code"`
					Message struct {
						Lang  string `json:"lang"`
						Value string `json:"value"`
					} `json:"message"`
				} `json:"error"`
			}

			err := json.Unmarshal(body, &dellerr)
			if err == nil {
				msg = fmt.Sprintf("Redfish Error Message: %s",
					dellerr.Error.Message.Value)
				break
			}
		}

		log.Printf("Notice: unable to unmarshal %s into known error response type",
			contentType)

		fallthrough
	default:
		if contentType != "application/json" {
			log.Printf("Notice: unrecognized Redfish response Content-Type: %s\n",
				contentType)
		}

		msg = fmt.Sprintf("unrecognized Redfish response (%s) format: %s\n",
			contentType, string(body))
	}

	// Not all HTTP errors return a json body. Those that don't will return
	// basic HTML. If the decode fails, just continue on and return the
	// response message.

	msg = fmt.Sprintf("%s %s HTTP %s, %s",
		ni.BmcType, ni.BmcFQDN, resp.Status, msg)

	log.Printf("Notice: %s\n", msg)

	return msg
}

// queueBmcJob queues up a Redfish API call to the global worker pool.
func (d *CapmcD) queueBmcCall(call bmcCall) {
	job := NewJobBmcPwr(call, d)
	// workerPoolQueue() returns 1 if the queue is full. We wait
	// here until all of our jobs can be queued.
	for d.WPool.Queue(job) == 1 {
		time.Sleep(1 * time.Second)
	}
}

// queueBmcCmd queues up a command that will concurrently make the same
// Redfish API call to each component by workers in the global worker pool.
func (d *CapmcD) queueBmcCmd(cmd bmcCmd, nodes []*NodeInfo) (int, <-chan bmcPowerRc) {
	waiters := len(nodes)
	rspChan := make(chan bmcPowerRc, waiters)

	for _, node := range nodes {
		var call = bmcCall{bmcCmd: cmd, ni: node, rspChan: rspChan}
		d.queueBmcCall(call)
	}

	return waiters, rspChan
}

// queueBmcCmds queues up commands that will concurrently make the Redfish API
// calls to each component by workers in the global worker pool.
func (d *CapmcD) queueBmcCmds(cmds map[*NodeInfo]bmcCmd, nodes []*NodeInfo) (int, <-chan bmcPowerRc) {
	waiters := len(nodes)
	rspChan := make(chan bmcPowerRc, waiters)

	for _, node := range nodes {
		cmd, ok := cmds[node]
		if !ok {
			log.Printf("Error: no BMC command for %s", node.Hostname)
			continue
		}
		var call = bmcCall{bmcCmd: cmd, ni: node, rspChan: rspChan}
		d.queueBmcCall(call)
	}

	return waiters, rspChan
}

// doBmcStatusCall - Handle the specific action of getting power status via
// Redfish for a target node identified by a NodeInfo structure. Returns a
// response that includes the NodeInfo structure, an error code, a message, and
// the power state of the node requested.
func (d *CapmcD) doBmcStatusCall(ni *NodeInfo) bmcPowerRc {
	var res = bmcPowerRc{ni: ni, rc: -1, state: "Unknown"}

	nodePath := "https://" + ni.BmcFQDN + ni.BmcPath

	// check for simulation only
	if d.simulationOnly {
		log.Printf("SIMULATION_ONLY: doBmcStatusCall with: GET %s", nodePath)
		res.state = "Simulation only - request not sent to hardware"
		return res
	}

	// create the request
	req, err := http.NewRequest("GET", nodePath, nil)
	req.SetBasicAuth(ni.BmcUser, ni.BmcPass)
	req.Header.Set("Accept", "*/*")
	req.Close = true

	// execute the reqest
	rfClientLock.RLock()
	rsp, err := d.rfClient.Do(req)
	rfClientLock.RUnlock()
	defer base.DrainAndCloseResponseBody(rsp)
	if err != nil {
		log.Printf("GET %s\n %s Network Error: %s",
			nodePath, ni.BmcType, err)
		res.msg = fmt.Sprintf("%s Communication Error", ni.BmcType)
		return res
	}

	stsBody, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Printf("Error: reading response body: %s\n", err)
		res.msg = "Internal Server Error"
		return res
	}

	if rsp.StatusCode >= http.StatusBadRequest {
		rsp.Body = ioutil.NopCloser(bytes.NewBuffer(stsBody))
		log.Printf("Error: HTTP %s GET %s", rsp.Status, nodePath)
		res.msg = d.decodeBmcResponse(ni, rsp)
		res.rc = rsp.StatusCode
		return res
	}

	var powerState string

	switch ni.RfType {
	case rf.ChassisType:
		var info rf.Chassis
		err = json.Unmarshal(stsBody, &info)
		if err == nil {
			if info.PowerState != "" {
				powerState = info.PowerState
			} else {
				log.Printf("Info: no power state for (%s/%s) %s %s; assuming 'On'",
					ni.RfType, ni.RfSubtype, ni.Type, ni.Hostname)
				powerState = rf.POWER_STATE_ON
			}
		}
	case rf.ComputerSystemType:
		var info rf.ComputerSystem
		err = json.Unmarshal(stsBody, &info)
		if err == nil {
			powerState = info.PowerState
		}
	case rf.OutletType:
		var info rf.Outlet
		err = json.Unmarshal(stsBody, &info)
		if err == nil {
			powerState = info.PowerState
		}
	case rf.ManagerType:
		var info rf.Manager
		err = json.Unmarshal(stsBody, &info)
		if err == nil {
			if d.debug {
				log.Printf("Debug: Status: %v\n", info.Status)
			}
			// Managers don't really have a power state so
			// we'll assume any reponse means 'On'
			powerState = rf.POWER_STATE_ON
		}
	default:
		log.Printf("Error: %s: unknown Redfish Type\n", ni.RfType)
		// punt
		return res
	}

	if err != nil {
		log.Printf("Error: decoding response body: %s\n", err)
		res.msg = fmt.Sprintf("%s Response Decode Error", ni.BmcType)
		return res
	}

	res.rc = 0
	res.state = powerState

	log.Printf("Info: %s %s [%s %s] Power State: %s\n",
		ni.Type, ni.Hostname, ni.BmcType, ni.BmcFQDN, res.state)

	return res
}

func (d *CapmcD) doBmcPowerCall(call bmcCall) bmcPowerRc {
	ni := call.ni
	var res = bmcPowerRc{ni: ni, rc: -1, state: "Unknown"}
	// NOTE The res.state isn't that important at this point. It is
	//      only used with the "status" command.  Guessing about Redfish
	//      PowerState isn't a good idea so if it is needed then querying
	//      for it makes more sense.

	resetType, err := d.cmdToResetType(call.cmd, ni.RfResetTypes)
	if err != nil {
		log.Printf("failed converting %s to Redfish ResetType: %s",
			call.cmd, err)
		log.Printf("%s %s: AllowableVaules: %s\n",
			ni.BmcType, ni.BmcFQDN, ni.RfResetTypes)
		res.msg = fmt.Sprintf("%s %s: %s", ni.BmcType, ni.BmcFQDN, err)
		return res
	}

	var body string
	var sessionAuthPath = ""
	var sessionAuthBody = ""
	switch ni.Type {
	// The CabinetPDUOutlet HMSType has been depricated in favor of
	// CabinetPDUPowerConnector. Support both for now.
	case "CabinetPDUOutlet":
		fallthrough
	case "CabinetPDUPowerConnector":
		var HPEPDU = true
		if strings.Contains(ni.BmcFQDN, "rts") {
			HPEPDU = false
		}
		if HPEPDU {
			outletNum := strings.Split(ni.Hostname, "v")
			if len(outletNum) < 2 {
				log.Printf("ERROR: Could not get outlet number")
				// Just return because it will not work
				return res
			}
			body = fmt.Sprintf(`{"OutletNumber":%s,"StartupState":"on","Outletname":"OUTLET%s","OnDelay":0,"OffDelay":0,"RebootDelay":5,"OutletStatus":"%s"}`, outletNum[1], outletNum[1], strings.ToLower(resetType))
			sessionAuthPath = "https://" + ni.BmcFQDN + "/redfish/v1/SessionService/Sessions"
			sessionAuthBody = fmt.Sprintf(`{"username":"%s","password":"%s"}`, ni.BmcUser, ni.BmcPass)
		} else {
			body = fmt.Sprintf(`{"PowerState": "%s"}`, resetType)
		}
	default:
		body = fmt.Sprintf(`{"ResetType": "%s"}`, resetType)
	}

	actionPath := "https://" + ni.BmcFQDN + ni.RfActionURI

	// check for simulation only
	if d.simulationOnly {
		log.Printf("SIMULATION_ONLY: doBmcPowerCall with: POST %s, Data: %s", actionPath, body)
		res.state = "Simulation only - request not sent to hardware"
		return res
	}

	if resetType == "PushPowerButton" {
		// Check current power state against requested power state
		status := d.doBmcStatusCall(ni)
		// If we are already in the desired power state, do nothing
		if status.state == call.cmd {
			res.rc = 0
			return res
		}
	}

	var sessionAuthToken string
	if len(sessionAuthPath) > 0 {
		req, err := http.NewRequest("POST", sessionAuthPath, bytes.NewBuffer([]byte(sessionAuthBody)))
		req.SetBasicAuth(ni.BmcUser, ni.BmcPass)
		req.Header.Set("Accept", "*/*")
		req.Header.Set("Content-Type", "application/json")
		// execute the request
		rfClientLock.RLock()
		rsp, err := d.rfClient.Do(req)
		rfClientLock.RUnlock()
		defer base.DrainAndCloseResponseBody(rsp)
		if err != nil {
			log.Printf("POST %s\n%s Network Error: %s",
				sessionAuthPath, ni.BmcType, err)
			res.msg = fmt.Sprintf("%s Communication Error", ni.BmcType)
			return res
		}
		sessionAuthToken = rsp.Header.Get("X-Auth-Token")
	}
	log.Printf("doBmcPowerCall with: POST %s, Data: %s", actionPath, body)
	// create the request
	req, err := http.NewRequest("POST", actionPath, bytes.NewBuffer([]byte(body)))
	req.SetBasicAuth(ni.BmcUser, ni.BmcPass)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Content-Type", "application/json")
	if len(sessionAuthToken) > 0 {
		req.Header.Set("X-Auth-Token", sessionAuthToken)
	}

	// execute the request
	rfClientLock.RLock()
	rsp, err := d.rfClient.Do(req)
	rfClientLock.RUnlock()
	defer base.DrainAndCloseResponseBody(rsp)
	if err != nil {
		log.Printf("POST %s\n Body           --> %s\n %s Network Error: %s",
			actionPath, body, ni.BmcType, err)
		res.msg = fmt.Sprintf("%s Communication Error", ni.BmcType)
		return res
	}

	cmdBody, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Printf("Error: reading response body: %s\n", err)
		res.msg = "Internal Server Error"
		return res
	}

	if rsp.StatusCode >= http.StatusBadRequest {
		rsp.Body = ioutil.NopCloser(bytes.NewBuffer(cmdBody))
		res.msg = d.decodeBmcResponse(ni, rsp)
		res.rc = rsp.StatusCode
		return res
	}

	res.rc = 0
	return res
}

func (d *CapmcD) doBmcGetCall(call bmcCall) bmcPowerRc {
	var (
		oid string
		res = bmcPowerRc{ni: call.ni, rc: -1, state: "Unknown"}
	)

	switch call.cmd {
	case bmcCmdGetPowerCap:
		// Just in case...
		if call.ni.RfPowerURL == "" {
			oid = warRedfishPowerOID(call.ni)
		} else {
			oid = call.ni.RfPowerURL
		}
	default:
		res.msg = fmt.Sprintf("Invalid command %s", call.cmd)
		log.Printf("Error: %s", res.msg)
		return res
	}

	if oid == "" {
		log.Printf("Error: missing URI path for %s", call.cmd)
		res.msg = fmt.Sprintf("Internal service error: no URI")
		return res
	}

	bmcURI := "https://" + call.ni.BmcFQDN + oid

	// check for simulation only
	if d.simulationOnly {
		log.Printf("SIMULATION_ONLY: doBmcGetCall with: GET %s", bmcURI)
		res.state = "Simulation only - request not sent to hardware"
		return res
	}

	// create the request
	req, err := http.NewRequest("GET", bmcURI, nil)
	req.SetBasicAuth(call.ni.BmcUser, call.ni.BmcPass)
	req.Header.Set("Accept", "*/*")

	// execute the request
	rfClientLock.RLock()
	rsp, err := d.rfClient.Do(req)
	rfClientLock.RUnlock()
	defer base.DrainAndCloseResponseBody(rsp)
	if err != nil {
		res.msg = fmt.Sprintf("%s Communication Error", call.ni.BmcType)
		log.Printf("GET %s\n %s: %s", bmcURI, res.msg, err)
		return res
	}

	stsBody, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Printf("Error: reading response body: %s\n", err)
		res.msg = "Internal Server Error"
		return res
	}

	if rsp.StatusCode >= http.StatusBadRequest {
		rsp.Body = ioutil.NopCloser(bytes.NewBuffer(stsBody))
		res.msg = d.decodeBmcResponse(call.ni, rsp)
		res.rc = rsp.StatusCode
		return res
	}

	res.rc = 0
	res.msg = string(stsBody)

	return res
}

func (d *CapmcD) doBmcPatchCall(call bmcCall) bmcPowerRc {
	var (
		oid    string
		res    = bmcPowerRc{ni: call.ni, rc: -1, state: "Unknown"}
		pcCall = bmcCall{bmcCmd: bmcCmd{cmd: bmcCmdGetPowerCap}, ni: call.ni}
	)
	// NOTE The res.state isn't that important at this point. It is
	//      only used with the "status" command.

	switch call.cmd {
	case bmcCmdSetPowerCap:
		// Just in case...
		if call.ni.RfPowerURL == "" {
			oid = warRedfishPowerOID(call.ni)
		} else {
			oid = call.ni.RfPowerURL
		}
	default:
		res.msg = fmt.Sprintf("Invalid command %s", call.cmd)
		log.Printf("Error: %s", res.msg)
		return res
	}

	if oid == "" {
		log.Printf("Error: missing URI path for %s", call.cmd)
		res.msg = fmt.Sprintf("Internal service error: no URI")
		return res
	}

	bmcURI := "https://" + call.ni.BmcFQDN + oid

	// check for simulation only
	if d.simulationOnly {
		log.Printf("SIMULATION_ONLY: doBmcPatchCall with: PATCH %s, DATA: %s", bmcURI, string(call.payload))
		res.state = "Simulation only - request not sent to hardware"
		return res
	}

	// create the request
	req, err := http.NewRequest("PATCH", bmcURI, bytes.NewBuffer(call.payload))
	req.SetBasicAuth(call.ni.BmcUser, call.ni.BmcPass)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Content-Type", "application/json")

	if strings.Index(oid, "Controls.Deep") < 0 {
		// Query the BMC to get the etag for the PATCH call
		pcRes := d.doBmcGetCall(pcCall)

		if pcRes.rc >= http.StatusBadRequest {
			return pcRes
		}

		var rfPower capmc.Power
		err := json.Unmarshal([]byte(pcRes.msg), &rfPower)
		if err != nil {
			res.msg = fmt.Sprintf("%s unable to unmarshal status request",
				call.bmcCmd)
			log.Printf("%s", res.msg)
			return res
		}
		req.Header.Set("If-Match", rfPower.Oetag)
	}

	// execute the request
	rfClientLock.RLock()
	rsp, err := d.rfClient.Do(req)
	rfClientLock.RUnlock()
	defer base.DrainAndCloseResponseBody(rsp)
	if err != nil {
		log.Printf("PATCH %s\n Body           --> %s\n %s Network Error: %s",
			bmcURI, call.payload, call.ni.BmcType, err)
		res.msg = fmt.Sprintf("%s Communication Error", call.ni.BmcType)
		return res
	}

	cmdBody, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Printf("Error: reading response body: %s\n", err)
		res.msg = "Internal Server Error"
		return res
	}

	if rsp.StatusCode >= http.StatusBadRequest {
		rsp.Body = ioutil.NopCloser(bytes.NewBuffer(cmdBody))
		res.msg = d.decodeBmcResponse(call.ni, rsp)
		res.rc = rsp.StatusCode
		return res
	}

	res.rc = 0
	res.msg = string(cmdBody)

	return res
}

func (d *CapmcD) doBmcPostCall(call bmcCall, path string) bmcPowerRc {
	ni := call.ni
	var res = bmcPowerRc{ni: ni, rc: -1, state: "Unknown"}
	// NOTE The res.state isn't that important at this point. It is
	//      only used with the "status" command.  Guessing about Redfish
	//      PowerState isn't a good idea so if it is needed then querying
	//      for it makes more sense.

	actionPath := "https://" + ni.BmcFQDN + path

	// check for simulation only
	if d.simulationOnly {
		log.Printf("SIMULATION_ONLY: doBmcPostCall with: POST %s, Data: %s", actionPath, call.payload)
		res.state = "Simulation only - request not sent to hardware"
		return res
	}

	// create the request
	// create the request
	req, err := http.NewRequest("POST", actionPath, bytes.NewBuffer(call.payload))
	req.SetBasicAuth(ni.BmcUser, ni.BmcPass)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Content-Type", "application/json")

	// execute the request
	rfClientLock.RLock()
	rsp, err := d.rfClient.Do(req)
	rfClientLock.RUnlock()
	defer base.DrainAndCloseResponseBody(rsp)
	if err != nil {
		log.Printf("POST %s\n Body           --> %s\n %s Network Error: %s",
			actionPath, call.payload, ni.BmcType, err)
		res.msg = fmt.Sprintf("%s Communication Error", ni.BmcType)
		return res
	}

	cmdBody, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Printf("Error: reading response body: %s\n", err)
		res.msg = "Internal Server Error"
		return res
	}

	if rsp.StatusCode >= http.StatusBadRequest {
		rsp.Body = ioutil.NopCloser(bytes.NewBuffer(cmdBody))
		res.msg = d.decodeBmcResponse(ni, rsp)
		res.rc = rsp.StatusCode
		return res
	}

	res.rc = 0
	return res
}

// This is intended to be started as a goroutine and handles calls down to the
// BMC that check power state status and control power on / off. This should
// probably be generalized for fanning out other types of requests to the BMCs.
func (d *CapmcD) doBmcCall(call bmcCall) {
	ni := call.ni
	log.Printf("Info: %s: '%s', %s: '%s', Command: '%s'\n",
		ni.Type, ni.Hostname, ni.BmcType, ni.BmcFQDN, call.cmd)

	var res = bmcPowerRc{ni: ni, rc: -1, state: "Unknown"}

	if ni.BmcFQDN == "" {
		// No BMC defined
		log.Printf("Error: %s %s no FQDN defined for %s\n",
			ni.Type, ni.Hostname, ni.BmcType)
		res.msg = fmt.Sprintf("Unknown %s (%s Controller)",
			ni.BmcType, ni.Type)
		call.rspChan <- res
		return
	}

	switch call.cmd {
	case bmcCmdPowerStatus:
		// TODO The "status" command should more closely follow the Cascade
		//      CAPMC API internals by using information from HSM.  The
		//      "status" should be the logical status of the node not power
		//      state. Reporting the Redfish PowerState of the node, if
		//      desired, should be an extension to the Shasta CAPMC API.
		res = d.doBmcStatusCall(ni)
	case bmcCmdGetPowerCap:
		res = d.doBmcGetCall(call)
	case bmcCmdSetPowerCap:
		if isHpeApollo6500(call.ni) {
			res = d.doBmcPostCall(call, ni.RfPowerTarget)
		} else {
			res = d.doBmcPatchCall(call)
		}
	default:
		res = d.doBmcPowerCall(call)
	}

	call.rspChan <- res
}
