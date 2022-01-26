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
	log "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
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
		destination: filepath.Join(destDir, "simple_receipt_test_file_tmp6_"),
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

	content, _ := ioutil.ReadFile(filepath.Join(destDir, "simple_receipt_test_file_tmp6_"))
	c.Assert(content, DeepEquals, []byte{0x42, 0x44, 0x46, 0x43, 0x45})
}

func (s *BytestreamsRecvSuite) Test_bytestreamWaitForCancel_removesInflightRecvOnCancel(c *C) {
	ctl := sdata.CreateFileTransferControl(nil, nil)
	ctx := &recvContext{
		control: ctl,
		sid:     "mytestsid001-42",
	}
	addInflightRecv(ctx)

	done := make(chan bool)
	go func() {
		bytestreamWaitForCancel(ctx)
		done <- true
	}()

	ctl.Cancel()

	<-done

	c.Assert(inflightRecvs.transfers["mytestsid001-42"], IsNil)
}

func (s *BytestreamsRecvSuite) Test_bytestreamWaitForCancel_sendsAMessageWhenCancelling(c *C) {
	ctl := sdata.CreateFileTransferControl(nil, nil)
	msg := make(chan bool, 1)
	ctx := &recvContext{
		control: ctl,
		sid:     "mytestsid001-42",
		opaque:  msg,
	}
	addInflightRecv(ctx)

	done := make(chan bool)
	go func() {
		bytestreamWaitForCancel(ctx)
		done <- true
	}()

	ctl.Cancel()

	<-done
	c.Assert(<-msg, Equals, true)
}

func (s *BytestreamsRecvSuite) Test_bytestreamInitialSetup_failsWhenDecodingQuery(c *C) {
	l, hook := test.NewNullLogger()
	stanza := &data.ClientIQ{
		Query: []byte(`<something bad`),
	}
	sess := new(mock.MockedSession)
	sess.On("Log").Return(l)
	sess.On("SendIQError", mck.Anything, iqErrorIBBBadRequest).Return()

	tag, ctx, early := bytestreamInitialSetup(sess, stanza)

	c.Assert(tag, DeepEquals, data.BytestreamQuery{})
	c.Assert(ctx, IsNil)
	c.Assert(early, Equals, true)

	c.Assert(hook.Entries, HasLen, 1)
	c.Assert(hook.Entries[0].Level, Equals, log.WarnLevel)
	c.Assert(hook.Entries[0].Message, Equals, "Failed to parse bytestream open")
	c.Assert(hook.Entries[0].Data, HasLen, 1)
	c.Assert(hook.Entries[0].Data["error"], ErrorMatches, "XML syntax error on line 1: unexpected EOF")
}

func (s *BytestreamsRecvSuite) Test_bytestreamInitialSetup_failsWhenNotFindingAnInflightReceive(c *C) {
	l, hook := test.NewNullLogger()
	stanza := &data.ClientIQ{
		Query: []byte(`<query xmlns="http://jabber.org/protocol/bytestreams" sid="testSID42" dstaddr="" mode="" activate="">
<streamhost jid="foo@bar.com" host="1.2.3.4"/>
<streamhost jid="humma@foo.com" host="15.111.61.11" port="1234"/>
</query>
`),
	}
	sess := new(mock.MockedSession)
	sess.On("Log").Return(l)
	sess.On("SendIQError", mck.Anything, iqErrorNotAcceptable).Return()

	tag, ctx, early := bytestreamInitialSetup(sess, stanza)

	c.Assert(tag, Not(DeepEquals), data.BytestreamQuery{})
	c.Assert(ctx, IsNil)
	c.Assert(early, Equals, true)

	c.Assert(hook.Entries, HasLen, 1)
	c.Assert(hook.Entries[0].Level, Equals, log.WarnLevel)
	c.Assert(hook.Entries[0].Message, Equals, "No file transfer associated with SID")
	c.Assert(hook.Entries[0].Data, HasLen, 1)
	c.Assert(hook.Entries[0].Data["SID"], Equals, "testSID42")
}

func (s *BytestreamsRecvSuite) Test_bytestreamInitialSetup_failsWhenAStartHasAlreadyHappened(c *C) {
	l, hook := test.NewNullLogger()
	ctx1 := &recvContext{
		sid:    "testSID42",
		opaque: make(chan bool),
	}
	addInflightRecv(ctx1)
	defer func() {
		delete(inflightRecvs.transfers, ctx1.sid)
	}()
	stanza := &data.ClientIQ{
		Query: []byte(`<query xmlns="http://jabber.org/protocol/bytestreams" sid="testSID42" dstaddr="" mode="" activate="">
<streamhost jid="foo@bar.com" host="1.2.3.4"/>
<streamhost jid="humma@foo.com" host="15.111.61.11" port="1234"/>
</query>
`),
	}
	sess := new(mock.MockedSession)
	sess.On("Log").Return(l)
	sess.On("SendIQError", mck.Anything, iqErrorNotAcceptable).Return()

	tag, ctx, early := bytestreamInitialSetup(sess, stanza)

	c.Assert(tag, Not(DeepEquals), data.BytestreamQuery{})
	c.Assert(ctx, Equals, ctx1)
	c.Assert(early, Equals, true)

	c.Assert(hook.Entries, HasLen, 1)
	c.Assert(hook.Entries[0].Level, Equals, log.WarnLevel)
	c.Assert(hook.Entries[0].Message, Equals, "No file transfer associated with SID")
	c.Assert(hook.Entries[0].Data, HasLen, 1)
	c.Assert(hook.Entries[0].Data["SID"], Equals, "testSID42")
}

func (s *BytestreamsRecvSuite) Test_bytestreamInitialSetup_failsWhenTheTagRequiresUDP(c *C) {
	l, hook := test.NewNullLogger()
	ctx1 := &recvContext{
		sid: "testSID43",
	}
	addInflightRecv(ctx1)
	defer func() {
		delete(inflightRecvs.transfers, ctx1.sid)
	}()
	stanza := &data.ClientIQ{
		Query: []byte(`<query xmlns="http://jabber.org/protocol/bytestreams" sid="testSID43" dstaddr="" mode="udp" activate="">
<streamhost jid="foo@bar.com" host="1.2.3.4"/>
<streamhost jid="humma@foo.com" host="15.111.61.11" port="1234"/>
</query>
`),
	}
	sess := new(mock.MockedSession)
	sess.On("Log").Return(l)
	sess.On("SendIQError", mck.Anything, iqErrorIBBBadRequest).Return()

	tag, ctx, early := bytestreamInitialSetup(sess, stanza)

	c.Assert(tag, Not(DeepEquals), data.BytestreamQuery{})
	c.Assert(ctx, Equals, ctx1)
	c.Assert(early, Equals, true)

	c.Assert(hook.Entries, HasLen, 1)
	c.Assert(hook.Entries[0].Level, Equals, log.WarnLevel)
	c.Assert(hook.Entries[0].Message, Equals, "Received a request for UDP, even though we don't support or advertize UDP - this means the peer is using a non-conforming application")
	c.Assert(hook.Entries[0].Data, HasLen, 0)
}
