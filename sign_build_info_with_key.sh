#!/usr/bin/env bash

KEYID=$1

gpg2 --armor --detach-sign -u $KEYID --output bin/build_info.$KEYID.asc bin/build_info
