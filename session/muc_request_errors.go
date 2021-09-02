package session

import (
	"errors"

	"github.com/coyim/coyim/coylog"
)

var (
	// ErrInvalidInformationQueryRequest is an invalid information query request error
	ErrInvalidInformationQueryRequest = errors.New("invalid information query request")
	// ErrUnexpectedResponse is an unexpected response from the server error
	ErrUnexpectedResponse = errors.New("received an unexpected response from the server")
	// ErrInformationQueryResponse contains an error received in the information query response
	ErrInformationQueryResponse = errors.New("received an error from the server")
	// ErrBadRequestResponse contains an error received in the information query response
	ErrBadRequestResponse = errors.New("received a bad request error from the server")
	// ErrInternalServerErrorResponse contains an internal server error received in the information query response
	ErrInternalServerErrorResponse = errors.New("received an internal server error")
	// ErrInformationQueryResponseWithGoneTag contains a gone tag inside of an error received through an information query
	ErrInformationQueryResponseWithGoneTag = errors.New("received a gone tag in the information query response")
)

var mucRequestErrorMessages = map[error]string{
	ErrUnexpectedResponse:             "Unexpected information query response",
	ErrInvalidInformationQueryRequest: "Unexpected information query reply",
	ErrInvalidReserveRoomRequest:      "An error occurred while reserving the room",
}

func mucRequestErrorMessage(err error) string {
	if message, ok := mucRequestErrorMessages[err]; ok {
		return message
	}
	return "Something wrong happened during the request"
}

type mucRequestError struct {
	err     error
	message string
	log     coylog.Logger
}

func (r *mucRequest) newMUCRoomRequestError(err error) *mucRequestError {
	return &mucRequestError{
		err,
		mucRequestErrorMessage(err),
		r.log,
	}
}

func (e *mucRequestError) logError() {
	e.log.WithError(e.err).Error(e.message)
}
