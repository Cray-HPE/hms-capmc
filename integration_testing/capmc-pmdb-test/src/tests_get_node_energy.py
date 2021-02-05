#!/usr/bin/env python3
# MIT License
#
# (C) Copyright [2020-2021] Hewlett Packard Enterprise Development LP
#
# Permission is hereby granted, free of charge, to any person obtaining a
# copy of this software and associated documentation files (the "Software"),
# to deal in the Software without restriction, including without limitation
# the rights to use, copy, modify, merge, publish, distribute, sublicense,
# and/or sell copies of the Software, and to permit persons to whom the
# Software is furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included
# in all copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
# THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
# OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
# ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
# OTHER DEALINGS IN THE SOFTWARE.

from src.config import get_capmc_url
from src.tests_base import register_test, execute_request

# Get the url to query for this set of tests
def get_url():
    return get_capmc_url() + "get_node_energy"

# Run all of the tests that check 'get_node_energy'
# NOTE: the time stamps in the test database run from
#  2020-01-09 10:00:00 to 2020-01-09 10:30:00
def register_tests():
    register_test(test_get_node_energy_basic)
    register_test(test_get_node_energy_outside_window)
    register_test(test_get_node_energy_no_args)
    register_test(test_get_node_energy_no_times)
    register_test(test_get_node_energy_no_nids)
    register_test(test_get_node_energy_by_app_id)
    register_test(test_get_node_energy_by_job_id)
    register_test(test_get_node_energy_by_job_app_nid)
    register_test(test_get_node_energy_get)
    register_test(test_get_node_energy_delete)
    register_test(test_get_node_energy_put)
    register_test(test_get_node_energy_patch)
    register_test(test_get_node_energy_start_end_same)
    register_test(test_get_node_energy_start_end_switched)
    register_test(test_get_node_energy_straddle_window)
    register_test(test_get_node_energy_pre_window)
    register_test(test_get_node_energy_post_window)
    register_test(test_get_node_energy_one_nid)
    register_test(test_get_node_energy_all_mtn_nids)
    register_test(test_get_node_energy_bad_nids)
    register_test(test_get_node_energy_mixed_nids)
    register_test(test_get_node_energy_river)
    register_test(test_get_node_energy_mtn_river_mix)
    register_test(test_get_node_energy_mtn_river_missing)
    register_test(test_get_node_energy_all_mtn_river)

# Test the basic net accumulated node energy query
def test_get_node_energy_basic():
    payload = '{"start_time":"2020-01-09 10:01:00", "end_time":"2020-01-09 10:06:00","nids":[1000,1001,1002,1003]}'
    expectedResults = {"http_rc":200, "e":0, "time":300, "nid_count":4, "nodes":[
        {'nid': 1000, 'energy': 45979},
        {'nid': 1001, 'energy': 50263},
        {'nid': 1002, 'energy': 53717},
        {'nid': 1003, 'energy': 46434}]}
    return execute_request("POST", get_url(), payload, expectedResults)

# Test the basic query outside the window there is data present
def test_get_node_energy_outside_window():
    payload = '{"start_time":"2020-01-09 09:31:00", "end_time":"2020-01-09 09:36:00","nids":[1000,1001,1002,1003]}'
    expectedResults = {"http_rc":200, "e":400, "err_msg":"Error: No data in time window, Start Time: 2020-01-09 09:31:00 +0000 UTC, End Time: 2020-01-09 09:36:00 +0000 UTC"}
    return execute_request("POST", get_url(), payload, expectedResults)

# Test the query with no parameters input
def test_get_node_energy_no_args():
    # must have a set of nodes and a time window to be valid
    payload = '{ }'
    expectedResults = {"http_rc":400, "e":400, "err_msg":"INVALID_ARGUMENTS"}
    return execute_request("POST", get_url(), payload, expectedResults)

# Test with no times present as input
def test_get_node_energy_no_times():
    payload = '{"nids":[1,2,3]}'
    expectedResults = {"http_rc":400, "e":400, "err_msg":"INVALID_ARGUMENTS"}
    return execute_request("POST", get_url(), payload, expectedResults)

# Test with no nids present as input
def test_get_node_energy_no_nids():
    payload = '{"start_time":"2020-01-09 10:01:00", "end_time":"2020-01-09 10:06:00"}'
    expectedResults = {"http_rc":400, "e":400, "err_msg":"INVALID_ARGUMENTS"}
    return execute_request("POST", get_url(), payload, expectedResults)

# Test the application id query - note: not implemented yet
# NOTE: CASMHMS-2622 - functionality not implmented yet
def test_get_node_energy_by_app_id():
    payload = '{"start_time":"2020-01-09 10:01:00", "end_time":"2020-01-09 10:06:00","apid":"8549334"}'
    expectedResults = {"http_rc":501, "e":501, "err_msg":"Argument not supported yet: APID, err: ARGUMENT_SUPPORT_NOT_IMPLEMENTED"}
    return execute_request("POST", get_url(), payload, expectedResults)

# Test the job id query - note: not implemented yet
# NOTE: CASMHMS-2622 - functionality not implmented yet
def test_get_node_energy_by_job_id():
    payload = '{"start_time":"2020-01-09 10:01:00", "end_time":"2020-01-09 10:06:00","job_id":"8549334"}'
    expectedResults = {"http_rc":501, "e":501, "err_msg":"Argument not supported yet: JOBID, err: ARGUMENT_SUPPORT_NOT_IMPLEMENTED"}
    return execute_request("POST", get_url(), payload, expectedResults)

# Test having nids, job id, and app id query - note: not implemented yet
# NOTE: CASMHMS-2622 - functionality not implmented yet
def test_get_node_energy_by_job_app_nid():
    payload = '{"start_time":"2020-01-09 10:01:00", "end_time":"2020-01-09 10:06:00","nids":[1000,1001,1002,1003],"job_id":"8957467", "apid":"847362"}'
    expectedResults = {"http_rc":501, "e":501, "err_msg":"Argument not supported yet: JOBID/APID, err: ARGUMENT_SUPPORT_NOT_IMPLEMENTED"}
    return execute_request("POST", get_url(), payload, expectedResults)

# Test GET method
def test_get_node_energy_get():
    payload = '{"start_time":"2020-01-09 10:01:00", "end_time":"2020-01-09 10:06:00","nids":[1000,1001,1002,1003]}'
    expectedResults = {"http_rc":405, "e":405, "err_msg":"(GET) Not Allowed"}
    return execute_request("GET", get_url(), payload, expectedResults)

# Test DELETE method
def test_get_node_energy_delete():
    payload = '{"start_time":"2020-01-09 10:01:00", "end_time":"2020-01-09 10:06:00","nids":[1000,1001,1002,1003]}'
    expectedResults = {"http_rc":405, "e":405, "err_msg":"(DELETE) Not Allowed"}
    return execute_request("DELETE", get_url(), payload, expectedResults)

# Test PUT method
def test_get_node_energy_put():
    payload = '{"start_time":"2020-01-09 10:01:00", "end_time":"2020-01-09 10:06:00","nids":[1000,1001,1002,1003]}'
    expectedResults = {"http_rc":405, "e":405, "err_msg":"(PUT) Not Allowed"}
    return execute_request("PUT", get_url(), payload, expectedResults)

# Test PATCH method
def test_get_node_energy_patch():
    payload = '{"start_time":"2020-01-09 10:01:00", "end_time":"2020-01-09 10:06:00","nids":[1000,1001,1002,1003]}'
    expectedResults = {"http_rc":405, "e":405, "err_msg":"(PATCH) Not Allowed"}
    return execute_request("PATCH", get_url(), payload, expectedResults)

# Test start time equal to end time
def test_get_node_energy_start_end_same():
    # hysteresis processing will widen this window to 15 sec
    payload = '{"start_time":"2020-01-09 10:01:00", "end_time":"2020-01-09 10:01:00","nids":[1000,1001,1002,1003]}'
    expectedResults = {"http_rc":200, "e": 0, "nid_count" : 4, "time":15,
        "nodes": [{'nid': 1000, 'energy': 2141}, {'nid': 1001, 'energy': 2425}, {'nid': 1002, 'energy': 2259}, {'nid': 1003, 'energy': 2204}]}
    return execute_request("POST", get_url(), payload, expectedResults)

# Test start time after end time
def test_get_node_energy_start_end_switched():
    payload = '{"start_time":"2020-01-09 10:05:00", "end_time":"2020-01-09 10:00:00","nids":[1000,1001,1002,1003]}'
    expectedResults = {"http_rc":400, "e":400, "err_msg":"Inverted time window - start time: 2020-01-09 10:05:00 +0000 UTC, end time: 2020-01-09 10:00:00 +0000 UTC"}
    return execute_request("POST", get_url(), payload, expectedResults)

# Test start time before data, end time after data
def test_get_node_energy_straddle_window():
    payload = '{"start_time":"2020-01-09 09:55:00", "end_time":"2020-01-09 10:45:00","nids":[1000,1001,1002,1003]}'
    expectedResults = {"http_rc":200, "e":0, "time":3000, "nid_count":4, "nodes":[
        {'nid': 1000, 'energy': 272753},
        {'nid': 1001, 'energy': 293317},
        {'nid': 1002, 'energy': 302245},
        {'nid': 1003, 'energy': 276351}]}
    return execute_request("POST", get_url(), payload, expectedResults)

# Test start time before data, end time in data
def test_get_node_energy_pre_window():
    payload = '{"start_time":"2020-01-09 09:55:00", "end_time":"2020-01-09 10:03:00","nids":[1000,1001,1002,1003]}'
    expectedResults = {"http_rc":200, "e":0, "time":480, "nid_count":4, "nodes":[
        {'nid': 1000, 'energy': 27401},
        {'nid': 1001, 'energy': 32330},
        {'nid': 1002, 'energy': 30245},
        {'nid': 1003, 'energy': 26685}]}
    return execute_request("POST", get_url(), payload, expectedResults)

# Test start time in data, end time beyond data
def test_get_node_energy_post_window():
    payload = '{"start_time":"2020-01-09 10:25:00", "end_time":"2020-01-09 10:45:00","nids":[1000,1001,1002,1003]}'
    expectedResults = {"http_rc":200, "e":0, "time":1200, "nid_count":4, "nodes":[
        {'nid': 1000, 'energy': 45526}, {'nid': 1001, 'energy': 48017},
        {'nid': 1002, 'energy': 48236}, {'nid': 1003, 'energy': 45784}]}
    return execute_request("POST", get_url(), payload, expectedResults)

# Test with only one nid
def test_get_node_energy_one_nid():
    payload = '{"start_time":"2020-01-09 10:01:00", "end_time":"2020-01-09 10:06:00","nids":[1000]}'
    expectedResults = {"http_rc":200, "e":0, "time":300, "nid_count":1, "nodes":[{"nid":1000, "energy":45979}]}
    return execute_request("POST", get_url(), payload, expectedResults)

# Test with all nids
def test_get_node_energy_all_mtn_nids():
    payload = '{"start_time":"2020-01-09 10:01:00", "end_time":"2020-01-09 10:06:00","nids":[1000,1001,1002,1003,1004,1005,1006,1007]}'
    expectedResults = {"http_rc":200, "e":0, "time":300, "nid_count":8, "nodes":[
        {'nid': 1000, 'energy': 45979}, {'nid': 1001, 'energy': 50263},
        {'nid': 1002, 'energy': 53717}, {'nid': 1003, 'energy': 46434},
        {'nid': 1004, 'energy': 51645}, {'nid': 1005, 'energy': 49637},
        {'nid': 1006, 'energy': 48935}, {'nid': 1007, 'energy': 57159} ]}
    return execute_request("POST", get_url(), payload, expectedResults)

# Test with only invalid nids
def test_get_node_energy_bad_nids():
    payload = '{"start_time":"2020-01-09 10:01:00", "end_time":"2020-01-09 10:06:00","nids":[5000,4567]}'
    expectedResults = {"http_rc":400, "e":400, "err_msg":"INVALID_ARGUMENTS"}
    # TODO - is this correct?
    return execute_request("POST", get_url(), payload, expectedResults)

# Test with both valid and invalid nids
def test_get_node_energy_mixed_nids():
    payload = '{"start_time":"2020-01-09 10:01:00", "end_time":"2020-01-09 10:06:00","nids":[1000,1001,1002,1003,1004,59998]}'
    expectedResults = {"http_rc":400, "e":400, "err_msg":"INVALID_ARGUMENTS"}
    # TODO - is this correct?
    return execute_request("POST", get_url(), payload, expectedResults)

# Test with valid river nids only
def test_get_node_energy_river():
    payload = '{"start_time":"2020-01-09 10:01:00", "end_time":"2020-01-09 10:06:00","nids":[1,2,3,4]}'
    expectedResults = {"http_rc":200, "e":0, "time":300, "nid_count":4, "nodes":[
        {'nid': 1, 'energy': 48060}, {'nid': 2, 'energy': 45140},
        {'nid': 3, 'energy': 44550}, {'nid': 4, 'energy': 44260}]}
    return execute_request("POST", get_url(), payload, expectedResults)

# Test with valid mountain and river nids
def test_get_node_energy_mtn_river_mix():
    payload = '{"start_time":"2020-01-09 10:01:00", "end_time":"2020-01-09 10:06:00","nids":[1,2,3,4,1001,1002,1003,1004]}'
    expectedResults = {"http_rc":200, "e":0, "time":300, "nid_count":8, "nodes":[
        {'nid': 1001, 'energy': 50263}, {'nid': 1002, 'energy': 53717},
        {'nid': 1003, 'energy': 46434}, {'nid': 1004, 'energy': 51645},
        {'nid': 1, 'energy': 48060}, {'nid': 2, 'energy': 45140},
        {'nid': 3, 'energy': 44550}, {'nid': 4, 'energy': 44260}]}
    return execute_request("POST", get_url(), payload, expectedResults)

# Test with valid mountain and river nids plus a few incorrect nids
# TODO - need to fix with correct expected results when river implemented
def test_get_node_energy_mtn_river_missing():
    payload = '{"start_time":"2020-01-09 10:01:00", "end_time":"2020-01-09 10:06:00","nids":[1,2,3,4,1001,1002,1003,1004,5999,5479]}'
    expectedResults = {"http_rc":400, "e":400, "err_msg":"INVALID_ARGUMENTS"}
    return execute_request("POST", get_url(), payload, expectedResults)

# Test with all mountain and river nodes
def test_get_node_energy_all_mtn_river():
    payload = '{"start_time":"2020-01-09 10:01:00", "end_time":"2020-01-09 10:06:00","nids":[1,2,3,4,1001,1002,1003,1004,1005,1005,1007]}'
    expectedResults = {"http_rc":200, "e":0, "time":300, "nid_count":10, "nodes":[
        {'nid': 1001, 'energy': 50263}, {'nid': 1002, 'energy': 53717},
        {'nid': 1003, 'energy': 46434}, {'nid': 1004, 'energy': 51645},
        {'nid': 1005, 'energy': 49637}, {'nid': 1007, 'energy': 57159},
        {'nid': 1, 'energy': 48060}, {'nid': 2, 'energy': 45140},
        {'nid': 3, 'energy': 44550}, {'nid': 4, 'energy': 44260}]}
    return execute_request("POST", get_url(), payload, expectedResults)
