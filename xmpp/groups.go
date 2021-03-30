package xmpp

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"time"

	"github.com/coyim/coyim/xmpp/data"
)

const requestDelimiterXML = `
<iq type='get' id='%x'>
  <query xmlns='jabber:iq:private'>
    <roster xmlns='roster:delimiter'/>
  </query>
</iq>
`

type rosterDelimiter struct {
	XMLName   xml.Name `xml:"roster:delimiter roster"`
	Delimiter string   `xml:",chardata"`
}

type rosterQuery struct {
	XMLName   xml.Name        `xml:"jabber:iq:private query"`
	Delimiter rosterDelimiter `xml:"roster:delimiter roster"`
}

var rosterRequestTimeout = 5 * time.Second

// GetRosterDelimiter blocks and waits for the roster delimiter to be delivered
func (c *conn) GetRosterDelimiter() (string, error) {
	rep, _, err := c.RequestRosterDelimiter()
	if err != nil {
		return "", err
	}

	select {
	case iqStanza := <-rep:
		stanza, ok := iqStanza.Value.(*data.ClientIQ)
		if ok {
			var rst rosterQuery
			if err := xml.NewDecoder(bytes.NewBuffer(stanza.Query)).Decode(&rst); err != nil || len(rst.Delimiter.Delimiter) == 0 {
				return "", nil
			}
			return rst.Delimiter.Delimiter, nil
		}
	case <-time.After(rosterRequestTimeout):
	}

	return "", nil
}

var createInflight = func(c *conn, cookie data.Cookie, to string) (<-chan data.Stanza, data.Cookie, error) {
	return c.createInflight(cookie, to)
}

// RequestRosterDelimiter will request the roster delimiter
func (c *conn) RequestRosterDelimiter() (<-chan data.Stanza, data.Cookie, error) {
	cookie := c.getCookie()

	var outb bytes.Buffer
	_, _ = fmt.Fprintf(&outb, requestDelimiterXML, cookie)

	_, err := c.safeWrite(outb.Bytes())
	if err != nil {
		return nil, 0, err
	}

	return createInflight(c, cookie, "")
}
