package gui

import (
	log "github.com/sirupsen/logrus"

	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/gotk3adapter/gtki"
)

type formField struct {
	field    interface{}
	label    gtki.Label
	widget   gtki.Widget
	required gtki.Label
}

func buildWidgetsForFields(fields []interface{}) []formField {
	ret := make([]formField, 0, len(fields))

	for _, f := range fields {
		switch field := f.(type) {
		case *data.TextFormField:
			l, _ := g.gtk.LabelNew(field.Label)
			l.SetHAlign(gtki.ALIGN_START)
			l.SetSelectable(true)

			w, _ := g.gtk.EntryNew()
			w.SetText(field.Default)
			w.SetVisibility(!field.Private)

			li, _ := g.gtk.LabelNew("")
			if field.Required {
				li, _ = g.gtk.LabelNew("*")
			}

			ret = append(ret, formField{field, l, w, li})
		case *data.BooleanFormField:
			l, _ := g.gtk.LabelNew(field.Label)
			l.SetHAlign(gtki.ALIGN_START)
			l.SetSelectable(true)

			w, _ := g.gtk.CheckButtonNew()
			w.SetActive(field.Result)

			li, _ := g.gtk.LabelNew("")
			if field.Required {
				li, _ = g.gtk.LabelNew("*")
			}

			ret = append(ret, formField{field, l, w, li})
		case *data.SelectionFormField:
			l, _ := g.gtk.LabelNew(field.Label)
			l.SetHAlign(gtki.ALIGN_START)
			l.SetSelectable(true)

			w, _ := g.gtk.ComboBoxTextNew()
			for _, opt := range field.Ids {
				w.AppendText(opt)
			}

			w.SetActive(field.Result)

			li, _ := g.gtk.LabelNew("")
			if field.Required {
				li, _ = g.gtk.LabelNew("*")
			}

			ret = append(ret, formField{field, l, w, li})
		case *data.FixedFormField:
			l, _ := g.gtk.LabelNew(field.Label)

			w, _ := g.gtk.LabelNew(field.Text)
			w.SetHAlign(gtki.ALIGN_START)

			li, _ := g.gtk.LabelNew("")
			if field.Required {
				li, _ = g.gtk.LabelNew("*")
			}
			ret = append(ret, formField{field, l, w, li})
		case *data.Media:
			pb := getPixbufFromBytes(field.Data)
			w, _ := g.gtk.ImageNewFromPixbuf(pb)

			ret = append(ret, formField{field, nil, w, nil})
		case *data.CaptchaFormField:
			pb := getPixbufFromBytes(field.MediaForm.Data)
			wi, _ := g.gtk.ImageNewFromPixbuf(pb)

			ret = append(ret, formField{field.MediaForm, nil, wi, nil})

			l, _ := g.gtk.LabelNew(field.TextForm.Label)
			l.SetHAlign(gtki.ALIGN_START)
			l.SetSelectable(true)

			wt, _ := g.gtk.EntryNew()
			wt.SetText(field.TextForm.Default)
			wt.SetVisibility(!field.TextForm.Private)

			ret = append(ret, formField{field.TextForm, l, wt, nil})
		default:
			log.Println("Missing to implement form field of:", field)
		}
	}

	return ret
}
