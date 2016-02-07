// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xmpp implements the XMPP IM protocol, as specified in RFC 6120 and
// 6121.
package xmpp

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/twstrike/coyim/xmpp/data"
)

// FormCallback is the type of a function called to process a form. The
// argument is a list of pointers to FormField types. The function should type
// cast the elements, prompt the user and fill in the result field in each
// struct.
type FormCallback func(title, instructions string, fields []interface{}) error

// processForm calls the callback with the given XMPP form and returns the
// result form. The datas argument contains any additional XEP-0231 blobs that
// might contain media for the questions in the form.
func processForm(form *data.Form, datas []data.BobData, callback FormCallback) (*data.Form, error) {
	var fields []interface{}

	for _, field := range form.Fields {
		base := data.FormField{
			Label:    field.Label,
			Type:     field.Type,
			Name:     field.Var,
			Required: field.Required != nil,
		}

		for _, media := range field.Media {
			var options []data.Media
			for _, uri := range media.URIs {
				media := data.Media{
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
			f := &data.FixedFormField{
				FormField: base,
				Text:      field.Values[0],
			}
			fields = append(fields, f)
		case "boolean":
			f := &data.BooleanFormField{
				FormField: base,
			}
			fields = append(fields, f)
		case "jid-multi", "text-multi":
			f := &data.MultiTextFormField{
				FormField: base,
				Defaults:  field.Values,
			}
			fields = append(fields, f)
		case "list-single":
			f := &data.SelectionFormField{
				FormField: base,
			}
			for _, opt := range field.Options {
				f.Ids = append(f.Ids, opt.Value)
				f.Values = append(f.Values, opt.Label)
			}
			fields = append(fields, f)
		case "list-multi":
			f := &data.MultiSelectionFormField{
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
			f := &data.TextFormField{
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

	result := &data.Form{
		Type: "submit",
	}

	// Copy the hidden fields across.
	for _, field := range form.Fields {
		if field.Type != "hidden" {
			continue
		}
		result.Fields = append(result.Fields, data.FormFieldX{
			Var:    field.Var,
			Values: field.Values,
		})
	}

	for _, field := range fields {
		switch field := field.(type) {
		case *data.BooleanFormField:
			value := "false"
			if field.Result {
				value = "true"
			}
			result.Fields = append(result.Fields, data.FormFieldX{
				Var:    field.Name,
				Values: []string{value},
			})
		case *data.TextFormField:
			result.Fields = append(result.Fields, data.FormFieldX{
				Var:    field.Name,
				Values: []string{field.Result},
			})
		case *data.MultiTextFormField:
			result.Fields = append(result.Fields, data.FormFieldX{
				Var:    field.Name,
				Values: field.Results,
			})
		case *data.SelectionFormField:
			result.Fields = append(result.Fields, data.FormFieldX{
				Var:    field.Name,
				Values: []string{field.Ids[field.Result]},
			})
		case *data.MultiSelectionFormField:
			var values []string
			for _, selected := range field.Results {
				values = append(values, field.Ids[selected])
			}

			result.Fields = append(result.Fields, data.FormFieldX{
				Var:    field.Name,
				Values: values,
			})
		case *data.FixedFormField:
			continue
		default:
			panic(fmt.Sprintf("unknown field type in result from callback: %T", field))
		}
	}

	return result, nil
}
