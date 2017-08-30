#!/bin/bash

set -xe

mkdir -p /root/deps

mkdir -p /pkgs && cd /pkgs &&\
    curl -L https://xorg.freedesktop.org/archive/individual/lib/libxshmfence-1.2.tar.bz2 -O

rm -rf /root/deps/libxshmfence* &&\
  tar xvf /pkgs/libxshmfence-1.2.tar.bz2 -C /root/deps


