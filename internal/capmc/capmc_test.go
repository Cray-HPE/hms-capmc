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

import (
	"testing"
)

// TestNodeStatusFilterParse
func TestNodeStatusFilterParse(t *testing.T) {
	const NodeFilterShowAllBit = FilterShowAllBit &^ (FilterShowEmptyBit | FilterShowPopulatedBit | FilterShowUndefinedBit | FilterShowUnknownBit)
	var tests = []struct {
		filterStr  string
		filterBMap uint
		errorOK    bool
	}{
		{"", NodeFilterShowAllBit, false},
		{FilterShowAll, NodeFilterShowAllBit, false},
		{FilterShowAlert, FilterShowAlertBit, false},
		{FilterShowDiag, FilterShowDiagBit, false},
		{FilterShowDisabled, FilterShowDisabledBit, false},
		{FilterShowHalt, FilterShowHaltBit, false},
		{FilterShowOff, FilterShowOffBit, false},
		{FilterShowOn, FilterShowOnBit, false},
		{FilterShowReady, FilterShowReadyBit, false},
		{FilterShowReserved, FilterShowReservedBit, false},
		{FilterShowStandby, FilterShowStandbyBit, false},
		{FilterShowWarning, FilterShowWarningBit, false},
		{FilterShowEmpty, 0, true},
		{"RandomJunk", 0, true},
	}

	fmt := "FAIL: Expected %v but got %v"

	for n, test := range tests {
		t.Run(test.filterStr, func(t *testing.T) {
			bMap, err := NodeStatusFilterParse(test.filterStr)
			if err != nil {
				if !test.errorOK {
					t.Errorf(fmt, n, test.filterBMap, bMap)
				}
			}

			if bMap != test.filterBMap {
				t.Errorf(fmt, n, test.filterBMap, bMap)
			}
		})
	}
}

// TestStatusFilterParse
func TestStatusFilterParse(t *testing.T) {
	var tests = []struct {
		filterStr  string
		filterBMap uint
		errorOK    bool
	}{
		{"", FilterShowAllBit, false},
		{FilterShowAll, FilterShowAllBit, false},
		{FilterShowAlert, FilterShowAlertBit, false},
		{FilterShowDiag, FilterShowDiagBit, false},
		{FilterShowDisabled, FilterShowDisabledBit, false},
		{FilterShowEmpty, FilterShowEmptyBit, false},
		{FilterShowHalt, FilterShowHaltBit, false},
		{FilterShowOff, FilterShowOffBit, false},
		{FilterShowOn, FilterShowOnBit, false},
		{FilterShowPopulated, FilterShowPopulatedBit, false},
		{FilterShowReady, FilterShowReadyBit, false},
		{FilterShowReserved, FilterShowReservedBit, false},
		{FilterShowStandby, FilterShowStandbyBit, false},
		{FilterShowUndefined, FilterShowUndefinedBit, false},
		{FilterShowUnknown, FilterShowUnknownBit, false},
		{FilterShowWarning, FilterShowWarningBit, false},
		{"RandomJunk", 0, true},
	}

	fmt := "FAIL: Expected %v but got %v"

	for _, test := range tests {
		t.Run(test.filterStr, func(t *testing.T) {
			bMap, err := StatusFilterParse(test.filterStr)
			if err != nil {
				if !test.errorOK {
					t.Errorf("FAIL: Unexpected error: %s", err)
				}
			} else {
				if bMap != test.filterBMap {
					t.Errorf(fmt, test.filterBMap, bMap)
				}
			}
		})
	}
}

// TestMakeXnameError
// TODO

// TestParseNidlist
// TODO
