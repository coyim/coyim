package filetransfer

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"

	"github.com/coyim/coyim/config"
	"github.com/coyim/coyim/coylog"
	sdata "github.com/coyim/coyim/session/data"
	"github.com/coyim/coyim/xmpp/data"
	xi "github.com/coyim/coyim/xmpp/interfaces"
	"github.com/coyim/coyim/xmpp/jid"
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

func (s *SendSuite) Test_sendContext_offerSend_failsIfNoFileExists(c *C) {
	ctx := &sendContext{
		file: "hopefully-a-file-that-doesnt-exist.txt",
	}

	res := ctx.offerSend()
	c.Assert(res, ErrorMatches, ".*(no such file or directory|cannot find the path specified).*")
}

func (s *SendSuite) Test_sendContext_offerSend_failsOnIncorrectSendingMechanism(c *C) {
	orgSupportedSendingMechanisms := supportedSendingMechanisms

	defer func() {
		supportedSendingMechanisms = orgSupportedSendingMechanisms
	}()

	supportedSendingMechanisms = map[string]func(*sendContext){}

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

	conn.sendIQ = func(v1 string, v2 string, v3 interface{}) (<-chan data.Stanza, data.Cookie, error) {
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

	done := make(chan error)
	ctrl := sdata.CreateFileTransferControl(func() bool { return false }, func(bool) {})
	ctx := &sendContext{
		file: tf.Name(),
		s:    mc,
		peer: "someone@some.where/foo",
		onErrorHook: func(ctx *sendContext, e error) {
			done <- e
		},
		control: ctrl,
	}

	ee := make(chan error)
	go func() {
		ctrl.WaitForError(func(eee error) {
			ee <- eee
		})
	}()

	res := ctx.offerSend()
	c.Assert(res, IsNil)

	e := <-done
	e2 := <-ee

	c.Assert(e, ErrorMatches, "Invalid sending mechanism sent from peer for file sending")
	c.Assert(e2, Equals, e)

	c.Assert(hook.Entries, HasLen, 0)
}

func (s *SendSuite) Test_sendContext_offerSend_failsOnInvalidSubmitForm(c *C) {
	orgSupportedSendingMechanisms := supportedSendingMechanisms

	defer func() {
		supportedSendingMechanisms = orgSupportedSendingMechanisms
	}()

	supportedSendingMechanisms = map[string]func(*sendContext){}

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

	conn.sendIQ = func(v1 string, v2 string, v3 interface{}) (<-chan data.Stanza, data.Cookie, error) {
		ch := make(chan data.Stanza, 1)
		ch <- data.Stanza{
			Value: &data.ClientIQ{
				Type: "result",
				Query: []byte(`<si xmlns="http://jabber.org/protocol/si">
  <feature xmlns="http://jabber.org/protocol/feature-neg">
    <x xmlns="jabber:x:data" type="submitfoo">
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

	done := make(chan error)
	ctrl := sdata.CreateFileTransferControl(func() bool { return false }, func(bool) {})
	ctx := &sendContext{
		file: tf.Name(),
		s:    mc,
		peer: "someone@some.where/foo",
		onErrorHook: func(ctx *sendContext, e error) {
			done <- e
		},
		control: ctrl,
	}

	ee := make(chan error)
	go func() {
		ctrl.WaitForError(func(eee error) {
			ee <- eee
		})
	}()

	res := ctx.offerSend()
	c.Assert(res, IsNil)

	e := <-done
	e2 := <-ee

	c.Assert(e, ErrorMatches, "Invalid data sent from peer for file sending")
	c.Assert(e2, Equals, e)

	c.Assert(hook.Entries, HasLen, 0)
}

func (s *SendSuite) Test_sendContext_offerSend_sendingIQFailsWithDecline(c *C) {
	orgSupportedSendingMechanisms := supportedSendingMechanisms

	defer func() {
		supportedSendingMechanisms = orgSupportedSendingMechanisms
	}()

	supportedSendingMechanisms = map[string]func(*sendContext){}

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

	conn.sendIQ = func(v1 string, v2 string, v3 interface{}) (<-chan data.Stanza, data.Cookie, error) {
		ch := make(chan data.Stanza, 1)
		ch <- data.Stanza{
			Value: &data.ClientIQ{
				Error: data.StanzaError{
					Code: "403",
				},
				Type: "error",
				Query: []byte(`<si xmlns="http://jabber.org/protocol/si">
  <feature xmlns="http://jabber.org/protocol/feature-neg">
    <x xmlns="jabber:x:data" type="submitfoo">
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

	done := make(chan bool, 1)
	ctrl := sdata.CreateFileTransferControl(func() bool { return false }, func(bool) {})

	notDeclined := make(chan bool)
	go func() {
		ctrl.WaitForFinish(func(v bool) {
			notDeclined <- v
		})
	}()

	ctx := &sendContext{
		file: tf.Name(),
		s:    mc,
		peer: "someone@some.where/foo",
		onDeclineHook: func(ctx *sendContext) {
			done <- true
		},
		control: ctrl,
	}

	res := ctx.offerSend()
	c.Assert(res, IsNil)

	<-done
	notDec := <-notDeclined

	c.Assert(notDec, Equals, false)
	c.Assert(hook.Entries, HasLen, 0)
}

func (s *SendSuite) Test_sendContext_offerSend_failsOnErrorIQ(c *C) {
	orgSupportedSendingMechanisms := supportedSendingMechanisms

	defer func() {
		supportedSendingMechanisms = orgSupportedSendingMechanisms
	}()

	supportedSendingMechanisms = map[string]func(*sendContext){}

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

	conn.sendIQ = func(v1 string, v2 string, v3 interface{}) (<-chan data.Stanza, data.Cookie, error) {
		ch := make(chan data.Stanza, 1)
		ch <- data.Stanza{
			Value: &data.ClientIQ{
				Type: "error",
			},
		}
		return ch, 0, nil
	}

	mc := &mockHasConnectionAndConfigAndLog{
		c: conn,
		l: l,
	}

	done := make(chan error)
	ctrl := sdata.CreateFileTransferControl(func() bool { return false }, func(bool) {})
	ctx := &sendContext{
		file: tf.Name(),
		s:    mc,
		peer: "someone@some.where/foo",
		onErrorHook: func(ctx *sendContext, e error) {
			done <- e
		},
		control: ctrl,
	}

	ee := make(chan error)
	go func() {
		ctrl.WaitForError(func(eee error) {
			ee <- eee
		})
	}()

	res := ctx.offerSend()
	c.Assert(res, IsNil)

	e := <-done
	e2 := <-ee

	c.Assert(e, ErrorMatches, "expected result IQ")
	c.Assert(e2, Equals, e)

	c.Assert(hook.Entries, HasLen, 0)
}

func (s *SendSuite) Test_sendContext_onFinish(c *C) {
	ctrl := sdata.CreateFileTransferControl(func() bool { return false }, func(bool) {})
	finished := false
	ctx := &sendContext{
		onFinishHook: func(ctx *sendContext) {
			finished = true
		},
		control: ctrl,
	}

	ff := make(chan bool, 1)
	go func() {
		ctrl.WaitForFinish(func(v bool) {
			ff <- v
		})
	}()

	ctx.onFinish()

	fval := <-ff
	c.Assert(finished, Equals, true)
	c.Assert(fval, Equals, true)
}

func (s *SendSuite) Test_sendContext_onUpdate(c *C) {
	ctrl := sdata.CreateFileTransferControl(func() bool { return false }, func(bool) {})
	vals := []int64{}
	ctx := &sendContext{
		onUpdateHook: func(ctx *sendContext, total int64) {
			vals = append(vals, total)
		},
		control:   ctrl,
		totalSize: 54354,
	}

	curr := make(chan int64, 2)
	tot := make(chan int64, 2)
	go func() {
		ctrl.WaitForUpdate(func(current, total int64) {
			curr <- current
			tot <- total
		})
	}()

	ctx.onUpdate(42)
	ctx.onUpdate(55)

	curr1 := <-curr
	curr2 := <-curr
	tot1 := <-tot
	tot2 := <-tot

	c.Assert(ctx.totalSent, Equals, int64(97))
	c.Assert(vals, DeepEquals, []int64{42, 97})
	c.Assert(curr1, Equals, int64(42))
	c.Assert(curr2, Equals, int64(97))
	c.Assert(tot1, Equals, int64(54354))
	c.Assert(tot2, Equals, int64(54354))
}

func (s *SendSuite) Test_sendContext_listenForCancellation(c *C) {
	ctrl := sdata.CreateFileTransferControl(func() bool { return false }, func(bool) {})
	ctx := &sendContext{
		control: ctrl,
	}

	done := make(chan bool, 1)
	go func() {
		ctx.listenForCancellation()
		done <- true
	}()
	ctrl.Cancel()
	<-done

	c.Assert(ctx.weWantToCancel, Equals, true)
}

func (s *SendSuite) Test_sendContext_initSend(c *C) {
	conn := &mockConnDiscoveryFeatures{
		f: func(string) ([]string, bool) {
			return []string{
				"foo",
				"http://jabber.org/protocol/si/profile/hubba",
				"http://jabber.org/protocol/si",
				"http://jabber.org/protocol/si/profile/bubba",
			}, true
		},
	}

	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	mc := &mockHasConnectionAndConfigAndLog{
		c: conn,
		l: l,
	}

	encDecCalled := false
	ctrl := sdata.CreateFileTransferControl(func() bool {
		return true
	}, func(bool) {
		encDecCalled = true
	})
	ctx := &sendContext{
		control: ctrl,
		s:       mc,
	}

	ctx.initSend()

	c.Assert(encDecCalled, Equals, true)

	c.Assert(hook.Entries, HasLen, 0)
}

func (s *SendSuite) Test_sendContext_initSend_failsWhenDiscoveringFeatures(c *C) {
	conn := &mockConnDiscoveryFeatures{
		f: func(string) ([]string, bool) {
			return nil, false
		},
	}

	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	mc := &mockHasConnectionAndConfigAndLog{
		c: conn,
		l: l,
	}

	ctrl := sdata.CreateFileTransferControl(func() bool {
		return true
	}, func(bool) {
	})
	ctx := &sendContext{
		control: ctrl,
		s:       mc,
	}

	ee := make(chan error)
	go func() {
		ctrl.WaitForError(func(eee error) {
			ee <- eee
		})
	}()

	ctx.initSend()

	e := <-ee

	c.Assert(e, ErrorMatches, "Problem discovering the features of the peer")
	c.Assert(hook.Entries, HasLen, 0)
}

func (s *SendSuite) Test_sendContext_initSend_wontSendUnencrypted(c *C) {
	conn := &mockConnDiscoveryFeatures{
		f: func(string) ([]string, bool) {
			return []string{
				"foo",
				"http://jabber.org/protocol/si/profile/hubba",
				"http://jabber.org/protocol/si",
				"http://jabber.org/protocol/si/profile/bubba",
			}, true
		},
	}

	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	mc := &mockHasConnectionAndConfigAndLog{
		c: conn,
		l: l,
	}

	ctrl := sdata.CreateFileTransferControl(func() bool {
		return false
	}, func(bool) {
	})
	ctx := &sendContext{
		control: ctrl,
		s:       mc,
	}

	ee := make(chan error)
	go func() {
		ctrl.WaitForError(func(eee error) {
			ee <- eee
		})
	}()

	ctx.initSend()

	e := <-ee

	c.Assert(e, ErrorMatches, "will not send unencrypted")
	c.Assert(hook.Entries, HasLen, 0)
}

type allSessionMock struct {
	*mockHasConnectionAndConfigAndLog

	createSymm func(jid.Any) []byte
	getAndWipe func(jid.Any) []byte
}

func (m *allSessionMock) CreateSymmetricKeyFor(v jid.Any) []byte {
	return m.createSymm(v)
}

func (m *allSessionMock) GetAndWipeSymmetricKeyFor(v jid.Any) []byte {
	return m.getAndWipe(v)
}

func (s *SendSuite) Test_InitSend(c *C) {
	conn := &mockConnDiscoveryFeatures{
		f: func(string) ([]string, bool) {
			return []string{
				"foo",
				"http://jabber.org/protocol/si/profile/hubba",
				"http://jabber.org/protocol/si",
				"http://jabber.org/protocol/si/profile/bubba",
			}, true
		},
	}

	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	var calledCreateSymm jid.Any

	mc := &allSessionMock{
		mockHasConnectionAndConfigAndLog: &mockHasConnectionAndConfigAndLog{
			c: conn,
			l: l,
		},

		createSymm: func(v jid.Any) []byte {
			calledCreateSymm = v
			return nil
		},
	}

	ctrl := InitSend(mc, jid.Parse("some@one.org"), "a file somewhere", nil, nil)

	c.Assert(ctrl, Not(IsNil))
	c.Assert(calledCreateSymm.String(), Equals, "some@one.org")

	c.Assert(hook.Entries, HasLen, 0)
}
