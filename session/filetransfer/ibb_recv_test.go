package filetransfer

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/coyim/coyim/coylog"
	sdata "github.com/coyim/coyim/session/data"
	"github.com/coyim/coyim/xmpp/data"
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

	go ctx.control.WaitForError(func(e error) {
		fmt.Printf("Had unexpected error: %v\n", e)
		finished <- false
	})

	ibbctx.recv = ctx.createReceiver()

	addInflightRecv(ctx)

	_, ee := ibbctx.recv.Write([]byte("hello world"))
	c.Assert(ee, IsNil)

	ret, iqtype, ignore := IbbClose(wl, stanza)

	c.Assert(<-finished, Equals, true)

	c.Assert(ret, DeepEquals, data.EmptyReply{})
	c.Assert(iqtype, Equals, "")
	c.Assert(ignore, Equals, false)

	c.Assert(len(hook.Entries), Equals, 0)

	content, _ := ioutil.ReadFile(ctx.destination)
	c.Assert(string(content), DeepEquals, "hello world")
}
