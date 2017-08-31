#!/bin/bash

set -xe

# args: NAME

BUILD_NAME=$1

cd /root/deps/$BUILD_NAME*
find . -name "*.so" -type f -delete
find . -name "*.lo" -type f -delete
find . -name "*.o" -type f -delete
#find . -name ".deps" -type d -exec rm -r \{\} \+
find . -name ".libs" -type d -exec rm -r \{\} \+
