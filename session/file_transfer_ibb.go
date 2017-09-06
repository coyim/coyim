package session

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/twstrike/coyim/xmpp/data"
)

// TODO: receiving packets via messages

func init() {
	registerKnownIQ("set", "http://jabber.org/protocol/ibb open", fileTransferIbbOpen)
	registerKnownIQ("set", "http://jabber.org/protocol/ibb data", fileTransferIbbData)
	registerKnownIQ("set", "http://jabber.org/protocol/ibb close", fileTransferIbbClose)
}

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

func fileTransferIbbOpen(s *session, stanza *data.ClientIQ) (ret interface{}, ignore bool) {
	var tag data.IBBOpen
	if err := xml.NewDecoder(bytes.NewBuffer(stanza.Query)).Decode(&tag); err != nil {
		s.warn(fmt.Sprintf("Failed to parse IBB open: %v", err))
		return nil, false
	}

	inflight, ok := getInflightFileTransfer(tag.Sid)

	if !ok || inflight.status.opaque != nil {
		s.warn(fmt.Sprintf("No file transfer associated with SID: %v", tag.Sid))
		s.sendIQError(stanza, iqErrorNotAcceptable)
		return nil, true
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

	ff, err := ioutil.TempFile("", "coyim_file_transfer")
	if err != nil {
		inflight.status.opaque = nil
		s.warn(fmt.Sprintf("Failed to open temporary file: %v", err))
		return nil, false
	}
	c.f = ff

	return data.EmptyReply{}, false
}

func fileTransferIbbData(s *session, stanza *data.ClientIQ) (ret interface{}, ignore bool) {
	var tag data.IBBData
	if err := xml.NewDecoder(bytes.NewBuffer(stanza.Query)).Decode(&tag); err != nil {
		s.warn(fmt.Sprintf("Failed to parse IBB data: %v", err))
		return nil, false
	}

	inflight, ok := getInflightFileTransfer(tag.Sid)

	if !ok || inflight.status.opaque == nil {
		s.warn(fmt.Sprintf("No file transfer associated with SID: %v", tag.Sid))
		s.sendIQError(stanza, iqErrorItemNotFound)
		return nil, true
	}

	ctx, ok := inflight.status.opaque.(*ibbContext)
	if !ok {
		s.warn(fmt.Sprintf("No IBB file transfer associated with SID: %v", tag.Sid))
		s.sendIQError(stanza, iqErrorItemNotFound)
		return nil, true
	}

	ctx.Lock()
	defer ctx.Unlock()

	// XEP-0047 wants us to keep track of previously used sequence numbers, and only do this error
	// when a sequence number is reused - otherwise we should immediately close the stream.
	// However, because of the wraparound behavior of "seq" - also specified in XEP-0047, for large
	// files we can't actually tell the difference between a reused sequence number or a number that
	// has just been wrapped around. Thus, we do this deviation from the spec here.
	if tag.Sequence != ctx.expectingSequence {
		s.warn(fmt.Sprintf("IBB expected sequence number %d, but got %d", ctx.expectingSequence, tag.Sequence))
		s.sendIQError(stanza, iqErrorUnexpectedRequest)
		return nil, true
	}

	ctx.expectingSequence++ // wraparound on purpose, to match uint16 spec behavior of the seq field

	result, err := base64.StdEncoding.DecodeString(tag.Base64)
	if err != nil {
		s.warn(fmt.Sprintf("IBB received corrupt data for sequence %d", tag.Sequence))
		s.sendIQError(stanza, iqErrorIBBBadRequest)
		return nil, true

	}

	var n int
	if n, err = ctx.f.Write(result); err != nil {
		s.warn(fmt.Sprintf("IBB had an error when writing to the file: %v", err))
		return nil, false
	}
	ctx.currentSize += int64(n)

	return nil, false
}

func fileTransferIbbClose(s *session, stanza *data.ClientIQ) (ret interface{}, ignore bool) {
	var tag data.IBBClose
	if err := xml.NewDecoder(bytes.NewBuffer(stanza.Query)).Decode(&tag); err != nil {
		s.warn(fmt.Sprintf("Failed to parse IBB close: %v", err))
		return nil, false
	}

	inflight, ok := getInflightFileTransfer(tag.Sid)

	if !ok || inflight.status.opaque == nil {
		s.warn(fmt.Sprintf("No file transfer associated with SID: %v", tag.Sid))
		s.sendIQError(stanza, iqErrorItemNotFound)
		return nil, true
	}

	ctx, ok := inflight.status.opaque.(*ibbContext)
	if !ok {
		s.warn(fmt.Sprintf("No IBB file transfer associated with SID: %v", tag.Sid))
		s.sendIQError(stanza, iqErrorItemNotFound)
		return nil, true
	}

	ctx.Lock()
	defer ctx.Unlock()

	defer ctx.f.Close()
	fstat, _ := ctx.f.Stat()

	// TODO: These checks ignore the range flags - we should think about how that would fit
	if ctx.currentSize != inflight.size || fstat.Size() != ctx.currentSize {
		s.warn(fmt.Sprintf("Expected sze of file to be %d, but was %d", inflight.size, fstat.Size()))
		return data.EmptyReply{}, false
	}

	// TODO: if there's a hash of the file in the inflight, we should calculate it on the file and check it

	fmt.Printf("WE HAVE A FILE AT: %s\n", ctx.f.Name())

	// TODO: move the file to its final name

	removeInflightFileTransfer(tag.Sid)

	return data.EmptyReply{}, false
}
