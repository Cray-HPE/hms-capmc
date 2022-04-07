/*
 * MIT License
 *
 * (C) Copyright [2019-2022] Hewlett Packard Enterprise Development LP
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

// This file is the definitions file for the capmcd application
// All definitions and constants should be moved into this module
//

package main

import (
	"net/url"
	"time"

	base "github.com/Cray-HPE/hms-base"
	"github.com/Cray-HPE/hms-capmc/internal/tsdb"
	"github.com/Cray-HPE/hms-certs/pkg/hms_certs"
	compcreds "github.com/Cray-HPE/hms-compcredentials"
	sstorage "github.com/Cray-HPE/hms-securestorage"
	reservation "github.com/Cray-HPE/hms-smd/pkg/service-reservations"
)

// Internal BMC command actions
const (
	bmcCmdNMI               = "NMI" // TODO future enhancement
	bmcCmdPowerForceOff     = "ForceOff"
	bmcCmdPowerForceOn      = "ForceOn"
	bmcCmdPowerForceRestart = "ForceRestart"
	bmcCmdPowerOff          = "Off"
	bmcCmdPowerOn           = "On"
	bmcCmdPowerRestart      = "Restart"
	bmcCmdPowerStatus       = "Status"
	bmcCmdGetPowerCap       = "GetPowerCap"
	bmcCmdSetPowerCap       = "SetPowerCap"
)

// Configuration values for OnUnsupportedAction
const (
	actionSimulate = "simulate"
	actionError    = "error"
	actionIgnore   = "ignore"
)

// Cascade error strings returned by the following xtremoted helpers:
// * xtremoted_get_system_power
// * xtremoted_get_system_power_details
// * xtremoted_get_node_energy
// * xtremoted_get_node_energy_counter
// * xtremoted_get_node_energy_stats
const (
	BadAPID                       = "BAD_APID"
	BadEndTime                    = "BAD_END_TIME"
	BadJobID                      = "BAD_JOB_ID"
	BadJSON                       = "BAD_JSON"
	BadNIDSFormat                 = "BAD_NIDS_FORMAT"
	BadOptions                    = "BAD_OPTIONS"
	BadStartTime                  = "BAD_START_TIME"
	BadWindowLen                  = "BAD_WINDOW_LEN"
	InvalidArguments              = "INVALID_ARGUMENTS"
	ArgumentSupportNotImplemented = "ARGUMENT_SUPPORT_NOT_IMPLEMENTED"
	MissingNidsApidJobid          = "MISSING_NIDS_APID_JOBID"
	NoData                        = "NO_DATA"
	NoResults                     = "NO_RESULTS"
	PSQLFailure                   = "PSQL_FAILURE"
	WindowLenOutOfRange           = "WINDOW_LEN_OUT_OF_RANGE"
)

// The Cascade API uses Python 'standard' time formating
// YYYY-MM-DD HH:MM:SS.MS which unfortunately isn't a default Go time format.
const intervalTimeFormat = "2006-01-02 15:04:05"

// DefaultSampleWindow is the default sampling window in seconds
const DefaultSampleWindow = 15 * time.Second

// DefaultHysteresis is the default hystersis window in seconds
const DefaultHysteresis = -15 * time.Second

// XXX Should there be defined types for the index and the list elements
// since both are really from an enumerated set of strings.  Example:
// type BmcCmdStr string
// type RedfishResetType string // should actually come from common source
// type BmcCmdToActionResetType map[string][]RedfishResetType

// MaxComponentQuery is the maximum number of components that we allow to be
// used in a query parameter when talking to HSM.
var MaxComponentQuery = 2048

// BmcCmdToActionResetTypeMap maps internal CAPMC BMC commands to the
// corresponding Redfish ResetType.
type BmcCmdToActionResetTypeMap map[string][]string

// CAPMC Configuration Defaults
// These values are used when there is no configuration file.
var (
	defaultActionMaxWorkers    = 1000
	defaultOnUnsupportedAction = actionSimulate
	defaultReinitActionSeq     = []string{bmcCmdPowerOff, bmcCmdPowerForceOff, bmcCmdPowerRestart, bmcCmdPowerForceRestart, bmcCmdPowerOn, bmcCmdPowerForceOn, bmcCmdNMI}
	defaultWaitForOffRetries   = 60
	defaultWaitForOffSleep     = 15
	// CompSeq:
	// The power sequencing list based on comments in CASMHMS-836
	// consists only of the following components:
	// Chassis       xXcC
	// RouterModule  xXcCrR
	// ComputeModule xXcCsS
	// Node          xXcCsSbBnN
	//
	// ResetType:
	// Default mapping for CAPMC `Off` operations to Redfish ResetType
	// GracefulSutdown - requires no additional checks
	// Off             - requires no additional checks (DTMF is adding)
	// PushPowerButton - requires additional checks (this is a toggle)
	// ForceOff        - requires no additional checks
	//
	// Default mapping for CAPMC `On` operations to Redfish ResetType
	// On              - requires no additional checks
	// PushPowerButton - requires additional checks (this is a toggle)
	// ForceOn         - requires no additional checks
	// PowerCycle      - requires additional checks (maybe)
	//
	// Default mapping for CAPMC `Reinit` operations to Redfish ResetType
	// GracefulRestart - requires no additional checks
	// ForceRestart    - requires no additional checks
	// PowerCycle      - requires no additional checks
	//                   (closer to ForceRestart)
	// PushPowerButton - requires additional checks and actions
	//                   (this is a toggle)
	// NOTE Order is from most to least preferred.
	defaultPowerControl = map[string]PowerCtl{
		// NOTE NMI is a future CAPMC enhancement.
		// Actual configuration TBD.
		bmcCmdNMI: {
			CompSeq:   []string{"Node"},
			ResetType: []string{"Nmi"},
		},
		bmcCmdPowerForceOff: {
			CompSeq:   []string{"Node", "ComputeModule", "HSNBoard", "RouterModule", "Chassis", "CabinetPDUOutlet", "CabinetPDUPowerConnector"},
			ResetType: []string{"ForceOff"},
		},
		bmcCmdPowerForceOn: {
			CompSeq:   []string{"CabinetPDUPowerConnector", "CabinetPDUOutlet", "Chassis", "RouterModule", "HSNBoard", "ComputeModule", "Node"},
			ResetType: []string{"ForceOn"},
		},
		bmcCmdPowerForceRestart: {
			CompSeq:   []string{"Node"},
			ResetType: []string{"ForceRestart", "PowerCycle"},
		},
		bmcCmdPowerOff: {
			CompSeq:   []string{"Node", "ComputeModule", "HSNBoard", "RouterModule", "Chassis", "CabinetPDUOutlet", "CabinetPDUPowerConnector"},
			ResetType: []string{"GracefulShutdown", "Off"},
		},
		bmcCmdPowerOn: {
			CompSeq:   []string{"CabinetPDUPowerConnector", "CabinetPDUOutlet", "Chassis", "RouterModule", "HSNBoard", "ComputeModule", "Node"},
			ResetType: []string{"On"},
		},
		bmcCmdPowerRestart: {
			CompSeq:   []string{"Node"},
			ResetType: []string{"GracefulRestart"},
		},
	}
	defaultSystemParameters = SystemParameters{
		PowerCapTarget: 0,
		PowerThreshold: 0,
		StaticPower:    0,
		RampLimited:    false,
		RampLimit:      2000000,
		PowerBandMin:   0,
		PowerBandMax:   0,
	}
	defaultCapmcConfiguration = CapmcConfiguration{
		ActionMaxWorkers:    defaultActionMaxWorkers,
		OnUnsupportedAction: defaultOnUnsupportedAction,
		ReinitActionSeq:     defaultReinitActionSeq,
		WaitForOffRetries:   defaultWaitForOffRetries,
		WaitForOffSleep:     defaultWaitForOffSleep,
	}
)

// TODO - figure out what is ideal (This struct is probably not ideal...)
type bmcPowerRc struct {
	ni    *NodeInfo
	rc    int
	msg   string
	state string
}

// The CAPMC command and Redfish payload
type bmcCmd struct {
	cmd     string
	payload []byte
}

// The single argument to the doBmcCall() function.
type bmcCall struct {
	bmcCmd
	ni      *NodeInfo
	rspChan chan bmcPowerRc
}

type nodeRsp struct {
	Name  string `json:"name"`
	Msg   string `json:"msg"`
	State string `json:"state"`
	Rc    int    `json:"rc"`
}

// CapmcD contains the communication information to enable talking with BMCs and
// Hardware State Manager (HSM)
type CapmcD struct {
	httpListen          string
	rfClient            *hms_certs.HTTPClientPair // HTTP Client for talking to BMCs
	smClient            *hms_certs.HTTPClientPair
	simulationOnly      bool // If true NO COMMANDS WILL BE SENT TO BMC's (defaults to false)
	debug               bool
	debugLevel          int
	hsmURL              *url.URL
	config              *Config
	ActionMaxWorkers    int
	OnUnsupportedAction string
	ReinitActionSeq     []string
	WPool               *base.WorkerPool
	db                  tsdb.TSDB
	ss                  sstorage.SecureStorage
	ccs                 *compcreds.CompCredStore
	reservation         reservation.Production
	reservationsEnabled bool
}

// TODO This maybe sub-optimal but it will do for now.  This is mainly
// just a copy of the old hardware management database NodeInfo structure
// with the boot services information dropped (hmdb/dpapi.go).
// All the JSON exporting is a hold over as well. Not clear on original
// intent but it can probably go away too.
// Ideally this wouldn't be required at all as the hardware state manager
// (HSM) would just return what is needed (as this structure is the synthesis
// of several different HSM "queries").

// NodeInfo encapsulates information about a node
type NodeInfo struct {
	Location      string
	Hostname      string
	Domain        string
	FQDN          string
	Nid           int
	Role          string
	State         string
	Enabled       bool
	BmcFQDN       string
	BmcPath       string
	RfActionURI   string
	RfResetTypes  []string
	RfEpoURI      string
	RfPowerURL    string
	RfPowerTarget string
	RfPwrCtlCnt   int
	RfControlsCnt int
	PowerCaps     map[string]PowerCap
	PowerETag     string
	Type          string
	RfType        string
	RfSubtype     string
	BmcUser       string
	BmcPass       string
	BmcProtocol   string
	BmcType       string
}

// PowerCap defines the values used for CAPMC power capping to Redfish
// PowerControl mapping and setting.
type PowerCap struct {
	// Name is the CAPMC 'control name' ("accel" or "node")
	Name string
	// Path is the URI for the Redfish resouce
	Path string
	// Min is the minimum acceptable value for the 'control'
	Min int
	// Max is the maximum acceptable value for the 'control'
	Max int
	// PwrCtlIndex is the zero-based array index in the Redfish
	// Power.PowerControl array this 'control' corresponds to.
	PwrCtlIndex int
}

// PowerOpRules defines a set of rules that apply to power operations.
type PowerOpRules struct {
	// MaxOffTime is the maximum time, in seconds, a component may be in
	// the `Off` state. -1 is unlimited.
	MaxOffTime int
	// MinOffTime is the mimimum time, in seconds, a component may be in
	// the `On` state. -1 is unlimited.
	MinOffTime int
	Off        OpRule
	On         OpRule
	Reinit     OpRule
}

// OpRule defines a set of constraints that apply to operations.
type OpRule struct {
	// Latency is the approximate time, in seconds, an operation is
	// expected to take for a component.
	Latency int
	// MaxReq is the maximum number of components which an opperation
	// can effect in one request.
	MaxReq int `toml:"MaxRequest"`
}

// SystemParameters defines the read-only expected worst case system power
// consumption, static power overhead, or administratively define dvalues such
// as a system wide power limit
type SystemParameters struct {
	// Administratively defined upper limit on system power
	PowerCapTarget int
	// System power level, which if crossed, will result in Cray management
	// software emitting over power budget warnings
	PowerThreshold int
	// Additional static system wide power overhead which is unreported,
	// specified in watts
	StaticPower int
	// True if out-of-band HSS power ramp rate limiting features are enabled
	RampLimited bool
	// Administratively defined maximum rate of change (increasing or
	// decreasing) in system wide power consumption, specified in watts per
	// minute
	RampLimit int
	// Administratively defined minimum allowable system power consumption,
	// specified in watts
	PowerBandMin int
	// Administratively defined maximum allowable system power consumption,
	// specified in watts
	PowerBandMax int
}

// CapmcConfiguration defines the internal knobs that can be modified inside
// of CAPMC to alter how certain portions of the code work.
type CapmcConfiguration struct {
	ActionMaxWorkers    int
	OnUnsupportedAction string
	ReinitActionSeq     []string
	WaitForOffRetries   int
	WaitForOffSleep     int
}

//PowerCapCapabilityMonikerType is consistent with the V3 XC moniker schema
//and is used as a template to separate nodes into groups of blade types
type PowerCapCapabilityMonikerType struct {
	Version          string `json:"version"`
	SSD              string `json:"ssd"`
	BaseBoardType    string `json:"base_board_type"`
	BaseBoardSubType string `json:"base_board_sub_type"`
	CPUID            string `json:"cpuid"`
	TDP              string `json:"tdp"`
	NumCores         string `json:"num_cores"`
	MemSizeGiB       string `json:"mem_size_gib"`
	MemSpeedMHZ      string `json:"mem_speed_mhz"`
	Accelerator      string `json:"accelerator"`
}

//PowerCapCapabilityMonikerGroup is used to represent a grouping of nodes by
//blade type, and is dynamically generated by the doPowerCapCabilities handler
type PowerCapCapabilityMonikerGroup struct {
	Name   string   `json:"name"`
	Desc   string   `json:"desc"`
	Xnames []string `json:"xnames"`
	Nids   []int    `json:"nids"`
}
