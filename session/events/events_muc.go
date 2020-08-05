package events

import (
	"github.com/coyim/coyim/xmpp/data"
)

// MUCPresence represents a muc session presence event
type MUCPresence struct {
	*data.ClientPresence
}
