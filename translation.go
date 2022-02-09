package main

import (
	"os"
	"path/filepath"

	"github.com/coyim/coyim/config"
	"github.com/coyim/coyim/i18n"
)

const localDirectoryForTranslations = "./i18n"

func dirExists(name string) bool {
	st, err := os.Stat(name)
	return err == nil && st.IsDir()
}

func shouldUseEmbeddedTranslations() bool {
	return !dirExists(localDirectoryForTranslations)
}

func translationsDirectory() string {
	if shouldUseEmbeddedTranslations() {
		return filepath.Join(config.SystemDataDir(), "coyim", "translations")
	}
	return localDirectoryForTranslations
}

func initTranslations() {
	if shouldUseEmbeddedTranslations() {
		i18n.UnpackTranslationFilesInto(translationsDirectory(), BuildCommit)
	}
}
