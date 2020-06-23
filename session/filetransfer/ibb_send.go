package filetransfer

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/coyim/coyim/session/access"
	"github.com/coyim/coyim/xmpp/data"
)

// IBBMethod contains the profile name for IBB
const IBBMethod = "http://jabber.org/protocol/ibb"

func init() {
	registerSendFileTransferMethod(IBBMethod, ibbSendDo, ibbSendCurrentlyValid)
}

func ibbSendCurrentlyValid(string, access.Session) bool {
	return true
}

const ibbDefaultBlockSize = 4096

func ibbSendDo(ctx *sendContext) {
	ibbSendDoWithBlockSize(ctx, ibbDefaultBlockSize)
}

func ibbSendDoWithBlockSize(ctx *sendContext, blocksize int) {
	nonblockIQ(ctx.s, ctx.peer, "set", data.IBBOpen{
		BlockSize: ibbDefaultBlockSize,
		Sid:       ctx.sid,
		Stanza:    "iq",
	}, nil, func(*data.ClientIQ) {
		go ibbSendStartTransfer(ctx, blocksize)
	}, func(ciq *data.ClientIQ, e error) {
		if ciq != nil &&
			ciq.Type == "error" &&
			ciq.Error.Type == "modify" &&
			ciq.Error.Condition.XMLName.Local == "resource-constraint" &&
			ciq.Error.Condition.XMLName.Space == "urn:ietf:params:xml:ns:xmpp-stanzas" {
			ibbSendDoWithBlockSize(ctx, blocksize/2)
			return
		}
		ctx.onError(e)
	})
}

func ibbReadAndWrite(ctx *sendContext) io.ReadCloser {
	f, err := os.Open(ctx.file)
	if err != nil {
		ctx.onError(err)
		return nil
	}

	r, w := io.Pipe()

	ctx.totalSize = ctx.enc.totalSize(ctx.size)

	go func() {
		w2, beforeFinish := ctx.enc.wrapForSending(w, w)
		_, err = io.Copy(w2, f)
		if err != nil && err != errLocalCancel {
			ctx.onError(err)
		}

		beforeFinish()

		_ = w2.Close()
		_ = w.Close()
	}()

	return r
}

func ibbSendChunk(ctx *sendContext, r io.ReadCloser, buffer []byte, seq uint16) bool {
	if ctx.weWantToCancel {
		_, _, _ = ctx.s.Conn().SendIQ(ctx.peer, "set", data.IBBClose{Sid: ctx.sid})
		removeInflightSend(ctx)
		return false
	} else if ctx.theyWantToCancel {
		ctx.control.CloseAll()
		removeInflightSend(ctx)
		return false
	}

	n, err := r.Read(buffer)
	if n > 0 {
		rpl, _, e := ctx.s.Conn().SendIQ(ctx.peer, "set", data.IBBData{
			Sid:      ctx.sid,
			Sequence: seq,
			Base64:   base64.StdEncoding.EncodeToString(buffer[:n]),
		})
		if e != nil {
			ctx.onError(e)
			return false
		}
		ctx.onUpdate(n)

		go trackResultOfSend(ctx, rpl)
	}
	if err == io.EOF {
		closeAndIgnore(r)
		addInflightMAC(ctx)
		_, _, _ = ctx.s.Conn().SendIQ(ctx.peer, "set", data.IBBClose{Sid: ctx.sid})
		ctx.onFinish()
		return false
	} else if err != nil {
		ctx.onError(err)
		return false
	}

	return true
}

func trackResultOfSend(ctx *sendContext, reply <-chan data.Stanza) {
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

var ibbScheduleSendLimit = time.Duration(500) * time.Millisecond

func ibbScheduleNextSend(ctx *sendContext, r io.ReadCloser, buffer []byte, seq uint16) bool {
	time.AfterFunc(ibbScheduleSendLimit, func() {
		ibbSendChunks(ctx, r, buffer, seq)
	})

	return true
}

func ibbSendChunks(ctx *sendContext, r io.ReadCloser, buffer []byte, seq uint16) {
	// The seq variable can wrap around here - THAT IS ON PURPOSE
	// See XEP-0047 for details
	_ = ibbSendChunk(ctx, r, buffer, seq) &&
		ibbScheduleNextSend(ctx, r, buffer, seq+1)
}

func ibbSendStartTransfer(ctx *sendContext, blockSize int) {
	seq := uint16(0)
	buffer := make([]byte, blockSize)
	r := ibbReadAndWrite(ctx)
	if r == nil {
		return
	}

	ibbSendChunks(ctx, r, buffer, seq)
}

func ibbReceivedClose(ctx *sendContext) {
	ctx.theyWantToCancel = true
}
