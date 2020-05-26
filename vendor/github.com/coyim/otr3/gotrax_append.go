package otr3

import "math/big"

// AppendShort will append the serialization of the given value to the byte array
// Data will be serialized big-endian
func AppendShort(l []byte, r uint16) []byte {
	return append(l, SerializeShort(r)...)
}

// AppendWord will append the serialization of the given value to the byte array
// Data will be serialized big-endian
func AppendWord(l []byte, r uint32) []byte {
	return append(l, SerializeWord(r)...)
}

// AppendLong will append the serialization of the given value to the byte array
// Data will be serialized big-endian
func AppendLong(l []byte, r uint64) []byte {
	return append(l, SerializeLong(r)...)
}

// AppendData will append the serialization of the given value to the byte array
// Data will be serialized big-endian
func AppendData(l, r []byte) []byte {
	return append(AppendWord(l, uint32(len(r))), r...)
}

// AppendMPI will append the serialization of the given value to the byte array
// Data will be serialized big-endian
func AppendMPI(l []byte, r *big.Int) []byte {
	return AppendData(l, r.Bytes())
}

// AppendMPIs will append the serialization of the given values to the byte array
// Data will be serialized big-endian
func AppendMPIs(l []byte, r ...*big.Int) []byte {
	for _, mpi := range r {
		l = AppendMPI(l, mpi)
	}
	return l
}
