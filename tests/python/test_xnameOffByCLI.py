#!/usr/bin/python3
# Copyright 2019 Cray Inc. All Rights Reserved
"""
Test case for xnameOffByCLI
"""
from sys import exit

from capmcLib import getXnameState, xnamePower, xnameWait

################################################################################
#
#   xnameOffByCLI
#
################################################################################
def xnameOffByCLI():
    TEST = "xnameOffByCLI"
    xname = "x0c0s28b0n0"
    OFF = "off"
    ON = "on"

    print("["+TEST+"] Test powering off of xname "+xname)

    # Get current xname status
    print("["+TEST+"] Getting initial xname status for "+xname)
    origState = getXnameState(xname)
    if origState["errcode"] != 0:
        print("["+TEST+"] FAIL: trying to get the status for xname "+xname+": "+origState["errstr"])
        return 1

    # Turn on the xname
    print("["+TEST+"] Powering on xname "+xname)
    ret = xnamePower(ON, xname)
    if ret["errcode"] != 0:
        print("["+TEST+"] FAIL: trying to power on xname "+xname+": "+ret["errstr"])
        return 1

    # Wait for xname to power on
    print("["+TEST+"] Waiting for xname "+xname+" to power on")
    ret = xnameWait(xname, "Ready", 240)
    if ret["errcode"] != 0:
        print("["+TEST+"] FAIL: trying to power on xname "+xname+": "+ret["errstr"])
        return 1

    # Power off xname
    print("["+TEST+"] Powering off xname "+xname)
    ret = xnamePower(OFF, xname)
    if ret["errcode"] != 0:
        print("["+TEST+"] FAIL: trying to power off xname "+xname+": "+ret["errstr"])
        return 1

    # Wait for xname to power off
    print("["+TEST+"] Waiting for xname "+xname+" to power off")
    ret = xnameWait(xname, "Standby", 60)
    if ret["errcode"] != 0:
        print("["+TEST+"] FAIL: trying to power off xname "+xname+": "+ret["errstr"])
        return 1

    curState = getXnameState(xname)
    if curState["errcode"] != 0:
        print("["+TEST+"] FAIL: trying to get status for xname "+xname+": "+curState["errstr"])
        return 1

    if curState["state"] != origState["state"]:
        print("["+TEST+"] Returning xname "+xname+" to the on state.")
        # Turn on the xname
        print("["+TEST+"] Powering on xname "+xname)
        ret = xnamePower(ON, xname)
        if ret["errcode"] != 0:
            print("["+TEST+"] FAIL: trying to power on xname "+xname+": "+ret["errstr"])
            return 1

        # Wait for xname to power on
        print("["+TEST+"] Waiting for xname "+xname+" to power on")
        ret = xnameWait(xname, "Ready", 240)
        if ret["errcode"] != 0:
            print("["+TEST+"] FAIL: trying to power on xname "+xname+": "+ret["errstr"])
            return 1

    print("["+TEST+"] PASS: Xname "+xname+" properly powered off.")
    return 0

def test_xnameOffByCLI():
    assert xnameOffByCLI() == 0

if __name__ == "__main__":
    ret = xnameOffByCLI()
    exit(ret)