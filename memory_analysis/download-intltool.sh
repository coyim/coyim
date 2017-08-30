#!/bin/bash

set -xe

mkdir -p /root/deps

mkdir -p /pkgs && cd /pkgs &&\
    curl -L https://launchpad.net/intltool/trunk/0.51.0/+download/intltool-0.51.0.tar.gz -O

rm -rf /root/deps/intltool* &&\
  tar xvf /pkgs/intltool-0.51.0.tar.gz -C /root/deps

