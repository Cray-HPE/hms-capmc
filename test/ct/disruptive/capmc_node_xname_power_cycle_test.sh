#!/bin/bash -l
#
# Copyright 2019 Cray Inc. All Rights Reserved.
#
###############################################################
#
#     CASM Test - Cray Inc.
#
#     TEST IDENTIFIER   : capmc_node_xname_power_cycle_test
#
#     DESCRIPTION       : Automated test for verifying CAPMC xname
#                         control APIs on Cray Shasta systems
#                         
#     AUTHOR            : Mitch Schooler
#
#     DATE STARTED      : 03/19/2019
#
#     LAST MODIFIED     : 10/07/2019
#
#     SYNOPSIS
#       This test verifies the CAPMC xname control APIs for Shasta including
#       the xname_status, xname_off, and xname_on APIs by making remote curl
#       calls to the APIs and verifying that the expected data is returned
#       by xname_status and that the intended power operation for the target
#       nodes actually works by polling until the expected state changes occur.
#
#     INPUT SPECIFICATIONS
#       Usage: capmc_node_xname_power_cycle_test [-h] [-d] [-s <num>] [-i <num>]
#       
#       Arguments:
#               -h         display this help message
#               -d         enable detailed debugging output
#               -s <num>   wait4boot polling iteration sleep time [default is 20s]
#               -i <num>   maximum number of wait4boot polling iterations [default is 30]
#
#     OUTPUT SPECIFICATIONS
#       Plaintext is printed to stdout and/or stderr. The script exits with
#       a status of '0' on success and '1' on failure.
#
#     DESIGN DESCRIPTION
#       The nodes targeted for the test are all node xnames returned by the
#       initial xname_status call. If an unexpected error occurs after nodes
#       have been brought down with xname_off, one attempt to power them back
#       up and reboot them is made. Sending a first kill signal to the test
#       during execution will also result in one attempt to bring the test
#       nodes back up before exiting if nodes have already been powered down.
#       Subsequent kill signals will result in the script exiting prematurely
#       and potentially leaving test nodes powered off.
#
#     SPECIAL REQUIREMENTS
#       Must be executed from the NCN.
#
#     UPDATE HISTORY
#       user       date         description
#       -------------------------------------------------------
#       schooler   03/29/2019   initial implementation
#       schooler   04/09/2019   add port update workaround for CASMCMS-2257
#       schooler   04/26/2019   add SLURM workaround for DST-1866
#       schooler   07/10/2019   add Keycloak AuthN support
#       schooler   07/10/2019   update smoke test library location
#                               from hms-services to hms-common
#       schooler   10/07/2019   update smoke test library location
#                               from hms-common to hms-test
#       schooler   10/07/2019   switch from SMS to NCN naming convention
#
#     DEPENDENCIES
#       - hms_smoke_test_lib_ncn-resources.sh which is expected to be packaged
#         in /opt/cray/tests/ncn-resources/hms/hms-test on the NCN.
#
#     BUGS/LIMITATIONS
#       None
#
###############################################################

# cleanup
cleanup()
{
    echo "cleaning up..."
    if ${RUN_POWER_ON_CLEANUP} ; then
        run_verify_xname_on
    fi
    if ${RUN_WAIT4BOOT_CLEANUP} ; then
        wait4boot
    fi
    if ${RUN_SLURM_WORKAROUND_CLEANUP} ; then
        slurm_workaround
    fi
}

# echoerr_exit error_message exit_code
echoerr_exit()
{
        >&2 echo -e "${1}"
        exit ${2}
}

# debug_print message
debug_print()
{
    if ${VERBOSE} ; then
        echo "DEBUG: ${1}"
    fi
}

# timestamp_print message
timestamp_print()
{
    echo "($(date +"%H:%M:%S")) ${1}"
}

# run_verify_xname_on
run_verify_xname_on()
{
    # generate json payload for xname_on operation
    XNAME_ON_REQUEST_BODY_XNAMES_ARG=""
    for XNAME in ${TEST_NODE_XNAMES} ; do
        XNAME_TOKEN="\"${XNAME}\","
        XNAME_ON_REQUEST_BODY_XNAMES_ARG="${XNAME_ON_REQUEST_BODY_XNAMES_ARG}${XNAME_TOKEN}"
        debug_print "${XNAME_ON_REQUEST_BODY_XNAMES_ARG}"
    done
    XNAME_ON_REQUEST_BODY_XNAMES_ARG=$(echo "${XNAME_ON_REQUEST_BODY_XNAMES_ARG}" | sed 's/,$//')
    debug_print "${XNAME_ON_REQUEST_BODY_XNAMES_ARG}"
    XNAME_ON_REQUEST_BODY_JSON="{\"reason\":\"capmc_node_xname_power_cycle_test.${TEST_RUN_TIMESTAMP}\",\"xnames\":[${XNAME_ON_REQUEST_BODY_XNAMES_ARG}]}"
    debug_print "${XNAME_ON_REQUEST_BODY_JSON}"

    # call xname_on API
    XNAME_ON_CMD="curl -k -X POST -i -H 'Authorization: Bearer ${TOKEN}' -d '${XNAME_ON_REQUEST_BODY_JSON}' https://${HOST}:${PORT}/${CAPMC_URI_PATH}/xname_on"
    timestamp_print "running '${XNAME_ON_CMD}'..."
    XNAME_ON_OUT=$(eval ${XNAME_ON_CMD} 2> /dev/null)
    XNAME_ON_RET=$?
    debug_print "${XNAME_ON_OUT}"
    if [[ ${XNAME_ON_RET} -eq 0 ]] ; then
        echo "PASS: '${XNAME_ON_CMD}' returned success code: ${XNAME_ON_RET}"
    else
        echo "FAIL: '${XNAME_ON_CMD}' failed with error code: ${XNAME_ON_RET}, expected: 0'"
        exit 1
    fi
    XNAME_ON_HTTP_STATUS_CHECK=$(echo "${XNAME_ON_OUT}" | grep -E "200 OK")
    if [[ -n "${XNAME_ON_HTTP_STATUS_CHECK}" ]] ; then
        echo "PASS: '${XNAME_ON_CMD}' returned \"200 OK\""
    else
        echo "${XNAME_ON_OUT}"
        echo "FAIL: '${XNAME_ON_CMD}' failed to return \"200 OK\""
        exit 1
    fi

    XNAME_ON_ITER=0
    XNAME_ON_COMPLETE=false
    # poll until all test nodes have been powered on 
    echo "polling until all test nodes have powered on..."
    echo
    while [[ ${XNAME_ON_COMPLETE} == "false" ]] && [[ ${XNAME_ON_ITER} -lt ${XNAME_ON_ITER_MAX} ]] ; do
    
        echo "ITER=${XNAME_ON_ITER}"
        GET_XNAME_STATUS_CMD="curl -k -X POST -i -H 'Authorization: Bearer ${TOKEN}' -d '{}' https://${HOST}:${PORT}/${CAPMC_URI_PATH}/get_xname_status"
        #GET_XNAME_STATUS_CMD="cat /tmp/schooler/fake"
        timestamp_print "running '${GET_XNAME_STATUS_CMD}'..."
        GET_XNAME_STATUS_OUT=$(eval ${GET_XNAME_STATUS_CMD} 2> /dev/null)
        GET_XNAME_STATUS_RET=$?
        debug_print "${GET_XNAME_STATUS_OUT}"
        if [[ ${GET_XNAME_STATUS_RET} -ne 0 ]] ; then
            echo "FAIL: '${GET_XNAME_STATUS_CMD}' failed with error code: ${GET_XNAME_STATUS_RET}, expected: 0'"
            exit 1
        fi
        GET_XNAME_STATUS_HTTP_STATUS_CHECK=$(echo "${GET_XNAME_STATUS_OUT}" | grep -E "200 OK")
        if [[ -z "${GET_XNAME_STATUS_HTTP_STATUS_CHECK}" ]] ; then
            echo "${GET_XNAME_STATUS_OUT}"
            echo "FAIL: '${GET_XNAME_STATUS_CMD}' failed to return \"200 OK\""
            exit 1
        fi

        # flatten JSON response
        GET_XNAME_STATUS_OUT_JSON=$(echo "${GET_XNAME_STATUS_OUT}" | sed -n -e "/{/,/}/p" | sed 's/^[[:space:]]*//g' | tr -d "\n")
        if [[ -z ${GET_XNAME_STATUS_OUT_JSON} ]] ; then
            echo "${GET_XNAME_STATUS_OUT}"
            echo "FAIL: '${GET_XNAME_STATUS_CMD}' failed to return reponse json"
            exit 1
        fi

        GET_XNAME_STATUS_STATES=$(echo "${GET_XNAME_STATUS_OUT_JSON}" | grep -E -o "\"[a-zA-Z0-9]+\":\[(\"[a-zA-Z0-9]+\")(,\"[a-zA-Z0-9]+\")*\]")
        GET_XNAME_STATUS_STATES_ON_CHECK=$(echo "${GET_XNAME_STATUS_STATES}" | grep -E "\"on\":\[")
        if [[ -n ${GET_XNAME_STATUS_STATES_ON_CHECK} ]] ; then
            # parse the get_xname_status output by state key
            while read -r line ; do
                echo "${line}"
                LINE_STATE=$(echo "${line}" | cut -d ":" -f 1 | grep -E -o "\".*\"" | tr -d "\"")
                LINE_NODE_XNAMES=$(echo "${line}" | grep -E -o "${NODE_XNAME_REGEX}" | tr "\n" " " | sed 's/[[:space:]]*$//')
                # we're interested when the nodes have powered on
                if [[ "${LINE_STATE}" == "on" ]] ; then
                    XNAME_NOT_ON_FLAG=false
                    # verify that all test nodes are on
                    for XNAME in ${TEST_NODE_XNAMES} ; do
                        XNAME_ON_CHECK=$(echo "${LINE_NODE_XNAMES}" | grep -E -w "${XNAME}")
                        if [[ -z "${XNAME_ON_CHECK}" ]] ; then
                            XNAME_NOT_ON_FLAG=true
                            echo "${XNAME} is not 'on' yet..."
                        else
                            echo "${XNAME} is 'on'..."
                        fi
                    done
                    echo
                    if ${XNAME_NOT_ON_FLAG} ; then
                        XNAME_ON_COMPLETE=false
                        sleep ${XNAME_ON_ITER_SLEEP_TIME}
                        ((XNAME_ON_ITER++))
                        break
                    else
                        XNAME_ON_COMPLETE=true
                        ((XNAME_ON_ITER++))
                        break
                    fi
                fi
            done <<< "${GET_XNAME_STATUS_STATES}"
        else
            echo "${GET_XNAME_STATUS_STATES}"
            echo "no nodes are 'on' yet..."
            echo
            XNAME_ON_COMPLETE=false
            sleep ${XNAME_ON_ITER_SLEEP_TIME}
            ((XNAME_ON_ITER++))
        fi

    done

    if ${XNAME_ON_COMPLETE} ; then
        echo "PASS: '${GET_XNAME_STATUS_CMD}' powered on all test nodes after ${XNAME_ON_ITER} poll iterations"
        RUN_POWER_ON_CLEANUP=false
    else
        echo "FAIL: '${GET_XNAME_STATUS_CMD}' failed to power on all test nodes after ${XNAME_ON_ITER} poll iterations"
        exit 1
    fi
    echo
}

# wait4boot
wait4boot()
{
    timestamp_print "entering wait4boot()..."
    echo
    for XNAME in ${TEST_NODE_XNAMES} ; do
        echo "${XNAME}"
        XNAME_SSHABLE=false
        XNAME_SSHABLE_ITER=0
        while [[ ${XNAME_SSHABLE} == "false" ]] && [[ ${XNAME_SSHABLE_ITER} -lt ${XNAME_BOOT_ITER_MAX} ]]  ; do
            SSH_OUT=$(ssh ${XNAME} hostname 2> /dev/null)
            SSH_RET=$?
            if [[ ${SSH_RET} -eq 0 ]] ; then
                timestamp_print "${XNAME} is booted"
                XNAME_SSHABLE=true
                continue
            else
                timestamp_print "${XNAME} is not booted..."
                sleep ${XNAME_BOOT_ITER_SLEEP_TIME}
                ((XNAME_SSHABLE_ITER++))
            fi
        done
        echo
        RUN_WAIT4BOOT_CLEANUP=false
        if [[ ${XNAME_SSHABLE} == "false" ]] ; then
            echo "FAIL: timed out waiting for ${XNAME} to boot"
            exit 1
        fi
    done
    echo "PASS: all test nodes rebooted successfully"
    echo
}

# slurm_workaround
slurm_workaround()
{
    echo "entering slurm_workaround()..."
    WORKAROUND_CMD="crayctl configure-managed-plane"
    timestamp_print "running '${WORKAROUND_CMD}'..."
    WORKAROUND_OUT=$(eval ${WORKAROUND_CMD})
    WORKAROUND_RET=$?
    RUN_SLURM_WORKAROUND_CLEANUP=false
    if [[ ${WORKAROUND_RET} -eq 0 ]] ; then
        echo "PASS: DST-1866 workaround '${WORKAROUND_CMD}' returned success code: ${WORKAROUND_RET}"
    else
        echo "${WORKAROUND_OUT}"
        echo "FAIL: DST-1866 workaround '${WORKAROUND_CMD}' failed with error code: ${WORKAROUND_RET}, expected: 0"
        exit 1
    fi
}

trap "cleanup; exit 1" SIGHUP SIGINT SIGTERM

# default options
VERBOSE=false
XNAME_OFF_ITER_SLEEP_TIME=5
XNAME_OFF_ITER_MAX=6
XNAME_ON_ITER_SLEEP_TIME=5
XNAME_ON_ITER_MAX=6
XNAME_BOOT_ITER_SLEEP_TIME=20
XNAME_BOOT_ITER_MAX=30
RUN_POWER_ON_CLEANUP=false
RUN_WAIT4BOOT_CLEANUP=false
RUN_SLURM_WORKAROUND_CLEANUP=false

# parse command-line options
while getopts "hds:i:" opt; do
        case ${opt} in
                h) echo
                   echo "Microservice Cray CAPMC XName Test" 
                   echo "----------------------------------"
                   echo "Automated test for verifying CAPMC xname"
                   echo "control APIs on Cray Shasta systems"
                   echo
                   echo "Usage: capmc_node_xname_power_cycle_test [-h] [-d] [-s <num>] [-i <num>]"
                   echo
                   echo "Arguments:"
                   echo "        -h         display this help message"
                   echo "        -d         enable detailed debugging output"
                   echo "        -s <num>   wait4boot polling iteration sleep time [default is ${XNAME_BOOT_ITER_SLEEP_TIME}s]"
                   echo "        -i <num>   maximum number of wait4boot polling iterations [default is ${XNAME_BOOT_ITER_MAX}]"
                   echo
                   exit 0
                   ;;
                d) VERBOSE=true
                   ;;
                s) XNAME_BOOT_ITER_SLEEP_TIME=$(echo "${OPTARG}" | grep -E "^[1-9]([0-9]+)?$")
                   if [[ -z ${XNAME_BOOT_ITER_SLEEP_TIME} ]] ; then
                        >&2 echo "invalid polling iteration sleep time specified: ${OPTARG}"
                        >&2 echo "argument must be a positive integer"
                        exit 1
                   fi
                   ;;
                i) XNAME_BOOT_ITER_MAX=$(echo "${OPTARG}" | grep -E "^[1-9]([0-9]+)?$")
                   if [[ -z ${XNAME_BOOT_ITER_MAX} ]] ; then
                        >&2 echo "invalid number of polling iterations specified: ${OPTARG}"
                        >&2 echo "argument must be a positive integer"
                        exit 1
                   fi
                   ;;
               \?) echoerr_exit "invalid argument: use -h for help" "1"
                   ;;
        esac
done

# initialize test variables
HOST="mgmt-plane-nmn.local"
PORT="30443"
CAPMC_URI_PATH="apis/capmc/capmc"
NODE_XNAME_REGEX="x[0-9]+c[0-9]+s[0-9]+b[0-9]+n[0-9]+"
TEST_RUN_TIMESTAMP=$(date +"%Y%m%dT%H%M%S")
TEST_NODE_XNAMES=""

# set path to expected location of HMS CT test shared library
HMS_CT_TEST_LIB="/opt/cray/tests/ncn-resources/hms/hms-test/hms_smoke_test_lib_ncn-resources.sh"

# source HMS CT test library file
if [[ -r ${HMS_CT_TEST_LIB} ]] ; then
    . ${HMS_CT_TEST_LIB}
else
    echo "FAIL: failed to source HMS CT test library: ${HMS_CT_TEST_LIB}"
    exit 1
fi

# retrieve Keycloak authentication token for session
TOKEN=$(get_auth_access_token)
TOKEN_RET=$?
if [[ ${TOKEN_RET} -ne 0 ]] ; then
    # get_auth_access_token errors are printed to stderr
    exit 1
fi

####################
# get_xname_status #
####################

GET_XNAME_STATUS_CMD="curl -k -X POST -i -H 'Authorization: Bearer ${TOKEN}' -d '{}' https://${HOST}:${PORT}/${CAPMC_URI_PATH}/get_xname_status"
timestamp_print "running '${GET_XNAME_STATUS_CMD}'..."
GET_XNAME_STATUS_OUT=$(eval ${GET_XNAME_STATUS_CMD} 2> /dev/null)
#GET_XNAME_STATUS_OUT=$(cat /tmp/schooler/fake)
GET_XNAME_STATUS_RET=$?
debug_print "${GET_XNAME_STATUS_OUT}"
if [[ ${GET_XNAME_STATUS_RET} -eq 0 ]] ; then
    echo "PASS: '${GET_XNAME_STATUS_CMD}' returned success code: ${GET_XNAME_STATUS_RET}"
else
    echo "FAIL: '${GET_XNAME_STATUS_CMD}' failed with error code: ${GET_XNAME_STATUS_RET}, expected: 0'"
    exit 1
fi
GET_XNAME_STATUS_HTTP_STATUS_CHECK=$(echo "${GET_XNAME_STATUS_OUT}" | grep -E "200 OK")
if [[ -n "${GET_XNAME_STATUS_HTTP_STATUS_CHECK}" ]] ; then
    echo "PASS: '${GET_XNAME_STATUS_CMD}' returned \"200 OK\""
else
    echo "${GET_XNAME_STATUS_OUT}"
    echo "FAIL: '${GET_XNAME_STATUS_CMD}' failed to return \"200 OK\""
    exit 1
fi

# flatten JSON response
GET_XNAME_STATUS_OUT_JSON=$(echo "${GET_XNAME_STATUS_OUT}" | sed -n -e "/{/,/}/p" | sed 's/^[[:space:]]*//g' | tr -d "\n")
if [[ -z ${GET_XNAME_STATUS_OUT_JSON} ]] ; then
    echo "${GET_XNAME_STATUS_OUT}"
    echo "FAIL: '${GET_XNAME_STATUS_CMD}' failed to return reponse json"
    exit 1
fi

GET_XNAME_STATUS_E=$(echo "${GET_XNAME_STATUS_OUT_JSON}" | grep -E -o "\"e\":(-)?[0-9]+" | grep -E -o "(-)?[0-9]+")
if [[ ${GET_XNAME_STATUS_E} -eq 0 ]] ; then
    echo "PASS: '${GET_XNAME_STATUS_CMD}' returned \"e\" code: ${GET_XNAME_STATUS_E}"
else
    echo "${GET_XNAME_STATUS_OUT}"
    echo "FAIL: '${GET_XNAME_STATUS_CMD}' failed with \"e\" code: ${GET_XNAME_STATUS_E}, expected: 0"
    exit 1
fi
GET_XNAME_STATUS_ERR_MSG=$(echo "${GET_XNAME_STATUS_OUT_JSON}" | grep -E -o "\"err_msg\":\"\"")
if [[ -n ${GET_XNAME_STATUS_ERR_MSG} ]] ; then
    echo "PASS: '${GET_XNAME_STATUS_CMD}' returned \"err_msg\":\"\""
else
    echo "${GET_XNAME_STATUS_OUT}"
    echo "FAIL: '${GET_XNAME_STATUS_CMD}' returned non-empty err_msg"
    exit 1
fi
echo

echo "verifying xname status json..."
echo
GET_XNAME_STATUS_STATES=$(echo "${GET_XNAME_STATUS_OUT_JSON}" | grep -E -o "\"[a-zA-Z0-9]+\":\[(\"[a-zA-Z0-9]+\")(,\"[a-zA-Z0-9]+\")*\]")
while read -r line ; do
    echo "${line}"
    LINE_STATE=$(echo "${line}" | cut -d ":" -f 1 | grep -E -o "\".*\"" | tr -d "\"")
    LINE_NODE_XNAMES=$(echo "${line}" | grep -E -o "${NODE_XNAME_REGEX}" | tr "\n" " " | sed 's/[[:space:]]*$//')
    if [[ -z "${LINE_NODE_XNAMES}" ]] ; then
        # no node xnames listed for this state, continue to look for some
        continue
    fi
    # collect the node xnames returned by get_xname_status to use as test nodes
    if [[ -z ${TEST_NODE_XNAMES} ]] ; then
        TEST_NODE_XNAMES=${LINE_NODE_XNAMES}
    else
        TEST_NODE_XNAMES="${TEST_NODE_XNAMES} ${LINE_NODE_XNAMES}"
    fi
    LINE_INVALID_STATE_FLAG=false
    # verify that valid states are returned
    LINE_STATE_CHECK=$(echo "${LINE_STATE}" | grep -E "on|off|halt|standby|ready|diag|disabled")
    if [[ -n ${LINE_STATE_CHECK} ]] ; then
        for XNAME in ${LINE_NODE_XNAMES} ; do
            echo "PASS: '${GET_XNAME_STATUS_CMD}' returned node xname=${XNAME} in valid state: ${LINE_STATE}"
        done
    else
        for XNAME in ${LINE_NODE_XNAMES} ; do
            echo "FAIL: '${GET_XNAME_STATUS_CMD}' returned node xname=${XNAME} in invalid state: ${LINE_STATE}"
            LINE_INVALID_STATE_FLAG=true
        done
    fi
    if ${LINE_INVALID_STATE_FLAG} ; then
        exit 1
    fi
    echo
done <<< "${GET_XNAME_STATUS_STATES}"

debug_print "TEST_NODE_XNAMES=${TEST_NODE_XNAMES}"
if [[ -z "${TEST_NODE_XNAMES}" ]] ; then
    echo "FAIL: could not find any node xnames to test"
    exit 1
fi
TEST_NODE_XNAMES_COMMAS=$(echo "${TEST_NODE_XNAMES}" | tr " " ",")
echo "test node xnames will be: ${TEST_NODE_XNAMES_COMMAS}..."
echo

#############
# xname_off #
#############

# generate json payload for xname_off operation
XNAME_OFF_REQUEST_BODY_XNAMES_ARG=""
for XNAME in ${TEST_NODE_XNAMES} ; do
    XNAME_TOKEN="\"${XNAME}\","
    XNAME_OFF_REQUEST_BODY_XNAMES_ARG="${XNAME_OFF_REQUEST_BODY_XNAMES_ARG}${XNAME_TOKEN}"
    debug_print "${XNAME_OFF_REQUEST_BODY_XNAMES_ARG}"
done
XNAME_OFF_REQUEST_BODY_XNAMES_ARG=$(echo "${XNAME_OFF_REQUEST_BODY_XNAMES_ARG}" | sed 's/,$//')
debug_print "${XNAME_OFF_REQUEST_BODY_XNAMES_ARG}"
XNAME_OFF_REQUEST_BODY_JSON="{\"reason\":\"capmc_node_xname_power_cycle_test.${TEST_RUN_TIMESTAMP}\",\"xnames\":[${XNAME_OFF_REQUEST_BODY_XNAMES_ARG}]}"
debug_print "${XNAME_OFF_REQUEST_BODY_JSON}"

# call xname_off API
XNAME_OFF_CMD="curl -k -X POST -i -H 'Authorization: Bearer ${TOKEN}' -d '${XNAME_OFF_REQUEST_BODY_JSON}' https://${HOST}:${PORT}/${CAPMC_URI_PATH}/xname_off"
timestamp_print "running '${XNAME_OFF_CMD}'..."
XNAME_OFF_OUT=$(eval ${XNAME_OFF_CMD} 2> /dev/null)
XNAME_OFF_RET=$?
debug_print "${XNAME_OFF_OUT}"
if [[ ${XNAME_OFF_RET} -eq 0 ]] ; then
    echo "PASS: '${XNAME_OFF_CMD}' returned success code: ${XNAME_OFF_RET}"
else
    echo "FAIL: '${XNAME_OFF_CMD}' failed with error code: ${XNAME_OFF_RET}, expected: 0'"
    exit 1
fi
XNAME_OFF_HTTP_STATUS_CHECK=$(echo "${XNAME_OFF_OUT}" | grep -E "200 OK")
if [[ -n "${XNAME_OFF_HTTP_STATUS_CHECK}" ]] ; then
    echo "PASS: '${XNAME_OFF_CMD}' returned \"200 OK\""
else
    echo "${XNAME_OFF_OUT}"
    echo "FAIL: '${XNAME_OFF_CMD}' failed to return \"200 OK\""
    exit 1
fi
RUN_POWER_ON_CLEANUP=true
RUN_WAIT4BOOT_CLEANUP=true
RUN_SLURM_WORKAROUND_CLEANUP=true

XNAME_OFF_ITER=0
XNAME_OFF_COMPLETE=false
# poll until all test nodes have been powered off
echo "polling until all test nodes have powered off..."
echo
while [[ ${XNAME_OFF_COMPLETE} == "false" ]] && [[ ${XNAME_OFF_ITER} -lt ${XNAME_OFF_ITER_MAX} ]] ; do
    
    echo "ITER=${XNAME_OFF_ITER}"
    GET_XNAME_STATUS_CMD="curl -k -X POST -i -H 'Authorization: Bearer ${TOKEN}' -d '{}' https://${HOST}:${PORT}/${CAPMC_URI_PATH}/get_xname_status"
    #GET_XNAME_STATUS_CMD="cat /tmp/schooler/fake"
    timestamp_print "running '${GET_XNAME_STATUS_CMD}'..."
    GET_XNAME_STATUS_OUT=$(eval ${GET_XNAME_STATUS_CMD} 2> /dev/null)
    GET_XNAME_STATUS_RET=$?
    debug_print "${GET_XNAME_STATUS_OUT}"
    if [[ ${GET_XNAME_STATUS_RET} -ne 0 ]] ; then
        echo "FAIL: '${GET_XNAME_STATUS_CMD}' failed with error code: ${GET_XNAME_STATUS_RET}, expected: 0'"
        cleanup
        exit 1
    fi
    GET_XNAME_STATUS_HTTP_STATUS_CHECK=$(echo "${GET_XNAME_STATUS_OUT}" | grep -E "200 OK")
    if [[ -z "${GET_XNAME_STATUS_HTTP_STATUS_CHECK}" ]] ; then
        echo "${GET_XNAME_STATUS_OUT}"
        echo "FAIL: '${GET_XNAME_STATUS_CMD}' failed to return \"200 OK\""
        cleanup
        exit 1
    fi

    # flatten JSON response
    GET_XNAME_STATUS_OUT_JSON=$(echo "${GET_XNAME_STATUS_OUT}" | sed -n -e "/{/,/}/p" | sed 's/^[[:space:]]*//g' | tr -d "\n")
    if [[ -z ${GET_XNAME_STATUS_OUT_JSON} ]] ; then
        echo "${GET_XNAME_STATUS_OUT}"
        echo "FAIL: '${GET_XNAME_STATUS_CMD}' failed to return reponse json"
        cleanup
        exit 1
    fi

    GET_XNAME_STATUS_STATES=$(echo "${GET_XNAME_STATUS_OUT_JSON}" | grep -E -o "\"[a-zA-Z0-9]+\":\[(\"[a-zA-Z0-9]+\")(,\"[a-zA-Z0-9]+\")*\]")
    GET_XNAME_STATUS_STATES_OFF_CHECK=$(echo "${GET_XNAME_STATUS_STATES}" | grep -E "\"off\":\[")
    if [[ -n ${GET_XNAME_STATUS_STATES_OFF_CHECK} ]] ; then
        # parse the get_xname_status output by state key
        while read -r line ; do
            echo "${line}"
            LINE_STATE=$(echo "${line}" | cut -d ":" -f 1 | grep -E -o "\".*\"" | tr -d "\"")
            LINE_NODE_XNAMES=$(echo "${line}" | grep -E -o "${NODE_XNAME_REGEX}" | tr "\n" " " | sed 's/[[:space:]]*$//')
            # we're interested when the nodes have powered off
            if [[ "${LINE_STATE}" == "off" ]] ; then
                XNAME_NOT_OFF_FLAG=false
                # verify that all test nodes are off
                for XNAME in ${TEST_NODE_XNAMES} ; do
                    XNAME_OFF_CHECK=$(echo "${LINE_NODE_XNAMES}" | grep -E -w "${XNAME}")
                    if [[ -z "${XNAME_OFF_CHECK}" ]] ; then
                        XNAME_NOT_OFF_FLAG=true
                        echo "${XNAME} is not 'off' yet..."
                    else
                        echo "${XNAME} is 'off'..."
                    fi
                done
                echo
                if ${XNAME_NOT_OFF_FLAG} ; then
                    XNAME_OFF_COMPLETE=false
                    sleep ${XNAME_OFF_ITER_SLEEP_TIME}
                    ((XNAME_OFF_ITER++))
                    break
                else
                    XNAME_OFF_COMPLETE=true
                    ((XNAME_OFF_ITER++))
                    break
                fi
            fi
        done <<< "${GET_XNAME_STATUS_STATES}"
    else
        echo "${GET_XNAME_STATUS_STATES}"
        echo "no nodes are 'off' yet..."
        echo
        XNAME_OFF_COMPLETE=false
        sleep ${XNAME_OFF_ITER_SLEEP_TIME}
        ((XNAME_OFF_ITER++))
    fi

done

if ${XNAME_OFF_COMPLETE} ; then
    echo "PASS: '${GET_XNAME_STATUS_CMD}' powered off all test nodes after ${XNAME_OFF_ITER} poll iterations"
else
    echo "FAIL: '${GET_XNAME_STATUS_CMD}' failed to power off all test nodes after ${XNAME_OFF_ITER} poll iterations"
    cleanup
    exit 1
fi
echo

############
# xname_on #
############

run_verify_xname_on

#############
# wait4boot #
#############

wait4boot

#######################
# DST-1866 workaround #
#######################

slurm_workaround

