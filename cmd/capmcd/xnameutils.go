// Copyright 2019 Cray Inc. All Rights Reserved.
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
