#!/bin/bash

set -xe

mkdir -p /root/deps

mkdir -p /pkgs && cd /pkgs &&\
    curl -L http://ftp.gnome.org/pub/gnome/sources/at-spi2-core/2.22/at-spi2-core-2.22.0.tar.xz -O &&\
    echo "415ea3af21318308798e098be8b3a17b2f0cf2fe16cecde5ad840cf4e0f2c80a  /pkgs/at-spi2-core-2.22.0.tar.xz" | sha256sum -c -

rm -rf /root/deps/at-spi2-core* &&\
  tar xvf /pkgs/at-spi2-core-2.22.0.tar.xz -C /root/deps

