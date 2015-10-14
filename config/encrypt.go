package config

func (c *Config) ShouldEncryptTo(uid string) bool {
	if c.AlwaysEncrypt {
		return true
	}

	for _, contact := range c.AlwaysEncryptWith {
		if contact == uid {
			return true
		}
	}
	return false
}
