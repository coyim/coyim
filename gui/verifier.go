package gui

import (
	"crypto/rand"
	"fmt"
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

	l withLog
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
		l:            conv.account,
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
			return conv.currentPeerForSending().(jid.WithResource)
		},
	}

	v.buildPinWindow()
	v.buildAnswerSMPDialog()
	// A function is used below because we cannot check whether a contact is verified
	// when newVerifier is called.
	v.buildUnverifiedWarning(func() bool {
		return conv.isEncrypted() && !conv.hasVerifiedKey()
	})
	v.buildWaitingForPeerNotification()
	v.buildPeerRequestsSMPNotification()
	v.buildSMPFailedDialog()

	return v
}

type pinWindow struct {
	b             *builder
	dialog        gtki.Dialog `gtk-widget:"dialog"`
	prompt        gtki.Label  `gtk-widget:"prompt"`
	pin           gtki.Label  `gtk-widget:"pin"`
	smpImage      gtki.Image  `gtk-widget:"smp_image"`
	padlockImage1 gtki.Image  `gtk-widget:"padlock_image1"`
	padlockImage2 gtki.Image  `gtk-widget:"padlock_image2"`
	alertImage    gtki.Image  `gtk-widget:"alert_image"`
}

func (v *verifier) buildPinWindow() {
	v.pinWindow = &pinWindow{
		b: newBuilder("GeneratePIN"),
	}

	panicOnDevError(v.pinWindow.b.bindObjects(v.pinWindow))

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

func (v *verifier) updateUnverifiedWarning() {
	v.unverifiedWarning.update()
}

func (v *verifier) hideUnverifiedWarning() {
	v.unverifiedWarning.infobar.Hide()
}

// TODO: check on linux
type unverifiedWarning struct {
	b                         *builder
	infobar                   gtki.Box    `gtk-widget:"verify-infobar"`
	closeInfobar              gtki.Box    `gtk-widget:"verify-close-infobar"`
	notification              gtki.Box    `gtk-widget:"verify-notification"`
	label                     gtki.Label  `gtk-widget:"verify-message"`
	image                     gtki.Image  `gtk-widget:"verify-image"`
	button                    gtki.Button `gtk-widget:"verify-button"`
	shouldShowVerificationBar func() bool
}

// TODO: how will this work with i18n?
var question = "Please enter the PIN that I shared with you."
var coyIMQuestion = regexp.MustCompile("Please enter the PIN that I shared with you.")

func (u *unverifiedWarning) update() {
	if u.shouldShowVerificationBar() {
		u.infobar.Show()
		u.label.Show()
		u.image.ShowAll()
	} else {
		u.infobar.Hide()
	}
}

func (v *verifier) buildUnverifiedWarning(shouldShowVerificationBar func() bool) {
	v.unverifiedWarning = &unverifiedWarning{
		b: newBuilder("UnverifiedWarning"),
	}

	v.unverifiedWarning.shouldShowVerificationBar = shouldShowVerificationBar
	panicOnDevError(v.unverifiedWarning.b.bindObjects(v.unverifiedWarning))

	v.unverifiedWarning.b.ConnectSignals(map[string]interface{}{
		"on_press_image": v.hideUnverifiedWarning,
	})

	prov := providerWithCSS("box { background-color: #fff3f3; color: #000000; border: 2px; }")
	updateWithStyle(v.unverifiedWarning.infobar, prov)

	prov = providerWithCSS("box { background-color: #e5d7d6; }")
	updateWithStyle(v.unverifiedWarning.closeInfobar, prov)

	setImageFromFile(v.unverifiedWarning.image, "warning.svg")
	v.unverifiedWarning.label.SetLabel(i18n.Local("Make sure no one else\nis reading your messages"))
	_, _ = v.unverifiedWarning.button.Connect("clicked", v.showPINDialog)

	v.notifier.notify(v.unverifiedWarning.infobar)
}

func (v *verifier) smpError(err error) {
	v.hideUnverifiedWarning()
	v.showCannotGeneratePINDialog(err)
}

func (v *verifier) showPINDialog() {
	v.unverifiedWarning.infobar.Hide()

	pin, err := v.createPIN()
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

func (v *verifier) createPIN() (string, error) {
	val, err := rand.Int(rand.Reader, big.NewInt(int64(1000000)))
	if err != nil {
		v.l.Log().WithError(err).Warn("Error encountered when creating a new PIN")
		return "", err
	}

	return fmt.Sprintf("%06d", val), err
}

type waitingForPeerNotification struct {
	b       *builder
	infobar gtki.InfoBar `gtk-widget:"smp-waiting-infobar"`
	label   gtki.Label   `gtk-widget:"smp-waiting-label"`
	image   gtki.Image   `gtk-widget:"smp-waiting-image"`
	button  gtki.Button  `gtk-widget:"smp-waiting-button"`
}

func (v *verifier) buildWaitingForPeerNotification() {
	v.waitingForPeer = &waitingForPeerNotification{
		b: newBuilder("WaitingSMPComplete"),
	}

	panicOnDevError(v.waitingForPeer.b.bindObjects(v.waitingForPeer))

	prov := providerWithCSS("box { background-color: #fff3f3; color: #000000; border: 2px; }")
	updateWithStyle(v.waitingForPeer.infobar, prov)

	v.waitingForPeer.label.SetText(i18n.Localf("Waiting for peer to finish \nsecuring the channel..."))
	setImageFromFile(v.waitingForPeer.image, "waiting.svg")

	_, _ = v.waitingForPeer.button.Connect("clicked", func() {
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
	v.l.Log().WithError(err).Warn("Cannot recover from error. Quitting SMP verification.")
	setImageFromFile(image, "failure.svg")
	_, _ = button.Connect("clicked", infobar.Destroy)

	infobar.ShowAll()

	v.notifier.notify(infobar)
}

type answerSMPWindow struct {
	b             *builder
	dialog        gtki.Dialog `gtk-widget:"dialog"`
	question      gtki.Label  `gtk-widget:"question_from_peer"`
	answer        gtki.Entry  `gtk-widget:"answer"`
	submitButton  gtki.Button `gtk-widget:"button_submit"`
	smpImage      gtki.Image  `gtk-widget:"smp_image"`
	padlockImage1 gtki.Image  `gtk-widget:"padlock_image1"`
	padlockImage2 gtki.Image  `gtk-widget:"padlock_image2"`
	alertImage    gtki.Image  `gtk-widget:"alert_image"`
}

func (v *verifier) buildAnswerSMPDialog() {
	v.answerSMPWindow = &answerSMPWindow{
		b: newBuilder("AnswerSMPQuestion"),
	}

	panicOnDevError(v.answerSMPWindow.b.bindObjects(v.answerSMPWindow))

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
	infobar      gtki.Box    `gtk-widget:"smp-requested-infobar"`
	closeInfobar gtki.Box    `gtk-widget:"smp-requested-close-infobar"`
	notification gtki.Box    `gtk-widget:"smp-requested-notification"`
	label        gtki.Label  `gtk-widget:"smp-requested-message"`
	image        gtki.Image  `gtk-widget:"smp-requested-image"`
	button       gtki.Button `gtk-widget:"smp-requested-button"`
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

	panicOnDevError(v.peerRequestsSMP.b.bindObjects(v.peerRequestsSMP))

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
	_, _ = v.peerRequestsSMP.button.Connect("clicked", func() {
		v.showAnswerSMPDialog(question)
	})

	v.peerRequestsSMP.show()
}

type verificationSuccessNotification struct {
	b      *builder
	dialog gtki.Dialog `gtk-widget:"verif-success-dialog"`
	label  gtki.Label  `gtk-widget:"verif-success-label"`
	image  gtki.Image  `gtk-widget:"verif-success-image"`
	button gtki.Button `gtk-widget:"verif-success-button"`
}

func (v *verifier) displayVerificationSuccess() {
	v.verificationSuccess = &verificationSuccessNotification{
		b: newBuilder("VerificationSucceeded"),
	}

	panicOnDevError(v.verificationSuccess.b.bindObjects(v.verificationSuccess))

	_, _ = v.verificationSuccess.button.Connect("clicked", v.verificationSuccess.dialog.Destroy)

	v.verificationSuccess.label.SetMarkup(i18n.Localf("Hooray! No one is listening to your conversations with <b>%s</b>", v.peerName()))
	setImageFromFile(v.verificationSuccess.image, "smpsuccess.svg")

	v.hideUnverifiedWarning()
	v.verificationSuccess.dialog.SetTransientFor(v.parentWindow)
	v.verificationSuccess.dialog.ShowAll()
}

type smpFailedNotification struct {
	dialog gtki.Dialog `gtk-widget:"verif-failure-dialog"`
	label  gtki.Label  `gtk-widget:"verif-failure-label"`
	button gtki.Button `gtk-widget:"verif-failure-button"`
}

// TODO: make this consistent
func (v *verifier) buildSMPFailedDialog() {
	builder := newBuilder("VerificationFailed")
	v.smpFailed = &smpFailedNotification{}
	panicOnDevError(builder.bindObjects(v.smpFailed))

	v.smpFailed.dialog.SetTransientFor(v.parentWindow)
	v.smpFailed.dialog.HideOnDelete()

	_, _ = v.smpFailed.dialog.Connect("response", func() {
		v.updateUnverifiedWarning()
		v.smpFailed.dialog.Hide()
	})
	_, _ = v.smpFailed.button.Connect("clicked", func() {
		v.updateUnverifiedWarning()
		v.smpFailed.dialog.Hide()
	})
}

func (v *verifier) displayVerificationFailure() {
	v.smpFailed.label.SetMarkup(i18n.Localf("We could not verify this channel with <b>%s</b>.", v.peerName()))
	v.smpFailed.dialog.ShowAll()
}

func (v *verifier) updateInProgressDialogs(encrypted bool) {
	if !encrypted {
		v.removeInProgressDialogs()
	}
}

func (v *verifier) removeInProgressDialogs() {
	v.peerRequestsSMP.infobar.Hide()
	v.waitingForPeer.infobar.Hide()
	v.pinWindow.dialog.Hide()
	v.answerSMPWindow.dialog.Hide()
}

func (v *verifier) cancelSMP() {
	v.removeInProgressDialogs()
	v.session.AbortSMP(v.peerToSendTo())
	v.updateUnverifiedWarning()
}
