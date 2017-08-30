#!/bin/bash

set -xe

mkdir -p /root/deps

mkdir -p /pkgs && cd /pkgs &&\
    curl -L https://xorg.freedesktop.org/archive/individual/lib/libX11-1.6.4.tar.bz2 -O &&\
    curl -L https://xorg.freedesktop.org/archive/individual/proto/xproto-7.0.31.tar.bz2 -O &&\
    curl -L https://xorg.freedesktop.org/archive/individual/proto/xextproto-7.3.0.tar.bz2 -O &&\
    curl -L https://xorg.freedesktop.org/archive/individual/lib/xtrans-1.3.5.tar.bz2 -O &&\
    curl -L https://xorg.freedesktop.org/archive/individual/xcb/xcb-proto-1.12.tar.bz2 -O &&\
    curl -L https://xorg.freedesktop.org/archive/individual/xcb/libxcb-1.12.tar.bz2 -O &&\
    curl -L https://xorg.freedesktop.org/archive/individual/proto/kbproto-1.0.7.tar.bz2 -O &&\
    curl -L https://xorg.freedesktop.org/archive/individual/xcb/libpthread-stubs-0.4.tar.bz2 -O &&\
    curl -L https://xorg.freedesktop.org/archive/individual/lib/libXau-1.0.8.tar.bz2 -O &&\
    curl -L https://xorg.freedesktop.org/archive/individual/proto/inputproto-2.3.2.tar.bz2 -O &&\
    curl -L https://xorg.freedesktop.org/archive/individual/lib/libXtst-1.2.3.tar.bz2 -O &&\
    curl -L https://www.x.org/archive/individual/lib/libXi-1.7.9.tar.bz2 -O &&\
    curl -L https://xorg.freedesktop.org/archive/individual/lib/libXext-1.3.3.tar.bz2 -O &&\
    curl -L https://xorg.freedesktop.org/archive/individual/proto/recordproto-1.14.2.tar.bz2 -O &&\
    curl -L https://xorg.freedesktop.org/archive/individual/lib/libXfixes-5.0.3.tar.bz2 -O &&\
    curl -L https://xorg.freedesktop.org/archive/individual/proto/fixesproto-5.0.tar.bz2 -O &&\
    curl -L https://xorg.freedesktop.org/archive/individual/proto/glproto-1.4.17.tar.bz2 -O

rm -rf /root/deps/libX11* &&\
  tar xvf /pkgs/libX11-1.6.4.tar.bz2 -C /root/deps

rm -rf /root/deps/xproto* &&\
  tar xvf /pkgs/xproto-7.0.31.tar.bz2 -C /root/deps

rm -rf /root/deps/xextproto* &&\
  tar xvf /pkgs/xextproto-7.3.0.tar.bz2 -C /root/deps

rm -rf /root/deps/xtrans* &&\
  tar xvf /pkgs/xtrans-1.3.5.tar.bz2 -C /root/deps

rm -rf /root/deps/libpthread-stubs* &&\
  tar xvf /pkgs/libpthread-stubs-0.4.tar.bz2 -C /root/deps

rm -rf /root/deps/libXau* &&\
  tar xvf /pkgs/libXau-1.0.8.tar.bz2 -C /root/deps

rm -rf /root/deps/xcb-proto* &&\
  tar xvf /pkgs/xcb-proto-1.12.tar.bz2 -C /root/deps

rm -rf /root/deps/libxcb* &&\
  tar xvf /pkgs/libxcb-1.12.tar.bz2 -C /root/deps

rm -rf /root/deps/kbproto* &&\
  tar xvf /pkgs/kbproto-1.0.7.tar.bz2 -C /root/deps

rm -rf /root/deps/inputproto* &&\
  tar xvf /pkgs/inputproto-2.3.2.tar.bz2 -C /root/deps

rm -rf /root/deps/libXtst* &&\
  tar xvf /pkgs/libXtst-1.2.3.tar.bz2 -C /root/deps

rm -rf /root/deps/libXi* &&\
  tar xvf /pkgs/libXi-1.7.9.tar.bz2 -C /root/deps

rm -rf /root/deps/libXext* &&\
  tar xvf /pkgs/libXext-1.3.3.tar.bz2 -C /root/deps

rm -rf /root/deps/recordproto* &&\
  tar xvf /pkgs/recordproto-1.14.2.tar.bz2 -C /root/deps

rm -rf /root/deps/libXfixes* &&\
  tar xvf /pkgs/libXfixes-5.0.3.tar.bz2 -C /root/deps

rm -rf /root/deps/fixesproto* &&\
  tar xvf /pkgs/fixesproto-5.0.tar.bz2 -C /root/deps

rm -rf /root/deps/glproto* &&\
  tar xvf /pkgs/glproto-1.4.17.tar.bz2 -C /root/deps
