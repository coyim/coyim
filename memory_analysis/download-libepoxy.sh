#!/bin/bash

set -xe

mkdir -p /root/deps

mkdir -p /pkgs && cd /pkgs &&\
    curl -L https://github.com/anholt/libepoxy/releases/download/v1.3.1/libepoxy-1.3.1.tar.bz2 -O

rm -rf /root/deps/libepoxy* &&\
  tar xvf /pkgs/libepoxy-1.3.1.tar.bz2 -C /root/deps
