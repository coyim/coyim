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

const bufSize = 16 * 4096

func init() {
	registerSendFileTransferMethod("http://jabber.org/protocol/bytestreams", bytestreamsSendDo, bytestreamsSendCurrentlyValid)
}

func bytestreamsSendCurrentlyValid(_ string, s access.Session) bool {
	return len(bytestreamsGetCurrentValidProxies(s)) > 0
}

var defaultBytestreamProxyTimeout = 2 * time.Hour

func bytestreamsGetStreamhostDataFor(s access.Session, jid string) (result *data.BytestreamStreamhost) {
	var q data.BytestreamQuery
	basicIQ(s, jid, "get", &data.BytestreamQuery{}, &q, func(*data.ClientIQ) {
		for _, sh := range q.Streamhosts {
			result = &sh
			return
		}
	})
	return
}

func bytestreamsCalculateValidProxies(s access.Session) func(key string) interface{} {
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
			return nil
		}

		results := make([]*data.BytestreamStreamhost, len(possibleProxies))
		wg := &sync.WaitGroup{}
		wg.Add(len(possibleProxies))
		for ix, pp := range possibleProxies {
			go func(index int, proxy string) {
				results[index] = bytestreamsGetStreamhostDataFor(ctx, proxy)
				wg.Done()
			}(ix, pp)
		}
		wg.Wait()
		return results
	}
}

func bytestreamsGetCurrentValidProxies(s access.Session) []*data.BytestreamStreamhost {
	proxies, _ := s.Conn().Cache().GetOrComputeTimed("http://jabber.org/protocol/bytestreams . proxies", defaultBytestreamProxyTimeout, bytestreamsCalculateValidProxies(s))
	return proxies.([]*data.BytestreamStreamhost)

}

func bytestreamsSendData(ctx *sendContext, c net.Conn) {
	fmt.Printf("bytestreamsSendData()\n")
	defer func() {
		fmt.Printf("  bytestreamsSendData()-before c.Close()\n")
		c.Close()
		fmt.Printf("  bytestreamsSendData()-c.Close()\n")
	}()

	buffer := make([]byte, bufSize)
	r, err := os.Open(ctx.file)
	if err != nil {
		ctx.onError(err)
		return
	}
	defer func() {
		fmt.Printf("  bytestreamsSendData()-before r.Close()\n")
		r.Close()
		fmt.Printf("  bytestreamsSendData()-r.Close()\n")
	}()
	for {
		if ctx.weWantToCancel {
			fmt.Printf("  bytestreamsSendData()-we want to cancel\n")
			removeInflightSend(ctx)
			fmt.Printf("  bytestreamsSendData()-after removing inflight send\n")
			return
		}

		fmt.Printf("  bytestreamsSendData()-r.read\n")
		n, rerr := r.Read(buffer)
		fmt.Printf("  bytestreamsSendData()-  have read %d\n", n)
		if n > 0 {
			fmt.Printf("  bytestreamsSendData()-  writing\n")
			n2, err := c.Write(buffer[:n])
			fmt.Printf("  bytestreamsSendData()-  have written: %d\n", n2)
			if err != nil {
				fmt.Printf("  bytestreamsSendData()-  had an error writing\n")
				ctx.onError(err)
				return
			}
			fmt.Printf("  bytestreamsSendData()-  sending on update\n")
			ctx.onUpdate(n)
		}

		if rerr == io.EOF {
			fmt.Printf("  bytestreamsSendData()-  finishing up\n")
			ctx.onFinish()
			fmt.Printf("  bytestreamsSendData()-  done finishing  up\n")
			return
		} else if rerr != nil {
			fmt.Printf("  bytestreamsSendData()-  had sending error\n")
			ctx.onError(rerr)
			return
		}
	}
}

func bytestreamsSendDo(ctx *sendContext) {
	go func() {
		proxies := bytestreamsGetCurrentValidProxies(ctx.s)
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
					go bytestreamsSendData(ctx, c)
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
