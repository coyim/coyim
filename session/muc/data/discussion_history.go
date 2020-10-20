package data

import (
	"sync"
	"time"
)

// DelayedMessage contains the information of a delayed message
type DelayedMessage struct {
	Nickname  string
	Message   string
	Timestamp time.Time
}

func newDelayedMessage(nickname, message string, timestamp time.Time) *DelayedMessage {
	return &DelayedMessage{
		Nickname:  nickname,
		Message:   message,
		Timestamp: serverTimestampInLocalTime(timestamp),
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
		date: serverTimestampInLocalTime(date),
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

	for _, m := range dm.messages {
		result = append(result, m)
	}

	return result
}

func (dm *DelayedMessages) add(nickname, message string, timestamp time.Time) {
	dm.lock.Lock()
	dm.messages = append(dm.messages, newDelayedMessage(nickname, message, timestamp))
	dm.lock.Unlock()
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

	result := []*DelayedMessages{}

	for _, h := range dh.history {
		result = append(result, h)
	}

	return result
}

// AddMessage add a new delayed message to the history
func (dh *DiscussionHistory) AddMessage(nickname, message string, timestamp time.Time) {
	for _, dm := range dh.GetHistory() {
		if areTheSameDate(dm.date, serverTimestampInLocalTime(timestamp)) {
			dm.add(nickname, message, timestamp)
			return
		}
	}

	dh.addNewMessagesGroup(nickname, message, timestamp)
}

func (dh *DiscussionHistory) addNewMessagesGroup(nickname, message string, timestamp time.Time) {
	dm := newDelayedMessages(timestamp)
	dm.add(nickname, message, timestamp)

	dh.lock.Lock()
	dh.history = append(dh.history, dm)
	dh.lock.Unlock()
}

func areTheSameDate(d1, d2 time.Time) bool {
	return d1.Format("02 Jan 06") == d2.Format("02 Jan 06")
}

func serverTimestampInLocalTime(timestamp time.Time) time.Time {
	return timestamp.In(time.Now().Location())
}
