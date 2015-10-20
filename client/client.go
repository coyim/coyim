package client

// Client represent the minimum necessary functionality for a client
type Client interface {
	LoadConfig(string) error
	Loop()
	Close()
}
