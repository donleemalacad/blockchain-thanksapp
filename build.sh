#!/bin/bash
set -e

# Copy All Binaries from /binary/ubuntu/ to Current directory
cp ./binary/ubuntu/* .

# So that configtxgen tool will be able to locate the configtx.yaml file
export FABRIC_CFG_PATH=${PWD}

# Run Fabric Network Setup
sh ./fabric.sh

# Start Network
sh ./startNetwork.sh

# Set Interval
#sleep 5

# Fire up docker
#docker-compose up -d