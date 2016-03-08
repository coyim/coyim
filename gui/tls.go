package gui

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"

	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/digests"
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
	if <-v.u.certificateFailedToVerify(v.account, certs) {
		return nil
	}
	return errors.New("tls: failed to verify TLS certificate: " + err.Error())
}

func (v *accountTLSVerifier) verifyHostnameFailure(certs []*x509.Certificate, origin string, err error) error {
	if <-v.u.certificateFailedToVerifyHostname(v.account, certs, origin) {
		return nil
	}
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

func (v *accountTLSVerifier) needToCheckPins() bool {
	return v.config != nil && v.config.PinningPolicy != "none"
}

func (v *accountTLSVerifier) addCert(cert *x509.Certificate) {
	fmt.Printf("Adding pinned cert: %s, %s, %X", cert.Subject.CommonName, cert.Issuer.CommonName, digests.Sha3_256(cert.Raw))
	v.account.session.GetConfig().SaveCert(cert.Subject.CommonName, cert.Issuer.CommonName, digests.Sha3_256(cert.Raw))
	v.u.SaveConfig()
}

func (v *accountTLSVerifier) askPinning(certs []*x509.Certificate) error {
	if <-v.u.validCertificateShouldBePinned(v.account, certs) {
		return nil
	}
	v.account.session.SetWantToBeOnline(false)
	return errors.New("tls: you manually denied the possibility of connecting using this certificate")
}

func (v *accountTLSVerifier) canConnectInPresenceOfPins(certs []*x509.Certificate) error {
	if v.config != nil {
		switch v.config.PinningPolicy {
		case "none": // No pinning policy. This means we will not even consider saving pins or checking them
			return nil
		case "deny": // We will never approve a new certificate and will fail immediately, not even asking the user
			v.account.session.SetWantToBeOnline(false)
			return errors.New("tls: you have a pinning policy that stops us from connecting using this certificate")
		case "add": // We will always add a new certificate to our list of certs
			v.addCert(certs[0])
			return nil
		case "", "add-first-ask-rest": // We will quietly add the first cert but ask for the rest. This is the default.
			if len(v.account.session.GetConfig().Certificates) == 0 {
				v.addCert(certs[0])
				return nil
			}
			return v.askPinning(certs)
		case "add-first-deny-rest": // We will quietly add the first cert but deny for the rest
			if len(v.account.session.GetConfig().Certificates) == 0 {
				v.addCert(certs[0])
				return nil
			}
			v.account.session.SetWantToBeOnline(false)
			return errors.New("tls: you have a pinning policy that stops us from connecting using this certificate")
		case "ask": // We will always ask
			return v.askPinning(certs)
		}

		v.account.session.SetWantToBeOnline(false)
		return errors.New("tls: you have a pinning policy that stops us from connecting using other certificates")
	}

	return nil
}

func (v *accountTLSVerifier) Verify(state tls.ConnectionState, conf tls.Config, originDomain string) error {
	if len(state.PeerCertificates) == 0 {
		v.account.session.SetWantToBeOnline(false)
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
		return v.verifyHostnameFailure(state.PeerCertificates, originDomain, err)
	}

	if v.needToCheckPins() {
		return v.canConnectInPresenceOfPins(state.PeerCertificates)
	}

	return nil
}
