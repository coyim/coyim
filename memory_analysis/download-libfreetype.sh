#!/bin/bash

set -xe

mkdir -p /root/deps

mkdir -p /pkgs && cd /pkgs &&\
    curl -L http://download.savannah.gnu.org/releases/freetype/freetype-2.6.3.tar.bz2 -O

rm -rf /root/deps/freetype* &&\
  tar xvf /pkgs/freetype-2.6.3.tar.bz2 -C /root/deps
