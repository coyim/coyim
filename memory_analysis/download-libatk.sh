#!/bin/bash

set -xe

mkdir -p /root/deps

mkdir -p /pkgs && cd /pkgs &&\
    curl -L http://ftp.gnome.org/pub/gnome/sources/atk/2.22/atk-2.22.0.tar.xz -O &&\
    echo "d349f5ca4974c9c76a4963e5b254720523b0c78672cbc0e1a3475dbd9b3d44b6  /pkgs/atk-2.22.0.tar.xz" | sha256sum -c -

rm -rf /root/deps/atk-2* &&\
  tar xvf /pkgs/atk-2.22.0.tar.xz -C /root/deps
