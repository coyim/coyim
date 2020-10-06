package gui

import (
	"sync"

	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/session/events"
	"github.com/coyim/coyim/xmpp/jid"
	log "github.com/sirupsen/logrus"
)

func (v *roomView) handleRoomEvent(ev events.MUC) {
	switch t := ev.(type) {
	case events.MUCSelfOccupantJoined:
		v.publishEvent(occupantSelfJoinedEvent{t.Nickname})
	case events.MUCOccupantUpdated:
		v.publishEvent(occupantUpdatedEvent{t.Nickname})
	case events.MUCOccupantJoined:
		v.publishEvent(occupantJoinedEvent{t.Nickname})
	case events.MUCOccupantLeft:
		v.publishEvent(occupantLeftEvent{t.Nickname})
	case events.MUCMessageReceived:
		v.publishMessageEvent("received", t.Nickname, t.Subject, t.Message)
	case events.MUCLoggingEnabled:
		v.publishEvent(loggingEnabledEvent{})
	case events.MUCLoggingDisabled:
		v.publishEvent(loggingDisabledEvent{})
	default:
		v.log.WithField("event", t).Warn("Unsupported room event received")
	}
}

type roomViewObserver struct {
	id       string
	onNotify func(roomViewEvent)
}

type roomViewSubscribers struct {
	observers     []*roomViewObserver
	observersLock sync.Mutex

	log coylog.Logger
}

func newRoomViewSubscribers(roomID jid.Bare, l coylog.Logger) *roomViewSubscribers {
	s := &roomViewSubscribers{}

	s.log = l.WithFields(log.Fields{
		"who":  "newRoomViewSubscribers",
		"room": roomID,
	})

	return s
}

func (s *roomViewSubscribers) subscribe(id string, onNotify func(roomViewEvent)) {
	s.observersLock.Lock()
	defer s.observersLock.Unlock()

	s.observers = append(s.observers, &roomViewObserver{
		id:       id,
		onNotify: onNotify,
	})
}

func (s *roomViewSubscribers) unsubscribe(id string) {
	s.observersLock.Lock()
	hasChanged := false
	observers := s.observers

	defer func() {
		if hasChanged {
			s.observers = observers
		}
		s.observersLock.Unlock()
	}()

	for i, o := range observers {
		if o.id == id {
			observers = append(
				observers[:i], observers[i+1:]...,
			)
			hasChanged = true
			return
		}
	}
}

func (s *roomViewSubscribers) publish(ev roomViewEvent) {
	s.observersLock.Lock()
	observers := s.observers
	s.observersLock.Unlock()

	for _, o := range observers {
		o.onNotify(ev)
	}
}

func (s *roomViewSubscribers) debug(m string, ev string) {
	s.log.Debug(m, ev)
}

func (v *roomView) subscribe(id string, onNotify func(roomViewEvent)) {
	v.subscribers.subscribe(id, onNotify)
}

func (v *roomView) unsubscribe(id string) {
	v.subscribers.unsubscribe(id)
}

func (v *roomView) publishEvent(ev roomViewEvent) {
	if v.opened {
		// We should analyze if at some point we need
		// to be able to publish an event in a closed view
		v.subscribers.publish(ev)
	}
}

func (v *roomView) publishMessageEvent(tp, nickname, subject, msg string) {
	v.publishEvent(messageEvent{
		tp:       tp,
		nickname: nickname,
		subject:  subject,
		message:  msg,
	})
}
