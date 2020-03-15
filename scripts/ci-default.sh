#!/bin/sh
# Login docker
echo "${DOCKER_PASSWORD}" | docker login -u "${DOCKER_USERNAME}" --password-stdin
# Build Golang Application
go build -o dist/schedule-microservice
# Build docker image
docker build . -t kainonly/schedule-microservice:latest
# Push docker image
docker push kainonly/schedule-microservice:latest