package otr

import "github.com/twstrike/otr3"

var (
	// QueryMessage can be sent to a peer to start an OTR conversation.
	QueryMessage = "?OTRv2?"

	// ErrorPrefix can be used to make an OTR error by appending an error message
	// to it.
	ErrorPrefix = "?OTR Error:"

	minFragmentSize = 18
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

type EventHandler struct {
	SmpQuestion      string
	securityChange   SecurityChange
	waitingForSecret bool
}

func (EventHandler) WishToHandleErrorMessage() bool {
	return true
}

func (EventHandler) HandleErrorMessage(error otr3.ErrorCode) []byte {
	return nil
}

func (e *EventHandler) HandleSecurityEvent(event otr3.SecurityEvent) {
	switch event {
	case otr3.GoneSecure, otr3.StillSecure:
		e.securityChange = NewKeys
	}
}

func (e *EventHandler) HandleSMPEvent(event otr3.SMPEvent, progressPercent int, question string) {
	switch event {
	case otr3.SMPEventAskForSecret, otr3.SMPEventAskForAnswer:
		e.securityChange = SMPSecretNeeded
		e.SmpQuestion = question
		e.waitingForSecret = true
	case otr3.SMPEventSuccess:
		if progressPercent == 100 {
			e.securityChange = SMPComplete
		}
	case otr3.SMPEventAbort, otr3.SMPEventFailure, otr3.SMPEventCheated:
		e.securityChange = SMPFailed
	}
}

func (e *EventHandler) HandleMessageEvent(event otr3.MessageEvent, message []byte, err error) {
	if event == otr3.MessageEventConnectionEnded {
		e.securityChange = ConversationEnded
	}
}

func (e *EventHandler) ConsumeSecurityChange() SecurityChange {
	ret := e.securityChange
	e.securityChange = NoChange
	return ret
}
