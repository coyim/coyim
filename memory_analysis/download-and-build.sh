#!/bin/bash

set -xe

# args: NAME URL EXTRA_CONFIGURE EXTRA_CFLAGS USE_PIC_CFLAGS DISTCLEAN

/root/installers/standard-download.sh "$1" "$2"
/root/installers/standard-build.sh "$1" "$3" "$4" "$5" "$6"
