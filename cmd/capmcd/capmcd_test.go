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
