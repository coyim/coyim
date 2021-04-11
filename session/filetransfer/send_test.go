package filetransfer

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"

	"github.com/coyim/coyim/config"
	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/xmpp/data"
	xi "github.com/coyim/coyim/xmpp/interfaces"
	"github.com/coyim/coyim/xmpp/mock"
	log "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	. "gopkg.in/check.v1"
)

type SendSuite struct{}

var _ = Suite(&SendSuite{})

type mockWithConn struct {
	c xi.Conn
}

func (m *mockWithConn) Conn() xi.Conn {
	return m.c
}

type mockConnDiscoveryFeatures struct {
	*mock.Conn

	f func(string) ([]string, bool)
}

func (m *mockConnDiscoveryFeatures) DiscoveryFeatures(v string) ([]string, bool) {
	return m.f(v)
}

func (s *SendSuite) Test_discoverSupport_fails(c *C) {
	mc := &mockWithConn{&mockConnDiscoveryFeatures{
		f: func(string) ([]string, bool) {
			return nil, false
		},
	}}

	res, err := discoverSupport(mc, "foo@bar.com/foo")
	c.Assert(res, HasLen, 0)
	c.Assert(err, ErrorMatches, "Problem discovering the features of the peer")
}

func (s *SendSuite) Test_discoverSupport_withNoSISupport(c *C) {
	mc := &mockWithConn{&mockConnDiscoveryFeatures{
		f: func(string) ([]string, bool) {
			return []string{}, true
		},
	}}

	res, err := discoverSupport(mc, "foo@bar.com/foo")
	c.Assert(res, HasLen, 0)
	c.Assert(err, ErrorMatches, "Peer doesn't support stream initiation")
}

func (s *SendSuite) Test_discoverSupport_withNoProfiles(c *C) {
	mc := &mockWithConn{&mockConnDiscoveryFeatures{
		f: func(string) ([]string, bool) {
			return []string{
				"foo",
				"http://jabber.org/protocol/si",
			}, true
		},
	}}

	res, err := discoverSupport(mc, "foo@bar.com/foo")
	c.Assert(res, HasLen, 0)
	c.Assert(err, ErrorMatches, "Peer doesn't support any stream initiation profiles")
}

func (s *SendSuite) Test_discoverSupport_withProfiles(c *C) {
	mc := &mockWithConn{&mockConnDiscoveryFeatures{
		f: func(string) ([]string, bool) {
			return []string{
				"foo",
				"http://jabber.org/protocol/si/profile/hubba",
				"http://jabber.org/protocol/si",
				"http://jabber.org/protocol/si/profile/bubba",
			}, true
		},
	}}

	res, err := discoverSupport(mc, "foo@bar.com/foo")
	c.Assert(res, DeepEquals, map[string]bool{
		"http://jabber.org/protocol/si/profile/hubba": true,
		"http://jabber.org/protocol/si/profile/bubba": true,
	})
	c.Assert(err, IsNil)
}

type mockConnRand struct {
	*mock.Conn

	f func() io.Reader
}

func (m *mockConnRand) Rand() io.Reader {
	return m.f()
}

func (s *SendSuite) Test_genSid_succeeds(c *C) {
	m := &mockConnRand{}

	m.f = func() io.Reader {
		return bytes.NewBufferString("abcdefesdfgsdfgsdfg")
	}

	res := genSid(m)
	c.Assert(res, Equals, "sid8315164872671453793")
}

func (s *SendSuite) Test_genSid_panics(c *C) {
	m := &mockConnRand{}

	m.f = func() io.Reader {
		return bytes.NewBufferString("")
	}

	c.Assert(func() { genSid(m) }, PanicMatches, "Failed to read random bytes: EOF")
}

func (s *SendSuite) Test_calculateAvailableSendOptions(c *C) {
	orgSupportedSendingMechanisms := supportedSendingMechanisms
	orgIsSendingMechanismCurrentlyValid := isSendingMechanismCurrentlyValid

	defer func() {
		supportedSendingMechanisms = orgSupportedSendingMechanisms
		isSendingMechanismCurrentlyValid = orgIsSendingMechanismCurrentlyValid
	}()

	supportedSendingMechanisms = map[string]func(*sendContext){}
	isSendingMechanismCurrentlyValid = map[string]func(string, hasConnectionAndConfig) bool{}

	supportedSendingMechanisms["foo"] = nil
	supportedSendingMechanisms["bar"] = nil

	isSendingMechanismCurrentlyValid["foo"] = func(string, hasConnectionAndConfig) bool {
		return true
	}

	isSendingMechanismCurrentlyValid["bar"] = func(string, hasConnectionAndConfig) bool {
		return false
	}

	res := calculateAvailableSendOptions(nil)

	c.Assert(res, DeepEquals, []data.FormFieldOptionX{
		data.FormFieldOptionX{Value: "foo"},
	})
}

type mockHasConnectionAndConfigAndLog struct {
	c    xi.Conn
	l    coylog.Logger
	conf *config.Account
}

func (m *mockHasConnectionAndConfigAndLog) Log() coylog.Logger {
	return m.l
}

func (m *mockHasConnectionAndConfigAndLog) GetConfig() *config.Account {
	return m.conf
}

func (m *mockHasConnectionAndConfigAndLog) Conn() xi.Conn {
	return m.c
}

type sendIQAndRandMock struct {
	*sendIQMock

	rand func() io.Reader
}

func (m *sendIQAndRandMock) Rand() io.Reader {
	return m.rand()
}

func (s *SendSuite) Test_sendContext_offerSend_works(c *C) {
	orgSupportedSendingMechanisms := supportedSendingMechanisms
	orgIsSendingMechanismCurrentlyValid := isSendingMechanismCurrentlyValid

	defer func() {
		supportedSendingMechanisms = orgSupportedSendingMechanisms
		isSendingMechanismCurrentlyValid = orgIsSendingMechanismCurrentlyValid
	}()

	supportedSendingMechanisms = map[string]func(*sendContext){}
	isSendingMechanismCurrentlyValid = map[string]func(string, hasConnectionAndConfig) bool{}

	done := make(chan bool)
	supportedSendingMechanisms["foo"] = func(*sendContext) {
		done <- true
	}

	isSendingMechanismCurrentlyValid["foo"] = func(string, hasConnectionAndConfig) bool {
		return true
	}

	tf, _ := ioutil.TempFile("", "")
	tf.Write([]byte(`hello again`))
	_ = tf.Close()
	defer os.Remove(tf.Name())

	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	conn := &sendIQAndRandMock{
		sendIQMock: &sendIQMock{},
		rand: func() io.Reader {
			return bytes.NewBufferString("abcdefesdfgsdfgsdfg")
		},
	}

	var sendIQv1 string
	var sendIQv2 string
	var sendIQv3 interface{}

	conn.sendIQ = func(v1 string, v2 string, v3 interface{}) (<-chan data.Stanza, data.Cookie, error) {
		sendIQv1 = v1
		sendIQv2 = v2
		sendIQv3 = v3
		ch := make(chan data.Stanza, 1)
		ch <- data.Stanza{
			Value: &data.ClientIQ{
				Type: "result",
				Query: []byte(`<si xmlns="http://jabber.org/protocol/si">
  <feature xmlns="http://jabber.org/protocol/feature-neg">
    <x xmlns="jabber:x:data" type="submit">
      <field var="stream-method">
        <value>foo</value>
      </field>
    </x>
  </feature>
</si>`),
			},
		}
		return ch, 0, nil
	}

	mc := &mockHasConnectionAndConfigAndLog{
		c: conn,
		l: l,
	}

	ctx := &sendContext{
		file: tf.Name(),
		s:    mc,
		peer: "someone@some.where/foo",
	}

	res := ctx.offerSend()
	c.Assert(res, IsNil)
	<-done

	c.Assert(sendIQv1, Equals, "someone@some.where/foo")
	c.Assert(sendIQv2, Equals, "set")
	c.Assert(sendIQv3, Not(IsNil))
	c.Assert(sendIQv3, FitsTypeOf, data.SI{})
	siqv3 := sendIQv3.(data.SI)
	c.Assert(siqv3.ID, Equals, "sid8315164872671453793")
	c.Assert(siqv3.Profile, Equals, "http://jabber.org/protocol/si/profile/file-transfer")
	c.Assert(siqv3.Feature.Form.Fields, HasLen, 1)
	c.Assert(siqv3.Feature.Form.Fields[0].Var, Equals, "stream-method")
	c.Assert(siqv3.Feature.Form.Fields[0].Type, Equals, "list-single")
	c.Assert(siqv3.Feature.Form.Fields[0].Options, HasLen, 1)
	c.Assert(siqv3.Feature.Form.Fields[0].Options[0].Value, Equals, "foo")

	c.Assert(hook.Entries, HasLen, 1)
	c.Assert(hook.Entries[0].Level, Equals, log.InfoLevel)
	c.Assert(hook.Entries[0].Message, Equals, "Started sending file to peer using method")
	c.Assert(hook.Entries[0].Data, HasLen, 3)
	c.Assert(hook.Entries[0].Data["file"], Equals, tf.Name())
	c.Assert(hook.Entries[0].Data["peer"], Equals, "someone@some.where/foo")
	c.Assert(hook.Entries[0].Data["method"], Equals, "foo")
}
