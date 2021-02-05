#!/usr/bin/env bash

set -e

found=0

echo "Checking for esc binary in $GOPATH"

while IFS=':' read -ra GOP; do
    for i in "${GOP[@]}"; do
        echo "  ... testing for esc in $i"
        if [ -f $i/bin/esc ]; then
            found=1
            cp $i/bin/esc $1/esc
        fi
    done
done <<< "$GOPATH"

if [ $found -eq 0 ]; then
    echo "The program 'esc' is required but not available. Please install it by running 'make deps'."
    exit 1
fi
