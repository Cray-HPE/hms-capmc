// MIT License
//
// (C) Copyright [2019, 2021] Hewlett Packard Enterprise Development LP
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
