package otr3

import "math/big"

type smp1State struct {
	a2, a3 *big.Int
	r2, r3 *big.Int
	msg    smp1Message
}

type smp1Message struct {
	g2a, g3a    *big.Int
	c2, c3      *big.Int
	d2, d3      *big.Int
	hasQuestion bool
	question    string
}

func (m smp1Message) tlv() tlv {
	t := genSMPTLV(uint16(tlvTypeSMP1), m.g2a, m.c2, m.d2, m.g3a, m.c3, m.d3)
	if m.hasQuestion {
		t.tlvType = tlvTypeSMP1WithQuestion
		t.tlvValue = append(append([]byte(m.question), 0), t.tlvValue...)
		t.tlvLength = uint16(len(t.tlvValue))
	}
	return t
}

func (c *Conversation) generateSMP1Parameters() (s smp1State, err error) {
	b := make([]byte, c.version.parameterLength())
	var err1, err2, err3, err4 error
	s.a2, err1 = c.randMPI(b)
	s.a3, err2 = c.randMPI(b)
	s.r2, err3 = c.randMPI(b)
	s.r3, err4 = c.randMPI(b)
	return s, firstError(err1, err2, err3, err4)
}

func generateSMP1Message(s smp1State, v otrVersion) (m smp1Message) {
	m.g2a = modExpP(g1, s.a2)
	m.g3a = modExpP(g1, s.a3)
	m.c2, m.d2 = generateZKP(s.r2, s.a2, 1, v)
	m.c3, m.d3 = generateZKP(s.r3, s.a3, 2, v)
	return
}

func (c *Conversation) generateSMP1() (s smp1State, err error) {
	if s, err = c.generateSMP1Parameters(); err != nil {
		return s, err
	}
	s.msg = generateSMP1Message(s, c.version)
	return
}

func (c *Conversation) verifySMP1(msg smp1Message) error {
	if !c.version.isGroupElement(msg.g2a) {
		return newOtrError("g2a is an invalid group element")
	}

	if !c.version.isGroupElement(msg.g3a) {
		return newOtrError("g3a is an invalid group element")
	}

	if !verifyZKP(msg.d2, msg.g2a, msg.c2, 1, c.version) {
		return newOtrError("c2 is not a valid zero knowledge proof")
	}

	if !verifyZKP(msg.d3, msg.g3a, msg.c3, 2, c.version) {
		return newOtrError("c3 is not a valid zero knowledge proof")
	}

	return nil
}
