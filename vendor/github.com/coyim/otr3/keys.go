package otr3

import (
	"bufio"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/dsa"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"
	"os"

	"github.com/coyim/otr3/sexp"
)

// PublicKey is a public key used to verify signed messages
type PublicKey interface {
	Parse([]byte) ([]byte, bool)
	Fingerprint() []byte
	Verify([]byte, []byte) ([]byte, bool)

	serialize() []byte

	IsSame(PublicKey) bool
}

// PrivateKey is a private key used to sign messages
type PrivateKey interface {
	Parse([]byte) ([]byte, bool)
	Serialize() []byte
	Sign(io.Reader, []byte) ([]byte, error)
	Generate(io.Reader) error
	PublicKey() PublicKey
	IsAvailableForVersion(uint16) bool
}

// GenerateMissingKeys will look through the existing serialized keys and generate new keys to ensure that the functioning of this version of OTR will work correctly. It will only return the newly generated keys, not the old ones
func GenerateMissingKeys(existing [][]byte) ([]PrivateKey, error) {
	var result []PrivateKey
	hasDSA := false

	for _, x := range existing {
		_, typeTag, ok := extractShort(x)
		if ok && typeTag == dsaKeyTypeValue {
			hasDSA = true
		}
	}

	if !hasDSA {
		var priv DSAPrivateKey
		if err := priv.Generate(rand.Reader); err != nil {
			return nil, err
		}
		result = append(result, &priv)
	}

	return result, nil
}

// DSAPublicKey is a DSA public key
type DSAPublicKey struct {
	dsa.PublicKey
}

// DSAPrivateKey is a DSA private key
type DSAPrivateKey struct {
	DSAPublicKey
	dsa.PrivateKey
}

// Account is a holder for the private key associated with an account
// It contains name, protocol and otr private key of an otr Account
type Account struct {
	Name     string
	Protocol string
	Key      PrivateKey
}

func readSymbolAndExpect(r *bufio.Reader, s string) bool {
	res, ok := readPotentialSymbol(r)
	return ok && res == s
}

func readPotentialBigNum(r *bufio.Reader) (*big.Int, bool) {
	res, _ := sexp.ReadValue(r)
	if res != nil {
		if tres, ok := res.(sexp.BigNum); ok {
			return tres.Value().(*big.Int), true
		}
	}
	return nil, false
}

func readPotentialSymbol(r *bufio.Reader) (string, bool) {
	res, _ := sexp.ReadValue(r)
	if res != nil {
		if tres, ok := res.(sexp.Symbol); ok {
			return tres.Value().(string), true
		}
	}
	return "", false
}

func readPotentialString(r *bufio.Reader) (string, bool) {
	res, _ := sexp.ReadValue(r)
	if res != nil {
		if tres, ok := res.(sexp.Sstring); ok {
			return tres.Value().(string), true
		}
	}
	return "", false
}

func readPotentialStringOrSymbol(r *bufio.Reader) (string, bool) {
	res, _ := sexp.ReadValue(r)
	if res != nil {
		if tres, ok := res.(sexp.Sstring); ok {
			return tres.Value().(string), true
		}
		if tres, ok := res.(sexp.Symbol); ok {
			return tres.Value().(string), true
		}
	}
	return "", false
}

// ImportKeysFromFile will read the libotr formatted file given and return all accounts defined in it
func ImportKeysFromFile(fname string) ([]*Account, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ImportKeys(f)
}

// ExportKeysToFile will create the named file (or truncate it) and write all the accounts to that file in libotr format.
func ExportKeysToFile(acs []*Account, fname string) error {
	f, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer f.Close()
	exportAccounts(acs, f)
	return nil
}

// ImportKeys will read the libotr formatted data given and return all accounts defined in it
func ImportKeys(r io.Reader) ([]*Account, error) {
	res, ok := readAccounts(bufio.NewReader(r))
	if !ok {
		return nil, newOtrError("couldn't import data into private key")
	}
	return res, nil
}

func assignParameter(k *dsa.PrivateKey, s string, v *big.Int) bool {
	switch s {
	case "g":
		k.G = v
	case "p":
		k.P = v
	case "q":
		k.Q = v
	case "x":
		k.X = v
	case "y":
		k.Y = v
	default:
		return false
	}
	return true
}

func readAccounts(r *bufio.Reader) ([]*Account, bool) {
	sexp.ReadListStart(r)
	ok1 := readSymbolAndExpect(r, "privkeys")
	ok2 := true
	var as []*Account
	for {
		a, ok, atEnd := readAccount(r)
		ok2 = ok2 && ok
		if atEnd {
			break
		}
		as = append(as, a)
	}
	ok3 := sexp.ReadListEnd(r)
	return as, ok1 && ok2 && ok3
}

func readAccountName(r *bufio.Reader) (string, bool) {
	sexp.ReadListStart(r)
	ok1 := readSymbolAndExpect(r, "name")
	nm, ok2 := readPotentialStringOrSymbol(r)
	ok3 := sexp.ReadListEnd(r)
	return nm, ok1 && ok2 && ok3
}

func readAccountProtocol(r *bufio.Reader) (string, bool) {
	sexp.ReadListStart(r)
	ok1 := readSymbolAndExpect(r, "protocol")
	nm, ok2 := readPotentialSymbol(r)
	ok3 := sexp.ReadListEnd(r)
	return nm, ok1 && ok2 && ok3
}

func readAccount(r *bufio.Reader) (a *Account, ok bool, atEnd bool) {
	if !sexp.ReadListStart(r) {
		return nil, true, true
	}
	ok1 := readSymbolAndExpect(r, "account")
	a = new(Account)
	var ok2, ok3, ok4 bool
	a.Name, ok2 = readAccountName(r)
	a.Protocol, ok3 = readAccountProtocol(r)
	a.Key, ok4 = readPrivateKey(r)
	ok5 := sexp.ReadListEnd(r)
	return a, ok1 && ok2 && ok3 && ok4 && ok5, false
}

func readPrivateKey(r *bufio.Reader) (PrivateKey, bool) {
	sexp.ReadListStart(r)
	ok1 := readSymbolAndExpect(r, "private-key")
	k := new(DSAPrivateKey)
	res, ok2 := readDSAPrivateKey(r)
	if ok2 {
		k.PrivateKey = *res
		k.DSAPublicKey.PublicKey = k.PrivateKey.PublicKey
	}
	ok3 := sexp.ReadListEnd(r)
	return k, ok1 && ok2 && ok3
}

func readDSAPrivateKey(r *bufio.Reader) (*dsa.PrivateKey, bool) {
	sexp.ReadListStart(r)
	ok1 := readSymbolAndExpect(r, "dsa")
	k := new(dsa.PrivateKey)
	for {
		tag, value, end, ok := readParameter(r)
		if !ok {
			return nil, false
		}
		if end {
			break
		}
		if !assignParameter(k, tag, value) {
			return nil, false
		}
	}
	ok2 := sexp.ReadListEnd(r)
	return k, ok1 && ok2
}

func readParameter(r *bufio.Reader) (tag string, value *big.Int, end bool, ok bool) {
	if !sexp.ReadListStart(r) {
		return "", nil, true, true
	}
	tag, ok1 := readPotentialSymbol(r)
	value, ok2 := readPotentialBigNum(r)
	ok = ok1 && ok2
	end = false
	if !sexp.ReadListEnd(r) {
		return "", nil, true, true
	}
	return
}

// IsAvailableForVersion returns true if this key is possible to use with the given version
func (pub *DSAPublicKey) IsAvailableForVersion(v uint16) bool {
	return v == 2 || v == 3
}

// IsSame returns true if the given public key is a DSA public key that is equal to this key
func (pub *DSAPublicKey) IsSame(other PublicKey) bool {
	oth, ok := other.(*DSAPublicKey)
	return ok && pub == oth
}

// ParsePrivateKey is an algorithm indepedent way of parsing private keys
func ParsePrivateKey(in []byte) (index []byte, ok bool, key PrivateKey) {
	var typeTag uint16
	index, typeTag, ok = extractShort(in)
	if !ok {
		return in, false, nil
	}

	switch typeTag {
	case dsaKeyTypeValue:
		key = &DSAPrivateKey{}
		index, ok = key.Parse(in)
		return
	}

	return in, false, nil
}

// ParsePublicKey is an algorithm independent way of parsing public keys
func ParsePublicKey(in []byte) (index []byte, ok bool, key PublicKey) {
	var typeTag uint16
	index, typeTag, ok = extractShort(in)
	if !ok {
		return in, false, nil
	}

	switch typeTag {
	case dsaKeyTypeValue:
		key = &DSAPublicKey{}
		index, ok = key.Parse(in)
		return
	}

	return in, false, nil
}

// Parse takes the given data and tries to parse it into the PublicKey receiver. It will return not ok if the data is malformed or not for a DSA key
func (pub *DSAPublicKey) Parse(in []byte) (index []byte, ok bool) {
	var typeTag uint16
	if index, typeTag, ok = extractShort(in); !ok || typeTag != dsaKeyTypeValue {
		return in, false
	}
	if index, pub.P, ok = extractMPI(index); !ok {
		return in, false
	}
	if index, pub.Q, ok = extractMPI(index); !ok {
		return in, false
	}
	if index, pub.G, ok = extractMPI(index); !ok {
		return in, false
	}
	if index, pub.Y, ok = extractMPI(index); !ok {
		return in, false
	}
	return
}

// Parse will parse a Private Key from the given data, by first parsing the public key components and then the private key component. It returns not ok for the same reasons as PublicKey.Parse.
func (priv *DSAPrivateKey) Parse(in []byte) (index []byte, ok bool) {
	if in, ok = priv.DSAPublicKey.Parse(in); !ok {
		return nil, false
	}

	priv.PrivateKey.PublicKey = priv.DSAPublicKey.PublicKey
	index, priv.X, ok = extractMPI(in)

	return index, ok
}

var dsaKeyType = []byte{0x00, 0x00}
var dsaKeyTypeValue = uint16(0x0000)

func (priv *DSAPrivateKey) serialize() []byte {
	result := priv.DSAPublicKey.serialize()
	return appendMPI(result, priv.PrivateKey.X)
}

// Serialize will return the serialization of the private key to a byte array
func (priv *DSAPrivateKey) Serialize() []byte {
	return priv.serialize()
}

func (pub *DSAPublicKey) serialize() []byte {
	if pub.P == nil || pub.Q == nil || pub.G == nil || pub.Y == nil {
		return nil
	}

	result := dsaKeyType
	result = appendMPI(result, pub.P)
	result = appendMPI(result, pub.Q)
	result = appendMPI(result, pub.G)
	result = appendMPI(result, pub.Y)
	return result
}

// Fingerprint will generate a fingerprint of the serialized version of the key using the provided hash.
func (pub *DSAPublicKey) Fingerprint() []byte {
	b := pub.serialize()
	if b == nil {
		return nil
	}

	h := fingerprintHashInstanceForVersion(3)

	h.Write(b[2:]) // if public key is DSA, ignore the leading 0x00 0x00 for the key type (according to spec)
	return h.Sum(nil)
}

// Sign will generate a signature of a hashed data using dsa Sign.
func (priv *DSAPrivateKey) Sign(rand io.Reader, hashed []byte) ([]byte, error) {
	r, s, err := dsa.Sign(rand, &priv.PrivateKey, hashed)
	if err == nil {
		rBytes := r.Bytes()
		sBytes := s.Bytes()

		out := make([]byte, 40)
		copy(out[20-len(rBytes):], rBytes)
		copy(out[len(out)-len(sBytes):], sBytes)
		return out, nil
	}
	return nil, err
}

// Verify will verify a signature of a hashed data using dsa Verify.
func (pub *DSAPublicKey) Verify(hashed, sig []byte) (nextPoint []byte, sigOk bool) {
	if len(sig) < 2*20 {
		return nil, false
	}
	r := new(big.Int).SetBytes(sig[:20])
	s := new(big.Int).SetBytes(sig[20:40])
	ok := dsa.Verify(&pub.PublicKey, hashed, r, s)
	return sig[20*2:], ok
}

func counterEncipher(key, iv, src, dst []byte) error {
	aesCipher, err := aes.NewCipher(key)

	if err != nil {
		return err
	}

	ctr := cipher.NewCTR(aesCipher, iv)
	ctr.XORKeyStream(dst, src)

	return nil
}

func encrypt(key, data []byte) (dst []byte, err error) {
	dst = make([]byte, len(data))
	err = counterEncipher(key, dst[:aes.BlockSize], data, dst)
	return
}

func decrypt(key, dst, src []byte) error {
	return counterEncipher(key, make([]byte, aes.BlockSize), src, dst)
}

// Import parses the contents of a libotr private key file.
func (priv *DSAPrivateKey) Import(in []byte) bool {
	mpiStart := []byte(" #")

	mpis := make([]*big.Int, 5)

	for i := 0; i < len(mpis); i++ {
		start := bytes.Index(in, mpiStart)
		if start == -1 {
			return false
		}
		in = in[start+len(mpiStart):]
		end := bytes.IndexFunc(in, notHex)
		if end == -1 {
			return false
		}
		hexBytes := in[:end]
		in = in[end:]

		if len(hexBytes)&1 != 0 {
			return false
		}

		mpiBytes := make([]byte, len(hexBytes)/2)
		if _, err := hex.Decode(mpiBytes, hexBytes); err != nil {
			return false
		}

		mpis[i] = new(big.Int).SetBytes(mpiBytes)
	}

	priv.PrivateKey.P = mpis[0]
	priv.PrivateKey.Q = mpis[1]
	priv.PrivateKey.G = mpis[2]
	priv.PrivateKey.Y = mpis[3]
	priv.PrivateKey.X = mpis[4]
	priv.DSAPublicKey.PublicKey = priv.PrivateKey.PublicKey

	a := new(big.Int).Exp(priv.PrivateKey.G, priv.PrivateKey.X, priv.PrivateKey.P)
	return a.Cmp(priv.PrivateKey.Y) == 0
}

// Generate will generate a new DSA Private Key with the randomness provided. The parameter size used is 1024 and 160.
func (priv *DSAPrivateKey) Generate(rand io.Reader) error {
	if err := dsa.GenerateParameters(&priv.PrivateKey.PublicKey.Parameters, rand, dsa.L1024N160); err != nil {
		return err
	}
	if err := dsa.GenerateKey(&priv.PrivateKey, rand); err != nil {
		return err
	}
	priv.DSAPublicKey.PublicKey = priv.PrivateKey.PublicKey
	return nil
}

// PublicKey returns the public key corresponding to this private key
func (priv *DSAPrivateKey) PublicKey() PublicKey {
	return &priv.DSAPublicKey
}

func notHex(r rune) bool {
	if r >= '0' && r <= '9' ||
		r >= 'a' && r <= 'f' ||
		r >= 'A' && r <= 'F' {
		return false
	}

	return true
}

func exportName(n string, w *bufio.Writer) {
	indent := "    "
	w.WriteString(indent)
	w.WriteString("(name \"")
	w.WriteString(n)
	w.WriteString("\")\n")
}

func exportProtocol(n string, w *bufio.Writer) {
	indent := "    "
	w.WriteString(indent)
	w.WriteString("(protocol ")
	w.WriteString(n)
	w.WriteString(")\n")
}

func exportPrivateKey(key PrivateKey, w *bufio.Writer) {
	indent := "    "
	w.WriteString(indent)
	w.WriteString("(private-key\n")
	exportDSAPrivateKey(key.(*DSAPrivateKey), w)
	w.WriteString(indent)
	w.WriteString(")\n")
}

func exportDSAPrivateKey(key *DSAPrivateKey, w *bufio.Writer) {
	indent := "      "
	w.WriteString(indent)
	w.WriteString("(dsa\n")
	exportParameter("p", key.PrivateKey.P, w)
	exportParameter("q", key.PrivateKey.Q, w)
	exportParameter("g", key.PrivateKey.G, w)
	exportParameter("y", key.PrivateKey.Y, w)
	exportParameter("x", key.PrivateKey.X, w)
	w.WriteString(indent)
	w.WriteString(")\n")
}

func exportParameter(name string, val *big.Int, w *bufio.Writer) {
	indent := "        "
	w.WriteString(indent)
	w.WriteString(fmt.Sprintf("(%s #%X#)\n", name, val))
}

func exportAccount(a *Account, w *bufio.Writer) {
	indent := "  "
	w.WriteString(indent)
	w.WriteString("(account\n")
	exportName(a.Name, w)
	exportProtocol(a.Protocol, w)
	exportPrivateKey(a.Key, w)
	w.WriteString(indent)
	w.WriteString(")\n")
}

func exportAccounts(as []*Account, w io.Writer) {
	bw := bufio.NewWriter(w)
	bw.WriteString("(privkeys\n")
	for _, a := range as {
		exportAccount(a, bw)
	}
	bw.WriteString(")\n")
	bw.Flush()
}
