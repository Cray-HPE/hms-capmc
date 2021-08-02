/*
 * MIT License
 *
 * (C) Copyright [2019-2021] Hewlett Packard Enterprise Development LP
 *
 * Permission is hereby granted, free of charge, to any person obtaining a
 * copy of this software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation
 * the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the
 * Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included
 * in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
 * THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
 * OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
 * ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
 * OTHER DEALINGS IN THE SOFTWARE.
 */

package tsdb

import (
	"github.com/Cray-HPE/hms-capmc/internal/logger"
	"github.com/stretchr/testify/suite"
	"reflect"
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
	tbnr.EndTime = endTime
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
	tbnr.EndTime = endTime
	tbnr.StartTime = startTime

	x, err := TSDBContext.GetNodeEnergyStats(tbnr)
	suite.Equal(nil, err)
	suite.True(reflect.TypeOf(x) == reflect.TypeOf(&NodeEnergyStats{}))

}

func (suite *TSDB_TS) TestGetNodeEnergyCounter_nilError() {
	startTime, _ := time.Parse("2001-04-01 00:00:00", "2001-04-01 00:00:00")
	endTime, _ := time.Parse("2001-04-01 00:00:00", "2001-04-01 12:00:00")
	tbnr := TimeBoundNodeRequest{}
	tbnr.EndTime = endTime
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
