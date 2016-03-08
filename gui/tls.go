package gui

import (
	"crypto/tls"
	"crypto/x509"
	"errors"

	"github.com/twstrike/coyim/config"
	ourtls "github.com/twstrike/coyim/tls"
)

func (u *gtkUI) verifierFor(account *account) ourtls.Verifier {
	return &accountTLSVerifier{u, account, account.session.GetConfig()}
}

// TODO: fill in here
func (u *gtkUI) unassociatedVerifier() ourtls.Verifier {
	return &accountTLSVerifier{u: u}
}

type accountTLSVerifier struct {
	u       *gtkUI
	account *account
	config  *config.Account
}

func (v *accountTLSVerifier) verifyCert(certs []*x509.Certificate, conf tls.Config) ([][]*x509.Certificate, error) {
	opts := x509.VerifyOptions{
		Intermediates: x509.NewCertPool(),
		Roots:         conf.RootCAs,
	}

	for _, cert := range certs[1:] {
		opts.Intermediates.AddCert(cert)
	}

	return certs[0].Verify(opts)
}

func (v *accountTLSVerifier) verifyFailure(certs []*x509.Certificate, err error) error {
	c := v.u.certificateFailedToVerify(v.account, certs)
	if <-c {
		return nil
	}
	return errors.New("tls: failed to verify TLS certificate: " + err.Error())
}

func (v *accountTLSVerifier) verifyHostnameFailure(err error) error {
	return errors.New("tls: failed to match TLS certificate to name: " + err.Error())
}

func (v *accountTLSVerifier) verifyHostName(leafCert *x509.Certificate, originDomain string) error {
	return leafCert.VerifyHostname(originDomain)
}

func (v *accountTLSVerifier) hasPinned(certs []*x509.Certificate) bool {
	if v.config != nil {
		for _, pin := range v.config.Certificates {
			if pin.Matches(certs[0]) {
				return true
			}
		}
	}

	return false
}

func (v *accountTLSVerifier) Verify(state tls.ConnectionState, conf tls.Config, originDomain string) error {
	if len(state.PeerCertificates) == 0 {
		return errors.New("tls: server has no certificates")
	}

	if v.hasPinned(state.PeerCertificates) {
		return nil
	}

	chains, err := v.verifyCert(state.PeerCertificates, conf)
	if err != nil {
		return v.verifyFailure(state.PeerCertificates, err)
	}

	if err = v.verifyHostName(chains[0][0], originDomain); err != nil {
		return v.verifyHostnameFailure(err)
	}

	return nil
}
