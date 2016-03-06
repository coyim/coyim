package otr3

type policies int

type policy int

const (
	allowV2 policy = 2 << iota
	allowV3
	requireEncryption
	sendWhitespaceTag
	whitespaceStartAKE
	errorStartAKE
)

func (p *policies) isOTREnabled() bool {
	return p.has(allowV2) || p.has(allowV3)
}

func (p *policies) has(c policy) bool {
	return int(*p)&int(c) == int(c)
}

func (p *policies) add(c policy) {
	*p = policies(int(*p) | int(c))
}

func (p *policies) AllowV2() {
	p.add(allowV2)
}

func (p *policies) AllowV3() {
	p.add(allowV3)
}

func (p *policies) RequireEncryption() {
	p.add(requireEncryption)
}

func (p *policies) SendWhitespaceTag() {
	p.add(sendWhitespaceTag)
}

func (p *policies) WhitespaceStartAKE() {
	p.add(whitespaceStartAKE)
}

func (p *policies) ErrorStartAKE() {
	p.add(errorStartAKE)
}
