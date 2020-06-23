package session

import (
	"fmt"

	"github.com/coyim/coyim/session/data"
	"github.com/coyim/coyim/session/filetransfer"
	"github.com/coyim/coyim/xmpp/jid"
)

func init() {
	registerKnownIQ("set", "http://jabber.org/protocol/bytestreams query", filetransfer.BytestreamQuery)

	registerKnownIQ("set", "http://jabber.org/protocol/ibb open", filetransfer.IbbOpen)
	registerKnownIQ("set", "http://jabber.org/protocol/ibb data", filetransfer.IbbData)
	registerKnownIQ("set", "http://jabber.org/protocol/ibb close", filetransfer.IbbClose)
	registerKnownExtension("http://jabber.org/protocol/ibb data", filetransfer.IbbMessageData)
}

func (s *session) SendFileTo(peer jid.Any, filename string, onNoEnc func() bool, encDecision func(bool)) *data.FileTransferControl {
	s.info(fmt.Sprintf("SendFileTo(%s, %s)", peer, filename))
	return filetransfer.InitSend(s, peer, filename, onNoEnc, encDecision)
}

func (s *session) SendDirTo(peer jid.Any, dirname string, onNoEnc func() bool, encDecision func(bool)) *data.FileTransferControl {
	s.info(fmt.Sprintf("SendDirTo(%s, %s)", peer, dirname))
	return filetransfer.InitSendDir(s, peer, dirname, onNoEnc, encDecision)
}
