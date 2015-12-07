package digestmd5

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/twstrike/coyim/sasl"
)

type digestState interface {
	challenge(sasl.Properties, sasl.AttributeValuePairs) (digestState, sasl.Token, error)
}

type clientChallenge struct{}

func (clientChallenge) challenge(props sasl.Properties, pairs sasl.AttributeValuePairs) (digestState, sasl.Token, error) {
	return digestChallenge{}, nil, nil
}

// RFC 2831, section 2.1.2
type digestChallenge struct{}

func (c digestChallenge) challenge(props sasl.Properties, pairs sasl.AttributeValuePairs) (digestState, sasl.Token, error) {
	var ok bool

	user, ok := props[sasl.AuthID]
	if !ok {
		return c, nil, sasl.PropertyMissingError{sasl.AuthID}
	}

	password, ok := props[sasl.Password]
	if !ok {
		return c, nil, sasl.PropertyMissingError{sasl.Password}
	}

	authorizationID := props[sasl.AuthZID]

	//realm, ok := props[sasl.Realm]
	//if !ok {
	//	return c, nil, sasl.PropertyMissingError{sasl.Realm}
	//}

	realm, ok := pairs["realm"]
	if !ok {
		return c, nil, sasl.ErrMissingParameter
	}

	service, ok := props[sasl.Service]
	if !ok {
		return c, nil, sasl.PropertyMissingError{sasl.Service}
	}

	clientNonce, ok := props[sasl.ClientNonce]
	if !ok {
		return c, nil, sasl.PropertyMissingError{sasl.ClientNonce}
	}

	qop, ok := props[sasl.QOP]
	if !ok {
		return c, nil, sasl.PropertyMissingError{sasl.QOP}
	}

	serverNonce, ok := pairs["nonce"]
	if !ok {
		return c, nil, sasl.ErrMissingParameter
	}

	//TODO: check if pairs[qop] contains props[QOP]

	digestURI := service + "/" + realm

	a1 := c.a1(user, realm, password, serverNonce, clientNonce, authorizationID)
	a2 := c.a2(digestURI, qop)
	nc := fmt.Sprintf("%08x", 1)

	kd := strings.Join([]string{
		hex.EncodeToString(a1[:]),
		serverNonce,
		nc,
		clientNonce,
		qop,
		hex.EncodeToString(a2[:]),
	}, ":")

	mdKd := md5.Sum([]byte(kd))
	responseValue := hex.EncodeToString(mdKd[:])

	ret := fmt.Sprintf(`charset=utf-8,username=%q,realm=%q,nonce=%q,nc=%s,cnonce=%q,digest-uri=%q,response=%s,qop=auth`,
		user, realm, serverNonce, nc, clientNonce, digestURI, responseValue,
	)

	if authorizationID != "" {
		ret = ret + ",authzid=\"" + authorizationID + "\""
	}

	return responseAuth{}, sasl.Token(ret), nil
}

func (c digestChallenge) a1(user, realm, password, serverNonce, clientNonce, authorizationID string) [md5.Size]byte {
	//TODO: they should be encoded using the charset specified by the server
	x := strings.Join([]string{
		user, realm, password,
	}, ":")
	y := md5.Sum([]byte(x))

	a1Values := []string{
		string(y[:]),
		serverNonce,
		clientNonce,
	}

	if authorizationID != "" {
		a1Values = append(a1Values, authorizationID)
	}

	a1 := strings.Join(a1Values, ":")
	return md5.Sum([]byte(a1))
}

func (c digestChallenge) digestURI(service, realm string) string {
	return service + "/" + realm
}

func (c digestChallenge) a2(digestURI, qop string) [md5.Size]byte {
	a2 := strings.Join([]string{
		"AUTHENTICATE",
		digestURI,
	}, ":")

	if qop == "auth-int" {
		a2 = a2 + ":00000000000000000000000000000000"
	}

	return md5.Sum([]byte(a2))
}

type responseAuth struct{}

func (c responseAuth) challenge(props sasl.Properties, pairs sasl.AttributeValuePairs) (digestState, sasl.Token, error) {
	if _, ok := pairs["rspauth"]; !ok {
		return c, nil, sasl.ErrMissingParameter
	}

	return responseAuthReply{}, nil, nil
}

type responseAuthReply struct{}

func (responseAuthReply) challenge(props sasl.Properties, pairs sasl.AttributeValuePairs) (digestState, sasl.Token, error) {
	// RFC 4422 section 3
	// Where the mechanism specifies that the server is to return additional
	// data to the client with a successful outcome and this field is
	// unavailable or unused, the additional data is sent as a challenge
	// whose response is empty.
	return finished{}, nil, nil
}

type finished struct{}

func (c finished) challenge(props sasl.Properties, pairs sasl.AttributeValuePairs) (digestState, sasl.Token, error) {
	return c, nil, nil
}
