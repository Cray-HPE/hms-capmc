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
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"

	log "github.com/sirupsen/logrus"
)

// IsDBErrorPg is error from postgres?
func IsDBErrorPg(err error) bool {
	_, ok := err.(*pq.Error)
	if !ok {
		return false
	}
	return true
}

// The Cascade API uses Python 'standard' time formating
// YYYY-MM-DD HH:MM:SS.MS which unfortunately isn't a default Go time format.
const psqlIntervalTimeFormat = "2006-01-02T15:04:05Z"

////////////////////////////////////////////////////////////////////////////
//
// TSDB PostgresqlDB Interface Implementations
//
////////////////////////////////////////////////////////////////////////////

// ImplementationName - Returns the database type
func (d PostgresqlDB) ImplementationName() string {
	return "Postgres"
}

// generateXnameString - put list of xnames together for sql query
func generateXnameString(req []NodeLookup, river bool) (str string) {
	var xnames []string
	for _, v := range req {
		//xnames = append(xnames, "'"+v.Xname+"'")
		// TODO: need to remove this when river location is fixed!!!
		if river {
			// need to strip the node info from xname to get nodeBMC for correct river location
			xnames = append(xnames, "'"+v.Xname[0:len(v.Xname)-2]+"'")
		} else {
			xnames = append(xnames, "'"+v.Xname+"'")
		}
	}

	str = strings.Join(xnames, ",")
	return str
}

// generateXnameStringMap - put list of xnames together for sql query from map
func generateXnameStringMap(req *map[string]int, river bool) (str string) {
	// NOTE: for the time being, the river list of xnames needs to pull off
	//  the node designation to get the nodeBMC xname since that is what is
	//  stored in the river_data table.  This may (hopefully will) change.
	//  This is ugly, but should be temporary.
	var xnames []string
	for key := range *req {
		// TODO: need to remove this when river location is fixed!!!
		if river {
			// need to strip the node info from xname to get nodeBMC for correct river location
			xnames = append(xnames, "'"+key[0:len(key)-2]+"'")
		} else {
			xnames = append(xnames, "'"+key+"'")
		}
	}

	str = strings.Join(xnames, ",")
	return str
}

// processMtnRvrErrors creates a single error from potentially mutiple error messages
func processMtnRvrErrors(merr, rerr error) error {
	// if there are no errors, then return nil
	if merr == nil && rerr == nil {
		return nil
	}

	// There is really only a no data error if both river and mountain queries
	// contain no data.  It is valid to have a river-only or mountain-only
	// syste in which case we expect one 'no data' error.
	var mtberr, rtberr *TimeBoundRequestError
	hasMtbErr := errors.As(merr, &mtberr)
	hasRtbErr := errors.As(rerr, &rtberr)
	if hasMtbErr && hasRtbErr {
		// both river and mountian have no data - return the error
		return mtberr
	}

	// If river OR mountain has no data, but the other is ok, suppress
	if (hasMtbErr && rerr == nil) || (hasRtbErr && merr == nil) {
		return nil
	}

	// some other kind of error - aggregate what we have and return it
	if merr != nil && rerr != nil {
		// report both errors
		return fmt.Errorf("TSDB errors- Mountain query: %s, River query: %s",
			merr.Error(), rerr.Error())
	}

	// if there is only a mountain error - return that one
	if merr != nil && rerr == nil {
		return merr
	}

	// only river error - report it
	return rerr
}

// splitMountainRiverNodes - splits out mountain, river, and unknown xnames
func splitMountainRiverNodes(tbr TimeBoundNodeRequest) (mtnNodes, rvrNodes, unknownNodes []NodeLookup) {
	// NOTE: this will just look at the database tables to find the xnames
	//  contained in the mountain tables and river tables for the given time
	//  interval.  If a name is not contained in either, we can't tell what
	//  type it is, but it won't have a result anyway...

	// TODO - do we need to do something with the river xnames to map to nodeBMC?
	// Or will this be fixed when the xnames are stored for location information
	//  instead of the ip address in the river_data table?

	// Convert array of NodeLookup objects to map for efficiency and ease of use
	m := make(map[string]int)
	for _, elem := range tbr.Nodes {
		m[elem.Xname] = elem.Nid
	}

	// split out the mountain nodes - removes found nodes from m
	mtnNodes, err := findMountainNodes(tbr.StartTime, tbr.EndTime, &m)
	if err != nil {
		log.Printf("Info: Mountain node split error: %s", err.Error())
	}
	// split out the river nodes - removes found nodes from m
	rvrNodes, err = findRiverNodes(tbr.StartTime, tbr.EndTime, &m)
	if err != nil {
		log.Printf("Info: River node split error: %s", err.Error())
	}

	// convert missing nodes back to array of NodeLookup objects
	for key, val := range m {
		unknownNodes = append(unknownNodes, NodeLookup{Xname: key, Nid: val})
	}

	// return the results
	log.Printf("Info: Node split found Mountain:%d, River:%d, Unknown:%d ",
		len(mtnNodes), len(rvrNodes), len(unknownNodes))
	return mtnNodes, rvrNodes, unknownNodes
}

// Find the mountain nodes from the input request
func findMountainNodes(stTime, endTime time.Time, m *map[string]int) (foundNodes []NodeLookup, err error) {
	// for the time being we will execute a time based query on the
	// mountain database table to find the xnames it contains

	xnameList := generateXnameStringMap(m, false)

	// make the sql query to search the mountain table
	// NOTE: we know details about what types of records are needed for all
	//  of the node based queries, so cook them into this search.  We may as
	//  well discard names without the data we want sooner rather than later.
	q := `
	SELECT distinct location
	FROM pmdb.nc_view
	WHERE physical_context = 'VoltageRegulator'
		AND parental_context = 'Chassis'
		AND physical_sub_context = 'Input'
		AND sensor_type = 'Energy'
		AND parental_index IS NULL
		AND device_specific_context IS NULL
		AND sub_index IS NULL
		AND timestamp between $1 and $2
		AND location IN (` + xnameList + `)`

	log.WithFields(log.Fields{"SQL": q,
		"$1": stTime,
		"$2": endTime,
		"$3": xnameList}).Trace("Mountain node query: Preparing statement")

	stmt, err := DB.Prepare(q)
	if err != nil {
		log.Error(err)
		return foundNodes, err
	}
	defer stmt.Close()

	// Execute the query
	log.Trace("Starting query")
	rows, err := stmt.Query(stTime, endTime)
	if err != nil {
		if err == sql.ErrNoRows {
			// convert to no data in time window error
			err = NewTimeBoundRequestError(stTime, endTime, "No data in the time window")
			log.Warningf("TSDB: %s", err.Error())
		} else {
			// log all other errors
			log.Error(err)
		}
	}

	// Scan results into the temporary array struct
	for rows.Next() {
		// grab the value
		var tmp string
		err = rows.Scan(&tmp)

		// process any error
		if err != nil {
			log.Error(err)
			return nil, err
		}

		// pull off the found value and remove from search map
		if nidVal, ok := (*m)[tmp]; ok {
			foundNodes = append(foundNodes, NodeLookup{Xname: tmp, Nid: nidVal})
			delete(*m, tmp)
		}
	}

	return foundNodes, nil
}

// Find the river nodes from the input request
func findRiverNodes(stTime, endTime time.Time, m *map[string]int) (foundNodes []NodeLookup, err error) {
	// if there is nothing left to look up bail early
	if len(*m) == 0 {
		return foundNodes, nil
	}

	// make the sql query to search the mountain table
	xnameList := generateXnameStringMap(m, true)

	// Really simple sql to grab all the location names in the time frame
	// NOTE: we may want to make this more specific to the types being used
	//  like the mountain based query does, but need to wait for the usage
	//  cases before doing so.
	q := `
	SELECT distinct location
	FROM pmdb.river_view
	WHERE timestamp between $1 and $2
		AND location IN (` + xnameList + `)`

	log.WithFields(log.Fields{"SQL": q,
		"$1": stTime,
		"$2": endTime,
		"$3": xnameList}).Trace("River node query: Preparing statement")

	stmt, err := DB.Prepare(q)
	if err != nil {
		log.Error(err)
		return foundNodes, err
	}
	defer stmt.Close()

	// Execute the query
	log.Trace("Starting query")
	rows, err := stmt.Query(stTime, endTime)
	if err != nil {
		if err == sql.ErrNoRows {
			// convert to no data in time window error
			err = NewTimeBoundRequestError(stTime, endTime, "No data in the time window")
			log.Warningf("TSDB: %s", err.Error())
		} else {
			// log all other errors
			log.Error(err)
		}
	}

	// Scan results into the temporary array struct
	for rows.Next() {
		// grab the value
		var tmp string
		err = rows.Scan(&tmp)

		// process any error
		if err != nil {
			log.Error(err)
			return nil, err
		}

		// pull off the found value and remove from search map
		if nidVal, ok := (*m)[tmp]; ok {
			foundNodes = append(foundNodes, NodeLookup{Xname: tmp, Nid: nidVal})
			delete(*m, tmp)
		} else {
			// TODO - look at this after location information fixed in river telemetry data
			//  this is a little brute-force now but re-address when river location fixed
			for key, val := range *m {
				if strings.Contains(key, tmp) {
					// NOTE: there may be more than one match so grab them all
					foundNodes = append(foundNodes, NodeLookup{Xname: key, Nid: val})
					delete(*m, key)
				}
			}
		}
	}

	return foundNodes, nil
}
