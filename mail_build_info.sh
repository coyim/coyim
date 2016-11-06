#!/usr/bin/env bash

FILE=$1
TAG=$2

echo "Hello, attached is a signed build info for tag: $TAG" | mail -a $FILE -s "Signed build info for tag: $TAG" coyim@thoughtworks.com
