#!/bin/bash

set -xe

export PATH="/root/go/bin:$GOPATH/bin:$PATH"

shasum /root/go/bin/*

/root/go/bin/go get -d $GO_PKG
cd $GOPATH/src && source /root/setup-reproducible

cd $GO_PKG
make ci

