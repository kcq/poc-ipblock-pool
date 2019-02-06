#!/usr/bin/env bash

set -e

SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
SDIR="$( cd -P "$( dirname "$SOURCE" )" && pwd )"
BDIR="$( cd -P "$( dirname "$SOURCE" )/.." && pwd )"

echo "Preparing local GOPATH..."

pushd ${BDIR}
rm -rf _gopath
mkdir _gopath
pushd ${BDIR}/_gopath
mkdir -p src/github.com/kcq
ln -sf $BDIR src/github.com/kcq/poc-ipblock-pool
popd
popd