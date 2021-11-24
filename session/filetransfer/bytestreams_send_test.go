package filetransfer

import (
	"encoding/xml"
	"fmt"

	"github.com/coyim/coyim/cache"
	"github.com/coyim/coyim/config"
	sdata "github.com/coyim/coyim/session/data"
	"github.com/coyim/coyim/session/mock"
	"github.com/coyim/coyim/xmpp/data"
	xmock "github.com/coyim/coyim/xmpp/mock"
	"github.com/prashantv/gostub"
	mck "github.com/stretchr/testify/mock"
	. "gopkg.in/check.v1"
)

type BytestreamsSendSuite struct {
	WithTempFileSuite
}

var _ = Suite(&BytestreamsSendSuite{})

func (s *BytestreamsSendSuite) Test_bytestreamsSendDo_works(c *C) {
	defer gostub.StubFunc(&createTorProxy, nil, nil).Reset()

	md := new(mockedDialer)
	defer gostub.StubFunc(&socks5XMPP, md, nil).Reset()

	smc := new(mockedConn)
	md.On("Dial", "tcp", "65e77469d4570a364de45c329f8a5e65ff5620b5:0").Return(smc, nil)

	smc.On("Write", []byte{0x73, 0x6f, 0x6d, 0x65, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x20, 0x6e, 0x65, 0x77}).Return(13, nil).Once()
	smc.On("Close").Return(nil).Once()

	sess := new(mock.MockedSession)
	mc := new(xmock.MockedConn)

	sess.On("Conn").Return(mc)
	sess.On("GetConfig").Return(&config.Account{})

	cac := cache.New()
	cac.Put("http://jabber.org/protocol/bytestreams . proxies", []*data.BytestreamStreamhost{
		{Jid: "foo@bar.com", Host: "1.2.3.4"},
		{Jid: "humma@foo.com", Host: "15.111.61.11", Port: 1234},
	})
	mc.On("Cache").Return(cac)

	ret1 := make(chan data.Stanza)
	var ret1read <-chan data.Stanza = ret1
	mc.On("SendIQ", "", "set", mck.Anything).Return(ret1read, data.Cookie(42), nil).Once()

	ret2 := make(chan data.Stanza)
	var ret2read <-chan data.Stanza = ret2
	mc.On("SendIQ", "foo@bar.com", "set", mck.Anything).Return(ret2read, data.Cookie(43), nil).Once()

	ctx := &sendContext{
		s:       sess,
		sid:     "something42",
		file:    s.file,
		control: sdata.CreateFileTransferControl(nil, nil),
	}

	done := make(chan bool)

	go ctx.control.WaitForFinish(func(v bool) {
		done <- v
	})

	go ctx.control.WaitForError(func(e error) {
		fmt.Printf("error: %v\n", e)
		done <- false
	})

	bytestreamsSendDo(ctx)

	ret1 <- data.Stanza{
		Value: &data.ClientIQ{
			From: "whoops@com.foo",
			Type: "result",
			Query: []byte(`
<query xmlns="http://jabber.org/protocol/bytestreams">
  <streamhost-used jid="foo@bar.com"/>
</query>`),
		},
	}

	ret2 <- data.Stanza{
		Value: &data.ClientIQ{
			From: "foo@bar.com",
			Type: "result",
		},
	}

	c.Assert(<-done, Equals, true)
}

func (s *BytestreamsSendSuite) Test_bytestreamsCalculateValidProxies_works(c *C) {
	sess := new(mock.MockedSession)
	mc := new(xmock.MockedConn)

	sess.On("GetConfig").Return(&config.Account{Account: "hello@bar.com"})
	sess.On("Conn").Return(mc)

	ret1 := make(chan data.Stanza)
	var ret1read <-chan data.Stanza = ret1
	mc.On("SendIQ", "bar.com", "get", mck.AnythingOfType("*data.DiscoveryItemsQuery")).Return(ret1read, data.Cookie(44), nil).Once()

	ret2 := make(chan data.Stanza)
	var ret2read <-chan data.Stanza = ret2
	mc.On("SendIQ", "anythingelse@stuff.com", "get", mck.AnythingOfType("*data.BytestreamQuery")).Return(ret2read, data.Cookie(45), nil).Once()

	ret3 := make(chan data.Stanza)
	var ret3read <-chan data.Stanza = ret3
	mc.On("SendIQ", "testbar.something.com", "get", mck.AnythingOfType("*data.BytestreamQuery")).Return(ret3read, data.Cookie(46), nil).Once()

	mc.On("DiscoveryFeaturesAndIdentities", "testbar.something.com").Return(
		[]data.DiscoveryIdentity{
			{Category: "stuff", Type: "something else", Name: "just for testing"},
			{Category: "proxy", Type: "bytestreams"}},
		[]string{"something", BytestreamMethod, "something else"},
		false,
	).Once()

	mc.On("DiscoveryFeaturesAndIdentities", "anythingelse@stuff.com").Return(
		[]data.DiscoveryIdentity{
			{Category: "proxy", Type: "bytestreams"},
			{Category: "stuff", Type: "something else", Name: "just for testing"},
		},
		[]string{BytestreamMethod, "something else"},
		false,
	).Once()

	mc.On("DiscoveryFeaturesAndIdentities", "final.kefka.com").Return(
		[]data.DiscoveryIdentity{
			{Category: "stuff", Type: "something else", Name: "just for testing"},
		},
		[]string{"hmm", "something else"},
		false,
	).Once()

	mc.On("DiscoveryFeaturesAndIdentities", "otherproxy.com").Return(
		[]data.DiscoveryIdentity{
			{Category: "stuff", Type: "something else", Name: "just for testing"},
			{Category: "proxy", Type: "ibb"}},
		[]string{"something", BytestreamMethod, "something else"},
		false,
	).Once()

	mc.On("DiscoveryFeaturesAndIdentities", "proxynobytestream.com").Return(
		[]data.DiscoveryIdentity{
			{Category: "stuff", Type: "something else", Name: "just for testing"},
			{Category: "proxy", Type: "bytestreams"}},
		[]string{"something", "something else"},
		false,
	).Once()

	mc.On("DiscoveryFeaturesAndIdentities", "bytestreammethodbutnotproxy.com.com").Return(
		[]data.DiscoveryIdentity{
			{Category: "stuff", Type: "something else", Name: "just for testing"},
		},
		[]string{"something", BytestreamMethod, "something else"},
		false,
	).Once()

	done := make(chan interface{})

	go func() {
		output := bytestreamsCalculateValidProxies(sess)("this argument isn't used and doesn't matter")
		done <- output
	}()

	ret1 <- data.Stanza{
		Value: &data.ClientIQ{
			From: "bar.com",
			Type: "result",
			Query: []byte(`
<query xmlns="http://jabber.org/protocol/disco#items">
  <node/>
  <item jid="testbar.something.com" name="Hello world"/>
  <item jid="anythingelse@stuff.com" name="Is this the droids you're looking for?"/>
  <item jid="final.kefka.com" name="Nope"/>
  <item jid="otherproxy.com" name="Another proxy type"/>
  <item jid="proxynobytestream.com" name="A proxy service with no bytestream method"/>
  <item jid="bytestreammethodbutnotproxy.com.com" name="Something that has a bytestream method but no proxy"/>
</query>`),
		},
	}

	ret2 <- data.Stanza{
		Value: &data.ClientIQ{
			From: "anythingelse@stuff.com",
			Type: "result",
			Query: []byte(`
<query xmlns="http://jabber.org/protocol/bytestreams">
  <streamhost jid="one@foo.bar" host="1.2.3.4" port="1234"/>
  <streamhost jid="two@foo.bar" host="somewhere.com"/>
  <streamhost jid="three@somewhere.com" host="4.3.2.1" port="5555"/>
</query>`),
		},
	}

	ret3 <- data.Stanza{
		Value: &data.ClientIQ{
			From: "testbar.something.com",
			Type: "result",
			Query: []byte(`
<query xmlns="http://jabber.org/protocol/bytestreams">
  <streamhost jid="five@foo.bar" host="4.4.4.4" port="1111"/>
</query>`),
		},
	}

	result := (<-done).([]*data.BytestreamStreamhost)
	c.Assert(result, HasLen, 2)
	c.Assert(result[0], DeepEquals, &data.BytestreamStreamhost{
		XMLName: xml.Name{Space: "http://jabber.org/protocol/bytestreams", Local: "streamhost"},
		Jid:     "five@foo.bar",
		Host:    "4.4.4.4",
		Port:    1111,
	})
	c.Assert(result[1], DeepEquals, &data.BytestreamStreamhost{
		XMLName: xml.Name{Space: "http://jabber.org/protocol/bytestreams", Local: "streamhost"},
		Jid:     "one@foo.bar",
		Host:    "1.2.3.4",
		Port:    1234,
	})
}
