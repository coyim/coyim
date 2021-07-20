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
	table, _ := g.gtk.TextTagTableNew()
	cs := u.currentMUCColorSet()

	tags := []func(mucColorSet) gtki.TextTag{
		c.createWarningTag,
		c.createLeftRoomTag,
		c.createJoinedRoomTag,
		c.createTimestampTag,
		c.createNicknameTag,
		c.createSubjectTag,
		c.createInfoMessageTag,
		c.createMessageTag,
		c.createGroupDateTag,
		c.createDividerTag,
		c.createErrorTag,
		c.createConfigurationChangeTag,
		c.createPasswordTag,
	}

	for _, t := range tags {
		table.Add(t(cs))
	}

	c.createTextFormatTags(table)

	return table
}
