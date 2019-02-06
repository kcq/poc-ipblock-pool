#!/usr/bin/env bash

set -e

SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
SDIR="$( cd -P "$( dirname "$SOURCE" )" && pwd )"
BDIR="$( cd -P "$( dirname "$SOURCE" )/.." && pwd )"

cd $BDIR/cmd
gofmt -l -w -s .
cd $BDIR/internal
gofmt -l -w -s .
cd $BDIR/pkg
gofmt -l -w -s .
