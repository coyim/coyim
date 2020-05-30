package ui

import (
	"io/ioutil"
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/gotk3adapter/glib_mock"

	. "gopkg.in/check.v1"
)

var escapingTests = []string{
	"",
	"foo",
	"foo\\",
	"foo\\x01",
	"العربية",
}

func Test(t *testing.T) { TestingT(t) }

func init() {
	log.SetOutput(ioutil.Discard)
	i18n.InitLocalization(&glib_mock.Mock{})
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
