package glib_mock

import "github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/glibi"

type Mock struct{}

func (*Mock) IdleAdd(f interface{}, args ...interface{}) (glibi.SourceHandle, error) {
	return glibi.SourceHandle(0), nil
}

func (*Mock) InitI18n(domain string, dir string) {
}

func (*Mock) Local(vx string) string {
	return vx
}

func (*Mock) MainDepth() int {
	return 0
}

func (*Mock) SignalNew(s string) (glibi.Signal, error) {
	return &MockSignal{}, nil
}
