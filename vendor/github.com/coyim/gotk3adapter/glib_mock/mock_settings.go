package glib_mock

import "github.com/coyim/gotk3adapter/glibi"

type MockSettings struct {
	MockObject
}

func (*MockSettings) IsWritable(string) bool {
	return false
}

func (*MockSettings) Delay() {
}

func (*MockSettings) Apply() {
}

func (*MockSettings) Revert() {
}

func (*MockSettings) GetHasUnapplied() bool {
	return false
}

func (*MockSettings) GetChild(string) glibi.Settings {
	return nil
}

func (*MockSettings) Reset(string) {
}

func (*MockSettings) ListChildren() []string {
	return nil
}

func (*MockSettings) GetBoolean(string) bool {
	return false
}

func (*MockSettings) SetBoolean(string, bool) bool {
	return false
}

func (*MockSettings) GetInt(string) int {
	return 0
}

func (*MockSettings) SetInt(string, int) bool {
	return false
}

func (*MockSettings) GetUInt(string) uint {
	return 0
}

func (*MockSettings) SetUInt(string, uint) bool {
	return false
}

func (*MockSettings) GetDouble(string) float64 {
	return 0
}

func (*MockSettings) SetDouble(string, float64) bool {
	return false
}

func (*MockSettings) GetString(string) string {
	return ""
}

func (*MockSettings) SetString(string, string) bool {
	return false
}

func (*MockSettings) GetEnum(string) int {
	return 0
}

func (*MockSettings) SetEnum(string, int) bool {
	return false
}

func (*MockSettings) GetFlags(string) uint {
	return 0
}

func (*MockSettings) SetFlags(string, uint) bool {
	return false
}
