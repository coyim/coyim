#!/bin/bash

set -xe

export PATH="/root/go/bin:$GOPATH/bin:$PATH"
export GOPATH="/gopath"

export TZ=UTC
export LC_ALL=C

mkdir -p /gopath/src/github.com/twstrike/coyim
cp -r /src/* /gopath/src/github.com/twstrike/coyim
cp -r /src/.git /gopath/src/github.com/twstrike/coyim

cd /gopath/src/github.com/twstrike/coyim

make build-gui-sanitize-address BUILD_DIR=/builds
