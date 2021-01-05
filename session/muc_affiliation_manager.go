package session

import (
	"errors"

	"github.com/coyim/coyim/session/muc/data"
	xmppData "github.com/coyim/coyim/xmpp/data"
	xi "github.com/coyim/coyim/xmpp/interfaces"
	"github.com/coyim/coyim/xmpp/jid"
)

type affiliationContext struct {
	occupant    jid.Bare
	roomID      jid.Bare
	reason      string
	affiliation string

	conn xi.Conn
}

func (s *session) newAffiliationContext(occupant, roomID jid.Bare, reason, affiliation string) *affiliationContext {
	return &affiliationContext{
		occupant:    occupant,
		roomID:      roomID,
		reason:      reason,
		affiliation: affiliation,
		conn:        s.conn,
	}
}

func (s *session) AssignAdminPrivilige(occupant, roomID jid.Bare, reason string) (chan bool, chan error) {
	rc := make(chan bool)
	ec := make(chan error)
	ctx := s.newAffiliationContext(occupant, roomID, reason, data.AffiliationAdmin)

	go func() {
		err := ctx.assingAffiliation()
		if err != nil {
			ec <- err
			return
		}
		rc <- true
	}()
	return rc, ec
}

func (ctx *affiliationContext) assingAffiliation() error {
	reply, _, err := ctx.conn.SendIQ(ctx.roomID.String(), "set", &xmppData.MUCAdmin{
		Item: &xmppData.Item{
			Affiliation: ctx.affiliation,
			Jid:         ctx.occupant.String(),
			Reason:      ctx.reason,
		},
	})

	if err != nil {
		return err
	}
	return ctx.validateStanza(<-reply)
}

func (ctx *affiliationContext) validateStanza(stanza xmppData.Stanza) error {
	iq, ok := stanza.Value.(*xmppData.ClientIQ)
	if !ok {
		return errors.New("error trying to parse the information query")
	}

	if iq.Type != "result" {
		return errors.New("response is not a result type")
	}

	return nil
}
