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

type mucRoomLobbyErr struct {
	nickname string
	errType  error
}

func (e *mucRoomLobbyErr) Error() string {
	return e.errType.Error()
}

func newMUCRoomLobbyErr(nickname string, errType error) error {
	return &mucRoomLobbyErr{
		nickname: nickname,
		errType:  errType,
	}
}

func newRoomLobbyInvalidNicknameError() error {
	return newMUCRoomLobbyErr("", errInvalidNickname)
}

func (l *roomViewLobby) joinRequestErrorEvent(err error) {
	l.finishJoinRequestWithError(newMUCRoomLobbyErr("", err))
}

func (l *roomViewLobby) nicknameConflictEvent(nickname string) {
	l.joinRequestErrorEvent(newMUCRoomLobbyErr(nickname, errJoinNicknameConflict))
}

func (l *roomViewLobby) registrationRequiredEvent() {
	l.joinRequestErrorEvent(errJoinOnlyMembers)
}

func (l *roomViewLobby) notAuthorizedEvent() {
	l.joinRequestErrorEvent(errJoinNotAuthorized)
}

func (l *roomViewLobby) serviceUnavailableEvent() {
	l.joinRequestErrorEvent(errServiceUnavailable)
}

func (l *roomViewLobby) unknownErrorEvent() {
	l.joinRequestErrorEvent(errUnknownError)
}

func (l *roomViewLobby) occupantForbiddenEvent() {
	l.joinRequestErrorEvent(errOccupantForbidden)
}

func (l *roomViewLobby) getUserErrorMessage(err *mucRoomLobbyErr) string {
	switch err.errType {
	case errInvalidNickname:
		return i18n.Local("You must provide a valid nickname")
	case errJoinNicknameConflict:
		return i18n.Local("Can't join the room using that nickname because it's already being used")
	case errJoinOnlyMembers:
		return i18n.Local("Sorry, this room only allows registered members")
	case errJoinNotAuthorized:
		return i18n.Local("Invalid password")
	case errServiceUnavailable:
		return i18n.Local("Can't join the room because the maximun number of occupants has been reached")
	case errUnknownError:
		return i18n.Local("An unknown error occurred while trying to join the room, please try again later")
	case errOccupantForbidden:
		return i18n.Local("Can't join the room because you are banned")
	default:
		return i18n.Local("An error occurred while trying to join the room, please check your connection or make sure the room exists")
	}
}
