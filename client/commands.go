package client

import "../config"

type AuthorizeFingerprintCmd struct {
	Account     *config.Account
	Peer        string
	Fingerprint []byte
}

type SaveInstanceTagCmd struct {
	Account     *config.Account
	InstanceTag uint32
}

type SaveApplicationConfigCmd struct{}

type CommandManager interface {
	ExecuteCmd(c interface{})
}
