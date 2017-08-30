#!/bin/bash

set -xe

mkdir -p /root/deps

mkdir -p /pkgs && cd /pkgs &&\
    curl -L http://ftp.gnome.org/pub/gnome/sources/at-spi2-atk/2.22/at-spi2-atk-2.22.0.tar.xz -O &&\
    echo "e8bdedbeb873eb229eb08c88e11d07713ec25ae175251648ad1a9da6c21113c1  /pkgs/at-spi2-atk-2.22.0.tar.xz" | sha256sum -c -

rm -rf /root/deps/at-spi2-atk* &&\
  tar xvf /pkgs/at-spi2-atk-2.22.0.tar.xz -C /root/deps
