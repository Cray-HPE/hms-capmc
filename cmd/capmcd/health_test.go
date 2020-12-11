// Copyright 2020 Hewlett Packard Enterprise Development LP

package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"stash.us.cray.com/HMS/hms-capmc/internal/capmc"
	compcreds "stash.us.cray.com/HMS/hms-compcredentials"
	sstorage "stash.us.cray.com/HMS/hms-securestorage"
)

func TestCapmcdLiveness(t *testing.T) {
	// configure an instance with synthetic hsm
	var tSvc CapmcD
	var err error
	tSvc.hsmURL, err = url.Parse("http://localhost:27779")
	if err != nil {
		t.Fatal(err)
	}
	testClient := NewTestClient(hsmHealthSynthTestFunc())
	tSvc.rfClient = testClient
	tSvc.smClient = testClient
	tSvc.config = loadConfig("")

	// create a sythetic vault - just needs to be there, not functional
	ss, adapter := sstorage.NewMockAdapter()
	ccs := compcreds.NewCompCredStore("secret/hms-cred", ss)
	tSvc.ss = ss
	tSvc.ccs = ccs
	adapter.LookupData = ssDataLiveness

	// set up liveness request
	handler := http.HandlerFunc(tSvc.doLiveness)

	// set up tests
	tests := []struct {
		name   string
		method string
		rc     int
	}{
		{
			"Valid liveness request",
			http.MethodGet,
			http.StatusNoContent,
		},
		{
			"Invalid liveness request - POST",
			http.MethodPost,
			http.StatusMethodNotAllowed,
		},
		{
			"Invalid liveness request - PUT",
			http.MethodPut,
			http.StatusMethodNotAllowed,
		},
	}

	// run the tests
	for _, tt := range tests {
		adapter.LookupNum = 0
		t.Run(tt.name, func(t *testing.T) {
			// create the http request
			req, err := http.NewRequest(tt.method, capmc.Liveness, nil)
			if err != nil {
				t.Fatal(err)
			}

			// run the request
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if tt.rc != rr.Code {
				t.Errorf("Expected return code: %d, got: %d", tt.rc, rr.Code)
			}
		})
	}
}

func TestCapmcdReadiness(t *testing.T) {
	// configure an instance with synthetic hsm
	var tSvc CapmcD
	var err error
	tSvc.hsmURL, err = url.Parse("http://localhost:27779")
	if err != nil {
		t.Fatal(err)
	}
	testClient := NewTestClient(hsmHealthSynthTestFunc())
	tSvc.rfClient = testClient
	tSvc.smClient = testClient
	tSvc.config = loadConfig("")

	// create a sythetic vault - just needs to be there, not functional
	ss, adapter := sstorage.NewMockAdapter()
	ccs := compcreds.NewCompCredStore("secret/hms-cred", ss)
	tSvc.ss = ss
	tSvc.ccs = ccs
	adapter.LookupData = ssDataLiveness

	// set up liveness request
	handler := http.HandlerFunc(tSvc.doReadiness)

	// set up tests
	tests := []struct {
		name   string
		method string
		rc     int
	}{
		{
			"Valid readiness request",
			http.MethodGet,
			http.StatusNoContent,
		},
		{
			"Invalid readiness request - POST",
			http.MethodPost,
			http.StatusMethodNotAllowed,
		},
		{
			"Invalid readiness request - PUT",
			http.MethodPut,
			http.StatusMethodNotAllowed,
		},
	}

	// run the tests
	for _, tt := range tests {
		adapter.LookupNum = 0
		t.Run(tt.name, func(t *testing.T) {
			// create the http request
			req, err := http.NewRequest(tt.method, capmc.Liveness, nil)
			if err != nil {
				t.Fatal(err)
			}

			// run the request
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if tt.rc != rr.Code {
				t.Errorf("Expected return code: %d, got: %d", tt.rc, rr.Code)
			}
		})
	}
}

func TestCapmcdHealth(t *testing.T) {
	// configure an instance with synthetic hsm
	var tSvc CapmcD
	var err error
	tSvc.hsmURL, err = url.Parse("http://localhost:27779")
	if err != nil {
		t.Fatal(err)
	}
	testClient := NewTestClient(hsmHealthSynthTestFunc())
	tSvc.rfClient = testClient
	tSvc.smClient = testClient
	tSvc.config = loadConfig("")

	// create a sythetic vault - just needs to be there, not functional
	ss, adapter := sstorage.NewMockAdapter()
	ccs := compcreds.NewCompCredStore("secret/hms-cred", ss)
	tSvc.ss = ss
	tSvc.ccs = ccs
	adapter.LookupData = ssDataLiveness

	// set up liveness request
	handler := http.HandlerFunc(tSvc.doHealth)

	// set up tests
	tests := []struct {
		name   string
		method string
		rc     int
	}{
		{
			"Valid health request",
			http.MethodGet,
			http.StatusOK,
		},
		{
			"Invalid health request - POST",
			http.MethodPost,
			http.StatusMethodNotAllowed,
		},
		{
			"Invalid health request - PUT",
			http.MethodPut,
			http.StatusMethodNotAllowed,
		},
	}

	// run the tests
	for _, tt := range tests {
		adapter.LookupNum = 0
		t.Run(tt.name, func(t *testing.T) {
			// create the http request
			req, err := http.NewRequest(tt.method, capmc.Liveness, nil)
			if err != nil {
				t.Fatal(err)
			}

			// run the request
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if tt.rc != rr.Code {
				t.Errorf("Expected return code: %d, got: %d", tt.rc, rr.Code)
			}
		})
	}
}

// hsmHealthSynthTestFunc - provides simplest implementation of synthetic state manager
func hsmHealthSynthTestFunc() RoundTripFunc {
	return func(req *http.Request) (*http.Response, error) {
		// should not need anything other than 'alive'
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBufferString("")),
			Header:     make(http.Header),
		}, nil
	}
}

// mock storage data for liveness tests
var ssDataLiveness = []sstorage.MockLookup{
	{
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname: "x1000c0",
			},
		},
	},
}
