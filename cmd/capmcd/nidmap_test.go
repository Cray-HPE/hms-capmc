// Copyright 2019 Cray Inc. All Rights Reserved.
//
// Except as permitted by contract or express written permission of Cray Inc.,
// no part of this work or its content may be modified, used, reproduced or
// disclosed in any form. Modifications made without express permission of
// Cray Inc. may damage the system the software is installed within, may
// disqualify the user from receiving support from Cray Inc. under support or
// maintenance contracts, or require additional support services outside the
// scope of those contracts to repair the software or system.

package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"stash.us.cray.com/HMS/hms-capmc/internal/capmc"
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
