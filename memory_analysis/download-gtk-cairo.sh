#!/bin/bash

set -xe

mkdir -p /root/gtk

mkdir -p /pkgs && cd /pkgs &&\
    curl -L https://www.cairographics.org/releases/cairo-1.14.8.tar.xz -O &&\
    echo "c6f7b99986f93c9df78653c3e6a3b5043f65145e  /pkgs/cairo-1.14.8.tar.xz" | sha1sum -c -

rm -rf /root/gtk/cairo* &&\
  tar xvf /pkgs/cairo-1.14.8.tar.xz -C /root/gtk
