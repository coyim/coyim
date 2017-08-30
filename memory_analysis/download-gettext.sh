#!/bin/bash

set -xe

mkdir -p /root/deps

mkdir -p /pkgs && cd /pkgs &&\
    curl -L https://ftp.gnu.org/pub/gnu/gettext/gettext-0.19.8.tar.xz -O

rm -rf /root/deps/gettext* &&\
  tar xvf /pkgs/gettext-0.19.8.tar.xz -C /root/deps
