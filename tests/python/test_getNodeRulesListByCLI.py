#!/usr/bin/python3
# Copyright 2019 Cray Inc. All Rights Reserved
"""
Test case for getNodeRulesListByCLI
"""
from subprocess import Popen, PIPE
from json import loads, dumps
from shlex import split
from sys import exit

################################################################################
#
#   getNodeRulesListByCLI
#
################################################################################
def getNodeRulesListByCLI():
    TEST = "getNodeRulesListByCLI"
    CMD = "cray capmc get_node_rules list"

    print("["+TEST+"] Queries for valid node rules: "+CMD)

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

    fields = ("e", "err_msg", "latency_node_off", "latency_node_on",
        "latency_node_reinit", "max_off_req_count", "max_off_time",
        "max_on_req_count", "max_reinit_req_count", "min_off_time")

    for f in fields:
        if f not in map:
            print("["+TEST+"] FAIL: Missing "+f+" field: "+dumps(map))
            return 1

    print("["+TEST+"] PASS: Valid node rules returned.")
    return 0

def test_getNodeRulesListByCLI():
    assert getNodeRulesListByCLI() == 0

if __name__ == "__main__":
    ret = getNodeRulesListByCLI()
    exit(ret)

