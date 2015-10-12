package client

type Client interface {
	LoadConfig(string) error
	Loop()
	Close()
}
