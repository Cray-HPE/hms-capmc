#!/usr/bin/python3
# Copyright 2019 Cray Inc. All Rights Reserved
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