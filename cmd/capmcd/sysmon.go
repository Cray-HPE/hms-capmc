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
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"strconv"
	"time"

	"github.com/Cray-HPE/hms-capmc/internal/capmc"
	"github.com/Cray-HPE/hms-capmc/internal/tsdb"
)

const (
	minInterval = 2
	maxInterval = 3600
)

func inRange(i, min, max int) bool {
	return (i >= min) && (i <= max)
}

// doSystemParameters handles a system power parameters
func (d *CapmcD) doSystemParams(w http.ResponseWriter, r *http.Request) {

	// NOTE Cascade CAPMC does not support a GET but there is no reason
	//      why it couldn't as it accepts an *empty* JSON object.
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		w.Header().Set("Allow", "GET, POST")
		sendJsonError(w, http.StatusMethodNotAllowed,
			fmt.Sprintf("(%s) Not Allowed", r.Method))
		return
	}

	switch r.Method {
	case http.MethodGet:
		err := r.ParseForm()
		if err != nil {
			sendJsonError(w, http.StatusBadRequest,
				fmt.Sprintf("Bad Request: %s", err))
			return
		}

		if len(r.Form) > 0 {
			log.Printf("Unexpected: %s %s with query %s",
				r.Method, r.URL.Path, r.URL.RawQuery)
		}
	case http.MethodPost:
		// NOTE Finding no good examples or documentation on what
		// should happen when a POST has query parameters we'll
		// follow the maxim "be liberal in what you accept, and
		// conservative in what you send".  Log it and move on.
		m, err := url.ParseQuery(r.URL.RawQuery)
		if len(m) > 0 {
			log.Printf("Warning: Unexpected %s %s with query %s",
				r.Method, r.URL.Path, r.URL.RawQuery)
		}

		// There are no supported body parameters at this time but
		// following the above "be liberal in what you accept, and
		// conservative in what you send", we'll only error out when
		// the request is malformed. Log it and move on.
		var v interface{}
		err = json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			log.Printf("Warning: Unexpected %s %s with body",
				r.Method, r.URL.Path)
		}
	}

	log.Printf("Info: Get System Power Parameters")

	var data = capmc.GetSystemParametersResponse{
		PowerCapTarget: d.config.SystemParams.PowerCapTarget,
		PowerThreshold: d.config.SystemParams.PowerThreshold,
		StaticPower:    d.config.SystemParams.StaticPower,
		RampLimited:    d.config.SystemParams.RampLimited,
		RampLimit:      d.config.SystemParams.RampLimit,
		PowerBandMin:   d.config.SystemParams.PowerBandMin,
		PowerBandMax:   d.config.SystemParams.PowerBandMax,
		ErrResponse:    capmc.ErrResponse{E: 0, ErrMsg: ""},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// doSystemPower handles a system power request
func (d *CapmcD) doSystemPower(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		w.Header().Set("Allow", "GET, POST")
		sendJsonError(w, http.StatusMethodNotAllowed,
			fmt.Sprintf("(%s) Not Allowed", r.Method))
		return
	}

	var args capmc.TimeWindowRequest

	switch r.Method {
	case http.MethodGet:
		err := r.ParseForm()
		if err != nil {
			sendJsonError(w, http.StatusBadRequest,
				fmt.Sprintf("Bad Request: %s", err))
			return
		}

		args.StartTime = r.FormValue("start_time")
		s := r.FormValue("window_len")
		if s != "" {
			i, err := strconv.Atoi(s)
			if err != nil {
				sendJsonError(w, http.StatusBadRequest, BadWindowLen)
				return
			}

			args.WindowLen = new(int)
			*args.WindowLen = i
		}
	case http.MethodPost:
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&args)
		if err != nil {
			var msg string

			switch err {
			case io.EOF:
				msg = "no request"
			default:
				e, ok := err.(*json.UnmarshalTypeError)
				if ok && e.Field == "window_len" {
					msg = BadWindowLen
				} else {
					msg = BadJSON
				}
			}

			sendJsonError(w, http.StatusBadRequest, msg)
			return
		}
	}

	// Default window_len is 15 (optional parameter)
	// Because of the hysteresis, this value has been expanded from 10
	// seconds to 15 seconds. The bleeding edge is not stable at Now(); it
	// will stabilize within 15 seconds. We have a guarantee that this
	// period size will make sure we dont wind up with empty periods, due
	// to hysteresis.
	var sampleInterval = DefaultSampleWindow

	// Validate input request parameters
	if args.WindowLen != nil {
		if !inRange(*args.WindowLen, minInterval, maxInterval) {
			// Error EINVAL
			sendJsonError(w, http.StatusBadRequest,
				WindowLenOutOfRange)
			return
		}

		sampleInterval = time.Duration(*args.WindowLen) * time.Second
		// while we can accept 2 - 14 as the window length; we are
		// going to actually enforce 15 (DefaultSampleWindow as the
		// minInterval.
		// Put this code here, b/c we have to honor the old Cascade API
		if sampleInterval < DefaultSampleWindow {
			sampleInterval = DefaultSampleWindow
		}
	}

	var endTime, startTime time.Time
	var err error

	switch args.StartTime {
	case "", "now":

		endTime = time.Now().UTC()
		startTime = endTime.Add(-sampleInterval)

	default:
		// validate input string
		startTime, err = time.Parse(intervalTimeFormat, args.StartTime)

		if err != nil {
			log.Printf("Notice: failed start time parse %s", err)
			sendJsonError(w, http.StatusBadRequest, BadStartTime)
			return
		}

		endTime = startTime.Add(sampleInterval)
	}
	log.Printf("Get System Power from %s to %s (%s window)",
		startTime, endTime, sampleInterval.String())
	startTime, endTime, _ = tsdb.AdjustTimestampsForHystersisAndMinimumWindow(startTime, endTime, time.Now().UTC(), DefaultHysteresis, DefaultSampleWindow)
	tbr := tsdb.TimeBoundRequest{startTime, endTime}
	// Hit PostgreSQL TimeScale DB here
	sysPow, err := tsdb.TSDBContext.GetSystemPower(tbr)
	if err != nil {
		// Use the better message if one is sent
		var terr *tsdb.TimeBoundRequestError
		if errors.As(err, &terr) {
			log.Printf("Info: %s", terr.Error())
			SendResponseJSON(w, http.StatusOK, capmc.ErrResponse{E: http.StatusBadRequest, ErrMsg: terr.Error()})
			return
		}

		// process generic error
		log.Printf("Error: Querying TSDB/PMDB: %s", err)
		sendJsonError(w, http.StatusInternalServerError,
			"failure retrieving system power data from power management database")
		return
	}
	var data = capmc.GetSystemPowerResponse{
		ErrResponse: capmc.ErrResponse{E: 0, ErrMsg: ""},
		WindowLen:   int(sampleInterval / time.Second),
		StartTime:   startTime.Format(intervalTimeFormat),
		Avg:         sysPow.Avg,
		Max:         sysPow.Max,
		Min:         sysPow.Min,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// doSystemPowerDetails handles a system power details request
func (d *CapmcD) doSystemPowerDetails(w http.ResponseWriter, r *http.Request) {

	// NOTE Cascade doesn't support GET but there is no reason it couldn't
	//      as the arguments are optional to the API.
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		w.Header().Set("Allow", "GET, POST")
		sendJsonError(w, http.StatusMethodNotAllowed,
			fmt.Sprintf("(%s) Not Allowed", r.Method))
		return
	}

	var args capmc.TimeWindowRequest

	switch r.Method {
	case http.MethodGet:
		err := r.ParseForm()
		if err != nil {
			sendJsonError(w, http.StatusBadRequest,
				fmt.Sprintf("Bad Request: %s", err))
			return
		}

		args.StartTime = r.FormValue("start_time")
		s := r.FormValue("window_len")
		if s != "" {
			i, err := strconv.Atoi(s)
			if err != nil {
				sendJsonError(w, http.StatusBadRequest, BadWindowLen)
				return
			}

			args.WindowLen = new(int)
			*args.WindowLen = i
		}
	case http.MethodPost:
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&args)
		if err != nil {
			var msg string

			switch err {
			case io.EOF:
				msg = "no request"
			default:
				e, ok := err.(*json.UnmarshalTypeError)
				if ok && e.Field == "window_len" {
					msg = BadWindowLen
				} else {
					msg = BadJSON
				}
			}

			sendJsonError(w, http.StatusBadRequest, msg)
			return
		}
	}

	// Default window_len is 15 (optional parameter)
	// Because of the hysteresis, this value has been expanded from 10
	// seconds to 15 seconds. The bleeding edge is not stable at Now(); it
	// will stabilize within 15 seconds. We have a guarantee that this
	// period size will make sure we dont wind up with empty periods, due
	// to hysteresis.
	var sampleInterval = DefaultSampleWindow

	// Validate input request window_len
	if args.WindowLen != nil {
		if !inRange(*args.WindowLen, minInterval, maxInterval) {
			// Error EINVAL
			sendJsonError(w, http.StatusBadRequest,
				WindowLenOutOfRange)
			return
		}

		sampleInterval = time.Duration(*args.WindowLen) * time.Second
		// while we can accept 2 - 14 as the window length; we are
		// going to actually enforce 15 (DefaultSampleWindow as the
		// minInterval
		// Put this code here, b/c we have to honor the old cascade api
		if sampleInterval < DefaultSampleWindow {
			sampleInterval = DefaultSampleWindow
		}
	}

	var endTime, startTime time.Time
	var err error

	switch args.StartTime {
	case "", "now":

		endTime = time.Now().UTC()
		startTime = endTime.Add(-sampleInterval)

	default:
		// validate input string
		startTime, err = time.Parse(intervalTimeFormat, args.StartTime)

		if err != nil {
			log.Printf("Notice: failed start time parse %s", err)
			sendJsonError(w, http.StatusBadRequest, BadStartTime)
			return
		}

		endTime = startTime.Add(sampleInterval)
	}

	log.Printf("Get System Power Details from %s to %s (%s window)",
		startTime, endTime, sampleInterval.String())
	startTime, endTime, _ = tsdb.AdjustTimestampsForHystersisAndMinimumWindow(startTime, endTime, time.Now().UTC(), DefaultHysteresis, DefaultSampleWindow)
	tbr := tsdb.TimeBoundRequest{startTime, endTime}

	// Hit PostgreSQL TimeScale DB here

	cabPower, err := tsdb.TSDBContext.GetSystemPowerDetails(tbr)
	if err != nil {
		// Use the better message if one is sent
		var terr *tsdb.TimeBoundRequestError
		if errors.As(err, &terr) {
			log.Printf("Info: %s", terr.Error())
			SendResponseJSON(w, http.StatusOK, capmc.ErrResponse{E: http.StatusBadRequest, ErrMsg: terr.Error()})
			return
		}

		// Process the generic error
		log.Printf("Error: Querying TSDB/PMDB: %s", err)
		sendJsonError(w, http.StatusInternalServerError,
			"failure retrieving cabinet-level system power data from power management database")
		return
	}

	// Note the published Cascade API has a cabinet x, y (col, row).
	// We'll stick the cabinet number in the x (col) and leave
	// y (row) as 0. Leaving this as a component id (xname)
	// doesn't convey any additional information so convert the
	// cabinet number.
	var cpi []*capmc.CabinetPowerInfo

	for _, v := range cabPower.Cabinets {
		tmp := capmc.CabinetPowerInfo{}
		tmp.Avg = *v.Avg
		tmp.Min = *v.Min
		tmp.Max = *v.Max
		tmp.Y = 0
		tmp.X, _ = strconv.Atoi(v.ComponentID[1:])
		cpi = append(cpi, &tmp)
	}

	var data = capmc.GetSystemPowerDetailsResponse{
		ErrResponse: capmc.ErrResponse{E: 0, ErrMsg: ""},
		WindowLen:   int(sampleInterval / time.Second),
		StartTime:   startTime.Format(intervalTimeFormat),
		Cabinets:    cpi,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
