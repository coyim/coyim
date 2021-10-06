package filetransfer

import (
	"bytes"
	"errors"
	"io"

	. "gopkg.in/check.v1"
)

var (
	testDataEncryptedContent = []byte{0x42, 0xf9, 0x6b, 0x9e, 0x70, 0x2d, 0xf8}

	testDataMac = []byte{
		0x2, 0x12, 0xac, 0x1b, 0xc3, 0xf6, 0x66, 0xe1,
		0x54, 0xb9, 0x95, 0xf9, 0xbd, 0x70, 0xf, 0x6a,
		0xad, 0x4a, 0xf3, 0x3c, 0x8d, 0x95, 0x6b, 0x26,
		0xe4, 0x78, 0x26, 0x77, 0x41, 0x81, 0x49, 0xfc,
	}

	testDataIV = []byte{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F,
	}

	testDataEncryptionKey = []byte{
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
		0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F,
	}

	testDataMacKey = []byte{
		0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27,
		0x28, 0x29, 0x2A, 0x2B, 0x2C, 0x2D, 0x2E, 0x2F,
		0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37,
		0x38, 0x39, 0x3A, 0x3B, 0x3C, 0x3D, 0x3E, 0x3F,
	}

	testDataContent = []byte{0xDE, 0xAD, 0xBE, 0xEF, 0x00, 0x01, 0x02}
)

type EncryptionSuite struct{}

var _ = Suite(&EncryptionSuite{})

func (s *EncryptionSuite) Test_generateEncryptionParameters_withNotEnabled(c *C) {
	enc := generateEncryptionParameters(false, nil, "")
	c.Assert(enc, IsNil)
}

var testDataRawKeyData = []byte{
	0x01, 0x02, 0x03, 0x04, 0x01, 0x02, 0x03, 0x04,
	0x02, 0x02, 0x03, 0x04, 0x01, 0x02, 0x03, 0x04,
	0x03, 0x02, 0x03, 0x04, 0x01, 0x02, 0x03, 0x04,
	0x04, 0x02, 0x03, 0x04, 0x01, 0x02, 0x03, 0x04,
}

func (s *EncryptionSuite) Test_generateEncryptionParameters_createsBasicEncryptionParameters(c *C) {
	keyGen := func() []byte {
		return testDataRawKeyData
	}
	enc := generateEncryptionParameters(true, keyGen, "something")
	c.Assert(enc, Not(IsNil))
	c.Assert(enc.keyType, Equals, "something")
	c.Assert(enc.key, DeepEquals, testDataRawKeyData)
	c.Assert(enc.encryptionKey, DeepEquals, []byte{
		0x91, 0x61, 0xa9, 0xaf, 0x35, 0xc6, 0x31, 0xb4,
		0x71, 0xcf, 0xd8, 0x53, 0x00, 0xde, 0xae, 0xf8})
	c.Assert(enc.macKey, DeepEquals, []byte{
		0x0d, 0xa7, 0x60, 0xbd, 0x58, 0x6d, 0x0f, 0x2d,
		0x8b, 0x81, 0xb9, 0xe6, 0xe3, 0x9b, 0x64, 0x67,
		0x62, 0x20, 0xf9, 0x88, 0x0e, 0x55, 0x24, 0x36,
		0xf4, 0x8c, 0xa5, 0x30, 0x12, 0xe2, 0x4f, 0x6c})
	c.Assert(enc.iv, Not(IsNil))
}

func (s *EncryptionSuite) Test_generateEncryptionParameters_wrongKeySizeReturnsNilObject(c *C) {
	keyGen := func() []byte {
		return testDataRawKeyData[:31]
	}
	enc := generateEncryptionParameters(true, keyGen, "something")
	c.Assert(enc, IsNil)
}

func (s *EncryptionSuite) Test_totalSize_returnsOriginalSizeForNilEncryptionObject(c *C) {
	var enc *encryptionParameters
	res := enc.totalSize(42)
	c.Assert(res, Equals, int64(42))
}

func (s *EncryptionSuite) Test_totalSize_returnsTotalSizeBasedOnEncryptionParameters(c *C) {
	enc := &encryptionParameters{
		macKey: []byte{1, 2, 3},
	}
	res := enc.totalSize(44)
	c.Assert(res, Equals, int64(92))
}

func (s *EncryptionSuite) Test_wrapForReceiving_failsIfNotEnoughDataReadForIV(c *C) {
	enc := &encryptionParameters{}
	orgIoReadFull := ioReadFull
	defer func() {
		ioReadFull = orgIoReadFull
	}()
	ioReadFull = func(io.Reader, []byte) (int, error) {
		return 42, nil
	}

	testReader := &bytes.Buffer{}

	res, f, err := enc.wrapForReceiving(testReader)
	c.Assert(err, ErrorMatches, "couldn't read the IV")
	c.Assert(res, DeepEquals, testReader)
	c.Assert(f, IsNil)
}

func (s *EncryptionSuite) Test_wrapForReceiving_failsIfReaderIsBroken(c *C) {
	enc := &encryptionParameters{}
	orgIoReadFull := ioReadFull
	defer func() {
		ioReadFull = orgIoReadFull
	}()
	ioReadFull = func(io.Reader, []byte) (int, error) {
		return 16, errors.New("something bad, right?")
	}

	testReader := &bytes.Buffer{}

	res, f, err := enc.wrapForReceiving(testReader)
	c.Assert(err, ErrorMatches, "something bad, right\\?")
	c.Assert(res, DeepEquals, testReader)
	c.Assert(f, IsNil)
}

func (s *EncryptionSuite) Test_wrapForReceiving_works(c *C) {
	enc := &encryptionParameters{
		macKey:        testDataMacKey,
		encryptionKey: testDataEncryptionKey,
	}

	var data []byte = nil
	data = append(data, testDataIV...)
	data = append(data, testDataEncryptedContent...)
	data = append(data, testDataMac...)

	testReader := bytes.NewBuffer(data)

	res, f, _ := enc.wrapForReceiving(testReader)
	c.Assert(res, Not(DeepEquals), testReader)
	c.Assert(f, Not(IsNil))

	output := [7]byte{}
	_, _ = res.Read(output[:])
	c.Assert(output[:], DeepEquals, testDataContent)

	more, e := f()

	c.Assert(e, IsNil)
	c.Assert(more, DeepEquals, testDataMacKey)
}

type readerWithError struct {
	data [][]byte
	n    int
	e    error
}

func (r *readerWithError) Read(out []byte) (n int, e error) {
	if len(r.data) == 0 {
		return r.n, r.e
	}
	current := r.data[0]
	r.data = r.data[1:]
	copy(out, current)
	return len(current), nil
}

func (s *EncryptionSuite) Test_wrapForReceiving_failsWhenReadingMacTag(c *C) {
	enc := &encryptionParameters{
		macKey:        testDataMacKey,
		encryptionKey: testDataEncryptionKey,
	}

	var data []byte = nil
	data = append(data, testDataIV...)
	data = append(data, testDataEncryptedContent...)

	testReader := bytes.NewBuffer(data)

	res, f, _ := enc.wrapForReceiving(testReader)
	c.Assert(res, Not(DeepEquals), testReader)
	c.Assert(f, Not(IsNil))

	output := [7]byte{}
	_, _ = res.Read(output[:])
	c.Assert(output[:], DeepEquals, testDataContent)

	more, e := f()

	c.Assert(e, ErrorMatches, "couldn't read MAC tag")
	c.Assert(more, IsNil)
}

func (s *EncryptionSuite) Test_wrapForReceiving_failsInAnotherWayWhenReadingMacTag(c *C) {
	enc := &encryptionParameters{
		macKey:        testDataMacKey,
		encryptionKey: testDataEncryptionKey,
	}

	testReader := &readerWithError{
		data: [][]byte{
			testDataIV,
			testDataEncryptedContent,
		},
		e: errors.New("oh no"),
		n: 32,
	}

	res, f, _ := enc.wrapForReceiving(testReader)
	c.Assert(res, Not(DeepEquals), testReader)
	c.Assert(f, Not(IsNil))

	output := [7]byte{}
	_, _ = res.Read(output[:])
	c.Assert(output[:], DeepEquals, testDataContent)

	more, e := f()

	c.Assert(e, ErrorMatches, "oh no")
	c.Assert(more, IsNil)
}

func (s *EncryptionSuite) Test_wrapForReceiving_incorrectMacTag(c *C) {
	enc := &encryptionParameters{
		macKey:        testDataMacKey,
		encryptionKey: testDataEncryptionKey,
	}

	var data []byte = nil
	data = append(data, testDataIV...)
	data = append(data, testDataEncryptedContent...)
	data = append(data, 42)
	data = append(data, testDataMac...)

	testReader := bytes.NewBuffer(data)

	res, f, _ := enc.wrapForReceiving(testReader)
	c.Assert(res, Not(DeepEquals), testReader)
	c.Assert(f, Not(IsNil))

	output := [7]byte{}
	_, _ = res.Read(output[:])
	c.Assert(output[:], DeepEquals, testDataContent)

	more, e := f()

	c.Assert(e, ErrorMatches, "bad MAC - transfer integrity broken")
	c.Assert(more, IsNil)
}

type nopWriterCloser struct {
	io.Writer
}

func (nopWriterCloser) Close() error {
	return nil
}

func (s *EncryptionSuite) Test_wrapForSending_doesntDoAnythingWithoutEncryptionParameters(c *C) {
	var enc *encryptionParameters

	data := &nopWriterCloser{&bytes.Buffer{}}
	d2, f := enc.wrapForSending(data, data)

	f()

	c.Assert(d2, Equals, data)
}

func (s *EncryptionSuite) Test_wrapForSending_works(c *C) {
	enc := &encryptionParameters{
		iv:            testDataIV,
		macKey:        testDataMacKey,
		encryptionKey: testDataEncryptionKey,
	}

	orgData := &bytes.Buffer{}
	macWriter := &bytes.Buffer{}

	d2, f := enc.wrapForSending(&nopWriterCloser{orgData}, macWriter)
	_, _ = d2.Write(testDataContent)
	_ = d2.Close()

	c.Assert(orgData.Bytes(), DeepEquals, testDataEncryptedContent)
	c.Assert(macWriter.Bytes(), DeepEquals, testDataIV)

	f()

	c.Assert(macWriter.Bytes(), DeepEquals, append(testDataIV, testDataMac...))
}
