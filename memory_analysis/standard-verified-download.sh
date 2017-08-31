#!/bin/bash

set -xe

# args: NAME URL SHA256SUM

NAME=$1
URL=$2
SUM=$3
FILENAME=${URL##*/}

mkdir -p /root/deps

mkdir -p /pkgs && cd /pkgs &&\
    curl -L $URL -O &&\
    echo "$SUM  /pkgs/$FILENAME" | sha256sum -c -

rm -rf /root/deps/$NAME* &&\
  tar xvf /pkgs/$FILENAME -C /root/deps
