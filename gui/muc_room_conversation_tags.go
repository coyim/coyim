package gui

import (
	"fmt"

	"github.com/coyim/gotk3adapter/gtki"
	"github.com/coyim/gotk3adapter/pangoi"
)

type conversationTag string

const (
	conversationTagTimestamp         conversationTag = "timestamp"
	conversationTagMessage           conversationTag = "message"
	conversationTagNickname          conversationTag = "nickname"
	conversationTagSomeoneLeftRoom   conversationTag = "occupantLeftRoom"
	conversationTagSomeoneJoinedRoom conversationTag = "occupantJoinedRoom"
	conversationTagRoomSubject       conversationTag = "subject"
	conversationTagRoomConfigChange  conversationTag = "roomConfigChange"
	conversationTagDateGroup         conversationTag = "dateGroup"
	conversationTagDivider           conversationTag = "divider"
	conversationTagPassword          conversationTag = "password"
	conversationTagInfo              conversationTag = "info"
	conversationTagWarning           conversationTag = "warning"
	conversationTagError             conversationTag = "error"
)

const converstationLineSpacing = 12

var conversationTagsPropertiesRegistry = map[conversationTag]pangoAttributes{
	conversationTagTimestamp: {
		"style": pangoFontStyleNormal,
	},
	conversationTagMessage: {
		"style": pangoFontStyleNormal,
	},
	conversationTagNickname: {
		"style":  pangoFontStyleNormal,
		"weight": pangoFontWeightBold,
	},
	conversationTagSomeoneLeftRoom: {
		"style": pangoFontStyleItalic,
	},
	conversationTagSomeoneJoinedRoom: {
		"style": pangoFontStyleItalic,
	},
	conversationTagRoomSubject: {
		"style": pangoFontStyleItalic,
	},
	conversationTagRoomConfigChange: {
		"style": pangoFontStyleItalic,
	},
	conversationTagDateGroup: {
		"justification":      pangoJustifyCenter,
		"pixels-above-lines": converstationLineSpacing,
		"pixels-below-lines": converstationLineSpacing,
		"style":              pangoFontStyleItalic,
	},
	conversationTagDivider: {
		"justification":      pangoJustifyCenter,
		"pixels-above-lines": converstationLineSpacing,
		"pixels-below-lines": converstationLineSpacing,
	},
	conversationTagPassword: {
		"style": pangoFontStyleNormal,
	},
	conversationTagInfo: {
		"style": pangoFontStyleItalic,
	},
	conversationTagWarning: {
		"style": pangoFontStyleNormal,
	},
	conversationTagError: {
		"style": pangoFontStyleNormal,
	},
}

func (c *roomViewConversation) createConversationTag(name conversationTag, properties pangoAttributes) gtki.TextTag {
	tag, _ := g.gtk.TextTagNew(string(name))
	for attribute, value := range properties {
		_ = tag.SetProperty(attribute, value)
	}
	return tag
}

func (c *roomViewConversation) createWarningTag(cs mucColorSet) gtki.TextTag {
	return c.createConversationTag(conversationTagWarning, pangoAttributes{
		"foreground": cs.warningForeground,
	})
}

func (c *roomViewConversation) createInfoMessageTag(cs mucColorSet) gtki.TextTag {
	return c.createConversationTag(conversationTagInfo, pangoAttributes{
		"foreground": cs.infoMessageForeground,
		"style":      pangoi.STYLE_ITALIC,
	})
}

func (c *roomViewConversation) createLeftRoomTag(cs mucColorSet) gtki.TextTag {
	return c.createConversationTag(conversationTagSomeoneLeftRoom, pangoAttributes{
		"foreground": cs.someoneLeftForeground,
		"style":      pangoi.STYLE_ITALIC,
	})
}

func (c *roomViewConversation) createJoinedRoomTag(cs mucColorSet) gtki.TextTag {
	return c.createConversationTag(conversationTagSomeoneJoinedRoom, pangoAttributes{
		"foreground": cs.someoneJoinedForeground,
		"style":      pangoi.STYLE_ITALIC,
	})
}

func (c *roomViewConversation) createTimestampTag(cs mucColorSet) gtki.TextTag {
	return c.createConversationTag(conversationTagTimestamp, pangoAttributes{
		"foreground": cs.timestampForeground,
		"style":      pangoi.STYLE_NORMAL,
	})
}

func (c *roomViewConversation) createNicknameTag(cs mucColorSet) gtki.TextTag {
	return c.createConversationTag(conversationTagNickname, pangoAttributes{
		"foreground": cs.nicknameForeground,
		"style":      pangoi.STYLE_NORMAL,
		"weight":     pangoi.WEIGHT_BOLD,
	})
}

func (c *roomViewConversation) createSubjectTag(cs mucColorSet) gtki.TextTag {
	return c.createConversationTag(conversationTagRoomSubject, pangoAttributes{
		"foreground": cs.subjectForeground,
		"style":      pangoi.STYLE_ITALIC,
	})
}

func (c *roomViewConversation) createMessageTag(cs mucColorSet) gtki.TextTag {
	return c.createConversationTag(conversationTagMessage, pangoAttributes{
		"foreground": cs.messageForeground,
		"style":      pangoi.STYLE_NORMAL,
	})
}

func (c *roomViewConversation) createGroupDateTag(cs mucColorSet) gtki.TextTag {
	return c.createConversationTag(conversationTagDateGroup, pangoAttributes{
		"justification":      pangoi.JUSTIFY_CENTER,
		"pixels-above-lines": 12,
		"pixels-below-lines": 12,
		"foreground":         cs.infoMessageForeground,
		"style":              pangoi.STYLE_ITALIC,
	})
}

func (c *roomViewConversation) createDividerTag(cs mucColorSet) gtki.TextTag {
	return c.createConversationTag(conversationTagDivider, pangoAttributes{
		"justification":      pangoi.JUSTIFY_CENTER,
		"pixels-above-lines": 12,
		"pixels-below-lines": 12,
	})
}

func (c *roomViewConversation) createErrorTag(cs mucColorSet) gtki.TextTag {
	return c.createConversationTag(conversationTagError, pangoAttributes{
		"foreground": cs.errorForeground,
	})
}

func (c *roomViewConversation) createConfigurationChangeTag(cs mucColorSet) gtki.TextTag {
	return c.createConversationTag(conversationTagRoomConfigChange, pangoAttributes{
		"foreground": cs.configurationForeground,
		"style":      pangoi.STYLE_ITALIC,
	})
}

func (c *roomViewConversation) createPasswordTag(cs mucColorSet) gtki.TextTag {
	return c.createConversationTag(conversationTagPassword, pangoAttributes{
		"foreground": cs.warningForeground,
		"background": cs.warningBackground,
	})
}

var uniqueHighlightFormatName = func(format string) conversationTag {
	return conversationTag(fmt.Sprintf("f%s", format))
}

var (
	conversationTagFormatNickame     conversationTag = uniqueHighlightFormatName(highlightFormatNickname)
	conversationTagFormatAffiliation conversationTag = uniqueHighlightFormatName(highlightFormatAffiliation)
	conversationTagFormatRole        conversationTag = uniqueHighlightFormatName(highlightFormatRole)
)

var conversationTagFormats = map[conversationTag]pangoAttributes{
	conversationTagFormatNickame: {
		"weight": pangoFontWeightBold,
		"style":  pangoFontStyleItalic,
	},
	conversationTagFormatAffiliation: {
		"weight": pangoFontWeightBold,
		"style":  pangoFontStyleItalic,
	},
	conversationTagFormatRole: {
		"weight": pangoFontWeightBold,
		"style":  pangoFontStyleItalic,
	},
}

// createTextFormatTags MUST be called from the UI thread
func (c *roomViewConversation) createTextFormatTags(table gtki.TextTagTable) {
	for tagFormat, properties := range conversationTagFormats {
		tag := c.createConversationTag(tagFormat, pangoAttributesNormalize(properties))
		table.Add(tag)
	}
}

func (c *roomViewConversation) newMUCTableStyleTags(u *gtkUI) gtki.TextTagTable {
	cs := u.currentMUCColorSet()

	table, _ := g.gtk.TextTagTableNew()
	for tagName, properties := range conversationTagsPropertiesRegistry {
		tag := c.createConversationTag(tagName, pangoAttributesNormalize(properties))
		c.applyTagColors(tagName, tag, cs)
		table.Add(tag)
	}

	c.createTextFormatTags(table)

	return table
}

// applyTagColors MUST be called from the UI thread
func (c *roomViewConversation) applyTagColors(tagName conversationTag, tag gtki.TextTag, cs mucColorSet) {
	colors := conversationTagColorDefinition(tagName, cs)
	colors.applyToTag(tag)
}

type conversationTagColor struct {
	foreground string
	background string
}

// applyToTag MUST be called from the UI thread
func (tc *conversationTagColor) applyToTag(tag gtki.TextTag) {
	tc.applyTagColor("foreground", tc.foreground, tag)
	tc.applyTagColor("background", tc.background, tag)
}

// applyTagColor MUST be called from the UI thread
func (tc *conversationTagColor) applyTagColor(property, color string, tag gtki.TextTag) {
	if color != "" {
		tag.SetProperty(property, color)
	}
}

var defaultConversationTagColor = &conversationTagColor{}

func conversationTagColorDefinition(tagName conversationTag, cs mucColorSet) *conversationTagColor {
	switch tagName {
	case conversationTagWarning:
		return &conversationTagColor{
			foreground: cs.warningForeground,
		}

	case conversationTagInfo:
		return &conversationTagColor{
			foreground: cs.infoMessageForeground,
		}

	case conversationTagSomeoneLeftRoom:
		return &conversationTagColor{
			foreground: cs.someoneLeftForeground,
		}

	case conversationTagSomeoneJoinedRoom:
		return &conversationTagColor{
			foreground: cs.someoneJoinedForeground,
		}

	case conversationTagTimestamp:
		return &conversationTagColor{
			foreground: cs.timestampForeground,
		}

	case conversationTagNickname:
		return &conversationTagColor{
			foreground: cs.nicknameForeground,
		}

	case conversationTagRoomSubject:
		return &conversationTagColor{
			foreground: cs.subjectForeground,
		}

	case conversationTagMessage:
		return &conversationTagColor{
			foreground: cs.messageForeground,
		}

	case conversationTagDateGroup:
		return &conversationTagColor{
			foreground: cs.infoMessageForeground,
		}

	case conversationTagError:
		return &conversationTagColor{
			foreground: cs.errorForeground,
		}

	case conversationTagRoomConfigChange:
		return &conversationTagColor{
			foreground: cs.configurationForeground,
		}

	case conversationTagPassword:
		return &conversationTagColor{
			foreground: cs.warningForeground,
			background: cs.warningBackground,
		}
	}

	return defaultConversationTagColor
}
