// Copyright 2019 Cray Inc. All Rights Reserved.
//
// Except as permitted by contract or express written permission of Cray Inc.,
// no part of this work or its content may be modified, used, reproduced or
// disclosed in any form. Modifications made without express permission of
// Cray Inc. may damage the system the software is installed within, may
// disqualify the user from receiving support from Cray Inc. under support or
// maintenance contracts, or require additional support services outside the
// scope of those contracts to repair the software or system.
//
// This file contains the unit test for CAPMC
//

package main

import (
	"testing"
)

var allowableValues = [][]string{
	// #0
	// ResetType@Redfish.AllowableValues from Intel River
	// /redfish/v1/Systems/{SystemId}
	{"On", "ForceOff", "GracefulShutdown", "GracefulRestart", "ForceRestart", "Nmi"},

	// #1
	// ResetType@Redfish.AllowableValues from Intel River
	// /redfish/v1/Systems/{SystemId} NOTE from later Intel firmware
	// This maybe a good eventual test case for PushPowerButton logic
	{"PushPowerButton", "On", "GracefulShutdown", "ForceRestart", "Nmi", "ForceOn", "ForceOff"},

	// #2
	// ResetType@Redfish.AllowableValues from Intel River
	// /redfish/v1/Managers/BMC
	{"ForceRestart"},

	// #3
	// ResetType@Redfish.AllowableValues from Cray Mountain
	// Switch or Node /redfish/v1/Managers/BMC
	// NOTE the Cray weirdness using a non-Redfish enumerated ResetType
	{"ForceRestart", "ForceEraseNetworkReload"},

	// #4
	// ResetType@Redfish.AllowableValues from Cray Mountain
	// Switch /redfish/v1/Chassis/Enclosure
	// NOTE DTMF is adding "Off" at Cray's request
	{"ForceOff", "GracefulShutdown", "Off", "On", "GracefulRestart", "PowerCycle", "ForceOn", "ForceRestart"},

	// #5
	// ResetType@Redfish.AllowableValues from Cray Mountain Node
	// /redfish/v1/Systems/{SystemId}
	{"ForceOff", "Off", "On"},

	// #6
	// ResetType@Redfish.AllowableValues from Cray Mountain Node
	// /redfish/v1/Chassis/Enclosure
	// NOTE Yes it is really empty
	{},
}

func TestCmdToResetType(t *testing.T) {
	var svc CapmcD
	svc.config = loadConfig("")

	var tests = []struct {
		cmd string
		av  []string
		r   string
		e   bool
	}{
		// NOTE Tests assume default mappings/preferences
		{bmcCmdPowerOff, allowableValues[0], "GracefulShutdown", false},
		{bmcCmdPowerForceOff, allowableValues[0], "ForceOff", false},
		{bmcCmdPowerOn, allowableValues[0], "On", false},
		{bmcCmdPowerForceOn, allowableValues[0], "", true},
		{bmcCmdPowerOn, allowableValues[1], "On", false},
		{bmcCmdPowerForceRestart, allowableValues[2], "FroceRestart", false},
		{bmcCmdPowerOff, allowableValues[2], "", true},
		{bmcCmdPowerOn, allowableValues[2], "", true},
		{bmcCmdPowerOff, allowableValues[3], "", true},
		{bmcCmdPowerOn, allowableValues[3], "", true},
		{bmcCmdPowerForceRestart, allowableValues[3], "FroceRestart", false},
		{bmcCmdPowerOff, allowableValues[4], "GracefulShutdown", false},
		{bmcCmdPowerOn, allowableValues[4], "GracefulShutdown", false},
		{bmcCmdPowerOff, allowableValues[5], "Off", false},
		{bmcCmdPowerOn, allowableValues[5], "On", false},
		{bmcCmdPowerOff, allowableValues[6], "", true},
		{bmcCmdPowerOn, allowableValues[6], "", true},
	}

	for n, test := range tests {
		rt, err := svc.cmdToResetType(test.cmd, test.av)

		if !test.e && err != nil {
			t.Errorf("TestCmdToResetType Test Case %d: FAIL: Expected %s but instead got %s: %s", n, test.r, rt, err)
		}
		if test.e && err == nil {
			t.Errorf("TestCmdToResetType Test Case %d: FAIL: Error Expected but instead got %s", n, rt)
		}
	}
}

func TestCmdBlockRole(t *testing.T) {
	var svc CapmcD
	svc.config = loadConfig("")

	tests := []struct {
		cmd string
		r   []string
		e   bool
	}{
		{bmcCmdNMI, []string{"Storage", "System"}, false},
		{bmcCmdPowerForceOff, []string{"Storage", "System"}, false},
		{bmcCmdPowerForceOn, []string{"Storage", "System"}, false},
		{bmcCmdPowerOff, []string{"Storage", "System"}, false},
		{bmcCmdPowerOn, []string{"Storage", "System"}, false},
		{"HCF", []string{}, true},
	}

	for n, test := range tests {
		rt, err := svc.cmdBlockRole(test.cmd)

		if !test.e && err != nil {
			t.Errorf("TestCmdBlockRole Test Case %d: FAIL: Expected %s but instead got %s: %s", n, test.r, rt, err)
		}
		if test.e && err == nil {
			t.Errorf("TestGetCmdPowerSequence Test Case %d: FAIL: Error Expected but instead got %s", n, rt)
		}
	}
}
