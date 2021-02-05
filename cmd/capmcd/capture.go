// Copyright 2019,2020 Hewlett Packard Enterprise Development LP
//
// The second info to capture is interaction with external services, namely,
// HSM and redfish.  By capturing this data, we are able to replay it later
// when running a test case, which allows us to verify that CAPMC is still
// functioning correctly given the same inputs from these external services.
//
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"time"
	"stash.us.cray.com/HMS/hms-certs/pkg/hms_certs"
)

type testData struct {
	// Request items...
	reqURL    string
	reqMethod string
	reqBody   string
	reqHeader map[string][]string

	// Response items...
	respStatus     string
	respStatusCode int
	respProto      string
	respProtoMajor int
	respProtoMinor int
	respHeader     map[string][]string
	// Might want to add a Trailer as well?
	respBody string
	// Perhaps need something for returning an error as well at some point...
}

const capturedTestsFile = "captured_test.go"

var (
	realClient   *hms_certs.HTTPClientPair
	captureFile  *os.File
	captureMutex sync.Mutex
	captureTag   string
)

// This function initializes the capture mechanism.
// It opens and sets up the capture file.
func captureInit(tag string) {
	if tag != "" {
		log.Printf("captureInit: tag: %s\n", tag)
		captureTag = tag
		f, err := os.OpenFile(capturedTestsFile, os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			// If file does not exist, the open will fail, so we try again
			f, err = os.OpenFile(capturedTestsFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
			fmt.Fprintf(f, "// This file is auto-generated via the capture facility\n\n")
			fmt.Fprintf(f, "package main\n\nimport \"testing\"\n")
		}
		if err != nil {
			log.Printf("Could not open capture file %s: %s", capturedTestsFile, err)
		} else {
			captureMutex.Lock()
			defer captureMutex.Unlock()
			captureFile = f
			hsmFlag := flag.Lookup("hsm")
			if hsmFlag != nil {
				val := hsmFlag.Value.String()
				if val == "" {
					val = hsmFlag.DefValue
				}
				fmt.Fprintf(f, "\nvar %sHSM = \"%s\"\n", tag, val)
			}
			fmt.Fprintf(f, "var %sReplayData = []testData{\n", tag)
			captureFile.Sync()
		}
	}
}

// Set up the handler to intercept the incoming request and also capture the response.
// The customized in place function is needed to pass the normal handler to captureHandler
func captureHandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		captureHandler(w, r, handler)
	})
}

// The captureHandler captures an incoming request to the CAPMC service.  It
// then calls the normal handler, and records the response.
func captureHandler(w http.ResponseWriter, r *http.Request, handler http.HandlerFunc) {
	var td testData
	td.SaveReq(r)
	if td.reqBody == "" {
		// Make sure we capture the body of the request.
		reqBody, _ := ioutil.ReadAll(r.Body)
		td.reqBody = string(reqBody)
		r.Body = ioutil.NopCloser(bytes.NewReader(reqBody))
	}
	wRecorder := httptest.NewRecorder()
	handler(wRecorder, r)
	resp := wRecorder.Result()
	td.SaveResp(resp)
	for hdr, vals := range resp.Header {
		for _, v := range vals {
			w.Header().Set(hdr, v)
		}
	}
	w.WriteHeader(resp.StatusCode)
	w.Write([]byte(td.respBody))
	if captureFile != nil {
		captureMutex.Lock()
		defer captureMutex.Unlock()
		fmt.Fprintf(captureFile, "}\n")
		fmt.Fprintf(captureFile, "\nfunc Test%s(t *testing.T) {\n", strings.Title(captureTag))
		fmt.Fprintf(captureFile, "\tvar testReq = testData")
		td.Write(captureFile, "\t")
		fmt.Fprintf(captureFile, "\n\trunTest(t, %sHSM, ", captureTag)
		fmt.Fprintf(captureFile, "&testReq, &%sReplayData)\n}\n", captureTag)
		captureFile.Sync()
	}
}

// Save the request data into a testData object
func (td *testData) SaveReq(req *http.Request) {
	td.reqURL = req.URL.String()
	td.reqMethod = req.Method
	td.reqHeader = req.Header
	if req.GetBody != nil {
		bodyReader, err := req.GetBody()
		if err == nil {
			body, err := ioutil.ReadAll(bodyReader)
			if err != nil {
				log.Printf("ReadAll failed: %s", err)
			} else if len(body) > 0 {
				td.reqBody = string(body)
			}
		}
	}
}

// Save the response data into a testData object
func (td *testData) SaveResp(resp *http.Response) {
	td.respStatus = resp.Status
	td.respStatusCode = resp.StatusCode
	td.respProto = resp.Proto
	td.respProtoMajor = resp.ProtoMajor
	td.respProtoMinor = resp.ProtoMinor
	td.respHeader = resp.Header
	// Might want to add a Trailer as well?
	body, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		// Replace the Body reader since we just read the body.
		resp.Body.Close()
		resp.Body = ioutil.NopCloser(bytes.NewReader(body))
		td.respBody = string(body)
	}
}

// This support function writes out the http header data to a file.
func (td *testData) writeHeader(w io.Writer, header *map[string][]string) {
	fmt.Fprintf(w, "map[string][]string{")
	sep := ""
	for hdr, vals := range *header {
		fmt.Fprintf(w, "%s\"%s\": []string{", sep, hdr)
		sep = ", "
		s := ""
		for _, v := range vals {
			fmt.Fprintf(w, "%s\"%s\"", s, v)
			s = ","
		}
		fmt.Fprintf(w, "}")
	}
	fmt.Fprintf(w, "},\n")
}

// Write out the test data to a file.
func (td *testData) Write(w io.Writer, pad string) {
	fmt.Fprintf(w, "{\"%s\",\n", td.reqURL)
	fmt.Fprintf(w, pad+"\t\"%s\",\n", td.reqMethod)
	fmt.Fprintf(w, pad+"\t`%s`,\n", td.reqBody)
	fmt.Fprintf(w, pad+"\t")
	td.writeHeader(w, &td.reqHeader)
	fmt.Fprintf(w, pad+"\t\"%s\", %d,\n", td.respStatus, td.respStatusCode)
	fmt.Fprintf(w, pad+"\t\"%s\", %d, %d,\n", td.respProto, td.respProtoMajor, td.respProtoMinor)
	fmt.Fprintf(w, pad+"\t")
	td.writeHeader(w, &td.respHeader)
	fmt.Fprintf(w, pad+"\t`%s`}", td.respBody)
}

type RoundTripFunc func(req *http.Request) (*http.Response, error)

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// Here is the function which will receive all requests destined for external
// clients.  We record the request, then forward the request to the external
// client.  Once we receive the response, we record that as well.  We then save
// the client interaction for our test case data, and also return it to the
// original caller.  So the caller appears to be interacting directly with the
// client, but we can record this interaction for later replay as a test case.
func capture(req *http.Request) (*http.Response, error) {
	// Test request parameters
	var td testData
	td.SaveReq(req)

	resp, err := realClient.Do(req)
	captureMutex.Lock()
	defer captureMutex.Unlock()
	td.SaveResp(resp)

	if captureFile != nil {
		fmt.Fprintf(captureFile, "\t")
		td.Write(captureFile, "\t")
		fmt.Fprintf(captureFile, ",\n")
		captureFile.Sync()
	}

	if err != nil {
		log.Printf("   Error: %s\n", err)
	}
	return resp, err
}

// captureClient returns *http.Client with Transport replaced to funnel all
// requests/responses through a the capture function.  The client to talk to
// the real services is also constructed.
func captureClient(flags clientFlags, timeout time.Duration) *hms_certs.HTTPClientPair {
	realClient,_ = makeClient((flags|clientInsecure), timeout)

	realClient.InsecureClient.HTTPClient = &http.Client{
        Transport: RoundTripFunc(capture),
    }
	realClient.SecureClient = realClient.InsecureClient
	return realClient
}
