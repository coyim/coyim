#!/bin/bash

set -xe

mkdir -p /root/deps

mkdir -p /pkgs && cd /pkgs &&\
    curl -L https://www.freedesktop.org/software/fontconfig/release/fontconfig-2.11.94.tar.bz2 -O

rm -rf /root/deps/fontconfig* &&\
  tar xvf /pkgs/fontconfig-2.11.94.tar.bz2 -C /root/deps
