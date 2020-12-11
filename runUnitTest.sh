#!/usr/bin/env bash

# Copyright 2020 Hewlett Packard Enterprise Development LP


# Build the build base image
docker build -t cray/hms-capmc-build-base -f Dockerfile.build-base .

docker build -t cray/hms-capmc-testing -f Dockerfile.testing .
