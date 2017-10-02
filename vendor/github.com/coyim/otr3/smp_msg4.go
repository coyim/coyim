package otr3

import "math/big"

type smp4State struct {
	y   *big.Int
	r7  *big.Int
	msg smp4Message
}

type smp4Message struct {
	cr *big.Int
	d7 *big.Int
	rb *big.Int
}

func (m smp4Message) tlv() tlv {
	return genSMPTLV(uint16(tlvTypeSMP4), m.rb, m.cr, m.d7)
}

func (c *Conversation) generateSMP4(secret *big.Int, s2 smp2State, msg3 smp3Message) (s smp4State, err error) {
	if s, err = c.generateSMP4Parameters(); err != nil {
		return s, err
	}
	s.y = secret
	s.msg = generateSMP4Message(s, s2, msg3, c.version)
	return
}

func (c *Conversation) verifySMP4(s3 *smp3State, msg smp4Message) error {
	if !c.version.isGroupElement(msg.rb) {
		return newOtrError("Rb is an invalid group element")
	}

	if !verifyZKP4(msg.cr, s3.g3b, msg.d7, s3.qaqb, msg.rb, 8, c.version) {
		return newOtrError("cR is not a valid zero knowledge proof")
	}

	return nil
}

func (c *Conversation) generateSMP4Parameters() (s smp4State, err error) {
	b := make([]byte, c.version.parameterLength())
	s.r7, err = c.randMPI(b)
	return
}

func generateSMP4Message(s smp4State, s2 smp2State, msg3 smp3Message, v otrVersion) smp4Message {
	var m smp4Message

	qaqb := divMod(msg3.qa, s2.qb, p)

	m.rb = modExp(qaqb, s2.b3)
	m.cr = hashMPIsBN(v.hash2Instance(), 8, modExp(g1, s.r7), modExp(qaqb, s.r7))
	m.d7 = subMod(s.r7, mul(s2.b3, m.cr), q)

	return m
}

func (c *Conversation) verifySMP4ProtocolSuccess(s1 *smp1State, s3 *smp3State, msg smp4Message) error {
	rab := modExp(msg.rb, s1.a3)
	if !eq(rab, s3.papb) {
		return newOtrError("protocol failed: x != y")
	}

	return nil
}
