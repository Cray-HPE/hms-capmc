#!/bin/bash
# MIT License
#
# (C) Copyright [2021] Hewlett Packard Enterprise Development LP
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

set -x
echo ${DB_HOSTNAME}
psql -h ${DB_HOSTNAME} pmdb pmdbuser -c "\copy pmdb.nc_data from '/csv/nc_data.csv' DELIMITER ',' CSV HEADER;"
psql -h ${DB_HOSTNAME} pmdb pmdbuser -c "\copy pmdb.cc_data from '/csv/cc_data.csv' DELIMITER ',' CSV HEADER;"
psql -h ${DB_HOSTNAME} pmdb pmdbuser -c "\copy pmdb.sc_data from '/csv/sc_data.csv' DELIMITER ',' CSV HEADER;"
psql -h ${DB_HOSTNAME} pmdb pmdbuser -c "\copy pmdb.river_data from '/csv/river_data.csv' DELIMITER ',' CSV HEADER;"
psql -h ${DB_HOSTNAME} pmdb pmdbuser -c "\copy pmdb.device_specific_contexts from '/csv/device_specific_contexts.csv' DELIMITER ',' CSV HEADER;"
psql -h ${DB_HOSTNAME} pmdb pmdbuser -c "\copy pmdb.parental_contexts from '/csv/parental_contexts.csv' DELIMITER ',' CSV HEADER;"
psql -h ${DB_HOSTNAME} pmdb pmdbuser -c "\copy pmdb.physical_contexts from '/csv/physical_contexts.csv' DELIMITER ',' CSV HEADER;"
psql -h ${DB_HOSTNAME} pmdb pmdbuser -c "\copy pmdb.physical_sub_contexts from '/csv/physical_sub_contexts.csv' DELIMITER ',' CSV HEADER;"
psql -h ${DB_HOSTNAME} pmdb pmdbuser -c "\copy pmdb.sensor_types from '/csv/sensor_types.csv' DELIMITER ',' CSV HEADER;"
psql -h ${DB_HOSTNAME} pmdb pmdbuser -c "\copy pmdb.version from '/csv/version.csv' DELIMITER ',' CSV ;"
