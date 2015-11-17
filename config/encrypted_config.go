package config

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"sync"

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

	if len(p.nonceInternal) == 0 || len(p.saltInternal) == 0 {
		return errDecryptionParamsEmpty
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

var errNoPasswordSupplied = errors.New("no password supplied, aborting")
var errDecryptionFailed = errors.New("decryption failed")
var errDecryptionParamsEmpty = errors.New("decryption params are empty")

func decryptConfiguration(content []byte, ks KeySupplier) ([]byte, *EncryptionParameters, error) {
	data, err := parseEncryptedData(content)
	if err != nil {
		return nil, nil, err
	}

	key, macKey, ok := ks.GenerateKey(data.Params)
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

func encryptConfiguration(content string, params *EncryptionParameters, ks KeySupplier) ([]byte, error) {
	key, macKey, ok := ks.GenerateKey(*params)
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

// KeySupplier is a function that can be used to get key data from a user
type KeySupplier interface {
	GenerateKey(params EncryptionParameters) ([]byte, []byte, bool)
	Invalidate()
}

type functionKeySupplier struct {
	getKeys func(params EncryptionParameters) ([]byte, []byte, bool)
}

// FunctionKeySupplier is a key supplier that wraps a function to ask for the password
func FunctionKeySupplier(getKeys func(params EncryptionParameters) ([]byte, []byte, bool)) KeySupplier {
	return &functionKeySupplier{getKeys}
}

func (fk *functionKeySupplier) Invalidate() {
}

func (fk *functionKeySupplier) GenerateKey(params EncryptionParameters) ([]byte, []byte, bool) {
	return fk.getKeys(params)
}

type cachingKeySupplier struct {
	sync.Mutex
	haveKeys    bool
	key, macKey []byte
	getKeys     func(params EncryptionParameters) ([]byte, []byte, bool)
}

func (ck *cachingKeySupplier) Invalidate() {
	ck.Lock()
	defer ck.Unlock()
	ck.haveKeys = false
	ck.key = []byte{}
	ck.macKey = []byte{}
}

func (ck *cachingKeySupplier) GenerateKey(params EncryptionParameters) ([]byte, []byte, bool) {
	var ok bool
	ck.Lock()
	defer ck.Unlock()
	if !ck.haveKeys {
		ck.key, ck.macKey, ok = ck.getKeys(params)
		if !ok {
			return nil, nil, false
		}
		ck.haveKeys = true
	}
	return ck.key, ck.macKey, true
}

// CachingKeySupplier is a key supplier that only asks the user for a password if it doesn't already have the key material
func CachingKeySupplier(getKeys func(params EncryptionParameters) ([]byte, []byte, bool)) KeySupplier {
	return &cachingKeySupplier{
		getKeys: getKeys,
	}
}
