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
// This file contains the utilities for sorting xname referenced
// components for capmcd
//

package main

import (
	"regexp"
	"strconv"
)

var splitRE = regexp.MustCompile("^([[:alpha:]])([[:digit:]]+)(.*)$")

// The xnameSlice struct is used to manage sorting xnames.
// It is used as a concrete instance of sort.Interface
type xnameSlice struct {
	slice *[]string
}

func (xs xnameSlice) Len() int {
	return len(*xs.slice)
}

// The Less function for xname comparison looks at each subsequent sequence of
// a letter and a set of digits.  If the letters differ, the return value
// reflects their comparison.  If they are the same, then the string of digits
// are compared numerically.  If they differ, the return value reflects this
// comparison.  Otherwise, the next set of a letter followed by a string of
// digits is compared.  This continues until the end of one of the xnames is
// reached.  For two identical items, this function will always return false.
func (xs xnameSlice) Less(i, j int) bool {
	split := func(str string) (string, string, string) {
		r := splitRE.FindSubmatch([]byte(str))
		l := len(r)
		var ret [3]string
		if l > 4 {
			l = 4
		}
		for i := 1; i < l; i++ {
			ret[i-1] = string(r[i])
		}
		return ret[0], ret[1], ret[2]
	}
	var ret int
	x1 := (*xs.slice)[i]
	x2 := (*xs.slice)[j]
	for ret == 0 && (x1 != "" || x2 != "") {
		letter1, digits1, remaining1 := split(x1)
		letter2, digits2, remaining2 := split(x2)
		if letter1 < letter2 {
			ret = -1
		} else if letter1 > letter2 {
			ret = 1
		} else if digits1 != digits2 {
			nbr1, _ := strconv.ParseInt(digits1, 10, 32)
			nbr2, _ := strconv.ParseInt(digits2, 10, 32)
			ret = int(nbr1) - int(nbr2)
		}
		x1 = remaining1
		x2 = remaining2
	}
	return ret < 0
}

func (xs xnameSlice) Swap(i, j int) {
	(*xs.slice)[i], (*xs.slice)[j] = (*xs.slice)[j], (*xs.slice)[i]
}
