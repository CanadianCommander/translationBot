#!/bin/bash

pushd $(dirname $0)/../

# Dont need anything fancy for just one simple service. If I add a DB, then dev in local cluster w/ skaffold.
go mod vendor
DOCKER_BUILDKIT=1 docker build . -t trasnlation-bot:latest
docker run -p 8080:8080 --name translation-bot-dev trasnlation-bot:latest
docker container rm translation-bot-dev
docker image rm trasnlation-bot:latest

popd