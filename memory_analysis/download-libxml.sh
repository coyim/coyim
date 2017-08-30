#!/bin/bash

set -xe

mkdir -p /root/deps

mkdir -p /pkgs && cd /pkgs &&\
    curl -L ftp://xmlsoft.org/libxml2/libxml2-2.9.4.tar.gz -O

rm -rf /root/deps/libxml2* &&\
  tar xvf /pkgs/libxml2-2.9.4.tar.gz -C /root/deps
