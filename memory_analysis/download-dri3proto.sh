#!/bin/bash

set -xe

mkdir -p /root/deps

mkdir -p /pkgs && cd /pkgs &&\
    curl -L https://xorg.freedesktop.org/archive/individual/proto/dri3proto-1.0.tar.bz2 -O

rm -rf /root/deps/dri3proto* &&\
  tar xvf /pkgs/dri3proto-1.0.tar.bz2 -C /root/deps
