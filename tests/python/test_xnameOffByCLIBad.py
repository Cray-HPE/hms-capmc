#!/usr/bin/python3
# Copyright 2019 Cray Inc. All Rights Reserved
"""
Test case for xnameOffByCLIBad
"""
from sys import exit

from capmcLib import xnamePower

################################################################################
#
#   xnameOffByCLIBad
#
################################################################################
def xnameOffByCLIBad():
    TEST = "xnameOffByCLIBad"
    xname = "x99999c9s9b0n9"
    OFF = "off"

    print("["+TEST+"] Test powering off of bad xname "+xname)

    # Power off xname
    print("["+TEST+"] Powering off xname "+xname)
    ret = xnamePower(OFF, xname)
    if ret["errcode"] != 0:
        print("["+TEST+"] PASS: expected failure when trying to power off xname "+xname+": "+ret["errstr"])
        return 0

    print("["+TEST+"] FAIL: Did not receive a failure when trying to power off xname "+xname)
    return 1

def test_xnameOffByCLIBad():
    assert xnameOffByCLIBad() == 0

if __name__ == "__main__":
    ret = xnameOffByCLIBad()
    exit(ret)