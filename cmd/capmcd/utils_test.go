// Copyright 2017-2019 Cray Inc. All Rights Reserved.
//
// Except as permitted by contract or express written permission of Cray Inc.,
// no part of this work or its content may be modified, used, reproduced or
// disclosed in any form. Modifications made without express permission of
// Cray Inc. may damage the system the software is installed within, may
// disqualify the user from receiving support from Cray Inc. under support or
// maintenance contracts, or require additional support services outside the
// scope of those contracts to repair the software or system.
//
// A grab bag of utility functions.
//
// TODO The functions here are not core CAPMC "business" logic. They should
//      be moved to a new hms-capmc/pkg or hsm-common.

package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestStringInSlice(t *testing.T) {
	var tests = []struct {
		s string
		a []string
		r bool
	}{
		{
			"Yes",
			[]string{"Maybe", "No", "Yes"},
			true,
		},
		{
			"Five",
			[]string{"Two", "Four", "Six", "Eight"},
			false,
		},
	}

	for n, test := range tests {
		r := stringInSlice(test.s, test.a)
		if r != test.r {
			t.Errorf("TestContains Test Case %d: FAIL: Expected %v but got %v", n, test.r, r)
		}
	}
}

func TestStringSliceMap(t *testing.T) {
	var tests = []struct {
		a []string
		f func(string) string
		r []string
	}{
		{
			a: []string{"apple", "boroughs", "cray"},
			f: strings.ToUpper,
			r: []string{"APPLE", "BOROUGHS", "CRAY"},
		},
	}

	for n, test := range tests {
		r := stringSliceMap(test.a, test.f)
		if !reflect.DeepEqual(r, test.r) {
			t.Errorf("StringSliceMap Test Case %d: FAIL: Expected %v but got %v", n, test.r, r)
		}
	}
}

func TestCompIdsToNids(t *testing.T) {
	var tests = []struct {
		cids  []string
		cnMap map[string]int
		nids  []int
		errOK bool
	}{
		{
			cids: []string{"x0c0s0b0n0", "x0c0s0b0n1", "x0c0s0b1n0", "x0c0s0b1n1"},
			cnMap: map[string]int{
				"x0c0s0b0n0": 1,
				"x0c0s0b0n1": 2,
				"x0c0s0b1n0": 3,
				"x0c0s0b1n1": 4,
			},
			nids:  []int{1, 2, 3, 4},
			errOK: false,
		}, {
			cids: []string{"x0c0s0b0n0", "x0c0s0b0n1", "x0c0s0b1n0", "x0c0s0b1n1"},
			cnMap: map[string]int{
				"x0c0s0b0n0": 1,
				"x0c0s0b0n1": 2,
				"x0c0s0b1n0": 3,
			},
			nids:  []int{1, 2, 3},
			errOK: true,
		},
	}

	for n, test := range tests {
		r, err := compIdsToNids(test.cids, test.cnMap)
		if err != nil && !test.errOK {
			t.Errorf("compIdsToNids Test Case %d: FAIL: Unexpected error: %s", n, err.Error())
		}

		if !reflect.DeepEqual(r, test.nids) {
			t.Errorf("compIdsToNids Test Case %d: FAIL: Expected %v but got %v", n, test.nids, r)
		}

	}
}
