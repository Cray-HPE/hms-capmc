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
	"strings"
	"testing"
	"time"

	"github.com/Cray-HPE/hms-capmc/internal/capmc"
	"github.com/Cray-HPE/hms-capmc/internal/logger"
	"github.com/Cray-HPE/hms-capmc/internal/tsdb"
)

const (
	BadJSONJSON             = `{"e":400,"err_msg":"BAD_JSON"}`
	BadStartTimeJSON        = `{"e":400,"err_msg":"BAD_START_TIME"}`
	BadWindowLenJSON        = `{"e":400,"err_msg":"BAD_WINDOW_LEN"}`
	WindowLenOutOfRangeJSON = `{"e":400,"err_msg":"WINDOW_LEN_OUT_OF_RANGE"}`
)

func TestInRange(t *testing.T) {
	var tests = []struct {
		i   int
		min int
		max int
		r   bool
	}{
		{
			i:   1,
			min: 1,
			max: 10,
			r:   true,
		},
		{
			i:   5,
			min: 1,
			max: 10,
			r:   true,
		},
		{
			i:   10,
			min: 1, max: 10, r: true},
		{
			i:   0,
			min: 1,
			max: 10,
			r:   false,
		},
		{
			i:   11,
			min: 1,
			max: 10,
			r:   false,
		},
		{
			i:   -100,
			min: 1,
			max: 10,
			r:   false,
		},
		{
			i:   100,
			min: 1,
			max: 10,
			r:   false,
		},
	}

	for n, test := range tests {
		r := inRange(test.i, test.min, test.max)
		if r != test.r {
			t.Errorf("InRange Test Case %d: FAIL: Expected %v but got %v", n, test.r, r)
		}
	}
}

func TestDoSystemParams(t *testing.T) {
	var svc CapmcD
	svc.config = loadConfig("")
	handler := http.HandlerFunc(svc.doSystemParams)

	tests := []struct {
		name     string
		method   string
		path     string
		body     io.Reader
		code     int
		expected string
	}{
		{
			name:     "GET",
			method:   http.MethodGet,
			path:     capmc.SystemParams,
			body:     nil,
			code:     http.StatusOK,
			expected: `{"e":0,"err_msg":"","power_cap_target":0,"power_threshold":0,"static_power":0,"ramp_limited":false,"ramp_limit":2000000,"power_band_min":0,"power_band_max":0}` + "\n",
		},
		{
			name:     "POST",
			method:   http.MethodPost,
			path:     capmc.SystemParams,
			body:     bytes.NewBufferString(""),
			code:     http.StatusOK,
			expected: `{"e":0,"err_msg":"","power_cap_target":0,"power_threshold":0,"static_power":0,"ramp_limited":false,"ramp_limit":2000000,"power_band_min":0,"power_band_max":0}` + "\n",
		},
		{
			name:     "POST empty JSON",
			method:   http.MethodPost,
			path:     capmc.SystemParams,
			body:     bytes.NewBuffer([]byte(`{}`)),
			code:     http.StatusOK,
			expected: `{"e":0,"err_msg":"","power_cap_target":0,"power_threshold":0,"static_power":0,"ramp_limited":false,"ramp_limit":2000000,"power_band_min":0,"power_band_max":0}` + "\n",
		},
		{
			name:     "PUT",
			method:   http.MethodPut,
			path:     capmc.SystemParams,
			body:     nil,
			code:     http.StatusMethodNotAllowed,
			expected: PutNotAllowedJSON + "\n",
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

func TestDoSystemPower(t *testing.T) {
	var svc CapmcD

	logger.SetupLogging()
	tsdb.ConfigureDataImplementation(tsdb.DUMMY)

	handler := http.HandlerFunc(svc.doSystemPower)

	tests := []struct {
		name     string
		method   string
		path     string
		body     io.Reader
		code     int
		expected string
	}{
		{
			name:     "GET",
			method:   http.MethodGet,
			path:     capmc.SystemPower,
			body:     nil,
			code:     http.StatusOK,
			expected: `{"e":0,"err_msg":"","window_len":15,"start_time":"@NOW@","avg":50,"max":100,"min":1}` + "\n",
		},
		{
			name:     "GET start_time valid",
			method:   http.MethodGet,
			path:     capmc.SystemPower + "?" + "start_time=2000-04-01 01:02:03",
			body:     nil,
			code:     http.StatusOK,
			expected: `{"e":0,"err_msg":"","window_len":15,"start_time":"2000-04-01 01:02:03","avg":50,"max":100,"min":1}` + "\n",
		},
		{
			name:     "GET start_time now",
			method:   http.MethodGet,
			path:     capmc.SystemPower + "?" + "start_time=now",
			body:     nil,
			code:     http.StatusOK,
			expected: `{"e":0,"err_msg":"","window_len":15,"start_time":"@NOW@","avg":50,"max":100,"min":1}` + "\n",
		},
		{
			name:     "GET start_time invalid",
			method:   http.MethodGet,
			path:     capmc.SystemPower + "?" + "start_time=2019-05-17 06:48:00 AND",
			body:     nil,
			code:     http.StatusBadRequest,
			expected: BadStartTimeJSON + "\n",
		},
		{
			name:     "GET window_len valid",
			method:   http.MethodGet,
			path:     capmc.SystemPower + "?" + "window_len=15",
			body:     nil,
			code:     http.StatusOK,
			expected: `{"e":0,"err_msg":"","window_len":15,"start_time":"@NOW@","avg":50,"max":100,"min":1}` + "\n",
		},
		{
			name:     "GET window_len below range",
			method:   http.MethodGet,
			path:     capmc.SystemPower + "?" + "window_len=1",
			body:     nil,
			code:     http.StatusBadRequest,
			expected: WindowLenOutOfRangeJSON + "\n",
		},
		{
			name:     "GET window_len above range",
			method:   http.MethodGet,
			path:     capmc.SystemPower + "?" + "window_len=3601",
			body:     nil,
			code:     http.StatusBadRequest,
			expected: WindowLenOutOfRangeJSON + "\n",
		},
		{
			name:     "GET window_len not a number",
			method:   http.MethodGet,
			path:     capmc.SystemPower + "?" + "window_len=\"abc\"",
			body:     nil,
			code:     http.StatusBadRequest,
			expected: BadWindowLenJSON + "\n",
		},
		{
			name:     "POST",
			method:   http.MethodPost,
			path:     capmc.SystemPower,
			body:     bytes.NewBufferString(""),
			code:     http.StatusBadRequest,
			expected: NoRequestJSON + "\n",
		},
		{
			name:     "POST empty JSON",
			method:   http.MethodPost,
			path:     capmc.SystemPower,
			body:     bytes.NewBuffer([]byte(`{}`)),
			code:     http.StatusOK,
			expected: `{"e":0,"err_msg":"","window_len":15,"start_time":"@NOW@","avg":50,"max":100,"min":1}` + "\n",
		},
		{
			name:     "POST invalid JSON",
			method:   http.MethodPost,
			path:     capmc.SystemPower,
			body:     bytes.NewBuffer([]byte(`{"window_len":,"start_time":}`)),
			code:     http.StatusBadRequest,
			expected: BadJSONJSON + "\n",
		},
		{
			name:     "POST start_time valid",
			method:   http.MethodPost,
			path:     capmc.SystemPower,
			body:     bytes.NewBuffer([]byte(`{"start_time":"2000-04-01 01:02:03"}`)),
			code:     http.StatusOK,
			expected: `{"e":0,"err_msg":"","window_len":15,"start_time":"2000-04-01 01:02:03","avg":50,"max":100,"min":1}` + "\n",
		},
		{
			name:     "POST start_time now",
			method:   http.MethodPost,
			path:     capmc.SystemPower,
			body:     bytes.NewBuffer([]byte(`{"start_time":"now"}`)),
			code:     http.StatusOK,
			expected: `{"e":0,"err_msg":"","window_len":15,"start_time":"@NOW@","avg":50,"max":100,"min":1}` + "\n",
		},
		{
			name:     "POST start_time invalid",
			method:   http.MethodPost,
			path:     capmc.SystemPower,
			body:     bytes.NewBuffer([]byte(`{"start_time":"2019-05-17 06:48:00 AND"}`)),
			code:     http.StatusBadRequest,
			expected: BadStartTimeJSON + "\n",
		},
		{
			name:     "POST window_len valid",
			method:   http.MethodPost,
			path:     capmc.SystemPower,
			body:     bytes.NewBuffer([]byte(`{"window_len":15}`)),
			code:     http.StatusOK,
			expected: `{"e":0,"err_msg":"","window_len":15,"start_time":"@NOW@","avg":50,"max":100,"min":1}` + "\n",
		},
		{
			name:     "POST window_len below range",
			method:   http.MethodPost,
			path:     capmc.SystemPower,
			body:     bytes.NewBuffer([]byte(`{"window_len":1}`)),
			code:     http.StatusBadRequest,
			expected: WindowLenOutOfRangeJSON + "\n",
		},
		{
			name:     "POST window_len above range",
			method:   http.MethodPost,
			path:     capmc.SystemPower,
			body:     bytes.NewBuffer([]byte(`{"window_len":3601}`)),
			code:     http.StatusBadRequest,
			expected: WindowLenOutOfRangeJSON + "\n",
		},
		{
			name:     "POST window_len not a number",
			method:   http.MethodPost,
			path:     capmc.SystemPower,
			body:     bytes.NewBuffer([]byte(`{"window_len":"NaN"}`)),
			code:     http.StatusBadRequest,
			expected: BadWindowLenJSON + "\n",
		},
		{
			name:     "PUT",
			method:   http.MethodPut,
			path:     capmc.SystemPower,
			body:     nil,
			code:     http.StatusMethodNotAllowed,
			expected: PutNotAllowedJSON + "\n",
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

			body := rr.Body.String()

			if strings.Contains(tc.expected, "@NOW@") {
				now := time.Now()
				nowStr := now.Format(intervalTimeFormat)
				tc.expected = strings.Replace(tc.expected, "@NOW@", nowStr, 1)
				i := strings.Index(body, "start_time\":")
				i = i + len("start_time\":") + 1
				dt := body[i : i+len(intervalTimeFormat)]
				body = strings.Replace(body, dt, nowStr, 1)
			}

			if body != tc.expected {
				t.Errorf("handler returned unexpected body: want '%v' but got '%v'",
					tc.expected, body)
			}
		})
	}
}

func TestDoSystemPowerDetails(t *testing.T) {
	var svc CapmcD
	logger.SetupLogging()
	tsdb.ConfigureDataImplementation(tsdb.DUMMY)
	handler := http.HandlerFunc(svc.doSystemPowerDetails)

	tests := []struct {
		name     string
		method   string
		path     string
		body     io.Reader
		code     int
		expected string
	}{
		{
			name:     "GET",
			method:   http.MethodGet,
			path:     capmc.SystemPowerDetails,
			body:     nil,
			code:     http.StatusOK,
			expected: `{"e":0,"err_msg":"","window_len":15,"start_time":"@NOW@","cabinets":[{"avg":0,"max":0,"min":0,"x":0,"y":0},{"avg":1.5,"max":3,"min":1,"x":1,"y":0}]}` + "\n",
		},
		{
			name:     "GET start_time valid",
			method:   http.MethodGet,
			path:     capmc.SystemPowerDetails + "?" + "start_time=2000-04-01 01:02:03",
			body:     nil,
			code:     http.StatusOK,
			expected: `{"e":0,"err_msg":"","window_len":15,"start_time":"2000-04-01 01:02:03","cabinets":[{"avg":0,"max":0,"min":0,"x":0,"y":0},{"avg":1.5,"max":3,"min":1,"x":1,"y":0}]}` + "\n",
		},
		{
			name:     "GET start_time now",
			method:   http.MethodGet,
			path:     capmc.SystemPowerDetails + "?" + "start_time=now",
			body:     nil,
			code:     http.StatusOK,
			expected: `{"e":0,"err_msg":"","window_len":15,"start_time":"@NOW@","cabinets":[{"avg":0,"max":0,"min":0,"x":0,"y":0},{"avg":1.5,"max":3,"min":1,"x":1,"y":0}]}` + "\n",
		},
		{
			name:     "GET start_time invalid",
			method:   http.MethodGet,
			path:     capmc.SystemPowerDetails + "?" + "start_time=2019-05-17 06:48:00 AND",
			body:     nil,
			code:     http.StatusBadRequest,
			expected: BadStartTimeJSON + "\n",
		},
		{
			name:     "GET window_len valid",
			method:   http.MethodGet,
			path:     capmc.SystemPowerDetails + "?" + "window_len=15",
			body:     nil,
			code:     http.StatusOK,
			expected: `{"e":0,"err_msg":"","window_len":15,"start_time":"@NOW@","cabinets":[{"avg":0,"max":0,"min":0,"x":0,"y":0},{"avg":1.5,"max":3,"min":1,"x":1,"y":0}]}` + "\n",
		},
		{
			name:     "GET window_len below range",
			method:   http.MethodGet,
			path:     capmc.SystemPowerDetails + "?" + "window_len=1",
			body:     nil,
			code:     http.StatusBadRequest,
			expected: WindowLenOutOfRangeJSON + "\n",
		},
		{
			name:     "GET window_len above range",
			method:   http.MethodGet,
			path:     capmc.SystemPowerDetails + "?" + "window_len=3601",
			body:     nil,
			code:     http.StatusBadRequest,
			expected: WindowLenOutOfRangeJSON + "\n",
		},
		{
			name:     "GET window_len not a number",
			method:   http.MethodGet,
			path:     capmc.SystemPowerDetails + "?" + "window_len=\"abc\"",
			body:     nil,
			code:     http.StatusBadRequest,
			expected: BadWindowLenJSON + "\n",
		},
		{
			name:     "POST",
			method:   http.MethodPost,
			path:     capmc.SystemPowerDetails,
			body:     bytes.NewBufferString(""),
			code:     http.StatusBadRequest,
			expected: NoRequestJSON + "\n",
		},
		{
			name:     "POST empty JSON",
			method:   http.MethodPost,
			path:     capmc.SystemPowerDetails,
			body:     bytes.NewBuffer([]byte(`{}`)),
			code:     http.StatusOK,
			expected: `{"e":0,"err_msg":"","window_len":15,"start_time":"@NOW@","cabinets":[{"avg":0,"max":0,"min":0,"x":0,"y":0},{"avg":1.5,"max":3,"min":1,"x":1,"y":0}]}` + "\n",
		},
		{
			name:     "POST invalid JSON",
			method:   http.MethodPost,
			path:     capmc.SystemPowerDetails,
			body:     bytes.NewBuffer([]byte(`{"window_len":,"start_time":}`)),
			code:     http.StatusBadRequest,
			expected: BadJSONJSON + "\n",
		},
		{
			name:     "POST start_time valid",
			method:   http.MethodPost,
			path:     capmc.SystemPowerDetails,
			body:     bytes.NewBuffer([]byte(`{"start_time":"2000-04-01 01:02:03"}`)),
			code:     http.StatusOK,
			expected: `{"e":0,"err_msg":"","window_len":15,"start_time":"2000-04-01 01:02:03","cabinets":[{"avg":0,"max":0,"min":0,"x":0,"y":0},{"avg":1.5,"max":3,"min":1,"x":1,"y":0}]}` + "\n",
		},
		{
			name:     "POST start_time now",
			method:   http.MethodPost,
			path:     capmc.SystemPowerDetails,
			body:     bytes.NewBuffer([]byte(`{"start_time":"now"}`)),
			code:     http.StatusOK,
			expected: `{"e":0,"err_msg":"","window_len":15,"start_time":"@NOW@","cabinets":[{"avg":0,"max":0,"min":0,"x":0,"y":0},{"avg":1.5,"max":3,"min":1,"x":1,"y":0}]}` + "\n",
		},
		{
			name:     "POST start_time invalid",
			method:   http.MethodPost,
			path:     capmc.SystemPowerDetails,
			body:     bytes.NewBuffer([]byte(`{"start_time":"2019-05-17 06:48:00 AND"}`)),
			code:     http.StatusBadRequest,
			expected: BadStartTimeJSON + "\n",
		},
		{
			name:     "POST window_len valid",
			method:   http.MethodPost,
			path:     capmc.SystemPowerDetails,
			body:     bytes.NewBuffer([]byte(`{"window_len":15}`)),
			code:     http.StatusOK,
			expected: `{"e":0,"err_msg":"","window_len":15,"start_time":"@NOW@","cabinets":[{"avg":0,"max":0,"min":0,"x":0,"y":0},{"avg":1.5,"max":3,"min":1,"x":1,"y":0}]}` + "\n",
		},
		{
			name:     "POST window_len below range",
			method:   http.MethodPost,
			path:     capmc.SystemPowerDetails,
			body:     bytes.NewBuffer([]byte(`{"window_len":1}`)),
			code:     http.StatusBadRequest,
			expected: WindowLenOutOfRangeJSON + "\n",
		},
		{
			name:     "POST window_len above range",
			method:   http.MethodPost,
			path:     capmc.SystemPowerDetails,
			body:     bytes.NewBuffer([]byte(`{"window_len":3601}`)),
			code:     http.StatusBadRequest,
			expected: WindowLenOutOfRangeJSON + "\n",
		},
		{
			name:     "POST window_len not a number",
			method:   http.MethodPost,
			path:     capmc.SystemPowerDetails,
			body:     bytes.NewBuffer([]byte(`{"window_len":"NaN"}`)),
			code:     http.StatusBadRequest,
			expected: BadWindowLenJSON + "\n",
		},
		{
			name:     "PUT",
			method:   http.MethodPut,
			path:     capmc.SystemPowerDetails,
			body:     nil,
			code:     http.StatusMethodNotAllowed,
			expected: PutNotAllowedJSON + "\n",
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

			body := rr.Body.String()

			if strings.Contains(tc.expected, "@NOW@") {
				now := time.Now()
				nowStr := now.Format(intervalTimeFormat)
				tc.expected = strings.Replace(tc.expected, "@NOW@", nowStr, 1)
				i := strings.Index(body, "start_time\":")
				i = i + len("start_time\":") + 1
				dt := body[i : i+len(intervalTimeFormat)]
				body = strings.Replace(body, dt, nowStr, 1)
			}

			if body != tc.expected {
				t.Errorf("handler returned unexpected body: want '%v' but got '%v'",
					tc.expected, body)
			}
		})
	}
}
