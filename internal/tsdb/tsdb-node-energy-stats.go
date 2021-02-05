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
// This file is the interface for CAPMC to communicate with a postgres database
// containing telemetry data.

package tsdb

import (
	"database/sql"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

// GetNodeEnergyStats - returns node energy stats for the input time period
func (d PostgresqlDB) GetNodeEnergyStats(req TimeBoundNodeRequest) (stat *NodeEnergyStats, err error) {
	log.WithFields(log.Fields{"start": req.StartTime, "end": req.EndTime}).Debug("Starting function GetNodeEnergyStats")

	// split the mountain and river queries and make the respective calls
	mtnNodes, rvrNodes, missingNodes := splitMountainRiverNodes(req)

	// if there are missing nodes and no river/mountain nodes create a 'no data' error
	if len(missingNodes) > 0 && len(rvrNodes) == 0 && len(mtnNodes) == 0 {
		log.Warningf("No nodes found in data tables for time window")
		return nil, NewTimeBoundRequestError(req.StartTime, req.EndTime, "No data in the time window")
	}

	// make the queries on the found nodes
	var merr, rerr error
	var mNec, rNec *NodeEnergyStats
	if len(mtnNodes) > 0 {
		// execute the mountain query
		req.Nodes = mtnNodes
		mNec, merr = d.getMountainNodeEnergyStats(req)
	} else {
		// record through error that no mountain nodes exist in time window
		// maybe an error, maybe not - don't know yet
		merr = NewTimeBoundRequestError(req.StartTime, req.EndTime, "No data in the time window")
	}

	if len(rvrNodes) > 0 {
		// execute the river query
		req.Nodes = rvrNodes
		rNec, rerr = d.getRiverNodeEnergyStats(req)
	} else {
		// record through error that no river nodes exist in time window
		// maybe an error, maybe not - don't know yet
		rerr = NewTimeBoundRequestError(req.StartTime, req.EndTime, "No data in the time window")
	}

	// Figuring out if there really is an error is a little complicated...
	err = processMtnRvrErrors(merr, rerr)

	// if there is no real error, then process the returned information
	// NOTE: if there is no error, there should be counter info present
	stat = nil
	if err == nil {
		// if there is mountain data, just take it
		if mNec != nil {
			stat = mNec
		}

		// if there is river data but no mountain data, take that
		if rNec != nil && stat == nil {
			stat = rNec
		} else if rNec != nil && stat != nil {
			// add the river data to the mountain data
			*stat.EnergyTotal += *rNec.EnergyTotal

			// for avg and std weight by the number of nodes
			totNodes := *stat.NodeCount + *rNec.NodeCount
			*stat.EnergyAvg = ((*stat.EnergyAvg)*float64(*stat.NodeCount) + (*rNec.EnergyAvg)*float64(*rNec.NodeCount)) / float64(totNodes)
			*stat.EnergyStd = ((*stat.EnergyStd)*float64(*stat.NodeCount) + (*rNec.EnergyStd)*float64(*rNec.NodeCount)) / float64(totNodes)

			// add the river nodes
			// NOTE: must do this AFTER weighted average combination since
			//  this value is used in the above calculation.
			*stat.NodeCount += *rNec.NodeCount

			// pick the min/max between the mountain/river systems
			if rNec.EnergyMax.Energy > stat.EnergyMax.Energy {
				stat.EnergyMax = rNec.EnergyMax
			}
			if rNec.EnergyMin.Energy < stat.EnergyMin.Energy {
				stat.EnergyMin = rNec.EnergyMin
			}
		}
	}

	return stat, err
}

// getRiverNodeEnergyStats returns node energy stats for the mountain hardware
func (d PostgresqlDB) getRiverNodeEnergyStats(req TimeBoundNodeRequest) (stat *NodeEnergyStats, err error) {
	log.WithFields(log.Fields{"start": req.StartTime, "end": req.EndTime}).Trace("Starting function getRiverNodeEnergyStats")

	// NOTE: see the documentation for getRiverNodeEnergy for details on calculation

	// get the duration of the request window in seconds
	dur := strconv.FormatFloat(req.EndTime.Sub(req.StartTime).Seconds(), 'f', -1, 64)

	// assemble sql query with xnames
	xnameList := generateXnameString(req.Nodes, true)
	q0 := `
	SELECT coalesce(SUM(delta_e),0)::INTEGER as total, coalesce(AVG(delta_e),0) as avg, coalesce(stddev(delta_e),0) as std
	FROM (
		SELECT location, avg(value) * ` + dur + ` as delta_e
		FROM pmdb.river_view
		WHERE
				(physical_context = 'Chassis' OR physical_context = 'PowerSupplyBay' OR physical_context = 'Intake')
			AND sensor_type = 'Power'
			AND timestamp between $1 and $2
			AND location IN (` + xnameList + `)
		group by river_view.location) t1;`

	log.WithFields(log.Fields{"SQL": q0,
		"$1": req.StartTime,
		"$2": req.EndTime,
		"$3": xnameList}).Trace("Preparing statement")

	// prepare the sql statement
	stmt, err := DB.Prepare(q0)
	if err != nil {
		log.Error(err)
		return stat, err
	}

	stat = new(NodeEnergyStats)

	// Execute the statement
	log.Trace("Starting query0")
	err = stmt.QueryRow(req.StartTime, req.EndTime).Scan(&stat.EnergyTotal, &stat.EnergyAvg, &stat.EnergyStd)
	if err != nil {
		log.Error(err)
	}
	stmt.Close()

	q1 := `
SELECT delta_e as val, location
FROM (
	SELECT location, (avg(value) *` + dur + `)::INTEGER as delta_e
	FROM pmdb.river_view
	WHERE
			(physical_context = 'Chassis' OR physical_context = 'PowerSupplyBay')
		AND sensor_type = 'Power'
		AND timestamp BETWEEN $1 and $2
		AND location IN (` + xnameList + `)
	GROUP BY river_view.location) as t1
ORDER BY val DESC`

	log.WithFields(log.Fields{"SQL": q1,
		"$1": req.StartTime,
		"$2": req.EndTime,
		"$3": xnameList}).Trace("Preparing statement")

	stmt, err = DB.Prepare(q1)
	if err != nil {
		log.Error(err)
		return stat, err
	}
	defer stmt.Close()

	log.Trace("Starting query1")
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
		// bail since the query is in a bad state
		return stat, err
	}

	log.Trace("Starting data scan query1")

	rowCount := 0
	stat.EnergyMax = new(NodeEnergyLevel)
	stat.EnergyMin = new(NodeEnergyLevel)

	// This is a bit hacky; but the list is ordered, so the first one is always max; and keep reseting the min to the current one.
	for rows.Next() {
		log.Trace("row found, retrieving data")
		if rowCount == 0 {
			_ = rows.Scan(&stat.EnergyMax.Energy, &stat.EnergyMax.Node.Xname)
		}
		_ = rows.Scan(&stat.EnergyMin.Energy, &stat.EnergyMin.Node.Xname)
		rowCount++
	}
	log.WithFields(log.Fields{"EnergyMax": stat.EnergyMax, "EnergyMin": stat.EnergyMin}).Trace("Values")

	// check if anything was found
	count := rowCount
	if rowCount == 0 {
		// gerenate an error for no data in the requested time window
		err = NewTimeBoundRequestError(req.StartTime, req.EndTime, "No data in the time window")
		log.Warningf("TSDB: %s", err.Error())
		return nil, err
	}

	// Fill NodeEnergyStat struct
	stat.NodeCount = &count
	stat.TimeDuration = req.EndTime.Sub(req.StartTime)

	// Fill in nid information
	for _, nodes := range req.Nodes {
		// TODO - this may have to change after the xnames are correctly flowing
		//  through the telemetry data
		if strings.Contains(nodes.Xname, stat.EnergyMax.Node.Xname) {
			stat.EnergyMax.Node.Nid = nodes.Nid
			stat.EnergyMax.Node.Xname = nodes.Xname
		}
		if strings.Contains(nodes.Xname, stat.EnergyMin.Node.Xname) {
			stat.EnergyMin.Node.Nid = nodes.Nid
			stat.EnergyMin.Node.Xname = nodes.Xname
		}
	}

	return stat, nil
}

// getMountainNodeEnergyStats returns node energy stats for the mountain hardware
func (d PostgresqlDB) getMountainNodeEnergyStats(req TimeBoundNodeRequest) (stat *NodeEnergyStats, err error) {
	log.WithFields(log.Fields{"start": req.StartTime, "end": req.EndTime}).Trace("Starting function getMountainNodeEnergyStats")

	// assemble sql query with xnames
	xnameList := generateXnameString(req.Nodes, false)
	q0 := `
SELECT coalesce(SUM(delta_e),0) as total, coalesce(AVG(delta_e),0) as avg, coalesce(stddev(delta_e),0) as std
FROM (
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
	group by nc_view.location) t1;`

	log.WithFields(log.Fields{"SQL": q0,
		"$1": req.StartTime,
		"$2": req.EndTime,
		"$3": xnameList}).Trace("Preparing statement")

	// prepare the sql statement
	stmt, err := DB.Prepare(q0)
	if err != nil {
		log.Error(err)
		return stat, err
	}

	stat = new(NodeEnergyStats)

	// Execute the statement
	log.Trace("Starting query0")
	err = stmt.QueryRow(req.StartTime, req.EndTime).Scan(&stat.EnergyTotal, &stat.EnergyAvg, &stat.EnergyStd)
	if err != nil {
		log.Error(err)
	}
	stmt.Close()

	q1 := `
SELECT delta_e as val, location
FROM (
	SELECT location, last(value, timestamp) - first(value, timestamp) as delta_e
	FROM pmdb.nc_view
	WHERE physical_context = 'VoltageRegulator'
		AND parental_context = 'Chassis'
		AND physical_sub_context = 'Input'
		AND sensor_type = 'Energy'
		AND parental_index IS NULL
		AND device_specific_context IS NULL
		AND sub_index IS NULL
		AND timestamp BETWEEN $1 and $2
		AND location IN (` + xnameList + `)
	GROUP BY nc_view.location) as t1
ORDER BY val DESC`

	log.WithFields(log.Fields{"SQL": q1,
		"$1": req.StartTime,
		"$2": req.EndTime,
		"$3": xnameList}).Trace("Preparing statement")

	stmt, err = DB.Prepare(q1)
	if err != nil {
		log.Error(err)
		return stat, err
	}
	defer stmt.Close()

	log.Trace("Starting query1")
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

	log.Trace("Starting data scan query1")

	rowCount := 0
	stat.EnergyMax = new(NodeEnergyLevel)
	stat.EnergyMin = new(NodeEnergyLevel)

	// This is a bit hacky; but the list is ordered, so the first one is always max; and keep reseting the min to the current one.
	for rows.Next() {
		log.Trace("row found, retrieving data")
		if rowCount == 0 {
			_ = rows.Scan(&stat.EnergyMax.Energy, &stat.EnergyMax.Node.Xname)
		}
		_ = rows.Scan(&stat.EnergyMin.Energy, &stat.EnergyMin.Node.Xname)
		rowCount++
	}
	log.WithFields(log.Fields{"EnergyMax": stat.EnergyMax, "EnergyMin": stat.EnergyMin}).Trace("Values")

	// check if anything was found
	count := rowCount
	if rowCount == 0 {
		// gerenate an error for no data in the requested time window
		err = NewTimeBoundRequestError(req.StartTime, req.EndTime, "No data in the time window")
		log.Warningf("TSDB: %s", err.Error())
		return nil, err
	}

	// Fill NodeEnergyStat struct
	stat.NodeCount = &count
	stat.TimeDuration = req.EndTime.Sub(req.StartTime)

	// Fill in nid information
	for _, nodes := range req.Nodes {
		if nodes.Xname == stat.EnergyMax.Node.Xname {
			stat.EnergyMax.Node.Nid = nodes.Nid
		}
		if nodes.Xname == stat.EnergyMin.Node.Xname {
			stat.EnergyMin.Node.Nid = nodes.Nid
		}
	}

	return stat, nil
}
