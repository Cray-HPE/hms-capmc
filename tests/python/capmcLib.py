#!/usr/bin/python3
# Copyright 2019 Cray Inc. All Rights Reserved
"""
TOOL IDENTIFIER : capmcLib.py

TOOL TITLE      : Common Library for testing CAPMC

AUTHOR          : Michael Jendrysik

DATE STARTED    : 03/2019
"""

from subprocess import Popen, PIPE
from json import loads
from shlex import split
from time import sleep

################################################################################
#
#   getSystemNids() returns string of command separated nids
#
################################################################################
#function getSystemNids {
    #HOST=$(cray config get core.hostname)
    #CMD="$CURL -k -X POST -d '{}' ${HOST}/apis/capmc/capmc/get_nid_map"
    #json=$(eval $CMD)
    #echo $json
#}

"""
    getNodeState(nid)

    Queries the state manager to get the current status of a single node

    Arguments:
        string nid

    Returns:
        A dictionary containing an error code, an error string, and a state
        string.
        Dictionary:
            int errcode   - Error returned from Popen's communicate() or -1 if
                            there was no status returned. Default: 0
            string errstr - Error indicated by the failure type. Default: ""
            string state  - State of the node: on, off, halt, etc. Default: ""
"""
def getNodeState(nid):
    #command = "curl -k -X GET https://slice-sms.us.cray.com/apis/smd/hsm/v1/State/Components/ByNID/"+nid
    command = "cray hsm State Components ByNID describe "+nid
    proc = Popen(split(command), stdout=PIPE, stderr=PIPE)
    proc.wait()
    stdout, stderr = proc.communicate()

    state = dict()
    state["errcode"] = 0
    state["errstr"] = ""
    state["state"] = ""

    if proc.returncode != 0:
        state["errcode"] = proc.returncode
        state["errstr"] = stderr.decode("utf-8").split(":")[2].lstrip()
        return state

    if stdout == None:
        state["errcode"] = -1
        state["errstr"] = "No status returned"
        return state

    stateMap = loads(stdout)
    state["state"] = stateMap["State"]

    return state

"""
    nodePower(opt, nid)

    Uses craycli to send a power off request to capmc for the specific node

    Arguments:
        string opt  - "on", "off", "reinit"
        string nid

    Returns:
        A dictinoary containing an error code and an error string.

        Dictionary:
            int errcode   - Error returned from Popen's communicate() or -1 if
                            there was no response. Default: 0
            string errstr - Error indicated by the failure type. Default: ""
"""
def nodePower(opt, nid):
    command = "cray capmc node_"+opt+" create --nids "+nid
    proc = Popen(split(command), stdout=PIPE, stderr=PIPE)
    proc.wait()
    stdout, stderr = proc.communicate()

    ret = dict()
    ret["errcode"] = 0
    ret["errstr"] = ""

    if proc.returncode != 0:
        ret["errcode"] = proc.returncode
        ret["errstr"] = stderr.decode("utf-8").split(":")[2].lstrip()
        return ret

    if stdout == None:
        ret["errcode"] = -1
        ret["errstr"] = "Nothing returned"
        return ret

    return ret

"""
    nodeWait(nid, targ, maxloops)

    Checks to see if the status of the nid matches targ maxloops times.

    Arguments:
        string nid
        string targ     - Standby, Ready, etc.
        string maxloops

    Returns:
        A dictionary containing an error code and an error string.

        Dictionary:
            int errcode   - Error returned from Popen's communicate() or -1 if
                            there was no response. Default: 0
            string errstr - Error indicated by the failure type. Default: ""
"""
def nodeWait(nid, targ, maxloops):
    curState = getNodeState(nid)
    ret = dict()
    ret["errcode"] = 0
    ret["errstr"] = ""
    if curState["errcode"] != 0:
        ret["errcode"] = curState["errcode"]
        ret["errstr"] = curState["errstr"]
        return ret

    loops = 1
    while curState["state"] != targ:
        sleep(1)

        curState = getNodeState(nid)
        if curState["errcode"] != 0:
            ret["errcode"] = curState["errcode"]
            ret["errstr"] = curState["errstr"]
            return ret

        if loops > maxloops:
            ret["errcode"] = -1
            ret["errstr"] = "Timeout exceeded"
            return ret

    return ret

"""
    getXnameState(xname)

    Queries the state manager to get the current status of a single xname

    Arguments:
        string xname

    Returns:
        A dictionary containing an error code, an error string, and a state
        string.
        Dictionary:
            int errcode   - Error returned from Popen's communicate() or -1 if
                            there was no status returned. Default: 0
            string errstr - Error indicated by the failure type. Default: ""
            string state  - State of the xname: on, off, halt, etc. Default: ""
"""
def getXnameState(xname):
    #command = "curl -k -X GET https://slice-sms.us.cray.com/apis/smd/hsm/v1/State/Components/"+xname
    command = "cray hsm State Components describe "+xname
    proc = Popen(split(command), stdout=PIPE, stderr=PIPE)
    proc.wait()
    stdout, stderr = proc.communicate()

    state = dict()
    state["errcode"] = 0
    state["errstr"] = ""
    state["state"] = ""

    if proc.returncode != 0:
        state["errcode"] = proc.returncode
        state["errstr"] = stderr.decode("utf-8").split(":")[2].lstrip()
        return state

    if stdout == None:
        state["errcode"] = -1
        state["errstr"] = "No status returned"
        return state

    stateMap = loads(stdout)
    state["state"] = stateMap["State"]

    return state

"""
    xnamePower(opt, xname)

    Uses craycli to send a power off request to capmc for the specific xname

    Arguments:
        string opt  - "on", "off", "reinit"
        string xname

    Returns:
        A dictinoary containing an error code and an error string.

        Dictionary:
            int errcode   - Error returned from Popen's communicate() or -1 if
                            there was no response. Default: 0
            string errstr - Error indicated by the failure type. Default: ""
"""
def xnamePower(opt, xname):
    command = "cray capmc xname_"+opt+" create --xnames "+xname
    proc = Popen(split(command), stdout=PIPE, stderr=PIPE)
    proc.wait()
    stdout, stderr = proc.communicate()

    ret = dict()
    ret["errcode"] = 0
    ret["errstr"] = ""

    if proc.returncode != 0:
        ret["errcode"] = proc.returncode
        ret["errstr"] = stderr.decode("utf-8").split(":")[2].lstrip()
        return ret

    if stdout == None:
        ret["errcode"] = -1
        ret["errstr"] = "Nothing returned"
        return ret

    return ret

"""
    xnameWait(xname, targ, maxloops)

    Checks to see if the status of the xname matches targ maxloops times.

    Arguments:
        string xname
        string targ     - Standby, Ready, etc.
        string maxloops

    Returns:
        A dictionary containing an error code and an error string.

        Dictionary:
            int errcode   - Error returned from Popen's communicate() or -1 if
                            there was no response. Default: 0
            string errstr - Error indicated by the failure type. Default: ""
"""
def xnameWait(xname, targ, maxloops):
    curState = getXnameState(xname)
    ret = dict()
    ret["errcode"] = 0
    ret["errstr"] = ""
    if curState["errcode"] != 0:
        ret["errcode"] = curState["errcode"]
        ret["errstr"] = curState["errstr"]
        return ret

    loops = 1
    while curState["state"] != targ:
        sleep(1)

        curState = getXnameState(xname)
        if curState["errcode"] != 0:
            ret["errcode"] = curState["errcode"]
            ret["errstr"] = curState["errstr"]
            return ret

        if loops > maxloops:
            ret["errcode"] = -1
            ret["errstr"] = "Timeout exceeded"
            return ret

    return ret
