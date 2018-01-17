package gui

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"regexp"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/access"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
)

type verifier struct {
	parentWindow        gtki.Window
	session             access.Session
	notifier            *notifier
	pinWindow           *pinWindow
	answerSMPWindow     *answerSMPWindow
	smpFailed           *smpFailedNotification
	waitingForPeer      *waitingForPeerNotification
	peerRequestsSMP     *peerRequestsSMPNotification
	unverifiedWarning   *unverifiedWarning
	verificationSuccess *verificationSuccessNotification
	peerName            func() string
	peerToSendTo        func() jid.WithResource
}

type notifier struct {
	notificationArea gtki.Box
}

func (n *notifier) notify(i gtki.InfoBar) {
	n.notificationArea.Add(i)
}

// TODO: unify repeated stuff
func newVerifier(u *gtkUI, conv *conversationPane) *verifier {
	v := &verifier{
		parentWindow: conv.transientParent,
		session:      conv.account.session,
		notifier:     &notifier{conv.notificationArea},
		peerName: func() string {
			p, ok := conv.currentPeer()
			if !ok {
				return ""
			}
			return p.NameForPresentation()
		},
		peerToSendTo: func() jid.WithResource {
			return conv.peerToSendTo().(jid.WithResource)
		},
	}

	v.buildPinWindow()
	v.buildAnswerSMPDialog()
	// A function is used below because we cannot check whether a contact is verified
	// when newVerifier is called.
	v.buildUnverifiedWarning(func() bool {
		return conv.isVerified(u)
	})
	v.buildWaitingForPeerNotification()
	v.buildPeerRequestsSMPNotification()
	v.buildSMPFailedDialog()

	return v
}

type pinWindow struct {
	b             *builder
	dialog        gtki.Dialog
	prompt        gtki.Label
	pin           gtki.Label
	smpImage      gtki.Image
	padlockImage1 gtki.Image
	padlockImage2 gtki.Image
	alertImage    gtki.Image
}

func (v *verifier) buildPinWindow() {
	v.pinWindow = &pinWindow{
		b: newBuilder("GeneratePIN"),
	}

	v.pinWindow.b.getItems(
		"dialog", &v.pinWindow.dialog,
		"prompt", &v.pinWindow.prompt,
		"pin", &v.pinWindow.pin,
		"smp_image", &v.pinWindow.smpImage,
		"padlock_image1", &v.pinWindow.padlockImage1,
		"padlock_image2", &v.pinWindow.padlockImage2,
		"alert_image", &v.pinWindow.alertImage,
	)

	v.pinWindow.dialog.HideOnDelete()
	v.pinWindow.dialog.SetTransientFor(v.parentWindow)
	addBoldHeaderStyle(v.pinWindow.pin)

	v.pinWindow.b.ConnectSignals(map[string]interface{}{
		"close_share_pin": func() {
			v.showWaitingForPeerToCompleteSMPDialog()
			v.pinWindow.dialog.Hide()
		},
	})

	setImageFromFile(v.pinWindow.smpImage, "smp.svg")
	setImageFromFile(v.pinWindow.padlockImage1, "padlock.svg")
	setImageFromFile(v.pinWindow.padlockImage2, "padlock.svg")
	setImageFromFile(v.pinWindow.alertImage, "alert.svg")
}

func (v *verifier) showUnverifiedWarning() {
	v.unverifiedWarning.show()
}

// TODO: check on linux
type unverifiedWarning struct {
	b              *builder
	infobar        gtki.Box
	closeInfobar   gtki.Box
	notification   gtki.Box
	label          gtki.Label
	image          gtki.Image
	button         gtki.Button
	peerIsVerified func() bool
}

// TODO: how will this work with i18n?
var question = "Please enter the PIN that I shared with you."
var coyIMQuestion = regexp.MustCompile("Please enter the PIN that I shared with you.")

func (u *unverifiedWarning) show() {
	if !u.peerIsVerified() {
		u.infobar.Show()
		u.label.Show()
		u.image.ShowAll()
	} else {
		log.Println("We already have a peer and a trusted fingerprint. No reason to show the unverified warning")
	}
}

func (v *verifier) buildUnverifiedWarning(peerIsVerified func() bool) {
	v.unverifiedWarning = &unverifiedWarning{
		b: newBuilder("UnverifiedWarning"),
	}

	v.unverifiedWarning.peerIsVerified = peerIsVerified

	v.unverifiedWarning.b.getItems(
		"verify-infobar", &v.unverifiedWarning.infobar,
		"verify-close-infobar", &v.unverifiedWarning.closeInfobar,
		"verify-notification", &v.unverifiedWarning.notification,
		"verify-message", &v.unverifiedWarning.label,
		"verify-image", &v.unverifiedWarning.image,
		"verify-button", &v.unverifiedWarning.button,
	)

	v.unverifiedWarning.b.ConnectSignals(map[string]interface{}{
		"on_press_image": v.hideUnverifiedWarning,
	})

	prov := providerWithCSS("box { background-color: #fff3f3; color: #000000; border: 2px; }")
	updateWithStyle(v.unverifiedWarning.infobar, prov)

	prov = providerWithCSS("box { background-color: #e5d7d6; }")
	updateWithStyle(v.unverifiedWarning.closeInfobar, prov)

	setImageFromFile(v.unverifiedWarning.image, "warning.svg")
	v.unverifiedWarning.label.SetLabel(i18n.Local("Make sure no one else\nis reading your messages"))
	v.unverifiedWarning.button.Connect("clicked", v.showPINDialog)

	v.notifier.notify(v.unverifiedWarning.infobar)
}

func (v *verifier) smpError(err error) {
	v.hideUnverifiedWarning()
	v.showCannotGeneratePINDialog(err)
}

func (v *verifier) showPINDialog() {
	v.unverifiedWarning.infobar.Hide()

	pin, err := createPIN()
	if err != nil {
		v.pinWindow.dialog.Hide()
		v.smpError(err)
		return
	}
	v.pinWindow.pin.SetText(pin)
	v.pinWindow.prompt.SetMarkup(i18n.Localf("Share this one-time PIN with <b>%s</b>", v.peerName()))

	v.session.StartSMP(v.peerToSendTo(), question, pin)
	v.pinWindow.dialog.ShowAll()

	v.waitingForPeer.label.SetLabel(i18n.Localf("Waiting for peer to finish \nsecuring the channel..."))
	v.waitingForPeer.infobar.ShowAll()
}

func createPIN() (string, error) {
	val, err := rand.Int(rand.Reader, big.NewInt(int64(1000000)))
	if err != nil {
		log.Printf("Error encountered when creating a new PIN: %v", err)
		return "", err
	}

	return fmt.Sprintf("%06d", val), err
}

type waitingForPeerNotification struct {
	b       *builder
	infobar gtki.InfoBar
	label   gtki.Label
	image   gtki.Image
	button  gtki.Button
}

func (v *verifier) buildWaitingForPeerNotification() {
	v.waitingForPeer = &waitingForPeerNotification{
		b: newBuilder("WaitingSMPComplete"),
	}

	v.waitingForPeer.b.getItems(
		"smp-waiting-infobar", &v.waitingForPeer.infobar,
		"smp-waiting-label", &v.waitingForPeer.label,
		"smp-waiting-image", &v.waitingForPeer.image,
		"smp-waiting-button", &v.waitingForPeer.button,
	)

	prov := providerWithCSS("box { background-color: #fff3f3; color: #000000; border: 2px; }")
	updateWithStyle(v.waitingForPeer.infobar, prov)

	v.waitingForPeer.label.SetText(i18n.Localf("Waiting for peer to finish \nsecuring the channel..."))
	setImageFromFile(v.waitingForPeer.image, "waiting.svg")

	v.waitingForPeer.button.Connect("clicked", func() {
		v.cancelSMP()
	})

	v.notifier.notify(v.waitingForPeer.infobar)
}

func (v *verifier) showWaitingForPeerToCompleteSMPDialog() {
	v.waitingForPeer.label.SetLabel(i18n.Localf("Waiting for peer to finish \nsecuring the channel..."))
	v.hideUnverifiedWarning()
	v.waitingForPeer.infobar.ShowAll()
}

func (v *verifier) showCannotGeneratePINDialog(err error) {
	b := newBuilder("CannotVerifyWithSMP")

	infobar := b.getObj("smp-error-infobar").(gtki.InfoBar)
	label := b.getObj("smp-error-label").(gtki.Label)
	image := b.getObj("smp-error-image").(gtki.Image)
	button := b.getObj("smp-error-button").(gtki.Button)

	prov := providerWithCSS("box { background-color: #fff3f3; color: #000000; border: 2px; }")
	updateWithStyle(infobar, prov)

	label.SetText(i18n.Local("Unable to verify at this time."))
	log.Printf("Cannot recover from error: %v. Quitting SMP verification.", err)
	setImageFromFile(image, "failure.svg")
	button.Connect("clicked", infobar.Destroy)

	infobar.ShowAll()

	v.notifier.notify(infobar)
}

type answerSMPWindow struct {
	b             *builder
	dialog        gtki.Dialog
	question      gtki.Label
	answer        gtki.Entry
	submitButton  gtki.Button
	smpImage      gtki.Image
	padlockImage1 gtki.Image
	padlockImage2 gtki.Image
	alertImage    gtki.Image
}

func (v *verifier) buildAnswerSMPDialog() {
	v.answerSMPWindow = &answerSMPWindow{
		b: newBuilder("AnswerSMPQuestion"),
	}

	v.answerSMPWindow.b.getItems(
		"dialog", &v.answerSMPWindow.dialog,
		"question_from_peer", &v.answerSMPWindow.question,
		"button_submit", &v.answerSMPWindow.submitButton,
		"answer", &v.answerSMPWindow.answer,
		"smp_image", &v.answerSMPWindow.smpImage,
		"padlock_image1", &v.answerSMPWindow.padlockImage1,
		"padlock_image2", &v.answerSMPWindow.padlockImage2,
		"alert_image", &v.answerSMPWindow.alertImage,
	)

	v.answerSMPWindow.dialog.HideOnDelete()
	v.answerSMPWindow.dialog.SetTransientFor(v.parentWindow)
	v.answerSMPWindow.submitButton.SetSensitive(false)

	setImageFromFile(v.answerSMPWindow.smpImage, "smp.svg")
	setImageFromFile(v.answerSMPWindow.padlockImage1, "padlock.svg")
	setImageFromFile(v.answerSMPWindow.padlockImage2, "padlock.svg")
	setImageFromFile(v.answerSMPWindow.alertImage, "alert.svg")

	v.answerSMPWindow.b.ConnectSignals(map[string]interface{}{
		"text_changing": func() {
			answer, _ := v.answerSMPWindow.answer.GetText()
			v.answerSMPWindow.submitButton.SetSensitive(len(answer) > 0)
		},
		"close_share_pin": func() {
			answer, _ := v.answerSMPWindow.answer.GetText()
			v.removeInProgressDialogs()
			v.session.FinishSMP(v.peerToSendTo(), answer)
			v.showWaitingForPeerToCompleteSMPDialog()
		},
	})

}

func (v *verifier) showAnswerSMPDialog(question string) {
	if "" == question {
		v.answerSMPWindow.question.SetMarkup(i18n.Localf("Enter the secret that <b>%s</b> shared with you", v.peerName()))
	} else if coyIMQuestion.MatchString(question) {
		v.answerSMPWindow.question.SetMarkup(i18n.Localf("Type the PIN that <b>%s</b> sent you. It can be used only once.", v.peerName()))
	} else {
		v.answerSMPWindow.question.SetMarkup(i18n.Localf("Enter the answer to\n<b>%s</b>", question))
	}

	v.answerSMPWindow.answer.SetText("")
	v.answerSMPWindow.dialog.ShowAll()
}

type peerRequestsSMPNotification struct {
	b            *builder
	infobar      gtki.Box
	closeInfobar gtki.Box
	notification gtki.Box
	label        gtki.Label
	image        gtki.Image
	button       gtki.Button
}

func (p *peerRequestsSMPNotification) show() {
	p.infobar.Show()
	p.closeInfobar.Show()
	p.label.Show()
}

func (v *verifier) buildPeerRequestsSMPNotification() {
	v.peerRequestsSMP = &peerRequestsSMPNotification{
		b: newBuilder("PeerRequestsSMP"),
	}

	v.peerRequestsSMP.b.getItems(
		"smp-requested-infobar", &v.peerRequestsSMP.infobar,
		"smp-requested-close-infobar", &v.peerRequestsSMP.closeInfobar,
		"smp-requested-notification", &v.peerRequestsSMP.notification,
		"smp-requested-message", &v.peerRequestsSMP.label,
		"smp-requested-image", &v.peerRequestsSMP.image,
		"smp-requested-button", &v.peerRequestsSMP.button,
	)

	prov := providerWithCSS("box { background-color: #fff3f3; color: #000000; border: 2px; }")
	updateWithStyle(v.peerRequestsSMP.infobar, prov)

	prov = providerWithCSS("box { background-color: #e5d7d6; }")
	updateWithStyle(v.peerRequestsSMP.closeInfobar, prov)

	v.peerRequestsSMP.b.ConnectSignals(map[string]interface{}{
		"on_press_close_image": v.cancelSMP,
	})

	setImageFromFile(v.peerRequestsSMP.image, "waiting.svg")

	v.notifier.notify(v.peerRequestsSMP.infobar)
}

func (v *verifier) displayRequestForSecret(question string) {
	v.hideUnverifiedWarning()

	v.peerRequestsSMP.label.SetLabel(i18n.Localf("Finish verifying the \nsecurity of this channel..."))
	v.peerRequestsSMP.button.Connect("clicked", func() {
		v.showAnswerSMPDialog(question)
	})

	v.peerRequestsSMP.show()
}

type verificationSuccessNotification struct {
	b      *builder
	dialog gtki.Dialog
	label  gtki.Label
	image  gtki.Image
	button gtki.Button
}

func (v *verifier) displayVerificationSuccess() {
	v.verificationSuccess = &verificationSuccessNotification{
		b: newBuilder("VerificationSucceeded"),
	}

	v.verificationSuccess.b.getItems(
		"verif-success-dialog", &v.verificationSuccess.dialog,
		"verif-success-label", &v.verificationSuccess.label,
		"verif-success-image", &v.verificationSuccess.image,
		"verif-success-button", &v.verificationSuccess.button,
	)

	v.verificationSuccess.button.Connect("clicked", v.verificationSuccess.dialog.Destroy)

	v.verificationSuccess.label.SetMarkup(i18n.Localf("Hooray! No one is listening to your conversations with <b>%s</b>", v.peerName()))
	setImageFromFile(v.verificationSuccess.image, "smpsuccess.svg")

	v.hideUnverifiedWarning()
	v.verificationSuccess.dialog.SetTransientFor(v.parentWindow)
	v.verificationSuccess.dialog.ShowAll()
}

type smpFailedNotification struct {
	dialog gtki.Dialog
	label  gtki.Label
	button gtki.Button
}

// TODO: make this consistent
func (v *verifier) buildSMPFailedDialog() {
	builder := newBuilder("VerificationFailed")
	v.smpFailed = &smpFailedNotification{
		dialog: builder.getObj("verif-failure-dialog").(gtki.Dialog),
		label:  builder.getObj("verif-failure-label").(gtki.Label),
		button: builder.getObj("verif-failure-button").(gtki.Button),
	}

	v.smpFailed.dialog.SetTransientFor(v.parentWindow)
	v.smpFailed.dialog.HideOnDelete()

	v.smpFailed.dialog.Connect("response", func() {
		v.showUnverifiedWarning()
		v.smpFailed.dialog.Hide()
	})
	v.smpFailed.button.Connect("clicked", func() {
		v.showUnverifiedWarning()
		v.smpFailed.dialog.Hide()
	})
}

func (v *verifier) displayVerificationFailure() {
	v.smpFailed.label.SetMarkup(i18n.Localf("We could not verify this channel with <b>%s</b>.", v.peerName()))
	v.smpFailed.dialog.ShowAll()
}

func (v *verifier) removeInProgressDialogs() {
	v.peerRequestsSMP.infobar.Hide()
	v.waitingForPeer.infobar.Hide()
	v.pinWindow.dialog.Hide()
	v.answerSMPWindow.dialog.Hide()
}

func (v *verifier) hideUnverifiedWarning() {
	v.unverifiedWarning.infobar.Hide()
}

func (v *verifier) cancelSMP() {
	v.removeInProgressDialogs()
	v.session.AbortSMP(v.peerToSendTo())
	v.showUnverifiedWarning()
}
