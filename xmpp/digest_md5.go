package xmpp

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

type digestMD5 struct {
	user            string
	password        string
	clientNonce     string
	servType        string
	authorizationID string

	serverNonce string
	realm       string
	qop         string

	encoder *base64.Encoding
}

var (
	errMissingNonce    = errors.New("missing server nonce")
	errAuthQOPRequired = errors.New("server does not support auth qop")
)

func (d *digestMD5) receive(msg string) error {
	dec, err := base64.StdEncoding.DecodeString(msg)
	if err != nil {
		return err
	}

	params := strings.Split(string(dec), ",")

	for _, p := range params {
		switch {
		case strings.HasPrefix(p, "realm="):
			realm := strings.TrimPrefix(p, "realm=")
			d.realm = realm[1 : len(realm)-1]
		case strings.HasPrefix(p, "nonce="):
			//TODO: should have only one nonce
			nonce := strings.TrimPrefix(p, "nonce=")
			d.serverNonce = nonce[1 : len(nonce)-1]
		case strings.HasPrefix(p, "qop="):
			qop := strings.TrimPrefix(p, "qop=")
			d.qop = qop[1 : len(qop)-1]
		}
	}

	if len(d.serverNonce) == 0 {
		return errMissingNonce
	}

	return nil
}

func (d *digestMD5) a1() [md5.Size]byte {
	//TODO: they should be encoded using the charset specified by the server
	x := strings.Join([]string{
		d.user, d.realm, d.password,
	}, ":")
	y := md5.Sum([]byte(x))

	a1Values := []string{
		string(y[:]),
		d.serverNonce,
		d.clientNonce,
	}

	if d.authorizationID != "" {
		a1Values = append(a1Values, d.authorizationID)
	}

	a1 := strings.Join(a1Values, ":")
	return md5.Sum([]byte(a1))
}

func (d *digestMD5) digestURI() string {
	return d.servType + "/" + d.realm
}

func (d *digestMD5) a2() [md5.Size]byte {
	a2 := strings.Join([]string{
		"AUTHENTICATE",
		d.digestURI(),
	}, ":")

	if d.qop == "auth-int" {
		a2 = a2 + ":00000000000000000000000000000000"
	}

	return md5.Sum([]byte(a2))
}

func (d *digestMD5) send() string {
	if d.clientNonce == "" {
		panic("missing client nonce")
	}

	a1 := d.a1()
	a2 := d.a2()
	nc := fmt.Sprintf("%08x", 1)

	kd := strings.Join([]string{
		hex.EncodeToString(a1[:]),
		d.serverNonce,
		nc,
		d.clientNonce,
		d.qop,
		hex.EncodeToString(a2[:]),
	}, ":")

	mdKd := md5.Sum([]byte(kd))
	responseValue := hex.EncodeToString(mdKd[:])

	ret := fmt.Sprintf(`charset=utf-8,username=%q,realm=%q,nonce=%q,nc=%s,cnonce=%q,digest-uri=%q,response=%s,qop=auth`,
		d.user, d.realm, d.serverNonce, nc, d.clientNonce, d.digestURI(), responseValue,
	)

	if d.authorizationID != "" {
		ret = ret + ",authzid=\"" + d.authorizationID + "\""
	}

	return base64.StdEncoding.EncodeToString([]byte(ret))
}

func (d *digestMD5) verifyResponse(response string) error {
	dec, err := base64.StdEncoding.DecodeString(response)
	if err != nil {
		return err
	}

	if !strings.HasPrefix(string(dec), "rspauth=") {
		return ErrAuthenticationFailed
	}

	return nil
}
