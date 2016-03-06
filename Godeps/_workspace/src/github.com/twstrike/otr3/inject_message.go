package otr3

type injections struct {
	messages []ValidMessage
}

// injectMessage will promise to send the messages now or later
// The Injected Messages are promised to be well formed valid messages
// including fragmentation and encoding
func (c *Conversation) injectMessage(vm ValidMessage) {
	c.injections.messages = append(c.injections.messages, vm)
}

func (c *Conversation) withInjects(vms []ValidMessage) []ValidMessage {
	msgs := c.injections.messages
	c.injections.messages = c.injections.messages[0:0]
	return append(vms, msgs...)
}

func (c *Conversation) withInjectionsPlain(plain MessagePlaintext, vms []ValidMessage, err error) (MessagePlaintext, []ValidMessage, error) {
	return plain, c.withInjects(vms), err
}

func (c *Conversation) withInjections(vms []ValidMessage, err error) ([]ValidMessage, error) {
	return c.withInjects(vms), err
}
