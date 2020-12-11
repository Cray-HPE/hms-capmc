#!/usr/bin/python3
# Copyright 2019 Cray Inc. All Rights Reserved
"""
Test case for nodeOnByCLIBad
"""
from sys import exit

from capmcLib import nodePower

################################################################################
#
#   nodeOnByCLIBad
#
################################################################################
def nodeOnByCLIBad():
    TEST = "nodeOnByCLIBad"
    nid = "99999"
    ON = "on"

    print("["+TEST+"] Test powering on of bad node "+nid)

    # Power on node
    print("["+TEST+"] Powering on nid "+nid)
    ret = nodePower(ON, nid)
    if ret["errcode"] != 0:
        print("["+TEST+"] PASS: expected failure when trying to power on nid "+nid+": "+ret["errstr"])
        return 0

    print("["+TEST+"] FAIL: Did not receive a failure when trying to power on node "+nid)
    return 1

def test_nodeOnByCLIBad():
    assert nodeOnByCLIBad() == 0

if __name__ == "__main__":
    ret = nodeOnByCLIBad()
    exit(ret)