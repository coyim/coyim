#!/bin/bash

set -xe

mkdir -p /root/deps

mkdir -p /pkgs && cd /pkgs &&\
    curl -L https://download.gnome.org/sources/librsvg/2.40/librsvg-2.40.16.tar.xz -O &&\
    echo "d48bcf6b03fa98f07df10332fb49d8c010786ddca6ab34cbba217684f533ff2e  /pkgs/librsvg-2.40.16.tar.xz" | sha256sum -c -

rm -rf /root/deps/librsvg* &&\
  tar xvf /pkgs/librsvg-2.40.16.tar.xz -C /root/deps
