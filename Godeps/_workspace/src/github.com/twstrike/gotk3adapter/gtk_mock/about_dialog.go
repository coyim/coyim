package gtk_mock

type MockAboutDialog struct {
	MockDialog
}

func (*MockAboutDialog) SetAuthors(v1 []string) {
}

func (*MockAboutDialog) SetProgramName(v1 string) {
}

func (*MockAboutDialog) SetVersion(v1 string) {
}

func (*MockAboutDialog) SetLicense(v1 string) {
}

func (*MockAboutDialog) SetWrapLicense(v1 bool) {
}
