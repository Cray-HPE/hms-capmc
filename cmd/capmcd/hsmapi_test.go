/*
 * MIT License
 *
 * (C) Copyright [2019-2022,2025] Hewlett Packard Enterprise Development LP
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

// This file contains the unit tests for CAPMC to hardware state manager (HSM)
// interfaces
//

package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/Cray-HPE/hms-certs/pkg/hms_certs"
	compcreds "github.com/Cray-HPE/hms-compcredentials"
	sstorage "github.com/Cray-HPE/hms-securestorage"
	rf "github.com/Cray-HPE/hms-smd/v2/pkg/redfish"
	"github.com/Cray-HPE/hms-smd/v2/pkg/sm"
)

// testEq tests for equality between two arrays/slices
func testEq(a, b []int) bool {
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

// TestValidateNIDs unit test function for validateNIDs
func TestValidateNIDs(t *testing.T) {
	var tests = []struct {
		input []int
		good  []int
		bad   []int
		dups  bool
	}{
		{
			[]int{0, 1, 2, 3, 4},
			[]int{0, 1, 2, 3, 4},
			[]int{},
			false,
		},
		{
			[]int{1, 2, 3, 4, 6, 8},
			[]int{1, 2, 3, 4, 6, 8},
			[]int{},
			false,
		},
		{
			[]int{-1, 0, 1},
			[]int{0, 1},
			[]int{-1},
			false,
		},
		{
			[]int{1, 2, 3, 4, 5, 4, 3, 2, 1},
			[]int{1, 2, 3, 4, 5},
			[]int{4, 3, 2, 1},
			false,
		},
		{
			[]int{1, 2, 3, 4, 5, 4, 3, 2, 1},
			[]int{1, 2, 3, 4, 5},
			[]int{},
			true,
		},
		{
			[]int{10, 12, 12, 14, 16, 18, 20, 20},
			[]int{10, 12, 14, 16, 18, 20},
			[]int{12, 20},
			false,
		},
		{
			[]int{10, 12, 12, 14, 16, 18, 20, 20},
			[]int{10, 12, 14, 16, 18, 20},
			[]int{},
			true,
		},
		{
			[]int{-1, 2, 2, 4, 20, 30, -8, 10, 18, 2},
			[]int{2, 4, 20, 30, 10, 18},
			[]int{-1, 2, -8, 2},
			false,
		},
		{
			[]int{-1, 2, 2, 4, 20, 30, -8, 10, 18, 2},
			[]int{2, 4, 20, 30, 10, 18},
			[]int{-1, -8},
			true,
		},
	}

	for n, test := range tests {
		good, bad := validateNIDs(test.dups, test.input)
		if !testEq(good, test.good) {
			t.Errorf("Test #%d validateNIDs(%v, %v) = %v, %v",
				n, test.dups, test.input, good, bad)
		}
		if !testEq(bad, test.bad) {
			t.Errorf("Test #%d validateNIDs(%v, %v) = %v, %v",
				n, test.dups, test.input, good, bad)
		}
	}
}

func TestGetRestrictStr(t *testing.T) {
	tests := []struct {
		in  HSMQuery
		out string
	}{{
		// Test forming an HSM query string with all fields
		in: HSMQuery{
			ComponentIDs: []string{"x0c0s0b0", "x0c0s0b0n0"},
			Types:        []string{"Node", "NodeBMC"},
			States:       []string{"!Empty"},
			NIDs:         []int{1, 2, 3},
			Enabled:      []bool{true},
		},
		out: "enabled=true&id=x0c0s0b0&id=x0c0s0b0n0&nid=1&nid=2&nid=3&state=%21Empty&type=Node&type=NodeBMC",
	}, {
		// Test forming an HSM query string with only one item
		in: HSMQuery{
			ComponentIDs: []string{"x0c0s0b0"},
		},
		out: "id=x0c0s0b0",
	}, {
		// Test forming an HSM query string with an empty query
		in:  HSMQuery{},
		out: "",
	}}
	for n, test := range tests {
		result := getRestrictStr(test.in)
		if test.out != result {
			t.Errorf("TestGetRestrictStr Test Case %d: FAIL: Expected %v but got %v", n, test.out, result)
		}
	}
}

func NewTestClient(f RoundTripFunc) *hms_certs.HTTPClientPair {
	hms_certs.ConfigParams.LogInsecureFailover = false
	rc, _ := makeClient(0, 5)
	rc.InsecureClient.HTTPClient = &http.Client{Transport: RoundTripFunc(f)}
	rc.SecureClient = rc.InsecureClient
	return rc
}

func TestDoRequest(t *testing.T) {
	var tests = []struct {
		sc   int
		body io.ReadCloser
		hct  string
		e    bool
		eMsg string
	}{
		{
			sc:   http.StatusOK,
			body: ioutil.NopCloser(bytes.NewBufferString(`OK`)),
			hct:  "application/json",
			e:    false,
			eMsg: "",
		},
		{
			sc:   http.StatusBadRequest,
			body: ioutil.NopCloser(bytes.NewBuffer([]byte(`{"type":"about:blank","title":"Bad Request","detail": "well, duh","status":400}`))),
			hct:  "application/problem+json",
			e:    true,
			eMsg: "Hardware State Manager error: Error: 400 Bad Request: Details: well, duh",
		},
		{
			sc:   http.StatusBadRequest,
			body: ioutil.NopCloser(bytes.NewBuffer([]byte(`{"type":"about:blank","title":"Bad Request","detail": "well, duh","status":400}`))),
			hct:  "",
			e:    true,
			eMsg: "{\"type\":\"about:blank\",\"title\":\"Bad Request\",\"detail\": \"well, duh\",\"status\":400}",
		},
	}

	for _, test := range tests {
		client := NewTestClient(func(req *http.Request) (*http.Response, error) {
			// Test request paramters here
			locHeader := make(http.Header)
			if test.hct != "" {
				locHeader.Add("Content-Type", test.hct)
			}
			return &http.Response{
				StatusCode: test.sc,
				Body:       test.body,
				Header:     locHeader,
			}, nil
		})

		svc := CapmcD{rfClient: client, smClient: client, config: loadConfig("")}
		req, err := http.NewRequest(http.MethodGet, "http://example.com", nil)
		if err != nil {
			t.Fatal(err)
		}

		body, err := svc.doRequest(req)

		// check err
		if err != nil {
			// see if we are expecting an error
			if !test.e {
				t.Fatal(err)
			} else if err.Error() != test.eMsg {
				// expected error, but message doesn't match
				t.Fatalf("Wrong error message, expected %s, got %s", test.eMsg, err.Error())
			}
		}

		// check body
		if body == nil {
			t.Errorf("incomplete")
		}
	}
}

const testNodeComponentEndpoint = `
{
   "ComponentEndpoints" : [
      {
         "RedfishURL" : "10.100.107.181/redfish/v1/Systems/Node1",
         "RedfishEndpointFQDN" : "10.100.107.181",
         "ID" : "x1000c7s1b1n1",
         "MACAddr" : "00:40:a6:82:f5:bd",
         "RedfishEndpointID" : "x1000c7s1b1",
         "RedfishSubtype" : "Physical",
         "ComponentEndpointType" : "ComponentEndpointComputerSystem",
         "RedfishSystemInfo" : {
            "EthernetNICInfo" : [
               {
                  "MACAddress" : "00:40:a6:82:f5:bd",
                  "Description" : "Node Maintenance Network",
                  "@odata.id" : "/redfish/v1/Systems/Node1/EthernetInterfaces/ManagementEthernet",
                  "RedfishId" : "ManagementEthernet",
                  "PermanentMACAddress" : "00:40:a6:82:f5:bd"
               }
            ],
            "Actions" : {
               "#ComputerSystem.Reset" : {
                  "ResetType@Redfish.AllowableValues" : [
						"ForceOff",
						"On",
						"Off"
				  ],
                  "target" : "/redfish/v1/Systems/Node1/Actions/ComputerSystem.Reset"
               }
            },
            "Name" : "Node1"
         },
         "Type" : "Node",
         "OdataID" : "/redfish/v1/Systems/Node1",
         "RedfishType" : "ComputerSystem"
      }
   ]
}
`

const testNodeComponentState = `
{
   "Components" : [
      {
         "State" : "Ready",
         "NID" : 1232,
         "NetType" : "Sling",
         "Flag" : "OK",
         "Arch" : "X86",
         "Enabled" : true,
         "ID" : "x1000c7s1b1n1",
         "Type" : "Node",
         "Role" : "Compute"
      }
   ]
}
`

const testPDUComponentEndpoint = `
{
    "ComponentEndpoints": [
        {
            "ID": "x0m0p0j19",
            "Type": "CabinetPDUOutlet",
            "RedfishType": "Outlet",
            "RedfishSubtype": "C13",
            "OdataID": "/redfish/v1/PowerEquipment/RackPDUs/A/Outlets/AA26",
            "RedfishEndpointID": "x0m0",
            "RedfishEndpointFQDN": "x0m0:8082",
            "RedfishURL": "x0m0:8082/redfish/v1/PowerEquipment/RackPDUs/A/Outlets/AA26",
            "ComponentEndpointType": "ComponentEndpointOutlet",
            "RedfishOutletInfo": {
                "Name": "Master_Outlet_26",
                "Actions": {
                    "#Outlet.PowerControl": {
                        "PowerState@Redfish.AllowableValues": [
                            "On",
                            "Off"
                        ],
                        "target": "/redfish/v1/PowerEquipment/RackPDUs/A/Outlets/AA26/Actions/Outlet.PowerControl"
                    }
                }
            }
        }
    ]
}
`

const testPDUComponentState = `
{
    "Components": [
        {
            "ID": "x0m0p0j19",
            "Type": "CabinetPDUOutlet",
            "State": "On",
            "Flag": "OK",
            "Enabled": true
        }
    ]
}
`

const testPDUComponentEndpoint2 = `
{
    "ComponentEndpoints": [
        {
            "ID": "x0m0p0v19",
            "Type": "CabinetPDUPowerConnector",
            "RedfishType": "Outlet",
            "RedfishSubtype": "C13",
            "OdataID": "/redfish/v1/PowerEquipment/RackPDUs/A/Outlets/AA26",
            "RedfishEndpointID": "x0m0",
            "RedfishEndpointFQDN": "x0m0:8082",
            "RedfishURL": "x0m0:8082/redfish/v1/PowerEquipment/RackPDUs/A/Outlets/AA26",
            "ComponentEndpointType": "ComponentEndpointOutlet",
            "RedfishOutletInfo": {
                "Name": "Master_Outlet_26",
                "Actions": {
                    "#Outlet.PowerControl": {
                        "PowerState@Redfish.AllowableValues": [
                            "On",
                            "Off"
                        ],
                        "target": "/redfish/v1/PowerEquipment/RackPDUs/A/Outlets/AA26/Actions/Outlet.PowerControl"
                    }
                }
            }
        }
    ]
}
`

const testPDUComponentState2 = `
{
    "Components": [
        {
            "ID": "x0m0p0v19",
            "Type": "CabinetPDUPowerConnector",
            "State": "On",
            "Flag": "OK",
            "Enabled": true
        }
    ]
}
`

func GetNodesTestFunc() RoundTripFunc {
	return func(req *http.Request) (*http.Response, error) {
		switch req.URL.String() {
		case "http://localhost:27779/State/Components?id=x1000c7s1b1n1":
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(testNodeComponentState)),
				Header:     make(http.Header),
			}, nil
		case "http://localhost:27779/Inventory/ComponentEndpoints?id=x1000c7s1b1n1":
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(testNodeComponentEndpoint)),
				Header:     make(http.Header),
			}, nil
		case "http://localhost:27779/State/Components?id=x0m0p0j19":
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(testPDUComponentState)),
				Header:     make(http.Header),
			}, nil
		case "http://localhost:27779/State/Components?id=x0m0p0v19":
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(testPDUComponentState2)),
				Header:     make(http.Header),
			}, nil
		case "http://localhost:27779/Inventory/ComponentEndpoints?id=x0m0p0j19":
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(testPDUComponentEndpoint)),
				Header:     make(http.Header),
			}, nil
		case "http://localhost:27779/Inventory/ComponentEndpoints?id=x0m0p0v19":
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(testPDUComponentEndpoint2)),
				Header:     make(http.Header),
			}, nil
		default:
			return &http.Response{
				StatusCode: 404,
				Body:       ioutil.NopCloser(bytes.NewBufferString("")),
				Header:     make(http.Header),
			}, nil
		}
	}
}

func TestGetNodes(t *testing.T) {
	var tSvc CapmcD
	var err error
	tSvc.hsmURL, err = url.Parse("http://localhost:27779")
	// Use the mock secure storage
	ss, adapter := sstorage.NewMockAdapter()
	ccs := compcreds.NewCompCredStore("secret/hms-cred", ss)
	tSvc.ss = ss
	tSvc.ccs = ccs
	testClient := NewTestClient(GetNodesTestFunc())
	tSvc.rfClient = testClient
	tSvc.smClient = testClient
	if err != nil {
		t.Fatal(err)
	}
	var ni []*NodeInfo
	tests := []struct {
		name      string
		in        HSMQuery
		out       []*NodeInfo
		ssData    []sstorage.MockLookup
		res       NodeInfo
		expectErr bool
	}{
		{
			"Good node",
			HSMQuery{ComponentIDs: []string{"x1000c7s1b1n1"}},
			ni,
			[]sstorage.MockLookup{
				{
					Output: sstorage.OutputLookup{
						Output: &compcreds.CompCredentials{
							Xname:    "x1000c7s1b1n1",
							URL:      "10.100.107.181/redfish/v1/Node1",
							Username: "root",
							Password: "********",
						},
						Err: nil,
					},
				},
			},
			NodeInfo{
				Hostname:     "x1000c7s1b1n1",
				Nid:          1232,
				Role:         "Compute",
				State:        "Ready",
				Enabled:      true,
				BmcFQDN:      "10.100.107.181",
				BmcPath:      "/redfish/v1/Systems/Node1",
				RfActionURI:  "/redfish/v1/Systems/Node1/Actions/ComputerSystem.Reset",
				RfResetTypes: []string{"ForceOff", "On", "Off"},
				Type:         "Node",
				RfType:       "ComputerSystem",
				RfSubtype:    "Physical",
				BmcUser:      "root",
				BmcPass:      "********",
				BmcType:      "NodeBMC",
			}, false,
		}, {
			"Missing node",
			HSMQuery{ComponentIDs: []string{"x1000c7s1b1n0"}},
			ni,
			[]sstorage.MockLookup{},
			NodeInfo{},
			true,
		}, {
			"Good PDU Outlet",
			HSMQuery{ComponentIDs: []string{"x0m0p0j19"}},
			ni,
			[]sstorage.MockLookup{
				{
					Output: sstorage.OutputLookup{
						Output: &compcreds.CompCredentials{
							Xname:    "x0m0p0j19",
							URL:      "x0m0:8082/redfish/v1/PowerEquipment/RackPDUs/A/Outlets/AA26",
							Username: "root",
							Password: "********",
						},
						Err: nil,
					},
				},
			},
			NodeInfo{
				Hostname:     "x0m0p0j19",
				Nid:          0,
				Role:         "",
				State:        "On",
				Enabled:      true,
				BmcFQDN:      "x0m0:8082",
				BmcPath:      "/redfish/v1/PowerEquipment/RackPDUs/A/Outlets/AA26",
				RfActionURI:  "/redfish/v1/PowerEquipment/RackPDUs/A/Outlets/AA26/Actions/Outlet.PowerControl",
				RfResetTypes: []string{"On", "Off"},
				Type:         "CabinetPDUOutlet",
				RfType:       "Outlet",
				RfSubtype:    "C13",
				BmcUser:      "root",
				BmcPass:      "********",
				BmcType:      "CabinetPDUController",
			},
			false,
		}, {
			"Good PDU Power Connector",
			HSMQuery{ComponentIDs: []string{"x0m0p0v19"}},
			ni,
			[]sstorage.MockLookup{
				{
					Output: sstorage.OutputLookup{
						Output: &compcreds.CompCredentials{
							Xname:    "x0m0p0v19",
							URL:      "x0m0:8082/redfish/v1/PowerEquipment/RackPDUs/A/Outlets/AA26",
							Username: "root",
							Password: "********",
						},
						Err: nil,
					},
				},
			},
			NodeInfo{
				Hostname:     "x0m0p0v19",
				Nid:          0,
				Role:         "",
				State:        "On",
				Enabled:      true,
				BmcFQDN:      "x0m0:8082",
				BmcPath:      "/redfish/v1/PowerEquipment/RackPDUs/A/Outlets/AA26",
				RfActionURI:  "/redfish/v1/PowerEquipment/RackPDUs/A/Outlets/AA26/Actions/Outlet.PowerControl",
				RfResetTypes: []string{"On", "Off"},
				Type:         "CabinetPDUPowerConnector",
				RfType:       "Outlet",
				RfSubtype:    "C13",
				BmcUser:      "root",
				BmcPass:      "********",
				BmcType:      "CabinetPDUController",
			},
			false,
		}, {
			"Missing PDU Outlet",
			HSMQuery{ComponentIDs: []string{"x0m0p9j19"}},
			ni,
			[]sstorage.MockLookup{},
			NodeInfo{},
			true,
		},
	}

	for _, tt := range tests {
		adapter.LookupNum = 0
		adapter.LookupData = tt.ssData
		t.Run(tt.name, func(t *testing.T) {
			ret, err := tSvc.GetNodes(tt.in)
			if err != nil && tt.expectErr == false {
				t.Fatal("FAIL - GetNodes failed to return expected information:", err)
			}

			if tt.expectErr == false {
				for _, e := range ret {
					if !reflect.DeepEqual(tt.res, *e) {
						t.Errorf("FAIL - Expected\n%v\nbut got\n%v", tt.res, *e)
					}
				}
			}
		})
	}
}

var testNodeInfoDisabledList1 = []NodeInfo{
	{
		Hostname:     "x1000c7s1b0n0",
		Nid:          1230,
		Role:         "Compute",
		State:        "Ready",
		Enabled:      true,
		BmcFQDN:      "10.100.107.180",
		BmcPath:      "/redfish/v1/Systems/Node1",
		RfActionURI:  "/redfish/v1/Systems/Node1/Actions/ComputerSystem.Reset",
		RfResetTypes: []string{"ForceOff", "On", "Off"},
		Type:         "Node",
		RfType:       "ComputerSystem",
		RfSubtype:    "Physical",
		BmcUser:      "root",
		BmcPass:      "********",
		BmcType:      "NodeBMC",
	},
	{
		Hostname:     "x1000c7s1b0n1",
		Nid:          1231,
		Role:         "Compute",
		State:        "Ready",
		Enabled:      true,
		BmcFQDN:      "10.100.107.181",
		BmcPath:      "/redfish/v1/Systems/Node1",
		RfActionURI:  "/redfish/v1/Systems/Node1/Actions/ComputerSystem.Reset",
		RfResetTypes: []string{"ForceOff", "On", "Off"},
		Type:         "Node",
		RfType:       "ComputerSystem",
		RfSubtype:    "Physical",
		BmcUser:      "root",
		BmcPass:      "********",
		BmcType:      "NodeBMC",
	},
	{
		Hostname:     "x1000c7s1b1n0",
		Nid:          1232,
		Role:         "Compute",
		State:        "Ready",
		Enabled:      false,
		BmcFQDN:      "10.100.107.182",
		BmcPath:      "/redfish/v1/Systems/Node1",
		RfActionURI:  "/redfish/v1/Systems/Node1/Actions/ComputerSystem.Reset",
		RfResetTypes: []string{"ForceOff", "On", "Off"},
		Type:         "Node",
		RfType:       "ComputerSystem",
		RfSubtype:    "Physical",
		BmcUser:      "root",
		BmcPass:      "********",
		BmcType:      "NodeBMC",
	},
	{
		Hostname:     "x1000c7s1b1n1",
		Nid:          1233,
		Role:         "Compute",
		State:        "Ready",
		Enabled:      true,
		BmcFQDN:      "10.100.107.183",
		BmcPath:      "/redfish/v1/Systems/Node1",
		RfActionURI:  "/redfish/v1/Systems/Node1/Actions/ComputerSystem.Reset",
		RfResetTypes: []string{"ForceOff", "On", "Off"},
		Type:         "Node",
		RfType:       "ComputerSystem",
		RfSubtype:    "Physical",
		BmcUser:      "root",
		BmcPass:      "********",
		BmcType:      "NodeBMC",
	},
}

var testNodeInfoDisabledList2 = []NodeInfo{
	{
		Hostname:     "x1000c7s1b0n0",
		Nid:          1230,
		Role:         "Compute",
		State:        "Ready",
		Enabled:      true,
		BmcFQDN:      "10.100.107.180",
		BmcPath:      "/redfish/v1/Systems/Node1",
		RfActionURI:  "/redfish/v1/Systems/Node1/Actions/ComputerSystem.Reset",
		RfResetTypes: []string{"ForceOff", "On", "Off"},
		Type:         "Node",
		RfType:       "ComputerSystem",
		RfSubtype:    "Physical",
		BmcUser:      "root",
		BmcPass:      "********",
		BmcType:      "NodeBMC",
	},
	{
		Hostname:     "x1000c7s1b0n1",
		Nid:          1231,
		Role:         "Compute",
		State:        "Ready",
		Enabled:      true,
		BmcFQDN:      "10.100.107.181",
		BmcPath:      "/redfish/v1/Systems/Node1",
		RfActionURI:  "/redfish/v1/Systems/Node1/Actions/ComputerSystem.Reset",
		RfResetTypes: []string{"ForceOff", "On", "Off"},
		Type:         "Node",
		RfType:       "ComputerSystem",
		RfSubtype:    "Physical",
		BmcUser:      "root",
		BmcPass:      "********",
		BmcType:      "NodeBMC",
	},
	{
		Hostname:     "x1000c7s1b1n0",
		Nid:          1232,
		Role:         "Compute",
		State:        "Ready",
		Enabled:      true,
		BmcFQDN:      "10.100.107.182",
		BmcPath:      "/redfish/v1/Systems/Node1",
		RfActionURI:  "/redfish/v1/Systems/Node1/Actions/ComputerSystem.Reset",
		RfResetTypes: []string{"ForceOff", "On", "Off"},
		Type:         "Node",
		RfType:       "ComputerSystem",
		RfSubtype:    "Physical",
		BmcUser:      "root",
		BmcPass:      "********",
		BmcType:      "NodeBMC",
	},
	{
		Hostname:     "x1000c7s1b1n1",
		Nid:          1233,
		Role:         "Compute",
		State:        "Ready",
		Enabled:      true,
		BmcFQDN:      "10.100.107.183",
		BmcPath:      "/redfish/v1/Systems/Node1",
		RfActionURI:  "/redfish/v1/Systems/Node1/Actions/ComputerSystem.Reset",
		RfResetTypes: []string{"ForceOff", "On", "Off"},
		Type:         "Node",
		RfType:       "ComputerSystem",
		RfSubtype:    "Physical",
		BmcUser:      "root",
		BmcPass:      "********",
		BmcType:      "NodeBMC",
	},
}

func TestCheckForDisabledComponents(t *testing.T) {
	var tSvc CapmcD
	var nl1 []*NodeInfo
	var nl2 []*NodeInfo

	for i := 0; i < len(testNodeInfoDisabledList1); i++ {
		nl1 = append(nl1, &testNodeInfoDisabledList1[i])
	}
	for i := 0; i < len(testNodeInfoDisabledList2); i++ {
		nl2 = append(nl2, &testNodeInfoDisabledList2[i])
	}

	tests := []struct {
		name    string
		nl      []*NodeInfo
		t       string
		e       string
		wantErr bool
	}{
		{
			"single disabled xname",
			nl1,
			"xname",
			"components disabled: [x1000c7s1b1n0]",
			true,
		},
		{
			"no disabled xname",
			nl2,
			"xname",
			"",
			false,
		},
		{
			"single disabled node",
			nl1,
			"nid",
			"nodes disabled: [1232]",
			true,
		},
		{
			"no disabled node",
			nl2,
			"nid",
			"",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tSvc.checkForDisabledComponents(tt.nl, tt.t)

			if (err != nil) != tt.wantErr {
				t.Errorf("checkForDisabledComponents() error = %v, wantErr %v", err, tt.wantErr)
			}

			if (err != nil) && err.Error() != tt.e {
				t.Errorf("checkForDisabledComponents() error = %v, wantErr %v", err, tt.e)
			}

			if (err == nil) && "" != tt.e {
				t.Errorf("checkForDisabledComponents() error = %v, wantErr %v", err, tt.e)
			}
		})
	}
}

func TestConvertControlsToPowerCaps(t *testing.T) {
	var emptyCtl []*rf.Control
	nodeOnlyWant := make(map[string]PowerCap)
	nodeOnlyWant["Node Power Limit"] = PowerCap{
		Name:        "Node Power Limit",
		Path:        "/redfish/v1/Chassis/Node0/Controls/NodePowerLimit",
		Min:         200,
		Max:         1000,
		PwrCtlIndex: 0,
	}
	nodesAccelsWant := make(map[string]PowerCap)
	nodesAccelsWant["Node Power Limit"] = PowerCap{
		Name:        "Node Power Limit",
		Path:        "/redfish/v1/Chassis/Node0/Controls/NodePowerLimit",
		Min:         200,
		Max:         1000,
		PwrCtlIndex: 0,
	}
	nodesAccelsWant["GPU0 Power Limit"] = PowerCap{
		Name:        "GPU0 Power Limit",
		Path:        "/redfish/v1/Chassis/Node0/Controls/GPU0PowerLimit",
		Min:         200,
		Max:         400,
		PwrCtlIndex: 1,
	}
	nodesAccelsWant["GPU1 Power Limit"] = PowerCap{
		Name:        "GPU1 Power Limit",
		Path:        "/redfish/v1/Chassis/Node0/Controls/GPU1PowerLimit",
		Min:         200,
		Max:         400,
		PwrCtlIndex: 2,
	}
	nodesAccelsWant["GPU2 Power Limit"] = PowerCap{
		Name:        "GPU2 Power Limit",
		Path:        "/redfish/v1/Chassis/Node0/Controls/GPU2PowerLimit",
		Min:         200,
		Max:         400,
		PwrCtlIndex: 3,
	}
	nodesAccelsWant["GPU3 Power Limit"] = PowerCap{
		Name:        "GPU3 Power Limit",
		Path:        "/redfish/v1/Chassis/Node0/Controls/GPU3PowerLimit",
		Min:         200,
		Max:         400,
		PwrCtlIndex: 4,
	}

	var nodeOnlyCtl []*rf.Control
	nodeOnlyCtl = append(nodeOnlyCtl, &rf.Control{
		URL: "/redfish/v1/Chassis/Node0/Controls/NodePowerLimit",
		Control: rf.RFControl{
			ControlDelaySeconds: 6,
			ControlMode:         "Automatic",
			ControlType:         "Power",
			Id:                  "NodePowerLimit",
			Name:                "Node Power Limit",
			PhysicalContext:     "Chassis",
			SetPoint:            1000,
			SetPointUnits:       "W",
			SettingRangeMax:     1000,
			SettingRangeMin:     200,
			Status: rf.StatusRF{
				Health: "OK",
				State:  "",
			},
		},
	})
	nodesAccelsCtl := nodeOnlyCtl
	nodesAccelsCtl = append(nodesAccelsCtl, &rf.Control{
		URL: "/redfish/v1/Chassis/Node0/Controls/GPU0PowerLimit",
		Control: rf.RFControl{
			ControlDelaySeconds: 6,
			ControlMode:         "Automatic",
			ControlType:         "Power",
			Id:                  "GPU0PowerLimit",
			Name:                "GPU0 Power Limit",
			PhysicalContext:     "GPU",
			SetPoint:            400,
			SetPointUnits:       "W",
			SettingRangeMax:     400,
			SettingRangeMin:     200,
			Status: rf.StatusRF{
				Health: "OK",
				State:  "",
			},
		},
	})
	nodesAccelsCtl = append(nodesAccelsCtl, &rf.Control{
		URL: "/redfish/v1/Chassis/Node0/Controls/GPU1PowerLimit",
		Control: rf.RFControl{
			ControlDelaySeconds: 6,
			ControlMode:         "Automatic",
			ControlType:         "Power",
			Id:                  "GPU1PowerLimit",
			Name:                "GPU1 Power Limit",
			PhysicalContext:     "GPU",
			SetPoint:            400,
			SetPointUnits:       "W",
			SettingRangeMax:     400,
			SettingRangeMin:     200,
			Status: rf.StatusRF{
				Health: "OK",
				State:  "",
			},
		},
	})
	nodesAccelsCtl = append(nodesAccelsCtl, &rf.Control{
		URL: "/redfish/v1/Chassis/Node0/Controls/GPU2PowerLimit",
		Control: rf.RFControl{
			ControlDelaySeconds: 6,
			ControlMode:         "Automatic",
			ControlType:         "Power",
			Id:                  "GPU2PowerLimit",
			Name:                "GPU2 Power Limit",
			PhysicalContext:     "GPU",
			SetPoint:            400,
			SetPointUnits:       "W",
			SettingRangeMax:     400,
			SettingRangeMin:     200,
			Status: rf.StatusRF{
				Health: "OK",
				State:  "",
			},
		},
	})
	nodesAccelsCtl = append(nodesAccelsCtl, &rf.Control{
		URL: "/redfish/v1/Chassis/Node0/Controls/GPU3PowerLimit",
		Control: rf.RFControl{
			ControlDelaySeconds: 6,
			ControlMode:         "Automatic",
			ControlType:         "Power",
			Id:                  "GPU3PowerLimit",
			Name:                "GPU3 Power Limit",
			PhysicalContext:     "GPU",
			SetPoint:            400,
			SetPointUnits:       "W",
			SettingRangeMax:     400,
			SettingRangeMin:     200,
			Status: rf.StatusRF{
				Health: "OK",
				State:  "",
			},
		},
	})
	tests := []struct {
		name     string
		ni       *NodeInfo
		controls []*rf.Control
		want     map[string]PowerCap
	}{
		{
			"empty caps",
			&NodeInfo{},
			emptyCtl,
			nil,
		},
		{
			"node only",
			&NodeInfo{},
			nodeOnlyCtl,
			nodeOnlyWant,
		},
		{
			"nodes and accels",
			&NodeInfo{},
			nodesAccelsCtl,
			nodesAccelsWant,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertControlsToPowerCaps(tt.ni, tt.controls); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertControlsToPowerCaps() = %v, want %v", got, tt.want)
			}
		})
	}
}

const riverNodeCompEndpoint = `
{
  "ID": "x3000c0s9b0n0",
  "Type": "Node",
  "RedfishType": "ComputerSystem",
  "RedfishSubtype": "Physical",
  "UUID": "36383150-3630-584D-5130-31393032304A",
  "OdataID": "/redfish/v1/Systems/1",
  "RedfishEndpointID": "x3000c0s9b0",
  "Enabled": true,
  "RedfishEndpointFQDN": "x3000c0s9b0",
  "RedfishURL": "x3000c0s9b0/redfish/v1/Systems/1",
  "ComponentEndpointType": "ComponentEndpointComputerSystem",
  "RedfishSystemInfo": {
    "Name": "Computer System",
    "Actions": {
      "#ComputerSystem.Reset": {
        "ResetType@Redfish.AllowableValues": [
          "On",
          "ForceOff",
          "GracefulShutdown",
          "ForceRestart",
          "Nmi",
          "PushPowerButton",
          "GracefulRestart"
        ],
        "@Redfish.ActionInfo": "",
        "target": "/redfish/v1/Systems/1/Actions/ComputerSystem.Reset"
      }
    },
    "EthernetNICInfo": [
      {
        "RedfishId": "1",
        "@odata.id": "/redfish/v1/Systems/1/EthernetInterfaces/1",
        "MACAddress": "94:40:c9:5b:e5:70"
      },
      {
        "RedfishId": "2",
        "@odata.id": "/redfish/v1/Systems/1/EthernetInterfaces/2",
        "InterfaceEnabled": true,
        "MACAddress": "94:40:c9:5b:e5:71"
      },
      {
        "RedfishId": "3",
        "@odata.id": "/redfish/v1/Systems/1/EthernetInterfaces/3",
        "InterfaceEnabled": true,
        "MACAddress": "14:02:ec:d9:3e:90"
      },
      {
        "RedfishId": "4",
        "@odata.id": "/redfish/v1/Systems/1/EthernetInterfaces/4",
        "InterfaceEnabled": true,
        "MACAddress": "14:02:ec:d9:3e:91"
      }
    ],
    "PowerURL": "/redfish/v1/Chassis/1/Power",
    "PowerControl": [
      {
        "@odata.id": "/redfish/v1/Chassis/1/Power#PowerControl/0",
        "MemberId": "0",
        "PowerCapacityWatts": 1000
      }
    ]
  }
}
`

const mountainBardPeakCompEndpoint = `
{
  "ID": "x9000c3s0b0n0",
  "Type": "Node",
  "RedfishType": "ComputerSystem",
  "RedfishSubtype": "Physical",
  "OdataID": "/redfish/v1/Systems/Node0",
  "RedfishEndpointID": "x9000c3s0b0",
  "Enabled": true,
  "RedfishEndpointFQDN": "x9000c3s0b0",
  "RedfishURL": "x9000c3s0b0/redfish/v1/Systems/Node0",
  "ComponentEndpointType": "ComponentEndpointComputerSystem",
  "RedfishSystemInfo": {
    "Name": "Node0",
    "Actions": {
      "#ComputerSystem.Reset": {
        "ResetType@Redfish.AllowableValues": [
          "ForceOff",
          "Off",
          "On"
        ],
        "@Redfish.ActionInfo": "/redfish/v1/Systems/Node0/ResetActionInfo",
        "target": "/redfish/v1/Systems/Node0/Actions/ComputerSystem.Reset"
      }
    },
    "EthernetNICInfo": [
      {
        "RedfishId": "HPCNet0",
        "@odata.id": "/redfish/v1/Systems/Node0/EthernetInterfaces/HPCNet0",
        "Description": "SS11 200Gb 2P NIC Mezz REV02 (HSN)",
        "MACAddress": "Not Available"
      },
      {
        "RedfishId": "HPCNet1",
        "@odata.id": "/redfish/v1/Systems/Node0/EthernetInterfaces/HPCNet1",
        "Description": "SS11 200Gb 2P NIC Mezz REV02 (HSN)",
        "MACAddress": "Not Available"
      },
      {
        "RedfishId": "HPCNet2",
        "@odata.id": "/redfish/v1/Systems/Node0/EthernetInterfaces/HPCNet2",
        "Description": "SS11 200Gb 2P NIC Mezz REV02 (HSN)",
        "MACAddress": "Not Available"
      },
      {
        "RedfishId": "HPCNet3",
        "@odata.id": "/redfish/v1/Systems/Node0/EthernetInterfaces/HPCNet3",
        "Description": "SS11 200Gb 2P NIC Mezz REV02 (HSN)",
        "MACAddress": "Not Available"
      },
      {
        "RedfishId": "ManagementEthernet",
        "@odata.id": "/redfish/v1/Systems/Node0/EthernetInterfaces/ManagementEthernet",
        "Description": "Node Maintenance Network",
        "MACAddress": "00:40:a6:83:3a:52",
        "PermanentMACAddress": "00:40:a6:83:3a:52"
      }
    ],
    "PowerURL": "/redfish/v1/Chassis/Node0/Power",
    "Controls": [
      {
        "URL": "/redfish/v1/Chassis/Node0/Controls/NodePowerLimit",
        "Control": {
          "ControlDelaySeconds": 6,
          "ControlMode": "Disabled",
          "ControlType": "Power",
          "Id": "NodePowerLimit",
          "Name": "Node Power Limit",
          "PhysicalContext": "Chassis",
          "SetPoint": 0,
          "SetPointUnits": "",
          "SettingRangeMax": 2754,
          "SettingRangeMin": 764,
          "Status": {
            "Health": "OK"
          }
        }
      },
      {
        "URL": "/redfish/v1/Chassis/Node0/Controls/Accelerator0PowerLimit",
        "Control": {
          "ControlDelaySeconds": 6,
          "ControlMode": "Disabled",
          "ControlType": "Power",
          "Id": "Accelerator0PowerLimit",
          "Name": "Accelerator0 Power Limit",
          "PhysicalContext": "Accelerator",
          "SetPoint": 0,
          "SetPointUnits": "",
          "SettingRangeMax": 560,
          "SettingRangeMin": 100,
          "Status": {
            "Health": "OK"
          }
        }
      },
      {
        "URL": "/redfish/v1/Chassis/Node0/Controls/Accelerator1PowerLimit",
        "Control": {
          "ControlDelaySeconds": 6,
          "ControlMode": "Disabled",
          "ControlType": "Power",
          "Id": "Accelerator1PowerLimit",
          "Name": "Accelerator1 Power Limit",
          "PhysicalContext": "Accelerator",
          "SetPoint": 0,
          "SetPointUnits": "",
          "SettingRangeMax": 560,
          "SettingRangeMin": 100,
          "Status": {
            "Health": "OK"
          }
        }
      },
      {
        "URL": "/redfish/v1/Chassis/Node0/Controls/Accelerator2PowerLimit",
        "Control": {
          "ControlDelaySeconds": 6,
          "ControlMode": "Disabled",
          "ControlType": "Power",
          "Id": "Accelerator2PowerLimit",
          "Name": "Accelerator2 Power Limit",
          "PhysicalContext": "Accelerator",
          "SetPoint": 0,
          "SetPointUnits": "",
          "SettingRangeMax": 560,
          "SettingRangeMin": 100,
          "Status": {
            "Health": "OK"
          }
        }
      },
      {
        "URL": "/redfish/v1/Chassis/Node0/Controls/Accelerator3PowerLimit",
        "Control": {
          "ControlDelaySeconds": 6,
          "ControlMode": "Disabled",
          "ControlType": "Power",
          "Id": "Accelerator3PowerLimit",
          "Name": "Accelerator3 Power Limit",
          "PhysicalContext": "Accelerator",
          "SetPoint": 0,
          "SetPointUnits": "",
          "SettingRangeMax": 560,
          "SettingRangeMin": 100,
          "Status": {
            "Health": "OK"
          }
        }
      }
    ]
  }
}
`

func TestExtractRedfishSystemInfo(t *testing.T) {
	var niRiver NodeInfo
	var niMtnBP NodeInfo
	tests := []struct {
		name string
		ni   *NodeInfo
		json string
		want NodeInfo
	}{
		{"River node", &niRiver, riverNodeCompEndpoint, NodeInfo{
			Nid:           0,
			Enabled:       false,
			RfActionURI:   "/redfish/v1/Systems/1/Actions/ComputerSystem.Reset",
			RfResetTypes:  []string{"On", "ForceOff", "GracefulShutdown", "ForceRestart", "Nmi", "PushPowerButton", "GracefulRestart"},
			RfPowerURL:    "/redfish/v1/Chassis/1/Power",
			RfPwrCtlCnt:   1,
			RfControlsCnt: 0,
			PowerCaps: map[string]PowerCap{
				"Node Power Control": {Name: "Node Power Control", Path: "/redfish/v1/Chassis/1/Power", Min: -1, Max: -1, PwrCtlIndex: 0},
			},
		}},
		{"Mtn BP node", &niMtnBP, mountainBardPeakCompEndpoint, NodeInfo{
			Nid:           0,
			Enabled:       false,
			RfActionURI:   "/redfish/v1/Systems/Node0/Actions/ComputerSystem.Reset",
			RfResetTypes:  []string{"ForceOff", "Off", "On"},
			RfPowerURL:    "/redfish/v1/Chassis/Node0/Power",
			RfPwrCtlCnt:   0,
			RfControlsCnt: 5,
			PowerCaps: map[string]PowerCap{
				"Accelerator0 Power Limit": {Name: "Accelerator0 Power Limit", Path: "/redfish/v1/Chassis/Node0/Controls/Accelerator0PowerLimit", Min: 100, Max: 560, PwrCtlIndex: 1},
				"Accelerator1 Power Limit": {Name: "Accelerator1 Power Limit", Path: "/redfish/v1/Chassis/Node0/Controls/Accelerator1PowerLimit", Min: 100, Max: 560, PwrCtlIndex: 2},
				"Accelerator2 Power Limit": {Name: "Accelerator2 Power Limit", Path: "/redfish/v1/Chassis/Node0/Controls/Accelerator2PowerLimit", Min: 100, Max: 560, PwrCtlIndex: 3},
				"Accelerator3 Power Limit": {Name: "Accelerator3 Power Limit", Path: "/redfish/v1/Chassis/Node0/Controls/Accelerator3PowerLimit", Min: 100, Max: 560, PwrCtlIndex: 4},
				"Node Power Limit":         {Name: "Node Power Limit", Path: "/redfish/v1/Chassis/Node0/Controls/NodePowerLimit", Min: 764, Max: 2754, PwrCtlIndex: 0},
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ce sm.ComponentEndpoint

			err := json.Unmarshal([]byte(tt.json), &ce)
			if err != nil {
				t.Errorf("Failed to unmarshal json %s\n", tt.json)
			}

			extractRedfishSystemInfo(tt.ni, &ce)

			if !reflect.DeepEqual(&tt.want, tt.ni) {
				t.Errorf("FAIL - Expected\n%v\nbut got\n%v", tt.want, tt.ni)
			}
		})
	}
}
