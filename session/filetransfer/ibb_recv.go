package filetransfer

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/coyim/coyim/session/access"
	"github.com/coyim/coyim/xmpp/data"
)

func init() {
	registerRecieveFileTransferMethod(IBBMethod, 0, ibbWaitForCancel)
}

type ibbContext struct {
	expectingSequence uint16

	recv *receiver
}

func (ctx *recvContext) ibbCleanup() {
	removeInflightRecv(ctx.sid)
}

func ibbWaitForCancel(ctx *recvContext) {
	ctx.control.WaitForCancel(func() {
		ictx, ok := ctx.opaque.(*ibbContext)
		if ok && ictx != nil && ictx.recv != nil {
			ictx.recv.cancel()
		}
		ctx.ibbCleanup()
		_, _, _ = ctx.s.Conn().SendIQ(ctx.peer, "set", data.IBBClose{Sid: ctx.sid})
	})
}

// IbbOpen is the hook function that will be called when we receive an ibb open IQ
func IbbOpen(s access.Session, stanza *data.ClientIQ) (ret interface{}, iqtype string, ignore bool) {
	var tag data.IBBOpen
	if err := xml.NewDecoder(bytes.NewBuffer(stanza.Query)).Decode(&tag); err != nil {
		s.Log().WithError(err).Warn("Failed to parse IBB open")
		return iqErrorNotAcceptable, "error", false
	}

	ctx, ok := getInflightRecv(tag.Sid)

	if !ok || ctx.opaque != nil {
		s.Log().WithField("SID", tag.Sid).Warn("No file transfer associated with SID")
		return iqErrorNotAcceptable, "error", false
	}

	c := &ibbContext{}
	ctx.opaque = c

	c.recv = ctx.createReceiver()

	return data.EmptyReply{}, "", false
}

func ibbParseXMLData(s access.Session, dt []byte) (tag data.IBBData, ctx *recvContext, ictx *ibbContext, ret interface{}, iqtype string, ignore bool) {
	if err := xml.NewDecoder(bytes.NewBuffer(dt)).Decode(&tag); err != nil {
		s.Log().WithError(err).Warn("Failed to parse IBB data")
		return tag, nil, nil, iqErrorNotAcceptable, "error", false
	}

	ctx, ok := getInflightRecv(tag.Sid)

	if !ok || ctx.opaque == nil {
		if hasAndRemoveInflightMAC(tag.Sid) {
			// This is a MAC key reveal sent to us, so we will ignore it.
			return tag, nil, nil, data.EmptyReply{}, "", false
		}

		s.Log().WithField("SID", tag.Sid).Warn("No file transfer associated with SID")
		return tag, nil, nil, iqErrorItemNotFound, "error", false
	}

	ictx, ok = ctx.opaque.(*ibbContext)
	if !ok {
		s.Log().WithField("SID", tag.Sid).Warn("No file transfer associated with SID")
		return tag, nil, nil, iqErrorItemNotFound, "error", false
	}

	return tag, ctx, ictx, nil, "", false
}

func ibbOnData(s access.Session, body []byte) (ret interface{}, iqtype string, ignore bool) {
	tag, ctx, ictx, ret, iqtype, ignore := ibbParseXMLData(s, body)
	if ret != nil {
		return ret, iqtype, ignore
	}

	// XEP-0047 wants us to keep track of previously used sequence numbers, and only do this error
	// when a sequence number is reused - otherwise we should immediately close the stream.
	// However, because of the wraparound behavior of "seq" - also specified in XEP-0047, for large
	// files we can't actually tell the difference between a reused sequence number or a number that
	// has just been wrapped around. Thus, we do this deviation from the spec here.
	if tag.Sequence != ictx.expectingSequence {
		s.Log().WithFields(log.Fields{
			"expected": ictx.expectingSequence,
			"current":  tag.Sequence,
		}).Warn("IBB unexpected sequence")
		ctx.control.ReportError(errors.New("Unexpected data sent from the peer"))
		ctx.ibbCleanup()
		return iqErrorUnexpectedRequest, "error", false

	}

	ictx.expectingSequence++ // wraparound on purpose, to match uint16 spec behavior of the seq field

	dt, err := base64.StdEncoding.DecodeString(tag.Base64)
	if err != nil {
		s.Log().WithError(err).Warn("IBB had an error when decoding")
		ctx.control.ReportError(errors.New("Couldn't decode incoming data"))
		ctx.ibbCleanup()
		return iqErrorNotAcceptable, "error", false
	}

	_, err = ictx.recv.Write(dt)
	if err != nil {
		s.Log().WithError(err).Warn("IBB had an error when writing")
		ctx.control.ReportError(errors.New("Couldn't write incoming data"))
		ctx.ibbCleanup()
		return iqErrorNotAcceptable, "error", false
	}

	return data.EmptyReply{}, "", false
}

// IbbData is the hook function that will be called when we receive an ibb data IQ
func IbbData(s access.Session, stanza *data.ClientIQ) (interface{}, string, bool) {
	return ibbOnData(s, stanza.Query)
}

// IbbMessageData is the hook function that will be called when we receive a message containing an ibb data
func IbbMessageData(s access.Session, stanza *data.ClientMessage, ext *data.Extension) {
	_, _, _ = ibbOnData(s, []byte(ext.Body))
}

// IbbClose is the hook function that will be called when we receive an ibb close IQ
func IbbClose(s access.Session, stanza *data.ClientIQ) (ret interface{}, iqtype string, ignore bool) {
	var tag data.IBBClose
	if err := xml.NewDecoder(bytes.NewBuffer(stanza.Query)).Decode(&tag); err != nil {
		s.Log().WithError(err).Warn("Failed to parse IBB close")
		return iqErrorNotAcceptable, "error", false
	}

	inflightSend, ok := getInflightSend(tag.Sid)
	if ok {
		ibbReceivedClose(inflightSend)
		return data.EmptyReply{}, "", false
	}

	ctx, ok := getInflightRecv(tag.Sid)

	if !ok || ctx.opaque == nil {
		s.Log().WithField("SID", tag.Sid).Warn("No file transfer associated with SID")
		return iqErrorItemNotFound, "error", false
	}

	ictx, ok := ctx.opaque.(*ibbContext)
	if !ok {
		s.Log().WithField("SID", tag.Sid).Warn("No IBB file transfer associated with SID")
		return iqErrorItemNotFound, "error", false
	}

	toSend, fname, ok, ee := ictx.recv.wait()
	if !ok {
		s.Log().WithError(ee).Warn("Had error when waiting for receiving")
		ctx.control.ReportError(errors.New("Couldn't recv final data"))
		ctx.ibbCleanup()
		return
	}

	if toSend != nil {
		encoded := base64.StdEncoding.EncodeToString(toSend)
		_, _, _ = ctx.s.Conn().SendIQ(ctx.peer, "set", data.IBBData{
			Sid:      ctx.sid,
			Sequence: 0,
			Base64:   encoded,
		})
	}

	if err := ctx.finalizeFileTransfer(fname); err != nil {
		s.Log().WithError(err).Warn("Had error when trying to move the final file")
		ctx.control.ReportError(errors.New("Couldn't move the final file"))
		ctx.ibbCleanup()
		_ = os.Remove(fname)
		return
	}

	return data.EmptyReply{}, "", false
}
