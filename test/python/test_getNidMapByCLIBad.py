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
Test case for getNidMapByCLIBad
"""
from subprocess import Popen, PIPE
from shlex import split
from sys import exit
from re import search

################################################################################
#
#   getNidMapByCLIBad
#
################################################################################
def getNidMapByCLIBad():
    TEST = "getNidMapByCLIBad"
    nids = "99999"
    CMD = "cray capmc get_nid_map create --nids "+nids

    print("["+TEST+"] Checks for error when requesting an invalid node: "+CMD)

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

def test_getNidMapByCLIBad():
    assert getNidMapByCLIBad() == 0

if __name__ == "__main__":
    ret = getNidMapByCLIBad()
    exit(ret)
