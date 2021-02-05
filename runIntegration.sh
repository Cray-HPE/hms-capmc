#!/usr/bin/env bash
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
