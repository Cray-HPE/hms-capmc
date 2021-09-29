NAME ?= cray-capmc
VERSION ?= $(shell cat .version)

# Helm Chart
CHART_PATH ?= kubernetes
CHART_NAME ?= cray-hms-capmc
CHART_VERSION ?= $(shell cat .version)

all : image chart unittest coverage integration

image:
		docker build --pull ${DOCKER_ARGS} --tag '${NAME}:${VERSION}' .

chart:
		helm repo add cray-algol60 https://artifactory.algol60.net/artifactory/csm-helm-charts
		helm dep up ${CHART_PATH}/${CHART_NAME}
		helm package ${CHART_PATH}/${CHART_NAME} -d ${CHART_PATH}/.packaged --version ${CHART_VERSION}

unittest:
		./runUnitTest.sh

coverage:
		./runCoverage.sh

integration:
		./runIntegration.sh
