#!/bin/bash

set -xe

mkdir -p /root/deps

mkdir -p /pkgs && cd /pkgs &&\
    curl -L https://www.kernel.org/pub/linux/utils/util-linux/v2.29/util-linux-2.29.1.tar.xz -O

rm -rf /root/deps/util-linux* &&\
  tar xvf /pkgs/util-linux-2.29.1.tar.xz -C /root/deps
