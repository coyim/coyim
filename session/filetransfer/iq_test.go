package filetransfer

import (
	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/session/events"
	"github.com/coyim/coyim/session/mock"
	"github.com/coyim/coyim/xmpp/data"
	log "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	. "gopkg.in/check.v1"
)

type IQSuite struct{}

var _ = Suite(&IQSuite{})

type mockForInitIQ struct {
	*mock.SessionMock

	log coylog.Logger

	publishEventValue interface{}
}

func (m *mockForInitIQ) Log() coylog.Logger {
	return m.log
}

func (m *mockForInitIQ) PublishEvent(v interface{}) {
	m.publishEventValue = v
}

func (s *IQSuite) Test_InitIQ(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	m := &mockForInitIQ{
		log: l,
	}

	ret, iqtype, ign := InitIQ(m, &data.ClientIQ{
		From: "hello@goodbye.com/foo",
	}, data.SI{
		Profile: dirTransferProfile,
		File:    &data.File{},
		Feature: data.FeatureNegotation{
			Form: data.Form{
				Type: "form",
				Fields: []data.FormFieldX{
					data.FormFieldX{
						Var:  "stream-method",
						Type: "list-single",
						Options: []data.FormFieldOptionX{
							data.FormFieldOptionX{Value: "one"},
							data.FormFieldOptionX{Value: "three"},
							data.FormFieldOptionX{Value: "four"},
						},
					},
				},
			},
		},
	})

	c.Assert(ret, IsNil)
	c.Assert(iqtype, Equals, "")
	c.Assert(ign, Equals, true)

	ft := m.publishEventValue.(events.FileTransfer)
	c.Assert(ft.IsDirectory, Equals, true)

	c.Assert(hook.Entries, HasLen, 0)
}

func (s *IQSuite) Test_InitIQ_withEncryptedProfile(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	m := &mockForInitIQ{
		log: l,
	}

	ret, iqtype, ign := InitIQ(m, &data.ClientIQ{
		From: "hello@goodbye.com/foo",
	}, data.SI{
		Profile:       encryptedTransferProfile,
		EncryptedData: &data.EncryptedData{},
		File:          &data.File{},
		Feature: data.FeatureNegotation{
			Form: data.Form{
				Type: "form",
				Fields: []data.FormFieldX{
					data.FormFieldX{
						Var:  "stream-method",
						Type: "list-single",
						Options: []data.FormFieldOptionX{
							data.FormFieldOptionX{Value: "one"},
							data.FormFieldOptionX{Value: "three"},
							data.FormFieldOptionX{Value: "four"},
						},
					},
				},
			},
		},
	})

	c.Assert(ret, IsNil)
	c.Assert(iqtype, Equals, "")
	c.Assert(ign, Equals, true)

	ft := m.publishEventValue.(events.FileTransfer)
	c.Assert(ft.Encrypted, Equals, true)

	c.Assert(hook.Entries, HasLen, 0)
}

func (s *IQSuite) Test_InitIQ_failsIfStartedWithoutResource(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	m := &mockForInitIQ{
		log: l,
	}

	ret, iqtype, ign := InitIQ(m, &data.ClientIQ{
		From: "hello@goodbye.com",
	}, data.SI{
		File: &data.File{},
		Feature: data.FeatureNegotation{
			Form: data.Form{
				Type: "form",
				Fields: []data.FormFieldX{
					data.FormFieldX{
						Var:  "stream-method",
						Type: "list-single",
						Options: []data.FormFieldOptionX{
							data.FormFieldOptionX{Value: "one"},
							data.FormFieldOptionX{Value: "three"},
							data.FormFieldOptionX{Value: "four"},
						},
					},
				},
			},
		},
	})

	c.Assert(ret, IsNil)
	c.Assert(iqtype, Equals, "")
	c.Assert(ign, Equals, false)

	c.Assert(hook.Entries, HasLen, 1)
	c.Assert(hook.Entries[0].Level, Equals, log.WarnLevel)
	c.Assert(hook.Entries[0].Message, Equals, "Stanza sender doesn't contain resource - this shouldn't happen")
	c.Assert(hook.Entries[0].Data, HasLen, 1)
	c.Assert(hook.Entries[0].Data["from"], Equals, "hello@goodbye.com")
}

func (s *IQSuite) Test_InitIQ_failsWhenExtractingFileTransferOptions(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	m := &mockForInitIQ{
		log: l,
	}

	ret, iqtype, ign := InitIQ(m, &data.ClientIQ{
		From: "hello@goodbye.com/foo",
	}, data.SI{
		Profile: dirTransferProfile,
		File:    &data.File{},
		Feature: data.FeatureNegotation{
			Form: data.Form{
				Type: "fluff",
			},
		},
	})

	c.Assert(ret, IsNil)
	c.Assert(iqtype, Equals, "")
	c.Assert(ign, Equals, false)

	c.Assert(hook.Entries, HasLen, 1)
	c.Assert(hook.Entries[0].Level, Equals, log.WarnLevel)
	c.Assert(hook.Entries[0].Message, Equals, "Failed to parse stream initiation")
	c.Assert(hook.Entries[0].Data, HasLen, 1)
	c.Assert(hook.Entries[0].Data["error"], ErrorMatches, "Invalid form for file transfer initiation.*")
}
