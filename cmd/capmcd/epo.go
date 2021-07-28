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

	base "github.com/Cray-HPE/hms-base"
	"github.com/Cray-HPE/hms-capmc/internal/capmc"
)

func (d *CapmcD) validateXnamesForEPO(xnames []string) error {
	for _, x := range xnames {
		switch base.GetHMSType(x) {
		case base.CabinetPDU, base.Chassis, base.Cabinet,
			base.System, base.HMSTypeAll:
			// Valid component types that we handle for EPO
			// Nothing to do here
		default:
			// I should really generate a bad component list and return that
			return errors.New("Invalid component")
		}
	}
	return nil
}

func (d *CapmcD) genEPOList(xnames []string) ([]string, error) {
	var powerList []string
	var squery HSMQuery
	for _, xname := range xnames {
		switch base.GetHMSType(xname) {
		case base.System, base.HMSTypeAll, base.Cabinet:
			squery = HSMQuery{
				Types: []string{
					"chassis",
				},
			}
			comps, err := d.GetComponentsQuery(xname, getRestrictStr(squery))
			if err != nil {
				return nil, err
			}
			for _, comp := range comps {
				powerList = append(powerList, comp.ID)
			}
		case base.Chassis:
			powerList = append(powerList, xname)
		}
	}
	return powerList, nil
}

func (d *CapmcD) doEmergencyPowerOff(w http.ResponseWriter, r *http.Request) {
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
		if err == io.EOF {
			sendJsonError(w, http.StatusBadRequest, "no request")
		} else {
			sendJsonError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	data := d.executeEmergencyPowerOff(args)

	if data.E >= http.StatusMultipleChoices {
		sendJsonError(w, data.E, data.ErrMsg)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
	return
}

func (d *CapmcD) executeEmergencyPowerOff(args capmc.XnameControl) capmc.XnameControlResponse {
	// We will be expecting one or more of the following xnames:
	//   s0 - everything
	//   all - everything
	//   x0 - cabinet
	//   x0c0 - chassis
	//   x0m0 - iPDU
	//
	// A single xname of another type will result in an error and no power off
	// action will be taken.

	var data capmc.XnameControlResponse
	err := d.validateXnamesForEPO(args.Xnames)
	if err != nil {
		data.E = http.StatusBadRequest
		data.ErrMsg = fmt.Sprintf("failure validating xnames: %s",
			err.Error())
		return data
	}

	// Gen master xname list, expand s0, all, and x0
	// if s0 (System) or 'all' (HMSTypeAll) get all Chassis
	// if x0 (Cabinet) get all Chassis in x0
	// if x0c0 (Chassis) add it to the list
	powerList, err := d.genEPOList(args.Xnames)
	if err != nil {
		data.E = http.StatusBadRequest
		data.ErrMsg =
			fmt.Sprintf("failure generating component list: %s", err.Error())
		return data
	}

	var query HSMQuery
	if len(powerList) == 0 {
		data.E = http.StatusBadRequest
		data.ErrMsg = fmt.Sprintf("missing valid target xname(s)")
		return data
	}

	// Due to the previous validateXnamesforEPO, all of the components
	// in powerList will be valid at this point and dups won't matter
	query.ComponentIDs, _ = base.ValidateCompIDs(powerList, true)

	var nl []*NodeInfo
	nl, err = d.GetNodesByXname(query)

	if err != nil {
		if _, ok := err.(*InvalidCompIDsError); ok {
			data.E = http.StatusBadRequest
		} else {
			log.Printf("Error: %s", err)
			data.E = http.StatusInternalServerError
		}
		data.ErrMsg = err.Error()

		return data
	}

	if nl == nil {
		data.E = http.StatusNotFound
		data.ErrMsg = fmt.Sprintf("No components in machine")
		return data
	}

	log.Printf("Info: CAPMC Emergency Power Off - %s %v\n", args.Reason,
		query.ComponentIDs)

	cmap := make(map[string][]*NodeInfo)
	cmap[bmcCmdPowerForceOff] = nl

	var (
		failures  int
		totalWait int
	)

	for op, list := range cmap {
		waitNum := len(list)
		totalWait += waitNum
		waitCh := make(chan bmcPowerRc, waitNum)
		for _, ni := range list {

			// In this case the action URI is the EPO target
			ni.RfActionURI = ni.RfEpoURI

			// If the firmware is out of date the EPO target uri will not be
			// provided - use a hard-coded location and log it.
			// NOTE: this may be removed at some future date when it is
			// determined we can rely on correct firmware being present.
			if ni.RfEpoURI == "" {
				ni.RfActionURI = "/redfish/v1/Chassis/Enclosure/Actions/Oem/Chassis.EmergencyPower"
				log.Printf("Notice: Firmware may be out of date, using default URI for EPO on node: %s\n", ni.Hostname)
			}

			go d.doBmcCall(bmcCall{
				bmcCmd:  bmcCmd{cmd: op},
				ni:      ni,
				rspChan: waitCh,
			})
		}

		for i := 0; i < waitNum; i++ {
			result := <-waitCh
			if result.rc != 0 {
				failures++
				xnameErr := capmc.MakeXnameError(result.ni.Hostname,
					result.rc, result.msg)
				data.Xnames = append(data.Xnames, xnameErr)
			}
		}
	}

	if failures > 0 {
		data.E = -1
		data.ErrMsg =
			fmt.Sprintf("Errors encountered with %d/%d components",
				failures, totalWait)
	}
	return data
}
