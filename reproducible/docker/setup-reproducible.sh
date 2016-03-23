#!/bin/bash

set -xe

# Does nothing if REFERENCE_DATETIME is missing
test -z "$REFERENCE_DATETIME" && return 0

# lucid
# export LD_PRELOAD=/usr/lib/faketime/libfaketime.so.1
# vivid
export LD_PRELOAD=/usr/lib/x86_64-linux-gnu/faketime/libfaketime.so.1
export FAKETIME=$REFERENCE_DATETIME
export TZ=UTC
export LC_ALL=C

find -type f -print0 | xargs -0 touch --date="$REFERENCE_DATETIME"

