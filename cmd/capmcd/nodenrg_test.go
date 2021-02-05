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

package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"stash.us.cray.com/HMS/hms-capmc/internal/logger"
	"stash.us.cray.com/HMS/hms-capmc/internal/tsdb"
	"testing"

	"stash.us.cray.com/HMS/hms-capmc/internal/capmc"
)

func TestValidateTimeBoundNidRequestiOptions(t *testing.T) {
	var apid string = "1"
	tests := []struct {
		req capmc.TimeBoundNidRequest
		res bool
	}{
		// 0
		{
			req: capmc.TimeBoundNidRequest{
				Apid: apid,
			},
			res: true,
		},
		// 1
		{
			req: capmc.TimeBoundNidRequest{
				JobId: "abcd1234",
			},
			res: true,
		},
		// 2
		{
			req: capmc.TimeBoundNidRequest{
				Nids:      []int{1, 2, 3},
				StartTime: "2000-04-01 01:02:03",
				EndTime:   "2019-09-25 01:02:03",
			},
			res: true,
		},
		// 3
		{
			req: capmc.TimeBoundNidRequest{
				Apid: apid,
				Nids: []int{1, 2, 3},
			},
			res: true,
		},
		// 4
		{
			req: capmc.TimeBoundNidRequest{
				JobId: "abcd1234",
				Nids:  []int{1, 2, 3},
			},
			res: true,
		},
		// 5
		{
			req: capmc.TimeBoundNidRequest{
				Apid:      apid,
				StartTime: "2000-04-01 01:02:03",
			},
			res: true,
		},
		// 6
		{
			req: capmc.TimeBoundNidRequest{
				JobId:     "abcd1234",
				StartTime: "2000-04-01 01:02:03",
			},
			res: true,
		},
		// 7
		{
			req: capmc.TimeBoundNidRequest{
				Nids:      []int{1, 2, 3},
				StartTime: "2000-04-01 01:02:03",
			},
			res: false,
		},
		// 8
		{
			req: capmc.TimeBoundNidRequest{
				Apid:    apid,
				EndTime: "2019-09-25 01:02:03",
			},
			res: true,
		},
		// 9
		{
			req: capmc.TimeBoundNidRequest{
				JobId:   "abcd1234",
				EndTime: "2019-09-25 01:02:03",
			},
			res: true,
		},
		// 10
		{
			req: capmc.TimeBoundNidRequest{
				Nids:    []int{1, 2, 3},
				EndTime: "2019-09-25 01:02:03",
			},
			res: false,
		},
		// 11
		{
			req: capmc.TimeBoundNidRequest{
				StartTime: "2000-04-01 01:02:03",
				EndTime:   "2019-09-25 01:02:03",
			},
			res: false,
		},
		// 12
		{
			req: capmc.TimeBoundNidRequest{
				StartTime: "2000-04-01 01:02:03",
			},
			res: false,
		},
		// 13
		{
			res: false,
		},
	}

	for n, test := range tests {
		isValid, _ := validateTimeBoundNidRequestOptions(test.req)
		if isValid != test.res {
			t.Errorf("TestValidTimeBoundNidRequestOptions Test Case %d: FAIL: Request %v", n, test.req)
		}
	}
}

func DoNodeEnergyTestFunc() RoundTripFunc {
	return func(req *http.Request) (*http.Response, error) {
		//fmt.Println(req.URL.String())
		switch req.URL.String() {
		case "http://localhost:27779/State/Components?type=node",
			"http://localhost:27779/State/Components?nid=0&nid=1",
			"http://localhost:27779/State/Components?nid=1&nid=0":
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(testQueryComponents)),
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

func DoNodeEnergyStatsTestFunc() RoundTripFunc {
	return func(req *http.Request) (*http.Response, error) {
		//fmt.Println(req.URL.String())
		switch req.URL.String() {
		case "http://localhost:27779/State/Components?type=node",
			"http://localhost:27779/State/Components?nid=0&nid=1",
			"http://localhost:27779/State/Components?nid=1&nid=0":
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(testQueryComponents)),
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

func DoNodeEnergyCounterTestFunc() RoundTripFunc {
	return func(req *http.Request) (*http.Response, error) {
		//fmt.Println(req.URL.String())
		switch req.URL.String() {
		case "http://localhost:27779/State/Components?type=node",
			"http://localhost:27779/State/Components?nid=0&nid=1",
			"http://localhost:27779/State/Components?nid=1&nid=0":
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(testQueryComponents)),
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

const testQueryComponents = `
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
      "ID": "x0c0s0b0n0"
    },
    {
      "Arch": "X86",
      "Enabled": true,
      "Flag": "OK",
      "NID": 1,
      "State": "Ready",
      "Role": "Compute",
      "NetType": "Sling",
      "Type": "Node",
      "ID": "x0c0s0b0n1"
    }
  ]
}
`

func TestDoNodeEnergy(t *testing.T) {
	var tSvc CapmcD
	var err error
	tSvc.hsmURL, err = url.Parse("http://localhost:27779")
	if err != nil {
		t.Fatal(err)
	}
	testClient := NewTestClient(DoNodeEnergyTestFunc())
	tSvc.rfClient = testClient
	tSvc.smClient = testClient
	tSvc.config = loadConfig("")

	logger.SetupLogging()

	tsdb.ConfigureDataImplementation(tsdb.DUMMY)
	handler := http.HandlerFunc(tSvc.doNodeEnergy)

	tests := []struct {
		name   string
		method string
		path   string
		body   io.Reader
		ret    int
	}{
		{
			"Get",
			http.MethodGet,
			capmc.NodeEnergy,
			nil,
			http.StatusMethodNotAllowed,
		}, {
			"Put",
			http.MethodPut,
			capmc.NodeEnergy,
			nil,
			http.StatusMethodNotAllowed,
		}, {
			"Post empty",
			http.MethodPost,
			capmc.NodeEnergy,
			bytes.NewBufferString(""),
			http.StatusBadRequest,
		}, {
			"Post empty body",
			http.MethodPost,
			capmc.NodeEnergy,
			bytes.NewBuffer(json.RawMessage(`{}`)),
			http.StatusBadRequest,
		}, {
			"Post w/body apid only",
			http.MethodPost,
			capmc.NodeEnergy,
			bytes.NewBuffer(json.RawMessage(
				`{
					"apid": "831138"
				}`)),
			http.StatusNotImplemented,
		}, {
			"Post w/body times + nids = happy path",
			http.MethodPost,
			capmc.NodeEnergy,
			bytes.NewBuffer(json.RawMessage(
				`{
					"start_time": "2019-07-10 11:36:32",
					"end_time": "2019-07-10 12:36:32",
					"nids": [0, 1]
				}`)),
			http.StatusOK,
		}, {
			"Post w/body include apid",
			http.MethodPost,
			capmc.NodeEnergy,
			bytes.NewBuffer(json.RawMessage(
				`{
					"start_time": "2019-07-10 11:36:32",
					"end_time": "2019-07-10 12:36:32",
					"nids": [0, 1],
					"apid": "831138"
				}`)),
			http.StatusNotImplemented,
		}, {
			"Post w/body include jobid",
			http.MethodPost,
			capmc.NodeEnergy,
			bytes.NewBuffer(json.RawMessage(
				`{
					"start_time": "2019-07-10 11:36:32",
					"end_time": "2019-07-10 12:36:32",
					"nids": [40, 41, 42, 43],
					"job_id": "1234321.sdb"
				}`)),
			http.StatusNotImplemented,
		}, {
			"Post w/body include apid/jobid",
			http.MethodPost,
			capmc.NodeEnergy,
			bytes.NewBuffer(json.RawMessage(
				`{
					"start_time": "2019-07-10 11:36:32",
					"end_time": "2019-07-10 12:36:32",
					"nids": [0, 1],
					"apid": "831138",
					"job_id": "1234321.sdb"
				}`)),
			http.StatusNotImplemented,
		}, {
			"Post w/body job_id only",
			http.MethodPost,
			capmc.NodeEnergy,
			bytes.NewBuffer(json.RawMessage(
				`{
					"job_id": "1234321.sdb"
				}`)),
			http.StatusNotImplemented,
		}, {
			"Delete",
			http.MethodDelete,
			capmc.NodeEnergy,
			nil,
			http.StatusMethodNotAllowed,
		}, {
			"Connect",
			http.MethodConnect,
			capmc.NodeEnergy,
			nil,
			http.StatusMethodNotAllowed,
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

			if tt.ret != w.Code {
				t.Errorf("Returned wrong status code: got %v want %v",
					w.Code, tt.ret)
			}
		})
	}
}

func TestDoNodeEnergyStats(t *testing.T) {
	var tSvc CapmcD
	var err error
	tSvc.hsmURL, err = url.Parse("http://localhost:27779")
	if err != nil {
		t.Fatal(err)
	}
	testClient := NewTestClient(DoNodeEnergyStatsTestFunc())
	tSvc.rfClient = testClient
	tSvc.smClient = testClient
	tSvc.config = loadConfig("")

	logger.SetupLogging()

	tsdb.ConfigureDataImplementation(tsdb.DUMMY)

	handler := http.HandlerFunc(tSvc.doNodeEnergyStats)

	tests := []struct {
		name   string
		method string
		path   string
		body   io.Reader
		ret    int
	}{
		{
			"Get",
			http.MethodGet,
			capmc.NodeEnergyStats,
			nil,
			http.StatusMethodNotAllowed,
		}, {
			"Put",
			http.MethodPut,
			capmc.NodeEnergyStats,
			nil,
			http.StatusMethodNotAllowed,
		}, {
			"Post empty",
			http.MethodPost,
			capmc.NodeEnergyStats,
			bytes.NewBufferString(""),
			http.StatusBadRequest,
		}, {
			"Post empty body",
			http.MethodPost,
			capmc.NodeEnergyStats,
			bytes.NewBuffer(json.RawMessage(`{}`)),
			http.StatusBadRequest,
		}, {
			"Post w/body apid only",
			http.MethodPost,
			capmc.NodeEnergyStats,
			bytes.NewBuffer(json.RawMessage(
				`{
					"apid": "831138"
				}`)),
			http.StatusNotImplemented,
		}, {
			"Post w/body times + nids = happy path",
			http.MethodPost,
			capmc.NodeEnergyStats,
			bytes.NewBuffer(json.RawMessage(
				`{
					"start_time": "2019-07-10 11:36:32",
					"end_time": "2019-07-10 12:36:32",
					"nids": [0, 1]
				}`)),
			http.StatusOK,
		}, {
			"Post w/body include apid",
			http.MethodPost,
			capmc.NodeEnergyStats,
			bytes.NewBuffer(json.RawMessage(
				`{
					"start_time": "2019-07-10 11:36:32",
					"end_time": "2019-07-10 12:36:32",
					"nids": [0, 1],
					"apid": "831138"
				}`)),
			http.StatusNotImplemented,
		}, {
			"Post w/body include jobid",
			http.MethodPost,
			capmc.NodeEnergyStats,
			bytes.NewBuffer(json.RawMessage(
				`{
					"start_time": "2019-07-10 11:36:32",
					"end_time": "2019-07-10 12:36:32",
					"nids": [40, 41, 42, 43],
					"job_id": "1234321.sdb"
				}`)),
			http.StatusNotImplemented,
		}, {
			"Post w/body include apid/jobid",
			http.MethodPost,
			capmc.NodeEnergyStats,
			bytes.NewBuffer(json.RawMessage(
				`{
					"start_time": "2019-07-10 11:36:32",
					"end_time": "2019-07-10 12:36:32",
					"nids": [0, 1],
					"apid": "831138",
					"job_id": "1234321.sdb"
				}`)),
			http.StatusNotImplemented,
		}, {
			"Post w/body job_id only",
			http.MethodPost,
			capmc.NodeEnergyStats,
			bytes.NewBuffer(json.RawMessage(
				`{
					"job_id": "1234321.sdb"
				}`)),
			http.StatusNotImplemented,
		}, {
			"Delete",
			http.MethodDelete,
			capmc.NodeEnergyStats,
			nil,
			http.StatusMethodNotAllowed,
		}, {
			"Connect",
			http.MethodConnect,
			capmc.NodeEnergyStats,
			nil,
			http.StatusMethodNotAllowed,
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

			//t.Logf("rec: %s", w.Body.String())

			if tt.ret != w.Code {
				t.Errorf("Returned wrong status code: got %v want %v",
					w.Code, tt.ret)
			}
		})
	}
}

func TestDoNodeEnergyCounter(t *testing.T) {
	var tSvc CapmcD
	var err error
	tSvc.hsmURL, err = url.Parse("http://localhost:27779")
	if err != nil {
		t.Fatal(err)
	}
	testClient := NewTestClient(DoNodeEnergyCounterTestFunc())
	tSvc.rfClient = testClient
	tSvc.smClient = testClient
	tSvc.config = loadConfig("")

	logger.SetupLogging()

	tsdb.ConfigureDataImplementation(tsdb.DUMMY)

	handler := http.HandlerFunc(tSvc.doNodeEnergyCounter)

	tests := []struct {
		name   string
		method string
		path   string
		body   io.Reader
		ret    int
	}{
		{
			"Get",
			http.MethodGet,
			capmc.NodeEnergyCounter,
			nil,
			http.StatusMethodNotAllowed,
		}, {
			"Put",
			http.MethodPut,
			capmc.NodeEnergyCounter,
			nil,
			http.StatusMethodNotAllowed,
		}, {
			"Post empty",
			http.MethodPost,
			capmc.NodeEnergyCounter,
			bytes.NewBufferString(""),
			http.StatusBadRequest,
		}, {
			"Post empty body",
			http.MethodPost,
			capmc.NodeEnergyCounter,
			bytes.NewBuffer(json.RawMessage(`{}`)),
			http.StatusBadRequest,
		}, {
			"Post w/body apid only",
			http.MethodPost,
			capmc.NodeEnergyCounter,
			bytes.NewBuffer(json.RawMessage(
				`{
					"apid": "831138"
				}`)),
			http.StatusNotImplemented,
		}, {
			"Post w/body times + nids = happy path",
			http.MethodPost,
			capmc.NodeEnergyCounter,
			bytes.NewBuffer(json.RawMessage(
				`{
					"time": "2019-07-10 11:36:32",
					"nids": [0, 1]
				}`)),
			http.StatusOK,
		}, {
			"Post w/body include apid",
			http.MethodPost,
			capmc.NodeEnergyCounter,
			bytes.NewBuffer(json.RawMessage(
				`{
					"time": "2019-07-10 11:36:32",
					"nids": [0, 1],
					"apid": "831138"
				}`)),
			http.StatusNotImplemented,
		}, {
			"Post w/body include jobid",
			http.MethodPost,
			capmc.NodeEnergyCounter,
			bytes.NewBuffer(json.RawMessage(
				`{
					"time": "2019-07-10 11:36:32",
					"nids": [40, 41, 42, 43],
					"job_id": "1234321.sdb"
				}`)),
			http.StatusNotImplemented,
		}, {
			"Post w/body include apid/jobid",
			http.MethodPost,
			capmc.NodeEnergyCounter,
			bytes.NewBuffer(json.RawMessage(
				`{
					"time": "2019-07-10 11:36:32",
					"nids": [0, 1],
					"apid": "831138",
					"job_id": "1234321.sdb"
				}`)),
			http.StatusNotImplemented,
		}, {
			"Post w/body job_id only",
			http.MethodPost,
			capmc.NodeEnergyCounter,
			bytes.NewBuffer(json.RawMessage(
				`{
					"job_id": "1234321.sdb"
				}`)),
			http.StatusNotImplemented,
		}, {
			"Delete",
			http.MethodDelete,
			capmc.NodeEnergyCounter,
			nil,
			http.StatusMethodNotAllowed,
		}, {
			"Connect",
			http.MethodConnect,
			capmc.NodeEnergyCounter,
			nil,
			http.StatusMethodNotAllowed,
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

			if tt.ret != w.Code {
				t.Errorf("Returned wrong status code: got %v want %v",
					w.Code, tt.ret)
			}
		})
	}
}
