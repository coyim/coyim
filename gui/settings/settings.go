package settings

import (
	"fmt"
	"os"

	"github.com/twstrike/coyim/gui/settings/definitions"
	"github.com/twstrike/gotk3adapter/glibi"
)

var g glibi.Glib

// InitSettings should be called before using settings
func InitSettings(gx glibi.Glib) {
	g = gx
}

var cachedSchema glibi.SettingsSchemaSource

func getSchemaSource() glibi.SettingsSchemaSource {
	if cachedSchema == nil {
		dir := definitions.SchemaInTempDir()
		defer os.Remove(dir)
		cachedSchema = g.SettingsSchemaSourceNewFromDirectory(dir, nil, true)
	}

	return cachedSchema
}

func getSchema() glibi.SettingsSchema {
	return getSchemaSource().Lookup("im.coy.coyim.MainSettings", false)
}

func getSettingsFor(s string) glibi.Settings {
	return g.SettingsNewFull(getSchema(), nil, fmt.Sprintf("/im/coy/coyim/%s/", s))
}

func getDefaultSettings() glibi.Settings {
	return g.SettingsNewFull(getSchema(), nil, "/im/coy/coyim/")
}

// Settings allow access to our configured settings
type Settings struct {
	def, spec glibi.Settings
}

// For will return a valid settings instance for the given ident, or the empty string
func For(ident string) *Settings {
	s := &Settings{}
	s.def = getDefaultSettings()
	if ident != "" {
		s.spec = getSettingsFor(ident)
	}
	return s
}

func hasSetStr(k string) string {
	return fmt.Sprintf("has-set-%s", k)
}

func hasSet(s glibi.Settings, k string) bool {
	return s.GetBoolean(hasSetStr(k))
}

func (s *Settings) settingsForGet(name string) glibi.Settings {
	if s.spec != nil && hasSet(s.spec, name) {
		return s.spec
	}
	return s.def
}

func (s *Settings) settingsForSet() glibi.Settings {
	if s.spec != nil {
		return s.spec
	}
	return s.def
}

func (s *Settings) getBooleanSetting(name string) bool {
	return s.settingsForGet(name).GetBoolean(name)
}

func (s *Settings) setBooleanSetting(name string, val bool) {
	sets := s.settingsForSet()
	sets.SetBoolean(hasSetStr(name), true)
	sets.SetBoolean(name, val)
}

func (s *Settings) getStringSetting(name string) string {
	return s.settingsForGet(name).GetString(name)
}

func (s *Settings) setStringSetting(name string, val string) {
	sets := s.settingsForSet()
	sets.SetBoolean(hasSetStr(name), true)
	sets.SetString(name, val)
}

// GetSingleWindow returns the single-window setting
func (s *Settings) GetSingleWindow() bool {
	return s.getBooleanSetting("single-window")
}

// SetSingleWindow sets the single-window setting
func (s *Settings) SetSingleWindow(val bool) {
	s.setBooleanSetting("single-window", val)
}

// GetNotificationUrgency returns the notification-urgency setting
func (s *Settings) GetNotificationUrgency() bool {
	return s.getBooleanSetting("notification-urgency")
}

// SetNotificationUrgency sets the notification-urgency setting
func (s *Settings) SetNotificationUrgency(val bool) {
	s.setBooleanSetting("notification-urgency", val)
}

// GetNotificationExpires returns the notification-expires setting
func (s *Settings) GetNotificationExpires() bool {
	return s.getBooleanSetting("notification-expires")
}

// SetNotificationExpires sets the notification-expires setting
func (s *Settings) SetNotificationExpires(val bool) {
	s.setBooleanSetting("notification-expires", val)
}

// GetNotificationStyle returns the notification-style setting
func (s *Settings) GetNotificationStyle() string {
	return s.getStringSetting("notification-style")
}

// SetNotificationStyle sets the notification-style setting
func (s *Settings) SetNotificationStyle(val string) {
	s.setStringSetting("notification-style", val)
}

// GetShiftEnterForSend returns the shift-enter-for-send setting
func (s *Settings) GetShiftEnterForSend() bool {
	return s.getBooleanSetting("shift-enter-for-send")
}

// SetShiftEnterForSend sets the shift-enter-for-send setting
func (s *Settings) SetShiftEnterForSend(val bool) {
	s.setBooleanSetting("shift-enter-for-send", val)
}

// GetEmacsKeyBindings returns the emacs-keyboard setting
func (s *Settings) GetEmacsKeyBindings() bool {
	return s.getBooleanSetting("emacs-keyboard")
}

// SetEmacsKeyBindings sets the emacs-keyboard setting
func (s *Settings) SetEmacsKeyBindings(val bool) {
	s.setBooleanSetting("emacs-keyboard", val)
}

// GetShowEmptyGroups returns the show-empty-groups setting
func (s *Settings) GetShowEmptyGroups() bool {
	return s.getBooleanSetting("show-empty-groups")
}

// SetShowEmptyGroups sets the show-empty-groups setting
func (s *Settings) SetShowEmptyGroups(val bool) {
	s.setBooleanSetting("show-empty-groups", val)
}

// GetShowAdvancedOptions returns the show-advanced-options setting
func (s *Settings) GetShowAdvancedOptions() bool {
	return s.getBooleanSetting("show-advanced-options")
}

// SetShowAdvancedOptions sets the show-advanced-options setting
func (s *Settings) SetShowAdvancedOptions(val bool) {
	s.setBooleanSetting("show-advanced-options", val)
}

// GetConnectAutomatically returns the connect-automatically setting
func (s *Settings) GetConnectAutomatically() bool {
	return s.getBooleanSetting("connect-automatically")
}

// SetConnectAutomatically sets the connect-automatically setting
func (s *Settings) SetConnectAutomatically(val bool) {
	s.setBooleanSetting("connect-automatically", val)
}

// GetCollapsed returns the collapsed setting
func (s *Settings) GetCollapsed() string {
	return s.getStringSetting("collapsed")
}

// SetCollapsed sets the collapsed setting
func (s *Settings) SetCollapsed(val string) {
	s.setStringSetting("collapsed", val)
}
