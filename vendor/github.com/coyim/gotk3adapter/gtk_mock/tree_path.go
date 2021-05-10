package gtk_mock

type MockTreePath struct {
}

func (*MockTreePath) GetDepth() int {
	return 0
}

func (*MockTreePath) String() string {
	return ""
}
