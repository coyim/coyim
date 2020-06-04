#!/bin/bash

set -xe

export GOROOT_BOOTSTRAP=/usr/lib/go-1.13

cd /root/go
source /root/setup-reproducible

cd src
./make.bash
