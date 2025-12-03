# CoyIM - a safe and secure chat client

[![Build Status](https://github.com/coyim/coyim/workflows/CoyIM%20CI/badge.svg)](https://github.com/coyim/coyim/actions?query=workflow%3A%22CoyIM+CI%22)
[![Coverage Status](https://coveralls.io/repos/coyim/coyim/badge.svg?branch=main&service=github)](https://coveralls.io/github/coyim/coyim?branch=main)
[![Translation status](https://hosted.weblate.org/widgets/coyim/-/main/svg-badge.svg)](https://hosted.weblate.org/engage/coyim/)
[![Go Report Card](https://goreportcard.com/badge/github.com/coyim/coyim)](https://goreportcard.com/report/github.com/coyim/coyim)

<p align="center">
  <img src="build/osx/mac-bundle/coyim.iconset/icon_256x256.png">
</p>

CoyIM is a chat client for the XMPP protocol. It is built upon https://github.com/agl/xmpp-client and
https://github.com/coyim/otr3. It adds a graphical user interface and defaults to safe and secure options. Our ambition
is that it should be possible for even the most high-risk people on the planet to safely use CoyIM, without having to
make any configuration changes.

To do this, CoyIM has OTR enabled and uses Tor by default. Besides that, it will only use the Tor Onion Service for a
known server and also uses TLS and TLS certificates to verify the connection - no configuration required. The
implementation is written in the Go language, to avoid many common types of vulnerabilities that come from using unsafe
languages.

We use Weblate (https://hosted.weblate.org/projects/coyim/) to help crowd-source translation.


## Security warning

CoyIM is currently under active development. The code for our OTR implementation has been audited, with good
results. You can find out more information on our website. However, the rest of CoyIM still has not received an
audit. This is worth keeping in mind if you are using CoyIM for something sensitive.


## Getting started

Using CoyIM is very simple: you just need to download the executable file from the project's [home
page](https://coy.im/) and then run it.

**If you are using Arch Linux, you can install CoyIM via the [AUR](https://aur.archlinux.org/packages/coyim).**

When you first launch CoyIM, a wizard will appear. If you already have a Jabber client installed and configured for OTR
encryption in your computer, you can use this wizard to import your account settings as well as your OTR keys, and your
contacts' fingerprints. By importing them, you won't have to do anything else to use CoyIM.

If you don't import your account settings, keys and fingerprints through the wizard that opens at the first launch, you
can still import them by going to Accounts -> Import at a later stage.

<p align="left">
  <img src="images/wizard.png" height="242" width="242">
</p>

If the client you have been using so far is Pidgin, you will find the files you need to import in the `.purple`
directory in your home.

If you want to know more about the features you will and will not find in CoyIM, read [this
page](https://coy.im/features/).

<p align="left">
  <img src="images/main_window.png">
</p>


## Building CoyIM

**Please note**: CoyIM requires Golang version 1.21 or higher to build. CoyIM also requires at least GTK+ version 3.12
or higher. The installation of this depends on your operating system:

**Ubuntu:**

```sh
sudo apt-get install gtk+3.0 libgtk-3-dev ruby
```

**MacOS:**

```sh
brew install gnome-icon-theme
brew install gtk+3 gtk-mac-integration
brew install ruby
```

In order to build CoyIM, you should check out the source code, and run:

```sh
make build
```

It might be possible to build CoyIM using `go get` but we currently do not support this method.


## Contributing to CoyIM

We have instructions here on how you [can get started contributing to CoyIM](CONTRIBUTING.md).


## Reproducibility

CoyIM supports reproducible builds for Linux on AMD64. See [REPRODUCIBILITY](REPRODUCIBILITY.md) for instructions on how
to build and verify these builds.


## License

The CoyIM project and all source files licensed under the [GPL version 3](https://www.gnu.org/licenses/gpl-3.0.html),
unless otherwise noted.
