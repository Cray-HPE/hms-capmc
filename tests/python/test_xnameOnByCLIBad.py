#!/usr/bin/python3
# Copyright 2019 Cray Inc. All Rights Reserved
"""
Test case for xnameOnByCLIBad
"""
from sys import exit

from capmcLib import xnamePower

################################################################################
#
#   xnameOnByCLIBad
#
################################################################################
def xnameOnByCLIBad():
    TEST = "xnameOnByCLIBad"
    xname = "x99999c9s9b0n9"
    ON = "on"

    print("["+TEST+"] Test powering on of bad xname "+xname)

    # Power on xname
    print("["+TEST+"] Powering on xname "+xname)
    ret = xnamePower(ON, xname)
    if ret["errcode"] != 0:
        print("["+TEST+"] PASS: expected failure when trying to power on xname "+xname+": "+ret["errstr"])
        return 0

    print("["+TEST+"] FAIL: Did not receive a failure when trying to power on xname "+xname)
    return 1

def test_xnameOnByCLIBad():
    assert xnameOnByCLIBad() == 0

if __name__ == "__main__":
    ret = xnameOnByCLIBad()
    exit(ret)