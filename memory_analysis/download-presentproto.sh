#!/bin/bash

set -xe

mkdir -p /root/deps

mkdir -p /pkgs && cd /pkgs &&\
    curl -L https://xorg.freedesktop.org/archive/individual/proto/presentproto-1.1.tar.bz2 -O

rm -rf /root/deps/presentproto* &&\
  tar xvf /pkgs/presentproto-1.1.tar.bz2 -C /root/deps
