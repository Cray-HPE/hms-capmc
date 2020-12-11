// Copyright 2019 Cray Inc. All Rights Reserved.
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
	base "stash.us.cray.com/HMS/hms-base"
)

var e = base.NewHMSError("tsdb", "GenericError")

// TSDBContext - the TSDB interface global; used for setting context
var TSDBContext TSDB

// PostgresqlDB - Production Database interface for SQL version
type PostgresqlDB struct{}

// DummyDB - the in memory version of the fox data interface
type DummyDB struct{}

// TSDB - This is the interface for the database; this defines all supported methods
type TSDB interface {

	// GetSystemPower retrieves the system power min, max, and average
	// between the specified start and end time.
	GetSystemPower(tbr TimeBoundRequest) (sysPow *SystemPower, err error)

	//TODO convert this return to be something from TSDB
	// GetSystemPowerDetails retrieves the by cabinet min, max, and
	// average (system) power between the specified start and end time.
	GetSystemPowerDetails(tbr TimeBoundRequest) (cabPower *SystemPowerByCabinet, err error)

	// GetNodeEnergy retrieves node energy usage by time, jobID and node, all logically ANDed together
	GetNodeEnergy(req TimeBoundNodeRequest) (nodeEnergy *NodeEnergy, err error)

	// GetNodeEnergyStats retrieves node energy stats by time, jobID and node, all logically ANDed together
	GetNodeEnergyStats(req TimeBoundNodeRequest) (stat *NodeEnergyStats, err error)

	// GetNodeEnergyCounter retrieves node energy counter at sample time  (can be for jobID) jobID and node
	GetNodeEnergyCounter(req TimeBoundNodeRequest) (counter *NodeEnergyCounters, err error)

	// ImplementationName
	// Returns the database type
	ImplementationName() string
}

//GetSystemPower - returns default dummy values
func (d DummyDB) GetSystemPower(tbr TimeBoundRequest) (sysPow *SystemPower, err error) {
	sysPow = new(SystemPower)

	sysPow.Min = new(int)
	*sysPow.Min = 1

	sysPow.Max = new(int)
	*sysPow.Max = 100

	sysPow.Avg = new(int)
	*sysPow.Avg = 50

	return sysPow, nil
}

//GetSystemPowerDetails - returns default dummy values
func (d DummyDB) GetSystemPowerDetails(tbr TimeBoundRequest) (cabPower *SystemPowerByCabinet, err error) {
	cabPower = new(SystemPowerByCabinet)

	avg0 := 0.0
	max0 := 0
	min0 := 0

	var c0 = CabinetPower{&avg0, &max0, &min0, "x0"}
	cabPower.Cabinets = append(cabPower.Cabinets, &c0)

	avg1 := 1.5
	max1 := 3
	min1 := 1
	var c1 = CabinetPower{&avg1, &max1, &min1, "x1"}
	cabPower.Cabinets = append(cabPower.Cabinets, &c1)

	return cabPower, nil
}

// GetNodeEnergy - returns dummy values
func (d DummyDB) GetNodeEnergy(req TimeBoundNodeRequest) (nodeEnergy *NodeEnergy, err error) {
	nodeEnergy = new(NodeEnergy)

	nodeEnergy.NodeCount = new(int)
	*nodeEnergy.NodeCount = 2

	nodeEnergy.TimeDuration = req.StartTime.Sub(req.EndTime)

	var c = NodeEnergyLevel{NodeLookup{"x0c0s0b0n0", 0}, 100}
	var c1 = NodeEnergyLevel{NodeLookup{"x0c0s0b0n1", 1}, 1000}

	nodeEnergy.NodeLevels = append(nodeEnergy.NodeLevels, c)
	nodeEnergy.NodeLevels = append(nodeEnergy.NodeLevels, c1)

	return nodeEnergy, nil
}

// GetNodeEnergyStats - returns dummy values
func (d DummyDB) GetNodeEnergyStats(req TimeBoundNodeRequest) (stat *NodeEnergyStats, err error) {

	stat = new(NodeEnergyStats)

	stat.EnergyTotal = new(int)
	*stat.EnergyTotal = 1100

	stat.NodeCount = new(int)
	*stat.NodeCount = 2

	stat.EnergyAvg = new(float64)
	*stat.EnergyAvg = 550.0

	stat.EnergyStd = new(float64)
	*stat.EnergyStd = 450.0

	stat.TimeDuration = req.EndTime.Sub(req.StartTime)

	stat.EnergyMax = &NodeEnergyLevel{NodeLookup{"x0c0s0b0n1", 1}, 1000}
	stat.EnergyMin = &NodeEnergyLevel{NodeLookup{"x0c0s0b0n0", 0}, 100}

	return stat, nil
}

// GetNodeEnergyCounter - returns dummy values
func (d DummyDB) GetNodeEnergyCounter(req TimeBoundNodeRequest) (counter *NodeEnergyCounters, err error) {

	counter = new(NodeEnergyCounters)
	e1 := 100
	e2 := 1000

	counter.NodeCount = new(int)
	*counter.NodeCount = 2

	counter.Nodes = append(counter.Nodes, NodeEnergyCounter{NodeLookup{"x0c0s0b0n0", 0}, &e1, req.EndTime})
	counter.Nodes = append(counter.Nodes, NodeEnergyCounter{NodeLookup{"x0c0s0b0n1", 1}, &e2, req.EndTime})

	return counter, nil
}

//ImplementationName - returns "DUMMY"
func (d DummyDB) ImplementationName() string {

	return "DUMMY"
}
