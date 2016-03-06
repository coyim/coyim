package otr3

type smpStateBase struct{}
type smpStateExpect1 struct{ smpStateBase }
type smpStateExpect2 struct{ smpStateBase }
type smpStateExpect3 struct{ smpStateBase }
type smpStateExpect4 struct{ smpStateBase }
type smpStateWaitingForSecret struct {
	smpStateBase
	msg smp1Message
}

type smpMessage interface {
	receivedMessage(*Conversation) (smpMessage, error)
	tlv() tlv
}

type smpState interface {
	startAuthenticate(*Conversation, string, []byte) ([]tlv, error)
	receiveMessage1(*Conversation, smp1Message) (smpState, smpMessage, error)
	continueMessage1(*Conversation, []byte) (smpState, smpMessage, error)
	receiveMessage2(*Conversation, smp2Message) (smpState, smpMessage, error)
	receiveMessage3(*Conversation, smp3Message) (smpState, smpMessage, error)
	receiveMessage4(*Conversation, smp4Message) (smpState, smpMessage, error)
	identity() int
	identityString() string
}

func (c *Conversation) restart() []byte {
	var ret smpMessage
	c.smp.state, ret, _ = sendSMPAbortAndRestartStateMachine()
	return ret.tlv().serialize()
}

func abortState(e error) (smpState, smpMessage, error) {
	return smpStateExpect1{}, smpMessageAbort{}, e
}

func sendSMPAbortAndRestartStateMachine() (smpState, smpMessage, error) {
	//must return nil error otherwise the abort message will be ignored
	return abortState(nil)
}

func (c *Conversation) abortStateMachineAndNotifyCheated() (smpState, smpMessage, error) {
	c.smpEvent(SMPEventCheated, 0)
	return sendSMPAbortAndRestartStateMachine()
}

func (c *Conversation) receiveSMP(m smpMessage) (*tlv, error) {
	toSend, err := m.receivedMessage(c)

	if err != nil {
		return nil, err
	}

	if toSend == nil {
		return nil, nil
	}

	result := toSend.tlv()

	return &result, nil
}

func (c *Conversation) continueSMP(mutualSecret []byte) (*tlv, error) {
	toSend, err := c.continueMessage(mutualSecret)

	if err != nil {
		return nil, err
	}

	result := toSend.tlv()

	return &result, nil
}

func abortStateMachineAndNotifyError(c *Conversation) (smpState, smpMessage, error) {
	c.smpEvent(SMPEventError, 0)
	return sendSMPAbortAndRestartStateMachine()
}

func (smpStateBase) receiveMessage1(c *Conversation, m smp1Message) (smpState, smpMessage, error) {
	return abortStateMachineAndNotifyError(c)
}

func (smpStateBase) continueMessage1(c *Conversation, mutualSecret []byte) (smpState, smpMessage, error) {
	return abortState(errNotWaitingForSMPSecret)
}

func (smpStateBase) receiveMessage2(c *Conversation, m smp2Message) (smpState, smpMessage, error) {
	return abortStateMachineAndNotifyError(c)
}

func (smpStateBase) receiveMessage3(c *Conversation, m smp3Message) (smpState, smpMessage, error) {
	return abortStateMachineAndNotifyError(c)
}

func (smpStateBase) receiveMessage4(c *Conversation, m smp4Message) (smpState, smpMessage, error) {
	return abortStateMachineAndNotifyError(c)
}

func (smpStateExpect1) receiveMessage1(c *Conversation, m smp1Message) (smpState, smpMessage, error) {
	err := c.verifySMP1(m)
	if err != nil {
		return c.abortStateMachineAndNotifyCheated()
	}

	if m.hasQuestion {
		c.smp.question = &m.question
		c.smpEventWithQuestion(SMPEventAskForAnswer, 25, m.question)
	} else {
		c.smpEvent(SMPEventAskForSecret, 25)
	}

	return smpStateWaitingForSecret{msg: m}, nil, nil
}

func (s smpStateWaitingForSecret) continueMessage1(c *Conversation, mutualSecret []byte) (smpState, smpMessage, error) {
	if !c.IsEncrypted() {
		return abortState(errCantAuthenticateWithoutEncryption)
	}

	// Using ssid here should always be safe - we can't be in an encrypted state without having gone through the AKE
	c.smp.secret = generateSMPSecret(c.theirKey.Fingerprint(), c.ourCurrentKey.PublicKey().Fingerprint(), c.ssid[:], mutualSecret, c.version)
	s2, err := c.generateSMP2(c.smp.secret, s.msg)
	if err != nil {
		return c.abortStateMachineAndNotifyCheated()
	}

	c.smp.s2 = &s2

	return smpStateExpect3{}, s2.msg, nil
}

func (smpStateExpect2) receiveMessage2(c *Conversation, m smp2Message) (smpState, smpMessage, error) {
	err := c.verifySMP2(c.smp.s1, m)
	if err != nil {
		return c.abortStateMachineAndNotifyCheated()
	}

	s3, err := c.generateSMP3(c.smp.secret, *c.smp.s1, m)
	if err != nil {
		return c.abortStateMachineAndNotifyCheated()
	}

	c.smpEvent(SMPEventInProgress, 60)

	c.smp.s3 = &s3

	return smpStateExpect4{}, s3.msg, nil
}

func (smpStateExpect3) receiveMessage3(c *Conversation, m smp3Message) (smpState, smpMessage, error) {
	err := c.verifySMP3(c.smp.s2, m)
	if err != nil {
		return c.abortStateMachineAndNotifyCheated()
	}

	err = c.verifySMP3ProtocolSuccess(c.smp.s2, m)
	if err != nil {
		c.smpEvent(SMPEventFailure, 100)
		return sendSMPAbortAndRestartStateMachine()
	}
	c.smpEvent(SMPEventSuccess, 100)

	ret, err := c.generateSMP4(c.smp.secret, *c.smp.s2, m)
	if err != nil {
		return c.abortStateMachineAndNotifyCheated()
	}

	c.smp.wipe()
	return smpStateExpect1{}, ret.msg, nil
}

func (smpStateExpect4) receiveMessage4(c *Conversation, m smp4Message) (smpState, smpMessage, error) {
	err := c.verifySMP4(c.smp.s3, m)
	if err != nil {
		return c.abortStateMachineAndNotifyCheated()
	}

	err = c.verifySMP4ProtocolSuccess(c.smp.s1, c.smp.s3, m)
	if err != nil {
		c.smpEvent(SMPEventFailure, 100)
		return sendSMPAbortAndRestartStateMachine()
	}
	c.smpEvent(SMPEventSuccess, 100)

	c.smp.wipe()
	return smpStateExpect1{}, nil, nil
}

func (m smp1Message) receivedMessage(c *Conversation) (ret smpMessage, err error) {
	c.smp.state, ret, err = c.smp.state.receiveMessage1(c, m)
	return
}

func (m smp2Message) receivedMessage(c *Conversation) (ret smpMessage, err error) {
	c.smp.state, ret, err = c.smp.state.receiveMessage2(c, m)
	return
}

func (m smp3Message) receivedMessage(c *Conversation) (ret smpMessage, err error) {
	c.smp.state, ret, err = c.smp.state.receiveMessage3(c, m)
	return
}

func (m smp4Message) receivedMessage(c *Conversation) (ret smpMessage, err error) {
	c.smp.state, ret, err = c.smp.state.receiveMessage4(c, m)
	return
}

func (m smpMessageAbort) receivedMessage(c *Conversation) (ret smpMessage, err error) {
	c.smp.state = smpStateExpect1{}
	c.smpEvent(SMPEventAbort, 0)
	return
}

func (c *Conversation) continueMessage(mutualSecret []byte) (ret smpMessage, err error) {
	c.smp.state, ret, err = c.smp.state.continueMessage1(c, mutualSecret)
	return
}

func (smpStateExpect1) String() string          { return "SMPSTATE_EXPECT1" }
func (smpStateExpect2) String() string          { return "SMPSTATE_EXPECT2" }
func (smpStateExpect3) String() string          { return "SMPSTATE_EXPECT3" }
func (smpStateExpect4) String() string          { return "SMPSTATE_EXPECT4" }
func (smpStateWaitingForSecret) String() string { return "SMPSTATE_WAITINGFORSECRET (internal)" }

func (smpStateBase) startAuthenticate(c *Conversation, question string, mutualSecret []byte) (tlvs []tlv, err error) {
	tlvs, err = smpStateExpect1{}.startAuthenticate(c, question, mutualSecret)
	tlvs = append([]tlv{smpMessageAbort{}.tlv()}, tlvs...)
	return
}

func (smpStateExpect1) startAuthenticate(c *Conversation, question string, mutualSecret []byte) (tlvs []tlv, err error) {
	if !c.IsEncrypted() {
		return nil, errCantAuthenticateWithoutEncryption
	}

	// Using ssid here should always be safe - we can't be in an encrypted state without having gone through the AKE
	c.smp.secret = generateSMPSecret(c.ourCurrentKey.PublicKey().Fingerprint(), c.theirKey.Fingerprint(), c.ssid[:], mutualSecret, c.version)

	s1, err := c.generateSMP1()
	if err != nil {
		return nil, errShortRandomRead
	}

	if question != "" {
		s1.msg.hasQuestion = true
		s1.msg.question = question
	}

	c.smp.s1 = &s1
	c.smp.state = smpStateExpect2{}

	return []tlv{s1.msg.tlv()}, nil
}
