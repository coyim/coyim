package gui

import (
	"time"

	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/session/muc/data"
)

type occupantSelfJoinedEvent struct {
	nickname string
	role     data.Role
}

type occupantLeftEvent struct {
	nickname string
}

type occupantJoinedEvent struct {
	nickname string
}

type occupantUpdatedEvent struct {
	nickname string
	role     data.Role
}

type nicknameConflictEvent struct {
	nickname string
}

type registrationRequiredEvent struct {
	nickname string
}

type roomInfoReceivedEvent struct {
	info *muc.RoomListing
}

type roomInfoTimeoutEvent struct{}

type loggingEnabledEvent struct{}

type loggingDisabledEvent struct{}

type nonAnonymousRoomEvent struct{}

type semiAnonymousRoomEvent struct{}

type messageEvent struct {
	tp        string
	nickname  string
	message   string
	timestamp time.Time
}

type subjectUpdatedEvent struct {
	nickname string
	subject  string
}

type subjectReceivedEvent struct {
	subject string
}

type messageForbidden struct{}

type messageNotAcceptable struct{}

type discussionHistoryEvent struct {
	history *data.DiscussionHistory
}

type roomViewEvent interface {
	markAsRoomViewEvent()
}

type roomConfigurationChanged struct {
	oldConfiguration *muc.RoomListing
	newConfiguration *muc.RoomListing
}

func (occupantLeftEvent) markAsRoomViewEvent()         {}
func (occupantJoinedEvent) markAsRoomViewEvent()       {}
func (occupantUpdatedEvent) markAsRoomViewEvent()      {}
func (occupantSelfJoinedEvent) markAsRoomViewEvent()   {}
func (messageEvent) markAsRoomViewEvent()              {}
func (subjectUpdatedEvent) markAsRoomViewEvent()       {}
func (subjectReceivedEvent) markAsRoomViewEvent()      {}
func (nicknameConflictEvent) markAsRoomViewEvent()     {}
func (registrationRequiredEvent) markAsRoomViewEvent() {}
func (roomInfoReceivedEvent) markAsRoomViewEvent()     {}
func (roomInfoTimeoutEvent) markAsRoomViewEvent()      {}
func (loggingEnabledEvent) markAsRoomViewEvent()       {}
func (loggingDisabledEvent) markAsRoomViewEvent()      {}
func (nonAnonymousRoomEvent) markAsRoomViewEvent()     {}
func (semiAnonymousRoomEvent) markAsRoomViewEvent()    {}
func (roomConfigurationChanged) markAsRoomViewEvent()  {}
func (messageForbidden) markAsRoomViewEvent()          {}
func (messageNotAcceptable) markAsRoomViewEvent()      {}
func (discussionHistoryEvent) markAsRoomViewEvent()    {}
