package event

import (
	"log"

	"../i18n"
	"github.com/twstrike/otr3"
)

// SecurityChange describes a change in the security state of a Conversation.
type SecurityChange int

const (
	// NoChange happened in the security status
	NoChange SecurityChange = iota
	// NewKeys indicates that a key exchange has completed. This occurs
	// when a conversation first becomes encrypted, and when the keys are
	// renegotiated within an encrypted conversation.
	NewKeys
	// SMPSecretNeeded indicates that the peer has started an
	// authentication and that we need to supply a secret. Call SMPQuestion
	// to get the optional, human readable challenge and then Authenticate
	// to supply the matching secret.
	SMPSecretNeeded
	// SMPComplete indicates that an authentication completed. The identity
	// of the peer has now been confirmed.
	SMPComplete
	// SMPFailed indicates that an authentication failed.
	SMPFailed
	// ConversationEnded indicates that the peer ended the secure
	// conversation.
	ConversationEnded
)

var (
	// QueryMessage can be sent to a peer to start an OTR conversation.
	QueryMessage = "?OTRv2?"

	// ErrorPrefix can be used to make an OTR error by appending an error message
	// to it.
	ErrorPrefix = "?OTR Error:"
)

// OtrEventHandler is used to contain information pertaining to the events of a specific OTR interaction
type OtrEventHandler struct {
	SmpQuestion      string
	securityChange   SecurityChange
	WaitingForSecret bool
}

// WishToHandleErrorMessage means this event handler wants to handle error messages
func (*OtrEventHandler) WishToHandleErrorMessage() bool {
	return true
}

// HandleErrorMessage is called when asked to handle a specific error message
func (e *OtrEventHandler) HandleErrorMessage(error otr3.ErrorCode) []byte {
	log.Printf("HandleErrorMessage(%s)", error.String())

	switch error {
	case otr3.ErrorCodeEncryptionError:
		return []byte(i18n.Local("Error occurred encrypting message."))
	case otr3.ErrorCodeMessageUnreadable:
		return []byte(i18n.Local("You transmitted an unreadable encrypted message."))
	case otr3.ErrorCodeMessageMalformed:
		return []byte(i18n.Local("You transmitted a malformed data message."))
	case otr3.ErrorCodeMessageNotInPrivate:
		return []byte(i18n.Local("You sent encrypted data to a peer, who wasn't expecting it."))
	}

	return nil
}

// HandleSecurityEvent is called to handle a specific security event
func (e *OtrEventHandler) HandleSecurityEvent(event otr3.SecurityEvent) {
	log.Printf("HandleSecurityEvent(%s)", event.String())
	switch event {
	case otr3.GoneSecure, otr3.StillSecure:
		e.securityChange = NewKeys
	case otr3.GoneInsecure:
		e.securityChange = ConversationEnded
	}
}

// HandleSMPEvent is called to handle a specific SMP event
func (e *OtrEventHandler) HandleSMPEvent(event otr3.SMPEvent, progressPercent int, question string) {
	log.Printf("HandleSMPEvent(%s, %d, %s)", event.String(), progressPercent, question)
	switch event {
	case otr3.SMPEventAskForSecret, otr3.SMPEventAskForAnswer:
		e.securityChange = SMPSecretNeeded
		e.SmpQuestion = question
		e.WaitingForSecret = true
	case otr3.SMPEventSuccess:
		if progressPercent == 100 {
			e.securityChange = SMPComplete
		}
	case otr3.SMPEventAbort, otr3.SMPEventFailure, otr3.SMPEventCheated:
		e.securityChange = SMPFailed
	}
}

// HandleMessageEvent is called to handle a specific message event
func (e *OtrEventHandler) HandleMessageEvent(event otr3.MessageEvent, message []byte, err error) {
	log.Printf("HandleMessageEvent(%s, %s, %v)", event.String(), message, err)
}

// ConsumeSecurityChange is called to get the current security change and forget the old one
func (e *OtrEventHandler) ConsumeSecurityChange() SecurityChange {
	ret := e.securityChange
	e.securityChange = NoChange
	return ret
}
