package gui

import (
	"github.com/coyim/coyim/session/events"
	log "github.com/sirupsen/logrus"
)

func (u *gtkUI) handleErrorMUCConflictEvent(a *account, ev events.MUCErrorEvent) {
	u.log.WithField("Event", ev).Debug("handleErrorMUCConflictEvent")
	a.log.WithFields(log.Fields{
		"from": ev.From,
	}).Info("Nickname conflict received")

}
