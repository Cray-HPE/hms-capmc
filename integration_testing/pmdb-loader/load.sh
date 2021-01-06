#!/bin/bash
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