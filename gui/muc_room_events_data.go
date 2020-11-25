package gui

import (
	"time"

	"github.com/coyim/coyim/xmpp/jid"

	"github.com/coyim/coyim/session/muc/data"
)

type selfOccupantRemovedEvent struct{}

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

type occupantRemovedEvent struct {
	nickname string
}

type nicknameConflictEvent struct {
	nickname string
}

type registrationRequiredEvent struct {
	nickname string
}

type notAuthorizedEvent struct{}

type loggingEnabledEvent struct{}

type loggingDisabledEvent struct{}

type roomAnonymityEvent struct {
	anonymityLevel string
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
	password    string
}

type messageForbidden struct{}

type messageNotAcceptable struct{}

type discussionHistoryEvent struct {
	history *data.DiscussionHistory
}

type roomViewEvent interface {
	markAsRoomViewEvent()
}

type roomDiscoInfoReceivedEvent struct {
	info data.RoomDiscoInfo
}

type roomConfigRequestTimeoutEvent struct{}

type roomConfigChangedEvent struct {
	changes   []data.RoomConfigType
	discoInfo data.RoomDiscoInfo
}

func (selfOccupantRemovedEvent) markAsRoomViewEvent()      {}
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
func (roomDiscoInfoReceivedEvent) markAsRoomViewEvent()    {}
func (roomConfigRequestTimeoutEvent) markAsRoomViewEvent() {}
func (roomDestroyedEvent) markAsRoomViewEvent()            {}
func (roomConfigChangedEvent) markAsRoomViewEvent()        {}
func (occupantRemovedEvent) markAsRoomViewEvent()          {}
func (notAuthorizedEvent) markAsRoomViewEvent()            {}
