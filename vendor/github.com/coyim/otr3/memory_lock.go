package otr3

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/awnumar/memcall"
)

var notifiedLockFailure sync.Once

func tryLock(buf []byte) {
	e := memcall.Lock(buf)
	if e != nil {
		notifiedLockFailure.Do(func() {
			fmt.Printf("Warning: couldn't lock memory pages containing sensitive material: %v\n", e)
		})
	}
}

func tryUnlock(buf []byte) {
	_ = memcall.Unlock(buf)
}
func tryLockBigInt(x *big.Int) {
	if x == nil {
		return
	}
	e := mlockWordSlice(x.Bits())
	if e != nil {
		notifiedLockFailure.Do(func() {
			fmt.Printf("Warning: couldn't lock memory pages containing sensitive material: %v\n", e)
		})
	}
}

func (s *sessionKeys) lock() {
	tryLock(s.sendingAESKey)
	tryLock(s.receivingAESKey)
	tryLock(s.sendingMACKey)
	tryLock(s.receivingMACKey)
	tryLock(s.extraKey)
}

func (s *sessionKeys) unlock() {
	tryUnlock(s.sendingAESKey)
	tryUnlock(s.receivingAESKey)
	tryUnlock(s.sendingMACKey)
	tryUnlock(s.receivingMACKey)
	tryUnlock(s.extraKey)
}

func (a *akeKeys) lock() {
	tryLock(a.c)
	tryLock(a.m1)
	tryLock(a.m2)
}

func (a *akeKeys) unlock() {
	tryUnlock(a.c)
	tryUnlock(a.m1)
	tryUnlock(a.m2)
}

func (priv *DSAPrivateKey) lock() {
	tryLockBigInt(priv.X)
}
