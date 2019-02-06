#!/usr/bin/env bash

set -e

SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
SDIR="$( cd -P "$( dirname "$SOURCE" )" && pwd )"
BDIR="$( cd -P "$( dirname "$SOURCE" )/.." && pwd )"

echo "Preparing..."

mkdir -p ${BDIR}/bin/

echo "Building PoC apps..."

pushd ${BDIR}/cmd/ipblock-pool
pwd
GOOS=linux GOARCH=amd64 go build -o "${BDIR}/bin/ipblock-pool"
popd

pushd ${BDIR}/cmd/ipblock-pool-server
pwd
GOOS=linux GOARCH=amd64 go build -o "${BDIR}/bin/ipblock-pool-server"
popd