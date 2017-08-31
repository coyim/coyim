#!/bin/bash

set -xe

# args: NAME URL

NAME=$1
URL=$2
FILENAME=${URL##*/}

mkdir -p /root/deps

mkdir -p /pkgs && cd /pkgs &&\
    curl -L $URL -O

rm -rf /root/deps/$NAME* &&\
  tar xvf /pkgs/$FILENAME -C /root/deps
