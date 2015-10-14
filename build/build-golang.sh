#!/bin/bash

set -xe

# Requires gccgo-X to be installed
export GOROOT_BOOTSTRAP=/usr

rm -rf /root/go &&\
  mkdir -p /root/go &&\
  tar xvf /pkgs/go1.5.1.linux-amd64.tar.gz -C /root

cd /root/go
source /root/setup-reproducible

cd src
./make.bash

