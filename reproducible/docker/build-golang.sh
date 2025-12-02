#!/bin/bash

set -xe

export GOROOT_BOOTSTRAP=/usr/lib/go-1.18

cd /root/go
source /root/setup-reproducible

cd src
./make.bash
