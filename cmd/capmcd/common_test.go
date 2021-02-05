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
