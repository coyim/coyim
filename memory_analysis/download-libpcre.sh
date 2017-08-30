#!/bin/bash

set -xe

mkdir -p /root/deps

mkdir -p /pkgs && cd /pkgs &&\
    curl -L ftp://ftp.csx.cam.ac.uk/pub/software/programming/pcre/pcre-8.39.tar.bz2 -O

rm -rf /root/deps/pcre* &&\
  tar xvf /pkgs/pcre-8.39.tar.bz2 -C /root/deps
