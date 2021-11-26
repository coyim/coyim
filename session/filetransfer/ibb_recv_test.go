package filetransfer

import (
	"encoding/xml"
	"io/ioutil"
	"path/filepath"

	"github.com/coyim/coyim/coylog"
	sdata "github.com/coyim/coyim/session/data"
	"github.com/coyim/coyim/xmpp/data"
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
		destination: filepath.Join(destDir, "simple_receipt_test_file"),
	}

	addInflightRecv(ctx)

	ret, iqtype, ignore := IbbOpen(wl, stanza)

	c.Assert(ret, DeepEquals, data.EmptyReply{})
	c.Assert(iqtype, Equals, "")
	c.Assert(ignore, Equals, false)

	c.Assert(ctx.opaque, Not(IsNil))
	ibbctx := ctx.opaque.(*ibbContext)
	c.Assert(ibbctx.recv, Not(IsNil))

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
		destination: filepath.Join(destDir, "simple_receipt_test_file"),
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
		destination: filepath.Join(destDir, "simple_receipt_test_file"),
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
