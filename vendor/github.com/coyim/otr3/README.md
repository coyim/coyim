# OTR3 

[![Build Status](https://github.com/coyim/otr3/workflows/OTR3%20CI/badge.svg)](https://github.com/coyim/otr3/actions?query=workflow%3A%22OTR3+CI%22)
[![Coverage Status](https://coveralls.io/repos/coyim/otr3/badge.svg?branch=main&service=github)](https://coveralls.io/github/coyim/otr3?branch=main)

Implements version 3 of the OTR standard. Implements feature parity with libotr 4.1.0.

## API Documentation

[![GoDoc](https://godoc.org/github.com/coyim/otr3?status.svg)](https://godoc.org/github.com/coyim/otr3)

## Developing

Before doing any work, if you want to separate out your GOPATH from other projects, install direnv
```
$ brew update
$ brew install direnv
$ echo 'eval "$(direnv hook bash)"' >> ~/.bashrc
```
Then, create a symbolic link to the OTR3 repository
```
ln -s /PathToMyGoPackages/.gopkgs/otr3/src/github.com/coyim/ .
```

Install all dependencies:

``
./deps.sh
``
