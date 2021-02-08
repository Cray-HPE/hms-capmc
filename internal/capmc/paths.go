// MIT License
//
// (C) Copyright [2021] Hewlett Packard Enterprise Development LP
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

package capmc

// We goofed in the original CAPMC API when we didn't put a version
// number in the URL. These constants define paths with and without
// the version number for backward compatibility with published
// documentation. Newer code should use the v1 paths.

// The Shasta implementation of the Cascade CAPMC APIs
const (
	ComputeNodeControlV1   = "/capmc/v1/cnctl"
	EmergencyPowerOffV1    = "/capmc/v1/emergency_power_off"
	GroupOnV1              = "/capmc/v1/group_on"
	GroupOffV1             = "/capmc/v1/group_off"
	GroupReinitV1          = "/capmc/v1/group_reinit"
	GroupStatusV1          = "/capmc/v1/get_group_status"
	HealthV1               = "/capmc/v1/health"
	LivenessV1             = "/capmc/v1/liveness"
	MCDRAMCapabilitiesV1   = "/capmc/v1/get_mcdram_capabilities"
	MCDRAMConfigClearV1    = "/capmc/v1/clr_mcdram_cfg"
	MCDRAMConfigGetV1      = "/capmc/v1/get_mcdram_cfg"
	MCDRAMConfigSetV1      = "/capmc/v1/set_mcdram_cfg"
	NUMACapabilitiesV1     = "/capmc/v1/get_numa_capabilities"
	NUMAConfigClearV1      = "/capmc/v1/clr_numa_cfg"
	NUMAConfigGetV1        = "/capmc/v1/get_numa_cfg"
	NUMAConfigSetV1        = "/capmc/v1/set_numa_cfg"
	NodeEnergyCounterV1    = "/capmc/v1/get_node_energy_counter"
	NodeEnergyStatsV1      = "/capmc/v1/get_node_energy_stats"
	NodeEnergyV1           = "/capmc/v1/get_node_energy"
	NodeIDMapV1            = "/capmc/v1/get_nid_map"
	NodeOffV1              = "/capmc/v1/node_off"
	NodeOnV1               = "/capmc/v1/node_on"
	NodeReinitV1           = "/capmc/v1/node_reinit"
	NodeRulesV1            = "/capmc/v1/get_node_rules"
	NodeStatusV1           = "/capmc/v1/get_node_status"
	PartitionMapV1         = "/capmc/v1/get_partition_map"
	PowerBiasClearV1       = "/capmc/v1/clr_power_bias"
	PowerBiasComputeV1     = "/capmc/v1/compute_power_bias"
	PowerBiasDataV1        = "/capmc/v1/set_power_bias_data"
	PowerBiasGetV1         = "/capmc/v1/get_power_bias"
	PowerBiasSetV1         = "/capmc/v1/set_power_bias"
	PowerCapCapabilitiesV1 = "/capmc/v1/get_power_cap_capabilities"
	PowerCapGetV1          = "/capmc/v1/get_power_cap"
	PowerCapSetV1          = "/capmc/v1/set_power_cap"
	ReadinessV1            = "/capmc/v1/readiness"
	SSDDiagsGetV1          = "/capmc/v1/get_ssd_diags"
	SSDEnableClearV1       = "/capmc/v1/clr_ssd_enable"
	SSDEnableGetV1         = "/capmc/v1/get_ssd_enable"
	SSDEnableSetV1         = "/capmc/v1/set_ssd_enable"
	SSDSGetV1              = "/capmc/v1/get_ssds"
	SystemParamsV1         = "/capmc/v1/get_system_parameters"
	SystemPowerDetailsV1   = "/capmc/v1/get_system_power_details"
	SystemPowerV1          = "/capmc/v1/get_system_power"
	XnameOffV1             = "/capmc/v1/xname_off"
	XnameOnV1              = "/capmc/v1/xname_on"
	XnameReinitV1          = "/capmc/v1/xname_reinit"
	XnameStatusV1          = "/capmc/v1/get_xname_status"
)

// The original (Cascade) CAPMC APIs
const (
	ComputeNodeControlV0   = "/capmc/v0/cnctl"
	MCDRAMCapabilitiesV0   = "/capmc/v0/get_mcdram_capabilities"
	MCDRAMConfigClearV0    = "/capmc/v0/clr_mcdram_cfg"
	MCDRAMConfigV0         = "/capmc/v0/get_mcdram_cfg"
	NUMACapabilitiesV0     = "/capmc/v0/get_numa_capabilities"
	NUMAConfigClearV0      = "/capmc/v0/clr_numa_cfg"
	NUMAConfigGetV0        = "/capmc/v0/get_numa_cfg"
	NodeEnergyCounterV0    = "/capmc/v0/get_node_energy_counter"
	NodeEnergyStatsV0      = "/capmc/v0/get_node_energy_stats"
	NodeEnergyV0           = "/capmc/v0/get_node_energy"
	NodeIDMapV0            = "/capmc/v0/get_nid_map"
	NodeOffV0              = "/capmc/v0/node_off"
	NodeOnV0               = "/capmc/v0/node_on"
	NodeReinitV0           = "/capmc/v0/node_reinit"
	NodeRulesV0            = "/capmc/v0/get_node_rules"
	NodeStatusV0           = "/capmc/v0/get_node_status"
	PartitionMapV0         = "/capmc/v0/get_partition_map"
	PowerBiasClearV0       = "/capmc/v0/clr_power_bias"
	PowerBiasComputeV0     = "/capmc/v0/compute_power_bias"
	PowerBiasDataV0        = "/capmc/v0/set_power_bias_data"
	PowerBiasGetV0         = "/capmc/v0/get_power_bias"
	PowerBiasSetV0         = "/capmc/v0/set_power_bias"
	PowerCapCapabilitiesV0 = "/capmc/v0/get_power_cap_capabilities"
	PowerCapGetV0          = "/capmc/v0/get_power_cap"
	PowerCapSetV0          = "/capmc/v0/set_power_cap"
	SSDDiagsGetV0          = "/capmc/v0/get_ssd_diags"
	SSDEnableClearV0       = "/capmc/v0/clr_ssd_enable"
	SSDEnableGetV0         = "/capmc/v0/get_ssd_enable"
	SSDEnableSetV0         = "/capmc/v0/get_ssd_enable"
	SystemParamsV0         = "/capmc/v0/get_system_parameters"
	SystemPowerDetailsV0   = "/capmc/v0/get_system_power_details"
	SystemPowerV0          = "/capmc/v0/get_system_power"
)

// The default CAPMC APIs
// NOTE This may not be the best...
// It is assumed that these point to the latest supported API version.
const (
	ComputeNodeControl   = "/capmc/cnctl"
	EmergencyPowerOff    = "/capmc/emergency_power_off"
	GroupOn              = "/capmc/group_on"
	GroupOff             = "/capmc/group_off"
	GroupReinit          = "/capmc/group_reinit"
	GroupStatus          = "/capmc/get_group_status"
	Health               = "/capmc/health"
	Liveness             = "/capmc/liveness"
	MCDRAMCapabilities   = "/capmc/get_mcdram_capabilities"
	MCDRAMConfigClear    = "/capmc/clr_mcdram_cfg"
	MCDRAMConfigGet      = "/capmc/get_mcdram_cfg"
	MCDRAMConfigSet      = "/capmc/set_mcdram_cfg"
	NUMACapabilities     = "/capmc/get_numa_capabilities"
	NUMAConfigClear      = "/capmc/clr_numa_cfg"
	NUMAConfigGet        = "/capmc/get_numa_cfg"
	NUMAConfigSet        = "/capmc/set_numa_cfg"
	NodeEnergyCounter    = "/capmc/get_node_energy_counter"
	NodeEnergyStats      = "/capmc/get_node_energy_stats"
	NodeEnergy           = "/capmc/get_node_energy"
	NodeIDMap            = "/capmc/get_nid_map"
	NodeOff              = "/capmc/node_off"
	NodeOn               = "/capmc/node_on"
	NodeReinit           = "/capmc/node_reinit"
	NodeRules            = "/capmc/get_node_rules"
	NodeStatus           = "/capmc/get_node_status"
	PartitionMap         = "/capmc/get_partition_map"
	PowerBiasClear       = "/capmc/clr_power_bias"
	PowerBiasCompute     = "/capmc/compute_power_bias"
	PowerBiasData        = "/capmc/set_power_bias_data"
	PowerBiasGet         = "/capmc/get_power_bias"
	PowerBiasSet         = "/capmc/set_power_bias"
	PowerCapCapabilities = "/capmc/get_power_cap_capabilities"
	PowerCapGet          = "/capmc/get_power_cap"
	PowerCapSet          = "/capmc/set_power_cap"
	Readiness            = "/capmc/readiness"
	SSDDiagsGet          = "/capmc/get_ssd_diags"
	SSDEnableClear       = "/capmc/clr_ssd_enable"
	SSDEnableGet         = "/capmc/get_ssd_enable"
	SSDEnableSet         = "/capmc/set_ssd_enable"
	SSDSGet              = "/capmc/get_ssds"
	SystemParams         = "/capmc/get_system_parameters"
	SystemPowerDetails   = "/capmc/get_system_power_details"
	SystemPower          = "/capmc/get_system_power"
	XnameOff             = "/capmc/xname_off"
	XnameOn              = "/capmc/xname_on"
	XnameReinit          = "/capmc/xname_reinit"
	XnameStatus          = "/capmc/get_xname_status"
)
