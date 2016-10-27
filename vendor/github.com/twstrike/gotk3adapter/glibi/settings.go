package glibi

type Settings interface {
	Object

	IsWritable(string) bool
	Delay()
	Apply()
	Revert()
	GetHasUnapplied() bool
	GetChild(string) Settings
	Reset(string)
	ListChildren() []string
	GetBoolean(string) bool
	SetBoolean(string, bool) bool
	GetInt(string) int
	SetInt(string, int) bool
	GetUInt(string) uint
	SetUInt(string, uint) bool
	GetDouble(string) float64
	SetDouble(string, float64) bool
	GetString(string) string
	SetString(string, string) bool
	GetEnum(string) int
	SetEnum(string, int) bool
	GetFlags(string) uint
	SetFlags(string, uint) bool
}

func AssertSettings(_ Settings) {}
