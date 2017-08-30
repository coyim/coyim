#!/bin/bash

set -xe

# Download Go 1.9.0
# SHA256: d70eadefce8e160638a9a6db97f7192d8463069ab33138893ad3bf31b0650a79
# https://golang.org/dl/
mkdir -p /pkgs && cd /pkgs &&\
    curl https://storage.googleapis.com/golang/go1.9.linux-amd64.tar.gz -O &&\
    echo "d70eadefce8e160638a9a6db97f7192d8463069ab33138893ad3bf31b0650a79 /pkgs/go1.9.linux-amd64.tar.gz" | sha256sum -c -

rm -rf /root/go &&\
  mkdir -p /root/go &&\
  tar xvf /pkgs/go1.9.linux-amd64.tar.gz -C /root
