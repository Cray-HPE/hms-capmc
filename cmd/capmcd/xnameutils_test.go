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
