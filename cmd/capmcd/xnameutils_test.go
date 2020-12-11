// Copyright 2019 Cray Inc. All Rights Reserved.
//
// This file contains the unit tests for CAPMC's xname utils
//

package main

import (
	"reflect"
	"testing"
)

func TestXnameSliceLen(t *testing.T) {
	var tests = []struct {
		xs xnameSlice
		r  int
	}{{
		// Length of populated xnameSlice.
		xs: xnameSlice{slice: &[]string{"x0c0s0b0n0", "x0c0s1b0n0", "x0c0s1"}},
		r:  3,
	}, {
		// Length of empty xnameSlice.
		xs: xnameSlice{slice: &[]string{}},
		r:  0,
	}}

	for n, test := range tests {
		rt := test.xs.Len()

		if test.r != rt {
			t.Errorf("TestXnameSliceLen Test Case %d: FAIL: Expected %v but instead got %v", n, test.r, rt)
		}
	}
}

func TestXnameSliceLess(t *testing.T) {
	xs := xnameSlice{slice: &[]string{"x0c0s0b0n0", "x0c0s1b0n0", "x0c0s1"}}
	var tests = []struct {
		i int
		j int
		r bool
	}{{
		// Comparison of lesser xname to greater xname of the same length.
		i: 0,
		j: 1,
		r: true,
	}, {
		// Comparison of greater xname to lesser xname of the same length.
		i: 1,
		j: 0,
		r: false,
	}, {
		// Comparison of the same.
		i: 0,
		j: 0,
		r: false,
	}, {
		// Comparison of lesser xname to greater xname of differing lengths.
		i: 0,
		j: 2,
		r: true,
	}, {
		// Comparison of similar xnames of differing lengths.
		i: 1,
		j: 2,
		r: false,
	}}

	for n, test := range tests {
		rt := xs.Less(test.i, test.j)

		if test.r != rt {
			t.Errorf("TestXnameSliceLess Test Case %d: FAIL: Expected %v but instead got %v", n, test.r, rt)
		}
	}
}

func TestXnameSliceSwap(t *testing.T) {
	var tests = []struct {
		xs xnameSlice
		i  int
		j  int
		r  xnameSlice
	}{{
		// Swap non-adjacent entries.
		xs: xnameSlice{slice: &[]string{"x0c0s1", "x0c0s1b0n0", "x0c0s0b0n0"}},
		i:  0,
		j:  2,
		r:  xnameSlice{slice: &[]string{"x0c0s0b0n0", "x0c0s1b0n0", "x0c0s1"}},
	}, {
		// Swap adjacent entries.
		xs: xnameSlice{slice: &[]string{"x0c0s0b0n0", "x0c0s1b0n0", "x0c0s1"}},
		i:  1,
		j:  2,
		r:  xnameSlice{slice: &[]string{"x0c0s0b0n0", "x0c0s1", "x0c0s1b0n0"}},
	}}

	for n, test := range tests {
		test.xs.Swap(test.i, test.j)

		if !reflect.DeepEqual(test.r, test.xs) {
			t.Errorf("TestXnameSliceSwap Test Case %d: FAIL: Expected %v but instead got %v", n, test.r, test.xs)
		}
	}
}
