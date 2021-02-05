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

# NOTE: This tests the following capmc functions:
#    get_system_power
#    get_system_power_details
#    get_node_energy
#    get_node_energy_statistics
#    get_node_energy_counter
#
#  It uses a current version of hms/hms-postgresql and hms/hms-postgresql-util
#  checked out from the git repo, captured telemetry data from jolt1 circa 
#  11/19/2019, and the current version of capmc.  The telemetry data is imported
#  into a running instance of hms-postgresql, a synthetic hsm that can just serve
#  up the information needed for these function calls is created and started,
#  and the current version of capmc is built and started.  The tests are 
#  contained in python scripts and run against the spun-up services.

# Debugging notes:
# There are several command line args that can help with debugging:
#    -s SUFFIX will use the input instead of a large randomly generated
#         string appended to the end of the image/container names.  This
#         helps if leaving services running.
#    -d will leave the services running and images in docker.  They are in
#         a state that capmc may be used.
#    -t will tear down and remove the services / images.  As of right now
#         you must supply the correct suffix for the container names.
#    -r will rebuild capmc (assuming it is already running) with the input
#         suffix applied.

# function to help with cleanup
do_teardown () {
  SUFFIX=$1
  echo "Tearing down and cleaning up... Suffix: ${SUFFIX}"
  docker stop hms-capmc_${SUFFIX}
  docker stop synthetic_hsm_${SUFFIX}
  docker stop hms-postgresql-pmdb_${SUFFIX}
  docker rm capmc-pmdb-tests_${RANDY}
  docker image rm capmc_pmdb_test/capmc_${SUFFIX}
  docker image rm capmc_pmdb_test/synthetic_hsm_${SUFFIX}
  docker image rm cray/hms/hms-postgresql_${SUFFIX}
  docker image rm capmc_pmdb_test/pmdb_loader_${SUFFIX}
  docker image rm cray/hms/hms-postgresql-util_${SUFFIX}
  docker image rm capmc_pmdb_test/capmc_pmdb_script_${SUFFIX}
  docker network rm ${SUFFIX}
}

# create a unique hash to use for things created in this process
RANDY=$(echo $RANDOM | md5sum | awk '{print $1}')
TEARDOWN_ONLY=0
DEBUG_MODE=0
DO_REBUILD=0

# add some input arguments for ease of debugging
# -s SUFFIX : different suffix for service names for ease of debugging
# -d : debug - leave the services running and open localhost ports
# -t : only do teardown with input suffix
# -r : rebuild/redeploy capmc
while getopts s:dtr option
do
case "$option"
in
s) RANDY=${OPTARG};;
d) DEBUG_MODE=1;;
t) TEARDOWN_ONLY=1;;
r) DO_REBUILD=1;;
\?)
  # Bail with something unknown to prevent a LOT of unintended work
  echo "Invalid Option: -$OPTARG" 1>&2
  exit 1
esac
done 


# keep track of the current directory
CURWD=$(pwd)
echo $CURWD

# check for request to just tear down running system
if [ ${TEARDOWN_ONLY} -eq 1 ]
then
  do_teardown ${RANDY} 
  exit 0
fi

# check if we are just rebuilding the redeploying campmc
# NOTE: this is taking all the commands scattered through the
#  rest of the script and condensing them here to tear down, 
#  rebuild, and restart capmc.  If any of those procedures change
#  this will need to be updated as well.
if [ ${DO_REBUILD} -eq 1 ]
then
  CAPMC_LOCAL_PORT="-p 27777:27777"
  cd ${CURWD}/..
  docker stop hms-capmc_${RANDY}
  docker build -t capmc_pmdb_test/capmc_${RANDY} .
  docker run -d --rm --name hms-capmc_${RANDY} -e DB_HOSTNAME=hms-postgresql-pmdb_${RANDY} -e HSM_URL=http://synthetic_hsm_${RANDY}:27779 -e LOG_LEVEL=Trace --network ${RANDY} ${CAPMC_LOCAL_PORT} capmc_pmdb_test/capmc_${RANDY}
  cd ${CURWD}
  exit 0
fi

# if running in debug mode, configure to open localhost ports
DO_CLEANUP=1
HSM_LOCAL_PORT=""
CAPMC_LOCAL_PORT=""
PSQL_LOCAL_PORT=""
if [ ${DEBUG_MODE} -eq 1 ]
then
  DO_CLEANUP=0
  PSQL_LOCAL_PORT="-p 5432:5432"
  HSM_LOCAL_PORT="-p 27779:27779"
  CAPMC_LOCAL_PORT="-p 27777:27777"
fi

echo "Running capmc-pmdb integration tests"

# Fail on error and print executions
set -x

# Set up branches to try and build against
CURBRANCH=$(git branch | grep \* | cut -d ' ' -f2)
echo $CURBRANCH

BRANCH_HIERARCHY=(
    ${CURBRANCH}
    develop
    master
)

REPO_DIR=$(mktemp -d)
REPOS=(
    hms/hms-postgresql-util
    hms/hms-postgresql
)

#Step 1) get the source code of the other containers
echo "Checking out all repos"

for repo in ${REPOS[@]} ; do
    echo "Cloning $repo into $REPO_DIR/$repo"
    git clone --depth 1 --no-single-branch https://stash.us.cray.com/scm/"$repo".git "${REPO_DIR}"/"${repo}"
done

echo "trying to checkout feature branch"
for repo in ${REPOS[@]} ; do
    echo "cd into $REPO_DIR/$repo"
    cd $REPO_DIR/$repo
    echo $(pwd)

    #step 2) make sure on the right branch
    for x in ${BRANCH_HIERARCHY[@]} ; do
        echo "attempting to checkout branch ${x}"
        git checkout ${x}
        if [ $? -eq 0 ]
        then
          echo "successfully checked out branch ${x}"
          break
        else
          echo "could not find branch ${x}..." >&2
            if [ "${x}" == "master" ]; then
                echo "all out of options... exiting"
                exit 1
            fi
        fi
    done

    #step 3) build that source into container image
    echo "loaded correct branch, proceeding to build"
    docker build -t cray/${repo}_${RANDY} -f Dockerfile .
done

 #step 4) stand up the new image with the network, wait for db to be done.

#need to set the hostname on these container!
docker network create ${RANDY}
docker run -d --rm --name hms-postgresql-pmdb_${RANDY} -e POSTGRES_HOST_AUTH_METHOD=trust --network ${RANDY} ${PSQL_LOCAL_PORT} cray/hms/hms-postgresql_${RANDY}
docker run -d --rm --name hms-postgresql-util_${RANDY} -e PMDB_HOSTNAME=hms-postgresql-pmdb_${RANDY} --network ${RANDY} cray/hms/hms-postgresql-util_${RANDY}
docker wait hms-postgresql-util_${RANDY}

#step 5) need to apply csv files to database!
cd ${CURWD}/pmdb-loader
docker build -t capmc_pmdb_test/pmdb_loader_${RANDY} .
docker run -d --rm --name hms-pmdb-loader_${RANDY} -e DB_HOSTNAME=hms-postgresql-pmdb_${RANDY} --network ${RANDY} capmc_pmdb_test/pmdb_loader_${RANDY}
#docker run -d --name hms-pmdb-loader_${RANDY} -e DB_HOSTNAME=hms-postgresql-pmdb_${RANDY} --network ${RANDY} capmc_pmdb_test/pmdb_loader_${RANDY}

#step 6) need to stand up synthetic hsm
cd ${CURWD}/synthetic-hsm
docker build -t capmc_pmdb_test/synthetic_hsm_${RANDY} .
docker run -d --rm --name synthetic_hsm_${RANDY} --network ${RANDY} ${HSM_LOCAL_PORT} capmc_pmdb_test/synthetic_hsm_${RANDY}

#step 7) need to stand up capmc
cd ${CURWD}/..
docker build -t capmc_pmdb_test/capmc_${RANDY} .
docker run -d --rm --name hms-capmc_${RANDY} -e DB_HOSTNAME=hms-postgresql-pmdb_${RANDY} -e HSM_URL=http://synthetic_hsm_${RANDY}:27779 -e LOG_LEVEL=Trace --network ${RANDY} ${CAPMC_LOCAL_PORT} capmc_pmdb_test/capmc_${RANDY}

#step 8) run the testing script, wait to complete, and pull the exit code
#  NOTE: when running in debug mode do not execute tests, just set up the
#  env for them to run locally
RESULT=0
if [ ${DEBUG_MODE} -eq 0 ]
then
  # NOTE: Need to wait for capmc to finish initializing before hitting it with
  #  tests or they will fail.  Without readiness probe, will just sleep a little...
  sleep 5

  # run the test from inside another docker image - wait for it to complete
  cd ${CURWD}/capmc-pmdb-test
  docker build -t capmc_pmdb_test/capmc_pmdb_script_${RANDY} .
  docker run --name capmc-pmdb-tests_${RANDY} -e HSM_URL=http://synthetic_hsm_${RANDY}:27779 -e CAPMC_URL=http://hms-capmc_${RANDY}:27777 --network ${RANDY} capmc_pmdb_test/capmc_pmdb_script_${RANDY}
  docker wait capmc-pmdb-tests_${RANDY}
  RESULT=$(docker inspect capmc-pmdb-tests_${RANDY} --format='{{.State.ExitCode}}')
fi

#step 9) Tear it all down and clean up the fingerprints
if [ ${DO_CLEANUP} -eq 1 ]
then
  do_teardown ${RANDY}
else
  echo "Requested no cleanup: ${DO_CLEANUP}"
fi

# NOTE: always remove the random dir where other repos were checked out
rm -rf ${REPO_DIR}

#step 10) Report results
if [ ${RESULT} -eq 0 ]
then
  echo "PASS: Successfully completed capmc-pmdb integration tests"
  exit 0
else
  echo "FAIL: Failing tests in capmc-pmdb integration tests" >&2
  exit 1
fi
