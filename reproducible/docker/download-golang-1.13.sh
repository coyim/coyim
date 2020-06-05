#!/bin/bash

set -xe

# Download Go 1.13.12
# SHA256: 9cacc6653563771b458c13056265aa0c21b8a23ca9408278484e4efde4160618
# https://golang.org/dl/
mkdir -p /pkgs && cd /pkgs &&\
    curl https://dl.google.com/go/go1.13.12.linux-amd64.tar.gz -O &&\
    echo "9cacc6653563771b458c13056265aa0c21b8a23ca9408278484e4efde4160618 /pkgs/go1.13.12.linux-amd64.tar.gz" | sha256sum -c -

rm -rf /root/go &&\
  mkdir -p /root/go &&\
  tar xvf /pkgs/go1.13.12.linux-amd64.tar.gz -C /root
