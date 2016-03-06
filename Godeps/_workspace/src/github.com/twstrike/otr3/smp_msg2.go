package otr3

import "math/big"

type smp2State struct {
	y                  *big.Int
	b2, b3             *big.Int
	r2, r3, r4, r5, r6 *big.Int
	g3a                *big.Int
	g2, g3             *big.Int
	pb, qb             *big.Int
	msg                smp2Message
}

type smp2Message struct {
	g2b, g3b *big.Int
	c2, c3   *big.Int
	d2, d3   *big.Int
	pb, qb   *big.Int
	cp       *big.Int
	d5, d6   *big.Int
}

func (m smp2Message) tlv() tlv {
	return genSMPTLV(uint16(tlvTypeSMP2), m.g2b, m.c2, m.d2, m.g3b, m.c3, m.d3, m.pb, m.qb, m.cp, m.d5, m.d6)
}

func (c *Conversation) generateSMP2Parameters() (s smp2State, err error) {
	b := make([]byte, c.version.parameterLength())
	var err1, err2, err3, err4, err5, err6, err7 error
	s.b2, err1 = c.randMPI(b)
	s.b3, err2 = c.randMPI(b)
	s.r2, err3 = c.randMPI(b)
	s.r3, err4 = c.randMPI(b)
	s.r4, err5 = c.randMPI(b)
	s.r5, err6 = c.randMPI(b)
	s.r6, err7 = c.randMPI(b)

	return s, firstError(err1, err2, err3, err4, err5, err6, err7)
}

func generateSMP2Message(s *smp2State, s1 smp1Message, v otrVersion) smp2Message {
	var m smp2Message

	m.g2b = modExp(g1, s.b2)
	m.g3b = modExp(g1, s.b3)

	m.c2, m.d2 = generateZKP(s.r2, s.b2, 3, v)
	m.c3, m.d3 = generateZKP(s.r3, s.b3, 4, v)

	s.g3a = s1.g3a
	s.g2 = modExp(s1.g2a, s.b2)
	s.g3 = modExp(s1.g3a, s.b3)

	s.pb = modExp(s.g3, s.r4)
	s.qb = mulMod(modExp(g1, s.r4), modExp(s.g2, s.y), p)

	m.pb = s.pb
	m.qb = s.qb

	m.cp = hashMPIsBN(v.hash2Instance(), 5,
		modExp(s.g3, s.r5),
		mulMod(modExp(g1, s.r5), modExp(s.g2, s.r6), p))

	m.d5 = subMod(s.r5, mul(s.r4, m.cp), q)
	m.d6 = subMod(s.r6, mul(s.y, m.cp), q)

	return m
}

func (c *Conversation) generateSMP2(secret *big.Int, s1 smp1Message) (s smp2State, err error) {
	if s, err = c.generateSMP2Parameters(); err != nil {
		return s, err
	}

	s.y = secret
	s.msg = generateSMP2Message(&s, s1, c.version)
	return
}

func (c *Conversation) verifySMP2(s1 *smp1State, msg smp2Message) error {
	if !c.version.isGroupElement(msg.g2b) {
		return newOtrError("g2b is an invalid group element")
	}

	if !c.version.isGroupElement(msg.g3b) {
		return newOtrError("g3b is an invalid group element")
	}

	if !c.version.isGroupElement(msg.pb) {
		return newOtrError("Pb is an invalid group element")
	}

	if !c.version.isGroupElement(msg.qb) {
		return newOtrError("Qb is an invalid group element")
	}

	if !verifyZKP(msg.d2, msg.g2b, msg.c2, 3, c.version) {
		return newOtrError("c2 is not a valid zero knowledge proof")
	}

	if !verifyZKP(msg.d3, msg.g3b, msg.c3, 4, c.version) {
		return newOtrError("c3 is not a valid zero knowledge proof")
	}

	g2 := modExp(msg.g2b, s1.a2)
	g3 := modExp(msg.g3b, s1.a3)

	if !verifyZKP2(g2, g3, msg.d5, msg.d6, msg.pb, msg.qb, msg.cp, 5, c.version) {
		return newOtrError("cP is not a valid zero knowledge proof")
	}

	return nil
}
