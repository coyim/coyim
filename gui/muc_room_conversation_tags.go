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

func (c *roomViewConversation) createConversationTag(name conversationTag, properties map[string]interface{}) gtki.TextTag {
	tag, _ := g.gtk.TextTagNew(string(name))
	for attribute, value := range properties {
		_ = tag.SetProperty(attribute, value)
	}
	return tag
}

func (c *roomViewConversation) createWarningTag(cs mucColorSet) gtki.TextTag {
	return c.createConversationTag(conversationTagWarning, map[string]interface{}{
		"foreground": cs.warningForeground,
	})
}

func (c *roomViewConversation) createInfoMessageTag(cs mucColorSet) gtki.TextTag {
	return c.createConversationTag(conversationTagInfo, map[string]interface{}{
		"foreground": cs.infoMessageForeground,
		"style":      pangoi.STYLE_ITALIC,
	})
}

func (c *roomViewConversation) createLeftRoomTag(cs mucColorSet) gtki.TextTag {
	return c.createConversationTag(conversationTagSomeoneLeftRoom, map[string]interface{}{
		"foreground": cs.someoneLeftForeground,
		"style":      pangoi.STYLE_ITALIC,
	})
}

func (c *roomViewConversation) createJoinedRoomTag(cs mucColorSet) gtki.TextTag {
	return c.createConversationTag(conversationTagSomeoneJoinedRoom, map[string]interface{}{
		"foreground": cs.someoneJoinedForeground,
		"style":      pangoi.STYLE_ITALIC,
	})
}

func (c *roomViewConversation) createTimestampTag(cs mucColorSet) gtki.TextTag {
	return c.createConversationTag(conversationTagTimestamp, map[string]interface{}{
		"foreground": cs.timestampForeground,
		"style":      pangoi.STYLE_NORMAL,
	})
}

func (c *roomViewConversation) createNicknameTag(cs mucColorSet) gtki.TextTag {
	return c.createConversationTag(conversationTagNickname, map[string]interface{}{
		"foreground": cs.nicknameForeground,
		"style":      pangoi.STYLE_NORMAL,
		"weight":     pangoi.WEIGHT_BOLD,
	})
}

func (c *roomViewConversation) createSubjectTag(cs mucColorSet) gtki.TextTag {
	return c.createConversationTag(conversationTagRoomSubject, map[string]interface{}{
		"foreground": cs.subjectForeground,
		"style":      pangoi.STYLE_ITALIC,
	})
}

func (c *roomViewConversation) createMessageTag(cs mucColorSet) gtki.TextTag {
	return c.createConversationTag(conversationTagMessage, map[string]interface{}{
		"foreground": cs.messageForeground,
		"style":      pangoi.STYLE_NORMAL,
	})
}

func (c *roomViewConversation) createGroupDateTag(cs mucColorSet) gtki.TextTag {
	return c.createConversationTag(conversationTagDateGroup, map[string]interface{}{
		"justification":      pangoi.JUSTIFY_CENTER,
		"pixels-above-lines": 12,
		"pixels-below-lines": 12,
		"foreground":         cs.infoMessageForeground,
		"style":              pangoi.STYLE_ITALIC,
	})
}

func (c *roomViewConversation) createDividerTag(cs mucColorSet) gtki.TextTag {
	return c.createConversationTag(conversationTagDivider, map[string]interface{}{
		"justification":      pangoi.JUSTIFY_CENTER,
		"pixels-above-lines": 12,
		"pixels-below-lines": 12,
	})
}

func (c *roomViewConversation) createErrorTag(cs mucColorSet) gtki.TextTag {
	return c.createConversationTag(conversationTagError, map[string]interface{}{
		"foreground": cs.errorForeground,
	})
}

func (c *roomViewConversation) createConfigurationChangeTag(cs mucColorSet) gtki.TextTag {
	return c.createConversationTag(conversationTagRoomConfigChange, map[string]interface{}{
		"foreground": cs.configurationForeground,
		"style":      pangoi.STYLE_ITALIC,
	})
}

func (c *roomViewConversation) createPasswordTag(cs mucColorSet) gtki.TextTag {
	return c.createConversationTag(conversationTagPassword, map[string]interface{}{
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

var conversationTagFormats = map[conversationTag]map[string]interface{}{
	conversationTagFormatNickame: {
		"foreground": "blue",
		"weight":     pangoFontWeightBold,
		"style":      pangoFontStyleItalic,
	},
	conversationTagFormatAffiliation: {
		"foreground": "red",
		"weight":     pangoFontWeightBold,
		"style":      pangoFontStyleItalic,
	},
	conversationTagFormatRole: {
		"foreground": "orange",
		"weight":     pangoFontWeightBold,
		"style":      pangoFontStyleItalic,
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
