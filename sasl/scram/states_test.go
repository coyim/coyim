package scram

import (
	"github.com/coyim/coyim/sasl"
	. "gopkg.in/check.v1"
)

func (s *ScramSuite) Test_states_start_finished(c *C) {
	c.Assert((start{}).finished(), Equals, false)
}

func (s *ScramSuite) Test_states_start_next_failsOnMissingAuthID(c *C) {
	_, _, e := (start{}).next(nil, sasl.Properties{}, nil, nil)
	c.Assert(e, ErrorMatches, "missing property '\\\\x00'")
}

func (s *ScramSuite) Test_states_start_next_failsOnMissingClientNonce(c *C) {
	_, _, e := (start{}).next(nil, sasl.Properties{
		sasl.AuthID: "foo",
	}, nil, nil)
	c.Assert(e, ErrorMatches, "missing property '\\\\x06'")
}

func (s *ScramSuite) Test_states_expectingServerFirstMessage_next_failsOnMissingNonce(c *C) {
	_, _, e := (expectingServerFirstMessage{}).next(nil, sasl.Properties{}, sasl.AttributeValuePairs{}, nil)
	c.Assert(e, ErrorMatches, "missing parameter in server challenge")
}

func (s *ScramSuite) Test_states_expectingServerFirstMessage_next_failsOnMissingSalt(c *C) {
	_, _, e := (expectingServerFirstMessage{}).next(nil, sasl.Properties{}, sasl.AttributeValuePairs{
		"r": "foo",
	}, nil)
	c.Assert(e, ErrorMatches, "missing parameter in server challenge")
}

func (s *ScramSuite) Test_states_expectingServerFirstMessage_next_failsOnBadSalt(c *C) {
	_, _, e := (expectingServerFirstMessage{}).next(nil, sasl.Properties{}, sasl.AttributeValuePairs{
		"r": "foo",
		"s": "&%*&^*(&(",
	}, nil)
	c.Assert(e, ErrorMatches, "illegal base64 data at input byte 0")
}

func (s *ScramSuite) Test_states_expectingServerFirstMessage_next_failsOnMissingCount(c *C) {
	_, _, e := (expectingServerFirstMessage{}).next(nil, sasl.Properties{}, sasl.AttributeValuePairs{
		"r": "foo",
		"s": "salt",
	}, nil)
	c.Assert(e, ErrorMatches, "missing parameter in server challenge")
}

func (s *ScramSuite) Test_states_expectingServerFirstMessage_next_failsOnBadCount(c *C) {
	_, _, e := (expectingServerFirstMessage{}).next(nil, sasl.Properties{}, sasl.AttributeValuePairs{
		"r": "foo",
		"s": "salt",
		"i": "not a number",
	}, nil)
	c.Assert(e, ErrorMatches, "strconv.Atoi.*")
}

func (s *ScramSuite) Test_states_expectingServerFirstMessage_next_failsOnMissingPassword(c *C) {
	_, _, e := (expectingServerFirstMessage{}).next(nil, sasl.Properties{}, sasl.AttributeValuePairs{
		"r": "foo",
		"s": "salt",
		"i": "123",
	}, nil)
	c.Assert(e, ErrorMatches, "missing property '\\\\x01'")
}

func (s *ScramSuite) Test_states_expectingServerFirstMessage_next_failsOnMissingClientNonce(c *C) {
	_, _, e := (expectingServerFirstMessage{}).next(nil, sasl.Properties{
		sasl.Password: "foo",
	}, sasl.AttributeValuePairs{
		"r": "foo",
		"s": "salt",
		"i": "123",
	}, nil)
	c.Assert(e, ErrorMatches, "missing property '\\\\x06'")
}

func (s *ScramSuite) Test_states_expectingServerFirstMessage_next_failsIfClientNonceIsNotPrefixOfServerNonce(c *C) {
	_, _, e := (expectingServerFirstMessage{}).next(nil, sasl.Properties{
		sasl.Password:    "foo",
		sasl.ClientNonce: "something else",
	}, sasl.AttributeValuePairs{
		"r": "foo",
		"s": "salt",
		"i": "123",
	}, nil)
	c.Assert(e, ErrorMatches, "nonce mismatch")
}

func (s *ScramSuite) Test_calculateChannelBindingPrefix_withPlus(c *C) {
	c.Assert(calculateChannelBindingPrefix(true, false), Equals, "p=tls-unique")
}

func (s *ScramSuite) Test_states_expectingServerFirstMessage_calculateChannelBinding(c *C) {
	st := expectingServerFirstMessage{
		plus: true,
	}
	res := st.calculateChannelBinding([]byte("something"))
	c.Assert(res, Equals, "cD10bHMtdW5pcXVlLCxzb21ldGhpbmc=")
}

func (s *ScramSuite) Test_states_expectingServerFinalMessage_next_failsOnComparingServerAuth(c *C) {
	st := expectingServerFinalMessage{
		serverAuthentication: []byte("foo bar"),
	}
	_, _, e := st.next(sasl.Token([]byte("something else")), nil, nil, nil)
	c.Assert(e, ErrorMatches, "server signature mismatch")
}

func (s *ScramSuite) Test_states_finished_next(c *C) {
	st := finished{}
	nst, _, _ := st.next(nil, nil, nil, nil)
	c.Assert(nst, DeepEquals, st)
}
