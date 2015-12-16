package cli

import "../client"

func (c *cliUI) ExecuteCmd(comm interface{}) {
	c.commands <- comm
}

func (c *cliUI) watchClientCommands() {
	for command := range c.commands {
		switch comm := command.(type) {
		case client.AuthorizeFingerprintCmd:
			account := comm.Account
			uid := comm.Peer
			fpr := comm.Fingerprint

			//TODO: it could be a different pointer,
			//find the account by ID()
			account.AuthorizeFingerprint(uid, fpr)
			c.ExecuteCmd(client.SaveApplicationConfigCmd{})
		case client.SaveInstanceTagCmd:
			account := comm.Account
			account.InstanceTag = comm.InstanceTag
			c.ExecuteCmd(client.SaveApplicationConfigCmd{})
		case client.SaveApplicationConfigCmd:
			c.SaveConf()
		}
	}
}
