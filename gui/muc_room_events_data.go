package gui

import (
	"time"

	"github.com/coyim/coyim/xmpp/jid"

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

type selfOccupantRemovedEvent struct{}

type occupantRemovedEvent struct {
	nickname string
}

type nicknameConflictEvent struct {
	nickname string
}

type registrationRequiredEvent struct {
	nickname string
}

type loggingEnabledEvent struct{}

type loggingDisabledEvent struct{}

type roomAnonymityEvent struct {
	anonymityType data.RoomAnonymityType
}

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

type roomDestroyedEvent struct {
	reason      string
	alternative jid.Bare
}

type messageForbidden struct{}

type messageNotAcceptable struct{}

type discussionHistoryEvent struct {
	history *data.DiscussionHistory
}

type roomViewEvent interface {
	markAsRoomViewEvent()
}

type roomConfigReceivedEvent struct {
	config data.RoomConfig
}

type roomConfigRequestTimeoutEvent struct{}

type roomConfigChangedEvent struct {
	changes []data.RoomConfigType
	config  data.RoomConfig
}

func (occupantLeftEvent) markAsRoomViewEvent()             {}
func (occupantJoinedEvent) markAsRoomViewEvent()           {}
func (occupantUpdatedEvent) markAsRoomViewEvent()          {}
func (occupantSelfJoinedEvent) markAsRoomViewEvent()       {}
func (messageEvent) markAsRoomViewEvent()                  {}
func (subjectUpdatedEvent) markAsRoomViewEvent()           {}
func (subjectReceivedEvent) markAsRoomViewEvent()          {}
func (nicknameConflictEvent) markAsRoomViewEvent()         {}
func (registrationRequiredEvent) markAsRoomViewEvent()     {}
func (loggingEnabledEvent) markAsRoomViewEvent()           {}
func (loggingDisabledEvent) markAsRoomViewEvent()          {}
func (roomAnonymityEvent) markAsRoomViewEvent()            {}
func (messageForbidden) markAsRoomViewEvent()              {}
func (messageNotAcceptable) markAsRoomViewEvent()          {}
func (discussionHistoryEvent) markAsRoomViewEvent()        {}
func (roomConfigReceivedEvent) markAsRoomViewEvent()       {}
func (roomConfigRequestTimeoutEvent) markAsRoomViewEvent() {}
func (roomDestroyedEvent) markAsRoomViewEvent()            {}
func (roomConfigChangedEvent) markAsRoomViewEvent()        {}
func (selfOccupantRemovedEvent) markAsRoomViewEvent()      {}
func (occupantRemovedEvent) markAsRoomViewEvent()          {}
