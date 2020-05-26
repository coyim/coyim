package otr3

import "encoding/binary"

// SerializeShort returns the big endian serialization of the given value
func SerializeShort(r uint16) []byte {
	d := make([]byte, 2)
	binary.BigEndian.PutUint16(d, r)
	return d
}

// SerializeWord returns the big endian serialization of the given value
func SerializeWord(r uint32) []byte {
	d := make([]byte, 4)
	binary.BigEndian.PutUint32(d, r)
	return d
}

// SerializeLong returns the big endian serialization of the given value
func SerializeLong(r uint64) []byte {
	d := make([]byte, 8)
	binary.BigEndian.PutUint64(d, r)
	return d
}

// DeserializeShort returns the big endian deserialization of the given value
// The buffer is expected to be long enough - otherwise a panic will happen
func DeserializeShort(d []byte) uint16 {
	return binary.BigEndian.Uint16(d)
}

// DeserializeWord returns the big endian deserialization of the given value
// The buffer is expected to be long enough - otherwise a panic will happen
func DeserializeWord(d []byte) uint32 {
	return binary.BigEndian.Uint32(d)
}

// DeserializeLong returns the big endian deserialization of the given value
// The buffer is expected to be long enough - otherwise a panic will happen
func DeserializeLong(d []byte) uint64 {
	return binary.BigEndian.Uint64(d)
}
