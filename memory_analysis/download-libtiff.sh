#!/bin/bash

set -xe

mkdir -p /root/deps

mkdir -p /pkgs && cd /pkgs &&\
    curl -L ftp://download.osgeo.org/libtiff/tiff-4.0.7.tar.gz -O

rm -rf /root/deps/tiff* &&\
  tar xvf /pkgs/tiff-4.0.7.tar.gz -C /root/deps
