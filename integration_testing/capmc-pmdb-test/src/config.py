#!/usr/bin/env python3
# Copyright Cray Inc 2019
import logging

# Base configuration for reaching the capmc and hsm service
CAPMC_URL = ""
CAPMC_BASE_PATH = ""
HSM_URL = ""
HSM_BASE_PATH = ""

# Initialize the url for communication with capmc
def capmc_init(init_CAPMC_URL, init_CAPMC_BASE_PATH):
    global CAPMC_URL
    global CAPMC_BASE_PATH

    CAPMC_URL = init_CAPMC_URL
    CAPMC_BASE_PATH = init_CAPMC_BASE_PATH

    con = {}
    con["CAPMC_URL"] = CAPMC_URL
    con["CAPMC_BASE_PATH"] = CAPMC_BASE_PATH
    logging.info("Configuring CAPMC connection: %s", con)

# Get the base url for communication with capmc
def get_capmc_url() :
    return CAPMC_URL + CAPMC_BASE_PATH

# Initialize the communication url for the hsm
def hsm_init(init_HSM_URL, init_HSM_BASE_PATH):
    global HSM_URL
    global HSM_BASE_PATH

    HSM_URL = init_HSM_URL
    HSM_BASE_PATH = init_HSM_BASE_PATH

    con = {}
    con["HSM_URL"] = HSM_URL
    con["HSM_BASE_PATH"] = HSM_BASE_PATH
    logging.info("Configuring HSM connection: %s", con)

# Get the base url for communication with the hsm
def get_hsm_url() :
    return HSM_URL + HSM_BASE_PATH
