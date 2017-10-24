package gliba

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/coyim/gotk3adapter/glibi"
)

type settings struct {
	*Object
	*glib.Settings
}

func wrapSettingsSimple(v *glib.Settings) *settings {
	if v == nil {
		return nil
	}
	return &settings{WrapObjectSimple(v.Object), v}
}

func unwrapSettings(v glibi.Settings) *glib.Settings {
	if v == nil {
		return nil
	}
	return v.(*settings).Settings
}

func (v *settings) IsWritable(v1 string) bool {
	return v.Settings.IsWritable(v1)
}

func (v *settings) Delay() {
	v.Settings.Delay()
}

func (v *settings) Apply() {
	v.Settings.Apply()
}

func (v *settings) Revert() {
	v.Settings.Revert()
}

func (v *settings) GetHasUnapplied() bool {
	return v.Settings.GetHasUnapplied()
}

func (v *settings) GetChild(v1 string) glibi.Settings {
	return wrapSettingsSimple(v.Settings.GetChild(v1))
}

func (v *settings) Reset(v1 string) {
	v.Settings.Reset(v1)
}

func (v *settings) ListChildren() []string {
	return v.Settings.ListChildren()
}

func (v *settings) GetBoolean(v1 string) bool {
	return v.Settings.GetBoolean(v1)
}

func (v *settings) SetBoolean(v1 string, v2 bool) bool {
	return v.Settings.SetBoolean(v1, v2)
}

func (v *settings) GetInt(v1 string) int {
	return v.Settings.GetInt(v1)
}

func (v *settings) SetInt(v1 string, v2 int) bool {
	return v.Settings.SetInt(v1, v2)
}

func (v *settings) GetUInt(v1 string) uint {
	return v.Settings.GetUInt(v1)
}

func (v *settings) SetUInt(v1 string, v2 uint) bool {
	return v.Settings.SetUInt(v1, v2)
}

func (v *settings) GetDouble(v1 string) float64 {
	return v.Settings.GetDouble(v1)
}

func (v *settings) SetDouble(v1 string, v2 float64) bool {
	return v.Settings.SetDouble(v1, v2)
}

func (v *settings) GetString(v1 string) string {
	return v.Settings.GetString(v1)
}

func (v *settings) SetString(v1 string, v2 string) bool {
	return v.Settings.SetString(v1, v2)
}

func (v *settings) GetEnum(v1 string) int {
	return v.Settings.GetEnum(v1)
}

func (v *settings) SetEnum(v1 string, v2 int) bool {
	return v.Settings.SetEnum(v1, v2)
}

func (v *settings) GetFlags(v1 string) uint {
	return v.Settings.GetFlags(v1)
}

func (v *settings) SetFlags(v1 string, v2 uint) bool {
	return v.Settings.SetFlags(v1, v2)
}
