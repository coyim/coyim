package glib_mock

type MockValue struct{}

func (*MockValue) GetString() (string, error) {
	return "", nil
}

func (*MockValue) GoValue() (interface{}, error) {
	return nil, nil
}
