#!/bin/bash

set -xe

mkdir -p /root/deps

mkdir -p /pkgs && cd /pkgs &&\
    curl -L https://xorg.freedesktop.org/archive/individual/lib/libpciaccess-0.13.4.tar.bz2 -O

rm -rf /root/deps/libpciaccess* &&\
  tar xvf /pkgs/libpciaccess-0.13.4.tar.bz2 -C /root/deps

