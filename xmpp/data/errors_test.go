package data

import (
	"bytes"
	"encoding/xml"

	. "gopkg.in/check.v1"
)

type ErrorsSuite struct{}

var _ = Suite(&ErrorsSuite{})

func (s *ErrorsSuite) Test_StreamError_String_returnsTextIfAvailable(c *C) {
	c.Assert((&StreamError{Text: "something"}).String(), Equals, "something")
}

func (s *ErrorsSuite) Test_StreamError_String_returnsAppSpecificConditionIfAvailable(c *C) {
	c.Assert((&StreamError{AppSpecificCondition: &Any{XMLName: xml.Name{Space: "foo", Local: "bar"}}}).String(), Equals, "{foo bar}")
}

func (s *ErrorsSuite) Test_StreamError_String_returnsEmptyIfNeitherNameOrCondition(c *C) {
	c.Assert((&StreamError{}).String(), Equals, "")
}

func (s *ErrorsSuite) Test_StreamErrorCondition_MarshalXML_marshalsProperly(c *C) {
	var result bytes.Buffer
	t := HostGone

	t.MarshalXML(xml.NewEncoder(&result), xml.StartElement{})

	c.Assert(result.String(), Equals, "")
}
