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

package capmc

// Structures that enable marshalling and unmarshalling of the Power schema as
// defined by the Redfish Power.v1_5_4.json and
// HpeServerAccPowerLimit.v1_0_0.HpeServerAccPowerLimit

// Power struct used to unmarshal the Redfish Power data that contains the power
// metrics, consumption, and limiting information
type Power struct {
	OContext string `json:"@odata.context,omitempty"`
	Oetag    string `json:"@odata.etag,omitempty"`
	Oid      string `json:"@odata.id,omitempty"`
	Otype    string `json:"@odata.type,omitempty"`
	Id       string `json:"Id,omitempty"`
	Name     string `json:"Name,omitempty"`

	// Redfish Power.v1_5_4.json
	Description string         `json:"Description,omitempty"`
	PowerCtl    []PowerControl `json:"PowerControl,omitempty"`
	PowerCtlCnt int            `json:"PowerControl@odata.count,omitempty"`

	// HpeServerAccPowerLimit.v1_0_0.HpeServerAccPowerLimit
	Error             *HpeError              `json:"error,omitempty"`
	ActualPowerLimits []HpeActualPowerLimits `json:"ActualPowerLimits,omitempty"`
	PowerLimitRanges  []HpePowerLimitRanges  `json:"PowerLimitRanges,omitempty"`
	PowerLimits       []HpePowerLimits       `json:"PowerLimits,omitempty"`

	// Redfish Control.v1_0_0.Control
	RFControl
}

// Structs used to [un]marshal HPE Redfish
// HpeServerAccPowerLimit.v1_0_0.HpeServerAccPowerLimit data
type HpeConfigurePowerLimit struct {
	PowerLimits []HpePowerLimits `json:"PowerLimits"`
}

type HpeActualPowerLimits struct {
	PowerLimitInWatts *int `json:"PowerLimitInWatts"`
	ZoneNumber        *int `json:"ZoneNumber"`
}

type HpePowerLimitRanges struct {
	MaximumPowerLimit *int `json:"MaximumPowerLimit"`
	MinimumPowerLimit *int `json:"MinimumPowerLimit"`
	ZoneNumber        *int `json:"ZoneNumber"`
}

type HpePowerLimits struct {
	PowerLimitInWatts *int `json:"PowerLimitInWatts"`
	ZoneNumber        *int `json:"ZoneNumber"`
}

type HpeError struct {
	Code         string                   `json:"code"`
	Message      string                   `json:"message"`
	ExtendedInfo []HpeMessageExtendedInfo `json:"@Message.ExtendedInfo"`
}

type HpeMessageExtendedInfo struct {
	MessageId string `json:"MessageId"`
}

// PowerControl struct used to unmarshal the Redfish Power.v1_5_4 data
type PowerControl struct {
	Oid                 string           `json:"@odata.id,omitempty"`
	Name                string           `json:"Name,omitempty"`
	OEM                 *PowerControlOEM `json:"Oem,omitempty"`
	PowerAllocatedWatts *int             `json:"PowerAllocatedWatts,omitempty"`
	PowerAvailableWatts *int             `json:"PowerAvailableWatts,omitempty"`
	PowerCapacityWatts  *int             `json:"PowerCapacityWatts,omitempty"`
	PowerConsumedWatts  *int             `json:"PowerConsumedWatts,omitempty"`
	PowerLimit          *PowerLimit      `json:"PowerLimit,omitempty"`
	PowerMetrics        *PowerMetric     `json:"PowerMetrics,omitempty"`
	PowerRequestedWatts *int             `json:"PowerRequestedWatts,omitempty"`
}

// StatusRF struct used to unmarshal health info from Redfish Control.v1_0_0
type StatusRF struct {
	Health       string `json:"Health,omitempty"`
	HealthRollUp string `json:"HealthRollUp,omitempty"`
	State        string `json:"State,omitempty"`
}

// Member array of RFControl structs for Deep Patch with Olympus hardare
type RFControlsDeep struct {
	Members []RFControl `json:"Members"`
}

// RFControl struct used to unmarshal the Redfish Control.v1_0_0.Control data
type RFControl struct {
	Oid                 string    `json:"@odata.id,omitempty"`
	ControlDelaySeconds *int      `json:"ControlDelaySeconds,omitempty"`
	ControlMode         string    `json:"ControlMode,omitempty"`
	ControlType         string    `json:"ControlType,omitempty"`
	Id                  string    `json:"Id,omitempty"`
	Name                string    `json:"Name,omitempty"`
	PhysicalContext     string    `json:"PhysicalContext,omitempty"`
	SetPoint            *int      `json:"SetPoint,omitempty"`
	SetPointUnits       string    `json:"SetPointUnits,omitempty"`
	SettingRangeMax     *int      `json:"SettingRangeMax,omitempty"`
	SettingRangeMin     *int      `json:"SettingRangeMin,omitempty"`
	Status              *StatusRF `json:"Status,omitempty"`
}

// PowerControlOEM contains a pointer to the OEM specific information
type PowerControlOEM struct {
	Cray *PowerControlOEMCray `json:"Cray,omitempty"`
}

// PowerControlOEMCray describes the Mountain specific power information
type PowerControlOEMCray struct {
	PowerAllocatedWatts   *int               `json:"PowerAllocatedWatts,omitempty"`
	PowerIdleWatts        *int               `json:"PowerIdleWatts,omitempty"`
	PowerLimit            *PowerLimitOEMCray `json:"PowerLimit,omitempty"`
	PowerFloorTargetWatts *int               `json:"PowerFloorTargetWatts,omitempty"`
	PowerResetWatts       *int               `json:"PowerResetWatts,omitempty"`
}

// PowerLimitOEMCray describes the power limit status and configuration for
// Mountain nodes
type PowerLimitOEMCray struct {
	Min    *int     `json:"Min,omitempty"`
	Max    *int     `json:"Max,omitempty"`
	Factor *float32 `json:"Factor,omitempty"`
}

// PowerLimit describes the power limit status and configuration for a
// compute module
type PowerLimit struct {
	CorrectionInMs *int   `json:"CorrectionInMs,omitempty"`
	LimitException string `json:"LimitException,omitempty"`
	LimitInWatts   *int   `json:"LimitInWatts"`
}

// PowerMetric describes the power readings for the compute module
type PowerMetric struct {
	AverageConsumedWatts *int `json:"AverageConsumedWatts,omitempty"`
	IntervalInMin        *int `json:"IntervalInMin,omitempty"`
	MaxConsumedWatts     *int `json:"MaxConsumedWatts,omitempty"`
	MinConsumedWatts     *int `json:"MinConsumedWatts,omitempty"`
}
