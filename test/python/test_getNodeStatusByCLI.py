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
Test case for getNodeStatusByCLI
"""
from subprocess import Popen, PIPE
from json import loads, dumps
from shlex import split
from sys import exit

################################################################################
#
#   getNodeStatusByCLI
#
################################################################################
def getNodeStatusByCLI():
    TEST = "getNodeStatusByCLI"
    nids = "1,2,3,4"
    CMD = "cray capmc get_node_status create --nids "+nids

    print("["+TEST+"] Checks the status return of nodes: "+CMD)

    process = Popen(split(CMD), stdout=PIPE, stderr=PIPE)
    process.wait()
    stdout, stderr = process.communicate()

    if process.returncode != 0:
        errstr = stderr.decode("utf-8").split(":")[2].lstrip()
        print("["+TEST+"] FAIL: "+errstr)
        return 1

    if stdout == None:
        print("["+TEST+"] FAIL: Nothing returned.")
        return 1

    map = loads(stdout)

    fields = ("e", "err_msg")
    for f in fields:
        if f not in map:
            print("["+TEST+"] FAIL: Missing "+f+" field: "+dumps(map))
            return 1

    stateFields = ("undefined", "on", "off", "halt", "standby", "ready", "diag", "disabled")
    valid=False
    for f in stateFields:
        if f in map:
            valid=True
            break

    if valid == False:
        print("["+TEST+"] FAIL: Missing state field: "+dumps(map))
        return 1

    print("["+TEST+"] PASS: Valid status returned.")
    return 0

def test_getNodeStatusByCLI():
    assert getNodeStatusByCLI() == 0

if __name__ == "__main__":
    ret = getNodeStatusByCLI()
    exit(ret)