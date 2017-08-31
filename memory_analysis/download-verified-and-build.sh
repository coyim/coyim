#!/bin/bash

set -xe

# args: NAME URL SHA256SUM EXTRA_CONFIGURE EXTRA_CFLAGS USE_PIC_CFLAGS DISTCLEAN

/root/installers/standard-verified-download.sh "$1" "$2" "$3"
/root/installers/standard-build.sh "$1" "$4" "$5" "$6" "$7"
