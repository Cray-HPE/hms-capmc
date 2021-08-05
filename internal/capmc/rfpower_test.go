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

package capmc

import (
	"encoding/json"
	"reflect"
	"testing"
)

const powerData = `{"PowerControl":[{"Name":"Node Power Control","Oem":{"Cray":{"PowerAllocatedWatts":900,"PowerIdleWatts":250,"PowerLimit":{"Min":350,"Max":850,"Factor":1.02},"PowerResetWatts":250}},"PowerCapacityWatts":900,"PowerLimit":{"CorrectionInMs":6000,"LimitException":"LogEventOnly","LimitInWatts":500}},{"Name":"Accelerator0 Power Control","Oem":{"Cray":{"PowerIdleWatts":100,"PowerLimit":{"Min":200,"Max":350,"Factor":1}}},"PowerLimit":{"CorrectionInMs":6000,"LimitException":"LogEventOnly","LimitInWatts":300}}]}`

var (
	oneHundred        = 100
	twoHundred        = 200
	twoHundredFifty   = 250
	threeHundred      = 300
	threeHundredFifty = 350
	fiveHundred       = 500
	eightHundredFifty = 850
	nineHundred       = 900
	sixThousand       = 6000
	zero              = 0
	thirtyTwoFifty    = 3250
	thirtyFiveHundred = 3500
	thirtySevenFifty  = 3750

	onePtZero    = float32(1.0)
	onePtZeroTwo = float32(1.02)

	rfPCNode = PowerControl{
		Name:               "Node Power Control",
		PowerCapacityWatts: &nineHundred,
		OEM: &PowerControlOEM{
			Cray: &PowerControlOEMCray{
				PowerAllocatedWatts: &nineHundred,
				PowerIdleWatts:      &twoHundredFifty,
				PowerResetWatts:     &twoHundredFifty,
				PowerLimit: &PowerLimitOEMCray{
					Min:    &threeHundredFifty,
					Max:    &eightHundredFifty,
					Factor: &onePtZeroTwo,
				},
			},
		},
		PowerLimit: &PowerLimit{
			CorrectionInMs: &sixThousand,
			LimitException: "LogEventOnly",
			LimitInWatts:   &fiveHundred,
		},
	}
	rfPCAccel = PowerControl{
		Name: "Accelerator0 Power Control",
		OEM: &PowerControlOEM{
			Cray: &PowerControlOEMCray{
				PowerIdleWatts: &oneHundred,
				PowerLimit: &PowerLimitOEMCray{
					Min:    &twoHundred,
					Max:    &threeHundredFifty,
					Factor: &onePtZero,
				},
			},
		},
		PowerLimit: &PowerLimit{
			CorrectionInMs: &sixThousand,
			LimitException: "LogEventOnly",
			LimitInWatts:   &threeHundred,
		},
	}
	rfErrorComp = HpeError{
		Code:    "iLO.0.10.ExtendedInfo",
		Message: "See @Message.ExtendedInfo for more information.",
		ExtendedInfo: []HpeMessageExtendedInfo{
			{
				MessageId: "iLO.2.14.LicenseKeyRequired",
			},
		},
	}
	rfHpeActualPowerComp = HpeActualPowerLimits{
		PowerLimitInWatts: &thirtyFiveHundred,
		ZoneNumber:        &zero,
	}
	rfHpePowerLimitRangeComp = HpePowerLimitRanges{
		MaximumPowerLimit: &thirtySevenFifty,
		MinimumPowerLimit: &thirtyTwoFifty,
		ZoneNumber:        &zero,
	}
	rfHpePowerLimitComp = HpePowerLimits{
		PowerLimitInWatts: &thirtyFiveHundred,
		ZoneNumber:        &zero,
	}
	rfHpeConfigurePowerLimitComp = HpeConfigurePowerLimit{
		PowerLimits: []HpePowerLimits{
			{
				PowerLimitInWatts: &thirtyFiveHundred,
				ZoneNumber:        &zero,
			},
		},
	}
)

const HpeErrorData = `{"error":{"code":"iLO.0.10.ExtendedInfo","message":"See @Message.ExtendedInfo for more information.","@Message.ExtendedInfo":[{"MessageId":"iLO.2.14.LicenseKeyRequired"}]}}`
const HpePowerLimitData = `{"ActualPowerLimits":[{"PowerLimitInWatts":3500,"ZoneNumber":0}],"PowerLimitRanges":[{"MaximumPowerLimit":3750,"MinimumPowerLimit":3250,"ZoneNumber":0}],"PowerLimits":[{"PowerLimitInWatts":3500,"ZoneNumber":0}]}`
const HpeConfigurePowerLimitData = `{"PowerLimits":[{"PowerLimitInWatts":3500,"ZoneNumber":0}]}`

func TestMarshalRFStruct(t *testing.T) {
	t.Run("Cray Power Control", func(t *testing.T) {
		var rfPower Power
		rfPower.PowerCtl = append(rfPower.PowerCtl, rfPCNode)
		rfPower.PowerCtl = append(rfPower.PowerCtl, rfPCAccel)

		body, err := json.Marshal(rfPower)
		if err != nil {
			t.Fatal("Failed to marshal Power structure")
		}

		if string(body) != powerData {
			t.Errorf("Marshalled structure does not match:\ngot  %s\nwant %s",
				string(body), powerData)
		}
	})

	t.Run("Hpe Power Control Error", func(t *testing.T) {
		var rfHpePowerError Power
		rfHpePowerError.Error = &rfErrorComp

		body, err := json.Marshal(rfHpePowerError)
		if err != nil {
			t.Fatal("Failed to marshal Hpe Power Error structure")
		}

		if string(body) != HpeErrorData {
			t.Errorf("Marshalled structure does not match:\ngot  %s\nwant %s",
				string(body), HpeErrorData)
		}
	})

	t.Run("Hpe Power Control", func(t *testing.T) {
		var rfHpePower Power
		rfHpePower.ActualPowerLimits = append(rfHpePower.ActualPowerLimits, rfHpeActualPowerComp)
		rfHpePower.PowerLimitRanges = append(rfHpePower.PowerLimitRanges, rfHpePowerLimitRangeComp)
		rfHpePower.PowerLimits = append(rfHpePower.PowerLimits, rfHpePowerLimitComp)

		body, err := json.Marshal(rfHpePower)
		if err != nil {
			t.Fatal("Failed to marshal Hpe Power structure")
		}

		if string(body) != HpePowerLimitData {
			t.Errorf("Marshalled structure does not match:\ngot  %s\nwant %s",
				string(body), HpePowerLimitData)
		}
	})
	t.Run("Hpe Configure Power Control", func(t *testing.T) {
		var rfHpeConfigPowerLimit HpeConfigurePowerLimit
		rfHpeConfigPowerLimit = rfHpeConfigurePowerLimitComp

		body, err := json.Marshal(rfHpeConfigPowerLimit)
		if err != nil {
			t.Fatal("Failed to marshal Hpe Configure Power structure")
		}

		if string(body) != HpeConfigurePowerLimitData {
			t.Errorf("Marshalled structure does not match:\ngot  %s\nwant %s",
				string(body), HpeConfigurePowerLimitData)
		}
	})
}

func TestUnMarshalRFStruct(t *testing.T) {
	t.Run("Cray Power Control", func(t *testing.T) {
		var rfData Power
		var rfPower Power
		err := json.Unmarshal([]byte(powerData), &rfData)
		if err != nil {
			t.Fatal("Failed to unmarshal PowerControl json", err.Error())
		}

		rfPower.PowerCtl = append(rfPower.PowerCtl, rfPCNode)
		rfPower.PowerCtl = append(rfPower.PowerCtl, rfPCAccel)

		for idx, comp := range rfPower.PowerCtl {
			if !reflect.DeepEqual(rfData.PowerCtl[idx], comp) {
				t.Errorf("Unmarshalled PowerControl does not match:\ngot  %v\nwant %v",
					rfData, rfPower)
			}
		}
	})

	t.Run("Hpe Power Control Error", func(t *testing.T) {
		var rfHpeErrorData Power
		err := json.Unmarshal([]byte(HpeErrorData), &rfHpeErrorData)
		if err != nil {
			t.Fatal("Failed to unmarshal HpeError json", err.Error())
		}

		if rfHpeErrorData.Error == nil {
			t.Errorf("Unmarshalled Error does not contain a valid error")
		} else if !reflect.DeepEqual(*rfHpeErrorData.Error, rfErrorComp) {
			t.Errorf("Unmarshalled Error does not match:\ngot  %v\nwant %v",
				rfHpeErrorData.Error, rfErrorComp)
		}
	})

	t.Run("Hpe Power Control", func(t *testing.T) {
		var rfHpePowerLimitData Power
		err := json.Unmarshal([]byte(HpePowerLimitData), &rfHpePowerLimitData)
		if err != nil {
			t.Fatal("Failed to unmarshal Hpe PowerLimit json", err.Error())
		}

		if !reflect.DeepEqual(rfHpePowerLimitData.ActualPowerLimits[0], rfHpeActualPowerComp) {
			t.Errorf("Unmarshalled Hpe ActualPowerLimits does not match:\ngot  %v\nwant %v",
				rfHpePowerLimitData.ActualPowerLimits[0], rfHpeActualPowerComp)
		}
		if !reflect.DeepEqual(rfHpePowerLimitData.PowerLimitRanges[0], rfHpePowerLimitRangeComp) {
			t.Errorf("Unmarshalled Hpe PowerLimitRanges does not match:\ngot  %v\nwant %v",
				rfHpePowerLimitData.PowerLimitRanges[0], rfHpePowerLimitRangeComp)
		}
		if !reflect.DeepEqual(rfHpePowerLimitData.PowerLimits[0], rfHpePowerLimitComp) {
			t.Errorf("Unmarshalled Hpe PowerLimits does not match:\ngot  %v\nwant %v",
				rfHpePowerLimitData.PowerLimits[0], rfHpePowerLimitComp)
		}
	})

	t.Run("Hpe Configure Power Control", func(t *testing.T) {
		var rfConfigPowerLimit HpeConfigurePowerLimit
		err := json.Unmarshal([]byte(HpeConfigurePowerLimitData), &rfConfigPowerLimit)
		if err != nil {
			t.Fatal("Failed to unmarshal HpeConfigurePowerLimit json", err.Error())
		}

		if !reflect.DeepEqual(rfConfigPowerLimit, rfHpeConfigurePowerLimitComp) {
			t.Errorf("Unmarshalled Configure Power Limit does not match:\ngot  %v\nwant %v",
				rfConfigPowerLimit, rfHpeConfigurePowerLimitComp)
		}
	})

}
