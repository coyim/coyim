#!/bin/bash

set -xe

mkdir -p /root/gtk

mkdir -p /pkgs && cd /pkgs &&\
    curl -L http://ftp.gnome.org/pub/gnome/sources/glib/2.53/glib-2.53.6.tar.xz -O &&\
    echo "e01296a9119c09d2dccb37ad09f5eaa1e0c5570a473d8fed04fc759ace6fb6cc  /pkgs/glib-2.53.6.tar.xz" | sha256sum -c -

rm -rf /root/gtk/glib* &&\
  tar xvf /pkgs/glib-2.53.6.tar.xz -C /root/gtk
