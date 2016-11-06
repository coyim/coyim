#!/bin/bash

set -xe

export PATH="/root/go/bin:$GOPATH/bin:$PATH"
export GOPATH="/gopath"

# Does nothing if REFERENCE_DATETIME is missing
test -z "$REFERENCE_DATETIME" && return 0

export LD_PRELOAD=/usr/lib/x86_64-linux-gnu/faketime/libfaketime.so.1
export FAKETIME=$REFERENCE_DATETIME
export TZ=UTC
export LC_ALL=C

mkdir -p /gopath/src/github.com/twstrike/coyim
cp -r /src/* /gopath/src/github.com/twstrike/coyim
cp -r /src/.git /gopath/src/github.com/twstrike/coyim
find /gopath/src/github.com/twstrike/coyim -type f -print0 | xargs -0 touch --date="$REFERENCE_DATETIME"

cd /gopath/src/github.com/twstrike/coyim

export SRCUID=`stat -c"%u" /src`
export SRCGID=`stat -c"%g" /src`

make build-cli BUILD_DIR=/builds
make build-gui BUILD_DIR=/builds

mkdir -p /src/bin
chown $SRCUID:$SRCGID /src/bin

cp /builds/coyim-cli /src/bin
cp /builds/coyim /src/bin

chown $SRCUID:$SRCGID /src/bin/coyim-cli
chown $SRCUID:$SRCGID /src/bin/coyim

export GTK_VERSION=`pkg-config --modversion gtk+-3.0 | tr . _ | cut -d '_' -f 1-2`
export GIT_VERSION=`git rev-parse HEAD`
export TAG_VERSION=`git tag -l --contains $GIT_VERSION | tail -1`
export GO_VERSION=$(go version | grep  -o 'go[[:digit:]]\.[[:digit:]]')
export SUM1=`sha256sum /builds/coyim-cli`
export SUM2=`sha256sum /builds/coyim`


cat <<EOF > /src/bin/build_info
CoyIM buildinfo
Revision: $GIT_VERSION
Tag: $TAG_VERSION

GTK: $GTK_VERSION
Go: $GO_VERSION

$SUM1
$SUM2
EOF

chown $SRCUID:$SRCGID /src/bin/build_info
