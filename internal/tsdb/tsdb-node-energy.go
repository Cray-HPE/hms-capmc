// Copyright 2019-2020 Hewlett Packard Enterprise Development LP
//
// This file is the interface for CAPMC to communicate with a postgres database
// containing telemetry data.

package tsdb

import (
	"database/sql"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

// GetNodeEnergy - returns based on times and location
// TODO when jobID and apid are added this will need to be changed to have a dynamic selection created
// when we do finally have those values I recommend always querying based on time and location; and use jobid/apid to generate that range
func (d PostgresqlDB) GetNodeEnergy(req TimeBoundNodeRequest) (nodeEnergy *NodeEnergy, err error) {
	log.WithFields(log.Fields{"start": req.StartTime, "end": req.EndTime}).Debug("Starting function GetNodeEnergy")

	// split the mountain and river queries and make the respective calls
	mtnNodes, rvrNodes, _ := splitMountainRiverNodes(req)

	// if there are missing nodes and no river/mountain nodes create a 'no data' error
	if len(rvrNodes) == 0 && len(mtnNodes) == 0 {
		log.Warningf("No nodes found in data tables for time window")
		return nil, NewTimeBoundRequestError(req.StartTime, req.EndTime, "No data in the time window")
	}

	// make the queries on the found nodes
	var merr, rerr error
	var mne, rne *NodeEnergy
	if len(mtnNodes) > 0 {
		// execute the mountain query
		req.Nodes = mtnNodes
		mne, merr = d.getMountainNodeEnergy(req)
	} else {
		// record through error that no mountain nodes exist in time window
		// maybe an error, maybe not - don't know yet
		merr = NewTimeBoundRequestError(req.StartTime, req.EndTime, "No data in the time window")
	}

	if len(rvrNodes) > 0 {
		// execute the river query
		req.Nodes = rvrNodes
		rne, rerr = d.getRiverNodeEnergy(req)
	} else {
		// record through error that no river nodes exist in time window
		// maybe an error, maybe not - don't know yet
		rerr = NewTimeBoundRequestError(req.StartTime, req.EndTime, "No data in the time window")
	}

	// Figuring out if there really is an error is a little complicated...
	err = processMtnRvrErrors(merr, rerr)

	// if there is no real error, then process the returned information
	// NOTE: if there is no error, there should be counter info present
	nodeEnergy = nil
	if err == nil {
		// if there is mountain data, just take it
		if mne != nil {
			nodeEnergy = mne
		}

		// if there is river data but no mountain data, take that
		if rne != nil && nodeEnergy == nil {
			nodeEnergy = rne
		} else if rne != nil && nodeEnergy != nil {
			// add the river data to the mountain data
			*nodeEnergy.NodeCount += *rne.NodeCount
			nodeEnergy.NodeLevels = append(nodeEnergy.NodeLevels, rne.NodeLevels...)
		}
	}

	return nodeEnergy, err
}

// getRiverNodeEnergy - returned node based energy for river hardware
func (d PostgresqlDB) getRiverNodeEnergy(req TimeBoundNodeRequest) (nodeEnergy *NodeEnergy, err error) {
	log.WithFields(log.Fields{"start": req.StartTime, "end": req.EndTime}).Trace("Starting function getRiverNodeEnergy")

	// NOTE: for the river telemetry data, all we get is the power (in watts) at time
	//  of the sample.  To turn this into energy, we need to multiply the average power
	//  by the timespan of the sample (fast and dirty integration of the power over
	//  the sample time) - the math is Watts * Seconds = Joules.

	// NOTE: at this time, the only available river data is showing sensor_type=Power and
	//  physical_context=Chassis as the filled in fields in the river_view table.  This
	//  still needs to be verified it is correct with both intel and gigabyte hardware.

	// get the duration of the request window in seconds
	dur := strconv.FormatFloat(req.EndTime.Sub(req.StartTime).Seconds(), 'f', -1, 64)

	// assemble sql statement with xnames
	xnameList := generateXnameString(req.Nodes, true)
	q := `
SELECT location, (avg(value) * ` + dur + `)::INTEGER as energy
	FROM pmdb.river_view
	WHERE sensor_type = 'Power'
		AND (physical_context = 'Chassis' OR physical_context = 'PowerSupplyBay' OR physical_context = 'Intake')
		AND timestamp between $1 and $2
		AND location IN (` + xnameList + `)
	GROUP BY river_view.location;
`
	log.WithFields(log.Fields{"SQL": q,
		"$1": req.StartTime,
		"$2": req.EndTime,
		"$3": xnameList}).Trace("Preparing statement")

	// prepare the query
	stmt, err := DB.Prepare(q)
	if err != nil {
		log.Error(err)
		return nodeEnergy, err
	}
	defer stmt.Close()

	// Execute the query
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
		// NOTE - query is not in a good state, so bail here
		return nodeEnergy, err
	}

	// Scan results into the temporary array struct
	nodeEnergy = new(NodeEnergy)
	for rows.Next() {
		var tmp NodeEnergyLevel
		err = rows.Scan(&tmp.Node.Xname, &tmp.Energy)
		if err != nil {
			// log the error on this scan line, then conintinue processing
			log.Warningf("Error retrieving query scan line: %s", err.Error())
			continue
		}
		log.Tracef("retrieving river data %s, %d", tmp.Node.Xname, tmp.Energy)

		// add to the list of nodes
		nodeEnergy.NodeLevels = append(nodeEnergy.NodeLevels, tmp)
	}

	// Make sure something was found
	count := len(nodeEnergy.NodeLevels)
	if count == 0 {
		// make error with no information in time window message
		err = NewTimeBoundRequestError(req.StartTime, req.EndTime, "No data in the time window")
		log.Warningf("TSDB: %s", err.Error())
		return nil, err
	}

	// Fill NodeEnergy struct
	nodeEnergy.NodeCount = &count
	nodeEnergy.TimeDuration = req.EndTime.Sub(req.StartTime)
	for i := range nodeEnergy.NodeLevels {
		for _, nodes := range req.Nodes {
			// TODO - this may have to change after the xnames are correctly flowing
			//  through the telemetry data
			if strings.Contains(nodes.Xname, nodeEnergy.NodeLevels[i].Node.Xname) {
				// swap out the truncated xname and nid info for the full version
				nodeEnergy.NodeLevels[i].Node.Nid = nodes.Nid
				nodeEnergy.NodeLevels[i].Node.Xname = nodes.Xname
			}
		}
	}

	return nodeEnergy, nil
}

// getMountainNodeEnergy - returned node based energy for mountain hardware
func (d PostgresqlDB) getMountainNodeEnergy(req TimeBoundNodeRequest) (nodeEnergy *NodeEnergy, err error) {
	log.WithFields(log.Fields{"start": req.StartTime, "end": req.EndTime}).Trace("Starting function getMountainNodeEnergy")

	// assemble sql statement with xnames
	xnameList := generateXnameString(req.Nodes, false)
	q := `
SELECT location, last(value, timestamp) - first(value, timestamp) as delta_e
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
GROUP BY nc_view.location;
`

	log.WithFields(log.Fields{"SQL": q,
		"$1": req.StartTime,
		"$2": req.EndTime,
		"$3": xnameList}).Trace("Preparing statement")

	// prepare the query
	stmt, err := DB.Prepare(q)
	if err != nil {
		log.Error(err)
		return nodeEnergy, err
	}
	defer stmt.Close()

	// Execute the query
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

	nodeEnergy = new(NodeEnergy)

	// Scan results into the temporary array struct
	for rows.Next() {
		var tmp NodeEnergyLevel
		err = rows.Scan(&tmp.Node.Xname, &tmp.Energy)
		log.Tracef("retrieving data %v", tmp)

		if err != nil {
			log.Error(err)
			return nil, err
		}

		nodeEnergy.NodeLevels = append(nodeEnergy.NodeLevels, tmp)
	}

	// Fill NodeEnergy struct
	count := len(nodeEnergy.NodeLevels)
	if count == 0 {
		// make error with no information in time window message
		err = NewTimeBoundRequestError(req.StartTime, req.EndTime, "No data in the time window")
		log.Warningf("TSDB: %s", err.Error())
		return nil, err
	}

	// fill out the return data
	nodeEnergy.NodeCount = &count
	nodeEnergy.TimeDuration = req.EndTime.Sub(req.StartTime)
	for i := range nodeEnergy.NodeLevels {
		for _, nodes := range req.Nodes {
			if nodes.Xname == nodeEnergy.NodeLevels[i].Node.Xname {
				nodeEnergy.NodeLevels[i].Node.Nid = nodes.Nid
			}
		}
	}

	return nodeEnergy, nil
}
