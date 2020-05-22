package gui

import (
	"errors"
	"fmt"
	"log"

	"github.com/coyim/coyim/config"
	"github.com/coyim/coyim/i18n"
	ourNet "github.com/coyim/coyim/net"
	"github.com/coyim/coyim/servers"
	"github.com/coyim/coyim/session"
	"github.com/coyim/coyim/tls"
	"github.com/coyim/coyim/xmpp"
	"github.com/coyim/coyim/xmpp/data"
	xmppErr "github.com/coyim/coyim/xmpp/errors"
	"github.com/coyim/coyim/xmpp/interfaces"
	"github.com/coyim/gotk3adapter/gtki"
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
		switch ff := field.field.(type) {
		case *data.FixedFormField:
			switch ff.Name {
			case "captcha-fallback-text":
				log.Printf("Captcha fallback text %s", ff.Label)
			default:
				log.Printf("Field %s", ff.Label)
			}
		case *data.TextFormField:
			w := field.widget.(gtki.Entry)
			ff.Result, _ = w.GetText()

			switch ff.Label {
			case "User", "Username":
				conf.Account = ff.Result + "@" + f.server
			case "Password":
				conf.Password = ff.Result
			case "Enter the text you see":
				conf.Password = ff.Result
			default:
				log.Printf("Field %s", ff.Label)
			}
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

		var i int
		for _, field := range f.fields {
			f.grid.Attach(field.label, 0, i+1, 1, 1)
			f.grid.Attach(field.widget, 1, i+1, 1, 1)
			f.grid.Attach(field.required, 2, i+1, 1, 1)
			i++
		}

		li, _ := g.gtk.LabelNew("* The field is required.")
		f.grid.Attach(li, 0, i+1, 1, i+1)

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

const (
	torErrorMessage = "The registration process currently requires Tor in order to ensure your safety.\n\n" +
		"You don't have Tor running. Please, start it.\n\n"
	torLogMessage            = "We had an error when trying to register your account: Tor is not running. %v"
	storeAccountInfoError    = "We had an error when trying to store your account information."
	storeAccountInfoLog      = "We had an error when trying to store your account information. Please, try again.%v"
	contactServerError       = "Could not contact the server.\n\nPlease, correct your server choice and try again."
	contactServerLog         = "Error when trying to get registration form: %v"
	timeOutError             = "We had an error:\n\nTimeout."
	timeOutLog               = "Error when trying to get registration form: %v"
	requiredFieldsError      = "We had an error:\n\nSome required fields are missing. Please, try again and fill all fields."
	requiredFieldsLog        = "Error when trying to get registration form: %v"
	wrongCaptchaError        = "We had an error:\n\nThe captcha entered is wrong"
	wrongCaptchaLog          = "We had an error when trying to create your account: %v"
	conflictingUserNameError = "We had an error:\n\nIncorrect Username"
	conflictingUserNameLog   = "We had an error when trying to create your account: %v"
	resourceConstraintError  = "We had an error:\n\ntoo many requests for creating account."
	resourceConstraintLog    = "We had an error when trying to create your account: %v"
)

func renderError(message gtki.Label, errorMessage, logMessage string, err error) {
	log.Printf(logMessage, err)
	message.SetLabel(i18n.Local(errorMessage))
}

func (w *serverSelectionWindow) renderConnectionErrorFor(err error) {
	w.spinner.Stop()
	setImageFromFile(w.formImage, "failure.svg")

	switch err {

	case ourNet.ErrTimeout:
		renderError(w.formMessage, timeOutError, timeOutLog, err)
	case config.ErrTorNotRunning:
		renderError(w.formMessage, torErrorMessage, torLogMessage, err)
	default:
		renderError(w.formMessage, contactServerError, contactServerLog, err)
	}
}

func (w *serverSelectionWindow) renderErrorFor(err error) {
	setImageFromFile(w.doneImage, "failure.svg")

	switch err {
	case xmpp.ErrMissingRequiredRegistrationInfo:
		renderError(w.doneMessage, requiredFieldsError, requiredFieldsLog, err)
	case xmpp.ErrUsernameConflict:
		renderError(w.doneMessage, conflictingUserNameError, conflictingUserNameLog, err)
	case xmpp.ErrWrongCaptcha:
		renderError(w.doneMessage, wrongCaptchaError, wrongCaptchaLog, err)
	case xmpp.ErrResourceConstraint:
		renderError(w.doneMessage, resourceConstraintError, resourceConstraintLog, err)
	default:
		renderError(w.doneMessage, contactServerError, contactServerLog, err)
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
	formImage   gtki.Image
	doneImage   gtki.Image

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
		"formImage", &w.formImage,
		"doneImage", &w.doneImage,
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

	w.grid.Hide()
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
	if err != nil {
		// TODO: refactor me!
		if err == config.ErrTorNotRunning || err == xmppErr.ErrAuthenticationFailed || err == xmpp.ErrRegistrationFailed || err == ourNet.ErrTimeout {
			w.assistant.SetPageType(pg, gtki.ASSISTANT_PAGE_SUMMARY)
			w.assistant.SetPageComplete(pg, true)
			w.renderConnectionErrorFor(err)
			return
		}
	}

	go w.assistant.SetCurrentPage(2)

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
	w.grid.Show()
	w.formSubmitted <- w.form.accepted()

	err := <-w.done
	w.spinner.Stop()

	if err != nil {
		w.renderErrorFor(err)
		return
	}

	// Save the account
	err = w.u.addAndSaveAccountConfig(w.form.conf)

	if err != nil {
		renderError(w.doneMessage, storeAccountInfoError, storeAccountInfoLog, err)
		return
	}

	if acc, ok := w.u.getAccountByID(w.form.conf.ID()); ok {
		acc.session.SetWantToBeOnline(true)
		acc.Connect()
	}

	setImageFromFile(w.doneImage, "success.svg")
	w.doneMessage.SetMarkup(i18n.Localf("<b>%s</b> successfully created.", w.form.conf.Account))
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

	entry := w.serverBox.GetChild().(gtki.Widget)
	entry.Connect("activate", func() {
		w.assistant.SetCurrentPage(1)
	})

	w.assistant.ShowAll()
}

func handleActiveSignal() {
	fmt.Printf("Activate signal connected")
}
