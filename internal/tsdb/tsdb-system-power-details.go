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

	log "github.com/sirupsen/logrus"
)

const (
	// GetSystemPowerDetailsSQL - returns system power by mountain cabinet
	//  NOTE: The mountain and river telemetry is very different from each
	//   other.  The mountain is collected per cabinet so is just averaged
	//   over a time bucket, then summed up.  The river has multiple
	//   controllers per cabinet.  The below query renames the location to
	//   the cabinet name when averaging the power per window, so they all
	//   get summed together during the 'group by' operation.
	GetSystemPowerDetailsSQL = `
SELECT min(value)::INTEGER
     , max(value)::INTEGER
     , avg(value)
     , split_part(location, 'c', 1) AS cid
FROM (
         SELECT base.time       AS time,
                'Power'         AS metric,
                sum(base.value) AS value,
                location
         FROM (
                  SELECT time_bucket('15s', timestamp) AS time,
                         location,
                         index,
                         avg(value)                    AS value
                  FROM pmdb.cc_view
                  WHERE physical_context = 'Rectifier'
                    AND sensor_type = 'Power'
                    AND physical_sub_context = 'Input'
                    AND timestamp between $1 AND $2
                  GROUP BY time, location, index
                  UNION ALL
                  SELECT time_bucket('15s', timestamp) AS time,
                         split_part(location, 'c', 1) as location,
                         index,
                         avg(value)                    AS value
                  FROM pmdb.river_view
				  WHERE
						  (physical_context = 'Chassis' OR physical_context = 'PowerSupplyBay' OR physical_context = 'Intake')
                    AND sensor_type = 'Power'
                    AND timestamp between $1 AND $2
                  GROUP BY time, location, index
              ) base
         GROUP BY base.time, location) vals
GROUP BY split_part(location, 'c', 1);
`
)

// GetSystemPowerDetails - return system power data per cabinet
func (d PostgresqlDB) GetSystemPowerDetails(tbr TimeBoundRequest) (cabPower *SystemPowerByCabinet, err error) {
	log.WithFields(log.Fields{"start": tbr.StartTime, "end": tbr.EndTime}).Trace("Starting function GetSystemPowerDetails")

	// load up the sql query
	stmt, err := DB.Prepare(GetSystemPowerDetailsSQL)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer stmt.Close()

	// execute the query
	rows, err := stmt.Query(tbr.StartTime, tbr.EndTime)
	if err != nil {
		if err == sql.ErrNoRows {
			// convert to no data in time window error
			err = NewTimeBoundRequestError(tbr.StartTime, tbr.EndTime, "No data in the time window")
			log.Warningf("TSDB: %s", err.Error())
		} else {
			// log all other errors
			log.Error(err)
		}
		return nil, err
	}

	// pull the data from the query
	cabPower = new(SystemPowerByCabinet)
	for rows.Next() {
		var (
			avg      *float64
			max, min *int
			cid      string
		)
		err := rows.Scan(&min, &max, &avg, &cid)
		if err != nil {
			log.Error(err)
			return nil, err
		}

		cabinet := new(CabinetPower)
		cabinet.Avg = avg
		cabinet.Max = max
		cabinet.Min = min
		cabinet.ComponentID = cid

		if err != nil {
			log.WithFields(log.Fields{"error": err, "component id": cid}).Warning("unexpected component id")
			// XXX skip it? shouldn't happen
			continue
		}
		cabPower.Cabinets = append(cabPower.Cabinets, cabinet)
	}

	// look for a database error
	err = rows.Err()
	if err != nil {
		log.Error(err)
		return cabPower, err
	}

	// if nothing was returned from the database log as error
	if len(cabPower.Cabinets) == 0 {
		err = NewTimeBoundRequestError(tbr.StartTime, tbr.EndTime, "No data in the time window")
		log.Warningf("TSDB: %s", err.Error())
		return nil, err
	}

	return cabPower, nil
}
