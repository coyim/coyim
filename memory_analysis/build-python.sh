#!/bin/bash

set -xe

export PATH=/usr/repackaged/bin:$PATH
export LD_LIBRARY_PATH=/usr/repackaged/lib
export CFLAGS="-I/usr/repackaged/include -I/usr/repackaged/include/ncurses -fno-omit-frame-pointer -fPIC"
export CPPFLAGS="-I/usr/repackaged/include -I/usr/repackaged/include/ncurses -fno-omit-frame-pointer -fPIC"
export LDFLAGS="-L/usr/repackaged/lib"
ldconfig

cd /root/deps/Python* &&
    ./configure --prefix=/usr/repackaged --with-system-expat &&
    make &&
    make install &&
    ldconfig
