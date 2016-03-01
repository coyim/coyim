#!/bin/bash

set -xe

sudo apt-get update && sudo apt-get -y --force-yes install -t testing golang
sudo apt-get install -y gcc libgtk2.0-dev

export GOPATH=$HOME/Persistent/go
mkdir -p $GOPATH

git clone https://github.com/golang/net.git $GOPATH/src/golang.org/x/net
git clone https://github.com/golang/crypto.git $GOPATH/src/golang.org/x/crypto

GTK_VERSION=$(pkg-config --modversion gtk+-3.0 | tr . _ | cut -d '_' -f 1-2)
go get -u -v -tags "gtk_${GTK_VERSION}" github.com/gotk3/gotk3
go get -u -v -tags "gtk_${GTK_VERSION}" github.com/twstrike/gotk3adapter
go get -u -v -tags github.com/twstrike/coyim

$GOPATH/bin/coyim
