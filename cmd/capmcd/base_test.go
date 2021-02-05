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
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"

	base "stash.us.cray.com/HMS/hms-base"
	compcreds "stash.us.cray.com/HMS/hms-compcredentials"
	sstorage "stash.us.cray.com/HMS/hms-securestorage"
	"stash.us.cray.com/HMS/hms-certs/pkg/hms_certs"
)

var replayList *[]testData
var mockVault *sstorage.MockAdapter

// The replay function is called to intercept external service requests and
// replay the previously recorded responses rather than contacting the actual
// external service.
func replay(req *http.Request) (*http.Response, error) {
	// Test request parameters
	var (
		body []byte
		resp *http.Response
		err  error
	)
	if req.GetBody != nil {
		bodyReader, err := req.GetBody()
		if err == nil {
			body, err = ioutil.ReadAll(bodyReader)
			if err != nil {
				log.Printf("ReadAll failed: %s", err)
			} else if len(body) > 0 {
				log.Printf("    Body: %s\n", body)
			}
		}
	}
	for _, v := range *replayList {
		if v.reqMethod == req.Method && v.reqURL == req.URL.String() {
			// We've got a match!
			resp = &http.Response{
				Status:     v.respStatus,
				StatusCode: v.respStatusCode,
				Proto:      v.respProto,
				ProtoMajor: v.respProtoMajor,
				ProtoMinor: v.respProtoMinor,
				Header:     v.respHeader,
				Body:       ioutil.NopCloser(strings.NewReader(v.respBody)),
				Request:    req,
			}
			break
		}
	}
	if resp == nil {
		// If we don't find a match, we should return an error.
		// This may indicate a test case failure, but could indicate a change
		// in the operation of CAPMC which invalidates the test case.  The
		// developer will need to determine which is the case.
		err = &url.Error{Op: req.Method, URL: req.URL.String(), Err: errors.New("test replay could not find a match")}
	}

	return resp, err
}

// replayClient returns *http.Client with Transport replaced to replay responses
// for given requests without an external connection.
func replayClient() *hms_certs.HTTPClientPair {
	rc,_ := makeClient(0,5)
	rc.InsecureClient.HTTPClient = &http.Client{ Transport: RoundTripFunc(replay), }
	rc.SecureClient = rc.InsecureClient
	return rc
}

func findHandler(url string) http.HandlerFunc {
	i := strings.Index(url, "/capmc/")
	base := url[i:]
	for _, vers := range capmcAPIs {
		for _, api := range vers {
			if base == api.pattern {
				return api.handler
			}
		}
	}
	return nil
}

func parseResponse(resp string) (map[string]interface{}, []string, error) {
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(resp), &m)
	if err != nil {
		log.Println("Error parsing result: ", err)
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return m, keys, err
}

func genericContains(startIndex int, elem interface{}, slice []interface{}) int {
	for i := startIndex; i < len(slice); i++ {
		// compare elem to the current slice element
		if reflect.DeepEqual(elem, slice[i]) {
			return i
		}
	}
	return -1 // nothing found
}

func compareSlices(t *testing.T, key string, expected, actual []interface{}) {
	if len(expected) != len(actual) {
		t.Errorf("Key %s mismatch:\n\tExpected: '%v'\n\tActual:   '%v'\n", key, expected, actual)
		return
	}
	// We want to compare the two slices to see if they contain
	// the same items even if they are not in the same order.
	m := make(map[int]bool)
	for _, expVal := range expected {
		i := genericContains(0, expVal, actual)
		for i >= 0 {
			if found, ok := m[i]; !ok || !found {
				// Mark this index as found
				m[i] = true
				break
			}
			i = genericContains(i+1, expVal, actual)
		}
		if i < 0 {
			t.Errorf("Mismatched value: cannot find expected value %v in %v\n", expVal, actual)
			return
		}
	}
}

func compareVals(t *testing.T, key string, expected, actual interface{}) {
	switch expVal := expected.(type) {
	case int:
		actVal, ok := actual.(int)
		if !ok {
			t.Errorf("Mismatched value type for key %s: exp: %d, act: %v\n", key, expVal, actual)
			return
		}
		if expVal != actVal {
			t.Errorf("Mismatched value for key %s: exp: %d, act: %d\n", key, expVal, actVal)
			return
		}
	case float64:
		actVal, ok := actual.(float64)
		if !ok {
			t.Errorf("Mismatched value type for key %s: exp: %f, act: %v\n", key, expVal, actual)
			return
		}
		if expVal != actVal {
			t.Errorf("Mismatched value for key %s: exp: %f, act: %f\n", key, expVal, actVal)
			return
		}
	case string:
		actVal, ok := actual.(string)
		if !ok {
			t.Errorf("Mismatched value type for key %s: exp: %s, act: %v\n", key, expVal, actual)
			return
		}
		if expVal != actVal {
			t.Errorf("Mismatched value for key %s: exp: %s, act: %s\n", key, expVal, actVal)
			return
		}
	case []interface{}:
		if actVal, ok := actual.([]interface{}); ok {
			compareSlices(t, key, expVal, actVal)
		} else {
			t.Errorf("Mismatched value type for key %s: exp: %v, act: %v\n", key, expVal, actual)
		}
		return
	default:
		t.Errorf("Unexpected value types for key %s:\n\texpected: '%v'\n\tActual:   '%s'",
			key, expected, actual)
		return
	}
}

// We are expecting all response results to be JSON.  This means we can't just
// do a straight string comparison.
func compareResults(t *testing.T, expected, actual string) {
	expMap, expKeys, expErr := parseResponse(expected)
	actMap, actKeys, actErr := parseResponse(actual)
	if expErr != nil {
		log.Println("Error parsing expected result: ", expErr)
	}
	if actErr != nil {
		log.Println("Error parsing actual result: ", actErr)
	}
	if len(expKeys) != len(actKeys) {
		// We've got a mismatch since the number of keys don't match.
		t.Errorf("Keys mismatch:\n\tExpKeys: '%v'\n\tActKeys: '%v'\n\tExpected: '%v'\n\tActual:   '%v'\n", expKeys, actKeys, expected, actual)
		return
	}
	for i := 0; i < len(expKeys); i++ {
		if expKeys[i] != actKeys[i] {
			// We've got a mismatch since the number of keys don't match.
			t.Errorf("Keys mismatch:\n\tExpKeys: '%v'\n\tActKeys: '%v'\n\tExpected: '%v'\n\tActual:   '%v'\n", expKeys, actKeys, expected, actual)
			return
		}
	}
	// If we get here, we know we have the same keys in our JSON response.
	// Now we need to look at individual values.  We are looking for
	// specific types.
	for _, k := range expKeys {
		compareVals(t, k, expMap[k], actMap[k])
	}

	if debug {
		log.Printf("Expected Results: %+v\n", expMap)
		log.Printf("Actual Results: %+v\n", actMap)
	}
}

// The runTest function is used to run the captured tests.
func runTest(t *testing.T, hsm string, testReq *testData, rlist *[]testData, vaultData []sstorage.MockLookup) {
	replayList = rlist
	var err error
	// Set the HSM URL for this test.
	if svc.hsmURL, err = url.Parse(hsm + "/hsm/v1"); err != nil {
		log.Fatalf("Invalid HSM URI specified: %s", err)
	}
	svc.config = loadConfig("")
	mockVault.LookupNum = 0
	mockVault.LookupData = vaultData
	handler := findHandler(testReq.reqURL)
	req := httptest.NewRequest(testReq.reqMethod, testReq.reqURL,
		ioutil.NopCloser(strings.NewReader(testReq.reqBody)))
	w := httptest.NewRecorder()
	handler(w, req)
	resp := w.Result()
	// Need to verify correctness, for now just print...
	body, _ := ioutil.ReadAll(resp.Body)
	compareResults(t, testReq.respBody, string(body))
	//log.Printf("Resp: %v\n", resp)
	//log.Printf("Resp Body: %s\n", body)
}

var (
	debug     bool
	enableLog bool
)

func init() {
	flag.BoolVar(&debug, "debug", false, "set internal debug flag")
	flag.BoolVar(&enableLog, "enable-log", false, "turns on standrd log output")
}

func TestMain(m *testing.M) {
	flag.Parse()
	// Need to get some base initialization done first.
	var err error
	svc.httpListen = ":27777"
	svc.rfClient = replayClient()
	svc.smClient = replayClient()
	// Spin up our global worker goroutine pool.
	svc.WPool = base.NewWorkerPool(10, 1000)
	svc.WPool.Run()
	// Use the mock secure storage
	ss, adapter := sstorage.NewMockAdapter()
	ccs := compcreds.NewCompCredStore("secret/hms-cred", ss)
	svc.ss = ss
	svc.ccs = ccs
	svc.debug = debug
	mockVault = adapter
	hsm := "https://localhost:27779/hsm/v1"
	if svc.hsmURL, err = url.Parse(hsm); err != nil {
		log.Fatalf("Invalid HSM URI specified: %s", err)
	}

	// debug is on so turn off log discard
	if debug && !enableLog {
		enableLog = true
	}

	if enableLog {
		// setup logger as it is in non-test main
		log.SetFlags(log.Lshortfile | log.LstdFlags)
	} else {
		// discard all log output
		log.SetOutput(ioutil.Discard)
	}
	excode := m.Run()
	os.Exit(excode)
}
