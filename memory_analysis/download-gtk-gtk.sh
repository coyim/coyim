#!/bin/bash

set -xe

mkdir -p /root/gtk

mkdir -p /pkgs && cd /pkgs &&\
    curl -L http://ftp.gnome.org/pub/gnome/sources/gtk+/3.22/gtk+-3.22.17.tar.xz -O &&\
    echo "a6c1fb8f229c626a3d9c0e1ce6ea138de7f64a5a6bc799d45fa286fe461c3437  /pkgs/gtk+-3.22.17.tar.xz" | sha256sum -c -

rm -rf /root/gtk/gtk* &&\
  tar xvf /pkgs/gtk+-3.22.17.tar.xz -C /root/gtk
