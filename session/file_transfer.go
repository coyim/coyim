package session

import (
	"fmt"

	"github.com/twstrike/coyim/session/filetransfer"
)

func init() {
	registerKnownIQ("set", "http://jabber.org/protocol/bytestreams query", filetransfer.BytestreamQuery)

	registerKnownIQ("set", "http://jabber.org/protocol/ibb open", filetransfer.IbbOpen)
	registerKnownIQ("set", "http://jabber.org/protocol/ibb data", filetransfer.IbbData)
	registerKnownIQ("set", "http://jabber.org/protocol/ibb close", filetransfer.IbbClose)
	registerKnownExtension("http://jabber.org/protocol/ibb data", filetransfer.IbbMessageData)
}

func (s *session) SendFileTo(peer, filename string) {
	s.info(fmt.Sprintf("SendFileTo(%s, %s)", peer, filename))
	if err := filetransfer.InitSend(s, peer, filename); err != nil {
		s.warn(fmt.Sprintf("Failed to send file to %s: %v", peer, err))
	}
}
