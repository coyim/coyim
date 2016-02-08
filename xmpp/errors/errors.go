package errors

import "errors"

var (
	// ErrAuthenticationFailed indicates a failure to authenticate to the server with the user and password provided.
	ErrAuthenticationFailed = errors.New("could not authenticate to the XMPP server")

	//ErrConnectionFailed indicates a failure to connect to the server provided.
	ErrConnectionFailed = errors.New("could not connect to XMPP server")

	//ErrTCPBindingFailed indicates a failure to determine a server address for the given origin domain
	ErrTCPBindingFailed = errors.New("failed to find a TCP address for XMPP server")
)
