#!/bin/bash

set -xe

mkdir -p /root/deps

mkdir -p /pkgs && cd /pkgs &&\
    curl -L https://xorg.freedesktop.org/archive/individual/lib/libXdamage-1.1.4.tar.bz2 -O

rm -rf /root/deps/libXdamage* &&\
  tar xvf /pkgs/libXdamage-1.1.4.tar.bz2 -C /root/deps


