package gtk_mock

type MockRevealer struct {
	MockBin
}

func (*MockRevealer) SetRevealChild(revealChild bool) {
}

func (*MockRevealer) GetRevealChild() bool {
	return false
}
