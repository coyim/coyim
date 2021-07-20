package gui

import (
	"github.com/coyim/gotk3adapter/pangoi"
)

type pangoAttributes map[string]interface{}

type pangoFontWeightType int

const (
	pangoFontWeightThin pangoFontWeightType = iota
	pangoFontWeightUltralight
	pangoFontWeightLight
	pangoFontWeightSemiLight
	pangoFontWeightBook
	pangoFontWeightNormal
	pangoFontWeightMedium
	pangoFontWeightSemibold
	pangoFontWeightBold
	pangoFontWeightUltraBold
	pangoFontWeightHeavy
	pangoFontWeightUltraHeavy
)

func pangoFontWeight(weight pangoFontWeightType) int {
	switch weight {
	case pangoFontWeightThin:
		return pangoi.WEIGHT_THIN
	case pangoFontWeightUltralight:
		return pangoi.WEIGHT_ULTRALIGHT
	case pangoFontWeightLight:
		return pangoi.WEIGHT_LIGHT
	case pangoFontWeightSemiLight:
		return pangoi.WEIGHT_SEMILIGHT
	case pangoFontWeightBook:
		return pangoi.WEIGHT_BOOK
	case pangoFontWeightMedium:
		return pangoi.WEIGHT_MEDIUM
	case pangoFontWeightSemibold:
		return pangoi.WEIGHT_SEMIBOLD
	case pangoFontWeightBold:
		return pangoi.WEIGHT_BOLD
	case pangoFontWeightUltraBold:
		return pangoi.WEIGHT_ULTRABOLD
	case pangoFontWeightHeavy:
		return pangoi.WEIGHT_HEAVY
	case pangoFontWeightUltraHeavy:
		return pangoi.WEIGHT_ULTRAHEAVY
	}
	return pangoi.WEIGHT_NORMAL
}

type pangoFontStyleType int

const (
	pangoFontStyleNormal pangoFontStyleType = iota
	pangoFontStyleOblique
	pangoFontStyleItalic
)

func pangoFontStyle(style pangoFontStyleType) int {
	switch style {
	case pangoFontStyleOblique:
		return pangoi.STYLE_OBLIQUE
	case pangoFontStyleItalic:
		return pangoi.STYLE_ITALIC
	}
	return pangoi.STYLE_NORMAL
}

func pangoAttributesNormalize(properties pangoAttributes) pangoAttributes {
	ret := pangoAttributes{}

	for prop, value := range properties {
		switch tp := value.(type) {
		case pangoFontWeightType:
			ret[prop] = pangoFontWeight(tp)
		case pangoFontStyleType:
			ret[prop] = pangoFontStyle(tp)
		default:
			ret[prop] = value
		}
	}

	return ret
}
