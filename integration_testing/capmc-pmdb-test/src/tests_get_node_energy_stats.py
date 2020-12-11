#!/usr/bin/env python3
# Copyright Cray Inc 2020

from src.config import get_capmc_url
from src.tests_base import register_test, execute_request

# Get the url to make this function call
def get_url():
    return get_capmc_url() + "get_node_energy_stats"

# Register all of the tests that check 'get_node_energy'
def register_tests():
    register_test(test_get_node_energy_stats_basic)
    register_test(test_get_node_energy_stats_outside_window)
    register_test(test_get_node_energy_stats_some_nids)
    register_test(test_get_node_energy_stats_no_args)
    register_test(test_get_node_energy_stats_no_times)
    register_test(test_get_node_energy_stats_no_nids)
    register_test(test_get_node_energy_stats_by_app_id)
    register_test(test_get_node_energy_stats_by_job_id)
    register_test(test_get_node_energy_stats_by_job_app_nid)
    register_test(test_get_node_energy_stats_get)
    register_test(test_get_node_energy_stats_delete)
    register_test(test_get_node_energy_stats_put)
    register_test(test_get_node_energy_stats_patch)
    register_test(test_get_node_energy_stats_start_end_same)
    register_test(test_get_node_energy_stats_start_end_switched)
    register_test(test_get_node_energy_stats_straddle_window)
    register_test(test_get_node_energy_stats_pre_window)
    register_test(test_get_node_energy_stats_post_window)
    register_test(test_get_node_energy_stats_one_nid)
    register_test(test_get_node_energy_stats_all_mtn_nids)
    register_test(test_get_node_energy_stats_bad_nids)
    register_test(test_get_node_energy_stats_mixed_nids)
    register_test(test_get_node_energy_stats_river)
    register_test(test_get_node_energy_stats_mtn_river_mix)
    register_test(test_get_node_energy_stats_mtn_river_missing)

    return 0
    
# Test the basic net accumulated node energy query
def test_get_node_energy_stats_basic():
    payload = '{"start_time":"2020-01-09 10:00:00", "end_time":"2020-01-09 10:05:00","nids":[1000,1001,1002,1003,1004]}'
    expectedResults = {"http_rc":200, "e":0, "time":300, "nid_count":5, "energy_total":244623, "energy_avg":48924.6, "energy_std":3187.223211511864, "energy_max": [1004, 51823], "energy_min": [1003, 45360]}
    return execute_request("POST", get_url(), payload, expectedResults)

# Test the basic net accumulated node energy query
def test_get_node_energy_stats_outside_window():
    # This time window is outside of the data that is present
    payload = '{"start_time":"2020-01-09 09:00:00", "end_time":"2020-01-09 09:30:00","nids":[1000,1001,1002,1003,1004]}'
    expectedResults = {"http_rc":200, "e":400, "err_msg":"Error: No data in time window, Start Time: 2020-01-09 09:00:00 +0000 UTC, End Time: 2020-01-09 09:30:00 +0000 UTC"}
    return execute_request("POST", get_url(), payload, expectedResults)

# Test call with only some of the nids present
def test_get_node_energy_stats_some_nids():
    payload = '{"start_time":"2020-01-09 10:00:00", "end_time":"2020-01-09 10:05:00","nids":[1000,1002,1004]}'
    expectedResults = {"http_rc":200, "e":0, "time":300, "nid_count":3, "energy_total":148181, "energy_avg":49393.666666666664, "energy_std":3362.495254023912, "energy_max": [1004, 51823], "energy_min": [1000, 45556]}
    return execute_request("POST", get_url(), payload, expectedResults)

# Test the query with no parameters input
def test_get_node_energy_stats_no_args():
    # must have a set of nodes and a time window to be valid
    payload = '{ }'
    expectedResults = {"http_rc":400, "e":400, "err_msg":"INVALID_ARGUMENTS"}
    return execute_request("POST", get_url(), payload, expectedResults)

# Test with no times present as input
def test_get_node_energy_stats_no_times():
    payload = '{"nids":[1000,1002,1003]}'
    expectedResults = {"http_rc":400, "e":400, "err_msg":"INVALID_ARGUMENTS"}
    return execute_request("POST", get_url(), payload, expectedResults)

# Test with no nids present as input
def test_get_node_energy_stats_no_nids():
    payload = '{"start_time":"2020-01-09 10:00:00", "end_time":"2020-01-09 10:30:00"}'
    expectedResults = {"http_rc":400, "e":400, "err_msg":"INVALID_ARGUMENTS"}
    return execute_request("POST", get_url(), payload, expectedResults)

# Test the application id query - note: not implemented yet
# NOTE: CASMHMS-2622 - functionality not implmented yet
def test_get_node_energy_stats_by_app_id():
    payload = '{"start_time":"2020-01-09 10:00:00", "end_time":"2020-01-09 10:30:00","apid":"8549334"}'
    expectedResults = {"http_rc":501, "e":501, "err_msg":"Argument not supported yet: APID, err: ARGUMENT_SUPPORT_NOT_IMPLEMENTED"}
    return execute_request("POST", get_url(), payload, expectedResults)

# Test the job id query - note: not implemented yet
# NOTE: CASMHMS-2622 - functionality not implmented yet
def test_get_node_energy_stats_by_job_id():
    payload = '{"start_time":"2020-01-09 10:00:00", "end_time":"2020-01-09 10:30:00","job_id":"8549334"}'
    expectedResults = {"http_rc":501, "e":501, "err_msg":"Argument not supported yet: JOBID, err: ARGUMENT_SUPPORT_NOT_IMPLEMENTED"}
    return execute_request("POST", get_url(), payload, expectedResults)

# Test having nids, job id, and app id query - note: not implemented yet
# NOTE: CASMHMS-2622 - functionality not implmented yet
def test_get_node_energy_stats_by_job_app_nid():
    payload = '{"start_time":"2020-01-09 10:00:00", "end_time":"2020-01-09 10:30:00","nids":[1000,1002,1003],"job_id":"8957467", "apid":"847362"}'
    expectedResults = {"http_rc":501, "e":501, "err_msg":"Argument not supported yet: JOBID/APID, err: ARGUMENT_SUPPORT_NOT_IMPLEMENTED"}
    return execute_request("POST", get_url(), payload, expectedResults)

# Test GET method
def test_get_node_energy_stats_get():
    payload = '{"start_time":"2020-01-09 10:00:00", "end_time":"2020-01-09 10:30:00","nids":[1000,1002,1003]}'
    expectedResults = {"http_rc":405, "e":405, "err_msg":"(GET) Not Allowed"}
    return execute_request("GET", get_url(), payload, expectedResults)

# Test DELETE method
def test_get_node_energy_stats_delete():
    payload = '{"start_time":"2020-01-09 10:00:00", "end_time":"2020-01-09 10:30:00","nids":[1000,1002,1003]}'
    expectedResults = {"http_rc":405, "e":405, "err_msg":"(DELETE) Not Allowed"}
    return execute_request("DELETE", get_url(), payload, expectedResults)

# Test PUT method
def test_get_node_energy_stats_put():
    payload = '{"start_time":"2020-01-09 10:00:00", "end_time":"2020-01-09 10:30:00","nids":[1000,1002,1003]}'
    expectedResults = {"http_rc":405, "e":405, "err_msg":"(PUT) Not Allowed"}
    return execute_request("PUT", get_url(), payload, expectedResults)

# Test PATCH method
def test_get_node_energy_stats_patch():
    payload = '{"start_time":"2020-01-09 10:00:00", "end_time":"2020-01-09 10:30:00","nids":[1000,1002,1003]}'
    expectedResults = {"http_rc":405, "e":405, "err_msg":"(PATCH) Not Allowed"}
    return execute_request("PATCH", get_url(), payload, expectedResults)

# Test start time equal to end time
def test_get_node_energy_stats_start_end_same():
    # hysteresis processing will widen this window to 15 sec
    payload = '{"start_time":"2020-01-09 10:30:00", "end_time":"2020-01-09 10:30:00","nids":[1000,1002,1003]}'
    expectedResults = {"http_rc":200, "e": 0, "nid_count" : 3, "time":15, "energy_total":6799,
        "energy_avg": 2266.3333333333335, "energy_std":118.196164630386, "energy_max": [1003, 2340], "energy_min" : [1000, 2130]}
    return execute_request("POST", get_url(), payload, expectedResults)

# Test start time after end time
def test_get_node_energy_stats_start_end_switched():
    payload = '{"start_time":"2020-01-09 10:05:00", "end_time":"2020-01-09 10:00:00","nids":[1000,1002,1003]}'
    expectedResults = {"http_rc":400, "e":400, "err_msg":"Inverted time window - start time: 2020-01-09 10:05:00 +0000 UTC, end time: 2020-01-09 10:00:00 +0000 UTC"}
    return execute_request("POST", get_url(), payload, expectedResults)

# Test start time before data, end time after data
def test_get_node_energy_stats_straddle_window():
    payload = '{"start_time":"2020-01-09 09:55:00", "end_time":"2020-01-09 10:45:00","nids":[1000,1002,1003]}'
    expectedResults = {"http_rc":200, "e":0, "time":3000, "nid_count":3, "energy_total":851349,
      "energy_avg":283783, "energy_std":16089.45256993, "energy_max":[1002, 302245], "energy_min":[1000, 272753]}
    return execute_request("POST", get_url(), payload, expectedResults)

# Test start time before data, end time in data
def test_get_node_energy_stats_pre_window():
    payload = '{"start_time":"2020-01-09 09:55:00", "end_time":"2020-01-09 10:15:00","nids":[1000,1002,1003]}'
    expectedResults = {"http_rc":200, "e":0, "time":1200, "nid_count":3, "energy_total":423286, "energy_avg":141095.33333333334,
      "energy_std":7957.71332314336, "energy_max":[1002, 150246], "energy_min":[1000, 135796]}
    return execute_request("POST", get_url(), payload, expectedResults)

# Test start time in data, end time beyond data
def test_get_node_energy_stats_post_window():
    payload = '{"start_time":"2020-01-09 10:15:00", "end_time":"2020-01-09 10:45:00","nids":[1000,1002,1003]}'
    expectedResults = {"http_rc":200, "e":0, "time":1800, "nid_count":3, "energy_total":427487, "energy_avg":142495.66666666666,
      "energy_std":8155.563397665016, "energy_max":[1002, 151834], "energy_min":[1000, 136773]}
    return execute_request("POST", get_url(), payload, expectedResults)

# Test with only one nid
def test_get_node_energy_stats_one_nid():
    payload = '{"start_time":"2020-01-09 10:00:00", "end_time":"2020-01-09 10:15:00","nids":[1000]}'
    expectedResults = {"http_rc":200, "e":0, "time":900, "nid_count":1, "energy_total":135796, "energy_avg":135796,
      "energy_std":0, "energy_max":[1000, 135796],"energy_min":[1000, 135796]}
    return execute_request("POST", get_url(), payload, expectedResults)

# Test with all nids
def test_get_node_energy_stats_all_mtn_nids():
    payload = '{"start_time":"2020-01-09 10:00:00", "end_time":"2020-01-09 10:30:00","nids":[1000,1001,1002,1003,1004,1005,1006,1007]}'
    expectedResults = {"http_rc":200, "e":0, "time":1800, "nid_count":8, "energy_total": 2422640,
      "energy_avg": 302830, "energy_std":22690.30672839, "energy_max":[1005, 339963], "energy_min":[1000, 272753]}
    return execute_request("POST", get_url(), payload, expectedResults)

# Test with only invalid nids
def test_get_node_energy_stats_bad_nids():
    payload = '{"start_time":"2020-01-09 10:00:00", "end_time":"2020-01-09 10:30:00","nids":[5000,4567]}'
    expectedResults = {"http_rc":400, "e":400, "err_msg":"INVALID_ARGUMENTS"}
    # TODO - is this correct?
    return execute_request("POST", get_url(), payload, expectedResults)

# Test with both valid and invalid nids
def test_get_node_energy_stats_mixed_nids():
    payload = '{"start_time":"2020-01-09 10:00:00", "end_time":"2020-01-09 10:30:00","nids":[1000,1001,1002,1003,1004,59998]}'
    expectedResults = {"http_rc":400, "e":400, "err_msg":"INVALID_ARGUMENTS"}
    # TODO - is this correct?
    return execute_request("POST", get_url(), payload, expectedResults)

# Test with valid river nids only
def test_get_node_energy_stats_river():
    payload = '{"start_time":"2020-01-09 10:00:00", "end_time":"2020-01-09 10:30:00","nids":[1,2,3,4]}'
    expectedResults = {"http_rc":200, "e":0,"time":1800, "nid_count":4, "energy_total": 1093060,
      "energy_avg": 273265, "energy_std":10763.179517843848, "energy_max":[1, 289340], "energy_min":[3, 266930]}
    return execute_request("POST", get_url(), payload, expectedResults)

# Test with valid mountain and river nids
def test_get_node_energy_stats_mtn_river_mix():
    payload = '{"start_time":"2020-01-09 10:00:00", "end_time":"2020-01-09 10:30:00","nids":[1,2,3,4,1001,1002,1003,1004]}'
    expectedResults = {"http_rc":200, "e":0,"time":1800, "nid_count":8, "energy_total": 2277176,
      "energy_avg": 284647, "energy_std":12990.879914161924, "energy_max":[1004, 312203], "energy_min":[3, 266930]}
    return execute_request("POST", get_url(), payload, expectedResults)

# Test with valid mountain and river nids plus a few incorrect nids
def test_get_node_energy_stats_mtn_river_missing():
    payload = '{"start_time":"2020-01-09 10:00:00", "end_time":"2020-01-09 10:30:00","nids":[1,2,3,4,1001,1002,1003,1004,5999,5479]}'
    expectedResults = {"http_rc":400, "e":400, "err_msg":"INVALID_ARGUMENTS"}
    return execute_request("POST", get_url(), payload, expectedResults)
