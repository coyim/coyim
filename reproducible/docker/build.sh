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

mkdir -p /gopath/src/github.com/coyim/coyim
cp -r /src/* /gopath/src/github.com/coyim/coyim
cp -r /src/.git /gopath/src/github.com/coyim/coyim
find /gopath/src/github.com/coyim/coyim -type f -print0 | xargs -0 touch --date="$REFERENCE_DATETIME"

cd /gopath/src/github.com/coyim/coyim

export SRCUID; SRCUID=$(stat -c"%u" /src)
export SRCGID; SRCGID=$(stat -c"%g" /src)

make deps
make build-gui BUILD_DIR=/builds

mkdir -p /src/bin
chown "$SRCUID:$SRCGID" /src/bin

cp /builds/coyim /src/bin

chown "$SRCUID:$SRCGID" /src/bin/coyim

export GLIB_VERSION;  GLIB_VERSION=$(pkg-config --modversion glib-2.0 | tr . _ | cut -d '_' -f 1-2)
export GTK_VERSION;   GTK_VERSION=$(pkg-config --modversion gtk+-3.0 | tr . _ | cut -d '_' -f 1-2)
export PANGO_VERSION; PANGO_VERSION=$(pkg-config --modversion pango | tr . _ | cut -d '_' -f 1-2)
export CAIRO_VERSION; CAIRO_VERSION=$(pkg-config --modversion cairo | tr . _ | cut -d '_' -f 1-2)
export GIT_VERSION; GIT_VERSION=$(git rev-parse HEAD)
export TAG_VERSION; TAG_VERSION=$(git tag -l --contains "$GIT_VERSION" | tail -1)
export GO_VERSION;   GO_VERSION=$(go version | egrep  -o 'go[[:digit:]]\.[[:digit:]]+')
export SUM; SUM=$(sha256sum /builds/coyim)

cat <<EOF > /src/bin/build_info
CoyIM buildinfo
Revision: $GIT_VERSION
Tag: $TAG_VERSION

GLib: $GLIB_VERSION
GTK: $GTK_VERSION
Pango: $PANGO_VERSION
Cairo: $CAIRO_VERSION
Go: $GO_VERSION

$SUM
EOF

chown "$SRCUID:$SRCGID" /src/bin/build_info
