package otr3

import "math/big"

type smp3State struct {
	x              *big.Int
	g3b            *big.Int
	r4, r5, r6, r7 *big.Int
	qaqb, papb     *big.Int
	msg            smp3Message
}

type smp3Message struct {
	pa, qa     *big.Int
	cp         *big.Int
	d5, d6, d7 *big.Int
	ra         *big.Int
	cr         *big.Int
}

func (m smp3Message) tlv() tlv {
	return genSMPTLV(uint16(tlvTypeSMP3), m.pa, m.qa, m.cp, m.d5, m.d6, m.ra, m.cr, m.d7)
}

func (c *Conversation) generateSMP3Parameters() (s smp3State, err error) {
	b := make([]byte, c.version.parameterLength())
	var err1, err2, err3, err4 error

	s.r4, err1 = c.randMPI(b)
	s.r5, err2 = c.randMPI(b)
	s.r6, err3 = c.randMPI(b)
	s.r7, err4 = c.randMPI(b)

	return s, firstError(err1, err2, err3, err4)
}

func generateSMP3Message(s *smp3State, s1 smp1State, m2 smp2Message, v otrVersion) smp3Message {
	var m smp3Message

	g2 := modExp(m2.g2b, s1.a2)
	g3 := modExp(m2.g3b, s1.a3)

	m.pa = modExp(g3, s.r4)
	m.qa = mulMod(modExp(g1, s.r4), modExp(g2, s.x), p)

	s.g3b = m2.g3b
	s.qaqb = divMod(m.qa, m2.qb, p)
	s.papb = divMod(m.pa, m2.pb, p)

	m.cp = hashMPIsBN(v.hash2Instance(), 6, modExp(g3, s.r5), mulMod(modExp(g1, s.r5), modExp(g2, s.r6), p))
	m.d5 = generateDZKP(s.r5, s.r4, m.cp)
	m.d6 = generateDZKP(s.r6, s.x, m.cp)

	m.ra = modExp(s.qaqb, s1.a3)

	m.cr = hashMPIsBN(v.hash2Instance(), 7, modExp(g1, s.r7), modExp(s.qaqb, s.r7))
	m.d7 = subMod(s.r7, mul(s1.a3, m.cr), q)

	return m
}

func (c *Conversation) generateSMP3(secret *big.Int, s1 smp1State, m2 smp2Message) (s smp3State, err error) {
	if s, err = c.generateSMP3Parameters(); err != nil {
		return s, err
	}
	s.x = secret
	s.msg = generateSMP3Message(&s, s1, m2, c.version)
	return
}

func (c *Conversation) verifySMP3(s2 *smp2State, msg smp3Message) error {
	if !c.version.isGroupElement(msg.pa) {
		return newOtrError("Pa is an invalid group element")
	}

	if !c.version.isGroupElement(msg.qa) {
		return newOtrError("Qa is an invalid group element")
	}

	if !c.version.isGroupElement(msg.ra) {
		return newOtrError("Ra is an invalid group element")
	}

	if !verifyZKP3(msg.cp, s2.g2, s2.g3, msg.d5, msg.d6, msg.pa, msg.qa, 6, c.version) {
		return newOtrError("cP is not a valid zero knowledge proof")
	}

	qaqb := divMod(msg.qa, s2.qb, p)

	if !verifyZKP4(msg.cr, s2.g3a, msg.d7, qaqb, msg.ra, 7, c.version) {
		return newOtrError("cR is not a valid zero knowledge proof")
	}

	return nil
}

func (c *Conversation) verifySMP3ProtocolSuccess(s2 *smp2State, msg smp3Message) error {
	papb := divMod(msg.pa, s2.pb, p)

	rab := modExp(msg.ra, s2.b3)
	if !eq(rab, papb) {
		return newOtrError("protocol failed: x != y")
	}

	return nil
}
