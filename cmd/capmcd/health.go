// Copyright 2020 Cray Inc. All Rights Reserved.
//
// Except as permitted by contract or express written permission of Cray Inc.,
// no part of this work or its content may be modified, used, reproduced or
// disclosed in any form. Modifications made without express permission of
// Cray Inc. may damage the system the software is installed within, may
// disqualify the user from receiving support from Cray Inc. under support or
// maintenance contracts, or require additional support services outside the
// scope of those contracts to repair the software or system.
//

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"stash.us.cray.com/HMS/hms-capmc/internal/tsdb"
)

// HealthResponse - used to report service health stats
type HealthResponse struct {
	Readiness      string `json:"readiness"`
	Vault          string `json:"vault"`
	HSMConnection  string `json:"hsm"`
	TMDBConnection string `json:"tmdb"`
}

// doHealth - returns useful information about the service to the user
func (d *CapmcD) doHealth(w http.ResponseWriter, r *http.Request) {
	// NOTE: this is provided as a debugging aid for administrators to
	//  find out what is going on with the system.  This should return
	//  information in a human-readable format that will help to
	//  determine the state of this service.

	// only allow 'GET' calls
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		sendJsonError(w, http.StatusMethodNotAllowed,
			fmt.Sprintf("(%s) Not Allowed", r.Method))
		return
	}

	// NOTE: this is a very rough first version.  It just does a quick check
	//  of the dependent services.  It is possible to add a lot more information
	//  here including:
	//  1) Regular polling of dependents with 'last known connection' time/date
	//  2) Log real usage stats and just report here
	//  3) Service uptime, usage stats
	//  4) Add config information to the health report???

	// collect health information
	var stats HealthResponse

	// d.ss / d.css - vault connection
	numDep := 0
	if d.ss == nil || d.ccs == nil {
		stats.Vault = "No connection established to vault"
	} else {
		// now test that the connection does something
		creds, cerr := d.ccs.GetAllCompCreds()
		if cerr != nil {
			stats.Vault = fmt.Sprintf("Error retrieving credentials:%s", cerr.Error())
		} else {
			stats.Vault = fmt.Sprintf("Vault connection established with %d credentials loaded", len(creds))
			numDep++
		}
	}

	// d.svc.hsmURL - is there another connection mechanism?
	var hsmR struct {
		Code    int
		Message string
	}
	err := d.GetFromHSM("/service/ready", "", &hsmR)
	if err != nil {
		log.Printf("Health hsm query error:%s", err.Error())
		stats.HSMConnection = "HSM queries result in error"
	} else if hsmR.Code != 0 {
		log.Printf("Health hsm query return: %d, %s", hsmR.Code, hsmR.Message)
		stats.HSMConnection = "HSM not ready"
	} else {
		stats.HSMConnection = "HSM Ready"
		numDep++
	}

	// d.tsdb - query telemetry database connection status
	if tsdb.DB == nil {
		stats.TMDBConnection = "Telemetry database connection not initialized"
	} else {
		perr := tsdb.DB.Ping()
		if perr != nil {
			stats.TMDBConnection = fmt.Sprintf("Telemetry database connection error:%s", perr.Error())
		} else {
			stats.TMDBConnection = "Telemetry database connection established"
			numDep++
		}
	}

	// Look at the overall readiness of the service.  If all dependencies are
	// good, call it 'Ready', if some are OK and others not, call it
	// 'Degraded', and if none are ok reply 'Not Ready'
	const (
		DependenciesAll  = 3
		DependenciesNone = 0
	)
	if numDep == DependenciesAll {
		stats.Readiness = "Ready"
	} else if numDep > DependenciesNone {
		stats.Readiness = "Service Degraded"
	} else {
		stats.Readiness = "Not Ready"
	}

	// write the output
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
	return
}

// doReadiness - used for k8s readiness check
func (d *CapmcD) doReadiness(w http.ResponseWriter, r *http.Request) {
	// NOTE: this is coded in accordance with kubernetes best practices
	//  for liveness/readiness checks.  This function should only be
	//  used to indicate if something is wrong with this service that
	//  prevents usage.  If this fails too many times, the instance
	//  will be killed and re-started.  Only fail this if restarting
	//  this service is likely to fix the problem.

	// only allow 'GET' calls
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		sendJsonError(w, http.StatusMethodNotAllowed,
			fmt.Sprintf("(%s) Not Allowed", r.Method))
		return
	}

	// NOTE: at this time there are no dependent connections that would
	//  benefit from a restart, so always return good.  There is no concept
	//  of degraded but still partly functional at this point.

	w.WriteHeader(http.StatusNoContent)
	return
}

// doLiveness - used for k8s liveness check
func (d *CapmcD) doLiveness(w http.ResponseWriter, r *http.Request) {
	// NOTE: this is coded in accordance with kubernetes best practices
	//  for liveness/readiness checks.  This function should only be
	//  used to indicate the server is still alive and processing requests.

	// only allow 'GET' calls
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		sendJsonError(w, http.StatusMethodNotAllowed,
			fmt.Sprintf("(%s) Not Allowed", r.Method))
		return
	}

	// return simple StatusOK response to indicate server is alive
	w.WriteHeader(http.StatusNoContent)
	return
}
