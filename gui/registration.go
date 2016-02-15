package gui

import (
	"errors"
	"log"

	"github.com/gotk3/gotk3/gtk"
	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/xmpp/data"
	"github.com/twstrike/coyim/xmpp/interfaces"
)

var (
	errRegistrationAborted = errors.New("registration cancelled")
)

type registrationForm struct {
	parent gtk.IWindow

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
		w := field.widget.(*gtk.Entry)
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
	f.addFields(fields)

	builder := builderForDefinition("RegistrationForm")

	obj, _ := builder.GetObject("dialog")
	dialog := obj.(*gtk.Dialog)
	dialog.SetTitle(title)

	obj, _ = builder.GetObject("instructions")
	label := obj.(*gtk.Label)
	label.SetText(instructions)
	label.SetSelectable(true)

	obj, _ = builder.GetObject("grid")
	grid := obj.(*gtk.Grid)

	for i, field := range f.fields {
		grid.Attach(field.label, 0, i+1, 1, 1)
		grid.Attach(field.widget, 1, i+1, 1, 1)
	}
	grid.ShowAll()

	dialog.SetTransientFor(f.parent)

	wait := make(chan error)
	doInUIThread(func() {
		resp := gtk.ResponseType(dialog.Run())
		switch resp {
		case gtk.RESPONSE_APPLY:
			wait <- f.accepted()
		default:
			wait <- errRegistrationAborted
		}

		dialog.Destroy()
	})

	return <-wait
}

func requestAndRenderRegistrationForm(server string, formHandler data.FormCallback, saveFn func(), df func() interfaces.Dialer) error {
	policy := config.ConnectionPolicy{DialerFactory: df}

	//TODO: this would not be necessary if RegisterAccount did not use it
	conf := &config.Account{
		Account: "@" + server,
		Proxies: []string{"tor-auto://"},
	}

	//TODO: this should receive only a JID domainpart
	_, err := policy.RegisterAccount(formHandler, conf)

	if err != nil {
		//TODO: show something in the UI
		log.Println("Registration failed:", err)
		return err
	}

	go saveFn()

	return nil
}

type formField struct {
	field  interface{}
	label  *gtk.Label
	widget gtk.IWidget
}

func buildWidgetsForFields(fields []interface{}) []formField {
	ret := make([]formField, 0, len(fields))

	for _, f := range fields {
		switch field := f.(type) {
		case *data.TextFormField:
			//TODO: notify if it is required
			l, _ := gtk.LabelNew(field.Label)
			l.SetSelectable(true)

			w, _ := gtk.EntryNew()
			w.SetText(field.Default)
			w.SetVisibility(!field.Private)

			ret = append(ret, formField{field, l, w})
		default:
			log.Println("Missing to implement form field:", field)
		}
	}

	return ret
}
