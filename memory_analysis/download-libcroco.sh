#!/bin/bash

set -xe

mkdir -p /root/deps

mkdir -p /pkgs && cd /pkgs &&\
    curl -L http://ftp.gnome.org/pub/GNOME/sources/libcroco/0.6/libcroco-0.6.11.tar.xz -O &&\
    echo "132b528a948586b0dfa05d7e9e059901bca5a3be675b6071a90a90b81ae5a056  /pkgs/libcroco-0.6.11.tar.xz" | sha256sum -c -

rm -rf /root/deps/libcroco* &&\
  tar xvf /pkgs/libcroco-0.6.11.tar.xz -C /root/deps
