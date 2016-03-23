#!/bin/bash

set -xe

export PATH="/root/go/bin:$GOPATH/bin:$PATH"

# Get package and setup a reproducible environment
/root/get-reproducibly $GO_PKG

cd ${GOPATH}/src/${GO_PKG}
ls -l
make -C ci

