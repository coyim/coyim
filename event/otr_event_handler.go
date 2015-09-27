package event

import "github.com/twstrike/otr3"

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

	minFragmentSize = 18
)

type OtrEventHandler struct {
	SmpQuestion      string
	securityChange   SecurityChange
	WaitingForSecret bool
}

func (OtrEventHandler) WishToHandleErrorMessage() bool {
	return true
}

func (OtrEventHandler) HandleErrorMessage(error otr3.ErrorCode) []byte {
	return nil
}

func (e *OtrEventHandler) HandleSecurityEvent(event otr3.SecurityEvent) {
	switch event {
	case otr3.GoneSecure, otr3.StillSecure:
		e.securityChange = NewKeys
	}
}

func (e *OtrEventHandler) HandleSMPEvent(event otr3.SMPEvent, progressPercent int, question string) {
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

func (e *OtrEventHandler) HandleMessageEvent(event otr3.MessageEvent, message []byte, err error) {
	if event == otr3.MessageEventConnectionEnded {
		e.securityChange = ConversationEnded
	}
}

func (e *OtrEventHandler) ConsumeSecurityChange() SecurityChange {
	ret := e.securityChange
	e.securityChange = NoChange
	return ret
}
