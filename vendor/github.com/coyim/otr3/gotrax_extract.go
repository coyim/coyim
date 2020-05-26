package otr3

import (
	"math/big"
	"time"
)

// ExtractByte will return the first byte, the rest, and a boolean indicating success
// Data is expected to be in big-endian format
func ExtractByte(d []byte) ([]byte, uint8, bool) {
	if len(d) < 1 {
		return nil, 0, false
	}

	return d[1:], uint8(d[0]), true
}

// ExtractShort will return the first short, the rest, and a boolean indicating success
// Data is expected to be in big-endian format
func ExtractShort(d []byte) ([]byte, uint16, bool) {
	if len(d) < 2 {
		return nil, 0, false
	}

	return d[2:], DeserializeShort(d), true
}

// ExtractWord will return the first word, the rest, and a boolean indicating success
// Data is expected to be in big-endian format
func ExtractWord(d []byte) ([]byte, uint32, bool) {
	if len(d) < 4 {
		return nil, 0, false
	}

	return d[4:], DeserializeWord(d), true
}

// ExtractLong will return the first long, the rest, and a boolean indicating success
// Data is expected to be in big-endian format
func ExtractLong(d []byte) ([]byte, uint64, bool) {
	if len(d) < 8 {
		return nil, 0, false
	}

	return d[8:], DeserializeLong(d), true
}

// ExtractData will return the first data, the rest, and a boolean indicating success
// A Data is serialized as a word indicating length, and then as many bytes as that length
// Data is expected to be in big-endian format
func ExtractData(d []byte) (newPoint []byte, data []byte, ok bool) {
	newPoint, length, ok := ExtractWord(d)
	if !ok || len(newPoint) < int(length) {
		return nil, nil, false
	}

	data = newPoint[:int(length)]
	newPoint = newPoint[int(length):]
	ok = true
	return
}

// ExtractTime will return the first time, the rest, and a boolean indicating success
// Time is encoded as a big-endian 64bit number
func ExtractTime(d []byte) (newPoint []byte, t time.Time, ok bool) {
	newPoint, tt, ok := ExtractLong(d)
	if !ok {
		return nil, time.Time{}, false
	}
	t = time.Unix(int64(tt), 0).In(time.UTC)
	ok = true
	return
}

// ExtractFixedData will return the first l bytes, the rest, and a boolean indicating success
func ExtractFixedData(d []byte, l int) (newPoint []byte, data []byte, ok bool) {
	if len(d) < l {
		return nil, nil, false
	}
	return d[l:], d[0:l], true
}

// ExtractMPI will return the first MPI, the rest, and a boolean indicating success
// The MPI is encoded as a word length, followed by length bytes indicating the minimal representation of the MPI in the obvious format.
// Data is expected to be in big-endian format
func ExtractMPI(d []byte) (newPoint []byte, mpi *big.Int, ok bool) {
	d, mpiLen, ok := ExtractWord(d)
	if !ok || len(d) < int(mpiLen) {
		return nil, nil, false
	}

	mpi = new(big.Int).SetBytes(d[:int(mpiLen)])
	newPoint = d[int(mpiLen):]
	ok = true
	return
}

// ExtractMPIs will return the first len MPIs, the rest, and a boolean indicating success
// The length is indicated by a word, followed by length MPIs
// Data is expected to be in big-endian format
func ExtractMPIs(d []byte) ([]byte, []*big.Int, bool) {
	current, mpiCount, ok := ExtractWord(d)
	if !ok {
		return nil, nil, false
	}
	result := make([]*big.Int, int(mpiCount))
	for i := 0; i < int(mpiCount); i++ {
		current, result[i], ok = ExtractMPI(current)
		if !ok {
			return nil, nil, false
		}
	}
	return current, result, true
}
