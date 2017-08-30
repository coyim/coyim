#!/bin/bash

set -xe

export PATH=/usr/repackaged/bin:$PATH
export LD_LIBRARY_PATH=/usr/repackaged/lib
export CFLAGS="-I/usr/repackaged/include -fno-omit-frame-pointer -fPIE"
export CPPFLAGS="-I/usr/repackaged/include -fno-omit-frame-pointer -fPIE"
export LDFLAGS="-L/usr/repackaged/lib"
ldconfig

cd /root/deps/at-spi2-core* &&
    ./configure --prefix=/usr/repackaged --x-includes=/usr/repackaged/include --x-libraries=/usr/repackaged/lib &&
    make &&
    make install &&
    ldconfig
