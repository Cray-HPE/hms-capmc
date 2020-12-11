#!/usr/bin/env bash

# This script is run each time a PR is created or updated.  If it fails,
# the build will fail.

# Generate a 'random' tag for the services created for this test
RANDY=$(echo $RANDOM | md5sum | awk '{print $1}')
CURWD=$(pwd)

# Run the capmc-pmdb integration tests
echo "Running integration tests..."
cd ${CURWD}/integration_testing
./run.sh -s ${RANDY}
RESULT=$?

# Return to current dir
cd ${CURWD}

# Report results
if [ ${RESULT} -eq 0 ]
then
  echo "PASS: Successfully completed integration tests"
  exit 0
else
  echo "FAIL: Failure in integration tests" >&2
  exit 1
fi
