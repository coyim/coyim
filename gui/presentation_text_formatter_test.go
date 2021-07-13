package gui

import (
	"github.com/coyim/coyim/i18n"
	. "gopkg.in/check.v1"
)

type PresentationTextFormatterSuite struct{}

var _ = Suite(&PresentationTextFormatterSuite{})

func (s *PresentationTextFormatterSuite) SetUpSuite(c *C) {
	initMUCI18n()
}

func (*PresentationTextFormatterSuite) Test_presentationTextFormatter_shouldReturnAllFormats(c *C) {
	ft := newPresentationTextFormatter(i18n.Local("This is a test"))

	c.Assert(ft.String(), Equals, "[localized] This is a test")
	c.Assert(ft.formats, HasLen, 0)
	c.Assert(ft.formats, IsNil)

	ft = newPresentationTextFormatter(i18n.Local("Hi, my nickname is {{ nickname \"bla\" }}"))
	c.Assert(ft.String(), Equals, "[localized] Hi, my nickname is bla")
	c.Assert(ft.formats, HasLen, 1)
	c.Assert(ft.formats, DeepEquals, []*presentationTextFormat{
		{
			typ:    "nickname",
			value:  "bla",
			start:  31,
			length: 3,
		},
	})

	ft = newPresentationTextFormatter(i18n.Local("{{ nickname \"Jhon\" }} is not an administrator anymore."))
	c.Assert(ft.String(), Equals, "[localized] Jhon is not an administrator anymore.")
	c.Assert(ft.formats, HasLen, 1)
	c.Assert(ft.formats, DeepEquals, []*presentationTextFormat{
		{
			typ:    "nickname",
			value:  "Jhon",
			start:  12,
			length: 4,
		},
	})

	ft = newPresentationTextFormatter(i18n.Local("Jhon is not an {{ affiliation \"administrator\" }} anymore."))
	c.Assert(ft.String(), Equals, "[localized] Jhon is not an administrator anymore.")
	c.Assert(ft.formats, HasLen, 1)
	c.Assert(ft.formats, DeepEquals, []*presentationTextFormat{
		{
			typ:    "affiliation",
			value:  "administrator",
			start:  27,
			length: 13,
		},
	})

	ft = newPresentationTextFormatter(i18n.Local("Jhon is not an {{ role \"moderator\" }} anymore."))
	c.Assert(ft.String(), Equals, "[localized] Jhon is not an moderator anymore.")
	c.Assert(ft.formats, HasLen, 1)
	c.Assert(ft.formats, DeepEquals, []*presentationTextFormat{
		{
			typ:    "role",
			value:  "moderator",
			start:  27,
			length: 9,
		},
	})

	ft = newPresentationTextFormatter(i18n.Local("The affiliation of {{ nickname \"Alberto\" }} was changed. Now he is a {{ role \"moderator\" }}."))
	c.Assert(ft.String(), Equals, "[localized] The affiliation of Alberto was changed. Now he is a moderator.")
	c.Assert(ft.formats, HasLen, 2)
	c.Assert(ft.formats, DeepEquals, []*presentationTextFormat{
		{
			typ:    "nickname",
			value:  "Alberto",
			start:  31,
			length: 7,
		},
		{
			typ:    "role",
			value:  "moderator",
			start:  64,
			length: 9,
		},
	})

	ft = newPresentationTextFormatter(i18n.Local("{{ nickname \"Alberto\" }} is not {{ affiliation \"an administrator\" }} anymore. Now he is a {{ role \"moderator\" }}."))
	c.Assert(ft.String(), Equals, "[localized] Alberto is not an administrator anymore. Now he is a moderator.")
	c.Assert(ft.formats, HasLen, 3)
	c.Assert(ft.formats, DeepEquals, []*presentationTextFormat{
		{
			typ:    "nickname",
			value:  "Alberto",
			start:  12,
			length: 7,
		},
		{
			typ:    "affiliation",
			value:  "an administrator",
			start:  27,
			length: 16,
		},
		{
			typ:    "role",
			value:  "moderator",
			start:  65,
			length: 9,
		},
	})
}
