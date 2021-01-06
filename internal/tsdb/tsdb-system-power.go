// Copyright 2019-2020 Hewlett Packard Enterprise Development LP
//
// This file is the interface for CAPMC to communicate with a postgres database
// containing telemetry data.

package tsdb

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
)

const (
	//GetSystemPowerSQL - sql query to GetSystemPower from Mountain data
	GetSystemPowerSQL = `
SELECT  min(value)::INTEGER , max(value)::INTEGER , avg(value)::INTEGER
FROM (
    SELECT
    base.time       AS time,
    'Power'         AS metric,
    sum(base.value) AS value
FROM
    (
        SELECT
            time_bucket('15s', timestamp) AS time,
            location,
            index,
            avg(value)                   AS value
        FROM
            pmdb.cc_view
        WHERE
              physical_context = 'Rectifier'
          AND sensor_type = 'Power'
          AND physical_sub_context = 'Input'
          AND timestamp between $1 and $2
        GROUP BY time, location, index
        UNION ALL
        SELECT
            time_bucket('15s', timestamp) AS time,
            location,
            index,
            avg(value)                   AS value
        FROM
            pmdb.river_view
        WHERE
			  (physical_context = 'Chassis' OR physical_context = 'PowerSupplyBay' OR physical_context = 'Intake')
          AND sensor_type = 'Power'
          AND timestamp between $1 and $2
        GROUP BY time, location, index
    ) base
GROUP BY
    base.time ) vals ;
`
)

// GetSystemPower gets the system power from telemetry database
func (d PostgresqlDB) GetSystemPower(tbr TimeBoundRequest) (sysPow *SystemPower, err error) {
	log.WithFields(log.Fields{"start": tbr.StartTime, "end": tbr.EndTime}).Trace("Starting function GetSystemPower")

	log.WithFields(log.Fields{"SQL": GetSystemPowerSQL,
		"$1": tbr.StartTime,
		"$2": tbr.EndTime}).Trace("Preparing statement")
	stmt, err := DB.Prepare(GetSystemPowerSQL)
	if err != nil {
		log.Error(err)
		return sysPow, err
	}
	defer stmt.Close()

	sysPow = new(SystemPower)
	err = stmt.QueryRow(tbr.StartTime, tbr.EndTime).Scan(&sysPow.Min, &sysPow.Max, &sysPow.Avg)
	if err != nil {
		if err == sql.ErrNoRows {
			// convert to no data in time window error
			err = NewTimeBoundRequestError(tbr.StartTime, tbr.EndTime, "No data in the time window")
			log.Warningf("TSDB: %s", err.Error())
		} else {
			// log all other errors
			log.Error(err)
		}
	} else if sysPow.Min == nil && sysPow.Max == nil && sysPow.Avg == nil {
		// if there is no data from the query, it will return null for the
		// data pointers - need to turn this into an error for correct api
		// response
		// TODO - find out why this isn't throwing sql.ErrNoRows
		err = NewTimeBoundRequestError(tbr.StartTime, tbr.EndTime, "No data in the time window")
		log.Warningf("TSDB: %s", err.Error())
	}

	return sysPow, err
}
