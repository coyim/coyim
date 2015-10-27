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

// look for encrypted file first
// if an encrypted file exists, ask for password
// if NO file exists, ask if we should encrypt the file
// if we get a password, keep track of the hash and the parameters for the hash
// when saving the file, generate a new IV and put the parameters first in the file

const encryptedFileEnding = ".enc"

type encryptedData struct {
	Params EncryptionParameters
	Data   string
}

// We will generate a new nonce every time we encrypt, but we will keep the salt the same. This way we can cache the scrypted password

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

func newEncryptionParameters() EncryptionParameters {
	res := EncryptionParameters{
		N: 262144, // 2 ** 18
		R: 8,
		P: 1,
	}
	res.nonceInternal = genRand(nonceLen)
	res.saltInternal = genRand(saltLen)
	return res
}

const aesKeyLen = 32
const macKeyLen = 16
const nonceLen = 12
const saltLen = 16

func GenerateKeys(password string, params EncryptionParameters) ([]byte, []byte) {
	res, _ := scrypt.Key([]byte(password), params.saltInternal, params.N, params.R, params.P, aesKeyLen+macKeyLen)
	return res[0:aesKeyLen], res[aesKeyLen+1:]
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
		return nil, e
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

type KeySupplier func(params EncryptionParameters) ([]byte, []byte, bool)

func decryptConfiguration(content []byte, getKeys KeySupplier) ([]byte, error) {
	data, err := parseEncryptedData(content)
	if err != nil {
		return nil, err
	}

	key, macKey, ok := getKeys(data.Params)
	if !ok {
		return nil, errors.New("no password supplied, aborting")
	}

	ctext, err := hex.DecodeString(data.Data)
	if err != nil {
		return nil, err
	}

	return decryptData(key, macKey, data.Params.nonceInternal, ctext)
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
