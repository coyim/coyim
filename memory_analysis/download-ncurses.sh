#!/bin/bash

set -xe

mkdir -p /root/deps

mkdir -p /pkgs && cd /pkgs &&\
    curl -L https://ftp.gnu.org/pub/gnu/ncurses/ncurses-6.0.tar.gz -O

rm -rf /root/deps/ncurses* &&\
  tar xvf /pkgs/ncurses-6.0.tar.gz -C /root/deps
