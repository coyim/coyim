#!/bin/bash

set -xe

export PATH=/usr/repackaged/bin:$PATH
export LD_LIBRARY_PATH=/usr/repackaged/lib
export CFLAGS="-I/usr/repackaged/include -fno-omit-frame-pointer -fPIE"
export CPPFLAGS="-I/usr/repackaged/include -fno-omit-frame-pointer -fPIE"
export LDFLAGS="-L/usr/repackaged/lib"
ldconfig

cd /root/deps/xextproto* &&
    ./configure  --prefix=/usr/repackaged  &&
    make &&
    make install &&
    ldconfig &&

    cd /root/deps/xtrans* &&
    ./configure  --prefix=/usr/repackaged  &&
    make &&
    make install &&
    ldconfig &&

    cd /root/deps/xcb-proto* &&
    ./configure  --prefix=/usr/repackaged  &&
    make &&
    make install &&
    ldconfig &&

    cd /root/deps/libpthread-stubs* &&
    ./configure  --prefix=/usr/repackaged  &&
    make &&
    make install &&
    ldconfig &&

    cd /root/deps/kbproto* &&
    ./configure  --prefix=/usr/repackaged  &&
    make &&
    make install &&
    ldconfig &&

    cd /root/deps/inputproto* &&
    ./configure  --prefix=/usr/repackaged  &&
    make &&
    make install &&
    ldconfig &&

    cd /root/deps/recordproto* &&
    ./configure  --prefix=/usr/repackaged  &&
    make &&
    make install &&
    ldconfig &&

    cd /root/deps/xproto* &&
    ./configure  --prefix=/usr/repackaged  &&
    make &&
    make install &&
    ldconfig &&

    cd /root/deps/libXau* &&
    ./configure  --prefix=/usr/repackaged  &&
    make &&
    make install &&
    ldconfig &&

    cd /root/deps/libxcb* &&
    ./configure  --prefix=/usr/repackaged  &&
    make &&
    make install &&
    ldconfig &&

    export LD_LIBRARY_PATH=/usr/repackaged/lib &&
    ldconfig &&

    cd /root/deps/libX11* &&
    ./configure  --prefix=/usr/repackaged  &&
    make &&
    make install &&
    ldconfig &&

    cd /root/deps/libXext* &&
    ./configure  --prefix=/usr/repackaged  &&
    make &&
    make install &&
    ldconfig &&

    cd /root/deps/fixesproto* &&
    ./configure  --prefix=/usr/repackaged  &&
    make &&
    make install &&
    ldconfig &&

    cd /root/deps/libXfixes* &&
    ./configure  --prefix=/usr/repackaged  &&
    make &&
    make install &&
    ldconfig &&

    cd /root/deps/libXi* &&
    ./configure  --prefix=/usr/repackaged  &&
    make &&
    make install &&
    ldconfig &&

    cd /root/deps/libXtst* &&
    ./configure  --prefix=/usr/repackaged  &&
    make &&
    make install &&
    ldconfig &&

    cd /root/deps/glproto* &&
    ./configure  --prefix=/usr/repackaged  &&
    make &&
    make install &&
    ldconfig
