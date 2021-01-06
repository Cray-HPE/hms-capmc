// Copyright 2019 Cray Inc. All Rights Reserved.
//
// Except as permitted by contract or express written permission of Cray Inc.,
// no part of this work or its content may be modified, used, reproduced or
// disclosed in any form. Modifications made without express permission of
// Cray Inc. may damage the system the software is installed within, may
// disqualify the user from receiving support from Cray Inc. under support or
// maintenance contracts, or require additional support services outside the
// scope of those contracts to repair the software or system.

package capmc

// Structures that enable marshalling and unmarshalling of the Power schema as
// defined by the Redfish Power.v1_5_4.json

// Power struct used to unmarshal the Redfish Power data that contains the power
// metrics, consumption, and limiting information
type Power struct {
	OContext    string `json:"@odata.context,omitempty"`
	Oetag       string `json:"@odata.etag,omitempty"`
	Oid         string `json:"@odata.id,omitempty"`
	Otype       string `json:"@odata.type,omitempty"`
	Description string `json:"Description,omitempty"`
	Id          string `json:"Id,omitempty"`
	Name        string `json:"Name,omitempty"`

	PowerCtl    []PowerControl `json:"PowerControl,omitempty"`
	PowerCtlCnt int            `json:"PowerControl@odata.count,omitempty"`
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
