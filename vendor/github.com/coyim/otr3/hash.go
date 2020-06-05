package otr3

import (
	/* #nosec G505*/
	"crypto/sha1"
	"hash"
)

func fingerprintHashInstanceForVersion(v int) hash.Hash {
	switch v {
	case 2, 3:
		/* #nosec G401*/
		return sha1.New()
	}

	return nil
}
