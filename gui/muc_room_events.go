package gui

import (
	"sync"

	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/session/events"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
	log "github.com/sirupsen/logrus"
)

func (v *roomView) handleRoomEvent(ev muc.MUC) {
	// TODO: Most of these handlers simply publish a new event
	// There is probably no point in doing this publishing in the UI thread
	switch t := ev.(type) {
	case events.MUCSelfOccupantJoined:
		v.publishWithInfo("occupantSelfJoined", roomViewEventInfo{
			"nickname": t.Nickname,
		})

	case events.MUCOccupantUpdated:
		v.publishWithInfo("occupantUpdated", roomViewEventInfo{
			"nickname": t.Nickname,
		})

	case events.MUCOccupantJoined:
		v.publishWithInfo("occupantJoined", roomViewEventInfo{
			"nickname": t.Nickname,
		})

	case events.MUCOccupantLeft:
		v.publishWithInfo("occupantLeft", roomViewEventInfo{
			"nickname": t.Nickname,
		})

	case events.MUCMessageReceived:
		v.publishWithInfo("messageReceived", roomViewEventInfo{
			"nickname": t.Nickname,
			"subject":  t.Subject,
			"message":  t.Message,
		})

	case events.MUCLoggingEnabled:
		v.publish("loggingEnabled")

	case events.MUCLoggingDisabled:
		v.publish("loggingDisabled")

	default:
		v.log.WithField("event", t).Warn("Unsupported room event received")
	}
}

// TODO: This pattern might be problematic once we start having more and more and more
// data. It might be better to have a map with information, so that each event can use
// whatever fields they need

// roomViewEventInfo contains information about any room view event
type roomViewEventInfo map[string]string

type roomViewEventObservers map[string]func(roomViewEventInfo)

type roomViewObserver struct {
	id       string
	onNotify func(roomViewEventInfo)
}

type roomViewSubscribers struct {
	observers map[string][]*roomViewObserver
	log       coylog.Logger
	sync.RWMutex
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

// TODO: I feel like the "roomViewObserver" type should not be exposed to users
// of the functionality. Better that it takes a name and the function directly
func (s *roomViewSubscribers) subscribe(ev string, o *roomViewObserver) {
	s.Lock()
	defer s.Unlock()

	s.observers[ev] = append(s.observers[ev], o)
}

func (s *roomViewSubscribers) subscribeAll(observers map[string]*roomViewObserver) {
	for ev, ob := range observers {
		s.subscribe(ev, ob)
	}
}

func (s *roomViewSubscribers) unsubscribe(ev string, id string) {
	// TODO: Unsafe
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

func (s *roomViewSubscribers) publish(ev string, ei roomViewEventInfo) {
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
