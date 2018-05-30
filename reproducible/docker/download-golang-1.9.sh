#!/bin/bash

set -xe

# Download Go 1.9.6
# SHA256: d1eb07f99ac06906225ac2b296503f06cc257b472e7d7817b8f822fe3766ebfe
# https://golang.org/dl/
mkdir -p /pkgs && cd /pkgs &&\
    curl https://storage.googleapis.com/golang/go1.9.6.linux-amd64.tar.gz -O &&\
    echo "d1eb07f99ac06906225ac2b296503f06cc257b472e7d7817b8f822fe3766ebfe /pkgs/go1.9.6.linux-amd64.tar.gz" | sha256sum -c -

rm -rf /root/go &&\
  mkdir -p /root/go &&\
  tar xvf /pkgs/go1.9.6.linux-amd64.tar.gz -C /root
