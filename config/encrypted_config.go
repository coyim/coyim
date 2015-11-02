package config

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"

	"golang.org/x/crypto/scrypt"
)

const encryptedFileEnding = ".enc"

type encryptedData struct {
	Params EncryptionParameters
	Data   string
}

// We will generate a new nonce every time we encrypt, but we will keep the salt the same. This way we can cache the scrypted password

// EncryptionParameters contains the parameters used for scrypting the password and encrypting the configuration file
type EncryptionParameters struct {
	Nonce string
	Salt  string
	N     int
	R     int
	P     int

	nonceInternal []byte `json:"-"`
	saltInternal  []byte `json:"-"`
}

func genRand(size int) []byte {
	buf := make([]byte, size)
	if _, err := rand.Reader.Read(buf[:]); err != nil {
		panic("Failed to read random bytes: " + err.Error())
	}
	return buf
}

func (p *EncryptionParameters) regenerateNonce() {
	p.nonceInternal = genRand(nonceLen)
}

func newEncryptionParameters() EncryptionParameters {
	res := EncryptionParameters{
		N: 262144, // 2 ** 18
		R: 8,
		P: 1,
	}
	res.regenerateNonce()
	res.saltInternal = genRand(saltLen)
	return res
}

const aesKeyLen = 32
const macKeyLen = 16
const nonceLen = 12
const saltLen = 16

// GenerateKeys takes a password and encryption parameters and generates an AES key and a MAC key using SCrypt
func GenerateKeys(password string, params EncryptionParameters) ([]byte, []byte) {
	res, _ := scrypt.Key([]byte(password), params.saltInternal, params.N, params.R, params.P, aesKeyLen+macKeyLen)
	return res[0:aesKeyLen], res[aesKeyLen:]
}

func encryptData(key, macKey, nonce []byte, plain string) []byte {
	c, _ := aes.NewCipher(key)
	block, _ := cipher.NewGCM(c)
	return block.Seal(nil, nonce, []byte(plain), macKey)
}

func decryptData(key, macKey, nonce, cipherText []byte) ([]byte, error) {
	c, _ := aes.NewCipher(key)
	block, _ := cipher.NewGCM(c)
	res, e := block.Open(nil, nonce, cipherText, macKey)
	if e != nil {
		return nil, errDecryptionFailed
	}
	return res, nil
}

func (p *EncryptionParameters) deserialize() (e error) {
	p.nonceInternal, e = hex.DecodeString(p.Nonce)
	if e != nil {
		return
	}

	p.saltInternal, e = hex.DecodeString(p.Salt)
	if e != nil {
		return
	}

	return nil
}

func (p *EncryptionParameters) serialize() {
	p.Nonce = hex.EncodeToString(p.nonceInternal)
	p.Salt = hex.EncodeToString(p.saltInternal)
}

func parseEncryptedData(content []byte) (ed *encryptedData, e error) {
	data := new(encryptedData)

	e = json.Unmarshal(content, data)
	if e != nil {
		return
	}

	e = data.Params.deserialize()

	return data, e
}

// KeySupplier is a function that can be used to get key data from a user
type KeySupplier func(params EncryptionParameters) ([]byte, []byte, bool)

var errNoPasswordSupplied = errors.New("no password supplied, aborting")
var errDecryptionFailed = errors.New("decryption failed")

func decryptConfiguration(content []byte, getKeys KeySupplier) ([]byte, *EncryptionParameters, error) {
	data, err := parseEncryptedData(content)
	if err != nil {
		return nil, nil, err
	}

	key, macKey, ok := getKeys(data.Params)
	if !ok {
		return nil, nil, errNoPasswordSupplied
	}

	ctext, err := hex.DecodeString(data.Data)
	if err != nil {
		return nil, nil, err
	}

	res, err := decryptData(key, macKey, data.Params.nonceInternal, ctext)
	return res, &data.Params, err
}

func encryptConfiguration(content string, params *EncryptionParameters, getKeys KeySupplier) ([]byte, error) {
	key, macKey, ok := getKeys(*params)
	if !ok {
		return nil, errors.New("no password supplied, aborting")
	}

	ctext := encryptData(key, macKey, params.nonceInternal, content)

	params.serialize()

	dd := encryptedData{
		Params: *params,
		Data:   hex.EncodeToString(ctext),
	}

	return json.MarshalIndent(dd, "", "\t")
}

// CachingKeySupplier is a key supplier that only asks the user for a password if it doesn't already have the key material
func CachingKeySupplier(getKeys KeySupplier) KeySupplier {
	haveKeys := false
	var key, macKey []byte

	return func(params EncryptionParameters) ([]byte, []byte, bool) {
		var ok bool
		if !haveKeys {
			key, macKey, ok = getKeys(params)
			if !ok {
				return nil, nil, false
			}
			haveKeys = true
		}
		return key, macKey, true
	}
}
