// Copyright (c) 2019 Cray Inc. All Rights Reserved.
//
// Except as permitted by contract or express written permission of Cray Inc.,
// no part of this work or its content may be modified, used, reproduced or
// disclosed in any form. Modifications made without express permission of
// Cray Inc. may damage the system the software is installed within, may
// disqualify the user from receiving support from Cray Inc. under support or
// maintenance contracts, or require additional support services outside the
// scope of those contracts to repair the software or system.
//
// This file is the interface for CAPMC to communicate with a postgres database
// containing telemetry data.

package tsdb

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
)

// GetNodeEnergyCounter - returns node engery counters for input request
func (d PostgresqlDB) GetNodeEnergyCounter(req TimeBoundNodeRequest) (counter *NodeEnergyCounters, err error) {
	log.WithFields(log.Fields{"start": req.StartTime, "end": req.EndTime}).Debug("Starting function GetNodeEnergyCounter")

	// split the mountain and river queries and make the respective calls
	mtnNodes, rvrNodes, missingNodes := splitMountainRiverNodes(req)

	// if there are missing nodes and no river/mountain nodes create a 'no data' error
	if len(missingNodes) > 0 && len(rvrNodes) == 0 && len(mtnNodes) == 0 {
		log.Warningf("No nodes found in data tables for time window")
		return nil, NewTimeBoundRequestError(req.StartTime, req.EndTime, "No data in the time window")
	}

	// make the queries on the found nodes
	var merr, rerr error
	var mNec, rNec *NodeEnergyCounters
	if len(mtnNodes) > 0 {
		// execute the mountain query
		req.Nodes = mtnNodes
		mNec, merr = d.getMountainNodeEnergyCounter(req)
	} else {
		// record through error that no mountain nodes exist in time window
		// maybe an error, maybe not - don't know yet
		merr = NewTimeBoundRequestError(req.StartTime, req.EndTime, "No data in the time window")
	}

	if len(rvrNodes) > 0 {
		// execute the river query
		req.Nodes = rvrNodes
		rNec, rerr = d.getRiverNodeEnergyCounter(req)
	} else {
		// record through error that no river nodes exist in time window
		// maybe an error, maybe not - don't know yet
		rerr = NewTimeBoundRequestError(req.StartTime, req.EndTime, "No data in the time window")
	}

	// Figuring out if there really is an error is a little complicated...
	err = processMtnRvrErrors(merr, rerr)

	// if there is no real error, then process the returned information
	// NOTE: if there is no error, there should be counter info present
	counter = nil
	if err == nil {
		// if there is mountain data, just take it
		if mNec != nil {
			counter = mNec
		}

		// if there is river data but no mountain data, take that
		if rNec != nil && counter == nil {
			counter = rNec
		} else if rNec != nil && counter != nil {
			// add the river data to the mountain data
			*counter.NodeCount += *rNec.NodeCount
			counter.Nodes = append(counter.Nodes, rNec.Nodes...)
		}
	}

	return counter, err
}

// getRiverNodeEnergyCounter - get the node energy counters for the river nodes
func (d PostgresqlDB) getRiverNodeEnergyCounter(req TimeBoundNodeRequest) (counter *NodeEnergyCounters, err error) {
	log.WithFields(log.Fields{"start": req.StartTime, "end": req.EndTime}).Trace("Starting function getRiverNodeEnergyCounter")

	// NOTE: Hmmmm - looks like the river telemetry data does not contain a
	//  counter or an equivalent.  Only instant power readings are stored, which
	//  means we can calculate energy usage over time, but there isn't a way to
	//  fudge a consistent counter like the mountain systems have.

	log.Infof("TSDB: getRiverNodeEnergyCounter not implemented")
	return nil, NewTimeBoundRequestError(req.StartTime, req.EndTime, "No data in the time window")
}

// getMountainNodeEnergyCounter - returns node energy counters for mountain hardware
func (d PostgresqlDB) getMountainNodeEnergyCounter(req TimeBoundNodeRequest) (counter *NodeEnergyCounters, err error) {
	log.WithFields(log.Fields{"start": req.StartTime, "end": req.EndTime}).Trace("Starting function getMountainNodeEnergyCounter")

	// assemble the sql query with generated xname list
	xnameList := generateXnameString(req.Nodes, false)
	q := `
	select sum(v), location, max(t)
	FROM (
		SELECT location, last(value, timestamp) as v, last(timestamp, timestamp) as t, index
		FROM pmdb.nc_view
		WHERE physical_context = 'VoltageRegulator'
			  AND parental_context = 'Chassis'
			  AND physical_sub_context = 'Input'
			  AND sensor_type = 'Energy'
			  AND parental_index IS NULL
			AND device_specific_context IS NULL
			AND sub_index IS NULL
			AND timestamp between $1 and $2
			AND location IN (` + xnameList + `)
		group by nc_view.location, nc_view.index) as t1
	group by location
	order by 2, 3
`
	log.WithFields(log.Fields{"SQL": q,
		"$1": req.StartTime,
		"$2": req.EndTime,
		"$3": xnameList}).Trace("Preparing statement")

	// prepare the query
	stmt, err := DB.Prepare(q)
	if err != nil {
		log.Error(err)
		return counter, err
	}
	defer stmt.Close()

	// Execute the query
	log.Trace("Starting query")
	rows, err := stmt.Query(req.StartTime, req.EndTime)
	if err != nil {
		if err == sql.ErrNoRows {
			// convert to no data in time window error
			err = NewTimeBoundRequestError(req.StartTime, req.EndTime, "No data in the time window")
			log.Warningf("TSDB: %s", err.Error())
		} else {
			// log all other errors
			log.Error(err)
		}
	}

	counter = new(NodeEnergyCounters)
	// Scan results into the temporary array struct
	for rows.Next() {
		var tmp NodeEnergyCounter
		err = rows.Scan(&tmp.EnergyCtr, &tmp.Node.Xname, &tmp.SampleTime)
		log.Tracef("retrieving data %v", tmp)

		if err != nil {
			log.Error(err)
			return nil, err
		}

		counter.Nodes = append(counter.Nodes, tmp)
	}

	// see if anything was found
	count := len(counter.Nodes)
	if count == 0 {
		// generate an error for no data in the requested window
		err = NewTimeBoundRequestError(req.StartTime, req.EndTime, "No data in the time window")
		log.Warningf("TSDB: %s", err.Error())
		return nil, err
	}

	// Fill NodeEnergy struct
	counter.NodeCount = &count

	// Fill in nid information
	for i := range counter.Nodes {
		for _, nodes := range req.Nodes {
			if nodes.Xname == counter.Nodes[i].Node.Xname {
				counter.Nodes[i].Node.Nid = nodes.Nid
			}
		}
	}

	return counter, nil
}
