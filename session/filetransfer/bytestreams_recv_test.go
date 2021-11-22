package filetransfer

import (
	"io"
	"io/ioutil"
	"path/filepath"

	"github.com/coyim/coyim/config"
	sdata "github.com/coyim/coyim/session/data"
	"github.com/coyim/coyim/session/mock"
	"github.com/coyim/coyim/xmpp/data"
	"github.com/prashantv/gostub"
	mck "github.com/stretchr/testify/mock"
	. "gopkg.in/check.v1"
)

type BytestreamsRecvSuite struct{}

var _ = Suite(&BytestreamsRecvSuite{})

func (s *BytestreamsRecvSuite) Test_BytestreamQuery_works(c *C) {
	destDir := c.MkDir()

	defer gostub.StubFunc(&createTorProxy, nil, nil).Reset()

	md := new(mockedDialer)
	defer gostub.StubFunc(&socks5XMPP, md, nil).Reset()

	mc := new(mockedConn)
	md.On("Dial", "tcp", "3923a798497112d7e2bcb1f40ce4bb0f664b2841:0").Return(mc, nil)

	ctl := sdata.CreateFileTransferControl(nil, nil)
	ctx := &recvContext{
		sid:         "testSID42",
		destination: filepath.Join(destDir, "simple_receipt_test_file"),
		control:     ctl,
		size:        5,
	}

	addInflightRecv(ctx)

	stanza := &data.ClientIQ{
		Query: []byte(`<query xmlns="http://jabber.org/protocol/bytestreams" sid="testSID42" dstaddr="" mode="" activate="">
<streamhost jid="foo@bar.com" host="1.2.3.4"/>
<streamhost jid="humma@foo.com" host="15.111.61.11" port="1234"/>
</query>
`),
	}

	conn := new(mock.MockedSession)
	conn.On("GetConfig").Return(&config.Account{})

	mc.On("Read", mck.Anything).Return(5, nil).Run(func(a mck.Arguments) {
		buf := a.Get(0).([]byte)
		copy(buf, []byte{0x42, 0x44, 0x46, 0x43, 0x45})
	}).Once().
		On("Read", mck.Anything).Return(0, io.EOF).Run(func(a mck.Arguments) {
	}).Once().
		On("Close").Return(nil)

	done := make(chan bool)

	go ctl.WaitForFinish(func(v bool) {
		done <- v
	})

	ret, iqt, ignore := BytestreamQuery(conn, stanza)

	q, ok := ret.(data.BytestreamQuery)
	c.Assert(ok, Equals, true)
	c.Assert(q.Sid, Equals, "testSID42")
	c.Assert(q.StreamhostUsed, Not(IsNil))
	c.Assert(q.StreamhostUsed.Jid, Equals, "foo@bar.com")
	c.Assert(iqt, Equals, "result")
	c.Assert(ignore, Equals, false)

	c.Assert(<-done, Equals, true)

	content, _ := ioutil.ReadFile(filepath.Join(destDir, "simple_receipt_test_file"))
	c.Assert(content, DeepEquals, []byte{0x42, 0x44, 0x46, 0x43, 0x45})
}
