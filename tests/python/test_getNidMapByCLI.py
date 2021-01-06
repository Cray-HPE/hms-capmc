#!/usr/bin/python3
# Copyright 2019 Cray Inc. All Rights Reserved
"""
Test case for getNidMapByCLI
"""
from subprocess import Popen, PIPE
from json import loads, dumps
from sys import exit
from shlex import split

################################################################################
#
#   getNidMapByCLI
#
################################################################################
def getNidMapByCLI():
    TEST = "getNidMapByCLI"
    nids = "1,2,3,4"
    CMD = "cray capmc get_nid_map create --nids "+nids

    print("["+TEST+"] Checks for valid nid map with valid nodes: "+CMD)

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

    fields = ("e", "err_msg", "nids")
    for f in fields:
        if f not in map:
            print("["+TEST+"] FAIL: Missing "+f+" field: "+dumps(map))
            return 1

    if len(map["nids"]) == 0:
        print("["+TEST+"] FAIL: No nids in map: "+dumps(map))
        return 1

    nfields = ("cname", "nid", "role")
    for n in map["nids"]:
        for f in nfields:
            if f not in n:
                print("["+TEST+"] FAIL: Missing "+f+" field: "+dumps(n))
                return 1

    print("["+TEST+"] PASS: Valid nid map returned.")
    return 0

def test_getNidMapByCLI():
    assert getNidMapByCLI() == 0

if __name__ == "__main__":
    ret = getNidMapByCLI()
    exit(ret)
