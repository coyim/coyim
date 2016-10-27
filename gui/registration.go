package gui

import (
	"errors"
	"log"

	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/tls"
	"github.com/twstrike/coyim/xmpp/data"
	"github.com/twstrike/coyim/xmpp/interfaces"
	"github.com/twstrike/gotk3adapter/gtki"
)

var (
	errRegistrationAborted = errors.New("registration cancelled")
)

type registrationForm struct {
	parent gtki.Window

	server string
	conf   *config.Account
	fields []formField
}

func (f *registrationForm) accepted() error {
	conf, err := config.NewAccount()
	if err != nil {
		return err
	}

	//Find the fields we need to copy from the form to the account
	for _, field := range f.fields {
		ff := field.field.(*data.TextFormField)
		w := field.widget.(gtki.Entry)
		ff.Result, _ = w.GetText()

		switch ff.Label {
		case "User", "Username":
			conf.Account = ff.Result + "@" + f.server
		case "Password":
			conf.Password = ff.Result
		default:
			log.Println("Field", ff.Label)
		}
	}

	f.conf = conf
	return nil
}

func (f *registrationForm) addFields(fields []interface{}) {
	f.fields = buildWidgetsForFields(fields)
}

func (f *registrationForm) renderForm(title, instructions string, fields []interface{}) error {
	wait := make(chan error)
	doInUIThread(func() {
		f.addFields(fields)

		builder := newBuilder("RegistrationForm")

		obj := builder.getObj("dialog")
		dialog := obj.(gtki.Dialog)
		dialog.SetTitle(title)

		obj = builder.getObj("instructions")
		label := obj.(gtki.Label)
		label.SetText(instructions)
		label.SetSelectable(true)

		obj = builder.getObj("grid")
		grid := obj.(gtki.Grid)

		for i, field := range f.fields {
			grid.Attach(field.label, 0, i+1, 1, 1)
			grid.Attach(field.widget, 1, i+1, 1, 1)
		}
		grid.ShowAll()

		dialog.SetTransientFor(f.parent)

		resp := gtki.ResponseType(dialog.Run())
		switch resp {
		case gtki.RESPONSE_APPLY:
			wait <- f.accepted()
		default:
			wait <- errRegistrationAborted
		}

		dialog.Destroy()
	})

	return <-wait
}

func requestAndRenderRegistrationForm(server string, formHandler data.FormCallback, saveFn func(), errorFn func(error), df interfaces.DialerFactory, verifier tls.Verifier) {
	policy := config.ConnectionPolicy{DialerFactory: df}

	//TODO: this would not be necessary if RegisterAccount did not use it
	//TODO: we should give the choice of using Tor to the user
	conf := &config.Account{
		Account: "@" + server,
		Proxies: []string{"tor-auto://"},
	}

	//TODO: this should receive only a JID domainpart
	_, err := policy.RegisterAccount(formHandler, conf, verifier)

	if err != nil {
		errorFn(err)
		return
	}

	go saveFn()
}

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
			//TODO: notify if it is required
			l, _ := g.gtk.LabelNew(field.Label)
			l.SetSelectable(true)

			w, _ := g.gtk.EntryNew()
			w.SetText(field.Default)
			w.SetVisibility(!field.Private)

			ret = append(ret, formField{field, l, w})
		default:
			log.Println("Missing to implement form field:", field)
		}
	}

	return ret
}
