package otr3

import (
	"hash"
	"math/big"
	"strconv"
)

func hashMPIs(h hash.Hash, magic byte, mpis ...*big.Int) []byte {
	h.Reset()
	h.Write([]byte{magic})
	for _, mpi := range mpis {
		h.Write(AppendMPI(nil, mpi))
	}
	return h.Sum(nil)
}

func hashMPIsBN(h hash.Hash, magic byte, mpis ...*big.Int) *big.Int {
	return new(big.Int).SetBytes(hashMPIs(h, magic, mpis...))
}

func bytesToUint16(d []byte) (uint16, error) {
	res, e := strconv.Atoi(string(d))
	return uint16(res), e
}

func makeCopy(i []byte) []byte {
	return append([]byte{}, i...)
}
