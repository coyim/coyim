package otr3

import (
	"math/big"
)

type smp struct {
	state    smpState
	question *string
	secret   *big.Int
	s1       *smp1State
	s2       *smp2State
	s3       *smp3State
}

const smpVersion = 1

func (s *smp) wipe() {
	s.state = nil
	s.question = nil
	wipeBigInt(s.secret)
	s.secret = nil
	s.s1 = nil
	s.s2 = nil
	s.s3 = nil
}

func (s *smp) ensureSMP() {
	if s.state != nil {
		return
	}

	s.state = smpStateExpect1{}
}

// SMPQuestion returns the current SMP question and ok if there is one, and not ok if there isn't one.
func (c *Conversation) SMPQuestion() (string, bool) {
	if c.smp.question == nil {
		return "", false
	}
	return *c.smp.question, true
}

func generateSMPSecret(initiatorFingerprint, recipientFingerprint, ssid, secret []byte, v otrVersion) *big.Int {
	h := v.hash2Instance()
	_, _ = h.Write([]byte{smpVersion})
	_, _ = h.Write(initiatorFingerprint)
	_, _ = h.Write(recipientFingerprint)
	_, _ = h.Write(ssid)
	_, _ = h.Write(secret)
	return new(big.Int).SetBytes(h.Sum(nil))
}

func generateDZKP(r, a, c *big.Int) *big.Int {
	return subMod(r, mul(a, c), q)
}

func generateZKP(r, a *big.Int, ix byte, v otrVersion) (c, d *big.Int) {
	c = hashMPIsBN(v.hash2Instance(), ix, modExpP(g1, r))
	d = generateDZKP(r, a, c)
	return
}

func verifyZKP(d, gen, c *big.Int, ix byte, v otrVersion) bool {
	r := modExpP(g1, d)
	s := modExpP(gen, c)
	t := hashMPIsBN(v.hash2Instance(), ix, mulMod(r, s, p))
	return eq(c, t)
}

func verifyZKP2(g2, g3, d5, d6, pb, qb, cp *big.Int, ix byte, v otrVersion) bool {
	l := mulMod(
		modExpP(g3, d5),
		modExpP(pb, cp),
		p)
	r := mulMod(mul(modExpP(g1, d5),
		modExpP(g2, d6)),
		modExpP(qb, cp),
		p)
	t := hashMPIsBN(v.hash2Instance(), ix, l, r)
	return eq(cp, t)
}

func verifyZKP3(cp, g2, g3, d5, d6, pa, qa *big.Int, ix byte, v otrVersion) bool {
	l := mulMod(modExpP(g3, d5), modExpP(pa, cp), p)
	r := mulMod(mul(modExpP(g1, d5), modExpP(g2, d6)), modExpP(qa, cp), p)
	t := hashMPIsBN(v.hash2Instance(), ix, l, r)
	return eq(cp, t)
}

func verifyZKP4(cr, g3a, d7, qaqb, ra *big.Int, ix byte, v otrVersion) bool {
	l := mulMod(modExpP(g1, d7), modExpP(g3a, cr), p)
	r := mulMod(modExpP(qaqb, d7), modExpP(ra, cr), p)
	t := hashMPIsBN(v.hash2Instance(), ix, l, r)
	return eq(cr, t)
}

func genSMPTLV(tp uint16, mpis ...*big.Int) tlv {
	data := make([]byte, 0, 1000)

	// TODO: is this really correct? It seems like it's adding the length for MPIs two times
	data = AppendWord(data, uint32(len(mpis)))
	data = AppendMPIs(data, mpis...)
	length := uint16(len(data))
	out := tlv{
		tlvType:   tp,
		tlvLength: length,
		tlvValue:  data,
	}

	return out
}
