package otr3

import (
	"crypto/sha1"
	"hash"
)

func fingerprintHashInstanceForVersion(v int) hash.Hash {
	switch v {
	case 2, 3:
		return sha1.New()
	}

	return nil
}
