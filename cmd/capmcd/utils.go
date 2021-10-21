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
	"fmt"
	"strings"
)

// stringInSlice checks the slice a contains the string s
// TODO - not a core CAPMC function, move
func stringInSlice(s string, a []string) bool {

	for _, n := range a {
		if s == n {
			return true
		}
	}

	return false
}

// stringSliceMap returns a copy of the string slice a with all its strings
// modified according to the mapping function f.
// TODO - not a core CAPMC function, move
func stringSliceMap(a []string, f func(string) string) []string {

	na := make([]string, len(a))
	for i, s := range a {
		na[i] = f(s)
	}

	return na
}

// getForceOption - Returns the force version of the command sent in.
func (d *CapmcD) getForceOption(cmd string) (string, error) {
	var ncmd string
	switch cmd {
	case bmcCmdPowerOn:
		ncmd = bmcCmdPowerForceOn
	case bmcCmdPowerOff:
		ncmd = bmcCmdPowerForceOff
	case bmcCmdPowerRestart:
		ncmd = bmcCmdPowerForceRestart
	default:
		return "", fmt.Errorf("Cannot force the %s operation", cmd)
	}
	return ncmd, nil
}

func cmdBlockRole(cmd string) ([]string, error) {
	return svc.cmdBlockRole(cmd)
}

func (d *CapmcD) cmdBlockRole(cmd string) ([]string, error) {

	pc, ok := d.config.PowerControls[cmd]
	if !ok {
		return []string{}, fmt.Errorf("no power controls for %s operation", cmd)
	}

	return pc.BlockRole, nil
}

func cmdToResetType(cmd string, allowable []string) (string, error) {
	return svc.cmdToResetType(cmd, allowable)
}

func (d *CapmcD) cmdToResetType(cmd string, allowable []string) (string, error) {

	pc, ok := d.config.PowerControls[cmd]
	if !ok {
		return "", fmt.Errorf("no power controls for %s operation", cmd)
	}

	// search
	for _, rt := range pc.ResetType {
		// found?
		if stringInSlice(rt, allowable) {
			return rt, nil
		}
	}

	return "", fmt.Errorf("no supported ResetType for %s operation", cmd)
}

// CompIdsToNids returns a list of Node IDs (NIDs) for the list of
// Component IDs (xnames) based upon the supplied component id to nid map.
// Assumes the mapping is complete.
func compIdsToNids(cids []string, cidToNid map[string]int) ([]int, error) {
	var (
		err  error
		nids []int
		bad  []string
	)

	for _, cid := range cids {
		nid, ok := cidToNid[cid]
		if !ok {
			bad = append(bad, cid)
			continue
		}
		nids = append(nids, nid)
	}

	if len(bad) > 0 {
		err = &InvalidCompIDsError{"Component IDs without NIDs", bad}
	}

	return nids, err
}

func isHpeServer(ni *NodeInfo) bool {
	if strings.Contains(ni.RfPowerURL, "Chassis/1/Power") {
		return true
	}

	return false
}

func isHpeApollo6500(ni *NodeInfo) bool {
	if strings.Contains(ni.RfPowerURL, "AccPowerService/PowerLimit") {
		return true
	}

	return false
}

func isControlPoint(ni *NodeInfo) bool {
	if strings.Contains(ni.RfPowerURL, "Controls") {
		return true
	}

	return false
}

func isGigabyte(ni *NodeInfo) bool {
	if strings.Contains(ni.RfPowerURL, "Chassis/Self") {
		return true
	}

	return false
}
