package ui

import (
	"errors"

	. "gopkg.in/check.v1"
)

var escapingTests = []string{
	"",
	"foo",
	"foo\\",
	"foo\\x01",
	"العربية",
}

type UISuite struct{}

var _ = Suite(&UISuite{})

func (s *UISuite) TestEscaping(t *C) {
	for _, test := range escapingTests {
		escaped := EscapeNonASCII(test)
		unescaped, err := UnescapeNonASCII(escaped)
		if err != nil {
			t.Errorf("Error unescaping '%s' (from '%s')", escaped, test)
			continue
		}
		if unescaped != test {
			t.Errorf("Unescaping didn't return the original value: '%s' -> '%s' -> '%s'", test, escaped, unescaped)
		}
	}
}

func (s *UISuite) TestHTMLStripping(t *C) {
	raw := []byte("<hr>This is some <font color='green'>html</font><br />.")
	exp := []byte("This is some html.")
	res := StripHTML(raw)

	t.Check(res, DeepEquals, exp)
}

func (s *UISuite) Test_StripSomeHTML(t *C) {
	raw := []byte("<p>This is <walloftext>some</walloftext> <FONT color='green'>html</font><br />.")
	exp := "This is <walloftext>some</walloftext> html."
	res := StripSomeHTML(raw)

	t.Check(string(res), DeepEquals, exp)
}

func (s *UISuite) Test_EscapeAllHTMLTags(t *C) {
	raw := "<p><This> <!--is--> < walloftext >some</walloftext> <FONT color='green'>html</font><!DOCTYPE html><br />."
	exp := "&lt;p&gt;&lt;This&gt; &lt;!--is--&gt; < walloftext >some&lt;/walloftext&gt; &lt;FONT color=&#39;green&#39;&gt;html&lt;/font&gt;&lt;!DOCTYPE html&gt;&lt;br /&gt;."
	res := EscapeAllHTMLTags(raw)

	t.Check(res, DeepEquals, exp)
}

func (s *UISuite) Test_UnescapeNonASCII_failsOnEscapeSequenceAtEnd(c *C) {
	res, e := UnescapeNonASCII("foo \\")
	c.Assert(res, Equals, "")
	c.Assert(e, ErrorMatches, "truncated escape sequence at end: .*")
}

func (s *UISuite) Test_UnescapeNonASCII_failsOnBadEscapeType(c *C) {
	res, e := UnescapeNonASCII("foo \\q01")
	c.Assert(res, Equals, "")
	c.Assert(e, ErrorMatches, "escape sequence didn't start with .*")
}

func (s *UISuite) Test_UnescapeNonASCII_failsOnBadEscapeValue(c *C) {
	res, e := UnescapeNonASCII("foo \\xQQ")
	c.Assert(res, Equals, "")
	c.Assert(e, ErrorMatches, "failed to parse value in .*")
}

type errReader struct {
	e error
}

func (r *errReader) Read([]byte) (int, error) {
	return 0, r.e
}

func (s *UISuite) Test_EscapeAllHTMLTags_fails(c *C) {
	r := &errReader{e: errors.New("something")}
	res := escapeAllHTMLTagsInternal([]byte("bla <foo something"), r)
	c.Assert(res, Equals, "bla <foo something")
}

func (s *UISuite) Test_StripHTML_fails(c *C) {
	r := &errReader{e: errors.New("something")}
	res := stripHTMLInternal([]byte("bla <foo something"), r)
	c.Assert(string(res), Equals, "bla <foo something")
}

func (s *UISuite) Test_StripSomeHTML_fails(c *C) {
	r := &errReader{e: errors.New("something")}
	res := stripSomeHTMLInternal([]byte("bla <foo something"), r)
	c.Assert(string(res), Equals, "bla <foo something")
}

func (s *UISuite) Test_StripSomeHTML_keepsComments(c *C) {
	res := StripSomeHTML([]byte("bla <!-- hello -->"))
	c.Assert(string(res), Equals, "bla <!-- hello -->")
}

func (s *UISuite) Test_StripSomeHTML_keepsDoctype(c *C) {
	res := StripSomeHTML([]byte("bla <!DOCTYPE bla [ ]>"))
	c.Assert(string(res), Equals, "bla <!DOCTYPE bla [ ]>")
}

func (s *UISuite) Test_UnescapeNewlineTags_fails(c *C) {
	r := &errReader{e: errors.New("something")}
	res := unescapeNewlineTagsInternal([]byte("bla foo"), r)
	c.Assert(string(res), Equals, "bla foo")
}

func (s *UISuite) Test_UnescapeNewlineTags_escapesEverything(c *C) {
	res := UnescapeNewlineTags([]byte("bla <br> something <!-- a comment --> hello <br/> <i>foo</i>something else <!DOCTYPE bla [ ] > eh <br>"))
	c.Assert(string(res), Equals, ""+
		"bla \n"+
		" something <!-- a comment --> hello \n"+
		" <i>foo</i>something else <!DOCTYPE bla [ ] > eh \n")
}
