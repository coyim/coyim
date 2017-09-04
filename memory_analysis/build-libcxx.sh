#!/bin/bash

set -xe

cd /root/deps/gcc*

#CLANG_HEADERS="/usr/lib/llvm-4.0/lib/clang/4.0.0/include"
#GCC_HEADERS=$(g++ -v -xc++ -c /dev/null |& grep -E '/[01-9.]+/include$' | head -1 | sed 's/ \+//')
MSAN_CFLAGS="-fsanitize=memory -g -O2 -fno-omit-frame-pointer"
MSAN_LDFLAGS="-fsanitize=memory"
mkdir build && cd build
CFLAGS="$MSAN_CFLAGS" CXXFLAGS="$MSAN_CFLAGS" LDFLAGS="$MSAN_LDFLAGS" ../libstdc++-v3/configure --enable-multilib=no --prefix=/usr/repackaged
make -j10
make install
