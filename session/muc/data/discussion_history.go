package data

import (
	"sync"
	"time"
)

// MessageType represents a message type
type MessageType int

const (
	// Chat represents a single chat message
	Chat MessageType = iota
	// Subject represents a subject message
	Subject
	// Password represents a password message
	Password
	// Joined is used when an occupant joined the room
	Joined
	// Left is used when an occupant left the room
	Left
	// Connected is used when an occupant connected to the room
	Connected
	// Disconnected is used when an occupant lost connection to the room
	Disconnected
	// OccupantInformationChanged is used when the occupant information changed
	OccupantInformationChanged
	// RoomConfigurationChanged is used when the room information changed
	RoomConfigurationChanged
)

// DelayedMessage contains the information of a delayed message
type DelayedMessage struct {
	Nickname    string
	Message     string
	Timestamp   time.Time
	MessageType MessageType
}

func newDelayedMessage(nickname, message string, timestamp time.Time, messageType MessageType) *DelayedMessage {
	return &DelayedMessage{
		Nickname:    nickname,
		Message:     message,
		Timestamp:   serverTimeInLocal(timestamp),
		MessageType: messageType,
	}
}

// DelayedMessages contains the delayed messages for specific date
type DelayedMessages struct {
	date     time.Time
	messages []*DelayedMessage
	lock     sync.RWMutex
}

func newDelayedMessages(date time.Time) *DelayedMessages {
	return &DelayedMessages{
		date: serverTimeInLocal(date),
	}
}

// GetDate returns the delayed messages group's date
func (dm *DelayedMessages) GetDate() time.Time {
	return dm.date
}

// GetMessages returns a list of delayed messages
func (dm *DelayedMessages) GetMessages() []*DelayedMessage {
	dm.lock.RLock()
	defer dm.lock.RUnlock()

	result := []*DelayedMessage{}
	result = append(result, dm.messages...)

	return result
}

func (dm *DelayedMessages) add(nickname, message string, timestamp time.Time, messageType MessageType) {
	dm.lock.Lock()
	defer dm.lock.Unlock()

	shouldAddDelayedMessage := true
	if len(dm.messages) > 0 {
		lastMessage := dm.messages[len(dm.messages)-1]
		shouldAddDelayedMessage = lastMessage.Timestamp.Before(timestamp)
	}

	if shouldAddDelayedMessage {
		dm.messages = append(dm.messages, newDelayedMessage(nickname, message, timestamp, messageType))
	}
}

// DiscussionHistory contains the rooms's discussion history
type DiscussionHistory struct {
	history []*DelayedMessages
	lock    sync.RWMutex
}

// NewDiscussionHistory creates a new discussion history instance
func NewDiscussionHistory() *DiscussionHistory {
	return &DiscussionHistory{}
}

// GetHistory returns the discussion history
func (dh *DiscussionHistory) GetHistory() []*DelayedMessages {
	dh.lock.RLock()
	defer dh.lock.RUnlock()

	return append([]*DelayedMessages{}, dh.history...)
}

// AddMessage add a new delayed message to the history
func (dh *DiscussionHistory) AddMessage(nickname, message string, timestamp time.Time, messageType MessageType) {
	t := serverTimeInLocal(timestamp)

	for _, dm := range dh.GetHistory() {
		if sameDate(dm.date, t) {
			dm.add(nickname, message, t, messageType)
			return
		}
	}

	dm := dh.addNewMessagesGroup(t)
	dm.add(nickname, message, t, messageType)
}

func (dh *DiscussionHistory) addNewMessagesGroup(date time.Time) *DelayedMessages {
	dh.lock.Lock()
	defer dh.lock.Unlock()

	dm := newDelayedMessages(date)
	dh.history = append(dh.history, dm)

	return dm
}

func sameDate(d1, d2 time.Time) bool {
	t1y, t1m, t1d := d1.In(time.Local).Date()
	t2y, t2m, t2d := d2.In(time.Local).Date()

	return t1d == t2d && t1m == t2m && t1y == t2y
}

func serverTimeInLocal(t time.Time) time.Time {
	return t.In(time.Local)
}
