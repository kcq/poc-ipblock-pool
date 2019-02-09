#!/usr/bin/env bash

SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
SDIR="$( cd -P "$( dirname "$SOURCE" )" && pwd )"
BDIR="$( cd -P "$( dirname "$SOURCE" )/.." && pwd )"

echo "Setting up env..."

if [ -z "${GOPATH}" ]; then
	echo "GOPATH: Not set yet"
	if [ -d $BDIR/_gopath ]; then
		echo "GOPATH: Using local _gopath dir..."
		export GOPATH=$BDIR/_gopath
	else
		echo "GOPATH: Using default value..."
		export GOPATH=$HOME/go
	fi
else
	echo "GOPATH: Already set..."
fi

export PATH=$PATH:$GOPATH/bin

echo "GOPATH: ${GOPATH}"
