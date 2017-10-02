package otr3

import (
	"hash"
	"math/big"
	"strconv"
)

func appendWord(l []byte, r uint32) []byte {
	return append(l, byte(r>>24), byte(r>>16), byte(r>>8), byte(r))
}

func appendShort(l []byte, r uint16) []byte {
	return append(l, byte(r>>8), byte(r))
}

func appendData(l, r []byte) []byte {
	return append(appendWord(l, uint32(len(r))), r...)
}

func appendMPI(l []byte, r *big.Int) []byte {
	return appendData(l, r.Bytes())
}

func appendMPIs(l []byte, r ...*big.Int) []byte {
	for _, mpi := range r {
		l = appendMPI(l, mpi)
	}
	return l
}

func hashMPIs(h hash.Hash, magic byte, mpis ...*big.Int) []byte {
	h.Reset()
	h.Write([]byte{magic})
	for _, mpi := range mpis {
		h.Write(appendMPI(nil, mpi))
	}
	return h.Sum(nil)
}

func hashMPIsBN(h hash.Hash, magic byte, mpis ...*big.Int) *big.Int {
	return new(big.Int).SetBytes(hashMPIs(h, magic, mpis...))
}

func extractWord(d []byte) ([]byte, uint32, bool) {
	if len(d) < 4 {
		return nil, 0, false
	}

	return d[4:], uint32(d[0])<<24 |
		uint32(d[1])<<16 |
		uint32(d[2])<<8 |
		uint32(d[3]), true
}

func extractMPI(d []byte) (newPoint []byte, mpi *big.Int, ok bool) {
	d, mpiLen, ok := extractWord(d)
	if !ok || len(d) < int(mpiLen) {
		return nil, nil, false
	}

	mpi = new(big.Int).SetBytes(d[:int(mpiLen)])
	newPoint = d[int(mpiLen):]
	ok = true
	return
}

func extractMPIs(d []byte) ([]byte, []*big.Int, bool) {
	current, mpiCount, ok := extractWord(d)
	if !ok {
		return nil, nil, false
	}
	result := make([]*big.Int, int(mpiCount))
	for i := 0; i < int(mpiCount); i++ {
		current, result[i], ok = extractMPI(current)
		if !ok {
			return nil, nil, false
		}
	}
	return current, result, true
}

func extractShort(d []byte) ([]byte, uint16, bool) {
	if len(d) < 2 {
		return nil, 0, false
	}

	return d[2:], uint16(d[0])<<8 |
		uint16(d[1]), true
}

func extractData(d []byte) (newPoint []byte, data []byte, ok bool) {
	newPoint, length, ok := extractWord(d)
	if !ok || len(newPoint) < int(length) {
		return d, nil, false
	}

	data = newPoint[:int(length)]
	newPoint = newPoint[int(length):]
	ok = true
	return
}

func bytesToUint16(d []byte) (uint16, error) {
	res, e := strconv.Atoi(string(d))
	return uint16(res), e
}

func makeCopy(i []byte) []byte {
	return append([]byte{}, i...)
}
