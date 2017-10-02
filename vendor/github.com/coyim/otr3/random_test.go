package otr3

import (
	"crypto/rand"
	"testing"
)

func Test_conversation_rand_returnsTheSetRandomIfThereIsOne(t *testing.T) {
	r := fixtureRand()
	c := &Conversation{Rand: r}
	assertEquals(t, c.rand(), r)
}

func Test_conversation_rand_returnsRandReaderIfNoRandomnessIsSet(t *testing.T) {
	c := &Conversation{}
	assertEquals(t, c.rand(), rand.Reader)
}

func Test_randMPI_returnsNilForARealRead(t *testing.T) {
	c := newConversation(otrV3{}, fixedRand([]string{"ABCD"}))
	var buf [2]byte

	_, err := c.randMPI(buf[:])
	assertEquals(t, err, nil)
}

func Test_randMPI_returnsShortRandomReadErrorIfFails(t *testing.T) {
	c := newConversation(otrV3{}, fixedRand([]string{"ABCD"}))
	var buf [3]byte
	_, err := c.randMPI(buf[:])

	assertEquals(t, err, errShortRandomRead)
}
