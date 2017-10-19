package filetransfer

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"

	"github.com/coyim/coyim/digests"
	"github.com/coyim/coyim/session/access"
	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/utils"
)

const bufSize = 64 * 4096

func init() {
	registerSendFileTransferMethod("http://jabber.org/protocol/bytestreams", bytestreamsSendDo, bytestreamsSendCurrentlyValid)
}

func bytestreamsSendCurrentlyValid(_ string, s access.Session, ctx *sendContext) bool {
	return len(bytestreamsGetCurrentValidProxies(s, ctx)) > 0
}

var defaultBytestreamProxyTimeout = 2 * time.Hour

func bytestreamsGetStreamhostDataFor(s access.Session, ctx *sendContext, jid string) (result *data.BytestreamStreamhost) {
	var q data.BytestreamQuery
	basicIQ(s, jid, "get", &data.BytestreamQuery{}, &q, func(*data.ClientIQ) {
		for _, s := range q.Streamhosts {
			result = &s
			return
		}
	})
	return
}

func bytestreamsCalculateValidProxies(s access.Session, ctx *sendContext) func(key string) interface{} {
	return func(key string) interface{} {
		var ditems data.DiscoveryItemsQuery
		possibleProxies := []string{}
		e := basicIQ(s, utils.DomainFromJid(s.GetConfig().Account), "get", &data.DiscoveryItemsQuery{}, &ditems, func(*data.ClientIQ) {
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
		})

		if e != nil {
			ctx.control.ReportError(e)
			removeInflightSend(ctx)
			return nil
		}

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

func (ctx *sendContext) bytestreamsSendData(s access.Session, c net.Conn) {
	defer c.Close()

	buffer := make([]byte, bufSize)
	r, err := os.Open(ctx.file)
	if err != nil {
		ctx.control.ReportError(err)
		removeInflightSend(ctx)
		return
	}
	defer r.Close()
	for {
		if ctx.weWantToCancel {
			if ctx.weWantToCancel {
				removeInflightSend(ctx)
				return
			}
		}
		n, err := r.Read(buffer)
		if err == io.EOF && n == 0 {
			ctx.control.ReportFinished()
			removeInflightSend(ctx)
			return
		} else if err != nil {
			ctx.control.ReportError(err)
			removeInflightSend(ctx)
			return
		}
		_, err = c.Write(buffer[0:n])
		if err != nil {
			ctx.control.ReportError(err)
			removeInflightSend(ctx)
			return
		}
		ctx.totalSent += int64(n)
		ctx.control.SendUpdate(ctx.totalSent)
	}
}

func bytestreamsSendDo(s access.Session, ctx *sendContext) {
	go func() {
		proxies := bytestreamsGetCurrentValidProxies(s, ctx)
		proxiesToSend := make([]data.BytestreamStreamhost, len(proxies))
		proxyMap := make(map[string]data.BytestreamStreamhost)
		for ix, p := range proxies {
			proxiesToSend[ix] = *p
			proxyMap[p.Jid] = *p
		}

		var bq data.BytestreamQuery
		if err := basicIQ(s, ctx.peer, "set", &data.BytestreamQuery{
			Sid:         ctx.sid,
			Streamhosts: proxiesToSend,
		}, &bq, func(ciq *data.ClientIQ) {
			sh, ok := proxyMap[bq.StreamhostUsed.Jid]
			if !ok {
				ctx.control.ReportError(errors.New("Invalid streamhost to use - this is likely a developer error from the peers side"))
				removeInflightSend(ctx)
				return
			}
			dstAddr := hex.EncodeToString(digests.Sha1([]byte(ctx.sid + ciq.To + ciq.From)))
			if !tryStreamhost(s, sh, dstAddr, func(c net.Conn) {
				e := basicIQ(s, bq.StreamhostUsed.Jid, "set", &data.BytestreamQuery{
					Sid:      ctx.sid,
					Activate: ciq.From,
				}, nil, func(*data.ClientIQ) {
					go ctx.bytestreamsSendData(s, c)
				})
				if e != nil {
					ctx.control.ReportError(e)
					removeInflightSend(ctx)
				}
			}) {
				ctx.control.ReportError(fmt.Errorf("Failed at connecting to streamhost: %#v", sh))
				removeInflightSend(ctx)
			}
		}); err != nil {
			ctx.control.ReportError(err)
			removeInflightSend(ctx)
		}
	}()
}
