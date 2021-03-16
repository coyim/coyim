package gui

import (
	"os"
	"sort"
	"strings"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
)

const (
	languageSelectorCodeIndex int = iota
	languageSelectorDescriptionIndex
)

type languageSelectorComponent struct {
	entry gtki.Entry
	combo gtki.ComboBoxText
	model gtki.ListStore

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

	for langCode, e := range lc.languages.sort() {
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
	code        string
	description string
	values      []string
}

func newlanguageSelectorEntry(code, description string) *languageSelectorEntry {
	return &languageSelectorEntry{
		code:        code,
		description: description,
		values:      []string{description},
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
	collator *collate.Collator
	list     []*languageSelectorEntry
}

func newLanguageSelectorValues(lang language.Tag) *languageSelectorValues {
	return &languageSelectorValues{
		collator: collate.New(lang, collate.OptionsFromTag(lang)),
	}
}

func (v *languageSelectorValues) languageBasedOnText(t string) string {
	if ix := v.indexOf(t); ix != -1 {
		return v.list[ix].code
	}

	if vix := v.valueIndexOf(t); vix != -1 {
		return v.list[vix].code
	}

	return t
}

func (v *languageSelectorValues) valueIndexOf(langDesc string) int {
	for ix, e := range v.list {
		if e.contains(langDesc) {
			return ix
		}
	}
	return -1
}

func (v *languageSelectorValues) indexOf(t string) int {
	for ix, e := range v.list {
		if e.code == t {
			return ix
		}
	}
	return -1
}

func (v *languageSelectorValues) add(langCode string, langDesc string, values ...string) {
	var entry *languageSelectorEntry

	ix := v.indexOf(langCode)
	if ix == -1 {
		entry = newlanguageSelectorEntry(langCode, langDesc)
		v.list = append(v.list, entry)
	} else {
		entry = v.list[ix]
	}

	entry.add(values...)
}

func (v *languageSelectorValues) sort() []*languageSelectorEntry {
	copy := make([]*languageSelectorEntry, len(v.list))
	for ix, e := range v.list {
		copy[ix] = e
	}

	sort.SliceStable(copy, func(i, j int) bool {
		return v.collator.CompareString(copy[i].description, copy[j].description) == -1
	})

	return copy
}

var knownLanguagesValues *languageSelectorValues

func getKnownLanguages() *languageSelectorValues {
	if knownLanguagesValues == nil {
		sl := systemDefaultLanguage()
		knownLanguagesValues = newLanguageSelectorValues(sl)

		for _, tag := range display.Supported.Tags() {
			langCode := tag.String()
			knownLanguagesValues.add(langCode, supportedLanguageDescription(langCode), display.Self.Name(tag))
		}
	}
	return knownLanguagesValues
}

func supportedLanguageDescription(langCode string) string {
	systemLangNamer := systemLanguageNamer()
	langTag := language.Make(langCode)

	friendlyName := systemLangNamer.Name(langTag)
	langName := display.Self.Name(langTag)

	if langName != "" && friendlyName != langName {
		return i18n.Localf("%s (%s)", friendlyName, langName)
	}

	return friendlyName
}

func systemLanguageNamer() display.Namer {
	namer := display.Tags(systemDefaultLanguage())
	if namer != nil {
		return namer
	}

	return display.Self
}

func systemDefaultLanguage() language.Tag {
	lang, isPresent := os.LookupEnv("LC_ALL")
	if isPresent {
		tag, err := language.Parse(lang)
		if err == nil {
			return tag
		}
	}
	return language.English
}
