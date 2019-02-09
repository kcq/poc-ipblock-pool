#!/usr/bin/env bash

set -e

SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
SDIR="$( cd -P "$( dirname "$SOURCE" )" && pwd )"
BDIR="$( cd -P "$( dirname "$SOURCE" )/.." && pwd )"

if ! which golint > /dev/null; then
    echo "No golint. Installing...."
    go get -u github.com/golang/lint/golint
fi

source $SDIR/env.sh
cd $BDIR/cmd
go tool vet .
golint ./...
cd $BDIR/internal
go tool vet .
golint ./...
cd $BDIR/pkg
go tool vet .
golint ./...
