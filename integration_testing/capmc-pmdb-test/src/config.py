#!/usr/bin/env python3
# MIT License
#
# (C) Copyright [2019, 2021] Hewlett Packard Enterprise Development LP
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
