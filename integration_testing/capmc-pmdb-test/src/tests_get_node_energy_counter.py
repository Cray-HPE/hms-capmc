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
from src.config import get_capmc_url
from src.tests_base import register_test, execute_request

# Get the endpoint to call for this function
def get_url():
    return get_capmc_url() + "get_node_energy_counter"

# Register all of the tests that check 'get_node_energy'
def register_tests():
    register_test(test_get_node_energy_counter_basic)
    register_test(test_get_node_energy_counter_outside_window)
    register_test(test_get_node_energy_counter_some_nids)
    register_test(test_get_node_energy_counter_no_args)
    register_test(test_get_node_energy_counter_no_time)
    register_test(test_get_node_energy_counter_no_nids)
    register_test(test_get_node_energy_counter_by_app_id)
    register_test(test_get_node_energy_counter_by_job_id)
    register_test(test_get_node_energy_counter_by_job_app_nid)
    register_test(test_get_node_energy_counter_get)
    register_test(test_get_node_energy_counter_delete)
    register_test(test_get_node_energy_counter_put)
    register_test(test_get_node_energy_counter_patch)
    register_test(test_get_node_energy_counter_pre_window)
    register_test(test_get_node_energy_counter_post_window)
    register_test(test_get_node_energy_counter_one_nid)
    register_test(test_get_node_energy_counter_all_mtn_nids)
    register_test(test_get_node_energy_counter_bad_nids)
    register_test(test_get_node_energy_counter_mixed_nids)
    register_test(test_get_node_energy_counter_river)
    register_test(test_get_node_energy_counter_mtn_river_mix)
    register_test(test_get_node_energy_counter_mtn_river_missing)
    return 0
    
# Test the basic get node energy counter function
def test_get_node_energy_counter_basic():
    payload = '{"time":"2020-01-09 10:06:00", "nids":[1000,1001,1002,1003]}'
    expectedResults = {"http_rc":200, "e":0, "nid_count": 4, "nodes":[
        {'nid': 1000, 'energy_ctr': 29445535, 'time': '2020-01-09 10:05:59'},
        {'nid': 1001, 'energy_ctr': 29206350, 'time': '2020-01-09 10:05:59'},
        {'nid': 1002, 'energy_ctr': 31714405, 'time': '2020-01-09 10:05:59'},
        {'nid': 1003, 'energy_ctr': 30327135, 'time': '2020-01-09 10:05:59'}] }
    return execute_request("POST", get_url(), payload, expectedResults)

# Test the basic net accumulated node energy query
# NOTE: CASMHMS-2613 - should always return the most recent counter
def test_get_node_energy_counter_outside_window():
    # This time window is outside of the data that is present - should
    # return the closest counter value in the DB
    payload = '{"time":"2020-01-09 09:55:00", "nids":[1000,1001,1002,1003]}'
    expectedResults = {"http_rc":200, "e":400,
    "err_msg":"Error: No data in time window, Start Time: 2020-01-09 09:54:45 +0000 UTC, End Time: 2020-01-09 09:55:00 +0000 UTC"}
    # correct expectedResults = {"http_rc":200, "e":0, "nid_count":3, "nodes": [{},{},{}]}
    return execute_request("POST", get_url(), payload, expectedResults)

# Test call with only some of the nids present
def test_get_node_energy_counter_some_nids():
    payload = '{"time":"2020-01-09 10:01:00", "nids":[1000,1002,1004]}'
    expectedResults = {"http_rc":200, "e":0, "nid_count": 3, "nodes": [
        {'nid': 1000, 'energy_ctr': 29399404, 'time': '2020-01-09 10:00:59'},
        {'nid': 1002, 'energy_ctr': 31660505, 'time': '2020-01-09 10:00:59'},
        {'nid': 1004, 'energy_ctr': 32057630, 'time': '2020-01-09 10:00:59'}]}
    return execute_request("POST", get_url(), payload, expectedResults)

# Test the query with no parameters input
def test_get_node_energy_counter_no_args():
    # must have a set of nodes and a time window to be valid
    payload = '{ }'
    expectedResults = {"http_rc":400, "e":400, "err_msg":"INVALID_ARGUMENTS"}
    return execute_request("POST", get_url(), payload, expectedResults)

# Test with no times present as input
# NOTE: CASMHMS-2613 - should always return the most recent counter
def test_get_node_energy_counter_no_time():
    payload = '{"nids":[1000,1001,1002,1003]}'
    expectedResults = {"http_rc":200, "e":400 }
    # correct expectedResults = {"http_rc":200, "e":0, "nid_count":3, "nodes": [{},{},{}]}
    return execute_request("POST", get_url(), payload, expectedResults)

# Test with no nids present as input
def test_get_node_energy_counter_no_nids():
    payload = '{"time":"2020-01-09 10:01:00"}'
    expectedResults = {"http_rc":400, "e":400, "err_msg":"INVALID_ARGUMENTS"}
    return execute_request("POST", get_url(), payload, expectedResults)

# Test the application id query - note: not implemented yet
# NOTE: CASMHMS-2622 - functionality not implmented yet
def test_get_node_energy_counter_by_app_id():
    payload = '{"time":"2020-01-09 10:01:00", "apid":"8549334"}'
    expectedResults = {"http_rc":501, "e":501, "err_msg":"Argument not supported yet: APID, err: ARGUMENT_SUPPORT_NOT_IMPLEMENTED"}
    return execute_request("POST", get_url(), payload, expectedResults)

# Test the job id query - note: not implemented yet
# NOTE: CASMHMS-2622 - functionality not implmented yet
def test_get_node_energy_counter_by_job_id():
    payload = '{"time":"2020-01-09 10:01:00","job_id":"8549334"}'
    expectedResults = {"http_rc":501, "e":501, "err_msg":"Argument not supported yet: JOBID, err: ARGUMENT_SUPPORT_NOT_IMPLEMENTED"}
    return execute_request("POST", get_url(), payload, expectedResults)

# Test having nids, job id, and app id query - note: not implemented yet
# NOTE: CASMHMS-2622 - functionality not implmented yet
def test_get_node_energy_counter_by_job_app_nid():
    payload = '{"time":"2020-01-09 10:01:00","nids":[1000,1001,1002,1003],"job_id":"8957467", "apid":"847362"}'
    expectedResults = {"http_rc":501, "e":501, "err_msg":"Argument not supported yet: JOBID/APID, err: ARGUMENT_SUPPORT_NOT_IMPLEMENTED"}
    return execute_request("POST", get_url(), payload, expectedResults)

# Test GET method
def test_get_node_energy_counter_get():
    payload = '{"time":"2020-01-09 10:01:00","nids":[1000,1001,1002,1003]}'
    expectedResults = {"http_rc":405, "e":405, "err_msg":"(GET) Not Allowed"}
    return execute_request("GET", get_url(), payload, expectedResults)

# Test DELETE method
def test_get_node_energy_counter_delete():
    payload = '{"time":"2020-01-09 10:01:00","nids":[1000,1001,1002,1003]}'
    expectedResults = {"http_rc":405, "e":405, "err_msg":"(DELETE) Not Allowed"}
    return execute_request("DELETE", get_url(), payload, expectedResults)

# Test PUT method
def test_get_node_energy_counter_put():
    payload = '{"time":"2020-01-09 10:01:00", "nids":[1000,1001,1002,1003]}'
    expectedResults = {"http_rc":405, "e":405, "err_msg":"(PUT) Not Allowed"}
    return execute_request("PUT", get_url(), payload, expectedResults)

# Test PATCH method
def test_get_node_energy_counter_patch():
    payload = '{"time":"2020-01-09 10:01:00","nids":[1000,1001,1002,1003]}'
    expectedResults = {"http_rc":405, "e":405, "err_msg":"(PATCH) Not Allowed"}
    return execute_request("PATCH", get_url(), payload, expectedResults)

# Test time before data
# NOTE: CASMHMS-2613 - should always return the most recent counter
def test_get_node_energy_counter_pre_window():
    payload = '{"time":"2020-01-09 09:55:00","nids":[1000,1001,1002,1003]}'
    expectedResults = {"http_rc":200, "e":400, "err_msg":"Error: No data in time window, Start Time: 2020-01-09 09:54:45 +0000 UTC, End Time: 2020-01-09 09:55:00 +0000 UTC"}
    return execute_request("POST", get_url(), payload, expectedResults)

# Test time beyond data
# NOTE: CASMHMS-2613 - should always return the most recent counter
def test_get_node_energy_counter_post_window():
    payload = '{"time":"2020-01-09 11:30:00","nids":[1000,1001,1002,1003]}'
    expectedResults = {"http_rc":200, "e":400, "err_msg":"Error: No data in time window, Start Time: 2020-01-09 11:29:45 +0000 UTC, End Time: 2020-01-09 11:30:00 +0000 UTC"}
    return execute_request("POST", get_url(), payload, expectedResults)

# Test with only one nid
def test_get_node_energy_counter_one_nid():
    payload = '{"time":"2020-01-09 10:01:00","nids":[1000]}'
    expectedResults = {"http_rc":200, "e":0, "nid_count":1,
        "nodes" : [{'nid': 1000, 'energy_ctr': 29399404, 'time': '2020-01-09 10:00:59'}]}
    return execute_request("POST", get_url(), payload, expectedResults)

# Test with all nids
def test_get_node_energy_counter_all_mtn_nids():
    payload = '{"time":"2020-01-09 10:01:00","nids":[1000,1001,1002,1003,1004,1005,1006,1007]}'
    expectedResults = {"http_rc":200, "e":0, "nid_count":8, "nodes":[
        {'nid': 1000, 'energy_ctr': 29399404, 'time': '2020-01-09 10:00:59'},
        {'nid': 1001, 'energy_ctr': 29155908, 'time': '2020-01-09 10:00:59'},
        {'nid': 1002, 'energy_ctr': 31660505, 'time': '2020-01-09 10:00:59'},
        {'nid': 1003, 'energy_ctr': 30280548, 'time': '2020-01-09 10:00:59'},
        {'nid': 1004, 'energy_ctr': 32057630, 'time': '2020-01-09 10:00:59'},
        {'nid': 1005, 'energy_ctr': 31292793, 'time': '2020-01-09 10:00:59'},
        {'nid': 1006, 'energy_ctr': 33558771, 'time': '2020-01-09 10:00:59'},
        {'nid': 1007, 'energy_ctr': 33096233, 'time': '2020-01-09 10:00:59'}] }
    return execute_request("POST", get_url(), payload, expectedResults)

# Test with only invalid nids
def test_get_node_energy_counter_bad_nids():
    payload = '{"time":"2020-01-09 10:01:00","nids":[5000,4567]}'
    expectedResults = {"http_rc":400, "e":400, "err_msg":"INVALID_ARGUMENTS"}
    # TODO - is this correct?
    return execute_request("POST", get_url(), payload, expectedResults)

# Test with both valid and invalid nids
def test_get_node_energy_counter_mixed_nids():
    payload = '{"time":"2020-01-09 10:01:00","nids":[1000,1001,1002,1003,59998]}'
    expectedResults = {"http_rc":400, "e":400, "err_msg":"INVALID_ARGUMENTS"}
    # TODO - is this correct?
    return execute_request("POST", get_url(), payload, expectedResults)

# Test with valid river nids only
# TODO - need to fix with correct expected results when river implemented
def test_get_node_energy_counter_river():
    payload = '{"time":"2020-01-09 10:01:00","nids":[1,2,3,4]}'
    expectedResults = {"http_rc":200, "e":400, "err_msg":"Error: No data in time window, Start Time: 2020-01-09 10:00:45 +0000 UTC, End Time: 2020-01-09 10:01:00 +0000 UTC"}
    return execute_request("POST", get_url(), payload, expectedResults)

# Test with valid mountain and river nids
# TODO - need to fix with correct expected results when river implemented
# NOTE - the below result is only the mountain nids - correct will include river
def test_get_node_energy_counter_mtn_river_mix():
    payload = '{"time":"2020-01-09 10:01:00","nids":[1,2,3,4,1001,1002,1003,1004]}'
    expectedResults = {"http_rc":200, "e":0, "nid_count":4, "nodes":[
        {'nid': 1001, 'energy_ctr': 29155908, 'time': '2020-01-09 10:00:59'},
        {'nid': 1002, 'energy_ctr': 31660505, 'time': '2020-01-09 10:00:59'},
        {'nid': 1003, 'energy_ctr': 30280548, 'time': '2020-01-09 10:00:59'},
        {'nid': 1004, 'energy_ctr': 32057630, 'time': '2020-01-09 10:00:59'}]}
    return execute_request("POST", get_url(), payload, expectedResults)

# Test with valid mountain and river nids plus a few incorrect nids
# TODO - need to fix with correct expected results when river implemented
def test_get_node_energy_counter_mtn_river_missing():
    payload = '{"time":"2020-01-09 10:01:00","nids":[1,2,3,4,1001,1002,1003,1004,5999,5479]}'
    expectedResults = {"http_rc":400, "e":400, "err_msg":"INVALID_ARGUMENTS"}
    return execute_request("POST", get_url(), payload, expectedResults)
