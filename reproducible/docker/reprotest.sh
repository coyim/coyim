#!/bin/sh


BUILDDIR='/tmp/reprotest/'

mkdir "$BUILDDIR"
cd "$BUILDDIR" || exit 1

apt-get update
apt-src install coyim
cd "${BUILDDIR}./coyim-0.3.8+ds" || exit 1
reprotest auto .

cd "${BUILDDIR}" || exit 1
pwd
ls -la

# Copy files to mounted volume
[ -d /src/reproducible/docker/.reprotest ] || mkdir /src/reproducible/docker/.reprotest
cp "${BUILDDIR}/*" /src/reproducible/docker/.reprotest

sleep 9000
