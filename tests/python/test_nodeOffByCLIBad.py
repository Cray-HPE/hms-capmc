#!/usr/bin/python3
# Copyright 2019 Cray Inc. All Rights Reserved
"""
Test case for nodeOffByCLIBad
"""
from sys import exit

from capmcLib import nodePower

################################################################################
#
#   nodeOffByCLIBad
#
################################################################################
def nodeOffByCLIBad():
    TEST = "nodeOffByCLIBad"
    nid = "99999"
    OFF = "off"

    print("["+TEST+"] Test powering off of bad node "+nid)

    # Power off node
    print("["+TEST+"] Powering off nid "+nid)
    ret = nodePower(OFF, nid)
    if ret["errcode"] != 0:
        print("["+TEST+"] PASS: expected failure when trying to power off nid "+nid+": "+ret["errstr"])
        return 0

    print("["+TEST+"] FAIL: Did not receive a failure when trying to power off node "+nid)
    return 1

def test_nodeOffByCLIBad():
    assert nodeOffByCLIBad() == 0

if __name__ == "__main__":
    ret = nodeOffByCLIBad()
    exit(ret)

