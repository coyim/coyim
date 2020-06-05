#!/bin/bash

set -xe

# Download Go 1.13.12
# SHA256: 17ba2c4de4d78793a21cc659d9907f4356cd9c8de8b7d0899cdedcef712eba34
# https://golang.org/dl/
mkdir -p /pkgs && cd /pkgs &&\
    curl https://dl.google.com/go/go1.13.12.linux-amd64.tar.gz -O &&\
    echo "17ba2c4de4d78793a21cc659d9907f4356cd9c8de8b7d0899cdedcef712eba34 /pkgs/go1.13.12.linux-amd64.tar.gz" | sha256sum -c -

rm -rf /root/go &&\
  mkdir -p /root/go &&\
  tar xvf /pkgs/go1.13.12.linux-amd64.tar.gz -C /root
