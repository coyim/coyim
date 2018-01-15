package cli

import "github.com/coyim/coyim/otr_client"

func (c *cliUI) ExecuteCmd(comm interface{}) {
	c.commands <- comm
}

func (c *cliUI) watchClientCommands() {
	for command := range c.commands {
		switch comm := command.(type) {
		case otr_client.AuthorizeFingerprintCmd:
			account := comm.Account
			uid := comm.Peer
			fpr := comm.Fingerprint

			//TODO: it could be a different pointer,
			//find the account by ID()
			account.AuthorizeFingerprint(uid.Representation(), fpr)
			c.ExecuteCmd(otr_client.SaveApplicationConfigCmd{})
		case otr_client.SaveInstanceTagCmd:
			account := comm.Account
			account.InstanceTag = comm.InstanceTag
			c.ExecuteCmd(otr_client.SaveApplicationConfigCmd{})
		case otr_client.SaveApplicationConfigCmd:
			c.SaveConf()
		}
	}
}
