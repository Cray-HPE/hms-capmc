//
// MIT License
//
// (C) Copyright [2019-2022] Hewlett Packard Enterprise Development LP
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

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path"
	"reflect"
	"testing"

	base "github.com/Cray-HPE/hms-base"
	"github.com/Cray-HPE/hms-capmc/internal/capmc"
	compcreds "github.com/Cray-HPE/hms-compcredentials"
	sstorage "github.com/Cray-HPE/hms-securestorage"
	"github.com/Cray-HPE/hms-smd/pkg/sm"
)

var vaultData = []sstorage.MockLookup{
	{
		Input: sstorage.InputLookup{
			Key: "secret/hms-cred/x0c0s0b0n0",
		},
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname: "x0c0s0b0n0",
			},
		},
	},
	{
		Input: sstorage.InputLookup{
			Key: "secret/hms-cred/x0c0s0b0n1",
		},
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname: "x0c0s0b0n1",
			},
		},
	},
	{
		Input: sstorage.InputLookup{
			Key: "secret/hms-cred/x0c0s0b1n0",
		},
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname: "x0c0s0b1n0",
			},
		},
	},
	{
		Input: sstorage.InputLookup{
			Key: "secret/hms-cred/x0c0s0b1n1",
		},
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname: "x0c0s0b1n1",
			},
		},
	},
	{
		Input: sstorage.InputLookup{
			Key: "secret/hms-cred/x0c0s1b0n0",
		},
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname: "x0c0s1b0n0",
			},
		},
	},
	{
		Input: sstorage.InputLookup{
			Key: "secret/hms-cred/x0c0s1b0n1",
		},
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname: "x0c0s1b0n1",
			},
		},
	},
	{
		Input: sstorage.InputLookup{
			Key: "secret/hms-cred/x0c0s1b1n0",
		},
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname: "x0c0s1b1n0",
			},
		},
	},
	{
		Input: sstorage.InputLookup{
			Key: "secret/hms-cred/x0c0s1b1n1",
		},
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname: "x0c0s1b1n1",
			},
		},
	},
	{
		Input: sstorage.InputLookup{
			Key: "secret/hms-cred/x0c0s2b0n0",
		},
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname: "x0c0s2b0n0",
			},
		},
	},
	{
		Input: sstorage.InputLookup{
			Key: "secret/hms-cred/x0c0s2b0n1",
		},
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname: "x0c0s2b0n1",
			},
		},
	},
	{
		Input: sstorage.InputLookup{
			Key: "secret/hms-cred/x0c0s2b1n0",
		},
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname: "x0c0s2b1n0",
			},
		},
	},
	{
		Input: sstorage.InputLookup{
			Key: "secret/hms-cred/x0c0s2b1n1",
		},
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname: "x0c0s2b1n1",
			},
		},
	},
	{
		Input: sstorage.InputLookup{
			Key: "secret/hms-cred/x0c0s3b0n0",
		},
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname: "x0c0s3b0n0",
			},
		},
	},
	{
		Input: sstorage.InputLookup{
			Key: "secret/hms-cred/x0c0s3b0n1",
		},
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname: "x0c0s3b0n1",
			},
		},
	},
	{
		Input: sstorage.InputLookup{
			Key: "secret/hms-cred/x0c0s3b1n0",
		},
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname: "x0c0s3b1n0",
			},
		},
	},
	{
		Input: sstorage.InputLookup{
			Key: "secret/hms-cred/x0c0s3b1n1",
		},
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname: "x0c0s3b1n1",
			},
		},
	},
	{
		Input: sstorage.InputLookup{
			Key: "secret/hms-cred/x0c0s4b0n0",
		},
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname: "x0c0s4b0n0",
			},
		},
	},
	{
		Input: sstorage.InputLookup{
			Key: "secret/hms-cred/x0c0s4b0n1",
		},
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname: "x0c0s4b0n1",
			},
		},
	},
	{
		Input: sstorage.InputLookup{
			Key: "secret/hms-cred/x0c0s4b1n0",
		},
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname: "x0c0s4b1n0",
			},
		},
	},
	{
		Input: sstorage.InputLookup{
			Key: "secret/hms-cred/x0c0s4b1n1",
		},
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname: "x0c0s4b1n1",
			},
		},
	},
	{
		Input: sstorage.InputLookup{
			Key: "secret/hms-cred/x0c0s5b0n0",
		},
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname: "x0c0s5b0n0",
			},
		},
	},
	{
		Input: sstorage.InputLookup{
			Key: "secret/hms-cred/x0c0s5b0n1",
		},
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname: "x0c0s5b0n1",
			},
		},
	},
	{
		Input: sstorage.InputLookup{
			Key: "secret/hms-cred/x0c0s5b1n0",
		},
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname: "x0c0s5b1n0",
			},
		},
	},
	{
		Input: sstorage.InputLookup{
			Key: "secret/hms-cred/x0c0s5b1n1",
		},
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname: "x0c0s5b1n1",
			},
		},
	},
	{
		Input: sstorage.InputLookup{
			Key: "secret/hms-cred/x0c0s6b0n0",
		},
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname: "x0c0s6b0n0",
			},
		},
	},
	{
		Input: sstorage.InputLookup{
			Key: "secret/hms-cred/x0c0s6b0n1",
		},
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname: "x0c0s6b0n1",
			},
		},
	},
	{
		Input: sstorage.InputLookup{
			Key: "secret/hms-cred/x0c0s6b1n0",
		},
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname: "x0c0s6b1n0",
			},
		},
	},
	{
		Input: sstorage.InputLookup{
			Key: "secret/hms-cred/x0c0s6b1n1",
		},
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname: "x0c0s6b1n1",
			},
		},
	},
	{
		Input: sstorage.InputLookup{
			Key: "secret/hms-cred/x0c0s7b0n0",
		},
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname: "x0c0s7b0n0",
			},
		},
	},
	{
		Input: sstorage.InputLookup{
			Key: "secret/hms-cred/x0c0s7b0n1",
		},
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname: "x0c0s7b0n1",
			},
		},
	},
	{
		Input: sstorage.InputLookup{
			Key: "secret/hms-cred/x0c0s7b1n0",
		},
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname: "x0c0s7b1n0",
			},
		},
	},
	{
		Input: sstorage.InputLookup{
			Key: "secret/hms-cred/x0c0s7b1n1",
		},
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname: "x0c0s7b1n1",
			},
		},
	},
}

type clientMock struct {
	Body       []byte
	StatusCode int
}

type hsmMock struct {
	Components         clientMock
	ComponentEndpoints clientMock
}

// hsmTestMock is a RoundTripper that mocks HMS API calls
func hsmTestMock(m *hsmMock) RoundTripFunc {
	return func(r *http.Request) (*http.Response, error) {
		var (
			statusCode int
			body       io.ReadCloser
		)

		if m != nil {
			switch path.Base(r.URL.Path) {
			case "Components":
				statusCode = m.Components.StatusCode
				body = ioutil.NopCloser(bytes.NewBuffer(m.Components.Body))
			case "ComponentEndpoints":
				statusCode = m.ComponentEndpoints.StatusCode
				body = ioutil.NopCloser(bytes.NewBuffer(m.ComponentEndpoints.Body))
			default:
				statusCode = http.StatusInternalServerError
				body = ioutil.NopCloser(bytes.NewBuffer([]byte{}))
			}

			return &http.Response{
				Status: fmt.Sprintf("%d %s",
					statusCode,
					http.StatusText(statusCode)),
				StatusCode: statusCode,
				Body:       body,
				Header:     make(http.Header),
				Request:    r,
			}, nil
		} else {
			return &http.Response{
				Status: fmt.Sprintf("%d %s",
					http.StatusNotFound,
					http.StatusText(http.StatusNotFound)),
				StatusCode: http.StatusNotFound,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"e": 2, "err_msg": "path not found"}`)),
				Header:     make(http.Header),
				Request:    r,
			}, nil
		}
	}
}

func DoPowerCapCapabilitiesTestFunc(t *testing.T) RoundTripFunc {
	return func(req *http.Request) (*http.Response, error) {
		switch req.URL.String() {
		case "http://localhost:27779/Inventory/Hardware/Query/all":
			testSystemHWInventoryAll := loadTestDataBytes(t, "system-hw-inventory-all.input")
			return &http.Response{
				Status: fmt.Sprintf("%d %s",
					http.StatusOK,
					http.StatusText(http.StatusOK)),
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader(testSystemHWInventoryAll)),
				Header:     make(http.Header),
				Request:    req,
			}, nil
		case "http://localhost:27779/Inventory/ComponentEndpoints?type=node":
			testNodeComponentEndpoints := loadTestDataBytes(t, "componentendpoints-nodes-full-system.input")
			return &http.Response{
				Status: fmt.Sprintf("%d %s",
					http.StatusOK,
					http.StatusText(http.StatusOK)),
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader(testNodeComponentEndpoints)),
				Header:     make(http.Header),
				Request:    req,
			}, nil
		case "http://localhost:27779/Inventory/ComponentEndpoints?nid=1&nid=2&nid=3&nid=4&nid=9&nid=1002&nid=1003&nid=1004&nid=1005":
			testNodeComponentEndpoints := loadTestDataBytes(t, "componentendpoints-nodes-full-system.input")
			return &http.Response{
				Status: fmt.Sprintf("%d %s",
					http.StatusOK,
					http.StatusText(http.StatusOK)),
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader(testNodeComponentEndpoints)),
				Header:     make(http.Header),
				Request:    req,
			}, nil
		case "http://localhost:27779/State/Components?type=node":
			testNodeComponents := loadTestDataBytes(t, "components-nodes-full-system.input")
			return &http.Response{
				Status: fmt.Sprintf("%d %s",
					http.StatusOK,
					http.StatusText(http.StatusOK)),
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader(testNodeComponents)),
				Header:     make(http.Header),
				Request:    req,
			}, nil
		case "http://localhost:27779/State/Components?nid=1&nid=2&nid=3&nid=4&nid=9&nid=1002&nid=1003&nid=1004&nid=1005":
			testNodeComponents := loadTestDataBytes(t, "components-nodes-full-system.input")
			return &http.Response{
				Status: fmt.Sprintf("%d %s",
					http.StatusOK,
					http.StatusText(http.StatusOK)),
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader(testNodeComponents)),
				Header:     make(http.Header),
				Request:    req,
			}, nil
		case "http://localhost:27779/State/Components?nid=5&nid=6&nid=7&nid=8&nid=1000&nid=1001&nid=1006&nid=1007":
			// bad node case (no valid nodes)
			testNodeComponents := loadTestDataBytes(t, "components-nodes-empty.input")
			return &http.Response{
				Status: fmt.Sprintf("%d %s",
					http.StatusOK,
					http.StatusText(http.StatusOK)),
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader(testNodeComponents)),
				Header:     make(http.Header),
				Request:    req,
			}, nil
		case "http://localhost:27779/Inventory/ComponentEndpoints?nid=5&nid=6&nid=7&nid=8&nid=1000&nid=1001&nid=1006&nid=1007":
			// bad node case
			testNodeComponentEndpoints := loadTestDataBytes(t, "componentendpoints-nodes-full-system.input")
			return &http.Response{
				Status: fmt.Sprintf("%d %s",
					http.StatusOK,
					http.StatusText(http.StatusOK)),
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader(testNodeComponentEndpoints)),
				Header:     make(http.Header),
				Request:    req,
			}, nil
		case "http://localhost:27779/State/Components?nid=1&nid=2&nid=5&nid=6&nid=1002&nid=1003&nid=1006&nid=1007":
			//mixed good and bad node case
			testNodeComponents := loadTestDataBytes(t, "components-nodes-mixed.input")
			return &http.Response{
				Status: fmt.Sprintf("%d %s",
					http.StatusOK,
					http.StatusText(http.StatusOK)),
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader(testNodeComponents)),
				Header:     make(http.Header),
				Request:    req,
			}, nil
		case "http://localhost:27779/Inventory/ComponentEndpoints?nid=1&nid=2&nid=5&nid=6&nid=1002&nid=1003&nid=1006&nid=1007":
			//mixed good and bad node case
			testNodeComponentEndpoints := loadTestDataBytes(t, "componentendpoints-nodes-full-system.input")
			return &http.Response{
				Status: fmt.Sprintf("%d %s",
					http.StatusOK,
					http.StatusText(http.StatusOK)),
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader(testNodeComponentEndpoints)),
				Header:     make(http.Header),
				Request:    req,
			}, nil
		default:
			//fmt.Printf("handling default case %s\n", req.URL.String())
			return &http.Response{
				Status: fmt.Sprintf("%d %s",
					http.StatusNotImplemented,
					http.StatusText(http.StatusNotImplemented)),
				StatusCode: http.StatusNotImplemented,
				Body:       ioutil.NopCloser(bytes.NewBufferString("")),
				Header:     make(http.Header),
				Request:    req,
			}, nil
		}
	}
}

func TestDoPowerCapCapabilities(t *testing.T) {
	var tSvc CapmcD
	var err error
	tSvc.hsmURL, err = url.Parse("http://localhost:27779")
	if err != nil {
		t.Fatal(err)
	}
	testClient := NewTestClient(DoPowerCapCapabilitiesTestFunc(t))
	tSvc.rfClient = testClient
	tSvc.smClient = testClient
	tSvc.config = loadConfig("")
	handler := http.HandlerFunc(tSvc.doPowerCapCapabilities)

	tests := []struct {
		name    string
		method  string
		path    string
		body    io.Reader
		retCode int
		retBody io.Reader
	}{
		{
			"Get", http.MethodGet, capmc.PowerCapCapabilities,
			nil,
			http.StatusMethodNotAllowed,
			bytes.NewBuffer(json.RawMessage(GetNotAllowedJSON)),
		}, {
			"Put", http.MethodPut, capmc.PowerCapCapabilities,
			nil,
			http.StatusMethodNotAllowed,
			bytes.NewBuffer(json.RawMessage(PutNotAllowedJSON)),
		}, {
			"Post empty", http.MethodPost, capmc.PowerCapCapabilities,
			bytes.NewBufferString(""),
			http.StatusBadRequest,
			bytes.NewBuffer(json.RawMessage(`{"e":400,"err_msg":"Bad Request: JSON: EOF"}`)),
		}, {
			"Post empty body", http.MethodPost, capmc.PowerCapCapabilities,
			bytes.NewBuffer(json.RawMessage(`{"nids":[]}`)),
			http.StatusOK,
			bytes.NewBuffer(json.RawMessage(loadTestDataBytes(t, "power_cap_cap_case3.input"))),
		}, {
			"Post w/nids", http.MethodPost, capmc.PowerCapCapabilities,
			bytes.NewBuffer(json.RawMessage(`{"nids":[1,2,3,4,9,1002,1003,1004,1005]}`)),
			http.StatusOK,
			bytes.NewBuffer(json.RawMessage(loadTestDataBytes(t, "power_cap_cap_case4.input"))),
		}, {
			"Post w/ only bad nids", http.MethodPost, capmc.PowerCapCapabilities,
			bytes.NewBuffer(json.RawMessage(`{"nids":[5,6,7,8,1000,1001,1006,1007]}`)),
			http.StatusOK,
			bytes.NewBuffer(json.RawMessage(loadTestDataBytes(t, "power_cap_cap_case5.input"))),
		}, {
			"Post w/ mixed good and bad nids", http.MethodPost, capmc.PowerCapCapabilities,
			bytes.NewBuffer(json.RawMessage(`{"nids":[1,2,5,6,1002,1003,1006,1007]}`)),
			http.StatusOK,
			bytes.NewBuffer(json.RawMessage(loadTestDataBytes(t, "power_cap_cap_case6.input"))),
		}, {
			"Delete", http.MethodDelete, capmc.PowerCapCapabilities,
			nil,
			http.StatusMethodNotAllowed,
			bytes.NewBuffer(json.RawMessage(DeleteNotAllowedJSON)),
		}, {
			"Connect", http.MethodConnect, capmc.PowerCapCapabilities,
			nil,
			http.StatusMethodNotAllowed,
			bytes.NewBuffer(json.RawMessage(ConnectNotAllowedJSON)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.path, tt.body)
			if err != nil {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			// check the error code of the transaction
			if tt.retCode != w.Code {
				t.Errorf("Returned wrong status code: got %v want %v",
					w.Code, tt.retCode)
			}

			// if we expected data back, verify it is correct
			if tt.retBody != nil {
				// pull the json data from the response and expected result
				var expRet []byte
				expRet, err1 := ioutil.ReadAll(tt.retBody)
				gotRet, err2 := ioutil.ReadAll(w.Body)
				if err1 != nil {
					t.Errorf("Read error1: %s", err1)
				}
				if err2 != nil {
					t.Errorf("Read error2: %s", err2)
				}

				// compare the results - careful they are json marshalled objects
				// so they may not come through in the same order all the time!
				compareResults(t, string(expRet), string(gotRet))
			}
		})
	}
}

func TestDoPowerCapGet(t *testing.T) {
	testQueryAllComponents := loadTestDataBytes(t,
		"components-nodes-only-one-chassis.input")
	testQueryAllComponentEndpoints := loadTestDataBytes(t,
		"componentendpoints-nodes-only-one-chassis.input")

	tests := []struct {
		name   string
		method string
		path   string
		body   io.Reader
		ret    int
		*hsmMock
	}{
		{
			name:    "Get",
			method:  http.MethodGet,
			path:    capmc.PowerCapGet,
			body:    nil,
			ret:     http.StatusMethodNotAllowed,
			hsmMock: nil,
		}, {
			name:    "Put",
			method:  http.MethodPut,
			path:    capmc.PowerCapGet,
			body:    nil,
			ret:     http.StatusMethodNotAllowed,
			hsmMock: nil,
		}, {
			name:    "Post empty",
			method:  http.MethodPost,
			path:    capmc.PowerCapGet,
			body:    bytes.NewBufferString(""),
			ret:     http.StatusBadRequest,
			hsmMock: nil,
		}, {
			name:   "Post empty body",
			method: http.MethodPost,
			path:   capmc.PowerCapGet,
			body:   bytes.NewBuffer(json.RawMessage(`{}`)),
			ret:    http.StatusOK,
			hsmMock: &hsmMock{
				Components: clientMock{
					Body:       testQueryAllComponents,
					StatusCode: http.StatusOK,
				},
				ComponentEndpoints: clientMock{
					Body:       testQueryAllComponentEndpoints,
					StatusCode: http.StatusOK,
				},
			},
		}, {
			name:   "Post empty nid list",
			method: http.MethodPost,
			path:   capmc.PowerCapGet,
			body:   bytes.NewBuffer(json.RawMessage(`{"nids":[]}`)),
			ret:    http.StatusOK,
			hsmMock: &hsmMock{
				Components: clientMock{
					Body:       testQueryAllComponents,
					StatusCode: http.StatusOK,
				},
				ComponentEndpoints: clientMock{
					Body:       testQueryAllComponentEndpoints,
					StatusCode: http.StatusOK,
				},
			},
		}, {
			name:   "Post w/nids",
			method: http.MethodPost,
			path:   capmc.PowerCapGet,
			body:   bytes.NewBuffer(json.RawMessage(`{"nids":[1,2,4,8,16,17,18,19]}`)),
			ret:    http.StatusOK,
			hsmMock: &hsmMock{
				Components: clientMock{
					Body:       testQueryAllComponents,
					StatusCode: http.StatusOK,
				},
				ComponentEndpoints: clientMock{
					Body:       testQueryAllComponentEndpoints,
					StatusCode: http.StatusOK,
				},
			},
		}, {
			name:    "Delete",
			method:  http.MethodDelete,
			path:    capmc.PowerCapGet,
			body:    nil,
			ret:     http.StatusMethodNotAllowed,
			hsmMock: nil,
		}, {
			name:    "Connect",
			method:  http.MethodConnect,
			path:    capmc.PowerCapGet,
			body:    nil,
			ret:     http.StatusMethodNotAllowed,
			hsmMock: nil,
		},
	}

	// Use the mock secure storage
	ss, adapter := sstorage.NewMockAdapter()
	ccs := compcreds.NewCompCredStore("secret/hms-cred", ss)
	adapter.LookupData = vaultData

	for _, test := range tests {
		adapter.LookupNum = -1 // use mockAdapter "search" mode
		t.Run(test.name, func(t *testing.T) {
			svc := CapmcD{
				smClient: NewTestClient(hsmTestMock(test.hsmMock)),
				rfClient: NewTestClient(hsmTestMock(test.hsmMock)),
				config:   loadConfig(""),
				ss:       ss,
				ccs:      ccs,
				WPool:    base.NewWorkerPool(100, 100*10),
				debug:    debug,
			}
			svc.WPool.Run()
			svc.hsmURL, _ = url.Parse("http://localhost")

			req, err := http.NewRequest(test.method, test.path, test.body)
			if err != nil {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()
			svc.doPowerCapGet(w, req)

			if test.ret != w.Code {
				t.Errorf("Returned wrong status code: got %v want %v",
					w.Code, test.ret)
			}
			// TODO add data validation here
		})
	}
}

func TestDoPowerCapSet(t *testing.T) {
	componentsData := loadTestDataBytes(t,
		"components-nodes-only-one-chassis.input")
	componentEndpointData := loadTestDataBytes(t,
		"componentendpoints-nodes-only-one-chassis.input")

	tests := []struct {
		name   string
		method string
		path   string
		body   io.Reader
		ret    int
		*hsmMock
	}{
		{
			name:    "Get",
			method:  http.MethodGet,
			path:    capmc.PowerCapSet,
			body:    nil,
			ret:     http.StatusMethodNotAllowed,
			hsmMock: nil,
		}, {
			name:    "Put",
			method:  http.MethodPut,
			path:    capmc.PowerCapSet,
			body:    nil,
			ret:     http.StatusMethodNotAllowed,
			hsmMock: nil,
		}, {
			name:    "Post empty",
			method:  http.MethodPost,
			path:    capmc.PowerCapSet,
			body:    bytes.NewBufferString(""),
			ret:     http.StatusBadRequest,
			hsmMock: nil,
		}, {
			name:    "Post empty body",
			method:  http.MethodPost,
			path:    capmc.PowerCapSet,
			body:    bytes.NewBuffer(json.RawMessage(`{"nids":[]}`)),
			ret:     http.StatusBadRequest,
			hsmMock: nil,
		}, {
			name:   "Post w/nids",
			method: http.MethodPost,
			path:   capmc.PowerCapSet,
			body:   bytes.NewBuffer(json.RawMessage(`{"nids": [ { "nid": 30, "controls": [ { "name": "node", "val": 42 } ] } ] }`)),
			ret:    http.StatusOK,
			hsmMock: &hsmMock{
				Components: clientMock{
					Body:       componentsData,
					StatusCode: http.StatusOK,
				},
				ComponentEndpoints: clientMock{
					Body:       componentEndpointData,
					StatusCode: http.StatusOK,
				},
			},
		}, {
			name:    "Delete",
			method:  http.MethodDelete,
			path:    capmc.PowerCapSet,
			body:    nil,
			ret:     http.StatusMethodNotAllowed,
			hsmMock: nil,
		}, {
			name:    "Connect",
			method:  http.MethodConnect,
			path:    capmc.PowerCapSet,
			body:    nil,
			ret:     http.StatusMethodNotAllowed,
			hsmMock: nil,
		},
	}

	ss, adapter := sstorage.NewMockAdapter()
	ccs := compcreds.NewCompCredStore("secret/hms-cred", ss)
	adapter.LookupData = vaultData

	for _, test := range tests {
		adapter.LookupNum = -1 // use mockAdapter "search" mode
		t.Run(test.name, func(t *testing.T) {
			svc := CapmcD{
				smClient: NewTestClient(hsmTestMock(test.hsmMock)),
				rfClient: NewTestClient(hsmTestMock(test.hsmMock)),
				config:   loadConfig(""),
				ss:       ss,
				ccs:      ccs,
				WPool:    base.NewWorkerPool(100, 100*10),
				debug:    debug,
			}
			svc.WPool.Run()
			svc.hsmURL, _ = url.Parse("http://localhost")

			req, err := http.NewRequest(test.method, test.path, test.body)
			if err != nil {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()
			svc.doPowerCapSet(w, req)

			if test.ret != w.Code {
				t.Errorf("Returned wrong status code: got %v want %v",
					w.Code, test.ret)
			}
			// TODO add data validation here
		})
	}
}

func TestBuildPowerCapCapabilitiesGroup(t *testing.T) {
	var componentEndpoints sm.ComponentEndpointArray
	var monikerGroups []PowerCapCapabilityMonikerGroup
	testNodeComponentEndpoints := loadTestDataBytes(t, "componentendpoints-nodes-full-system.input")
	ceErr := json.Unmarshal(testNodeComponentEndpoints, &componentEndpoints)
	if ceErr != nil {
		t.Error("Unable to unmarshall testNodeComponentEndpoints data: " + ceErr.Error())
	}

	//create an xname to componentEndpoint lookup to be used when querying for ComponentEndpoint PowerControl data
	xnameComponentLookup := make(map[string]*sm.ComponentEndpoint)
	for _, componentEndpoint := range componentEndpoints.ComponentEndpoints {
		xnameComponentLookup[componentEndpoint.ID] = componentEndpoint
	}
	testMonikerGroups := loadTestDataBytes(t, "moniker-groups.input")
	mgErr := json.Unmarshal(testMonikerGroups, &monikerGroups)
	if mgErr != nil {
		t.Error("Unable to unmarshall testMonikerGroups data" + mgErr.Error())
	}

	if len(monikerGroups) > 0 && len(xnameComponentLookup) > 0 {
		for _, mg := range monikerGroups {
			powerCapGroup, pcgErr := buildPowerCapCapabilitiesGroup(mg, xnameComponentLookup)
			if pcgErr != nil {
				t.Error("buildPowerCapCapabilitiesGroup returned an error: " + pcgErr.Error())
			}
			//do some checking of powerCapGroup therefore
			if powerCapGroup.Controls == nil &&
				mg.Name != "3_ssd_bbt_bbst_cpuid_tdp_64_244_0_accelerator" &&
				mg.Name != "3_ssd_bbt_bbst_cpuid_tdp_64_244_3200_accelerator" {
				t.Error("powerCapGroup controls is nil")
			}
		}
	} else {
		t.Error("Unable to unmarshall testMonikerGroups data")
	}
}

func TestConvertSystemHWInventoryToUniqueMonikerGroups(t *testing.T) {
	var testHWInventory sm.SystemHWInventory

	expectedUniqueMonikerGroupCount := 6
	testSystemHWInventoryAll := loadTestDataBytes(t, "system-hw-inventory-all.input")
	err := json.Unmarshal(testSystemHWInventoryAll, &testHWInventory)
	if err != nil {
		t.Error("Unable to unmarshall testSystemHWInventoryAll data")
	}
	uniqueMonikerGroups := convertSystemHWInventoryToUniqueMonikerGroups(testHWInventory)
	if expectedUniqueMonikerGroupCount != len(uniqueMonikerGroups) {
		t.Errorf("expected %d uniqueMonikerGroups, got %d", expectedUniqueMonikerGroupCount, len(uniqueMonikerGroups))
	}
}

func TestGeneratePayload(t *testing.T) {
	var (
		pCtl       []capmc.PowerControl
		goodPwrCtl []capmc.PowerControl
		pLimit     capmc.HpeConfigurePowerLimit
		goodPwrLim capmc.HpeConfigurePowerLimit
		rfCtl      capmc.RFControl
		goodRfCtl  capmc.RFControl
		zeroCtl    capmc.RFControl
	)

	var fiveHundred int = 500
	var oneThousand int = 1000

	goodPwrCtl = []capmc.PowerControl{
		{
			PowerLimit: &capmc.PowerLimit{
				LimitInWatts: &oneThousand,
			},
		},
		{
			PowerLimit: &capmc.PowerLimit{
				LimitInWatts: &fiveHundred,
			},
		},
	}

	var zero int = 0

	goodPwrLim = capmc.HpeConfigurePowerLimit{
		PowerLimits: []capmc.HpePowerLimits{
			{
				PowerLimitInWatts: &oneThousand,
				ZoneNumber:        &zero,
			},
		},
	}

	goodRfCtl = capmc.RFControl{
		SetPoint:    &oneThousand,
		ControlMode: "Automatic",
	}

	zeroCtl = capmc.RFControl{
		ControlMode: "Disabled",
	}

	type args struct {
		node       *NodeInfo
		powerCtl   []capmc.PowerControl
		powerLimit capmc.HpeConfigurePowerLimit
		rfCtl      capmc.RFControl
	}

	mtnNode := NodeInfo{
		RfPowerURL: "/redfish/v1/Chassis/Node0/Power",
	}
	stdNode := NodeInfo{
		RfPowerURL: "/redfish/v1/Chassis/Self/Power",
	}
	A6500Node := NodeInfo{
		RfPowerURL: "/redfish/v1/Chassis/1/Power/AccPowerService/PowerLimit",
	}
	bardNode := NodeInfo{
		RfPowerURL:    "/redfish/v1/Chassis/Node/Controls/NodePowerLimit",
		RfControlsCnt: 5,
	}
	GPU0Node := NodeInfo{
		RfPowerURL:    "/redfish/v1/Chassis/Node/Controls/GPU0PowerLimit",
		RfControlsCnt: 5,
	}

	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "bad mtnNode",
			args: args{
				node:       &mtnNode,
				powerCtl:   pCtl,
				powerLimit: pLimit,
				rfCtl:      rfCtl,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "bad stdNode",
			args: args{
				node:       &stdNode,
				powerCtl:   pCtl,
				powerLimit: pLimit,
				rfCtl:      rfCtl,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "bad A6500Node",
			args: args{
				node:       &A6500Node,
				powerCtl:   pCtl,
				powerLimit: pLimit,
				rfCtl:      rfCtl,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "bad bard node",
			args: args{
				node:       &bardNode,
				powerCtl:   pCtl,
				powerLimit: pLimit,
				rfCtl:      rfCtl,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "good mtnNode",
			args: args{
				node:       &mtnNode,
				powerCtl:   goodPwrCtl,
				powerLimit: pLimit,
				rfCtl:      rfCtl,
			},
			want:    json.RawMessage(`{"PowerControl":[{"PowerLimit":{"LimitInWatts":1000}},{"PowerLimit":{"LimitInWatts":500}}]}`),
			wantErr: false,
		},
		{
			name: "good stdNode",
			args: args{
				node:       &stdNode,
				powerCtl:   goodPwrCtl,
				powerLimit: pLimit,
				rfCtl:      rfCtl,
			},
			want:    json.RawMessage(`{"PowerControl":[{"PowerLimit":{"LimitInWatts":1000}},{"PowerLimit":{"LimitInWatts":500}}]}`),
			wantErr: false,
		},
		{
			name: "good A6500Node",
			args: args{
				node:       &A6500Node,
				powerCtl:   pCtl,
				powerLimit: goodPwrLim,
				rfCtl:      rfCtl,
			},
			want:    json.RawMessage(`{"PowerLimits":[{"PowerLimitInWatts":1000,"ZoneNumber":0}]}`),
			wantErr: false,
		},
		{
			name: "good mtnNode node",
			args: args{
				node:       &bardNode,
				powerCtl:   pCtl,
				powerLimit: pLimit,
				rfCtl:      goodRfCtl,
			},
			want:    json.RawMessage(`{"ControlMode":"Automatic","SetPoint":1000}`),
			wantErr: false,
		},
		{
			name: "good mtnNode GPU",
			args: args{
				node:       &GPU0Node,
				powerCtl:   pCtl,
				powerLimit: pLimit,
				rfCtl:      goodRfCtl,
			},
			want:    json.RawMessage(`{"ControlMode":"Automatic","SetPoint":1000}`),
			wantErr: false,
		},
		{
			name: "good mtnNode zero",
			args: args{
				node:       &bardNode,
				powerCtl:   pCtl,
				powerLimit: pLimit,
				rfCtl:      zeroCtl,
			},
			want:    json.RawMessage(`{"ControlMode":"Disabled"}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//got, err := generatePayload(tt.args.node, tt.args.powerCtl, tt.args.powerLimit)
			got, err := generatePayload(tt.args.node,
				powerGen{
					powerCtl:   tt.args.powerCtl,
					powerLimit: tt.args.powerLimit,
					controls:   tt.args.rfCtl,
				})
			if (err != nil) != tt.wantErr {
				t.Errorf("generatePayload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generatePayload() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestExpandNodeListForControlStruct(t *testing.T) {
	var pc1 = make(map[string]PowerCap)
	pc1["Node Power Limit"] = PowerCap{
		Name:        "Node Power Limit",
		Path:        "/redfish/v1/Chassis/Node0/Power",
		Min:         200,
		Max:         1000,
		PwrCtlIndex: 0,
	}

	var node1 = NodeInfo{
		Hostname:      "x1000c0s0b0n0",
		RfPowerURL:    "/redfish/v1/Chassis/Node0/Power",
		RfControlsCnt: 0,
		RfPwrCtlCnt:   1,
		PowerCaps:     pc1,
	}

	wantNode1 := []string{"/redfish/v1/Chassis/Node0/Power"}

	var pc2 = make(map[string]PowerCap)
	pc2["Node Power Limit"] = PowerCap{
		Name:        "Node Power Limit",
		Path:        "/redfish/v1/Chassis/Node0/Controls/NodePowerLimit",
		Min:         200,
		Max:         1000,
		PwrCtlIndex: 0,
	}
	pc2["GPU0 Power Limit"] = PowerCap{
		Name:        "GPU0 Power Limit",
		Path:        "/redfish/v1/Chassis/Node0/Controls/GPU0PowerLimit",
		Min:         100,
		Max:         200,
		PwrCtlIndex: 1,
	}
	pc2["GPU1 Power Limit"] = PowerCap{
		Name:        "GPU1 Power Limit",
		Path:        "/redfish/v1/Chassis/Node0/Controls/GPU1PowerLimit",
		Min:         100,
		Max:         200,
		PwrCtlIndex: 2,
	}
	pc2["GPU2 Power Limit"] = PowerCap{
		Name:        "GPU2 Power Limit",
		Path:        "/redfish/v1/Chassis/Node0/Controls/GPU2PowerLimit",
		Min:         100,
		Max:         200,
		PwrCtlIndex: 3,
	}
	pc2["GPU3 Power Limit"] = PowerCap{
		Name:        "GPU3 Power Limit",
		Path:        "/redfish/v1/Chassis/Node0/Controls/GPU3PowerLimit",
		Min:         100,
		Max:         200,
		PwrCtlIndex: 4,
	}

	var node2 = NodeInfo{
		Hostname:      "x1000c0s0b0n0",
		RfPowerURL:    "/redfish/v1/Chassis/Node0/Power",
		RfControlsCnt: 5,
		RfPwrCtlCnt:   0,
		PowerCaps:     pc2,
	}

	wantNode2 := []string{
		"/redfish/v1/Chassis/Node0/Controls/NodePowerLimit",
		"/redfish/v1/Chassis/Node0/Controls/GPU0PowerLimit",
		"/redfish/v1/Chassis/Node0/Controls/GPU1PowerLimit",
		"/redfish/v1/Chassis/Node0/Controls/GPU2PowerLimit",
		"/redfish/v1/Chassis/Node0/Controls/GPU3PowerLimit",
	}

	var nlNoExpand []*NodeInfo
	var nlExpand []*NodeInfo

	nlNoExpand = append(nlNoExpand, &node1)
	nlExpand = append(nlExpand, &node2)

	tests := []struct {
		name string
		nl   []*NodeInfo
		want []string
	}{
		{
			"No expand",
			nlNoExpand,
			wantNode1,
		},
		{
			"Expand",
			nlExpand,
			wantNode2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNodeListForControlStruct(tt.nl)
			if len(got) != len(tt.want) {
				var paths []string
				for _, e := range got {
					paths = append(paths, e.RfPowerURL)
				}
				t.Errorf("List different length than expected, got %v, want %v",
					paths, tt.want)
			}
			found := 0
			for _, g := range got {
				for _, w := range tt.want {
					if g.RfPowerURL == w {
						found += 1
						continue
					}
				}
			}
			if found != got[0].RfControlsCnt && found != got[0].RfPwrCtlCnt {
				var paths []string
				for _, e := range got {
					paths = append(paths, e.RfPowerURL)
				}
				t.Errorf("List mismatch, got %v, want %v", paths, tt.want)
			}
		})
	}
}

func TestGenerateControls(t *testing.T) {
	var fiveHundred int = 500
	var twoHundred int = 200
	var fiveThousand int = 5000
	var one int = 1
	var zero int = 0

	// Old style Redfish power cap for single control
	var pc1 = make(map[string]PowerCap)
	pc1["Node Power Limit"] = PowerCap{Name: "Node Power Limit", Path: "/redfish/v1/Chassis/Node0/Power", Min: 200, Max: 1000, PwrCtlIndex: 0}
	node1 := NodeInfo{RfPowerURL: "/redfish/v1/Chassis/Node0/Power", RfPwrCtlCnt: 1, RfControlsCnt: 0, PowerCaps: pc1}
	wantNode1 := make(map[*NodeInfo]powerGen)
	pc1Ctl := make([]capmc.PowerControl, 1)
	pc1Ctl[0] = capmc.PowerControl{PowerLimit: &capmc.PowerLimit{LimitInWatts: &fiveHundred}}
	wantNode1[&node1] = powerGen{powerCtl: pc1Ctl}
	ctl1 := make([]capmc.PowerCapControl, 1)
	ctl1[0] = capmc.PowerCapControl{Name: "Node Power Limit", Val: &fiveHundred}

	// Old style Redfish power cap for two controls
	var pc2 = make(map[string]PowerCap)
	pc2["Node Power Limit"] = PowerCap{Name: "Node Power Limit", Path: "/redfish/v1/Chassis/Node0/Power", Min: 200, Max: 1000, PwrCtlIndex: 0}
	pc2["Accelerator Power Limit"] = PowerCap{Name: "Accelerator Power Limit", Path: "/redfish/v1/Chassis/Node0/Power", Min: 100, Max: 400, PwrCtlIndex: 1}
	node2 := NodeInfo{RfPowerURL: "/redfish/v1/Chassis/Node0/Power", RfPwrCtlCnt: 2, RfControlsCnt: 0, PowerCaps: pc2}
	wantNode2 := make(map[*NodeInfo]powerGen)
	pc2Ctl := make([]capmc.PowerControl, 2)
	pc2Ctl[0] = capmc.PowerControl{PowerLimit: &capmc.PowerLimit{LimitInWatts: &fiveHundred}}
	pc2Ctl[1] = capmc.PowerControl{PowerLimit: &capmc.PowerLimit{LimitInWatts: &twoHundred}}
	wantNode2[&node2] = powerGen{powerCtl: pc2Ctl}
	ctl2 := make([]capmc.PowerCapControl, 2)
	ctl2[0] = capmc.PowerCapControl{Name: "Node Power Limit", Val: &fiveHundred}
	ctl2[1] = capmc.PowerCapControl{Name: "Accelerator Power Limit", Val: &twoHundred}

	// HPE Apollo 6500 style power cap
	var pc3 = make(map[string]PowerCap)
	pc3["HPE Node Power Limit"] = PowerCap{Name: "HPE Node Power Limit", Path: "/redfish/v1/Chassis/1/Power/AccPowerService/PowerLimit", Min: 200, Max: 1000, PwrCtlIndex: 0}
	node3 := NodeInfo{RfPowerURL: "/redfish/v1/Chassis/1/Power/AccPowerService/PowerLimit", RfPwrCtlCnt: 1, RfControlsCnt: 0, PowerCaps: pc3}
	wantNode3 := make(map[*NodeInfo]powerGen)
	wantNode3[&node3] = powerGen{
		powerLimit: capmc.HpeConfigurePowerLimit{
			PowerLimits: []capmc.HpePowerLimits{
				{
					PowerLimitInWatts: &fiveHundred,
					ZoneNumber:        &zero,
				},
			},
		},
	}
	ctl3 := make([]capmc.PowerCapControl, 1)
	ctl3[0] = capmc.PowerCapControl{Name: "HPE Node Power Limit", Val: &fiveHundred}

	// New Redfish controls power capping (Bard Peak), single control
	var pc4 = make(map[string]PowerCap)
	pc4["Node Power Limit"] = PowerCap{Name: "Node Power Limit", Path: "/redfish/v1/Chassis/Node/Controls/NodePowerLimit", Min: 200, Max: 1000, PwrCtlIndex: 0}
	pc4["GPU0 Power Limit"] = PowerCap{Name: "GPU0 Power Limit", Path: "/redfish/v1/Chassis/Node/Controls/GPU0PowerLimit", Min: 100, Max: 400, PwrCtlIndex: 0}
	pc4["GPU1 Power Limit"] = PowerCap{Name: "GPU1 Power Limit", Path: "/redfish/v1/Chassis/Node/Controls/GPU1PowerLimit", Min: 100, Max: 400, PwrCtlIndex: 0}
	pc4["GPU2 Power Limit"] = PowerCap{Name: "GPU2 Power Limit", Path: "/redfish/v1/Chassis/Node/Controls/GPU2PowerLimit", Min: 100, Max: 400, PwrCtlIndex: 0}
	pc4["GPU3 Power Limit"] = PowerCap{Name: "GPU3 Power Limit", Path: "/redfish/v1/Chassis/Node/Controls/GPU3PowerLimit", Min: 100, Max: 400, PwrCtlIndex: 0}
	node4 := NodeInfo{RfPowerURL: "/redfish/v1/Chassis/Node/Controls/NodePowerLimit", RfPwrCtlCnt: 0, RfControlsCnt: 5, PowerCaps: pc4}
	wantNode4 := make(map[*NodeInfo]powerGen)
	wantNode4[&node4] = powerGen{
		controls: capmc.RFControl{
			SetPoint:    &fiveHundred,
			ControlMode: "Automatic",
		},
	}
	ctl4 := make([]capmc.PowerCapControl, 1)
	ctl4[0] = capmc.PowerCapControl{Name: "Node Power Limit", Val: &fiveHundred}

	// New Redfish controls power capping (Bard Peak), multi control
	var pc5 = make(map[string]PowerCap)
	pc5["Node Power Limit"] = PowerCap{Name: "Node Power Limit", Path: "/redfish/v1/Chassis/Node/Controls/NodePowerLimit", Min: 200, Max: 1000, PwrCtlIndex: 0}
	pc5["GPU0 Power Limit"] = PowerCap{Name: "GPU0 Power Limit", Path: "/redfish/v1/Chassis/Node/Controls/GPU0PowerLimit", Min: 100, Max: 400, PwrCtlIndex: 0}
	pc5["GPU1 Power Limit"] = PowerCap{Name: "GPU1 Power Limit", Path: "/redfish/v1/Chassis/Node/Controls/GPU1PowerLimit", Min: 100, Max: 400, PwrCtlIndex: 0}
	pc5["GPU2 Power Limit"] = PowerCap{Name: "GPU2 Power Limit", Path: "/redfish/v1/Chassis/Node/Controls/GPU2PowerLimit", Min: 100, Max: 400, PwrCtlIndex: 0}
	pc5["GPU3 Power Limit"] = PowerCap{Name: "GPU3 Power Limit", Path: "/redfish/v1/Chassis/Node/Controls/GPU3PowerLimit", Min: 100, Max: 400, PwrCtlIndex: 0}
	node5 := NodeInfo{RfPowerURL: "/redfish/v1/Chassis/Node/Controls/NodePowerLimit", RfPwrCtlCnt: 0, RfControlsCnt: 5, PowerCaps: pc5}
	wantNode5 := make(map[*NodeInfo]powerGen)
	wantNode5[&node5] = powerGen{
		controls: capmc.RFControl{
			SetPoint:    &twoHundred,
			ControlMode: "Automatic",
		},
	}
	ctl5 := make([]capmc.PowerCapControl, 5)
	ctl5[0] = capmc.PowerCapControl{Name: "Node Power Limit", Val: &twoHundred}
	ctl5[1] = capmc.PowerCapControl{Name: "GPU0 Power Limit", Val: &twoHundred}
	ctl5[2] = capmc.PowerCapControl{Name: "GPU1 Power Limit", Val: &twoHundred}
	ctl5[3] = capmc.PowerCapControl{Name: "GPU2 Power Limit", Val: &twoHundred}
	ctl5[4] = capmc.PowerCapControl{Name: "GPU3 Power Limit", Val: &twoHundred}

	// Duplicate controls to produce an error
	dupCtl := make([]capmc.PowerCapControl, 2)
	dupCtl[0] = capmc.PowerCapControl{Name: "Node Power Limit", Val: &twoHundred}
	dupCtl[1] = capmc.PowerCapControl{Name: "Node Power Limit", Val: &twoHundred}

	// Value to set too large
	tooBigCtl := make([]capmc.PowerCapControl, 1)
	tooBigCtl[0] = capmc.PowerCapControl{Name: "Node Power Limit", Val: &fiveThousand}

	// Value to set too small but not zero
	tooSmCtl := make([]capmc.PowerCapControl, 1)
	tooSmCtl[0] = capmc.PowerCapControl{Name: "Node Power Limit", Val: &one}

	// Set a zero value
	var pc6 = make(map[string]PowerCap)
	pc6["Node Power Limit"] = PowerCap{Name: "Node Power Limit", Path: "/redfish/v1/Chassis/Node0/Power", Min: 200, Max: 1000, PwrCtlIndex: 0}
	node6 := NodeInfo{RfPowerURL: "/redfish/v1/Chassis/Node0/Power", RfPwrCtlCnt: 1, RfControlsCnt: 0, PowerCaps: pc6}
	wantNode6 := make(map[*NodeInfo]powerGen)
	pc6Ctl := make([]capmc.PowerControl, 1)
	pc6Ctl[0] = capmc.PowerControl{PowerLimit: &capmc.PowerLimit{LimitInWatts: &zero}}
	wantNode6[&node6] = powerGen{powerCtl: pc6Ctl}
	ctl6 := make([]capmc.PowerCapControl, 1)
	ctl6[0] = capmc.PowerCapControl{Name: "Node Power Limit", Val: &zero}

	// Bad control name
	var pc7 = make(map[string]PowerCap)
	pc7["Node Power Limit"] = PowerCap{Name: "Node Power Limit", Path: "/redfish/v1/Chassis/Node0/Power", Min: 200, Max: 1000, PwrCtlIndex: 0}
	node7 := NodeInfo{RfPowerURL: "/redfish/v1/Chassis/Node0/Power", RfPwrCtlCnt: 1, RfControlsCnt: 0, PowerCaps: pc7}
	wantNode7 := make(map[*NodeInfo]powerGen)
	pc7Ctl := make([]capmc.PowerControl, 1)
	pc7Ctl[0] = capmc.PowerControl{PowerLimit: &capmc.PowerLimit{LimitInWatts: &zero}}
	ctl7 := make([]capmc.PowerCapControl, 1)
	ctl7[0] = capmc.PowerCapControl{Name: "Bad control name", Val: &zero}

	// New Redfish controls power capping (Bard Peak), multi control, disable
	var pc8 = make(map[string]PowerCap)
	pc8["Node Power Limit"] = PowerCap{Name: "Node Power Limit", Path: "/redfish/v1/Chassis/Node/Controls/NodePowerLimit", Min: 200, Max: 1000, PwrCtlIndex: 0}
	pc8["GPU0 Power Limit"] = PowerCap{Name: "GPU0 Power Limit", Path: "/redfish/v1/Chassis/Node/Controls/GPU0PowerLimit", Min: 100, Max: 400, PwrCtlIndex: 0}
	pc8["GPU1 Power Limit"] = PowerCap{Name: "GPU1 Power Limit", Path: "/redfish/v1/Chassis/Node/Controls/GPU1PowerLimit", Min: 100, Max: 400, PwrCtlIndex: 0}
	pc8["GPU2 Power Limit"] = PowerCap{Name: "GPU2 Power Limit", Path: "/redfish/v1/Chassis/Node/Controls/GPU2PowerLimit", Min: 100, Max: 400, PwrCtlIndex: 0}
	pc8["GPU3 Power Limit"] = PowerCap{Name: "GPU3 Power Limit", Path: "/redfish/v1/Chassis/Node/Controls/GPU3PowerLimit", Min: 100, Max: 400, PwrCtlIndex: 0}
	node8 := NodeInfo{RfPowerURL: "/redfish/v1/Chassis/Node/Controls/NodePowerLimit", RfPwrCtlCnt: 0, RfControlsCnt: 5, PowerCaps: pc8}
	wantNode8 := make(map[*NodeInfo]powerGen)
	wantNode8[&node8] = powerGen{
		controls: capmc.RFControl{
			ControlMode: "Disabled",
		},
	}
	ctl8 := make([]capmc.PowerCapControl, 5)
	ctl8[0] = capmc.PowerCapControl{Name: "Node Power Limit", Val: &zero}
	ctl8[1] = capmc.PowerCapControl{Name: "GPU0 Power Limit", Val: &zero}
	ctl8[2] = capmc.PowerCapControl{Name: "GPU1 Power Limit", Val: &zero}
	ctl8[3] = capmc.PowerCapControl{Name: "GPU2 Power Limit", Val: &zero}
	ctl8[4] = capmc.PowerCapControl{Name: "GPU3 Power Limit", Val: &zero}

	type args struct {
		node     *NodeInfo
		controls []capmc.PowerCapControl
	}
	tests := []struct {
		name    string
		args    args
		want    map[*NodeInfo]powerGen
		wantErr bool
	}{
		{
			name: "oldRedfishSingle",
			args: args{
				node:     &node1,
				controls: ctl1,
			},
			want:    wantNode1,
			wantErr: false,
		},
		{
			name: "oldRedfishMulti",
			args: args{
				node:     &node2,
				controls: ctl2,
			},
			want:    wantNode2,
			wantErr: false,
		},
		{
			name: "apollo6500",
			args: args{
				node:     &node3,
				controls: ctl3,
			},
			want:    wantNode3,
			wantErr: false,
		},
		{
			name: "newRedfishSingle",
			args: args{
				node:     &node4,
				controls: ctl4,
			},
			want:    wantNode4,
			wantErr: false,
		},
		{
			name: "newRedfishMulti",
			args: args{
				node:     &node5,
				controls: ctl5,
			},
			want:    wantNode5,
			wantErr: false,
		},
		{
			name: "dupCtl",
			args: args{
				node:     &node1,
				controls: dupCtl,
			},
			want:    wantNode1,
			wantErr: true,
		},
		{
			name: "tooBigCtl",
			args: args{
				node:     &node1,
				controls: tooBigCtl,
			},
			want:    wantNode1,
			wantErr: true,
		},
		{
			name: "tooSmCtl",
			args: args{
				node:     &node1,
				controls: tooSmCtl,
			},
			want:    wantNode1,
			wantErr: true,
		},
		{
			name: "setZero",
			args: args{
				node:     &node6,
				controls: ctl6,
			},
			want:    wantNode6,
			wantErr: false,
		},
		{
			name: "badCtlName",
			args: args{
				node:     &node7,
				controls: ctl7,
			},
			want:    wantNode7,
			wantErr: false,
		},
		{
			name: "disable",
			args: args{
				node:     &node8,
				controls: ctl8,
			},
			want:    wantNode8,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateControls(tt.args.node, tt.args.controls)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateControls() error = %+v, wantErr %+v", err, tt.wantErr)
				return
			}

			if (tt.wantErr == true) && (err != nil) {
				return
			}

			if len(got) == 0 {
				if len(got) != len(tt.want) {
					t.Errorf("generateControls() = %+v, want %+v", got, tt.want)
				}
				return
			}

			for _, eg := range got {
				if !reflect.DeepEqual(eg, tt.want[tt.args.node]) {
					t.Errorf("generateControls() = %+v, want %+v", eg,
						tt.want[tt.args.node])
				}
			}
		})
	}
}
