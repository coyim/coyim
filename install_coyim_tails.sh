#!/bin/bash

set -xe

sudo apt-get update && sudo apt-get -y --force-yes install -t testing golang
sudo apt-get install -y gcc libgtk2.0-dev

export GOPATH=$HOME/Persistent/go
mkdir -p $GOPATH

git clone https://github.com/golang/net.git $GOPATH/src/golang.org/x/net
git clone https://github.com/golang/crypto.git $GOPATH/src/golang.org/x/crypto

go get -u -v github.com/twstrike/go-gtk/gdk
go get -u -v github.com/twstrike/go-gtk/glib
go get -u -v github.com/twstrike/go-gtk/gtk
go get -u -v -tags nocli github.com/twstrike/coyim

$GOPATH/bin/coyim
