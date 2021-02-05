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
