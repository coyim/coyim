package muc

import (
	"github.com/coyim/coyim/session/muc/data"
	"github.com/coyim/coyim/xmpp/jid"
)

type RoomBanListEntry struct {
	Jid         jid.Any
	Affiliation data.Affiliation
	Reason      string
}
