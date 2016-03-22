package tls

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
)

type BasicVerifier struct {
	OnNoPeerCertificates  func()
	OnPinDeny             func()
	HasPinned             func([]*x509.Certificate) bool
	VerifyFailure         func([]*x509.Certificate, error) error
	VerifyHostnameFailure func([]*x509.Certificate, string, error) error
	AddCert               func(*x509.Certificate)
	AskPinning            func([]*x509.Certificate) error
	HasCertificates       func() bool
	NeedToCheckPins       bool
	PinningPolicy         string
}

func verifyHostName(leafCert *x509.Certificate, originDomain string) error {
	return leafCert.VerifyHostname(originDomain)
}

func (v *BasicVerifier) verifyCert(certs []*x509.Certificate, conf tls.Config) ([][]*x509.Certificate, error) {
	opts := x509.VerifyOptions{
		Intermediates: x509.NewCertPool(),
		Roots:         conf.RootCAs,
	}

	for _, cert := range certs[1:] {
		opts.Intermediates.AddCert(cert)
	}

	return certs[0].Verify(opts)
}

func (v *BasicVerifier) canConnectInPresenceOfPins(certs []*x509.Certificate) error {
	switch v.PinningPolicy {
	case "none": // No pinning policy. This means we will not even consider saving pins or checking them
		return nil
	case "deny": // We will never approve a new certificate and will fail immediately, not even asking the user
		v.OnPinDeny()
		return errors.New("tls: you have a pinning policy that stops us from connecting using this certificate")
	case "add": // We will always add a new certificate to our list of certs
		v.AddCert(certs[0])
		return nil
	case "", "add-first-ask-rest": // We will quietly add the first cert but ask for the rest. This is the default.
		if !v.HasCertificates() {
			v.AddCert(certs[0])
			return nil
		}
		return v.AskPinning(certs)
	case "add-first-deny-rest": // We will quietly add the first cert but deny for the rest
		if !v.HasCertificates() {
			v.AddCert(certs[0])
			return nil
		}
		v.OnPinDeny()
		return errors.New("tls: you have a pinning policy that stops us from connecting using this certificate")
	case "ask": // We will always ask
		return v.AskPinning(certs)
	}

	v.OnPinDeny()
	return errors.New("tls: you have a pinning policy that stops us from connecting using other certificates")
}

func (v *BasicVerifier) Verify(state tls.ConnectionState, conf tls.Config, originDomain string) error {
	if len(state.PeerCertificates) == 0 {
		v.OnNoPeerCertificates()
		return errors.New("tls: server has no certificates")
	}

	if v.HasPinned(state.PeerCertificates) {
		return nil
	}

	chains, err := v.verifyCert(state.PeerCertificates, conf)
	if err != nil {
		return v.VerifyFailure(state.PeerCertificates, err)
	}

	if err = verifyHostName(chains[0][0], originDomain); err != nil {
		return v.VerifyHostnameFailure(state.PeerCertificates, originDomain, err)
	}

	if v.NeedToCheckPins {
		return v.canConnectInPresenceOfPins(state.PeerCertificates)
	}

	return nil
}
