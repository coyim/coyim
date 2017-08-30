#!/bin/bash

set -xe

cd /root/gtk/cairo* &&
    make CC=clang CFLAGS="-fsanitize=memory -fsanitize-memory-track-origins -fsanitize-recover=all -fno-omit-frame-pointer -fPIE" &&
    make install
