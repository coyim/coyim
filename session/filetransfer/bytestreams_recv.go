package filetransfer

import (
	"bytes"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/coyim/coyim/digests"
	"github.com/coyim/coyim/session/access"
	"github.com/coyim/coyim/xmpp/data"
)

// TODO: at some point this should be refactored away into a pure socks5 bytestream implementation and a small piece that is file transfer specific

func bytestreamWaitForCancel(ctx *recvContext) {
	ctx.control.WaitForCancel(func() {
		ctx.opaque.(chan bool) <- true
		removeInflightRecv(ctx.sid)
	})
}

func bytestreamInitialSetup(s access.Session, stanza *data.ClientIQ) (tag data.BytestreamQuery, ctx *recvContext, earlyReturn bool) {
	if err := xml.NewDecoder(bytes.NewBuffer(stanza.Query)).Decode(&tag); err != nil || tag.Sid == "" {
		s.Warn(fmt.Sprintf("Failed to parse bytestream open: %v", err))
		s.SendIQError(stanza, iqErrorIBBBadRequest)
		return tag, ctx, true
	}

	ctx, ok := getInflightRecv(tag.Sid)

	if !ok || ctx.opaque != nil {
		s.Warn(fmt.Sprintf("No file transfer associated with SID: %v", tag.Sid))
		s.SendIQError(stanza, iqErrorNotAcceptable)
		return tag, ctx, true
	}

	if tag.Mode == "udp" {
		// This shouldn't really be possible, since we don't advertise udp support
		// But we can always register the error anyway.
		s.Warn("Received a request for UDP, even though we don't support or advertize UDP - this means the peer is using a non-conforming application")
		s.SendIQError(stanza, iqErrorIBBBadRequest)
		return tag, ctx, true
	}

	ctx.opaque = make(chan bool)

	return tag, ctx, false
}

func bytestreamCalculateDestinationAddress(tag data.BytestreamQuery, stanza *data.ClientIQ) string {
	if tag.DestinationAddress != "" {
		return tag.DestinationAddress
	}
	return hex.EncodeToString(digests.Sha1([]byte(tag.Sid + stanza.From + stanza.To)))
}

const chunkSize = 4096
const cancelCheckFrequency = 100

func (ctx *recvContext) bytestreamCleanup(conn net.Conn, ff *os.File) {
	conn.Close()
	os.Remove(ff.Name())
	removeInflightRecv(ctx.sid)
}

func (ctx *recvContext) bytestreamDoReceive(conn net.Conn) {
	ff, err := ctx.openDestinationTempFile()
	if err != nil {
		ctx.s.Warn(fmt.Sprintf("Failed to open temporary file: %v", err))
		return
	}

	defer ff.Close()

	cancel := ctx.opaque.(chan bool)

	buf := make([]byte, chunkSize)
	totalWritten := int64(0)
	writes := 0
	for {
		if writes%cancelCheckFrequency == 0 {
			select {
			case <-cancel:
				ctx.bytestreamCleanup(conn, ff)
				return
			default:
				// Fall through, since we are not going to cancel
			}
		}

		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				ctx.s.Warn(fmt.Sprintf("Had error when trying to read from connection: %v", err))
				ctx.control.ReportError(errors.New("Error reading from peer"))
				ctx.bytestreamCleanup(conn, ff)
				return
			}
			break
		}
		_, err = ff.Write(buf[:n])
		if err != nil {
			ctx.s.Warn(fmt.Sprintf("Had error when trying to write to file: %v", err))
			ctx.control.ReportError(errors.New("Error writing to file"))
			ctx.bytestreamCleanup(conn, ff)
			return
		}
		totalWritten += int64(n)
		writes++
		ctx.control.SendUpdate(totalWritten)
	}

	fstat, _ := ff.Stat()

	// TODO[LATER]: These checks ignore the range flags - we should think about how that would fit
	if totalWritten != ctx.size || fstat.Size() != totalWritten {
		ctx.s.Warn(fmt.Sprintf("Expected size of file to be %d, but was %d - this probably means the transfer was cancelled", ctx.size, fstat.Size()))
		ctx.control.ReportError(errors.New("Incorrect final size of file - this implies the transfer was cancelled"))
		ctx.bytestreamCleanup(conn, ff)
		return
	}

	// TODO[LATER]: if there's a hash of the file in the inflight, we should calculate it on the file and check it
	if err := ctx.finalizeFileTransfer(ff.Name()); err != nil {
		ctx.s.Warn(fmt.Sprintf("Had error when trying to move the final file: %v", err))
		ctx.bytestreamCleanup(conn, ff)
	}
}

// BytestreamQuery is the hook function that will be called when we receive a bytestream query IQ
func BytestreamQuery(s access.Session, stanza *data.ClientIQ) (ret interface{}, iqtype string, ignore bool) {
	tag, ctx, earlyReturn := bytestreamInitialSetup(s, stanza)
	if earlyReturn {
		return nil, "", true
	}

	dstAddr := bytestreamCalculateDestinationAddress(tag, stanza)

	k := func(c net.Conn) {
		go ctx.bytestreamDoReceive(c)
	}

	for _, sh := range tag.Streamhosts {
		if tryStreamhost(s, sh, dstAddr, k) {
			return data.BytestreamQuery{
				Sid:            tag.Sid,
				StreamhostUsed: &data.BytestreamStreamhostUsed{Jid: sh.Jid},
			}, "result", false
		}
	}

	s.SendIQError(stanza, iqErrorItemNotFound)
	return nil, "", true
}
