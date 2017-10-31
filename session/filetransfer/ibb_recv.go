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
	expectingSequence uint16
	currentSize       int64
	f                 *os.File
}

func (ctx *recvContext) ibbCleanup(lock bool) {
	ictx, ok := ctx.opaque.(*ibbContext)
	if ok {
		if lock {
			ictx.Lock()
			defer ictx.Unlock()
		}

		if ictx.f != nil {
			ictx.f.Close() // we ignore any errors here - if the file is already closed, that's OK
			os.Remove(ictx.f.Name())
		}
	}
	removeInflightRecv(ctx.sid)
}

func ibbWaitForCancel(ctx *recvContext) {
	ctx.control.WaitForCancel(func() {
		ctx.ibbCleanup(true)
		ctx.s.Conn().SendIQ(ctx.peer, "set", data.IBBClose{Sid: ctx.sid})
	})
}

// IbbOpen is the hook function that will be called when we receive an ibb open IQ
func IbbOpen(s access.Session, stanza *data.ClientIQ) (ret interface{}, iqtype string, ignore bool) {
	var tag data.IBBOpen
	if err := xml.NewDecoder(bytes.NewBuffer(stanza.Query)).Decode(&tag); err != nil {
		s.Warn(fmt.Sprintf("Failed to parse IBB open: %v", err))
		return iqErrorNotAcceptable, "error", false
	}

	ctx, ok := getInflightRecv(tag.Sid)

	if !ok || ctx.opaque != nil {
		s.Warn(fmt.Sprintf("No file transfer associated with SID: %v", tag.Sid))
		return iqErrorNotAcceptable, "error", false
	}

	c := &ibbContext{}
	ctx.opaque = c

	ff, err := ctx.openDestinationTempFile()
	if err != nil {
		s.Warn(fmt.Sprintf("Failed to open temporary file: %v", err))
		return iqErrorNotAcceptable, "error", false
	}
	c.f = ff

	return data.EmptyReply{}, "", false
}

// TODO: continue refactoring the handling of data to be generic for all

func ibbHandleData(s access.Session, data []byte) (tag data.IBBData, ctx *recvContext, ictx *ibbContext, ret interface{}, iqtype string, ignore bool) {
	if err := xml.NewDecoder(bytes.NewBuffer(data)).Decode(&tag); err != nil {
		s.Warn(fmt.Sprintf("Failed to parse IBB data: %v", err))
		return tag, nil, nil, iqErrorNotAcceptable, "error", false
	}

	ctx, ok := getInflightRecv(tag.Sid)

	if !ok || ctx.opaque == nil {
		s.Warn(fmt.Sprintf("No file transfer associated with SID: %v", tag.Sid))
		return tag, nil, nil, iqErrorItemNotFound, "error", false
	}

	ictx, ok = ctx.opaque.(*ibbContext)
	if !ok {
		s.Warn(fmt.Sprintf("No IBB file transfer associated with SID: %v", tag.Sid))
		return tag, nil, nil, iqErrorItemNotFound, "error", false
	}

	ictx.Lock()
	return tag, ctx, ictx, nil, "", false
}

// IbbData is the hook function that will be called when we receive an ibb data IQ
func IbbData(s access.Session, stanza *data.ClientIQ) (ret interface{}, iqtype string, ignore bool) {
	tag, ctx, ictx, ret, iqtype, ignore := ibbHandleData(s, stanza.Query)
	if ret != nil {
		return ret, iqtype, ignore
	}
	defer ictx.Unlock()

	// XEP-0047 wants us to keep track of previously used sequence numbers, and only do this error
	// when a sequence number is reused - otherwise we should immediately close the stream.
	// However, because of the wraparound behavior of "seq" - also specified in XEP-0047, for large
	// files we can't actually tell the difference between a reused sequence number or a number that
	// has just been wrapped around. Thus, we do this deviation from the spec here.
	if tag.Sequence != ictx.expectingSequence {
		s.Warn(fmt.Sprintf("IBB expected sequence number %d, but got %d", ictx.expectingSequence, tag.Sequence))
		ctx.control.ReportError(errors.New("Unexpected data sent from the peer"))
		ctx.ibbCleanup(false)
		return iqErrorUnexpectedRequest, "error", false
	}

	ictx.expectingSequence++ // wraparound on purpose, to match uint16 spec behavior of the seq field

	result, err := base64.StdEncoding.DecodeString(tag.Base64)
	if err != nil {
		s.Warn(fmt.Sprintf("IBB received corrupt data for sequence %d", tag.Sequence))
		ctx.control.ReportError(errors.New("Corrupt data sent by the peer"))
		ctx.ibbCleanup(false)
		return iqErrorIBBBadRequest, "error", false

	}

	var n int
	if n, err = ictx.f.Write(result); err != nil {
		s.Warn(fmt.Sprintf("IBB had an error when writing to the file: %v", err))
		ctx.control.ReportError(errors.New("Couldn't write data to the file system"))
		ctx.ibbCleanup(false)
		return iqErrorNotAcceptable, "error", false
	}
	ictx.currentSize += int64(n)
	ctx.control.SendUpdate(ictx.currentSize, ctx.size)

	return data.EmptyReply{}, "", false
}

// IbbMessageData is the hook function that will be called when we receive a message containing an ibb data
func IbbMessageData(s access.Session, stanza *data.ClientMessage, ext *data.Extension) {
	tag, ctx, ictx, ret, _, _ := ibbHandleData(s, []byte(ext.Body))
	if ret != nil {
		return
	}
	defer ictx.Unlock()

	// XEP-0047 wants us to keep track of previously used sequence numbers, and only do this error
	// when a sequence number is reused - otherwise we should immediately close the stream.
	// However, because of the wraparound behavior of "seq" - also specified in XEP-0047, for large
	// files we can't actually tell the difference between a reused sequence number or a number that
	// has just been wrapped around. Thus, we do this deviation from the spec here.
	if tag.Sequence != ictx.expectingSequence {
		s.Warn(fmt.Sprintf("IBB expected sequence number %d, but got %d", ictx.expectingSequence, tag.Sequence))
		// we can't actually send anything back to indicate this problem...
		ctx.control.ReportError(errors.New("Unexpected data sent from the peer"))
		ctx.ibbCleanup(false)
		return
	}

	ictx.expectingSequence++ // wraparound on purpose, to match uint16 spec behavior of the seq field

	result, err := base64.StdEncoding.DecodeString(tag.Base64)
	if err != nil {
		s.Warn(fmt.Sprintf("IBB received corrupt data for sequence %d", tag.Sequence))
		// we can't actually send anything back to indicate this problem...
		ctx.control.ReportError(errors.New("Corrupt data sent by the peer"))
		ctx.ibbCleanup(false)
		return

	}

	var n int
	if n, err = ictx.f.Write(result); err != nil {
		s.Warn(fmt.Sprintf("IBB had an error when writing to the file: %v", err))
		ctx.control.ReportError(errors.New("Couldn't write data to the file system"))
		ctx.ibbCleanup(false)
		return
	}
	ictx.currentSize += int64(n)
	ctx.control.SendUpdate(ictx.currentSize, ctx.size)
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
		ibbReceivedClose(inflightSend)
		return data.EmptyReply{}, "", false
	}

	ctx, ok := getInflightRecv(tag.Sid)

	if !ok || ctx.opaque == nil {
		s.Warn(fmt.Sprintf("No file transfer associated with SID: %v", tag.Sid))
		return iqErrorItemNotFound, "error", false
	}

	ictx, ok := ctx.opaque.(*ibbContext)
	if !ok {
		s.Warn(fmt.Sprintf("No IBB file transfer associated with SID: %v", tag.Sid))
		return iqErrorItemNotFound, "error", false
	}

	ictx.Lock()
	defer ictx.Unlock()

	defer ictx.f.Close()
	fstat, _ := ictx.f.Stat()

	// TODO[LATER]: These checks ignore the range flags - we should think about how that would fit
	if ictx.currentSize != ctx.size || fstat.Size() != ictx.currentSize {
		s.Warn(fmt.Sprintf("Expected size of file to be %d, but was %d - this probably means the transfer was cancelled", ctx.size, fstat.Size()))
		ctx.control.ReportError(errors.New("Incorrect final size of file - this implies the transfer was cancelled"))
		ctx.ibbCleanup(false)
		return data.EmptyReply{}, "", false
	}

	// TODO[LATER]: if there's a hash of the file in the inflight, we should calculate it on the file and check it

	if err := ctx.finalizeFileTransfer(ictx.f.Name()); err != nil {
		s.Warn(fmt.Sprintf("Had error when trying to move the final file: %v", err))
		ctx.ibbCleanup(false)
	}

	return data.EmptyReply{}, "", false
}
