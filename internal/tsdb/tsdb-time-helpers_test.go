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
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"stash.us.cray.com/HMS/hms-capmc/internal/logger"
)

type TSDB_Helpers_TS struct {
	suite.Suite
}

// SetupSuit is run ONCE
func (suite *TSDB_Helpers_TS) SetupSuite() {
}

//Both start and end are < NOW - hystersis
// start is <, end is > hystersis , but  less than now
//start is inside hystersis, end is > now
func (suite *TSDB_Helpers_TS) TestAdjustTimestampsForHystersisAndMinimumWindow_StartLessHyst_EndLessHyst() {

	var startTime, endTime, rightNow time.Time
	var hysteresis, window time.Duration

	hysteresis = -15 * time.Second
	window = 15 * time.Second

	layout := "2006-01-02 15:04:05Z"

	startTime, _ = time.Parse(layout, "2018-04-01 13:00:00Z")
	endTime, _ = time.Parse(layout, "2018-04-01 14:00:00Z")
	rightNow, _ = time.Parse(layout, "2018-04-01 15:00:00Z")

	actualStart, actualEnd, _ := AdjustTimestampsForHystersisAndMinimumWindow(startTime, endTime, rightNow, hysteresis, window)
	suite.Equal(startTime, actualStart)
	suite.Equal(endTime, actualEnd)
}

func (suite *TSDB_Helpers_TS) TestAdjustTimestampsForHystersisAndMinimumWindow_StartLessHyst_EndGreatHystLessNow() {

	var startTime, endTime, rightNow time.Time
	var hysteresis, window time.Duration

	hysteresis = -15 * time.Second
	window = 15 * time.Second

	layout := "2006-01-02 15:04:05Z"

	startTime, _ = time.Parse(layout, "2018-04-01 13:00:00Z")
	endTime, _ = time.Parse(layout, "2018-04-01 14:00:00Z")
	rightNow, _ = time.Parse(layout, "2018-04-01 14:00:05Z")

	intendedStart, _ := time.Parse(layout, "2018-04-01 12:59:50Z")
	intendedEnd, _ := time.Parse(layout, "2018-04-01 13:59:50Z")
	actualStart, actualEnd, _ := AdjustTimestampsForHystersisAndMinimumWindow(startTime, endTime, rightNow, hysteresis, window)
	suite.Equal(intendedStart, actualStart)
	suite.Equal(intendedEnd, actualEnd)
}

func (suite *TSDB_Helpers_TS) TestAdjustTimestampsForHystersisAndMinimumWindow_StartLessHyst_EndGreatNow() {

	var startTime, endTime, rightNow time.Time
	var hysteresis, window time.Duration

	hysteresis = -15 * time.Second
	window = 15 * time.Second

	layout := "2006-01-02 15:04:05Z"

	startTime, _ = time.Parse(layout, "2018-04-01 13:00:00Z")
	endTime, _ = time.Parse(layout, "2018-04-01 14:00:00Z")
	rightNow, _ = time.Parse(layout, "2018-04-01 13:30:00Z")

	intendedStart, _ := time.Parse(layout, "2018-04-01 12:29:45Z")
	intendedEnd, _ := time.Parse(layout, "2018-04-01 13:29:45Z")
	actualStart, actualEnd, _ := AdjustTimestampsForHystersisAndMinimumWindow(startTime, endTime, rightNow, hysteresis, window)
	suite.Equal(intendedStart, actualStart)
	suite.Equal(intendedEnd, actualEnd)
}

func (suite *TSDB_Helpers_TS) TestAdjustTimestampsForHystersisAndMinimumWindow_StartGreatHyst_EndGreatNow() {

	var startTime, endTime, rightNow time.Time
	var hysteresis, window time.Duration

	hysteresis = -15 * time.Second
	window = 15 * time.Second

	layout := "2006-01-02 15:04:05Z"

	startTime, _ = time.Parse(layout, "2018-04-01 13:29:50Z")
	endTime, _ = time.Parse(layout, "2018-04-01 14:00:00Z")
	rightNow, _ = time.Parse(layout, "2018-04-01 13:30:00Z")

	intendedStart, _ := time.Parse(layout, "2018-04-01 12:59:35Z")
	intendedEnd, _ := time.Parse(layout, "2018-04-01 13:29:45Z")
	actualStart, actualEnd, _ := AdjustTimestampsForHystersisAndMinimumWindow(startTime, endTime, rightNow, hysteresis, window)
	suite.Equal(intendedStart, actualStart)
	suite.Equal(intendedEnd, actualEnd)
}

func (suite *TSDB_Helpers_TS) TestAdjustTimestampsForHystersisAndMinimumWindow_StartGreatNow_EndGreatNow() {

	var startTime, endTime, rightNow time.Time
	var hysteresis, window time.Duration

	hysteresis = -15 * time.Second
	window = 15 * time.Second

	layout := "2006-01-02 15:04:05Z"

	startTime, _ = time.Parse(layout, "2018-04-01 13:01:00Z")
	endTime, _ = time.Parse(layout, "2018-04-01 13:02:00Z")
	rightNow, _ = time.Parse(layout, "2018-04-01 13:00:00Z")

	intendedStart, _ := time.Parse(layout, "2018-04-01 12:58:45Z")
	intendedEnd, _ := time.Parse(layout, "2018-04-01 12:59:45Z")
	actualStart, actualEnd, _ := AdjustTimestampsForHystersisAndMinimumWindow(startTime, endTime, rightNow, hysteresis, window)
	suite.Equal(intendedStart, actualStart)
	suite.Equal(intendedEnd, actualEnd)
}

func (suite *TSDB_Helpers_TS) TestAdjustTimestampsForHystersisAndMinimumWindow_StartGreatHys_EndGreatHysLessNow() {

	var startTime, endTime, rightNow time.Time
	var hysteresis, window time.Duration

	hysteresis = -15 * time.Second
	window = 15 * time.Second

	layout := "2006-01-02 15:04:05Z"

	startTime, _ = time.Parse(layout, "2018-04-01 12:59:50Z")
	endTime, _ = time.Parse(layout, "2018-04-01 12:59:55Z")
	rightNow, _ = time.Parse(layout, "2018-04-01 13:00:00Z")

	intendedStart, _ := time.Parse(layout, "2018-04-01 12:59:30Z")
	intendedEnd, _ := time.Parse(layout, "2018-04-01 12:59:45Z")
	actualStart, actualEnd, _ := AdjustTimestampsForHystersisAndMinimumWindow(startTime, endTime, rightNow, hysteresis, window)
	suite.Equal(intendedStart, actualStart)
	suite.Equal(intendedEnd, actualEnd)
}

func (suite *TSDB_Helpers_TS) TestValidateTimeBoundNodeRequest_HappyPath() {

	var startTime, endTime time.Time

	layout := "2006-01-02 15:04:05Z"

	startTime, _ = time.Parse(layout, "2018-04-01 12:59:50Z")
	endTime, _ = time.Parse(layout, "2018-04-01 12:59:55Z")

	var tbnr TimeBoundNodeRequest

	tbnr.StartTime = startTime
	tbnr.EndTime = endTime
	tbnr.Nodes = append(tbnr.Nodes, NodeLookup{"x0c0s0b0n0", 0})
	isValid, err := ValidateTimeBoundNodeRequest(tbnr)
	suite.Equal(true, isValid)
	suite.True(err == nil)
}

func (suite *TSDB_Helpers_TS) TestValidateTimeBoundNodeRequest_MissingNodes() {

	var startTime, endTime time.Time

	layout := "2006-01-02 15:04:05Z"

	startTime, _ = time.Parse(layout, "2018-04-01 12:59:50Z")
	endTime, _ = time.Parse(layout, "2018-04-01 12:59:55Z")

	var tbnr TimeBoundNodeRequest

	tbnr.StartTime = startTime
	tbnr.EndTime = endTime

	isValid, err := ValidateTimeBoundNodeRequest(tbnr)
	suite.Equal(false, isValid)
	suite.True(err != nil)
}

func (suite *TSDB_Helpers_TS) TestValidateTimeBoundNodeRequest_MissingStartTime() {

	var endTime time.Time

	layout := "2006-01-02 15:04:05Z"

	endTime, _ = time.Parse(layout, "2018-04-01 12:59:55Z")

	var tbnr TimeBoundNodeRequest

	tbnr.EndTime = endTime
	tbnr.Nodes = append(tbnr.Nodes, NodeLookup{"x0c0s0b0n0", 0})
	isValid, err := ValidateTimeBoundNodeRequest(tbnr)
	suite.Equal(false, isValid)
	suite.True(err != nil)
}

func (suite *TSDB_Helpers_TS) TestValidateTimeBoundNodeRequest_MissingEndTime() {

	var startTime time.Time
	layout := "2006-01-02 15:04:05Z"
	startTime, _ = time.Parse(layout, "2018-04-01 12:59:55Z")

	var tbnr TimeBoundNodeRequest
	tbnr.StartTime = startTime
	tbnr.Nodes = append(tbnr.Nodes, NodeLookup{"x0c0s0b0n0", 0})

	isValid, err := ValidateTimeBoundNodeRequest(tbnr)
	suite.Equal(false, isValid)
	suite.True(err != nil)
}

func (suite *TSDB_Helpers_TS) TestValidateTimeBoundNodeRequest_MissingTime() {
	var tbnr TimeBoundNodeRequest

	tbnr.Nodes = append(tbnr.Nodes, NodeLookup{"x0c0s0b0n0", 0})
	isValid, err := ValidateTimeBoundNodeRequest(tbnr)
	suite.Equal(false, isValid)
	suite.True(err != nil)
}


func TestTSDBTimeHelpers(t *testing.T) {
	logger.SetupLogging()

	suite.Run(t, new(TSDB_Helpers_TS))
}
