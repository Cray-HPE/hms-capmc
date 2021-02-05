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

import requests
import inspect
from src.config import get_capmc_url
from src.tests_base import register_test, execute_request

# get the url for use in this test suite
def get_url():
    return get_capmc_url() + "get_system_power"

# Register all the tests in this module to be run
def register_tests():
    register_test(test_get_system_power_basic)
    register_test(test_get_system_power_outside_window)
    register_test(test_get_system_power_after_window)
    register_test(test_get_system_power_no_start)
    register_test(test_get_system_power_no_win_len)
    register_test(test_get_system_power_no_params)
    register_test(test_get_system_power_extra_params)
    register_test(test_get_system_power_put)
    register_test(test_get_system_power_patch)
    register_test(test_get_system_power_delete)
    register_test(test_get_system_power_get_json)
    register_test(test_get_system_power_get)
    register_test(test_get_system_power_long_win)
    register_test(test_get_system_power_short_win)
    register_test(test_get_system_power_neg_win)
    register_test(test_get_system_power_bad_time)

# Test a basic call within the window where there is data
def test_get_system_power_basic():
    payload = '{"start_time":"2020-01-09 10:00:00", "window_len": 300}'
    expectedResults = {"http_rc":200, "e":0, "window_len":300, "start_time":"2020-01-09 10:00:00","avg":4037,"max":4148,"min":3923}
    return execute_request("POST", get_url(), payload, expectedResults)

# Test a basic call outside (before) the window where there is data
def test_get_system_power_outside_window():
    payload = '{"start_time":"2020-01-09 09:30:00", "window_len": 300}'
    expectedResults = {"http_rc":200, "e":400, "err_msg":"Error: No data in time window, Start Time: 2020-01-09 09:30:00 +0000 UTC, End Time: 2020-01-09 09:35:00 +0000 UTC"}
    return execute_request("POST", get_url(), payload, expectedResults)

# Test a basic call after the window where there is data
def test_get_system_power_after_window():
    payload = '{"start_time":"2020-01-09 11:30:00", "window_len": 30}'
    expectedResults = {"http_rc":200, "e":400, "err_msg":"Error: No data in time window, Start Time: 2020-01-09 11:30:00 +0000 UTC, End Time: 2020-01-09 11:30:30 +0000 UTC"}
    return execute_request("POST", get_url(), payload, expectedResults)

# Test with no start time, but a defined window length
def test_get_system_power_no_start():
    # NOTE - with no start time supplied, 'now' is used so error message will
    #  change depending on when test is run.
    payload = '{"window_len": 300}'
    expectedResults = {"http_rc":200, "e":400}
    return execute_request("POST", get_url(), payload, expectedResults)

# Test with a start time, but no defined window length
def test_get_system_power_no_win_len():
    payload = '{"start_time":"2020-01-09 10:00:00"}'
    expectedResults = {"http_rc":200, "e":0, "window_len":15, "start_time":"2020-01-09 10:00:00", "avg":4138, "max":4138, "min":4138}
    return execute_request("POST", get_url(), payload, expectedResults)

# Test with no parameters defined at all
def test_get_system_power_no_params():
    # NOTE - with no start time supplied, 'now' is used so error message will
    #  change depending on when test is run.
    payload = '{ }'
    expectedResults = {"http_rc":200, "e":400}
    return execute_request("POST", get_url(), payload, expectedResults)

# Test with an extra parameter defined
def test_get_system_power_extra_params():
    payload = '{"start_time":"2020-01-09 10:00:00", "window_len": 30, "nonsense":200}'
    expectedResults = {"http_rc":200, "e":0, "window_len":30, "start_time":"2020-01-09 10:00:00", "avg":4086, "max":4138, "min":4035}
    return execute_request("POST", get_url(), payload, expectedResults)

# Test with a PUT method - not allowed
def test_get_system_power_put():
    payload = '{"start_time":"2020-01-09 10:00:00", "window_len": 30}'
    expectedResults = {"http_rc":405, "e":405, "err_msg":"(PUT) Not Allowed"}
    return execute_request("PUT", get_url(), payload, expectedResults)

# Test with a PATCH method - not allowed
def test_get_system_power_patch():
    payload = '{"start_time":"2020-01-09 10:00:00", "window_len": 30}'
    expectedResults = {"http_rc":405, "e":405, "err_msg":"(PATCH) Not Allowed"}
    return execute_request("PATCH", get_url(), payload, expectedResults)

# Test with a DELETE method - not allowed
def test_get_system_power_delete():
    payload = '{"start_time":"2020-01-09 10:00:00", "window_len": 30}'
    expectedResults = {"http_rc":405, "e":405, "err_msg":"(DELETE) Not Allowed"}
    return execute_request("DELETE", get_url(), payload, expectedResults)

# Test a 'GET' call with json data
# NOTE: CASMHMS-2580 - the mixing of 'GET' with json input data/headers should
#  be returning a 'Bad Request' code, not 'Internal Server Error'.
def test_get_system_power_get_json():
    # I expect a 'Bad Request' return since the wrong kind of input is sent for
    #  the request type
    # NOTE: the current return is happening because the mis-read json input is
    #  interpreted as no input arguments (which would be valid) and there is
    #  no data for that window.  Still need better error here.
    payload = '{"start_time":"2020-01-09 10:00:00", "window_len": 300}'
    expectedResults = {"http_rc":200, "e":400 }
    # correct -> expectedResults = {"http_rc":400, "e":400, "err_msg":"Something about json payload with form data expected"}
    return execute_request("GET", get_url(), payload, expectedResults)

# Test the basic function call through the 'GET' method
def test_get_system_power_get():
    url = get_url() + "?window_len=300&start_time=2020-01-09 10:00:00"
    payload = ""
    expectedResults = {"http_rc":200, "e":0, "window_len":300, "start_time":"2020-01-09 10:00:00","avg":4037,"max":4148,"min":3923}
    return execute_request("GET", url, payload, expectedResults)

# Check for a long time window request
def test_get_system_power_long_win():
    payload = '{"start_time":"2020-01-09 10:00:00", "window_len": 3601}'
    expectedResults = {"http_rc":400,"e":400, "err_msg":"WINDOW_LEN_OUT_OF_RANGE"}
    # is this correct???
    return execute_request("POST", get_url(), payload, expectedResults)

# Check for a short time window request
def test_get_system_power_short_win():
    payload = '{"start_time":"2020-01-09 10:00:00", "window_len": 1}'
    expectedResults = {"http_rc":400,"e":400, "err_msg":"WINDOW_LEN_OUT_OF_RANGE"}
    # is this correct???
    return execute_request("POST", get_url(), payload, expectedResults)

# Check for a negative time window request
def test_get_system_power_neg_win():
    payload = '{"start_time":"2020-01-09 10:00:00", "window_len": -50}'
    expectedResults = {"http_rc":400,"e":400, "err_msg":"WINDOW_LEN_OUT_OF_RANGE"}
    # is this correct???
    return execute_request("POST", get_url(), payload, expectedResults)

# Check for a bad start time
def test_get_system_power_bad_time():
    payload = '{"start_time":"01/09/2020 10:00:00", "window_len": 30}'
    expectedResults = {"http_rc":400,"e":400, "err_msg":"BAD_START_TIME"}
    return execute_request("POST", get_url(), payload, expectedResults)

