#!/bin/bash

set -xe

# args: NAME EXTRA_CONFIGURE EXTRA_CFLAGS USE_PIC_CFLAGS DISTCLEAN

BUILD_NAME=$1
EXTRA_CONFIGURE_ARGS=$2
EXTRA_CFLAGS=$3
USE_PIC_CFLAGS=$4
DISTCLEAN=$5

CFLAGS_TO_USE=$STANDARD_CFLAGS

if [ "x$USE_PIC_CFLAGS" == "xtrue" ]
then
    CFLAGS_TO_USE=$PIC_CFLAGS
fi

export CFLAGS="$CFLAGS_TO_USE $EXTRA_CFLAGS"
export CPPFLAGS="$CFLAGS_TO_USE $EXTRA_CFLAGS"
ldconfig

cd /root/deps/$BUILD_NAME*

if [ "x$DISTCLEAN" == "xtrue" ]
then
    make distclean
fi

./configure --prefix=/usr/repackaged $EXTRA_CONFIGURE_ARGS &&
    make &&
    make install
