package gui

import (
	"strings"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

const (
	languageSelectorCodeIndex int = iota
	languageSelectorDescriptionIndex
)

type languageSelectorComponent struct {
	entry     gtki.Entry
	combo     gtki.ComboBoxText
	model     gtki.ListStore
	languages *languageSelectorValues
}

func (u *gtkUI) createLanguageSelectorComponent(entry gtki.Entry, combo gtki.ComboBoxText) *languageSelectorComponent {
	lc := &languageSelectorComponent{
		entry:     entry,
		combo:     combo,
		languages: getKnownLanguages(),
	}

	lc.initModel()
	lc.initLanguageCombo()
	lc.initLanguageEntry()

	return lc
}

func (lc *languageSelectorComponent) initModel() {
	model, _ := g.gtk.ListStoreNew(
		// language code (like "en" or "es")
		glibi.TYPE_STRING,
		// language friendly description (like "English" or "Español")
		glibi.TYPE_STRING,
	)

	for langCode, e := range lc.languages.list {
		iter := model.Append()

		_ = model.SetValue(iter, languageSelectorCodeIndex, langCode)
		_ = model.SetValue(iter, languageSelectorDescriptionIndex, e.description)
	}

	lc.model = model
}

func (lc *languageSelectorComponent) initLanguageCombo() {
	lc.combo.SetModel(lc.model)
	lc.combo.SetIDColumn(languageSelectorCodeIndex)
	lc.combo.SetEntryTextColumn(languageSelectorDescriptionIndex)
}

func (lc *languageSelectorComponent) initLanguageEntry() {
	ec, _ := g.gtk.EntryCompletionNew()
	ec.SetModel(lc.model)
	ec.SetMinimumKeyLength(1)
	ec.SetTextColumn(languageSelectorDescriptionIndex)

	lc.entry.SetCompletion(ec)
}

func (lc *languageSelectorComponent) setLanguage(t string) {
	setEntryText(lc.entry, supportedLanguageDescription(lc.languages.languageBasedOnText(t)))
}

func (lc *languageSelectorComponent) currentLanguage() string {
	return lc.languages.languageBasedOnText(getEntryText(lc.entry))
}

type languageSelectorEntry struct {
	description string
	values      []string
}

func newlanguageSelectorEntry(langDesc string) *languageSelectorEntry {
	return &languageSelectorEntry{
		description: langDesc,
		values:      []string{langDesc},
	}
}

func (e *languageSelectorEntry) contains(t string) bool {
	for _, v := range e.values {
		if strings.EqualFold(v, t) {
			return true
		}
	}
	return false
}

func (e *languageSelectorEntry) add(t ...string) {
	values := e.values
	for _, vv := range t {
		values = append(values, vv)
	}
	e.values = values
}

type languageSelectorValues struct {
	list map[string]*languageSelectorEntry
}

func newLanguageSelectorValues() *languageSelectorValues {
	return &languageSelectorValues{
		list: make(map[string]*languageSelectorEntry),
	}
}

func (v *languageSelectorValues) languageBasedOnText(t string) string {
	_, ok := v.list[t]
	if !ok {
		for tt, e := range v.list {
			if e.contains(t) {
				return tt
			}
		}
	}
	return t
}

func (v *languageSelectorValues) add(langCode string, langDesc string, values ...string) {
	e, ok := v.list[langCode]
	if !ok {
		e = newlanguageSelectorEntry(langDesc)
	}
	e.add(values...)
	v.list[langCode] = e
}

var knownLanguagesValues *languageSelectorValues

func getKnownLanguages() *languageSelectorValues {
	if knownLanguagesValues == nil {
		knownLanguagesValues = newLanguageSelectorValues()
		for _, tag := range display.Supported.Tags() {
			langCode := tag.String()
			knownLanguagesValues.add(langCode, supportedLanguageDescription(langCode), display.Self.Name(tag))
		}
	}
	return knownLanguagesValues
}

func supportedLanguageDescription(langCode string) string {
	tag, _ := language.Parse(langCode)
	langName := display.Self.Name(tag)
	if langName != "" {
		return i18n.Localf("%s (%s)", langName, langCode)
	}
	return langCode
}
