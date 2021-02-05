#!/usr/bin/env python3

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

#  MIT License
#
#
#  Permission is hereby granted, free of charge, to any person obtaining a
#  copy of this software and associated documentation files (the "Software"),
#  to deal in the Software without restriction, including without limitation
#  the rights to use, copy, modify, merge, publish, distribute, sublicense,
#  and/or sell copies of the Software, and to permit persons to whom the
#  Software is furnished to do so, subject to the following conditions:
#
#
import requests
import logging
import collections
import inspect
from src.config import get_capmc_url, get_hsm_url

# Common defines:
headerData = {
    'content-type'  : 'application/json',
    'accept'        : 'application/json',
}

# global dictionary of all the tests - add to it using register_test
all_tests = {}

# globals to keep track of results
failingTests = []
passingTests = []

# helper function to add a test to the input dictionary
def register_test(test_func):
    global all_tests
    all_tests[test_func.__name__] = test_func

# execute all the tests registered
def run_all_tests():
    for key,value in all_tests.items():
        run_test(value)

# execute a single registered test with the input name
def run_single_test(testName):
    if testName in all_tests:
        run_test(all_tests[testName])
    else:
        global failingTests
        logging.error("Test: {} not a registered test".format(testName))
        failingTests.append(testName)

# Collection of really basic smoke tests to insure testing environment is live
def run_smoke_tests():
    run_test(test_capmc_connection)
    run_test(test_hsm_live)
    run_test(test_hsm_ready)
    #run_test(test_pmdb_present) - not sure how to do this???

# Test that the basic connection to capmc service is live
def test_capmc_connection():
    # check basic system parameters are returned as expected
    url = get_capmc_url() + "get_system_parameters"
    payload = ""
    expectedResults = {"http_rc":200, "e":0,"err_msg":"","power_cap_target":0,"power_threshold":0}
    return execute_request("GET", url, payload, expectedResults)

# Test that the synthetic hsm service is live and ready
def test_hsm_live():
    # only has the one path that is active, so try it
    url = get_hsm_url() + "State/Components"
    r = requests.request("GET", url, data="", headers=headerData)
    return r.status_code != requests.codes.ok

# Test that the synthetic hsm service is live and ready
def test_hsm_ready():
    # only has the one path that is active, so try it
    url = get_hsm_url() + "service/ready"
    payload = ""
    expectedResults = {"http_rc":200, "code":0, "message":"Ready"}
    return execute_request("GET", url, payload, expectedResults)

# Function to run a test and record the results
def run_test(testFunc):
    # global vars to keep track of the passing/failing tests
    global failingTests
    global passingTests

    # run the tests and record the result
    retVal = testFunc()
    if retVal==0:
        passingTests.append(testFunc.__name__)
    else:
        failingTests.append(testFunc.__name__)

# Quick function to find if there are any failing tests
def has_failing_tests():
    global failingTests
    return len(failingTests)>0
    
# Function to display the results of the test run
def report_results():
    # global vars to keep track of the passing/failing tests
    global failingTests
    global passingTests

    # log the failed tests
    exitCode = 0
    if len(failingTests) > 0:
        exitCode = 1
        logging.error("There were {} failed tests".format( len(failingTests) ))
        [logging.error("Failed test: {}".format(i)) for i in failingTests]

    # log the passing tests
    if len(passingTests) > 0:
        logging.info("There were {} passing tests".format(len(passingTests)))
        [logging.info("Passed: {}".format(i)) for i in passingTests]

    return exitCode

# Check actual resonse against expected response
def compare_with_expected(testName, results, expected):
    # NOTE: for now this only checks for exact values, a future enhancement
    #       would be to allow the expected values to either contain a range
    #       or have a tol to allow for slight variations
    
    # inputs are two dictionaries, check the expected contents against
    # the results
    retVal = 0
    for key in expected:
        if key not in results.keys():
            # the http return code is a special case - already processed
            if key!="http_rc":
                # if the required key is not present, log an error
                logging.error("Test:{}, required value:{} not in result".format(testName, key))
                retVal = 1
        #elif isinstance(results[key],collections.Mapping) and isinstance(expected[key],collections.Mapping):
        #    # comparing two dictionaries, recurse adding the key to the test name
        #    if not compare_with_expected(testName+key, results[key], expected[key]):
        #        retVal = 1
        elif results[key]!=expected[key]:
            # value is present, but the result is not correct
            logging.error("Test:{}, for {}\nexpected:{}\n     got:{}".format(testName, key, expected[key], results[key]))
            retVal = 1
    return retVal

# Test the response of capmc to the request
def execute_request(method, url, payload, expResult):
    # Little trick to get the name of the funciton being run
    testName = inspect.stack()[1][3]

    # see if we just want to see the results of the execution
    resultIsEmpty = isinstance(expResult, str) and not expResult
    if resultIsEmpty:
        # make the request and log the return code
        r = requests.request(method, url, data=payload, headers=headerData)
        logging.info("Http return code:{}".format(r.status_code))

        # pull apart and print the json return object
        try:
            jContent = r.json()
            logging.info("Test:{} returned results are:".format(testName))
            for key in jContent:
                print("  {} -> {}".format(key,jContent[key]))
        except ValueError as e:
            logging.info("Test:{} json parsing error:{}".format(testName,e))
        return 0

    # expected result must be a dictionary and contain expected http return code
    resultIsMap = isinstance(expResult,collections.Mapping)
    mapHasRC = resultIsMap and "http_rc" in expResult
    if not mapHasRC:
        logging.error("Test: {} has mis-formed expected results".format(testName))
        for key in expResult:
            print("  {} -> {}".format(key,expResult[key]))
        return 1

    # keep track of the expected return code for the request
    http_rc = expResult["http_rc"]

    # make the request
    retVal = 0
    try:
        # make the http request and check for good reply
        r = requests.request(method, url, data=payload, headers=headerData)
        if r.status_code != http_rc:
            # unexpected http return code from request - log as error
            logging.error("Test: {} http failed: return code: {}, expected: {}".format(testName, r.status_code, http_rc))
            retVal = 1
        else:
            # expected return code, compare the results
            try:
                # pull response as json data
                jContent = r.json()
                
                # compare return to the expected results
                retVal = compare_with_expected(testName, jContent, expResult)
            except ValueError as e:
                logging.error("Test:{} json parsing error:{}".format(testName,e))
                retVal = 1
    except requests.exceptions.RequestException as e:
        # log the request exception
        logging.error("Test:{} Failed to connect, return code: {}".format(testName, e))
        retVal = 1

    # only log the failures here
    if retVal!=0:
        logging.error("FAILED: {}".format(testName))
    return retVal
