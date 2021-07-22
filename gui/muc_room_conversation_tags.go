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
	conversationTagUnd               conversationTag = ""
)

const formattingTagText = "%sFormatting"

func formattingTagName(format string) conversationTag {
	return conversationTag(fmt.Sprintf(formattingTagText, format))
}

var (
	conversationTagFormatNickame     conversationTag = formattingTagName(highlightFormatNickname)
	conversationTagFormatAffiliation conversationTag = formattingTagName(highlightFormatAffiliation)
	conversationTagFormatRole        conversationTag = formattingTagName(highlightFormatRole)
)

type conversationTagsFormatsList []conversationTag

func (list conversationTagsFormatsList) tagForFormat(format string) (conversationTag, bool) {
	ids := []conversationTag{conversationTag(format), formattingTagName(format)}
	for _, id := range ids {
		if list.includes(id) {
			return id, true
		}
	}
	return conversationTagUnd, false
}

func (list conversationTagsFormatsList) includes(format conversationTag) bool {
	for _, f := range list {
		if f == format {
			return true
		}
	}
	return false
}

var conversationTagFormats = conversationTagsFormatsList{
	conversationTagFormatNickame,
	conversationTagFormatAffiliation,
	conversationTagFormatRole,
	conversationTagPassword,
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

func (ctc conversationTagColors) includes(tag conversationTag) bool {
	_, ok := ctc[tag]
	return ok
}

type conversationTags struct {
	table  gtki.TextTagTable
	colors conversationTagColors
}

func (c *roomViewConversation) newConversationTags() *conversationTags {
	ct := &conversationTags{
		colors: conversationTagColors{},
	}

	ct.initTagTable(c.u)

	return ct
}

func (ct *conversationTags) initTagTable(u *gtkUI) {
	table, _ := g.gtk.TextTagTableNew()
	ct.table = table

	cs := u.currentMUCColorSet()
	for tagName, properties := range conversationTagsPropertiesRegistry {
		tag := ct.createTag(tagName, pangoAttributesNormalize(properties))
		ct.applyTagColors(tagName, tag, cs)
		ct.table.Add(tag)
	}
}

// createTag MUST be called from the UI thread
func (ct *conversationTags) createTag(name conversationTag, properties pangoAttributes) gtki.TextTag {
	tag, _ := g.gtk.TextTagNew(string(name))
	for attribute, value := range properties {
		_ = tag.SetProperty(attribute, value)
	}
	return tag
}

// applyTagColors MUST be called from the UI thread
func (ct *conversationTags) applyTagColors(tagName conversationTag, tag gtki.TextTag, cs mucColorSet) {
	c := conversationTagColorDefinition(tagName, cs)
	c.applyToTag(tag)
	ct.colors[tagName] = c
}

func (ct *conversationTags) colorDefinitionFor(tag conversationTag) *conversationTagColor {
	ret := &conversationTagColor{}
	if ct.colors.includes(tag) {
		ret = ct.colors[tag]
	}
	return ret
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

	return &conversationTagColor{}
}
