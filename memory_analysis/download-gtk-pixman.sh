#!/bin/bash

set -xe

mkdir -p /root/gtk

mkdir -p /pkgs && cd /pkgs &&\
    curl -L https://www.cairographics.org/releases/pixman-0.34.0.tar.gz -O &&\
    echo "a1b1683c1a55acce9d928fea1ab6ceb79142ddc7  /pkgs/pixman-0.34.0.tar.gz" | sha1sum -c -

rm -rf /root/gtk/pixman* &&\
  tar xvf /pkgs/pixman-0.34.0.tar.gz -C /root/gtk
