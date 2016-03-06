package otr3

import "fmt"

// SMPEvent define the events used to indicate status of SMP to the UI
type SMPEvent int

const (
	// SMPEventError means abort the current auth and update the auth progress dialog with progress_percent. This event is only sent when we receive a message for another message state than we are in
	SMPEventError SMPEvent = iota
	// SMPEventAbort means update the auth progress dialog with progress_percent
	SMPEventAbort
	// SMPEventCheated means abort the current auth and update the auth progress dialog with progress_percent
	SMPEventCheated
	// SMPEventAskForAnswer means ask the user to answer the secret question
	SMPEventAskForAnswer
	// SMPEventAskForSecret means prompt the user to enter a shared secret
	SMPEventAskForSecret
	// SMPEventInProgress means update the auth progress dialog with progress_percent
	SMPEventInProgress
	// SMPEventSuccess means update the auth progress dialog with progress_percent
	SMPEventSuccess
	// SMPEventFailure means update the auth progress dialog with progress_percent
	SMPEventFailure
)

// SMPEventHandler handles SMPEvents
type SMPEventHandler interface {
	// HandleSMPEvent should update the authentication UI with respect to SMP events
	HandleSMPEvent(event SMPEvent, progressPercent int, question string)
}

type dynamicSMPEventHandler struct {
	eh func(event SMPEvent, progressPercent int, question string)
}

func (d dynamicSMPEventHandler) HandleSMPEvent(event SMPEvent, pp int, question string) {
	d.eh(event, pp, question)
}

func (c *Conversation) smpEvent(e SMPEvent, percent int) {
	if c.smpEventHandler != nil {
		c.smpEventHandler.HandleSMPEvent(e, percent, "")
	}
}

func (c *Conversation) smpEventWithQuestion(e SMPEvent, percent int, question string) {
	if c.smpEventHandler != nil {
		c.smpEventHandler.HandleSMPEvent(e, percent, question)
	}
}

func (s SMPEvent) String() string {
	switch s {
	case SMPEventError:
		return "SMPEventError"
	case SMPEventAbort:
		return "SMPEventAbort"
	case SMPEventCheated:
		return "SMPEventCheated"
	case SMPEventAskForAnswer:
		return "SMPEventAskForAnswer"
	case SMPEventAskForSecret:
		return "SMPEventAskForSecret"
	case SMPEventInProgress:
		return "SMPEventInProgress"
	case SMPEventSuccess:
		return "SMPEventSuccess"
	case SMPEventFailure:
		return "SMPEventFailure"
	default:
		return "SMP EVENT: (THIS SHOULD NEVER HAPPEN)"
	}
}

type combinedSMPEventHandler struct {
	handlers []SMPEventHandler
}

func (c combinedSMPEventHandler) HandleSMPEvent(event SMPEvent, progressPercent int, question string) {
	for _, h := range c.handlers {
		if h != nil {
			h.HandleSMPEvent(event, progressPercent, question)
		}
	}
}

// CombineSMPEventHandlers creates a SMPEventHandler that will call all handlers
// given to this function. It ignores nil entries.
func CombineSMPEventHandlers(handlers ...SMPEventHandler) SMPEventHandler {
	return combinedSMPEventHandler{handlers}
}

// DebugSMPEventHandler is an SMPEventHandler that dumps all SMPEvents to standard error
type DebugSMPEventHandler struct{}

// HandleSMPEvent dumps all SMP events
func (DebugSMPEventHandler) HandleSMPEvent(event SMPEvent, progressPercent int, question string) {
	fmt.Fprintf(standardErrorOutput, "%sHandleSMPEvent(%s, %d, %#v)\n", debugPrefix, event, progressPercent, question)
}
