// MIT License
//
// (C) Copyright [2017-2019, 2021] Hewlett Packard Enterprise Development LP
//
// Permission is hereby granted, free of charge, to any person obtaining a
// copy of this software and associated documentation files (the "Software"),
// to deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included
// in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
// THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
// OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
// ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.
//

package capmc

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"syscall"
)

// CAPMC API Structs
// ========================================================

type ErrResponse struct {
	E      int    `json:"e"` // Error code
	ErrMsg string `json:"err_msg"`
}

// Common Error Respsonses
// Cascade CAPMC used errno and strerror(errno) extensively for the error
// response. These are the most common cases.
var (
	ErrResponseEINVAL = ErrResponse{E: int(syscall.EINVAL), ErrMsg: syscall.EINVAL.Error()}
)

// Node Power and Status
// --------------------------------------------------------

type GetNodeRulesResponse struct {
	ErrResponse
	LatencyNodeOff    int `json:"latency_node_off"`
	LatencyNodeOn     int `json:"latency_node_on"`
	LatencyNodeReinit int `json:"latency_node_reinit"`
	MaxOffReqCount    int `json:"max_off_req_count"`
	MaxOffTime        int `json:"max_off_time"`
	MaxOnReqCount     int `json:"max_on_req_count"`
	MaxReinitReqCount int `json:"max_reinit_req_count"`
	MinOffTime        int `json:"min_off_time"`
}

// NodeStatusRequest is the API POST body for get_node_status.
type NodeStatusRequest struct {
	Filter string `json:"filter,omitempty"`
	Source string `json:"source,omitempty"`
	Nids   []int  `json:"nids,omitempty"`
}

// NodeStatusFlags contains arrays of the possible flags from the original
// undocumented Cascade CAPMC API. The only documentation is the capmc(1)
// command `capmc help --verbose`.
type NodeStatusFlags struct {
	Warning  []int `json:"warn,omitempty"`
	Alert    []int `json:"alert,omitempty"`
	Reserved []int `json:"resvd,omitempty"`
}

// NodeStatusResponse contains arrays of the possible statuses (states)
// from the original Cascade CAPMC API. Empirically, this struct also contains
// a "flags" member, though this isn't mentioned anywhere in the API
// (specification) documentation. There are additional Shasta HSM states, but
// not reported (for Cascade API compatibility) as node status: Populated,
// Unknown. Empty is a Cascade state which CAPMC does not report. It is
// here for completeness and possible inclusion in internal interfaces..
// The published Cascade API specification does not cover the "diag"
// state/status but it was supported by the implementation. Shasta does not
// have a "diag" state so there will not be a "Diag" node list.
// The original API didn't have "indeterminate" or "error" states.
// It's likely that with the Redfish backend handling we'll encounter
// errors trying to talk to the hardware. Status Undefined is for such errors.
type NodeStatusResponse struct {
	ErrResponse
	Diag      []int            `json:"diag,omitempty"`
	Disabled  []int            `json:"disabled,omitempty"`
	Empty     []int            `json:"empty,omitempty"`
	Halt      []int            `json:"halt,omitempty"`
	Off       []int            `json:"off,omitempty"`
	On        []int            `json:"on,omitempty"`
	Populated []int            `json:"populated,omitempty"`
	Ready     []int            `json:"ready,omitempty"`
	Standby   []int            `json:"standby,omitempty"`
	Undefined []int            `json:"undefined,omitempty"`
	Unknown   []int            `json:"unknown,omitempty"`
	Flags     *NodeStatusFlags `json:"flags,omitempty"`
}

// The following Xname structs mostly mirror the Node structs above.

// XnameStatusRequest is the API POST body for a get_xname_status request.
type XnameStatusRequest struct {
	Filter string   `json:"filter,omitempty"`
	Source string   `json:"source,omitempty"`
	Xnames []string `json:"xnames,omitempty"`
}

// XnameStatusRequest contains arrays of the possible HMSFlags.
type XnameStatusFlags struct {
	// The valid Shasta HSM Component Flag values
	Alert   []string `json:"alert,omitempty"`
	Warning []string `json:"warning,omitempty"`
	Locked  []string `json:"locked,omitempty"`
	OK      []string `json:"ok,omitempty"`
	Unknown []string `json:"unknown,omitempty"`

	// Shasta doesn't have Disabled as a State or Flag. The boolean Enabled
	// is a separate indicator for each Component. This is a synthetic
	// list of components where Enabled == false regardless of State.
	Disabled []string `json:"disabled,omitempty"`
}

// XnameStatusResponse contains the common HMS States used by the HSM
// as CAPMC Status. There are two statuses defined by Shasta CAPMC as
// the original Cascade API (node_status) this was modeled after didn't
// have "indeterminate" or "error" states: undefined and unresponsive.
// It's likely that with the Redfish backend handling we'll encounter
// errors trying to talk to the hardware. We need an extra state.
// Early versions of CAPMC used Undefined for any "indeterminate" or "error"
// states. Possibly should have used HMS commonly defined "unknown". The
// new status "unresponsive" indicates that CAPMC was not able to communicate
// with the node for some reason.
type XnameStatusResponse struct {
	ErrResponse
	Empty        []string          `json:"empty,omitempty"`
	Halt         []string          `json:"halt,omitempty"`
	Off          []string          `json:"off,omitempty"`
	On           []string          `json:"on,omitempty"`
	Populated    []string          `json:"populated,omitempty"`
	Ready        []string          `json:"ready,omitempty"`
	Standby      []string          `json:"standby,omitempty"`
	Undefined    []string          `json:"undefined,omitempty"`
	Unknown      []string          `json:"unknown,omitempty"`
	Unresponsive []string          `json:"unresponsive,omitempty"`
	Flags        *XnameStatusFlags `json:"flags,omitempty"`
}

// The original node status API uses a pipe delimited string to pass
// fliter arguments. This is non-ideal. It would have been better to use an
// array of state names instead.
//
// States with the initial Whitebox Redfish will be off/on/ready/undefined.
// Filter names & states don't accout for "undefined". This will need some
// more effort.
//
// The Shasta HMS States are not a one-for-one match between the Cascade HSS
// States that the Cascade CAPMC API specified as 'status'. The following
// filters represent all of the possible Shasta states along with those that
// were present in Cascade. Here are the differences:
// * The show_empty filter is a Shasta extension. This filter is not valid
//   for the 'node_status' API.
// * The show_diag filter is deprecated for Shasta as HMS no longer has a
//   diag(nostic) state. The filter is accepted but will never return anything
//   in diag.
// * The show_disabled filter does not directly translate into a state as
//   Shasta HMS no longer has a disabled state. Components are enabled/disabled
//   via a separate Enabled field which does not affect Component State.
// * The show_populated filter is a Shasta extension. This state is not
//   currently set by HMS software. This filter is not valid for the
//   'node_status' API.
// * The show_undefined filter is not a true HMS state. This filter is not
//   valid for the node_status' API. The undefined 'state' is only used
//   when CAPMC uses Redfish to determine status.
// * The show_unknown filter is a Shasta extension. This is both a state
//   and a flag. Neither is currently used by HMS software.
// * The show_locked filter is the same as the Cascade show_resvd flag
//   filter. The show_locked filter is not valid for the 'node_status' API.
// * All others are present in Cascade. Note not all these filters are
//   documented by the Cascade API document but were documented for the
//   capmc(1) command.

const (
	FilterDelimiter = "|"

	// States
	FilterShowAll       = "show_all"
	FilterShowDiag      = "show_diag"
	FilterShowDisabled  = "show_disabled" // not a Shasta HSM state
	FilterShowEmpty     = "show_empty"
	FilterShowHalt      = "show_halt"
	FilterShowOff       = "show_off"
	FilterShowOn        = "show_on"
	FilterShowPopulated = "show_populated"
	FilterShowReady     = "show_ready"
	FilterShowStandby   = "show_standby"
	FilterShowUndefined = "show_undefined" // not a true "state"
	FilterShowUnknown   = "show_unknown"

	// Flags
	FilterShowAlert    = "show_alert"
	FilterShowLocked   = "show_locked"
	FilterShowReserved = "show_reserved"
	FilterShowResvd    = "show_resvd"
	FilterShowWarn     = "show_warn"
	FilterShowWarning  = "show_warning"
)

// These bits are sort of arbitrary at the moment. Apps should use the
// bits instead of strings.
const (
	// States
	FilterShowDiagBit = (1 << iota)
	FilterShowDisabledBit
	FilterShowEmptyBit
	FilterShowHaltBit
	FilterShowOffBit
	FilterShowOnBit
	FilterShowPopulatedBit
	FilterShowReadyBit
	FilterShowStandbyBit
	FilterShowUndefinedBit
	FilterShowUnknownBit

	// Flags
	FilterShowAlertBit
	FilterShowReservedBit
	FilterShowWarningBit
	FilterShowOKBit

	// Includes all flags *except* OK as generally OK isn't interesting
	FilterShowAllBit = (FilterShowAlertBit | FilterShowDiagBit |
		FilterShowDisabledBit | FilterShowEmptyBit |
		FilterShowHaltBit | FilterShowOffBit |
		FilterShowOnBit | FilterShowPopulatedBit |
		FilterShowReadyBit | FilterShowReservedBit |
		FilterShowStandbyBit | FilterShowUndefinedBit |
		FilterShowUnknownBit | FilterShowWarningBit)

	// alias
	FilterShowLockedBit = FilterShowReservedBit
)

// NodeStatusFilterParse parses the incoming pipe ('|') delimited string
// and returns its bitmask representation. This is a specialized version
// of the more general StatusFilterPase that only accepts filters valid for
// the get_node_status API.
func NodeStatusFilterParse(filter string) (uint, error) {
	// List of filters that aren't applicable to Node Status as per the
	// original Cascade CAPMC API (and command capmc(1)) document.
	var invalidNodeFilters = []string{
		FilterShowEmpty,
		FilterShowPopulated,
		FilterShowUndefined,
		FilterShowUnknown,
	}
	const invalidNodeFilterBits = (FilterShowEmptyBit | FilterShowPopulatedBit | FilterShowUndefinedBit | FilterShowUnknownBit)

	if len(filter) < 1 {
		return FilterShowAllBit &^ invalidNodeFilterBits, nil
	}

	// Check for invalid filters
	for _, token := range invalidNodeFilters {
		if strings.Contains(filter, token) {
			return 0, errors.New("invalid filter string: " + token)
		}
	}

	bitmap, err := StatusFilterParse(filter)

	return bitmap &^ invalidNodeFilterBits, err
}

// StatusFilterParse parses the incoming pipe delimited filter string and
// returns a bit mask representation. An empty string matches "show_all".
// Return non-null error if input is invalid.
func StatusFilterParse(filter string) (uint, error) {
	var bitmap uint

	if len(filter) < 1 {
		return FilterShowAllBit, nil
	}

	tokens := strings.Split(strings.ToLower(filter), FilterDelimiter)

	for _, token := range tokens {
		switch strings.TrimSpace(token) {
		case FilterShowAlert:
			bitmap |= FilterShowAlertBit
		case FilterShowAll:
			bitmap |= FilterShowAllBit
		case FilterShowDiag:
			bitmap |= FilterShowDiagBit
		case FilterShowDisabled:
			bitmap |= FilterShowDisabledBit
		case FilterShowEmpty:
			bitmap |= FilterShowEmptyBit
		case FilterShowHalt:
			bitmap |= FilterShowHaltBit
		case FilterShowOff:
			bitmap |= FilterShowOffBit
		case FilterShowOn:
			bitmap |= FilterShowOnBit
		case FilterShowPopulated:
			bitmap |= FilterShowPopulatedBit
		case FilterShowReady:
			bitmap |= FilterShowReadyBit
		case FilterShowReserved, FilterShowResvd:
			bitmap |= FilterShowReservedBit
		case FilterShowStandby:
			bitmap |= FilterShowStandbyBit
		case FilterShowUndefined:
			bitmap |= FilterShowUndefinedBit
		case FilterShowUnknown:
			bitmap |= FilterShowUnknownBit
		case FilterShowWarning, FilterShowWarn:
			bitmap |= FilterShowWarningBit
		default:
			return 0, errors.New("invalid filter string: " + token)
		}
	}

	return bitmap, nil
}

// NodePowerRequest - Same for node_on, node_off, node_reinit
type NodePowerRequest struct {
	Nids   []int  `json:"nids"`
	Force  bool   `json:"force,omitempty"`
	Reason string `json:"reason,omitempty"`
}

type NodePowerNidErr struct {
	Nid int `json:"nid"`
	ErrResponse
}

type NodePowerResponse struct {
	ErrResponse
	Nids []*NodePowerNidErr `json:"nids,omitempty"`
}

type XnameControlErr struct {
	Xname string `json:"xname"`
	ErrResponse
}

// Creates an XnameControlErr
func MakeXnameError(xname string, err int, msg string) *XnameControlErr {
	xnameErr := new(XnameControlErr)
	xnameErr.Xname = xname
	xnameErr.ErrResponse.E = err
	xnameErr.ErrResponse.ErrMsg = msg
	return xnameErr
}

type XnameControlResponse struct {
	ErrResponse
	Xnames []*XnameControlErr `json:"xnames,omitempty"`
}

// Node Capabilities and Power Control
// --------------------------------------------------------

// Same for get_power_cap_capabilties, get_power_cap, get_power_bias,
// clr_power_bias, get_nid_map, get_mcdram_capabilities, get_mcdram_cfg,
// clr_mcdram_cfg, get_numa_capabilities, get_numa_cfg, clr_numa_cfg,
// get_ssd_enable, clr_ssd_enable
type NidlistRequest struct {
	Nids []int `json:"nids"`
}

type PowerCapCapabilityControl struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
	Max  int    `json:"max"`
	Min  int    `json:"min"`
}

type PowerCapGroup struct {
	Name         string                      `json:"name"`
	Desc         string                      `json:"desc"`
	HostLimitMax int                         `json:"host_limit_max"`
	HostLimitMin int                         `json:"host_limit_min"`
	Static       int                         `json:"static"`
	Supply       int                         `json:"supply"`
	Powerup      int                         `json:"powerup"`
	Nids         []int                       `json:"nids"`
	Controls     []PowerCapCapabilityControl `json:"controls"`
}

type GetPowerCapCapabilitiesResponse struct {
	ErrResponse
	Groups []PowerCapGroup `json:"groups,omitempty"`
	Nids   []PowerCapNid   `json:"nids,omitempty"`
}

type PowerCapControl struct {
	Name string `json:"name"`
	Val  *int   `json:"val"`
}

// On error, this struct contains Nid, E, and ErrMsg. Otherwise, it contains
// Nid and Controls.
type PowerCapNid struct {
	Nid      int               `json:"nid"`
	Controls []PowerCapControl `json:"controls,omitempty"`
	E        int               `json:"e,omitempty"` // Error code
	ErrMsg   string            `json:"err_msg,omitempty"`
}

// Same for get_power_cap, set_power_cap
type PowerCapResponse struct {
	ErrResponse
	Nids []PowerCapNid `json:"nids"`
}

type SetPowerCapRequest struct {
	Nids []PowerCapNid `json:"nids"`
}

type PowerBiasNid struct {
	Nid       int     `json:"nid"`
	PowerBias float64 `json:"power-bias"`
}

type GetPowerBiasResponse struct {
	ErrResponse
	Nids []PowerBiasNid `json:"nids"`
}

type SetPowerBiasRequest struct {
	Nids []PowerBiasNid `json:"nids"`
}

type PowerBiasDataNid struct {
	Avgpwr int `json:"avgpwr"`
	Nid    int `json:"nid"`
}

type SetPowerBiasDataRequest struct {
	App  string             `json:"app"`
	Nids []PowerBiasDataNid `json:"nids"`
}

type ComputePowerBiasRequest struct {
	App  string `json:"app"`
	Nids []int  `json:"nids"`
}

type ComputePowerBiasResponse struct {
	ErrResponse
	Nids []PowerBiasNid `json:"nids"`
}

// Node Frequency and Sleep State Control
// --------------------------------------------------------

// TODO: review this section when actually implementing code that exercises it.
//       These requests have a different format than most others, and it maybe makes
//       sense to format these structs differently.

type PwrAttrName struct {
	PwrAttrName string `json:"PWR_AttrName"`
}

// TODO: better name
type Pwr struct {
	PwrFunction     string        `json:"PWR_Function"`
	PwrAttrs        []PwrAttrName `json:"PWR_Attrs"`
	PwrMajorVersion int           `json:"PWR_MajorVersion"`
	PwrMinorVersion int           `json:"PWR_MinorVersion"`
}

// Shared with various PWR_* requests
type CnctlRequest struct {
	Nids []int `json:"nids"`
	Data Pwr   `json:"data"`
}

type PwrAttr struct {
	PwrAttrName        string `json:"PWR_AttrName"`
	PwrAttrValue       int    `json:"PWR_AttrValue"`
	PwrReturnCode      int    `json:"PWR_ReturnCode"`
	PwrTimeNanoseconds int    `json:"PWR_TimeNanoseconds"`
	PwrTimeSeconds     int    `json:"PWR_TimeSeconds"`
}

type PwrAttrData struct {
	PwrAttrs         []PwrAttr `json:"PWR_Attrs"`
	PwrErrorMessages string    `json:"PWR_ErrorMessages"`
	PwrMajorVersion  int       `json:"PWR_MajorVersion"`
	PwrMessages      string    `json:"PWR_Messages"`
	PwrMinorVersion  int       `json:"PWR_MinorVersion"`
	PwrReturnCode    int       `json:"PWR_ReturnCode"`
}

type NidPwrAttrData struct {
	Nid  int         `json:"nid"`
	Data PwrAttrData `json:"data"`
}

type GetFreqCapabilitiesResponse struct {
	ErrResponse
	Nids []NidPwrAttrData `json:"nids"`
}

// TODO: too many variations of this struct...
type PwrAttrNameValue struct {
	PwrAttrName  string `json:"PWR_AttrName"`
	PwrAttrValue string `json:"PWR_AttrValue"`
}

// Node MCDRAM Capabilities and Control
// --------------------------------------------------------

type McdramCfgNid struct {
	ErrResponse
	McdramCfg string `json:"mcdram_cfg"`
	Nid       int    `json:"nid"`
}

type GetMcdramCapabilitiesResponse struct {
	ErrResponse
	Nids []McdramCfgNid `json:"nids"`
}

type McdramFullCfgNid struct {
	McdramCfg  string `json:"mcdram_cfg"`
	McdramPct  int    `json:"mcdram_pct"`
	McdramSize string `json:"mcdram_size"`
	Nid        int    `json:"nid"`
	DramSize   string `json:"dram_size"`
}

type GetMcdramCfgResponse struct {
	ErrResponse
	Nids []McdramFullCfgNid `json:"nids"`
}

type SetMcdramCfgRequest struct {
	Nids []McdramCfgNid `json:"nids"`
}

// Node NUMA Capabilities and Control
// --------------------------------------------------------

type NumaCfgCapNid struct {
	ErrResponse
	NumaCfg string `json:"numa_cfg"`
	Nid     int    `json:"nid"`
}

type GetNumaCapabilitiesResponse struct {
	ErrResponse
	Nids []NumaCfgCapNid `json:"nids"`
}

type NumaCfgNid struct {
	NumaCfg string `json:"numa_cfg"`
	Nid     int    `json:"nid"`
}

type GetNumaCfgResponse struct {
	ErrResponse
	Nids []NumaCfgNid `json:"nids"`
}

type SetNumaCfgRequest struct {
	Nids []NumaCfgNid `json:"nids"`
}

// Node SSD Control and Diagnostics
// --------------------------------------------------------

type SsdEnableNid struct {
	SsdEnable int `json:"ssd_enable"`
	Nid       int `json:"nid"`
}

type GetSsdEnableResponse struct {
	ErrResponse
	Nids []SsdEnableNid `json:"nids"`
}

type SetSsdEnableRequest struct {
	Nids []SsdEnableNid `json:"nids"`
}

// Used by get_ssds, get_ssd_diags
type NidCnamelist struct {
	Nids  []int    `json:"nids,omitempty"`
	Cname []string `json:"cname,omitempty"`
}

// Used with "get_ssds" response structure. Unlike other API functions,
// get_ssds returns a map with variable-named keys (and SsdInfo structs as
// values). It's therefore not possible to represent the data as a simple
// struct type.
type SsdInfo struct {
	Bus          int    `json:"bus"`
	Device       int    `json:"device"`
	Func         int    `json:"func"`
	ModelNumber  string `json:"model_number"`
	Nid          int    `json:"nid"`
	SerialNumber string `json:"serial_number"`
	Size         int    `json:"size"`
	SsdId        string `json:"ssd_id"`
	SubId        string `json:"sub_id"`
}

type SsdDiag struct {
	LifeRemaining float64 `json:"life_remaining"`
	Firmware      string  `json:"firmware"`
	Ts            string  `json:"ts"`
	Nid           int     `json:"nid"`
	SerialNum     string  `json:"serial_num"`
	Cname         string  `json:"cname"`
	ManuId        string  `json:"manu_id"`
	PercentUsed   float64 `json:"percent_used"`
	PartId        string  `json:"part_id"`
	CompOrd       string  `json:"comp_ord"`
	Size          int     `json:"size"`
}

type GetSsdDiagsResponse struct {
	ErrResponse
	SsdDiags []SsdDiag `json:"ssd_diags"`
}

// Node Energy Reporting
// --------------------------------------------------------

// Same for get_node_energy, get_node_energy_stats
type TimeBoundNidRequest struct {
	StartTime string `json:"start_time,omitempty"`
	EndTime   string `json:"end_time,omitempty"`
	Nids      []int  `json:"nids,omitempty"`
	Apid      string `json:"apid,omitempty"`
	JobId     string `json:"job_id,omitempty"`
}

type NidEnergy struct {
	Nid    int `json:"nid"`
	Energy int `json:"energy"`
}

type GetNodeEnergyResponse struct {
	ErrResponse
	NidCount *int         `json:"nid_count,omitempty"`
	Time     *float64     `json:"time,omitempty"`
	Nodes    []*NidEnergy `json:"nodes,omitempty"`
}

type GetNodeEnergyStatsResponse struct {
	ErrResponse
	Time        *float64 `json:"time"`
	NidCount    *int     `json:"nid_count"`
	EnergyTotal *int     `json:"energy_total"`
	EnergyAvg   *float64 `json:"energy_avg"`
	EnergyStd   *float64 `json:"energy_std"`
	EnergyMax   []*int   `json:"energy_max"`
	EnergyMin   []*int   `json:"energy_min"`
}

type GetNodeEnergyCounterRequest struct {
	Time  string `json:"time,omitempty"`
	Nids  []int  `json:"nids,omitempty"`
	Apid  string `json:"apid,omitempty"`
	JobId string `json:"job_id,omitempty"`
}

type NidEnergyCounter struct {
	Nid       *int   `json:"nid"`
	EnergyCtr *int   `json:"energy_ctr"`
	Time      string `json:"time"`
}

type GetNodeEnergyCounterResponse struct {
	ErrResponse
	NidCount *int               `json:"nid_count"`
	Nodes    []NidEnergyCounter `json:"nodes"`
}

// System Level Monitoring
// --------------------------------------------------------

type GetSystemParametersResponse struct {
	ErrResponse
	PowerCapTarget int  `json:"power_cap_target"`
	PowerThreshold int  `json:"power_threshold"`
	StaticPower    int  `json:"static_power"`
	RampLimited    bool `json:"ramp_limited"`
	RampLimit      int  `json:"ramp_limit"`
	PowerBandMin   int  `json:"power_band_min"`
	PowerBandMax   int  `json:"power_band_max"`
}

// Same for get_system_power_request, get_system_power_details
type TimeWindowRequest struct {
	StartTime string `json:"start_time"`
	WindowLen *int   `json:"window_len"`
}

type GetSystemPowerResponse struct {
	ErrResponse
	WindowLen int    `json:"window_len"`
	StartTime string `json:"start_time"`
	Avg       *int   `json:"avg"`
	Max       *int   `json:"max"`
	Min       *int   `json:"min"`
}

type CabinetPowerInfo struct {
	Avg float64 `json:"avg"`
	Max int     `json:"max"`
	Min int     `json:"min"`
	X   int     `json:"x"`
	Y   int     `json:"y"`
}

type GetSystemPowerDetailsResponse struct {
	ErrResponse
	WindowLen int                 `json:"window_len,omitempty"`
	StartTime string              `json:"start_time,omitempty"`
	Cabinets  []*CabinetPowerInfo `json:"cabinets,omitempty"`
}

// Xname Component Capabilities and Control
// --------------------------------------------------------

// XnameControl - Same for xname_on, xname_off
// Also used by emergency_power_off but Force, Recurse, and Prereq are ignored
type XnameControl struct {
	Xnames   []string `json:"xnames"`
	Force    bool     `json:"force,omitempty"`
	Reason   string   `json:"reason,omitempty"`
	Recurse  bool     `json:"recursive,omitempty"`
	Prereq   bool     `json:"prereq,omitempty"`
	Continue bool     `json:"continue,omitempty"`
}

// Group Component Capabilities and Control
// --------------------------------------------------------

// GroupControl - Same for group_on, group_off
type GroupControl struct {
	Groups []string `json:"groups"`
	Filter string   `json:"filter,omitempty"`
	Force  bool     `json:"force,omitempty"`
	Reason string   `json:"reason,omitempty"`
}

// Utility Functions
// --------------------------------------------------------

// NidInfoV0 defines the original Cascade CAPMC API NID Info structure.
type NidInfoV0 struct {
	Cname  string `json:"cname,omitempty"`
	Nid    int    `json:"nid"`
	Role   string `json:"role,omitempty"`
	E      int    `json:"e,omitempty"` // Error code
	ErrMsg string `json:"err_msg,omitempty"`
}

type NidInfo struct {
	Nid    int    `json:"nid"`
	Xname  string `json:"xname,omitempty"`
	Role   string `json:"role,omitempty"`
	E      int    `json:"e,omitempty"` // Error code
	ErrMsg string `json:"err_msg,omitempty"`
}

// GetNidMapResponseV0 is the original Cascade CAPMC API NID Info API response.
type GetNidMapResponseV0 struct {
	ErrResponse
	NIds []*NidInfoV0 `json:"nids,omitempty"`
}

type GetNidMapResponse struct {
	ErrResponse
	Nids []*NidInfo `json:"nids,omitempty"`
}

type PartitionInfo struct {
	Partition int   `json:"partition"`
	Nids      []int `json:"nids,omitempty"`
}

type GetPartitionMapResponse struct {
	ErrResponse
	Partitions []PartitionInfo `json:"partitions,omitempty"`
}

// CAPMC helper functions
// ========================================================

// ParseNidlist accepts a string of comma separated ints (or int ranges, e.g.,
// "23-29") and produces a slice of ints containing all those specified in the
// string. E.g., "1, 2,3-5,6" yields [1, 2, 3, 4, 5, 6].
func ParseNidlist(nidliststr string) ([]int, error) {
	nidlist := []int{}
	if nidliststr != "" {
		nids := strings.Split(nidliststr, ",")
		for _, s := range nids {
			if strings.Contains(s, "-") {
				rangebounds := strings.Split(s, "-")
				if len(rangebounds) != 2 {
					return nil, fmt.Errorf("invalid nid range: %s", s)
				}
				lowerBound, err1 := strconv.Atoi(rangebounds[0])
				upperBound, err2 := strconv.Atoi(rangebounds[1])
				if err1 != nil || err2 != nil {
					return nil, fmt.Errorf("invalid nid range: %s", s)
				}
				if lowerBound > upperBound {
					return nil, fmt.Errorf("invalid nid range: %s", s)
				}
				for i := lowerBound; i <= upperBound; i++ {
					nidlist = append(nidlist, i)
				}
			} else {
				n, err := strconv.Atoi(s)
				if err != nil {
					return nil, fmt.Errorf("invalid nid: %s", s)
				}
				nidlist = append(nidlist, n)
			}
		}
	}
	return nidlist, nil
}
