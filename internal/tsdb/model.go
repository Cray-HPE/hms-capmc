// Copyright (c) 2019 Cray Inc. All Rights Reserved.
//
// Except as permitted by contract or express written permission of Cray Inc.,
// no part of this work or its content may be modified, used, reproduced or
// disclosed in any form. Modifications made without express permission of
// Cray Inc. may damage the system the software is installed within, may
// disqualify the user from receiving support from Cray Inc. under support or
// maintenance contracts, or require additional support services outside the
// scope of those contracts to repair the software or system.

package tsdb

import (
	"fmt"
	"time"
)

const (
	BadEndTime          = "BAD_END_TIME"
	BadOptions          = "BAD_OPTIONS"
	BadStartTime        = "BAD_START_TIME"
	BadWindowLen        = "BAD_WINDOW_LEN"
	InvalidArguments    = "INVALID_ARGUMENTS"
	NoData              = "NO_DATA"
	NoResults           = "NO_RESULTS"
	PSQLFailure         = "PSQL_FAILURE"
	WindowLenOutOfRange = "WINDOW_LEN_OUT_OF_RANGE"
)

// TimeBoundRequestError - Error response for no data from time bound query
type TimeBoundRequestError struct {
	TimeBoundRequest
	ErrMsg string
}

// Implement the Error interface
func (e *TimeBoundRequestError) Error() string {
	return fmt.Sprintf("Error: %s, Start Time: %s, End Time: %s", e.ErrMsg,
		e.StartTime.String(), e.EndTime.String())
}

// NewTimeBoundRequestError - make a new object
func NewTimeBoundRequestError(st, et time.Time, eMsg string) *TimeBoundRequestError {
	var retVal TimeBoundRequestError = TimeBoundRequestError{
		TimeBoundRequest: TimeBoundRequest{StartTime: st, EndTime: et},
		ErrMsg:           "No data in time window",
	}
	return &retVal
}

// NodeLookup - pair of Xname and Nid
type NodeLookup struct {
	Xname string
	Nid   int
}

// TimeBoundNodeRequest - the general time based/job query params for TSDB
// for v1.1 Apid and JobId will always be empty; Once SMA starts putting these fields in PMDB, then we can use them
type TimeBoundNodeRequest struct {
	StartTime time.Time
	EndTime   time.Time
	Nodes     []NodeLookup
	JobId     string
	Apid      string
}

// TimeBoundRequest - the general time based/job query params for TSDB
type TimeBoundRequest struct {
	StartTime time.Time
	EndTime   time.Time
}

// NodeEnergyLevel - Pair of Node and Energy
type NodeEnergyLevel struct {
	Node   NodeLookup
	Energy int
}

// NodeEnergyStats - TSDB return struct for node energy stats
type NodeEnergyStats struct {
	EnergyTotal  *int
	EnergyAvg    *float64
	EnergyStd    *float64
	EnergyMax    *NodeEnergyLevel
	EnergyMin    *NodeEnergyLevel
	TimeDuration time.Duration
	NodeCount    *int
}

// NodeEnergy - TSDB return struct for node energy
type NodeEnergy struct {
	NodeCount    *int
	TimeDuration time.Duration
	NodeLevels   []NodeEnergyLevel
}

// NodeEnergyCounter - Used by TSDB to create NodeEnergyCounters return struct
type NodeEnergyCounter struct {
	Node       NodeLookup
	EnergyCtr  *int
	SampleTime time.Time
}

// NodeEnergyCounters - TSDB return struct for node energy counters
type NodeEnergyCounters struct {
	NodeCount *int
	Nodes     []NodeEnergyCounter
}

// SystemPower - TSDB return struct for GetSystemPower
type SystemPower struct {
	Min *int
	Max *int
	Avg *int
}

type CabinetPower struct {
	Avg         *float64
	Max         *int
	Min         *int
	ComponentID string
}

type SystemPowerByCabinet struct {
	Cabinets []*CabinetPower
}
