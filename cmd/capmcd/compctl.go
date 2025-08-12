/*
 * MIT License
 *
 * (C) Copyright [2019-2023] Hewlett Packard Enterprise Development LP
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
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/Cray-HPE/hms-capmc/internal/capmc"

	rf "github.com/Cray-HPE/hms-smd/pkg/redfish"
)

func cmdCompPowerSeq(cmd string) ([]string, error) {
	return svc.cmdCompPowerSeq(cmd)
}

func (d *CapmcD) cmdCompPowerSeq(cmd string) ([]string, error) {

	pc, ok := d.config.PowerControls[cmd]
	if !ok {
		return []string{}, fmt.Errorf("no power controls for %s operation", cmd)
	}

	return pc.CompSeq, nil
}

// hasCompPowerSupport - Returns true if the component type exists in the power
// sequence configuration for the specified command type.
func (d *CapmcD) hasCompPowerSupport(cmd, ctype string) (bool, error) {
	types, err := d.cmdCompPowerSeq(cmd)
	if err != nil {
		return false, err
	}
	return stringInSlice(ctype, types), nil
}

// PCS Support structures
type PCSLocation struct {
	Xname     string `json:"xname"`
	DeputyKey string `json:"deputyKey"`
}

type PCSTransition struct {
	Operation string        `json:"operation"`
	Location  []PCSLocation `json:"location"`
}

type PCSTransitionResponse struct {
	TransitionID string `json:"transitionID"`
	Operation    string `json:"operation"`
}

type PCSTaskCounts struct {
	Total       int `json:"total"`
	New         int `json:"new"`
	InProgress  int `json:"in-progress"`
	Failed      int `json:"failed"`
	Succeeded   int `json:"succeeded"`
	UnSupported int `json:"un-supported"`
}

type PCSTasks struct {
	Xname          string `json:"xname"`
	TaskStatus     string `json:"taskStatus"`
	TaskStatusDesc string `json:"taskStatusDescription"`
}

type PCSTransitionGet struct {
	TransitionID     string        `json:"transitionID"`
	Operation        string        `json:"operation"`
	TransitionStatus string        `json:"transitionStatus"`
	TaskCounts       PCSTaskCounts `json:"taskCounts"`
	Tasks            []PCSTasks    `json:"tasks"`
}

type PCSPowerStatus struct {
	Xname                     string   `json:"xname"`
	PowerState                string   `json:"powerState"`
	ManagementState           string   `json:"managementState"`
	Error                     string   `json:"error"`
	SupportedPowerTransitions []string `json:"suportedPowerTransitions"`
	LastUpdated               string   `json:"lastUpdated"`
}

type PCSStatusGet struct {
	Status []PCSPowerStatus `json:"status"`
}

func (d *CapmcD) doCompOnOffCtrl(nl []*NodeInfo, command string) capmc.XnameControlResponse {
	var data capmc.XnameControlResponse
	data.Xnames = make([]*capmc.XnameControlErr, 0, 1)
	var failures int

	var targetedXname []string
	for _, v := range nl {
		supported, err := d.hasCompPowerSupport(command, v.Type)
		if err != nil {
			failures++
			msg := fmt.Sprintf("%s", err)
			log.Printf("Error: %s.", msg)
			xnameErr := capmc.MakeXnameError(v.Hostname, -1, msg)
			data.Xnames = append(data.Xnames, xnameErr)
			continue
		}
		if !supported {
			// Skip components not in the power action sequencing list.
			msg := fmt.Sprintf("Skipping %s: Type, '%s', not defined in power sequence for '%s'", v.Hostname, v.Type, command)
			log.Printf("Info: %s.", msg)
			failures++
			xnameErr := capmc.MakeXnameError(v.Hostname, -1, msg)
			data.Xnames = append(data.Xnames, xnameErr)
			continue
		}
		targetedXname = append(targetedXname, v.Hostname)
	}

	if failures > 0 {
		data.ErrResponse.E = -1
		data.ErrResponse.ErrMsg = fmt.Sprintf("Errors encountered with %d components for %s",
			failures, command)

		return data
	}
	// Grab the new list just in case a power off was done on a chassis
	// or compute module
	targetedXname, err := d.reserveComponents(targetedXname, command)
	defer d.releaseComponents(targetedXname)

	if err != nil {
		errstr := fmt.Sprintf("Error: Failed to reserve components while performing a %s.", command)
		log.Printf(errstr)
		data.ErrResponse.E = 37 // ENOLCK
		data.ErrResponse.ErrMsg = errstr
		return data
	}

	// Get lock status information
	res := d.reservation.Status()

	var tReq PCSTransition

	// Build PCS payload with deputy keys
	for _, x := range targetedXname {
		comp := PCSLocation{
			Xname:     x,
			DeputyKey: res[x].DeputyKey,
		}
		tReq.Location = append(tReq.Location, comp)
	}

	totalWait := len(tReq.Location)

	// If Off or Reinit, call PCS Off
	if command == bmcCmdPowerOff || command == bmcCmdPowerRestart ||
		command == bmcCmdPowerForceOff || command == bmcCmdPowerForceRestart {
		tReq.Operation = "off"
		if command == bmcCmdPowerForceOff || command == bmcCmdPowerForceRestart {
			tReq.Operation = "force-off"
		}

		failures, data = powerFunction(tReq, data, d, command, failures)
	}

	// If On or Reinit, call PCS On
	if (command == bmcCmdPowerOn || command == bmcCmdPowerRestart ||
		command == bmcCmdPowerForceOn || command == bmcCmdPowerForceRestart) &&
		data.E == 0 {
		tReq.Operation = "on"

		failures, data = powerFunction(tReq, data, d, command, failures)
	}

	if failures > 0 {
		data.ErrResponse.E = -1
		data.ErrResponse.ErrMsg = fmt.Sprintf("Errors encountered with %d/%d Xnames issued %s",
			failures, totalWait, command)
	}

	return data
}

func powerFunction(tReq PCSTransition, data capmc.XnameControlResponse, d *CapmcD, command string, failures int) (int, capmc.XnameControlResponse) {
	payload, err := json.Marshal(tReq)
	if err != nil {
		errstr := fmt.Sprintf("Error: Failed to marshal power request for PCS.")
		log.Printf(errstr)
		data.ErrResponse.E = http.StatusInternalServerError
		data.ErrResponse.ErrMsg = errstr
		return 0, data
	}

	url := d.pcsURL.String() + "/transitions"
	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		errstr := fmt.Sprintf("Error: Failed to create new request for power operation.")
		log.Printf(errstr)
		data.ErrResponse.E = http.StatusInternalServerError
		data.ErrResponse.ErrMsg = errstr
		return 0, data
	}

	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Content-Type", "application/json")

	body, err := d.doRequest(httpReq)
	if err != nil {
		errstr := fmt.Sprintf("Error: Failed to send power request to PCS.")
		log.Printf(errstr)
		data.ErrResponse.E = http.StatusInternalServerError
		data.ErrResponse.ErrMsg = errstr
		return 0, data
	}

	var tRsp PCSTransitionResponse
	err = json.Unmarshal(body, &tRsp)

	tID := tRsp.TransitionID

	url = d.pcsURL.String() + "/transitions/" + tID
	httpReq, err = http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		errstr := fmt.Sprintf("Error: Failed to create new request for power operation.")
		log.Printf(errstr)
		data.ErrResponse.E = http.StatusInternalServerError
		data.ErrResponse.ErrMsg = errstr
		return 0, data
	}

	httpReq.Header.Set("Accept", "application/json")

	// Arbitrary delay of 2 seconds
	time.Sleep(2 * time.Second)

	var tGet PCSTransitionGet

	retry := true

	// Wait for all to be !"New"
	for retry {
		body, err = d.doRequest(httpReq)
		if err != nil {
			errstr := fmt.Sprintf("Error: Failed to get transition from PCS.")
			log.Printf(errstr)
			data.ErrResponse.E = http.StatusInternalServerError
			data.ErrResponse.ErrMsg = errstr
			return 0, data
		}

		err = json.Unmarshal(body, &tGet)
		if err != nil {
			errstr := fmt.Sprintf("Error: Failed to unmarshal transition from PCS.")
			log.Printf(errstr)
			data.ErrResponse.E = http.StatusInternalServerError
			data.ErrResponse.ErrMsg = errstr
			return 0, data
		}

		if tGet.TransitionStatus != "new" && tGet.TaskCounts.New == 0 {
			retry = false
		} else {
			// Arbitrary delay of 2 seconds
			time.Sleep(2 * time.Second)
		}
	}

	// If the off portion of Reinit, wait for Off
	counts := tGet.TaskCounts
	if (command == bmcCmdPowerRestart || command == bmcCmdPowerForceRestart) &&
		(tReq.Operation == "off" || tReq.Operation == "force-off") {
		if (counts.Failed + counts.Succeeded + counts.UnSupported) < counts.Total {
			retry = true
		}
		for retry {
			body, err = d.doRequest(httpReq)
			if err != nil {
				errstr := fmt.Sprintf("Error: Failed to get transition from PCS.")
				log.Printf(errstr)
				data.ErrResponse.E = http.StatusInternalServerError
				data.ErrResponse.ErrMsg = errstr
				return 0, data
			}

			err = json.Unmarshal(body, &tGet)
			if err != nil {
				errstr := fmt.Sprintf("Error: Failed to unmarshal transition from PCS.")
				log.Printf(errstr)
				data.ErrResponse.E = http.StatusInternalServerError
				data.ErrResponse.ErrMsg = errstr
				return 0, data
			}

			counts = tGet.TaskCounts

			if (counts.Failed + counts.Succeeded + counts.UnSupported) == counts.Total {
				retry = false
			} else {
				// Arbitrary delay of 2 seconds
				time.Sleep(2 * time.Second)
			}

		}
	}

	// Check for failed components
	failures = counts.Failed + counts.UnSupported
	if failures > 0 {
		for _, task := range tGet.Tasks {
			if task.TaskStatus == "Failed" || task.TaskStatus == "Un-supported" {
				xnameErr := capmc.MakeXnameError(task.Xname, -1, task.TaskStatusDesc)
				data.Xnames = append(data.Xnames, xnameErr)
			}
		}
	}
	return failures, data
}

func (d *CapmcD) doCompStatus(nl []*NodeInfo, command string, filter uint) capmc.XnameStatusResponse {
	// The JSON encoder omits empty lists. The Cascade CAPMC API response
	// contains more lists than this, but at this time these are the only
	// ones that are reported.
	var data capmc.XnameStatusResponse
	data.On = make([]string, 0, 1)
	data.Off = make([]string, 0, 1)
	data.Undefined = make([]string, 0, 1)

	// Need a slice of xnames first

	xnames := make([]string, len(nl))
	for i, n := range nl {
		xnames[i] = n.Hostname
	}
	sort.Strings(xnames)

	// Then wrap it in a map

	payload := map[string][]string{
		"xname": xnames,
	}

	// Now marshal up the POST body

	postBody, merr := json.Marshal(&payload)
	if merr != nil {
		errstr := fmt.Sprintf("Error: Failed to marshal xnames to JSON.")
		log.Printf("%s", errstr)
		data.ErrResponse.E = http.StatusInternalServerError
		data.ErrResponse.ErrMsg = errstr
		return data
	}

	// And create the request

	url := d.pcsURL.String() + "/power-status"
	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer([]byte(postBody)))
	if err != nil {
		errstr := fmt.Sprintf("Error: Failed to create new request for power operation.")
		log.Printf(errstr)
		data.ErrResponse.E = http.StatusInternalServerError
		data.ErrResponse.ErrMsg = errstr
		return data
	}

	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Content-Type", "application/json")

	// Do the request

	var sGet PCSStatusGet

	body, err := d.doRequest(httpReq)
	if err != nil {
		errstr := fmt.Sprintf("Error: Failed to get status from PCS.")
		log.Printf(errstr)
		data.ErrResponse.E = http.StatusInternalServerError
		data.ErrResponse.ErrMsg = errstr
		return data
	}

	err = json.Unmarshal(body, &sGet)
	if err != nil {
		errstr := fmt.Sprintf("Error: Failed to unmarshal status from PCS.")
		log.Printf(errstr)
		data.ErrResponse.E = http.StatusInternalServerError
		data.ErrResponse.ErrMsg = errstr
		return data
	}

	var failures int
	for _, s := range sGet.Status {
		if s.Error != "" {
			data.Undefined = append(data.Undefined, s.Xname)
			failures++
			continue
		}
		switch s.PowerState {
		case strings.ToLower(rf.POWER_STATE_ON):
			if filter&capmc.FilterShowOnBit != 0 {
				data.On = append(data.On, s.Xname)
			}
		case strings.ToLower(rf.POWER_STATE_OFF):
			if filter&capmc.FilterShowOffBit != 0 {
				data.Off = append(data.Off, s.Xname)
			}
		// Other hardware states are not implemented
		default:
			// These are reported if there's a problem communicating
			// with the BMC.
			data.Undefined = append(data.Undefined, s.Xname)
		}

	}

	if failures > 0 {
		data.ErrResponse.E = -1
		data.ErrResponse.ErrMsg =
			fmt.Sprintf("Errors encountered with %d/%d Xnames for %s",
				failures, len(sGet.Status), command)
	}

	// Sorting is mostly for convenience; not strictly needed for capmc.
	sort.Sort(xnameSlice{&data.On})
	sort.Sort(xnameSlice{&data.Off})
	sort.Sort(xnameSlice{&data.Undefined})

	return data
}
