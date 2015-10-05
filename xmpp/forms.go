// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xmpp implements the XMPP IM protocol, as specified in RFC 6120 and
// 6121.
package xmpp

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"strings"
)

type Form struct {
	XMLName      xml.Name    `xml:"jabber:x:data x"`
	Type         string      `xml:"type,attr"`
	Title        string      `xml:"title,omitempty"`
	Instructions string      `xml:"instructions,omitempty"`
	Fields       []formField `xml:"field"`
}

type formField struct {
	XMLName  xml.Name           `xml:"field"`
	Desc     string             `xml:"desc,omitempty"`
	Var      string             `xml:"var,attr"`
	Type     string             `xml:"type,attr,omitempty"`
	Label    string             `xml:"label,attr,omitempty"`
	Required *formFieldRequired `xml:"required"`
	Values   []string           `xml:"value"`
	Options  []formFieldOption  `xml:"option"`
	Media    []formFieldMedia   `xml:"media"`
}

type formFieldMedia struct {
	XMLName xml.Name   `xml:"urn:xmpp:media-element media"`
	URIs    []mediaURI `xml:"uri"`
}

type mediaURI struct {
	XMLName  xml.Name `xml:"urn:xmpp:media-element uri"`
	MIMEType string   `xml:"type,attr,omitempty"`
	URI      string   `xml:",chardata"`
}

type formFieldRequired struct {
	XMLName xml.Name `xml:"required"`
}

type formFieldOption struct {
	Label string `xml:"var,attr,omitempty"`
	Value string `xml:"value"`
}

// FormField is the type of a generic form field. One should type cast to a
// specific type of field before processing.
type FormField struct {
	// Label is a human readable label for this field.
	Label string
	// Type is the XMPP-internal type of this field. One should type cast
	// rather than inspect this.
	Type string
	// Name gives the internal name of the field.
	Name     string
	Required bool
	// Media contains one of more items of media associated with this
	// field and, for each item, one or more representations of it.
	Media [][]Media
}

type Media struct {
	MIMEType string
	// URI contains a URI to the data. It may be empty if Data is not.
	URI string
	// Data contains the raw data itself. It may be empty if URI is not.
	Data []byte
}

// FixedFormField is used to indicate a section heading. It's for the form to
// send data to the user rather than the other way around.
type FixedFormField struct {
	FormField

	Text string
}

// BooleanFormField is for a yes/no answer. The Result member should be set to
// the user's answer.
type BooleanFormField struct {
	FormField

	Result bool
}

// TextFormField is for the entry of a single textual item. The Result member
// should be set to the data entered.
type TextFormField struct {
	FormField

	Default string
	Result  string
	// Private is true if this is a password or other sensitive entry.
	Private bool
}

// MultiTextFormField is for the entry of a several textual items. The Results
// member should be set to the data entered.
type MultiTextFormField struct {
	FormField

	Defaults []string
	Results  []string
}

// SelectionFormField asks the user to pick a single element from a set of
// choices. The Result member should be set to an index of the Values array.
type SelectionFormField struct {
	FormField

	Values []string
	Ids    []string
	Result int
}

// MultiSelectionFormField asks the user to pick a subset of possible choices.
// The Result member should be set to a series of indexes of the Results array.
type MultiSelectionFormField struct {
	FormField

	Values  []string
	Ids     []string
	Results []int
}

// FormCallback is the type of a function called to process a form. The
// argument is a list of pointers to FormField types. The function should type
// cast the elements, prompt the user and fill in the result field in each
// struct.
type FormCallback func(title, instructions string, fields []interface{}) error

// processForm calls the callback with the given XMPP form and returns the
// result form. The datas argument contains any additional XEP-0231 blobs that
// might contain media for the questions in the form.
func processForm(form *Form, datas []bobData, callback FormCallback) (*Form, error) {
	var fields []interface{}

	for _, field := range form.Fields {
		base := FormField{
			Label:    field.Label,
			Type:     field.Type,
			Name:     field.Var,
			Required: field.Required != nil,
		}

		for _, media := range field.Media {
			var options []Media
			for _, uri := range media.URIs {
				media := Media{
					MIMEType: uri.MIMEType,
					URI:      uri.URI,
				}
				if strings.HasPrefix(media.URI, "cid:") {
					// cid URIs are references to data
					// blobs that, hopefully, were sent
					// along with the form.
					cid := media.URI[4:]
					media.URI = ""

					for _, data := range datas {
						if data.CID == cid {
							var err error
							if media.Data, err = base64.StdEncoding.DecodeString(data.Base64); err != nil {
								media.Data = nil
							}
						}
					}
				}
				if len(media.URI) > 0 || len(media.Data) > 0 {
					options = append(options, media)
				}
			}

			base.Media = append(base.Media, options)
		}

		switch field.Type {
		case "fixed":
			if len(field.Values) < 1 {
				continue
			}
			f := &FixedFormField{
				FormField: base,
				Text:      field.Values[0],
			}
			fields = append(fields, f)
		case "boolean":
			f := &BooleanFormField{
				FormField: base,
			}
			fields = append(fields, f)
		case "jid-multi", "text-multi":
			f := &MultiTextFormField{
				FormField: base,
				Defaults:  field.Values,
			}
			fields = append(fields, f)
		case "list-single":
			f := &SelectionFormField{
				FormField: base,
			}
			for _, opt := range field.Options {
				f.Ids = append(f.Ids, opt.Value)
				f.Values = append(f.Values, opt.Label)
			}
			fields = append(fields, f)
		case "list-multi":
			f := &MultiSelectionFormField{
				FormField: base,
			}
			for _, opt := range field.Options {
				f.Ids = append(f.Ids, opt.Value)
				f.Values = append(f.Values, opt.Label)
			}
			fields = append(fields, f)
		case "hidden":
			continue
		default:
			f := &TextFormField{
				FormField: base,
				Private:   field.Type == "text-private",
			}
			if len(field.Values) > 0 {
				f.Default = field.Values[0]
			}
			fields = append(fields, f)
		}
	}

	if err := callback(form.Title, form.Instructions, fields); err != nil {
		return nil, err
	}

	result := &Form{
		Type: "submit",
	}

	// Copy the hidden fields across.
	for _, field := range form.Fields {
		if field.Type != "hidden" {
			continue
		}
		result.Fields = append(result.Fields, formField{
			Var:    field.Var,
			Values: field.Values,
		})
	}

	for _, field := range fields {
		switch field := field.(type) {
		case *BooleanFormField:
			value := "false"
			if field.Result {
				value = "true"
			}
			result.Fields = append(result.Fields, formField{
				Var:    field.Name,
				Values: []string{value},
			})
		case *TextFormField:
			result.Fields = append(result.Fields, formField{
				Var:    field.Name,
				Values: []string{field.Result},
			})
		case *MultiTextFormField:
			result.Fields = append(result.Fields, formField{
				Var:    field.Name,
				Values: field.Results,
			})
		case *SelectionFormField:
			result.Fields = append(result.Fields, formField{
				Var:    field.Name,
				Values: []string{field.Ids[field.Result]},
			})
		case *MultiSelectionFormField:
			var values []string
			for _, selected := range field.Results {
				values = append(values, field.Ids[selected])
			}

			result.Fields = append(result.Fields, formField{
				Var:    field.Name,
				Values: values,
			})
		case *FixedFormField:
			continue
		default:
			panic(fmt.Sprintf("unknown field type in result from callback: %T", field))
		}
	}

	return result, nil
}
