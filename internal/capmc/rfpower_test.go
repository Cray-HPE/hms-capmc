// MIT License
//
// (C) Copyright [2019, 2021] Hewlett Packard Enterprise Development LP
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
