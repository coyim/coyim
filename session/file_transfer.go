package session

import "github.com/twstrike/coyim/session/filetransfer"

func init() {
	registerKnownIQ("set", "http://jabber.org/protocol/bytestreams query", filetransfer.BytestreamQuery)

	registerKnownIQ("set", "http://jabber.org/protocol/ibb open", filetransfer.IbbOpen)
	registerKnownIQ("set", "http://jabber.org/protocol/ibb data", filetransfer.IbbData)
	registerKnownIQ("set", "http://jabber.org/protocol/ibb close", filetransfer.IbbClose)
	registerKnownExtension("http://jabber.org/protocol/ibb data", filetransfer.IbbMessageData)
}
