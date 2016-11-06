#!/usr/bin/env bash

FILE=$1
TAG=$2
UPLOAD_NAME=$3

BINTRAY_API_USER=$(cat ~/.bintray_api_user)
BINTRAY_API_KEY=$(cat ~/.bintray_api_key)

curl -T $FILE -u$BINTRAY_API_USER:$BINTRAY_API_KEY -H X-Bintray-Package:coyim-bin -H X-Bintray-Version:$TAG https://api.bintray.com/content/twstrike/coyim/$TAG/linux/amd64/$UPLOAD_NAME\?override=1\&publish=1
