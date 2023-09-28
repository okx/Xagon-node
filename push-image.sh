#!/bin/bash

export DOCKER_DEFAULT_PLATFORM=linux/amd64

make build-docker
docker tag zkevm-node:latest ${SENTIO_DOCKER_REPO}/zkevm-node:v0.3.1
docker push ${SENTIO_DOCKER_REPO}/zkevm-node:v0.3.1
