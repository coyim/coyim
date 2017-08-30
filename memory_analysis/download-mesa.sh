#!/bin/bash

set -xe

mkdir -p /root/deps

mkdir -p /pkgs && cd /pkgs &&\
    curl -L ftp://ftp.freedesktop.org/pub/mesa/mesa-17.0.3.tar.xz -O

rm -rf /root/deps/mesa-* &&\
  tar xvf /pkgs/mesa-17.0.3.tar.xz -C /root/deps

