package nettest

import "../../net"

// MockTorState returns a mocked TorState
func MockTorState(addr string) net.TorState {
	return &torStateMock{addr}
}

type torStateMock struct {
	addr string
}

func (s *torStateMock) Address() string {
	return s.addr
}

func (s *torStateMock) Detect() bool {
	return len(s.addr) > 0
}
