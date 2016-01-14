#!/usr/bin/env bash

cat <<EOF >| version_generated.go
package main

func init() {
	coyimVersion = "$1"
}
EOF
