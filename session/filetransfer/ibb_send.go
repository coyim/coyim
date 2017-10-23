package filetransfer

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/coyim/coyim/xmpp/data"
)

func init() {
	registerSendFileTransferMethod("http://jabber.org/protocol/ibb", ibbSendDo, ibbSendCurrentlyValid)
}

func ibbSendCurrentlyValid(string, *sendContext) bool {
	return true
}

const ibbDefaultBlockSize = 4096

func ibbSendDo(ctx *sendContext) {
	ctx.ibbSendDoWithBlockSize(ibbDefaultBlockSize)
}

func (ctx *sendContext) ibbSendDoWithBlockSize(blocksize int) {
	nonblockIQ(ctx.s, ctx.peer, "set", data.IBBOpen{
		BlockSize: ibbDefaultBlockSize,
		Sid:       ctx.sid,
		Stanza:    "iq",
	}, nil, func(*data.ClientIQ) {
		go ctx.ibbSendStartTransfer(blocksize)
	}, func(ciq *data.ClientIQ, e error) {
		if ciq != nil &&
			ciq.Type == "error" &&
			ciq.Error.Type == "modify" &&
			ciq.Error.Any.XMLName.Local == "resource-constraint" &&
			ciq.Error.Any.XMLName.Space == "urn:ietf:params:xml:ns:xmpp-stanzas" {
			ctx.ibbSendDoWithBlockSize(blocksize / 2)
			return
		}
		ctx.control.ReportError(e)
		removeInflightSend(ctx)
	})
}

func (ctx *sendContext) ibbSendChunk(r io.ReadCloser, buffer []byte, seq uint16) bool {
	if ctx.weWantToCancel {
		ctx.s.Conn().SendIQ(ctx.peer, "set", data.IBBClose{Sid: ctx.sid})
		removeInflightSend(ctx)
		return false
	} else if ctx.theyWantToCancel {
		ctx.control.CloseAll()
		removeInflightSend(ctx)
		return false
	}

	n, err := r.Read(buffer)
	if err == io.EOF && n == 0 {
		r.Close()
		// TODO[LATER]: we ignore the result of this close - maybe we should react to it in some way, if it reports failure from the other side
		ctx.s.Conn().SendIQ(ctx.peer, "set", data.IBBClose{Sid: ctx.sid})
		ctx.control.ReportFinished()
		removeInflightSend(ctx)
		return false
	} else if err != nil {
		ctx.control.ReportError(err)
		removeInflightSend(ctx)
		return false
	}
	encdata := base64.StdEncoding.EncodeToString(buffer[:n])

	rpl, _, e := ctx.s.Conn().SendIQ(ctx.peer, "set", data.IBBData{
		Sid:      ctx.sid,
		Sequence: seq,
		Base64:   encdata,
	})
	if e != nil {
		ctx.control.ReportError(e)
		removeInflightSend(ctx)
		return false
	}
	ctx.totalSent += int64(n)
	ctx.control.SendUpdate(ctx.totalSent)

	go ctx.trackResultOfSend(rpl)

	return true
}

func (ctx *sendContext) trackResultOfSend(reply <-chan data.Stanza) {
	select {
	case r := <-reply:
		switch ciq := r.Value.(type) {
		case *data.ClientIQ:
			if ciq.Type == "result" {
				return
			}
		}
		ctx.s.Info(fmt.Sprintf("Received unhappy response to IBB data sent: %#v", r))
		ctx.theyWantToCancel = true
	case <-time.After(time.Minute * 5):
		// Ignore timeout
	}
}

func (ctx *sendContext) ibbScheduleNextSend(r io.ReadCloser, buffer []byte, seq uint16) bool {
	time.AfterFunc(time.Duration(200)*time.Millisecond, func() {
		ctx.ibbSendChunks(r, buffer, seq)
	})

	return true
}

func (ctx *sendContext) ibbSendChunks(r io.ReadCloser, buffer []byte, seq uint16) {
	// The seq variable can wrap around here - THAT IS ON PURPOSE
	// See XEP-0047 for details
	ignore := ctx.ibbSendChunk(r, buffer, seq) &&
		ctx.ibbSendChunk(r, buffer, seq+1) &&
		ctx.ibbSendChunk(r, buffer, seq+2) &&
		ctx.ibbSendChunk(r, buffer, seq+3) &&
		ctx.ibbSendChunk(r, buffer, seq+4) &&
		ctx.ibbScheduleNextSend(r, buffer, seq+5)
	ignore = ignore
}

func (ctx *sendContext) ibbSendStartTransfer(blockSize int) {
	seq := uint16(0)
	buffer := make([]byte, blockSize)
	f, err := os.Open(ctx.file)
	if err != nil {
		ctx.control.ReportError(err)
		removeInflightSend(ctx)
		return
	}
	ctx.ibbSendChunks(f, buffer, seq)
}

func (ctx *sendContext) ibbReceivedClose() {
	ctx.theyWantToCancel = true
}
