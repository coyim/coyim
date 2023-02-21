#!/usr/bin/env bash

set -e

onpath=`which esc` || found=0

if [ -x "$onpath" ]; then
    found=1
    cp $onpath $1/esc
fi

if [ $found -eq 0 ]; then
    GP=`go env GOPATH`

    while IFS=':' read -ra GOP; do
        for i in "${GOP[@]}"; do
            if [ -f $i/bin/esc ]; then
                found=1
                cp $i/bin/esc $1/esc
            fi
        done
    done <<< "$GP"
fi

if [ $found -eq 0 ]; then
    echo "The program 'esc' is required but not available. Please install it by running 'make deps'."
    exit 1
fi
