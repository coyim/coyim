package otr3

type smpMessageAbort struct{}

func (m smpMessageAbort) tlv() tlv {
	return genSMPTLV(tlvTypeSMPAbort)
}
