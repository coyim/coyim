#!/bin/bash

set -xe

export PATH="/root/go/bin:$GOPATH/bin:$PATH"
export GOPATH="/gopath"
export CC="clang"
export CXX="clang++"
export PKG_CONFIG_PATH="/usr/repackaged/lib/pkgconfig:/usr/repackaged/share/pkgconfig"
export LD_LIBRARY_PATH="/usr/repackaged/lib"
export CFLAGS="-I/usr/repackaged/include -fno-omit-frame-pointer -fPIE -fsanitize=memory -fsanitize-memory-track-origins -fsanitize-recover=all"
export LDFLAGS="-L/usr/repackaged/lib -fsanitize=memory -fsanitize-memory-track-origins -fsanitize-recover=all"
ldconfig
export MSAN_OPTIONS=""

export TZ=UTC
export LC_ALL=C

mkdir -p /gopath/src/github.com/twstrike/coyim
cp -r /src/* /gopath/src/github.com/twstrike/coyim
cp -r /src/.git /gopath/src/github.com/twstrike/coyim

cd /gopath/src/github.com/twstrike/coyim

make build-gui-memory-analyzer BUILD_DIR=/src/bin
