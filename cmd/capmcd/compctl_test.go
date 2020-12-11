// Copyright 2019-2020 Hewlett Packard Enterprise Development LP
//

package main

import (
	"reflect"
	"testing"

	"stash.us.cray.com/HMS/hms-capmc/internal/capmc"
)

func TestCmdCompPowerSeq(t *testing.T) {
	var svc CapmcD
	svc.config = loadConfig("")

	var tests = []struct {
		cmd string
		r   []string
		e   bool
	}{
		{bmcCmdNMI, []string{"Node"}, false},
		{bmcCmdPowerForceOff, []string{"Node", "ComputeModule", "HSNBoard", "RouterModule", "Chassis"}, false},
		{bmcCmdPowerForceOn, []string{"Chassis", "RouterModule", "HSNBoard", "ComputeModule", "Node"}, false},
		{bmcCmdPowerOff, []string{"Node", "ComputeModule", "HSNBoard", "RouterModule", "Chassis"}, false},
		{bmcCmdPowerOn, []string{"Chassis", "RouterModule", "HSNBoard", "ComputeModule", "Node"}, false},
		{"HCF", []string{}, true},
	}

	for n, test := range tests {
		rt, err := svc.cmdCompPowerSeq(test.cmd)

		if !test.e && err != nil {
			t.Errorf("TestGetCmdPowerSequence Test Case %d: FAIL: Expected %s but instead got %s: %s", n, test.r, rt, err)
		}
		if test.e && err == nil {
			t.Errorf("TestGetCmdPowerSequence Test Case %d: FAIL: Error Expected but instead got %s", n, rt)
		}
	}
}

func TestDoRemoveComp(t *testing.T) {
	comp1 := &NodeInfo{Hostname: "x0c0s0b0n0"}
	comp2 := &NodeInfo{Hostname: "x0c0s0b0n1"}
	comp3 := &NodeInfo{Hostname: "x0c0s0b0n2"}
	tests := []struct {
		cmapIn     map[string]map[string][]*NodeInfo
		xname      string
		ctype      string
		actionList []string
		cmapOut    map[string]map[string][]*NodeInfo
	}{{
		// Test removing a component for a single action
		cmapIn: map[string]map[string][]*NodeInfo{
			bmcCmdPowerOff: {
				"Node": {comp1, comp2, comp3},
			},
			bmcCmdPowerOn: {
				"Node": {comp1, comp2, comp3},
			},
		},
		xname:      comp2.Hostname,
		ctype:      "Node",
		actionList: []string{bmcCmdPowerOn},
		cmapOut: map[string]map[string][]*NodeInfo{
			bmcCmdPowerOff: {
				"Node": {comp1, comp2, comp3},
			},
			bmcCmdPowerOn: {
				"Node": {comp1, comp3},
			},
		},
	}, {
		// Test removing a component for a multiple actions
		cmapIn: map[string]map[string][]*NodeInfo{
			bmcCmdPowerOff: {
				"Node": {comp1, comp2, comp3},
			},
			bmcCmdPowerOn: {
				"Node": {comp1, comp2, comp3},
			},
		},
		xname:      comp2.Hostname,
		ctype:      "Node",
		actionList: []string{bmcCmdPowerOff, bmcCmdPowerOn},
		cmapOut: map[string]map[string][]*NodeInfo{
			bmcCmdPowerOff: {
				"Node": {comp1, comp3},
			},
			bmcCmdPowerOn: {
				"Node": {comp1, comp3},
			},
		},
	}, {
		// Test removing a component for a non-existing type
		cmapIn: map[string]map[string][]*NodeInfo{
			bmcCmdPowerOff: {
				"Node": {comp1, comp2, comp3},
			},
			bmcCmdPowerOn: {
				"Node": {comp1, comp2, comp3},
			},
		},
		xname:      comp2.Hostname,
		ctype:      "NodeBMC",
		actionList: []string{bmcCmdPowerOn},
		cmapOut: map[string]map[string][]*NodeInfo{
			bmcCmdPowerOff: {
				"Node": {comp1, comp2, comp3},
			},
			bmcCmdPowerOn: {
				"Node": {comp1, comp2, comp3},
			},
		},
	}, {
		// Test removing a component for a non-existing action
		cmapIn: map[string]map[string][]*NodeInfo{
			bmcCmdPowerOff: {
				"Node": {comp1, comp2, comp3},
			},
			bmcCmdPowerOn: {
				"Node": {comp1, comp2, comp3},
			},
		},
		xname:      comp2.Hostname,
		ctype:      "Node",
		actionList: []string{bmcCmdPowerRestart},
		cmapOut: map[string]map[string][]*NodeInfo{
			bmcCmdPowerOff: {
				"Node": {comp1, comp2, comp3},
			},
			bmcCmdPowerOn: {
				"Node": {comp1, comp2, comp3},
			},
		},
	}}
	for n, test := range tests {
		doRemoveComp(test.cmapIn, test.xname, test.ctype, test.actionList)
		if !reflect.DeepEqual(test.cmapOut, test.cmapIn) {
			t.Errorf("TestDoRemoveComp Test Case %d: FAIL: Expected %v but got %v", n, test.cmapOut, test.cmapIn)
		}
	}
}

func TestDoSortActionMap(t *testing.T) {
	var d CapmcD
	var emptyData capmc.XnameControlResponse
	comp1 := &NodeInfo{Hostname: "x0c0s0b0n0", Type: "Node", RfResetTypes: allowableValues[0]}
	comp2 := &NodeInfo{Hostname: "x0c0", Type: "Chassis", RfResetTypes: allowableValues[4]}
	comp3 := &NodeInfo{Hostname: "x0c0s1b0n0", Type: "Node", RfResetTypes: allowableValues[5]}
	comp4 := &NodeInfo{Hostname: "x0c0s1r0b0", Type: "RouterBMC", RfResetTypes: allowableValues[3]}
	comp1Err1 := capmc.MakeXnameError(comp1.Hostname, -1, "no power controls for badCommand operation")
	comp2Err1 := capmc.MakeXnameError(comp2.Hostname, -1, "Skipping x0c0: Type, 'Chassis', not defined in power sequence for 'Restart'")
	comp3Err1 := capmc.MakeXnameError(comp3.Hostname, -1, " : no supported ResetType for Restart operation")
	comp3Err2 := capmc.MakeXnameError(comp3.Hostname, 0, "Ignored:  : no supported ResetType for Restart operation")
	comp4Err1 := capmc.MakeXnameError(comp4.Hostname, -1, "Skipping x0c0s1r0b0: Type, 'RouterBMC', not defined in power sequence for 'On'")
	emptyData.Xnames = make([]*capmc.XnameControlErr, 0, 1)
	tests := []struct {
		nlIn     []*NodeInfo
		cmdIn    string
		oua      string
		cmapOut  map[string]map[string][]*NodeInfo
		dataOut  capmc.XnameControlResponse
		countOut int
		badOut   int
	}{{
		// 0: Test sorting components into the action map
		nlIn:  []*NodeInfo{comp1, comp2, comp3},
		cmdIn: bmcCmdPowerOn,
		oua:   actionIgnore,
		cmapOut: map[string]map[string][]*NodeInfo{
			bmcCmdPowerOn: {
				"Node":    {comp1, comp3},
				"Chassis": {comp2},
			},
		},
		dataOut:  emptyData,
		countOut: 3,
		badOut:   0,
	}, {
		// 1: Test sorting components into the action map with simulation
		nlIn:  []*NodeInfo{comp1, comp3},
		cmdIn: bmcCmdPowerRestart,
		oua:   actionSimulate,
		cmapOut: map[string]map[string][]*NodeInfo{
			bmcCmdPowerOff: {
				"Node": {comp3},
			},
			bmcCmdPowerRestart: {
				"Node": {comp1},
			},
			bmcCmdPowerOn: {
				"Node": {comp3},
			},
		},
		dataOut:  emptyData,
		countOut: 2,
		badOut:   0,
	}, {
		// 2: Test sorting components into the action map with simulation (force)
		nlIn:  []*NodeInfo{comp1, comp3},
		cmdIn: bmcCmdPowerForceRestart,
		oua:   actionSimulate,
		cmapOut: map[string]map[string][]*NodeInfo{
			bmcCmdPowerForceOff: {
				"Node": {comp3},
			},
			bmcCmdPowerForceRestart: {
				"Node": {comp1},
			},
			bmcCmdPowerForceOn: {
				"Node": {comp3},
			},
		},
		dataOut:  emptyData,
		countOut: 2,
		badOut:   0,
	}, {
		// 3: Test sorting components into the action map with errors
		nlIn:  []*NodeInfo{comp1, comp3},
		cmdIn: bmcCmdPowerRestart,
		oua:   actionError,
		cmapOut: map[string]map[string][]*NodeInfo{
			bmcCmdPowerRestart: {
				"Node": {comp1},
			},
		},
		dataOut:  capmc.XnameControlResponse{Xnames: []*capmc.XnameControlErr{comp3Err1}},
		countOut: 2,
		badOut:   1,
	}, {
		// 4: Test sorting components into the action map with ignored errors
		nlIn:  []*NodeInfo{comp1, comp3},
		cmdIn: bmcCmdPowerRestart,
		oua:   actionIgnore,
		cmapOut: map[string]map[string][]*NodeInfo{
			bmcCmdPowerRestart: {
				"Node": {comp1},
			},
		},
		dataOut:  capmc.XnameControlResponse{Xnames: []*capmc.XnameControlErr{comp3Err2}},
		countOut: 2,
		badOut:   0,
	}, {
		// 5: Test sorting components into the action map with unset OnUnsupportedAction
		nlIn:  []*NodeInfo{comp1, comp3},
		cmdIn: bmcCmdPowerRestart,
		oua:   "",
		cmapOut: map[string]map[string][]*NodeInfo{
			bmcCmdPowerRestart: {
				"Node": {comp1, comp3},
			},
		},
		dataOut:  emptyData,
		countOut: 2,
		badOut:   0,
	}, {
		// 6: Test sorting components into the action map with unsupported components Restart
		nlIn:  []*NodeInfo{comp1, comp2},
		cmdIn: bmcCmdPowerRestart,
		oua:   actionError,
		cmapOut: map[string]map[string][]*NodeInfo{
			bmcCmdPowerRestart: {
				"Node": {comp1},
			},
		},
		dataOut:  capmc.XnameControlResponse{Xnames: []*capmc.XnameControlErr{comp2Err1}},
		countOut: 1,
		badOut:   1,
	}, {
		// 7: Test sorting components with an empty nodelist
		nlIn:  []*NodeInfo{},
		cmdIn: bmcCmdPowerRestart,
		oua:   actionIgnore,
		cmapOut: map[string]map[string][]*NodeInfo{
			bmcCmdPowerRestart: {},
		},
		dataOut:  emptyData,
		countOut: 0,
		badOut:   0,
	}, {
		// 8: Test sorting components into the action map with unsupported components On
		nlIn:  []*NodeInfo{comp1, comp4},
		cmdIn: bmcCmdPowerOn,
		oua:   actionError,
		cmapOut: map[string]map[string][]*NodeInfo{
			bmcCmdPowerOn: {
				"Node": {comp1},
			},
		},
		dataOut:  capmc.XnameControlResponse{Xnames: []*capmc.XnameControlErr{comp4Err1}},
		countOut: 1,
		badOut:   1,
	}, {
		// 9: Test sorting components into the action map with unsupported command
		nlIn:  []*NodeInfo{comp1},
		cmdIn: "badCommand",
		oua:   actionError,
		cmapOut: map[string]map[string][]*NodeInfo{
			"badCommand": {},
		},
		dataOut:  capmc.XnameControlResponse{Xnames: []*capmc.XnameControlErr{comp1Err1}},
		countOut: 0,
		badOut:   1,
	}}
	d.config = new(Config)
	d.config.PowerControls = defaultPowerControl
	for n, test := range tests {
		var data capmc.XnameControlResponse
		data.Xnames = make([]*capmc.XnameControlErr, 0, 1)
		d.OnUnsupportedAction = test.oua
		cmap := make(map[string]map[string][]*NodeInfo)
		cmap[test.cmdIn] = make(map[string][]*NodeInfo)
		count, bad := d.doSortActionMap(test.nlIn, cmap, test.cmdIn, &data)
		if !reflect.DeepEqual(test.cmapOut, cmap) {
			t.Errorf("TestDoSortActionMap Test Case %d: cmap FAIL: Expected %v but got %v", n, test.cmapOut, cmap)
		}
		if !reflect.DeepEqual(test.dataOut, data) {
			t.Errorf("TestDoSortActionMap Test Case %d: data FAIL: Expected %v but got %v", n, test.dataOut, data)
		}
		if test.countOut != count {
			t.Errorf("TestDoSortActionMap Test Case %d: count FAIL: Expected %v but got %v", n, test.countOut, count)
		}
		if test.badOut != bad {
			t.Errorf("TestDoSortActionMap Test Case %d: bad count FAIL: Expected %v but got %v", n, test.badOut, bad)
		}
	}
}
