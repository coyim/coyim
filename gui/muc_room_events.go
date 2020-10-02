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
		v.publishWithInfo("occupantSelfJoinedEvent", roomViewEventInfo{
			"nickname": t.Nickname,
		})

	case events.MUCOccupantUpdated:
		v.publishWithInfo("occupantUpdatedEvent", roomViewEventInfo{
			"nickname": t.Nickname,
		})

	case events.MUCOccupantJoined:
		v.publishWithInfo("occupantJoinedEvent", roomViewEventInfo{
			"nickname": t.Nickname,
		})

	case events.MUCOccupantLeft:
		v.publishWithInfo("occupantLeftEvent", roomViewEventInfo{
			"nickname": t.Nickname,
		})

	case events.MUCMessageReceived:
		v.publishWithInfo("messageReceivedEvent", roomViewEventInfo{
			"nickname": t.Nickname,
			"subject":  t.Subject,
			"message":  t.Message,
		})

	case events.MUCLoggingEnabled:
		v.publish("loggingEnabledEvent")

	case events.MUCLoggingDisabled:
		v.publish("loggingDisabledEvent")

	default:
		v.log.WithField("event", t).Warn("Unsupported room event received")
	}
}

// roomViewEventInfo contains information about any room view event
type roomViewEventInfo map[string]string

type roomViewEventObservers map[string]func(roomViewEventInfo)

type roomViewObserver struct {
	id       string
	onNotify func(roomViewEventInfo)
}

type roomViewSubscribers struct {
	observers     map[string][]*roomViewObserver
	observersLock sync.Mutex

	log coylog.Logger
}

func newRoomViewSubscribers(room jid.Bare, l coylog.Logger) *roomViewSubscribers {
	return &roomViewSubscribers{
		log: l.WithFields(log.Fields{
			"who":  "newRoomViewSubscribers",
			"room": room,
		}),
		observers: make(map[string][]*roomViewObserver),
	}
}

func (s *roomViewSubscribers) subscribe(ev string, o *roomViewObserver) {
	s.observersLock.Lock()
	defer s.observersLock.Unlock()

	s.observers[ev] = append(s.observers[ev], o)
}

func (s *roomViewSubscribers) unsubscribe(ev string, id string) {
	s.observersLock.Lock()
	hasChanged := false
	observers, ok := s.observers[ev]

	defer func() {
		if hasChanged {
			s.observers[ev] = observers
		}
		s.observersLock.Unlock()
	}()

	if !ok {
		s.debug("unsubscribe(): trying to unsubscribe from a not initialized event", ev)
		return
	}

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

func (s *roomViewSubscribers) publish(ev string, ei roomViewEventInfo) {
	s.observersLock.Lock()
	observers, ok := s.observers[ev]
	s.observersLock.Unlock()

	if !ok {
		s.debug("publish(): trying to publish a not initialized event", ev)
		return
	}

	for _, o := range observers {
		o.onNotify(ei)
	}
}

func (s *roomViewSubscribers) debug(m string, ev string) {
	s.log.Debug(m, ev)
}

func (v *roomView) subscribe(id string, ev string, onNotify func(roomViewEventInfo)) {
	v.subscribers.subscribe(ev, &roomViewObserver{
		id:       id,
		onNotify: onNotify,
	})
}

func (v *roomView) subscribeAll(id string, o roomViewEventObservers) {
	for ev, onNotify := range o {
		v.subscribe(id, ev, onNotify)
	}
}

func (v *roomView) unsubscribe(id string, ev string) {
	v.subscribers.unsubscribe(ev, id)
}

func (v *roomView) publish(ev string) {
	v.subscribers.publish(ev, roomViewEventInfo{})
}

func (v *roomView) publishWithInfo(ev string, ei roomViewEventInfo) {
	v.subscribers.publish(ev, ei)
}
