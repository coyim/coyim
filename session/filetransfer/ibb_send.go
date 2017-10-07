package filetransfer

import (
	"encoding/base64"
	"errors"
	"io"
	"os"
	"time"

	"github.com/coyim/coyim/session/access"
	"github.com/coyim/coyim/xmpp/data"
)

// TODO - 2 - make cancel work based on different failures
// TODO - 4 - multiplex the close tag

const ibbDefaultBlockSize = 4096

func ibbSendDoWithBlockSize(s access.Session, ctx *sendContext, blocksize int) {
	res, _, e := s.Conn().SendIQ(ctx.peer, "set", data.IBBOpen{
		BlockSize: ibbDefaultBlockSize,
		Sid:       ctx.sid,
		Stanza:    "iq",
	})
	if e != nil {
		ctx.control.ReportError(e)
		return
	}
	go ibbSendWaitForConfirmationOfOpen(s, ctx, res, blocksize)
}

func ibbSendDo(s access.Session, ctx *sendContext) {
	ibbSendDoWithBlockSize(s, ctx, ibbDefaultBlockSize)
}

// TODO: Print everything received and sent
// TODO: We need to multiplex the Close IBB tag...

func ibbSendChunk(s access.Session, ctx *sendContext, r io.ReadCloser, buffer []byte, seq uint16) bool {
	// TODO: check for cancel here

	n, err := r.Read(buffer)
	if err == io.EOF && n == 0 {
		r.Close()
		// TODO[LATER]: we ignore the result of this close - maybe we should react to it in some way, if it reports failure from the other side
		s.Conn().SendIQ(ctx.peer, "set", data.IBBClose{Sid: ctx.sid})
		return false
	} else if err != nil {
		ctx.control.ReportError(err)
		return false
	}
	encdata := base64.StdEncoding.EncodeToString(buffer[:n])

	// TODO: we should keep track of each response here
	_, _, e := s.Conn().SendIQ(ctx.peer, "set", data.IBBData{
		Sid:      ctx.sid,
		Sequence: seq,
		Base64:   encdata,
	})
	if e != nil {
		ctx.control.ReportError(e)
		return false
	}
	return true
}

func ibbScheduleNextSend(s access.Session, ctx *sendContext, r io.ReadCloser, buffer []byte, seq uint16) bool {
	time.AfterFunc(time.Duration(200)*time.Millisecond, func() {
		ibbSendChunks(s, ctx, r, buffer, seq)
	})

	return true
}

func ibbSendChunks(s access.Session, ctx *sendContext, r io.ReadCloser, buffer []byte, seq uint16) {
	// The seq variable can wrap around here - THAT IS ON PURPOSE
	// See XEP-0047 for details
	ignore := ibbSendChunk(s, ctx, r, buffer, seq) &&
		ibbSendChunk(s, ctx, r, buffer, seq+1) &&
		ibbSendChunk(s, ctx, r, buffer, seq+2) &&
		ibbSendChunk(s, ctx, r, buffer, seq+3) &&
		ibbSendChunk(s, ctx, r, buffer, seq+4) &&
		ibbScheduleNextSend(s, ctx, r, buffer, seq+5)
	ignore = ignore
}

func ibbSendStartTransfer(s access.Session, ctx *sendContext, blockSize int) {
	//	cancel := make(chan bool)
	seq := uint16(0)
	buffer := make([]byte, blockSize)
	f, err := os.Open(ctx.file)
	if err != nil {
		ctx.control.ReportError(err)
		return
	}
	ibbSendChunks(s, ctx, f, buffer, seq)
}

func ibbSendWaitForConfirmationOfOpen(s access.Session, ctx *sendContext, reply <-chan data.Stanza, blockSize int) {
	r, ok := <-reply
	if !ok {
		ctx.control.ReportError(errors.New("We didn't receive a response when trying to open IBB file transfer with peer"))
		return
	}

	switch ciq := r.Value.(type) {
	case *data.ClientIQ:
		if ciq.Type == "result" {
			go ibbSendStartTransfer(s, ctx, blockSize)
		} else if ciq.Type == "error" {
			if ciq.Error.Type == "cancel" {
				ctx.control.ReportErrorNonblocking(errors.New("The peer canceled the file transfer"))
			} else if ciq.Error.Type == "modify" &&
				ciq.Error.Any.XMLName.Local == "resource-constraint" &&
				ciq.Error.Any.XMLName.Space == "urn:ietf:params:xml:ns:xmpp-stanzas" {
				ibbSendDoWithBlockSize(s, ctx, blockSize/2)
			} else {
				ctx.control.ReportErrorNonblocking(errors.New("Invalid error type - this shouldn't happen"))
			}
		} else {
			ctx.control.ReportErrorNonblocking(errors.New("Invalid IQ type - this shouldn't happen"))
		}
	default:
		ctx.control.ReportErrorNonblocking(errors.New("Invalid stanza type - this shouldn't happen"))
	}
}
