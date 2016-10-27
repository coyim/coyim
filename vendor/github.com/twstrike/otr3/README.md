# OTR3 [![Build Status](https://travis-ci.org/twstrike/otr3.svg?branch=master)](https://travis-ci.org/twstrike/otr3)
Implements version 3 of the OTR standard. Implements feature parity with libotr 4.1.0.

## API Documentation

[![GoDoc](https://godoc.org/github.com/twstrike/otr3?status.svg)](https://godoc.org/github.com/twstrike/otr3)

## Developing

Before doing any work, if you want to separate out your GOPATH from other projects, install direnv
```
$ brew update
$ brew install direnv
$ echo 'eval "$(direnv hook bash)"' >> ~/.bashrc
```
Then, create a symbolic link to the OTR3 repository
```
ln -s /PathToMyGoPackages/.gopkgs/otr3/src/github.com/twstrike/ .
```

Install all dependencies:

``
./deps.sh
``
