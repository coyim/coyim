#!/bin/bash

set -xe

mkdir -p /root/gtk

mkdir -p /pkgs && cd /pkgs &&\
    curl -L http://ftp.gnome.org/pub/gnome/sources/pango/1.40/pango-1.40.7.tar.xz -O &&\
    echo "517645c00c4554e82c0631e836659504d3fd3699c564c633fccfdfd37574e278  /pkgs/pango-1.40.7.tar.xz" | sha256sum -c -

rm -rf /root/gtk/pango &&\
  tar xvf /pkgs/pango-1.40.7.tar.xz -C /root/gtk
