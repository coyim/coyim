package gui

import (
	"fmt"

	"github.com/coyim/gotk3adapter/gtki"
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

var (
	conversationTagFormatInfoNickname        = formattingTagName(highlightFormatNickname, conversationTagInfo)
	conversationTagFormatInfoAffiliation     = formattingTagName(highlightFormatAffiliation, conversationTagInfo)
	conversationTagFormatInfoRole            = formattingTagName(highlightFormatRole, conversationTagInfo)
	conversationTagFormatLeftRoomNickname    = formattingTagName(highlightFormatNickname, conversationTagSomeoneLeftRoom)
	conversationTagFormatJoinedRoomNickname  = formattingTagName(highlightFormatNickname, conversationTagSomeoneJoinedRoom)
	conversationTagFormatRoomSubjectNickname = formattingTagName(highlightFormatNickname, conversationTagRoomSubject)
)

func formattingTagName(format string, tag conversationTag) conversationTag {
	return conversationTag(fmt.Sprintf("%s_%s", format, tag))
}

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
	conversationTagFormatLeftRoomNickname: {
		"weight": pangoFontWeightBold,
		"style":  pangoFontStyleItalic,
	},
	conversationTagSomeoneJoinedRoom: {
		"style": pangoFontStyleItalic,
	},
	conversationTagFormatJoinedRoomNickname: {
		"weight": pangoFontWeightBold,
		"style":  pangoFontStyleItalic,
	},
	conversationTagRoomSubject: {
		"style": pangoFontStyleItalic,
	},
	conversationTagFormatRoomSubjectNickname: {
		"weight": pangoFontWeightBold,
		"style":  pangoFontStyleItalic,
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
	conversationTagFormatInfoNickname: {
		"weight": pangoFontWeightBold,
		"style":  pangoFontStyleItalic,
	},
	conversationTagFormatInfoAffiliation: {
		"weight": pangoFontWeightBold,
		"style":  pangoFontStyleItalic,
	},
	conversationTagFormatInfoRole: {
		"weight": pangoFontWeightBold,
		"style":  pangoFontStyleItalic,
	},
	conversationTagWarning: {
		"style": pangoFontStyleNormal,
	},
	conversationTagError: {
		"style": pangoFontStyleNormal,
	},
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

type conversationTagColors map[conversationTag]*conversationTagColor

type conversationTags struct {
	table         gtki.TextTagTable
	temporaryTags map[conversationTag]gtki.TextTag
}

func (c *roomViewConversation) newConversationTags() *conversationTags {
	ct := &conversationTags{
		temporaryTags: map[conversationTag]gtki.TextTag{},
	}

	ct.initTagTable(c.u)

	return ct
}

func (ct *conversationTags) initTagTable(u *gtkUI) {
	table, _ := g.gtk.TextTagTableNew()
	ct.table = table

	cs := u.currentMUCColorSet()
	for tagName, properties := range conversationTagsPropertiesRegistry {
		tag := ct.createTag(tagName, properties)
		ct.applyTagColors(tagName, tag, cs)
		ct.table.Add(tag)
	}
}

// createTag MUST be called from the UI thread
func (ct *conversationTags) createTag(name conversationTag, properties pangoAttributes) gtki.TextTag {
	tag, _ := g.gtk.TextTagNew(string(name))
	for attribute, value := range pangoAttributesNormalize(properties) {
		_ = tag.SetProperty(attribute, value)
	}
	return tag
}

// applyTagColors MUST be called from the UI thread
func (ct *conversationTags) applyTagColors(tagName conversationTag, tag gtki.TextTag, cs mucColorSet) {
	c := conversationTagColorDefinition(tagName, cs)
	c.applyToTag(tag)
}

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

	case conversationTagSomeoneLeftRoom,
		conversationTagFormatLeftRoomNickname:
		return &conversationTagColor{
			foreground: cs.someoneLeftForeground,
		}

	case conversationTagSomeoneJoinedRoom,
		conversationTagFormatJoinedRoomNickname:
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

	case conversationTagRoomSubject,
		conversationTagFormatRoomSubjectNickname:
		return &conversationTagColor{
			foreground: cs.subjectForeground,
		}

	case conversationTagMessage:
		return &conversationTagColor{
			foreground: cs.messageForeground,
		}

	case conversationTagDateGroup,
		conversationTagFormatInfoNickname,
		conversationTagFormatInfoAffiliation,
		conversationTagFormatInfoRole:
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

	return &conversationTagColor{}
}
