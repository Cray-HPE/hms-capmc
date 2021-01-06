#!/usr/bin/python3
# Copyright 2019 Cray Inc. All Rights Reserved
"""
Test case for getNodeStatusByCLIBad
"""
from subprocess import Popen, PIPE
from shlex import split
from sys import exit
from re import search

################################################################################
#
#   getNodeStatusByCLIBad
#
################################################################################
def getNodeStatusByCLIBad():
    TEST = "getNodeStatusByCLIBad"
    nids = "99999"
    CMD = "cray capmc get_node_status create --nids "+nids

    print("["+TEST+"] Checks the status of invalid node: "+CMD)

    process = Popen(split(CMD), stdout=PIPE, stderr=PIPE)
    process.wait()
    _, stderr = process.communicate()

    if process.returncode != 0:
        errstr = stderr.decode("utf-8")
        if search("400 Client Error", errstr):
            print("["+TEST+"] PASS: Received expected 400 Client Error.")
            return 0
        print("["+TEST+"] FAIL: "+errstr+", expected 400 Client Error.")
        return 1

    print("["+TEST+"] FAIL: No error, expected 400 Client Error")
    return 1

def test_getNodeStatusByCLIBad():
    assert getNodeStatusByCLIBad() == 0

if __name__ == "__main__":
    ret = getNodeStatusByCLIBad()
    exit(ret)