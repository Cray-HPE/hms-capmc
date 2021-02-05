#!/usr/bin/python3
#  MIT License
#
#  (C) Copyright [2019-2021] Hewlett Packard Enterprise Development LP
#
#  Permission is hereby granted, free of charge, to any person obtaining a
#  copy of this software and associated documentation files (the "Software"),
#  to deal in the Software without restriction, including without limitation
#  the rights to use, copy, modify, merge, publish, distribute, sublicense,
#  and/or sell copies of the Software, and to permit persons to whom the
#  Software is furnished to do so, subject to the following conditions:
#
#  The above copyright notice and this permission notice shall be included
#  in all copies or substantial portions of the Software.
#
#  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
#  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
#  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
#  THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
#  OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
#  ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
#  OTHER DEALINGS IN THE SOFTWARE.
"""
Test case for nodeOnByCLI
"""
from sys import exit

from capmcLib import getNodeState, nodePower, nodeWait

################################################################################
#
#   nodeOnByCLI
#
################################################################################
def nodeOnByCLI():
    TEST = "nodeOnByCLI"
    nid = "1"
    OFF = "off"
    ON = "on"

    print("["+TEST+"] Test powering on of node "+nid)

    # Get current node status
    print("["+TEST+"] Getting initial node status for "+nid)
    origState = getNodeState(nid)
    if origState["errcode"] != 0:
        print("["+TEST+"] FAIL: trying to get the status for nid "+nid+": "+origState["errstr"])
        return 1

    # Turn off the node
    print("["+TEST+"] Powering off nid "+nid)
    ret = nodePower(OFF, nid)
    if ret["errcode"] != 0:
        print("["+TEST+"] FAIL: trying to power off nid "+nid+": "+ret["errstr"])
        return 1

    # Wait for node to power off
    print("["+TEST+"] Waiting for nid "+nid+" to power off")
    ret = nodeWait(nid, "Standby", 60)
    if ret["errcode"] != 0:
        print("["+TEST+"] FAIL: trying to power off nid "+nid+": "+ret["errstr"])
        return 1

    # Power on node
    print("["+TEST+"] Powering on nid "+nid)
    ret = nodePower(ON, nid)
    if ret["errcode"] != 0:
        print("["+TEST+"] FAIL: trying to power on nid "+nid+": "+ret["errstr"])
        return 1

    # Wait for node to power on
    print("["+TEST+"] Waiting for nid "+nid+" to power on")
    ret = nodeWait(nid, "Ready", 240)
    if ret["errcode"] != 0:
        print("["+TEST+"] FAIL: trying to power on nid "+nid+": "+ret["errstr"])
        return 1

    curState = getNodeState(nid)
    if curState["errcode"] != 0:
        print("["+TEST+"] FAIL: trying to get status for nid "+nid+": "+curState["errstr"])
        return 1

    if curState["state"] != origState["state"]:
        print("["+TEST+"] Returning nid "+nid+" to the off state.")
        # Turn off the node
        print("["+TEST+"] Powering off nid "+nid)
        ret = nodePower(OFF, nid)
        if ret["errcode"] != 0:
            print("["+TEST+"] FAIL: trying to power off nid "+nid+": "+ret["errstr"])
            return 1

        # Wait for node to power off
        print("["+TEST+"] Waiting for nid "+nid+" to power off")
        ret = nodeWait(nid, "Standby", 60)
        if ret["errcode"] != 0:
            print("["+TEST+"] FAIL: trying to power off nid "+nid+": "+ret["errstr"])
            return 1

    print("["+TEST+"] PASS: Node "+nid+" properly powered on.")
    return 0

def test_nodeOnByCLI():
    assert nodeOnByCLI() == 0

if __name__ == "__main__":
    ret = nodeOnByCLI()
    exit(ret)