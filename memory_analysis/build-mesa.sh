#!/bin/bash

set -xe

export PATH=/usr/repackaged/bin:$PATH
export LD_LIBRARY_PATH=/usr/repackaged/lib
export CFLAGS="-I/usr/repackaged/include -fno-omit-frame-pointer -fPIE"
export CPPFLAGS="-I/usr/repackaged/include -fno-omit-frame-pointer -fPIE"
export LDFLAGS="-L/usr/repackaged/lib"
ldconfig

export CXX=clang++

cd /root/deps/mesa* &&
    ./configure --prefix=/usr/repackaged &&
    make &&
    make install &&
    ldconfig
