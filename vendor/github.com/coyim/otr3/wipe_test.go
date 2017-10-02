package otr3

import (
	"math/big"
	"testing"
)

func Test_zeroes_generateZeroes(t *testing.T) {
	z := zeroes(5)
	assertDeepEquals(t, z, []byte{0, 0, 0, 0, 0})
}

func Test_wipeBytes_zeroesTheSlice(t *testing.T) {
	b := []byte{1, 2, 3, 4, 5}
	wipeBytes(b)

	assertDeepEquals(t, b, zeroes(len(b)))
}

func Test_wipeBigInt_numberIsZeroed(t *testing.T) {
	n := big.NewInt(3)
	wipeBigInt(n)
	assertEquals(t, n.Cmp(big.NewInt(0)), 0)
}

func Test_setBigInt_numberIsSet(t *testing.T) {
	n := big.NewInt(3)
	n = setBigInt(n, big.NewInt(5))
	assertEquals(t, n.Cmp(big.NewInt(5)), 0)
}

func Test_setBigInt_setWhenSourceIsNull(t *testing.T) {
	var n *big.Int
	n = setBigInt(n, big.NewInt(5))
	assertEquals(t, n.Cmp(big.NewInt(5)), 0)
}

func Test_wipe_macKey(t *testing.T) {
	k := macKey{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8, 7, 6, 5}
	k.wipe()

	assertDeepEquals(t, k, macKey{})
}

func Test_wipe_keyManagementContext(t *testing.T) {
	keys := keyManagementContext{
		ourKeyID:   2,
		theirKeyID: 3,
		ourCurrentDHKeys: dhKeyPair{
			priv: big.NewInt(1),
			pub:  big.NewInt(2),
		},
		ourPreviousDHKeys: dhKeyPair{
			priv: big.NewInt(3),
			pub:  big.NewInt(4),
		},
		theirCurrentDHPubKey:  big.NewInt(5),
		theirPreviousDHPubKey: big.NewInt(6),
		counterHistory: counterHistory{
			counters: []*keyPairCounter{
				&keyPairCounter{1, 1, 1, 1},
			},
		},
		macKeyHistory: macKeyHistory{
			items: []macKeyUsage{
				macKeyUsage{
					ourKeyID:     2,
					theirKeyID:   3,
					receivingKey: macKey{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 3, 4},
				},
			},
		},
		oldMACKeys: []macKey{
			macKey{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8, 7, 6, 5},
		},
	}

	keys.wipe()

	assertDeepEquals(t, keys, keyManagementContext{})
}

func Test_dhKeyPair_wipe_HandlesNilWell(t *testing.T) {
	(*dhKeyPair)(nil).wipe()
}

func Test_akeKeys_wipe_HandlesNilWell(t *testing.T) {
	(*akeKeys)(nil).wipe()
}

func Test_ake_wipe_HandlesNilWell(t *testing.T) {
	(*ake)(nil).wipe(true)
}

func Test_ake_wipeGX_HandlesNilWell(t *testing.T) {
	(*ake)(nil).wipeGX()
}

func Test_keyManagementContext_wipeKeys_HandlesNilWell(t *testing.T) {
	(*keyManagementContext)(nil).wipeKeys()
}

func Test_keyManagementContext_wipe_HandlesNilWell(t *testing.T) {
	(*keyManagementContext)(nil).wipe()
}

func Test_counterHistory_wipe_HandlesNilWell(t *testing.T) {
	(*counterHistory)(nil).wipe()
}

func Test_keyPairCounter_wipe_HandlesNilWell(t *testing.T) {
	(*keyPairCounter)(nil).wipe()
}

func Test_macKeyHistory_wipe_HandlesNilWell(t *testing.T) {
	(*macKeyHistory)(nil).wipe()
}

func Test_macKeyUsage_wipe_HandlesNilWell(t *testing.T) {
	(*macKeyUsage)(nil).wipe()
}

func Test_macKey_wipe_HandlesNilWell(t *testing.T) {
	(*macKey)(nil).wipe()
}
