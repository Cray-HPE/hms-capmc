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
	"io/ioutil"
	"path/filepath"
	"testing"
)

const (
	// 200 - Success
	SuccessJSON = `{"e":0,"err_msg":""}`

	// 400 - Bad Request response
	NoRequestJSON           = `{"e":400,"err_msg":"no request"}`
	MissingXnameJSON        = `{"e":400,"err_msg":"missing valid target xname(s)"}`
	InvalidComponentJSON    = `{"e":400,"err_msg":"failure validating xnames: Invalid component"}`
	FailedCompGenJSON       = `{"e":400,"err_msg":"failure generating component list: "}`
	UnexpectedEOFJSON       = `{"e":400,"err_msg":"unexpected EOF"}`
	InvalidCharJSON         = `{"e":400,"err_msg":"invalid character '}' looking for beginning of value"}`
	XnameControlUnmarshJSON = `{"e":400,"err_msg":"json: cannot unmarshal string into Go value of type capmc.XnameControl"}`
	InvalidFilterStringJSON = `{"e":400,"err_msg":"invalid filter string: "}`

	// 405 - Method Not Allowed response
	GetNotAllowedJSON     = `{"e":405,"err_msg":"(GET) Not Allowed"}`
	PatchNotAllowedJSON   = `{"e":405,"err_msg":"(PATCH) Not Allowed"}`
	PostNotAllowedJSON    = `{"e":405,"err_msg":"(POST) Not Allowed"}`
	PutNotAllowedJSON     = `{"e":405,"err_msg":"(PUT) Not Allowed"}`
	DeleteNotAllowedJSON  = `{"e":405,"err_msg":"(DELETE) Not Allowed"}`
	ConnectNotAllowedJSON = `{"e":405,"err_msg":"(CONNECT) Not Allowed"}`

	// 501 - Not Implemented
	UnavailableJSON = `{"e":501,"err_msg":"Unavailable - Not Implemented"}`
)

// helperLoadBytes loads test data from name file as an array of bytes
// Idea/code from https://medium.com/@povilasve/go-advanced-tips-tricks-a872503ac859
func loadTestDataBytes(t *testing.T, name string) []byte {
	path := filepath.Join("testdata", name) // relative path
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return bytes
}
