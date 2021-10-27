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
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/Cray-HPE/hms-capmc/internal/capmc"
	"github.com/Cray-HPE/hms-capmc/internal/logger"
	"github.com/Cray-HPE/hms-capmc/internal/tsdb"
	"github.com/Cray-HPE/hms-certs/pkg/hms_certs"
	"github.com/sirupsen/logrus"

	base "github.com/Cray-HPE/hms-base"
	compcreds "github.com/Cray-HPE/hms-compcredentials"
	sstorage "github.com/Cray-HPE/hms-securestorage"
)

const clientTimeout = time.Duration(180) * time.Second

// API defines the relationship between the REST URL and the backend server
// function
type API struct {
	pattern string
	handler http.HandlerFunc
}

// APIs contains groupings of REST URL to backend server functionality for the
// different versions of the API
type APIs []API

var serviceName string
var svc CapmcD
var hms_ca_uri string
var rfClientLock sync.RWMutex

// TODO Move these to new file router.go
var capmcAPIs = []APIs{
	// Default API
	{
		API{capmc.EmergencyPowerOff, svc.doEmergencyPowerOff},
		API{capmc.GroupOff, svc.doGroupOff},
		API{capmc.GroupOn, svc.doGroupOn},
		API{capmc.GroupReinit, svc.doGroupReinit},
		API{capmc.GroupStatus, svc.doGroupStatus},
		API{capmc.Health, svc.doHealth},
		API{capmc.Liveness, svc.doLiveness},
		API{capmc.NodeEnergyCounter, svc.doNodeEnergyCounter},
		API{capmc.NodeEnergyStats, svc.doNodeEnergyStats},
		API{capmc.NodeEnergy, svc.doNodeEnergy},
		API{capmc.NodeIDMap, svc.doNidMap},
		API{capmc.NodeOff, svc.doNodeOff},
		API{capmc.NodeOn, svc.doNodeOn},
		API{capmc.NodeReinit, svc.doNodeRestart},
		API{capmc.NodeRules, svc.doNodeRules},
		API{capmc.NodeStatus, svc.doNodeStatus},
		API{capmc.PowerCapCapabilities, svc.doPowerCapCapabilities},
		API{capmc.PowerCapGet, svc.doPowerCapGet},
		API{capmc.PowerCapSet, svc.doPowerCapSet},
		API{capmc.Readiness, svc.doReadiness},
		API{capmc.SystemParams, svc.doSystemParams},
		API{capmc.SystemPower, svc.doSystemPower},
		API{capmc.SystemPowerDetails, svc.doSystemPowerDetails},
		API{capmc.XnameOff, svc.doXnameOff},
		API{capmc.XnameOn, svc.doXnameOn},
		API{capmc.XnameReinit, svc.doXnameReinit},
		API{capmc.XnameStatus, svc.doXnameStatus},
		API{capmc.SSDDiagsGet, svc.NotImplemented},
		API{capmc.SSDEnableClear, svc.NotImplemented},
		API{capmc.SSDEnableGet, svc.NotImplemented},
		API{capmc.SSDEnableSet, svc.NotImplemented},
		API{capmc.SSDSGet, svc.NotImplemented},
		API{capmc.MCDRAMCapabilities, svc.NotImplemented},
		API{capmc.MCDRAMConfigClear, svc.NotImplemented},
		API{capmc.MCDRAMConfigGet, svc.NotImplemented},
		API{capmc.MCDRAMConfigSet, svc.NotImplemented},
		API{capmc.NUMACapabilities, svc.NotImplemented},
		API{capmc.NUMAConfigClear, svc.NotImplemented},
		API{capmc.NUMAConfigGet, svc.NotImplemented},
		API{capmc.NUMAConfigSet, svc.NotImplemented},
	},
	// V0 API
	{
		API{capmc.NodeEnergyCounterV0, svc.doNodeEnergyCounter},
		API{capmc.NodeEnergyStatsV0, svc.doNodeEnergyStats},
		API{capmc.NodeEnergyV0, svc.doNodeEnergy},
		API{capmc.NodeIDMapV0, svc.NotImplemented},
		API{capmc.NodeOffV0, svc.doNodeOff},
		API{capmc.NodeOnV0, svc.doNodeOn},
		API{capmc.NodeReinitV0, svc.doNodeRestart},
		API{capmc.NodeRulesV0, svc.doNodeRules},
		API{capmc.NodeStatusV0, svc.doNodeStatus},
		API{capmc.PowerCapCapabilitiesV0, svc.doPowerCapCapabilities},
		API{capmc.PowerCapGetV0, svc.doPowerCapGet},
		API{capmc.PowerCapSetV0, svc.doPowerCapSet},
		API{capmc.SystemParamsV0, svc.doSystemParams},
		API{capmc.SystemPowerV0, svc.doSystemPower},
		API{capmc.SystemPowerDetailsV0, svc.doSystemPowerDetails},
	},
	// V1 API
	{
		API{capmc.EmergencyPowerOffV1, svc.doEmergencyPowerOff},
		API{capmc.GroupOffV1, svc.doGroupOff},
		API{capmc.GroupOnV1, svc.doGroupOn},
		API{capmc.GroupReinitV1, svc.doGroupReinit},
		API{capmc.GroupStatusV1, svc.doGroupStatus},
		API{capmc.HealthV1, svc.doHealth},
		API{capmc.LivenessV1, svc.doLiveness},
		API{capmc.NodeEnergyCounterV1, svc.doNodeEnergyCounter},
		API{capmc.NodeEnergyStatsV1, svc.doNodeEnergyStats},
		API{capmc.NodeEnergyV1, svc.doNodeEnergy},
		API{capmc.NodeIDMapV1, svc.doNidMap},
		API{capmc.NodeOffV1, svc.doNodeOff},
		API{capmc.NodeOnV1, svc.doNodeOn},
		API{capmc.NodeReinitV1, svc.doNodeRestart},
		API{capmc.NodeRulesV1, svc.doNodeRules},
		API{capmc.NodeStatusV1, svc.doNodeStatus},
		API{capmc.PowerCapCapabilitiesV1, svc.doPowerCapCapabilities},
		API{capmc.PowerCapGetV1, svc.doPowerCapGet},
		API{capmc.PowerCapSetV1, svc.doPowerCapSet},
		API{capmc.ReadinessV1, svc.doReadiness},
		API{capmc.SystemParamsV1, svc.doSystemParams},
		API{capmc.SystemPowerV1, svc.doSystemPower},
		API{capmc.SystemPowerDetailsV1, svc.doSystemPowerDetails},
		API{capmc.XnameOffV1, svc.doXnameOff},
		API{capmc.XnameOnV1, svc.doXnameOn},
		API{capmc.XnameReinitV1, svc.doXnameReinit},
		API{capmc.XnameStatusV1, svc.doXnameStatus},
		API{capmc.SSDDiagsGetV1, svc.NotImplemented},
		API{capmc.SSDEnableClearV1, svc.NotImplemented},
		API{capmc.SSDEnableGetV1, svc.NotImplemented},
		API{capmc.SSDEnableSetV1, svc.NotImplemented},
		API{capmc.SSDSGetV1, svc.NotImplemented},
		API{capmc.MCDRAMCapabilitiesV1, svc.NotImplemented},
		API{capmc.MCDRAMConfigClearV1, svc.NotImplemented},
		API{capmc.MCDRAMConfigGetV1, svc.NotImplemented},
		API{capmc.MCDRAMConfigSetV1, svc.NotImplemented},
		API{capmc.NUMACapabilitiesV1, svc.NotImplemented},
		API{capmc.NUMAConfigClearV1, svc.NotImplemented},
		API{capmc.NUMAConfigGetV1, svc.NotImplemented},
		API{capmc.NUMAConfigSetV1, svc.NotImplemented},
	},
}

// ResponseWriter contains the http ResponseWriter function to allow for
// wrapping of HTTP calls, as a server, for logging purposes.
type ResponseWriter struct {
	status int
	length int
	data   string
	http.ResponseWriter
}

// WriteHeader wraps HTTP calls, as a server, to enable logging of requests and
// responses.
func (w *ResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

// Write wrapps HTTP calls, as a server, to enable logging of requests and
// responses.
func (w *ResponseWriter) Write(b []byte) (int, error) {
	if svc.debug && svc.debugLevel > 1 {
		// This is simpler than using httptest.NewRecorder in the
		// middleware logger for a response to and incoming request.
		w.data = string(b)
	}

	n, err := w.ResponseWriter.Write(b)
	w.length += n

	return n, err
}

var suppressLogPaths map[string]bool

// Add a path to the list suppress for logging
func suppressLoggingForPath(p string) {
	// make the map if it doesn't already exist
	if suppressLogPaths == nil {
		suppressLogPaths = make(map[string]bool)
	}

	// add the entry
	suppressLogPaths[p] = true
}

// find if the input path has logging suppressed
func isPathSuppressed(p string) bool {
	// if the map wasn't created, false
	if suppressLogPaths == nil {
		return false
	}

	// query the map
	_, retVal := suppressLogPaths[p]
	return retVal
}

// logRequest is a middleware wrapper that handles logging server HTTP
// inbound requests and outbound responses.
func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			sendFmts = []string{
				"Debug: --> Response %s:\n%s\n",
				"Debug: --> Response %s: %q\n", // "safe"
			}
			recvFmts = []string{
				"Debug: <-- Request %s:\n%s\n",
				"Debug: <-- Request %s: %q\n", // "safe"
			}
		)

		rw := &ResponseWriter{
			status:         http.StatusOK,
			ResponseWriter: w,
		}

		// see if this is a path to be suppressed
		suppressLog := isPathSuppressed(r.URL.Path)

		start := time.Now()
		if !suppressLog {
			log.Printf("Info: <-- %s HTTP %s %s\n",
				r.RemoteAddr, r.Method, r.URL)
		}

		if svc.debug && svc.debugLevel > 1 {
			sendFmt := sendFmts[svc.debugLevel%2]
			dump, err := httputil.DumpRequest(r, true)
			if err != nil {
				log.Printf("Debug: failed to dump request: %s",
					err)
				http.Error(w, fmt.Sprint(err),
					http.StatusInternalServerError)
				return
			}

			s := bytes.SplitAfterN(dump, []byte("\r\n\r\n"), 2)
			if svc.debugLevel > 3 {
				log.Printf(sendFmt, "Header", s[0])
			}
			log.Printf(sendFmt, "Body", dump)
		}

		handler.ServeHTTP(rw, r)

		if !suppressLog {
			log.Printf("Info: --> %s HTTP %d %s %s %s (%s)",
				r.RemoteAddr, rw.status, http.StatusText(rw.status),
				r.Method, r.URL.String(), time.Since(start))
		}
		if svc.debug && svc.debugLevel > 1 {
			// Capturing the reponse headers requires more
			// work using httptest.NewRecorder. Skip for now.
			recvFmt := recvFmts[svc.debugLevel%2]
			log.Printf(recvFmt, "Body", rw.data)
		}
	})
}

// Transport contains the http RoundTripper function to allow for wrapping of
// HTTP calls, as a client, for logging purposes.
type Transport struct {
	Transport http.RoundTripper
}

// RoundTrip wraps HTTP calls, as a client, to enable logging of requests and
// responses.
func (t *Transport) RoundTrip(r *http.Request) (*http.Response, error) {
	var (
		sendFmts = []string{
			"Debug: --> Request %s:\n%s\n",
			"Debug: --> Request %s: %q\n", // "safe"
		}
		recvFmts = []string{
			"Debug: <-- Response %s:\n%s\n",
			"Debug: <-- Response %s: %q\n", // "safe"
		}
	)

	log.Printf("Info: --> HTTP %s %s\n", r.Method, r.URL)
	if svc.debug && svc.debugLevel > 1 {
		sendFmt := sendFmts[svc.debugLevel%2]
		dump, err := httputil.DumpRequestOut(r, true)
		if err != nil {
			log.Printf("Debug: failed to dump request: %s", err)
		} else {
			s := bytes.SplitAfterN(dump, []byte("\r\n\r\n"), 2)
			if svc.debugLevel > 3 {
				log.Printf(sendFmt, "Header", s[0])
			}
			log.Printf(sendFmt, "Body", s[1])
		}
	}

	start := time.Now()
	resp, err := t.Transport.RoundTrip(r)
	if err != nil {
		return resp, err
	}

	log.Printf("Info: <-- HTTP %s %s %s (%s)\n",
		resp.Status, resp.Request.Method, resp.Request.URL,
		time.Since(start))
	if svc.debug && svc.debugLevel > 1 {
		recvFmt := recvFmts[svc.debugLevel%2]
		dump, err := httputil.DumpResponse(resp, true)
		s := bytes.SplitAfterN(dump, []byte("\r\n\r\n"), 2)
		if err != nil {
			log.Printf("Debug: failed to dump response: %s", err)
		} else {
			if svc.debugLevel > 3 {
				log.Printf(recvFmt, "Header", s[0])
			}
			log.Printf(recvFmt, "Body", s[1])
		}
	}

	return resp, err
}

type clientFlags int

const (
	clientInsecure clientFlags = 1 << iota
	// Other flags could be added here...
)

// Construct an http client object.
func makeClient(flags clientFlags, timeout time.Duration) (*hms_certs.HTTPClientPair, error) {
	var pair *hms_certs.HTTPClientPair
	var err error

	uri := hms_ca_uri
	if (flags & clientInsecure) != 0 {
		uri = ""
	}
	toSec := int(timeout.Seconds())
	pair, err = hms_certs.CreateHTTPClientPair(uri, toSec)
	if err != nil {
		log.Printf("ERROR: Can't set up cert-secured HTTP client: %v", err)
		return nil, err
	}

	return pair, nil
}

// This will initialize the global CapmcD struct with default values upon startup
func init() {
}

// Called when CA cert bundle is rolled.

func caCB(caData string) {
	log.Printf("INFO: Updating CA bundle for Redfish HTTP transports.")

	//Wait for all reader locks to release, prevent new reader locks.  Once
	//acquired, all RF calls are blocked.

	rfClientLock.Lock()
	log.Printf("INFO: All RF threads are paused.")

	//Update the the transports.

	if captureTag != "" {
		svc.rfClient = captureClient(0, clientTimeout)
	} else {
		cl, err := makeClient(0, clientTimeout)
		if err != nil {
			log.Printf("ERROR: can't create Redfish HTTP client after CA roll: %v", err)
			log.Printf("    Using previous HTTP client (CA bundle may not work.)")
		} else {
			svc.rfClient = cl
		}
	}
	log.Printf("Redfish transport clients updated with new CA bundle.")
	rfClientLock.Unlock()
}

// Set up Redfish and non-redfish HTTP clients.

func setupRedfishHTTPClients(captureTag string) error {
	var err error

	if hms_ca_uri != "" {
		log.Printf("INFO: Creating Redfish HTTP transport using CA bundle '%s'",
			hms_ca_uri)
	} else {
		log.Printf("INFO: Creating non-validated Redfish HTTP transport (no CA bundle)")
	}

	if captureTag != "" {
		svc.rfClient = captureClient(0, clientTimeout)
	} else {
		svc.rfClient, err = makeClient(0, clientTimeout)
	}
	if err != nil {
		log.Printf("ERROR setting up Redfish HTTP transport: '%v'", err)
		return err
	}
	if hms_ca_uri != "" {
		err := hms_certs.CAUpdateRegister(hms_ca_uri, caCB)
		if err != nil {
			log.Printf("ERROR: can't register CA bundle watcher: %v", err)
			log.Printf("   This means CA bundle updates will not update Redfish HTTP transports.")
		} else {
			log.Printf("INFO: CA bundle watcher registered for '%s'.", hms_ca_uri)
		}
	} else {
		log.Printf("INFO: CA bundle URI is empty, no CA bundle watcher registered.")
	}
	return nil
}

func main() {
	var (
		configFile string
		hsm        string
		err        error
	)

	// TODO Add debug levels at some point
	flag.BoolVar(&svc.debug, "debug", false, "enable debug messages")
	flag.IntVar(&svc.debugLevel, "debug-level", 0, "increase debug verbosity")

	// Simulation only
	//  NOTE: if this is set to 'true' then all of the calls to BMC's will only
	//   be logged and not executed.  This is provided as a mechanism for testing
	//   on real hardware without actually turning systems on and off.
	flag.BoolVar(&svc.simulationOnly, "simulateOnly", false, "Only log calls to BMC's instead of executing them")

	// LOGGING // SETUP LOGRUS GLOBAL!
	// using name logrus, as to not override the log stuff.
	logger.SetupLogging()

	// TODO Add support for specifying http/https with the latter as default
	//      It might make sense to use the URI format here too.
	flag.StringVar(&svc.httpListen, "http-listen", "0.0.0.0:27777", "HTTP server IP + port binding")
	flag.StringVar(&hsm, "hsm", "http://localhost:27779",
		"Hardware State Manager location as URI, e.g. [scheme]://[host[:port]]")
	flag.StringVar(&hms_ca_uri, "ca_uri", "",
		"Certificate Authority CA bundle URI")

	// The "default" is installed with the service. The intent is
	// ConfigPath/ConfigFile is a customized config and
	// ConfigPath/default/ConfigFile contains the installed (and internal)
	// default values.
	// TODO Add a development ConfigPath allowing for non-install,
	// non-container development without needed to specify -config <file>.
	flag.StringVar(&configFile, "config",
		filepath.Join(ConfigPath, "default", ConfigFile),
		"Configuration file")

	var captureTag string
	flag.StringVar(&captureTag, "capture", "",
		"Capture client traffic for test case using TAG")
	flag.Parse()

	log.SetFlags(log.Lshortfile | log.LstdFlags)

	serviceName, err = base.GetServiceInstanceName()
	if err != nil {
		serviceName = "CAPMC"
		log.Printf("WARNING: can't get service/instance name, using: '%s'",
			serviceName)
	}
	log.Printf("Service name/instance: '%s'", serviceName)

	svc.config = loadConfig(configFile)
	conf := svc.config.CapmcConf

	log.Printf("Configuration loaded:\n")
	log.Printf("\tMax workers: %d\n", conf.ActionMaxWorkers)
	log.Printf("\tOn unsupported action: %s\n", conf.OnUnsupportedAction)
	log.Printf("\tReinit seq: %v\n", conf.ReinitActionSeq)
	log.Printf("\tWait for off retries: %d\n", conf.WaitForOffRetries)
	log.Printf("\tWait for off sleep: %d\n", conf.WaitForOffSleep)

	svc.ActionMaxWorkers = conf.ActionMaxWorkers
	svc.OnUnsupportedAction = conf.OnUnsupportedAction
	svc.ReinitActionSeq = conf.ReinitActionSeq

	// log the hostname of this instance - mostly useful for pod name in
	// multi-replica k8s envinronment
	hostname, hostErr := os.Hostname()
	if hostErr != nil {
		log.Printf("Error getting hostname:%s", hostErr.Error())
	} else {
		log.Printf("Starting on host: %s", hostname)
	}

	// log if this is in simulate only mode
	if svc.simulationOnly {
		log.Printf("WARNING: Started in SIMULATION ONLY mode - no commands will be sent to BMC hardware")
	}

	// CapmcD is both HTTP server and client
	if svc.hsmURL, err = url.Parse(hsm); err != nil {
		log.Fatalf("Invalid HSM URI specified: %s", err)
	}

	// Set up the hsm information before connecting to any external
	// resources since we may bail if we don't find what we want
	// Check for non-empty URL (URI) scheme
	if !svc.hsmURL.IsAbs() {
		log.Fatal("WARNING: HSM URL not absolute\n")
	}
	switch svc.hsmURL.Scheme {
	case "http", "https":
		log.Printf("Info: hardware state manager (HSM) --> %s\n",
			svc.hsmURL.String())
		// Stash the HSM Base version in the URL.Path (default)
		// XXX Should the HSM API (default) version be configurable?
		switch svc.hsmURL.Path {
		case "":
			svc.hsmURL.Path = "/hsm/v1"
		case "/hsm/v1":
			// do nothing
		default:
			if !strings.HasSuffix(svc.hsmURL.Path, "/hsm/v1") {
				svc.hsmURL.Path += "/hsm/v1"
			}
		}
	default:
		log.Fatalf("Unexpected HSM URL scheme: %s", svc.hsmURL.Scheme)
	}

	//CA/cert stuff

	vurl := os.Getenv("CAPMC_VAULT_CA_URL")
	if vurl != "" {
		log.Printf("Replacing default Vault CA URL with: '%s'", vurl)
		hms_certs.ConfigParams.VaultCAUrl = vurl
	}
	vurl = os.Getenv("CAPMC_VAULT_PKI_URL")
	if vurl != "" {
		log.Printf("Replacing default Vault PKI URL with: '%s'", vurl)
		hms_certs.ConfigParams.VaultPKIUrl = vurl
	}
	if hms_ca_uri == "" {
		vurl = os.Getenv("CAPMC_CA_URI")
		if vurl != "" {
			log.Printf("Using CA URI: '%s'", vurl)
			hms_ca_uri = vurl
		}
	}
	vurl = os.Getenv("CAPMC_LOG_INSECURE_FAILOVER")
	if vurl != "" {
		yn, _ := strconv.ParseBool(vurl)
		if yn == false {
			hms_certs.ConfigParams.LogInsecureFailover = false
		}
	}
	hms_certs.InitInstance(nil, serviceName)

	log.Printf("CAPMC serivce starting (debug=%v)\n", svc.debug)

	// set up a channel to wait for the os to tell us to stop
	// NOTE - must be set up before initializing anything that needs
	//  to be cleaned up.  This will trap any signals and wait to
	//  process them until the channel is read.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	//initialize the service-reservation pkg
	svc.reservation.Init(svc.hsmURL.Scheme+"://"+svc.hsmURL.Host, "", 3, nil)
	svc.reservationsEnabled = true

	//Configure TSDB connection; MAY USE DUMMY
	// Spin this in another thread until the connection is successful
	go func() {
		const (
			initBackoff time.Duration = 5
			maxBackoff  time.Duration = 60
		)

		backoff := initBackoff
		for {
			if err := tsdb.ConfigureDataImplementation(tsdb.UNKNOWN); err != nil {
				logrus.Warning("DB Failed to connect")

				time.Sleep(backoff * time.Second)
			} else {
				logrus.Info("DB Connection established")
				break
			}
			if backoff < maxBackoff {
				backoff += backoff
			}
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}
	}()

	// Create deferred function to close the database connection
	// This is important if any other 'startup' sections below
	// fail before entering normal serving mode.
	defer func() {
		if tsdb.ImplementedConnection == tsdb.POSTGRES {
			err = tsdb.DB.Close()
			if err != nil {
				logrus.WithField("error", err).Error("Could not cleanly close connection to DB!")
			} else {
				logrus.Info("Connection to DB closed.")
			}
		}
	}()

	// Spin a thread for connecting to Vault
	go func() {
		const (
			initBackoff time.Duration = 5
			maxBackoff  time.Duration = 60
		)
		var err error

		backoff := initBackoff
		for {
			log.Printf("Info: Connecting to secure store (Vault)...")
			// Start a connection to Vault
			if svc.ss, err = sstorage.NewVaultAdapter(""); err != nil {
				log.Printf("Info: Secure Store connection failed - %s", err)
				time.Sleep(backoff * time.Second)
			} else {
				log.Printf("Info: Connection to secure store (Vault) succeeded")
				svc.ccs = compcreds.NewCompCredStore("secret/hms-creds", svc.ss)
				break
			}
			if backoff < maxBackoff {
				backoff += backoff
			}
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}
	}()

	//Set up non-secure HTTP client for HSM, etc.

	log.Printf("INFO: Creating inter-service HTTP transport (not TLS validated).")
	if captureTag != "" {
		captureInit(captureTag)
		svc.smClient = captureClient(clientInsecure, clientTimeout)
	} else {
		svc.smClient, _ = makeClient(clientInsecure, clientTimeout)
	}

	//Set up secure HTTP clients for Redfish.  Keep trying in the background until success.

	go func() {
		const (
			initBackoff time.Duration = 2
			maxBackoff  time.Duration = 30
		)

		backoff := initBackoff
		for {
			log.Printf("Info: Creating Redfish secure HTTP client transports...")
			err := setupRedfishHTTPClients(captureTag)
			if err == nil {
				log.Printf("Info: Success creating Redfish HTTP clients.")
				break
			} else {
				log.Printf("Redfish HTTP client creation error: %v", err)
			}
			if backoff < maxBackoff {
				backoff += backoff
			}
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}
	}()

	// These are registrations for the CAPMC API calls.
	for _, vers := range capmcAPIs {
		for _, api := range vers {
			if captureTag == "" {
				http.HandleFunc(api.pattern, api.handler)
			} else {
				captureHandleFunc(api.pattern, api.handler)
			}
		}
	}

	// Do not log the calls for liveness/readiness
	suppressLoggingForPath(capmc.Liveness)
	suppressLoggingForPath(capmc.LivenessV1)
	suppressLoggingForPath(capmc.Readiness)
	suppressLoggingForPath(capmc.ReadinessV1)

	// Spin up our global worker goroutine pool.
	svc.WPool = base.NewWorkerPool(svc.ActionMaxWorkers, svc.ActionMaxWorkers*10)
	svc.WPool.Run()

	// The following thread talks about limiting the max post body size...
	// https://stackoverflow.com/questions/28282370/is-it-advisable-to-further-limit-the-size-of-forms-when-using-golang

	// spin the server in a separate thread so main can wait on an os
	// signal to cleanly shut down
	httpSrv := http.Server{
		Addr:    svc.httpListen,
		Handler: logRequest(http.DefaultServeMux),
	}
	go func() {
		// NOTE: do not use log.Fatal as that will immediately exit
		// the program and short-circuit the shutdown logic below
		log.Printf("Info: Server %s\n", httpSrv.ListenAndServe())
	}()
	log.Printf("Info: CAPMC API listening on: %v\n", svc.httpListen)

	//////////////////
	// Clean shutdown section
	//////////////////

	// wait here for a signal from the os that we are shutting down
	sig := <-sigs
	log.Printf("Info: Detected signal to close service: %s", sig)

	// The service is being killed, so release all active locks in hsm
	// NOTE: this happens when k8s kills a pod
	svc.removeAllActiveReservations()

	// stop the server from taking requests
	// NOTE: this waits for active connections to finish
	log.Printf("Info: Server shutting down")
	httpSrv.Shutdown(context.Background())

	// terminate worker pool
	//  NOTE: have to do this the hard way since there isn't a
	//  clean shutdown implemented on the WorkerPool
	log.Printf("Info: Finishing current jobs in the pool")
	waitTime := time.Second * 10 // hard code to 10 sec wait time for now
	poolTimeout := time.Now().Add(waitTime)
	for {
		// if we have hit the timeout or all the jobs are out of the queue
		// we can bail
		if len(svc.WPool.JobQueue) == 0 || time.Now().After(poolTimeout) {
			break
		}
		// wait another second then check again
		time.Sleep(time.Second)
	}
	// with no jobs left in the queue, we can stop the service
	// this waits until currently running jobs are complete before exiting
	svc.WPool.Stop()

	// NOTE: This is where we should terminate our connection to the
	//  vault, but it looks like there is no way to do so at this time.

	// NOTE: The db connection is being shut down via a deferred function call

	// shut down hsm client connections
	log.Printf("Info: Closing idle client connections...")
	svc.smClient.CloseIdleConnections()
	svc.rfClient.CloseIdleConnections()

	log.Printf("Info: Service Exiting.")
}

// TODO https://blog.golang.org/error-handling-and-go

// Send error or empty OK response. This format matches existing CAPMC API.
// Error code is the http status response.
func sendJsonError(w http.ResponseWriter, ecode int, message string) {
	// If HTTP call is success, put zero in returned json error field.
	// This is what Cascade capmc does today.
	httpCode := ecode
	if ecode >= 200 && ecode <= 299 {
		ecode = 0
	}

	data := capmc.ErrResponse{
		E:      ecode,
		ErrMsg: message,
	}

	SendResponseJSON(w, httpCode, data)
}

// SendResponseJSON sends data marshalled as a JSON body and sets the HTTP
// status code to sc.
func SendResponseJSON(w http.ResponseWriter, sc int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(sc)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("Error: encoding/sending JSON response: %s\n", err)
		return
	}
}

// NotImplemented is used as a placeholder API entry point.
func (d *CapmcD) NotImplemented(w http.ResponseWriter, r *http.Request) {
	var body = capmc.ErrResponse{
		E:      http.StatusNotImplemented,
		ErrMsg: fmt.Sprintf("%s API Unavailable/Not Implemented", r.URL.Path),
	}

	SendResponseJSON(w, http.StatusNotImplemented, body)

	return
}
