package gui

import (
	"crypto/x509"
	"errors"

	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/digests"
	ourtls "github.com/twstrike/coyim/tls"
)

func (u *gtkUI) verifierFor(account *account) ourtls.Verifier {
	conf := account.session.GetConfig()
	return &ourtls.BasicVerifier{
		OnNoPeerCertificates: func() { account.session.SetWantToBeOnline(false) },
		OnPinDeny:            func() { account.session.SetWantToBeOnline(false) },
		HasPinned:            func(certs []*x509.Certificate) bool { return checkPinned(conf, certs) },
		VerifyFailure: func(certs []*x509.Certificate, err error) error {
			if <-u.certificateFailedToVerify(account, certs) {
				return nil
			}
			return errors.New("tls: failed to verify TLS certificate: " + err.Error())
		},
		VerifyHostnameFailure: func(certs []*x509.Certificate, origin string, err error) error {
			if <-u.certificateFailedToVerifyHostname(account, certs, origin) {
				return nil
			}
			return errors.New("tls: failed to match TLS certificate to name: " + err.Error())
		},
		AddCert: func(cert *x509.Certificate) {
			conf.SaveCert(cert.Subject.CommonName, cert.Issuer.CommonName, digests.Sha3_256(cert.Raw))
			u.SaveConfig()
		},
		AskPinning: func(certs []*x509.Certificate) error {
			if <-u.validCertificateShouldBePinned(account, certs) {
				return nil
			}
			account.session.SetWantToBeOnline(false)
			return errors.New("tls: you manually denied the possibility of connecting using this certificate")
		},
		HasCertificates: func() bool { return len(conf.Certificates) > 0 },
		NeedToCheckPins: conf.PinningPolicy != "none",
		PinningPolicy:   conf.PinningPolicy,
	}
}

func (u *gtkUI) unassociatedVerifier() ourtls.Verifier {
	return &ourtls.BasicVerifier{
		OnNoPeerCertificates: func() {},
		OnPinDeny:            func() {},
		HasPinned:            func(certs []*x509.Certificate) bool { return false },
		VerifyFailure: func(certs []*x509.Certificate, err error) error {
			if <-u.certificateFailedToVerify(nil, certs) {
				return nil
			}
			return errors.New("tls: failed to verify TLS certificate: " + err.Error())
		},
		VerifyHostnameFailure: func(certs []*x509.Certificate, origin string, err error) error {
			if <-u.certificateFailedToVerifyHostname(nil, certs, origin) {
				return nil
			}
			return errors.New("tls: failed to match TLS certificate to name: " + err.Error())
		},
		AddCert: func(cert *x509.Certificate) {},
		AskPinning: func(certs []*x509.Certificate) error {
			if <-u.validCertificateShouldBePinned(nil, certs) {
				return nil
			}
			return errors.New("tls: you manually denied the possibility of connecting using this certificate")
		},
		HasCertificates: func() bool { return false },
		NeedToCheckPins: false,
		PinningPolicy:   "",
	}
}

func checkPinned(c *config.Account, certs []*x509.Certificate) bool {
	if c != nil {
		for _, pin := range c.Certificates {
			if pin.Matches(certs[0]) {
				return true
			}
		}
	}

	return false
}
