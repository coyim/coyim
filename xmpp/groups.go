package xmpp

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"time"
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
	delimiter string   `xml:",chardata"`
}

type rosterQuery struct {
	XMLName   xml.Name        `xml:"jabber:iq:private query"`
	delimiter rosterDelimiter `xml:"roster:delimiter roster"`
}

// GetRosterDelimiter blocks and waits for the roster delimiter to be delivered
func (c *Conn) GetRosterDelimiter() (string, error) {
	rep, _, err := c.RequestRosterDelimiter()
	if err != nil {
		return "", err
	}

	select {
	case iqStanza := <-rep:
		stanza, ok := iqStanza.Value.(*ClientIQ)
		if ok {
			var rst rosterQuery
			if err := xml.NewDecoder(bytes.NewBuffer(stanza.Query)).Decode(&rst); err != nil || len(rst.delimiter.delimiter) == 0 {
				return "", nil
			}
			return rst.delimiter.delimiter, nil
		}
	case <-time.After(5000 * time.Millisecond):
	}

	return "", nil
}

// RequestRosterDelimiter will request the roster delimiter
func (c *Conn) RequestRosterDelimiter() (<-chan Stanza, Cookie, error) {
	cookie := c.getCookie()
	if _, err := fmt.Fprintf(c.out, requestDelimiterXML, cookie); err != nil {
		return nil, 0, err
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	ch := make(chan Stanza, 1)
	c.inflights[cookie] = inflight{ch, ""}
	return ch, cookie, nil
}
