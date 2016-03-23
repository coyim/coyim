#!/bin/bash

set -xe

PACKAGE=${1}

export PATH="/root/go/bin:$GOPATH/bin:$PATH"

# Build golang
#/root/download-golang && /root/build-golang
shasum /root/go/bin/*

# Get packages at a specific revision
/root/go/bin/go get -d $PACKAGE
cd ${GOPATH}/src/${PACKAGE}
git checkout $REVISION

# Alternatively, we could use source code from current dir if the volume is mounted
#mkdir -p $(echo "${GOPATH}/src/${PACKAGE}" | rev | cut -d '/' -f 2- | rev) &&\
#  ln -s /src "${GOPATH}/src/${PACKAGE}"

# Fetch dependencies (this should not build dependencies)
make deps

# Prepare the reproducible environment
cd $GOPATH/src && source /root/setup-reproducible

cd ${GOPATH}/src/${PACKAGE}

