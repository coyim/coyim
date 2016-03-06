package otr3

import "fmt"

// SecurityEvent define the events used to indicate changes in security status. In comparison with libotr, this library does not take trust levels into concern for security events
type SecurityEvent int

const (
	// GoneInsecure is signalled when we have gone from a secure state to an insecure state
	GoneInsecure SecurityEvent = iota
	// GoneSecure is signalled when we have gone from an insecure state to a secure state
	GoneSecure
	// StillSecure is signalled when we have refreshed the security state but is still in a secure state
	StillSecure
)

// SecurityEventHandler is an interface for events that are related to changes of security status
type SecurityEventHandler interface {
	// HandleSecurityEvent is called when a change in security status happens
	HandleSecurityEvent(event SecurityEvent)
}

type dynamicSecurityEventHandler struct {
	eh func(event SecurityEvent)
}

func (d dynamicSecurityEventHandler) HandleSecurityEvent(event SecurityEvent) {
	d.eh(event)
}

func (c *Conversation) securityEvent(e SecurityEvent) {
	if c.securityEventHandler != nil {
		c.securityEventHandler.HandleSecurityEvent(e)
	}
}

// String returns the string representation of the SecurityEvent
func (s SecurityEvent) String() string {
	switch s {
	case GoneInsecure:
		return "GoneInsecure"
	case GoneSecure:
		return "GoneSecure"
	case StillSecure:
		return "StillSecure"
	default:
		return "SECURITY EVENT: (THIS SHOULD NEVER HAPPEN)"
	}
}

type combinedSecurityEventHandler struct {
	handlers []SecurityEventHandler
}

func (c combinedSecurityEventHandler) HandleSecurityEvent(event SecurityEvent) {
	for _, h := range c.handlers {
		if h != nil {
			h.HandleSecurityEvent(event)
		}
	}
}

// CombineSecurityEventHandlers creates a SecurityEventHandler that will call all handlers
// given to this function. It ignores nil entries.
func CombineSecurityEventHandlers(handlers ...SecurityEventHandler) SecurityEventHandler {
	return combinedSecurityEventHandler{handlers}
}

func (c *Conversation) signalSecurityEventIf(cond bool, event SecurityEvent) {
	if cond {
		c.securityEvent(event)
	}
}

// DebugSecurityEventHandler is a SecurityEventHandler that dumps all SecurityEvents to standard error
type DebugSecurityEventHandler struct{}

// HandleSecurityEvent dumps all security events
func (DebugSecurityEventHandler) HandleSecurityEvent(event SecurityEvent) {
	fmt.Fprintf(standardErrorOutput, "%sHandleSecurityEvent(%s)\n", debugPrefix, event)
}
