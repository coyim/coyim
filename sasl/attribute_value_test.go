package sasl

import . "gopkg.in/check.v1"

type AttributeValueSuite struct{}

var _ = Suite(&AttributeValueSuite{})

func (s *AttributeValueSuite) Test_ParseAttributeValuePairs_parsesAnEmptyString(c *C) {
	result := ParseAttributeValuePairs([]byte(""))
	c.Assert(result, DeepEquals, make(AttributeValuePairs))
}

func (s *AttributeValueSuite) Test_ParseAttributeValuePairs_parsesASimplePair(c *C) {
	result := ParseAttributeValuePairs([]byte("foo=bar"))
	c.Assert(result, DeepEquals, AttributeValuePairs{
		"foo": "bar",
	})
}

func (s *AttributeValueSuite) Test_ParseAttributeValuePairs_parsesMoreThanOnePair(c *C) {
	result := ParseAttributeValuePairs([]byte("foo=bar,bar=foo"))
	c.Assert(result, DeepEquals, AttributeValuePairs{
		"foo": "bar",
		"bar": "foo",
	})
}

func (s *AttributeValueSuite) Test_ParseAttributeValuePairs_returnsOnlyTheLastValueWithTheSameKey(c *C) {
	result := ParseAttributeValuePairs([]byte("foo=bar,foo=quux"))
	c.Assert(result, DeepEquals, AttributeValuePairs{
		"foo": "quux",
	})
}

func (s *AttributeValueSuite) Test_ParseAttributeValuePairs_triesToMatchQuotedValue(c *C) {
	result := ParseAttributeValuePairs([]byte("foo=\"bar\""))
	c.Assert(result, DeepEquals, AttributeValuePairs{
		"foo": "bar",
	})
}

func (s *AttributeValueSuite) Test_ParseAttributeValuePairs_ignoresBadValue(c *C) {
	result := ParseAttributeValuePairs([]byte("foo=bar,hmm,bar=foo"))
	c.Assert(result, DeepEquals, AttributeValuePairs{
		"foo": "bar",
		"bar": "foo",
	})
}

func (s *AttributeValueSuite) Test_ParseAttributeValuePairs_doesntAcceptEmptyKey(c *C) {
	result := ParseAttributeValuePairs([]byte("=bar,=\"flux\",bar=foo"))
	c.Assert(result, DeepEquals, AttributeValuePairs{
		"bar": "foo",
	})
}

func (s *AttributeValueSuite) Test_ParseAttributeValuePairs_doesntAcceptEmptyValueExceptForQuotes(c *C) {
	result := ParseAttributeValuePairs([]byte("bb=hey,ab=\"\",bar="))
	c.Assert(result, DeepEquals, AttributeValuePairs{
		"bb": "hey",
		"ab": "\"\"",
	})
}
