#!/bin/bash

set -xe

mkdir -p /root/deps

mkdir -p /pkgs && cd /pkgs &&\
    curl -L https://dbus.freedesktop.org/releases/dbus/dbus-1.10.10.tar.gz -O

rm -rf /root/deps/dbus* &&\
  tar xvf /pkgs/dbus-1.10.10.tar.gz -C /root/deps
