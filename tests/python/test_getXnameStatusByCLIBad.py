#!/usr/bin/python3
# Copyright 2019 Cray Inc. All Rights Reserved
"""
Test case for getXnameStatusByCLIBad
"""
from subprocess import Popen, PIPE
from shlex import split
from sys import exit
from re import search

################################################################################
#
#   getXnameStatusByCLIBad
#
################################################################################
def getXnameStatusByCLIBad():
    TEST = "getNodeStatusByCLIBad"
    xnames = "x99999c9s9b0n9"
    CMD = "cray capmc get_xname_status create --xnames "+xnames

    print("["+TEST+"] Checks the status of invalid xname: "+CMD)

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

def test_getXnameStatusByCLIBad():
    assert getXnameStatusByCLIBad() == 0

if __name__ == "__main__":
    ret = getXnameStatusByCLIBad()
    exit(ret)
