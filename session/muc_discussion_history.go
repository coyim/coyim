package session

import (
	"sync"

	"github.com/coyim/coyim/xmpp/jid"

	"github.com/coyim/coyim/session/muc/data"
	xmppData "github.com/coyim/coyim/xmpp/data"
)

type discussionHistoryManager struct {
	history              map[string]*data.DiscussionHistory
	handleClientPresence func(*xmppData.ClientMessage)
	lock                 sync.Mutex
}

func newDiscussionHistoryManager(handleHistory func(*xmppData.ClientMessage)) *discussionHistoryManager {
	return &discussionHistoryManager{
		history:              make(map[string]*data.DiscussionHistory),
		handleClientPresence: doOnceWithStanza(handleHistory),
	}
}

// getHistory returns the discussion history for the given room and
// a boolean indicating if it was found or not
func (dm *discussionHistoryManager) getHistory(roomID jid.Bare) (*data.DiscussionHistory, bool) {
	h, ok := dm.history[roomID.String()]
	return h, ok
}

func (dm *discussionHistoryManager) addHistory(roomID jid.Bare) *data.DiscussionHistory {
	dm.lock.Lock()
	defer dm.lock.Unlock()

	h := data.NewDiscussionHistory()
	dm.history[roomID.String()] = h

	return h
}
