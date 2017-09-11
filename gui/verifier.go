package gui

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"regexp"

	"github.com/twstrike/coyim/i18n"
	rosters "github.com/twstrike/coyim/roster"
	"github.com/twstrike/coyim/session/access"
	"github.com/twstrike/gotk3adapter/gtki"
)

type verifier struct {
	parentWindow                        gtki.Window
	currentResource                     string
	session                             access.Session
	notifier                            *notifier
	peerName                            string
	peerJid                             string
	pinWindow                           *pinWindow
	answerSMPWindow                     *answerSMPWindow
	smpFailed                           gtki.Dialog
	waitingForPeer                      *waitingForPeerNotification
	peerRequestsSMP                     *peerRequestsSMPNotification
	unverifiedWarning                   *unverifiedWarning
	verificationSuccess                 *verificationSuccessNotification
	showingUnverifiedWarning            bool
	widthBeforeShowingUnverifiedWarning int
}

type notifier struct {
	notificationArea gtki.Box
}

func (n *notifier) notify(i gtki.InfoBar) {
	n.notificationArea.Add(i)
}

func newVerifier(u *gtkUI, conv *conversationPane) *verifier {
	w, _ := conv.transientParent.GetSize()
	v := &verifier{
		parentWindow:    conv.transientParent,
		currentResource: conv.currentResource(),
		session:         conv.account.session,
		notifier:        &notifier{conv.notificationArea},
		peerName: conv.mapCurrentPeer("", func(p *rosters.Peer) string {
			return p.NameForPresentation()
		}),
		peerJid: conv.mapCurrentPeer("", func(p *rosters.Peer) string {
			return p.Jid
		}),
		showingUnverifiedWarning:            false,
		widthBeforeShowingUnverifiedWarning: w,
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
	d             gtki.Dialog
	prompt        gtki.Label
	pin           gtki.Label
	smpImage      gtki.Image
	padlockImage1 gtki.Image
	padlockImage2 gtki.Image
	alertImage    gtki.Image
}

func (v *verifier) buildPinWindow() {
	v.pinWindow = &pinWindow{b: newBuilder("GeneratePIN")}
	v.pinWindow.b.getItems(
		"dialog", &v.pinWindow.d,
		"prompt", &v.pinWindow.prompt,
		"pin", &v.pinWindow.pin,
		"smp_image", &v.pinWindow.smpImage,
		"padlock_image1", &v.pinWindow.padlockImage1,
		"padlock_image2", &v.pinWindow.padlockImage2,
		"alert_image", &v.pinWindow.alertImage,
	)
	v.pinWindow.d.HideOnDelete()
	v.pinWindow.d.SetTransientFor(v.parentWindow)
	v.pinWindow.prompt.SetText(i18n.Localf("Share the one-time PIN below with %s", v.peerName))
	addBoldHeaderStyle(v.pinWindow.pin)
	v.pinWindow.b.ConnectSignals(map[string]interface{}{
		"close_share_pin": func() {
			v.showWaitingForPeerToCompleteSMPDialog()
			v.pinWindow.d.Hide()
		},
	})
	setImageFromFile(v.pinWindow.smpImage, "smp.svg")
	setImageFromFile(v.pinWindow.padlockImage1, "padlock.svg")
	setImageFromFile(v.pinWindow.padlockImage2, "padlock.svg")
	setImageFromFile(v.pinWindow.alertImage, "alert.svg")
}

func (v *verifier) showUnverifiedWarning() {
	v.unverifiedWarning.show()
	v.showingUnverifiedWarning = true
	v.widthBeforeShowingUnverifiedWarning, _ = v.parentWindow.GetSize()
}

type unverifiedWarning struct {
	b                 *builder
	box               gtki.Box
	infobar           gtki.InfoBar
	msg               gtki.Label
	verifyButtonVert  gtki.Button
	verifyButtonHoriz gtki.Button
	alertImage        gtki.Image
	peerIsVerified    func() bool
}

func (u *unverifiedWarning) show() {
	if !u.peerIsVerified() {
		u.infobar.ShowAll()
	} else {
		log.Println("We have a peer and a trusted fingerprint already, so no reason to show the unverified warning")
	}
}

func (v *verifier) buildUnverifiedWarning(peerIsVerified func() bool) {
	v.unverifiedWarning = &unverifiedWarning{b: newBuilder("UnverifiedWarning")}
	v.unverifiedWarning.peerIsVerified = peerIsVerified
	v.unverifiedWarning.b.getItems(
		"infobar", &v.unverifiedWarning.infobar,
		"box", &v.unverifiedWarning.box,
		"message", &v.unverifiedWarning.msg,
		"button_verify_vertical", &v.unverifiedWarning.verifyButtonVert,
		"button_verify_horizontal", &v.unverifiedWarning.verifyButtonHoriz,
		"alert_image", &v.unverifiedWarning.alertImage,
	)
	setImageFromFile(v.unverifiedWarning.alertImage, "alert.svg")
	v.unverifiedWarning.b.ConnectSignals(map[string]interface{}{
		"close_verification": func() {
			v.hideUnverifiedWarning()
		},
	})
	v.unverifiedWarning.msg.SetText(i18n.Local("Make sure no one else is reading your messages"))
	v.unverifiedWarning.verifyButtonVert.Connect("clicked", v.showPINDialog)
	v.unverifiedWarning.verifyButtonHoriz.Connect("clicked", v.showPINDialog)
	v.notifier.notify(v.unverifiedWarning.infobar)
}

func (v *verifier) smpError(err error) {
	v.hideUnverifiedWarning()
	v.showNotificationWhenWeCannotGeneratePINForSMP(err)
}

func (v *verifier) showPINDialog() {
	pin, err := createPIN()
	if err != nil {
		v.pinWindow.d.Hide()
		v.smpError(err)
		return
	}
	v.pinWindow.pin.SetText(pin)
	v.session.StartSMP(v.peerJid, v.currentResource, question, pin)
	v.pinWindow.d.ShowAll()
}

func createPIN() (string, error) {
	val, err := rand.Int(rand.Reader, big.NewInt(int64(1000000)))
	if err != nil {
		log.Printf("Error encountered when creating a new PIN: %v", err)
		return "", err
	}
	return fmt.Sprintf("%06d", val), nil
}

type waitingForPeerNotification struct {
	b            *builder
	bar          gtki.InfoBar
	msg          gtki.Label
	cancelButton gtki.Button
}

func (v *verifier) buildWaitingForPeerNotification() {
	v.waitingForPeer = &waitingForPeerNotification{b: newBuilder("WaitingSMPComplete")}
	v.waitingForPeer.b.getItems(
		"smp_waiting_infobar", &v.waitingForPeer.bar,
		"message", &v.waitingForPeer.msg,
		"cancel_button", &v.waitingForPeer.cancelButton,
	)
	v.waitingForPeer.msg.SetText(i18n.Localf("Waiting for %s to finish securing the channel...", v.peerName))
	v.waitingForPeer.cancelButton.Connect("clicked", func() {
		v.session.AbortSMP(v.peerJid, v.currentResource)
		v.removeInProgressDialogs()
		v.showUnverifiedWarning()
	})
	v.notifier.notify(v.waitingForPeer.bar)
}

func (v *verifier) showWaitingForPeerToCompleteSMPDialog() {
	v.hideUnverifiedWarning()
	v.waitingForPeer.bar.ShowAll()
}

func (v *verifier) showNotificationWhenWeCannotGeneratePINForSMP(err error) {
	log.Printf("Cannot recover from error: %v. Quitting verification using SMP.", err)
	errBuilder := newBuilder("CannotVerifyWithSMP")
	errInfoBar := errBuilder.getObj("error_verifying_smp").(gtki.InfoBar)
	message := errBuilder.getObj("message").(gtki.Label)
	message.SetText(i18n.Local("Unable to verify the channel at this time."))
	button := errBuilder.getObj("try_later_button").(gtki.Button)
	button.Connect("clicked", errInfoBar.Destroy)
	errInfoBar.ShowAll()
	v.notifier.notify(errInfoBar)
}

type answerSMPWindow struct {
	b             *builder
	d             gtki.Dialog
	question      gtki.Label
	answer        gtki.Entry
	submitButton  gtki.Button
	smpImage      gtki.Image
	padlockImage1 gtki.Image
	padlockImage2 gtki.Image
	alertImage    gtki.Image
}

func (v *verifier) buildAnswerSMPDialog() {
	v.answerSMPWindow = &answerSMPWindow{b: newBuilder("AnswerSMPQuestion")}
	v.answerSMPWindow.b.getItems(
		"dialog", &v.answerSMPWindow.d,
		"question_from_peer", &v.answerSMPWindow.question,
		"button_submit", &v.answerSMPWindow.submitButton,
		"answer", &v.answerSMPWindow.answer,
		"smp_image", &v.answerSMPWindow.smpImage,
		"padlock_image1", &v.answerSMPWindow.padlockImage1,
		"padlock_image2", &v.answerSMPWindow.padlockImage2,
		"alert_image", &v.answerSMPWindow.alertImage,
	)
	v.answerSMPWindow.d.SetTransientFor(v.parentWindow)
	v.answerSMPWindow.d.HideOnDelete()
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
			v.showWaitingForPeerToCompleteSMPDialog()
			v.session.FinishSMP(v.peerJid, v.currentResource, answer)
		},
	})

}

var question = "Please enter the PIN that I shared with you."
var coyIMQuestion = regexp.MustCompile("Please enter the PIN that I shared with you.")

func (v *verifier) showAnswerSMPDialog(question string) {
	if "" == question {
		v.answerSMPWindow.question.SetText(i18n.Localf("Enter the secret that %s has shared with you", v.peerName))
	} else if coyIMQuestion.MatchString(question) {
		v.answerSMPWindow.question.SetText(i18n.Localf("Type the PIN that %s sent you. It can be used only once.", v.peerName))
	} else {
		v.answerSMPWindow.question.SetText(question)
	}
	v.answerSMPWindow.answer.SetText("")
	v.answerSMPWindow.d.ShowAll()
}

type peerRequestsSMPNotification struct {
	b                 *builder
	box               gtki.Box
	infobar           gtki.InfoBar
	msg               gtki.Label
	vertActionButtons gtki.Box
	verifyButtonVert  gtki.Button
	cancelButtonVert  gtki.Button
	verifyButtonHoriz gtki.Button
	cancelButtonHoriz gtki.Button
}

func (v *verifier) buildPeerRequestsSMPNotification() {
	v.peerRequestsSMP = &peerRequestsSMPNotification{b: newBuilder("PeerRequestsSMP")}
	v.peerRequestsSMP.b.getItems(
		"box", &v.peerRequestsSMP.box,
		"vert_action_buttons", &v.peerRequestsSMP.vertActionButtons,
		"peer_requests_smp", &v.peerRequestsSMP.infobar,
		"message", &v.peerRequestsSMP.msg,
		"verification_button_vertical", &v.peerRequestsSMP.verifyButtonVert,
		"cancel_button_vertical", &v.peerRequestsSMP.cancelButtonVert,
		"verification_button_horizontal", &v.peerRequestsSMP.verifyButtonHoriz,
		"cancel_button_horizontal", &v.peerRequestsSMP.cancelButtonHoriz,
	)
	v.peerRequestsSMP.cancelButtonVert.Connect("clicked", func() {
		v.removeInProgressDialogs()
		v.session.AbortSMP(v.peerJid, v.currentResource)
		v.unverifiedWarning.show()
	})
	v.peerRequestsSMP.cancelButtonHoriz.Connect("clicked", func() {
		v.removeInProgressDialogs()
		v.session.AbortSMP(v.peerJid, v.currentResource)
		v.showUnverifiedWarning()
	})

	v.peerRequestsSMP.msg.SetText(i18n.Localf("%s is waiting for you to finish verifying the security of this channel...", v.peerName))
	v.notifier.notify(v.peerRequestsSMP.infobar)
}

func (v *verifier) displayRequestForSecret(question string) {
	v.hideUnverifiedWarning()
	v.peerRequestsSMP.verifyButtonVert.Connect("clicked", func() {
		v.showAnswerSMPDialog(question)
	})
	v.peerRequestsSMP.verifyButtonHoriz.Connect("clicked", func() {
		v.showAnswerSMPDialog(question)
	})
	v.peerRequestsSMP.infobar.ShowAll()
}

type verificationSuccessNotification struct {
	b      *builder
	d      gtki.Dialog
	msg    gtki.Label
	img    gtki.Image
	button gtki.Button
}

func (v *verifier) displayVerificationSuccess() {
	v.verificationSuccess = &verificationSuccessNotification{b: newBuilder("VerificationSucceeded")}
	v.verificationSuccess.b.getItems(
		"dialog", &v.verificationSuccess.d,
		"verification_message", &v.verificationSuccess.msg,
		"success_image", &v.verificationSuccess.img,
		"button_ok", &v.verificationSuccess.button,
	)
	v.verificationSuccess.msg.SetText(i18n.Localf("Hooray! No one is listening in on your conversations with %s", v.peerName))
	v.verificationSuccess.button.Connect("clicked", v.verificationSuccess.d.Destroy)
	setImageFromFile(v.verificationSuccess.img, "smpsuccess.svg")

	v.verificationSuccess.d.SetTransientFor(v.parentWindow)
	v.verificationSuccess.d.ShowAll()

	v.hideUnverifiedWarning()
}

func (v *verifier) buildSMPFailedDialog() {
	builder := newBuilder("VerificationFailed")
	v.smpFailed = builder.getObj("dialog").(gtki.Dialog)
	v.smpFailed.SetTransientFor(v.parentWindow)
	v.smpFailed.HideOnDelete()
	addBoldHeaderStyle(builder.getObj("header").(gtki.Label))
	msg := builder.getObj("verification_message").(gtki.Label)
	msg.SetText(i18n.Localf("We could not verify this channel with %s.", v.peerName))
	tryLaterButton := builder.getObj("try_later").(gtki.Button)
	tryLaterButton.Connect("clicked", func() {
		v.showUnverifiedWarning()
		v.smpFailed.Hide()
	})
}

func (v *verifier) displayVerificationFailure() {
	v.smpFailed.ShowAll()
}

func (v *verifier) removeInProgressDialogs() {
	v.peerRequestsSMP.infobar.Hide()
	v.waitingForPeer.bar.Hide()
	v.pinWindow.d.Hide()
	v.answerSMPWindow.d.Hide()
}

func (v *verifier) hideUnverifiedWarning() {
	v.unverifiedWarning.infobar.Hide()
}

func (v *verifier) defaultLayoutForNotifications() {
	v.peerRequestsSMP.box.SetOrientation(gtki.VerticalOrientation)
	v.peerRequestsSMP.vertActionButtons.Show()
	v.peerRequestsSMP.verifyButtonHoriz.Hide()
	v.peerRequestsSMP.cancelButtonHoriz.Hide()
	addStyle(v.peerRequestsSMP.cancelButtonHoriz, "cancelButton", `.cancelButton {
		margin-left: 0.5em;
	}`)

	v.unverifiedWarning.box.SetOrientation(gtki.VerticalOrientation)
	v.unverifiedWarning.verifyButtonVert.Show()
	v.unverifiedWarning.verifyButtonHoriz.Hide()
	addStyle(v.unverifiedWarning.alertImage, "alert-icon", `.alert-icon {
		margin-left: 1em;
		margin-right: 0.5em;
	}`)
}

func (v *verifier) layoutForNotificationsInWiderWindow() {
	v.peerRequestsSMP.box.SetOrientation(gtki.HorizontalOrientation)
	v.peerRequestsSMP.vertActionButtons.Hide()
	v.peerRequestsSMP.verifyButtonHoriz.Show()
	v.peerRequestsSMP.cancelButtonHoriz.Show()
	addStyle(v.peerRequestsSMP.cancelButtonHoriz, "cancelButton", `.cancelButton {
		margin-left: 0.5em;
	}`)

	v.unverifiedWarning.box.SetOrientation(gtki.HorizontalOrientation)
	v.unverifiedWarning.verifyButtonVert.Hide()
	v.unverifiedWarning.verifyButtonHoriz.Show()
	addStyle(v.unverifiedWarning.alertImage, "alert-icon", `.alert-icon {
			margin-left: 3em;
			margin-right: 0.5em;
		}`)
}

var bestTransitionWidth = 600

func (v *verifier) chooseBestLayout(currentWindow gtki.Window) {
	w, l := currentWindow.GetSize()

	// This call to Resize() makes the first SMP notification display smoothly
	if v.showingUnverifiedWarning {
		currentWindow.Resize(v.widthBeforeShowingUnverifiedWarning, l)
		v.showingUnverifiedWarning = false
	}

	if w > bestTransitionWidth {
		v.layoutForNotificationsInWiderWindow()
	} else {
		v.defaultLayoutForNotifications()
	}
}
