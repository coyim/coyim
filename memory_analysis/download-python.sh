#!/bin/bash

set -xe

mkdir -p /root/deps

mkdir -p /pkgs && cd /pkgs &&\
    curl -L https://www.python.org/ftp/python/2.7.13/Python-2.7.13.tar.xz -O

rm -rf /root/deps/Python* &&\
  tar xvf /pkgs/Python-2.7.13.tar.xz -C /root/deps
