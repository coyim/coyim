package gui

import (
	"crypto/x509"
	"fmt"
	"strings"
	"time"

	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/gtki"
	"github.com/twstrike/coyim/digests"
)

func (u *gtkUI) certificateFailedToVerify(a *account, certs []*x509.Certificate) <-chan bool {
	c := make(chan bool)

	u.certificateFailedToVerifyDisplayDialog(a, certs, c)

	return c
}

const chunkingDefaultGrouping = 8

func splitStringEvery(s string, n int) []string {
	result := []string{}
	str := s
	for len(str) > 0 {
		result = append(result, str[0:n])
		str = str[n:]
	}
	return result
}

func join20(s []string) string {
	return joinOther(s)
}

func join32(s []string) string {
	return fmt.Sprintf("%s %s %s %s  %s %s %s %s", s[0], s[1], s[2], s[3], s[4], s[5], s[6], s[7])
}

func joinOther(s []string) string {
	return strings.Join(s, " ")
}

func displayChunked(fpr []byte) string {
	str := fmt.Sprintf("%X", fpr)
	spl := splitStringEvery(str, chunkingDefaultGrouping)
	switch len(fpr) {
	case 20:
		return join20(spl)
	case 32:
		return join32(spl)
	default:
		return joinOther(spl)
	}
}

func (u *gtkUI) certificateFailedToVerifyDisplayDialog(a *account, certs []*x509.Certificate, c chan<- bool) {
	doInUIThread(func() {
		builder := newBuilder("CertificateDialog")

		var md gtki.Dialog
		var message gtki.Label
		var issuedToCN, issuedToO, issuedToOU, serial gtki.Label
		var issuedByCN, issuedByO, issuedByOU gtki.Label
		var issuedOn, expiresOn gtki.Label
		var sha1Fingerprint, sha256Fingerprint, sha3_256Fingerprint gtki.Label

		builder.getItems(
			"dialog", &md,
			"message", &message,
			"issuedToCnValue", &issuedToCN,
			"issuedToOValue", &issuedToO,
			"issuedToOUValue", &issuedToOU,
			"SNValue", &serial,
			"issuedByCnValue", &issuedByCN,
			"issuedByOValue", &issuedByO,
			"issuedByOUValue", &issuedByOU,
			"issuedOnValue", &issuedOn,
			"expiresOnValue", &expiresOn,
			"sha1FingerprintValue", &sha1Fingerprint,
			"sha256FingerprintValue", &sha256Fingerprint,
			"sha3_256FingerprintValue", &sha3_256Fingerprint,
		)

		issuedToCN.SetLabel(certs[0].Subject.CommonName)
		issuedToO.SetLabel(strings.Join(certs[0].Subject.Organization, ", "))
		issuedToOU.SetLabel(strings.Join(certs[0].Subject.OrganizationalUnit, ", "))
		serial.SetLabel(certs[0].SerialNumber.String())

		issuedByCN.SetLabel(certs[0].Issuer.CommonName)
		issuedByO.SetLabel(strings.Join(certs[0].Issuer.Organization, ", "))
		issuedByOU.SetLabel(strings.Join(certs[0].Issuer.OrganizationalUnit, ", "))

		issuedOn.SetLabel(certs[0].NotBefore.Format(time.RFC822))
		expiresOn.SetLabel(certs[0].NotAfter.Format(time.RFC822))

		sha1Fingerprint.SetLabel(displayChunked(digests.Sha1(certs[0].Raw)))
		sha256Fingerprint.SetLabel(displayChunked(digests.Sha256(certs[0].Raw)))
		sha3_256Fingerprint.SetLabel(displayChunked(digests.Sha3_256(certs[0].Raw)))

		accountName := "this account"
		if a != nil {
			accountName = a.session.GetConfig().Account
		}

		md.SetTitle(strings.Replace(md.GetTitle(), "ACCOUNT_NAME", accountName, -1))
		message.SetLabel(strings.Replace(message.GetLabel(), "ACCOUNT_NAME", accountName, -1))

		md.SetTransientFor(u.window)

		md.ShowAll()

		switch gtki.ResponseType(md.Run()) {
		case gtki.RESPONSE_OK:
			c <- true
		case gtki.RESPONSE_ACCEPT:
			if a != nil {
				a.session.GetConfig().SaveCert(certs[0].Subject.CommonName, certs[0].Issuer.CommonName, digests.Sha3_256(certs[0].Raw))
				u.SaveConfig()
			}
			c <- true
		case gtki.RESPONSE_CANCEL:
			if a != nil {
				a.session.SetWantToBeOnline(false)
			}
			c <- false
		default:
			a.session.SetWantToBeOnline(false)
			c <- false
		}

		md.Destroy()
	})
}
