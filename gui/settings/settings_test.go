package settings

import (
	"github.com/coyim/gotk3adapter/glib_mock"
	"github.com/coyim/gotk3adapter/glibi"
	. "gopkg.in/check.v1"
)

type SettingsSuite struct{}

var _ = Suite(&SettingsSuite{})

type mockGlib struct {
	glib_mock.Mock

	schemaSourceArg1 string
	schemaSourceArg2 glibi.SettingsSchemaSource
	schemaSourceArg3 bool

	schemaSourceToReturn glibi.SettingsSchemaSource

	settingsNewFullArg1   glibi.SettingsSchema
	settingsNewFullArg2   glibi.SettingsBackend
	settingsNewFullArg3   string
	settingsNewFullReturn glibi.Settings
}

func (m *mockGlib) SettingsSchemaSourceNewFromDirectory(v1 string, v2 glibi.SettingsSchemaSource, v3 bool) glibi.SettingsSchemaSource {
	m.schemaSourceArg1 = v1
	m.schemaSourceArg2 = v2
	m.schemaSourceArg3 = v3
	return m.schemaSourceToReturn
}

func (m *mockGlib) SettingsNewFull(v1 glibi.SettingsSchema, v2 glibi.SettingsBackend, v3 string) glibi.Settings {
	m.settingsNewFullArg1 = v1
	m.settingsNewFullArg2 = v2
	m.settingsNewFullArg3 = v3
	return m.settingsNewFullReturn
}

type mockSettingsSchema struct {
	glib_mock.MockSettingsSchema
}

type mockSettingsSchemaSource struct {
	glib_mock.MockSettingsSchemaSource

	lookupArg1   string
	lookupArg2   bool
	lookupReturn glibi.SettingsSchema
}

func (m *mockSettingsSchemaSource) Lookup(v1 string, v2 bool) glibi.SettingsSchema {
	m.lookupArg1 = v1
	m.lookupArg2 = v2
	return m.lookupReturn
}

func (s *SettingsSuite) Test_InitSettings_initializesGlib(c *C) {
	orgG := g
	defer func() {
		g = orgG
	}()

	v := &mockGlib{}
	InitSettings(v)

	c.Assert(v, Equals, g)
}

func (s *SettingsSuite) Test_getSchemaSource_createsANewSchemaSourceIfNeeded(c *C) {
	orgG := g
	defer func() {
		g = orgG
	}()

	ss := &mockSettingsSchemaSource{}
	gv := &mockGlib{
		schemaSourceToReturn: ss,
	}
	g = gv

	ret := getSchemaSource()
	c.Assert(ret, Equals, ss)
	c.Assert(gv.schemaSourceArg1, Not(Equals), "")
	c.Assert(gv.schemaSourceArg2, IsNil)
	c.Assert(gv.schemaSourceArg3, Equals, true)
}

func (s *SettingsSuite) Test_getSchemaSource_returnsCachedSourceIfAvailable(c *C) {
	defer func() {
		cachedSchema = nil
	}()

	ss := &mockSettingsSchemaSource{}
	cachedSchema = ss

	ret := getSchemaSource()
	c.Assert(ret, Equals, ss)
}

func (s *SettingsSuite) Test_getSchema(c *C) {
	defer func() {
		cachedSchema = nil
	}()

	ss := &mockSettingsSchemaSource{}
	cachedSchema = ss
	sch := &mockSettingsSchema{}
	ss.lookupReturn = sch

	ret := getSchema()
	c.Assert(ret, Equals, sch)
	c.Assert(ss.lookupArg1, Equals, "im.coy.coyim.MainSettings")
	c.Assert(ss.lookupArg2, Equals, false)
}

func (s *SettingsSuite) Test_getSettingsFor(c *C) {
	orgG := g
	defer func() {
		g = orgG
	}()

	gv := &mockGlib{}
	g = gv

	defer func() {
		cachedSchema = nil
	}()

	ss := &mockSettingsSchemaSource{}
	cachedSchema = ss
	sch := &mockSettingsSchema{}
	ss.lookupReturn = sch

	set := &mockSettings{}
	gv.settingsNewFullReturn = set

	res := getSettingsFor("bla")
	c.Assert(res, Equals, set)
	c.Assert(gv.settingsNewFullArg1, Equals, sch)
	c.Assert(gv.settingsNewFullArg2, IsNil)
	c.Assert(gv.settingsNewFullArg3, Equals, "/im/coy/coyim/bla/")
}

func (s *SettingsSuite) Test_getDefaultSettings(c *C) {
	orgG := g
	defer func() {
		g = orgG
	}()

	gv := &mockGlib{}
	g = gv

	defer func() {
		cachedSchema = nil
	}()

	ss := &mockSettingsSchemaSource{}
	cachedSchema = ss
	sch := &mockSettingsSchema{}
	ss.lookupReturn = sch

	set := &mockSettings{}
	gv.settingsNewFullReturn = set

	res := getDefaultSettings()
	c.Assert(res, Equals, set)
	c.Assert(gv.settingsNewFullArg1, Equals, sch)
	c.Assert(gv.settingsNewFullArg2, IsNil)
	c.Assert(gv.settingsNewFullArg3, Equals, "/im/coy/coyim/")
}

func (s *SettingsSuite) Test_For_returnsOnlyDefaultSettings(c *C) {
	orgG := g
	defer func() {
		g = orgG
	}()

	gv := &mockGlib{}
	g = gv

	defer func() {
		cachedSchema = nil
	}()

	ss := &mockSettingsSchemaSource{}
	cachedSchema = ss
	sch := &mockSettingsSchema{}
	ss.lookupReturn = sch

	set := &mockSettings{}
	gv.settingsNewFullReturn = set

	res := For("")
	c.Assert(res.def, Equals, set)
	c.Assert(res.spec, IsNil)
}

type mockGlibSettings struct {
	glib_mock.Mock

	settingsNewFullReturn []glibi.Settings
}

func (m *mockGlibSettings) SettingsNewFull(v1 glibi.SettingsSchema, v2 glibi.SettingsBackend, v3 string) glibi.Settings {
	ret := m.settingsNewFullReturn[0]
	m.settingsNewFullReturn = m.settingsNewFullReturn[1:]
	return ret
}

func (s *SettingsSuite) Test_For_returnsBothGeneralAndSpecificSettings(c *C) {
	orgG := g
	defer func() {
		g = orgG
	}()

	gv := &mockGlibSettings{}
	g = gv

	defer func() {
		cachedSchema = nil
	}()

	ss := &mockSettingsSchemaSource{}
	cachedSchema = ss
	sch := &mockSettingsSchema{}
	ss.lookupReturn = sch

	set1 := &mockSettings{}
	set2 := &mockSettings{}
	gv.settingsNewFullReturn = []glibi.Settings{set1, set2}

	res := For("something")
	c.Assert(res.def, Equals, set1)
	c.Assert(res.spec, Equals, set2)
}

func (s *SettingsSuite) Test_hasSetStr(c *C) {
	c.Assert(hasSetStr("foo"), Equals, "has-set-foo")
	c.Assert(hasSetStr("bla"), Equals, "has-set-bla")
}

type mockSettings struct {
	glib_mock.MockSettings

	getBooleanArg []string
	getBooleanRet []bool

	setBooleanArg1 []string
	setBooleanArg2 []bool

	getStringArg []string
	getStringRet []string

	setStringArg1 []string
	setStringArg2 []string
}

func (m *mockSettings) GetBoolean(v string) bool {
	m.getBooleanArg = append(m.getBooleanArg, v)
	ret := m.getBooleanRet[0]
	m.getBooleanRet = m.getBooleanRet[1:]
	return ret
}

func (m *mockSettings) GetString(v string) string {
	m.getStringArg = append(m.getStringArg, v)
	ret := m.getStringRet[0]
	m.getStringRet = m.getStringRet[1:]
	return ret
}

func (m *mockSettings) SetString(v1 string, v2 string) bool {
	m.setStringArg1 = append(m.setStringArg1, v1)
	m.setStringArg2 = append(m.setStringArg2, v2)
	return true
}

func (m *mockSettings) SetBoolean(v1 string, v2 bool) bool {
	m.setBooleanArg1 = append(m.setBooleanArg1, v1)
	m.setBooleanArg2 = append(m.setBooleanArg2, v2)
	return true
}

func (s *SettingsSuite) Test_hasSet(c *C) {
	set := &mockSettings{
		getBooleanRet: []bool{false, true},
	}
	c.Assert(hasSet(set, "foo"), Equals, false)
	c.Assert(hasSet(set, "bla"), Equals, true)
	c.Assert(set.getBooleanArg, DeepEquals, []string{"has-set-foo", "has-set-bla"})
}

func (s *SettingsSuite) Test_settingsForGet_returnsGenericIfNoSpecialSet(c *C) {
	defset := &mockSettings{}
	set := &Settings{
		def: defset,
	}

	res := set.settingsForGet("something")
	c.Assert(res, Equals, defset)
}

func (s *SettingsSuite) Test_settingsForGet_returnsGenericIfSpecialSetButNoValueForConfig(c *C) {
	defset := &mockSettings{}
	specset := &mockSettings{
		getBooleanRet: []bool{false},
	}
	set := &Settings{
		def:  defset,
		spec: specset,
	}

	res := set.settingsForGet("something")
	c.Assert(res, Equals, defset)
	c.Assert(specset.getBooleanArg, DeepEquals, []string{"has-set-something"})
}

func (s *SettingsSuite) Test_settingsForGet_returnsSpecificIfValueSet(c *C) {
	defset := &mockSettings{}
	specset := &mockSettings{
		getBooleanRet: []bool{true},
	}
	set := &Settings{
		def:  defset,
		spec: specset,
	}

	res := set.settingsForGet("something")
	c.Assert(res, Equals, specset)
	c.Assert(specset.getBooleanArg, DeepEquals, []string{"has-set-something"})
}

func (s *SettingsSuite) Test_settingsForSet_returnsGenericIfNoSpecific(c *C) {
	defset := &mockSettings{}
	set := &Settings{
		def: defset,
	}

	res := set.settingsForSet()
	c.Assert(res, Equals, defset)
}

func (s *SettingsSuite) Test_settingsForSet_returnsSpecificIfExists(c *C) {
	defset := &mockSettings{}
	specset := &mockSettings{}
	set := &Settings{
		def:  defset,
		spec: specset,
	}

	res := set.settingsForSet()
	c.Assert(res, Equals, specset)
}

func (s *SettingsSuite) Test_Settings_GetCollapsed(c *C) {
	defset := &mockSettings{
		getStringRet: []string{"foo"},
	}
	set := &Settings{
		def: defset,
	}

	ret := set.GetCollapsed()
	c.Assert(ret, Equals, "foo")
}

func (s *SettingsSuite) Test_Settings_SetCollapsed(c *C) {
	defset := &mockSettings{}
	specset := &mockSettings{}
	set := &Settings{
		def:  defset,
		spec: specset,
	}

	set.SetCollapsed("something")
	c.Assert(specset.setStringArg1, DeepEquals, []string{"collapsed"})
	c.Assert(specset.setStringArg2, DeepEquals, []string{"something"})
}

func (s *SettingsSuite) Test_Settings_GetNotificationStyle(c *C) {
	defset := &mockSettings{
		getStringRet: []string{"geh"},
	}
	set := &Settings{
		def: defset,
	}

	ret := set.GetNotificationStyle()
	c.Assert(ret, Equals, "geh")
	c.Assert(defset.getStringArg, DeepEquals, []string{"notification-style"})
}

func (s *SettingsSuite) Test_Settings_SetNotificationStyle(c *C) {
	defset := &mockSettings{}
	specset := &mockSettings{}
	set := &Settings{
		def:  defset,
		spec: specset,
	}

	set.SetNotificationStyle("heh")
	c.Assert(specset.setStringArg1, DeepEquals, []string{"notification-style"})
	c.Assert(specset.setStringArg2, DeepEquals, []string{"heh"})
}

func (s *SettingsSuite) Test_Settings_GetSingleWindow(c *C) {
	defset := &mockSettings{
		getBooleanRet: []bool{false},
	}
	set := &Settings{
		def: defset,
	}

	ret := set.GetSingleWindow()
	c.Assert(ret, Equals, false)
	c.Assert(defset.getBooleanArg, DeepEquals, []string{"single-window"})
}

func (s *SettingsSuite) Test_Settings_SetSingleWindow(c *C) {
	defset := &mockSettings{}
	specset := &mockSettings{}
	set := &Settings{
		def:  defset,
		spec: specset,
	}

	set.SetSingleWindow(true)
	c.Assert(specset.setBooleanArg1, DeepEquals, []string{"has-set-single-window", "single-window"})
	c.Assert(specset.setBooleanArg2, DeepEquals, []bool{true, true})
}

func (s *SettingsSuite) Test_Settings_GetNotificationUrgency(c *C) {
	defset := &mockSettings{
		getBooleanRet: []bool{false},
	}
	set := &Settings{
		def: defset,
	}

	ret := set.GetNotificationUrgency()
	c.Assert(ret, Equals, false)
	c.Assert(defset.getBooleanArg, DeepEquals, []string{"notification-urgency"})
}

func (s *SettingsSuite) Test_Settings_SetNotificationUrgency(c *C) {
	defset := &mockSettings{}
	specset := &mockSettings{}
	set := &Settings{
		def:  defset,
		spec: specset,
	}

	set.SetNotificationUrgency(true)
	c.Assert(specset.setBooleanArg1, DeepEquals, []string{"has-set-notification-urgency", "notification-urgency"})
	c.Assert(specset.setBooleanArg2, DeepEquals, []bool{true, true})
}

func (s *SettingsSuite) Test_Settings_GetNotificationExpires(c *C) {
	defset := &mockSettings{
		getBooleanRet: []bool{false},
	}
	set := &Settings{
		def: defset,
	}

	ret := set.GetNotificationExpires()
	c.Assert(ret, Equals, false)
	c.Assert(defset.getBooleanArg, DeepEquals, []string{"notification-expires"})
}

func (s *SettingsSuite) Test_Settings_SetNotificationExpires(c *C) {
	defset := &mockSettings{}
	specset := &mockSettings{}
	set := &Settings{
		def:  defset,
		spec: specset,
	}

	set.SetNotificationExpires(true)
	c.Assert(specset.setBooleanArg1, DeepEquals, []string{"has-set-notification-expires", "notification-expires"})
	c.Assert(specset.setBooleanArg2, DeepEquals, []bool{true, true})
}

func (s *SettingsSuite) Test_Settings_GetShiftEnterForSend(c *C) {
	defset := &mockSettings{
		getBooleanRet: []bool{false},
	}
	set := &Settings{
		def: defset,
	}

	ret := set.GetShiftEnterForSend()
	c.Assert(ret, Equals, false)
	c.Assert(defset.getBooleanArg, DeepEquals, []string{"shift-enter-for-send"})
}

func (s *SettingsSuite) Test_Settings_SetShiftEnterForSend(c *C) {
	defset := &mockSettings{}
	specset := &mockSettings{}
	set := &Settings{
		def:  defset,
		spec: specset,
	}

	set.SetShiftEnterForSend(true)
	c.Assert(specset.setBooleanArg1, DeepEquals, []string{"has-set-shift-enter-for-send", "shift-enter-for-send"})
	c.Assert(specset.setBooleanArg2, DeepEquals, []bool{true, true})
}

func (s *SettingsSuite) Test_Settings_GetEmacsKeyBindings(c *C) {
	defset := &mockSettings{
		getBooleanRet: []bool{false},
	}
	set := &Settings{
		def: defset,
	}

	ret := set.GetEmacsKeyBindings()
	c.Assert(ret, Equals, false)
	c.Assert(defset.getBooleanArg, DeepEquals, []string{"emacs-keyboard"})
}

func (s *SettingsSuite) Test_Settings_SetEmacsKeyBindings(c *C) {
	defset := &mockSettings{}
	specset := &mockSettings{}
	set := &Settings{
		def:  defset,
		spec: specset,
	}

	set.SetEmacsKeyBindings(true)
	c.Assert(specset.setBooleanArg1, DeepEquals, []string{"has-set-emacs-keyboard", "emacs-keyboard"})
	c.Assert(specset.setBooleanArg2, DeepEquals, []bool{true, true})
}

func (s *SettingsSuite) Test_Settings_GetShowEmptyGroups(c *C) {
	defset := &mockSettings{
		getBooleanRet: []bool{false},
	}
	set := &Settings{
		def: defset,
	}

	ret := set.GetShowEmptyGroups()
	c.Assert(ret, Equals, false)
	c.Assert(defset.getBooleanArg, DeepEquals, []string{"show-empty-groups"})
}

func (s *SettingsSuite) Test_Settings_SetShowEmptyGroups(c *C) {
	defset := &mockSettings{}
	specset := &mockSettings{}
	set := &Settings{
		def:  defset,
		spec: specset,
	}

	set.SetShowEmptyGroups(true)
	c.Assert(specset.setBooleanArg1, DeepEquals, []string{"has-set-show-empty-groups", "show-empty-groups"})
	c.Assert(specset.setBooleanArg2, DeepEquals, []bool{true, true})
}

func (s *SettingsSuite) Test_Settings_GetShowAdvancedOptions(c *C) {
	defset := &mockSettings{
		getBooleanRet: []bool{false},
	}
	set := &Settings{
		def: defset,
	}

	ret := set.GetShowAdvancedOptions()
	c.Assert(ret, Equals, false)
	c.Assert(defset.getBooleanArg, DeepEquals, []string{"show-advanced-options"})
}

func (s *SettingsSuite) Test_Settings_SetShowAdvancedOptions(c *C) {
	defset := &mockSettings{}
	specset := &mockSettings{}
	set := &Settings{
		def:  defset,
		spec: specset,
	}

	set.SetShowAdvancedOptions(true)
	c.Assert(specset.setBooleanArg1, DeepEquals, []string{"has-set-show-advanced-options", "show-advanced-options"})
	c.Assert(specset.setBooleanArg2, DeepEquals, []bool{true, true})
}

func (s *SettingsSuite) Test_Settings_GetConnectAutomatically(c *C) {
	defset := &mockSettings{
		getBooleanRet: []bool{false},
	}
	set := &Settings{
		def: defset,
	}

	ret := set.GetConnectAutomatically()
	c.Assert(ret, Equals, false)
	c.Assert(defset.getBooleanArg, DeepEquals, []string{"connect-automatically"})
}

func (s *SettingsSuite) Test_Settings_SetConnectAutomatically(c *C) {
	defset := &mockSettings{}
	specset := &mockSettings{}
	set := &Settings{
		def:  defset,
		spec: specset,
	}

	set.SetConnectAutomatically(true)
	c.Assert(specset.setBooleanArg1, DeepEquals, []string{"has-set-connect-automatically", "connect-automatically"})
	c.Assert(specset.setBooleanArg2, DeepEquals, []bool{true, true})
}
