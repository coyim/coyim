package gui

import (
	"errors"

	"github.com/coyim/coyim/i18n"
)

var (
	errJoinRequestFailed    = errors.New("the request to join the room has failed")
	errInvalidNickname      = errors.New("not valid nickname")
	errJoinNoConnection     = errors.New("join request failed because maybe no connection available")
	errJoinNicknameConflict = errors.New("join failed because the nickname is being used")
	errJoinOnlyMembers      = errors.New("join failed because only registered members are allowed")
	errJoinNotAuthorized    = errors.New("join failed because doesn't have authorization")
	errServiceUnavailable   = errors.New("join failed because the service is unavailable")
	errUnknownError         = errors.New("join failed because an unknown error occurred")
	errOccupantForbidden    = errors.New("join failed because the occupant is banned")
)

type roomViewCustomError struct {
	nickname        string
	friendlyMessage string
	errType         error
}

// Error implements the `error` interface
func (e *roomViewCustomError) Error() string {
	return e.errType.Error()
}

func (v *roomView) newCustomRoomViewError(nickname string, errType error) error {
	return &roomViewCustomError{
		nickname:        nickname,
		friendlyMessage: v.userFriendlyRoomErrorMessage(errType),
		errType:         errType,
	}
}

func (v *roomView) invalidNicknameError() error {
	return v.newCustomRoomViewError("", errInvalidNickname)
}

func (v *roomView) joinRequestError(err error) {
	joinErr := v.newCustomRoomViewError("", err)
	v.finishJoinRequestWithError(joinErr)
}

func (v *roomView) nicknameConflictEvent(nickname string) {
	nicknameErr := v.newCustomRoomViewError(nickname, errJoinNicknameConflict)
	v.finishJoinRequestWithError(nicknameErr)
}

func (v *roomView) registrationRequiredEvent() {
	v.joinRequestError(errJoinOnlyMembers)
}

func (v *roomView) notAuthorizedEvent() {
	v.joinRequestError(errJoinNotAuthorized)
}

func (v *roomView) serviceUnavailableEvent() {
	v.joinRequestError(errServiceUnavailable)
}

func (v *roomView) unknownErrorEvent() {
	v.joinRequestError(errUnknownError)
}

func (v *roomView) occupantForbiddenEvent() {
	v.joinRequestError(errOccupantForbidden)
}

func (v *roomView) userFriendlyRoomErrorMessage(err error) string {
	switch err {
	case errInvalidNickname:
		return i18n.Local("You must provide a valid nickname.")
	case errJoinNicknameConflict:
		return i18n.Local("You can't join the room using that nickname because it's already being used.")
	case errJoinOnlyMembers:
		return i18n.Local("Sorry, this room only allows registered members.")
	case errJoinNotAuthorized:
		return i18n.Local("You can't join the room because the password is not valid.")
	case errServiceUnavailable:
		return i18n.Local("You can't join the room because the maximum number of occupants has been reached.")
	case errUnknownError:
		return i18n.Local("An unknown error occurred when trying to join the room. Please try again later.")
	case errOccupantForbidden:
		return i18n.Local("You can't join the room because your account is currently banned.")
	}
	return i18n.Local("An error occurred when trying to join the room. Please check your connection or make sure the room exists.")
}
