package config

import (
	"encoding/hex"

	. "gopkg.in/check.v1"
)

type EncryptedConfigXmppSuite struct{}

var _ = Suite(&EncryptedConfigXmppSuite{})

func bytesFromHex(s string) []byte {
	val, _ := hex.DecodeString(s)
	return val
}

func byteStringFromHex(s string) string {
	val, _ := hex.DecodeString(s)
	return string(val)
}

const testPassword = "vella extol bowel gnome 34th"

var testSalt = bytesFromHex("E18CB93A823465D2797539EBC5F3C0FD")

var testN = 262144 // 2 ** 18

const testR = 8
const testP = 1

var testKey = bytesFromHex("a9af9b3684e680e9ef6b7986c142783f5bf9da26fde00b69dd550451df240998")
var testKeyWrong = bytesFromHex("a9af9b3684e680e9ef6b7986c142783f5bf9da26fde00b69dd560451df240998")
var testMacKey = bytesFromHex("d1d2f87e9280a8a338a1d0e180bd19")
var testMacKeyWrong = bytesFromHex("d1d2f87e9281a8a338a1d0e180bd19")
var testNonce = bytesFromHex("dbd8f7642b05349123d59d1b")
var testNonceWrong = bytesFromHex("dcd8f7642b05349123d59d1b")
var testEncryptedData = bytesFromHex("2b356e6939d1ff771eaa9f9f8866e2aa3732c96913e0d7fa6b4a05d54667c44f78c714af73152b52e2d28cb9b5a2e4e029923465f2cf794d5e9f")
var testEncryptedDataFlip = bytesFromHex("2b356e6939d1ff771eaa9f9f8866e2aa3732c96913e0d7fa6b4a05d54668c44f78c714af73152b52e2d28cb9b5a2e4e029923465f2cf794d5e9f")

func (s *EncryptedConfigXmppSuite) Test_generateKeys(c *C) {
	params := EncryptionParameters{
		saltInternal: testSalt,
		N:            testN,
		R:            testR,
		P:            testP,
	}

	res, res2 := GenerateKeys(testPassword, params)
	c.Assert(res, DeepEquals, testKey)
	c.Assert(res2, DeepEquals, testMacKey)
}

func (s *EncryptedConfigXmppSuite) Test_encryptData(c *C) {
	res := encryptData(testKey, testMacKey, testNonce, "this is some data I want to have encrypted")
	c.Assert(res, DeepEquals, testEncryptedData)
}

func (s *EncryptedConfigXmppSuite) Test_decryptData(c *C) {
	res, e := decryptData(testKey, testMacKey, testNonce, testEncryptedData)
	c.Assert(e, IsNil)
	c.Assert(string(res), DeepEquals, "this is some data I want to have encrypted")

	_, e = decryptData(testKey, testMacKey, testNonce, testEncryptedDataFlip)
	c.Assert(e.Error(), Equals, "cipher: message authentication failed")

	_, e = decryptData(testKeyWrong, testMacKey, testNonce, testEncryptedData)
	c.Assert(e.Error(), Equals, "cipher: message authentication failed")

	_, e = decryptData(testKey, testMacKeyWrong, testNonce, testEncryptedData)
	c.Assert(e.Error(), Equals, "cipher: message authentication failed")

	_, e = decryptData(testKey, testMacKey, testNonceWrong, testEncryptedData)
	c.Assert(e.Error(), Equals, "cipher: message authentication failed")
}

var encryptedDataContent = []byte(`
{
  "Params": {
    "Nonce": "dbd8f7642b05349123d59d1b",
    "Salt": "E18CB93A823465D2797539EBC5F3C0FD",
    "N": 262144,
    "R": 8,
    "P": 1
  },

  "Data": "2b356e6939d1ff771eaa9f9f8866e2aa3732c96913e0d7fa6b4a05d54667c44f78c714af73152b52e2d28cb9b5a2e4e029923465f2cf794d5e9f"
}
`)

func (s *EncryptedConfigXmppSuite) Test_decryptConfiguration(c *C) {
	res, e := decryptConfiguration(encryptedDataContent, func(params EncryptionParameters) ([]byte, []byte, bool) {
		return testKey, testMacKey, true
	})
	c.Assert(e, IsNil)
	c.Assert(string(res), Equals, "this is some data I want to have encrypted")
}

func (s *EncryptedConfigXmppSuite) Test_encryptConfiguration(c *C) {
	p := &EncryptionParameters{
		Nonce: "dbd8f7642b05349123d59d1b",
		Salt:  "E18CB93A823465D2797539EBC5F3C0FD",
		N:     262144,
		R:     8,
		P:     1,
	}
	p.deserialize()

	res, e := encryptConfiguration("this is some data I want to have encrypted", p, func(params EncryptionParameters) ([]byte, []byte, bool) {
		return testKey, testMacKey, true
	})
	c.Assert(e, IsNil)
	c.Assert(string(res), Equals, `{
	"Params": {
		"Nonce": "dbd8f7642b05349123d59d1b",
		"Salt": "e18cb93a823465d2797539ebc5f3c0fd",
		"N": 262144,
		"R": 8,
		"P": 1
	},
	"Data": "2b356e6939d1ff771eaa9f9f8866e2aa3732c96913e0d7fa6b4a05d54667c44f78c714af73152b52e2d28cb9b5a2e4e029923465f2cf794d5e9f"
}`)
}
