package otr3

import "fmt"

// ErrorCode represents an error that can happen during OTR processing
type ErrorCode int

const (
	// ErrorCodeEncryptionError means an error occured while encrypting a message
	ErrorCodeEncryptionError ErrorCode = iota

	// ErrorCodeMessageUnreadable means we received an unreadable encrypted message
	ErrorCodeMessageUnreadable

	// ErrorCodeMessageMalformed means the message sent is malformed
	ErrorCodeMessageMalformed

	// ErrorCodeMessageNotInPrivate means we received an encrypted message when not expecting it
	ErrorCodeMessageNotInPrivate
)

// ErrorMessageHandler generates error messages for error codes
type ErrorMessageHandler interface {
	// HandleErrorMessage should return a string according to the error event. This string will be concatenated to an OTR header to produce an OTR protocol error message
	HandleErrorMessage(error ErrorCode) []byte
}

type dynamicErrorMessageHandler struct {
	eh func(error ErrorCode) []byte
}

func (d dynamicErrorMessageHandler) HandleErrorMessage(error ErrorCode) []byte {
	return d.eh(error)
}

func (c *Conversation) generatePotentialErrorMessage(ec ErrorCode) {
	if c.errorMessageHandler != nil {
		msg := c.errorMessageHandler.HandleErrorMessage(ec)
		c.injectMessage(append(append(errorMarker, ' '), msg...))
	}
}

func (s ErrorCode) String() string {
	switch s {
	case ErrorCodeEncryptionError:
		return "ErrorCodeEncryptionError"
	case ErrorCodeMessageUnreadable:
		return "ErrorCodeMessageUnreadable"
	case ErrorCodeMessageMalformed:
		return "ErrorCodeMessageMalformed"
	case ErrorCodeMessageNotInPrivate:
		return "ErrorCodeMessageNotInPrivate"
	default:
		return "ERROR CODE: (THIS SHOULD NEVER HAPPEN)"
	}
}

type combinedErrorMessageHandler struct {
	handlers []ErrorMessageHandler
}

func (c combinedErrorMessageHandler) HandleErrorMessage(error ErrorCode) []byte {
	var result []byte
	for _, h := range c.handlers {
		if h != nil {
			result = h.HandleErrorMessage(error)
		}
	}
	return result
}

// CombineErrorMessageHandlers creates an ErrorMessageHandler that will call all handlers
// given to this function. It returns the result of the final handler called
func CombineErrorMessageHandlers(handlers ...ErrorMessageHandler) ErrorMessageHandler {
	return combinedErrorMessageHandler{handlers}
}

// DebugErrorMessageHandler is an ErrorMessageHandler that dumps all error message requests to standard error. It returns nil
type DebugErrorMessageHandler struct{}

// HandleErrorMessage dumps all error messages and returns nil
func (DebugErrorMessageHandler) HandleErrorMessage(error ErrorCode) []byte {
	fmt.Fprintf(standardErrorOutput, "%sHandleErrorMessage(%s)\n", debugPrefix, error)
	return nil
}
