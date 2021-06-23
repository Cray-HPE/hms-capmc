NAME ?= hms-capmc 
VERSION ?= $(shell cat .version)

all : image unittest coverage

image:
		docker build --pull ${DOCKER_ARGS} --tag '${NAME}:${VERSION}' .

unittest: buildbase
		docker build --no-cache -t cray/hms-capmc-testing -f Dockerfile.testing .

coverage: buildbase
		docker build -t cray/hms-capmc-coverage -f Dockerfile.coverage .

buildbase: 
		docker build -t cray/hms-capmc-build-base -f Dockerfile.build-base .
		
