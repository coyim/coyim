package otr3

import (
	"fmt"
	"math/big"
	"reflect"
	"runtime"
	"unsafe"
)

func (p *dhKeyPair) wipe() {
	if p == nil {
		return
	}

	wipeBigInt(p.pub)

	tryUnlock(p.priv)
	wipeSecretKeyValue(p.priv)
	p.pub = nil
	p.priv = nil
}

func (k *akeKeys) wipe() {
	if k == nil {
		return
	}
	wipeBytes(k.c)
	k.c = nil
	wipeBytes(k.m1)
	k.m1 = nil
	wipeBytes(k.m2)
	k.m2 = nil
}

func (a *ake) wipe(wipeKeys bool) {
	if a == nil {
		return
	}

	tryUnlock(a.secretExponent)
	wipeSecretKeyValue(a.secretExponent)
	a.secretExponent = nil

	wipeBigInt(a.ourPublicValue)
	a.ourPublicValue = nil

	wipeBigInt(a.theirPublicValue)
	a.theirPublicValue = nil

	wipeBytes(a.r[:])

	a.wipeGX()
	a.revealKey.unlock()
	a.revealKey.wipe()
	a.sigKey.unlock()
	a.sigKey.wipe()

	if wipeKeys {
		a.keys.wipe()
	} else {
		a.keys = keyManagementContext{}
	}
}

func (a *ake) wipeGX() {
	if a == nil {
		return
	}

	wipeBytes(a.xhashedGx)
	wipeBytes(a.encryptedGx[:])
	a.xhashedGx = nil
	a.encryptedGx = nil
}

func (c *keyManagementContext) wipeKeys() {
	if c == nil {
		return
	}

	c.ourCurrentDHKeys.wipe()
	c.ourPreviousDHKeys.wipe()

	wipeBigInt(c.theirCurrentDHPubKey)
	c.theirCurrentDHPubKey = nil

	wipeBigInt(c.theirPreviousDHPubKey)
	c.theirPreviousDHPubKey = nil
}

func (c *keyManagementContext) wipe() {
	if c == nil {
		return
	}

	c.wipeKeys()
	c.ourKeyID = 0
	c.theirKeyID = 0

	for i := range c.oldMACKeys {
		c.oldMACKeys[i].wipe()
		c.oldMACKeys[i] = []byte{}
	}
	c.oldMACKeys = nil

	c.counterHistory.wipe()
	c.macKeyHistory.wipe()
}

func (c *keyManagementContext) wipeAndKeepRevealKeys() keyManagementContext {
	ret := keyManagementContext{}
	ret.oldMACKeys = make([]macKey, len(c.oldMACKeys))
	copy(ret.oldMACKeys, c.oldMACKeys)

	c.wipe()

	return ret
}

func (h *counterHistory) wipe() {
	if h == nil {
		return
	}

	for i := range h.counters {
		h.counters[i].wipe()
		h.counters[i] = nil
	}

	h.counters = nil
}

func (c *keyPairCounter) wipe() {
	if c == nil {
		return
	}

	c.ourKeyID = 0
	c.theirKeyID = 0
	c.ourCounter = 0
	c.theirCounter = 0
}

func (h *macKeyHistory) wipe() {
	if h == nil {
		return
	}

	for i := range h.items {
		h.items[i].wipe()
		h.items[i] = macKeyUsage{} //prevent memory leak
	}
	h.items = nil
}

func (u *macKeyUsage) wipe() {
	if u == nil {
		return
	}

	u.receivingKey.wipe()
	u.receivingKey = macKey{}
}

func (k *macKey) wipe() {
	if k == nil {
		return
	}

	wipeBytes(*k)
	*k = []byte{}
}

func zeroes(n int) []byte {
	return make([]byte, n)
}

func zeroesUint32(n int) []uint32 {
	return make([]uint32, n)
}

func wipeBytes(b []byte) {
	copy(b, zeroes(len(b)))

	runtime.KeepAlive(b)
}

func wipeUint32(b []uint32) {
	copy(b, zeroesUint32(len(b)))

	runtime.KeepAlive(b)
}

func wipeBigInt(k *big.Int) {
	if k == nil {
		return
	}

	k.SetBytes(zeroes(len(k.Bytes())))

	runtime.KeepAlive(k)
}

func wipeSecretKeyValue(k secretKeyValue) {
	if k == nil {
		return
	}

	copy(k, zeroes(len(k)))

	runtime.KeepAlive(k)
}

func setBigInt(dst *big.Int, src *big.Int) *big.Int {
	wipeBigInt(dst)

	ret := big.NewInt(0)
	ret.Set(src)
	return ret
}

func setSecretKeyValue(dst secretKeyValue, src secretKeyValue) secretKeyValue {
	tryUnlock(dst)
	wipeSecretKeyValue(dst)

	ret := make(secretKeyValue, len(src))
	copy(ret, src)
	tryLock(ret)

	return ret
}

func unwrapToStruct(val interface{}) reflect.Value {
	current := reflect.ValueOf(val)
	for current.Kind() != reflect.Struct {
		current = current.Elem()
	}
	return current
}

func unsafeWipeStruct(val reflect.Value) {
	fields := val.NumField()
	for i := 0; i < fields; i++ {
		unsafeWipeField(val.Field(i))
	}
}

var uint8Type = reflect.ValueOf(uint8(1)).Type()
var uint8SliceType = reflect.SliceOf(uint8Type)

var uint32Type = reflect.ValueOf(uint32(1)).Type()
var uint32SliceType = reflect.SliceOf(uint32Type)

func unsafeWipeSliceUint8(val reflect.Value) {
	/* #nosec G103*/
	ss := *(*[]uint8)(unsafe.Pointer(val.UnsafeAddr()))
	wipeBytes(ss)
}

func unsafeWipeSliceUint32(val reflect.Value) {
	/* #nosec G103*/
	ss := *(*[]uint32)(unsafe.Pointer(val.UnsafeAddr()))
	wipeUint32(ss)
}

func unsafeWipeSlice(val reflect.Value) {
	switch val.Type() {
	case uint8SliceType:
		unsafeWipeSliceUint8(val)
	case uint32SliceType:
		unsafeWipeSliceUint32(val)
	default:
		fmt.Printf("Unsupported wipe type: %v\n", val.Type().String())
	}
}

func unsafeWipeField(val reflect.Value) {
	switch val.Kind() {
	case reflect.Struct:
		unsafeWipeStruct(val)
	case reflect.Slice:
		unsafeWipeSlice(val)
	default:
		fmt.Printf("Unsupported wipe kind: %v\n", val.Kind().String())
	}
}

func unsafeWipe(val interface{}) {
	v := unwrapToStruct(val)

	unsafeWipeStruct(v)
}
