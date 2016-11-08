package config

import (
	"encoding/hex"

	. "gopkg.in/check.v1"
)

type EncryptedConfigXMPPSuite struct{}

var _ = Suite(&EncryptedConfigXMPPSuite{})

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
var testMacKey = bytesFromHex("78d1d2f87e9280a8a338a1d0e180bd19")
var testMacKeyWrong = bytesFromHex("78d1d2f87e9281a8a338a1d0e180bd19")
var testNonce = bytesFromHex("dbd8f7642b05349123d59d1b")
var testNonceWrong = bytesFromHex("dcd8f7642b05349123d59d1b")
var testEncryptedData = bytesFromHex("2b356e6939d1ff771eaa9f9f8866e2aa3732c96913e0d7fa6b4a05d54667c44f78c714af73152b52e2d2a61e79b8671f85a27505cb9c5477ed75")
var testEncryptedDataFlip = bytesFromHex("2b356e6939d1ff771eaa9f9f8866e2aa3732c96913e0d7fa6b4a05d54667c44f78c714af73152b52e2d2a61e79b8671f85a27505cb9c5477ed76")

func (s *EncryptedConfigXMPPSuite) Test_generateKeys(c *C) {
	params := EncryptionParameters{
		saltInternal: testSalt,
		N:            testN,
		R:            testR,
		P:            testP,
	}

	res, res2 := GenerateKeys(testPassword, params)
	c.Assert(res, DeepEquals, testKey)
	c.Assert(res2, DeepEquals, testMacKey)
	c.Assert(len(res), Equals, aesKeyLen)
	c.Assert(len(res2), Equals, macKeyLen)
}

func (s *EncryptedConfigXMPPSuite) Test_encryptData(c *C) {
	res := encryptData(testKey, testMacKey, testNonce, "this is some data I want to have encrypted")
	c.Assert(res, DeepEquals, testEncryptedData)
}

func (s *EncryptedConfigXMPPSuite) Test_decryptData(c *C) {
	res, e := decryptData(testKey, testMacKey, testNonce, testEncryptedData)
	c.Assert(e, IsNil)
	c.Assert(string(res), DeepEquals, "this is some data I want to have encrypted")

	_, e = decryptData(testKey, testMacKey, testNonce, testEncryptedDataFlip)
	c.Assert(e, Equals, errDecryptionFailed)

	_, e = decryptData(testKeyWrong, testMacKey, testNonce, testEncryptedData)
	c.Assert(e, Equals, errDecryptionFailed)

	_, e = decryptData(testKey, testMacKeyWrong, testNonce, testEncryptedData)
	c.Assert(e, Equals, errDecryptionFailed)

	_, e = decryptData(testKey, testMacKey, testNonceWrong, testEncryptedData)
	c.Assert(e, Equals, errDecryptionFailed)
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

  "Data": "2b356e6939d1ff771eaa9f9f8866e2aa3732c96913e0d7fa6b4a05d54667c44f78c714af73152b52e2d2a61e79b8671f85a27505cb9c5477ed75"
}
`)

func (s *EncryptedConfigXMPPSuite) Test_decryptConfiguration(c *C) {
	res, _, e := decryptConfiguration(encryptedDataContent, FunctionKeySupplier(func(params EncryptionParameters, _ bool) ([]byte, []byte, bool) {
		return testKey, testMacKey, true
	}))
	c.Assert(e, IsNil)
	c.Assert(string(res), Equals, "this is some data I want to have encrypted")
}

func (s *EncryptedConfigXMPPSuite) Test_encryptConfiguration(c *C) {
	p := &EncryptionParameters{
		Nonce: "dbd8f7642b05349123d59d1b",
		Salt:  "E18CB93A823465D2797539EBC5F3C0FD",
		N:     262144,
		R:     8,
		P:     1,
	}
	p.deserialize()

	res, e := encryptConfiguration("this is some data I want to have encrypted", p, FunctionKeySupplier(func(params EncryptionParameters, _ bool) ([]byte, []byte, bool) {
		return testKey, testMacKey, true
	}))
	c.Assert(e, IsNil)
	c.Assert(string(res), Equals, `{
	"Params": {
		"Nonce": "dbd8f7642b05349123d59d1b",
		"Salt": "e18cb93a823465d2797539ebc5f3c0fd",
		"N": 262144,
		"R": 8,
		"P": 1
	},
	"Data": "2b356e6939d1ff771eaa9f9f8866e2aa3732c96913e0d7fa6b4a05d54667c44f78c714af73152b52e2d2a61e79b8671f85a27505cb9c5477ed75"
}`)
}

func (s *EncryptedConfigXMPPSuite) Test_deserializeConfigurationEmptyErr(c *C) {
	p := &EncryptionParameters{
		Nonce: "",
		Salt:  "",
		N:     262144,
		R:     8,
		P:     1,
	}
	e := p.deserialize()
	c.Assert(e, Equals, errDecryptionParamsEmpty)
}
