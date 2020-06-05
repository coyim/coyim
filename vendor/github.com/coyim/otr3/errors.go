package otr3

import "fmt"

var errCantAuthenticateWithoutEncryption = newOtrError("can't authenticate a peer without a secure conversation established")
var errCorruptEncryptedSignature = newOtrError("corrupt encrypted signature")
var errInvalidOTRMessage = newOtrError("invalid OTR message")
var errInvalidVersion = newOtrError("no valid version agreement could be found") //libotr ignores this situation
var errNotWaitingForSMPSecret = newOtrError("not expected SMP secret to be provided now")
var errReceivedMessageForOtherInstance = newOtrError("received message for other OTR instance") //not exactly an error - we should ignore these messages by default
var errShortRandomRead = newOtrError("short read from random source")
var errUnsupportedOTRVersion = newOtrError("unsupported OTR version")
var errWrongProtocolVersion = newOtrError("wrong protocol version")
var errMessageNotInPrivate = newOtrError("message not in private")
var errCannotSendUnencrypted = newOtrConflictError("cannot send message in unencrypted state")

// OtrError is an error in the OTR library
type OtrError struct {
	msg      string
	conflict bool
}

func newOtrError(s string) error {
	return OtrError{msg: s, conflict: false}
}

func newOtrConflictError(s string) error {
	return OtrError{msg: s, conflict: true}
}

func newOtrErrorf(format string, a ...interface{}) error {
	return OtrError{msg: fmt.Sprintf(format, a...), conflict: false}
}

func (oe OtrError) Error() string {
	return "otr: " + oe.msg
}

func firstError(es ...error) error {
	for _, e := range es {
		if e != nil {
			return e
		}
	}
	return nil
}

func isConflict(e error) bool {
	if oe, ok := e.(OtrError); ok {
		return oe.conflict
	}
	return false
}
