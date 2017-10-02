package cli

import (
	"crypto/x509"
	"errors"

	"github.com/coyim/coyim/config"
	"github.com/coyim/coyim/digests"
	ourtls "github.com/coyim/coyim/tls"
)

func (c *cliUI) verifier() ourtls.Verifier {
	conf := c.session.GetConfig()
	return &ourtls.BasicVerifier{
		OnNoPeerCertificates: func() { c.session.SetWantToBeOnline(false) },
		OnPinDeny:            func() { c.session.SetWantToBeOnline(false) },
		HasPinned:            func(certs []*x509.Certificate) bool { return checkPinned(conf, certs) },
		VerifyFailure: func(certs []*x509.Certificate, err error) error {
			return errors.New("tls: failed to verify TLS certificate: " + err.Error())
		},
		VerifyHostnameFailure: func(certs []*x509.Certificate, origin string, err error) error {
			return errors.New("tls: failed to match TLS certificate to name: " + err.Error())
		},
		AddCert: func(cert *x509.Certificate) {
			conf.SaveCert(cert.Subject.CommonName, cert.Issuer.CommonName, digests.Sha3_256(cert.Raw))
			c.SaveConf()
		},
		AskPinning: func(certs []*x509.Certificate) error {
			c.session.SetWantToBeOnline(false)
			return errors.New("tls: you manually denied the possibility of connecting using this certificate")
		},
		HasCertificates: func() bool { return len(conf.Certificates) > 0 },
		NeedToCheckPins: conf.PinningPolicy != "none",
		PinningPolicy:   conf.PinningPolicy,
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
