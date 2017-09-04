#!/bin/bash

set -xe

export CC=
export CXX=
export PKG_CONFIG_PATH=
export CFLAGS=
export CPPFLAGS=
export LDFLAGS=
export LD_LIBRARY_PATH=
ldconfig

# Requires gccgo-X to be installed
export GOROOT_BOOTSTRAP=/usr

cd /root/go

cd src
./make.bash
