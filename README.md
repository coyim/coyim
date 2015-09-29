# coyim - a safe and secure chat client

CoyIM is a new client for the XMPP protocol. It is built upon https://github.com/agl/xmpp-client. It adds a graphical user interface and tries to be safe and secure by default. Our ambition is that it should be possible for even the most high-risk people on the planet to safely use CoyIM, without having to make any configuration changes.

To do this, we enable OTR by default, we try hard to use Tor and Tor hidden services, and also to use TLS and TLS certificates to verify the connection. The implementation is written in the Go language, to avoid many common types of vulnerabilities that come from using unsafe languages.

## Security warning

CoyIM is currently under active development. There have been no security audits of the code, and you should currently not use this for anything sensitive.
