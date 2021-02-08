// MIT License
//
// (C) Copyright [2019-2021] Hewlett Packard Enterprise Development LP
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
	"testing"

	"stash.us.cray.com/HMS/hms-capmc/internal/capmc"
	"stash.us.cray.com/HMS/hms-smd/pkg/sm"

	base "stash.us.cray.com/HMS/hms-base"
	compcreds "stash.us.cray.com/HMS/hms-compcredentials"
	sstorage "stash.us.cray.com/HMS/hms-securestorage"
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
		case "http://localhost:27779/Inventory/ComponentEndpoints?nid=1&nid=2&nid=3&nid=4&nid=1002&nid=1003&nid=1004&nid=1005":
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
		case "http://localhost:27779/State/Components?nid=1&nid=2&nid=3&nid=4&nid=1002&nid=1003&nid=1004&nid=1005":
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
			bytes.NewBuffer(json.RawMessage(`{"nids":[1,2,3,4,1002,1003,1004,1005]}`)),
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
				config: loadConfig(""),
				ss:     ss,
				ccs:    ccs,
				WPool:  base.NewWorkerPool(100, 100*10),
				debug:  debug,
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
				config: loadConfig(""),
				ss:     ss,
				ccs:    ccs,
				WPool:  base.NewWorkerPool(100, 100*10),
				debug:  debug,
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
		powerCapGroup, pcgErr := buildPowerCapCapabilitiesGroup(monikerGroups[0], xnameComponentLookup)
		if pcgErr != nil {
			t.Error("buildPowerCapCapabilitiesGroup returned an error: " + pcgErr.Error())
		}
		//do some checking of powerCapGroup therefore
		if powerCapGroup.Controls == nil {
			t.Error("powerCapGroup controls is nil")
		}
	} else {
		t.Error("Unable to unmarshall testMonikerGroups data")
	}
}

func TestConvertSystemHWInventoryToUniqueMonikerGroups(t *testing.T) {
	var testHWInventory sm.SystemHWInventory

	expectedUniqueMonikerGroupCount := 5
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
