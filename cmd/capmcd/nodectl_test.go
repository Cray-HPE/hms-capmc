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
	"net/http"
	"net/http/httptest"
	"testing"

	"stash.us.cray.com/HMS/hms-capmc/internal/capmc"
)

const (
	defaultNodeRulesJSON = `{"e":0,"err_msg":"","latency_node_off":60,"latency_node_on":120,"latency_node_reinit":180,"max_off_req_count":-1,"max_off_time":-1,"max_on_req_count":-1,"max_reinit_req_count":-1,"min_off_time":-1}`
)

func TestDoNodeRules(t *testing.T) {
	var svc CapmcD
	svc.config = loadConfig("")
	handler := http.HandlerFunc(svc.doNodeRules)

	tests := []struct {
		name     string
		method   string
		path     string
		body     io.Reader
		code     int
		expected string
	}{
		{
			"Get",
			http.MethodGet,
			capmc.NodeRules,
			nil,
			http.StatusOK,
			defaultNodeRulesJSON + "\n",
		},
		{
			"Put",
			http.MethodPut,
			capmc.NodeRules,
			nil,
			http.StatusMethodNotAllowed,
			"{\"e\":405,\"err_msg\":\"(PUT) Not Allowed\"}\n",
		},
		{
			"Post (empty)",
			http.MethodPost,
			capmc.NodeRules,
			bytes.NewBufferString(""),
			http.StatusBadRequest,
			"{\"e\":400,\"err_msg\":\"Bad Request: JSON: EOF\"}\n",
		},
		{
			"Post (w/ empty body)",
			http.MethodPost,
			capmc.NodeRules,
			bytes.NewBuffer(json.RawMessage(`{}`)),
			http.StatusOK,
			defaultNodeRulesJSON + "\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, tc.path, tc.body)
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
