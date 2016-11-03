#!/bin/bash

set -xe

# Download Go 1.7.3
# SHA256: 508028aac0654e993564b6e2014bf2d4a9751e3b286661b0b0040046cf18028e
# https://golang.org/dl/
mkdir -p /pkgs && cd /pkgs &&\
    curl https://storage.googleapis.com/golang/go1.7.3.linux-amd64.tar.gz -O &&\
    echo "508028aac0654e993564b6e2014bf2d4a9751e3b286661b0b0040046cf18028e /pkgs/go1.7.3.linux-amd64.tar.gz" | sha256sum -c -

rm -rf /root/go &&\
  mkdir -p /root/go &&\
  tar xvf /pkgs/go1.7.3.linux-amd64.tar.gz -C /root
