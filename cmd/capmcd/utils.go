// Copyright 2017-2020 Hewlett Packard Enterprise Development LP
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
// TODO - Not all the functions here are core CAPMC "business" logic. Those
//        that are not should be moved to a new hms-capmc/pkg or hsm-common.

package main

import "fmt"

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
