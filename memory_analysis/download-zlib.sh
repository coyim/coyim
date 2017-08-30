#!/bin/bash

set -xe

mkdir -p /root/deps

mkdir -p /pkgs && cd /pkgs &&\
    curl -L https://zlib.net/zlib-1.2.11.tar.gz -O &&\
    echo "c3e5e9fdd5004dcb542feda5ee4f0ff0744628baf8ed2dd5d66f8ca1197cb1a1  /pkgs/zlib-1.2.11.tar.gz" | sha256sum -c -

rm -rf /root/deps/zlib* &&\
  tar xvf /pkgs/zlib-1.2.11.tar.gz -C /root/deps
