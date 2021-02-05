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
	"log"
	"net/http"

	"time"

	"stash.us.cray.com/HMS/hms-capmc/internal/capmc"
	"stash.us.cray.com/HMS/hms-capmc/internal/tsdb"
)

// validateTimeBoundNidRequestOptions checks that the TimeBoundNidRequest contains valid options
// Checks for minimally required optional fields/parameters. (Yes that is
// an oxymoron.)
// The parameters are all optional but there must be a Temporal element and a Geographical element:
// Therefore TEMPORALLY at least one of ...
//       * apid
//       * job_id
//       * start_time and end_time
// Therefore GEOGRAPHICALLY at least one of ...
//       * nids
//       * apid
//       * job_id
// as there are no defaults.
// This **does not** check that start_time and end_time are valid time format.
func validateTimeBoundNidRequestOptions(r capmc.TimeBoundNidRequest) (isValid bool, err error) {
	if (r.Apid != "" || r.JobId != "" || len(r.Nids) != 0) &&
		(r.Apid != "" || r.JobId != "" || (r.StartTime != "" && r.EndTime != "")) {
		return true, nil
	}
	return false, errors.New(InvalidArguments)
}

func temporarilyRestrictJobIDApid(r capmc.TimeBoundNidRequest) (isValid bool, problemArg string, err error) {
	if r.Apid != "" && r.JobId != "" {
		return false, "JOBID/APID", errors.New(ArgumentSupportNotImplemented)
	} else if r.Apid != "" {
		return false, "APID", errors.New(ArgumentSupportNotImplemented)
	} else if r.JobId != "" {
		return false, "JOBID", errors.New(ArgumentSupportNotImplemented)
	}
	return true, "", nil
}

// transformTimeBoundNidRequest - transforms a capmc.TimeBoundNidRequest to a tsdb.TimeBoundNodeRequest;
// does API validation, adjust timestamps for hysteresis,
// converts nids to NodeLookup via HSM query, and final TimeBoundNodeRequest validation
func (d *CapmcD) transformTimeBoundNidRequest(r capmc.TimeBoundNidRequest) (queryData *tsdb.TimeBoundNodeRequest, statusCode int, err error) {
	queryData = new(tsdb.TimeBoundNodeRequest)

	// Validate input request against what the doc says it should be.
	// We will later validate the data before it goes to the db.
	_, err = validateTimeBoundNidRequestOptions(r)
	if err != nil {
		// Error EINVAL
		return queryData, http.StatusBadRequest, err
	}

	// Right now Apid/JobId cannot be honored.  If the user sends them; we tell them NOT implemented
	_, problemArgs, err := temporarilyRestrictJobIDApid(r)
	if err != nil {
		return queryData, http.StatusNotImplemented, fmt.Errorf("Argument not supported yet: %s, err: %s", problemArgs, err)
	}

	// Build a TimeBoundNodeRequest in preparation for final validation and query exec
	queryData.StartTime, _ = time.Parse(intervalTimeFormat, r.StartTime)
	queryData.EndTime, _ = time.Parse(intervalTimeFormat, r.EndTime)
	// Make sure the timestamps are compliant with the hysteresis
	if !queryData.EndTime.IsZero() && !queryData.StartTime.IsZero() {
		// check for inverted start/end times
		if queryData.StartTime.After(queryData.EndTime) {
			log.Printf("Error: inverted time window - start time: %s, end time: %s",
				queryData.StartTime.String(), queryData.EndTime.String())
			return queryData, http.StatusBadRequest, fmt.Errorf("Inverted time window - start time: %s, end time: %s",
				queryData.StartTime.String(), queryData.EndTime.String())
		}

		// Adjust for hysteresis
		queryData.StartTime, queryData.EndTime, _ = tsdb.AdjustTimestampsForHystersisAndMinimumWindow(queryData.StartTime, queryData.EndTime, time.Now(), DefaultHysteresis, DefaultSampleWindow)
	}

	// Load Node data
	if len(r.Nids) > 0 {
		var query HSMQuery

		// strip out any duplicates; the JSON unmarshal will have
		// already handled any non-integer values
		query.NIDs, _ = validateNIDs(true, r.Nids)
		nodes, err := d.GetNidInfo(query)
		if err != nil {
			var ierr *InvalidNIDsError
			if errors.As(err, &ierr) {
				return queryData, http.StatusBadRequest, errors.New(InvalidArguments)
			} else {
				log.Printf("Error: %s", err)
				return queryData, http.StatusInternalServerError, errors.New("HSM_FAILURE")
			}
		}

		// copy nid/node to node lookup
		for _, node := range nodes {
			var nl = tsdb.NodeLookup{
				Xname: node.Xname,
				Nid:   node.Nid,
			}
			queryData.Nodes = append(queryData.Nodes, nl)
		}
	}

	isValid, err := tsdb.ValidateTimeBoundNodeRequest(*queryData)
	if isValid == false {
		return queryData, http.StatusBadRequest, err
	}

	return queryData, http.StatusOK, nil

}

// doNodeEnergy handles a node energy request
func (d *CapmcD) doNodeEnergy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		sendJsonError(w, http.StatusMethodNotAllowed,
			fmt.Sprintf("(%s) Not Allowed", r.Method))
		return
	}

	var args capmc.TimeBoundNidRequest

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&args)
	if err != nil {
		sendJsonError(w, http.StatusBadRequest,
			fmt.Sprintf("Bad Request: JSON: %s", err))
		return
	}

	queryData, statusCode, err := d.transformTimeBoundNidRequest(args)
	if err != nil {
		sendJsonError(w, statusCode, err.Error())
		return
	}

	// Build log message
	log.Printf("Info: Get Node Energy %v", queryData)

	nodeEnergy, err := tsdb.TSDBContext.GetNodeEnergy(*queryData)
	if err != nil {
		// Use the better message if one is sent
		var terr *tsdb.TimeBoundRequestError
		if errors.As(err, &terr) {
			log.Printf("Info: %s", terr.Error())
			SendResponseJSON(w, http.StatusOK, capmc.ErrResponse{E: http.StatusBadRequest, ErrMsg: terr.Error()})
			return
		}

		// stuck processing generic error
		log.Printf("Error: failed to retrieve data %s", err)
		sendJsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Attempted to set via '&' directly, but it doesnt work
	time := nodeEnergy.TimeDuration.Seconds()

	var data = capmc.GetNodeEnergyResponse{
		Time:        &time,
		NidCount:    nodeEnergy.NodeCount,
		ErrResponse: capmc.ErrResponse{E: 0, ErrMsg: ""},
	}

	for _, v := range nodeEnergy.NodeLevels {
		var nidEnergy capmc.NidEnergy
		nidEnergy.Energy = v.Energy
		nidEnergy.Nid = v.Node.Nid
		data.Nodes = append(data.Nodes, &nidEnergy)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
	return
}

// doNodeEnergyStats handles a node energy statistics request
func (d *CapmcD) doNodeEnergyStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		sendJsonError(w, http.StatusMethodNotAllowed,
			fmt.Sprintf("(%s) Not Allowed", r.Method))
		return
	}

	var args capmc.TimeBoundNidRequest

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&args)
	if err != nil {
		sendJsonError(w, http.StatusBadRequest,
			fmt.Sprintf("Bad Request: JSON: %s", err))
		return
	}

	queryData, statusCode, err := d.transformTimeBoundNidRequest(args)
	if err != nil {
		sendJsonError(w, statusCode, err.Error())
		return
	}

	// Build log message
	log.Printf("Info: Get Node Energy Stats %v", queryData)

	nodeEnergyStats, err := tsdb.TSDBContext.GetNodeEnergyStats(*queryData)
	if err != nil {
		// Use the better message if one is sent
		var terr *tsdb.TimeBoundRequestError
		if errors.As(err, &terr) {
			log.Printf("Info: %s", terr.Error())
			SendResponseJSON(w, http.StatusOK, capmc.ErrResponse{E: http.StatusBadRequest, ErrMsg: terr.Error()})
			return
		}

		// process the generic error
		log.Printf("Error: failed to retrieve data %s", err)
		sendJsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Attempted to set via '&' directly, but it doesnt work
	time := nodeEnergyStats.TimeDuration.Seconds()

	var data = capmc.GetNodeEnergyStatsResponse{
		EnergyTotal: nodeEnergyStats.EnergyTotal,
		EnergyStd:   nodeEnergyStats.EnergyStd,
		EnergyAvg:   nodeEnergyStats.EnergyAvg,
		NidCount:    nodeEnergyStats.NodeCount,
		Time:        &time,
		ErrResponse: capmc.ErrResponse{E: 0, ErrMsg: ""},
	}

	// DO NOT TOUCH THIS! The array MUST be ordered in this fashion;
	// This is a CAPMC holdover, so we have to honor the contract
	if nodeEnergyStats.EnergyMin != nil {
		data.EnergyMin = append(data.EnergyMin, &nodeEnergyStats.EnergyMin.Node.Nid)
		data.EnergyMin = append(data.EnergyMin, &nodeEnergyStats.EnergyMin.Energy)
	} else {
		var zero int = 0
		data.EnergyMin = append(data.EnergyMin, &zero)
		data.EnergyMin = append(data.EnergyMin, &zero)
	}
	if nodeEnergyStats.EnergyMax != nil {
		data.EnergyMax = append(data.EnergyMax, &nodeEnergyStats.EnergyMax.Node.Nid)
		data.EnergyMax = append(data.EnergyMax, &nodeEnergyStats.EnergyMax.Energy)
	} else {
		var zero int = 0
		data.EnergyMax = append(data.EnergyMax, &zero)
		data.EnergyMax = append(data.EnergyMax, &zero)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
	return
}

// doNodeEnergyCounter handles a node energy statistics request
func (d *CapmcD) doNodeEnergyCounter(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		sendJsonError(w, http.StatusMethodNotAllowed,
			fmt.Sprintf("(%s) Not Allowed", r.Method))
		return
	}

	var ncr capmc.GetNodeEnergyCounterRequest

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&ncr)
	if err != nil {
		sendJsonError(w, http.StatusBadRequest,
			fmt.Sprintf("Bad Request: JSON: %s", err))
		return
	}

	var nidRequest capmc.TimeBoundNidRequest

	// If time is empty, then set the time
	if ncr.Time == "" {
		log.Printf("no sample time, setting one")
		nidRequest.EndTime = time.Now().Format(intervalTimeFormat)
		nidRequest.StartTime = nidRequest.EndTime
	} else {
		nidRequest.EndTime = ncr.Time
		nidRequest.StartTime = ncr.Time
	}

	nidRequest.Nids = ncr.Nids
	nidRequest.Apid = ncr.Apid
	nidRequest.JobId = ncr.JobId

	queryData, statusCode, err := d.transformTimeBoundNidRequest(nidRequest)
	if err != nil {
		sendJsonError(w, statusCode, err.Error())
		return
	}
	// Build log message
	log.Printf("Info: Get Node Energy Counter %v", queryData)

	nodeEnergyCounters, err := tsdb.TSDBContext.GetNodeEnergyCounter(*queryData)
	if err != nil {
		// Use the better message if one is sent
		var terr *tsdb.TimeBoundRequestError
		if errors.As(err, &terr) {
			log.Printf("Info: %s", terr.Error())
			SendResponseJSON(w, http.StatusOK, capmc.ErrResponse{E: http.StatusBadRequest, ErrMsg: terr.Error()})
			return
		}

		// process the generic error
		log.Printf("Error: failed to retrieve data %s", err)
		sendJsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var data = capmc.GetNodeEnergyCounterResponse{
		NidCount:    nodeEnergyCounters.NodeCount,
		ErrResponse: capmc.ErrResponse{E: 0, ErrMsg: ""},
	}

	for _, v := range nodeEnergyCounters.Nodes {
		tmp := capmc.NidEnergyCounter{}
		tmp.Time = v.SampleTime.Format(intervalTimeFormat)
		tmp.EnergyCtr = v.EnergyCtr
		tmpNid := v.Node.Nid
		tmp.Nid = &tmpNid
		data.Nodes = append(data.Nodes, tmp)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
	return
}
