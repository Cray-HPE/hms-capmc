#!/usr/bin/python3
#  MIT License
#
#  (C) Copyright [2019-2022] Hewlett Packard Enterprise Development LP
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
