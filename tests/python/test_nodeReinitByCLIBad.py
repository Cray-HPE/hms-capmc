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