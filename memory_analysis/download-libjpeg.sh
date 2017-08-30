#!/bin/bash

set -xe

mkdir -p /root/deps

mkdir -p /pkgs && cd /pkgs &&\
    curl -L http://www.ijg.org/files/jpegsrc.v8c.tar.gz -O

rm -rf /root/deps/jpeg* &&\
  tar xvf /pkgs/jpegsrc.v8c.tar.gz -C /root/deps

