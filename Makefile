NAME ?= cray-capmc
VERSION ?= $(shell cat .version)

# Helm Chart
CHART_PATH ?= kubernetes
CHART_NAME ?= cray-hms-capmc
CHART_VERSION ?= $(shell cat .version)

all : image chart unittest integration

image:
	docker build ${NO_CACHE} --pull ${DOCKER_ARGS} --tag '${NAME}:${VERSION}' .

chart:
	helm repo add cray-algol60 https://artifactory.algol60.net/artifactory/csm-helm-charts
	helm dep up ${CHART_PATH}/${CHART_NAME}
	helm package ${CHART_PATH}/${CHART_NAME} -d ${CHART_PATH}/.packaged --version ${CHART_VERSION}

unittest:
	./runUnitTest.sh

integration:
	./runIntegration.sh

