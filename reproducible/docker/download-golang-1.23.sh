#!/bin/bash

set -xe

# Download Go 1.23.12
# SHA256: d3847fef834e9db11bf64e3fb34db9c04db14e068eeb064f49af747010454f90
# https://golang.org/dl/
mkdir -p /pkgs && cd /pkgs &&\
    curl https://dl.google.com/go/go1.23.12.linux-amd64.tar.gz -O &&\
    echo "d3847fef834e9db11bf64e3fb34db9c04db14e068eeb064f49af747010454f90 /pkgs/go1.23.12.linux-amd64.tar.gz" | sha256sum -c -

rm -rf /root/go &&\
  mkdir -p /root/go &&\
  tar xvf /pkgs/go1.23.12.linux-amd64.tar.gz -C /root
