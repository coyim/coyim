package filetransfer

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"errors"
	"hash"
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

	_, _ = ivMacWriter.Write(enc.iv)

	ww := &cipher.StreamWriter{S: blockc, W: ioMultiWriter(data, mac)}
	beforeFinish := func() {
		sum := mac.Sum(nil)
		_, _ = ivMacWriter.Write(sum)
	}

	return ww, beforeFinish
}

func (enc *encryptionParameters) wrapForReceiving(r io.Reader) (io.Reader, func() ([]byte, error), error) {
	if enc == nil {
		return r, func() ([]byte, error) { return nil, nil }, nil
	}

	er := &encryptedReceiver{enc: enc, r: r}
	return er.receive()
}

type encryptedReceiver struct {
	enc             *encryptionParameters
	r               io.Reader
	iv              [16]byte
	mac             hash.Hash
	encryptedReader *cipher.StreamReader
}

func (er *encryptedReceiver) readIV() error {
	n, err := ioReadFull(er.r, er.iv[:])
	if n != 16 {
		return errors.New("couldn't read the IV")
	}

	return err
}

func (er *encryptedReceiver) createMac() {
	er.mac = hmac.New(sha256.New, er.enc.macKey)
}

func (er *encryptedReceiver) createBlockCipher() cipher.Stream {
	aesc, _ := aes.NewCipher(er.enc.encryptionKey)
	return cipher.NewCTR(aesc, er.iv[:])
}

func (er *encryptedReceiver) createEncryptedReader() {
	er.createMac()

	er.encryptedReader = &cipher.StreamReader{
		S: er.createBlockCipher(),
		R: ioTeeReader(er.r, er.mac),
	}
}

func (er *encryptedReceiver) receive() (io.Reader, func() ([]byte, error), error) {
	if err := er.readIV(); err != nil {
		return er.r, nil, err
	}

	er.createEncryptedReader()

	return er.encryptedReader, er.macVerifier, nil
}

func (er *encryptedReceiver) readMacFromSender() ([]byte, error) {
	readMAC := make([]byte, er.mac.Size())
	n, err := er.r.Read(readMAC)
	if n != er.mac.Size() {
		err = errors.New("couldn't read MAC tag")
	}

	return readMAC, err
}

func (er *encryptedReceiver) macVerifier() ([]byte, error) {
	readMAC, err := er.readMacFromSender()
	if err != nil {
		return nil, err
	}

	sum := er.mac.Sum(nil)

	// It's not strictly necessary to use constant time compare for MACs based on hashes, due to the random nature
	// of hashes and also the fact that this MAC tag is a public value. But we do it anyway - there's no harm to it
	// here.
	if subtle.ConstantTimeCompare(readMAC, sum) == 0 {
		return nil, errors.New("bad MAC - transfer integrity broken")
	}

	return er.enc.macKey, nil
}
