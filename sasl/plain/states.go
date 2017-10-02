package plain

import "github.com/coyim/coyim/sasl"

type state interface {
	challenge(string, string) (state, sasl.Token, error)
}

type replyChallenge struct{}

func (replyChallenge) challenge(user string, password string) (state, sasl.Token, error) {
	ret := sasl.Token("\x00" + user + "\x00" + password)
	return authenticateServer{}, ret, nil
}

type authenticateServer struct{}

func (authenticateServer) challenge(user string, pass string) (state, sasl.Token, error) {
	//Server is always authenticated
	return finished{}, nil, nil
}

type finished struct{}

func (finished) challenge(string, string) (state, sasl.Token, error) {
	return finished{}, nil, nil
}
