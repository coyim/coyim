#!/bin/bash

set -xe

mkdir -p /root/deps

mkdir -p /pkgs && cd /pkgs &&\
    curl -L https://xorg.freedesktop.org/archive/individual/proto/dri2proto-2.8.tar.bz2 -O

rm -rf /root/deps/dri2proto* &&\
  tar xvf /pkgs/dri2proto-2.8.tar.bz2 -C /root/deps

