#!/usr/bin/env bash

set -e

SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
SDIR="$( cd -P "$( dirname "$SOURCE" )" && pwd )"
BDIR="$( cd -P "$( dirname "$SOURCE" )/.." && pwd )"

echo "Checking binaries..."

if [ ! -d $BDIR/bin ]; then
	echo "Missing binaries (run build first)..."
	exit 1
fi

echo "Building container..."

pushd ${BDIR}
pwd
docker build -t poc/ipblock-pool .
popd
