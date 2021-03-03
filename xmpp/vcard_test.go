package xmpp

import (
	"bytes"
	"errors"

	"github.com/coyim/coyim/xmpp/data"
	. "gopkg.in/check.v1"
)

type VCardSuite struct{}

var _ = Suite(&VCardSuite{})

type errWriter struct {
	e error
}

func (w *errWriter) Write([]byte) (int, error) {
	return 0, w.e
}

func (s *VCardSuite) Test_conn_RequestVCard_failsIfWritingFails(c *C) {
	cn := &conn{
		out: &errWriter{errors.New("error marker")},
	}
	ch, cookie, e := cn.RequestVCard()
	c.Assert(e, ErrorMatches, "error marker")
	c.Assert(ch, IsNil)
	c.Assert(cookie, Equals, data.Cookie(0))
}

func (s *VCardSuite) Test_conn_RequestVCard_succeds(c *C) {
	var out bytes.Buffer

	cn := &conn{
		out:       &out,
		inflights: make(map[data.Cookie]inflight),
	}

	ch, _, e := cn.RequestVCard()
	c.Assert(e, IsNil)
	c.Assert(out.String(), Matches, "<iq type='get' id='.*?'><vCard xmlns='vcard-temp'/></iq>")
	c.Assert(ch, Not(IsNil))
}
