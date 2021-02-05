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
Test case for xnameOnByCLI
"""
from sys import exit

from capmcLib import getXnameState, xnamePower, xnameWait

################################################################################
#
#   xnameOnByCLI
#
################################################################################
def xnameOnByCLI():
    TEST = "xnameOnByCLI"
    xname = "x0c0s28b0n0"
    OFF = "off"
    ON = "on"

    print("["+TEST+"] Test powering on of xname "+xname)

    # Get current xname status
    print("["+TEST+"] Getting initial xname status for "+xname)
    origState = getXnameState(xname)
    if origState["errcode"] != 0:
        print("["+TEST+"] FAIL: trying to get the status for xname "+xname+": "+origState["errstr"])
        return 1

    # Turn off the xname
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

    # Power on xname
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

    curState = getXnameState(xname)
    if curState["errcode"] != 0:
        print("["+TEST+"] FAIL: trying to get status for xname "+xname+": "+curState["errstr"])
        return 1

    if curState["state"] != origState["state"]:
        print("["+TEST+"] Returning xname "+xname+" to the off state.")
        # Turn off the xname
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

    print("["+TEST+"] PASS: Xname "+xname+" properly powered on.")
    return 0

def test_xnameOnByCLI():
    assert xnameOnByCLI() == 0

if __name__ == "__main__":
    ret = xnameOnByCLI()
    exit(ret)