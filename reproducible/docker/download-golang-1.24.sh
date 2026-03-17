#!/bin/bash

set -xe

# Download Go 1.24.13
# SHA256: 1fc94b57134d51669c72173ad5d49fd62afb0f1db9bf3f798fd98ee423f8d730
# https://golang.org/dl/
mkdir -p /pkgs && cd /pkgs &&\
    curl https://dl.google.com/go/go1.24.13.linux-amd64.tar.gz -O &&\
    echo "1fc94b57134d51669c72173ad5d49fd62afb0f1db9bf3f798fd98ee423f8d730 /pkgs/go1.24.13.linux-amd64.tar.gz" | sha256sum -c -

rm -rf /root/go &&\
  mkdir -p /root/go &&\
  tar xvf /pkgs/go1.24.13.linux-amd64.tar.gz -C /root
