#!/bin/bash

set -xe

mkdir -p /root/deps

mkdir -p /pkgs && cd /pkgs &&\
    curl -L https://xorg.freedesktop.org/archive/individual/proto/damageproto-1.2.1.tar.bz2 -O

rm -rf /root/deps/damageproto* &&\
  tar xvf /pkgs/damageproto-1.2.1.tar.bz2 -C /root/deps
