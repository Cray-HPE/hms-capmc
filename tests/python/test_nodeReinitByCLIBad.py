#!/usr/bin/python3
# Copyright 2019 Cray Inc. All Rights Reserved
"""
Test case for nodeReinitByCLIBad
"""
from sys import exit

from capmcLib import nodePower

################################################################################
#
#   nodeReinitByCLIBad
#
################################################################################
def nodeReinitByCLIBad():
    TEST = "nodeReinitByCLIBad"
    nid = "99999"
    REINIT = "reinit"

    print("["+TEST+"] Test reinit of bad node "+nid)

    # Reinit node
    print("["+TEST+"] Reiniting nid "+nid)
    ret = nodePower(REINIT, nid)
    if ret["errcode"] != 0:
        print("["+TEST+"] PASS: expected failure when trying to reinit nid "+nid+": "+ret["errstr"])
        return 0

    print("["+TEST+"] FAIL: Did not receive a failure when trying to reinit node "+nid)
    return 1

def test_nodeReinitByCLIBad():
    assert nodeReinitByCLIBad() == 0

if __name__ == "__main__":
    ret = nodeReinitByCLIBad()
    exit(ret)