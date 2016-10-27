package sasl

import . "gopkg.in/check.v1"

type TokenSuite struct{}

var _ = Suite(&TokenSuite{})

func (s *TokenSuite) Test_Token_String_returnsTheStringValue(c *C) {
	result := Token("foo bar").String()
	c.Assert(result, DeepEquals, "foo bar")
}

func (s *TokenSuite) Test_Token_Encode_willEncode(c *C) {
	result := Token("foo bar").Encode()
	c.Assert(string(result), DeepEquals, "Zm9vIGJhcg==")
}

func (s *TokenSuite) Test_DecodeToken_willDecode(c *C) {
	result, err := DecodeToken([]byte("Zm9vIGJhcg=="))
	c.Assert(err, IsNil)
	c.Assert(result.String(), DeepEquals, "foo bar")
}

func (s *TokenSuite) Test_DecodeToken_willReturnErrorOnFailure(c *C) {
	_, err := DecodeToken([]byte("****"))
	c.Assert(err.Error(), Equals, "illegal base64 data at input byte 0")
}
