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
	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/utils"
)

const bufSize = 64 * 4096

func init() {
	registerSendFileTransferMethod("http://jabber.org/protocol/bytestreams", bytestreamsSendDo, bytestreamsSendCurrentlyValid)
}

func bytestreamsSendCurrentlyValid(_ string, ctx *sendContext) bool {
	return len(bytestreamsGetCurrentValidProxies(ctx)) > 0
}

var defaultBytestreamProxyTimeout = 2 * time.Hour

func bytestreamsGetStreamhostDataFor(ctx *sendContext, jid string) (result *data.BytestreamStreamhost) {
	var q data.BytestreamQuery
	basicIQ(ctx.s, jid, "get", &data.BytestreamQuery{}, &q, func(*data.ClientIQ) {
		for _, s := range q.Streamhosts {
			result = &s
			return
		}
	})
	return
}

func bytestreamsCalculateValidProxies(ctx *sendContext) func(key string) interface{} {
	return func(key string) interface{} {
		var ditems data.DiscoveryItemsQuery
		possibleProxies := []string{}
		e := basicIQ(ctx.s, utils.DomainFromJid(ctx.s.GetConfig().Account), "get", &data.DiscoveryItemsQuery{}, &ditems, func(*data.ClientIQ) {
			for _, di := range ditems.DiscoveryItems {
				ids, feats, _ := ctx.s.Conn().DiscoveryFeaturesAndIdentities(di.Jid)
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
			ctx.onError(e)
			return nil
		}

		results := make([]*data.BytestreamStreamhost, len(possibleProxies))
		wg := &sync.WaitGroup{}
		wg.Add(len(possibleProxies))
		for ix, pp := range possibleProxies {
			go func() {
				results[ix] = bytestreamsGetStreamhostDataFor(ctx, pp)
				wg.Done()
			}()
		}
		wg.Wait()
		return results
	}
}

func bytestreamsGetCurrentValidProxies(ctx *sendContext) []*data.BytestreamStreamhost {
	proxies, _ := ctx.s.Conn().Cache().GetOrComputeTimed("http://jabber.org/protocol/bytestreams . proxies", defaultBytestreamProxyTimeout, bytestreamsCalculateValidProxies(ctx))
	return proxies.([]*data.BytestreamStreamhost)

}

func (ctx *sendContext) bytestreamsSendData(c net.Conn) {
	defer c.Close()

	buffer := make([]byte, bufSize)
	r, err := os.Open(ctx.file)
	if err != nil {
		ctx.onError(err)
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
			ctx.onFinish()
			return
		} else if err != nil {
			ctx.onError(err)
			return
		}
		_, err = c.Write(buffer[0:n])
		if err != nil {
			ctx.onError(err)
			return
		}
		ctx.onUpdate(n)
	}
}

func bytestreamsSendDo(ctx *sendContext) {
	go func() {
		proxies := bytestreamsGetCurrentValidProxies(ctx)
		proxiesToSend := make([]data.BytestreamStreamhost, len(proxies))
		proxyMap := make(map[string]data.BytestreamStreamhost)
		for ix, p := range proxies {
			proxiesToSend[ix] = *p
			proxyMap[p.Jid] = *p
		}

		var bq data.BytestreamQuery
		if err := basicIQ(ctx.s, ctx.peer, "set", &data.BytestreamQuery{
			Sid:         ctx.sid,
			Streamhosts: proxiesToSend,
		}, &bq, func(ciq *data.ClientIQ) {
			sh, ok := proxyMap[bq.StreamhostUsed.Jid]
			if !ok {
				ctx.onError(errors.New("Invalid streamhost to use - this is likely a developer error from the peers side"))
				return
			}
			dstAddr := hex.EncodeToString(digests.Sha1([]byte(ctx.sid + ciq.To + ciq.From)))
			if !tryStreamhost(ctx.s, sh, dstAddr, func(c net.Conn) {
				e := basicIQ(ctx.s, bq.StreamhostUsed.Jid, "set", &data.BytestreamQuery{
					Sid:      ctx.sid,
					Activate: ciq.From,
				}, nil, func(*data.ClientIQ) {
					go ctx.bytestreamsSendData(c)
				})
				if e != nil {
					ctx.onError(e)
				}
			}) {
				ctx.onError(fmt.Errorf("Failed at connecting to streamhost: %#v", sh))
			}
		}); err != nil {
			ctx.onError(err)
		}
	}()
}
