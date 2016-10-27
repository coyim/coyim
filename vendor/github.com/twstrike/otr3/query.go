package otr3

import (
	"bytes"
	"strconv"
	"time"
)

func isQueryMessage(msg ValidMessage) bool {
	return bytes.HasPrefix(msg, queryMarker)
}

func parseOTRQueryMessage(msg ValidMessage) []int {
	ret := []int{}

	if bytes.HasPrefix(msg, queryMarker) && len(msg) > len(queryMarker) {
		versions := msg[len(queryMarker):]

		if versions[0] == '?' {
			ret = append(ret, 1)
			versions = versions[1:]
		}

		if len(versions) > 0 && versions[0] == 'v' {
			for _, c := range versions {
				if v, err := strconv.Atoi(string(c)); err == nil {
					ret = append(ret, v)
				}
			}
		}
	}

	return ret
}

func extractVersionsFromQueryMessage(p policies, msg ValidMessage) int {
	versions := 0
	for _, v := range parseOTRQueryMessage(msg) {
		switch {
		case v == 3 && p.has(allowV3):
			versions |= (1 << 3)
		case v == 2 && p.has(allowV2):
			versions |= (1 << 2)
		}
	}

	return versions
}

var timeoutLength = time.Duration(1) * time.Minute

func isWithinTimeToIgnoreQueryMessage(t time.Time) bool {
	return t.Add(timeoutLength).After(time.Now())

}

func (c *Conversation) receiveQueryMessage(msg ValidMessage) ([]messageWithHeader, error) {
	versions := extractVersionsFromQueryMessage(c.Policies, msg)
	err := c.commitToVersionFrom(versions)
	if err != nil {
		return nil, err
	}

	if (c.msgState == encrypted && isWithinTimeToIgnoreQueryMessage(c.lastMessageStateChange)) ||
		(c.ake != nil && isWithinTimeToIgnoreQueryMessage(c.ake.lastStateChange)) {
		return nil, nil
	}

	ts, err := c.sendDHCommit()
	return c.potentialAuthError(compactMessagesWithHeader(ts), err)
}

//QueryMessage will return a QueryMessage determined by Conversation Policies
func (c Conversation) QueryMessage() ValidMessage {
	queryMessage := []byte("?OTRv")

	if c.Policies.has(allowV2) {
		queryMessage = append(queryMessage, '2')
	}

	if c.Policies.has(allowV3) {
		queryMessage = append(queryMessage, '3')
	}

	suffix := "?"
	if c.friendlyQueryMessage != "" {
		suffix = "? " + c.friendlyQueryMessage
	}

	return append(queryMessage, suffix...)
}

func (c *Conversation) SetFriendlyQueryMessage(msg string) {
	c.friendlyQueryMessage = msg
}
