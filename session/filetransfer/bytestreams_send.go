package filetransfer

import (
	"bytes"
	"encoding/xml"
	"fmt"

	"github.com/coyim/coyim/session/access"
	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/utils"
)

func init() {
	registerSendFileTransferMethod("http://jabber.org/protocol/bytestreams", bytestreamsSendDo)
}

func bytestreamsSendDo(s access.Session, ctx *sendContext) {
	fmt.Printf("HOLA, WHATS UP?\n")

	rp, _, err := s.Conn().SendIQ(utils.DomainFromJid(s.GetConfig().Account), "get", &data.DiscoveryItemsQuery{})
	if err != nil {
		// TODO, fix
	}
	go ctx.bytestreamsWaitForDiscoveryItems(s, rp)

	//   for each of these, we need to do a <query xmlns='http://jabber.org/protocol/bytestreams'/> this should return streamhost information

}

func (ctx *sendContext) unpackIQData(s access.Session, d data.Stanza, res interface{}) bool {
	switch ciq := d.Value.(type) {
	case *data.ClientIQ:
		if ciq.Type == "result" {
			if err := xml.NewDecoder(bytes.NewBuffer(ciq.Query)).Decode(res); err != nil {
				// TODO: blah
				return false
			}
			return true
		}
	}
	return false
}

func (ctx *sendContext) bytestreamsWaitForDiscoveryItems(s access.Session, reply <-chan data.Stanza) {
	r, ok := <-reply
	if !ok {
		// TODO: report here
		return
	}
	var ditems data.DiscoveryItemsQuery
	if ctx.unpackIQData(s, r, &ditems) {
		possibleProxies := []string{}
		for _, di := range ditems.DiscoveryItems {
			ids, feats, _ := s.Conn().DiscoveryFeaturesAndIdentities(di.Jid)
			hasCorrectIdentity := false
			hasBytestreamsFeature := false
			for _, id := range ids {
				if id.Category == "proxy" && id.Type == "bytestreams" {
					hasCorrectIdentity = true
				}
			}
			for _, feat := range feats {
				if feat == "http://jabber.org/protocol/bytestreams" {
					hasBytestreamsFeature = true
				}
			}
			if hasCorrectIdentity && hasBytestreamsFeature {
				possibleProxies = append(possibleProxies, di.Jid)
			}

		}
		fmt.Printf("Possible proxies: %#v\n", possibleProxies)
	}
}
