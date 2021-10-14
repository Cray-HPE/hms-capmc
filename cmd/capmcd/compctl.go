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
	"log"
	"sort"
	"time"

	"github.com/Cray-HPE/hms-capmc/internal/capmc"

	rf "github.com/Cray-HPE/hms-smd/pkg/redfish"
)

func cmdCompPowerSeq(cmd string) ([]string, error) {
	return svc.cmdCompPowerSeq(cmd)
}

func (d *CapmcD) cmdCompPowerSeq(cmd string) ([]string, error) {

	pc, ok := d.config.PowerControls[cmd]
	if !ok {
		return []string{}, fmt.Errorf("no power controls for %s operation", cmd)
	}

	return pc.CompSeq, nil
}

// hasCompPowerSupport - Returns true if the component type exists in the power
// sequence configuration for the specified command type.
func (d *CapmcD) hasCompPowerSupport(cmd, ctype string) (bool, error) {
	types, err := d.cmdCompPowerSeq(cmd)
	if err != nil {
		return false, err
	}
	return stringInSlice(ctype, types), nil
}

// Removes a component from the component map for the specified actions
func doRemoveComp(
	cmap map[string]map[string][]*NodeInfo,
	xname string,
	ctype string,
	actionList []string) {

	for _, action := range actionList {
		if _, ok := cmap[action]; !ok {
			continue
		}
		for idx, comp := range cmap[action][ctype] {
			if xname == comp.Hostname {
				// Remove the component from the future action
				cmap[action][ctype] = append(cmap[action][ctype][:idx], cmap[action][ctype][idx+1:]...)
				break
			}
		}
	}
}

// Takes a list of components and inserts them into a map sorted by type and
// action type. This returns the number of unique components inserted into the
// map and the number of components with errors.
// NOTE: This modifies the specified 'data' and 'cmap' args.
func (d *CapmcD) doSortActionMap(
	nl []*NodeInfo,
	cmap map[string]map[string][]*NodeInfo,
	cmd string,
	data *capmc.XnameControlResponse) (int, int) {

	var (
		failures  int
		totalWait int
	)
	// Simulate restart by doing an 'Off' then 'On'.
	simulate := map[string][]string{
		bmcCmdPowerRestart:      {bmcCmdPowerOff, bmcCmdPowerOn},
		bmcCmdPowerForceRestart: {bmcCmdPowerForceOff, bmcCmdPowerForceOn},
	}
	for _, c := range nl {
		// Check if the component type is supported for this power action.
		supported, err := d.hasCompPowerSupport(cmd, c.Type)
		if err != nil {
			failures++
			msg := fmt.Sprintf("%s", err)
			log.Printf("Error: %s.", msg)
			xnameErr := capmc.MakeXnameError(c.Hostname, -1, msg)
			data.Xnames = append(data.Xnames, xnameErr)
			continue
		}
		if !supported {
			// Skip components not in the power action sequencing list.
			msg := fmt.Sprintf("Skipping %s: Type, '%s', not defined in power sequence for '%s'", c.Hostname, c.Type, cmd)
			log.Printf("Info: %s.", msg)
			failures++
			xnameErr := capmc.MakeXnameError(c.Hostname, -1, msg)
			data.Xnames = append(data.Xnames, xnameErr)
			continue
		}
		// Set totalWait here so it will count of all components only once in
		// the case of more than one operation being done on a component.
		totalWait++
		if cmd == bmcCmdPowerRestart || cmd == bmcCmdPowerForceRestart {
			// Checking for an error from cmdToResetType() will tell us
			// if the component has the actions we need. If not, we check
			// OnUnsupportedAction for what we should do next.
			_, err := d.cmdToResetType(cmd, c.RfResetTypes)
			if err != nil {
				switch d.OnUnsupportedAction {
				case actionSimulate:
					for _, simAction := range simulate[cmd] {
						if _, ok := cmap[simAction]; !ok {
							cmap[simAction] = make(map[string][]*NodeInfo)
						}
						cmap[simAction][c.Type] = append(cmap[simAction][c.Type], c)
					}
					//TODO: Find dependant components and add to cmap when we
					// add reinit support for more than nodes.
				case actionError:
					// Fail the operation.
					failures++
					msg := fmt.Sprintf("%s %s: %s", c.BmcType, c.BmcFQDN, err)
					xnameErr := capmc.MakeXnameError(c.Hostname, -1, msg)
					data.Xnames = append(data.Xnames, xnameErr)
				case actionIgnore:
					// Leave out of the list but still report as ignored. Don't
					// increment 'failures' because these are "report only".
					msg := fmt.Sprintf("Ignored: %s %s: %s", c.BmcType, c.BmcFQDN, err)
					xnameErr := capmc.MakeXnameError(c.Hostname, 0, msg)
					data.Xnames = append(data.Xnames, xnameErr)
				default:
					// Leave in the list. This will just let doBmcCall() handle
					// the error as normal.
					cmap[cmd][c.Type] = append(cmap[cmd][c.Type], c)
				}
			} else {
				cmap[cmd][c.Type] = append(cmap[cmd][c.Type], c)
			}
		} else {
			cmap[cmd][c.Type] = append(cmap[cmd][c.Type], c)
		}
	}
	return totalWait, failures
}

func (d *CapmcD) waitForOff(ni *NodeInfo) bmcPowerRc {
	var res = bmcPowerRc{ni: ni, rc: -1, state: "Unknown"}
	var nl []*NodeInfo
	var retries int

	nl = append(nl, ni)

	offRetries := d.config.CapmcConf.WaitForOffRetries
	offSleep := d.config.CapmcConf.WaitForOffSleep

	for retries = 0; retries < offRetries; retries++ {
		rsp := d.doCompStatus(nl, bmcCmdPowerStatus, capmc.FilterShowOffBit)
		if len(rsp.Off) > 0 {
			res.rc = 0
			res.state = "Off"
			break
		} else {
			time.Sleep(time.Duration(offSleep) * time.Second)
		}
	}

	// One last attempt to get the power state after the last sleep
	if retries == offRetries {
		rsp := d.doCompStatus(nl, bmcCmdPowerStatus, capmc.FilterShowOffBit)
		if len(rsp.Off) > 0 {
			res.rc = 0
			res.state = "Off"
		} else {
			res.msg = "exceeded retries waiting for component to be Off"
		}
	}

	return res
}

func (d *CapmcD) doCompOnOffCtrl(nl []*NodeInfo, command string) capmc.XnameControlResponse {
	var data capmc.XnameControlResponse
	data.Xnames = make([]*capmc.XnameControlErr, 0, 1)

	// Build the sorted xname list here. The nl contains all of the
	// components that exist and should be powered on/off or restarted.
	// cmap maps lists of components, sorted by component type, to the
	// power action to be taken on them.
	// cmap[on/off/reinit][ComponentType]:[xname,xname,xname]
	// A single component may end up in multiple lists if multiple power
	// actions are being taken on it. For instance if a reinit is requested
	// on a node but the node doesn't support restart. A restart may be
	// simulated by doing an 'off' action then a power 'on' action. The
	// component would then end up in cmap[off][node]:[xname] and
	// cmap[on][node]:[xname].
	cmap := make(map[string]map[string][]*NodeInfo)
	cmap[command] = make(map[string][]*NodeInfo)
	totalWait, failures := d.doSortActionMap(nl, cmap, command, &data)

	if failures > 0 {
		// Fail the operation before doBmcCall()
		data.ErrResponse.E = -1
		data.ErrResponse.ErrMsg = fmt.Sprintf("Errors encountered with %d components for %s",
			failures, command)

		return data
	}

	var targetedXname []string
	for _, v := range nl {
		targetedXname = append(targetedXname, v.Hostname)
	}

	// Grab the new list just in case a power off was done on a chassis
	// or compute module
	targetedXname, err := d.reserveComponents(targetedXname, command)
	defer d.releaseComponents(targetedXname)

	if err != nil {
		errstr := fmt.Sprintf("Error: Failed to reserve components while performing a %s.", command)
		log.Printf(errstr)
		data.ErrResponse.E = 37 // ENOLCK
		data.ErrResponse.ErrMsg = errstr
		return data
	}

	// The order of keys when ranging over a map is undefined in newer
	// versions of Go. For this reason, range over an ordered list of
	// power actions, 'ReinitActionSeq', to control the order.
	for actionNum, cmd := range d.ReinitActionSeq {
		// Skip the power action if it has no power sequence defined.
		compPowerSeq, err := d.cmdCompPowerSeq(cmd)
		if err != nil {
			continue
		}

		// Skip the power action if there is no map of components
		// defined or the map is empty.
		ct, ok := cmap[cmd]
		if !ok || len(ct) == 0 {
			continue
		}

		for _, t := range compPowerSeq {
			if list, ok := cmap[cmd][t]; ok {
				waitNum, waitCh := d.queueBmcCmd(bmcCmd{cmd: cmd}, list)

				for i := 0; i < waitNum; i++ {
					// Wait for each task to complete.
					result := <-waitCh
					// If any BMC/Node op fails then record overall error code for the
					// call. The HTTP code returned will still be 200.
					if result.rc != 0 {
						failures++
						xnameErr := capmc.MakeXnameError(result.ni.Hostname, result.rc, result.msg)
						data.Xnames = append(data.Xnames, xnameErr)
						// Check any future actions for the same xname and remove it.
						// For example, if a component fails for 'Off' we shouldn't
						// try 'On'
						doRemoveComp(cmap, t, result.ni.Hostname, d.ReinitActionSeq[actionNum+1:])
					}

					if (cmd == bmcCmdPowerOff || cmd == bmcCmdPowerForceOff) &&
						result.rc == 0 {
						offResult := d.waitForOff(result.ni)
						if offResult.rc != 0 {
							failures++
							xnameErr := capmc.MakeXnameError(offResult.ni.Hostname, offResult.rc, offResult.msg)
							data.Xnames = append(data.Xnames, xnameErr)
							// Check any future actions for the same xname and remove it.
							doRemoveComp(cmap, t, result.ni.Hostname, d.ReinitActionSeq[actionNum+1:])
						}
					}
				}

				// Do a simple wait to allow time for the components to be powered
				// on. Can't start the next batch of power ons until the previous
				// batch is done. No need to wait for the Nodes or RouterModules
				// since there is nothing below them currently that has power
				// control. Time estimates pulled from demo scripts for Q1'19.
				switch cmd {
				case bmcCmdPowerOn, bmcCmdPowerForceOn:
					switch t {
					case "Chassis":
						time.Sleep(90 * time.Second)
					case "ComputeModule":
						time.Sleep(15 * time.Second)
					default:
					}
				}
			}
		}
	}

	if failures > 0 {
		data.ErrResponse.E = -1
		data.ErrResponse.ErrMsg = fmt.Sprintf("Errors encountered with %d/%d Xnames issued %s",
			failures, totalWait, command)
	}

	return data
}

func (d *CapmcD) doCompStatus(nl []*NodeInfo, command string, filter uint) capmc.XnameStatusResponse {
	// The JSON encoder omits empty lists. The Cascade CAPMC API response
	// contains more lists than this, but at this time these are the only
	// ones that are reported.
	var data capmc.XnameStatusResponse
	data.On = make([]string, 0, 1)
	data.Off = make([]string, 0, 1)
	data.Undefined = make([]string, 0, 1)

	waitNum, waitCh := d.queueBmcCmd(bmcCmd{cmd: bmcCmdPowerStatus}, nl)

	var failures int
	for i := 0; i < waitNum; i++ {
		// Wait for each task to complete.
		result := <-waitCh
		// If any BMC/Node op fails then record overall error code for
		// the call. The HTTP code returned will still be 200.
		if result.rc != 0 {
			failures++
			// TODO develop a method to return more detailed failure info
		}

		switch result.state {
		case rf.POWER_STATE_ON:
			if filter&capmc.FilterShowOnBit != 0 {
				data.On = append(data.On, result.ni.Hostname)
			}
		case rf.POWER_STATE_OFF:
			if filter&capmc.FilterShowOffBit != 0 {
				data.Off = append(data.Off, result.ni.Hostname)
			}
		// Other hardware states are not implemented
		default:
			// These are reported if there's a problem communicating
			// with the BMC.
			data.Undefined = append(data.Undefined, result.ni.Hostname)
		}
	}

	if failures > 0 {
		data.ErrResponse.E = -1
		data.ErrResponse.ErrMsg =
			fmt.Sprintf("Errors encountered with %d/%d Xnames for %s",
				failures, waitNum, command)
	}

	// Sorting is mostly for convenience; not strictly needed for capmc.
	sort.Sort(xnameSlice{&data.On})
	sort.Sort(xnameSlice{&data.Off})
	sort.Sort(xnameSlice{&data.Undefined})

	return data
}
