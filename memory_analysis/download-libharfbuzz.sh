#!/bin/bash

set -xe

mkdir -p /root/deps

mkdir -p /pkgs && cd /pkgs &&\
    curl -L https://www.freedesktop.org/software/harfbuzz/release/harfbuzz-1.4.2.tar.bz2 -O &&\
    echo "8f234dcfab000fdec24d43674fffa2fdbdbd654eb176afbde30e8826339cb7b3  /pkgs/harfbuzz-1.4.2.tar.bz2" | sha256sum -c -

rm -rf /root/deps/harfbuzz* &&\
  tar xvf /pkgs/harfbuzz-1.4.2.tar.bz2 -C /root/deps
