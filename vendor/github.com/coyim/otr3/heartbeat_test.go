package otr3

import (
	"crypto/rand"
	"testing"
	"time"
)

func Test_potentialHeartbeat_returnsNothingIfThereWasntPlaintext(t *testing.T) {
	c := newConversation(otrV3{}, rand.Reader)
	c.msgState = encrypted
	var plain []byte
	ret, err := c.potentialHeartbeat(plain)
	assertNil(t, ret)
	assertNil(t, err)
}

func Test_potentialHeartbeat_returnsNothingIfLastSentWasRecently(t *testing.T) {
	c := newConversation(otrV3{}, rand.Reader)
	c.msgState = encrypted
	c.heartbeat.lastSent = time.Now().Add(-10 * time.Second)
	plain := []byte("Foo plain")
	ret, err := c.potentialHeartbeat(plain)
	assertNil(t, ret)
	assertNil(t, err)
}

func Test_potentialHeartbeat_doesntUpdateLastSentIfLastSentWasRecently(t *testing.T) {
	c := newConversation(otrV3{}, rand.Reader)
	c.msgState = encrypted
	tt := time.Now().Add(-10 * time.Second)
	c.heartbeat.lastSent = tt
	plain := []byte("Foo plain")
	c.potentialHeartbeat(plain)
	assertEquals(t, c.heartbeat.lastSent, tt)
}

func Test_potentialHeartbeat_updatesLastSentIfWeNeedToSendAHeartbeat(t *testing.T) {
	c := bobContextAfterAKE()
	c.msgState = encrypted
	tt := time.Now().Add(-61 * time.Second)
	c.heartbeat.lastSent = tt
	plain := []byte("Foo plain")
	c.potentialHeartbeat(plain)
	assertEquals(t, c.heartbeat.lastSent.After(tt), true)
}

func Test_potentialHeartbeat_logsTheHeartbeatWhenWeSendIt(t *testing.T) {
	c := bobContextAfterAKE()
	c.msgState = encrypted
	tt := time.Now().Add(-61 * time.Second)
	c.heartbeat.lastSent = tt
	plain := []byte("Foo plain")

	c.expectMessageEvent(t, func() {
		c.potentialHeartbeat(plain)
	}, MessageEventLogHeartbeatSent, nil, nil)
}

func Test_potentialHeartbeat_putsTogetherAMessageForAHeartbeat(t *testing.T) {
	c := bobContextAfterAKE()
	c.msgState = encrypted
	tt := time.Now().Add(-61 * time.Second)
	c.heartbeat.lastSent = tt
	plain := []byte("Foo plain")

	msg, err := c.potentialHeartbeat(plain)

	assertNil(t, err)
	assertDeepEquals(t, msg, messageWithHeader(bytesFromHex("0003030000010100000101010000000100000001000000c0075dfab5a1eab059052d0ad881c4938d52669630d61833a367155d67d03a457f619683d0fa829781e974fd24f6865e8128a9312a167b77326a87dea032fc31784d05b18b9cbafebe162ae9b5369f8b0c5911cf1be757f45f2a674be5126a714a6366c28086b3c7088911dcc4e5fb1481ad70a5237b8e4a6aff4954c2ca6df338b9f08691e4c0defe12689b37d4df30ddef2687f789fcf623c5d0cf6f09b7e5e69f481d5fd1b24a77636fb676e6d733d129eb93e81189340233044766a36eb07d000000000000000100000100dec63c53470b5a82b4f0176eb0a3aed6da242675dc16c3c473f2b14c2dda1cd52bfd20559eedf51d275b049fdefd93af2325d28d5d2f0fb05e8524842e32d4275c69a621e5fa133977563345055fded5511a78337a6d9a213bc5319de11a578818c2edb21b510595157feea3ed93a1178021571aa21765fd974c89cdcbda8ec0afce0c0ea5901021657b959f842df47224edd5dd50d9e736ed8982580373dcd0e2f06a5421472ae2bc58cc4ea7cb2b054e22c1781b72595909b37640e28f435df98b16410c76969fa9112a114b4ab7fb5b3265aa5efa0a99b9c47097d6d42a232a223d03b7d4a8fd5e57a748d1e06ef106e265f70421b708ca85b89e92f020823cb571fd1efe80922982fdf1fb23f62060f49b1a00000000")))
}

func Test_potentialHeartbeat_returnsAnErrorIfWeCantPutTogetherAMessage(t *testing.T) {
	c := bobContextAfterAKE()
	c.msgState = encrypted
	c.keys.ourKeyID = 0
	tt := time.Now().Add(-61 * time.Second)
	c.heartbeat.lastSent = tt
	plain := []byte("Foo plain")

	_, err := c.potentialHeartbeat(plain)
	assertDeepEquals(t, err, newOtrConflictError("invalid key id for local peer"))
}
