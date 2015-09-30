#!/bin/bash
 
sudo apt-get update && sudo apt-get -y --force-yes install -t testing golang
sudo apt-get install -y libgtk2.0-dev
sudo apt-get install -y gcc

cd $HOME/Persistent && mkdir go
export GOPATH=$HOME/Persistent/go
cd $GOPATH && mkdir src 
cd src
mkdir golang.org && cd golang.org && mkdir x && cd x
git clone https://github.com/golang/net.git
git clone https://github.com/golang/crypto.git

go get -v github.com/twstrike/go-gtk/gdk
go get -v github.com/twstrike/go-gtk/glib
go get -v github.com/twstrike/go-gtk/gtk
go get -v github.com/twstrike/coyim

cd $GOPATH/src/github.com/twstrike/coyim
go build -tags nocli

./coyim
