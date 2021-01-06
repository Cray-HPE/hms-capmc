// Copyright 2019,2020 Hewlett Packard Enterprise Development LP

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

const testBmcGetPowerCapCall = `
{
	"odata.context": "/redfish/v1/$metadata#Power.Power(*)",
	"@odata.etag": "W/\"946641656\"",
	"@odata.id": "/redfish/v1/Chassis/Node0/Power",
	"@odata.type": "#Power.v1_4_0.Power",
	"Description": "Power sensor readings",
	"Id": "Power",
	"Name": "Power",
	"PowerControl": [
		{
			"Name":"Node Power Control",
			"Oem": {
				"Cray": {
					"PowerAllocatedWatts":900,
					"PowerIdleWatts":250,
					"PowerLimit": {
						"Min":350,
						"Max":850,
						"Factor":1.02
					},
					"PowerResetWatts":250
				}
			},
			"PowerCapacityWatts":900,
			"PowerLimit": {
				"CorrectionInMs":6000,
				"LimitException":"LogEventOnly",
				"LimitInWatts":500
			}
		},
		{
			"Name":"Accelerator0 Power Control",
			"Oem": {
				"Cray": {
					"PowerIdleWatts":100,
					"PowerLimit": {
						"Min":200,
						"Max":350,
						"Factor":1
					}
				}
			},
			"PowerLimit": {
				"CorrectionInMs":6000,
				"LimitException":"LogEventOnly",
				"LimitInWatts":300
			}
		}
	],
	"PowerControl@odata.count": 1
}
`
const testGBBmcGetPowerCapCall = `
{
  "@odata.context": "/redfish/v1/$metadata#Power.Power",
  "@odata.etag": "W/\"947088808\"",
  "@odata.id": "/redfish/v1/Chassis/Self/Power",
  "@odata.type": "#Power.v1_5_3.Power",
  "Description": "Power sensor readings",
  "Id": "Power",
  "Name": "Power",
  "PowerControl": [
    {
      "@odata.id": "/redfish/v1/Chassis/Self/Power#/PowerControl/0",
      "MemberId": "0",
      "Name": "Chassis Power Control",
      "Oem": {
        "Vendor": {
          "@odata.type": "#GbtPowerLimit.v1_0_0.Vendor",
          "PowerIdleWatts": 61,
          "PowerLimit": {
            "Factor": 10,
            "Max": 900,
            "Min": 61
          },
          "PowerMetrics": {
            "AccumulatedEnergyJoules": 719,
            "Timestamp": "2000-01-02T01:06:29"
          },
          "PowerResetWatts": 0
        }
      },
      "PhysicalContext": "Intake",
      "PowerCapacityWatts": 900,
      "PowerConsumedWatts": 95,
      "PowerLimit": {
        "CorrectionInMs": 1000,
        "LimitException": "HardPowerOff",
        "LimitInWatts": 500
      },
      "PowerMetrics": {
        "AverageConsumedWatts": 40,
        "IntervalInMin": 0,
        "MaxConsumedWatts": 602,
        "MinConsumedWatts": 4
      },
      "Status": {
        "Health": "Critical",
        "State": "Disabled"
      }
    }
  ],
  "PowerControl@odata.count": 1
}
`

const testRedfishErrorNotFound = `
{
   "error" : {
      "code" : "Base.1.0.ResourceMissingAtURI",
      "@Message.ExtendedInfo" : [
         {
            "MessageArgs" : [
               "/redfish/v1/Chassis/Node0/cower"
            ],
            "Resolution" : "Place a valid resource at the URI or correct the URI and resubmit the request.",
            "MessageId" : "Base.1.0.ResourceMissingAtURI",
            "Severity" : "Critical",
            "@odata.type" : "#Message.v1_0_5.Message",
            "Message" : "The resource at the URI /redfish/v1/Chassis/Node0/cower was not found."
         }
      ],
      "message" : "The resource at the URI /redfish/v1/Chassis/Node0/cower was not found."
   }
}
`

func BmcTestFunc() RoundTripFunc {
	return func(req *http.Request) (*http.Response, error) {
		// fmt.Printf("BMCTESTFUNC(): %s %s\n", req.Method, req.URL.String())
		// fmt.Printf("%+v\n", req)
		header := make(http.Header)
		header.Set("Content-Type", "application/json")
		switch req.Method {
		case "GET":
			switch req.URL.String() {
			case "https://x5000c0s0b0/redfish/v1/Chassis/Node0/Power":
				return &http.Response{
					Status: fmt.Sprintf("%d %s",
						http.StatusOK,
						http.StatusText(http.StatusOK)),
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(testBmcGetPowerCapCall)),
					Header:     header,
					Request:    req,
				}, nil
			case "https://x3000c0s0b0/redfish/v1/Chassis/Node0/Power":
				return &http.Response{
					Status: fmt.Sprintf("%d %s",
						http.StatusOK,
						http.StatusText(http.StatusOK)),
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(testGBBmcGetPowerCapCall)),
					Header:     header,
					Request:    req,
				}, nil
			case "https://x5000c0s0b0/redfish/v1/Chassis/Node0/cower":
				return &http.Response{
					Status: fmt.Sprintf("%d %s",
						http.StatusNotFound,
						http.StatusText(http.StatusNotFound)),
					StatusCode: http.StatusNotFound,
					Body:       ioutil.NopCloser(bytes.NewBufferString(testRedfishErrorNotFound)),
					Header:     header,
					Request:    req,
				}, nil
			}
		case "PATCH":
			switch req.URL.String() {
			case "https://x5000c0s0b0/redfish/v1/Chassis/Node0/Power",
				"https://x3000c0s0b0/redfish/v1/Chassis/Node0/Power":
				return &http.Response{
					Status: fmt.Sprintf("%d %s",
						http.StatusCreated,
						http.StatusText(http.StatusCreated)),
					StatusCode: http.StatusCreated,
					Body:       ioutil.NopCloser(bytes.NewBufferString("")),
					Header:     header,
					Request:    req,
				}, nil
			case "https://x5000c0s0b0/redfish/v1/Chassis/Node0/cower":
				return &http.Response{
					Status: fmt.Sprintf("%d %s",
						http.StatusNotFound,
						http.StatusText(http.StatusNotFound)),
					StatusCode: http.StatusNotFound,
					Body:       ioutil.NopCloser(bytes.NewBufferString(testRedfishErrorNotFound)),
					Header:     header,
					Request:    req,
				}, nil
			}
		}

		return &http.Response{
			Status: fmt.Sprintf("%d %s",
				http.StatusBadRequest,
				http.StatusText(http.StatusBadRequest)),
			StatusCode: http.StatusBadRequest,
			Body:       ioutil.NopCloser(bytes.NewBufferString(`{"status":400}`)),
			Header:     header,
			Request:    req,
		}, nil
	}
}

const testBmcMissingURI = " x5000c0s0b0 HTTP 404 Not Found, Redfish Error " +
	"Message: The resource at the URI /redfish/v1/Chassis/Node0/cower " +
	"was not found. ExtendedInfo: Message: The resource at the URI " +
	"/redfish/v1/Chassis/Node0/cower was not found. Resolution: Place " +
	"a valid resource at the URI or correct the URI and resubmit the " +
	"request."

func TestDoBmcGetCall(t *testing.T) {
	var tSvc CapmcD
	var err error
	testClient := NewTestClient(BmcTestFunc())
	tSvc.smClient = testClient
	tSvc.rfClient = testClient
	if err != nil {
		t.Fatal(err)
	}

	comp1 := &NodeInfo{
		Hostname:   "x5000c0s0b0n0",
		BmcFQDN:    "x5000c0s0b0",
		RfPowerURL: "/redfish/v1/Chassis/Node0/cower",
	}
	comp2 := &NodeInfo{
		Hostname:   "x5000c0s0b0n0",
		BmcFQDN:    "x5000c0s0b0",
		RfPowerURL: "/redfish/v1/Chassis/Node0/Power",
	}

	tests := []struct {
		name    string
		ni      *NodeInfo
		command string
		want    bmcPowerRc
	}{
		{
			"bad command",
			comp2,
			bmcCmdPowerOff,
			bmcPowerRc{comp2, -1, "Invalid command Off", "Unknown"},
		},
		{
			"bad URL",
			comp1,
			bmcCmdGetPowerCap,
			bmcPowerRc{comp1, http.StatusNotFound, testBmcMissingURI, "Unknown"},
		},
		{
			"good request",
			comp2,
			bmcCmdGetPowerCap,
			bmcPowerRc{comp2, 0, testBmcGetPowerCapCall, "Unknown"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tSvc.doBmcGetCall(bmcCall{ni: tt.ni, bmcCmd: bmcCmd{cmd: tt.command}}); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CapmcD.doBmcGetCall() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDoBmcPatchCall(t *testing.T) {
	var tSvc CapmcD
	var err error
	testClient := NewTestClient(BmcTestFunc())
	tSvc.smClient = testClient
	tSvc.rfClient = testClient
	if err != nil {
		t.Fatal(err)
	}

	comp1 := &NodeInfo{
		Hostname:   "x5000c0s0b0n0",
		BmcFQDN:    "x5000c0s0b0",
		RfPowerURL: "/redfish/v1/Chassis/Node0/cower",
	}
	comp2 := &NodeInfo{
		Hostname:   "x5000c0s0b0n0",
		BmcFQDN:    "x5000c0s0b0",
		RfPowerURL: "/redfish/v1/Chassis/Node0/Power",
	}
	comp3 := &NodeInfo{
		Hostname:   "x3000c0s0b0n0",
		BmcFQDN:    "x3000c0s0b0",
		RfPowerURL: "/redfish/v1/Chassis/Node0/Power",
	}

	tests := []struct {
		name    string
		ni      *NodeInfo
		command string
		want    bmcPowerRc
	}{
		{
			"bad command",
			comp2,
			bmcCmdPowerOff,
			bmcPowerRc{comp2, -1, "Invalid command Off", "Unknown"},
		},
		{
			"bad URL",
			comp1,
			bmcCmdSetPowerCap,
			bmcPowerRc{comp1, http.StatusNotFound, testBmcMissingURI, "Unknown"},
		},
		{
			"good MTN request",
			comp2,
			bmcCmdSetPowerCap,
			bmcPowerRc{comp2, 0, "", "Unknown"},
		},
		{
			"good RVR request",
			comp3,
			bmcCmdSetPowerCap,
			bmcPowerRc{comp3, 0, "", "Unknown"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tSvc.doBmcPatchCall(bmcCall{ni: tt.ni, bmcCmd: bmcCmd{cmd: tt.command}}); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CapmcD.doBmcPatchCall() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
