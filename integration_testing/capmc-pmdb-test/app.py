# Copyright Cray Inc 2019

import logging
import argparse
import os
from src.config import capmc_init, hsm_init
import src.tests_base as tests
import src.tests_get_system_power as gsp_tests
import src.tests_get_system_power_details as gspd_tests
import src.tests_get_node_energy as gne_tests
import src.tests_get_node_energy_stats as gnes_tests
import src.tests_get_node_energy_counter as gnec_tests
import src.tests_not_implemented as tni_tests

# The actual tests are contained in the above imported files.  Each python file
# has a 'run_all' function that will run all the of tests contained in that
# file.  The tests_base file has test framework utilities and a few smoke
# tests to insure the system is properly set up to run the tests.

# Prior to running this script, the pmdb database, synthetic hsm, and capmc
# all need to be spun up and the database needs to be loaded with the
# test data from jolt1 from 11/19/2019.  These scripts assume that all
# these servers are available for http calls.

# The telemetry data exists in a time window between:
#   2019-11-19 21:30:00
#   2019-11-19 22:00:00

# NOTE: At this time, there is no way to set/override the 'now' portion of
#  the time based queries.  This means that the hysteresis functionality
#  is not ablet to be testes, since all telemetry data is historical.

# Main loop for the test run
if __name__ == '__main__':
    # There are two possible command-line arguments:
    # --logLevel will get how much information is output by logging
    # --debug will set up the default urls for hsm and capmc to allow
    #         local access to the servers.  Use this is the servers
    #         have been started and are already running.  If this is
    #         not set, the env vars containing urls must be defined.
    # --test will allow the user to run a single test of the given 
    #        name instead of all tests.
    parser = argparse.ArgumentParser()
    parser.add_argument("--logLevel", help="Level of logging to process")
    parser.add_argument("--debug",action='store_true', help="Use default service urls")
    parser.add_argument("--test", help="Name of the single test to run")
    args = parser.parse_args()
    if args.logLevel:
        numeric_log_level = getattr(logging, args.logLevel.upper(), None)
        logging.basicConfig(level=numeric_log_level)
    if args.debug:
        os.environ['CAPMC_URL'] = "http://localhost:27777"
        os.environ['HSM_URL'] = "http://localhost:27779"

    # Load the configuration to enable communication with capmc
    try:
        capmc_init(os.environ['CAPMC_URL'], "/capmc/v1/")
        hsm_init(os.environ['HSM_URL'], "/hsm/v1/")
    except KeyError:
        logging.error("The environment variables CAPMC_URL and HSM_URL must be set")
        exit(1)

    # run the smoke tests - if they don't pass no need to do anything else
    tests.run_smoke_tests()
    if not tests.has_failing_tests():
        # register all the tests
        gsp_tests.register_tests()
        gspd_tests.register_tests()
        gne_tests.register_tests()
        gnes_tests.register_tests()
        gnec_tests.register_tests()
        tni_tests.register_tests()

        # run all the tests or a single input test name
        if args.test:
            tests.run_single_test(args.test)
        else:
            tests.run_all_tests()

    # report the results to the log and exit using the application return code
    exitCode = tests.report_results()
    exit(exitCode)
