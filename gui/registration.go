package gui

import (
	"errors"
	"log"

	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/i18n"
	"github.com/twstrike/coyim/servers"
	"github.com/twstrike/coyim/session"
	"github.com/twstrike/coyim/tls"
	"github.com/twstrike/coyim/xmpp"
	"github.com/twstrike/coyim/xmpp/data"
	"github.com/twstrike/coyim/xmpp/interfaces"
	"github.com/twstrike/gotk3adapter/gtki"
)

var (
	errRegistrationAborted = errors.New("registration cancelled")
)

type registrationForm struct {
	grid gtki.Grid

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

func (f *registrationForm) renderForm(title string, fields []interface{}) error {
	doInUIThread(func() {
		f.addFields(fields)

		for i, field := range f.fields {
			f.grid.Attach(field.label, 0, i+1, 1, 1)
			f.grid.Attach(field.widget, 1, i+1, 1, 1)
		}

		f.grid.ShowAll()
	})

	return nil
}

func requestAndRenderRegistrationForm(server string, formHandler data.FormCallback, df interfaces.DialerFactory, verifier tls.Verifier, c *config.ApplicationConfig) error {
	_, xmppLog := session.CreateXMPPLogger(c.RawLogFile)
	policy := config.ConnectionPolicy{DialerFactory: df, XMPPLogger: xmppLog, Logger: session.LogToDebugLog()}

	//TODO: this would not be necessary if RegisterAccount did not use it
	//TODO: we should give the choice of using Tor to the user
	conf := &config.Account{
		Account: "@" + server,
		Proxies: []string{"tor-auto://"},
	}

	//TODO: this should receive only a JID domainpart
	conn, err := policy.RegisterAccount(formHandler, conf, verifier)
	if conn != nil {
		defer conn.Close()
	}

	return err
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

const (
	torErrorMessage = "The registration process currently requires Tor in order to ensure your safety.\n\n" +
		"You don't have Tor running. Please, start it.\n\n"
	torLogMessage         = "We had an error when trying to register your account: Tor is not running. %v"
	storeAccountInfoError = "We had an error when trying to store your account information."
	storeAccountInfoLog   = "We had an error when trying to store your account information. %v"
	contactServerError    = "Could not contact the server.\n\n Please correct your server choice and try again."
	contactServerLog      = "Error when trying to get registration form: %v"
	requiredFieldsError   = "We had an error:\n\nSome required fields are missing."
	requiredFieldsLog     = "Error when trying to get registration form: %v"
)

// TODO: check rendering of images
func renderError(doneMessage gtki.Label, errorMessage, logMessage string, err error) {
	log.Printf(logMessage, err)
	//doneImage.SetFromIconName("software-update-urgent", gtki.ICON_SIZE_DIALOG)
	doneMessage.SetLabel(i18n.Local(errorMessage))
}

// TODO: currently this shows up even when Tor is running but either:
// a timeout occured, xmpp: account creation failed, or could not authenticate to the XMPP server
func renderTorError(assistant gtki.Assistant, pg gtki.Widget, formMessage gtki.Label, err error) {
	log.Printf(torLogMessage, err)
	assistant.SetPageType(pg, gtki.ASSISTANT_PAGE_SUMMARY)
	formMessage.SetLabel(i18n.Local(torErrorMessage))
	//formImage.Clear()
	//formImage.SetFromIconName("software-update-urgent", gtki.ICON_SIZE_DIALOG)
}

func (w *serverSelectionWindow) renderErrorFor(err error) {
	if err != xmpp.ErrMissingRequiredRegistrationInfo {
		renderError(w.doneMessage, contactServerError, contactServerLog, err)
	} else {
		renderError(w.doneMessage, requiredFieldsError, requiredFieldsLog, err)
	}
}

type serverSelectionWindow struct {
	b           *builder
	assistant   gtki.Assistant
	formMessage gtki.Label
	doneMessage gtki.Label
	serverBox   gtki.ComboBoxText
	spinner     gtki.Spinner
	grid        gtki.Grid
	// formImage := builder.getObj("formImage").(gtki.Image)
	// doneImage := builder.getObj("doneImage").(gtki.Image)

	formSubmitted chan error
	done          chan error

	form *registrationForm

	u *gtkUI
}

func createServerSelectionWindow(u *gtkUI) *serverSelectionWindow {
	w := &serverSelectionWindow{b: newBuilder("AccountRegistration"), u: u}

	w.b.getItems(
		"assistant", &w.assistant,
		"formMessage", &w.formMessage,
		"doneMessage", &w.doneMessage,
		"server", &w.serverBox,
		"spinner", &w.spinner,
		"formGrid", &w.grid,
	)

	w.assistant.SetTransientFor(u.window)

	w.formSubmitted = make(chan error)
	w.done = make(chan error)

	w.form = &registrationForm{grid: w.grid}

	return w
}

func (w *serverSelectionWindow) initializeServers() {
	for _, s := range servers.GetServersForRegistration() {
		w.serverBox.AppendText(s.Name)
	}
	w.serverBox.SetActive(0)
}

func (w *serverSelectionWindow) initialPage() {
	w.serverBox.SetSensitive(true)
	w.form.server = ""

	//TODO: Destroy everything in the grid on page 1?
}

func (w *serverSelectionWindow) renderForm(pg gtki.Widget) func(string, string, []interface{}) error {
	return func(title, instructions string, fields []interface{}) error {
		w.spinner.Stop()
		w.formMessage.SetLabel("")
		w.doneMessage.SetLabel("")

		w.form.renderForm(title, fields)
		w.assistant.SetPageComplete(pg, true)

		return <-w.formSubmitted
	}
}

func (w *serverSelectionWindow) doRendering(pg gtki.Widget) {
	err := requestAndRenderRegistrationForm(w.form.server, w.renderForm(pg), w.u.dialerFactory, w.u.unassociatedVerifier(), w.u.config)
	if err != nil && w.assistant.GetCurrentPage() != 2 {
		if err != config.ErrTorNotRunning {
			go w.assistant.SetCurrentPage(2)
		}
		w.spinner.Stop()
		renderTorError(w.assistant, pg, w.formMessage, err)
		return
	}

	w.done <- err
}

func (w *serverSelectionWindow) serverChosenPage(pg gtki.Widget) {
	w.serverBox.SetSensitive(false)
	w.form.server = w.serverBox.GetActiveText()
	w.spinner.Start()
	w.formMessage.SetLabel(i18n.Local("Connecting to server for registration... \n\n" +
		"This might take a while."))

	go w.doRendering(pg)
}

func (w *serverSelectionWindow) formSubmittedPage() {
	w.formSubmitted <- w.form.accepted()
	err := <-w.done
	w.spinner.Stop()

	if err != nil {
		w.renderErrorFor(err)
		return
	}

	//Save the account
	err = w.u.addAndSaveAccountConfig(w.form.conf)

	if err != nil {
		renderError(w.doneMessage, storeAccountInfoError, storeAccountInfoLog, err)
		return
	}

	if acc, ok := w.u.getAccountByID(w.form.conf.ID()); ok {
		acc.session.SetWantToBeOnline(true)
		acc.Connect()
	}

	// doneImage.SetFromIconName("emblem-default", gtki.ICON_SIZE_DIALOG)
	w.doneMessage.SetLabel(i18n.Localf("%s successfully created.", w.form.conf.Account))
}

func (w *serverSelectionWindow) onPageChange(_ gtki.Assistant, pg gtki.Widget) {
	switch w.assistant.GetCurrentPage() {
	case 0:
		w.initialPage()
	case 1:
		w.serverChosenPage(pg)
	case 2:
		w.formSubmittedPage()
	}
}

func (u *gtkUI) showServerSelectionWindow() {
	w := createServerSelectionWindow(u)
	w.initializeServers()

	w.b.ConnectSignals(map[string]interface{}{
		"on_prepare":       w.onPageChange,
		"on_cancel_signal": w.assistant.Destroy,
	})

	w.assistant.ShowAll()
}
