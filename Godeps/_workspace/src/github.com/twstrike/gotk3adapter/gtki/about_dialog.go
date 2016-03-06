package gtki

type AboutDialog interface {
	Dialog

	SetAuthors([]string)
	SetProgramName(string)
	SetVersion(string)
	SetLicense(string)
	SetWrapLicense(bool)
}

func AssertAboutDialog(_ AboutDialog) {}
