package filetransfer

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"sync"
	"time"

	"github.com/coyim/coyim/session/access"
	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/utils"
)

func init() {
	registerSendFileTransferMethod("http://jabber.org/protocol/bytestreams", bytestreamsSendDo, bytestreamsSendCurrentlyValid)
}

func bytestreamsSendCurrentlyValid(_ string, s access.Session, ctx *sendContext) bool {
	return len(bytestreamsGetCurrentValidProxies(s, ctx)) > 0
}

var defaultBytestreamProxyTimeout = 2 * time.Hour

func bytestreamsGetStreamhostDataFor(s access.Session, ctx *sendContext, jid string) *data.BytestreamStreamhost {
	rp, _, err := s.Conn().SendIQ(jid, "get", &data.BytestreamQuery{})
	if err != nil {
		return nil
	}
	r, ok := <-rp
	if !ok {
		return nil
	}
	var q data.BytestreamQuery
	if ctx.unpackIQData(s, r, &q) {
		for _, s := range q.Streamhosts {
			return &s
		}
	}
	return nil
}

func bytestreamsCalculateValidProxies(s access.Session, ctx *sendContext) func(key string) interface{} {
	return func(key string) interface{} {
		rp, _, err := s.Conn().SendIQ(utils.DomainFromJid(s.GetConfig().Account), "get", &data.DiscoveryItemsQuery{})
		if err != nil {
			// TODO, fix
			return nil
		}
		possibleProxies := ctx.bytestreamsWaitForDiscoveryItems(s, rp)
		results := make([]*data.BytestreamStreamhost, len(possibleProxies))
		wg := &sync.WaitGroup{}
		wg.Add(len(possibleProxies))
		for ix, pp := range possibleProxies {
			go func() {
				results[ix] = bytestreamsGetStreamhostDataFor(s, ctx, pp)
				wg.Done()
			}()
		}
		wg.Wait()
		return results
	}
}

func bytestreamsGetCurrentValidProxies(s access.Session, ctx *sendContext) []*data.BytestreamStreamhost {
	proxies, _ := s.Conn().Cache().GetOrComputeTimed("http://jabber.org/protocol/bytestreams . proxies", defaultBytestreamProxyTimeout, bytestreamsCalculateValidProxies(s, ctx))
	return proxies.([]*data.BytestreamStreamhost)

}
func bytestreamsSendDo(s access.Session, ctx *sendContext) {
	go func() {
		proxies := bytestreamsGetCurrentValidProxies(s, ctx)
		proxiesToSend := make([]data.BytestreamStreamhost, len(proxies))
		for ix, p := range proxies {
			proxiesToSend[ix] = *p
		}

		rp, _, err := s.Conn().SendIQ(ctx.peer, "set", &data.BytestreamQuery{
			Sid:         ctx.sid,
			Streamhosts: proxiesToSend,
		})
		if err != nil {
			// TODO something here
			return
		}
		r, ok := <-rp
		if !ok {
			// TODO something here
			return
		}

		var bq data.BytestreamQuery
		if ctx.unpackIQData(s, r, &bq) {
			fmt.Printf("Got streamhost to use: %#v\n", *bq.StreamhostUsed)
			// TODO: Continue HERE to send data - I get the JID here, and have to use that to lookup the date
			return
		}
		// TODO something here
		return
	}()
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

func (ctx *sendContext) bytestreamsWaitForDiscoveryItems(s access.Session, reply <-chan data.Stanza) []string {
	r, ok := <-reply
	if !ok {
		return []string{}
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
		return possibleProxies
	}
	return []string{}
}
