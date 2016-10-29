#!/bin/bash

set -xe

# Download Go 1.6.3
# SHA1: cdde5e08530c0579255d6153b08fdb3b8e47caabbe717bc7bcd7561275a87aeb
# https://golang.org/dl/
mkdir -p /pkgs && cd /pkgs &&\
    curl https://storage.googleapis.com/golang/go1.6.3.linux-amd64.tar.gz -O &&\
    echo "cdde5e08530c0579255d6153b08fdb3b8e47caabbe717bc7bcd7561275a87aeb /pkgs/go1.6.3.linux-amd64.tar.gz" | sha256sum -c -

rm -rf /root/go &&\
  mkdir -p /root/go &&\
  tar xvf /pkgs/go1.6.3.linux-amd64.tar.gz -C /root

