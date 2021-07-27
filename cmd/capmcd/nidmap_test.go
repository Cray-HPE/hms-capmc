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
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/Cray-HPE/hms-capmc/internal/capmc"
)

func TestDoNidMap(t *testing.T) {
	var svc CapmcD
	svc.hsmURL, _ = url.Parse("http://localhost")
	svc.config = loadConfig("")
	handler := http.HandlerFunc(svc.doNidMap)

	tests := []struct {
		name    string
		method  string
		path    string
		bodyIn  io.Reader
		code    int
		bodyOut string
	}{
		{
			name:    "GET",
			method:  http.MethodGet,
			path:    capmc.NodeIDMap,
			bodyIn:  nil,
			code:    http.StatusMethodNotAllowed,
			bodyOut: GetNotAllowedJSON + "\n",
		},
		{
			name:    "POST",
			method:  http.MethodPost,
			path:    capmc.NodeIDMap,
			bodyIn:  bytes.NewBufferString(""),
			code:    http.StatusBadRequest,
			bodyOut: NoRequestJSON + "\n",
		},
		/*
			   This currently won't work with out additional mocking
			{
				name:    "POST empty JSON",
				method:  http.MethodPost,
				path:    capmc.NodeIDMap,
				bodyIn:  bytes.NewBuffer([]byte(`{}`)),
				code:    http.StatusBadRequest,
				bodyOut: NoRequestJSON + "\n",
			},
		*/
		{
			name:    "POST invalid JSON",
			method:  http.MethodPost,
			path:    capmc.NodeIDMap,
			bodyIn:  bytes.NewBufferString("{nids: []}"),
			code:    http.StatusBadRequest,
			bodyOut: "", // don't care
		},
		{
			name:    "PUT",
			method:  http.MethodPut,
			path:    capmc.NodeIDMap,
			bodyIn:  nil,
			code:    http.StatusMethodNotAllowed,
			bodyOut: PutNotAllowedJSON + "\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, tc.path, tc.bodyIn)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if rr.Code != tc.code {
				t.Errorf("handler returned wrong status code: want %v but got %v",
					tc.code, rr.Code)
			}

			if tc.bodyOut != "" && rr.Body.String() != tc.bodyOut {
				t.Errorf("handler returned unexpected body: want '%v' but got '%v'",
					tc.bodyOut, rr.Body.String())
			}
		})
	}
}
