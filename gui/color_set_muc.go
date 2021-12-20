package gui

type mucColorSet struct {
	warningForeground                      cssColor
	warningBackground                      cssColor
	someoneJoinedForeground                cssColor
	someoneLeftForeground                  cssColor
	timestampForeground                    cssColor
	nicknameForeground                     cssColor
	subjectForeground                      cssColor
	infoMessageForeground                  cssColor
	messageForeground                      cssColor
	errorForeground                        cssColor
	configurationForeground                cssColor
	roomMessagesBackground                 cssColor
	roomMessagesBoxShadow                  cssColor
	roomNameDisabledForeground             cssColor
	roomSubjectForeground                  cssColor
	roomOverlaySolidBackground             cssColor
	roomOverlayContentSolidBackground      cssColor
	roomOverlayContentBackground           cssColor
	roomOverlayBackground                  cssColor
	roomOverlayContentForeground           cssColor
	roomOverlayContentBoxShadow            cssColor
	roomWarningsDialogBackground           cssColor
	roomWarningsDialogDecorationBackground cssColor
	roomWarningsDialogDecorationShadow     cssColor
	roomWarningsDialogHeaderBackground     cssColor
	roomWarningsDialogContentBackground    cssColor
	roomWarningsCurrentInfoForeground      cssColor
	roomNotificationsBackground            cssColor
	rosterGroupBackground                  cssColor
	rosterGroupForeground                  cssColor
	rosterOccupantRoleForeground           cssColor
	occupantStatusAvailableForeground      cssColor
	occupantStatusAvailableBackground      cssColor
	occupantStatusAvailableBorder          cssColor
	occupantStatusNotAvailableForeground   cssColor
	occupantStatusNotAvailableBackground   cssColor
	occupantStatusNotAvailableBorder       cssColor
	occupantStatusAwayForeground           cssColor
	occupantStatusAwayBackground           cssColor
	occupantStatusAwayBorder               cssColor
	occupantStatusBusyForeground           cssColor
	occupantStatusBusyBackground           cssColor
	occupantStatusBusyBorder               cssColor
	occupantStatusFreeForChatForeground    cssColor
	occupantStatusFreeForChatBackground    cssColor
	occupantStatusFreeForChatBorder        cssColor
	occupantStatusExtendedAwayForeground   cssColor
	occupantStatusExtendedAwayBackground   cssColor
	occupantStatusExtendedAwayBorder       cssColor
	infoBarDefaultBorderColor              cssColor
	infoBarTypeInfoBackgroundStart         cssColor
	infoBarTypeInfoBackgroundStop          cssColor
	infoBarTypeInfoTitle                   cssColor
	infoBarTypeInfoTime                    cssColor
	infoBarTypeWarningBackgroundStart      cssColor
	infoBarTypeWarningBackgroundStop       cssColor
	infoBarTypeWarningTitle                cssColor
	infoBarTypeWarningTime                 cssColor
	infoBarTypeQuestionBackgroundStart     cssColor
	infoBarTypeQuestionBackgroundStop      cssColor
	infoBarTypeQuestionTitle               cssColor
	infoBarTypeQuestionTime                cssColor
	infoBarTypeErrorBackgroundStart        cssColor
	infoBarTypeErrorBackgroundStop         cssColor
	infoBarTypeErrorTitle                  cssColor
	infoBarTypeErrorTime                   cssColor
	infoBarTypeOtherBackgroundStart        cssColor
	infoBarTypeOtherBackgroundStop         cssColor
	infoBarTypeOtherTitle                  cssColor
	infoBarTypeOtherTime                   cssColor
	infoBarButtonBackground                cssColor
	infoBarButtonForeground                cssColor
	infoBarButtonHoverBackground           cssColor
	infoBarButtonHoverForeground           cssColor
	infoBarButtonActiveBackground          cssColor
	infoBarButtonActiveForeground          cssColor
	entryErrorBackground                   cssColor
	entryErrorBorderShadow                 cssColor
	entryErrorBorder                       cssColor
	entryErrorLabel                        cssColor
	occupantLostConnection                 cssColor
	occupantRestablishConnection           cssColor
}

func (cm *hasColorManagement) currentMUCColorSet() mucColorSet {
	if cm.isDarkThemeVariant() {
		return defaultMUCDarkColorSet
	}
	return defaultMUCLightColorSet
}

var defaultMUCLightColorSet = mucColorSet{
	warningForeground:                      rgbFrom(194, 23, 29),
	warningBackground:                      rgbFrom(254, 202, 202),
	errorForeground:                        rgbFrom(163, 7, 7),
	someoneJoinedForeground:                rgbFrom(41, 115, 22),
	someoneLeftForeground:                  rgbFrom(115, 22, 41),
	timestampForeground:                    colorThemeInsensitiveForeground,
	nicknameForeground:                     rgbFrom(57, 91, 163),
	subjectForeground:                      rgbFrom(0, 0, 128),
	infoMessageForeground:                  rgbFrom(57, 91, 163),
	messageForeground:                      rgbFrom(0, 0, 0),
	configurationForeground:                rgbFrom(154, 4, 191),
	roomMessagesBackground:                 colorThemeBase,
	roomMessagesBoxShadow:                  rgbaFrom(0, 0, 0, 0.35),
	roomNameDisabledForeground:             colorThemeInsensitiveForeground,
	roomSubjectForeground:                  colorThemeInsensitiveForeground,
	roomOverlaySolidBackground:             colorThemeBackground,
	roomOverlayContentSolidBackground:      colorTransparent,
	roomOverlayContentBackground:           colorThemeBackground,
	roomOverlayBackground:                  rgbaFrom(0, 0, 0, 0.5),
	roomOverlayContentForeground:           colorThemeForeground,
	roomOverlayContentBoxShadow:            rgbaFrom(0, 0, 0, 0.5),
	roomWarningsDialogBackground:           colorNone,
	roomWarningsDialogDecorationBackground: colorThemeBackground,
	roomWarningsDialogDecorationShadow:     rgbaFrom(0, 0, 0, 0.15),
	roomWarningsDialogHeaderBackground:     colorNone,
	roomWarningsDialogContentBackground:    colorNone,
	roomWarningsCurrentInfoForeground:      colorThemeInsensitiveForeground,
	roomNotificationsBackground:            colorThemeBackground,
	rosterGroupBackground:                  rgbFrom(245, 245, 244),
	rosterGroupForeground:                  rgbFrom(28, 25, 23),
	rosterOccupantRoleForeground:           rgbFrom(168, 162, 158),
	occupantStatusAvailableForeground:      rgbFrom(22, 101, 52),
	occupantStatusAvailableBackground:      rgbFrom(240, 253, 244),
	occupantStatusAvailableBorder:          rgbFrom(22, 163, 74),
	occupantStatusNotAvailableForeground:   rgbFrom(30, 41, 59),
	occupantStatusNotAvailableBackground:   rgbFrom(248, 250, 252),
	occupantStatusNotAvailableBorder:       rgbFrom(71, 85, 105),
	occupantStatusAwayForeground:           rgbFrom(154, 52, 18),
	occupantStatusAwayBackground:           rgbFrom(255, 247, 237),
	occupantStatusAwayBorder:               rgbFrom(234, 88, 12),
	occupantStatusBusyForeground:           rgbFrom(159, 18, 57),
	occupantStatusBusyBackground:           rgbFrom(255, 241, 242),
	occupantStatusBusyBorder:               rgbFrom(190, 18, 60),
	occupantStatusFreeForChatForeground:    rgbFrom(30, 64, 175),
	occupantStatusFreeForChatBackground:    rgbFrom(239, 246, 255),
	occupantStatusFreeForChatBorder:        rgbFrom(29, 78, 216),
	occupantStatusExtendedAwayForeground:   rgbFrom(146, 64, 14),
	occupantStatusExtendedAwayBackground:   rgbFrom(255, 251, 235),
	occupantStatusExtendedAwayBorder:       rgbFrom(217, 119, 6),
	infoBarDefaultBorderColor:              colorThemeBackground,
	infoBarTypeInfoBackgroundStart:         rgbFrom(63, 98, 18),
	infoBarTypeInfoBackgroundStop:          rgbFrom(77, 124, 15),
	infoBarTypeInfoTitle:                   rgbFrom(236, 254, 255),
	infoBarTypeInfoTime:                    rgbaFrom(236, 254, 255, 0.5),
	infoBarTypeWarningBackgroundStart:      rgbFrom(195, 149, 7),
	infoBarTypeWarningBackgroundStop:       rgbFrom(222, 173, 20),
	infoBarTypeWarningTitle:                rgbFrom(255, 247, 237),
	infoBarTypeWarningTime:                 rgbaFrom(255, 247, 237, 0.5),
	infoBarTypeQuestionBackgroundStart:     rgbFrom(234, 88, 12),
	infoBarTypeQuestionBackgroundStop:      rgbFrom(249, 115, 22),
	infoBarTypeQuestionTitle:               rgbFrom(254, 252, 232),
	infoBarTypeQuestionTime:                rgbaFrom(254, 252, 232, 0.5),
	infoBarTypeErrorBackgroundStart:        rgbFrom(185, 28, 28),
	infoBarTypeErrorBackgroundStop:         rgbFrom(203, 35, 35),
	infoBarTypeErrorTitle:                  rgbFrom(255, 241, 242),
	infoBarTypeErrorTime:                   rgbaFrom(255, 241, 242, 0.5),
	infoBarTypeOtherBackgroundStart:        rgbFrom(7, 89, 133),
	infoBarTypeOtherBackgroundStop:         rgbFrom(3, 105, 161),
	infoBarTypeOtherTitle:                  rgbFrom(240, 253, 250),
	infoBarTypeOtherTime:                   rgbaFrom(240, 253, 250, 0.5),
	infoBarButtonBackground:                rgbaFrom(0, 0, 0, 0.25),
	infoBarButtonForeground:                rgbFrom(255, 255, 255),
	infoBarButtonHoverBackground:           rgbaFrom(0, 0, 0, 0.35),
	infoBarButtonHoverForeground:           rgbFrom(255, 255, 255),
	infoBarButtonActiveBackground:          rgbaFrom(0, 0, 0, 0.45),
	infoBarButtonActiveForeground:          rgbFrom(255, 255, 255),
	entryErrorBackground:                   rgbFrom(255, 245, 246),
	entryErrorBorderShadow:                 rgbFrom(255, 127, 80),
	entryErrorBorder:                       rgbFrom(228, 70, 53),
	entryErrorLabel:                        rgbFrom(228, 70, 53),
	occupantLostConnection:                 rgbFrom(115, 22, 41),
	occupantRestablishConnection:           rgbFrom(41, 115, 22),
}

var defaultMUCDarkColorSet = mucColorSet{
	warningForeground:                      rgbFrom(194, 23, 29),
	warningBackground:                      rgbFrom(254, 202, 202),
	errorForeground:                        rgbFrom(209, 104, 96),
	someoneJoinedForeground:                rgbFrom(41, 115, 22),
	someoneLeftForeground:                  rgbFrom(115, 22, 41),
	timestampForeground:                    colorThemeInsensitiveForeground,
	nicknameForeground:                     rgbFrom(57, 91, 163),
	subjectForeground:                      rgbFrom(0, 0, 128),
	infoMessageForeground:                  rgbFrom(227, 66, 103),
	messageForeground:                      rgbFrom(0, 0, 0),
	configurationForeground:                rgbFrom(154, 4, 191),
	roomMessagesBackground:                 colorThemeBase,
	roomMessagesBoxShadow:                  rgbaFrom(0, 0, 0, 0.35),
	roomNameDisabledForeground:             colorThemeInsensitiveForeground,
	roomSubjectForeground:                  colorThemeInsensitiveForeground,
	roomOverlaySolidBackground:             colorThemeBase,
	roomOverlayContentSolidBackground:      colorTransparent,
	roomOverlayContentBackground:           colorThemeBase,
	roomOverlayBackground:                  rgbaFrom(0, 0, 0, 0.5),
	roomOverlayContentForeground:           rgbFrom(51, 51, 51),
	roomOverlayContentBoxShadow:            rgbaFrom(0, 0, 0, 0.5),
	roomWarningsDialogBackground:           colorNone,
	roomWarningsDialogDecorationBackground: colorThemeBackground,
	roomWarningsDialogDecorationShadow:     rgbaFrom(0, 0, 0, 0.15),
	roomWarningsDialogHeaderBackground:     colorNone,
	roomWarningsDialogContentBackground:    colorNone,
	roomWarningsCurrentInfoForeground:      colorThemeInsensitiveForeground,
	roomNotificationsBackground:            colorThemeBackground,
	rosterGroupBackground:                  rgbFrom(28, 25, 23),
	rosterGroupForeground:                  rgbFrom(250, 250, 249),
	rosterOccupantRoleForeground:           rgbFrom(231, 229, 228),
	occupantStatusAvailableForeground:      rgbFrom(22, 101, 52),
	occupantStatusAvailableBackground:      rgbFrom(240, 253, 244),
	occupantStatusAvailableBorder:          rgbFrom(22, 163, 74),
	occupantStatusNotAvailableForeground:   rgbFrom(30, 41, 59),
	occupantStatusNotAvailableBackground:   rgbFrom(248, 250, 252),
	occupantStatusNotAvailableBorder:       rgbFrom(71, 85, 105),
	occupantStatusAwayForeground:           rgbFrom(154, 52, 18),
	occupantStatusAwayBackground:           rgbFrom(255, 247, 237),
	occupantStatusAwayBorder:               rgbFrom(234, 88, 12),
	occupantStatusBusyForeground:           rgbFrom(159, 18, 57),
	occupantStatusBusyBackground:           rgbFrom(255, 241, 242),
	occupantStatusBusyBorder:               rgbFrom(190, 18, 60),
	occupantStatusFreeForChatForeground:    rgbFrom(30, 64, 175),
	occupantStatusFreeForChatBackground:    rgbFrom(239, 246, 255),
	occupantStatusFreeForChatBorder:        rgbFrom(29, 78, 216),
	occupantStatusExtendedAwayForeground:   rgbFrom(146, 64, 14),
	occupantStatusExtendedAwayBackground:   rgbFrom(255, 251, 235),
	occupantStatusExtendedAwayBorder:       rgbFrom(217, 119, 6),
	infoBarDefaultBorderColor:              colorThemeBackground,
	infoBarTypeInfoBackgroundStart:         rgbFrom(63, 98, 18),
	infoBarTypeInfoBackgroundStop:          rgbFrom(77, 124, 15),
	infoBarTypeInfoTitle:                   rgbFrom(236, 254, 255),
	infoBarTypeInfoTime:                    rgbaFrom(236, 254, 255, 0.5),
	infoBarTypeWarningBackgroundStart:      rgbFrom(195, 149, 7),
	infoBarTypeWarningBackgroundStop:       rgbFrom(222, 173, 20),
	infoBarTypeWarningTitle:                rgbFrom(255, 247, 237),
	infoBarTypeWarningTime:                 rgbaFrom(255, 247, 237, 0.5),
	infoBarTypeQuestionBackgroundStart:     rgbFrom(234, 88, 12),
	infoBarTypeQuestionBackgroundStop:      rgbFrom(249, 115, 22),
	infoBarTypeQuestionTitle:               rgbFrom(254, 252, 232),
	infoBarTypeQuestionTime:                rgbaFrom(254, 252, 232, 0.5),
	infoBarTypeErrorBackgroundStart:        rgbFrom(185, 28, 28),
	infoBarTypeErrorBackgroundStop:         rgbFrom(203, 35, 35),
	infoBarTypeErrorTitle:                  rgbFrom(255, 241, 242),
	infoBarTypeErrorTime:                   rgbaFrom(255, 241, 242, 0.5),
	infoBarTypeOtherBackgroundStart:        rgbFrom(7, 89, 133),
	infoBarTypeOtherBackgroundStop:         rgbFrom(3, 105, 161),
	infoBarTypeOtherTitle:                  rgbFrom(240, 253, 250),
	infoBarTypeOtherTime:                   rgbaFrom(240, 253, 250, 0.5),
	infoBarButtonBackground:                rgbaFrom(0, 0, 0, 0.25),
	infoBarButtonForeground:                rgbFrom(255, 255, 255),
	infoBarButtonHoverBackground:           rgbaFrom(0, 0, 0, 0.35),
	infoBarButtonHoverForeground:           rgbFrom(255, 255, 255),
	infoBarButtonActiveBackground:          rgbaFrom(0, 0, 0, 0.45),
	infoBarButtonActiveForeground:          rgbFrom(255, 255, 255),
	entryErrorBackground:                   rgbFrom(255, 245, 246),
	entryErrorBorderShadow:                 rgbFrom(255, 127, 80),
	entryErrorBorder:                       rgbFrom(228, 70, 53),
	entryErrorLabel:                        rgbFrom(228, 70, 53),
	occupantLostConnection:                 rgbFrom(115, 22, 41),
	occupantRestablishConnection:           rgbFrom(41, 115, 22),
}
