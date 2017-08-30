#!/bin/bash

set -xe

mkdir -p /root/deps

mkdir -p /pkgs && cd /pkgs &&\
    curl -L https://dri.freedesktop.org/libdrm/libdrm-2.4.76.tar.bz2 -O

rm -rf /root/deps/libdrm* &&\
  tar xvf /pkgs/libdrm-2.4.76.tar.bz2 -C /root/deps
