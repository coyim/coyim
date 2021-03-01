package digestmd5

import (
	"github.com/coyim/coyim/sasl"
	. "gopkg.in/check.v1"
)

func (s *DigestMD5Suite) Test_states_digestChallenge_challenge_failsOnMissingAuthID(c *C) {
	st := digestChallenge{}
	_, _, e := st.challenge(sasl.Properties{}, nil)
	c.Assert(e, ErrorMatches, "missing property '\\\\x00'")
}

func (s *DigestMD5Suite) Test_states_digestChallenge_challenge_failsOnMissingPassword(c *C) {
	st := digestChallenge{}
	_, _, e := st.challenge(sasl.Properties{
		sasl.AuthID: "something",
	}, nil)
	c.Assert(e, ErrorMatches, "missing property '\\\\x01'")
}

func (s *DigestMD5Suite) Test_states_digestChallenge_challenge_failsOnMissingRealm(c *C) {
	st := digestChallenge{}
	_, _, e := st.challenge(sasl.Properties{
		sasl.AuthID:   "something",
		sasl.Password: "bla",
	}, sasl.AttributeValuePairs{})
	c.Assert(e, ErrorMatches, "missing parameter in server challenge")
}

func (s *DigestMD5Suite) Test_states_digestChallenge_challenge_failsOnMissingService(c *C) {
	st := digestChallenge{}
	_, _, e := st.challenge(sasl.Properties{
		sasl.AuthID:   "something",
		sasl.Password: "bla",
	}, sasl.AttributeValuePairs{
		"realm": "something",
	})
	c.Assert(e, ErrorMatches, "missing property '\\\\x04'")
}

func (s *DigestMD5Suite) Test_states_digestChallenge_challenge_failsOnMissingClientNonce(c *C) {
	st := digestChallenge{}
	_, _, e := st.challenge(sasl.Properties{
		sasl.AuthID:   "something",
		sasl.Password: "bla",
		sasl.Service:  "foo",
	}, sasl.AttributeValuePairs{
		"realm": "something",
	})
	c.Assert(e, ErrorMatches, "missing property '\\\\x06'")
}

func (s *DigestMD5Suite) Test_states_digestChallenge_challenge_failsOnMissingQOP(c *C) {
	st := digestChallenge{}
	_, _, e := st.challenge(sasl.Properties{
		sasl.AuthID:      "something",
		sasl.Password:    "bla",
		sasl.Service:     "foo",
		sasl.ClientNonce: "hmm",
	}, sasl.AttributeValuePairs{
		"realm": "something",
	})
	c.Assert(e, ErrorMatches, "missing property '\\\\x05'")
}

func (s *DigestMD5Suite) Test_states_digestChallenge_challenge_failsOnMissingNonce(c *C) {
	st := digestChallenge{}
	_, _, e := st.challenge(sasl.Properties{
		sasl.AuthID:      "something",
		sasl.Password:    "bla",
		sasl.Service:     "foo",
		sasl.ClientNonce: "hmm",
		sasl.QOP:         "name",
	}, sasl.AttributeValuePairs{
		"realm": "something",
	})
	c.Assert(e, ErrorMatches, "missing parameter in server challenge")
}

func (s *DigestMD5Suite) Test_states_digestChallenge_challenge_succeeds(c *C) {
	st := digestChallenge{}
	nst, t, e := st.challenge(sasl.Properties{
		sasl.AuthID:      "something",
		sasl.Password:    "bla",
		sasl.Service:     "foo",
		sasl.ClientNonce: "hmm",
		sasl.QOP:         "name",
	}, sasl.AttributeValuePairs{
		"realm": "something",
		"nonce": "bla",
	})
	c.Assert(e, IsNil)
	c.Assert(nst, DeepEquals, responseAuth{})
	c.Assert(t, DeepEquals, sasl.Token("charset=utf-8,username=\"something\",realm=\"something\",nonce=\"bla\",nc=00000001,cnonce=\"hmm\",digest-uri=\"foo/something\",response=3b246d774f9cf11ce20eb75b64c3092e,qop=auth"))
}

func (s *DigestMD5Suite) Test_states_digestChallenge_challenge_succeedsWithAuthorizationID(c *C) {
	st := digestChallenge{}
	nst, t, e := st.challenge(sasl.Properties{
		sasl.AuthID:      "something",
		sasl.Password:    "bla",
		sasl.Service:     "foo",
		sasl.ClientNonce: "hmm",
		sasl.QOP:         "name",
		sasl.AuthZID:     "blub",
	}, sasl.AttributeValuePairs{
		"realm": "something",
		"nonce": "bla",
	})
	c.Assert(e, IsNil)
	c.Assert(nst, DeepEquals, responseAuth{})
	c.Assert(t, DeepEquals, sasl.Token("charset=utf-8,username=\"something\",realm=\"something\",nonce=\"bla\",nc=00000001,cnonce=\"hmm\",digest-uri=\"foo/something\",response=911d54f9a312286c5eac7283110fc37e,qop=auth,authzid=\"blub\""))
}

func (s *DigestMD5Suite) Test_states_digestChallenge_challenge_succeedsWithAuthInt(c *C) {
	st := digestChallenge{}
	nst, t, e := st.challenge(sasl.Properties{
		sasl.AuthID:      "something",
		sasl.Password:    "bla",
		sasl.Service:     "foo",
		sasl.ClientNonce: "hmm",
		sasl.QOP:         "auth-int",
	}, sasl.AttributeValuePairs{
		"realm": "something",
		"nonce": "bla",
	})
	c.Assert(e, IsNil)
	c.Assert(nst, DeepEquals, responseAuth{})
	c.Assert(t, DeepEquals, sasl.Token("charset=utf-8,username=\"something\",realm=\"something\",nonce=\"bla\",nc=00000001,cnonce=\"hmm\",digest-uri=\"foo/something\",response=ec2e845596bc500305a99f434fc205b2,qop=auth"))
}

func (s *DigestMD5Suite) Test_states_responseAuth_challenge_failsIfNoRSPauthProvided(c *C) {
	st := responseAuth{}
	_, _, e := st.challenge(nil, nil)
	c.Assert(e, ErrorMatches, "missing parameter in server challenge")
}

func (s *DigestMD5Suite) Test_states_finished_challenge_returnsItself(c *C) {
	st := finished{}
	nst, _, e := st.challenge(nil, nil)
	c.Assert(nst, Equals, st)
	c.Assert(e, IsNil)
}
