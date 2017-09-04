#!/bin/bash

set -xe

# args: NAME EXTRA_CFLAGS USE_PIC_CFLAGS

BUILD_NAME=$1
EXTRA_CONFIGURE_ARGS=$2
EXTRA_CFLAGS=$3
USE_PIC_CFLAGS=$4

CFLAGS_TO_USE=$MSAN_CFLAGS

if [ "x$USE_PIC_CFLAGS" == "xtrue" ]
then
    CFLAGS_TO_USE=$MSAN_PIC_CFLAGS
fi

export CFLAGS="$CFLAGS_TO_USE $EXTRA_CFLAGS"
export CPPFLAGS="$CFLAGS_TO_USE $EXTRA_CFLAGS"
export LDFLAGS="$MSAN_LDFLAGS"
ldconfig

cd /root/deps/$BUILD_NAME*

/root/installers/clean-binaries.sh $BUILD_NAME

make &&
    make install
