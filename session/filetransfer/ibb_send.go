package filetransfer

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/coyim/coyim/session/access"
	"github.com/coyim/coyim/xmpp/data"
)

const ibbDefaultBlockSize = 4096

func ibbSendDo(s access.Session, ctx *sendContext) {
	ctx.ibbSendDoWithBlockSize(s, ibbDefaultBlockSize)
}

func (ctx *sendContext) ibbSendDoWithBlockSize(s access.Session, blocksize int) {
	res, _, e := s.Conn().SendIQ(ctx.peer, "set", data.IBBOpen{
		BlockSize: ibbDefaultBlockSize,
		Sid:       ctx.sid,
		Stanza:    "iq",
	})
	if e != nil {
		ctx.control.ReportError(e)
		removeInflightSend(ctx)
		return
	}
	go ctx.ibbSendWaitForConfirmationOfOpen(s, res, blocksize)
}

func (ctx *sendContext) ibbSendChunk(s access.Session, r io.ReadCloser, buffer []byte, seq uint16) bool {
	if ctx.weWantToCancel {
		s.Conn().SendIQ(ctx.peer, "set", data.IBBClose{Sid: ctx.sid})
		removeInflightSend(ctx)
		return false
	} else if ctx.theyWantToCancel {
		close(ctx.control.TransferFinished)
		close(ctx.control.Update)
		close(ctx.control.ErrorOccurred)
		removeInflightSend(ctx)
		return false
	}

	n, err := r.Read(buffer)
	if err == io.EOF && n == 0 {
		r.Close()
		// TODO[LATER]: we ignore the result of this close - maybe we should react to it in some way, if it reports failure from the other side
		s.Conn().SendIQ(ctx.peer, "set", data.IBBClose{Sid: ctx.sid})
		ctx.control.ReportFinished()
		removeInflightSend(ctx)
		return false
	} else if err != nil {
		ctx.control.ReportError(err)
		removeInflightSend(ctx)
		return false
	}
	encdata := base64.StdEncoding.EncodeToString(buffer[:n])

	rpl, _, e := s.Conn().SendIQ(ctx.peer, "set", data.IBBData{
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
	ctx.control.Update <- ctx.totalSent

	go ctx.trackResultOfSend(s, rpl)

	return true
}

func (ctx *sendContext) trackResultOfSend(s access.Session, reply <-chan data.Stanza) {
	select {
	case r := <-reply:
		switch ciq := r.Value.(type) {
		case *data.ClientIQ:
			if ciq.Type == "result" {
				return
			}
		}
		s.Info(fmt.Sprintf("Received unhappy response to IBB data sent: %#v", r))
		ctx.theyWantToCancel = true
	case <-time.After(time.Minute * 5):
		// Ignore timeout
	}
}

func (ctx *sendContext) ibbScheduleNextSend(s access.Session, r io.ReadCloser, buffer []byte, seq uint16) bool {
	time.AfterFunc(time.Duration(200)*time.Millisecond, func() {
		ctx.ibbSendChunks(s, r, buffer, seq)
	})

	return true
}

func (ctx *sendContext) ibbSendChunks(s access.Session, r io.ReadCloser, buffer []byte, seq uint16) {
	// The seq variable can wrap around here - THAT IS ON PURPOSE
	// See XEP-0047 for details
	ignore := ctx.ibbSendChunk(s, r, buffer, seq) &&
		ctx.ibbSendChunk(s, r, buffer, seq+1) &&
		ctx.ibbSendChunk(s, r, buffer, seq+2) &&
		ctx.ibbSendChunk(s, r, buffer, seq+3) &&
		ctx.ibbSendChunk(s, r, buffer, seq+4) &&
		ctx.ibbScheduleNextSend(s, r, buffer, seq+5)
	ignore = ignore
}

func (ctx *sendContext) ibbSendStartTransfer(s access.Session, blockSize int) {
	seq := uint16(0)
	buffer := make([]byte, blockSize)
	f, err := os.Open(ctx.file)
	if err != nil {
		ctx.control.ReportError(err)
		removeInflightSend(ctx)
		return
	}
	ctx.ibbSendChunks(s, f, buffer, seq)
}

func (ctx *sendContext) ibbSendWaitForConfirmationOfOpen(s access.Session, reply <-chan data.Stanza, blockSize int) {
	r, ok := <-reply
	if !ok {
		ctx.control.ReportError(errors.New("We didn't receive a response when trying to open IBB file transfer with peer"))
		removeInflightSend(ctx)
		return
	}

	switch ciq := r.Value.(type) {
	case *data.ClientIQ:
		if ciq.Type == "result" {
			go ctx.ibbSendStartTransfer(s, blockSize)
			return
		} else if ciq.Type == "error" {
			if ciq.Error.Type == "cancel" {
				ctx.control.ReportErrorNonblocking(errors.New("The peer canceled the file transfer"))
			} else if ciq.Error.Type == "modify" &&
				ciq.Error.Any.XMLName.Local == "resource-constraint" &&
				ciq.Error.Any.XMLName.Space == "urn:ietf:params:xml:ns:xmpp-stanzas" {
				ctx.ibbSendDoWithBlockSize(s, blockSize/2)
			} else {
				ctx.control.ReportErrorNonblocking(errors.New("Invalid error type - this shouldn't happen"))
			}
		} else {
			ctx.control.ReportErrorNonblocking(errors.New("Invalid IQ type - this shouldn't happen"))
		}
	default:
		ctx.control.ReportErrorNonblocking(errors.New("Invalid stanza type - this shouldn't happen"))
	}
	removeInflightSend(ctx)
}

func (ctx *sendContext) ibbReceivedClose(s access.Session) {
	ctx.theyWantToCancel = true
}
