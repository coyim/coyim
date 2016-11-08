package xmpp

import (
	"encoding/xml"
	"errors"

	"github.com/twstrike/coyim/xmpp/data"

	. "gopkg.in/check.v1"
)

type FormsXMPPSuite struct{}

var _ = Suite(&FormsXMPPSuite{})

func (s *FormsXMPPSuite) Test_processForm_returnsErrorFromCallback(c *C) {
	e := errors.New("some kind of error")
	f := &data.Form{}
	_, err := processForm(f, nil, func(title, instructions string, fields []interface{}) error {
		return e
	})

	c.Assert(err, Equals, e)
}

func (s *FormsXMPPSuite) Test_processForm_returnsEmptySubmitFormForEmptyForm(c *C) {
	f := &data.Form{}
	f2, err := processForm(f, nil, func(title, instructions string, fields []interface{}) error {
		return nil
	})

	c.Assert(err, IsNil)
	c.Assert(*f2, DeepEquals, data.Form{Type: "submit"})
}

func (s *FormsXMPPSuite) Test_processForm_returnsFixedFields(c *C) {
	f := &data.Form{}
	f.Fields = []data.FormFieldX{
		data.FormFieldX{
			Label:  "hello",
			Type:   "fixed",
			Values: []string{"Something"},
		},
		data.FormFieldX{
			Label: "hello2",
			Type:  "fixed",
		},
	}
	f2, err := processForm(f, nil, func(title, instructions string, fields []interface{}) error {
		return nil
	})

	c.Assert(err, IsNil)
	c.Assert(*f2, DeepEquals, data.Form{
		XMLName:      xml.Name{Space: "", Local: ""},
		Type:         "submit",
		Title:        "",
		Instructions: "",
		Fields:       nil})
}

func (s *FormsXMPPSuite) Test_processForm_returnsBooleanFields(c *C) {
	f := &data.Form{}
	f.Fields = []data.FormFieldX{
		data.FormFieldX{
			Label: "hello3",
			Type:  "boolean",
		},
	}
	f2, err := processForm(f, nil, func(title, instructions string, fields []interface{}) error {
		return nil
	})

	c.Assert(err, IsNil)
	c.Assert(*f2, DeepEquals, data.Form{
		XMLName:      xml.Name{Space: "", Local: ""},
		Type:         "submit",
		Title:        "",
		Instructions: "",
		Fields: []data.FormFieldX{
			data.FormFieldX{
				XMLName:  xml.Name{Space: "", Local: ""},
				Desc:     "",
				Var:      "",
				Type:     "",
				Label:    "",
				Required: (*data.FormFieldRequiredX)(nil),
				Values:   []string{"false"},
				Options:  []data.FormFieldOptionX(nil),
				Media:    []data.FormFieldMediaX(nil)}}})
}

func (s *FormsXMPPSuite) Test_processForm_returnsMultiFields(c *C) {
	f := &data.Form{}
	f.Fields = []data.FormFieldX{
		data.FormFieldX{
			Label: "hello4",
			Type:  "jid-multi",
		},
		data.FormFieldX{
			Label: "hello5",
			Type:  "text-multi",
		},
	}
	f2, err := processForm(f, nil, func(title, instructions string, fields []interface{}) error {
		return nil
	})

	c.Assert(err, IsNil)
	c.Assert(*f2, DeepEquals, data.Form{
		XMLName:      xml.Name{Space: "", Local: ""},
		Type:         "submit",
		Title:        "",
		Instructions: "",
		Fields: []data.FormFieldX{
			data.FormFieldX{
				XMLName:  xml.Name{Space: "", Local: ""},
				Desc:     "",
				Var:      "",
				Type:     "",
				Label:    "",
				Required: (*data.FormFieldRequiredX)(nil),
				Values:   []string(nil),
				Options:  []data.FormFieldOptionX(nil),
				Media:    []data.FormFieldMediaX(nil)},
			data.FormFieldX{
				XMLName:  xml.Name{Space: "", Local: ""},
				Desc:     "",
				Var:      "",
				Type:     "",
				Label:    "",
				Required: (*data.FormFieldRequiredX)(nil),
				Values:   []string(nil),
				Options:  []data.FormFieldOptionX(nil),
				Media:    []data.FormFieldMediaX(nil)}}})
}

func (s *FormsXMPPSuite) Test_processForm_returnsListSingle(c *C) {
	f := &data.Form{}
	f.Fields = []data.FormFieldX{
		data.FormFieldX{
			Label: "hello7",
			Type:  "list-single",
			Options: []data.FormFieldOptionX{
				data.FormFieldOptionX{Label: "One", Value: "Two"},
				data.FormFieldOptionX{Label: "Three", Value: "Four"},
			},
		},
	}
	f2, err := processForm(f, nil, func(title, instructions string, fields []interface{}) error {
		return nil
	})

	c.Assert(err, IsNil)
	c.Assert(*f2, DeepEquals, data.Form{
		XMLName:      xml.Name{Space: "", Local: ""},
		Type:         "submit",
		Title:        "",
		Instructions: "",
		Fields: []data.FormFieldX{
			data.FormFieldX{
				XMLName:  xml.Name{Space: "", Local: ""},
				Desc:     "",
				Var:      "",
				Type:     "",
				Label:    "",
				Required: (*data.FormFieldRequiredX)(nil),
				Values:   []string{"Two"},
				Options:  []data.FormFieldOptionX(nil), Media: []data.FormFieldMediaX(nil)}}})
}

func (s *FormsXMPPSuite) Test_processForm_returnsListMulti(c *C) {
	f := &data.Form{}
	f.Fields = []data.FormFieldX{
		data.FormFieldX{
			Label: "hello1o7",
			Type:  "list-multi",
			Options: []data.FormFieldOptionX{
				data.FormFieldOptionX{Label: "One", Value: "Two"},
				data.FormFieldOptionX{Label: "Three", Value: "Four"},
			},
		},
	}
	f2, err := processForm(f, nil, func(title, instructions string, fields []interface{}) error {
		return nil
	})

	c.Assert(err, IsNil)
	c.Assert(*f2, DeepEquals, data.Form{
		XMLName:      xml.Name{Space: "", Local: ""},
		Type:         "submit",
		Title:        "",
		Instructions: "",
		Fields: []data.FormFieldX{
			data.FormFieldX{
				XMLName:  xml.Name{Space: "", Local: ""},
				Desc:     "",
				Var:      "",
				Type:     "",
				Label:    "",
				Required: (*data.FormFieldRequiredX)(nil),
				Values:   []string(nil),
				Options:  []data.FormFieldOptionX(nil),
				Media:    []data.FormFieldMediaX(nil)}}})
}

func (s *FormsXMPPSuite) Test_processForm_returnsHidden(c *C) {
	f := &data.Form{}
	f.Fields = []data.FormFieldX{
		data.FormFieldX{
			Label: "hello1o71",
			Type:  "hidden",
		},
	}
	f2, err := processForm(f, nil, func(title, instructions string, fields []interface{}) error {
		return nil
	})

	c.Assert(err, IsNil)
	c.Assert(*f2, DeepEquals, data.Form{
		XMLName:      xml.Name{Space: "", Local: ""},
		Type:         "submit",
		Title:        "",
		Instructions: "",
		Fields: []data.FormFieldX{
			data.FormFieldX{
				XMLName:  xml.Name{Space: "", Local: ""},
				Desc:     "",
				Var:      "",
				Type:     "",
				Label:    "",
				Required: (*data.FormFieldRequiredX)(nil),
				Values:   []string(nil),
				Options:  []data.FormFieldOptionX(nil),
				Media:    []data.FormFieldMediaX(nil)}}})
}

func (s *FormsXMPPSuite) Test_processForm_returnsUnknown(c *C) {
	f := &data.Form{}
	f.Fields = []data.FormFieldX{
		data.FormFieldX{
			Label: "hello1o71",
			Type:  "another-fancy-type",
		},
		data.FormFieldX{
			Label:  "hello1o73",
			Type:   "another-fancy-type",
			Values: []string{"another one"},
		},
	}
	f2, err := processForm(f, nil, func(title, instructions string, fields []interface{}) error {
		return nil
	})

	c.Assert(err, IsNil)
	c.Assert(*f2, DeepEquals, data.Form{
		XMLName:      xml.Name{Space: "", Local: ""},
		Type:         "submit",
		Title:        "",
		Instructions: "",
		Fields: []data.FormFieldX{
			data.FormFieldX{
				XMLName:  xml.Name{Space: "", Local: ""},
				Desc:     "",
				Var:      "",
				Type:     "",
				Label:    "",
				Required: (*data.FormFieldRequiredX)(nil),
				Values:   []string{""},
				Options:  []data.FormFieldOptionX(nil),
				Media:    []data.FormFieldMediaX(nil)},
			data.FormFieldX{
				XMLName:  xml.Name{Space: "", Local: ""},
				Desc:     "",
				Var:      "",
				Type:     "",
				Label:    "",
				Required: (*data.FormFieldRequiredX)(nil),
				Values:   []string{""},
				Options:  []data.FormFieldOptionX(nil),
				Media:    []data.FormFieldMediaX(nil)}}})
}

type testOtherFormType struct{}

func (s *FormsXMPPSuite) Test_processForm_panicsWhenGivenAWeirdFormType(c *C) {
	f := &data.Form{}
	f.Fields = []data.FormFieldX{
		data.FormFieldX{
			Label: "hello1o71",
			Type:  "another-fancy-type",
		},
	}
	c.Assert(func() {
		processForm(f, nil, func(title, instructions string, fields []interface{}) error {
			fields[0] = testOtherFormType{}
			return nil
		})
	}, PanicMatches, "unknown field type in result from callback: xmpp.testOtherFormType")
}

func (s *FormsXMPPSuite) Test_processForm_setsAValidBooleanReturnValue(c *C) {
	f := &data.Form{}
	f.Fields = []data.FormFieldX{
		data.FormFieldX{
			Label: "hello1o71",
			Type:  "boolean",
		},
	}
	f2, _ := processForm(f, nil, func(title, instructions string, fields []interface{}) error {
		fields[0].(*data.BooleanFormField).Result = true
		return nil
	})
	c.Assert(*f2, DeepEquals, data.Form{
		XMLName:      xml.Name{Space: "", Local: ""},
		Type:         "submit",
		Title:        "",
		Instructions: "",
		Fields: []data.FormFieldX{
			data.FormFieldX{
				XMLName:  xml.Name{Space: "", Local: ""},
				Desc:     "",
				Var:      "",
				Type:     "",
				Label:    "",
				Required: (*data.FormFieldRequiredX)(nil),
				Values:   []string{"true"},
				Options:  []data.FormFieldOptionX(nil),
				Media:    []data.FormFieldMediaX(nil)}}})
}

func (s *FormsXMPPSuite) Test_processForm_returnsListMultiWithResults(c *C) {
	f := &data.Form{}
	f.Fields = []data.FormFieldX{
		data.FormFieldX{
			Label: "hello1o7",
			Type:  "list-multi",
			Options: []data.FormFieldOptionX{
				data.FormFieldOptionX{Label: "One", Value: "Two"},
				data.FormFieldOptionX{Label: "Three", Value: "Four"},
			},
		},
	}
	f2, err := processForm(f, nil, func(title, instructions string, fields []interface{}) error {
		fields[0].(*data.MultiSelectionFormField).Results = []int{1}
		return nil
	})

	c.Assert(err, IsNil)
	c.Assert(*f2, DeepEquals, data.Form{
		XMLName:      xml.Name{Space: "", Local: ""},
		Type:         "submit",
		Title:        "",
		Instructions: "",
		Fields: []data.FormFieldX{
			data.FormFieldX{
				XMLName:  xml.Name{Space: "", Local: ""},
				Desc:     "",
				Var:      "",
				Type:     "",
				Label:    "",
				Required: (*data.FormFieldRequiredX)(nil),
				Values:   []string{"Four"},
				Options:  []data.FormFieldOptionX(nil), Media: []data.FormFieldMediaX(nil)}}})
}

func (s *FormsXMPPSuite) Test_processForm_dealsWithMediaCorrectly(c *C) {
	f := &data.Form{}
	datas := []data.BobData{
		data.BobData{
			CID:    "foobax",
			Base64: ".....",
		},
		data.BobData{
			CID:    "foobar",
			Base64: "aGVsbG8=",
		},
	}
	f.Fields = []data.FormFieldX{
		data.FormFieldX{
			Label: "hello1o7",
			Type:  "hidden",
			Media: []data.FormFieldMediaX{
				data.FormFieldMediaX{
					URIs: []data.MediaURIX{
						data.MediaURIX{
							MIMEType: "",
							URI:      "",
						},
						data.MediaURIX{
							MIMEType: "",
							URI:      "hello:world",
						},
						data.MediaURIX{
							MIMEType: "",
							URI:      "cid:foobar",
						},
						data.MediaURIX{
							MIMEType: "",
							URI:      "cid:foobax",
						},
					},
				},
			},
		},
	}
	f2, err := processForm(f, datas, func(title, instructions string, fields []interface{}) error {
		return nil
	})

	c.Assert(err, IsNil)
	c.Assert(*f2, DeepEquals, data.Form{
		XMLName:      xml.Name{Space: "", Local: ""},
		Type:         "submit",
		Title:        "",
		Instructions: "",
		Fields: []data.FormFieldX{
			data.FormFieldX{
				XMLName:  xml.Name{Space: "", Local: ""},
				Desc:     "",
				Var:      "",
				Type:     "",
				Label:    "",
				Required: nil,
				Values:   nil,
				Options:  nil,
				Media:    nil}}})
}
