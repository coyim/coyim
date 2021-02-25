package config

import (
	"encoding/hex"
	"errors"

	. "gopkg.in/check.v1"
)

type EncryptedConfigXMPPSuite struct{}

var _ = Suite(&EncryptedConfigXMPPSuite{})

func bytesFromHex(s string) []byte {
	val, _ := hex.DecodeString(s)
	return val
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
		N: testN,
		R: testR,
		P: testP,
	}
	params.saltInternal = testSalt

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
	_ = p.deserialize()

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

func (s *EncryptedConfigXMPPSuite) Test_genRand_panicsOnIOFailure(c *C) {
	origRandReaderRead := randReaderRead
	defer func() {
		randReaderRead = origRandReaderRead
	}()

	randReaderRead = func([]byte) (int, error) {
		return 0, errors.New("problem with ioooooooooo")
	}

	c.Assert(func() { genRand(10) }, PanicMatches, "Failed to read random bytes: problem with io+")
}

func (s *EncryptedConfigXMPPSuite) Test_EncryptionParameters_deserialize_errorsOnInvalidNonce(c *C) {
	p := &EncryptionParameters{}
	p.Nonce = "qqq"
	e := p.deserialize()
	c.Assert(e, ErrorMatches, "encoding/hex:.*")
}

func (s *EncryptedConfigXMPPSuite) Test_EncryptionParameters_deserialize_errorsOnInvalidSalt(c *C) {
	p := &EncryptionParameters{}
	p.Nonce = "abab"
	p.Salt = "qqqq"
	e := p.deserialize()
	c.Assert(e, ErrorMatches, "encoding/hex:.*")
}

func (s *EncryptedConfigXMPPSuite) Test_decryptConfiguration_failsOnBadJSON(c *C) {
	content, params, e := decryptConfiguration([]byte("{"), nil)
	c.Assert(e, ErrorMatches, "unexpected end of JSON input")
	c.Assert(content, IsNil)
	c.Assert(params, IsNil)
}

func (s *EncryptedConfigXMPPSuite) Test_decryptConfiguration_failsOnBadHexInData(c *C) {
	content, params, e := decryptConfiguration([]byte(`{
	"Params": {
		"Nonce": "dbd8f7642b05349123d59d1b",
		"Salt": "e18cb93a823465d2797539ebc5f3c0fd",
		"N": 262144,
		"R": 8,
		"P": 1
	},
	"Data": "qqq"
}`), FunctionKeySupplier(func(params EncryptionParameters, _ bool) ([]byte, []byte, bool) {
		return testKey, testMacKey, true
	}))
	c.Assert(e, ErrorMatches, "encoding/hex: invalid byte.*")
	c.Assert(content, IsNil)
	c.Assert(params, IsNil)
}

func (s *EncryptedConfigXMPPSuite) Test_functionKeySupplier_LastAttemptFailed_setsFlag(c *C) {
	fk := &functionKeySupplier{}
	fk.LastAttemptFailed()
	c.Assert(fk.lastAttemptFailed, Equals, true)
}

func (s *EncryptedConfigXMPPSuite) Test_cachingKeySupplier_LastAttemptFailed_setsFlag(c *C) {
	fk := &cachingKeySupplier{}
	fk.LastAttemptFailed()
	c.Assert(fk.lastAttemptFailed, Equals, true)
}

func (s *EncryptedConfigXMPPSuite) Test_cachingKeySupplier_Invalidate_removesSavedData(c *C) {
	fk := &cachingKeySupplier{
		haveKeys: true,
		key:      []byte{0x01},
		macKey:   []byte{0x02},
	}
	fk.Invalidate()
	c.Assert(fk.haveKeys, Equals, false)
	c.Assert(fk.key, DeepEquals, []byte{})
	c.Assert(fk.macKey, DeepEquals, []byte{})
}

func (s *EncryptedConfigXMPPSuite) Test_CachingKeySupplier_works(c *C) {
	called := 0
	fk := CachingKeySupplier(func(EncryptionParameters, bool) ([]byte, []byte, bool) {
		called++
		return []byte{0x03, 0x05}, []byte{0x42, 0x53}, true
	})

	res1k, res1mk, res1b := fk.GenerateKey(EncryptionParameters{})
	c.Assert(called, Equals, 1)
	res2k, res2mk, res2b := fk.GenerateKey(EncryptionParameters{})
	c.Assert(called, Equals, 1)

	c.Assert(res1k, DeepEquals, res2k)
	c.Assert(res1mk, DeepEquals, res2mk)
	c.Assert(res1b, Equals, res2b)
}

func (s *EncryptedConfigXMPPSuite) Test_CachingKeySupplier_returnsFailureIfGenerationFails(c *C) {
	called := 0
	fk := CachingKeySupplier(func(EncryptionParameters, bool) ([]byte, []byte, bool) {
		called++
		return []byte{0x03, 0x05}, []byte{0x42, 0x53}, false
	})

	res1k, res1mk, res1b := fk.GenerateKey(EncryptionParameters{})
	c.Assert(called, Equals, 1)
	res2k, res2mk, res2b := fk.GenerateKey(EncryptionParameters{})
	c.Assert(called, Equals, 2)

	c.Assert(res1k, DeepEquals, res2k)
	c.Assert(res1mk, DeepEquals, res2mk)
	c.Assert(res1b, Equals, res2b)
}
