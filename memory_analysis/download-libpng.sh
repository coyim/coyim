#!/bin/bash

set -xe

mkdir -p /root/deps

mkdir -p /pkgs && cd /pkgs &&\
    curl -L ftp://ftp-osl.osuosl.org/pub/libpng/src/libpng16/libpng-1.6.32.tar.xz -O

rm -rf /root/deps/libpng* &&\
  tar xvf /pkgs/libpng-1.6.32.tar.xz -C /root/deps
