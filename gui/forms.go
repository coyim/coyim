package gui

import (
	"log"

	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/gotk3adapter/gtki"
)

type formField struct {
	field  interface{}
	label  gtki.Label
	widget gtki.Widget
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

			ret = append(ret, formField{field, l, w})
		case *data.BooleanFormField:
			//TODO: notify if it is required
			l, _ := g.gtk.LabelNew(field.Label)
			l.SetHAlign(gtki.ALIGN_START)
			l.SetSelectable(true)

			w, _ := g.gtk.CheckButtonNew()
			w.SetActive(field.Result)

			ret = append(ret, formField{field, l, w})
		case *data.SelectionFormField:
			//TODO: notify if it is required
			l, _ := g.gtk.LabelNew(field.Label)
			l.SetHAlign(gtki.ALIGN_START)
			l.SetSelectable(true)

			w, _ := g.gtk.ComboBoxTextNew()
			for _, opt := range field.Ids {
				w.AppendText(opt)
			}

			w.SetActive(field.Result)

			ret = append(ret, formField{field, l, w})
		case *data.FixedFormField:
			l, _ := g.gtk.LabelNew(field.Label)

			w, _ := g.gtk.LabelNew(field.Text)
			w.SetHAlign(gtki.ALIGN_START)

			ret = append(ret, formField{field, l, w})
		default:
			log.Println("Missing to implement form field:", field)
		}
	}

	return ret
}
