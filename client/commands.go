package client

import "github.com/twstrike/coyim/config"

// AuthorizeFingerprintCmd is a command that represents a request to authorize a fingerprint
type AuthorizeFingerprintCmd struct {
	Account     *config.Account
	Session     interface{}
	Peer        string
	Fingerprint []byte
}

// SaveInstanceTagCmd is a command that represents a request to save an instance tag
type SaveInstanceTagCmd struct {
	Account     *config.Account
	InstanceTag uint32
}

// SaveApplicationConfigCmd is a command that represents a request to save the application configuration
type SaveApplicationConfigCmd struct{}

// CommandManager is anything that can execute commands
type CommandManager interface {
	ExecuteCmd(c interface{})
}
