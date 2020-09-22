package muc

import "github.com/coyim/coyim/xmpp/jid"

// MUC is a marker interface that is used to differentiate MUC "things"
type MUC interface {
	WhichRoom() jid.Bare
}
