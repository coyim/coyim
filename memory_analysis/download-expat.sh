#!/bin/bash

set -xe

mkdir -p /root/deps

mkdir -p /pkgs && cd /pkgs &&\
    curl -L https://downloads.sourceforge.net/project/expat/expat/2.2.4/expat-2.2.4.tar.bz2 -O

rm -rf /root/deps/expat* &&\
  tar xvf /pkgs/expat-2.2.4.tar.bz2 -C /root/deps
