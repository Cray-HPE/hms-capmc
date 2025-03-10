/*
 * MIT License
 *
 * (C) Copyright [2019-2023,2025] Hewlett Packard Enterprise Development LP
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

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"sort"
	"testing"

	base "github.com/Cray-HPE/hms-base/v2"
	"github.com/Cray-HPE/hms-capmc/internal/capmc"
	compcreds "github.com/Cray-HPE/hms-compcredentials"
	sstorage "github.com/Cray-HPE/hms-securestorage"
)

func TestMakeXnameErrors(t *testing.T) {
	tests := []struct {
		cidErr *InvalidCompIDsError
		result []*capmc.XnameControlErr
	}{
		{
			cidErr: &InvalidCompIDsError{
				err:     "foo xnames test",
				CompIDs: []string{"c0-0c0s5n2", "x3000c0s19b4n0"},
			},
			result: []*capmc.XnameControlErr{
				{
					Xname: "c0-0c0s5n2",
					ErrResponse: capmc.ErrResponse{
						E:      22,
						ErrMsg: "foo xname test",
					},
				},
				{
					Xname: "x3000c0s19b4n0",
					ErrResponse: capmc.ErrResponse{
						E:      22,
						ErrMsg: "foo xname test",
					},
				},
			},
		},
	}

	for n, test := range tests {
		r := MakeXnameErrors(test.cidErr)
		if len(r) != len(test.result) {
			t.Errorf("MakeXnameErrors Test Case %d: FAIL: expected %d but got %d", n, len(test.result), len(r))
		}

		for i := 0; i < len(r); i++ {
			if *r[i] != *test.result[i] {
				t.Errorf("MakeXnameErrors Test Case %d: FAIL: expected [%d] %+v but got %+v", n, i, test.result[i], r[i])
			}
		}
	}
}

const mountainHSNBoardStateComponents = `
{
  "Components": [
    {
      "Arch": "X86",
      "Enabled": true,
      "Flag": "OK",
      "NID": 0,
      "State": "Ready",
      "Role": "Compute",
      "NetType": "Sling",
	  "Type": "Node",
	  "Class": "Mountain",
      "ID": "x1002c0r0e0"
    }
  ]
}
`

const mountainHSNBoardComponentEndpoint = `
{
  "ComponentEndpoints": [
	{
	"RedfishType": "Chassis",
	"RedfishSubtype": "Enclosure",
	"Enabled": true,
	"OdataID": "/redfish/v1/Chassis/Enclosure",
	"RedfishEndpointID": "x1002c0r0b0",
	"RedfishChassisInfo": {
		"Name": "Enclosure",
		"Actions": {
		"#Chassis.Reset": {
			"ResetType@Redfish.AllowableValues": [
			"PowerCycle",
			"ForceRestart",
			"On",
			"ForceOn",
			"GracefulShutdown",
			"GracefulRestart",
			"Off",
			"ForceOff"
			],
			"target": "/redfish/v1/Chassis/Enclosure/Actions/Chassis.Reset",
			"@Redfish.ActionInfo": ""
		}
		}
	},
	"ComponentEndpointType": "ComponentEndpointChassis",
	"RedfishURL": "x1002c0r0b0/redfish/v1/Chassis/Enclosure",
	"Type": "HSNBoard",
	"ID": "x1002c0r0e0",
	"RedfishEndpointFQDN": "x1002c0r0b0"
	}
  ]
}
`

var ssDataOff = []sstorage.MockLookup{
	{Output: sstorage.OutputLookup{Output: &compcreds.CompCredentials{Xname: "x1002c0r0e0"}}},
	{Output: sstorage.OutputLookup{Output: &compcreds.CompCredentials{Xname: "x1002c0r0"}}},
}

func powerOffFunc() RoundTripFunc {
	return func(req *http.Request) (*http.Response, error) {
		//fmt.Println(req.URL.String())
		switch req.URL.String() {
		case "http://localhost:27779/State/Components?id=x1002c0r0e0",
			"http://localhost:27779/State/Components?id=x1002c0r0e0&role=%21Management":
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(mountainHSNBoardStateComponents)),
				Header:     make(http.Header),
			}, nil
		case "http://localhost:27779/Inventory/ComponentEndpoints?id=x1002c0r0e0",
			"http://localhost:27779/Inventory/ComponentEndpoints?id=x1002c0r0e0&role=%21Management":
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(mountainHSNBoardComponentEndpoint)),
				Header:     make(http.Header),
			}, nil
		case "http://locxalhost:27779/State/Components?id=x1000c0r0e0":
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(mountainHSNBoard)),
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

func TestDoXnameOff(t *testing.T) {
	var tSvc CapmcD
	var err error
	tSvc.hsmURL, err = url.Parse("http://localhost:27779")
	if err != nil {
		t.Fatal(err)
	}
	testClient := NewTestClient(powerOffFunc())
	tSvc.rfClient = testClient
	tSvc.smClient = testClient
	if err != nil {
		t.Fatal(err)
	}
	tSvc.config = loadConfig("")
	ss, adapter := sstorage.NewMockAdapter()
	ccs := compcreds.NewCompCredStore("secret/hms-cred", ss)
	tSvc.ss = ss
	tSvc.ccs = ccs
	tSvc.WPool = base.NewWorkerPool(1, 10)
	tSvc.WPool.Run()
	handler := http.HandlerFunc(tSvc.doXnameOff)

	log.SetOutput(os.Stdout)

	adapter.LookupData = ssDataOff

	tests := []struct {
		name     string
		body     io.Reader
		code     int
		expected string
	}{
		{
			"Missing xnames parameter",
			bytes.NewBufferString("{ \"reason\": \"power save, need less capacity\", \"nids\": [ \"x3000c0s38b1n0\", \"x0c0s04b0n0\" ]}\n"),
			http.StatusBadRequest,
			"{\"e\":400,\"err_msg\":\"Bad Request: Missing required xnames parameter\"}\n",
		},
		{
			"Missing xnames value list",
			bytes.NewBuffer(json.RawMessage("{ \"reason\": \"power save, need less capacity\", \"xnames\": [ ]}\n")),
			http.StatusBadRequest,
			"{\"e\":400,\"err_msg\":\"Bad Request: Required xnames list is empty\"}\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, capmc.XnameOffV1, tc.body)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if rr.Code != tc.code {
				t.Errorf("handler returned wrong status code: want %v but got %v",
					tc.code, rr.Code)
			}

			if rr.Body.String() != tc.expected {
				t.Errorf("handler returned unexpected body: want '%v' but got '%v'",
					tc.expected, rr.Body.String())
			}
		})
	}
	log.SetOutput(ioutil.Discard)
}

var ssDataStatus = []sstorage.MockLookup{
	{Output: sstorage.OutputLookup{Output: &compcreds.CompCredentials{Xname: "x1002c0s0b0n0"}}},
	{Output: sstorage.OutputLookup{Output: &compcreds.CompCredentials{Xname: "x1002c0s0b0n1"}}},
	{Output: sstorage.OutputLookup{Output: &compcreds.CompCredentials{Xname: "x1002c0s0b1n0"}}},
	{Output: sstorage.OutputLookup{Output: &compcreds.CompCredentials{Xname: "x1002c0s0b1n1"}}},
	{Output: sstorage.OutputLookup{Output: &compcreds.CompCredentials{Xname: "x1002c0s1b0n0"}}},
	{Output: sstorage.OutputLookup{Output: &compcreds.CompCredentials{Xname: "x1002c0s1b0n1"}}},
	{Output: sstorage.OutputLookup{Output: &compcreds.CompCredentials{Xname: "x1002c0s1b1n0"}}},
	{Output: sstorage.OutputLookup{Output: &compcreds.CompCredentials{Xname: "x1002c0s1b1n1"}}},
}

const compStatusEnabledReadyOK = `{"Components":[{"ID":"x1002c0s0b0n0","Type":"Node","State":"Ready","Flag":"OK","Enabled":true,"Role":"Compute","NID":1512,"NetType":"Sling","Arch":"X86","Class":"Mountain"}]}`
const x1002c0s0b0n0CompEndpoint = `{"ComponentEndpoints":[{"ID":"x1002c0s0b0n0","Type":"Node","RedfishType":"ComputerSystem","RedfishSubtype":"Physical","OdataID":"/redfish/v1/Systems/Node0","RedfishEndpointID":"x1002c0s0b0","Enabled":true,"RedfishEndpointFQDN":"10.104.8.11","RedfishURL":"10.104.8.11/redfish/v1/Systems/Node0","ComponentEndpointType":"ComponentEndpointComputerSystem","RedfishSystemInfo":{"Name":"Node0","Actions":{"#ComputerSystem.Reset":{"ResetType@Redfish.AllowableValues":["ForceOff","On","Off"],"@Redfish.ActionInfo":"/redfish/v1/Systems/Node0/ResetActionInfo","target":"/redfish/v1/Systems/Node0/Actions/ComputerSystem.Reset"}},"EthernetNICInfo":[{"RedfishId":"HPCNet0","@odata.id":"/redfish/v1/Systems/Node0/EthernetInterfaces/HPCNet0","Description":"Shasta Timms NMC REV04 (HSN)","MACAddress":"Not Available","PermanentMACAddress":" 00:40:a6:83:c5:30"},{"RedfishId":"HPCNet1","@odata.id":"/redfish/v1/Systems/Node0/EthernetInterfaces/HPCNet1","Description":"Shasta Timms NMC REV04 (HSN)","MACAddress":"Not Available","PermanentMACAddress":" 00:40:a6:83:c5:30"},{"RedfishId":"ManagementEthernet","@odata.id":"/redfish/v1/Systems/Node0/EthernetInterfaces/ManagementEthernet","Description":"Node Maintenance Network","MACAddress":"00:40:a6:83:b6:37","PermanentMACAddress":"00:40:a6:83:b6:37"}],"PowerURL":"/redfish/v1/Chassis/Node0/Power","PowerControl":[{"@odata.id":"/redfish/v1/Chassis/Node0/Power#/PowerControl/Node","MemberId":"Node","Name":"Node Power Control","OEM":{"Cray":{"PowerIdleWatts":105,"PowerLimit":{"Min":350,"Max":925},"PowerResetWatts":250}}}]}}]}`

const compStatusEnabledOnOK = `{"Components":[{"ID":"x1002c0s0b0n1","Type":"Node","State":"On","Flag":"OK","Enabled":true,"Role":"Compute","NID":1512,"NetType":"Sling","Arch":"X86","Class":"Mountain"}]}`
const x1002c0s0b0n1CompEndpoint = `{"ComponentEndpoints":[{"ID":"x1002c0s0b0n1","Type":"Node","RedfishType":"ComputerSystem","RedfishSubtype":"Physical","OdataID":"/redfish/v1/Systems/Node0","RedfishEndpointID":"x1002c0s0b0","Enabled":true,"RedfishEndpointFQDN":"10.104.8.11","RedfishURL":"10.104.8.11/redfish/v1/Systems/Node0","ComponentEndpointType":"ComponentEndpointComputerSystem","RedfishSystemInfo":{"Name":"Node0","Actions":{"#ComputerSystem.Reset":{"ResetType@Redfish.AllowableValues":["ForceOff","On","Off"],"@Redfish.ActionInfo":"/redfish/v1/Systems/Node0/ResetActionInfo","target":"/redfish/v1/Systems/Node0/Actions/ComputerSystem.Reset"}},"EthernetNICInfo":[{"RedfishId":"HPCNet0","@odata.id":"/redfish/v1/Systems/Node0/EthernetInterfaces/HPCNet0","Description":"Shasta Timms NMC REV04 (HSN)","MACAddress":"Not Available","PermanentMACAddress":" 00:40:a6:83:c5:30"},{"RedfishId":"HPCNet1","@odata.id":"/redfish/v1/Systems/Node0/EthernetInterfaces/HPCNet1","Description":"Shasta Timms NMC REV04 (HSN)","MACAddress":"Not Available","PermanentMACAddress":" 00:40:a6:83:c5:30"},{"RedfishId":"ManagementEthernet","@odata.id":"/redfish/v1/Systems/Node0/EthernetInterfaces/ManagementEthernet","Description":"Node Maintenance Network","MACAddress":"00:40:a6:83:b6:37","PermanentMACAddress":"00:40:a6:83:b6:37"}],"PowerURL":"/redfish/v1/Chassis/Node0/Power","PowerControl":[{"@odata.id":"/redfish/v1/Chassis/Node0/Power#/PowerControl/Node","MemberId":"Node","Name":"Node Power Control","OEM":{"Cray":{"PowerIdleWatts":105,"PowerLimit":{"Min":350,"Max":925},"PowerResetWatts":250}}}]}}]}`

const compStatusEnabledOffOK = `{"Components":[{"ID":"x1002c0s0b1n0","Type":"Node","State":"Off","Flag":"OK","Enabled":true,"Role":"Compute","NID":1512,"NetType":"Sling","Arch":"X86","Class":"Mountain"}]}`
const x1002c0s0b1n0CompEndpoint = `{"ComponentEndpoints":[{"ID":"x1002c0s0b1n0","Type":"Node","RedfishType":"ComputerSystem","RedfishSubtype":"Physical","OdataID":"/redfish/v1/Systems/Node0","RedfishEndpointID":"x1002c0s0b1","Enabled":true,"RedfishEndpointFQDN":"10.104.8.12","RedfishURL":"10.104.8.12/redfish/v1/Systems/Node0","ComponentEndpointType":"ComponentEndpointComputerSystem","RedfishSystemInfo":{"Name":"Node0","Actions":{"#ComputerSystem.Reset":{"ResetType@Redfish.AllowableValues":["ForceOff","On","Off"],"@Redfish.ActionInfo":"/redfish/v1/Systems/Node0/ResetActionInfo","target":"/redfish/v1/Systems/Node0/Actions/ComputerSystem.Reset"}},"EthernetNICInfo":[{"RedfishId":"HPCNet0","@odata.id":"/redfish/v1/Systems/Node0/EthernetInterfaces/HPCNet0","Description":"Shasta Timms NMC REV04 (HSN)","MACAddress":"Not Available","PermanentMACAddress":" 00:40:a6:83:c5:30"},{"RedfishId":"HPCNet1","@odata.id":"/redfish/v1/Systems/Node0/EthernetInterfaces/HPCNet1","Description":"Shasta Timms NMC REV04 (HSN)","MACAddress":"Not Available","PermanentMACAddress":" 00:40:a6:83:c5:30"},{"RedfishId":"ManagementEthernet","@odata.id":"/redfish/v1/Systems/Node0/EthernetInterfaces/ManagementEthernet","Description":"Node Maintenance Network","MACAddress":"00:40:a6:83:b6:37","PermanentMACAddress":"00:40:a6:83:b6:37"}],"PowerURL":"/redfish/v1/Chassis/Node0/Power","PowerControl":[{"@odata.id":"/redfish/v1/Chassis/Node0/Power#/PowerControl/Node","MemberId":"Node","Name":"Node Power Control","OEM":{"Cray":{"PowerIdleWatts":105,"PowerLimit":{"Min":350,"Max":925},"PowerResetWatts":250}}}]}}]}`

const compStatusEnabledStandbyOK = `{"Components":[{"ID":"x1002c0s0b1n1","Type":"Node","State":"Standby","Flag":"OK","Enabled":true,"Role":"Compute","NID":1512,"NetType":"Sling","Arch":"X86","Class":"Mountain"}]}`
const x1002c0s0b1n1CompEndpoint = `{"ComponentEndpoints":[{"ID":"x1002c0s0b1n1","Type":"Node","RedfishType":"ComputerSystem","RedfishSubtype":"Physical","OdataID":"/redfish/v1/Systems/Node0","RedfishEndpointID":"x1002c0s0b1","Enabled":true,"RedfishEndpointFQDN":"10.104.8.12","RedfishURL":"10.104.8.12/redfish/v1/Systems/Node0","ComponentEndpointType":"ComponentEndpointComputerSystem","RedfishSystemInfo":{"Name":"Node0","Actions":{"#ComputerSystem.Reset":{"ResetType@Redfish.AllowableValues":["ForceOff","On","Off"],"@Redfish.ActionInfo":"/redfish/v1/Systems/Node0/ResetActionInfo","target":"/redfish/v1/Systems/Node0/Actions/ComputerSystem.Reset"}},"EthernetNICInfo":[{"RedfishId":"HPCNet0","@odata.id":"/redfish/v1/Systems/Node0/EthernetInterfaces/HPCNet0","Description":"Shasta Timms NMC REV04 (HSN)","MACAddress":"Not Available","PermanentMACAddress":" 00:40:a6:83:c5:30"},{"RedfishId":"HPCNet1","@odata.id":"/redfish/v1/Systems/Node0/EthernetInterfaces/HPCNet1","Description":"Shasta Timms NMC REV04 (HSN)","MACAddress":"Not Available","PermanentMACAddress":" 00:40:a6:83:c5:30"},{"RedfishId":"ManagementEthernet","@odata.id":"/redfish/v1/Systems/Node0/EthernetInterfaces/ManagementEthernet","Description":"Node Maintenance Network","MACAddress":"00:40:a6:83:b6:37","PermanentMACAddress":"00:40:a6:83:b6:37"}],"PowerURL":"/redfish/v1/Chassis/Node0/Power","PowerControl":[{"@odata.id":"/redfish/v1/Chassis/Node0/Power#/PowerControl/Node","MemberId":"Node","Name":"Node Power Control","OEM":{"Cray":{"PowerIdleWatts":105,"PowerLimit":{"Min":350,"Max":925},"PowerResetWatts":250}}}]}}]}`

const compStatusNotEnabledReadyOK = `{"Components":[{"ID":"x1002c0s1b0n0","Type":"Node","State":"Ready","Flag":"OK","Enabled":false,"Role":"Compute","NID":1512,"NetType":"Sling","Arch":"X86","Class":"Mountain"}]}`
const x1002c0s1b0n0CompEndpoint = `{"ComponentEndpoints":[{"ID":"x1002c0s1b0n0","Type":"Node","RedfishType":"ComputerSystem","RedfishSubtype":"Physical","OdataID":"/redfish/v1/Systems/Node0","RedfishEndpointID":"x1002c0s0b0","Enabled":true,"RedfishEndpointFQDN":"10.104.8.11","RedfishURL":"10.104.8.11/redfish/v1/Systems/Node0","ComponentEndpointType":"ComponentEndpointComputerSystem","RedfishSystemInfo":{"Name":"Node0","Actions":{"#ComputerSystem.Reset":{"ResetType@Redfish.AllowableValues":["ForceOff","On","Off"],"@Redfish.ActionInfo":"/redfish/v1/Systems/Node0/ResetActionInfo","target":"/redfish/v1/Systems/Node0/Actions/ComputerSystem.Reset"}},"EthernetNICInfo":[{"RedfishId":"HPCNet0","@odata.id":"/redfish/v1/Systems/Node0/EthernetInterfaces/HPCNet0","Description":"Shasta Timms NMC REV04 (HSN)","MACAddress":"Not Available","PermanentMACAddress":" 00:40:a6:83:c5:30"},{"RedfishId":"HPCNet1","@odata.id":"/redfish/v1/Systems/Node0/EthernetInterfaces/HPCNet1","Description":"Shasta Timms NMC REV04 (HSN)","MACAddress":"Not Available","PermanentMACAddress":" 00:40:a6:83:c5:30"},{"RedfishId":"ManagementEthernet","@odata.id":"/redfish/v1/Systems/Node0/EthernetInterfaces/ManagementEthernet","Description":"Node Maintenance Network","MACAddress":"00:40:a6:83:b6:37","PermanentMACAddress":"00:40:a6:83:b6:37"}],"PowerURL":"/redfish/v1/Chassis/Node0/Power","PowerControl":[{"@odata.id":"/redfish/v1/Chassis/Node0/Power#/PowerControl/Node","MemberId":"Node","Name":"Node Power Control","OEM":{"Cray":{"PowerIdleWatts":105,"PowerLimit":{"Min":350,"Max":925},"PowerResetWatts":250}}}]}}]}`

const compStatusNotEnabledOnOK = `{"Components":[{"ID":"x1002c0s1b0n1","Type":"Node","State":"On","Flag":"OK","Enabled":false,"Role":"Compute","NID":1512,"NetType":"Sling","Arch":"X86","Class":"Mountain"}]}`
const x1002c0s1b0n1CompEndpoint = `{"ComponentEndpoints":[{"ID":"x1002c0s1b0n1","Type":"Node","RedfishType":"ComputerSystem","RedfishSubtype":"Physical","OdataID":"/redfish/v1/Systems/Node0","RedfishEndpointID":"x1002c0s0b0","Enabled":true,"RedfishEndpointFQDN":"10.104.8.11","RedfishURL":"10.104.8.11/redfish/v1/Systems/Node0","ComponentEndpointType":"ComponentEndpointComputerSystem","RedfishSystemInfo":{"Name":"Node0","Actions":{"#ComputerSystem.Reset":{"ResetType@Redfish.AllowableValues":["ForceOff","On","Off"],"@Redfish.ActionInfo":"/redfish/v1/Systems/Node0/ResetActionInfo","target":"/redfish/v1/Systems/Node0/Actions/ComputerSystem.Reset"}},"EthernetNICInfo":[{"RedfishId":"HPCNet0","@odata.id":"/redfish/v1/Systems/Node0/EthernetInterfaces/HPCNet0","Description":"Shasta Timms NMC REV04 (HSN)","MACAddress":"Not Available","PermanentMACAddress":" 00:40:a6:83:c5:30"},{"RedfishId":"HPCNet1","@odata.id":"/redfish/v1/Systems/Node0/EthernetInterfaces/HPCNet1","Description":"Shasta Timms NMC REV04 (HSN)","MACAddress":"Not Available","PermanentMACAddress":" 00:40:a6:83:c5:30"},{"RedfishId":"ManagementEthernet","@odata.id":"/redfish/v1/Systems/Node0/EthernetInterfaces/ManagementEthernet","Description":"Node Maintenance Network","MACAddress":"00:40:a6:83:b6:37","PermanentMACAddress":"00:40:a6:83:b6:37"}],"PowerURL":"/redfish/v1/Chassis/Node0/Power","PowerControl":[{"@odata.id":"/redfish/v1/Chassis/Node0/Power#/PowerControl/Node","MemberId":"Node","Name":"Node Power Control","OEM":{"Cray":{"PowerIdleWatts":105,"PowerLimit":{"Min":350,"Max":925},"PowerResetWatts":250}}}]}}]}`

const compStatusNotEnabledOffOK = `{"Components":[{"ID":"x1002c0s1b1n0","Type":"Node","State":"Off","Flag":"OK","Enabled":false,"Role":"Compute","NID":1512,"NetType":"Sling","Arch":"X86","Class":"Mountain"}]}`
const x1002c0s1b1n0CompEndpoint = `{"ComponentEndpoints":[{"ID":"x1002c0s1b1n0","Type":"Node","RedfishType":"ComputerSystem","RedfishSubtype":"Physical","OdataID":"/redfish/v1/Systems/Node0","RedfishEndpointID":"x1002c0s0b1","Enabled":true,"RedfishEndpointFQDN":"10.104.8.12","RedfishURL":"10.104.8.12/redfish/v1/Systems/Node0","ComponentEndpointType":"ComponentEndpointComputerSystem","RedfishSystemInfo":{"Name":"Node0","Actions":{"#ComputerSystem.Reset":{"ResetType@Redfish.AllowableValues":["ForceOff","On","Off"],"@Redfish.ActionInfo":"/redfish/v1/Systems/Node0/ResetActionInfo","target":"/redfish/v1/Systems/Node0/Actions/ComputerSystem.Reset"}},"EthernetNICInfo":[{"RedfishId":"HPCNet0","@odata.id":"/redfish/v1/Systems/Node0/EthernetInterfaces/HPCNet0","Description":"Shasta Timms NMC REV04 (HSN)","MACAddress":"Not Available","PermanentMACAddress":" 00:40:a6:83:c5:30"},{"RedfishId":"HPCNet1","@odata.id":"/redfish/v1/Systems/Node0/EthernetInterfaces/HPCNet1","Description":"Shasta Timms NMC REV04 (HSN)","MACAddress":"Not Available","PermanentMACAddress":" 00:40:a6:83:c5:30"},{"RedfishId":"ManagementEthernet","@odata.id":"/redfish/v1/Systems/Node0/EthernetInterfaces/ManagementEthernet","Description":"Node Maintenance Network","MACAddress":"00:40:a6:83:b6:37","PermanentMACAddress":"00:40:a6:83:b6:37"}],"PowerURL":"/redfish/v1/Chassis/Node0/Power","PowerControl":[{"@odata.id":"/redfish/v1/Chassis/Node0/Power#/PowerControl/Node","MemberId":"Node","Name":"Node Power Control","OEM":{"Cray":{"PowerIdleWatts":105,"PowerLimit":{"Min":350,"Max":925},"PowerResetWatts":250}}}]}}]}`

const compStatusNotEnabledStandbyOK = `{"Components":[{"ID":"x1002c0s1b1n1","Type":"Node","State":"Standby","Flag":"OK","Enabled":false,"Role":"Compute","NID":1512,"NetType":"Sling","Arch":"X86","Class":"Mountain"}]}`
const x1002c0s1b1n1CompEndpoint = `{"ComponentEndpoints":[{"ID":"x1002c0s1b1n1","Type":"Node","RedfishType":"ComputerSystem","RedfishSubtype":"Physical","OdataID":"/redfish/v1/Systems/Node0","RedfishEndpointID":"x1002c0s0b1","Enabled":true,"RedfishEndpointFQDN":"10.104.8.12","RedfishURL":"10.104.8.12/redfish/v1/Systems/Node0","ComponentEndpointType":"ComponentEndpointComputerSystem","RedfishSystemInfo":{"Name":"Node0","Actions":{"#ComputerSystem.Reset":{"ResetType@Redfish.AllowableValues":["ForceOff","On","Off"],"@Redfish.ActionInfo":"/redfish/v1/Systems/Node0/ResetActionInfo","target":"/redfish/v1/Systems/Node0/Actions/ComputerSystem.Reset"}},"EthernetNICInfo":[{"RedfishId":"HPCNet0","@odata.id":"/redfish/v1/Systems/Node0/EthernetInterfaces/HPCNet0","Description":"Shasta Timms NMC REV04 (HSN)","MACAddress":"Not Available","PermanentMACAddress":" 00:40:a6:83:c5:30"},{"RedfishId":"HPCNet1","@odata.id":"/redfish/v1/Systems/Node0/EthernetInterfaces/HPCNet1","Description":"Shasta Timms NMC REV04 (HSN)","MACAddress":"Not Available","PermanentMACAddress":" 00:40:a6:83:c5:30"},{"RedfishId":"ManagementEthernet","@odata.id":"/redfish/v1/Systems/Node0/EthernetInterfaces/ManagementEthernet","Description":"Node Maintenance Network","MACAddress":"00:40:a6:83:b6:37","PermanentMACAddress":"00:40:a6:83:b6:37"}],"PowerURL":"/redfish/v1/Chassis/Node0/Power","PowerControl":[{"@odata.id":"/redfish/v1/Chassis/Node0/Power#/PowerControl/Node","MemberId":"Node","Name":"Node Power Control","OEM":{"Cray":{"PowerIdleWatts":105,"PowerLimit":{"Min":350,"Max":925},"PowerResetWatts":250}}}]}}]}`

const RedfishPowerStateOn = `{"@odata.context":"/redfish/v1/$metadata#ComputerSystem.ComputerSystem","@odata.etag":"W/\"1592257487\"","@odata.id":"/redfish/v1/Systems/Node0","@odata.type":"#ComputerSystem.v1_5_0.ComputerSystem","Actions":{"#ComputerSystem.Reset":{"@Redfish.ActionInfo":"/redfish/v1/Systems/Node0/ResetActionInfo","target":"/redfish/v1/Systems/Node0/Actions/ComputerSystem.Reset"},"#ComputerSystem.SetDefaultBootOrder":{"@Redfish.ActionInfo":"/redfish/v1/Systems/Node0/SetDefaultBootOrderActionInfo","target":"/redfish/v1/Systems/Node0/Actions/ComputerSystem.SetDefaultBootOrder"}},"Bios":{"@odata.id":"/redfish/v1/Systems/Node0/Bios"},"BiosVersion":"wnc.bios-1.1.1-SBIOS-1625-vm4","Boot":{"BootOptions":{"@odata.id":"/redfish/v1/Systems/Node0/BootOptions"},"BootOrder":["Boot0001","Boot0002"]},"Description":"WNC","EthernetInterfaces":{"@odata.id":"/redfish/v1/Systems/Node0/EthernetInterfaces"},"Id":"Node0","Manufacturer":"Cray Inc","Memory":{"@odata.id":"/redfish/v1/Systems/Node0/Memory"},"MemorySummary":{"TotalSystemMemoryGiB":256},"Model":"WNC-Rome","Name":"Node0","PartNumber":"101920703","PowerState":"On","ProcessorSummary":{"Count":2,"Model":"AMD EPYC 7742 64-Core Processor"},"Processors":{"@odata.id":"/redfish/v1/Systems/Node0/Processors"},"SerialNumber":"HS19510077","Status":{"State":"Enabled"},"SystemType":"Physical"}`
const RedfishPowerStateOff = `{"@odata.context":"/redfish/v1/$metadata#ComputerSystem.ComputerSystem","@odata.etag":"W/\"1592257487\"","@odata.id":"/redfish/v1/Systems/Node0","@odata.type":"#ComputerSystem.v1_5_0.ComputerSystem","Actions":{"#ComputerSystem.Reset":{"@Redfish.ActionInfo":"/redfish/v1/Systems/Node0/ResetActionInfo","target":"/redfish/v1/Systems/Node0/Actions/ComputerSystem.Reset"},"#ComputerSystem.SetDefaultBootOrder":{"@Redfish.ActionInfo":"/redfish/v1/Systems/Node0/SetDefaultBootOrderActionInfo","target":"/redfish/v1/Systems/Node0/Actions/ComputerSystem.SetDefaultBootOrder"}},"Bios":{"@odata.id":"/redfish/v1/Systems/Node0/Bios"},"BiosVersion":"wnc.bios-1.1.1-SBIOS-1625-vm4","Boot":{"BootOptions":{"@odata.id":"/redfish/v1/Systems/Node0/BootOptions"},"BootOrder":["Boot0001","Boot0002"]},"Description":"WNC","EthernetInterfaces":{"@odata.id":"/redfish/v1/Systems/Node0/EthernetInterfaces"},"Id":"Node0","Manufacturer":"Cray Inc","Memory":{"@odata.id":"/redfish/v1/Systems/Node0/Memory"},"MemorySummary":{"TotalSystemMemoryGiB":256},"Model":"WNC-Rome","Name":"Node0","PartNumber":"101920703","PowerState":"Off","ProcessorSummary":{"Count":2,"Model":"AMD EPYC 7742 64-Core Processor"},"Processors":{"@odata.id":"/redfish/v1/Systems/Node0/Processors"},"SerialNumber":"HS19510077","Status":{"State":"Enabled"},"SystemType":"Physical"}`
const PCSPowerStatex1002c0s0b0n0 = `{"status":[{"xname":"x1002c0s0b0n0","powerState":"on","managementState":"available","error":"","supportedPowerTransitions":["on","off"],"lastUpdated":"2022-08-24T16:45:53.953811137Z"}]}`
const PCSPowerStatex1002c0s0b0n1 = `{"status":[{"xname":"x1002c0s0b0n1","powerState":"on","managementState":"available","error":"","supportedPowerTransitions":["on","off"],"lastUpdated":"2022-08-24T16:45:53.953811137Z"}]}`
const PCSPowerStatex1002c0s0b1n1 = `{"status":[{"xname":"x1002c0s0b1n1","powerState":"off","managementState":"available","error":"","supportedPowerTransitions":["on","off"],"lastUpdated":"2022-08-24T16:45:53.953811137Z"}]}`
const PCSPowerStatex1002c0s1b0n0 = `{"status":[{"xname":"x1002c0s1b0n0","powerState":"on","managementState":"available","error":"","supportedPowerTransitions":["on","off"],"lastUpdated":"2022-08-24T16:45:53.953811137Z"}]}`
const PCSPowerStatex1002c0s1b0n1 = `{"status":[{"xname":"x1002c0s1b0n1","powerState":"on","managementState":"available","error":"","supportedPowerTransitions":["on","off"],"lastUpdated":"2022-08-24T16:45:53.953811137Z"}]}`
const PCSPowerStatex1002c0s1b1n1 = `{"status":[{"xname":"x1002c0s1b1n1","powerState":"off","managementState":"available","error":"","supportedPowerTransitions":["on","off"],"lastUpdated":"2022-08-24T16:45:53.953811137Z"}]}`

func StatusFunc() RoundTripFunc {
	return func(req *http.Request) (*http.Response, error) {
		//fmt.Printf("StatusFunc %s\n", req.URL.String())
		switch req.URL.String() {
		case "http://localhost:27779/State/Components?enabled=true&state=Off&state=On&state=Ready&type=CabinetPDUPowerConnector&type=CabinetPDUOutlet&type=Chassis&type=RouterModule&type=HSNBoard&type=ComputeModule&type=Node":
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(compStatusEnabledReadyOK)),
				Header:     make(http.Header),
			}, nil
		case "http://localhost:27779/Inventory/ComponentEndpoints?id=x1002c0s0b0n0":
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(x1002c0s0b0n0CompEndpoint)),
				Header:     make(http.Header),
			}, nil
		case "http://localhost:27779/State/Components?enabled=true&state=Off&state=On&type=CabinetPDUPowerConnector&type=CabinetPDUOutlet&type=Chassis&type=RouterModule&type=HSNBoard&type=ComputeModule&type=Node":
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(compStatusEnabledOnOK)),
				Header:     make(http.Header),
			}, nil
		case "http://localhost:27779/Inventory/ComponentEndpoints?id=x1002c0s0b0n1":
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(x1002c0s0b0n1CompEndpoint)),
				Header:     make(http.Header),
			}, nil
		case "http://localhost:27779/Inventory/ComponentEndpoints?id=x1002c0s0b1n0":
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(x1002c0s0b1n0CompEndpoint)),
				Header:     make(http.Header),
			}, nil
		case "http://localhost:27779/State/Components?enabled=true&state=Off&state=On&state=Standby&type=CabinetPDUPowerConnector&type=CabinetPDUOutlet&type=Chassis&type=RouterModule&type=HSNBoard&type=ComputeModule&type=Node":
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(compStatusEnabledStandbyOK)),
				Header:     make(http.Header),
			}, nil
		case "http://localhost:27779/Inventory/ComponentEndpoints?id=x1002c0s0b1n1":
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(x1002c0s0b1n1CompEndpoint)),
				Header:     make(http.Header),
			}, nil
		case "http://localhost:27779/State/Components?enabled=false&state=Off&state=On&state=Ready&type=CabinetPDUPowerConnector&type=CabinetPDUOutlet&type=Chassis&type=RouterModule&type=HSNBoard&type=ComputeModule&type=Node":
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(compStatusNotEnabledReadyOK)),
				Header:     make(http.Header),
			}, nil
		case "http://localhost:27779/Inventory/ComponentEndpoints?id=x1002c0s1b0n0":
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(x1002c0s1b0n0CompEndpoint)),
				Header:     make(http.Header),
			}, nil
		case "http://localhost:27779/State/Components?enabled=false&state=Off&state=On&type=CabinetPDUPowerConnector&type=CabinetPDUOutlet&type=Chassis&type=RouterModule&type=HSNBoard&type=ComputeModule&type=Node":
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(compStatusNotEnabledOnOK)),
				Header:     make(http.Header),
			}, nil
		case "http://localhost:27779/Inventory/ComponentEndpoints?id=x1002c0s1b0n1":
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(x1002c0s1b0n1CompEndpoint)),
				Header:     make(http.Header),
			}, nil
		case "http://localhost:27779/Inventory/ComponentEndpoints?id=x1002c0s1b1n0":
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(x1002c0s1b1n0CompEndpoint)),
				Header:     make(http.Header),
			}, nil
		case "http://localhost:27779/State/Components?enabled=false&state=Off&state=On&state=Standby&type=CabinetPDUPowerConnector&type=CabinetPDUOutlet&type=Chassis&type=RouterModule&type=HSNBoard&type=ComputeModule&type=Node":
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(compStatusNotEnabledStandbyOK)),
				Header:     make(http.Header),
			}, nil
		case "http://localhost:27779/Inventory/ComponentEndpoints?id=x1002c0s1b1n1":
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(x1002c0s1b1n1CompEndpoint)),
				Header:     make(http.Header),
			}, nil
		case "https://10.104.8.11/redfish/v1/Systems/Node0":
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(RedfishPowerStateOn)),
				Header:     make(http.Header),
			}, nil
		case "http://localhost:28007/power-status?xname=x1002c0s0b0n0":
			fmt.Print("GOT IT GOT IT GOT IT x1002c0s0b0n0\n")
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(PCSPowerStatex1002c0s0b0n0)),
				Header:     make(http.Header),
			}, nil
		case "http://localhost:28007/power-status?xname=x1002c0s0b0n1":
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(PCSPowerStatex1002c0s0b0n1)),
				Header:     make(http.Header),
			}, nil
		case "http://localhost:28007/power-status?xname=x1002c0s1b0n0":
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(PCSPowerStatex1002c0s1b0n0)),
				Header:     make(http.Header),
			}, nil
		case "http://localhost:28007/power-status?xname=x1002c0s1b0n1":
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(PCSPowerStatex1002c0s1b0n1)),
				Header:     make(http.Header),
			}, nil
		case "https://10.104.8.12/redfish/v1/Systems/Node0":
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(RedfishPowerStateOff)),
				Header:     make(http.Header),
			}, nil
		case "http://localhost:28007/power-status?xname=x1002c0s0b1n1":
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(PCSPowerStatex1002c0s0b1n1)),
				Header:     make(http.Header),
			}, nil
		case "http://localhost:28007/power-status?xname=x1002c0s1b1n1":
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(PCSPowerStatex1002c0s1b1n1)),
				Header:     make(http.Header),
			}, nil
		default:
			return &http.Response{
				StatusCode: 400,
				Body:       ioutil.NopCloser(bytes.NewBufferString("missing case")),
				Header:     make(http.Header),
			}, nil
		}
	}

}

func TestDoXnameStatus(t *testing.T) {
	var tSvc CapmcD
	var err error
	tSvc.hsmURL, err = url.Parse("http://localhost:27779")
	tSvc.pcsURL, err = url.Parse("http://localhost:28007")
	if err != nil {
		t.Fatal(err)
	}
	testClient := NewTestClient(StatusFunc())
	tSvc.rfClient = testClient
	tSvc.smClient = testClient
	if err != nil {
		t.Fatal(err)
	}
	tSvc.config = loadConfig("")
	ss, adapter := sstorage.NewMockAdapter()
	ccs := compcreds.NewCompCredStore("secret/hms-cred", ss)
	tSvc.ss = ss
	tSvc.ccs = ccs
	tSvc.WPool = base.NewWorkerPool(1, 10)
	tSvc.WPool.Run()

	handler := http.HandlerFunc(tSvc.doXnameStatus)

	tests := []struct {
		name     string
		method   string
		body     io.Reader
		code     int
		expected string
	}{
		{
			"Post empty EOF",
			http.MethodPost,
			bytes.NewBufferString(""),
			http.StatusBadRequest,
			NoRequestJSON + "\n",
		},
		{
			"Enabled Ready OK",
			http.MethodPost,
			bytes.NewBuffer(json.RawMessage(`{"filter":"show_ready"}`)),
			http.StatusOK,
			"{\"e\":0,\"err_msg\":\"\",\"on\":[\"x1002c0s0b0n0\"]}\n",
		},
		{
			"Enabled On OK",
			http.MethodPost,
			bytes.NewBuffer(json.RawMessage(`{"filter":"show_on"}`)),
			http.StatusOK,
			"{\"e\":0,\"err_msg\":\"\",\"on\":[\"x1002c0s0b0n1\"]}\n",
		},
		{
			"Enabled Standby OK",
			http.MethodPost,
			bytes.NewBuffer(json.RawMessage(`{"filter":"show_standby"}`)),
			http.StatusOK,
			"{\"e\":0,\"err_msg\":\"\",\"off\":[\"x1002c0s0b1n1\"]}\n",
		},
		{
			"Not Enabled Ready OK",
			http.MethodPost,
			bytes.NewBuffer(json.RawMessage(`{"filter":"show_ready|show_disabled"}`)),
			http.StatusOK,
			"{\"e\":0,\"err_msg\":\"\",\"on\":[\"x1002c0s1b0n0\"]}\n",
		},
		{
			"Not Enabled On OK",
			http.MethodPost,
			bytes.NewBuffer(json.RawMessage(`{"filter":"show_on|show_disabled"}`)),
			http.StatusOK,
			"{\"e\":0,\"err_msg\":\"\",\"on\":[\"x1002c0s1b0n1\"]}\n",
		},
		{
			"Not Enabled Standby OK",
			http.MethodPost,
			bytes.NewBuffer(json.RawMessage(`{"filter":"show_standby|show_disabled"}`)),
			http.StatusOK,
			"{\"e\":0,\"err_msg\":\"\",\"off\":[\"x1002c0s1b1n1\"]}\n",
		},
	}

	adapter.LookupData = ssDataStatus
	for _, tc := range tests {
		adapter.LookupNum = 0
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, capmc.XnameStatusV1, tc.body)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if rr.Code != tc.code {
				t.Errorf("handler returned wrong status code: want %v but got %v",
					tc.code, rr.Code)
			}

			if rr.Body.String() != tc.expected {
				t.Errorf("handler returned unexpected body: want '%v' but got '%v'",
					tc.expected, rr.Body.String())
			}
		})
	}
}

const mountainHSNBoard = `
{
  "Components": [
    {
      "Arch": "X86",
      "Enabled": true,
      "Flag": "OK",
      "NID": 0,
      "State": "Ready",
      "Role": "Compute",
      "NetType": "Sling",
	  "Type": "Node",
	  "Class": "Mountain",
      "ID": "x1000c0r0e0"
    }
  ]
}
`
const hillHSNBoard = `
{
  "Components": [
    {
      "Arch": "X86",
      "Enabled": true,
      "Flag": "OK",
      "NID": 0,
      "State": "Ready",
      "Role": "Compute",
      "NetType": "Sling",
	  "Type": "Node",
	  "Class": "Hill",
      "ID": "x9000c0r0e0"
    }
  ]
}
`
const riverHSNBoard = `
{
  "Components": [
    {
      "Arch": "X86",
      "Enabled": true,
      "Flag": "OK",
      "NID": 0,
      "State": "Ready",
      "Role": "Compute",
      "NetType": "Sling",
	  "Type": "Node",
	  "Class": "River",
      "ID": "x3000c0r24e0"
    }
  ]
}
`

func callHSMFunc() RoundTripFunc {
	return func(req *http.Request) (*http.Response, error) {
		//fmt.Println(req.URL.String())
		switch req.URL.String() {
		case "http://localhost:27779/State/Components?id=x1000c0r0e0":
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(mountainHSNBoard)),
				Header:     make(http.Header),
			}, nil
		case "http://localhost:27779/State/Components?id=x9000c0r0e0":
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(hillHSNBoard)),
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

func TestHandleDependentComponents(t *testing.T) {
	var tSvc CapmcD
	var err error
	tSvc.hsmURL, err = url.Parse("http://localhost:27779")
	if err != nil {
		t.Fatal(err)
	}
	testClient := NewTestClient(callHSMFunc())
	tSvc.rfClient = testClient
	tSvc.smClient = testClient
	tSvc.config = loadConfig("")

	tests := []struct {
		name   string
		xnames []string
		cmd    string
		want   []string
	}{
		{"MTN Off HSNBoard only", []string{"x1000c0r0e0"}, bmcCmdPowerOff, []string{"x1000c0r0e0"}},
		{"MTN On HSNBoard only", []string{"x1000c0r0e0"}, bmcCmdPowerOn, []string{"x1000c0r0e0"}},
		{"MTN Off RouterModule only", []string{"x1000c0r0"}, bmcCmdPowerOff, []string{"x1000c0r0", "x1000c0r0e0"}},
		{"MTN On RouterModule only", []string{"x1000c0r0"}, bmcCmdPowerOn, []string{"x1000c0r0"}},
		{"MTN Off RouterModule first", []string{"x1000c0r0", "x1000c0r0e0"}, bmcCmdPowerOff, []string{"x1000c0r0", "x1000c0r0e0"}},
		{"MTN On RouterModule first", []string{"x1000c0r0", "x1000c0r0e0"}, bmcCmdPowerOn, []string{"x1000c0r0"}},
		{"MTN Off HSBBoard first", []string{"x1000c0r0e0", "x1000c0r0"}, bmcCmdPowerOff, []string{"x1000c0r0", "x1000c0r0e0"}},
		{"MTN On HSNBoard first", []string{"x1000c0r0e0", "x1000c0r0"}, bmcCmdPowerOn, []string{"x1000c0r0"}},
		{"Hill Off HSNBoard only", []string{"x9000c0r0e0"}, bmcCmdPowerOff, []string{"x9000c0r0e0"}},
		{"Hill On HSNBoard only", []string{"x9000c0r0e0"}, bmcCmdPowerOn, []string{"x9000c0r0e0"}},
		{"Hill Off RouterModule only", []string{"x9000c0r0"}, bmcCmdPowerOff, []string{"x9000c0r0", "x9000c0r0e0"}},
		{"Hill On RouterModule only", []string{"x9000c0r0"}, bmcCmdPowerOn, []string{"x9000c0r0"}},
		{"Hill Off RouterModule first", []string{"x9000c0r0", "x9000c0r0e0"}, bmcCmdPowerOff, []string{"x9000c0r0", "x9000c0r0e0"}},
		{"Hill On RouterModule first", []string{"x9000c0r0", "x9000c0r0e0"}, bmcCmdPowerOn, []string{"x9000c0r0"}},
		{"Hill Off HSBBoard first", []string{"x9000c0r0e0", "x9000c0r0"}, bmcCmdPowerOff, []string{"x9000c0r0", "x9000c0r0e0"}},
		{"Hill On HSNBoard first", []string{"x9000c0r0e0", "x9000c0r0"}, bmcCmdPowerOn, []string{"x9000c0r0"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tSvc.handleDependentComponents(tt.xnames, tt.cmd)
			sort.Strings(got)
			sort.Strings(tt.want)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("handleDependentComponents() = %v, want %v", got, tt.want)
			}
		})
	}
}
