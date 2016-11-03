#!/bin/bash

set -xe

# Requires gccgo-X to be installed
export GOROOT_BOOTSTRAP=/usr

cd /root/go
source /root/setup-reproducible

cd src
./make.bash
