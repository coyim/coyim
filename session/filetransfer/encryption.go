package filetransfer

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"errors"
	"io"
)

type encryptionParameters struct {
	keyType       string
	key           []byte
	iv            []byte
	encryptionKey []byte
	macKey        []byte
}

var encryptionKeyPrefix = []byte("xmpp-encryption-key")
var macKeyPrefix = []byte("xmpp-mac-key")

func deriveKeyWithPrefix(prefix, key []byte, l int) []byte {
	total := append(prefix, key...)
	result := sha256.Sum256(total)
	return result[0:l]
}

func generateSafeRandomBytes(l int) []byte {
	b := make([]byte, l)
	_, _ = rand.Read(b)
	return b
}

func generateEncryptionParameters(enabled bool, genKey func() []byte, keyType string) *encryptionParameters {
	if !enabled {
		return nil
	}
	key := genKey()
	if len(key) == 32 {
		enc := &encryptionParameters{
			keyType:       keyType,
			key:           key,
			iv:            generateSafeRandomBytes(16),
			encryptionKey: deriveKeyWithPrefix(encryptionKeyPrefix, key, 16),
			macKey:        deriveKeyWithPrefix(macKeyPrefix, key, 32),
		}
		return enc
	}
	return nil
}

func (enc *encryptionParameters) totalSize(fileSize int64) int64 {
	if enc == nil {
		return fileSize
	}
	// Size of IV, file size and MAC size
	return int64(16) + fileSize + int64(hmac.New(sha256.New, enc.macKey).Size())
}

func (enc *encryptionParameters) wrapForSending(data io.WriteCloser, ivMacWriter io.Writer) (io.WriteCloser, func()) {
	if enc == nil {
		return data, func() {}
	}

	mac := hmac.New(sha256.New, enc.macKey)
	aesc, _ := aes.NewCipher(enc.encryptionKey)
	blockc := cipher.NewCTR(aesc, enc.iv)

	ivMacWriter.Write(enc.iv)

	ww := &cipher.StreamWriter{S: blockc, W: io.MultiWriter(data, mac)}
	beforeFinish := func() {
		sum := mac.Sum(nil)
		_, _ = ivMacWriter.Write(sum)
	}

	return ww, beforeFinish
}

func (enc *encryptionParameters) wrapForReceiving(r io.Reader) (io.Reader, func() ([]byte, error)) {
	if enc == nil {
		return r, func() ([]byte, error) { return nil, nil }
	}

	hadError := false
	var errorHad *error

	var iv [16]byte
	n, err := io.ReadFull(r, iv[:])
	if n != 16 {
		err = errors.New("couldn't read the IV")
	}
	if err != nil {
		hadError = true
		errorHad = &err
		return r, nil
	}

	mac := hmac.New(sha256.New, enc.macKey)
	aesc, _ := aes.NewCipher(enc.encryptionKey)

	blockc := cipher.NewCTR(aesc, iv[:])

	rr := &cipher.StreamReader{S: blockc, R: io.TeeReader(r, mac)}

	return rr, func() ([]byte, error) {
		if hadError {
			return nil, *errorHad
		}

		readMAC := make([]byte, mac.Size())
		n, err := r.Read(readMAC)
		if n != mac.Size() {
			err = errors.New("couldn't read MAC tag")
		}
		if err != nil {
			return nil, err
		}

		sum := mac.Sum(nil)

		// It's not strictly necessary to use constant time compare for MACs based on hashes, due to the random nature
		// of hashes and also the fact that this MAC tag is a public value. But we do it anyway - there's no harm to it
		// here.
		if subtle.ConstantTimeCompare(readMAC, sum) == 0 {
			return nil, errors.New("bad MAC - transfer integrity broken")
		}

		return enc.macKey, nil
	}
}
