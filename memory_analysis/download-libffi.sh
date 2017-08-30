#!/bin/bash

set -xe

mkdir -p /root/deps

mkdir -p /pkgs && cd /pkgs &&\
    curl -L ftp://sourceware.org/pub/libffi/libffi-3.2.1.tar.gz -O

rm -rf /root/deps/libffi* &&\
  tar xvf /pkgs/libffi-3.2.1.tar.gz -C /root/deps
