# MIT License
#
# (C) Copyright [2019-2022] Hewlett Packard Enterprise Development LP
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

# Dockerfile for building hms-capmc.

# Build base just has the packages installed we need.
FROM artifactory.algol60.net/docker.io/library/golang:1.16-alpine AS build-base

RUN set -eux \
    && apk -U upgrade \
    && apk add build-base

# Base copies in the files we need to test/build.
FROM build-base AS base

RUN go env -w GO111MODULE=auto

# Copy all the necessary files to the image.
COPY cmd $GOPATH/src/github.com/Cray-HPE/hms-capmc/cmd
COPY internal $GOPATH/src/github.com/Cray-HPE/hms-capmc/internal
COPY vendor $GOPATH/src/github.com/Cray-HPE/hms-capmc/vendor


### Build Stage ###

FROM base AS builder

RUN set -ex && go build -v -i -o /usr/local/bin/capmc-service github.com/Cray-HPE/hms-capmc/cmd/capmcd

### Final Stage ###

FROM artifactory.algol60.net/docker.io/library/alpine:3.15

RUN set -eux \
    && apk -U upgrade

LABEL maintainer="Hewlett Packard Enterprise"
EXPOSE 27777
STOPSIGNAL SIGTERM

# Get the CAPMC service from the builder stage.
# Note: The name used here must match that used in the builder stage.
COPY --from=builder /usr/local/bin/capmc-service /usr/local/bin

COPY files/config.toml /usr/local/etc/capmc-service/default/config.toml

# Setup environment variables.
ENV HSM_URL=https://api-gateway.default.svc.cluster.local/apis/smd
ENV CAPMC_CONFIG=/usr/local/etc/capmc-service/default/config.toml
ENV PCS_URL=https://api-gateway.default.svc.cluster.local/apis/power-control/v1
#ENV CAPMC_CA_URI= #
ENV LOG_LEVEL="INFO"

# Used by the HMS secure storage pkg
ENV VAULT_ADDR="http://cray-vault.vault:8200"
ENV VAULT_SKIP_VERIFY="true"

COPY configs configs

# nobody 65534:65534
USER 65534:65534

# Start the service.
CMD ["sh", "-c", "capmc-service -config=$CAPMC_CONFIG -hsm=$HSM_URL -pcs=$PCS_URL"]
