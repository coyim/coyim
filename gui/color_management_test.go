package gui

import (
	"errors"

	. "gopkg.in/check.v1"

	"github.com/coyim/gotk3adapter/gdk_mock"
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/prashantv/gostub"
)

type ColorManagementSuite struct{}

var _ = Suite(&ColorManagementSuite{})

func (s *ColorManagementSuite) Test_hasColorManagement_setsTheThemeVariantToDarkBasedOnGTK_THEMEEnvironmentVariable(c *C) {
	defer gostub.New().SetEnv("GTK_THEME", "foo-bla-theme:dark").Reset()

	hcm := &hasColorManagement{}
	c.Assert(hcm.isDarkThemeVariant(), Equals, true)
	c.Assert(hcm.themeVariant, Equals, "dark")
}

func (s *ColorManagementSuite) Test_hasColorManagement_panicsIfCantGetGTKSettings(c *C) {
	defer gostub.New().SetEnv("GTK_THEME", "foo-bla-theme:light").Reset()
	mg := &mockedGTK{}
	defer gostub.Stub(&g, CreateGraphics(mg, nil, nil, nil, nil)).Reset()

	mg.mm.On("SettingsGetDefault").Return(nil, errors.New("something bad")).Once()

	hcm := &hasColorManagement{}
	c.Assert(hcm.actuallyCalculateThemeVariant, PanicMatches, "something bad")
}

func (s *ColorManagementSuite) Test_hasColorManagement_setsTheThemeVariantToDarkBasedOnGTKSettings(c *C) {
	defer gostub.New().SetEnv("GTK_THEME", "foo-bla-theme2").Reset()
	mg := &mockedGTK{}
	defer gostub.Stub(&g, CreateGraphics(mg, nil, nil, nil, nil)).Reset()

	ms := &mockedSettings{}
	mg.mm.On("SettingsGetDefault").Return(ms, nil).Once()

	ms.mm.On("GetProperty", "gtk-application-prefer-dark-theme").Return(true, nil).Once()

	hcm := &hasColorManagement{}
	c.Assert(hcm.isDarkThemeVariant(), Equals, true)
	c.Assert(hcm.themeVariant, Equals, "dark")
}

func (s *ColorManagementSuite) Test_hasColorManagement_setsTheThemeVariantToDarkDetectingOnInvisibleBox(c *C) {
	defer gostub.New().SetEnv("GTK_THEME", "").Reset()
	mg := &mockedGTK{}
	defer gostub.Stub(&g, CreateGraphics(mg, nil, nil, nil, nil)).Reset()

	ms := &mockedSettings{}
	mg.mm.On("SettingsGetDefault").Return(ms, nil)
	ms.mm.On("GetProperty", "gtk-application-prefer-dark-theme").Return(false, nil).Once()
	ms.mm.On("GetProperty", "gtk-theme-name").Return("", nil).Once()

	mb := &mockedBuilder{}
	mg.mm.On("BuilderNew").Return(mb, nil).Once()

	mlb := &mockedListBox{}
	mb.mm.On("GetObject", "bg-color-detection-invisible-listbox").Return(mlb, nil).Once()

	msc := &mockedStyleContext{}
	mlb.mm.On("GetStyleContext").Return(msc, nil).Once()

	msc.mm.On("GetProperty2", "background-color", gtki.STATE_FLAG_NORMAL).Return(&gdk_mock.MockRgba{}, nil).Once()

	hcm := &hasColorManagement{}
	c.Assert(hcm.isDarkThemeVariant(), Equals, true)
	c.Assert(hcm.themeVariant, Equals, "dark")
}

func (s *ColorManagementSuite) Test_hasColorManagement_setsTheThemeVariantToLightIfNoStrategiesLeadToIndicationOfDarkTheme(c *C) {
	defer gostub.New().SetEnv("GTK_THEME", "").Reset()
	mg := &mockedGTK{}
	defer gostub.Stub(&g, CreateGraphics(mg, nil, nil, nil, nil)).Reset()

	ms := &mockedSettings{}
	mg.mm.On("SettingsGetDefault").Return(ms, nil)
	ms.mm.On("GetProperty", "gtk-application-prefer-dark-theme").Return(nil, nil).Once()
	ms.mm.On("GetProperty", "gtk-theme-name").Return("fluffy", nil).Once()

	mb := &mockedBuilder{}
	mg.mm.On("BuilderNew").Return(mb, nil).Once()

	mlb := &mockedListBox{}
	mb.mm.On("GetObject", "bg-color-detection-invisible-listbox").Return(mlb, nil).Once()

	msc := &mockedStyleContext{}
	mlb.mm.On("GetStyleContext").Return(msc, nil).Once()

	msc.mm.On("GetProperty2", "background-color", gtki.STATE_FLAG_NORMAL).Return(&mockRGBAWithValues{r: 1, g: 1, b: 1}, nil).Once()

	hcm := &hasColorManagement{}
	c.Assert(hcm.isDarkThemeVariant(), Equals, false)
	c.Assert(hcm.themeVariant, Equals, "light")
}

func (s *ColorManagementSuite) Test_doesThemeNameIndicateDarkness_checksForVariantBasedOnColon(c *C) {
	c.Assert(doesThemeNameIndicateDarkness("foo:dark"), Equals, true)
	c.Assert(doesThemeNameIndicateDarkness(":dark"), Equals, true)
	c.Assert(doesThemeNameIndicateDarkness("dark"), Equals, false)
	c.Assert(doesThemeNameIndicateDarkness("blaa"), Equals, false)
	c.Assert(doesThemeNameIndicateDarkness("blaa:light"), Equals, false)
	c.Assert(doesThemeNameIndicateDarkness(""), Equals, false)
	c.Assert(doesThemeNameIndicateDarkness("something:dark:light"), Equals, false)
}

func (s *ColorManagementSuite) Test_doesThemeNameIndicateDarkness_checksForVariantBasedOnUnderscore(c *C) {
	c.Assert(doesThemeNameIndicateDarkness("dark"), Equals, false)
	c.Assert(doesThemeNameIndicateDarkness("blaa"), Equals, false)
	c.Assert(doesThemeNameIndicateDarkness("blaa_light"), Equals, false)
	c.Assert(doesThemeNameIndicateDarkness("something_dark_light"), Equals, false)
	c.Assert(doesThemeNameIndicateDarkness("something_dark:light"), Equals, false)
	c.Assert(doesThemeNameIndicateDarkness("something:dark_light"), Equals, false)
	c.Assert(doesThemeNameIndicateDarkness("foo_dark"), Equals, true)
	c.Assert(doesThemeNameIndicateDarkness("_dark"), Equals, true)
	c.Assert(doesThemeNameIndicateDarkness("something:light_dark"), Equals, true)
}

func (s *ColorManagementSuite) Test_doesThemeNameIndicateDarkness_checksForVariantBasedOnDash(c *C) {
	c.Assert(doesThemeNameIndicateDarkness("dark"), Equals, false)
	c.Assert(doesThemeNameIndicateDarkness("blaa"), Equals, false)
	c.Assert(doesThemeNameIndicateDarkness("blaa-light"), Equals, false)
	c.Assert(doesThemeNameIndicateDarkness("something-dark-light"), Equals, false)
	c.Assert(doesThemeNameIndicateDarkness("something-dark:light"), Equals, false)
	c.Assert(doesThemeNameIndicateDarkness("something:dark-light"), Equals, false)
	c.Assert(doesThemeNameIndicateDarkness("foo-dark"), Equals, true)
	c.Assert(doesThemeNameIndicateDarkness("-dark"), Equals, true)
	c.Assert(doesThemeNameIndicateDarkness("something:light-dark"), Equals, true)
}

func (s *ColorManagementSuite) Test_hasColorManagement_getThemeNameFromGTKSettings_returnsThemeName(c *C) {
	mg := &mockedGTK{}
	defer gostub.Stub(&g, CreateGraphics(mg, nil, nil, nil, nil)).Reset()

	ms := &mockedSettings{}
	mg.mm.On("SettingsGetDefault").Return(ms, nil).Once()

	ms.mm.On("GetProperty", "gtk-theme-name").Return("something nice", nil).Once()

	hcm := &hasColorManagement{}
	c.Assert(hcm.getThemeNameFromGTKSettings(), Equals, "something nice")
}

func (s *ColorManagementSuite) Test_hasColorManagement_setsTheThemeVariantToDarkBasedOnThemeNameInGTKSettings(c *C) {
	defer gostub.New().SetEnv("GTK_THEME", "foo-bla-theme2").Reset()
	mg := &mockedGTK{}
	defer gostub.Stub(&g, CreateGraphics(mg, nil, nil, nil, nil)).Reset()

	ms := &mockedSettings{}
	mg.mm.On("SettingsGetDefault").Return(ms, nil)

	ms.mm.On("GetProperty", "gtk-application-prefer-dark-theme").Return(false, nil).Once()
	ms.mm.On("GetProperty", "gtk-theme-name").Return("adwaita_dark", nil).Once()

	hcm := &hasColorManagement{}
	c.Assert(hcm.isDarkThemeVariant(), Equals, true)
	c.Assert(hcm.themeVariant, Equals, "dark")
}

// func (s *ColorManagementSuite) Test_hasColorManagement_setsTheThemeVariantToDarkBasedOnThemeNameInGTKSettings(c *C) {
// 	defer gostub.New().SetEnv("GTK_THEME", "foo-bla-theme2").Reset()
// 	mg := &mockedGTK{}
// 	defer gostub.Stub(&g, CreateGraphics(mg, nil, nil, nil, nil)).Reset()

// 	ms := &mockedSettings{}
// 	mg.mm.On("SettingsGetDefault").Return(ms, nil)

// 	ms.mm.On("GetProperty", "gtk-application-prefer-dark-theme").Return(false, nil).Once()
// 	ms.mm.On("GetProperty", "gtk-theme-name").Return("adwaita", nil).Once()

// 	hcm := &hasColorManagement{}
// 	c.Assert(hcm.isDarkThemeVariant(), Equals, true)
// 	c.Assert(hcm.themeVariant, Equals, "dark")
// }
