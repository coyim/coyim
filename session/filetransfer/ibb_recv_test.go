package filetransfer

import (
	"encoding/xml"
	"io/ioutil"
	"path/filepath"
	"sync"

	"github.com/coyim/coyim/coylog"
	sdata "github.com/coyim/coyim/session/data"
	"github.com/coyim/coyim/session/mock"
	"github.com/coyim/coyim/xmpp/data"
	xmock "github.com/coyim/coyim/xmpp/mock"
	log "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	. "gopkg.in/check.v1"
)

type IBBReceiverSuite struct{}

var _ = Suite(&IBBReceiverSuite{})

type mockHasLog struct {
	l coylog.Logger
}

func (m *mockHasLog) Log() coylog.Logger {
	return m.l
}

func (s *IBBReceiverSuite) Test_IbbOpen_works(c *C) {
	destDir := c.MkDir()
	l, hook := test.NewNullLogger()

	wl := &mockHasLog{l}

	stanza := &data.ClientIQ{
		Query: []byte(`<open xmlns="http://jabber.org/protocol/ibb" sid="testSID" block-size="4096"/>`),
	}

	ctx := &recvContext{
		sid:         "testSID",
		destination: filepath.Join(destDir, "simple_receipt_test_file_tmp1_"),
	}

	addInflightRecv(ctx)

	ret, iqtype, ignore := IbbOpen(wl, stanza)

	c.Assert(ret, DeepEquals, data.EmptyReply{})
	c.Assert(iqtype, Equals, "")
	c.Assert(ignore, Equals, false)

	c.Assert(ctx.opaque, Not(IsNil))
	ibbctx := ctx.opaque.(*ibbContext)
	c.Assert(ibbctx.recv, Not(IsNil))
	ibbctx.recv.cleanupAfterRun()

	c.Assert(len(hook.Entries), Equals, 0)
}

func (s *IBBReceiverSuite) Test_IbbData_works(c *C) {
	destDir := c.MkDir()
	l, hook := test.NewNullLogger()
	sess := &sessionMockWithCustomLog{
		log: l,
	}

	wl := &mockHasLog{l}

	stanza := &data.ClientIQ{
		Query: []byte(
			`
<data xmlns='http://jabber.org/protocol/ibb' sid='testSID' seq='0'>
aGVsbG8gd29ybGQuIHRoaXMgaXMgYSB0ZXN0IG9mIGZpbGUgZGVjcnlwdGlvbiBz
dHVmZiwgc28gdGhlIGNvbnRlbnQgZG9lc24ndCBtYXR0ZXIgc28gbXVjaC4=
</data>

`),
	}

	ibbctx := &ibbContext{}

	ctx := &recvContext{
		s:           sess,
		sid:         "testSID",
		opaque:      ibbctx,
		size:        92,
		control:     sdata.CreateFileTransferControl(nil, nil),
		destination: filepath.Join(destDir, "simple_receipt_test_file_tmp2_"),
	}

	ibbctx.recv = ctx.createReceiver()

	addInflightRecv(ctx)

	ret, iqtype, ignore := IbbData(wl, stanza)

	c.Assert(ret, DeepEquals, data.EmptyReply{})
	c.Assert(iqtype, Equals, "")
	c.Assert(ignore, Equals, false)

	c.Assert(len(hook.Entries), Equals, 0)
}

func (s *IBBReceiverSuite) Test_IbbClose_works(c *C) {
	destDir := c.MkDir()
	l, hook := test.NewNullLogger()
	sess := &sessionMockWithCustomLog{
		log: l,
	}

	wl := &mockHasLog{l}

	stanza := &data.ClientIQ{
		Query: []byte(`<close xmlns='http://jabber.org/protocol/ibb' sid='testSID'/>`),
	}

	ibbctx := &ibbContext{}

	ctx := &recvContext{
		s:           sess,
		sid:         "testSID",
		opaque:      ibbctx,
		size:        11,
		control:     sdata.CreateFileTransferControl(nil, nil),
		destination: filepath.Join(destDir, "simple_receipt_test_file_tmp3_"),
	}

	finished := make(chan bool)
	go ctx.control.WaitForFinish(func(ok bool) {
		finished <- ok
	})

	var unexpectedError error
	go ctx.control.WaitForError(func(e error) {
		unexpectedError = e
		finished <- false
	})

	ibbctx.recv = ctx.createReceiver()

	addInflightRecv(ctx)

	_, ee := ibbctx.recv.Write([]byte("hello world"))
	c.Assert(ee, IsNil)

	ret, iqtype, ignore := IbbClose(wl, stanza)

	c.Assert(<-finished, Equals, true)

	c.Assert(unexpectedError, IsNil)
	c.Assert(ret, DeepEquals, data.EmptyReply{})
	c.Assert(iqtype, Equals, "")
	c.Assert(ignore, Equals, false)

	c.Assert(len(hook.Entries), Equals, 0)

	content, _ := ioutil.ReadFile(ctx.destination)
	c.Assert(string(content), DeepEquals, "hello world")
}

func (s *IBBReceiverSuite) Test_ibbParseXMLData_failsOnParsingData(c *C) {
	l, hook := test.NewNullLogger()
	wl := &mockHasLog{l}

	tag, ctx, ictx, ret, iqtype, ignore := ibbParseXMLData(wl, []byte{0x42})
	c.Assert(tag, DeepEquals, data.IBBData{})
	c.Assert(ctx, IsNil)
	c.Assert(ictx, IsNil)
	c.Assert(ret, Equals, iqErrorNotAcceptable)
	c.Assert(iqtype, Equals, "error")
	c.Assert(ignore, Equals, false)

	c.Assert(hook.Entries, HasLen, 1)
	c.Assert(hook.Entries[0].Level, Equals, log.WarnLevel)
	c.Assert(hook.Entries[0].Message, Equals, "Failed to parse IBB data")
	c.Assert(hook.Entries[0].Data, HasLen, 1)
	c.Assert(hook.Entries[0].Data["error"], ErrorMatches, "EOF")
}

func (s *IBBReceiverSuite) Test_ibbParseXMLData_noFileTransferAssociatedWithTag(c *C) {
	l, hook := test.NewNullLogger()
	wl := &mockHasLog{l}

	tag, ctx, ictx, ret, iqtype, ignore := ibbParseXMLData(wl, []byte(`
<data xmlns="http://jabber.org/protocol/ibb" sid="123"/>
`))
	c.Assert(tag, DeepEquals, data.IBBData{
		XMLName: xml.Name{Space: "http://jabber.org/protocol/ibb", Local: "data"},
		Sid:     "123"})
	c.Assert(ctx, IsNil)
	c.Assert(ictx, IsNil)
	c.Assert(ret, Equals, iqErrorItemNotFound)
	c.Assert(iqtype, Equals, "error")
	c.Assert(ignore, Equals, false)

	c.Assert(hook.Entries, HasLen, 1)
	c.Assert(hook.Entries[0].Level, Equals, log.WarnLevel)
	c.Assert(hook.Entries[0].Message, Equals, "No file transfer associated with SID")
	c.Assert(hook.Entries[0].Data, HasLen, 1)
	c.Assert(hook.Entries[0].Data["SID"], Equals, "123")
}

func (s *IBBReceiverSuite) Test_ibbParseXMLData_waitingForMAC(c *C) {
	addInflightMAC(&sendContext{sid: "123"})

	l, hook := test.NewNullLogger()
	wl := &mockHasLog{l}

	tag, ctx, ictx, ret, iqtype, ignore := ibbParseXMLData(wl, []byte(`
<data xmlns="http://jabber.org/protocol/ibb" sid="123"/>
`))
	c.Assert(tag, DeepEquals, data.IBBData{
		XMLName: xml.Name{Space: "http://jabber.org/protocol/ibb", Local: "data"},
		Sid:     "123"})
	c.Assert(ctx, IsNil)
	c.Assert(ictx, IsNil)
	c.Assert(ret, Equals, data.EmptyReply{})
	c.Assert(iqtype, Equals, "")
	c.Assert(ignore, Equals, false)

	c.Assert(hook.Entries, HasLen, 0)

	c.Assert(inflightMACs.transfers["123"], Equals, false)
}

func (s *IBBReceiverSuite) Test_ibbParseXMLData_incorrectIBBContext(c *C) {
	l, hook := test.NewNullLogger()
	wl := &mockHasLog{l}

	rctx := &recvContext{sid: "123"}
	rctx.opaque = "hello"
	addInflightRecv(rctx)
	defer func() {
		removeInflightRecv("123")
	}()

	tag, ctx, ictx, ret, iqtype, ignore := ibbParseXMLData(wl, []byte(`
<data xmlns="http://jabber.org/protocol/ibb" sid="123"/>
`))
	c.Assert(tag, DeepEquals, data.IBBData{
		XMLName: xml.Name{Space: "http://jabber.org/protocol/ibb", Local: "data"},
		Sid:     "123"})
	c.Assert(ctx, IsNil)
	c.Assert(ictx, IsNil)
	c.Assert(ret, Equals, iqErrorItemNotFound)
	c.Assert(iqtype, Equals, "error")
	c.Assert(ignore, Equals, false)

	c.Assert(hook.Entries, HasLen, 1)
	c.Assert(hook.Entries[0].Level, Equals, log.WarnLevel)
	c.Assert(hook.Entries[0].Message, Equals, "No file transfer associated with SID")
	c.Assert(hook.Entries[0].Data, HasLen, 1)
	c.Assert(hook.Entries[0].Data["SID"], Equals, "123")
}

func (s *IBBReceiverSuite) Test_recvContext_ibbCleanup_removesTheInflight(c *C) {
	ctx := &recvContext{sid: "hello123"}
	addInflightRecv(ctx)

	ctx.ibbCleanup()

	c.Assert(inflightRecvs.transfers["hello123"], IsNil)
}

func (s *IBBReceiverSuite) Test_ibbWaitForCancel_waitsForCancel(c *C) {
	sess := new(mock.MockedSession)
	mc := new(xmock.MockedConn)

	sess.On("Conn").Return(mc)
	mc.On("SendIQ", "", "set", data.IBBClose{Sid: "hello444"}).Return(make(<-chan data.Stanza), data.Cookie(0), nil).Once()

	rcv := &receiver{}
	rcv.newData = sync.NewCond(rcv)

	ibbctx := &ibbContext{
		recv: rcv,
	}
	ctx := &recvContext{
		s:       sess,
		sid:     "hello444",
		opaque:  ibbctx,
		control: sdata.CreateFileTransferControl(nil, nil),
	}
	addInflightRecv(ctx)

	done := make(chan bool)
	go func() {
		ibbWaitForCancel(ctx)
		done <- true
	}()

	ctx.control.Cancel()

	<-done

	c.Assert(rcv.hadError, Equals, true)
	c.Assert(rcv.err, ErrorMatches, "local cancel")
}

func (s *IBBReceiverSuite) Test_IbbOpen_failsIfUnableToParseQuery(c *C) {
	l, hook := test.NewNullLogger()
	sess := &sessionMockWithCustomLog{
		log: l,
	}

	ret, iqtype, ignore := IbbOpen(sess, &data.ClientIQ{
		Query: []byte("hello world <"),
	})

	c.Assert(ret, Equals, iqErrorNotAcceptable)
	c.Assert(iqtype, Equals, "error")
	c.Assert(ignore, Equals, false)

	c.Assert(hook.Entries, HasLen, 1)
	c.Assert(hook.Entries[0].Level, Equals, log.WarnLevel)
	c.Assert(hook.Entries[0].Message, Equals, "Failed to parse IBB open")
	c.Assert(hook.Entries[0].Data, HasLen, 1)
	c.Assert(hook.Entries[0].Data["error"], ErrorMatches, "XML syntax error on line 1: unexpected EOF")
}

func (s *IBBReceiverSuite) Test_IbbOpen_failsIfReceiverAlreadyHaveAnIBBContext(c *C) {
	l, hook := test.NewNullLogger()
	sess := &sessionMockWithCustomLog{
		log: l,
	}

	ctx := &recvContext{
		sid:    "testSID42",
		opaque: "already",
	}
	addInflightRecv(ctx)

	ret, iqtype, ignore := IbbOpen(sess, &data.ClientIQ{
		Query: []byte(`<open xmlns="http://jabber.org/protocol/ibb" sid="testSID42" block-size="4096"/>`),
	})

	c.Assert(ret, Equals, iqErrorNotAcceptable)
	c.Assert(iqtype, Equals, "error")
	c.Assert(ignore, Equals, false)

	c.Assert(hook.Entries, HasLen, 1)
	c.Assert(hook.Entries[0].Level, Equals, log.WarnLevel)
	c.Assert(hook.Entries[0].Message, Equals, "No file transfer associated with SID")
	c.Assert(hook.Entries[0].Data, HasLen, 1)
	c.Assert(hook.Entries[0].Data["SID"], Equals, "testSID42")
}

func (s *IBBReceiverSuite) Test_ibbOnData_failsOnBadBodyData(c *C) {
	l, hook := test.NewNullLogger()
	wl := &mockHasLog{l}

	ret, iqtype, ignore := ibbOnData(wl, []byte{0x42})
	c.Assert(ret, Equals, iqErrorNotAcceptable)
	c.Assert(iqtype, Equals, "error")
	c.Assert(ignore, Equals, false)

	c.Assert(hook.Entries, HasLen, 1)
	c.Assert(hook.Entries[0].Level, Equals, log.WarnLevel)
	c.Assert(hook.Entries[0].Message, Equals, "Failed to parse IBB data")
	c.Assert(hook.Entries[0].Data, HasLen, 1)
	c.Assert(hook.Entries[0].Data["error"], ErrorMatches, "EOF")
}

func (s *IBBReceiverSuite) Test_ibbOnData_failsOnIncorrectSequenceNumber(c *C) {
	destDir := c.MkDir()
	l, hook := test.NewNullLogger()
	sess := &sessionMockWithCustomLog{
		log: l,
	}

	wl := &mockHasLog{l}

	ibbctx := &ibbContext{
		expectingSequence: 35,
	}

	ctx := &recvContext{
		s:           sess,
		sid:         "testSID",
		opaque:      ibbctx,
		size:        92,
		control:     sdata.CreateFileTransferControl(nil, nil),
		destination: filepath.Join(destDir, "simple_receipt_test_file_tmp4_"),
	}

	ee := make(chan error)
	go ctx.control.WaitForError(func(e error) {
		ee <- e
	})

	ibbctx.recv = ctx.createReceiver()

	addInflightRecv(ctx)

	ret, iqtype, ignore := ibbOnData(wl, []byte(`
<data xmlns='http://jabber.org/protocol/ibb' sid='testSID' seq='42'>
aGVsbG8gd29ybGQuIHRoaXMgaXMgYSB0ZXN0IG9mIGZpbGUgZGVjcnlwdGlvbiBz
dHVmZiwgc28gdGhlIGNvbnRlbnQgZG9lc24ndCBtYXR0ZXIgc28gbXVjaC4=
</data>
`))

	c.Assert(<-ee, ErrorMatches, "Unexpected data sent from the peer")

	ibbctx.recv.cleanupAfterRun()

	c.Assert(ret, DeepEquals, iqErrorUnexpectedRequest)
	c.Assert(iqtype, Equals, "error")
	c.Assert(ignore, Equals, false)

	c.Assert(len(hook.Entries), Equals, 1)
	c.Assert(hook.Entries[0].Level, Equals, log.WarnLevel)
	c.Assert(hook.Entries[0].Message, Equals, "IBB unexpected sequence")
	c.Assert(hook.Entries[0].Data, HasLen, 2)
	c.Assert(hook.Entries[0].Data["expected"], Equals, uint16(35))
	c.Assert(hook.Entries[0].Data["current"], Equals, uint16(42))
}

func (s *IBBReceiverSuite) Test_ibbOnData_failsOnDecodingBase64(c *C) {
	destDir := c.MkDir()
	l, hook := test.NewNullLogger()
	sess := &sessionMockWithCustomLog{
		log: l,
	}

	wl := &mockHasLog{l}

	ibbctx := &ibbContext{}

	ctx := &recvContext{
		s:           sess,
		sid:         "testSID",
		opaque:      ibbctx,
		size:        92,
		control:     sdata.CreateFileTransferControl(nil, nil),
		destination: filepath.Join(destDir, "simple_receipt_test_file_tmp5_"),
	}

	ee := make(chan error)
	go ctx.control.WaitForError(func(e error) {
		ee <- e
	})

	ibbctx.recv = ctx.createReceiver()

	addInflightRecv(ctx)

	ret, iqtype, ignore := ibbOnData(wl, []byte(
		`
<data xmlns='http://jabber.org/protocol/ibb' sid='testSID' seq='0'>
****biBz
dHVmZiwgc28gdGhlIGNvbnRlbnQgZG9lc24ndCBtYXR0ZXIgc28gbXVjaC4=
</data>

`))

	c.Assert(<-ee, ErrorMatches, "Couldn't decode incoming data")

	ibbctx.recv.cleanupAfterRun()

	c.Assert(ret, DeepEquals, iqErrorNotAcceptable)
	c.Assert(iqtype, Equals, "error")
	c.Assert(ignore, Equals, false)

	c.Assert(len(hook.Entries), Equals, 1)
	c.Assert(hook.Entries[0].Level, Equals, log.WarnLevel)
	c.Assert(hook.Entries[0].Message, Equals, "IBB had an error when decoding")
	c.Assert(hook.Entries[0].Data, HasLen, 1)
	c.Assert(hook.Entries[0].Data["error"], ErrorMatches, "illegal base64 data at input byte 1")
}
