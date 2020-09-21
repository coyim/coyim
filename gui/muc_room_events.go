package gui

import (
	"sync"

	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/xmpp/jid"
	log "github.com/sirupsen/logrus"
)

type roomViewEventInfo struct {
	nickname string
	subject  string
	message  string
}

func (e *roomViewEventInfo) whichNickname() string {
	return e.nickname
}

func (e *roomViewEventInfo) whichSubject() string {
	return e.subject
}

func (e *roomViewEventInfo) whichMessage() string {
	return e.message
}

type roomViewEventType int

const (
	occupantSelfJoined roomViewEventType = iota
	occupantJoined
	occupantUpdated
	occupantLeft

	roomInfoReceived

	messageReceived

	previousToSwitchToLobby
	previousToSwitchToMain
)

type roomViewObserver struct {
	id       string
	onNotify func(roomViewEventInfo)
}

type roomViewSubscribers struct {
	sync.RWMutex
	log       coylog.Logger
	observers map[roomViewEventType][]*roomViewObserver
}

func newRoomViewSubscribers(room jid.Bare, l coylog.Logger) *roomViewSubscribers {
	return &roomViewSubscribers{
		log: l.WithFields(log.Fields{
			"who":  "newRoomViewSubscribers",
			"room": room,
		}),
		observers: make(map[roomViewEventType][]*roomViewObserver),
	}
}

func (s *roomViewSubscribers) subscribe(ev roomViewEventType, o *roomViewObserver) {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.observers[ev]; !ok {
		s.observers[ev] = make([]*roomViewObserver, 0)
	}

	s.observers[ev] = append(s.observers[ev], o)
}

func (s *roomViewSubscribers) unsubscribe(ev roomViewEventType, id string) {
	s.RLock()
	defer s.RUnlock()

	observers, ok := s.observers[ev]
	if !ok {
		s.debug("unsubscribe(): trying to unsubscribe from a not initialized event", ev)
		return
	}

	for i, o := range observers {
		if o.id == id {
			s.observers[ev] = append(
				s.observers[ev][:i], s.observers[ev][i+1:]...,
			)
			return
		}
	}
}

func (s *roomViewSubscribers) publish(ev roomViewEventType, ei roomViewEventInfo) {
	s.RLock()
	defer s.RUnlock()

	observers, ok := s.observers[ev]
	if !ok {
		s.debug("publish(): trying to publish a not initialized event", ev)
		return
	}

	for _, o := range observers {
		o.onNotify(ei)
	}
}

func (s *roomViewSubscribers) debug(m string, ev roomViewEventType) {
	s.log.Debug(m, ev)
}

func (v *roomView) subscribe(id string, ev roomViewEventType, onNotify func(roomViewEventInfo)) {
	v.subscribers.subscribe(ev, &roomViewObserver{
		id:       id,
		onNotify: onNotify,
	})
}

func (v *roomView) unsubscribe(id string, ev roomViewEventType) {
	v.subscribers.unsubscribe(ev, id)
}

func (v *roomView) publish(ev roomViewEventType) {
	v.subscribers.publish(ev, roomViewEventInfo{})
}

func (v *roomView) publishWithInfo(ev roomViewEventType, ei roomViewEventInfo) {
	v.subscribers.publish(ev, ei)
}
