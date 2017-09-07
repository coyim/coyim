package session

import (
	"bytes"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/twstrike/coyim/digests"
	"github.com/twstrike/coyim/xmpp/data"
	"github.com/twstrike/coyim/xmpp/socks5"
	"golang.org/x/net/proxy"
)

func init() {
	registerKnownIQ("set", "http://jabber.org/protocol/bytestreams query", fileTransferBytestreamQuery)
}

type bytestreamContext struct {
	sync.Mutex
	sid        string
	conn       net.Conn
	streamhost string
}

func fileTransferBytestreamWaitForCancel(s *session, ift inflightFileTransfer) {
	// TODO: finish this
	if cancel, ok := <-ift.cancelChannel; ok && cancel {
		//		ift.fileTransferIbbCleanup()
		close(ift.finishedChannel)
		close(ift.updateChannel)
		close(ift.errorChannel)
		removeInflightFileTransfer(ift.id)
		//		s.conn.SendIQ(ift.peer, "set", data.IBBClose{Sid: ift.id})
	}
}

func fileTransferBytestreamQuery(s *session, stanza *data.ClientIQ) (ret interface{}, ignore bool) {
	var tag data.BytestreamQuery
	fmt.Printf("BLARG: %s\n", string(stanza.Query))
	if err := xml.NewDecoder(bytes.NewBuffer(stanza.Query)).Decode(&tag); err != nil {
		s.warn(fmt.Sprintf("Failed to parse bytestream open: %v", err))
		s.sendIQError(stanza, iqErrorNotAcceptable)
		return nil, true
	}
	fmt.Printf("Got query: %#v\n", tag)
	inflight, ok := getInflightFileTransfer(tag.Sid)

	if !ok || inflight.status.opaque != nil {
		s.warn(fmt.Sprintf("No file transfer associated with SID: %v", tag.Sid))
		s.sendIQError(stanza, iqErrorNotAcceptable)
		return nil, true
	}

	ctx := &bytestreamContext{
		sid: tag.Sid,
	}
	inflight.status.opaque = ctx

	dstAddr := hex.EncodeToString(digests.Sha1([]byte(tag.Sid + stanza.From + stanza.To)))

	// TODO: check if udp is asked
	// TODO: if we already have a destination address
	// TODO: make sure we use Tor or whatever proxies we have

	for _, sh := range tag.Streamhosts {
		fmt.Printf("Trying streamhost: %#v\n", sh)
		port := sh.Port
		if port == 0 {
			port = 1080
		}

		dialer, e := socks5.XMPP("tcp", net.JoinHostPort(sh.Host, strconv.Itoa(port)), nil, proxy.Direct)
		if e != nil {
			s.warn(fmt.Sprintf("Error setting up socks5 for %v: %v", sh, e))
		} else {
			conn, e2 := dialer.Dial("tcp", net.JoinHostPort(dstAddr, "0"))
			if e2 != nil {
				s.warn(fmt.Sprintf("Error connecting socks5 for %v: %v", sh, e2))
			} else {
				fmt.Printf("We have a connection to: %#v\n", sh)
				ctx.conn = conn
				ctx.streamhost = sh.Jid

				reply := data.BytestreamQuery{
					Sid:            tag.Sid,
					StreamhostUsed: &data.BytestreamStreamhostUsed{Jid: sh.Jid},
				}

				go func() {
					dest := inflight.status.destination
					// TODO: we probably don't want to use io.Copy in the real case, since we want to be able to give feedback/updates
					ff, _ := ioutil.TempFile(filepath.Dir(dest), filepath.Base(dest))
					defer ff.Close()
					fmt.Printf("Starting to read from the connection to: %#v\n", sh)

					chunkSize := 4096
					buf := make([]byte, chunkSize)
					totalWritten := int64(0)
					for {
						n, err := conn.Read(buf)
						if err != nil {
							if err != io.EOF {
								fmt.Println("read error:", err)
							}
							break
						}
						_, err = ff.Write(buf[:n])
						if err != nil {
							fmt.Println("write error:", err)
						}
						totalWritten += int64(n)
						inflight.updateChannel <- totalWritten

					}

					os.Rename(ff.Name(), inflight.status.destination)
					inflight.reportFinished()
					removeInflightFileTransfer(tag.Sid)
				}()

				fmt.Printf("Returning stuff: %#v\n", reply)
				return reply, false
			}
		}
	}

	// here we have failed to find a working stream host. =(

	return nil, false
}
