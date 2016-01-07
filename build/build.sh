#!/bin/bash

set -xe

export PATH="/root/go/bin:$GOPATH/bin:$PATH"

shasum /root/go/bin/*

# Use source code from current volume
# /root/go/bin/go get -d $GO_PKG
mkdir -p $(echo "${GOPATH}/src/${GO_PKG}" | rev | cut -d '/' -f 2- | rev) &&\
  ln -s /src "${GOPATH}/src/${GO_PKG}"

cd $GOPATH/src && source /root/setup-reproducible

# make build
/root/go/bin/go build -o /builds/coyim \
  $GO_PKG

/root/go/bin/go build -o /builds/coyim-gui \
  $GO_PKG

shasum /builds/*

