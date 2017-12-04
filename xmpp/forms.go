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

	"github.com/coyim/coyim/xmpp/data"
)

func mediaForPresentation(field data.FormFieldX, datas []data.BobData) [][]data.Media {
	if len(field.Media) == 0 {
		return nil
	}

	ret := make([][]data.Media, 0, len(field.Media))

	for _, media := range field.Media {
		options := make([]data.Media, 0, len(media.URIs))
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

		ret = append(ret, options)
	}

	return ret
}

func toFormField(field data.FormFieldX, media [][]data.Media) interface{} {
	base := data.FormField{
		Label:    field.Label,
		Type:     field.Type,
		Name:     field.Var,
		Required: field.Required != nil,
		Media:    media,
	}

	switch field.Type {
	case "fixed":
		if len(field.Values) < 1 {
			return nil
		}
		f := &data.FixedFormField{
			FormField: base,
			Text:      field.Values[0],
		}
		return f
	case "boolean":
		result := false
		if len(field.Values) > 0 {
			//See: XEP-0040, Appendix G, item 10.
			result = field.Values[0] == "true" || field.Values[0] == "1"
		}

		f := &data.BooleanFormField{
			FormField: base,
			Result:    result,
		}
		return f
	case "jid-multi", "text-multi":
		f := &data.MultiTextFormField{
			FormField: base,
			Defaults:  field.Values,
		}
		return f
	case "list-single":
		f := &data.SelectionFormField{
			FormField: base,
		}

		for i, opt := range field.Options {
			f.Ids = append(f.Ids, opt.Value)
			f.Values = append(f.Values, opt.Label)

			if field.Values[0] == opt.Value {
				f.Result = i
			}
		}
		return f
	case "list-multi":
		f := &data.MultiSelectionFormField{
			FormField: base,
		}
		for i, opt := range field.Options {
			f.Ids = append(f.Ids, opt.Value)
			f.Values = append(f.Values, opt.Label)

			if len(f.Results) < len(field.Values) {
				for _, v := range field.Values {
					if v == opt.Value {
						f.Results = append(f.Results, i)
						break
					}
				}
			}
		}
		return f
	case "hidden":
		return nil
	default:
		f := &data.TextFormField{
			FormField: base,
			Private:   field.Type == "text-private",
		}
		if len(field.Values) > 0 {
			f.Default = field.Values[0]
		}
		return f
	}

	return nil
}

func toFormFieldX(field interface{}) *data.FormFieldX {
	switch field := field.(type) {
	case *data.BooleanFormField:
		value := "false"
		if field.Result {
			value = "true"
		}
		return &data.FormFieldX{
			Var:    field.Name,
			Values: []string{value},
		}
	case *data.TextFormField:
		return &data.FormFieldX{
			Var:    field.Name,
			Values: []string{field.Result},
		}
	case *data.MultiTextFormField:
		return &data.FormFieldX{
			Var:    field.Name,
			Values: field.Results,
		}
	case *data.SelectionFormField:
		return &data.FormFieldX{
			Var:    field.Name,
			Values: []string{field.Ids[field.Result]},
		}
	case *data.MultiSelectionFormField:
		var values []string
		for _, selected := range field.Results {
			values = append(values, field.Ids[selected])
		}

		return &data.FormFieldX{
			Var:    field.Name,
			Values: values,
		}
	case *data.FixedFormField:
		return nil
	default:
		panic(fmt.Sprintf("unknown field type in result from callback: %T", field))
	}

	return nil
}

// processForm calls the callback with the given XMPP form and returns the
// result form. The datas argument contains any additional XEP-0231 blobs that
// might contain media for the questions in the form.
func processForm(form *data.Form, datas []data.BobData, callback data.FormCallback) (*data.Form, error) {
	fields := make([]interface{}, 0, len(form.Fields))
	result := &data.Form{
		Type: "submit",
	}

	for _, field := range form.Fields {
		// Copy the hidden fields across.
		if field.Type == "hidden" {
			//skipping hidden fields has a consequence of not processing their media
			result.Fields = append(result.Fields, data.FormFieldX{
				Var:    field.Var,
				Values: field.Values,
			})
			continue
		}

		media := mediaForPresentation(field, datas)
		formField := toFormField(field, media)
		if formField != nil {
			fields = append(fields, formField)
		}
	}

	if err := callback(form.Title, form.Instructions, fields); err != nil {
		return nil, err
	}

	for _, field := range fields {
		formFieldX := toFormFieldX(field)
		if formFieldX != nil {
			result.Fields = append(result.Fields, *formFieldX)
		}
	}

	return result, nil
}
