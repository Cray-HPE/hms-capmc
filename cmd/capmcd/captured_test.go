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

// This file is auto-generated via the capture facility

package main

import (
	"testing"

	compcreds "stash.us.cray.com/HMS/hms-compcredentials"
	sstorage "stash.us.cray.com/HMS/hms-securestorage"
)

var getNodeRulesPostHSM = "https://localhost:27779"
var getNodeRulesPostReplayData = []testData{}
var getNodeRulesPostSSData = []sstorage.MockLookup{}

func TestGetNodeRulesPost(t *testing.T) {
	var testReq = testData{"/capmc/get_node_rules",
		"POST",
		`{}`,
		map[string][]string{"Accept": []string{"*/*"}, "Content-Length": []string{"2"}, "Content-Type": []string{"application/x-www-form-urlencoded"}, "User-Agent": []string{"curl/7.47.0"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}},
		`{"e":0,"err_msg":"","latency_node_off":60,"latency_node_on":120,"latency_node_reinit":180,"max_off_req_count":-1,"max_off_time":-1,"max_on_req_count":-1,"max_reinit_req_count":-1,"min_off_time":-1}
`}
	runTest(t, getNodeRulesPostHSM, &testReq, &getNodeRulesPostReplayData, getNodeRulesPostSSData)
}

var getNodeRulesHSM = "https://localhost:27779"
var getNodeRulesReplayData = []testData{}
var getNodeRulesSSData = []sstorage.MockLookup{}

func TestGetNodeRules(t *testing.T) {
	var testReq = testData{"/capmc/get_node_rules",
		"GET",
		``,
		map[string][]string{"User-Agent": []string{"curl/7.47.0"}, "Accept": []string{"*/*"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}},
		`{"e":0,"err_msg":"","latency_node_off":60,"latency_node_on":120,"latency_node_reinit":180,"max_off_req_count":-1,"max_off_time":-1,"max_on_req_count":-1,"max_reinit_req_count":-1,"min_off_time":-1}
`}
	runTest(t, getNodeRulesHSM, &testReq, &getNodeRulesReplayData, getNodeRulesSSData)
}

var nodeOffHSM = "https://localhost:27779"
var nodeOffReplayData = []testData{
	{"https://localhost:27779/hsm/v1/State/Components?nid=16640",
		"GET",
		``,
		map[string][]string{"Content-Type": []string{"application/json"}, "Accept": []string{"application/json"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}, "Date": []string{"Wed, 27 Feb 2019 04:37:13 GMT"}, "Content-Length": []string{"689"}},
		`{"Components":[{"ID":"x0c0s8b0n0","Type":"Node","State":"On","Flag":"OK","Enabled":true,"Role":"Compute","NID":16640,"NetType":"Sling","Arch":"X86"}]}
`},
	{"https://localhost:27779/hsm/v1/Inventory/ComponentEndpoints?id=x0c0s8b0n0",
		"GET",
		``,
		map[string][]string{"Accept": []string{"application/json"}, "Content-Type": []string{"application/json"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}, "Date": []string{"Wed, 27 Feb 2019 04:37:13 GMT"}},
		`{"ComponentEndpoints":[{"ID":"x0c0s8b0n0","Type":"Node","Domain":"ice.next.cray.com","FQDN":"x0c0s8b0n0.ice.next.cray.com","RedfishType":"ComputerSystem","RedfishSubtype":"Physical","MACAddr":"44:a8:42:21:a1:fd","UUID":"4c4c4544-0057-3210-8038-b8c04f463432","OdataID":"/redfish/v1/Systems/System.Embedded.1","RedfishEndpointID":"x0c0s8b0","RedfishEndpointFQDN":"x0c0s8.ice.next.cray.com","RedfishURL":"x0c0s8.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1","ComponentEndpointType":"ComponentEndpointComputerSystem","RedfishSystemInfo":{"Name":"System","Actions":{"#ComputerSystem.Reset":{"ResetType@Redfish.AllowableValues":["On","GracefulShutdown","ForceOff","GracefulRestart","PushPowerButton","Nmi"],"target":"/redfish/v1/Systems/System.Embedded.1/Actions/ComputerSystem.Reset"}},"EthernetNICInfo":[{"RedfishId":"NIC.Integrated.1-3-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-3-1","Description":"Integrated NIC 1 Port 3 Partition 1","MACAddress":"44:a8:42:21:a1:fd","PermanentMACAddress":"44:a8:42:21:a1:fd"},{"RedfishId":"NIC.Integrated.1-4-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-4-1","Description":"Integrated NIC 1 Port 4 Partition 1","MACAddress":"44:a8:42:21:a1:fe","PermanentMACAddress":"44:a8:42:21:a1:fe"},{"RedfishId":"NIC.Integrated.1-1-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-1-1","Description":"Integrated NIC 1 Port 1 Partition 1","MACAddress":"44:a8:42:21:a1:fb","PermanentMACAddress":"44:a8:42:21:a1:fb"},{"RedfishId":"NIC.Integrated.1-2-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-2-1","Description":"Integrated NIC 1 Port 2 Partition 1","MACAddress":"44:a8:42:21:a1:fc","PermanentMACAddress":"44:a8:42:21:a1:fc"}]}}]}
`},
	{"https://x0c0s8.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1/Actions/ComputerSystem.Reset",
		"POST",
		`{"ResetType": "ForceOff"}`,
		map[string][]string{"Content-Type": []string{"application/json"}, "Authorization": []string{"Basic cm9vdDppbml0aWFsMA=="}, "Accept": []string{"*/*"}},
		"204 No Content", 204,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Cache-Control": []string{"no-cache"}, "Connection": []string{"Keep-Alive"}, "Odata-Version": []string{"4.0"}, "Content-Type": []string{"application/json;odata.metadata=minimal;charset=utf-8"}, "Date": []string{"Wed, 27 Feb 2019 10:30:48 GMT"}, "Content-Length": []string{"0"}, "Access-Control-Allow-Origin": []string{"*"}, "Accept-Ranges": []string{"bytes"}, "Odata-Entityid": []string{"/redfish/v1/Systems/System.Embedded.1"}, "Keep-Alive": []string{"timeout=60, max=199"}, "Server": []string{"Appweb/4.5.4"}},
		``},
	{"https://localhost:27779/hsm/v1/locks",
		"POST",
		`{"reason":"Performing a power Off","owner":"CAPMC","lifetime":180,"xnames":["x0c0s8b0n0"]}`,
		map[string][]string{"Content-Type": []string{"application/json"}, "Accept": []string{"application/json"}},
		"201 OK", 201,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}, "Date": []string{"Wed, 27 Feb 2019 04:37:13 GMT"}, "Content-Length": []string{"62"}},
		`{"uri":"/hsm/v1/locks/08c2e624-ded3-11e9-8a34-2a2ae2dbcce4"}]}
`},
}
var nodeOffSSData = []sstorage.MockLookup{
	{
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname:    "x0c0s8b0n0",
				URL:      "x0c0s8b0.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1",
				Username: "root",
				Password: "********",
			},
			Err: nil,
		},
	},
}

func TestNodeOff(t *testing.T) {
	var testReq = testData{"/capmc/node_off",
		"POST",
		`{"nids" : [ 16640]}`,
		map[string][]string{"User-Agent": []string{"curl/7.47.0"}, "Accept": []string{"*/*"}, "Content-Length": []string{"19"}, "Content-Type": []string{"application/x-www-form-urlencoded"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}},
		`{"e":0,"err_msg":""}
`}
	runTest(t, nodeOffHSM, &testReq, &nodeOffReplayData, nodeOffSSData)
}

var nodeOnHSM = "https://localhost:27779"
var nodeOnReplayData = []testData{
	{"https://localhost:27779/hsm/v1/State/Components?nid=16640",
		"GET",
		``,
		map[string][]string{"Accept": []string{"application/json"}, "Content-Type": []string{"application/json"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Date": []string{"Wed, 27 Feb 2019 04:33:18 GMT"}, "Content-Length": []string{"689"}, "Content-Type": []string{"application/json"}},
		`{"Components":[{"ID":"x0c0s8b0n0","Type":"Node","State":"On","Flag":"OK","Enabled":true,"Role":"Compute","NID":16640,"NetType":"Sling","Arch":"X86"}]}
`},
	{"https://localhost:27779/hsm/v1/Inventory/ComponentEndpoints?id=x0c0s8b0n0",
		"GET",
		``,
		map[string][]string{"Accept": []string{"application/json"}, "Content-Type": []string{"application/json"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}, "Date": []string{"Wed, 27 Feb 2019 04:33:18 GMT"}},
		`{"ComponentEndpoints":[{"ID":"x0c0s8b0n0","Type":"Node","Domain":"ice.next.cray.com","FQDN":"x0c0s8b0n0.ice.next.cray.com","RedfishType":"ComputerSystem","RedfishSubtype":"Physical","MACAddr":"44:a8:42:21:a1:fd","UUID":"4c4c4544-0057-3210-8038-b8c04f463432","OdataID":"/redfish/v1/Systems/System.Embedded.1","RedfishEndpointID":"x0c0s8b0","RedfishEndpointFQDN":"x0c0s8.ice.next.cray.com","RedfishURL":"x0c0s8.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1","ComponentEndpointType":"ComponentEndpointComputerSystem","RedfishSystemInfo":{"Name":"System","Actions":{"#ComputerSystem.Reset":{"ResetType@Redfish.AllowableValues":["On","ForceOff","GracefulRestart","PushPowerButton","Nmi"],"target":"/redfish/v1/Systems/System.Embedded.1/Actions/ComputerSystem.Reset"}},"EthernetNICInfo":[{"RedfishId":"NIC.Integrated.1-3-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-3-1","Description":"Integrated NIC 1 Port 3 Partition 1","MACAddress":"44:a8:42:21:a1:fd","PermanentMACAddress":"44:a8:42:21:a1:fd"},{"RedfishId":"NIC.Integrated.1-4-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-4-1","Description":"Integrated NIC 1 Port 4 Partition 1","MACAddress":"44:a8:42:21:a1:fe","PermanentMACAddress":"44:a8:42:21:a1:fe"},{"RedfishId":"NIC.Integrated.1-1-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-1-1","Description":"Integrated NIC 1 Port 1 Partition 1","MACAddress":"44:a8:42:21:a1:fb","PermanentMACAddress":"44:a8:42:21:a1:fb"},{"RedfishId":"NIC.Integrated.1-2-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-2-1","Description":"Integrated NIC 1 Port 2 Partition 1","MACAddress":"44:a8:42:21:a1:fc","PermanentMACAddress":"44:a8:42:21:a1:fc"}]}}]}
`},
	{"https://x0c0s8.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1/Actions/ComputerSystem.Reset",
		"POST",
		`{"ResetType": "On"}`,
		map[string][]string{"Authorization": []string{"Basic cm9vdDppbml0aWFsMA=="}, "Accept": []string{"*/*"}, "Content-Type": []string{"application/json"}},
		"204 No Content", 204,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Keep-Alive": []string{"timeout=60, max=199"}, "Content-Type": []string{"application/json;odata.metadata=minimal;charset=utf-8"}, "Connection": []string{"Keep-Alive"}, "Access-Control-Allow-Origin": []string{"*"}, "Odata-Version": []string{"4.0"}, "Odata-Entityid": []string{"/redfish/v1/Systems/System.Embedded.1"}, "Server": []string{"Appweb/4.5.4"}, "Date": []string{"Wed, 27 Feb 2019 10:26:53 GMT"}, "Cache-Control": []string{"no-cache"}, "Content-Length": []string{"0"}, "Accept-Ranges": []string{"bytes"}},
		``},
	{"https://localhost:27779/hsm/v1/locks",
		"POST",
		`{"reason":"Performing a power Off","owner":"CAPMC","lifetime":180,"xnames":["x0c0s8b0n0"]}`,
		map[string][]string{"Content-Type": []string{"application/json"}, "Accept": []string{"application/json"}},
		"201 OK", 201,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}, "Date": []string{"Wed, 27 Feb 2019 04:37:13 GMT"}, "Content-Length": []string{"62"}},
		`{"uri":"/hsm/v1/locks/08c2e624-ded3-11e9-8a34-2a2ae2dbcce4"}]}
`},
}
var nodeOnSSData = []sstorage.MockLookup{
	{
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname:    "x0c0s8b0n0",
				URL:      "x0c0s8b0.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1",
				Username: "root",
				Password: "********",
			},
			Err: nil,
		},
	},
}

func TestNodeOn(t *testing.T) {
	var testReq = testData{"/capmc/node_on",
		"POST",
		`{"nids" : [ 16640]}`,
		map[string][]string{"User-Agent": []string{"curl/7.47.0"}, "Accept": []string{"*/*"}, "Content-Length": []string{"19"}, "Content-Type": []string{"application/x-www-form-urlencoded"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}},
		`{"e":0,"err_msg":""}
`}
	runTest(t, nodeOnHSM, &testReq, &nodeOnReplayData, nodeOnSSData)
}

var nodeReinitHSM = "https://localhost:27779"
var nodeReinitReplayData = []testData{
	{"https://localhost:27779/hsm/v1/State/Components?nid=16640",
		"GET",
		``,
		map[string][]string{"Accept": []string{"application/json"}, "Content-Type": []string{"application/json"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Date": []string{"Wed, 27 Feb 2019 04:40:30 GMT"}, "Content-Length": []string{"689"}, "Content-Type": []string{"application/json"}},
		`{"Components":[{"ID":"x0c0s8b0n0","Type":"Node","State":"On","Flag":"OK","Enabled":true,"Role":"Compute","NID":16640,"NetType":"Sling","Arch":"X86"}]}
`},
	{"https://localhost:27779/hsm/v1/Inventory/ComponentEndpoints?id=x0c0s8b0n0",
		"GET",
		``,
		map[string][]string{"Accept": []string{"application/json"}, "Content-Type": []string{"application/json"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}, "Date": []string{"Wed, 27 Feb 2019 04:40:30 GMT"}},
		`{"ComponentEndpoints":[{"ID":"x0c0s8b0n0","Type":"Node","Domain":"ice.next.cray.com","FQDN":"x0c0s8b0n0.ice.next.cray.com","RedfishType":"ComputerSystem","RedfishSubtype":"Physical","MACAddr":"44:a8:42:21:a1:fd","UUID":"4c4c4544-0057-3210-8038-b8c04f463432","OdataID":"/redfish/v1/Systems/System.Embedded.1","RedfishEndpointID":"x0c0s8b0","RedfishEndpointFQDN":"x0c0s8.ice.next.cray.com","RedfishURL":"x0c0s8.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1","ComponentEndpointType":"ComponentEndpointComputerSystem","RedfishSystemInfo":{"Name":"System","Actions":{"#ComputerSystem.Reset":{"ResetType@Redfish.AllowableValues":["On","ForceOff","GracefulRestart","PushPowerButton","Nmi"],"target":"/redfish/v1/Systems/System.Embedded.1/Actions/ComputerSystem.Reset"}},"EthernetNICInfo":[{"RedfishId":"NIC.Integrated.1-3-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-3-1","Description":"Integrated NIC 1 Port 3 Partition 1","MACAddress":"44:a8:42:21:a1:fd","PermanentMACAddress":"44:a8:42:21:a1:fd"},{"RedfishId":"NIC.Integrated.1-4-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-4-1","Description":"Integrated NIC 1 Port 4 Partition 1","MACAddress":"44:a8:42:21:a1:fe","PermanentMACAddress":"44:a8:42:21:a1:fe"},{"RedfishId":"NIC.Integrated.1-1-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-1-1","Description":"Integrated NIC 1 Port 1 Partition 1","MACAddress":"44:a8:42:21:a1:fb","PermanentMACAddress":"44:a8:42:21:a1:fb"},{"RedfishId":"NIC.Integrated.1-2-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-2-1","Description":"Integrated NIC 1 Port 2 Partition 1","MACAddress":"44:a8:42:21:a1:fc","PermanentMACAddress":"44:a8:42:21:a1:fc"}]}}]}
`},
	{"https://x0c0s8.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1/Actions/ComputerSystem.Reset",
		"POST",
		`{"ResetType": "GracefulRestart"}`,
		map[string][]string{"Accept": []string{"*/*"}, "Content-Type": []string{"application/json"}, "Authorization": []string{"Basic cm9vdDppbml0aWFsMA=="}},
		"204 No Content", 204,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Odata-Version": []string{"4.0"}, "Odata-Entityid": []string{"/redfish/v1/Systems/System.Embedded.1"}, "Keep-Alive": []string{"timeout=60, max=199"}, "Server": []string{"Appweb/4.5.4"}, "Date": []string{"Wed, 27 Feb 2019 10:34:04 GMT"}, "Connection": []string{"Keep-Alive"}, "Accept-Ranges": []string{"bytes"}, "Content-Type": []string{"application/json;odata.metadata=minimal;charset=utf-8"}, "Cache-Control": []string{"no-cache"}, "Content-Length": []string{"0"}, "Access-Control-Allow-Origin": []string{"*"}},
		``},
	{"https://localhost:27779/hsm/v1/locks",
		"POST",
		`{"reason":"Performing a power Off","owner":"CAPMC","lifetime":180,"xnames":["x0c0s8b0n0"]}`,
		map[string][]string{"Content-Type": []string{"application/json"}, "Accept": []string{"application/json"}},
		"201 OK", 201,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}, "Date": []string{"Wed, 27 Feb 2019 04:37:13 GMT"}, "Content-Length": []string{"62"}},
		`{"uri":"/hsm/v1/locks/08c2e624-ded3-11e9-8a34-2a2ae2dbcce4"}]}
`},
}
var nodeReinitSSData = []sstorage.MockLookup{
	{
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname:    "x0c0s8b0n0",
				URL:      "x0c0s8b0.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1",
				Username: "root",
				Password: "********",
			},
			Err: nil,
		},
	},
}

func TestNodeReinit(t *testing.T) {
	var testReq = testData{"/capmc/node_reinit",
		"POST",
		`{"nids" : [ 16640]}`,
		map[string][]string{"User-Agent": []string{"curl/7.47.0"}, "Accept": []string{"*/*"}, "Content-Length": []string{"19"}, "Content-Type": []string{"application/x-www-form-urlencoded"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}},
		`{"e":0,"err_msg":""}
`}
	runTest(t, nodeReinitHSM, &testReq, &nodeReinitReplayData, nodeReinitSSData)
}

var nodeStatusAllHSM = "https://localhost:27779"
var nodeStatusAllReplayData = []testData{
	{"https://localhost:27779/hsm/v1/State/Components?state=%21Empty&state=%21Populated&state=%21Unknown&type=node",
		"GET",
		``,
		map[string][]string{"Accept": []string{"application/json"}, "Content-Type": []string{"application/json"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}, "Date": []string{"Wed, 27 Feb 2019 04:22:16 GMT"}, "Content-Length": []string{"689"}},
		`{"Components":[{"ID":"x0c0s10b0n0","Type":"Node","State":"On","Flag":"OK","Enabled":true,"Role":"Compute","NID":16704,"NetType":"Sling","Arch":"X86"},{"ID":"x0c0s12b0n0","Type":"Node","State":"On","Flag":"OK","Enabled":true,"Role":"Compute","NID":16768,"NetType":"Sling","Arch":"X86"},{"ID":"x0c0s7b0n0","Type":"Node","State":"On","Flag":"OK","Enabled":true,"Role":"Compute","NID":16608,"NetType":"Sling","Arch":"X86"},{"ID":"x0c0s8b0n0","Type":"Node","State":"Off","Flag":"OK","Enabled":true,"Role":"Compute","NID":16640,"NetType":"Sling","Arch":"X86"},{"ID":"x0c0s9b0n0","Type":"Node","State":"On","Flag":"OK","Enabled":true,"Role":"Compute","NID":16672,"NetType":"Sling","Arch":"X86"}]}
`},
	{"https://localhost:27779/hsm/v1/Inventory/ComponentEndpoints?id=x0c0s10b0n0&id=x0c0s12b0n0&id=x0c0s7b0n0&id=x0c0s8b0n0&id=x0c0s9b0n0",
		"GET",
		``,
		map[string][]string{"Accept": []string{"application/json"}, "Content-Type": []string{"application/json"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Date": []string{"Wed, 27 Feb 2019 04:22:16 GMT"}, "Content-Type": []string{"application/json"}},
		`{"ComponentEndpoints":[{"ID":"x0c0s10b0n0","Type":"Node","Domain":"ice.next.cray.com","FQDN":"x0c0s10b0n0.ice.next.cray.com","RedfishType":"ComputerSystem","RedfishSubtype":"Physical","MACAddr":"14:18:77:50:05:93","UUID":"4c4c4544-004c-3710-8043-b9c04f333832","OdataID":"/redfish/v1/Systems/System.Embedded.1","RedfishEndpointID":"x0c0s10b0","RedfishEndpointFQDN":"x0c0s10.ice.next.cray.com","RedfishURL":"x0c0s10.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1","ComponentEndpointType":"ComponentEndpointComputerSystem","RedfishSystemInfo":{"Name":"System","Actions":{"#ComputerSystem.Reset":{"ResetType@Redfish.AllowableValues":["On","ForceOff","GracefulRestart","PushPowerButton","Nmi"],"target":"/redfish/v1/Systems/System.Embedded.1/Actions/ComputerSystem.Reset"}},"EthernetNICInfo":[{"RedfishId":"NIC.Integrated.1-1-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-1-1","Description":"Integrated NIC 1 Port 1 Partition 1","MACAddress":"14:18:77:50:05:91","PermanentMACAddress":"14:18:77:50:05:91"},{"RedfishId":"NIC.Integrated.1-2-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-2-1","Description":"Integrated NIC 1 Port 2 Partition 1","MACAddress":"14:18:77:50:05:92","PermanentMACAddress":"14:18:77:50:05:92"},{"RedfishId":"NIC.Integrated.1-3-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-3-1","Description":"Integrated NIC 1 Port 3 Partition 1","MACAddress":"14:18:77:50:05:93","PermanentMACAddress":"14:18:77:50:05:93"},{"RedfishId":"NIC.Integrated.1-4-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-4-1","Description":"Integrated NIC 1 Port 4 Partition 1","MACAddress":"14:18:77:50:05:94","PermanentMACAddress":"14:18:77:50:05:94"}]}},{"ID":"x0c0s12b0n0","Type":"Node","Domain":"ice.next.cray.com","FQDN":"x0c0s12b0n0.ice.next.cray.com","RedfishType":"ComputerSystem","RedfishSubtype":"Physical","MACAddr":"24:6e:96:91:62:3a","UUID":"4c4c4544-005a-4c10-8048-b8c04f384d32","OdataID":"/redfish/v1/Systems/System.Embedded.1","RedfishEndpointID":"x0c0s12b0","RedfishEndpointFQDN":"x0c0s12.ice.next.cray.com","RedfishURL":"x0c0s12.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1","ComponentEndpointType":"ComponentEndpointComputerSystem","RedfishSystemInfo":{"Name":"System","Actions":{"#ComputerSystem.Reset":{"ResetType@Redfish.AllowableValues":["On","ForceOff","GracefulRestart","GracefulShutdown","PushPowerButton","Nmi"],"target":"/redfish/v1/Systems/System.Embedded.1/Actions/ComputerSystem.Reset"}},"EthernetNICInfo":[{"RedfishId":"NIC.Integrated.1-3-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-3-1","Description":"Integrated NIC 1 Port 3 Partition 1","MACAddress":"24:6e:96:91:62:3a","PermanentMACAddress":"24:6e:96:91:62:3a"},{"RedfishId":"NIC.Integrated.1-2-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-2-1","Description":"Integrated NIC 1 Port 2 Partition 1","MACAddress":"24:6e:96:91:62:39","PermanentMACAddress":"24:6e:96:91:62:39"},{"RedfishId":"NIC.Integrated.1-1-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-1-1","Description":"Integrated NIC 1 Port 1 Partition 1","MACAddress":"24:6e:96:91:62:38","PermanentMACAddress":"24:6e:96:91:62:38"},{"RedfishId":"NIC.Integrated.1-4-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-4-1","Description":"Integrated NIC 1 Port 4 Partition 1","MACAddress":"24:6e:96:91:62:3b","PermanentMACAddress":"24:6e:96:91:62:3b"}]}},{"ID":"x0c0s7b0n0","Type":"Node","Domain":"ice.next.cray.com","FQDN":"x0c0s7b0n0.ice.next.cray.com","RedfishType":"ComputerSystem","RedfishSubtype":"Physical","MACAddr":"44:a8:42:21:a8:af","UUID":"4c4c4544-0057-3410-8037-b8c04f463432","OdataID":"/redfish/v1/Systems/System.Embedded.1","RedfishEndpointID":"x0c0s7b0","RedfishEndpointFQDN":"x0c0s7.ice.next.cray.com","RedfishURL":"x0c0s7.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1","ComponentEndpointType":"ComponentEndpointComputerSystem","RedfishSystemInfo":{"Name":"System","Actions":{"#ComputerSystem.Reset":{"ResetType@Redfish.AllowableValues":["On","ForceOff","GracefulRestart","PushPowerButton","Nmi"],"target":"/redfish/v1/Systems/System.Embedded.1/Actions/ComputerSystem.Reset"}},"EthernetNICInfo":[{"RedfishId":"NIC.Integrated.1-1-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-1-1","Description":"Integrated NIC 1 Port 1 Partition 1","MACAddress":"44:a8:42:21:a8:ad","PermanentMACAddress":"44:a8:42:21:a8:ad"},{"RedfishId":"NIC.Integrated.1-2-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-2-1","Description":"Integrated NIC 1 Port 2 Partition 1","MACAddress":"44:a8:42:21:a8:ae","PermanentMACAddress":"44:a8:42:21:a8:ae"},{"RedfishId":"NIC.Integrated.1-3-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-3-1","Description":"Integrated NIC 1 Port 3 Partition 1","MACAddress":"44:a8:42:21:a8:af","PermanentMACAddress":"44:a8:42:21:a8:af"},{"RedfishId":"NIC.Integrated.1-4-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-4-1","Description":"Integrated NIC 1 Port 4 Partition 1","MACAddress":"44:a8:42:21:a8:b0","PermanentMACAddress":"44:a8:42:21:a8:b0"}]}},{"ID":"x0c0s8b0n0","Type":"Node","Domain":"ice.next.cray.com","FQDN":"x0c0s8b0n0.ice.next.cray.com","RedfishType":"ComputerSystem","RedfishSubtype":"Physical","MACAddr":"44:a8:42:21:a1:fd","UUID":"4c4c4544-0057-3210-8038-b8c04f463432","OdataID":"/redfish/v1/Systems/System.Embedded.1","RedfishEndpointID":"x0c0s8b0","RedfishEndpointFQDN":"x0c0s8.ice.next.cray.com","RedfishURL":"x0c0s8.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1","ComponentEndpointType":"ComponentEndpointComputerSystem","RedfishSystemInfo":{"Name":"System","Actions":{"#ComputerSystem.Reset":{"ResetType@Redfish.AllowableValues":["On","ForceOff","GracefulRestart","PushPowerButton","Nmi"],"target":"/redfish/v1/Systems/System.Embedded.1/Actions/ComputerSystem.Reset"}},"EthernetNICInfo":[{"RedfishId":"NIC.Integrated.1-3-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-3-1","Description":"Integrated NIC 1 Port 3 Partition 1","MACAddress":"44:a8:42:21:a1:fd","PermanentMACAddress":"44:a8:42:21:a1:fd"},{"RedfishId":"NIC.Integrated.1-4-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-4-1","Description":"Integrated NIC 1 Port 4 Partition 1","MACAddress":"44:a8:42:21:a1:fe","PermanentMACAddress":"44:a8:42:21:a1:fe"},{"RedfishId":"NIC.Integrated.1-1-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-1-1","Description":"Integrated NIC 1 Port 1 Partition 1","MACAddress":"44:a8:42:21:a1:fb","PermanentMACAddress":"44:a8:42:21:a1:fb"},{"RedfishId":"NIC.Integrated.1-2-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-2-1","Description":"Integrated NIC 1 Port 2 Partition 1","MACAddress":"44:a8:42:21:a1:fc","PermanentMACAddress":"44:a8:42:21:a1:fc"}]}},{"ID":"x0c0s9b0n0","Type":"Node","Domain":"ice.next.cray.com","FQDN":"x0c0s9b0n0.ice.next.cray.com","RedfishType":"ComputerSystem","RedfishSubtype":"Physical","MACAddr":"f4:8e:38:c0:b1:0e","UUID":"4c4c4544-005a-5610-804a-b6c04f474232","OdataID":"/redfish/v1/Systems/System.Embedded.1","RedfishEndpointID":"x0c0s9b0","RedfishEndpointFQDN":"x0c0s9.ice.next.cray.com","RedfishURL":"x0c0s9.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1","ComponentEndpointType":"ComponentEndpointComputerSystem","RedfishSystemInfo":{"Name":"System","Actions":{"#ComputerSystem.Reset":{"ResetType@Redfish.AllowableValues":["On","ForceOff","GracefulRestart","PushPowerButton","Nmi"],"target":"/redfish/v1/Systems/System.Embedded.1/Actions/ComputerSystem.Reset"}},"EthernetNICInfo":[{"RedfishId":"NIC.Integrated.1-4-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-4-1","Description":"Integrated NIC 1 Port 4 Partition 1","MACAddress":"f4:8e:38:c0:b1:0f","PermanentMACAddress":"f4:8e:38:c0:b1:0f"},{"RedfishId":"NIC.Integrated.1-1-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-1-1","Description":"Integrated NIC 1 Port 1 Partition 1","MACAddress":"f4:8e:38:c0:b1:0c","PermanentMACAddress":"f4:8e:38:c0:b1:0c"},{"RedfishId":"NIC.Integrated.1-2-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-2-1","Description":"Integrated NIC 1 Port 2 Partition 1","MACAddress":"f4:8e:38:c0:b1:0d","PermanentMACAddress":"f4:8e:38:c0:b1:0d"},{"RedfishId":"NIC.Integrated.1-3-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-3-1","Description":"Integrated NIC 1 Port 3 Partition 1","MACAddress":"f4:8e:38:c0:b1:0e","PermanentMACAddress":"f4:8e:38:c0:b1:0e"}]}}]}
`},
	{"https://x0c0s12.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1",
		"GET",
		``,
		map[string][]string{"Authorization": []string{"Basic cm9vdDppbml0aWFsMA=="}, "Accept": []string{"*/*"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"X-Frame-Options": []string{"DENY"}, "Content-Length": []string{"4036"}, "Content-Type": []string{"application/json;odata.metadata=minimal;charset=utf-8"}, "Date": []string{"Wed, 27 Feb 2019 10:22:58 GMT"}, "Odata-Version": []string{"4.0"}, "Link": []string{"</redfish/v1/Schemas/ComputerSystem.v1_1_0.json>;rel=describedby"}, "Access-Control-Allow-Origin": []string{"*"}, "Cache-Control": []string{"no-cache"}, "Server": []string{"Apache/2.4"}, "Allow": []string{"POST,PATCH"}},
		`{"@odata.context":"/redfish/v1/$metadata#ComputerSystem.ComputerSystem","@odata.id":"/redfish/v1/Systems/System.Embedded.1","@odata.type":"#ComputerSystem.v1_1_0.ComputerSystem","Actions":{"#ComputerSystem.Reset":{"ResetType@Redfish.AllowableValues":["On","ForceOff","GracefulRestart","GracefulShutdown","PushPowerButton","Nmi"],"target":"/redfish/v1/Systems/System.Embedded.1/Actions/ComputerSystem.Reset"}},"AssetTag":"","Bios":{"@odata.id":"/redfish/v1/Systems/System.Embedded.1/Bios"},"BiosVersion":"1.1.7","Boot":{"BootSourceOverrideEnabled":"Once","BootSourceOverrideMode":"UEFI","BootSourceOverrideTarget":"None","BootSourceOverrideTarget@Redfish.AllowableValues":["None","Pxe","Floppy","Cd","Hdd","BiosSetup","Utilities","UefiTarget","SDCard","UefiHttp"],"UefiTargetBootSourceOverride":""},"Description":"Computer System which represents a machine (physical or virtual) and the local resources such as memory, cpu and other devices that can be accessed from that machine.","EthernetInterfaces":{"@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces"},"HostName":"","Id":"System.Embedded.1","IndicatorLED":"Off","Links":{"Chassis":[{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1"}],"Chassis@odata.count":1,"CooledBy":[{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.1A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.1B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.2A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.2B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.3A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.3B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.4A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.4B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.5A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.5B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.6A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.6B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.7A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.7B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.8A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.8B"}],"CooledBy@odata.count":16,"ManagedBy":[{"@odata.id":"/redfish/v1/Managers/iDRAC.Embedded.1"}],"ManagedBy@odata.count":1,"Oem":{"Dell":{"@odata.type":"#DellComputerSystem.v1_0_0.DellComputerSystem","BootOrder":{"@odata.id":"/redfish/v1/Systems/System.Embedded.1/BootSources"}}},"PoweredBy":[{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Power/PowerSupplies/PSU.Slot.1"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Power/PowerSupplies/PSU.Slot.2"}],"PoweredBy@odata.count":2},"Manufacturer":" ","MemorySummary":{"MemoryMirroring":"System","Status":{"Health":"OK","HealthRollup":"OK","State":"Enabled"},"TotalSystemMemoryGiB":64.0},"Model":" ","Name":"System","PartNumber":"0CRT1GA03","PowerState":"On","ProcessorSummary":{"Count":2,"Model":"Intel(R) Xeon(R) Silver 4112 CPU @ 2.60GHz","Status":{"Health":"OK","HealthRollup":"OK","State":"Enabled"}},"Processors":{"@odata.id":"/redfish/v1/Systems/System.Embedded.1/Processors"},"SKU":"8ZLH8M2","SecureBoot":{"@odata.id":"/redfish/v1/Systems/System.Embedded.1/SecureBoot"},"SerialNumber":"CNIVC007B60269","SimpleStorage":{"@odata.id":"/redfish/v1/Systems/System.Embedded.1/Storage/Controllers"},"Status":{"Health":"OK","HealthRollup":"OK","State":"Enabled"},"SystemType":"Physical","TrustedModules":[{"InterfaceType":"TPM1_2","Status":{"State":"Disabled"}}],"UUID":"4c4c4544-005a-4c10-8048-b8c04f384d32"}
`},
	{"https://x0c0s8.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1",
		"GET",
		``,
		map[string][]string{"Authorization": []string{"Basic cm9vdDppbml0aWFsMA=="}, "Accept": []string{"*/*"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Date": []string{"Wed, 27 Feb 2019 10:15:50 GMT"}, "Link": []string{"</redfish/v1/Schemas/ComputerSystem.v1_0_2.json>;rel=describedby"}, "Access-Control-Allow-Origin": []string{"*"}, "Accept-Ranges": []string{"bytes"}, "Odata-Version": []string{"4.0"}, "Keep-Alive": []string{"timeout=60, max=199"}, "Content-Type": []string{"application/json;odata.metadata=minimal;charset=utf-8"}, "Server": []string{"Appweb/4.5.4"}, "Cache-Control": []string{"no-cache"}, "Content-Length": []string{"2154"}, "Allow": []string{"POST,PATCH"}, "Connection": []string{"Keep-Alive"}},
		`{"@odata.context":"/redfish/v1/$metadata#ComputerSystem.ComputerSystem","@odata.id":"/redfish/v1/Systems/System.Embedded.1","@odata.type":"#ComputerSystem.v1_0_2.ComputerSystem","Actions":{"#ComputerSystem.Reset":{"ResetType@Redfish.AllowableValues":["On","ForceOff","GracefulRestart","PushPowerButton","Nmi"],"target":"/redfish/v1/Systems/System.Embedded.1/Actions/ComputerSystem.Reset"}},"AssetTag":"","BiosVersion":"2.1.7","Boot":{"BootSourceOverrideEnabled":"Once","BootSourceOverrideTarget":"None","BootSourceOverrideTarget@Redfish.AllowableValues":["None","Pxe","Cd","Floppy","Hdd","BiosSetup","Utilities","UefiTarget","SDCard"],"UefiTargetBootSourceOverride":""},"Description":"Computer System which represents a machine (physical or virtual) and the local resources such as memory, cpu and other devices that can be accessed from that machine.","EthernetInterfaces":{"@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces"},"HostName":"","Id":"System.Embedded.1","IndicatorLED":"Off","Links":{"Chassis":[{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1"}],"Chassis@odata.count":1,"CooledBy":[],"CooledBy@odata.count":0,"ManagedBy":[{"@odata.id":"/redfish/v1/Managers/iDRAC.Embedded.1"}],"ManagedBy@odata.count":1,"PoweredBy":[{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Power/PowerSupplies/PSU.Slot.1"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Power/PowerSupplies/PSU.Slot.2"}],"PoweredBy@odata.count":2},"Manufacturer":" ","MemorySummary":{"Status":{"Health":null,"HealthRollUp":null,"State":"Enabled"},"TotalSystemMemoryGiB":128.0},"Model":" ","Name":"System","PartNumber":"0CNCJWA05","PowerState":"Off","ProcessorSummary":{"Count":2,"Model":"Intel(R) Xeon(R) CPU E5-2640 v3 @ 2.60GHz","Status":{"Health":null,"HealthRollUp":null,"State":"Enabled"}},"Processors":{"@odata.id":"/redfish/v1/Systems/System.Embedded.1/Processors"},"SKU":"8W28F42","SerialNumber":"CN7475153B0564","SimpleStorage":{"@odata.id":"/redfish/v1/Systems/System.Embedded.1/Storage/Controllers"},"Status":{"Health":"OK","HealthRollUp":"OK","State":"StandbyOffline"},"SystemType":"Physical","UUID":"4c4c4544-0057-3210-8038-b8c04f463432"}
`},
	{"https://x0c0s9.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1",
		"GET",
		``,
		map[string][]string{"Authorization": []string{"Basic cm9vdDppbml0aWFsMA=="}, "Accept": []string{"*/*"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Odata-Version": []string{"4.0"}, "Content-Type": []string{"application/json;odata.metadata=minimal;charset=utf-8"}, "Server": []string{"Appweb/4.5.4"}, "Connection": []string{"Keep-Alive"}, "Keep-Alive": []string{"timeout=60, max=199"}, "Date": []string{"Wed, 27 Feb 2019 05:27:34 GMT"}, "Link": []string{"</redfish/v1/Schemas/ComputerSystem.v1_0_2.json>;rel=describedby"}, "Cache-Control": []string{"no-cache"}, "Content-Length": []string{"3400"}, "Allow": []string{"POST,PATCH"}, "Access-Control-Allow-Origin": []string{"*"}, "Accept-Ranges": []string{"bytes"}},
		`{"@odata.context":"/redfish/v1/$metadata#ComputerSystem.ComputerSystem","@odata.id":"/redfish/v1/Systems/System.Embedded.1","@odata.type":"#ComputerSystem.v1_0_2.ComputerSystem","Actions":{"#ComputerSystem.Reset":{"ResetType@Redfish.AllowableValues":["On","ForceOff","GracefulRestart","PushPowerButton","Nmi"],"target":"/redfish/v1/Systems/System.Embedded.1/Actions/ComputerSystem.Reset"}},"AssetTag":"","BiosVersion":"2.0.1","Boot":{"BootSourceOverrideEnabled":"Once","BootSourceOverrideTarget":"None","BootSourceOverrideTarget@Redfish.AllowableValues":["None","Pxe","Cd","Floppy","Hdd","BiosSetup","Utilities","UefiTarget","SDCard"],"UefiTargetBootSourceOverride":""},"Description":"Computer System which represents a machine (physical or virtual) and the local resources such as memory, cpu and other devices that can be accessed from that machine.","EthernetInterfaces":{"@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces"},"HostName":"MINWINPC","Id":"System.Embedded.1","IndicatorLED":"Off","Links":{"Chassis":[{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1"}],"Chassis@odata.count":1,"CooledBy":[{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.1A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.2A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.3A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.4A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.5A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.6A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.7A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.1B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.2B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.3B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.4B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.5B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.6B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.7B"}],"CooledBy@odata.count":14,"ManagedBy":[{"@odata.id":"/redfish/v1/Managers/iDRAC.Embedded.1"}],"ManagedBy@odata.count":1,"PoweredBy":[{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Power/PowerSupplies/PSU.Slot.1"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Power/PowerSupplies/PSU.Slot.2"}],"PoweredBy@odata.count":2},"Manufacturer":" ","MemorySummary":{"Status":{"Health":"OK","HealthRollUp":"OK","State":"Enabled"},"TotalSystemMemoryGiB":128.0},"Model":" ","Name":"System","PartNumber":"0CNCJWA10","PowerState":"On","ProcessorSummary":{"Count":2,"Model":"Intel(R) Xeon(R) CPU E5-2640 v3 @ 2.60GHz","Status":{"Health":"OK","HealthRollUp":"OK","State":"Enabled"}},"Processors":{"@odata.id":"/redfish/v1/Systems/System.Embedded.1/Processors"},"SKU":"6ZVJGB2","SerialNumber":"CN747515C31036","SimpleStorage":{"@odata.id":"/redfish/v1/Systems/System.Embedded.1/Storage/Controllers"},"Status":{"Health":"OK","HealthRollUp":"OK","State":"Enabled"},"SystemType":"Physical","UUID":"4c4c4544-005a-5610-804a-b6c04f474232"}
`},
	{"https://x0c0s7.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1",
		"GET",
		``,
		map[string][]string{"Authorization": []string{"Basic cm9vdDppbml0aWFsMA=="}, "Accept": []string{"*/*"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Keep-Alive": []string{"timeout=60, max=199"}, "Link": []string{"</redfish/v1/Schemas/ComputerSystem.v1_0_2.json>;rel=describedby"}, "Cache-Control": []string{"no-cache"}, "Content-Length": []string{"3391"}, "Accept-Ranges": []string{"bytes"}, "Connection": []string{"Keep-Alive"}, "Access-Control-Allow-Origin": []string{"*"}, "Odata-Version": []string{"4.0"}, "Content-Type": []string{"application/json;odata.metadata=minimal;charset=utf-8"}, "Server": []string{"Appweb/4.5.4"}, "Date": []string{"Wed, 27 Feb 2019 10:21:44 GMT"}, "Allow": []string{"POST,PATCH"}},
		`{"@odata.context":"/redfish/v1/$metadata#ComputerSystem.ComputerSystem","@odata.id":"/redfish/v1/Systems/System.Embedded.1","@odata.type":"#ComputerSystem.v1_0_2.ComputerSystem","Actions":{"#ComputerSystem.Reset":{"ResetType@Redfish.AllowableValues":["On","ForceOff","GracefulRestart","PushPowerButton","Nmi"],"target":"/redfish/v1/Systems/System.Embedded.1/Actions/ComputerSystem.Reset"}},"AssetTag":"","BiosVersion":"2.4.3","Boot":{"BootSourceOverrideEnabled":"Once","BootSourceOverrideTarget":"None","BootSourceOverrideTarget@Redfish.AllowableValues":["None","Pxe","Cd","Floppy","Hdd","BiosSetup","Utilities","UefiTarget","SDCard"],"UefiTargetBootSourceOverride":""},"Description":"Computer System which represents a machine (physical or virtual) and the local resources such as memory, cpu and other devices that can be accessed from that machine.","EthernetInterfaces":{"@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces"},"HostName":"","Id":"System.Embedded.1","IndicatorLED":"Off","Links":{"Chassis":[{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1"}],"Chassis@odata.count":1,"CooledBy":[{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.1A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.2A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.3A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.4A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.5A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.6A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.7A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.1B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.2B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.3B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.4B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.5B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.6B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.7B"}],"CooledBy@odata.count":14,"ManagedBy":[{"@odata.id":"/redfish/v1/Managers/iDRAC.Embedded.1"}],"ManagedBy@odata.count":1,"PoweredBy":[{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Power/PowerSupplies/PSU.Slot.1"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Power/PowerSupplies/PSU.Slot.2"}],"PoweredBy@odata.count":2},"Manufacturer":" ","MemorySummary":{"Status":{"Health":"OK","HealthRollUp":"OK","State":"Enabled"},"TotalSystemMemoryGiB":64.0},"Model":" ","Name":"System","PartNumber":"0CNCJWA05","PowerState":"On","ProcessorSummary":{"Count":2,"Model":"Intel(R) Xeon(R) CPU E5-2623 v3 @ 3.00GHz","Status":{"Health":"OK","HealthRollUp":"OK","State":"Enabled"}},"Processors":{"@odata.id":"/redfish/v1/Systems/System.Embedded.1/Processors"},"SKU":"8W47F42","SerialNumber":"CN7475153B0496","SimpleStorage":{"@odata.id":"/redfish/v1/Systems/System.Embedded.1/Storage/Controllers"},"Status":{"Health":"OK","HealthRollUp":"OK","State":"Enabled"},"SystemType":"Physical","UUID":"4c4c4544-0057-3410-8037-b8c04f463432"}
`},
	{"https://x0c0s10.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1",
		"GET",
		``,
		map[string][]string{"Authorization": []string{"Basic cm9vdDppbml0aWFsMA=="}, "Accept": []string{"*/*"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Length": []string{"3391"}, "Accept-Ranges": []string{"bytes"}, "Odata-Version": []string{"4.0"}, "Keep-Alive": []string{"timeout=60, max=199"}, "Date": []string{"Wed, 27 Feb 2019 10:21:44 GMT"}, "Link": []string{"</redfish/v1/Schemas/ComputerSystem.v1_0_2.json>;rel=describedby"}, "Connection": []string{"Keep-Alive"}, "Access-Control-Allow-Origin": []string{"*"}, "Content-Type": []string{"application/json;odata.metadata=minimal;charset=utf-8"}, "Server": []string{"Appweb/4.5.4"}, "Cache-Control": []string{"no-cache"}, "Allow": []string{"POST,PATCH"}},
		`{"@odata.context":"/redfish/v1/$metadata#ComputerSystem.ComputerSystem","@odata.id":"/redfish/v1/Systems/System.Embedded.1","@odata.type":"#ComputerSystem.v1_0_2.ComputerSystem","Actions":{"#ComputerSystem.Reset":{"ResetType@Redfish.AllowableValues":["On","ForceOff","GracefulRestart","PushPowerButton","Nmi"],"target":"/redfish/v1/Systems/System.Embedded.1/Actions/ComputerSystem.Reset"}},"AssetTag":"","BiosVersion":"2.4.3","Boot":{"BootSourceOverrideEnabled":"Once","BootSourceOverrideTarget":"None","BootSourceOverrideTarget@Redfish.AllowableValues":["None","Pxe","Cd","Floppy","Hdd","BiosSetup","Utilities","UefiTarget","SDCard"],"UefiTargetBootSourceOverride":""},"Description":"Computer System which represents a machine (physical or virtual) and the local resources such as memory, cpu and other devices that can be accessed from that machine.","EthernetInterfaces":{"@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces"},"HostName":"","Id":"System.Embedded.1","IndicatorLED":"Off","Links":{"Chassis":[{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1"}],"Chassis@odata.count":1,"CooledBy":[{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.1A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.2A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.3A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.4A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.5A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.6A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.7A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.1B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.2B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.3B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.4B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.5B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.6B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.7B"}],"CooledBy@odata.count":14,"ManagedBy":[{"@odata.id":"/redfish/v1/Managers/iDRAC.Embedded.1"}],"ManagedBy@odata.count":1,"PoweredBy":[{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Power/PowerSupplies/PSU.Slot.1"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Power/PowerSupplies/PSU.Slot.2"}],"PoweredBy@odata.count":2},"Manufacturer":" ","MemorySummary":{"Status":{"Health":"OK","HealthRollUp":"OK","State":"Enabled"},"TotalSystemMemoryGiB":64.0},"Model":" ","Name":"System","PartNumber":"0CNCJWA08","PowerState":"On","ProcessorSummary":{"Count":2,"Model":"Intel(R) Xeon(R) CPU E5-2623 v3 @ 3.00GHz","Status":{"Health":"OK","HealthRollUp":"OK","State":"Enabled"}},"Processors":{"@odata.id":"/redfish/v1/Systems/System.Embedded.1/Processors"},"SKU":"9L7C382","SerialNumber":"CN7475159F0098","SimpleStorage":{"@odata.id":"/redfish/v1/Systems/System.Embedded.1/Storage/Controllers"},"Status":{"Health":"OK","HealthRollUp":"OK","State":"Enabled"},"SystemType":"Physical","UUID":"4c4c4544-004c-3710-8043-b9c04f333832"}
`},
}
var nodeStatusAllSSData = []sstorage.MockLookup{
	{
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname:    "x0c0s10b0n0",
				URL:      "x0c0s10b0.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1",
				Username: "root",
				Password: "********",
			},
			Err: nil,
		},
	}, {
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname:    "x0c0s12b0n0",
				URL:      "x0c0s12b0.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1",
				Username: "root",
				Password: "********",
			},
			Err: nil,
		},
	}, {
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname:    "x0c0s7b0n0",
				URL:      "x0c0s7b0.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1",
				Username: "root",
				Password: "********",
			},
			Err: nil,
		},
	}, {
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname:    "x0c0s8b0n0",
				URL:      "x0c0s8b0.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1",
				Username: "root",
				Password: "********",
			},
			Err: nil,
		},
	}, {
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname:    "x0c0s9b0n0",
				URL:      "x0c0s9b0.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1",
				Username: "root",
				Password: "********",
			},
			Err: nil,
		},
	},
}

func TestNodeStatusAll(t *testing.T) {
	var testReq = testData{"/capmc/get_node_status",
		"POST",
		`{}`,
		map[string][]string{"User-Agent": []string{"curl/7.47.0"}, "Accept": []string{"*/*"}, "Content-Length": []string{"2"}, "Content-Type": []string{"application/x-www-form-urlencoded"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}},
		`{"e":0,"err_msg":"","on":[16608,16672,16704,16768],"off":[16640]}
`}
	runTest(t, nodeStatusAllHSM, &testReq, &nodeStatusAllReplayData, nodeStatusAllSSData)
}

var nodeStatusError1HSM = "https://localhost:27779"
var nodeStatusError1ReplayData = []testData{
	{"https://localhost:27779/hsm/v1/State/Components?nid=16608&nid=16704&nid=18800&state=%21Empty&state=%21Populated&state=%21Unknown",
		"GET",
		``,
		map[string][]string{"Accept": []string{"application/json"}, "Content-Type": []string{"application/json"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}, "Date": []string{"Wed, 27 Feb 2019 04:26:44 GMT"}, "Content-Length": []string{"689"}},
		`{"Components":[{"ID":"x0c0s10b0n0","Type":"Node","State":"On","Flag":"OK","Enabled":true,"Role":"Compute","NID":16704,"NetType":"Sling","Arch":"X86"},{"ID":"x0c0s7b0n0","Type":"Node","State":"On","Flag":"OK","Enabled":true,"Role":"Compute","NID":16608,"NetType":"Sling","Arch":"X86"}]}
`},
	{"https://localhost:27779/hsm/v1/Inventory/ComponentEndpoints?id=x0c0s10b0n0&id=x0c0s7b0n0",
		"GET",
		``,
		map[string][]string{"Accept": []string{"application/json"}, "Content-Type": []string{"application/json"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Date": []string{"Wed, 27 Feb 2019 04:26:44 GMT"}, "Content-Type": []string{"application/json"}},
		`{"ComponentEndpoints":[{"ID":"x0c0s10b0n0","Type":"Node","Domain":"ice.next.cray.com","FQDN":"x0c0s10b0n0.ice.next.cray.com","RedfishType":"ComputerSystem","RedfishSubtype":"Physical","MACAddr":"14:18:77:50:05:93","UUID":"4c4c4544-004c-3710-8043-b9c04f333832","OdataID":"/redfish/v1/Systems/System.Embedded.1","RedfishEndpointID":"x0c0s10b0","RedfishEndpointFQDN":"x0c0s10.ice.next.cray.com","RedfishURL":"x0c0s10.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1","ComponentEndpointType":"ComponentEndpointComputerSystem","RedfishSystemInfo":{"Name":"System","Actions":{"#ComputerSystem.Reset":{"ResetType@Redfish.AllowableValues":["On","ForceOff","GracefulRestart","PushPowerButton","Nmi"],"target":"/redfish/v1/Systems/System.Embedded.1/Actions/ComputerSystem.Reset"}},"EthernetNICInfo":[{"RedfishId":"NIC.Integrated.1-1-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-1-1","Description":"Integrated NIC 1 Port 1 Partition 1","MACAddress":"14:18:77:50:05:91","PermanentMACAddress":"14:18:77:50:05:91"},{"RedfishId":"NIC.Integrated.1-2-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-2-1","Description":"Integrated NIC 1 Port 2 Partition 1","MACAddress":"14:18:77:50:05:92","PermanentMACAddress":"14:18:77:50:05:92"},{"RedfishId":"NIC.Integrated.1-3-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-3-1","Description":"Integrated NIC 1 Port 3 Partition 1","MACAddress":"14:18:77:50:05:93","PermanentMACAddress":"14:18:77:50:05:93"},{"RedfishId":"NIC.Integrated.1-4-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-4-1","Description":"Integrated NIC 1 Port 4 Partition 1","MACAddress":"14:18:77:50:05:94","PermanentMACAddress":"14:18:77:50:05:94"}]}},{"ID":"x0c0s7b0n0","Type":"Node","Domain":"ice.next.cray.com","FQDN":"x0c0s7b0n0.ice.next.cray.com","RedfishType":"ComputerSystem","RedfishSubtype":"Physical","MACAddr":"44:a8:42:21:a8:af","UUID":"4c4c4544-0057-3410-8037-b8c04f463432","OdataID":"/redfish/v1/Systems/System.Embedded.1","RedfishEndpointID":"x0c0s7b0","RedfishEndpointFQDN":"x0c0s7.ice.next.cray.com","RedfishURL":"x0c0s7.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1","ComponentEndpointType":"ComponentEndpointComputerSystem","RedfishSystemInfo":{"Name":"System","Actions":{"#ComputerSystem.Reset":{"ResetType@Redfish.AllowableValues":["On","ForceOff","GracefulRestart","PushPowerButton","Nmi"],"target":"/redfish/v1/Systems/System.Embedded.1/Actions/ComputerSystem.Reset"}},"EthernetNICInfo":[{"RedfishId":"NIC.Integrated.1-1-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-1-1","Description":"Integrated NIC 1 Port 1 Partition 1","MACAddress":"44:a8:42:21:a8:ad","PermanentMACAddress":"44:a8:42:21:a8:ad"},{"RedfishId":"NIC.Integrated.1-2-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-2-1","Description":"Integrated NIC 1 Port 2 Partition 1","MACAddress":"44:a8:42:21:a8:ae","PermanentMACAddress":"44:a8:42:21:a8:ae"},{"RedfishId":"NIC.Integrated.1-3-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-3-1","Description":"Integrated NIC 1 Port 3 Partition 1","MACAddress":"44:a8:42:21:a8:af","PermanentMACAddress":"44:a8:42:21:a8:af"},{"RedfishId":"NIC.Integrated.1-4-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-4-1","Description":"Integrated NIC 1 Port 4 Partition 1","MACAddress":"44:a8:42:21:a8:b0","PermanentMACAddress":"44:a8:42:21:a8:b0"}]}}]}
`},
}
var nodeStatusError1SSData = []sstorage.MockLookup{
	{
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname:    "x0c0s10b0n0",
				URL:      "x0c0s10b0.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1",
				Username: "root",
				Password: "********",
			},
			Err: nil,
		},
	}, {
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname:    "x0c0s7b0n0",
				URL:      "x0c0s7b0.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1",
				Username: "root",
				Password: "********",
			},
			Err: nil,
		},
	},
}

func TestNodeStatusError1(t *testing.T) {
	var testReq = testData{"/capmc/get_node_status",
		"POST",
		`{"nids" : [ 16608, 16704, 18800 ]}`,
		map[string][]string{"Content-Length": []string{"34"}, "Content-Type": []string{"application/x-www-form-urlencoded"}, "User-Agent": []string{"curl/7.47.0"}, "Accept": []string{"*/*"}},
		"400 Bad Request", 400,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}},
		`{"e":400,"err_msg":"nids not found: [18800]"}
`}
	runTest(t, nodeStatusError1HSM, &testReq, &nodeStatusError1ReplayData, nodeStatusError1SSData)
}

var nodeStatusError2HSM = "https://localhost:27779"
var nodeStatusError2ReplayData = []testData{}
var nodeStatusError2SSData = []sstorage.MockLookup{}

func TestNodeStatusError2(t *testing.T) {
	var testReq = testData{"/capmc/get_node_status",
		"POST",
		`{"nids" : [ "16608" ]}`,
		map[string][]string{"Content-Type": []string{"application/x-www-form-urlencoded"}, "User-Agent": []string{"curl/7.47.0"}, "Accept": []string{"*/*"}, "Content-Length": []string{"22"}},
		"400 Bad Request", 400,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}},
		`{"e":400,"err_msg":"Bad Request: JSON: json: cannot unmarshal string into Go struct field NodeStatusRequest.nids of type int"}
`}
	runTest(t, nodeStatusError2HSM, &testReq, &nodeStatusError2ReplayData, nodeStatusError2SSData)
}

var nodeStatusSomeHSM = "https://localhost:27779"
var nodeStatusSomeReplayData = []testData{
	{"https://localhost:27779/hsm/v1/State/Components?nid=16608&nid=16704&state=%21Empty&state=%21Populated&state=%21Unknown",
		"GET",
		``,
		map[string][]string{"Accept": []string{"application/json"}, "Content-Type": []string{"application/json"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Length": []string{"689"}, "Content-Type": []string{"application/json"}, "Date": []string{"Wed, 27 Feb 2019 04:25:07 GMT"}},
		`{"Components":[{"ID":"x0c0s10b0n0","Type":"Node","State":"On","Flag":"OK","Enabled":true,"Role":"Compute","NID":16704,"NetType":"Sling","Arch":"X86"},{"ID":"x0c0s7b0n0","Type":"Node","State":"On","Flag":"OK","Enabled":true,"Role":"Compute","NID":16608,"NetType":"Sling","Arch":"X86"}]}
`},
	{"https://localhost:27779/hsm/v1/Inventory/ComponentEndpoints?id=x0c0s10b0n0&id=x0c0s7b0n0",
		"GET",
		``,
		map[string][]string{"Accept": []string{"application/json"}, "Content-Type": []string{"application/json"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}, "Date": []string{"Wed, 27 Feb 2019 04:25:07 GMT"}},
		`{"ComponentEndpoints":[{"ID":"x0c0s10b0n0","Type":"Node","Domain":"ice.next.cray.com","FQDN":"x0c0s10b0n0.ice.next.cray.com","RedfishType":"ComputerSystem","RedfishSubtype":"Physical","MACAddr":"14:18:77:50:05:93","UUID":"4c4c4544-004c-3710-8043-b9c04f333832","OdataID":"/redfish/v1/Systems/System.Embedded.1","RedfishEndpointID":"x0c0s10b0","RedfishEndpointFQDN":"x0c0s10.ice.next.cray.com","RedfishURL":"x0c0s10.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1","ComponentEndpointType":"ComponentEndpointComputerSystem","RedfishSystemInfo":{"Name":"System","Actions":{"#ComputerSystem.Reset":{"ResetType@Redfish.AllowableValues":["On","ForceOff","GracefulRestart","PushPowerButton","Nmi"],"target":"/redfish/v1/Systems/System.Embedded.1/Actions/ComputerSystem.Reset"}},"EthernetNICInfo":[{"RedfishId":"NIC.Integrated.1-1-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-1-1","Description":"Integrated NIC 1 Port 1 Partition 1","MACAddress":"14:18:77:50:05:91","PermanentMACAddress":"14:18:77:50:05:91"},{"RedfishId":"NIC.Integrated.1-2-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-2-1","Description":"Integrated NIC 1 Port 2 Partition 1","MACAddress":"14:18:77:50:05:92","PermanentMACAddress":"14:18:77:50:05:92"},{"RedfishId":"NIC.Integrated.1-3-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-3-1","Description":"Integrated NIC 1 Port 3 Partition 1","MACAddress":"14:18:77:50:05:93","PermanentMACAddress":"14:18:77:50:05:93"},{"RedfishId":"NIC.Integrated.1-4-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-4-1","Description":"Integrated NIC 1 Port 4 Partition 1","MACAddress":"14:18:77:50:05:94","PermanentMACAddress":"14:18:77:50:05:94"}]}},{"ID":"x0c0s7b0n0","Type":"Node","Domain":"ice.next.cray.com","FQDN":"x0c0s7b0n0.ice.next.cray.com","RedfishType":"ComputerSystem","RedfishSubtype":"Physical","MACAddr":"44:a8:42:21:a8:af","UUID":"4c4c4544-0057-3410-8037-b8c04f463432","OdataID":"/redfish/v1/Systems/System.Embedded.1","RedfishEndpointID":"x0c0s7b0","RedfishEndpointFQDN":"x0c0s7.ice.next.cray.com","RedfishURL":"x0c0s7.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1","ComponentEndpointType":"ComponentEndpointComputerSystem","RedfishSystemInfo":{"Name":"System","Actions":{"#ComputerSystem.Reset":{"ResetType@Redfish.AllowableValues":["On","ForceOff","GracefulRestart","PushPowerButton","Nmi"],"target":"/redfish/v1/Systems/System.Embedded.1/Actions/ComputerSystem.Reset"}},"EthernetNICInfo":[{"RedfishId":"NIC.Integrated.1-1-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-1-1","Description":"Integrated NIC 1 Port 1 Partition 1","MACAddress":"44:a8:42:21:a8:ad","PermanentMACAddress":"44:a8:42:21:a8:ad"},{"RedfishId":"NIC.Integrated.1-2-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-2-1","Description":"Integrated NIC 1 Port 2 Partition 1","MACAddress":"44:a8:42:21:a8:ae","PermanentMACAddress":"44:a8:42:21:a8:ae"},{"RedfishId":"NIC.Integrated.1-3-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-3-1","Description":"Integrated NIC 1 Port 3 Partition 1","MACAddress":"44:a8:42:21:a8:af","PermanentMACAddress":"44:a8:42:21:a8:af"},{"RedfishId":"NIC.Integrated.1-4-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-4-1","Description":"Integrated NIC 1 Port 4 Partition 1","MACAddress":"44:a8:42:21:a8:b0","PermanentMACAddress":"44:a8:42:21:a8:b0"}]}}]}
`},
	{"https://x0c0s7.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1",
		"GET",
		``,
		map[string][]string{"Authorization": []string{"Basic cm9vdDppbml0aWFsMA=="}, "Accept": []string{"*/*"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json;odata.metadata=minimal;charset=utf-8"}, "Link": []string{"</redfish/v1/Schemas/ComputerSystem.v1_0_2.json>;rel=describedby"}, "Connection": []string{"Keep-Alive"}, "Accept-Ranges": []string{"bytes"}, "Odata-Version": []string{"4.0"}, "Server": []string{"Appweb/4.5.4"}, "Date": []string{"Wed, 27 Feb 2019 10:24:36 GMT"}, "Cache-Control": []string{"no-cache"}, "Content-Length": []string{"3391"}, "Allow": []string{"POST,PATCH"}, "Access-Control-Allow-Origin": []string{"*"}, "Keep-Alive": []string{"timeout=60, max=199"}},
		`{"@odata.context":"/redfish/v1/$metadata#ComputerSystem.ComputerSystem","@odata.id":"/redfish/v1/Systems/System.Embedded.1","@odata.type":"#ComputerSystem.v1_0_2.ComputerSystem","Actions":{"#ComputerSystem.Reset":{"ResetType@Redfish.AllowableValues":["On","ForceOff","GracefulRestart","PushPowerButton","Nmi"],"target":"/redfish/v1/Systems/System.Embedded.1/Actions/ComputerSystem.Reset"}},"AssetTag":"","BiosVersion":"2.4.3","Boot":{"BootSourceOverrideEnabled":"Once","BootSourceOverrideTarget":"None","BootSourceOverrideTarget@Redfish.AllowableValues":["None","Pxe","Cd","Floppy","Hdd","BiosSetup","Utilities","UefiTarget","SDCard"],"UefiTargetBootSourceOverride":""},"Description":"Computer System which represents a machine (physical or virtual) and the local resources such as memory, cpu and other devices that can be accessed from that machine.","EthernetInterfaces":{"@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces"},"HostName":"","Id":"System.Embedded.1","IndicatorLED":"Off","Links":{"Chassis":[{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1"}],"Chassis@odata.count":1,"CooledBy":[{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.1A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.2A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.3A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.4A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.5A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.6A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.7A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.1B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.2B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.3B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.4B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.5B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.6B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.7B"}],"CooledBy@odata.count":14,"ManagedBy":[{"@odata.id":"/redfish/v1/Managers/iDRAC.Embedded.1"}],"ManagedBy@odata.count":1,"PoweredBy":[{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Power/PowerSupplies/PSU.Slot.1"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Power/PowerSupplies/PSU.Slot.2"}],"PoweredBy@odata.count":2},"Manufacturer":" ","MemorySummary":{"Status":{"Health":"OK","HealthRollUp":"OK","State":"Enabled"},"TotalSystemMemoryGiB":64.0},"Model":" ","Name":"System","PartNumber":"0CNCJWA05","PowerState":"On","ProcessorSummary":{"Count":2,"Model":"Intel(R) Xeon(R) CPU E5-2623 v3 @ 3.00GHz","Status":{"Health":"OK","HealthRollUp":"OK","State":"Enabled"}},"Processors":{"@odata.id":"/redfish/v1/Systems/System.Embedded.1/Processors"},"SKU":"8W47F42","SerialNumber":"CN7475153B0496","SimpleStorage":{"@odata.id":"/redfish/v1/Systems/System.Embedded.1/Storage/Controllers"},"Status":{"Health":"OK","HealthRollUp":"OK","State":"Enabled"},"SystemType":"Physical","UUID":"4c4c4544-0057-3410-8037-b8c04f463432"}
`},
	{"https://x0c0s10.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1",
		"GET",
		``,
		map[string][]string{"Accept": []string{"*/*"}, "Authorization": []string{"Basic cm9vdDppbml0aWFsMA=="}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Link": []string{"</redfish/v1/Schemas/ComputerSystem.v1_0_2.json>;rel=describedby"}, "Cache-Control": []string{"no-cache"}, "Content-Length": []string{"3391"}, "Allow": []string{"POST,PATCH"}, "Odata-Version": []string{"4.0"}, "Keep-Alive": []string{"timeout=60, max=199"}, "Server": []string{"Appweb/4.5.4"}, "Date": []string{"Wed, 27 Feb 2019 10:24:35 GMT"}, "Access-Control-Allow-Origin": []string{"*"}, "Accept-Ranges": []string{"bytes"}, "Content-Type": []string{"application/json;odata.metadata=minimal;charset=utf-8"}, "Connection": []string{"Keep-Alive"}},
		`{"@odata.context":"/redfish/v1/$metadata#ComputerSystem.ComputerSystem","@odata.id":"/redfish/v1/Systems/System.Embedded.1","@odata.type":"#ComputerSystem.v1_0_2.ComputerSystem","Actions":{"#ComputerSystem.Reset":{"ResetType@Redfish.AllowableValues":["On","ForceOff","GracefulRestart","PushPowerButton","Nmi"],"target":"/redfish/v1/Systems/System.Embedded.1/Actions/ComputerSystem.Reset"}},"AssetTag":"","BiosVersion":"2.4.3","Boot":{"BootSourceOverrideEnabled":"Once","BootSourceOverrideTarget":"None","BootSourceOverrideTarget@Redfish.AllowableValues":["None","Pxe","Cd","Floppy","Hdd","BiosSetup","Utilities","UefiTarget","SDCard"],"UefiTargetBootSourceOverride":""},"Description":"Computer System which represents a machine (physical or virtual) and the local resources such as memory, cpu and other devices that can be accessed from that machine.","EthernetInterfaces":{"@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces"},"HostName":"","Id":"System.Embedded.1","IndicatorLED":"Off","Links":{"Chassis":[{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1"}],"Chassis@odata.count":1,"CooledBy":[{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.1A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.2A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.3A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.4A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.5A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.6A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.7A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.1B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.2B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.3B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.4B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.5B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.6B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.7B"}],"CooledBy@odata.count":14,"ManagedBy":[{"@odata.id":"/redfish/v1/Managers/iDRAC.Embedded.1"}],"ManagedBy@odata.count":1,"PoweredBy":[{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Power/PowerSupplies/PSU.Slot.1"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Power/PowerSupplies/PSU.Slot.2"}],"PoweredBy@odata.count":2},"Manufacturer":" ","MemorySummary":{"Status":{"Health":"OK","HealthRollUp":"OK","State":"Enabled"},"TotalSystemMemoryGiB":64.0},"Model":" ","Name":"System","PartNumber":"0CNCJWA08","PowerState":"On","ProcessorSummary":{"Count":2,"Model":"Intel(R) Xeon(R) CPU E5-2623 v3 @ 3.00GHz","Status":{"Health":"OK","HealthRollUp":"OK","State":"Enabled"}},"Processors":{"@odata.id":"/redfish/v1/Systems/System.Embedded.1/Processors"},"SKU":"9L7C382","SerialNumber":"CN7475159F0098","SimpleStorage":{"@odata.id":"/redfish/v1/Systems/System.Embedded.1/Storage/Controllers"},"Status":{"Health":"OK","HealthRollUp":"OK","State":"Enabled"},"SystemType":"Physical","UUID":"4c4c4544-004c-3710-8043-b9c04f333832"}
`},
}
var nodeStatusSomeSSData = []sstorage.MockLookup{
	{
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname:    "x0c0s10b0n0",
				URL:      "x0c0s10b0.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1",
				Username: "root",
				Password: "********",
			},
			Err: nil,
		},
	}, {
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname:    "x0c0s7b0n0",
				URL:      "x0c0s7b0.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1",
				Username: "root",
				Password: "********",
			},
			Err: nil,
		},
	},
}

func TestNodeStatusSome(t *testing.T) {
	var testReq = testData{"/capmc/get_node_status",
		"POST",
		`{"nids" : [ 16608, 16704 ]}`,
		map[string][]string{"User-Agent": []string{"curl/7.47.0"}, "Accept": []string{"*/*"}, "Content-Length": []string{"27"}, "Content-Type": []string{"application/x-www-form-urlencoded"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}},
		`{"e":0,"err_msg":"","on":[16608,16704]}
`}
	runTest(t, nodeStatusSomeHSM, &testReq, &nodeStatusSomeReplayData, nodeStatusSomeSSData)
}

var xnameOffHSM = "https://localhost:27779"
var xnameOffReplayData = []testData{
	{"https://localhost:27779/hsm/v1/State/Components?id=x0c0s8b0n0",
		"GET",
		``,
		map[string][]string{"Accept": []string{"application/json"}, "Content-Type": []string{"application/json"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}, "Date": []string{"Wed, 27 Feb 2019 04:10:31 GMT"}, "Content-Length": []string{"1801"}},
		`{"Components":[{"ID":"x0c0s8b0n0","Type":"Node","State":"On","Flag":"OK","Enabled":true,"Role":"Compute","NID":16640,"NetType":"Sling","Arch":"X86"}]}
`},
	{"https://localhost:27779/hsm/v1/Inventory/ComponentEndpoints?id=x0c0s8b0n0",
		"GET",
		``,
		map[string][]string{"Accept": []string{"application/json"}, "Content-Type": []string{"application/json"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}, "Date": []string{"Wed, 27 Feb 2019 04:10:31 GMT"}},
		`{"ComponentEndpoints":[{"ID":"x0c0s8b0n0","Type":"Node","Domain":"ice.next.cray.com","FQDN":"x0c0s8b0n0.ice.next.cray.com","RedfishType":"ComputerSystem","RedfishSubtype":"Physical","MACAddr":"44:a8:42:21:a1:fd","UUID":"4c4c4544-0057-3210-8038-b8c04f463432","OdataID":"/redfish/v1/Systems/System.Embedded.1","RedfishEndpointID":"x0c0s8b0","RedfishEndpointFQDN":"x0c0s8.ice.next.cray.com","RedfishURL":"x0c0s8.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1","ComponentEndpointType":"ComponentEndpointComputerSystem","RedfishSystemInfo":{"Name":"System","Actions":{"#ComputerSystem.Reset":{"ResetType@Redfish.AllowableValues":["On","GracefulShutdown","ForceOff","GracefulRestart","PushPowerButton","Nmi"],"target":"/redfish/v1/Systems/System.Embedded.1/Actions/ComputerSystem.Reset"}},"EthernetNICInfo":[{"RedfishId":"NIC.Integrated.1-3-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-3-1","Description":"Integrated NIC 1 Port 3 Partition 1","MACAddress":"44:a8:42:21:a1:fd","PermanentMACAddress":"44:a8:42:21:a1:fd"},{"RedfishId":"NIC.Integrated.1-4-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-4-1","Description":"Integrated NIC 1 Port 4 Partition 1","MACAddress":"44:a8:42:21:a1:fe","PermanentMACAddress":"44:a8:42:21:a1:fe"},{"RedfishId":"NIC.Integrated.1-1-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-1-1","Description":"Integrated NIC 1 Port 1 Partition 1","MACAddress":"44:a8:42:21:a1:fb","PermanentMACAddress":"44:a8:42:21:a1:fb"},{"RedfishId":"NIC.Integrated.1-2-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-2-1","Description":"Integrated NIC 1 Port 2 Partition 1","MACAddress":"44:a8:42:21:a1:fc","PermanentMACAddress":"44:a8:42:21:a1:fc"}]}}]}
`},
	{"https://x0c0s8.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1/Actions/ComputerSystem.Reset",
		"POST",
		`{"ResetType": "ForceOff"}`,
		map[string][]string{"Content-Type": []string{"application/json"}, "Authorization": []string{"Basic cm9vdDppbml0aWFsMA=="}, "Accept": []string{"*/*"}},
		"204 No Content", 204,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Odata-Entityid": []string{"/redfish/v1/Systems/System.Embedded.1"}, "Keep-Alive": []string{"timeout=60, max=199"}, "Content-Type": []string{"application/json;odata.metadata=minimal;charset=utf-8"}, "Date": []string{"Wed, 27 Feb 2019 10:04:06 GMT"}, "Access-Control-Allow-Origin": []string{"*"}, "Odata-Version": []string{"4.0"}, "Server": []string{"Appweb/4.5.4"}, "Cache-Control": []string{"no-cache"}, "Content-Length": []string{"0"}, "Connection": []string{"Keep-Alive"}, "Accept-Ranges": []string{"bytes"}},
		``},
	{"https://localhost:27779/hsm/v1/locks",
		"POST",
		`{"reason":"Performing a power Off","owner":"CAPMC","lifetime":180,"xnames":["x0c0s8b0n0"]}`,
		map[string][]string{"Content-Type": []string{"application/json"}, "Accept": []string{"application/json"}},
		"201 OK", 201,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}, "Date": []string{"Wed, 27 Feb 2019 04:37:13 GMT"}, "Content-Length": []string{"62"}},
		`{"uri":"/hsm/v1/locks/08c2e624-ded3-11e9-8a34-2a2ae2dbcce4"}]}
`},
}
var xnameOffSSData = []sstorage.MockLookup{
	{
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname:    "x0c0s8b0n0",
				URL:      "x0c0s8b0.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1",
				Username: "root",
				Password: "********",
			},
			Err: nil,
		},
	},
}

func TestXnameOff(t *testing.T) {
	var testReq = testData{"/capmc/xname_off",
		"POST",
		`{"xnames" : [ "x0c0s8b0n0" ]}`,
		map[string][]string{"User-Agent": []string{"curl/7.47.0"}, "Accept": []string{"*/*"}, "Content-Length": []string{"29"}, "Content-Type": []string{"application/x-www-form-urlencoded"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}},
		`{"e":0,"err_msg":""}
`}
	if testing.Short() {
		t.Skip("skipping test XnameOff in short mode")
	}
	runTest(t, xnameOffHSM, &testReq, &xnameOffReplayData, xnameOffSSData)
}

var xnameOnHSM = "https://localhost:27779"
var xnameOnReplayData = []testData{
	{"https://localhost:27779/hsm/v1/State/Components?id=x0c0s8b0n0",
		"GET",
		``,
		map[string][]string{"Accept": []string{"application/json"}, "Content-Type": []string{"application/json"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}, "Date": []string{"Wed, 27 Feb 2019 04:10:31 GMT"}, "Content-Length": []string{"1801"}},
		`{"Components":[{"ID":"x0c0s8b0n0","Type":"Node","State":"On","Flag":"OK","Enabled":true,"Role":"Compute","NID":16640,"NetType":"Sling","Arch":"X86"}]}
`},
	{"https://localhost:27779/hsm/v1/Inventory/ComponentEndpoints?id=x0c0s8b0n0",
		"GET",
		``,
		map[string][]string{"Accept": []string{"application/json"}, "Content-Type": []string{"application/json"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}, "Date": []string{"Wed, 27 Feb 2019 04:10:31 GMT"}},
		`{"ComponentEndpoints":[{"ID":"x0c0s8b0n0","Type":"Node","Domain":"ice.next.cray.com","FQDN":"x0c0s8b0n0.ice.next.cray.com","RedfishType":"ComputerSystem","RedfishSubtype":"Physical","MACAddr":"44:a8:42:21:a1:fd","UUID":"4c4c4544-0057-3210-8038-b8c04f463432","OdataID":"/redfish/v1/Systems/System.Embedded.1","RedfishEndpointID":"x0c0s8b0","RedfishEndpointFQDN":"x0c0s8.ice.next.cray.com","RedfishURL":"x0c0s8.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1","ComponentEndpointType":"ComponentEndpointComputerSystem","RedfishSystemInfo":{"Name":"System","Actions":{"#ComputerSystem.Reset":{"ResetType@Redfish.AllowableValues":["On","GracefulShutdown","ForceOff","GracefulRestart","PushPowerButton","Nmi"],"target":"/redfish/v1/Systems/System.Embedded.1/Actions/ComputerSystem.Reset"}},"EthernetNICInfo":[{"RedfishId":"NIC.Integrated.1-3-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-3-1","Description":"Integrated NIC 1 Port 3 Partition 1","MACAddress":"44:a8:42:21:a1:fd","PermanentMACAddress":"44:a8:42:21:a1:fd"},{"RedfishId":"NIC.Integrated.1-4-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-4-1","Description":"Integrated NIC 1 Port 4 Partition 1","MACAddress":"44:a8:42:21:a1:fe","PermanentMACAddress":"44:a8:42:21:a1:fe"},{"RedfishId":"NIC.Integrated.1-1-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-1-1","Description":"Integrated NIC 1 Port 1 Partition 1","MACAddress":"44:a8:42:21:a1:fb","PermanentMACAddress":"44:a8:42:21:a1:fb"},{"RedfishId":"NIC.Integrated.1-2-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-2-1","Description":"Integrated NIC 1 Port 2 Partition 1","MACAddress":"44:a8:42:21:a1:fc","PermanentMACAddress":"44:a8:42:21:a1:fc"}]}}]}
`},
	{"https://x0c0s8.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1/Actions/ComputerSystem.Reset",
		"POST",
		`{"ResetType": "On"}`,
		map[string][]string{"Authorization": []string{"Basic cm9vdDppbml0aWFsMA=="}, "Accept": []string{"*/*"}, "Content-Type": []string{"application/json"}},
		"204 No Content", 204,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Connection": []string{"Keep-Alive"}, "Accept-Ranges": []string{"bytes"}, "Odata-Entityid": []string{"/redfish/v1/Systems/System.Embedded.1"}, "Cache-Control": []string{"no-cache"}, "Content-Type": []string{"application/json;odata.metadata=minimal;charset=utf-8"}, "Server": []string{"Appweb/4.5.4"}, "Date": []string{"Wed, 27 Feb 2019 09:48:13 GMT"}, "Content-Length": []string{"0"}, "Access-Control-Allow-Origin": []string{"*"}, "Odata-Version": []string{"4.0"}, "Keep-Alive": []string{"timeout=60, max=199"}},
		``},
	{"https://localhost:27779/hsm/v1/locks",
		"POST",
		`{"reason":"Performing a power Off","owner":"CAPMC","lifetime":180,"xnames":["x0c0s8b0n0"]}`,
		map[string][]string{"Content-Type": []string{"application/json"}, "Accept": []string{"application/json"}},
		"201 OK", 201,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}, "Date": []string{"Wed, 27 Feb 2019 04:37:13 GMT"}, "Content-Length": []string{"62"}},
		`{"uri":"/hsm/v1/locks/08c2e624-ded3-11e9-8a34-2a2ae2dbcce4"}]}
`},
}
var xnameOnSSData = []sstorage.MockLookup{
	{
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname:    "x0c0s8b0n0",
				URL:      "x0c0s8b0.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1",
				Username: "root",
				Password: "********",
			},
			Err: nil,
		},
	},
}

func TestXnameOn(t *testing.T) {
	var testReq = testData{"/capmc/xname_on",
		"POST",
		`{"xnames" : [ "x0c0s8b0n0" ]}`,
		map[string][]string{"User-Agent": []string{"curl/7.47.0"}, "Accept": []string{"*/*"}, "Content-Length": []string{"29"}, "Content-Type": []string{"application/x-www-form-urlencoded"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}},
		`{"e":0,"err_msg":""}
`}
	runTest(t, xnameOnHSM, &testReq, &xnameOnReplayData, xnameOnSSData)
}

var xnameStatusAllHSM = "https://localhost:27779"
var xnameStatusAllReplayData = []testData{
	{"https://localhost:27779/hsm/v1/State/Components?type=CabinetPDUPowerConnector&type=CabinetPDUOutlet&type=Chassis&type=RouterModule&type=HSNBoard&type=ComputeModule&type=Node",
		"GET",
		``,
		map[string][]string{"Accept": []string{"application/json"}, "Content-Type": []string{"application/json"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}, "Date": []string{"Tue, 26 Feb 2019 22:24:37 GMT"}, "Content-Length": []string{"1801"}},
		`{"Components":[{"ID":"x0c0s10b0n0","Type":"Node","State":"On","Flag":"OK","Enabled":true,"Role":"Compute","NID":16704,"NetType":"Sling","Arch":"X86"},{"ID":"x0c0s12b0n0","Type":"Node","State":"On","Flag":"OK","Enabled":true,"Role":"Compute","NID":16768,"NetType":"Sling","Arch":"X86"},{"ID":"x0c0s7b0n0","Type":"Node","State":"On","Flag":"OK","Enabled":true,"Role":"Compute","NID":16608,"NetType":"Sling","Arch":"X86"},{"ID":"x0c0s8b0n0","Type":"Node","State":"On","Flag":"OK","Enabled":true,"Role":"Compute","NID":16640,"NetType":"Sling","Arch":"X86"},{"ID":"x0c0s9b0n0","Type":"Node","State":"On","Flag":"OK","Enabled":true,"Role":"Compute","NID":16672,"NetType":"Sling","Arch":"X86"}]}
`},
	{"https://localhost:27779/hsm/v1/Inventory/ComponentEndpoints?id=x0c0s10b0n0&id=x0c0s12b0n0&id=x0c0s7b0n0&id=x0c0s8b0n0&id=x0c0s9b0n0",
		"GET",
		``,
		map[string][]string{"Accept": []string{"application/json"}, "Content-Type": []string{"application/json"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}, "Date": []string{"Tue, 26 Feb 2019 22:24:37 GMT"}},
		`{"ComponentEndpoints":[{"ID":"x0c0s10b0","Type":"NodeBMC","RedfishType":"Manager","RedfishSubtype":"BMC","MACAddr":"64:00:6a:c2:e5:30","UUID":"3238334f-c0b9-4380-3710-004c4c4c4544","OdataID":"/redfish/v1/Managers/iDRAC.Embedded.1","RedfishEndpointID":"x0c0s10b0","RedfishEndpointFQDN":"x0c0s10.ice.next.cray.com","RedfishURL":"x0c0s10.ice.next.cray.com/redfish/v1/Managers/iDRAC.Embedded.1","ComponentEndpointType":"ComponentEndpointManager","RedfishManagerInfo":{"Name":"Manager","Actions":{"#Manager.Reset":{"ResetType@Redfish.AllowableValues":["GracefulRestart"],"target":"/redfish/v1/Managers/iDRAC.Embedded.1/Actions/Manager.Reset"}},"EthernetNICInfo":[{"RedfishId":"iDRAC.Embedded.1%23NIC.1","@odata.id":"/redfish/v1/Managers/iDRAC.Embedded.1/EthernetInterfaces/iDRAC.Embedded.1%23NIC.1","Description":"Management Network Interface","HostName":"cplus5-drac","InterfaceEnabled":true,"MACAddress":"64:00:6a:c2:e5:30","PermanentMACAddress":"64:00:6a:c2:e5:30"}]}},{"ID":"x0c0s10b0n0","Type":"Node","Domain":"ice.next.cray.com","FQDN":"x0c0s10b0n0.ice.next.cray.com","RedfishType":"ComputerSystem","RedfishSubtype":"Physical","MACAddr":"14:18:77:50:05:93","UUID":"4c4c4544-004c-3710-8043-b9c04f333832","OdataID":"/redfish/v1/Systems/System.Embedded.1","RedfishEndpointID":"x0c0s10b0","RedfishEndpointFQDN":"x0c0s10.ice.next.cray.com","RedfishURL":"x0c0s10.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1","ComponentEndpointType":"ComponentEndpointComputerSystem","RedfishSystemInfo":{"Name":"System","Actions":{"#ComputerSystem.Reset":{"ResetType@Redfish.AllowableValues":["On","ForceOff","GracefulRestart","PushPowerButton","Nmi"],"target":"/redfish/v1/Systems/System.Embedded.1/Actions/ComputerSystem.Reset"}},"EthernetNICInfo":[{"RedfishId":"NIC.Integrated.1-1-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-1-1","Description":"Integrated NIC 1 Port 1 Partition 1","MACAddress":"14:18:77:50:05:91","PermanentMACAddress":"14:18:77:50:05:91"},{"RedfishId":"NIC.Integrated.1-2-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-2-1","Description":"Integrated NIC 1 Port 2 Partition 1","MACAddress":"14:18:77:50:05:92","PermanentMACAddress":"14:18:77:50:05:92"},{"RedfishId":"NIC.Integrated.1-3-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-3-1","Description":"Integrated NIC 1 Port 3 Partition 1","MACAddress":"14:18:77:50:05:93","PermanentMACAddress":"14:18:77:50:05:93"},{"RedfishId":"NIC.Integrated.1-4-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-4-1","Description":"Integrated NIC 1 Port 4 Partition 1","MACAddress":"14:18:77:50:05:94","PermanentMACAddress":"14:18:77:50:05:94"}]}},{"ID":"x0c0s10e0","Type":"NodeEnclosure","RedfishType":"Chassis","RedfishSubtype":"Enclosure","OdataID":"/redfish/v1/Chassis/System.Embedded.1","RedfishEndpointID":"x0c0s10b0","RedfishEndpointFQDN":"x0c0s10.ice.next.cray.com","RedfishURL":"x0c0s10.ice.next.cray.com/redfish/v1/Chassis/System.Embedded.1","ComponentEndpointType":"ComponentEndpointChassis","RedfishChassisInfo":{"Name":"Computer System Chassis","Actions":{"#Chassis.Reset":{"ResetType@Redfish.AllowableValues":["On","ForceOff"],"target":"/redfish/v1/Chassis/System.Embedded.1/Actions/Chassis.Reset"}}}},{"ID":"x0c0s12b0","Type":"NodeBMC","RedfishType":"Manager","RedfishSubtype":"BMC","MACAddr":"58:8a:5a:e6:ec:32","UUID":"324d384f-c0b8-4880-4c10-005a4c4c4544","OdataID":"/redfish/v1/Managers/iDRAC.Embedded.1","RedfishEndpointID":"x0c0s12b0","RedfishEndpointFQDN":"x0c0s12.ice.next.cray.com","RedfishURL":"x0c0s12.ice.next.cray.com/redfish/v1/Managers/iDRAC.Embedded.1","ComponentEndpointType":"ComponentEndpointManager","RedfishManagerInfo":{"Name":"Manager","Actions":{"#Manager.Reset":{"ResetType@Redfish.AllowableValues":["GracefulRestart"],"target":"/redfish/v1/Managers/iDRAC.Embedded.1/Actions/Manager.Reset"}},"EthernetNICInfo":[{"RedfishId":"iDRAC.Embedded.1%23NIC.1","@odata.id":"/redfish/v1/Managers/iDRAC.Embedded.1/EthernetInterfaces/iDRAC.Embedded.1%23NIC.1","Description":"Management Network Interface","HostName":"iDRAC-8ZLH8M2","InterfaceEnabled":true,"MACAddress":"58:8a:5a:e6:ec:32","PermanentMACAddress":"58:8a:5a:e6:ec:32"}]}},{"ID":"x0c0s12b0n0","Type":"Node","Domain":"ice.next.cray.com","FQDN":"x0c0s12b0n0.ice.next.cray.com","RedfishType":"ComputerSystem","RedfishSubtype":"Physical","MACAddr":"24:6e:96:91:62:3a","UUID":"4c4c4544-005a-4c10-8048-b8c04f384d32","OdataID":"/redfish/v1/Systems/System.Embedded.1","RedfishEndpointID":"x0c0s12b0","RedfishEndpointFQDN":"x0c0s12.ice.next.cray.com","RedfishURL":"x0c0s12.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1","ComponentEndpointType":"ComponentEndpointComputerSystem","RedfishSystemInfo":{"Name":"System","Actions":{"#ComputerSystem.Reset":{"ResetType@Redfish.AllowableValues":["On","ForceOff","GracefulRestart","GracefulShutdown","PushPowerButton","Nmi"],"target":"/redfish/v1/Systems/System.Embedded.1/Actions/ComputerSystem.Reset"}},"EthernetNICInfo":[{"RedfishId":"NIC.Integrated.1-3-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-3-1","Description":"Integrated NIC 1 Port 3 Partition 1","MACAddress":"24:6e:96:91:62:3a","PermanentMACAddress":"24:6e:96:91:62:3a"},{"RedfishId":"NIC.Integrated.1-2-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-2-1","Description":"Integrated NIC 1 Port 2 Partition 1","MACAddress":"24:6e:96:91:62:39","PermanentMACAddress":"24:6e:96:91:62:39"},{"RedfishId":"NIC.Integrated.1-1-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-1-1","Description":"Integrated NIC 1 Port 1 Partition 1","MACAddress":"24:6e:96:91:62:38","PermanentMACAddress":"24:6e:96:91:62:38"},{"RedfishId":"NIC.Integrated.1-4-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-4-1","Description":"Integrated NIC 1 Port 4 Partition 1","MACAddress":"24:6e:96:91:62:3b","PermanentMACAddress":"24:6e:96:91:62:3b"}]}},{"ID":"x0c0s12e0","Type":"NodeEnclosure","RedfishType":"Chassis","RedfishSubtype":"Enclosure","OdataID":"/redfish/v1/Chassis/System.Embedded.1","RedfishEndpointID":"x0c0s12b0","RedfishEndpointFQDN":"x0c0s12.ice.next.cray.com","RedfishURL":"x0c0s12.ice.next.cray.com/redfish/v1/Chassis/System.Embedded.1","ComponentEndpointType":"ComponentEndpointChassis","RedfishChassisInfo":{"Name":"Computer System Chassis","Actions":{"#Chassis.Reset":{"ResetType@Redfish.AllowableValues":["On","ForceOff"],"target":"/redfish/v1/Chassis/System.Embedded.1/Actions/Chassis.Reset"}}}},{"ID":"x0c0s7b0","Type":"NodeBMC","RedfishType":"Manager","RedfishSubtype":"BMC","MACAddr":"74:e6:e2:fb:61:14","UUID":"3234464f-c0b8-3780-3410-00574c4c4544","OdataID":"/redfish/v1/Managers/iDRAC.Embedded.1","RedfishEndpointID":"x0c0s7b0","RedfishEndpointFQDN":"x0c0s7.ice.next.cray.com","RedfishURL":"x0c0s7.ice.next.cray.com/redfish/v1/Managers/iDRAC.Embedded.1","ComponentEndpointType":"ComponentEndpointManager","RedfishManagerInfo":{"Name":"Manager","Actions":{"#Manager.Reset":{"ResetType@Redfish.AllowableValues":["GracefulRestart"],"target":"/redfish/v1/Managers/iDRAC.Embedded.1/Actions/Manager.Reset"}},"EthernetNICInfo":[{"RedfishId":"iDRAC.Embedded.1%23NIC.1","@odata.id":"/redfish/v1/Managers/iDRAC.Embedded.1/EthernetInterfaces/iDRAC.Embedded.1%23NIC.1","Description":"Management Network Interface","HostName":"idrac","InterfaceEnabled":true,"MACAddress":"74:e6:e2:fb:61:14","PermanentMACAddress":"74:e6:e2:fb:61:14"}]}},{"ID":"x0c0s7b0n0","Type":"Node","Domain":"ice.next.cray.com","FQDN":"x0c0s7b0n0.ice.next.cray.com","RedfishType":"ComputerSystem","RedfishSubtype":"Physical","MACAddr":"44:a8:42:21:a8:af","UUID":"4c4c4544-0057-3410-8037-b8c04f463432","OdataID":"/redfish/v1/Systems/System.Embedded.1","RedfishEndpointID":"x0c0s7b0","RedfishEndpointFQDN":"x0c0s7.ice.next.cray.com","RedfishURL":"x0c0s7.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1","ComponentEndpointType":"ComponentEndpointComputerSystem","RedfishSystemInfo":{"Name":"System","Actions":{"#ComputerSystem.Reset":{"ResetType@Redfish.AllowableValues":["On","ForceOff","GracefulRestart","PushPowerButton","Nmi"],"target":"/redfish/v1/Systems/System.Embedded.1/Actions/ComputerSystem.Reset"}},"EthernetNICInfo":[{"RedfishId":"NIC.Integrated.1-1-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-1-1","Description":"Integrated NIC 1 Port 1 Partition 1","MACAddress":"44:a8:42:21:a8:ad","PermanentMACAddress":"44:a8:42:21:a8:ad"},{"RedfishId":"NIC.Integrated.1-2-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-2-1","Description":"Integrated NIC 1 Port 2 Partition 1","MACAddress":"44:a8:42:21:a8:ae","PermanentMACAddress":"44:a8:42:21:a8:ae"},{"RedfishId":"NIC.Integrated.1-3-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-3-1","Description":"Integrated NIC 1 Port 3 Partition 1","MACAddress":"44:a8:42:21:a8:af","PermanentMACAddress":"44:a8:42:21:a8:af"},{"RedfishId":"NIC.Integrated.1-4-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-4-1","Description":"Integrated NIC 1 Port 4 Partition 1","MACAddress":"44:a8:42:21:a8:b0","PermanentMACAddress":"44:a8:42:21:a8:b0"}]}},{"ID":"x0c0s7e0","Type":"NodeEnclosure","RedfishType":"Chassis","RedfishSubtype":"Enclosure","OdataID":"/redfish/v1/Chassis/System.Embedded.1","RedfishEndpointID":"x0c0s7b0","RedfishEndpointFQDN":"x0c0s7.ice.next.cray.com","RedfishURL":"x0c0s7.ice.next.cray.com/redfish/v1/Chassis/System.Embedded.1","ComponentEndpointType":"ComponentEndpointChassis","RedfishChassisInfo":{"Name":"Computer System Chassis","Actions":{"#Chassis.Reset":{"ResetType@Redfish.AllowableValues":["On","ForceOff"],"target":"/redfish/v1/Chassis/System.Embedded.1/Actions/Chassis.Reset"}}}},{"ID":"x0c0s8b0","Type":"NodeBMC","RedfishType":"Manager","RedfishSubtype":"BMC","MACAddr":"74:e6:e2:fb:5f:ce","UUID":"3234464f-c0b8-3880-3210-00574c4c4544","OdataID":"/redfish/v1/Managers/iDRAC.Embedded.1","RedfishEndpointID":"x0c0s8b0","RedfishEndpointFQDN":"x0c0s8.ice.next.cray.com","RedfishURL":"x0c0s8.ice.next.cray.com/redfish/v1/Managers/iDRAC.Embedded.1","ComponentEndpointType":"ComponentEndpointManager","RedfishManagerInfo":{"Name":"Manager","Actions":{"#Manager.Reset":{"ResetType@Redfish.AllowableValues":["GracefulRestart"],"target":"/redfish/v1/Managers/iDRAC.Embedded.1/Actions/Manager.Reset"}},"EthernetNICInfo":[{"RedfishId":"iDRAC.Embedded.1%23NIC.1","@odata.id":"/redfish/v1/Managers/iDRAC.Embedded.1/EthernetInterfaces/iDRAC.Embedded.1%23NIC.1","Description":"Management Network Interface","HostName":"crystal2-idrac","InterfaceEnabled":true,"MACAddress":"74:e6:e2:fb:5f:ce","PermanentMACAddress":"74:e6:e2:fb:5f:ce"}]}},{"ID":"x0c0s8b0n0","Type":"Node","Domain":"ice.next.cray.com","FQDN":"x0c0s8b0n0.ice.next.cray.com","RedfishType":"ComputerSystem","RedfishSubtype":"Physical","MACAddr":"44:a8:42:21:a1:fd","UUID":"4c4c4544-0057-3210-8038-b8c04f463432","OdataID":"/redfish/v1/Systems/System.Embedded.1","RedfishEndpointID":"x0c0s8b0","RedfishEndpointFQDN":"x0c0s8.ice.next.cray.com","RedfishURL":"x0c0s8.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1","ComponentEndpointType":"ComponentEndpointComputerSystem","RedfishSystemInfo":{"Name":"System","Actions":{"#ComputerSystem.Reset":{"ResetType@Redfish.AllowableValues":["On","ForceOff","GracefulRestart","PushPowerButton","Nmi"],"target":"/redfish/v1/Systems/System.Embedded.1/Actions/ComputerSystem.Reset"}},"EthernetNICInfo":[{"RedfishId":"NIC.Integrated.1-3-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-3-1","Description":"Integrated NIC 1 Port 3 Partition 1","MACAddress":"44:a8:42:21:a1:fd","PermanentMACAddress":"44:a8:42:21:a1:fd"},{"RedfishId":"NIC.Integrated.1-4-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-4-1","Description":"Integrated NIC 1 Port 4 Partition 1","MACAddress":"44:a8:42:21:a1:fe","PermanentMACAddress":"44:a8:42:21:a1:fe"},{"RedfishId":"NIC.Integrated.1-1-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-1-1","Description":"Integrated NIC 1 Port 1 Partition 1","MACAddress":"44:a8:42:21:a1:fb","PermanentMACAddress":"44:a8:42:21:a1:fb"},{"RedfishId":"NIC.Integrated.1-2-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-2-1","Description":"Integrated NIC 1 Port 2 Partition 1","MACAddress":"44:a8:42:21:a1:fc","PermanentMACAddress":"44:a8:42:21:a1:fc"}]}},{"ID":"x0c0s8e0","Type":"NodeEnclosure","RedfishType":"Chassis","RedfishSubtype":"Enclosure","OdataID":"/redfish/v1/Chassis/System.Embedded.1","RedfishEndpointID":"x0c0s8b0","RedfishEndpointFQDN":"x0c0s8.ice.next.cray.com","RedfishURL":"x0c0s8.ice.next.cray.com/redfish/v1/Chassis/System.Embedded.1","ComponentEndpointType":"ComponentEndpointChassis","RedfishChassisInfo":{"Name":"Computer System Chassis","Actions":{"#Chassis.Reset":{"ResetType@Redfish.AllowableValues":["On","ForceOff"],"target":"/redfish/v1/Chassis/System.Embedded.1/Actions/Chassis.Reset"}}}},{"ID":"x0c0s9b0","Type":"NodeBMC","RedfishType":"Manager","RedfishSubtype":"BMC","MACAddr":"64:00:6a:c4:a6:0a","UUID":"3242474f-c0b6-4a80-5610-005a4c4c4544","OdataID":"/redfish/v1/Managers/iDRAC.Embedded.1","RedfishEndpointID":"x0c0s9b0","RedfishEndpointFQDN":"x0c0s9.ice.next.cray.com","RedfishURL":"x0c0s9.ice.next.cray.com/redfish/v1/Managers/iDRAC.Embedded.1","ComponentEndpointType":"ComponentEndpointManager","RedfishManagerInfo":{"Name":"Manager","Actions":{"#Manager.Reset":{"ResetType@Redfish.AllowableValues":["GracefulRestart"],"target":"/redfish/v1/Managers/iDRAC.Embedded.1/Actions/Manager.Reset"}},"EthernetNICInfo":[{"RedfishId":"iDRAC.Embedded.1%23NIC.1","@odata.id":"/redfish/v1/Managers/iDRAC.Embedded.1/EthernetInterfaces/iDRAC.Embedded.1%23NIC.1","Description":"Management Network Interface","HostName":"idrac-6ZVJGB2","InterfaceEnabled":true,"MACAddress":"64:00:6a:c4:a6:0a","PermanentMACAddress":"64:00:6a:c4:a6:0a"}]}},{"ID":"x0c0s9b0n0","Type":"Node","Domain":"ice.next.cray.com","FQDN":"x0c0s9b0n0.ice.next.cray.com","RedfishType":"ComputerSystem","RedfishSubtype":"Physical","MACAddr":"f4:8e:38:c0:b1:0e","UUID":"4c4c4544-005a-5610-804a-b6c04f474232","OdataID":"/redfish/v1/Systems/System.Embedded.1","RedfishEndpointID":"x0c0s9b0","RedfishEndpointFQDN":"x0c0s9.ice.next.cray.com","RedfishURL":"x0c0s9.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1","ComponentEndpointType":"ComponentEndpointComputerSystem","RedfishSystemInfo":{"Name":"System","Actions":{"#ComputerSystem.Reset":{"ResetType@Redfish.AllowableValues":["On","ForceOff","GracefulRestart","PushPowerButton","Nmi"],"target":"/redfish/v1/Systems/System.Embedded.1/Actions/ComputerSystem.Reset"}},"EthernetNICInfo":[{"RedfishId":"NIC.Integrated.1-4-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-4-1","Description":"Integrated NIC 1 Port 4 Partition 1","MACAddress":"f4:8e:38:c0:b1:0f","PermanentMACAddress":"f4:8e:38:c0:b1:0f"},{"RedfishId":"NIC.Integrated.1-1-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-1-1","Description":"Integrated NIC 1 Port 1 Partition 1","MACAddress":"f4:8e:38:c0:b1:0c","PermanentMACAddress":"f4:8e:38:c0:b1:0c"},{"RedfishId":"NIC.Integrated.1-2-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-2-1","Description":"Integrated NIC 1 Port 2 Partition 1","MACAddress":"f4:8e:38:c0:b1:0d","PermanentMACAddress":"f4:8e:38:c0:b1:0d"},{"RedfishId":"NIC.Integrated.1-3-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-3-1","Description":"Integrated NIC 1 Port 3 Partition 1","MACAddress":"f4:8e:38:c0:b1:0e","PermanentMACAddress":"f4:8e:38:c0:b1:0e"}]}},{"ID":"x0c0s9e0","Type":"NodeEnclosure","RedfishType":"Chassis","RedfishSubtype":"Enclosure","OdataID":"/redfish/v1/Chassis/System.Embedded.1","RedfishEndpointID":"x0c0s9b0","RedfishEndpointFQDN":"x0c0s9.ice.next.cray.com","RedfishURL":"x0c0s9.ice.next.cray.com/redfish/v1/Chassis/System.Embedded.1","ComponentEndpointType":"ComponentEndpointChassis","RedfishChassisInfo":{"Name":"Computer System Chassis","Actions":{"#Chassis.Reset":{"ResetType@Redfish.AllowableValues":["On","ForceOff"],"target":"/redfish/v1/Chassis/System.Embedded.1/Actions/Chassis.Reset"}}}}]}
`},
	{"https://x0c0s12.ice.next.cray.com/redfish/v1/Managers/iDRAC.Embedded.1",
		"GET",
		``,
		map[string][]string{"Authorization": []string{"Basic cm9vdDppbml0aWFsMA=="}, "Accept": []string{"*/*"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Cache-Control": []string{"no-cache"}, "X-Frame-Options": []string{"DENY"}, "Content-Length": []string{"5045"}, "Date": []string{"Wed, 27 Feb 2019 04:25:19 GMT"}, "Server": []string{"Apache/2.4"}, "Link": []string{"</redfish/v1/Schemas/Manager.v1_0_2.json>;rel=describedby"}, "Odata-Version": []string{"4.0"}, "Access-Control-Allow-Origin": []string{"*"}, "Allow": []string{"POST"}, "Content-Type": []string{"application/json;odata.metadata=minimal;charset=utf-8"}},
		`{"@odata.context":"/redfish/v1/$metadata#Manager.Manager","@odata.id":"/redfish/v1/Managers/iDRAC.Embedded.1","@odata.type":"#Manager.v1_0_2.Manager","Actions":{"#Manager.Reset":{"ResetType@Redfish.AllowableValues":["GracefulRestart"],"target":"/redfish/v1/Managers/iDRAC.Embedded.1/Actions/Manager.Reset"},"Oem":{"DellManager.v1_0_0#DellManager.ResetToDefaults":{"ResetType@Redfish.AllowableValues":["All","ResetAllWithRootDefaults","Default"],"target":"/redfish/v1/Managers/iDRAC.Embedded.1/Actions/Oem/DellManager.ResetToDefaults"},"OemManager.v1_0_0#OemManager.ExportSystemConfiguration":{"ExportFormat@Redfish.AllowableValues":["XML","JSON"],"ExportUse@Redfish.AllowableValues":["Default","Clone","Replace"],"IncludeInExport@Redfish.AllowableValues":["Default","IncludeReadOnly","IncludePasswordHashValues","IncludeReadOnly,IncludePasswordHashValues"],"ShareParameters":{"IgnoreCertificateWarning@Redfish.AllowableValues":["Disabled","Enabled"],"ProxySupport@Redfish.AllowableValues":["Disabled","EnabledProxyDefault","Enabled"],"ProxyType@Redfish.AllowableValues":["HTTP","SOCKS4"],"ShareParameters@Redfish.AllowableValues":["IPAddress","ShareName","FileName","UserName","Password","Workgroup","ProxyServer","ProxyUserName","ProxyPassword","ProxyPort"],"ShareType@Redfish.AllowableValues":["NFS","CIFS","HTTP","HTTPS"],"Target@Redfish.AllowableValues":["ALL","IDRAC","BIOS","NIC","RAID"]},"target":"/redfish/v1/Managers/iDRAC.Embedded.1/Actions/Oem/EID_674_Manager.ExportSystemConfiguration"},"OemManager.v1_0_0#OemManager.ImportSystemConfiguration":{"HostPowerState@Redfish.AllowableValues":["On","Off"],"ImportSystemConfiguration@Redfish.AllowableValues":["TimeToWait","ImportBuffer"],"ShareParameters":{"IgnoreCertificateWarning@Redfish.AllowableValues":["Disabled","Enabled"],"ProxySupport@Redfish.AllowableValues":["Disabled","EnabledProxyDefault","Enabled"],"ProxyType@Redfish.AllowableValues":["HTTP","SOCKS4"],"ShareParameters@Redfish.AllowableValues":["IPAddress","ShareName","FileName","UserName","Password","Workgroup","ProxyServer","ProxyUserName","ProxyPassword","ProxyPort"],"ShareType@Redfish.AllowableValues":["NFS","CIFS","HTTP","HTTPS"],"Target@Redfish.AllowableValues":["ALL","IDRAC","BIOS","NIC","RAID"]},"ShutdownType@Redfish.AllowableValues":["Graceful","Forced","NoReboot"],"target":"/redfish/v1/Managers/iDRAC.Embedded.1/Actions/Oem/EID_674_Manager.ImportSystemConfiguration"},"OemManager.v1_0_0#OemManager.ImportSystemConfigurationPreview":{"ImportSystemConfigurationPreview@Redfish.AllowableValues":["ImportBuffer"],"ShareParameters":{"IgnoreCertificateWarning@Redfish.AllowableValues":["Disabled","Enabled"],"ProxySupport@Redfish.AllowableValues":["Disabled","EnabledProxyDefault","Enabled"],"ProxyType@Redfish.AllowableValues":["HTTP","SOCKS4"],"ShareParameters@Redfish.AllowableValues":["IPAddress","ShareName","FileName","UserName","Password","Workgroup","ProxyServer","ProxyUserName","ProxyPassword","ProxyPort"],"ShareType@Redfish.AllowableValues":["NFS","CIFS","HTTP","HTTPS"],"Target@Redfish.AllowableValues":["ALL"]},"target":"/redfish/v1/Managers/iDRAC.Embedded.1/Actions/Oem/EID_674_Manager.ImportSystemConfigurationPreview"}}},"CommandShell":{"ConnectTypesSupported":["SSH","Telnet","IPMI"],"ConnectTypesSupported@odata.count":3,"MaxConcurrentSessions":5,"ServiceEnabled":true},"DateTime":"2019-02-26T22:25:20-06:00","DateTimeLocalOffset":"-06:00","Description":"BMC","EthernetInterfaces":{"@odata.id":"/redfish/v1/Managers/iDRAC.Embedded.1/EthernetInterfaces"},"FirmwareVersion":"3.00.00.00","GraphicalConsole":{"ConnectTypesSupported":["KVMIP"],"ConnectTypesSupported@odata.count":1,"MaxConcurrentSessions":6,"ServiceEnabled":true},"Id":"iDRAC.Embedded.1","Links":{"ManagerForChassis":[{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1"}],"ManagerForChassis@odata.count":1,"ManagerForServers":[{"@odata.id":"/redfish/v1/Systems/System.Embedded.1"}],"ManagerForServers@odata.count":1,"ManagerInChassis":{"@odata.id":"/redfish/v1/Chassis"},"Oem":{"Dell":{"@odata.type":"#DellManager.v1_0_0.DellManager","DellAttributes":[{"@odata.id":"/redfish/v1/Managers/iDRAC.Embedded.1/Attributes"},{"@odata.id":"/redfish/v1/Managers/System.Embedded.1/Attributes"},{"@odata.id":"/redfish/v1/Managers/LifecycleController.Embedded.1/Attributes"}],"DellAttributes@odata.count":3,"Jobs":{"@odata.id":"/redfish/v1/Managers/iDRAC.Embedded.1/Jobs"}}}},"LogServices":{"@odata.id":"/redfish/v1/Managers/iDRAC.Embedded.1/LogServices"},"ManagerType":"BMC","Model":"14G Monolithic","Name":"Manager","NetworkProtocol":{"@odata.id":"/redfish/v1/Managers/iDRAC.Embedded.1/NetworkProtocol"},"Redundancy":[],"Redundancy@odata.count":0,"SerialConsole":{"ConnectTypesSupported":[],"ConnectTypesSupported@odata.count":0,"MaxConcurrentSessions":0,"ServiceEnabled":false},"SerialInterfaces":{"@odata.id":"/redfish/v1/Managers/iDRAC.Embedded.1/SerialInterfaces"},"Status":{"Health":"OK","State":"Enabled"},"UUID":"324d384f-c0b8-4880-4c10-005a4c4c4544","VirtualMedia":{"@odata.id":"/redfish/v1/Managers/iDRAC.Embedded.1/VirtualMedia"}}
`},
	{"https://x0c0s12.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1",
		"GET",
		``,
		map[string][]string{"Authorization": []string{"Basic cm9vdDppbml0aWFsMA=="}, "Accept": []string{"*/*"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Date": []string{"Wed, 27 Feb 2019 04:25:19 GMT"}, "Link": []string{"</redfish/v1/Schemas/ComputerSystem.v1_1_0.json>;rel=describedby"}, "Odata-Version": []string{"4.0"}, "Access-Control-Allow-Origin": []string{"*"}, "X-Frame-Options": []string{"DENY"}, "Content-Length": []string{"4036"}, "Server": []string{"Apache/2.4"}, "Allow": []string{"POST,PATCH"}, "Cache-Control": []string{"no-cache"}, "Content-Type": []string{"application/json;odata.metadata=minimal;charset=utf-8"}},
		`{"@odata.context":"/redfish/v1/$metadata#ComputerSystem.ComputerSystem","@odata.id":"/redfish/v1/Systems/System.Embedded.1","@odata.type":"#ComputerSystem.v1_1_0.ComputerSystem","Actions":{"#ComputerSystem.Reset":{"ResetType@Redfish.AllowableValues":["On","ForceOff","GracefulRestart","GracefulShutdown","PushPowerButton","Nmi"],"target":"/redfish/v1/Systems/System.Embedded.1/Actions/ComputerSystem.Reset"}},"AssetTag":"","Bios":{"@odata.id":"/redfish/v1/Systems/System.Embedded.1/Bios"},"BiosVersion":"1.1.7","Boot":{"BootSourceOverrideEnabled":"Once","BootSourceOverrideMode":"UEFI","BootSourceOverrideTarget":"None","BootSourceOverrideTarget@Redfish.AllowableValues":["None","Pxe","Floppy","Cd","Hdd","BiosSetup","Utilities","UefiTarget","SDCard","UefiHttp"],"UefiTargetBootSourceOverride":""},"Description":"Computer System which represents a machine (physical or virtual) and the local resources such as memory, cpu and other devices that can be accessed from that machine.","EthernetInterfaces":{"@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces"},"HostName":"","Id":"System.Embedded.1","IndicatorLED":"Off","Links":{"Chassis":[{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1"}],"Chassis@odata.count":1,"CooledBy":[{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.1A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.1B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.2A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.2B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.3A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.3B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.4A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.4B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.5A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.5B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.6A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.6B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.7A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.7B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.8A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.8B"}],"CooledBy@odata.count":16,"ManagedBy":[{"@odata.id":"/redfish/v1/Managers/iDRAC.Embedded.1"}],"ManagedBy@odata.count":1,"Oem":{"Dell":{"@odata.type":"#DellComputerSystem.v1_0_0.DellComputerSystem","BootOrder":{"@odata.id":"/redfish/v1/Systems/System.Embedded.1/BootSources"}}},"PoweredBy":[{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Power/PowerSupplies/PSU.Slot.1"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Power/PowerSupplies/PSU.Slot.2"}],"PoweredBy@odata.count":2},"Manufacturer":" ","MemorySummary":{"MemoryMirroring":"System","Status":{"Health":"OK","HealthRollup":"OK","State":"Enabled"},"TotalSystemMemoryGiB":64.0},"Model":" ","Name":"System","PartNumber":"0CRT1GA03","PowerState":"On","ProcessorSummary":{"Count":2,"Model":"Intel(R) Xeon(R) Silver 4112 CPU @ 2.60GHz","Status":{"Health":"OK","HealthRollup":"OK","State":"Enabled"}},"Processors":{"@odata.id":"/redfish/v1/Systems/System.Embedded.1/Processors"},"SKU":"8ZLH8M2","SecureBoot":{"@odata.id":"/redfish/v1/Systems/System.Embedded.1/SecureBoot"},"SerialNumber":"CNIVC007B60269","SimpleStorage":{"@odata.id":"/redfish/v1/Systems/System.Embedded.1/Storage/Controllers"},"Status":{"Health":"OK","HealthRollup":"OK","State":"Enabled"},"SystemType":"Physical","TrustedModules":[{"InterfaceType":"TPM1_2","Status":{"State":"Disabled"}}],"UUID":"4c4c4544-005a-4c10-8048-b8c04f384d32"}
`},
	{"https://x0c0s10.ice.next.cray.com/redfish/v1/Managers/iDRAC.Embedded.1",
		"GET",
		``,
		map[string][]string{"Authorization": []string{"Basic cm9vdDppbml0aWFsMA=="}, "Accept": []string{"*/*"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Cache-Control": []string{"no-cache"}, "Allow": []string{"POST"}, "Access-Control-Allow-Origin": []string{"*"}, "Content-Type": []string{"application/json;odata.metadata=minimal;charset=utf-8"}, "Link": []string{"</redfish/v1/Schemas/Manager.v1_0_2.json>;rel=describedby"}, "Server": []string{"Appweb/4.5.4"}, "Date": []string{"Wed, 27 Feb 2019 04:21:01 GMT"}, "Content-Length": []string{"3521"}, "Connection": []string{"Keep-Alive"}, "Accept-Ranges": []string{"bytes"}, "Odata-Version": []string{"4.0"}, "Keep-Alive": []string{"timeout=60, max=199"}},
		`{"@odata.context":"/redfish/v1/$metadata#Manager.Manager","@odata.id":"/redfish/v1/Managers/iDRAC.Embedded.1","@odata.type":"#Manager.v1_0_2.Manager","Actions":{"#Manager.Reset":{"ResetType@Redfish.AllowableValues":["GracefulRestart"],"target":"/redfish/v1/Managers/iDRAC.Embedded.1/Actions/Manager.Reset"},"Oem":{"OemManager.v1_0_0#OemManager.ExportSystemConfiguration":{"ExportFormat@Redfish.AllowableValues":["XML"],"ExportUse@Redfish.AllowableValues":["Default","Clone","Replace"],"IncludeInExport@Redfish.AllowableValues":["Default","IncludeReadOnly","IncludePasswordHashValues"],"ShareParameters":{"ShareParameters@Redfish.AllowableValues":["IPAddress","ShareName","FileName","UserName","Password","Workgroup"],"ShareType@Redfish.AllowableValues":["NFS","CIFS"],"Target@Redfish.AllowableValues":["ALL","IDRAC","BIOS","NIC","RAID"]},"target":"/redfish/v1/Managers/iDRAC.Embedded.1/Actions/Oem/EID_674_Manager.ExportSystemConfiguration"},"OemManager.v1_0_0#OemManager.ImportSystemConfiguration":{"HostPowerState@Redfish.AllowableValues":["On","Off"],"ImportSystemConfiguration@Redfish.AllowableValues":["TimeToWait","ImportBuffer"],"ShareParameters":{"ShareParameters@Redfish.AllowableValues":["IPAddress","ShareName","FileName","UserName","Password","Workgroup"],"ShareType@Redfish.AllowableValues":["NFS","CIFS"],"Target@Redfish.AllowableValues":["ALL","IDRAC","BIOS","NIC","RAID"]},"ShutdownType@Redfish.AllowableValues":["Graceful","Forced","NoReboot"],"target":"/redfish/v1/Managers/iDRAC.Embedded.1/Actions/Oem/EID_674_Manager.ImportSystemConfiguration"},"OemManager.v1_0_0#OemManager.ImportSystemConfigurationPreview":{"ImportSystemConfigurationPreview@Redfish.AllowableValues":["ImportBuffer"],"ShareParameters":{"ShareParameters@Redfish.AllowableValues":["IPAddress","ShareName","FileName","UserName","Password","Workgroup"],"ShareType@Redfish.AllowableValues":["NFS","CIFS"],"Target@Redfish.AllowableValues":["ALL"]},"target":"/redfish/v1/Managers/iDRAC.Embedded.1/Actions/Oem/EID_674_Manager.ImportSystemConfigurationPreview"}}},"CommandShell":{"ConnectTypesSupported":["SSH","Telnet","IPMI"],"ConnectTypesSupported@odata.count":3,"MaxConcurrentSessions":5,"ServiceEnabled":true},"DateTime":"2019-02-26T22:21:01-06:00","DateTimeLocalOffset":"-06:00","Description":"BMC","EthernetInterfaces":{"@odata.id":"/redfish/v1/Managers/iDRAC.Embedded.1/EthernetInterfaces"},"FirmwareVersion":"2.41.40.40","GraphicalConsole":{"ConnectTypesSupported":["KVMIP"],"ConnectTypesSupported@odata.count":1,"MaxConcurrentSessions":6,"ServiceEnabled":true},"Id":"iDRAC.Embedded.1","Links":{"ManagerForChassis":[{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1"}],"ManagerForChassis@odata.count":1,"ManagerForServers":[{"@odata.id":"/redfish/v1/Systems/System.Embedded.1"}],"ManagerForServers@odata.count":1},"LogServices":{"@odata.id":"/redfish/v1/Managers/iDRAC.Embedded.1/LogServices"},"ManagerType":"BMC","Model":"13G Monolithic","Name":"Manager","NetworkProtocol":{"@odata.id":"/redfish/v1/Managers/iDRAC.Embedded.1/NetworkProtocol"},"Redundancy":[],"Redundancy@odata.count":0,"RedundancySet":[],"RedundancySet@odata.count":0,"SerialConsole":{"ConnectTypesSupported":[],"ConnectTypesSupported@odata.count":0,"MaxConcurrentSessions":0,"ServiceEnabled":false},"SerialInterfaces":{"@odata.id":"/redfish/v1/Managers/iDRAC.Embedded.1/SerialInterfaces"},"Status":{"Health":"Ok","State":"Enabled"},"UUID":"3238334f-c0b9-4380-3710-004c4c4c4544","VirtualMedia":{"@odata.id":"/redfish/v1/Managers/iDRAC.Embedded.1/VirtualMedia"}}
`},
	{"https://x0c0s12.ice.next.cray.com/redfish/v1/Chassis/System.Embedded.1",
		"GET",
		``,
		map[string][]string{"Authorization": []string{"Basic cm9vdDppbml0aWFsMA=="}, "Accept": []string{"*/*"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Date": []string{"Wed, 27 Feb 2019 04:25:20 GMT"}, "Server": []string{"Apache/2.4"}, "Link": []string{"</redfish/v1/Schemas/Chassis.v1_0_2.json>;rel=describedby"}, "Odata-Version": []string{"4.0"}, "X-Frame-Options": []string{"DENY"}, "Allow": []string{"POST,PATCH"}, "Access-Control-Allow-Origin": []string{"*"}, "Cache-Control": []string{"no-cache"}, "Content-Length": []string{"3249"}, "Content-Type": []string{"application/json;odata.metadata=minimal;charset=utf-8"}},
		`{"@odata.context":"/redfish/v1/$metadata#Chassis.Chassis","@odata.id":"/redfish/v1/Chassis/System.Embedded.1","@odata.type":"#Chassis.v1_2_0.Chassis","Actions":{"#Chassis.Reset":{"ResetType@Redfish.AllowableValues":["On","ForceOff"],"target":"/redfish/v1/Chassis/System.Embedded.1/Actions/Chassis.Reset"}},"AssetTag":null,"ChassisType":"Enclosure","Description":"It represents the properties for physical components for any system.It represent racks, rackmount servers, blades, standalone, modular systems,enclosures, and all other containers.The non-cpu/device centric parts of the schema are all accessed either directly or indirectly through this resource.","Id":"System.Embedded.1","IndicatorLED":"Off","Links":{"ComputerSystems":[{"@odata.id":"/redfish/v1/Systems/System.Embedded.1"}],"ComputerSystems@odata.count":1,"Contains":[],"Contains@odata.count":0,"CooledBy":[{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.1A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.1B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.2A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.2B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.3A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.3B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.4A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.4B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.5A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.5B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.6A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.6B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.7A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.7B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.8A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.8B"}],"CooledBy@odata.count":16,"ManagedBy":[{"@odata.id":"/redfish/v1/Managers/iDRAC.Embedded.1"}],"ManagedBy@odata.count":1,"ManagersInChassis":[{"@odata.id":"/redfish/v1/Managers/iDRAC.Embedded.1"}],"ManagersInChassis@odata.count":1,"PoweredBy":[{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Power/PowerSupplies/PSU.Slot.1"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Power/PowerSupplies/PSU.Slot.2"}],"PoweredBy@odata.count":2},"Location":{"Info":null,"InfoFormat":null},"Manufacturer":" ","Model":" ","Name":"Computer System Chassis","PartNumber":"0CRT1GA03","PhysicalSecurity":{"IntrusionSensor":"Normal","IntrusionSensorNumber":115,"IntrusionSensorReArm":"Manual"},"Power":{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Power"},"PowerState":"On","SKU":"8ZLH8M2","SerialNumber":"CNIVC007B60269","Status":{"Health":"OK","HealthRollup":"OK","State":"Enabled"},"Thermal":{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Thermal"}}
`},
	{"https://x0c0s8.ice.next.cray.com/redfish/v1/Managers/iDRAC.Embedded.1",
		"GET",
		``,
		map[string][]string{"Authorization": []string{"Basic cm9vdDppbml0aWFsMA=="}, "Accept": []string{"*/*"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Allow": []string{"POST"}, "Connection": []string{"Keep-Alive"}, "Odata-Version": []string{"4.0"}, "Keep-Alive": []string{"timeout=60, max=199"}, "Server": []string{"Appweb/4.5.4"}, "Link": []string{"</redfish/v1/Schemas/Manager.v1_0_2.json>;rel=describedby"}, "Cache-Control": []string{"no-cache"}, "Content-Length": []string{"3521"}, "Access-Control-Allow-Origin": []string{"*"}, "Accept-Ranges": []string{"bytes"}, "Content-Type": []string{"application/json;odata.metadata=minimal;charset=utf-8"}, "Date": []string{"Wed, 27 Feb 2019 04:18:06 GMT"}},
		`{"@odata.context":"/redfish/v1/$metadata#Manager.Manager","@odata.id":"/redfish/v1/Managers/iDRAC.Embedded.1","@odata.type":"#Manager.v1_0_2.Manager","Actions":{"#Manager.Reset":{"ResetType@Redfish.AllowableValues":["GracefulRestart"],"target":"/redfish/v1/Managers/iDRAC.Embedded.1/Actions/Manager.Reset"},"Oem":{"OemManager.v1_0_0#OemManager.ExportSystemConfiguration":{"ExportFormat@Redfish.AllowableValues":["XML"],"ExportUse@Redfish.AllowableValues":["Default","Clone","Replace"],"IncludeInExport@Redfish.AllowableValues":["Default","IncludeReadOnly","IncludePasswordHashValues"],"ShareParameters":{"ShareParameters@Redfish.AllowableValues":["IPAddress","ShareName","FileName","UserName","Password","Workgroup"],"ShareType@Redfish.AllowableValues":["NFS","CIFS"],"Target@Redfish.AllowableValues":["ALL","IDRAC","BIOS","NIC","RAID"]},"target":"/redfish/v1/Managers/iDRAC.Embedded.1/Actions/Oem/EID_674_Manager.ExportSystemConfiguration"},"OemManager.v1_0_0#OemManager.ImportSystemConfiguration":{"HostPowerState@Redfish.AllowableValues":["On","Off"],"ImportSystemConfiguration@Redfish.AllowableValues":["TimeToWait","ImportBuffer"],"ShareParameters":{"ShareParameters@Redfish.AllowableValues":["IPAddress","ShareName","FileName","UserName","Password","Workgroup"],"ShareType@Redfish.AllowableValues":["NFS","CIFS"],"Target@Redfish.AllowableValues":["ALL","IDRAC","BIOS","NIC","RAID"]},"ShutdownType@Redfish.AllowableValues":["Graceful","Forced","NoReboot"],"target":"/redfish/v1/Managers/iDRAC.Embedded.1/Actions/Oem/EID_674_Manager.ImportSystemConfiguration"},"OemManager.v1_0_0#OemManager.ImportSystemConfigurationPreview":{"ImportSystemConfigurationPreview@Redfish.AllowableValues":["ImportBuffer"],"ShareParameters":{"ShareParameters@Redfish.AllowableValues":["IPAddress","ShareName","FileName","UserName","Password","Workgroup"],"ShareType@Redfish.AllowableValues":["NFS","CIFS"],"Target@Redfish.AllowableValues":["ALL"]},"target":"/redfish/v1/Managers/iDRAC.Embedded.1/Actions/Oem/EID_674_Manager.ImportSystemConfigurationPreview"}}},"CommandShell":{"ConnectTypesSupported":["SSH","Telnet","IPMI"],"ConnectTypesSupported@odata.count":3,"MaxConcurrentSessions":5,"ServiceEnabled":true},"DateTime":"2019-02-26T22:18:07-06:00","DateTimeLocalOffset":"-06:00","Description":"BMC","EthernetInterfaces":{"@odata.id":"/redfish/v1/Managers/iDRAC.Embedded.1/EthernetInterfaces"},"FirmwareVersion":"2.41.40.40","GraphicalConsole":{"ConnectTypesSupported":["KVMIP"],"ConnectTypesSupported@odata.count":1,"MaxConcurrentSessions":6,"ServiceEnabled":true},"Id":"iDRAC.Embedded.1","Links":{"ManagerForChassis":[{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1"}],"ManagerForChassis@odata.count":1,"ManagerForServers":[{"@odata.id":"/redfish/v1/Systems/System.Embedded.1"}],"ManagerForServers@odata.count":1},"LogServices":{"@odata.id":"/redfish/v1/Managers/iDRAC.Embedded.1/LogServices"},"ManagerType":"BMC","Model":"13G Monolithic","Name":"Manager","NetworkProtocol":{"@odata.id":"/redfish/v1/Managers/iDRAC.Embedded.1/NetworkProtocol"},"Redundancy":[],"Redundancy@odata.count":0,"RedundancySet":[],"RedundancySet@odata.count":0,"SerialConsole":{"ConnectTypesSupported":[],"ConnectTypesSupported@odata.count":0,"MaxConcurrentSessions":0,"ServiceEnabled":false},"SerialInterfaces":{"@odata.id":"/redfish/v1/Managers/iDRAC.Embedded.1/SerialInterfaces"},"Status":{"Health":"Ok","State":"Enabled"},"UUID":"3234464f-c0b8-3880-3210-00574c4c4544","VirtualMedia":{"@odata.id":"/redfish/v1/Managers/iDRAC.Embedded.1/VirtualMedia"}}
`},
	{"https://x0c0s10.ice.next.cray.com/redfish/v1/Chassis/System.Embedded.1",
		"GET",
		``,
		map[string][]string{"Authorization": []string{"Basic cm9vdDppbml0aWFsMA=="}, "Accept": []string{"*/*"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Odata-Version": []string{"4.0"}, "Keep-Alive": []string{"timeout=60, max=199"}, "Date": []string{"Wed, 27 Feb 2019 04:21:05 GMT"}, "Link": []string{"</redfish/v1/Schemas/Chassis.v1_0_2.json>;rel=describedby"}, "Cache-Control": []string{"no-cache"}, "Connection": []string{"Keep-Alive"}, "Access-Control-Allow-Origin": []string{"*"}, "Accept-Ranges": []string{"bytes"}, "Content-Type": []string{"application/json;odata.metadata=minimal;charset=utf-8"}, "Server": []string{"Appweb/4.5.4"}, "Content-Length": []string{"2813"}, "Allow": []string{"POST,PATCH"}},
		`{"@odata.context":"/redfish/v1/$metadata#Chassis.Chassis","@odata.id":"/redfish/v1/Chassis/System.Embedded.1","@odata.type":"#Chassis.v1_0_2.Chassis","Actions":{"#Chassis.Reset":{"ResetType@Redfish.AllowableValues":["On","ForceOff"],"Target":"/redfish/v1/Chassis/System.Embedded.1/Actions/Chassis.Reset"}},"AssetTag":"","ChassisType":"Enclosure","Description":"It represents the properties for physical components for any system.It represent racks, rackmount servers, blades, standalone, modular systems,enclosures, and all other containers.The non-cpu/device centric parts of the schema are all accessed either directly or indirectly through this resource.","Id":"System.Embedded.1","IndicatorLED":"Blinking","Links":{"ComputerSystems":[{"@odata.id":"/redfish/v1/Systems/System.Embedded.1"}],"ComputerSystems@odata.count":1,"Contains":[],"Contains@odata.count":0,"CooledBy":[{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.1A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.2A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.3A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.4A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.5A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.6A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.7A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.1B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.2B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.3B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.4B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.5B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.6B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.7B"}],"CooledBy@odata.count":14,"ManagedBy":[{"@odata.id":"/redfish/v1/Managers/iDRAC.Embedded.1"}],"ManagedBy@odata.count":1,"PoweredBy":[{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Power/PowerSupplies/PSU.Slot.1"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Power/PowerSupplies/PSU.Slot.2"}],"PoweredBy@odata.count":2},"Manufacturer":" ","Model":" ","Name":"Computer System Chassis","PartNumber":"0CNCJWA08","Power":{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Power"},"PowerState":"On","SKU":"9L7C382","SerialNumber":"CN7475159F0098","Status":{"Health":"Ok","HealthRollUp":"Ok","State":"Enabled"},"Thermal":{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Thermal"}}
`},
	{"https://x0c0s10.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1",
		"GET",
		``,
		map[string][]string{"Authorization": []string{"Basic cm9vdDppbml0aWFsMA=="}, "Accept": []string{"*/*"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Link": []string{"</redfish/v1/Schemas/ComputerSystem.v1_0_2.json>;rel=describedby"}, "Access-Control-Allow-Origin": []string{"*"}, "Accept-Ranges": []string{"bytes"}, "Keep-Alive": []string{"timeout=60, max=199"}, "Server": []string{"Appweb/4.5.4"}, "Date": []string{"Wed, 27 Feb 2019 04:21:05 GMT"}, "Cache-Control": []string{"no-cache"}, "Content-Length": []string{"3396"}, "Allow": []string{"POST,PATCH"}, "Connection": []string{"Keep-Alive"}, "Odata-Version": []string{"4.0"}, "Content-Type": []string{"application/json;odata.metadata=minimal;charset=utf-8"}},
		`{"@odata.context":"/redfish/v1/$metadata#ComputerSystem.ComputerSystem","@odata.id":"/redfish/v1/Systems/System.Embedded.1","@odata.type":"#ComputerSystem.v1_0_2.ComputerSystem","Actions":{"#ComputerSystem.Reset":{"ResetType@Redfish.AllowableValues":["On","ForceOff","GracefulRestart","PushPowerButton","Nmi"],"target":"/redfish/v1/Systems/System.Embedded.1/Actions/ComputerSystem.Reset"}},"AssetTag":"","BiosVersion":"2.4.3","Boot":{"BootSourceOverrideEnabled":"Once","BootSourceOverrideTarget":"None","BootSourceOverrideTarget@Redfish.AllowableValues":["None","Pxe","Cd","Floppy","Hdd","BiosSetup","Utilities","UefiTarget","SDCard"],"UefiTargetBootSourceOverride":""},"Description":"Computer System which represents a machine (physical or virtual) and the local resources such as memory, cpu and other devices that can be accessed from that machine.","EthernetInterfaces":{"@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces"},"HostName":"","Id":"System.Embedded.1","IndicatorLED":"Blinking","Links":{"Chassis":[{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1"}],"Chassis@odata.count":1,"CooledBy":[{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.1A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.2A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.3A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.4A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.5A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.6A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.7A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.1B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.2B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.3B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.4B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.5B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.6B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.7B"}],"CooledBy@odata.count":14,"ManagedBy":[{"@odata.id":"/redfish/v1/Managers/iDRAC.Embedded.1"}],"ManagedBy@odata.count":1,"PoweredBy":[{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Power/PowerSupplies/PSU.Slot.1"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Power/PowerSupplies/PSU.Slot.2"}],"PoweredBy@odata.count":2},"Manufacturer":" ","MemorySummary":{"Status":{"Health":"OK","HealthRollUp":"OK","State":"Enabled"},"TotalSystemMemoryGiB":64.0},"Model":" ","Name":"System","PartNumber":"0CNCJWA08","PowerState":"On","ProcessorSummary":{"Count":2,"Model":"Intel(R) Xeon(R) CPU E5-2623 v3 @ 3.00GHz","Status":{"Health":"OK","HealthRollUp":"OK","State":"Enabled"}},"Processors":{"@odata.id":"/redfish/v1/Systems/System.Embedded.1/Processors"},"SKU":"9L7C382","SerialNumber":"CN7475159F0098","SimpleStorage":{"@odata.id":"/redfish/v1/Systems/System.Embedded.1/Storage/Controllers"},"Status":{"Health":"OK","HealthRollUp":"OK","State":"Enabled"},"SystemType":"Physical","UUID":"4c4c4544-004c-3710-8043-b9c04f333832"}
`},
	{"https://x0c0s8.ice.next.cray.com/redfish/v1/Chassis/System.Embedded.1",
		"GET",
		``,
		map[string][]string{"Authorization": []string{"Basic cm9vdDppbml0aWFsMA=="}, "Accept": []string{"*/*"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Odata-Version": []string{"4.0"}, "Keep-Alive": []string{"timeout=60, max=199"}, "Content-Type": []string{"application/json;odata.metadata=minimal;charset=utf-8"}, "Server": []string{"Appweb/4.5.4"}, "Link": []string{"</redfish/v1/Schemas/Chassis.v1_0_2.json>;rel=describedby"}, "Content-Length": []string{"2808"}, "Access-Control-Allow-Origin": []string{"*"}, "Date": []string{"Wed, 27 Feb 2019 04:18:08 GMT"}, "Cache-Control": []string{"no-cache"}, "Allow": []string{"POST,PATCH"}, "Connection": []string{"Keep-Alive"}, "Accept-Ranges": []string{"bytes"}},
		`{"@odata.context":"/redfish/v1/$metadata#Chassis.Chassis","@odata.id":"/redfish/v1/Chassis/System.Embedded.1","@odata.type":"#Chassis.v1_0_2.Chassis","Actions":{"#Chassis.Reset":{"ResetType@Redfish.AllowableValues":["On","ForceOff"],"Target":"/redfish/v1/Chassis/System.Embedded.1/Actions/Chassis.Reset"}},"AssetTag":"","ChassisType":"Enclosure","Description":"It represents the properties for physical components for any system.It represent racks, rackmount servers, blades, standalone, modular systems,enclosures, and all other containers.The non-cpu/device centric parts of the schema are all accessed either directly or indirectly through this resource.","Id":"System.Embedded.1","IndicatorLED":"Off","Links":{"ComputerSystems":[{"@odata.id":"/redfish/v1/Systems/System.Embedded.1"}],"ComputerSystems@odata.count":1,"Contains":[],"Contains@odata.count":0,"CooledBy":[{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.1A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.2A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.3A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.4A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.5A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.6A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.7A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.1B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.2B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.3B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.4B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.5B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.6B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.7B"}],"CooledBy@odata.count":14,"ManagedBy":[{"@odata.id":"/redfish/v1/Managers/iDRAC.Embedded.1"}],"ManagedBy@odata.count":1,"PoweredBy":[{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Power/PowerSupplies/PSU.Slot.1"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Power/PowerSupplies/PSU.Slot.2"}],"PoweredBy@odata.count":2},"Manufacturer":" ","Model":" ","Name":"Computer System Chassis","PartNumber":"0CNCJWA05","Power":{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Power"},"PowerState":"On","SKU":"8W28F42","SerialNumber":"CN7475153B0564","Status":{"Health":"Ok","HealthRollUp":"Ok","State":"Enabled"},"Thermal":{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Thermal"}}
`},
	{"https://x0c0s7.ice.next.cray.com/redfish/v1/Managers/iDRAC.Embedded.1",
		"GET",
		``,
		map[string][]string{"Authorization": []string{"Basic cm9vdDppbml0aWFsMA=="}, "Accept": []string{"*/*"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Length": []string{"3521"}, "Connection": []string{"Keep-Alive"}, "Accept-Ranges": []string{"bytes"}, "Odata-Version": []string{"4.0"}, "Keep-Alive": []string{"timeout=60, max=199"}, "Cache-Control": []string{"no-cache"}, "Link": []string{"</redfish/v1/Schemas/Manager.v1_0_2.json>;rel=describedby"}, "Allow": []string{"POST"}, "Access-Control-Allow-Origin": []string{"*"}, "Content-Type": []string{"application/json;odata.metadata=minimal;charset=utf-8"}, "Server": []string{"Appweb/4.5.4"}, "Date": []string{"Wed, 27 Feb 2019 04:24:06 GMT"}},
		`{"@odata.context":"/redfish/v1/$metadata#Manager.Manager","@odata.id":"/redfish/v1/Managers/iDRAC.Embedded.1","@odata.type":"#Manager.v1_0_2.Manager","Actions":{"#Manager.Reset":{"ResetType@Redfish.AllowableValues":["GracefulRestart"],"target":"/redfish/v1/Managers/iDRAC.Embedded.1/Actions/Manager.Reset"},"Oem":{"OemManager.v1_0_0#OemManager.ExportSystemConfiguration":{"ExportFormat@Redfish.AllowableValues":["XML"],"ExportUse@Redfish.AllowableValues":["Default","Clone","Replace"],"IncludeInExport@Redfish.AllowableValues":["Default","IncludeReadOnly","IncludePasswordHashValues"],"ShareParameters":{"ShareParameters@Redfish.AllowableValues":["IPAddress","ShareName","FileName","UserName","Password","Workgroup"],"ShareType@Redfish.AllowableValues":["NFS","CIFS"],"Target@Redfish.AllowableValues":["ALL","IDRAC","BIOS","NIC","RAID"]},"target":"/redfish/v1/Managers/iDRAC.Embedded.1/Actions/Oem/EID_674_Manager.ExportSystemConfiguration"},"OemManager.v1_0_0#OemManager.ImportSystemConfiguration":{"HostPowerState@Redfish.AllowableValues":["On","Off"],"ImportSystemConfiguration@Redfish.AllowableValues":["TimeToWait","ImportBuffer"],"ShareParameters":{"ShareParameters@Redfish.AllowableValues":["IPAddress","ShareName","FileName","UserName","Password","Workgroup"],"ShareType@Redfish.AllowableValues":["NFS","CIFS"],"Target@Redfish.AllowableValues":["ALL","IDRAC","BIOS","NIC","RAID"]},"ShutdownType@Redfish.AllowableValues":["Graceful","Forced","NoReboot"],"target":"/redfish/v1/Managers/iDRAC.Embedded.1/Actions/Oem/EID_674_Manager.ImportSystemConfiguration"},"OemManager.v1_0_0#OemManager.ImportSystemConfigurationPreview":{"ImportSystemConfigurationPreview@Redfish.AllowableValues":["ImportBuffer"],"ShareParameters":{"ShareParameters@Redfish.AllowableValues":["IPAddress","ShareName","FileName","UserName","Password","Workgroup"],"ShareType@Redfish.AllowableValues":["NFS","CIFS"],"Target@Redfish.AllowableValues":["ALL"]},"target":"/redfish/v1/Managers/iDRAC.Embedded.1/Actions/Oem/EID_674_Manager.ImportSystemConfigurationPreview"}}},"CommandShell":{"ConnectTypesSupported":["SSH","Telnet","IPMI"],"ConnectTypesSupported@odata.count":3,"MaxConcurrentSessions":5,"ServiceEnabled":true},"DateTime":"2019-02-26T22:24:06-06:00","DateTimeLocalOffset":"-06:00","Description":"BMC","EthernetInterfaces":{"@odata.id":"/redfish/v1/Managers/iDRAC.Embedded.1/EthernetInterfaces"},"FirmwareVersion":"2.41.40.40","GraphicalConsole":{"ConnectTypesSupported":["KVMIP"],"ConnectTypesSupported@odata.count":1,"MaxConcurrentSessions":6,"ServiceEnabled":true},"Id":"iDRAC.Embedded.1","Links":{"ManagerForChassis":[{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1"}],"ManagerForChassis@odata.count":1,"ManagerForServers":[{"@odata.id":"/redfish/v1/Systems/System.Embedded.1"}],"ManagerForServers@odata.count":1},"LogServices":{"@odata.id":"/redfish/v1/Managers/iDRAC.Embedded.1/LogServices"},"ManagerType":"BMC","Model":"13G Monolithic","Name":"Manager","NetworkProtocol":{"@odata.id":"/redfish/v1/Managers/iDRAC.Embedded.1/NetworkProtocol"},"Redundancy":[],"Redundancy@odata.count":0,"RedundancySet":[],"RedundancySet@odata.count":0,"SerialConsole":{"ConnectTypesSupported":[],"ConnectTypesSupported@odata.count":0,"MaxConcurrentSessions":0,"ServiceEnabled":false},"SerialInterfaces":{"@odata.id":"/redfish/v1/Managers/iDRAC.Embedded.1/SerialInterfaces"},"Status":{"Health":"Ok","State":"Enabled"},"UUID":"3234464f-c0b8-3780-3410-00574c4c4544","VirtualMedia":{"@odata.id":"/redfish/v1/Managers/iDRAC.Embedded.1/VirtualMedia"}}
`},
	{"https://x0c0s9.ice.next.cray.com/redfish/v1/Managers/iDRAC.Embedded.1",
		"GET",
		``,
		map[string][]string{"Authorization": []string{"Basic cm9vdDppbml0aWFsMA=="}, "Accept": []string{"*/*"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Server": []string{"Appweb/4.5.4"}, "Date": []string{"Tue, 26 Feb 2019 23:29:56 GMT"}, "Link": []string{"</redfish/v1/Schemas/Manager.v1_0_2.json>;rel=describedby"}, "Allow": []string{"POST"}, "Accept-Ranges": []string{"bytes"}, "Access-Control-Allow-Origin": []string{"*"}, "Odata-Version": []string{"4.0"}, "Keep-Alive": []string{"timeout=60, max=199"}, "Content-Type": []string{"application/json;odata.metadata=minimal;charset=utf-8"}, "Cache-Control": []string{"no-cache"}, "Content-Length": []string{"3521"}, "Connection": []string{"Keep-Alive"}},
		`{"@odata.context":"/redfish/v1/$metadata#Manager.Manager","@odata.id":"/redfish/v1/Managers/iDRAC.Embedded.1","@odata.type":"#Manager.v1_0_2.Manager","Actions":{"#Manager.Reset":{"ResetType@Redfish.AllowableValues":["GracefulRestart"],"target":"/redfish/v1/Managers/iDRAC.Embedded.1/Actions/Manager.Reset"},"Oem":{"OemManager.v1_0_0#OemManager.ExportSystemConfiguration":{"ExportFormat@Redfish.AllowableValues":["XML"],"ExportUse@Redfish.AllowableValues":["Default","Clone","Replace"],"IncludeInExport@Redfish.AllowableValues":["Default","IncludeReadOnly","IncludePasswordHashValues"],"ShareParameters":{"ShareParameters@Redfish.AllowableValues":["IPAddress","ShareName","FileName","UserName","Password","Workgroup"],"ShareType@Redfish.AllowableValues":["NFS","CIFS"],"Target@Redfish.AllowableValues":["ALL","IDRAC","BIOS","NIC","RAID"]},"target":"/redfish/v1/Managers/iDRAC.Embedded.1/Actions/Oem/EID_674_Manager.ExportSystemConfiguration"},"OemManager.v1_0_0#OemManager.ImportSystemConfiguration":{"HostPowerState@Redfish.AllowableValues":["On","Off"],"ImportSystemConfiguration@Redfish.AllowableValues":["TimeToWait","ImportBuffer"],"ShareParameters":{"ShareParameters@Redfish.AllowableValues":["IPAddress","ShareName","FileName","UserName","Password","Workgroup"],"ShareType@Redfish.AllowableValues":["NFS","CIFS"],"Target@Redfish.AllowableValues":["ALL","IDRAC","BIOS","NIC","RAID"]},"ShutdownType@Redfish.AllowableValues":["Graceful","Forced","NoReboot"],"target":"/redfish/v1/Managers/iDRAC.Embedded.1/Actions/Oem/EID_674_Manager.ImportSystemConfiguration"},"OemManager.v1_0_0#OemManager.ImportSystemConfigurationPreview":{"ImportSystemConfigurationPreview@Redfish.AllowableValues":["ImportBuffer"],"ShareParameters":{"ShareParameters@Redfish.AllowableValues":["IPAddress","ShareName","FileName","UserName","Password","Workgroup"],"ShareType@Redfish.AllowableValues":["NFS","CIFS"],"Target@Redfish.AllowableValues":["ALL"]},"target":"/redfish/v1/Managers/iDRAC.Embedded.1/Actions/Oem/EID_674_Manager.ImportSystemConfigurationPreview"}}},"CommandShell":{"ConnectTypesSupported":["SSH","Telnet","IPMI"],"ConnectTypesSupported@odata.count":3,"MaxConcurrentSessions":5,"ServiceEnabled":true},"DateTime":"2019-02-26T17:29:57-06:00","DateTimeLocalOffset":"-06:00","Description":"BMC","EthernetInterfaces":{"@odata.id":"/redfish/v1/Managers/iDRAC.Embedded.1/EthernetInterfaces"},"FirmwareVersion":"2.41.40.40","GraphicalConsole":{"ConnectTypesSupported":["KVMIP"],"ConnectTypesSupported@odata.count":1,"MaxConcurrentSessions":6,"ServiceEnabled":true},"Id":"iDRAC.Embedded.1","Links":{"ManagerForChassis":[{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1"}],"ManagerForChassis@odata.count":1,"ManagerForServers":[{"@odata.id":"/redfish/v1/Systems/System.Embedded.1"}],"ManagerForServers@odata.count":1},"LogServices":{"@odata.id":"/redfish/v1/Managers/iDRAC.Embedded.1/LogServices"},"ManagerType":"BMC","Model":"13G Monolithic","Name":"Manager","NetworkProtocol":{"@odata.id":"/redfish/v1/Managers/iDRAC.Embedded.1/NetworkProtocol"},"Redundancy":[],"Redundancy@odata.count":0,"RedundancySet":[],"RedundancySet@odata.count":0,"SerialConsole":{"ConnectTypesSupported":[],"ConnectTypesSupported@odata.count":0,"MaxConcurrentSessions":0,"ServiceEnabled":false},"SerialInterfaces":{"@odata.id":"/redfish/v1/Managers/iDRAC.Embedded.1/SerialInterfaces"},"Status":{"Health":"Ok","State":"Enabled"},"UUID":"3242474f-c0b6-4a80-5610-005a4c4c4544","VirtualMedia":{"@odata.id":"/redfish/v1/Managers/iDRAC.Embedded.1/VirtualMedia"}}
`},
	{"https://x0c0s9.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1",
		"GET",
		``,
		map[string][]string{"Authorization": []string{"Basic cm9vdDppbml0aWFsMA=="}, "Accept": []string{"*/*"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json;odata.metadata=minimal;charset=utf-8"}, "Content-Length": []string{"3400"}, "Connection": []string{"Keep-Alive"}, "Odata-Version": []string{"4.0"}, "Keep-Alive": []string{"timeout=60, max=199"}, "Link": []string{"</redfish/v1/Schemas/ComputerSystem.v1_0_2.json>;rel=describedby"}, "Cache-Control": []string{"no-cache"}, "Allow": []string{"POST,PATCH"}, "Access-Control-Allow-Origin": []string{"*"}, "Accept-Ranges": []string{"bytes"}, "Server": []string{"Appweb/4.5.4"}, "Date": []string{"Tue, 26 Feb 2019 23:29:57 GMT"}},
		`{"@odata.context":"/redfish/v1/$metadata#ComputerSystem.ComputerSystem","@odata.id":"/redfish/v1/Systems/System.Embedded.1","@odata.type":"#ComputerSystem.v1_0_2.ComputerSystem","Actions":{"#ComputerSystem.Reset":{"ResetType@Redfish.AllowableValues":["On","ForceOff","GracefulRestart","PushPowerButton","Nmi"],"target":"/redfish/v1/Systems/System.Embedded.1/Actions/ComputerSystem.Reset"}},"AssetTag":"","BiosVersion":"2.0.1","Boot":{"BootSourceOverrideEnabled":"Once","BootSourceOverrideTarget":"None","BootSourceOverrideTarget@Redfish.AllowableValues":["None","Pxe","Cd","Floppy","Hdd","BiosSetup","Utilities","UefiTarget","SDCard"],"UefiTargetBootSourceOverride":""},"Description":"Computer System which represents a machine (physical or virtual) and the local resources such as memory, cpu and other devices that can be accessed from that machine.","EthernetInterfaces":{"@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces"},"HostName":"MINWINPC","Id":"System.Embedded.1","IndicatorLED":"Off","Links":{"Chassis":[{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1"}],"Chassis@odata.count":1,"CooledBy":[{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.1A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.2A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.3A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.4A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.5A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.6A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.7A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.1B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.2B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.3B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.4B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.5B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.6B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.7B"}],"CooledBy@odata.count":14,"ManagedBy":[{"@odata.id":"/redfish/v1/Managers/iDRAC.Embedded.1"}],"ManagedBy@odata.count":1,"PoweredBy":[{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Power/PowerSupplies/PSU.Slot.1"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Power/PowerSupplies/PSU.Slot.2"}],"PoweredBy@odata.count":2},"Manufacturer":" ","MemorySummary":{"Status":{"Health":"OK","HealthRollUp":"OK","State":"Enabled"},"TotalSystemMemoryGiB":128.0},"Model":" ","Name":"System","PartNumber":"0CNCJWA10","PowerState":"On","ProcessorSummary":{"Count":2,"Model":"Intel(R) Xeon(R) CPU E5-2640 v3 @ 2.60GHz","Status":{"Health":"OK","HealthRollUp":"OK","State":"Enabled"}},"Processors":{"@odata.id":"/redfish/v1/Systems/System.Embedded.1/Processors"},"SKU":"6ZVJGB2","SerialNumber":"CN747515C31036","SimpleStorage":{"@odata.id":"/redfish/v1/Systems/System.Embedded.1/Storage/Controllers"},"Status":{"Health":"OK","HealthRollUp":"OK","State":"Enabled"},"SystemType":"Physical","UUID":"4c4c4544-005a-5610-804a-b6c04f474232"}
`},
	{"https://x0c0s9.ice.next.cray.com/redfish/v1/Chassis/System.Embedded.1",
		"GET",
		``,
		map[string][]string{"Authorization": []string{"Basic cm9vdDppbml0aWFsMA=="}, "Accept": []string{"*/*"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Connection": []string{"Keep-Alive"}, "Access-Control-Allow-Origin": []string{"*"}, "Odata-Version": []string{"4.0"}, "Keep-Alive": []string{"timeout=60, max=199"}, "Server": []string{"Appweb/4.5.4"}, "Date": []string{"Tue, 26 Feb 2019 23:29:57 GMT"}, "Link": []string{"</redfish/v1/Schemas/Chassis.v1_0_2.json>;rel=describedby"}, "Cache-Control": []string{"no-cache"}, "Content-Type": []string{"application/json;odata.metadata=minimal;charset=utf-8"}, "Content-Length": []string{"2808"}, "Allow": []string{"POST,PATCH"}, "Accept-Ranges": []string{"bytes"}},
		`{"@odata.context":"/redfish/v1/$metadata#Chassis.Chassis","@odata.id":"/redfish/v1/Chassis/System.Embedded.1","@odata.type":"#Chassis.v1_0_2.Chassis","Actions":{"#Chassis.Reset":{"ResetType@Redfish.AllowableValues":["On","ForceOff"],"Target":"/redfish/v1/Chassis/System.Embedded.1/Actions/Chassis.Reset"}},"AssetTag":"","ChassisType":"Enclosure","Description":"It represents the properties for physical components for any system.It represent racks, rackmount servers, blades, standalone, modular systems,enclosures, and all other containers.The non-cpu/device centric parts of the schema are all accessed either directly or indirectly through this resource.","Id":"System.Embedded.1","IndicatorLED":"Off","Links":{"ComputerSystems":[{"@odata.id":"/redfish/v1/Systems/System.Embedded.1"}],"ComputerSystems@odata.count":1,"Contains":[],"Contains@odata.count":0,"CooledBy":[{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.1A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.2A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.3A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.4A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.5A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.6A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.7A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.1B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.2B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.3B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.4B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.5B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.6B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.7B"}],"CooledBy@odata.count":14,"ManagedBy":[{"@odata.id":"/redfish/v1/Managers/iDRAC.Embedded.1"}],"ManagedBy@odata.count":1,"PoweredBy":[{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Power/PowerSupplies/PSU.Slot.1"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Power/PowerSupplies/PSU.Slot.2"}],"PoweredBy@odata.count":2},"Manufacturer":" ","Model":" ","Name":"Computer System Chassis","PartNumber":"0CNCJWA10","Power":{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Power"},"PowerState":"On","SKU":"6ZVJGB2","SerialNumber":"CN747515C31036","Status":{"Health":"Ok","HealthRollUp":"Ok","State":"Enabled"},"Thermal":{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Thermal"}}
`},
	{"https://x0c0s7.ice.next.cray.com/redfish/v1/Chassis/System.Embedded.1",
		"GET",
		``,
		map[string][]string{"Authorization": []string{"Basic cm9vdDppbml0aWFsMA=="}, "Accept": []string{"*/*"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Date": []string{"Wed, 27 Feb 2019 04:24:09 GMT"}, "Link": []string{"</redfish/v1/Schemas/Chassis.v1_0_2.json>;rel=describedby"}, "Content-Length": []string{"2808"}, "Odata-Version": []string{"4.0"}, "Keep-Alive": []string{"timeout=60, max=199"}, "Cache-Control": []string{"no-cache"}, "Allow": []string{"POST,PATCH"}, "Connection": []string{"Keep-Alive"}, "Access-Control-Allow-Origin": []string{"*"}, "Accept-Ranges": []string{"bytes"}, "Content-Type": []string{"application/json;odata.metadata=minimal;charset=utf-8"}, "Server": []string{"Appweb/4.5.4"}},
		`{"@odata.context":"/redfish/v1/$metadata#Chassis.Chassis","@odata.id":"/redfish/v1/Chassis/System.Embedded.1","@odata.type":"#Chassis.v1_0_2.Chassis","Actions":{"#Chassis.Reset":{"ResetType@Redfish.AllowableValues":["On","ForceOff"],"Target":"/redfish/v1/Chassis/System.Embedded.1/Actions/Chassis.Reset"}},"AssetTag":"","ChassisType":"Enclosure","Description":"It represents the properties for physical components for any system.It represent racks, rackmount servers, blades, standalone, modular systems,enclosures, and all other containers.The non-cpu/device centric parts of the schema are all accessed either directly or indirectly through this resource.","Id":"System.Embedded.1","IndicatorLED":"Off","Links":{"ComputerSystems":[{"@odata.id":"/redfish/v1/Systems/System.Embedded.1"}],"ComputerSystems@odata.count":1,"Contains":[],"Contains@odata.count":0,"CooledBy":[{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.1A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.2A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.3A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.4A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.5A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.6A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.7A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.1B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.2B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.3B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.4B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.5B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.6B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.7B"}],"CooledBy@odata.count":14,"ManagedBy":[{"@odata.id":"/redfish/v1/Managers/iDRAC.Embedded.1"}],"ManagedBy@odata.count":1,"PoweredBy":[{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Power/PowerSupplies/PSU.Slot.1"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Power/PowerSupplies/PSU.Slot.2"}],"PoweredBy@odata.count":2},"Manufacturer":" ","Model":" ","Name":"Computer System Chassis","PartNumber":"0CNCJWA05","Power":{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Power"},"PowerState":"On","SKU":"8W47F42","SerialNumber":"CN7475153B0496","Status":{"Health":"Ok","HealthRollUp":"Ok","State":"Enabled"},"Thermal":{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Thermal"}}
`},
	{"https://x0c0s7.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1",
		"GET",
		``,
		map[string][]string{"Accept": []string{"*/*"}, "Authorization": []string{"Basic cm9vdDppbml0aWFsMA=="}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Access-Control-Allow-Origin": []string{"*"}, "Content-Type": []string{"application/json;odata.metadata=minimal;charset=utf-8"}, "Server": []string{"Appweb/4.5.4"}, "Date": []string{"Wed, 27 Feb 2019 04:24:09 GMT"}, "Link": []string{"</redfish/v1/Schemas/ComputerSystem.v1_0_2.json>;rel=describedby"}, "Content-Length": []string{"3391"}, "Allow": []string{"POST,PATCH"}, "Connection": []string{"Keep-Alive"}, "Accept-Ranges": []string{"bytes"}, "Odata-Version": []string{"4.0"}, "Keep-Alive": []string{"timeout=60, max=199"}, "Cache-Control": []string{"no-cache"}},
		`{"@odata.context":"/redfish/v1/$metadata#ComputerSystem.ComputerSystem","@odata.id":"/redfish/v1/Systems/System.Embedded.1","@odata.type":"#ComputerSystem.v1_0_2.ComputerSystem","Actions":{"#ComputerSystem.Reset":{"ResetType@Redfish.AllowableValues":["On","ForceOff","GracefulRestart","PushPowerButton","Nmi"],"target":"/redfish/v1/Systems/System.Embedded.1/Actions/ComputerSystem.Reset"}},"AssetTag":"","BiosVersion":"2.4.3","Boot":{"BootSourceOverrideEnabled":"Once","BootSourceOverrideTarget":"None","BootSourceOverrideTarget@Redfish.AllowableValues":["None","Pxe","Cd","Floppy","Hdd","BiosSetup","Utilities","UefiTarget","SDCard"],"UefiTargetBootSourceOverride":""},"Description":"Computer System which represents a machine (physical or virtual) and the local resources such as memory, cpu and other devices that can be accessed from that machine.","EthernetInterfaces":{"@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces"},"HostName":"","Id":"System.Embedded.1","IndicatorLED":"Off","Links":{"Chassis":[{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1"}],"Chassis@odata.count":1,"CooledBy":[{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.1A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.2A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.3A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.4A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.5A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.6A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.7A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.1B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.2B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.3B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.4B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.5B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.6B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.7B"}],"CooledBy@odata.count":14,"ManagedBy":[{"@odata.id":"/redfish/v1/Managers/iDRAC.Embedded.1"}],"ManagedBy@odata.count":1,"PoweredBy":[{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Power/PowerSupplies/PSU.Slot.1"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Power/PowerSupplies/PSU.Slot.2"}],"PoweredBy@odata.count":2},"Manufacturer":" ","MemorySummary":{"Status":{"Health":"OK","HealthRollUp":"OK","State":"Enabled"},"TotalSystemMemoryGiB":64.0},"Model":" ","Name":"System","PartNumber":"0CNCJWA05","PowerState":"On","ProcessorSummary":{"Count":2,"Model":"Intel(R) Xeon(R) CPU E5-2623 v3 @ 3.00GHz","Status":{"Health":"OK","HealthRollUp":"OK","State":"Enabled"}},"Processors":{"@odata.id":"/redfish/v1/Systems/System.Embedded.1/Processors"},"SKU":"8W47F42","SerialNumber":"CN7475153B0496","SimpleStorage":{"@odata.id":"/redfish/v1/Systems/System.Embedded.1/Storage/Controllers"},"Status":{"Health":"OK","HealthRollUp":"OK","State":"Enabled"},"SystemType":"Physical","UUID":"4c4c4544-0057-3410-8037-b8c04f463432"}
`},
	{"https://x0c0s8.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1",
		"GET",
		``,
		map[string][]string{"Authorization": []string{"Basic cm9vdDppbml0aWFsMA=="}, "Accept": []string{"*/*"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Odata-Version": []string{"4.0"}, "Content-Type": []string{"application/json;odata.metadata=minimal;charset=utf-8"}, "Date": []string{"Wed, 27 Feb 2019 04:18:12 GMT"}, "Cache-Control": []string{"no-cache"}, "Access-Control-Allow-Origin": []string{"*"}, "Keep-Alive": []string{"timeout=60, max=199"}, "Server": []string{"Appweb/4.5.4"}, "Link": []string{"</redfish/v1/Schemas/ComputerSystem.v1_0_2.json>;rel=describedby"}, "Content-Length": []string{"3392"}, "Allow": []string{"POST,PATCH"}, "Connection": []string{"Keep-Alive"}, "Accept-Ranges": []string{"bytes"}},
		`{"@odata.context":"/redfish/v1/$metadata#ComputerSystem.ComputerSystem","@odata.id":"/redfish/v1/Systems/System.Embedded.1","@odata.type":"#ComputerSystem.v1_0_2.ComputerSystem","Actions":{"#ComputerSystem.Reset":{"ResetType@Redfish.AllowableValues":["On","ForceOff","GracefulRestart","PushPowerButton","Nmi"],"target":"/redfish/v1/Systems/System.Embedded.1/Actions/ComputerSystem.Reset"}},"AssetTag":"","BiosVersion":"2.1.7","Boot":{"BootSourceOverrideEnabled":"Once","BootSourceOverrideTarget":"None","BootSourceOverrideTarget@Redfish.AllowableValues":["None","Pxe","Cd","Floppy","Hdd","BiosSetup","Utilities","UefiTarget","SDCard"],"UefiTargetBootSourceOverride":""},"Description":"Computer System which represents a machine (physical or virtual) and the local resources such as memory, cpu and other devices that can be accessed from that machine.","EthernetInterfaces":{"@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces"},"HostName":"","Id":"System.Embedded.1","IndicatorLED":"Off","Links":{"Chassis":[{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1"}],"Chassis@odata.count":1,"CooledBy":[{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.1A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.2A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.3A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.4A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.5A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.6A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.7A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.1B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.2B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.3B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.4B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.5B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.6B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.7B"}],"CooledBy@odata.count":14,"ManagedBy":[{"@odata.id":"/redfish/v1/Managers/iDRAC.Embedded.1"}],"ManagedBy@odata.count":1,"PoweredBy":[{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Power/PowerSupplies/PSU.Slot.1"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Power/PowerSupplies/PSU.Slot.2"}],"PoweredBy@odata.count":2},"Manufacturer":" ","MemorySummary":{"Status":{"Health":"OK","HealthRollUp":"OK","State":"Enabled"},"TotalSystemMemoryGiB":128.0},"Model":" ","Name":"System","PartNumber":"0CNCJWA05","PowerState":"On","ProcessorSummary":{"Count":2,"Model":"Intel(R) Xeon(R) CPU E5-2640 v3 @ 2.60GHz","Status":{"Health":"OK","HealthRollUp":"OK","State":"Enabled"}},"Processors":{"@odata.id":"/redfish/v1/Systems/System.Embedded.1/Processors"},"SKU":"8W28F42","SerialNumber":"CN7475153B0564","SimpleStorage":{"@odata.id":"/redfish/v1/Systems/System.Embedded.1/Storage/Controllers"},"Status":{"Health":"OK","HealthRollUp":"OK","State":"Enabled"},"SystemType":"Physical","UUID":"4c4c4544-0057-3210-8038-b8c04f463432"}
`},
}
var xnameStatusAllSSData = []sstorage.MockLookup{
	{
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname:    "x0c0s10b0",
				URL:      "x0c0s10b0.ice.next.cray.com/redfish/v1/Managers/iDRAC.Embedded.1",
				Username: "root",
				Password: "********",
			},
			Err: nil,
		},
	}, {
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname:    "x0c0s10b0n0",
				URL:      "x0c0s10b0.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1",
				Username: "root",
				Password: "********",
			},
			Err: nil,
		},
	}, {
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname:    "x0c0s10e0",
				URL:      "x0c0s10b0.ice.next.cray.com/redfish/v1/Chassis/System.Embedded.1",
				Username: "root",
				Password: "********",
			},
			Err: nil,
		},
	}, {
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname:    "x0c0s12b0",
				URL:      "x0c0s12b0.ice.next.cray.com/redfish/v1/Managers/iDRAC.Embedded.1",
				Username: "root",
				Password: "********",
			},
			Err: nil,
		},
	}, {
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname:    "x0c0s12b0n0",
				URL:      "x0c0s12b0.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1",
				Username: "root",
				Password: "********",
			},
			Err: nil,
		},
	}, {
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname:    "x0c0s12e0",
				URL:      "x0c0s12b0.ice.next.cray.com/redfish/v1/Chassis/System.Embedded.1",
				Username: "root",
				Password: "********",
			},
			Err: nil,
		},
	}, {
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname:    "x0c0s7b0",
				URL:      "x0c0s7b0.ice.next.cray.com/redfish/v1/Managers/iDRAC.Embedded.1",
				Username: "root",
				Password: "********",
			},
			Err: nil,
		},
	}, {
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname:    "x0c0s7b0n0",
				URL:      "x0c0s7b0.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1",
				Username: "root",
				Password: "********",
			},
			Err: nil,
		},
	}, {
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname:    "x0c0s7e0",
				URL:      "x0c0s7b0.ice.next.cray.com/redfish/v1/Chassis/System.Embedded.1",
				Username: "root",
				Password: "********",
			},
			Err: nil,
		},
	}, {
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname:    "x0c0s8b0",
				URL:      "x0c0s8b0.ice.next.cray.com/redfish/v1/Managers/iDRAC.Embedded.1",
				Username: "root",
				Password: "********",
			},
			Err: nil,
		},
	}, {
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname:    "x0c0s8b0n0",
				URL:      "x0c0s8b0.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1",
				Username: "root",
				Password: "********",
			},
			Err: nil,
		},
	}, {
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname:    "x0c0s8e0",
				URL:      "x0c0s8b0.ice.next.cray.com/redfish/v1/Chassis/System.Embedded.1",
				Username: "root",
				Password: "********",
			},
			Err: nil,
		},
	}, {
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname:    "x0c0s9b0",
				URL:      "x0c0s9b0.ice.next.cray.com/redfish/v1/Managers/iDRAC.Embedded.1",
				Username: "root",
				Password: "********",
			},
			Err: nil,
		},
	}, {
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname:    "x0c0s9b0n0",
				URL:      "x0c0s9b0.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1",
				Username: "root",
				Password: "********",
			},
			Err: nil,
		},
	}, {
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname:    "x0c0s9e0",
				URL:      "x0c0s9b0.ice.next.cray.com/redfish/v1/Chassis/System.Embedded.1",
				Username: "root",
				Password: "********",
			},
			Err: nil,
		},
	},
}

func TestXnameStatusAll(t *testing.T) {
	var testReq = testData{"/capmc/get_xname_status",
		"POST",
		`{}`,
		map[string][]string{"Content-Type": []string{"application/x-www-form-urlencoded"}, "User-Agent": []string{"curl/7.47.0"}, "Accept": []string{"*/*"}, "Content-Length": []string{"2"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}},
		`{"e":0,"err_msg":"","on":["x0c0s7b0n0","x0c0s8b0n0","x0c0s9b0n0","x0c0s10b0n0", "x0c0s12b0n0"]}
`}
	runTest(t, xnameStatusAllHSM, &testReq, &xnameStatusAllReplayData, xnameStatusAllSSData)
}

var xnameStatusError1HSM = "https://localhost:27779"
var xnameStatusError1ReplayData = []testData{
	{"https://localhost:27779/hsm/v1/State/Components?id=x0c0s0b0n0&type=CabinetPDUPowerConnector&type=CabinetPDUOutlet&type=Chassis&type=RouterModule&type=HSNBoard&type=ComputeModule&type=Node",
		"GET",
		``,
		map[string][]string{"Accept": []string{"application/json"}, "Content-Type": []string{"application/json"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}, "Date": []string{"Tue, 26 Feb 2019 22:27:09 GMT"}, "Content-Length": []string{"1801"}},
		`{"Components":[]}
`},
}
var xnameStatusError1SSData = []sstorage.MockLookup{}

func TestXnameStatusError1(t *testing.T) {
	var testReq = testData{"/capmc/get_xname_status",
		"POST",
		`{"xnames" : ["x0c0s0b0n0" ]}`,
		map[string][]string{"Content-Type": []string{"application/x-www-form-urlencoded"}, "User-Agent": []string{"curl/7.47.0"}, "Accept": []string{"*/*"}, "Content-Length": []string{"28"}},
		"400 Bad Request", 400,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}},
		`{"e":400,"err_msg":"xnames not found: [x0c0s0b0n0]"}
`}
	runTest(t, xnameStatusError1HSM, &testReq, &xnameStatusError1ReplayData, xnameStatusError1SSData)
}

var xnameStatusError2HSM = "https://localhost:27779"
var xnameStatusError2ReplayData = []testData{}
var xnameStatusError2SSData = []sstorage.MockLookup{}

func TestXnameStatusError2(t *testing.T) {
	var testReq = testData{"/capmc/get_xname_status",
		"POST",
		`{"xnames" : [ 4 ]}`,
		map[string][]string{"User-Agent": []string{"curl/7.47.0"}, "Accept": []string{"*/*"}, "Content-Length": []string{"18"}, "Content-Type": []string{"application/x-www-form-urlencoded"}},
		"400 Bad Request", 400,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}},
		`{"e":400,"err_msg":"json: cannot unmarshal number into Go struct field XnameStatusRequest.xnames of type string"}
`}
	runTest(t, xnameStatusError2HSM, &testReq, &xnameStatusError2ReplayData, xnameStatusError2SSData)
}

var xnameStatusSomeHSM = "https://localhost:27779"
var xnameStatusSomeReplayData = []testData{
	{"https://localhost:27779/hsm/v1/State/Components?id=x0c0s7b0n0&id=x0c0s8b0n0&id=x0c0s12b0n0&type=CabinetPDUPowerConnector&type=CabinetPDUOutlet&type=Chassis&type=RouterModule&type=HSNBoard&type=ComputeModule&type=Node",
		"GET",
		``,
		map[string][]string{"Accept": []string{"application/json"}, "Content-Type": []string{"application/json"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}, "Date": []string{"Tue, 26 Feb 2019 22:26:11 GMT"}, "Content-Length": []string{"1801"}},
		`{"Components":[{"ID":"x0c0s12b0n0","Type":"Node","State":"On","Flag":"OK","Enabled":true,"Role":"Compute","NID":16768,"NetType":"Sling","Arch":"X86"},{"ID":"x0c0s7b0n0","Type":"Node","State":"On","Flag":"OK","Enabled":true,"Role":"Compute","NID":16608,"NetType":"Sling","Arch":"X86"},{"ID":"x0c0s8b0n0","Type":"Node","State":"On","Flag":"OK","Enabled":true,"Role":"Compute","NID":16640,"NetType":"Sling","Arch":"X86"}]}
`},
	{"https://localhost:27779/hsm/v1/Inventory/ComponentEndpoints?id=x0c0s7b0n0&id=x0c0s8b0n0&id=x0c0s12b0n0&type=CabinetPDUPowerConnector&type=CabinetPDUOutlet&type=Chassis&type=RouterModule&type=HSNBoard&type=ComputeModule&type=Node",
		"GET",
		``,
		map[string][]string{"Accept": []string{"application/json"}, "Content-Type": []string{"application/json"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}, "Date": []string{"Tue, 26 Feb 2019 22:26:11 GMT"}},
		`{"ComponentEndpoints":[{"ID":"x0c0s12b0n0","Type":"Node","Domain":"ice.next.cray.com","FQDN":"x0c0s12b0n0.ice.next.cray.com","RedfishType":"ComputerSystem","RedfishSubtype":"Physical","MACAddr":"24:6e:96:91:62:3a","UUID":"4c4c4544-005a-4c10-8048-b8c04f384d32","OdataID":"/redfish/v1/Systems/System.Embedded.1","RedfishEndpointID":"x0c0s12b0","RedfishEndpointFQDN":"x0c0s12.ice.next.cray.com","RedfishURL":"x0c0s12.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1","ComponentEndpointType":"ComponentEndpointComputerSystem","RedfishSystemInfo":{"Name":"System","Actions":{"#ComputerSystem.Reset":{"ResetType@Redfish.AllowableValues":["On","ForceOff","GracefulRestart","GracefulShutdown","PushPowerButton","Nmi"],"target":"/redfish/v1/Systems/System.Embedded.1/Actions/ComputerSystem.Reset"}},"EthernetNICInfo":[{"RedfishId":"NIC.Integrated.1-3-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-3-1","Description":"Integrated NIC 1 Port 3 Partition 1","MACAddress":"24:6e:96:91:62:3a","PermanentMACAddress":"24:6e:96:91:62:3a"},{"RedfishId":"NIC.Integrated.1-2-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-2-1","Description":"Integrated NIC 1 Port 2 Partition 1","MACAddress":"24:6e:96:91:62:39","PermanentMACAddress":"24:6e:96:91:62:39"},{"RedfishId":"NIC.Integrated.1-1-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-1-1","Description":"Integrated NIC 1 Port 1 Partition 1","MACAddress":"24:6e:96:91:62:38","PermanentMACAddress":"24:6e:96:91:62:38"},{"RedfishId":"NIC.Integrated.1-4-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-4-1","Description":"Integrated NIC 1 Port 4 Partition 1","MACAddress":"24:6e:96:91:62:3b","PermanentMACAddress":"24:6e:96:91:62:3b"}]}},{"ID":"x0c0s7b0n0","Type":"Node","Domain":"ice.next.cray.com","FQDN":"x0c0s7b0n0.ice.next.cray.com","RedfishType":"ComputerSystem","RedfishSubtype":"Physical","MACAddr":"44:a8:42:21:a8:af","UUID":"4c4c4544-0057-3410-8037-b8c04f463432","OdataID":"/redfish/v1/Systems/System.Embedded.1","RedfishEndpointID":"x0c0s7b0","RedfishEndpointFQDN":"x0c0s7.ice.next.cray.com","RedfishURL":"x0c0s7.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1","ComponentEndpointType":"ComponentEndpointComputerSystem","RedfishSystemInfo":{"Name":"System","Actions":{"#ComputerSystem.Reset":{"ResetType@Redfish.AllowableValues":["On","ForceOff","GracefulRestart","PushPowerButton","Nmi"],"target":"/redfish/v1/Systems/System.Embedded.1/Actions/ComputerSystem.Reset"}},"EthernetNICInfo":[{"RedfishId":"NIC.Integrated.1-1-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-1-1","Description":"Integrated NIC 1 Port 1 Partition 1","MACAddress":"44:a8:42:21:a8:ad","PermanentMACAddress":"44:a8:42:21:a8:ad"},{"RedfishId":"NIC.Integrated.1-2-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-2-1","Description":"Integrated NIC 1 Port 2 Partition 1","MACAddress":"44:a8:42:21:a8:ae","PermanentMACAddress":"44:a8:42:21:a8:ae"},{"RedfishId":"NIC.Integrated.1-3-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-3-1","Description":"Integrated NIC 1 Port 3 Partition 1","MACAddress":"44:a8:42:21:a8:af","PermanentMACAddress":"44:a8:42:21:a8:af"},{"RedfishId":"NIC.Integrated.1-4-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-4-1","Description":"Integrated NIC 1 Port 4 Partition 1","MACAddress":"44:a8:42:21:a8:b0","PermanentMACAddress":"44:a8:42:21:a8:b0"}]}},{"ID":"x0c0s8b0n0","Type":"Node","Domain":"ice.next.cray.com","FQDN":"x0c0s8b0n0.ice.next.cray.com","RedfishType":"ComputerSystem","RedfishSubtype":"Physical","MACAddr":"44:a8:42:21:a1:fd","UUID":"4c4c4544-0057-3210-8038-b8c04f463432","OdataID":"/redfish/v1/Systems/System.Embedded.1","RedfishEndpointID":"x0c0s8b0","RedfishEndpointFQDN":"x0c0s8.ice.next.cray.com","RedfishURL":"x0c0s8.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1","ComponentEndpointType":"ComponentEndpointComputerSystem","RedfishSystemInfo":{"Name":"System","Actions":{"#ComputerSystem.Reset":{"ResetType@Redfish.AllowableValues":["On","ForceOff","GracefulRestart","PushPowerButton","Nmi"],"target":"/redfish/v1/Systems/System.Embedded.1/Actions/ComputerSystem.Reset"}},"EthernetNICInfo":[{"RedfishId":"NIC.Integrated.1-3-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-3-1","Description":"Integrated NIC 1 Port 3 Partition 1","MACAddress":"44:a8:42:21:a1:fd","PermanentMACAddress":"44:a8:42:21:a1:fd"},{"RedfishId":"NIC.Integrated.1-4-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-4-1","Description":"Integrated NIC 1 Port 4 Partition 1","MACAddress":"44:a8:42:21:a1:fe","PermanentMACAddress":"44:a8:42:21:a1:fe"},{"RedfishId":"NIC.Integrated.1-1-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-1-1","Description":"Integrated NIC 1 Port 1 Partition 1","MACAddress":"44:a8:42:21:a1:fb","PermanentMACAddress":"44:a8:42:21:a1:fb"},{"RedfishId":"NIC.Integrated.1-2-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-2-1","Description":"Integrated NIC 1 Port 2 Partition 1","MACAddress":"44:a8:42:21:a1:fc","PermanentMACAddress":"44:a8:42:21:a1:fc"}]}}]}
`},
	{"https://x0c0s12.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1",
		"GET",
		``,
		map[string][]string{"Authorization": []string{"Basic cm9vdDppbml0aWFsMA=="}, "Accept": []string{"*/*"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Length": []string{"4036"}, "Date": []string{"Wed, 27 Feb 2019 04:26:52 GMT"}, "Server": []string{"Apache/2.4"}, "Allow": []string{"POST,PATCH"}, "Link": []string{"</redfish/v1/Schemas/ComputerSystem.v1_1_0.json>;rel=describedby"}, "X-Frame-Options": []string{"DENY"}, "Odata-Version": []string{"4.0"}, "Access-Control-Allow-Origin": []string{"*"}, "Cache-Control": []string{"no-cache"}, "Content-Type": []string{"application/json;odata.metadata=minimal;charset=utf-8"}},
		`{"@odata.context":"/redfish/v1/$metadata#ComputerSystem.ComputerSystem","@odata.id":"/redfish/v1/Systems/System.Embedded.1","@odata.type":"#ComputerSystem.v1_1_0.ComputerSystem","Actions":{"#ComputerSystem.Reset":{"ResetType@Redfish.AllowableValues":["On","ForceOff","GracefulRestart","GracefulShutdown","PushPowerButton","Nmi"],"target":"/redfish/v1/Systems/System.Embedded.1/Actions/ComputerSystem.Reset"}},"AssetTag":"","Bios":{"@odata.id":"/redfish/v1/Systems/System.Embedded.1/Bios"},"BiosVersion":"1.1.7","Boot":{"BootSourceOverrideEnabled":"Once","BootSourceOverrideMode":"UEFI","BootSourceOverrideTarget":"None","BootSourceOverrideTarget@Redfish.AllowableValues":["None","Pxe","Floppy","Cd","Hdd","BiosSetup","Utilities","UefiTarget","SDCard","UefiHttp"],"UefiTargetBootSourceOverride":""},"Description":"Computer System which represents a machine (physical or virtual) and the local resources such as memory, cpu and other devices that can be accessed from that machine.","EthernetInterfaces":{"@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces"},"HostName":"","Id":"System.Embedded.1","IndicatorLED":"Off","Links":{"Chassis":[{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1"}],"Chassis@odata.count":1,"CooledBy":[{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.1A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.1B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.2A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.2B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.3A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.3B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.4A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.4B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.5A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.5B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.6A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.6B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.7A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.7B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.8A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.8B"}],"CooledBy@odata.count":16,"ManagedBy":[{"@odata.id":"/redfish/v1/Managers/iDRAC.Embedded.1"}],"ManagedBy@odata.count":1,"Oem":{"Dell":{"@odata.type":"#DellComputerSystem.v1_0_0.DellComputerSystem","BootOrder":{"@odata.id":"/redfish/v1/Systems/System.Embedded.1/BootSources"}}},"PoweredBy":[{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Power/PowerSupplies/PSU.Slot.1"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Power/PowerSupplies/PSU.Slot.2"}],"PoweredBy@odata.count":2},"Manufacturer":" ","MemorySummary":{"MemoryMirroring":"System","Status":{"Health":"OK","HealthRollup":"OK","State":"Enabled"},"TotalSystemMemoryGiB":64.0},"Model":" ","Name":"System","PartNumber":"0CRT1GA03","PowerState":"On","ProcessorSummary":{"Count":2,"Model":"Intel(R) Xeon(R) Silver 4112 CPU @ 2.60GHz","Status":{"Health":"OK","HealthRollup":"OK","State":"Enabled"}},"Processors":{"@odata.id":"/redfish/v1/Systems/System.Embedded.1/Processors"},"SKU":"8ZLH8M2","SecureBoot":{"@odata.id":"/redfish/v1/Systems/System.Embedded.1/SecureBoot"},"SerialNumber":"CNIVC007B60269","SimpleStorage":{"@odata.id":"/redfish/v1/Systems/System.Embedded.1/Storage/Controllers"},"Status":{"Health":"OK","HealthRollup":"OK","State":"Enabled"},"SystemType":"Physical","TrustedModules":[{"InterfaceType":"TPM1_2","Status":{"State":"Disabled"}}],"UUID":"4c4c4544-005a-4c10-8048-b8c04f384d32"}
`},
	{"https://x0c0s8.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1",
		"GET",
		``,
		map[string][]string{"Authorization": []string{"Basic cm9vdDppbml0aWFsMA=="}, "Accept": []string{"*/*"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Length": []string{"3392"}, "Access-Control-Allow-Origin": []string{"*"}, "Accept-Ranges": []string{"bytes"}, "Odata-Version": []string{"4.0"}, "Keep-Alive": []string{"timeout=60, max=199"}, "Content-Type": []string{"application/json;odata.metadata=minimal;charset=utf-8"}, "Server": []string{"Appweb/4.5.4"}, "Date": []string{"Wed, 27 Feb 2019 04:19:40 GMT"}, "Link": []string{"</redfish/v1/Schemas/ComputerSystem.v1_0_2.json>;rel=describedby"}, "Cache-Control": []string{"no-cache"}, "Allow": []string{"POST,PATCH"}, "Connection": []string{"Keep-Alive"}},
		`{"@odata.context":"/redfish/v1/$metadata#ComputerSystem.ComputerSystem","@odata.id":"/redfish/v1/Systems/System.Embedded.1","@odata.type":"#ComputerSystem.v1_0_2.ComputerSystem","Actions":{"#ComputerSystem.Reset":{"ResetType@Redfish.AllowableValues":["On","ForceOff","GracefulRestart","PushPowerButton","Nmi"],"target":"/redfish/v1/Systems/System.Embedded.1/Actions/ComputerSystem.Reset"}},"AssetTag":"","BiosVersion":"2.1.7","Boot":{"BootSourceOverrideEnabled":"Once","BootSourceOverrideTarget":"None","BootSourceOverrideTarget@Redfish.AllowableValues":["None","Pxe","Cd","Floppy","Hdd","BiosSetup","Utilities","UefiTarget","SDCard"],"UefiTargetBootSourceOverride":""},"Description":"Computer System which represents a machine (physical or virtual) and the local resources such as memory, cpu and other devices that can be accessed from that machine.","EthernetInterfaces":{"@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces"},"HostName":"","Id":"System.Embedded.1","IndicatorLED":"Off","Links":{"Chassis":[{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1"}],"Chassis@odata.count":1,"CooledBy":[{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.1A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.2A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.3A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.4A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.5A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.6A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.7A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.1B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.2B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.3B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.4B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.5B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.6B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.7B"}],"CooledBy@odata.count":14,"ManagedBy":[{"@odata.id":"/redfish/v1/Managers/iDRAC.Embedded.1"}],"ManagedBy@odata.count":1,"PoweredBy":[{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Power/PowerSupplies/PSU.Slot.1"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Power/PowerSupplies/PSU.Slot.2"}],"PoweredBy@odata.count":2},"Manufacturer":" ","MemorySummary":{"Status":{"Health":"OK","HealthRollUp":"OK","State":"Enabled"},"TotalSystemMemoryGiB":128.0},"Model":" ","Name":"System","PartNumber":"0CNCJWA05","PowerState":"On","ProcessorSummary":{"Count":2,"Model":"Intel(R) Xeon(R) CPU E5-2640 v3 @ 2.60GHz","Status":{"Health":"OK","HealthRollUp":"OK","State":"Enabled"}},"Processors":{"@odata.id":"/redfish/v1/Systems/System.Embedded.1/Processors"},"SKU":"8W28F42","SerialNumber":"CN7475153B0564","SimpleStorage":{"@odata.id":"/redfish/v1/Systems/System.Embedded.1/Storage/Controllers"},"Status":{"Health":"OK","HealthRollUp":"OK","State":"Enabled"},"SystemType":"Physical","UUID":"4c4c4544-0057-3210-8038-b8c04f463432"}
`},
	{"https://x0c0s7.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1",
		"GET",
		``,
		map[string][]string{"Authorization": []string{"Basic cm9vdDppbml0aWFsMA=="}, "Accept": []string{"*/*"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Odata-Version": []string{"4.0"}, "Date": []string{"Wed, 27 Feb 2019 04:25:39 GMT"}, "Cache-Control": []string{"no-cache"}, "Content-Length": []string{"3391"}, "Connection": []string{"Keep-Alive"}, "Access-Control-Allow-Origin": []string{"*"}, "Keep-Alive": []string{"timeout=60, max=199"}, "Content-Type": []string{"application/json;odata.metadata=minimal;charset=utf-8"}, "Server": []string{"Appweb/4.5.4"}, "Link": []string{"</redfish/v1/Schemas/ComputerSystem.v1_0_2.json>;rel=describedby"}, "Allow": []string{"POST,PATCH"}, "Accept-Ranges": []string{"bytes"}},
		`{"@odata.context":"/redfish/v1/$metadata#ComputerSystem.ComputerSystem","@odata.id":"/redfish/v1/Systems/System.Embedded.1","@odata.type":"#ComputerSystem.v1_0_2.ComputerSystem","Actions":{"#ComputerSystem.Reset":{"ResetType@Redfish.AllowableValues":["On","ForceOff","GracefulRestart","PushPowerButton","Nmi"],"target":"/redfish/v1/Systems/System.Embedded.1/Actions/ComputerSystem.Reset"}},"AssetTag":"","BiosVersion":"2.4.3","Boot":{"BootSourceOverrideEnabled":"Once","BootSourceOverrideTarget":"None","BootSourceOverrideTarget@Redfish.AllowableValues":["None","Pxe","Cd","Floppy","Hdd","BiosSetup","Utilities","UefiTarget","SDCard"],"UefiTargetBootSourceOverride":""},"Description":"Computer System which represents a machine (physical or virtual) and the local resources such as memory, cpu and other devices that can be accessed from that machine.","EthernetInterfaces":{"@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces"},"HostName":"","Id":"System.Embedded.1","IndicatorLED":"Off","Links":{"Chassis":[{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1"}],"Chassis@odata.count":1,"CooledBy":[{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.1A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.2A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.3A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.4A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.5A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.6A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.7A"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.1B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.2B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.3B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.4B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.5B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.6B"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17||Fan.Embedded.7B"}],"CooledBy@odata.count":14,"ManagedBy":[{"@odata.id":"/redfish/v1/Managers/iDRAC.Embedded.1"}],"ManagedBy@odata.count":1,"PoweredBy":[{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Power/PowerSupplies/PSU.Slot.1"},{"@odata.id":"/redfish/v1/Chassis/System.Embedded.1/Power/PowerSupplies/PSU.Slot.2"}],"PoweredBy@odata.count":2},"Manufacturer":" ","MemorySummary":{"Status":{"Health":"OK","HealthRollUp":"OK","State":"Enabled"},"TotalSystemMemoryGiB":64.0},"Model":" ","Name":"System","PartNumber":"0CNCJWA05","PowerState":"On","ProcessorSummary":{"Count":2,"Model":"Intel(R) Xeon(R) CPU E5-2623 v3 @ 3.00GHz","Status":{"Health":"OK","HealthRollUp":"OK","State":"Enabled"}},"Processors":{"@odata.id":"/redfish/v1/Systems/System.Embedded.1/Processors"},"SKU":"8W47F42","SerialNumber":"CN7475153B0496","SimpleStorage":{"@odata.id":"/redfish/v1/Systems/System.Embedded.1/Storage/Controllers"},"Status":{"Health":"OK","HealthRollUp":"OK","State":"Enabled"},"SystemType":"Physical","UUID":"4c4c4544-0057-3410-8037-b8c04f463432"}
`},
}
var xnameStatusSomeSSData = []sstorage.MockLookup{
	{
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname:    "x0c0s7b0n0",
				URL:      "x0c0s7b0.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1",
				Username: "root",
				Password: "********",
			},
			Err: nil,
		},
	}, {
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname:    "x0c0s8b0n0",
				URL:      "x0c0s8b0.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1",
				Username: "root",
				Password: "********",
			},
			Err: nil,
		},
	}, {
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname:    "x0c0s12b0n0",
				URL:      "x0c0s12b0.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1",
				Username: "root",
				Password: "********",
			},
			Err: nil,
		},
	},
}

func TestXnameStatusSome(t *testing.T) {
	var testReq = testData{"/capmc/get_xname_status",
		"POST",
		`{"xnames" : ["x0c0s7b0n0", "x0c0s8b0n0", "x0c0s12b0n0" ]}`,
		map[string][]string{"User-Agent": []string{"curl/7.47.0"}, "Accept": []string{"*/*"}, "Content-Length": []string{"57"}, "Content-Type": []string{"application/x-www-form-urlencoded"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}},
		`{"e":0,"err_msg":"","on":["x0c0s7b0n0","x0c0s8b0n0","x0c0s12b0n0"]}
`}
	runTest(t, xnameStatusSomeHSM, &testReq, &xnameStatusSomeReplayData, xnameStatusSomeSSData)
}

var statusNotAuthorizedHSM = "https://localhost:27779"
var statusNotAuthorizedReplayData = []testData{
	{"https://localhost:27779/hsm/v1/State/Components?id=x0c0s7b0n0&type=CabinetPDUPowerConnector&type=CabinetPDUOutlet&type=Chassis&type=RouterModule&type=HSNBoard&type=ComputeModule&type=Node",
		"GET",
		``,
		map[string][]string{"Accept": []string{"application/json"}, "Content-Type": []string{"application/json"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}, "Date": []string{"Thu, 04 Apr 2019 20:30:12 GMT"}, "Content-Length": []string{"1801"}},
		`{"Components":[{"ID":"x0c0s7b0n0","Type":"Node","State":"On","Flag":"OK","Enabled":true,"Role":"Compute","NID":16608,"NetType":"Sling","Arch":"X86"}]}
`},
	{"https://localhost:27779/hsm/v1/Inventory/ComponentEndpoints?id=x0c0s7b0n0&type=CabinetPDUPowerConnector&type=CabinetPDUOutlet&type=Chassis&type=RouterModule&type=HSNBoard&type=ComputeModule&type=Node",
		"GET",
		``,
		map[string][]string{"Accept": []string{"application/json"}, "Content-Type": []string{"application/json"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}, "Date": []string{"Thu, 04 Apr 2019 20:30:12 GMT"}},
		`{"ComponentEndpoints":[{"ID":"x0c0s7b0n0","Type":"Node","Domain":"ice.next.cray.com","FQDN":"x0c0s7b0n0.ice.next.cray.com","RedfishType":"ComputerSystem","RedfishSubtype":"Physical","MACAddr":"44:a8:42:21:a8:af","UUID":"4c4c4544-0057-3410-8037-b8c04f463432","OdataID":"/redfish/v1/Systems/System.Embedded.1","RedfishEndpointID":"x0c0s7b0","RedfishEndpointFQDN":"x0c0s7.ice.next.cray.com","RedfishURL":"x0c0s7.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1","ComponentEndpointType":"ComponentEndpointComputerSystem","RedfishSystemInfo":{"Name":"System","Actions":{"#ComputerSystem.Reset":{"ResetType@Redfish.AllowableValues":["On","ForceOff","GracefulRestart","PushPowerButton","Nmi"],"target":"/redfish/v1/Systems/System.Embedded.1/Actions/ComputerSystem.Reset"}},"EthernetNICInfo":[{"RedfishId":"NIC.Integrated.1-1-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-1-1","Description":"Integrated NIC 1 Port 1 Partition 1","MACAddress":"44:a8:42:21:a8:ad","PermanentMACAddress":"44:a8:42:21:a8:ad"},{"RedfishId":"NIC.Integrated.1-2-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-2-1","Description":"Integrated NIC 1 Port 2 Partition 1","MACAddress":"44:a8:42:21:a8:ae","PermanentMACAddress":"44:a8:42:21:a8:ae"},{"RedfishId":"NIC.Integrated.1-3-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-3-1","Description":"Integrated NIC 1 Port 3 Partition 1","MACAddress":"44:a8:42:21:a8:af","PermanentMACAddress":"44:a8:42:21:a8:af"},{"RedfishId":"NIC.Integrated.1-4-1","@odata.id":"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces/NIC.Integrated.1-4-1","Description":"Integrated NIC 1 Port 4 Partition 1","MACAddress":"44:a8:42:21:a8:b0","PermanentMACAddress":"44:a8:42:21:a8:b0"}]}}]}
`},
	{"https://x0c0s7.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1",
		"GET",
		``,
		map[string][]string{"Authorization": []string{"Basic cm9vdDpiYWRwYXNz"}, "Accept": []string{"*/*"}},
		"401 Unauthorized", 401,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Keep-Alive": []string{"timeout=60, max=199"}, "Date": []string{"Fri, 05 Apr 2019 02:28:56 GMT"}, "Connection": []string{"Keep-Alive"}, "Www-Authenticate": []string{"Basic realm=\"RedfishService\""}, "Accept-Ranges": []string{"bytes"}},
		``},
}
var statusNotAuthorizedSSData = []sstorage.MockLookup{
	{
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname:    "x0c0s7b0n0",
				URL:      "x0c0s7b0.ice.next.cray.com/redfish/v1/Systems/System.Embedded.1",
				Username: "root",
				Password: "********",
			},
			Err: nil,
		},
	},
}

func TestStatusNotAuthorized(t *testing.T) {
	var testReq = testData{"/capmc/get_xname_status",
		"POST",
		`{"xnames":["x0c0s7b0n0"]}`,
		map[string][]string{"User-Agent": []string{"curl/7.60.0"}, "Accept": []string{"*/*"}, "Content-Length": []string{"25"}, "Content-Type": []string{"application/x-www-form-urlencoded"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}},
		`{"e":-1,"err_msg":"Errors encountered with 1/1 Xnames for Status","undefined":["x0c0s7b0n0"]}
`}
	runTest(t, statusNotAuthorizedHSM, &testReq, &statusNotAuthorizedReplayData, statusNotAuthorizedSSData)
}

var StatusComponentNotFoundHSM = "https://frosty-sms.us.cray.com:30443/apis/smd"
var StatusComponentNotFoundReplayData = []testData{
	{"https://frosty-sms.us.cray.com:30443/apis/smd/hsm/v1/State/Components?id=x0c0s21b0n0&type=CabinetPDUPowerConnector&type=CabinetPDUOutlet&type=Chassis&type=RouterModule&type=HSNBoard&type=ComputeModule&type=Node",
		"GET",
		``,
		map[string][]string{"Accept": []string{"application/json"}, "Content-Type": []string{"application/json"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Via": []string{"kong/0.14.1"}, "Content-Type": []string{"application/json"}, "Content-Length": []string{"1457"}, "Connection": []string{"keep-alive"}, "Date": []string{"Wed, 10 Apr 2019 18:05:11 GMT"}, "X-Kong-Upstream-Latency": []string{"2"}, "X-Kong-Proxy-Latency": []string{"1"}},
		`{"Components":[{"ID":"x0c0s21b0n0","Type":"Node","State":"Ready","Flag":"OK","Enabled":true,"Role":"Compute","NID":4,"NetType":"Sling","Arch":"X86"}]}
`},
	{"https://frosty-sms.us.cray.com:30443/apis/smd/hsm/v1/Inventory/ComponentEndpoints?id=x0c0s21b0n0&type=CabinetPDUPowerConnector&type=CabinetPDUOutlet&type=Chassis&type=RouterModule&type=HSNBoard&type=ComputeModule&type=Node",
		"GET",
		``,
		map[string][]string{"Accept": []string{"application/json"}, "Content-Type": []string{"application/json"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Connection": []string{"keep-alive"}, "Date": []string{"Wed, 10 Apr 2019 18:05:11 GMT"}, "X-Kong-Upstream-Latency": []string{"3"}, "X-Kong-Proxy-Latency": []string{"0"}, "Via": []string{"kong/0.14.1"}, "Content-Type": []string{"application/json"}},
		`{"ComponentEndpoints":[{"ID":"x0c0s21b0n0","Type":"Node","RedfishType":"ComputerSystem","RedfishSubtype":"Physical","MACAddr":"a4:bf:01:2b:68:7c","UUID":"7f1d0db3-e077-11e7-ab21-a4bf012b687c","OdataID":"/redfish/v1/Systems/QSBP75001594","RedfishEndpointID":"x0c0s21b0","RedfishEndpointFQDN":"10.4.0.8","RedfishURL":"10.4.0.8/redfish/v1/Systems/QSBP75001594","ComponentEndpointType":"ComponentEndpointComputerSystem","RedfishSystemInfo":{"Name":"S2600BPB","Actions":{"#ComputerSystem.Reset":{"ResetType@Redfish.AllowableValues":["On","ForceOff","GracefulShutdown","GracefulRestart","ForceRestart","Nmi"],"target":"/redfish/v1/Systems/QSBP75001594/Actions/ComputerSystem.Reset"}},"EthernetNICInfo":[{"RedfishId":"","@odata.id":"","Description":"Missing interface 1, MAC computed via workaround","MACAddress":"a4:bf:01:2b:68:7c"},{"RedfishId":"","@odata.id":"","Description":"Missing interface 2, MAC computed via workaround","MACAddress":"a4:bf:01:2b:68:7d"}]}}]}
`},
	{"https://10.4.0.8/redfish/v1/Systems/QSBP75001594bad",
		"GET",
		``,
		map[string][]string{"Authorization": []string{"Basic cm9vdDppbml0aWFsMA=="}, "Accept": []string{"*/*"}},
		"404 Not Found", 404,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Strict-Transport-Security": []string{"max-age=31536000; includeSubdomains; preload"}, "X-Ua-Compatible": []string{"IE=11"}, "X-Xss-Protection": []string{"1; mode=block"}, "Content-Length": []string{"497"}, "Date": []string{"Wed, 10 Apr 2019 18:05:11 GMT"}, "X-Frame-Options": []string{"SAMEORIGIN"}, "Status": []string{"404"}, "Content-Type": []string{"application/json"}, "Server": []string{"lighttpd/1.4.45"}},
		`{"error":{"code":"Base.1.1.0.GeneralError","message":"A general error has occurred. See ExtendedInfo for more information.","@Message.ExtendedInfo":[{"@odata.type":"#Message.v1_0_4.Message","MessageId":"Base.1.1.0.ResourceMissingAtURI","Message":"The resource at the URI /redfish/v1/Systems/QSBP75001594bad was not found.","MessageArgs":["/redfish/v1/Systems/QSBP75001594bad"],"Severity":"Critical","Resolution":"Place a valid resource at the URI or correct the URI and resubmit the request."}]}}
`},
}
var StatusComponentNotFoundSSData = []sstorage.MockLookup{
	{
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname:    "x0c0s21b0n0",
				URL:      "10.4.0.8/redfish/v1/Systems/QSBP75001594",
				Username: "root",
				Password: "********",
			},
			Err: nil,
		},
	},
}

func TestStatusComponentNotFound(t *testing.T) {
	var testReq = testData{"/capmc/get_xname_status",
		"POST",
		`{"xnames":["x0c0s21b0n0"]}`,
		map[string][]string{"User-Agent": []string{"curl/7.37.0"}, "Accept": []string{"*/*"}, "Content-Length": []string{"26"}, "Content-Type": []string{"application/x-www-form-urlencoded"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}},
		`{"e":-1,"err_msg":"Errors encountered with 1/1 Xnames for Status","undefined":["x0c0s21b0n0"]}
`}
	runTest(t, StatusComponentNotFoundHSM, &testReq, &StatusComponentNotFoundReplayData, StatusComponentNotFoundSSData)
}

var StatusNodeNotFoundHSM = "https://frosty-sms.us.cray.com:30443/apis/smd"
var StatusNodeNotFoundReplayData = []testData{
	{"https://frosty-sms.us.cray.com:30443/apis/smd/hsm/v1/State/Components?nid=4&state=%21Empty&state=%21Populated&state=%21Unknown",
		"GET",
		``,
		map[string][]string{"Accept": []string{"application/json"}, "Content-Type": []string{"application/json"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}, "Content-Length": []string{"553"}, "Connection": []string{"keep-alive"}, "Date": []string{"Wed, 10 Apr 2019 18:05:28 GMT"}, "X-Kong-Upstream-Latency": []string{"2"}, "X-Kong-Proxy-Latency": []string{"1"}, "Via": []string{"kong/0.14.1"}},
		`{"Components":[]}
`},
	{"https://frosty-sms.us.cray.com:30443/apis/smd/hsm/v1/Inventory/ComponentEndpoints?id=x0c0s21b0n0",
		"GET",
		``,
		map[string][]string{"Accept": []string{"application/json"}, "Content-Type": []string{"application/json"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"X-Kong-Upstream-Latency": []string{"2"}, "X-Kong-Proxy-Latency": []string{"1"}, "Via": []string{"kong/0.14.1"}, "Content-Type": []string{"application/json"}, "Connection": []string{"keep-alive"}, "Date": []string{"Wed, 10 Apr 2019 18:05:28 GMT"}},
		`{"ComponentEndpoints":[{"ID":"x0c0s21b0n0","Type":"Node","RedfishType":"ComputerSystem","RedfishSubtype":"Physical","MACAddr":"a4:bf:01:2b:68:7c","UUID":"7f1d0db3-e077-11e7-ab21-a4bf012b687c","OdataID":"/redfish/v1/Systems/QSBP75001594","RedfishEndpointID":"x0c0s21b0","RedfishEndpointFQDN":"10.4.0.8","RedfishURL":"10.4.0.8/redfish/v1/Systems/QSBP75001594","ComponentEndpointType":"ComponentEndpointComputerSystem","RedfishSystemInfo":{"Name":"S2600BPB","Actions":{"#ComputerSystem.Reset":{"ResetType@Redfish.AllowableValues":["On","ForceOff","GracefulShutdown","GracefulRestart","ForceRestart","Nmi"],"target":"/redfish/v1/Systems/QSBP75001594/Actions/ComputerSystem.Reset"}},"EthernetNICInfo":[{"RedfishId":"","@odata.id":"","Description":"Missing interface 1, MAC computed via workaround","MACAddress":"a4:bf:01:2b:68:7c"},{"RedfishId":"","@odata.id":"","Description":"Missing interface 2, MAC computed via workaround","MACAddress":"a4:bf:01:2b:68:7d"}]}}]}
`},
	{"https://10.4.0.8/redfish/v1/Systems/QSBP75001594bad",
		"GET",
		``,
		map[string][]string{"Authorization": []string{"Basic cm9vdDppbml0aWFsMA=="}, "Accept": []string{"*/*"}},
		"404 Not Found", 404,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Status": []string{"404"}, "Content-Length": []string{"497"}, "Strict-Transport-Security": []string{"max-age=31536000; includeSubdomains; preload"}, "X-Ua-Compatible": []string{"IE=11"}, "X-Xss-Protection": []string{"1; mode=block"}, "Server": []string{"lighttpd/1.4.45"}, "X-Frame-Options": []string{"SAMEORIGIN"}, "Content-Type": []string{"application/json"}, "Date": []string{"Wed, 10 Apr 2019 18:05:28 GMT"}},
		`{"error":{"code":"Base.1.1.0.GeneralError","message":"A general error has occurred. See ExtendedInfo for more information.","@Message.ExtendedInfo":[{"@odata.type":"#Message.v1_0_4.Message","MessageId":"Base.1.1.0.ResourceMissingAtURI","Message":"The resource at the URI /redfish/v1/Systems/QSBP75001594bad was not found.","MessageArgs":["/redfish/v1/Systems/QSBP75001594bad"],"Severity":"Critical","Resolution":"Place a valid resource at the URI or correct the URI and resubmit the request."}]}}
`},
}
var StatusNodeNotFoundSSData = []sstorage.MockLookup{
	{
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname:    "x0c0s21b0n0",
				URL:      "10.4.0.8/redfish/v1/Systems/QSBP75001594",
				Username: "root",
				Password: "********",
			},
			Err: nil,
		},
	},
}

func TestStatusNodeNotFound(t *testing.T) {
	var testReq = testData{"/capmc/get_node_status",
		"POST",
		`{"nids":[4]}`,
		map[string][]string{"Accept": []string{"*/*"}, "Content-Length": []string{"12"}, "Content-Type": []string{"application/x-www-form-urlencoded"}, "User-Agent": []string{"curl/7.37.0"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}},
		`{"e":400,"err_msg":"nids not found: [4]"}
`}
	runTest(t, StatusNodeNotFoundHSM, &testReq, &StatusNodeNotFoundReplayData, StatusNodeNotFoundSSData)
}

var PowerNodeNotFoundHSM = "https://frosty-sms.us.cray.com:30443/apis/smd"
var PowerNodeNotFoundReplayData = []testData{
	{"https://frosty-sms.us.cray.com:30443/apis/smd/hsm/v1/State/Components?nid=4",
		"GET",
		``,
		map[string][]string{"Accept": []string{"application/json"}, "Content-Type": []string{"application/json"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Connection": []string{"keep-alive"}, "Date": []string{"Wed, 10 Apr 2019 18:14:11 GMT"}, "X-Kong-Upstream-Latency": []string{"2"}, "X-Kong-Proxy-Latency": []string{"1"}, "Via": []string{"kong/0.14.1"}, "Content-Type": []string{"application/json"}, "Content-Length": []string{"553"}},
		`{"Components":[{"ID":"x0c0s21b0n0","Type":"Node","State":"Ready","Flag":"OK","Enabled":true,"Role":"Compute","NID":4,"NetType":"Sling","Arch":"X86"}]}
`},
	{"https://frosty-sms.us.cray.com:30443/apis/smd/hsm/v1/Inventory/ComponentEndpoints?id=x0c0s21b0n0",
		"GET",
		``,
		map[string][]string{"Accept": []string{"application/json"}, "Content-Type": []string{"application/json"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"X-Kong-Proxy-Latency": []string{"0"}, "Via": []string{"kong/0.14.1"}, "Content-Type": []string{"application/json"}, "Connection": []string{"keep-alive"}, "Date": []string{"Wed, 10 Apr 2019 18:14:11 GMT"}, "X-Kong-Upstream-Latency": []string{"2"}},
		`{"ComponentEndpoints":[{"ID":"x0c0s21b0n0","Type":"Node","RedfishType":"ComputerSystem","RedfishSubtype":"Physical","MACAddr":"a4:bf:01:2b:68:7c","UUID":"7f1d0db3-e077-11e7-ab21-a4bf012b687c","OdataID":"/redfish/v1/Systems/QSBP75001595","RedfishEndpointID":"x0c0s21b0","RedfishEndpointFQDN":"10.4.0.8","RedfishURL":"10.4.0.8/redfish/v1/Systems/QSBP75001595","ComponentEndpointType":"ComponentEndpointComputerSystem","RedfishSystemInfo":{"Name":"S2600BPB","Actions":{"#ComputerSystem.Reset":{"ResetType@Redfish.AllowableValues":["On","ForceOff","GracefulShutdown","GracefulRestart","ForceRestart","Nmi"],"target":"/redfish/v1/Systems/QSBP75001595/Actions/ComputerSystem.Reset"}},"EthernetNICInfo":[{"RedfishId":"","@odata.id":"","Description":"Missing interface 1, MAC computed via workaround","MACAddress":"a4:bf:01:2b:68:7c"},{"RedfishId":"","@odata.id":"","Description":"Missing interface 2, MAC computed via workaround","MACAddress":"a4:bf:01:2b:68:7d"}]}}]}
`},
	{"https://10.4.0.8/redfish/v1/Systems/QSBP75001595/Actions/ComputerSystem.Reset",
		"POST",
		`{"ResetType": "On"}`,
		map[string][]string{"Authorization": []string{"Basic cm9vdDppbml0aWFsMA=="}, "Accept": []string{"*/*"}, "Content-Type": []string{"application/json"}},
		"404 Not Found", 404,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Length": []string{"549"}, "Date": []string{"Wed, 10 Apr 2019 18:14:11 GMT"}, "Server": []string{"lighttpd/1.4.45"}, "Strict-Transport-Security": []string{"max-age=31536000; includeSubdomains; preload"}, "X-Frame-Options": []string{"SAMEORIGIN"}, "X-Xss-Protection": []string{"1; mode=block"}, "Content-Type": []string{"application/json"}, "X-Ua-Compatible": []string{"IE=11"}, "Status": []string{"404"}},
		`{"error":{"code":"Base.1.1.0.GeneralError","message":"A general error has occurred. See ExtendedInfo for more information.","@Message.ExtendedInfo":[{"@odata.type":"#Message.v1_0_4.Message","MessageId":"Base.1.1.0.ResourceMissingAtURI","Message":"The resource at the URI /redfish/v1/Systems/QSBP75001595/Actions/ComputerSystem.Reset was not found.","MessageArgs":["/redfish/v1/Systems/QSBP75001595/Actions/ComputerSystem.Reset"],"Severity":"Critical","Resolution":"Place a valid resource at the URI or correct the URI and resubmit the request."}]}}
`},
	{"https://frosty-sms.us.cray.com:30443/apis/smd/hsm/v1/locks",
		"POST",
		`{"reason":"Performing a power Off","owner":"CAPMC","lifetime":180,"xnames":["x0c0s21b0n0"]}`,
		map[string][]string{"Content-Type": []string{"application/json"}, "Accept": []string{"application/json"}},
		"201 OK", 201,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}, "Date": []string{"Wed, 27 Feb 2019 04:37:13 GMT"}, "Content-Length": []string{"62"}},
		`{"uri":"/hsm/v1/locks/08c2e624-ded3-11e9-8a34-2a2ae2dbcce4"}]}
`},
}
var PowerNodeNotFoundSSData = []sstorage.MockLookup{
	{
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname:    "x0c0s21b0n0",
				URL:      "10.4.0.8/redfish/v1/Systems/QSBP75001595",
				Username: "root",
				Password: "********",
			},
			Err: nil,
		},
	},
}

func TestPowerNodeNotFound(t *testing.T) {
	var testReq = testData{"/capmc/node_on",
		"POST",
		`{"nids":[4]}`,
		map[string][]string{"Accept": []string{"*/*"}, "Content-Length": []string{"12"}, "Content-Type": []string{"application/x-www-form-urlencoded"}, "User-Agent": []string{"curl/7.37.0"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}},
		`{"e":-1,"err_msg":"Errors encountered with 1/1 NIDs","nids":[{"nid":4,"e":404,"err_msg":"NodeBMC 10.4.0.8 HTTP 404 Not Found, Redfish Error Message: A general error has occurred. See ExtendedInfo for more information. ExtendedInfo: Message: The resource at the URI /redfish/v1/Systems/QSBP75001595/Actions/ComputerSystem.Reset was not found. Resolution: Place a valid resource at the URI or correct the URI and resubmit the request."}]}
`}
	runTest(t, PowerNodeNotFoundHSM, &testReq, &PowerNodeNotFoundReplayData, PowerNodeNotFoundSSData)
}

var PowerNodeDisabledHSM = "https://frosty-sms.us.cray.com:30443/apis/smd"
var PowerNodeDisabledReplayData = []testData{
	{"https://frosty-sms.us.cray.com:30443/apis/smd/hsm/v1/State/Components?nid=4",
		"GET",
		``,
		map[string][]string{"Accept": []string{"application/json"}, "Content-Type": []string{"application/json"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Connection": []string{"keep-alive"}, "Date": []string{"Wed, 10 Apr 2019 18:14:11 GMT"}, "X-Kong-Upstream-Latency": []string{"2"}, "X-Kong-Proxy-Latency": []string{"1"}, "Via": []string{"kong/0.14.1"}, "Content-Type": []string{"application/json"}, "Content-Length": []string{"553"}},
		`{"Components":[{"ID":"x0c0s21b0n0","Type":"Node","State":"Ready","Flag":"OK","Enabled":false,"Role":"Compute","NID":4,"NetType":"Sling","Arch":"X86"}]}
`},
	{"https://frosty-sms.us.cray.com:30443/apis/smd/hsm/v1/Inventory/ComponentEndpoints?id=x0c0s21b0n0",
		"GET",
		``,
		map[string][]string{"Accept": []string{"application/json"}, "Content-Type": []string{"application/json"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"X-Kong-Proxy-Latency": []string{"0"}, "Via": []string{"kong/0.14.1"}, "Content-Type": []string{"application/json"}, "Connection": []string{"keep-alive"}, "Date": []string{"Wed, 10 Apr 2019 18:14:11 GMT"}, "X-Kong-Upstream-Latency": []string{"2"}},
		`{"ComponentEndpoints":[{"ID":"x0c0s21b0n0","Type":"Node","RedfishType":"ComputerSystem","RedfishSubtype":"Physical","MACAddr":"a4:bf:01:2b:68:7c","UUID":"7f1d0db3-e077-11e7-ab21-a4bf012b687c","OdataID":"/redfish/v1/Systems/QSBP75001595","RedfishEndpointID":"x0c0s21b0","RedfishEndpointFQDN":"10.4.0.8","RedfishURL":"10.4.0.8/redfish/v1/Systems/QSBP75001595","ComponentEndpointType":"ComponentEndpointComputerSystem","RedfishSystemInfo":{"Name":"S2600BPB","Actions":{"#ComputerSystem.Reset":{"ResetType@Redfish.AllowableValues":["On","ForceOff","GracefulShutdown","GracefulRestart","ForceRestart","Nmi"],"target":"/redfish/v1/Systems/QSBP75001595/Actions/ComputerSystem.Reset"}},"EthernetNICInfo":[{"RedfishId":"","@odata.id":"","Description":"Missing interface 1, MAC computed via workaround","MACAddress":"a4:bf:01:2b:68:7c"},{"RedfishId":"","@odata.id":"","Description":"Missing interface 2, MAC computed via workaround","MACAddress":"a4:bf:01:2b:68:7d"}]}}]}
`},
}
var PowerNodeDisabledSSData = []sstorage.MockLookup{
	{
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname:    "x0c0s21b0n0",
				URL:      "10.4.0.8/redfish/v1/Systems/QSBP75001595",
				Username: "root",
				Password: "********",
			},
			Err: nil,
		},
	},
}

func TestPowerNodeDisabled(t *testing.T) {
	var testReq = testData{"/capmc/node_on",
		"POST",
		`{"nids":[4]}`,
		map[string][]string{"Accept": []string{"*/*"}, "Content-Length": []string{"12"}, "Content-Type": []string{"application/x-www-form-urlencoded"}, "User-Agent": []string{"curl/7.37.0"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}},
		`{"e":400,"err_msg":"nodes disabled: [4]"}
`}
	runTest(t, PowerNodeDisabledHSM, &testReq, &PowerNodeDisabledReplayData, PowerNodeDisabledSSData)
}

var PowerComponentNotFoundHSM = "https://frosty-sms.us.cray.com:30443/apis/smd"
var PowerComponentNotFoundReplayData = []testData{
	{"https://frosty-sms.us.cray.com:30443/apis/smd/hsm/v1/State/Components?id=x0c0s21b0n0",
		"GET",
		``,
		map[string][]string{"Accept": []string{"application/json"}, "Content-Type": []string{"application/json"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Date": []string{"Wed, 10 Apr 2019 18:15:42 GMT"}, "X-Kong-Upstream-Latency": []string{"3"}, "X-Kong-Proxy-Latency": []string{"1"}, "Via": []string{"kong/0.14.1"}, "Content-Type": []string{"application/json"}, "Content-Length": []string{"1457"}, "Connection": []string{"keep-alive"}},
		`{"Components":[{"ID":"x0c0s21b0n0","Type":"Node","State":"Ready","Flag":"OK","Enabled":true,"Role":"Compute","NID":4,"NetType":"Sling","Arch":"X86"}]}
`},
	{"https://frosty-sms.us.cray.com:30443/apis/smd/hsm/v1/Inventory/ComponentEndpoints?id=x0c0s21b0n0",
		"GET",
		``,
		map[string][]string{"Accept": []string{"application/json"}, "Content-Type": []string{"application/json"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"X-Kong-Upstream-Latency": []string{"3"}, "X-Kong-Proxy-Latency": []string{"0"}, "Via": []string{"kong/0.14.1"}, "Content-Type": []string{"application/json"}, "Connection": []string{"keep-alive"}, "Date": []string{"Wed, 10 Apr 2019 18:15:42 GMT"}},
		`{"ComponentEndpoints":[{"ID":"x0c0s21b0n0","Type":"Node","RedfishType":"ComputerSystem","RedfishSubtype":"Physical","MACAddr":"a4:bf:01:2b:68:7c","UUID":"7f1d0db3-e077-11e7-ab21-a4bf012b687c","OdataID":"/redfish/v1/Systems/QSBP75001595","RedfishEndpointID":"x0c0s21b0","RedfishEndpointFQDN":"10.4.0.8","RedfishURL":"10.4.0.8/redfish/v1/Systems/QSBP75001595","ComponentEndpointType":"ComponentEndpointComputerSystem","RedfishSystemInfo":{"Name":"S2600BPB","Actions":{"#ComputerSystem.Reset":{"ResetType@Redfish.AllowableValues":["On","ForceOff","GracefulShutdown","GracefulRestart","ForceRestart","Nmi"],"target":"/redfish/v1/Systems/QSBP75001595/Actions/ComputerSystem.Reset"}},"EthernetNICInfo":[{"RedfishId":"","@odata.id":"","Description":"Missing interface 1, MAC computed via workaround","MACAddress":"a4:bf:01:2b:68:7c"},{"RedfishId":"","@odata.id":"","Description":"Missing interface 2, MAC computed via workaround","MACAddress":"a4:bf:01:2b:68:7d"}]}}]}
`},
	{"https://10.4.0.8/redfish/v1/Systems/QSBP75001595/Actions/ComputerSystem.Reset",
		"POST",
		`{"ResetType": "On"}`,
		map[string][]string{"Content-Type": []string{"application/json"}, "Authorization": []string{"Basic cm9vdDppbml0aWFsMA=="}, "Accept": []string{"*/*"}},
		"404 Not Found", 404,
		"HTTP/1.1", 1, 1,
		map[string][]string{"X-Frame-Options": []string{"SAMEORIGIN"}, "X-Xss-Protection": []string{"1; mode=block"}, "Status": []string{"404"}, "Content-Type": []string{"application/json"}, "Content-Length": []string{"549"}, "Server": []string{"lighttpd/1.4.45"}, "Strict-Transport-Security": []string{"max-age=31536000; includeSubdomains; preload"}, "X-Ua-Compatible": []string{"IE=11"}, "Date": []string{"Wed, 10 Apr 2019 18:15:42 GMT"}},
		`{"error":{"code":"Base.1.1.0.GeneralError","message":"A general error has occurred. See ExtendedInfo for more information.","@Message.ExtendedInfo":[{"@odata.type":"#Message.v1_0_4.Message","MessageId":"Base.1.1.0.ResourceMissingAtURI","Message":"The resource at the URI /redfish/v1/Systems/QSBP75001595/Actions/ComputerSystem.Reset was not found.","MessageArgs":["/redfish/v1/Systems/QSBP75001595/Actions/ComputerSystem.Reset"],"Severity":"Critical","Resolution":"Place a valid resource at the URI or correct the URI and resubmit the request."}]}}
`},
	{"https://frosty-sms.us.cray.com:30443/apis/smd/hsm/v1/locks",
		"POST",
		`{"reason":"Performing a power Off","owner":"CAPMC","lifetime":180,"xnames":["x0c0s21b0n0"]}`,
		map[string][]string{"Content-Type": []string{"application/json"}, "Accept": []string{"application/json"}},
		"201 OK", 201,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}, "Date": []string{"Wed, 27 Feb 2019 04:37:13 GMT"}, "Content-Length": []string{"62"}},
		`{"uri":"/hsm/v1/locks/08c2e624-ded3-11e9-8a34-2a2ae2dbcce4"}]}
`},
}
var PowerComponentNotFoundSSData = []sstorage.MockLookup{
	{
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname:    "x0c0s21b0n0",
				URL:      "10.4.0.8/redfish/v1/Systems/QSBP75001595",
				Username: "root",
				Password: "********",
			},
			Err: nil,
		},
	},
}

func TestPowerComponentNotFound(t *testing.T) {
	var testReq = testData{"/capmc/xname_on",
		"POST",
		`{"xnames":["x0c0s21b0n0"]}`,
		map[string][]string{"User-Agent": []string{"curl/7.37.0"}, "Accept": []string{"*/*"}, "Content-Length": []string{"26"}, "Content-Type": []string{"application/x-www-form-urlencoded"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}},
		`{"e":-1,"err_msg":"Errors encountered with 1/1 Xnames issued On","xnames":[{"xname":"x0c0s21b0n0","e":404,"err_msg":"NodeBMC 10.4.0.8 HTTP 404 Not Found, Redfish Error Message: A general error has occurred. See ExtendedInfo for more information. ExtendedInfo: Message: The resource at the URI /redfish/v1/Systems/QSBP75001595/Actions/ComputerSystem.Reset was not found. Resolution: Place a valid resource at the URI or correct the URI and resubmit the request."}]}
`}
	runTest(t, PowerComponentNotFoundHSM, &testReq, &PowerComponentNotFoundReplayData, PowerComponentNotFoundSSData)
}

var PowerComponentDisabledHSM = "https://frosty-sms.us.cray.com:30443/apis/smd"
var PowerComponentDisabledReplayData = []testData{
	{"https://frosty-sms.us.cray.com:30443/apis/smd/hsm/v1/State/Components?id=x0c0s21b0n0",
		"GET",
		``,
		map[string][]string{"Accept": []string{"application/json"}, "Content-Type": []string{"application/json"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Date": []string{"Wed, 10 Apr 2019 18:15:42 GMT"}, "X-Kong-Upstream-Latency": []string{"3"}, "X-Kong-Proxy-Latency": []string{"1"}, "Via": []string{"kong/0.14.1"}, "Content-Type": []string{"application/json"}, "Content-Length": []string{"1457"}, "Connection": []string{"keep-alive"}},
		`{"Components":[{"ID":"x0c0s21b0n0","Type":"Node","State":"Ready","Flag":"OK","Enabled":false,"Role":"Compute","NID":4,"NetType":"Sling","Arch":"X86"}]}
`},
	{"https://frosty-sms.us.cray.com:30443/apis/smd/hsm/v1/Inventory/ComponentEndpoints?id=x0c0s21b0n0",
		"GET",
		``,
		map[string][]string{"Accept": []string{"application/json"}, "Content-Type": []string{"application/json"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"X-Kong-Upstream-Latency": []string{"3"}, "X-Kong-Proxy-Latency": []string{"0"}, "Via": []string{"kong/0.14.1"}, "Content-Type": []string{"application/json"}, "Connection": []string{"keep-alive"}, "Date": []string{"Wed, 10 Apr 2019 18:15:42 GMT"}},
		`{"ComponentEndpoints":[{"ID":"x0c0s21b0n0","Type":"Node","RedfishType":"ComputerSystem","RedfishSubtype":"Physical","MACAddr":"a4:bf:01:2b:68:7c","UUID":"7f1d0db3-e077-11e7-ab21-a4bf012b687c","OdataID":"/redfish/v1/Systems/QSBP75001595","RedfishEndpointID":"x0c0s21b0","RedfishEndpointFQDN":"10.4.0.8","RedfishURL":"10.4.0.8/redfish/v1/Systems/QSBP75001595","ComponentEndpointType":"ComponentEndpointComputerSystem","RedfishSystemInfo":{"Name":"S2600BPB","Actions":{"#ComputerSystem.Reset":{"ResetType@Redfish.AllowableValues":["On","ForceOff","GracefulShutdown","GracefulRestart","ForceRestart","Nmi"],"target":"/redfish/v1/Systems/QSBP75001595/Actions/ComputerSystem.Reset"}},"EthernetNICInfo":[{"RedfishId":"","@odata.id":"","Description":"Missing interface 1, MAC computed via workaround","MACAddress":"a4:bf:01:2b:68:7c"},{"RedfishId":"","@odata.id":"","Description":"Missing interface 2, MAC computed via workaround","MACAddress":"a4:bf:01:2b:68:7d"}]}}]}
`},
	{"https://10.4.0.8/redfish/v1/Systems/QSBP75001595/Actions/ComputerSystem.Reset",
		"POST",
		`{"ResetType": "On"}`,
		map[string][]string{"Content-Type": []string{"application/json"}, "Authorization": []string{"Basic cm9vdDppbml0aWFsMA=="}, "Accept": []string{"*/*"}},
		"404 Not Found", 404,
		"HTTP/1.1", 1, 1,
		map[string][]string{"X-Frame-Options": []string{"SAMEORIGIN"}, "X-Xss-Protection": []string{"1; mode=block"}, "Status": []string{"404"}, "Content-Type": []string{"application/json"}, "Content-Length": []string{"549"}, "Server": []string{"lighttpd/1.4.45"}, "Strict-Transport-Security": []string{"max-age=31536000; includeSubdomains; preload"}, "X-Ua-Compatible": []string{"IE=11"}, "Date": []string{"Wed, 10 Apr 2019 18:15:42 GMT"}},
		`{"error":{"code":"Base.1.1.0.GeneralError","message":"A general error has occurred. See ExtendedInfo for more information.","@Message.ExtendedInfo":[{"@odata.type":"#Message.v1_0_4.Message","MessageId":"Base.1.1.0.ResourceMissingAtURI","Message":"The resource at the URI /redfish/v1/Systems/QSBP75001595/Actions/ComputerSystem.Reset was not found.","MessageArgs":["/redfish/v1/Systems/QSBP75001595/Actions/ComputerSystem.Reset"],"Severity":"Critical","Resolution":"Place a valid resource at the URI or correct the URI and resubmit the request."}]}}
`},
	{"https://frosty-sms.us.cray.com:30443/apis/smd/hsm/v1/locks",
		"POST",
		`{"reason":"Performing a power Off","owner":"CAPMC","lifetime":180,"xnames":["x0c0s21b0n0"]}`,
		map[string][]string{"Content-Type": []string{"application/json"}, "Accept": []string{"application/json"}},
		"201 OK", 201,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}, "Date": []string{"Wed, 27 Feb 2019 04:37:13 GMT"}, "Content-Length": []string{"62"}},
		`{"uri":"/hsm/v1/locks/08c2e624-ded3-11e9-8a34-2a2ae2dbcce4"}]}
`},
}
var PowerComponentDisabledSSData = []sstorage.MockLookup{
	{
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname:    "x0c0s21b0n0",
				URL:      "10.4.0.8/redfish/v1/Systems/QSBP75001595",
				Username: "root",
				Password: "********",
			},
			Err: nil,
		},
	},
}

func TestPowerComponentDisabled(t *testing.T) {
	var testReq = testData{"/capmc/xname_on",
		"POST",
		`{"xnames":["x0c0s21b0n0"]}`,
		map[string][]string{"User-Agent": []string{"curl/7.37.0"}, "Accept": []string{"*/*"}, "Content-Length": []string{"26"}, "Content-Type": []string{"application/x-www-form-urlencoded"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}},
		`{"e":400,"err_msg":"components disabled: [x0c0s21b0n0]"}
`}
	runTest(t, PowerComponentDisabledHSM, &testReq, &PowerComponentDisabledReplayData, PowerComponentDisabledSSData)
}

var PowerComponentNotAuthorizedHSM = "https://frosty-sms.us.cray.com:30443/apis/smd"
var PowerComponentNotAuthorizedReplayData = []testData{
	{"https://frosty-sms.us.cray.com:30443/apis/smd/hsm/v1/State/Components?id=x0c0s21b0n0",
		"GET",
		``,
		map[string][]string{"Accept": []string{"application/json"}, "Content-Type": []string{"application/json"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Date": []string{"Wed, 10 Apr 2019 18:47:06 GMT"}, "X-Kong-Upstream-Latency": []string{"3"}, "X-Kong-Proxy-Latency": []string{"1"}, "Via": []string{"kong/0.14.1"}, "Content-Type": []string{"application/json"}, "Content-Length": []string{"1457"}, "Connection": []string{"keep-alive"}},
		`{"Components":[{"ID":"x0c0s21b0n0","Type":"Node","State":"Ready","Flag":"OK","Enabled":true,"Role":"Compute","NID":4,"NetType":"Sling","Arch":"X86"}]}
`},
	{"https://frosty-sms.us.cray.com:30443/apis/smd/hsm/v1/Inventory/ComponentEndpoints?id=x0c0s21b0n0",
		"GET",
		``,
		map[string][]string{"Accept": []string{"application/json"}, "Content-Type": []string{"application/json"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Connection": []string{"keep-alive"}, "Date": []string{"Wed, 10 Apr 2019 18:47:06 GMT"}, "X-Kong-Upstream-Latency": []string{"3"}, "X-Kong-Proxy-Latency": []string{"0"}, "Via": []string{"kong/0.14.1"}, "Content-Type": []string{"application/json"}},
		`{"ComponentEndpoints":[{"ID":"x0c0s21b0n0","Type":"Node","RedfishType":"ComputerSystem","RedfishSubtype":"Physical","MACAddr":"a4:bf:01:2b:68:7c","UUID":"7f1d0db3-e077-11e7-ab21-a4bf012b687c","OdataID":"/redfish/v1/Systems/QSBP75001594","RedfishEndpointID":"x0c0s21b0","RedfishEndpointFQDN":"10.4.0.8","RedfishURL":"10.4.0.8/redfish/v1/Systems/QSBP75001594","ComponentEndpointType":"ComponentEndpointComputerSystem","RedfishSystemInfo":{"Name":"S2600BPB","Actions":{"#ComputerSystem.Reset":{"ResetType@Redfish.AllowableValues":["On","ForceOff","GracefulShutdown","GracefulRestart","ForceRestart","Nmi"],"target":"/redfish/v1/Systems/QSBP75001594/Actions/ComputerSystem.Reset"}},"EthernetNICInfo":[{"RedfishId":"","@odata.id":"","Description":"Missing interface 1, MAC computed via workaround","MACAddress":"a4:bf:01:2b:68:7c"},{"RedfishId":"","@odata.id":"","Description":"Missing interface 2, MAC computed via workaround","MACAddress":"a4:bf:01:2b:68:7d"}]}}]}
`},
	{"https://10.4.0.8/redfish/v1/Systems/QSBP75001594/Actions/ComputerSystem.Reset",
		"POST",
		`{"ResetType": "On"}`,
		map[string][]string{"Accept": []string{"*/*"}, "Content-Type": []string{"application/json"}, "Authorization": []string{"Basic cm9vdDpiYWRwYXNz"}},
		"401 Unauthorized", 401,
		"HTTP/1.1", 1, 1,
		map[string][]string{"X-Ua-Compatible": []string{"IE=11"}, "X-Frame-Options": []string{"SAMEORIGIN"}, "X-Xss-Protection": []string{"1; mode=block"}, "Date": []string{"Wed, 10 Apr 2019 18:47:06 GMT"}, "Server": []string{"lighttpd/1.4.45"}, "Strict-Transport-Security": []string{"max-age=31536000; includeSubdomains; preload"}, "Status": []string{"401"}, "Www-Authenticate": []string{"Basic realm=\"Redfish\""}, "Content-Type": []string{"application/json"}, "Content-Length": []string{"877"}},
		`{"error":{"code":"Base.1.1.0.GeneralError","message":"A general error has occurred. See ExtendedInfo for more information.","@Message.ExtendedInfo":[{"@odata.type":"#Message.v1_0_4.Message","MessageId":"Base.1.1.0.NoValidSession","Message":"There is no valid session established with the implementation.","Severity":"Critical","Resolution":"Establish as session before attempting any operations."},{"@odata.type":"#Message.v1_0_4.Message","MessageId":"Base.1.1.0.ResourceAtUriUnauthorized","Message":"While accessing the resource at /redfish/v1/Systems/QSBP75001594/Actions/ComputerSystem.Reset, the service received an authorization error failed.","MessageArgs":["/redfish/v1/Systems/QSBP75001594/Actions/ComputerSystem.Reset","failed"],"Severity":"Critical","Resolution":"Ensure that the appropriate access is provided for the service in order for it to access the URI."}]}}
`},
	{"https://frosty-sms.us.cray.com:30443/apis/smd/hsm/v1/locks",
		"POST",
		`{"reason":"Performing a power Off","owner":"CAPMC","lifetime":180,"xnames":["x0c0s21b0n0"]}`,
		map[string][]string{"Content-Type": []string{"application/json"}, "Accept": []string{"application/json"}},
		"201 OK", 201,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}, "Date": []string{"Wed, 27 Feb 2019 04:37:13 GMT"}, "Content-Length": []string{"62"}},
		`{"uri":"/hsm/v1/locks/08c2e624-ded3-11e9-8a34-2a2ae2dbcce4"}]}
`},
}
var PowerComponentNotAuthorizedSSData = []sstorage.MockLookup{
	{
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname:    "x0c0s21b0n0",
				URL:      "10.4.0.8/redfish/v1/Systems/QSBP75001594",
				Username: "root",
				Password: "********",
			},
			Err: nil,
		},
	},
}

func TestPowerComponentNotAuthorized(t *testing.T) {
	var testReq = testData{"/capmc/xname_on",
		"POST",
		`{"xnames":["x0c0s21b0n0"]}`,
		map[string][]string{"User-Agent": []string{"curl/7.37.0"}, "Accept": []string{"*/*"}, "Content-Length": []string{"26"}, "Content-Type": []string{"application/x-www-form-urlencoded"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}},
		`{"e":-1,"err_msg":"Errors encountered with 1/1 Xnames issued On","xnames":[{"xname":"x0c0s21b0n0","e":401,"err_msg":"NodeBMC 10.4.0.8 HTTP 401 Unauthorized, Redfish Error Message: A general error has occurred. See ExtendedInfo for more information. ExtendedInfo: Message: There is no valid session established with the implementation. Resolution: Establish as session before attempting any operations. Message: While accessing the resource at /redfish/v1/Systems/QSBP75001594/Actions/ComputerSystem.Reset, the service received an authorization error failed. Resolution: Ensure that the appropriate access is provided for the service in order for it to access the URI."}]}
`}
	runTest(t, PowerComponentNotAuthorizedHSM, &testReq, &PowerComponentNotAuthorizedReplayData, PowerComponentNotAuthorizedSSData)
}

var PowerNodeNotAuthorizedHSM = "https://frosty-sms.us.cray.com:30443/apis/smd"
var PowerNodeNotAuthorizedReplayData = []testData{
	{"https://frosty-sms.us.cray.com:30443/apis/smd/hsm/v1/State/Components?nid=4",
		"GET",
		``,
		map[string][]string{"Accept": []string{"application/json"}, "Content-Type": []string{"application/json"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Date": []string{"Wed, 10 Apr 2019 18:47:22 GMT"}, "X-Kong-Upstream-Latency": []string{"2"}, "X-Kong-Proxy-Latency": []string{"1"}, "Via": []string{"kong/0.14.1"}, "Content-Type": []string{"application/json"}, "Content-Length": []string{"553"}, "Connection": []string{"keep-alive"}},
		`{"Components":[{"ID":"x0c0s21b0n0","Type":"Node","State":"Ready","Flag":"OK","Enabled":true,"Role":"Compute","NID":4,"NetType":"Sling","Arch":"X86"}]}
`},
	{"https://frosty-sms.us.cray.com:30443/apis/smd/hsm/v1/Inventory/ComponentEndpoints?id=x0c0s21b0n0",
		"GET",
		``,
		map[string][]string{"Accept": []string{"application/json"}, "Content-Type": []string{"application/json"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}, "Connection": []string{"keep-alive"}, "Date": []string{"Wed, 10 Apr 2019 18:47:22 GMT"}, "X-Kong-Upstream-Latency": []string{"2"}, "X-Kong-Proxy-Latency": []string{"1"}, "Via": []string{"kong/0.14.1"}},
		`{"ComponentEndpoints":[{"ID":"x0c0s21b0n0","Type":"Node","RedfishType":"ComputerSystem","RedfishSubtype":"Physical","MACAddr":"a4:bf:01:2b:68:7c","UUID":"7f1d0db3-e077-11e7-ab21-a4bf012b687c","OdataID":"/redfish/v1/Systems/QSBP75001594","RedfishEndpointID":"x0c0s21b0","RedfishEndpointFQDN":"10.4.0.8","RedfishURL":"10.4.0.8/redfish/v1/Systems/QSBP75001594","ComponentEndpointType":"ComponentEndpointComputerSystem","RedfishSystemInfo":{"Name":"S2600BPB","Actions":{"#ComputerSystem.Reset":{"ResetType@Redfish.AllowableValues":["On","ForceOff","GracefulShutdown","GracefulRestart","ForceRestart","Nmi"],"target":"/redfish/v1/Systems/QSBP75001594/Actions/ComputerSystem.Reset"}},"EthernetNICInfo":[{"RedfishId":"","@odata.id":"","Description":"Missing interface 1, MAC computed via workaround","MACAddress":"a4:bf:01:2b:68:7c"},{"RedfishId":"","@odata.id":"","Description":"Missing interface 2, MAC computed via workaround","MACAddress":"a4:bf:01:2b:68:7d"}]}}]}
`},
	{"https://frosty-sms.us.cray.com:30443/apis/smd/hsm/v1/Inventory/RedfishEndpoints?id=x0c0s21b0",
		"GET",
		``,
		map[string][]string{"Accept": []string{"application/json"}, "Content-Type": []string{"application/json"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Via": []string{"kong/0.14.1"}, "Content-Type": []string{"application/json"}, "Content-Length": []string{"1491"}, "Connection": []string{"keep-alive"}, "Date": []string{"Wed, 10 Apr 2019 18:47:22 GMT"}, "X-Kong-Upstream-Latency": []string{"2"}, "X-Kong-Proxy-Latency": []string{"0"}},
		`{"RedfishEndpoints":[{"ID":"x0c0s21b0","Type":"NodeBMC","Hostname":"10.4.0.8","Domain":"","FQDN":"10.4.0.8","Enabled":true,"UUID":"d146230d-8b78-4da1-b06a-7363d526757f","User":"root","Password":"********","MACAddr":"a4bf012b6880","RediscoverOnUpdate":true,"DiscoveryInfo":{"LastDiscoveryAttempt":"2019-04-04T21:09:59.136106Z","LastDiscoveryStatus":"DiscoverOK","RedfishVersion":"1.1.0"}}]}
`},
	{"https://10.4.0.8/redfish/v1/Systems/QSBP75001594/Actions/ComputerSystem.Reset",
		"POST",
		`{"ResetType": "On"}`,
		map[string][]string{"Content-Type": []string{"application/json"}, "Authorization": []string{"Basic cm9vdDpiYWRwYXNz"}, "Accept": []string{"*/*"}},
		"401 Unauthorized", 401,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Status": []string{"401"}, "Www-Authenticate": []string{"Basic realm=\"Redfish\""}, "Content-Type": []string{"application/json"}, "Content-Length": []string{"877"}, "Server": []string{"lighttpd/1.4.45"}, "Strict-Transport-Security": []string{"max-age=31536000; includeSubdomains; preload"}, "X-Ua-Compatible": []string{"IE=11"}, "X-Frame-Options": []string{"SAMEORIGIN"}, "X-Xss-Protection": []string{"1; mode=block"}, "Date": []string{"Wed, 10 Apr 2019 18:47:22 GMT"}},
		`{"error":{"code":"Base.1.1.0.GeneralError","message":"A general error has occurred. See ExtendedInfo for more information.","@Message.ExtendedInfo":[{"@odata.type":"#Message.v1_0_4.Message","MessageId":"Base.1.1.0.NoValidSession","Message":"There is no valid session established with the implementation.","Severity":"Critical","Resolution":"Establish as session before attempting any operations."},{"@odata.type":"#Message.v1_0_4.Message","MessageId":"Base.1.1.0.ResourceAtUriUnauthorized","Message":"While accessing the resource at /redfish/v1/Systems/QSBP75001594/Actions/ComputerSystem.Reset, the service received an authorization error failed.","MessageArgs":["/redfish/v1/Systems/QSBP75001594/Actions/ComputerSystem.Reset","failed"],"Severity":"Critical","Resolution":"Ensure that the appropriate access is provided for the service in order for it to access the URI."}]}}
`},
	{"https://frosty-sms.us.cray.com:30443/apis/smd/hsm/v1/locks",
		"POST",
		`{"reason":"Performing a power Off","owner":"CAPMC","lifetime":180,"xnames":["x0c0s21b0n0"]}`,
		map[string][]string{"Content-Type": []string{"application/json"}, "Accept": []string{"application/json"}},
		"201 OK", 201,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}, "Date": []string{"Wed, 27 Feb 2019 04:37:13 GMT"}, "Content-Length": []string{"62"}},
		`{"uri":"/hsm/v1/locks/08c2e624-ded3-11e9-8a34-2a2ae2dbcce4"}]}
`},
}
var PowerNodeNotAuthorizedSSData = []sstorage.MockLookup{
	{
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname:    "x0c0s21b0n0",
				URL:      "10.4.0.8/redfish/v1/Systems/QSBP75001594",
				Username: "root",
				Password: "********",
			},
			Err: nil,
		},
	},
}

func TestPowerNodeNotAuthorized(t *testing.T) {
	var testReq = testData{"/capmc/node_on",
		"POST",
		`{"nids":[4]}`,
		map[string][]string{"Accept": []string{"*/*"}, "Content-Length": []string{"12"}, "Content-Type": []string{"application/x-www-form-urlencoded"}, "User-Agent": []string{"curl/7.37.0"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}},
		`{"e":-1,"err_msg":"Errors encountered with 1/1 NIDs","nids":[{"nid":4,"e":401,"err_msg":"NodeBMC 10.4.0.8 HTTP 401 Unauthorized, Redfish Error Message: A general error has occurred. See ExtendedInfo for more information. ExtendedInfo: Message: There is no valid session established with the implementation. Resolution: Establish as session before attempting any operations. Message: While accessing the resource at /redfish/v1/Systems/QSBP75001594/Actions/ComputerSystem.Reset, the service received an authorization error failed. Resolution: Ensure that the appropriate access is provided for the service in order for it to access the URI."}]}
`}
	runTest(t, PowerNodeNotAuthorizedHSM, &testReq, &PowerNodeNotAuthorizedReplayData, PowerNodeNotAuthorizedSSData)
}

var XnameOffRecursiveHSM = "https://slice-sms.us.cray.com:30443/apis/smd"
var XnameOffRecursiveReplayData = []testData{
	{"https://slice-sms.us.cray.com:30443/apis/smd/hsm/v1/State/Components/Query/x0c0s28?state=%21Empty&type=chassis&type=cabinetpdu&type=routermodule&type=hsnboard&type=computemodule&type=node",
		"GET",
		``,
		map[string][]string{"Accept": []string{"application/json"}, "Content-Type": []string{"application/json"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Date": []string{"Tue, 14 May 2019 11:38:17 GMT"}, "X-Kong-Upstream-Latency": []string{"3"}, "X-Kong-Proxy-Latency": []string{"1"}, "Via": []string{"kong/0.14.1"}, "Content-Type": []string{"application/json"}, "Content-Length": []string{"156"}, "Connection": []string{"keep-alive"}},
		`{"Components":[{"ID":"x0c0s28b0n0","Type":"Node","State":"Standby","Flag":"Alert","Enabled":true,"Role":"Compute","NID":1,"NetType":"Sling","Arch":"X86"}]}
`},
	{"https://slice-sms.us.cray.com:30443/apis/smd/hsm/v1/State/Components?id=x0c0s28b0n0",
		"GET",
		``,
		map[string][]string{"Accept": []string{"application/json"}, "Content-Type": []string{"application/json"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Length": []string{"1465"}, "Connection": []string{"keep-alive"}, "Date": []string{"Tue, 14 May 2019 11:38:17 GMT"}, "X-Kong-Upstream-Latency": []string{"2"}, "X-Kong-Proxy-Latency": []string{"1"}, "Via": []string{"kong/0.14.1"}, "Content-Type": []string{"application/json"}},
		`{"Components":[{"ID":"x0c0s28b0n0","Type":"Node","State":"Standby","Flag":"Alert","Enabled":true,"Role":"Compute","NID":1,"NetType":"Sling","Arch":"X86"}]}
`},
	{"https://slice-sms.us.cray.com:30443/apis/smd/hsm/v1/Inventory/ComponentEndpoints?id=x0c0s28b0n0",
		"GET",
		``,
		map[string][]string{"Accept": []string{"application/json"}, "Content-Type": []string{"application/json"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Via": []string{"kong/0.14.1"}, "Content-Type": []string{"application/json"}, "Connection": []string{"keep-alive"}, "Date": []string{"Tue, 14 May 2019 11:38:17 GMT"}, "X-Kong-Upstream-Latency": []string{"3"}, "X-Kong-Proxy-Latency": []string{"0"}},
		`{"ComponentEndpoints":[{"ID":"x0c0s28b0n0","Type":"Node","RedfishType":"ComputerSystem","RedfishSubtype":"Physical","MACAddr":"a4:bf:01:28:92:0b","UUID":"3504bf56-ae5b-11e7-ab21-a4bf0128920b","OdataID":"/redfish/v1/Systems/QSBP74100092","RedfishEndpointID":"x0c0s28b0","RedfishEndpointFQDN":"10.4.0.5","RedfishURL":"10.4.0.5/redfish/v1/Systems/QSBP74100092","ComponentEndpointType":"ComponentEndpointComputerSystem","RedfishSystemInfo":{"Name":"S2600BPB","Actions":{"#ComputerSystem.Reset":{"ResetType@Redfish.AllowableValues":["On","ForceOff","GracefulShutdown","GracefulRestart","ForceRestart","Nmi"],"target":"/redfish/v1/Systems/QSBP74100092/Actions/ComputerSystem.Reset"}},"EthernetNICInfo":[{"RedfishId":"","@odata.id":"","Description":"Missing interface 1, MAC computed via workaround","MACAddress":"a4:bf:01:28:92:0b"},{"RedfishId":"","@odata.id":"","Description":"Missing interface 2, MAC computed via workaround","MACAddress":"a4:bf:01:28:92:0c"}]}}]}
`},
	{"https://10.4.0.5/redfish/v1/Systems/QSBP74100092/Actions/ComputerSystem.Reset",
		"POST",
		`{"ResetType": "ForceOff"}`,
		map[string][]string{"Authorization": []string{"Basic cm9vdDppbml0aWFsMA=="}, "Accept": []string{"*/*"}, "Content-Type": []string{"application/json"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}, "Content-Length": []string{"154"}, "Date": []string{"Tue, 14 May 2019 11:38:25 GMT"}, "Server": []string{"lighttpd/1.4.45"}, "Strict-Transport-Security": []string{"max-age=31536000; includeSubdomains; preload"}, "X-Ua-Compatible": []string{"IE=11"}, "X-Frame-Options": []string{"SAMEORIGIN"}, "X-Xss-Protection": []string{"1; mode=block"}},
		`{"@odata.type":"#Message.v1_0_4.Message","MessageId":"Base.1.1.0.Success","Message":"Successfully Completed Request","Severity":"OK","Resolution":"None"}
`},
	{"https://slice-sms.us.cray.com:30443/apis/smd/hsm/v1/locks",
		"POST",
		`{"reason":"Performing a power Off","owner":"CAPMC","lifetime":180,"xnames":["x0c0s28b0n0"]}`,
		map[string][]string{"Content-Type": []string{"application/json"}, "Accept": []string{"application/json"}},
		"201 OK", 201,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}, "Date": []string{"Wed, 27 Feb 2019 04:37:13 GMT"}, "Content-Length": []string{"62"}},
		`{"uri":"/hsm/v1/locks/08c2e624-ded3-11e9-8a34-2a2ae2dbcce4"}]}
`},
}
var XnameOffRecursiveSSData = []sstorage.MockLookup{
	{
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname:    "x0c0s28b0n0",
				URL:      "10.4.0.5/redfish/v1/Systems/QSBP74100092",
				Username: "root",
				Password: "********",
			},
			Err: nil,
		},
	},
}

func TestXnameOffRecursive(t *testing.T) {
	var testReq = testData{"/capmc/xname_off",
		"POST",
		`{ "force": false, "recursive": true, "reason": "", "xnames": [ "x0c0s28" ] }`,
		map[string][]string{"Content-Type": []string{"application/json"}, "Accept": []string{"application/json"}, "Content-Length": []string{"76"}, "User-Agent": []string{"curl/7.37.0"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}},
		`{"e":0,"err_msg":""}
`}
	if testing.Short() {
		t.Skip("skipping test XnameOffRecursive in short mode")
	}
	runTest(t, XnameOffRecursiveHSM, &testReq, &XnameOffRecursiveReplayData, XnameOffRecursiveSSData)
}

var XnameOnRecursiveHSM = "https://slice-sms.us.cray.com:30443/apis/smd"
var XnameOnRecursiveReplayData = []testData{
	{"https://slice-sms.us.cray.com:30443/apis/smd/hsm/v1/State/Components/Query/x0c0s28?state=%21Empty&type=chassis&type=cabinetpdu&type=routermodule&type=hsnboard&type=computemodule&type=node",
		"GET",
		``,
		map[string][]string{"Accept": []string{"application/json"}, "Content-Type": []string{"application/json"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Length": []string{"156"}, "Connection": []string{"keep-alive"}, "Date": []string{"Tue, 14 May 2019 11:44:45 GMT"}, "X-Kong-Upstream-Latency": []string{"2"}, "X-Kong-Proxy-Latency": []string{"1"}, "Via": []string{"kong/0.14.1"}, "Content-Type": []string{"application/json"}},
		`{"Components":[{"ID":"x0c0s28b0n0","Type":"Node","State":"Standby","Flag":"Alert","Enabled":true,"Role":"Compute","NID":1,"NetType":"Sling","Arch":"X86"}]}
`},
	{"https://slice-sms.us.cray.com:30443/apis/smd/hsm/v1/State/Components?id=x0c0s28b0n0",
		"GET",
		``,
		map[string][]string{"Accept": []string{"application/json"}, "Content-Type": []string{"application/json"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}, "Content-Length": []string{"1465"}, "Connection": []string{"keep-alive"}, "Date": []string{"Tue, 14 May 2019 11:44:45 GMT"}, "X-Kong-Upstream-Latency": []string{"2"}, "X-Kong-Proxy-Latency": []string{"1"}, "Via": []string{"kong/0.14.1"}},
		`{"Components":[{"ID":"x0c0s28b0n0","Type":"Node","State":"Standby","Flag":"Alert","Enabled":true,"Role":"Compute","NID":1,"NetType":"Sling","Arch":"X86"}]}
`},
	{"https://slice-sms.us.cray.com:30443/apis/smd/hsm/v1/Inventory/ComponentEndpoints?id=x0c0s28b0n0",
		"GET",
		``,
		map[string][]string{"Accept": []string{"application/json"}, "Content-Type": []string{"application/json"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Connection": []string{"keep-alive"}, "Date": []string{"Tue, 14 May 2019 11:44:45 GMT"}, "X-Kong-Upstream-Latency": []string{"3"}, "X-Kong-Proxy-Latency": []string{"0"}, "Via": []string{"kong/0.14.1"}, "Content-Type": []string{"application/json"}},
		`{"ComponentEndpoints":[{"ID":"x0c0s28b0n0","Type":"Node","RedfishType":"ComputerSystem","RedfishSubtype":"Physical","MACAddr":"a4:bf:01:28:92:0b","UUID":"3504bf56-ae5b-11e7-ab21-a4bf0128920b","OdataID":"/redfish/v1/Systems/QSBP74100092","RedfishEndpointID":"x0c0s28b0","RedfishEndpointFQDN":"10.4.0.5","RedfishURL":"10.4.0.5/redfish/v1/Systems/QSBP74100092","ComponentEndpointType":"ComponentEndpointComputerSystem","RedfishSystemInfo":{"Name":"S2600BPB","Actions":{"#ComputerSystem.Reset":{"ResetType@Redfish.AllowableValues":["On","ForceOff","GracefulShutdown","GracefulRestart","ForceRestart","Nmi"],"target":"/redfish/v1/Systems/QSBP74100092/Actions/ComputerSystem.Reset"}},"EthernetNICInfo":[{"RedfishId":"","@odata.id":"","Description":"Missing interface 1, MAC computed via workaround","MACAddress":"a4:bf:01:28:92:0b"},{"RedfishId":"","@odata.id":"","Description":"Missing interface 2, MAC computed via workaround","MACAddress":"a4:bf:01:28:92:0c"}]}}]}
`},
	{"https://10.4.0.5/redfish/v1/Systems/QSBP74100092/Actions/ComputerSystem.Reset",
		"POST",
		`{"ResetType": "On"}`,
		map[string][]string{"Authorization": []string{"Basic cm9vdDppbml0aWFsMA=="}, "Accept": []string{"*/*"}, "Content-Type": []string{"application/json"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Strict-Transport-Security": []string{"max-age=31536000; includeSubdomains; preload"}, "X-Ua-Compatible": []string{"IE=11"}, "X-Frame-Options": []string{"SAMEORIGIN"}, "X-Xss-Protection": []string{"1; mode=block"}, "Content-Type": []string{"application/json"}, "Content-Length": []string{"154"}, "Date": []string{"Tue, 14 May 2019 11:44:54 GMT"}, "Server": []string{"lighttpd/1.4.45"}},
		`{"@odata.type":"#Message.v1_0_4.Message","MessageId":"Base.1.1.0.Success","Message":"Successfully Completed Request","Severity":"OK","Resolution":"None"}
`},
	{"https://slice-sms.us.cray.com:30443/apis/smd/hsm/v1/locks",
		"POST",
		`{"reason":"Performing a power Off","owner":"CAPMC","lifetime":180,"xnames":["x0c0s28b0n0"]}`,
		map[string][]string{"Content-Type": []string{"application/json"}, "Accept": []string{"application/json"}},
		"201 OK", 201,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}, "Date": []string{"Wed, 27 Feb 2019 04:37:13 GMT"}, "Content-Length": []string{"62"}},
		`{"uri":"/hsm/v1/locks/08c2e624-ded3-11e9-8a34-2a2ae2dbcce4"}]}
`},
}
var XnameOnRecursiveSSData = []sstorage.MockLookup{
	{
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname:    "x0c0s28b0n0",
				URL:      "10.4.0.5/redfish/v1/Systems/QSBP74100092",
				Username: "root",
				Password: "********",
			},
			Err: nil,
		},
	},
}

func TestXnameOnRecursive(t *testing.T) {
	var testReq = testData{"/capmc/xname_on",
		"POST",
		`{ "force": false, "recursive": true, "reason": "", "xnames": [ "x0c0s28" ] }`,
		map[string][]string{"User-Agent": []string{"curl/7.37.0"}, "Content-Type": []string{"application/json"}, "Accept": []string{"application/json"}, "Content-Length": []string{"76"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}},
		`{"e":0,"err_msg":""}
`}
	runTest(t, XnameOnRecursiveHSM, &testReq, &XnameOnRecursiveReplayData, XnameOnRecursiveSSData)
}

var XnameOffRecursiveChassisHSM = "https://slice-sms.us.cray.com:30443/apis/smd"
var XnameOffRecursiveChassisReplayData = []testData{
	{"https://slice-sms.us.cray.com:30443/apis/smd/hsm/v1/State/Components/Query/x0c0?state=%21Empty&type=chassis&type=cabinetpdu&type=routermodule&type=hsnboard&type=computemodule&type=node",
		"GET",
		``,
		map[string][]string{"Accept": []string{"application/json"}, "Content-Type": []string{"application/json"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"X-Kong-Upstream-Latency": []string{"1"}, "X-Kong-Proxy-Latency": []string{"2"}, "Via": []string{"kong/0.14.1"}, "Content-Type": []string{"application/json"}, "Content-Length": []string{"573"}, "Connection": []string{"keep-alive"}, "Date": []string{"Tue, 14 May 2019 13:12:22 GMT"}},
		`{"Components":[{"ID":"x0c0s21b0n0","Type":"Node","State":"Standby","Flag":"Alert","Enabled":true,"Role":"Compute","NID":4,"NetType":"Sling","Arch":"X86"}]}
`},
	{"https://slice-sms.us.cray.com:30443/apis/smd/hsm/v1/State/Components?id=x0c0s21b0n0",
		"GET",
		``,
		map[string][]string{"Accept": []string{"application/json"}, "Content-Type": []string{"application/json"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Connection": []string{"keep-alive"}, "Date": []string{"Tue, 14 May 2019 13:12:22 GMT"}, "X-Kong-Upstream-Latency": []string{"2"}, "X-Kong-Proxy-Latency": []string{"0"}, "Via": []string{"kong/0.14.1"}, "Content-Type": []string{"application/json"}, "Content-Length": []string{"1465"}},
		`{"Components":[{"ID":"x0c0s21b0n0","Type":"Node","State":"Standby","Flag":"Alert","Enabled":true,"Role":"Compute","NID":4,"NetType":"Sling","Arch":"X86"}]}
`},
	{"https://slice-sms.us.cray.com:30443/apis/smd/hsm/v1/Inventory/ComponentEndpoints?id=x0c0s21b0n0",
		"GET",
		``,
		map[string][]string{"Content-Type": []string{"application/json"}, "Accept": []string{"application/json"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"X-Kong-Proxy-Latency": []string{"1"}, "Via": []string{"kong/0.14.1"}, "Content-Type": []string{"application/json"}, "Connection": []string{"keep-alive"}, "Date": []string{"Tue, 14 May 2019 13:12:22 GMT"}, "X-Kong-Upstream-Latency": []string{"2"}},
		`{"ComponentEndpoints":[{"ID":"x0c0s21b0n0","Type":"Node","RedfishType":"ComputerSystem","RedfishSubtype":"Physical","MACAddr":"a4:bf:01:2c:f8:19","UUID":"f5f9d054-bb10-11e7-ab21-a4bf012cf819","OdataID":"/redfish/v1/Systems/QSBP74304730","RedfishEndpointID":"x0c0s21b0","RedfishEndpointFQDN":"10.4.0.8","RedfishURL":"10.4.0.8/redfish/v1/Systems/QSBP74304730","ComponentEndpointType":"ComponentEndpointComputerSystem","RedfishSystemInfo":{"Name":"S2600BPB","Actions":{"#ComputerSystem.Reset":{"ResetType@Redfish.AllowableValues":["On","ForceOff","GracefulShutdown","GracefulRestart","ForceRestart","Nmi"],"target":"/redfish/v1/Systems/QSBP74304730/Actions/ComputerSystem.Reset"}},"EthernetNICInfo":[{"RedfishId":"","@odata.id":"","Description":"Missing interface 1, MAC computed via workaround","MACAddress":"a4:bf:01:2c:f8:19"},{"RedfishId":"","@odata.id":"","Description":"Missing interface 2, MAC computed via workaround","MACAddress":"a4:bf:01:2c:f8:1a"}]}}]}
`},
	{"https://slice-sms.us.cray.com:30443/apis/smd/hsm/v1/Inventory/RedfishEndpoints?id=x0c0s21b0",
		"GET",
		``,
		map[string][]string{"Accept": []string{"application/json"}, "Content-Type": []string{"application/json"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}, "Content-Length": []string{"1491"}, "Connection": []string{"keep-alive"}, "Date": []string{"Tue, 14 May 2019 13:12:22 GMT"}, "X-Kong-Upstream-Latency": []string{"2"}, "X-Kong-Proxy-Latency": []string{"1"}, "Via": []string{"kong/0.14.1"}},
		`{"RedfishEndpoints":[{"ID":"x0c0s21b0","Type":"NodeBMC","Hostname":"10.4.0.8","Domain":"","FQDN":"10.4.0.8","Enabled":true,"UUID":"4d3c9478-8d43-482f-8371-a7838e45a674","User":"root","Password":"********","MACAddr":"a4bf012cf81d","RediscoverOnUpdate":true,"DiscoveryInfo":{"LastDiscoveryAttempt":"2019-05-09T13:08:32.318980Z","LastDiscoveryStatus":"DiscoverOK","RedfishVersion":"1.1.0"}}]}
`},
	{"https://10.4.0.8/redfish/v1/Systems/QSBP74304730/Actions/ComputerSystem.Reset",
		"POST",
		`{"ResetType": "ForceOff"}`,
		map[string][]string{"Authorization": []string{"Basic cm9vdDppbml0aWFsMA=="}, "Accept": []string{"*/*"}, "Content-Type": []string{"application/json"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"X-Xss-Protection": []string{"1; mode=block"}, "Content-Type": []string{"application/json"}, "Content-Length": []string{"154"}, "Date": []string{"Tue, 14 May 2019 13:12:22 GMT"}, "Server": []string{"lighttpd/1.4.45"}, "Strict-Transport-Security": []string{"max-age=31536000; includeSubdomains; preload"}, "X-Ua-Compatible": []string{"IE=11"}, "X-Frame-Options": []string{"SAMEORIGIN"}},
		`{"@odata.type":"#Message.v1_0_4.Message","MessageId":"Base.1.1.0.Success","Message":"Successfully Completed Request","Severity":"OK","Resolution":"None"}
`},
	{"https://slice-sms.us.cray.com:30443/apis/smd/hsm/v1/locks",
		"POST",
		`{"reason":"Performing a power Off","owner":"CAPMC","lifetime":180,"xnames":["x0c0s21b0n0"]}`,
		map[string][]string{"Content-Type": []string{"application/json"}, "Accept": []string{"application/json"}},
		"201 OK", 201,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}, "Date": []string{"Wed, 27 Feb 2019 04:37:13 GMT"}, "Content-Length": []string{"62"}},
		`{"uri":"/hsm/v1/locks/08c2e624-ded3-11e9-8a34-2a2ae2dbcce4"}]}
`},
}
var XnameOffRecursiveChassisSSData = []sstorage.MockLookup{
	{
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname:    "x0c0s21b0n0",
				URL:      "10.4.0.8/redfish/v1/Systems/QSBP74304730",
				Username: "root",
				Password: "********",
			},
			Err: nil,
		},
	},
}

func TestXnameOffRecursiveChassis(t *testing.T) {
	var testReq = testData{"/capmc/xname_off",
		"POST",
		`{ "force": false, "recursive": true, "reason": "", "xnames": [ "x0c0" ] }`,
		map[string][]string{"User-Agent": []string{"curl/7.37.0"}, "Content-Type": []string{"application/json"}, "Accept": []string{"application/json"}, "Content-Length": []string{"73"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}},
		`{"e":0,"err_msg":""}
`}
	runTest(t, XnameOffRecursiveChassisHSM, &testReq, &XnameOffRecursiveChassisReplayData, XnameOffRecursiveChassisSSData)
}

var XnameOnRecursiveChassisHSM = "https://slice-sms.us.cray.com:30443/apis/smd"
var XnameOnRecursiveChassisReplayData = []testData{
	{"https://slice-sms.us.cray.com:30443/apis/smd/hsm/v1/State/Components/Query/x0c0?state=%21Empty&type=chassis&type=cabinetpdu&type=routermodule&type=hsnboard&type=computemodule&type=node",
		"GET",
		``,
		map[string][]string{"Content-Type": []string{"application/json"}, "Accept": []string{"application/json"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Connection": []string{"keep-alive"}, "Date": []string{"Tue, 14 May 2019 13:15:23 GMT"}, "X-Kong-Upstream-Latency": []string{"7"}, "X-Kong-Proxy-Latency": []string{"1"}, "Via": []string{"kong/0.14.1"}, "Content-Type": []string{"application/json"}, "Content-Length": []string{"573"}},
		`{"Components":[{"ID":"x0c0s21b0n0","Type":"Node","State":"Standby","Flag":"Alert","Enabled":true,"Role":"Compute","NID":4,"NetType":"Sling","Arch":"X86"}]}
`},
	{"https://slice-sms.us.cray.com:30443/apis/smd/hsm/v1/State/Components?id=x0c0s21b0n0",
		"GET",
		``,
		map[string][]string{"Content-Type": []string{"application/json"}, "Accept": []string{"application/json"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Length": []string{"1465"}, "Connection": []string{"keep-alive"}, "Date": []string{"Tue, 14 May 2019 13:15:23 GMT"}, "X-Kong-Upstream-Latency": []string{"2"}, "X-Kong-Proxy-Latency": []string{"1"}, "Via": []string{"kong/0.14.1"}, "Content-Type": []string{"application/json"}},
		`{"Components":[{"ID":"x0c0s21b0n0","Type":"Node","State":"Standby","Flag":"Alert","Enabled":true,"Role":"Compute","NID":4,"NetType":"Sling","Arch":"X86"}]}
`},
	{"https://slice-sms.us.cray.com:30443/apis/smd/hsm/v1/Inventory/ComponentEndpoints?id=x0c0s21b0n0",
		"GET",
		``,
		map[string][]string{"Accept": []string{"application/json"}, "Content-Type": []string{"application/json"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Connection": []string{"keep-alive"}, "Date": []string{"Tue, 14 May 2019 13:15:23 GMT"}, "X-Kong-Upstream-Latency": []string{"2"}, "X-Kong-Proxy-Latency": []string{"1"}, "Via": []string{"kong/0.14.1"}, "Content-Type": []string{"application/json"}},
		`{"ComponentEndpoints":[{"ID":"x0c0s21b0n0","Type":"Node","RedfishType":"ComputerSystem","RedfishSubtype":"Physical","MACAddr":"a4:bf:01:2c:f8:19","UUID":"f5f9d054-bb10-11e7-ab21-a4bf012cf819","OdataID":"/redfish/v1/Systems/QSBP74304730","RedfishEndpointID":"x0c0s21b0","RedfishEndpointFQDN":"10.4.0.8","RedfishURL":"10.4.0.8/redfish/v1/Systems/QSBP74304730","ComponentEndpointType":"ComponentEndpointComputerSystem","RedfishSystemInfo":{"Name":"S2600BPB","Actions":{"#ComputerSystem.Reset":{"ResetType@Redfish.AllowableValues":["On","ForceOff","GracefulShutdown","GracefulRestart","ForceRestart","Nmi"],"target":"/redfish/v1/Systems/QSBP74304730/Actions/ComputerSystem.Reset"}},"EthernetNICInfo":[{"RedfishId":"","@odata.id":"","Description":"Missing interface 1, MAC computed via workaround","MACAddress":"a4:bf:01:2c:f8:19"},{"RedfishId":"","@odata.id":"","Description":"Missing interface 2, MAC computed via workaround","MACAddress":"a4:bf:01:2c:f8:1a"}]}}]}
`},
	{"https://slice-sms.us.cray.com:30443/apis/smd/hsm/v1/Inventory/RedfishEndpoints?id=x0c0s21b0",
		"GET",
		``,
		map[string][]string{"Accept": []string{"application/json"}, "Content-Type": []string{"application/json"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Via": []string{"kong/0.14.1"}, "Content-Type": []string{"application/json"}, "Content-Length": []string{"1491"}, "Connection": []string{"keep-alive"}, "Date": []string{"Tue, 14 May 2019 13:15:23 GMT"}, "X-Kong-Upstream-Latency": []string{"1"}, "X-Kong-Proxy-Latency": []string{"0"}},
		`{"RedfishEndpoints":[{"ID":"x0c0s21b0","Type":"NodeBMC","Hostname":"10.4.0.8","Domain":"","FQDN":"10.4.0.8","Enabled":true,"UUID":"4d3c9478-8d43-482f-8371-a7838e45a674","User":"root","Password":"********","MACAddr":"a4bf012cf81d","RediscoverOnUpdate":true,"DiscoveryInfo":{"LastDiscoveryAttempt":"2019-05-09T13:08:32.318980Z","LastDiscoveryStatus":"DiscoverOK","RedfishVersion":"1.1.0"}}]}
`},
	{"https://10.4.0.8/redfish/v1/Systems/QSBP74304730/Actions/ComputerSystem.Reset",
		"POST",
		`{"ResetType": "On"}`,
		map[string][]string{"Authorization": []string{"Basic cm9vdDppbml0aWFsMA=="}, "Accept": []string{"*/*"}, "Content-Type": []string{"application/json"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"X-Ua-Compatible": []string{"IE=11"}, "X-Frame-Options": []string{"SAMEORIGIN"}, "X-Xss-Protection": []string{"1; mode=block"}, "Content-Type": []string{"application/json"}, "Content-Length": []string{"154"}, "Date": []string{"Tue, 14 May 2019 13:15:23 GMT"}, "Server": []string{"lighttpd/1.4.45"}, "Strict-Transport-Security": []string{"max-age=31536000; includeSubdomains; preload"}},
		`{"@odata.type":"#Message.v1_0_4.Message","MessageId":"Base.1.1.0.Success","Message":"Successfully Completed Request","Severity":"OK","Resolution":"None"}
`},
	{"https://slice-sms.us.cray.com:30443/apis/smd/hsm/v1/locks",
		"POST",
		`{"reason":"Performing a power Off","owner":"CAPMC","lifetime":180,"xnames":["x0c0s21b0n0"]}`,
		map[string][]string{"Content-Type": []string{"application/json"}, "Accept": []string{"application/json"}},
		"201 OK", 201,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}, "Date": []string{"Wed, 27 Feb 2019 04:37:13 GMT"}, "Content-Length": []string{"62"}},
		`{"uri":"/hsm/v1/locks/08c2e624-ded3-11e9-8a34-2a2ae2dbcce4"}]}
`},
}
var XnameOnRecursiveChassisSSData = []sstorage.MockLookup{
	{
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname:    "x0c0s21b0n0",
				URL:      "10.4.0.8/redfish/v1/Systems/QSBP74304730",
				Username: "root",
				Password: "********",
			},
			Err: nil,
		},
	},
}

func TestXnameOnRecursiveChassis(t *testing.T) {
	var testReq = testData{"/capmc/xname_on",
		"POST",
		`{ "force": false, "recursive": true, "reason": "", "xnames": [ "x0c0" ] }`,
		map[string][]string{"User-Agent": []string{"curl/7.37.0"}, "Content-Type": []string{"application/json"}, "Accept": []string{"application/json"}, "Content-Length": []string{"73"}},
		"200 OK", 200,
		"HTTP/1.1", 1, 1,
		map[string][]string{"Content-Type": []string{"application/json"}},
		`{"e":0,"err_msg":""}
`}
	runTest(t, XnameOnRecursiveChassisHSM, &testReq, &XnameOnRecursiveChassisReplayData, XnameOnRecursiveChassisSSData)
}
