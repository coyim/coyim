#!/bin/bash

set -xe

# Download Go 1.5.1
# SHA1: 46eecd290d8803887dec718c691cc243f2175fe0
# https://golang.org/dl/
mkdir -p /pkgs && cd /pkgs &&\
    curl https://storage.googleapis.com/golang/go1.5.1.linux-amd64.tar.gz -O &&\
    echo "46eecd290d8803887dec718c691cc243f2175fe0 /pkgs/go1.5.1.linux-amd64.tar.gz" | sha1sum -c -

rm -rf /root/go &&\
  mkdir -p /root/go &&\
  tar xvf /pkgs/go1.5.1.linux-amd64.tar.gz -C /root

