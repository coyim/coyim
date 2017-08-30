#!/bin/bash

set -xe

mkdir -p /root/deps

mkdir -p /pkgs && cd /pkgs &&\
    curl -L https://ftp.gnome.org/pub/gnome/sources/gobject-introspection/1.52/gobject-introspection-1.52.0.tar.xz -O &&\
    echo "9fc6d1ebce5ad98942cb21e2fe8dd67b722dcc01981840632a1b233f7d0e2c1e  /pkgs/gobject-introspection-1.52.0.tar.xz" | sha256sum -c -

rm -rf /root/deps/gobject-introspection* &&\
  tar xvf /pkgs/gobject-introspection-1.52.0.tar.xz -C /root/deps
