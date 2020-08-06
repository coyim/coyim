package gui

import (
	"errors"
	"sync"

	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
)

type roomViewsManager struct {
	manager *muc.RoomManager
	views   map[string]*roomView
	events  chan interface{}
	log     coylog.Logger
	sync.RWMutex
}

func (m *roomViewsManager) hasRoom(ident jid.Bare) bool {
	_, v := m.manager.GetRoom(ident)
	return v
}

func (m *roomViewsManager) addRoom(ident jid.Bare, r *roomView) error {
	if !m.manager.AddRoom(r.room) {
		return errors.New("the room is already in the manager")
	}

	_, ok := m.views[r.id()]
	if ok {
		return errors.New("the room is already in the manager")
	}

	m.views[r.id()] = r

	return nil
}

func newRoomManager() *roomViewsManager {
	return &roomViewsManager{
		manager: muc.NewRoomManager(),
		events:  make(chan interface{}, 10),
		views:   make(map[string]*roomView),
	}
}

func (a *account) joinRoom(u *gtkUI, rjid jid.Bare) (*roomView, error) {
	return u.addRoom(a, rjid)
}

func (u *gtkUI) addRoom(a *account, ident jid.Bare) (*roomView, error) {
	a.roomManager.Lock()
	defer a.roomManager.Unlock()

	if a.roomManager.hasRoom(ident) {
		return nil, errors.New("the room is already opened")
	}

	r := u.newRoomView(a, ident)
	r.log = u.log

	err := a.roomManager.addRoom(ident, r)
	if err != nil {
		return nil, err
	}

	go a.roomManager.observeRoomEvents(r)

	a.session.Subscribe(r.events)

	return r, nil
}
