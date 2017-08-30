#!/bin/bash

set -xe

cd /root/gtk/glib* &&
    make CC=clang CFLAGS="-fsanitize=memory -fsanitize-memory-track-origins -fsanitize-recover=all -fno-omit-frame-pointer -fPIE" &&
    make install
