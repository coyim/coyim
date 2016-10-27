package otr3

import "math/big"

func (p *dhKeyPair) wipe() {
	if p == nil {
		return
	}

	wipeBigInt(p.pub)
	wipeBigInt(p.priv)
	p.pub = nil
	p.priv = nil
}

func (k *akeKeys) wipe() {
	if k == nil {
		return
	}
	wipeBytes(k.c)
	k.c = nil
	wipeBytes(k.m1)
	k.m1 = nil
	wipeBytes(k.m2)
	k.m2 = nil
}

func (a *ake) wipe(wipeKeys bool) {
	if a == nil {
		return
	}

	wipeBigInt(a.secretExponent)
	a.secretExponent = nil

	wipeBigInt(a.ourPublicValue)
	a.ourPublicValue = nil

	wipeBigInt(a.theirPublicValue)
	a.theirPublicValue = nil

	wipeBytes(a.r[:])

	a.wipeGX()
	a.revealKey.wipe()
	a.sigKey.wipe()

	if wipeKeys {
		a.keys.wipe()
	} else {
		a.keys = keyManagementContext{}
	}
}

func (a *ake) wipeGX() {
	if a == nil {
		return
	}

	wipeBytes(a.xhashedGx)
	wipeBytes(a.encryptedGx[:])
	a.xhashedGx = nil
	a.encryptedGx = nil
}

func (c *keyManagementContext) wipeKeys() {
	if c == nil {
		return
	}

	c.ourCurrentDHKeys.wipe()
	c.ourPreviousDHKeys.wipe()

	wipeBigInt(c.theirCurrentDHPubKey)
	c.theirCurrentDHPubKey = nil

	wipeBigInt(c.theirPreviousDHPubKey)
	c.theirPreviousDHPubKey = nil
}

func (c *keyManagementContext) wipe() {
	if c == nil {
		return
	}

	c.wipeKeys()
	c.ourKeyID = 0
	c.theirKeyID = 0

	for i := range c.oldMACKeys {
		c.oldMACKeys[i].wipe()
		c.oldMACKeys[i] = []byte{}
	}
	c.oldMACKeys = nil

	c.counterHistory.wipe()
	c.macKeyHistory.wipe()
}

func (c *keyManagementContext) wipeAndKeepRevealKeys() keyManagementContext {
	ret := keyManagementContext{}
	ret.oldMACKeys = make([]macKey, len(c.oldMACKeys))
	copy(ret.oldMACKeys, c.oldMACKeys)

	c.wipe()

	return ret
}

func (h *counterHistory) wipe() {
	if h == nil {
		return
	}

	for i := range h.counters {
		h.counters[i].wipe()
		h.counters[i] = nil
	}

	h.counters = nil
}

func (c *keyPairCounter) wipe() {
	if c == nil {
		return
	}

	c.ourKeyID = 0
	c.theirKeyID = 0
	c.ourCounter = 0
	c.theirCounter = 0
}

func (h *macKeyHistory) wipe() {
	if h == nil {
		return
	}

	for i := range h.items {
		h.items[i].wipe()
		h.items[i] = macKeyUsage{} //prevent memory leak
	}
	h.items = nil
}

func (u *macKeyUsage) wipe() {
	if u == nil {
		return
	}

	u.receivingKey.wipe()
	u.receivingKey = macKey{}
}

func (k *macKey) wipe() {
	if k == nil {
		return
	}

	wipeBytes(*k)
	*k = []byte{}
}

func zeroes(n int) []byte {
	return make([]byte, n)
}

func wipeBytes(b []byte) {
	copy(b, zeroes(len(b)))
}

func wipeBigInt(k *big.Int) {
	if k == nil {
		return
	}

	k.SetBytes(zeroes(len(k.Bytes())))
}

func setBigInt(dst *big.Int, src *big.Int) *big.Int {
	wipeBigInt(dst)

	ret := big.NewInt(0)
	ret.Set(src)
	return ret
}
