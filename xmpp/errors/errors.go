package errors

import (
	"errors"
	"fmt"
)

var (
	// ErrAuthenticationFailed indicates a failure to authenticate to the server with the user and password provided.
	ErrAuthenticationFailed = errors.New("could not authenticate to the XMPP server")

	//ErrConnectionFailed indicates a failure to connect to the server provided.
	ErrConnectionFailed = errors.New("could not connect to XMPP server")

	//ErrTCPBindingFailed indicates a failure to determine a server address for the given origin domain
	ErrTCPBindingFailed = errors.New("failed to find a TCP address for XMPP server")
)

// ErrFailedToConnect is an error representing connection failure
type ErrFailedToConnect struct {
	Addr string
	Err  error
}

func (e *ErrFailedToConnect) Error() string {
	return fmt.Sprintf("Failed to connect to %s: %s", e.Addr, e.Err.Error())
}

// CreateErrFailedToConnect will create a ErrFailedToConnect from the given data
func CreateErrFailedToConnect(addr string, err error) *ErrFailedToConnect {
	return &ErrFailedToConnect{
		Addr: addr,
		Err:  err,
	}
}
