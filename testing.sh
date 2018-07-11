#!/bin/bash

# Remove Containers
docker rm -f $(docker ps -a -q)
rm -rf /tmp/thanksapp-* thanksapp
docker rm -f -v `docker ps -a --no-trunc | grep "thanksapp" | cut -d ' ' -f 1` 2>/dev/null
docker rmi `docker images --no-trunc | grep "thanksapp" | cut -d ' ' -f 1` 2>/dev/null

# Remove Binary
rm thanksapp

# Fire docker
docker-compose up --force-recreate -d

# go build
go build

# Fire binary
./thanksapp