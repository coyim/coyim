#!/usr/bin/env bash

KEYID=$1
DST=${2:-bin}

gpg2 --armor --detach-sign -u $KEYID --output $DST/build_info.$KEYID.asc $DST/build_info
