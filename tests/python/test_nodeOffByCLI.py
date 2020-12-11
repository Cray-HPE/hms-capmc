#!/usr/bin/python3
# Copyright 2019 Cray Inc. All Rights Reserved
"""
Test case for nodeOffByCLI
"""
from sys import exit

from capmcLib import getNodeState, nodePower, nodeWait

################################################################################
#
#   nodeOffByCLI
#
################################################################################
def nodeOffByCLI():
    TEST = "nodeOffByCLI"
    nid = "1"
    OFF = "off"
    ON = "on"

    print("["+TEST+"] Test powering off of node "+nid)

    # Get current node status
    print("["+TEST+"] Getting initial node status for "+nid)
    origState = getNodeState(nid)
    if origState["errcode"] != 0:
        print("["+TEST+"] FAIL: trying to get the status for nid "+nid+": "+origState["errstr"])
        return 1

    # Turn on the node
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

    # Power off node
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

    curState = getNodeState(nid)
    if curState["errcode"] != 0:
        print("["+TEST+"] FAIL: trying to get status for nid "+nid+": "+curState["errstr"])
        return 1

    if curState["state"] != origState["state"]:
        print("["+TEST+"] Returning nid "+nid+" to the on state.")
        # Turn on the node
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

    print("["+TEST+"] PASS: Node "+nid+" properly powered off.")
    return 0

def test_nodeOffByCLI():
    assert nodeOffByCLI() == 0

if __name__ == "__main__":
    ret = nodeOffByCLI()
    exit(ret)