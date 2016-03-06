package otr3

import "fmt"

// MessageEvent define the events used to indicate the messages that need to be sent
type MessageEvent int

const (
	// MessageEventEncryptionRequired is signaled when our policy requires encryption but we are trying to send an unencrypted message.
	MessageEventEncryptionRequired MessageEvent = iota

	// MessageEventEncryptionError is signaled when an error occured while encrypting a message and the message was not sent.
	MessageEventEncryptionError

	// MessageEventConnectionEnded is signaled when we are asked to send a message but the peer has ended the private conversation.
	// At this point the connection should be closed or refreshed.
	MessageEventConnectionEnded

	// MessageEventSetupError will be signaled when a private conversation could not be established. The reason for this will be communicated with the attached error instance.
	MessageEventSetupError

	// MessageEventMessageReflected will be signaled if we received our own OTR messages.
	MessageEventMessageReflected

	// MessageEventMessageResent is signaled when a message is resent
	MessageEventMessageResent

	// MessageEventReceivedMessageNotInPrivate will be signaled when we receive an encrypted message that we cannot read, because we don't have an established private connection
	MessageEventReceivedMessageNotInPrivate

	// MessageEventReceivedMessageUnreadable will be signaled when we cannot read the received message.
	MessageEventReceivedMessageUnreadable

	// MessageEventReceivedMessageMalformed is signaled when we receive a message that contains malformed data.
	MessageEventReceivedMessageMalformed

	// MessageEventLogHeartbeatReceived is triggered when we received a heartbeat.
	MessageEventLogHeartbeatReceived

	// MessageEventLogHeartbeatSent is triggered when we have sent a heartbeat.
	MessageEventLogHeartbeatSent

	// MessageEventReceivedMessageGeneralError will be signaled when we receive an OTR error from the peer.
	// The message parameter will be passed, containing the error message
	MessageEventReceivedMessageGeneralError

	// MessageEventReceivedMessageUnencrypted is triggered when we receive a message that was sent in the clear when it should have been encrypted.
	// The actual message received will also be passed.
	MessageEventReceivedMessageUnencrypted

	// MessageEventReceivedMessageUnrecognized is triggered when we receive an OTR message whose type we cannot recognize
	MessageEventReceivedMessageUnrecognized

	// MessageEventReceivedMessageForOtherInstance is triggered when we receive and discard a message for another instance
	MessageEventReceivedMessageForOtherInstance
)

// MessageEventHandler handles MessageEvents
type MessageEventHandler interface {
	// HandleMessageEvent should handle and send the appropriate message(s) to the sender/recipient depending on the message events
	HandleMessageEvent(event MessageEvent, message []byte, err error)
}

type dynamicMessageEventHandler struct {
	eh func(event MessageEvent, message []byte, err error)
}

func (d dynamicMessageEventHandler) HandleMessageEvent(event MessageEvent, message []byte, err error) {
	d.eh(event, message, err)
}

func (c *Conversation) messageEvent(e MessageEvent) {
	if c.messageEventHandler != nil {
		c.messageEventHandler.HandleMessageEvent(e, nil, nil)
	}
}

func (c *Conversation) messageEventWithError(e MessageEvent, err error) {
	if c.messageEventHandler != nil {
		c.messageEventHandler.HandleMessageEvent(e, nil, err)
	}
}

func (c *Conversation) messageEventWithMessage(e MessageEvent, msg []byte) {
	if c.messageEventHandler != nil {
		c.messageEventHandler.HandleMessageEvent(e, msg, nil)
	}
}

// String returns the string representation of the MessageEvent
func (s MessageEvent) String() string {
	switch s {
	case MessageEventEncryptionRequired:
		return "MessageEventEncryptionRequired"
	case MessageEventEncryptionError:
		return "MessageEventEncryptionError"
	case MessageEventConnectionEnded:
		return "MessageEventConnectionEnded"
	case MessageEventSetupError:
		return "MessageEventSetupError"
	case MessageEventMessageReflected:
		return "MessageEventMessageReflected"
	case MessageEventMessageResent:
		return "MessageEventMessageResent"
	case MessageEventReceivedMessageNotInPrivate:
		return "MessageEventReceivedMessageNotInPrivate"
	case MessageEventReceivedMessageUnreadable:
		return "MessageEventReceivedMessageUnreadable"
	case MessageEventReceivedMessageMalformed:
		return "MessageEventReceivedMessageMalformed"
	case MessageEventLogHeartbeatReceived:
		return "MessageEventLogHeartbeatReceived"
	case MessageEventLogHeartbeatSent:
		return "MessageEventLogHeartbeatSent"
	case MessageEventReceivedMessageGeneralError:
		return "MessageEventReceivedMessageGeneralError"
	case MessageEventReceivedMessageUnencrypted:
		return "MessageEventReceivedMessageUnencrypted"
	case MessageEventReceivedMessageUnrecognized:
		return "MessageEventReceivedMessageUnrecognized"
	case MessageEventReceivedMessageForOtherInstance:
		return "MessageEventReceivedMessageForOtherInstance"
	default:
		return "MESSAGE EVENT: (THIS SHOULD NEVER HAPPEN)"
	}
}

type combinedMessageEventHandler struct {
	handlers []MessageEventHandler
}

func (c combinedMessageEventHandler) HandleMessageEvent(event MessageEvent, message []byte, err error) {
	for _, h := range c.handlers {
		if h != nil {
			h.HandleMessageEvent(event, message, err)
		}
	}
}

// CombineMessageEventHandlers creates a MessageEventHandler that will call all handlers
// given to this function. It ignores nil entries.
func CombineMessageEventHandlers(handlers ...MessageEventHandler) MessageEventHandler {
	return combinedMessageEventHandler{handlers}
}

// DebugMessageEventHandler is a MessageEventHandler that dumps all MessageEvents to standard error
type DebugMessageEventHandler struct{}

// HandleMessageEvent dumps all message events
func (DebugMessageEventHandler) HandleMessageEvent(event MessageEvent, message []byte, err error) {
	fmt.Fprintf(standardErrorOutput, "%sHandleMessageEvent(%s, message: %#v, error: %v)\n", debugPrefix, event, string(message), err)
}
