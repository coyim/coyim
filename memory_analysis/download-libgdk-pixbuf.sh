#!/bin/bash

set -xe

mkdir -p /root/deps

mkdir -p /pkgs && cd /pkgs &&\
    curl -L https://ftp.gnome.org/pub/gnome/sources/gdk-pixbuf/2.36/gdk-pixbuf-2.36.5.tar.xz -O &&\
    echo "7ace06170291a1f21771552768bace072ecdea9bd4a02f7658939b9a314c40fc  /pkgs/gdk-pixbuf-2.36.5.tar.xz" | sha256sum -c -

rm -rf /root/deps/gdk-pixbuf* &&\
  tar xvf /pkgs/gdk-pixbuf-2.36.5.tar.xz -C /root/deps
