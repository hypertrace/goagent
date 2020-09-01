#!/bin/bash

# This script compares the container ID obtained from docker ps and a test
# app writen in golang using the GetContainerID function.

docker rmi --force goagent-test
docker rm goagent-test
docker build --no-cache -f tests/docker/Dockerfile -t goagent-test .

# ACTUAL_CONTAINER_ID is obtained from the golang application running inside the
# docker file
ACTUAL_CONTAINER_ID=$(docker run --name goagent-test goagent-test)

# EXPECTED_SHORTEN_CONTAINER_ID is obtained from `docker ps`
EXPECTED_SHORTEN_CONTAINER_ID=$(docker ps -q -a --filter "name=goagent-test")

SHORTEN_LENGTH=${#EXPECTED_SHORTEN_CONTAINER_ID} # the length of the short format
ACTUAL_SHORTEN_CONTAINER_ID="${ACTUAL_CONTAINER_ID:0:$SHORTEN_LENGTH}"

if [[ -z "$EXPECTED_SHORTEN_CONTAINER_ID" ]]; then
    echo "Failed to obtain the container ID, empty value"
    exit 1
fi

echo ""
if [[ "$EXPECTED_SHORTEN_CONTAINER_ID" == "$ACTUAL_SHORTEN_CONTAINER_ID" ]]; then
    echo "Container ID successfully obtained."
    exit 0
else 
    echo -n "Failed to obtain the container ID, expected something starting"
    echo -n " with \"$EXPECTED_SHORTEN_CONTAINER_ID\""
    echo " got \"$ACTUAL_CONTAINER_ID\"."
    exit 1
fi
