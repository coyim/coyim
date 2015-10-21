# coyim - a safe and secure chat client

[![Build Status](https://travis-ci.org/twstrike/coyim.svg?branch=master)](https://travis-ci.org/twstrike/coyim)
[![Coverage Status](https://coveralls.io/repos/twstrike/coyim/badge.svg?branch=master&service=github)](https://coveralls.io/github/twstrike/coyim?branch=master)

CoyIM is a new client for the XMPP protocol. It is built upon https://github.com/agl/xmpp-client. It adds a graphical user interface and tries to be safe and secure by default. Our ambition is that it should be possible for even the most high-risk people on the planet to safely use CoyIM, without having to make any configuration changes.

To do this, we enable OTR by default, we try hard to use Tor and Tor hidden services, and also to use TLS and TLS certificates to verify the connection. The implementation is written in the Go language, to avoid many common types of vulnerabilities that come from using unsafe languages.

## Security warning

CoyIM is currently under active development. There have been no security audits of the code, and you should currently not use this for anything sensitive.

## Installation instructions

### GUI version

The GUI version requires GTK+ >= 3.6.16, which installation depends on your OS:

**Ubuntu:**

	sudo apt-get install -qq -y gtk+3.0 libgtk-3-dev

**Mac:**

	brew install gtk+3

Then install coyim:

```
GTK_VERSION=$(pkg-config --modversion gtk+-3.0 | tr . _ | cut -d '_' -f 1-2)
go get -u -tags "nocli gtk_${GTK_VERSION}" github.com/twstrike/coyim
```

### CLI version (xmpp-client)

```
go get -u github.com/twstrike/coyim
```

