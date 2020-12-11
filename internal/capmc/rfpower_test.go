// Copyright 2019 Cray Inc. All Rights Reserved.
//
// Except as permitted by contract or express written permission of Cray Inc.,
// no part of this work or its content may be modified, used, reproduced or
// disclosed in any form. Modifications made without express permission of
// Cray Inc. may damage the system the software is installed within, may
// disqualify the user from receiving support from Cray Inc. under support or
// maintenance contracts, or require additional support services outside the
// scope of those contracts to repair the software or system.
//

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
)

func TestMarshalRFStruct(t *testing.T) {
	var rfPower Power

	rfPower.PowerCtl = append(rfPower.PowerCtl, rfPCNode)
	rfPower.PowerCtl = append(rfPower.PowerCtl, rfPCAccel)

	body, err := json.Marshal(rfPower)
	if err != nil {
		t.Fatal("Nope")
	}

	if string(body) != powerData {
		t.Errorf("Marshalled structure does not match:\ngot  %s\nwant %s",
			string(body), powerData)
	}
}

func TestUnMarshalRFStruct(t *testing.T) {
	var (
		rfData  Power
		rfPower Power
	)

	err := json.Unmarshal([]byte(powerData), &rfData)
	if err != nil {
		t.Fatal("Nope", err.Error())
	}

	rfPower.PowerCtl = append(rfPower.PowerCtl, rfPCNode)
	rfPower.PowerCtl = append(rfPower.PowerCtl, rfPCAccel)

	for idx, comp := range rfPower.PowerCtl {
		if !reflect.DeepEqual(rfData.PowerCtl[idx], comp) {
			t.Errorf("Unmarshalled structure does not match:\ngot  %v\nwant %v",
				rfData, rfPower)
		}
	}
}
