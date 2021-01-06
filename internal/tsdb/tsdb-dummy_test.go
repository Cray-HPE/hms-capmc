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
	"github.com/stretchr/testify/suite"
	"reflect"
	"stash.us.cray.com/HMS/hms-capmc/internal/logger"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
)

type TSDB_TS struct {
	suite.Suite
}

// SetupSuit is run ONCE
func (suite *TSDB_TS) SetupSuite() {
}

func (suite *TSDB_TS) TestGetSystemPower_nilError() {
	tbr := TimeBoundRequest{}
	tbr.StartTime, _ = time.Parse("2001-04-01 00:00:00", "2001-04-01 00:00:00")
	tbr.EndTime, _ = time.Parse("2001-04-01 00:00:00", "2001-04-01 12:00:00")
	sysPow, err := TSDBContext.GetSystemPower(tbr)

	suite.Equal(1, *sysPow.Min)
	suite.Equal(100, *sysPow.Max)
	suite.Equal(50, *sysPow.Avg)
	suite.Equal(nil, err)
	suite.True(reflect.TypeOf(sysPow) == reflect.TypeOf(&SystemPower{}))

}

func (suite *TSDB_TS) TestGetSystemPowerDetails_nilError() {
	tbr := TimeBoundRequest{}
	tbr.StartTime, _ = time.Parse("2001-04-01 00:00:00", "2001-04-01 00:00:00")
	tbr.EndTime, _ = time.Parse("2001-04-01 00:00:00", "2001-04-01 12:00:00")
	x, err := TSDBContext.GetSystemPowerDetails(tbr)
	suite.Equal(nil, err)
	suite.True(reflect.TypeOf(x) == reflect.TypeOf(&SystemPowerByCabinet{}))

}

func (suite *TSDB_TS) TestGetNodeEnergy_nilError() {

	startTime, _ := time.Parse("2001-04-01 00:00:00", "2001-04-01 00:00:00")
	endTime, _ := time.Parse("2001-04-01 00:00:00", "2001-04-01 12:00:00")
	tbnr := TimeBoundNodeRequest{}
	tbnr.EndTime= endTime
	tbnr.StartTime = startTime

	x, err := TSDBContext.GetNodeEnergy(tbnr)
	suite.Equal(nil, err)
	log.Info(reflect.TypeOf(x))
	suite.True(reflect.TypeOf(x) == reflect.TypeOf(&NodeEnergy{}))

}

func (suite *TSDB_TS) TestGetNodeEnergyStats_nilError() {
	startTime, _ := time.Parse("2001-04-01 00:00:00", "2001-04-01 00:00:00")
	endTime, _ := time.Parse("2001-04-01 00:00:00", "2001-04-01 12:00:00")
	tbnr := TimeBoundNodeRequest{}
	tbnr.EndTime= endTime
	tbnr.StartTime = startTime

	x, err := TSDBContext.GetNodeEnergyStats(tbnr)
	suite.Equal(nil, err)
	suite.True(reflect.TypeOf(x) == reflect.TypeOf(&NodeEnergyStats{}))

}

func (suite *TSDB_TS) TestGetNodeEnergyCounter_nilError() {
	startTime, _ := time.Parse("2001-04-01 00:00:00", "2001-04-01 00:00:00")
	endTime, _ := time.Parse("2001-04-01 00:00:00", "2001-04-01 12:00:00")
	tbnr := TimeBoundNodeRequest{}
	tbnr.EndTime= endTime
	tbnr.StartTime = startTime

	x, err := TSDBContext.GetNodeEnergyCounter(tbnr)
	suite.Equal(nil, err)
	suite.True(reflect.TypeOf(x) == reflect.TypeOf(&NodeEnergyCounters{}))
}


func TestTSDBDummySuite(t *testing.T) {
	logger.SetupLogging()
	ConfigureDataImplementation(DUMMY)

	suite.Run(t, new(TSDB_TS))
}
