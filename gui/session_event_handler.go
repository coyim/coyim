package gui

import (
	"fmt"

	"github.com/twstrike/coyim/xmpp"
	"github.com/twstrike/otr3"
)

//I will ignore this class for now

type guiSessionEventHandler struct {
	u *gtkUI
}

func (seh guiSessionEventHandler) Info(m string) {
	seh.u.Info(m)
}

func (seh guiSessionEventHandler) Warn(m string) {
	seh.u.Warn(m)
}

func (seh guiSessionEventHandler) Alert(m string) {
	seh.u.Alert(m)
}

func (seh guiSessionEventHandler) NewOTRKeys(uid string, conversation *otr3.Conversation) {
	seh.u.Info(fmt.Sprintf("TODO: notify new keys from %s", uid))
}

func (seh guiSessionEventHandler) OTREnded(uid string) {
	//TODO: conversation ended
}

//TODO: we should update periodically (like Pidgin does) if we include the status (online/offline/away) on the label
func (seh guiSessionEventHandler) RosterReceived(roster []xmpp.RosterEntry) {
	//TODO: outdated
	//glib.IdleAdd(func() bool {
	//	seh.u.roster.Update(roster)
	//	return false
	//})
}

func (seh guiSessionEventHandler) IQReceived(string) {
	//TODO
}

func (seh guiSessionEventHandler) MessageReceived(from, timestamp string, encrypted bool, message []byte) {
	seh.u.roster.MessageReceived(from, timestamp, encrypted, message)
}

func (seh guiSessionEventHandler) ProcessPresence(stanza *xmpp.ClientPresence, ignore, gone bool) {
	//jid := xmpp.RemoveResourceFromJid(stanza.From)
	//state, ok := seh.u.session.KnownStates[jid]
	//if !ok || len(state) == 0 {
	//	state = "unknown"
	//}

	////TODO: Notify via UI
	//fmt.Println(jid, "is", state)
}

func (seh guiSessionEventHandler) Disconnected() {
	//TODO: remove everybody from the roster
	fmt.Println("TODO: Should disconnect the account")
}

func (seh guiSessionEventHandler) RegisterCallback() xmpp.FormCallback {
	//if !*createAccount {
	//  return nil
	//}

	return func(title, instructions string, fields []interface{}) error {
		//TODO: should open a registration window
		fmt.Println("TODO")
		return nil
	}
}
