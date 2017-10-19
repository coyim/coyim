package filetransfer

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/coyim/coyim/session/access"
	"github.com/coyim/coyim/xmpp/data"
)

// TOOD: at some point this should be refactored away into a pure IBB implementation and a small piece that is file transfer specific

type ibbContext struct {
	sync.Mutex
	sid               string
	blockSize         int
	stanza            string
	expectingSequence uint16
	currentSize       int64
	f                 *os.File
}

var iqErrorNotAcceptable = data.ErrorReply{
	Type:  "cancel",
	Error: data.ErrorNotAcceptable{},
}

var iqErrorItemNotFound = data.ErrorReply{
	Type:  "cancel",
	Error: data.ErrorItemNotFound{},
}

var iqErrorUnexpectedRequest = data.ErrorReply{
	Type:  "cancel",
	Error: data.ErrorUnexpectedRequest{},
}

var iqErrorIBBBadRequest = data.ErrorReply{
	Type:  "cancel",
	Error: data.ErrorBadRequest{},
}

func (ift inflight) ibbCleanup(lock bool) {
	ctx, ok := ift.status.opaque.(*ibbContext)
	if ok {
		if lock {
			ctx.Lock()
			defer ctx.Unlock()
		}

		if ctx.f != nil {
			ctx.f.Close() // we ignore any errors here - if the file is already closed, that's OK
			os.Remove(ctx.f.Name())
		}
	}
	removeInflight(ift.id)
}

func ibbWaitForCancel(s access.Session, ift inflight) {
	ift.control.WaitForCancel(func() {
		ift.ibbCleanup(true)
		s.Conn().SendIQ(ift.peer, "set", data.IBBClose{Sid: ift.id})
	})
}

// IbbOpen is the hook function that will be called when we receive an ibb open IQ
func IbbOpen(s access.Session, stanza *data.ClientIQ) (ret interface{}, iqtype string, ignore bool) {
	var tag data.IBBOpen
	if err := xml.NewDecoder(bytes.NewBuffer(stanza.Query)).Decode(&tag); err != nil {
		s.Warn(fmt.Sprintf("Failed to parse IBB open: %v", err))
		return iqErrorNotAcceptable, "error", false
	}

	inflight, ok := getInflight(tag.Sid)

	if !ok || inflight.status.opaque != nil {
		s.Warn(fmt.Sprintf("No file transfer associated with SID: %v", tag.Sid))
		return iqErrorNotAcceptable, "error", false
	}

	c := &ibbContext{
		sid:               tag.Sid,
		blockSize:         tag.BlockSize,
		stanza:            tag.Stanza,
		expectingSequence: 0,
	}
	if c.stanza == "" {
		c.stanza = "iq"
	}
	inflight.status.opaque = c

	ff, err := inflight.openDestinationTempFile()
	if err != nil {
		s.Warn(fmt.Sprintf("Failed to open temporary file: %v", err))
		return iqErrorNotAcceptable, "error", false
	}
	c.f = ff

	return data.EmptyReply{}, "", false
}

// IbbData is the hook function that will be called when we receive an ibb data IQ
func IbbData(s access.Session, stanza *data.ClientIQ) (ret interface{}, iqtype string, ignore bool) {
	var tag data.IBBData
	if err := xml.NewDecoder(bytes.NewBuffer(stanza.Query)).Decode(&tag); err != nil {
		s.Warn(fmt.Sprintf("Failed to parse IBB data: %v", err))
		return iqErrorNotAcceptable, "error", false
	}

	inflight, ok := getInflight(tag.Sid)

	if !ok || inflight.status.opaque == nil {
		s.Warn(fmt.Sprintf("No file transfer associated with SID: %v", tag.Sid))
		return iqErrorItemNotFound, "error", false
	}

	ctx, ok := inflight.status.opaque.(*ibbContext)
	if !ok {
		s.Warn(fmt.Sprintf("No IBB file transfer associated with SID: %v", tag.Sid))
		return iqErrorItemNotFound, "error", false
	}

	ctx.Lock()
	defer ctx.Unlock()

	// XEP-0047 wants us to keep track of previously used sequence numbers, and only do this error
	// when a sequence number is reused - otherwise we should immediately close the stream.
	// However, because of the wraparound behavior of "seq" - also specified in XEP-0047, for large
	// files we can't actually tell the difference between a reused sequence number or a number that
	// has just been wrapped around. Thus, we do this deviation from the spec here.
	if tag.Sequence != ctx.expectingSequence {
		s.Warn(fmt.Sprintf("IBB expected sequence number %d, but got %d", ctx.expectingSequence, tag.Sequence))
		inflight.control.ReportError(errors.New("Unexpected data sent from the peer"))
		inflight.ibbCleanup(false)
		return iqErrorUnexpectedRequest, "error", false
	}

	ctx.expectingSequence++ // wraparound on purpose, to match uint16 spec behavior of the seq field

	result, err := base64.StdEncoding.DecodeString(tag.Base64)
	if err != nil {
		s.Warn(fmt.Sprintf("IBB received corrupt data for sequence %d", tag.Sequence))
		inflight.control.ReportError(errors.New("Corrupt data sent by the peer"))
		inflight.ibbCleanup(false)
		return iqErrorIBBBadRequest, "error", false

	}

	var n int
	if n, err = ctx.f.Write(result); err != nil {
		s.Warn(fmt.Sprintf("IBB had an error when writing to the file: %v", err))
		inflight.control.ReportError(errors.New("Couldn't write data to the file system"))
		inflight.ibbCleanup(false)
		return iqErrorNotAcceptable, "error", false
	}
	ctx.currentSize += int64(n)
	inflight.control.SendUpdate(ctx.currentSize)

	return data.EmptyReply{}, "", false
}

// IbbMessageData is the hook function that will be called when we receive a message containing an ibb data
func IbbMessageData(s access.Session, stanza *data.ClientMessage, ext *data.Extension) {
	var tag data.IBBData
	if err := xml.NewDecoder(bytes.NewBuffer([]byte(ext.Body))).Decode(&tag); err != nil {
		s.Warn(fmt.Sprintf("Failed to parse IBB data: %v", err))
		return
	}

	inflight, ok := getInflight(tag.Sid)

	if !ok || inflight.status.opaque == nil {
		s.Warn(fmt.Sprintf("No file transfer associated with SID: %v", tag.Sid))
		// we can't actually send anything back to indicate this problem...
		return
	}

	ctx, ok := inflight.status.opaque.(*ibbContext)
	if !ok {
		s.Warn(fmt.Sprintf("No IBB file transfer associated with SID: %v", tag.Sid))
		// we can't actually send anything back to indicate this problem...
		return
	}

	ctx.Lock()
	defer ctx.Unlock()

	// XEP-0047 wants us to keep track of previously used sequence numbers, and only do this error
	// when a sequence number is reused - otherwise we should immediately close the stream.
	// However, because of the wraparound behavior of "seq" - also specified in XEP-0047, for large
	// files we can't actually tell the difference between a reused sequence number or a number that
	// has just been wrapped around. Thus, we do this deviation from the spec here.
	if tag.Sequence != ctx.expectingSequence {
		s.Warn(fmt.Sprintf("IBB expected sequence number %d, but got %d", ctx.expectingSequence, tag.Sequence))
		// we can't actually send anything back to indicate this problem...
		inflight.control.ReportError(errors.New("Unexpected data sent from the peer"))
		inflight.ibbCleanup(false)
		return
	}

	ctx.expectingSequence++ // wraparound on purpose, to match uint16 spec behavior of the seq field

	result, err := base64.StdEncoding.DecodeString(tag.Base64)
	if err != nil {
		s.Warn(fmt.Sprintf("IBB received corrupt data for sequence %d", tag.Sequence))
		// we can't actually send anything back to indicate this problem...
		inflight.control.ReportError(errors.New("Corrupt data sent by the peer"))
		inflight.ibbCleanup(false)
		return

	}

	var n int
	if n, err = ctx.f.Write(result); err != nil {
		s.Warn(fmt.Sprintf("IBB had an error when writing to the file: %v", err))
		inflight.control.ReportError(errors.New("Couldn't write data to the file system"))
		inflight.ibbCleanup(false)
		return
	}
	ctx.currentSize += int64(n)
	inflight.control.SendUpdate(ctx.currentSize)
}

// IbbClose is the hook function that will be called when we receive an ibb close IQ
func IbbClose(s access.Session, stanza *data.ClientIQ) (ret interface{}, iqtype string, ignore bool) {
	var tag data.IBBClose
	if err := xml.NewDecoder(bytes.NewBuffer(stanza.Query)).Decode(&tag); err != nil {
		s.Warn(fmt.Sprintf("Failed to parse IBB close: %v", err))
		return iqErrorNotAcceptable, "error", false
	}

	inflightSend, ok := getInflightSend(tag.Sid)
	if ok {
		inflightSend.ibbReceivedClose(s)
		return data.EmptyReply{}, "", false
	}

	inflight, ok := getInflight(tag.Sid)

	if !ok || inflight.status.opaque == nil {
		s.Warn(fmt.Sprintf("No file transfer associated with SID: %v", tag.Sid))
		return iqErrorItemNotFound, "error", false
	}

	ctx, ok := inflight.status.opaque.(*ibbContext)
	if !ok {
		s.Warn(fmt.Sprintf("No IBB file transfer associated with SID: %v", tag.Sid))
		return iqErrorItemNotFound, "error", false
	}

	ctx.Lock()
	defer ctx.Unlock()

	defer ctx.f.Close()
	fstat, _ := ctx.f.Stat()

	// TODO[LATER]: These checks ignore the range flags - we should think about how that would fit
	if ctx.currentSize != inflight.size || fstat.Size() != ctx.currentSize {
		s.Warn(fmt.Sprintf("Expected size of file to be %d, but was %d - this probably means the transfer was cancelled", inflight.size, fstat.Size()))
		inflight.control.ReportError(errors.New("Incorrect final size of file - this implies the transfer was cancelled"))
		inflight.ibbCleanup(false)
		return data.EmptyReply{}, "", false
	}

	// TODO[LATER]: if there's a hash of the file in the inflight, we should calculate it on the file and check it

	if err := inflight.finalizeFileTransfer(ctx.f.Name()); err != nil {
		s.Warn(fmt.Sprintf("Had error when trying to move the final file: %v", err))
		inflight.ibbCleanup(false)
	}

	return data.EmptyReply{}, "", false
}
