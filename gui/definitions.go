package gui

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

type _escLocalFS struct{}

var _escLocal _escLocalFS

type _escStaticFS struct{}

var _escStatic _escStaticFS

type _escDirectory struct {
	fs   http.FileSystem
	name string
}

type _escFile struct {
	compressed string
	size       int64
	modtime    int64
	local      string
	isDir      bool

	once sync.Once
	data []byte
	name string
}

func (_escLocalFS) Open(name string) (http.File, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	return os.Open(f.local)
}

func (_escStaticFS) prepare(name string) (*_escFile, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	var err error
	f.once.Do(func() {
		f.name = path.Base(name)
		if f.size == 0 {
			return
		}
		var gr *gzip.Reader
		b64 := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(f.compressed))
		gr, err = gzip.NewReader(b64)
		if err != nil {
			return
		}
		f.data, err = ioutil.ReadAll(gr)
	})
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (fs _escStaticFS) Open(name string) (http.File, error) {
	f, err := fs.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.File()
}

func (dir _escDirectory) Open(name string) (http.File, error) {
	return dir.fs.Open(dir.name + name)
}

func (f *_escFile) File() (http.File, error) {
	type httpFile struct {
		*bytes.Reader
		*_escFile
	}
	return &httpFile{
		Reader:   bytes.NewReader(f.data),
		_escFile: f,
	}, nil
}

func (f *_escFile) Close() error {
	return nil
}

func (f *_escFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

func (f *_escFile) Stat() (os.FileInfo, error) {
	return f, nil
}

func (f *_escFile) Name() string {
	return f.name
}

func (f *_escFile) Size() int64 {
	return f.size
}

func (f *_escFile) Mode() os.FileMode {
	return 0
}

func (f *_escFile) ModTime() time.Time {
	return time.Unix(f.modtime, 0)
}

func (f *_escFile) IsDir() bool {
	return f.isDir
}

func (f *_escFile) Sys() interface{} {
	return f
}

// FS returns a http.Filesystem for the embedded assets. If useLocal is true,
// the filesystem's contents are instead used.
func FS(useLocal bool) http.FileSystem {
	if useLocal {
		return _escLocal
	}
	return _escStatic
}

// Dir returns a http.Filesystem for the embedded assets on a given prefix dir.
// If useLocal is true, the filesystem's contents are instead used.
func Dir(useLocal bool, name string) http.FileSystem {
	if useLocal {
		return _escDirectory{fs: _escLocal, name: name}
	}
	return _escDirectory{fs: _escStatic, name: name}
}

// FSByte returns the named file from the embedded assets. If useLocal is
// true, the filesystem's contents are instead used.
func FSByte(useLocal bool, name string) ([]byte, error) {
	if useLocal {
		f, err := _escLocal.Open(name)
		if err != nil {
			return nil, err
		}
		b, err := ioutil.ReadAll(f)
		f.Close()
		return b, err
	}
	f, err := _escStatic.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.data, nil
}

// FSMustByte is the same as FSByte, but panics if name is not present.
func FSMustByte(useLocal bool, name string) []byte {
	b, err := FSByte(useLocal, name)
	if err != nil {
		panic(err)
	}
	return b
}

// FSString is the string version of FSByte.
func FSString(useLocal bool, name string) (string, error) {
	b, err := FSByte(useLocal, name)
	return string(b), err
}

// FSMustString is the string version of FSMustByte.
func FSMustString(useLocal bool, name string) string {
	return string(FSMustByte(useLocal, name))
}

var _escData = map[string]*_escFile{

	"/definitions/AccountDetails.xml": {
		local:   "definitions/AccountDetails.xml",
		size:    32827,
		modtime: 1489449600,
		compressed: `
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
ZlXYO4AAAA==
`,
	},

	"/definitions/AccountRegistration.xml": {
		local:   "definitions/AccountRegistration.xml",
		size:    2956,
		modtime: 1489449600,
		compressed: `
H4sIAAAJbogA/8RWzW7bPBC8f09B8K7PTQ89yQKStHXcBnGRqAXai7Gm1hITmhTIlX/evrDoxJZFKU7S
ojdJ2BnuzgxFxlIT2jkITP5jLDazexTEhALnhnxEDx8lKJNzJrMhz/zztpCxuLSmREsbpmGBQ06SFHJG
FrRTQDBTOOQbdDy5LIxxyIA5tEu0jAyzmEtHaNnGVJaBEKbSFA8eKXcrOJlrUDt+oYxDzgrQmUI75EZP
BWiBaurLOBuEYBkqJIxwiZqeQ3u4KKTKWC2LBhXVr0O+nJn1bvKQTBdm7TX6cVjXVmkBNpeaJ2fvjqcN
VbsShNR5Z3nd2/491Ng1zFD51hboHOTIDwHtNZUHBHz8phAcMuHtpAIfDV0VaHFrJVuBpoa/oB/d/T80
QaiBlYWSJ2QrDCPigR+y8a0E8SB13s+M6xJ0xpM5KNdB3gbNpVL97bTWjgc7Y17k1MjKzBuVb5/6u/I5
imaGyCy68hECWrOK9rF6fypMGFUt9B754RDYRLZGfTaYPkaRT94x8iUJvauJumYKcd1XjuR8w5NR+nX6
5ftdOv78c3o7Hl2l3SyhALLOEAYnwDlFQASi4EmndSEkmfIUYCCU7DCYrzHs0ixm5sKsU1zToW8nOFaA
i1CT3fAk7dxIf0Has38qbXuYrh9D67QBQdLoKVgE3v/XuKiIjH46fmbN4yc0p7ESNcF2AR/7ye340016
no4nN9Orye341+QmPb/u/DOcGBff2K4r0pE/bt+0vS9ril5Pj64MUjxgdsqloc+1DntfLwWUpdq8SYnz
LcNLwi1ARxnOoVLUd5idrKKDJf4BDXu2SKiguX9iv0uilcxyJPcEaXxmFl1ptPMa1AncpzEeNGqfZfDG
sZ2QQ74V0vOBNyRAd/SxbnM/VTw4uH//DgAA///B9LfOjAsAAA==
`,
	},

	"/definitions/AddContact.xml": {
		local:   "definitions/AddContact.xml",
		size:    9384,
		modtime: 1489449600,
		compressed: `
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
`,
	},

	"/definitions/AskForPassword.xml": {
		local:   "definitions/AskForPassword.xml",
		size:    4788,
		modtime: 1489449600,
		compressed: `
H4sIAAAJbogA/+RY33PiNhB+71+h0buHpjPtE3iGyxGONoUM8bXTvnjW8mJ0CImR1hD++w4YihPLPhMn
/TH3Fsf+Pq2+3U/apS81oV2AwPA7xvom+YKCmFDg3ICPafVRgjIZZzId8KFb3Rn7AM7tjE35AcBYf2PN
Bi3tmYY1DvhO6tTsgo1xkqTRPBxHv8S/T6bxw+wxvh1No9G83ztj/BQkSSFnZEE7BQSJwgHfo+Ph6BAs
25vcss0pjAqZk5kGdaJKUSFhgFvUxNkSdKrQDrjRsQAtUMVnmriAcdY70YilVCk7qqNBBcfHAd8m5um0
cZ9aH8xTIdVv5e+qO1yDzaTm4c33L8O/rH159q10DwmqYq01OgcZ8jKguqQqAB5RfSH4GJagZHZK6PB+
Mp7Gj9FwHrWFf8kdycW+wP/8+TGa3P0R34/uWhPsLGx4SDbHtog1PAU7mdIyEEuwjoc/evX2QR0qFEeV
eBjNP4/aL3nIbJAYIrOuSzBj/V6R0VLSey+y3qYKxlamRRFkh7/eJLQq0Jpd4DYgpM54ePNDW5gwKl/r
C/KnWmBlp/7dlmoehDC5pl9Ppc9egq+p/2HBVRedj8xby/PJ+FNtMftYqoYaTT82EgCRlUlO6Kovy6/P
jkGZLYmzLagcBzwxKuU9D22vnrdaqOe9gFhJnbUQHhcUABGIJQ9ra86HJLNpA+z3vLFUDMU6ldoU1q3q
7MpT8s0Vvvn/KXy4g9/EyQ81PUET279i5f+KqxqK5T1yPtJk95ecl7rIZm1dsDAid82Xvw+6lU4mUkna
83AByjWDn3WOIEhugfB514iHfhDT2rbx3VL8alt3SHFXXzvY4tmSxb+7mPsRtljb9TdRftMOr+3W3sfh
t0sUqw85kdHVGqj2Hv+0GTqo8ZrmvTI+Hg4Vo2OwCLy5sy80/HueTI6P8fOp0rdbYyVqgsvoPZtPRtNo
GE1m0/jTbD75czaNhvddW/Fyik+xFQN1a4uHt8fv2x/IQkmxwvSqKb4pZ968ddTArDodcbdGaxRXTSIC
dJDiAnJFLa7Hrwva9YLzKPpVuzz/oPTy8qLfK/1O9VcAAAD//+EBjDW0EgAA
`,
	},

	"/definitions/AskToEncrypt.xml": {
		local:   "definitions/AskToEncrypt.xml",
		size:    650,
		modtime: 1489449600,
		compressed: `
H4sIAAAJbogA/4SS0WrbQBBF3/MVwz638Q/IKm1qQiixS60S8iRGq5E99XpG7Iyq6u+LohoaMPhV3Hvu
QTsFi1PuMFJ5B1Bo84uiQ0xotg6PfnomMzzQV8akhwDcrsNnO1W6kZin3sPcAij6rD1ln0DwTOswsrQ6
fuzV2FkllI/Vt/rlaVt/3+3rh8222vwoVpfOdcRZW0yh9DzQraizJwrgGcUSOjaJ1mEiC+U/S4gqHR+G
jLMNdJxuM+mPX0W+6JBamHSAxCcCVzD8TfOHfGUGWAAFaPGgFjrNZ/RPUB3ZIKJAQ2B8EO44onia4KyZ
wCgOmT5AM/jb1sgpgajP8Vlm3sUYyezaMneLoBqBHwl6NBs1t/fwemGppOkNZqfFauFckjAeSeBBp6dn
MMfsdn/rjzWDu4otb/3lZ1Xttvv6dbOvt7v31WK1XFl5V6z+u76/AQAA///tme3oigIAAA==
`,
	},

	"/definitions/AuthorizeSubscription.xml": {
		local:   "definitions/AuthorizeSubscription.xml",
		size:    399,
		modtime: 1489449600,
		compressed: `
H4sIAAAJbogA/4SRTWrDMBCF9z6F0D74ArKgP8KEErutHEpXQranQa0iqZoxwbcvQRQSCGQ3i/d9w8wT
LhDkLzuBrBgTcfyGidjkLWLDW/rZAaI9wLOzPh44c3PD5zKf84yJlGOCTCsL9ggNJ0ceOKNsA3pLdvTQ
8BWQS72MOGWXyMXAMvwugCTqf1xWN3UnF+Z42qSI7sxx2Q4v5mPbmddemyfVDer90nFLcYyz9VxSXuBu
tBy7oTVBWbVTWj+0yrztlR62fXfPMC5EMWCBH/fD0HfafCptuv4aFXV5taxEfVHBXwAAAP//0Csv4o8B
AAA=
`,
	},

	"/definitions/BadUsernameNotification.xml": {
		local:   "definitions/BadUsernameNotification.xml",
		size:    538,
		modtime: 1489449600,
		compressed: `
H4sIAAAJbogA/3TSz2qzQBAA8LtPsew9+AKrYED85EsVjFDoRUYzJtvqjuxOaPL2pZqmG2KPw/z7zbJK
G0bbQ4dxIISi9h07Ft0AzkUy44/c9LQFK4U+RFKbnlqw8rtUCDVZmtDyVRgYMZIjOgdH3PB1Qhln9f/m
Jd3vkyxtXpOqyItMhT8dcbBM6E56OIiZYGDYzGEkOzKMhhuwCLdda7QtXRZWS5d72TPrRCMd0SCdnYx7
GBz6jr+6yGo0DKzJLLeUVZ4WdVLnZdH8K6v8rSzqZLc6aj7jN16z76DFYdHfnk36Dc+eTwuTjNmeV/VC
qHDZ4THCB8djwT0Z+BkVer/hKwAA//+UzHMGGgIAAA==
`,
	},

	"/definitions/CaptureInitialMasterPassword.xml": {
		local:   "definitions/CaptureInitialMasterPassword.xml",
		size:    2728,
		modtime: 1489449600,
		compressed: `
H4sIAAAJbogA/+xWX2/aMBB/36ew/DxG272GSB1DHWoLFUOdtJfocC7hVmNH9gHl209J2PiTEFImTZq0
Ryf3s+73xz4HZBhdAgrDd0IEdvYDFQulwfuevOOXzwTaplJQ3JN9yHjpcGiICfQjeEb3BN6vrYtlDhci
yJzN0PFGGFhgT67JxHbdyawnJmtkeDe9j74NR9HT+GvUH4ymg0nQ/YWp34KJNUrBDozXwDDT2JMb9DLs
W5NQunQoFkUvIts2c27LGBNYau6sKea5DD9eXVUQnlIDeluvtPUoxRxMrNH1pDWRAqNQR2WZFN0tTM1J
x6LQ1IDuFMueXM3s61agOo0/2ddS4Of9umrbC3ApGRleV9qtq/YZKDJp2/K5XdgUDdqll2EC2mMblHWE
hmHn7XgyHIymt9PheBQ9DybTYf/2oXajQprduk6XB5ihLpXJnX1E7yFFKfZR1ZZ0iarJy5NG8Cgwd0fw
vBIbkdjyu9oGq+AlEtL4XvDaCqYF+g91dGqz7yCTIbtlrZJCBN2S8cG3DNQLmbR5Z3zNwMQNNtWBEtK6
qZ0qYndqaxOUM6i0G3SPjG1j9MCw2+yM3rtPTqYVfCexKs/qWxityNOMNPHmnHgH5x8U0wr46ArwsMLj
C+Cfc/b6rzt7c8ba/x4de3RzqUeHHA9+nphUuYrWROAQmgbWktma32NrViyj5uHVOCm+jCfD7+PR9PJZ
UbZU9lPO5jMxOz0m+gW8XeqUJvWCcZu3wanYXXSu9unmGb+Y7Pj+D4i+5XS1jebuR9Dde5v+DAAA//88
MtEZqAoAAA==
`,
	},

	"/definitions/CertificateDialog.xml": {
		local:   "definitions/CertificateDialog.xml",
		size:    21609,
		modtime: 1489449600,
		compressed: `
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
ZdTyaDDf78I8tepSc9F/AgAA//+nIUHEaVQAAA==
`,
	},

	"/definitions/ChooseKeyToImport.xml": {
		local:   "definitions/ChooseKeyToImport.xml",
		size:    2051,
		modtime: 1489449600,
		compressed: `
H4sIAAAJbogA/6xVwW7bMAy97ysI3dNsdztA2w1dsKIBimCHXQJaZhzNsmhITJP8/RDLbZpEcVusN0sg
nx7f05Mz44T8EjVNvgBkXPwlLaAthpCrO6m/G7RcKTBlrsr4vS8EyFrPLXnZgcOGciVGLCkQjy5YFCws
5WpHQU1uV8yBAKGmHQiDaVr2ko2fAXq8YCqHtkfTlgMpWKErLflcsVtodJrsIpYpGKfaSrIkNKIncvJW
d2zXK2NL6ERwaEfdMldPBW/7OVOi3PA2KvL7dd25Jg36yjg1+fb1dNpUdWhRG1ddLO+4HdYpYvdYkI3U
GgoBK1KvG87PtLEh4dp8RbA0lmDHa9DRQs1O0LgADXsCWaEDdgStN08otPf3Cnq3l54bkBWBNUGgIMub
brm/A3vEDa9tCdbUdLgSV6mxU6w3Hls1Eb+mdEc2jsoc7bWoa+OqYWTatuhKNVmiDRfAz5uWxtphOmdn
Z+PezQ/Ze8tNwTe8ndNWoss17YIaHv3SUWeXHrUYdgv0hGqYx81ahN1LCorjFKQkYm/ICe4PUJO7+a/F
7HH642F+PZ/OHhY/Z4/TP7OH+fX9RcXPxBki1rMSN4qpV6eNH8nBbQdxiRgkXi6jayrf83YNuXbk3GdI
ESP2X1JMkw/3EJZGNyppiWsrQwF5t45xiE/QcSAmqYLjDGUxKaONKSuS8NJytA2eQssuRBW6W3i4kdn4
qPZNBK4V9Drmaq9jBHv+kybATjY7koeZsvGr3/6/AAAA///yWqc8AwgAAA==
`,
	},

	"/definitions/ConfirmAccountRemoval.xml": {
		local:   "definitions/ConfirmAccountRemoval.xml",
		size:    520,
		modtime: 1489449600,
		compressed: `
H4sIAAAJbogA/4SRQWrzMBBG9znFoH3IBWz/5E9NCCV2GzuUrsxEnqRqZY2RxnV9+xLUQAsB77TQew/m
S4wT8mfUlC0AEj69kxbQFkNI1VY+9hQCXujBoOWLAtOm6kAdf9Jaax6cqCsGkPSee/IygcOOUiVGLCkQ
jy5YFDxZStVEQWUbdmfjO8DIg7/a0Carm+G+MJBm16KflkJfcte89gQTDxCGn8eITkA4JgjkzYRb9t9c
bzSu5XHZczBi2KlsWz82L7uieSqrZpMXdX6YU3TcolWZ+IFmv8YrL2XqKab2eVWtt3nzfMyrelcWc4bT
IMIuRPj/sa7Lompe86opyr9osoobZ4tk9Wv77wAAAP//NvgUiAgCAAA=
`,
	},

	"/definitions/ConnectingAccountInfo.xml": {
		local:   "definitions/ConnectingAccountInfo.xml",
		size:    873,
		modtime: 1489449600,
		compressed: `
H4sIAAAJbogA/6STQW6zMBCF9/8pLO8jLgBIRD+lVilEgVU2aCBD4taxLdtpk56+aqANTayoVZcjv5n5
3hsIuXRoeugw/kdIqNon7BzpBFgb0cw9M9mrORhK+DqiXPaqBUM/pISE2iiNxh2JhB1GdIfWwgZn7qiR
xln90DymVZVkacOKuzIMPuVjd7flYk1O6yWI2amMaKekQ+kaMAjjHh/WXB0GpFYdvmTXSFu1UxuUqPaW
xj0Ii5cYvi5lOEoHjis5+CiXLC3qpGZl0dyXS7YqizrJvaNONs61jz2HFsVAP0ZGpw3XPCgE15a/IY0X
SZGVTZrnbFGxVdqkxX8fhm/KqwFNY2f23gwICYOBdGImuHDzE3eV5lLi+MHYsbhNBp3jL3iLzXNZEHwz
nifJWVb8JootHjTI9Z/S+C6YPJ4fwmDyc70HAAD//6cDQYZpAwAA
`,
	},

	"/definitions/ConnectionFailureNotification.xml": {
		local:   "definitions/ConnectionFailureNotification.xml",
		size:    3111,
		modtime: 1489449600,
		compressed: `
H4sIAAAJbogA/9RXT3PaPhC9/z6FRncPv/RsPEOmLvWUQobQSy4aWV7MNrLkkZQS+uk7QoQQkMHpZNL2
5j/7tLtP763lFJUDs+QCsv8ISXX5HYQjQnJrh3Ts7gu11NfcUILVkKJa6pIb6kMJSVujWzBuQxRvYEgb
sJbXkLhNCzQbL76wr/nt7Wics3w+n83TwVN8HG5Xep0IqS0k5YNzWtFsyaWFQ1wAihXKimwLV1wm29sh
FVo5UI5xA3xXYayha/0Ymin14z4sUk3LBaqaZv8f1x2LXulG16BAP9hI1V0obRCU4w59r56v2bzIp4vR
ophN2efZvLibTRejSa8C4LHlqqLZwjz0yt1wU6NKrOPG9WxyBwGfJg7Y7sTzfYz+CS9Bhg3Y6YUeAiIb
ARKE46UEmi3m3/JY4hgOpMTW4k+g2c1oOp6xfDIpbm6Lu5zl0499V7nIawy0Nrx9HWJHrcF65Wh2FWWX
kHQQ6HzxrOXiHlV9gYzf6KLVFoM0O8s5yZ0OjiTQRxLXwe3BlNtr1mgDzE+bC+LYb8+nLs9FUVxirWgG
quoLMSARljRTWvVO82OXRoCfVH1RpTYVGLbGyq06qScktW4j4eUzT7dn9almjhYqOjgbtB24xzHpILJ6
arFWXO5xKO6homTFVSXBDOnhtrEQS8nxwidyiEuiaPxgCFNCG0j8mknpFD2GnvKHQqvEX9LMo7r4i7sp
ouC3sN2r1Pnsu6t3911Qw/n6BFdsqYX/0nUUeEnUH/4C172/ff7U6OpLyEV3b3t9Q2eH9bA5PQNE20CL
4RBw5hsWA/qBwMJAWKOq9odMu2lKLVH8yxOiw0mXJ8TLig9ePr9IBwe/Br8CAAD//4tvKTAnDAAA
`,
	},

	"/definitions/ConnectionInformation.xml": {
		local:   "definitions/ConnectionInformation.xml",
		size:    6434,
		modtime: 1489449600,
		compressed: `
H4sIAAAJbogA/+yZXW/aPhTG7/+fwvI9yr/brqaARLuWsTGYaFppu4mMcwhujR3ZB9Z++4mELoG88JJ0
bTXuiu3nsX3800N1cIVCMFPGofMfIa6e3AFHwiWztk17eP9JMKlDSkTQphdaKeAotOqrqTZztvqTrnSE
uJHRERh8JIrNoU1RoARK0DBlJUM2kdCmj2BpJ3UhGRvXeTJY+1kRKibXblxqC5TMmAokmDbVyo+H/GQV
Jc5axWdCBiS+k2KyFX9s0+VEP6zPWXTHc/2QXPA2uy5/pzkzoVC0c/Z/9rTp8ni3VF60Vc+IINmLbxQz
Hs9Ky3ZvoY5o5912uXaIJhpRz7dPvofQIjN4hA5UcIjK6F8tGzEuVHjA5biWi7lKhWflytzTFD/PgE1A
Ju9jwSzBJAPbwvxRZCIswP069iHr54aAoP5Ydswi5xmTIlS00/O++t1Bvzf0r73u2DvE4m5hUUwfE48v
N9de/+qHP7i8OsjEggQeX4x2vPHNZbnWdZK65sYjxu+FCveoJkyxxRAZn9FOKURFStTRPkLXKTyL6xRw
cjQ7t0wuoBY7J052VC7LydkLcFIXFJT2FowVWjURNHxhgKTfK2SZWJ/iZldNG4mbCv6eOW5Sik6R8zYi
pwYrDcDSlaF+jrxhMtRG4Gx+SpxdVW0kcUr/3f0bibOC6BQ3byNuaoDSAClXQoVgIiMUPkfqTFP7g3Jn
WR+nf57Io6Pr/UtGVwbIU4JV80KIyxCNmCwQbH4yO732nmqFrQAsp2S5Km6bftNK24hxoE6BvVPu/1ri
swasR9F6vkDUKsE1EqrFwWAtRr8LRVYmYio4Q6isw1bfVfB7CDY7r5FQT33X3IO+lif70ES+5C+ztWhz
wcZkVTuaxV9cPjPAqrrSMQV/etOT+KNf3aHWRoDCdV9+lRCjcf9y6HW9/mjofx6N+z9HQ687KCrOPt3r
LJhJS766QVtO5MVKXdqt3c3gRvffaebR0gnXyfwq8jsAAP//kCX5riIZAAA=
`,
	},

	"/definitions/ContactPopupMenu.xml": {
		local:   "definitions/ContactPopupMenu.xml",
		size:    1776,
		modtime: 1489449600,
		compressed: `
H4sIAAAJbogA/7yVy2rDMBBF9/2KQR8QQ7u1A6Uvsuim6d6MpXGjWpaMZpySvy9+pEkLgTZ2ujOW53Dv
wYNS64ViiZqWVwBpKN5JC2iHzJl6kuqZfKvAmkzp4AW19C+6bwFSvbHODM+nZldC9TBPxsrdgdEf7GcB
0iaGhqLswGNNmXJYkFMgET07FCwcZWpHrJYPxspisUiT/cQRhO2bRzciUIvdopCCDXrjKGYq+LyLkY9d
FCRf6ZMh/tgsOar2t5qR6rClOYq+9KRzew455m26pgYjSojfKzM11yqZbA65GrW9hjXRWlBanmTwlisY
BYAEYCKQDdkI3LPPVYtc7b3mEnImygfgJf4ndC58zO2lY/40U++maumo/yemDLGwZmYzjz10djVD1ku6
Ob2aNzOspmnrZuXLMMntfVs3YH0ZzrXYpcg7wG+UHQ7S5OiO+wwAAP//q8ATP/AGAAA=
`,
	},

	"/definitions/Conversation.xml": {
		local:   "definitions/Conversation.xml",
		size:    664,
		modtime: 1489449600,
		compressed: `
H4sIAAAJbogA/4zSwWrCQBAG4HP7FMveQzy0tyRgi0hoSYoN9RjW7GimjTNhdzT69sVGa2ly8Law/7fD
z2yEJODWpoLkXqmIV59Qiaoa432s5/K1RLLcaYU21hXTHpw3gkz6FFcqah234OSoyGwh1t1PPGjZY5+a
Fy/lMs3Kt/y9zPJsFoUXMf6AhbXZNRLUgJtadPI4mdxKOrRS6+ThFuHF8THoUOqgNQ5IdCJuBwPocUOm
OTMgs2pAq9qQbcDFmqmsmAgq0SocExb9kFj0A1XV2Nj+PLaFJz70K1jxQV9iw1579HgaN9plDNS85Q0Q
8M7rZG0aP1B3/wk7BBJz3W6+SGdZMS3SPCs/ZosifZ6+DmdHYV/q3Df8LXy9iMI/f/E7AAD//z34o1eY
AgAA
`,
	},

	"/definitions/ConversationPane.xml": {
		local:   "definitions/ConversationPane.xml",
		size:    7721,
		modtime: 1489449600,
		compressed: `
H4sIAAAJbogA/+xZS4/iOBC+76+wfI/Yh/YGSGmG7UHTDSPIdmv2EjlOBbxtypHt8Pj3qyT0dAMOJLxm
NDs3UKrKVV9VvvqStAVa0Anj0P2FkLaK/gVuCZfMmA69ty93akWJiDs0UiuamxDSTrVKQds1QTaHDl0I
IyIJtGt1Bu3W61W38UzN1RQQVGZoN2HSHPVQWgBaZoVC2r0PPoWj8aA/DPxgMBqGT/1xMOj5D3tB+EzI
uPztKusRMLtjuixtDphFTNNX84Yl7p9XfebAwrw8lCtcgDZFXfkV+t75hAxcTpJFICmxmqGRzLJIQoeu
wdBu2Eeu12l+eN1gmQEvwxi0FAi0GxzKo0CD2HUKHWqyaL5fXxVAbx3ZczgRFHd76rbJWKbtKBg7WnR2
Vs1aNskzIVD2DWLCZ8wejW7EFJncxGbcigWzQMmMYSxBd6jCsKgwVFaHpTElLSdQrRIpB7itCnRPRx0w
/i4w72N8BcQB4+8M7wVokaz/EjgF/e1RfyqyIUmRTqoFXgT1ssYwSS+Mu9vBYbxvuGO0a9BOGX8ROK3e
S7BKGcYVi9TlkAgpD2yxHev8fC9n8XLxfvZ7n8JJ4I+DWs7KiHJp/7pv3m5t1bYFxLHd/VWSoLIiEbxY
oR7TwM7a4JeQHfvpu0sYYKK+yg8DPNPCrr0l0yhwWlcJVDbd5TUHY9gU3nXzsT+Z+Pf98NkfDwfD+yO7
vJCJyKRX/C3kiwW04Q7sR3u2W6z3pitrFX100zfWm4e8Dw7Bx9F48M9oGLjHoHocqlF6KCkxx+kVnk3f
bkHIS83SS1P5F5URpoFYJvO7nagFaMKQZJhqZYFv1ioiyMOHXo+WK0c8XyMKG094lFmrMLzyXN94Mu+K
orYKLPfpN9bEOZtAPWG2E5Qz9GJIWCZtzay2NAaXgr9AfBspfYbEcKkIFxxHlYTLqVQTjVwaSQpngEOy
guxLCwdSN9NZtc0vJLR+u6zQmnCtpIT4WWCsluXNPxPGKr0uL11Sby1METJi2kuVFHy9gWL0MOh9Cf2/
g9GjHwx6dULNzghVQ7wFsLJPArYBucbrm3whe3MVA+3+XtcHYlEQZMPbkmfaKO2dpi1TsQJpvAikWnpS
YM7Mf9Z+RwWJ9eZMTwU28NJiOjvqdsUHrrpTff7zVn/4oRkJOEblEAmczAQpYCxw+pMJYBuQn0zwv2KC
G756aUwFDfTAyTSweUD9YWlgJ5SZsVgtN50Rzo8nzZnD/ZB/EeaYMeMlimfmZMLJkXoe+5/D59H4Q9j7
6I+bM1Cjj1bvaeCP02igwq0BDfxQPNDoBewbBu3Wu4/D/wUAAP//5I3Y9ykeAAA=
`,
	},

	"/definitions/EditProxy.xml": {
		local:   "definitions/EditProxy.xml",
		size:    6725,
		modtime: 1489449600,
		compressed: `
H4sIAAAJbogA/+xZS2/iMBC+91dEvqNs93VYhUh9sJTdClZteti9RINjErfGjuyBwr9fJaElAQdly6Na
xK0mM+Px9/nLTKYel8j0CCjzzxzHU8NHRtGhAoxpky4+XXMQKiYOj9qkE3H8pdVsTjJbx/FSrVKmce5I
GLM2QY6CEQc1SCMAYShYm8yZIX7m6eSunvviZI8xVhEI4gd6wtZMDY8liIUhFcow4iQgI8F0mygZUpCU
ibAwI467cKMJF5GTn1OCaOXLNpkO1WxxDtu5L9WsOHTZzJIu6JhL4p9/WM12ufVybduoq3lU7BRnf5WN
17ebcsOHghEfLfhszrE1VIhqXJeqzVGr55ZJgXIZE//8Y1M3qsRkLJeeX2sd1+CxQ3QLQyYKjHCesmK5
6raehyjcLNcxv4lOFutbXW62gI8Tg3w0J343+Bn+eLgPet9/h3e97k3wL1ESEDyWRZCL2163H3b61/UB
PLfAY+33FOgTl3EDHNgIW4AINCF+Lfk2T1RpE0fPtebiuRZ+m3J+pcZDdalmAZthQX2qFSqqRCvjzUY/
RzY2lt/dugfV90kCMmZR9Y3ysmeY7Rm+mphnSFMWtYlUxN0zYef/CWElkU4M09uL9MEwfZJnU7Y3XJN9
sN2RqOdLtm1EV9QFFPkUcKVgG5iy1XK9N2jfLKQDQ1sSUgrG7KDagTHPSkcnMTVlvLbV2buY0gVXDQjP
e0EuOM6JPwJhavtB56jUeGBuSmo0TE93Udju8zgnNTZl/NO7qbFg/KiL24HBLRc3pXEHxU1pPEmpKduf
36+wKY1HLaQDQ1vpEjHZRZeIyUlITdn+8o4dIiZHLaSdQLt+mFejsyr4a4PhDDslQ9AMKjhbZsQTRCWX
k+J8GVYHxrbTKs2ZRMi2KS7+4K7X6QcXQW/QD28Gd70/g35wcbvt2LRIr8itGI1v9YK4ykM0/8iggtOn
1YFadUTf8A5uJZ0yCpkAtmvdYbr5M2t1FA6yFbERTARuntk3hK+s4LeDt0EcNoPSw+UDzy394+hvAAAA
//+qBvpkRRoAAA==
`,
	},

	"/definitions/Feedback.xml": {
		local:   "definitions/Feedback.xml",
		size:    707,
		modtime: 1489449600,
		compressed: `
H4sIAAAJbogA/4xSTWvbQBC951cMKvRUW/dWNqSuakxrO1QqOYrR6tnaer0jdkdW9O+LEIUWUpzbHN7n
8DLrFeHEBusHokzqXzBKxnGMq2Srlz1i5DO+WHZyTsg2q6SZ7wlPlHVBOgQdyfMVq0StOiSkgX10rFw7
rJIRMVk/gwbpXUPOXkAqFGBgb6BR+kAnoKnZXLL0j+Dr+hFGfMNhXChe9FWjJweO+EAK56iP1MpAVslG
Oov1ZzpJmDynCH0EbWTc7Zfv3708fipbGyectiDxbqSBRxpAhj2ZAFYQUw1VBFIRR9qy0mCdowvQzU2M
+BtCZLXiI3XB3lixvNfrv21a1S5+TFMj49Je07f+abC+kWHRSbRTkGS9Lb9Vz7tD9XQsqk1+KPMf9ySu
0rBL1hp63IXOK1no2GG22udF8bjNq93h6/Eeu+5VxceZ+PlnWR4PRbX5fizyf5lZOs9z/ZClf832dwAA
AP//h3toRMMCAAA=
`,
	},

	"/definitions/FeedbackInfo.xml": {
		local:   "definitions/FeedbackInfo.xml",
		size:    1128,
		modtime: 1489449600,
		compressed: `
H4sIAAAJbogA/7SUTY6cMBCF9zmF5T3iAkBEK4RB6QGJYZUNKkwBzrhtZBtluH3ET6Yb4URRlOywq571
Xn0lAi4t6g4YRh8ICVTzDZklTIAxIU3tayY7dQFNCW9D2iG2DbD1ki79hASjViNqOxMJNwzpDY2BHj07
j0ijtPpSPycvL3Ga1EX1lJSB/7PfLTeD+u4xoQx6zWStkjSyesKTzPBegthFGs2opEFKBpCtQB3S7aN8
L/i7jg1ctGTNLEF46zGkTEmL0tagEfZcrllc1Ns2h0a9vbedMwzqpnqUqCZDow6EOdl3qZTmKC1YvmRe
5laUWZJXcZUVef1UlNnXIq/iq/OpNcb97PJ+hQbFkeLzhoo+Cs++xCa0GqQRYKERGNIZDY1ijWRWExH8
lcuecPvRZY6QwN/cPBj2D46PDYfiL4ABW+b057zWVap/j+0/A7hs63wgsN/9LYDP+zPusTsWEwTv93Dx
NUvzOsk//Wtm90LgP/xbfgQAAP//RkKAqGgEAAA=
`,
	},

	"/definitions/FirstAccountDialog.xml": {
		local:   "definitions/FirstAccountDialog.xml",
		size:    2368,
		modtime: 1489449600,
		compressed: `
H4sIAAAJbogA/8SW32/aQAzH3/krrHtuYXsPkdoOdWgVTBRt2l4iczHh1sMX3TkF/vsJAoX8AFWiUt+4
xF/f1x9sK5FhIT9HTXEHIHKzf6QFtMUQ+upRXr4ZtC5TYNK+Ssvf20CAKPcuJy8bYFxSX60Mp251m7tg
xDhW8eP0R/J7OEp+jp+Th8FoOphEvYOmPYUYsaRAPHKwKDiz1FcbCip+Jily2LjCw9z4IIBau4KlkTGY
jNHu86VkSeiWXolFwQI5teT7ynGikTXZpIxW0Nur9cLYFHZEGO3t7thXrzO33hfdRujerUs8v07jmtUt
0WeGVfz1S9318e7jue2mJ5yRLe9aUgiYkToVNK+0paANqGFNW6CQOmAnsMBXAgTteG6ywlN6QHwDKwIm
SkEcBJKdqMhhZWQBjqnb6fxxBWhk0AvnAm3jzDJ3XgAZaG2CGM62oTfgKTNByAMC02r7DJyHJXKB1m6O
t1eUeyPdTrXWFoitfekxV7H4gtoVUa/kXHmWo34xnF3OTOscOVXxHG04k7wpmhtrL9lpKo4T1do4OxJ1
u1Gv1k5nWhv1NnOCnlBd7r37QsTxW6/Pdsek2vFt7p03xILHlTCeDAej6d10OB4l38eT4d/xaHr3dBZG
Yywu2at4K2dc1bXnpiR+2MWfMwKN7aKt0S+UvmexXO62lr/rytLL8Xt36W0LYrif4P3sdbtXgiktfTaY
w/65Cs3kbYl9FJyDrc/Gc1i5V+EZ1Pf21XgOtj4ATzOwFlQNOHl5fBH1Tr6Z/gcAAP//zkDdNUAJAAA=
`,
	},

	"/definitions/GlobalPreferences.xml": {
		local:   "definitions/GlobalPreferences.xml",
		size:    19914,
		modtime: 1489449600,
		compressed: `
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
`,
	},

	"/definitions/GroupDetails.xml": {
		local:   "definitions/GroupDetails.xml",
		size:    2087,
		modtime: 1489449600,
		compressed: `
H4sIAAAJbogA/9SWwW7bMAyG73sKQXcj2M62gTYLsqBFPGRGB+wSMDLjaFZFQ5Kb5O0H21ld24rXoKfd
YoKkvv8nCSSU2qHZg8D4E2Mh7X6jcEwosDbiS1d8laAo50xmEc/a33UiY2FpqETjzkzDM0b8KHVGx6Ak
K50kzeNl+rD9uVpvvyc/tvPFOl1swtnfGn8LJ51CzpwBbRU42CmM+Bktj++yjC0NVeWohThIlbFGhAYV
NJ8Rf9nR6cLpE3VPp1bR09u8MdCBnilHjVRZHu9BWRy+76siI1E76FxINqvFOr1LV8l6+7TYpKv53aO3
UUPfffvQH2GHqoXPa0OC+s1AtdG3pWOuS5LH3sbaJsuHxVg4azl6sRJEIXU+/SaeStDZhHm+or1UisfO
VO+u6Lbu8zUFI9xwNrD7PfYvtDPnof18Gu4ANtiTqHfoFk0gnHwBhzbIcA+VcjxOr5b/RxP6cvuEBiMa
nXttFektGAQ+Pb/7yjnSr/e/69+/D3zymr8lm9WvZJ3679m/VFNgFyqnAwFaoOLDwluuet60uAbm3xnP
WXxEBRUfUpA8XKcf9xGgu1OZ2s0bpI8TB0n9hP7Ghu1eBkeZ5ejsa0kvzAzakrRtBTQz7+Yfznq5/+xA
BWcXCyJeW9A2o8LbaBC0DXYnKJx1/w3+BAAA///oq0opJwgAAA==
`,
	},

	"/definitions/Importer.xml": {
		local:   "definitions/Importer.xml",
		size:    4786,
		modtime: 1489449600,
		compressed: `
H4sIAAAJbogA/9xY33PiNhB+v79Co3cfSaft9MH4huPSHNMcZAi9TPviEfZiVITWldYJ/Pcd2/ywsYy5
knY6fUNod/V9q0+rlf0Pm7ViL2CsRN3nt+9vOAMdYSx10ucZLbyf+IfgnS81gVmICIJ3jPk4/wMiYpES
1vb5Pa0epKUnQgOcybjP5TpFQ4MowkyTLSdyR8b8CFW21rYcHcaMtin0eRIthRHGiC3v/W2LOaICoffz
fq+ypN8robtZfJJCYVKlAOZZ6hhf9+hTgykY2jIt1tDnr8Wkl6KVJFHz4H72S/g8GoePk6dweDee3U39
3t7HHYIkKeCMjNBWCRJzBX2+BcuDUYGA7bPYFehVxrQMDfyZgSUefP/DTZfHEmSypKPLjzedLnM0MZiw
WIsHt037aClVzAq1aKG8YtjnL3Pc8MN2NfL+9SNuyqxX7RyAcY0JaMDM8mAhlIXT9Y8YjmOnXsUcVLmk
Kn5WzZsLlzauTRouES0wWgITu41iW8zYK2YqZkqugBGyUkxM6JhZEoZYZqVO2MLgmg1xO/ry3kXEuctG
pDwgkzmpt3l4a4yBB4+D8f0kfJ4OHsPnyfRTOPw8aOhzF6RyTo6BRbSSOjm/HGxSoeMz2+NyWkilvo3V
8cg1NLhn0IDr906UcYlSniKDSkG8qwO12lDOdWjnxRZWc2G8FJWMtmWVeJw8jIa/hYNfZ5Mvg9loeCnx
5ZXhGpTdtGcG4KuEE8KHf08DNHHmilM8cNwEbdDa4Z2HOCzKexWoR0tpvd159Mry74Dsgt1ZjvPQ5xi4
glo0FJYwQhnz4LvOAK1pcKdiCEpNQcdgwMwwSRS0p8Ps7FoS4sIvIpIvZTLOn9JKCEcBqc0LIiPnGYE9
g+NgVAUCRQIPU+0IupZo1INL4Lc6XSvd/DLwRJoqGYm8sr2Vbn/OL5lK3GvF21JvL0mEOxk18cKGWtJx
0G3vbTVF+ZI5q/+dovYnPuf5VmraVfLC7lol3f47Sqql4Z9V0e1/TkVuB4fx9R3ft7Rv1zR8LbK5vOFr
PE/yawV1KAwIfr4b/JgRoT4+V4phWH+1uNCjkaBJHB+Jk+nobjwbzEaTcfh5Mh39PhnPBg+XJmMtTCK1
R5i6HmF1tt29XkmqxigSOmq8ilxIzryMihDtR/xiZV5FBFdXkSibvrMd60msSGgvhoXIFHX1SleczROj
ukFt0i/F7b3KOAE6fnep/c0M2BS1LfEXO18Tgt+rmXcGwRVnuyT0eZ6EQzxcOWOd/HnyrcbvVT4//RUA
AP//YKiX67ISAAA=
`,
	},

	"/definitions/Main.xml": {
		local:   "definitions/Main.xml",
		size:    9171,
		modtime: 1489449600,
		compressed: `
H4sIAAAJbogA/+xaTY/iOBC976+wfO9mNCvtCSI1aHtnDr2z0qxmjpZjV4J3jB3ZDun8+5VJ+AgkIXw0
bUZzg+Cqen6vKJcpxkI5MAllEP2G0FjH/wFziElq7QT/5X48ZZkUjDqh1XehuC4wEnyCF1Ss33s7hMaZ
0RkYVyJFFzDBjCqSaJZbHD1TaWE8Wi9oX++Ek4CjmS4/vxxbyyGhuXSkENzNcfTxw4ehFnMQ6dzh6I8W
EytSReXGwDqjS4zmVHEJZoK1IkxqC6RYbZtUyzGyBc0y4BOsNB7VrthcSF69buN0ql8rFj/F+hWv1x2i
XgorYs/KvyY/IPAc0g/BHQH4rQnw3JhtdtoIUG6VWDhagnGCUdlpfAC6HfgLqHxKTZ2joPKYGrxvdv4e
uqF0w/nsYFHhmWnlKHPWP20BdSmwNntJY5AYOUOVldTRWMIEl2BxRNZoTnWZW3jIFQcjhepJzUPKkCsz
mGCbx4tuArpY3CraaXg5e6hX3n6AW5kp5z0KXxNrm59uxZ84f3x8HOy3UQ4pc2JJHTTrIeWcsCqNBlTF
jjijis4eRUZHJLlcMgXFTKtlkLL9DQXy4MDYVbG8soYKCq/hxv07KNjvoMe427DT6Pzq/cSYzlUo1XuN
5s2r9005/iagCINfj+TeTsaPAR6Nszmw/TbIP/JvX8CkEFy9XaFCQ79fG/+NUut0mkrgzUpbPyTM7594
Asgq1t2cmX1afp3r4kuSVJkfmKIeG6rBoaFt7ybMWcL6kKQO+dPo+50KJ1Qapr41uFvqW4f8OfTVxn11
1HmZQpNXG4emJargvbGu2jgyLUkV61cv3NGnfcn8VeHEVvjhjVq1Gsw9dWu6gkyO+Qiqb/tTMVNmbqZV
ItLcrG6Lz0JCc2lo1aNGjdgubJQIOTzW8EoCVTDSCEZ8sLs5JJqaZwYSMKAY2GAVJv9sQZ7y48jl9WDj
asDvLDtM/jpWOrLtE8gsiOu/B3LTq/vvAR4BTW0SAB7TqtAHVwKea3DX/MKuN0y4oFLfT5e/N4yIde6C
1OzJI7vqNMI7fD+1bldiuwzGGWU/hEqPzzvhNaOKDxl27hkmQsrjZ+O+VaatqOa8B1P33V21wm/lYOg4
eDPHVtqJpP4bwwM1QK89FL5kuB2AoiebeWDEn304AsXfMRnaqGunrYuyvmw+Pf9P3m7LVve22dzizofb
D8ajnb/x/B8AAP//ctSZa9MjAAA=
`,
	},

	"/definitions/MasterPassword.xml": {
		local:   "definitions/MasterPassword.xml",
		size:    2357,
		modtime: 1489449600,
		compressed: `
H4sIAAAJbogA/9RWT2/aThC9/z7FaO+ERL+rsZRSlKIkEFGUqr1YY3tspiy71u4YwrevbCcBgvkTeupx
7XmjN2/ePjtgI+QyTCj8DyCw8W9KBBKN3vfUncy/MmqbK+C0px7RC7kn9H5lXaoqAEBQOFuQkzUYXFBP
rdikdtUprGdha1R4N72PfgxH0dP4e9QfjKaDSdB9w7S3EBZNCsSh8RoFY009tSavwkFFFhY1DyheiZxq
l1KGpZbOilOZqfD/6+s9hOfcoH6v1yTUoSUZUTBDk2pyPWVNlKBJSEdNtYLuKzqZsU6hFtKg7tTHnlrG
9uVVozZhv9iXRtXn7bp99gt0ORsV3uyxbqv2BSZs8nPLZ3ZhczJkS6/CDLWnc1DWMRnBzXrHk+FgNL2d
Dsej6HkwmQ77tw+tjWppNuc2XR4wJt0oUy34kbzHnBRso/Yp6QbVYpknTegJqHaOzOijeyCzzfPEmozz
0tVzQcaaruCnLWHFWoOxAjEB+jm9IdhvemCObKA0whrWtgRHXtAJ9O16+HjVJkXr1XFYqFBc2boFgKDb
qLXzrMBkziY/3pleCjTpkRW3gTLW+hidfcTm0re6r5pgj27Q/WCKc0wyMOLWG5NsxdFBp6PvZDapfP6Z
iZbsOWbNsj4l3k6EYCK8RKHd+PC4pI/h8c9t9ubSze7OuPPyQIJWKloToSM8FqSliDXvcRrXx+h4qB5N
sG/jyfDXeDS9PMMaSg2f5ptxwp2H46tfw89zXaI5mVN6zjfrkO0uuo3b41Yev3jY8f1fDPqZ23WuNTcv
gu7Wj9KfAAAA//8qcNWVNQkAAA==
`,
	},

	"/definitions/NewCustomConversation.xml": {
		local:   "definitions/NewCustomConversation.xml",
		size:    4733,
		modtime: 1489449600,
		compressed: `
H4sIAAAJbogA/+xYTXPbNhC951eguKOKM5MeOhSnjqIoaj1SR2aTaS8cEFyRsEGAAywtq7++Q1FflCCJ
SuxpD72Zxr7F7r7dB0CB1Ah2zgWEbwgJTPIAAolQ3Lk+HeHjnXR4j8YCJTLtUy6EqTQ6VpgUFK0xhATC
qKrQrvkiJPiBMbK2JJoXQBjbrjW2BJcl9Gkmcm65tXxJe160TK/ABr29QIJek0v4xpvXR8mVyZqkJrAY
VA5NMTD6CazjKI3e5FZaU4LF5SqRPl1InZoFK42TjdUo+i3+Op7Ev0/v48FwEg1nQW+D8btIjE3BxguZ
Yk7Dny6Zo0QFlKDl2imOPFHQp0twNLxHbpFoWBCxF/glfxac/Lv2QsPIVnDJPIU5rxSyHGSWIw3fvX3b
FbLO8H0XhENrlmwhMWclt6CRhuiLTuRSpWTVtJortvrs06fEPNNtmxyR/cE8N0x/2bc7DiM3hclAg6kc
DedcuaP9fShjJWjku3aYzsbDSXQbjaeT+MtwFo0Ht3ddHLmSC6mzg57Yma+S3cF9mY6sTJtUs/qvfePj
/QpuM6kZmpKGN+99EZ5BJQbRFDS8OSL3AtDVXfsNONDpNShrFmxb0Jt3XWGNgrALVPjp8FNyxxNQbfVc
/+sQexyNagw9s//JmuLnU0n5XD1UDuV82XTor3/cR+NPf8az8ehzdI2XnCuZrdv89m48msTDycfTDnYi
fOiWi0epsw4VgDkyjshFTsOT5PuQaMouwKDnjSXoeejtSvnAFInZis6GdXo55vWp2j5lr+PHsbkRtX75
5PMsFJ5LXs+Y71S4XIUTlQClZqBTsGAjeMZWRVi9MbPNcrq7AbRcckQrkwrB+Zb3DTa819vUjG8XvH57
5xx72X/xdr75F9r5eyWsBLAvIF9fJeb/y1dXvs80ymvI11CjXa4nNU0tuC7SVSouIDcqBcuaEXSmAKPh
lweeJGB/NDZ7HTlyMtNcrWFcoHziCJTkXKcKbJ8aHa/uHHFjSMmR0PxX5vpFePYl40/kMJRNxc+dHIeY
uVTqOsTu6XQiXU+qx/J14iFQ0290zC1wev6m/KFCNHp7SCerz7j9PvCFf/a2/3k6G/81nUT++z65Ygab
8FqxCa7F9s19LsRGdsN4sAJc04yVA1bVx7GS+sT7sAVvjZ5QUjxC2p48oYyDayfvu9TLUznz2LlqJ9/Z
r1nGwzcI12z9hu5wjbvMwTepXydlOTBqG+wt7haC3t4PTv8EAAD//3OIjat9EgAA
`,
	},

	"/definitions/PeerDetails.xml": {
		local:   "definitions/PeerDetails.xml",
		size:    9977,
		modtime: 1489449600,
		compressed: `
H4sIAAAJbogA/+xaTXPbNhO+51fg5R3R637k0JE148iqqqkjzchsMtOLBwJWFGISYAHQsvvrOyT0RQmg
+CGnSdubKGIXu/vsswCW6HNhQC0JhcEbhPpy8RmoQTQmWl8HY/N4x7W5N1JBgDi7DmimFAiDIyWzVAe5
DEJ9KuMsEdo+7Z6ReUnhOojoiiiiFHkJepvxvQOBfs/OOXjjnP8DiGxiILHTE8bG+czbiVMlU1DmBQmS
wHXwxDVfxBAMQpVBv7d96x4ckwXEATKKCB0TQxYxXAcvoIPBFNaomObt27cnWjSPBIk3Ogg1/IkYCNCK
CBaDug6kwIQxLGBtgxSgXk03rYs2sjjJ/2js5ZlZbjmJZWTnYfa3ewrDTQzO2IwYN4hKYQg11bGhsdRH
gaFE0Dzm2zygKx4zVGSgIDEuHq+Dp4V8Dna5dOLEe/lsPfh4OO7Uh4SoiItgcPX/Qzv3w4vZ9uKuqcaK
sy0qnAWHg33zYSPTYHD143FszkgtpDEyOTa2hqA2RJkWciBYEykl11inhHIRBYOr7+qKWZ7vJd85oXDD
4YbkzrK2KAaUykwYbIl8LNqE8DdWk88rl7IViXkkgsE4/PXh5m4ynj6MprdNFGiIgRZGBINw/tvIL7sn
9bFGQh+5iGp4DkuDiTGEroKBF3OXpJFpHcF+z2lLv+dAtQPSuU01gD7F5j68mYffBDpXfwM6XeH5zNkF
SDh0LyrNgP7Hk7AiP16ZhJ9PVsB6kPwLuNcBlK6oCE4fc4MuQMDpRtUrMvBrYZF3A/M6LBoJo17KeNVA
anu40JjBkmSxcZ9rXi26renQIbpd6aDgj4wrwCCoekkNl+ICxJhbpWivFK25WSGz4hqlAOq/NasyH77/
smwbroA+vs+MkcKXFNtOROeYtOZIh5h0JkmpeVPpqpcSRXtEd8v7b2dj0Drzf2giuObMXHBt+h/GKJzd
zn5CYV6n9EpmMUMLQFwggjRVMo6BoTUXTK5tQUu46CXkGWn+JyCEcTv+hQrgI4d1qZ/1lP9xPgaJZHmx
LrcZmwSxsgVYJUiJwEtJM91cdAWEgcp93Mz9M4l1Iw16JdcYnlMickUtFCiQioEi9Vx3990svfLqeCri
w/l+J1OAXa3CDbctCfeju9EwnMymD/eT6fiugtsbTeVe44qICFi521gkD96ZhHeD9JqkKeSbMXm6EKCK
8uGhGvKSwx+4nCDDojt2QJOiwYFt06xmBG2ndrAXPxs4r6FuY4cQx3PI8xJUCM/GYy5WdgxzR7TQTYxR
fJEZ0L4hh4O2DuZz5hV498Krv3duAi96LSAvtdpLL77wGuNtNbsEu68xrZYEuxvbde43a8Ki+LfORkQq
DsIQW1nycjGbT0bT8KYoGL/M5pPfZ9Pw5u580atH0A8gssMNJGEML0xdSu6WAXOmDruEU5lm6YbP9iuQ
t2XuU8G42hbhYrM1n88+PdzOPk0vXRomCYmgFCHMk8gTJZelnEqxaerGXBtMGDtnYxVT0YUp3ixpyieO
RD5BVc4cfSrj9PF4+drosB8RPUW1A2J7EzuCZhV9Tbh9NaX53SW28qfun9Ziz46OFGXggSggJYDPFmhb
mR/KH1hd3nYszM3Wj41tZvchucsJdlioqET2PEc3dpx2E5wp26mrcRwCTZ7qNBH9Abgn1ax1HZB2jchz
S1uN2BUOtI9cDVpc9Ah8jFPFxYRj2hxHclGc1fBmL+b7CH9yypSJjECAzI+oS/8BsRFFP47m4WToI6gr
Nys7S0suIlCp4sLoiVhKlRDnabCqs1OVWq7scOVGLbv3tzsOzR6f3vRoe9ejw22P9vc92t74aH3no+Gt
j8ZwVhA9f9iyvTzsYMj+Rb93cNntrwAAAP//EEdQ2fkmAAA=
`,
	},

	"/definitions/RegistrationForm.xml": {
		local:   "definitions/RegistrationForm.xml",
		size:    2189,
		modtime: 1489449600,
		compressed: `
H4sIAAAJbogA/8RWQY/aPBC9f7/C8j3Ktz305ETabSuKugIJoR56QRNnCO4aO7KHZfn3FXGWEGK6S1qp
t2SY53nv+Q0glCF0a5CY/8eYsMVPlMSkBu8zPqGnzwq0rThTZcbL8HxsZEx4VRnQzMAWM+7Q19Z45GwD
ptToulLSVjhLY0ipbQ/WvPcwASQ3SpesYWtAJ81rxp8L+9ISirF/sC+B+vfzPsZE7WyNjg4tiWflVaGR
5+R2KNLXT68DtuAqZZLCEtktz+/+fz/osrtrbzR18JigiVNlUFQdn86bR6garSwGdHaf+BqkMhXP7z68
Fyat3m1Nh/x4FTiwJ27RIxSog0fKeHI7Scoazy+RI/1iTKRh5qBeg3xSpnp7kMY1JUAEcsPzqwbHkGTr
ccC9KmnD86vXcpQV5S/SYSwHBpyaLsI82FdoLmMFDoH/PukPOyJrTgtc9Bd49PVdgqxTaAiOrHg+WX5b
zRfTL7Pl/XI6n62+zhfTH/PZ8v7xT0MZ1LRSyCQSjET99zIZyVhYA3JgvAaCQmPGD+h5/qmZfWu+IzkY
L99hpTyh+zcGLNrpt5wmwSQlrmGnadz3w217FGvoL5kIq5TsVVkh+ROkV2avP8KNgCZxXfpE2ut98wSo
a33grHUh40cXwnnuZGjkxItiw7QTJtKzPyC/AgAA///6Pm5BjQgAAA==
`,
	},

	"/definitions/Roster.xml": {
		local:   "definitions/Roster.xml",
		size:    2558,
		modtime: 1489449600,
		compressed: `
H4sIAAAJbogA/7xWTY/aMBC991e4vltRtf24hEgrihDqdlkt6VY9IccewLvGjuwJLP++SgLZfIE2SO0x
Hr8343lvBkJlENyKC4g+EBLa5BkEEqG59yM6xZfYASzQOqBEyRF11iM4trUSNM0RhITC6mxrfPlFSPiR
MfKsJGGsOiqvEDykMKJrseGOO8cPNGiApPKp5gdi+BaGo7kQNjNIrsnskWPmibDauuHohIuXtbOZkdcy
7EGtN0jsiqyswfN4ZbCFRGs1qnR4SiWsOYuaypcH9ZpkqxbI2X1xY0C6MKj5IwxKg/V7bSGc1Rrkb2Wk
3dcNd7Ja6mwKDkuLjOjGF4iEO5ZarcSBRtP4x/Jhfjcb/1neT54mj2FwwvRT7C5S3P6K5z9v49m4QyM2
SsuqB71j86Sg8Qi2yw9OmG4px6Gqj1g7bW8TgEtwnu2UV4kGGq249vAepN/YPYPXlJucYABQww40U0aC
QY7KGhrd9OJyz8Tz73PiYGt3QLZ8rQQx2TaB+pR0MxxtzUrv0OhrL71Xa8P1EeLsnnGBascRJCUbbqQG
N6LWLE/HyyST8kBJUOMohCTFFjRcs+JzRD1oEOXLqqtn12N1txC7H9ovd2m2xeRuMo5n8/vlYnY/vZv0
vbU5O6eThgvbrrzszHHZ2aLkfBdUrW4W3WHsZx2D1o+Q+wjccXO0mZkr45IGHT6O6FSSIfh2qB48Ni4t
+aNvYVCF3gEToDV7W9U0+nwBXwu1Suo0/Z8qk5d+Rpn2xCjMhz//6PfPdWLG8IqdUq6W8tI+8I2FcFZG
zAuKPg3TfmUdnGS/GQZ9r2N6oOWvOo2+/DejNS/Ugm+BMKj94/sbAAD//y8D5/r+CQAA
`,
	},

	"/definitions/SimpleNotification.xml": {
		local:   "definitions/SimpleNotification.xml",
		size:    312,
		modtime: 1489449600,
		compressed: `
H4sIAAAJbogA/4TQz2oDIRAG8HueQryHfQFX6B+7hBAt1dKjmHUabDeO6ISQty9FCi0U9jaH7zczfCJl
gvoeZpAbxgQeP2AmNi+htZFP9HmA1sIJHlNY8MRZiiOPff7OMyZKxQKVbiyHM4z8mnLE67ZgS5Qwczm5
vX/baf9srH9Q2qkXMfyY/1ecMYaFS6oXWI3277Z0K9BPHZS1d5PyO/1k1vTxQoS5dXj/6pzR1pv9XyaG
3onciOFXV18BAAD//5eiD8A4AQAA
`,
	},

	"/definitions/Test.xml": {
		local:   "definitions/Test.xml",
		size:    261,
		modtime: 1489449600,
		compressed: `
H4sIAAAJbogA/3yPPW7DMAxGd5+C4F7ISzvJGrr0Bs0sS3TExBENmbGd2weB8zcY2QjwPeB9leWsVDof
yFUAVtoDBYXQ+3Fs8E+PO85RZgSODQbJE5XRK0vGGw5ghyIDFb1A9idqMFLnz71+JeJ9UnQ/dW3Ng/ms
zBw1ofveMELiPq73VuP/ryxr4dTKgs6aFbnb5qm/Hta87b4GAAD//0zV260FAQAA
`,
	},

	"/definitions/TorNotRunningNotification.xml": {
		local:   "definitions/TorNotRunningNotification.xml",
		size:    715,
		modtime: 1489449600,
		compressed: `
H4sIAAAJbogA/2ySwY6jMAyG7/sUUe6IFwAkqkUsWgYqyqkXZKgpmUljlKRqO08/amFaKnKL5f+3/s9O
IJRF3UOH0R/GAmo/sbOsk2BMyFP7lameNqA5E4eQC9VTC5rfpYwFo6YRtb0xBScM+QmNgSN69jYij9L6
f/OR7HZxmjRJVZVV4P/q3XYz0MXrJBn02rO1pHhk9RlXNiOOCuRs0mhGUgY5G0AdJOqQT4/q2fBnXzcI
eWAPXAXSe5Qh70hZVLYBjTBzudawoeu0gpauT9maYaATHVEhnQ2PepBmFd/lIi1QWbDiznzfW1llSVHH
dVYWzb+yyvZlUce5c9QD41W7sufQopzSzyfiS8M6D0opRiO+kUfbuEjLJsnzbLvL9kmTFH9dMVxTLhpG
9wlnvT8lXcD4bzTvgkXz1Qj8xff9CQAA//+w2rQNywIAAA==
`,
	},

	"/definitions/UnifiedLayout.xml": {
		local:   "definitions/UnifiedLayout.xml",
		size:    8750,
		modtime: 1489449600,
		compressed: `
H4sIAAAJbogA/+xaTY/bNhO+v79CL++Cu0mL9mDrsAt0GyAtiu6iV2FEjSXGNCmQlD/y6wuKlizblGRZ
TjYBclksJc5oPp4h5yE9Z8KgWgLF6H9BMJfJJ6QmoBy0XpBns/rItHkxUiEJWLognGmjq6GdHgRzKnm5
FtqNgmD+/zAM3LNQwBqDAjKMmUhxF4RhM8vNCMy+wAXJmDBk1qmh+tMpS3NQoBTsezR8Yuk0BVRyqaap
SICuMiVLMdGULbIsN7fG0kjJDSummaANmFJ36nhOV3+zXVIue1SUIkXFmehJbMuR+ayFsvnMYdQP1z8Q
UlSPoBxc82qYgKrhWihZoDL7ClULsmGaJRxJ9KpKnM/qt/7Jhhk79UnuP/w5NFfnchtTLjXGSWmMFL5P
9LvyKHfOCWUzfg8HKIh4KWmpSfQ7cD04X8i48gM47/gAzRlPmwReuPCCBSgw8pANXQ8fSC0z0qNbvPLJ
SMVQGDDMJmaDyjAK/FKwnSCnBuiKiaxbMe4KEOkIS5asK7i+2YXUzBn9k8/aE/Pms1Z2hjLVgC2RuzfP
Tl5H8dqPXJnO8zAMhMItHnEid6QtcUNQbg2M32S/2R8hQX5iOK+enEveaP4UFypZMEaxpDSoL1+2Xx90
u32OBBvgJVpY8vS4n7TkZt16z6v36MZ5FXc5OITDLrn+ou6S6i3uo1Ne60+KvXl4JXQe3Q5VYed0z/oq
2BkrqpAi26COU1xCyc0tGjjDJYmEFP1ymmUCeG0yZ3SFKQlyEClHtSBSxM1DvYWiwNTuml6genPhz8eH
NWSHfpvZfx88iZiUjIGEDBazT55RKWL7L4m2TKRyG1ZQCvV+nUjOaL++rlrtQPY9i3t46bpPdQNdxba3
JRGK9LZ14eEe64Ivcv6o3RSx8dEavQ56fL3w85qd/i9pMJFy5UpN1KMvu9mPEdNUSc5h9Oeqrt1AMthb
+OQSqVJUIyXXoDImYm1AGRL9MlIMLaq6hU5WYb1lhuZhUS2R7ZXYvYjdi87VeDr6xyRiCvg7in0Y/N8w
a/H4NIW1dPHLd2/NYL5/fvnuy/HL999Jdq7bRFq8qyJc583a1yWLZ3KOAgZGgdAcjN1JFmSPmkRPUmxQ
6SoIulNdD7m6gbB1kbUfHUkNpleF+C/DrcOTUYgbO/p2OpLctSRxYVv6PYkEmFL5i8cnvpYpchI1twhX
f7Y60dBx4+qohKOwwI81gqL5SFmghm3AYGw7DCYyjo7yDYTtpGVRchvWes7oY/24lz865ATVTY0AHlZD
u91xpA6WgxTfAuulmX+CrrDR83BZsBNOFWosP1Vn9wc6S686VXAn7b5lyyro5U4jaPYTcv4PihQVqsON
RWNkeLitUO699+xp8FDLs0pS5Dw83gWR6OfWoniljsLZGv06INp3NjaSV98dCDX4XJh9BxwX3KQAykRG
ot/GcOc1E/GWpSYn0cO7nmO1MRh0No/RdfU54iGicP5iCL6vuDMusCVLz7Cr7oXdibBdSoW1+Pvx4sa6
aAnEWMFDY2L55VjR5pqyQt0bFtsP1vqWXOj9GC7Uf6n7wj7js5JlcWCs7HNVE0XH9a7tl7opynzL0gzN
8ccXbnzSM1XXVzP/jANfaW7ZW+qOXsxnrV+I/BcAAP//PXkyFS4iAAA=
`,
	},

	"/definitions/VerifyFingerprint.xml": {
		local:   "definitions/VerifyFingerprint.xml",
		size:    1378,
		modtime: 1489449600,
		compressed: `
H4sIAAAJbogA/6xUwY7aMBC99yusuaO09yTSLkXbqKukotFW6iWaJAO4NR5km2X5+wqHFhzMlqp7yzjz
3rz3xkkqtSOzwI7yd0Kk3P6gzolOobUZPLifHyUqXoKQfQb98HxoFCLdGN6QcXuhcU0Z7KTueTfZsJVO
sob8of7cfCvK5kv1tZnOyno2T5PfmCNFt5KqF16BRjXxZQbPLb8ch8QU3fPLIOe87VJPy6YnM9nJ3q0g
//B+PDyGWfGal6SJtxbyBSpLt6DYSNIOT66reTEr67u6qMrmaTavi+ndY5TIGz7VMbeP2JIa/K7JWlwS
JGcMyYjiSqTYHdQ1aAjh9Xn3W+dYnzL2ZRNG/c8hfKrmxfeqrOMxxKN4TV6grUPdkYIx9lJih3rS0wK3
ykHuzDa63GtgNezBGdRWocNWUQZ7spBP/fzrXGkyuBgZHm/uP1N4JiMX+xtSuG7kyVO8gZHLxlFT2BC8
TIerevhul+TsH0RwLAzZDWtLGWgGcVxqBoelQh7cizQJkH/l80kEmUYZRode58lVmpz9VX8FAAD//3BQ
WgRiBQAA
`,
	},

	"/definitions/VerifyFingerprintUnknown.xml": {
		local:   "definitions/VerifyFingerprintUnknown.xml",
		size:    1102,
		modtime: 1489449600,
		compressed: `
H4sIAAAJbogA/5RUTW+bQBC991es5m7R3hekJLVSlAgqF7VSL2iAsb31esfaXdfJv68wTszyYTU3j5n3
Me8hpDKe7BprSj4JIbn6Q7UXtUbnYnj0u68KNW9AqCaGpvvdLgohD5YPZP2rMLinGE7KNHxaHNgpr9hA
8lg8lb/SrPye/ygfllmxXMnoDXOhqLdKN+LswKBenMcY/lb8chGZcnTPL52d/trYT8W2Ibs4qcZvIfny
eSg+hdnynjdkiI8OkjVqR/+DYqvIeLxena/SZVbcFWmelT+XqyJ9uHueJDoffJ2nrn3GinR3756cww2B
iHoU0YBjJlOsW3slWkK4LXh/9J7NNeTzWIZZfziFb/kq/Z1nRZhDyDcK45a/wBzvYIgb+6vRLBpa41F7
SLw9TlY7B9ZdC96icRo9VppieCUHSf40zyOjzv0guGFl/cW5YkOq4KHsum3f9A15944I/haW3IGNa4va
gbgEEUMbBCTvOcooQL2pjRWufmTU+4L8CwAA///rojW2TgQAAA==
`,
	},

	"/definitions/VerifyIdentityNotification.xml": {
		local:   "definitions/VerifyIdentityNotification.xml",
		size:    965,
		modtime: 1489449600,
		compressed: `
H4sIAAAJbogA/7ST32qzQBDF77+nWPY++AIqGAh+0lQhkRZ6I6MZk203O7I7tvHtSzRtYmP/3ORymTkz
vzmH9ZVhtDVUGP4TwqfyGSsWlQbnAhnzS2JqmoOVQm0CqUxNJVh5bBXCbyw1aLkTBvYYyD06B1uccdeg
DOP8rrhfrNdRvCgeo1WapLHvfShOA6qd0hvRExjQs/4ZyIoMo+ECLMJp1RTZnA4DVUmHz7Zrqh3taYsG
qXUyrEE7/IoxpSKr0DCwIjOckq2SRZpHeZKlxf9slTxlaR4tJ0f1Z5zfU+xLKFEP9CfX5KXgmufNQiND
tu0kvRC+N+y4wPBGHOOGUfGbGKA6Xv/3FFpmMsXPYdzY1nnPMOJ5Ravq7hd39RAHWzBOA0OpMZAdOhk+
9PJpz6/nVGBmG6yh1XybsM4F37v4uO8BAAD//wE+QbLFAwAA
`,
	},

	"/definitions/XMLConsole.xml": {
		local:   "definitions/XMLConsole.xml",
		size:    3090,
		modtime: 1489449600,
		compressed: `
H4sIAAAJbogA/7RWXXOrNhB9z6/Q6J067fR2Oh3gjuN67vX0xs44NMn0RSOkBdTIEiMJY/fXdwAn/gBs
k/Y+Ant2j/bsHuF/3qwkWoOxQqsA//jDLUagmOZCpQEuXOL9ij+HN75QDkxCGYQ3CPk6/huYQ0xSawP8
xb1GsHF3RZKAwUjwADOtrJYw0cqBcrgCIeTnRudg3BYpuoIAO9g4jJyhykrqaCwhwFuwOPwKUmpUaiO5
P3oD1YVHTeXwppPF74JKnTYMXu6/TRoSPdVLobguvVxb4YRWOPwS/UGeZ3PysHgkk+k8mi6Pi7dTxNpw
MKQU3GU4/OVSuBNOQueBX+4fHtCO7m9oPJks/pxHZD6+n17KacCKf6pMOIxMAZfCOSS0kM7LQKSZw+HP
t7fXQnan/HQNwjqjt14pXObl1NQT4LrYsUxIjurZUlR69WOA17He7ETrUvlObxqJnw7j2jQyvdIpKNCF
xWFCpW3V70JpI0A5uh+JxXI2nUfjaLaYk6fpMppNxt+uSWRzyoRKT+ZiH14fdg/vOukjM1pK4M/1qDaH
XoG1NAVid9+aMcaHidpc1sKK3hnpAzGqSKJZ1b0hMKWJzXRJqJTDgBe2qbfPGeW6JG6bAw7BsQy4J9S1
6JVQhDUuRd6W4lMvuKVZt26VGz4JOFGssrt19fY0wRm9urbmHLA0NPdWmgMOfxqCAy5c4yK9i9I7JoWx
2njvlAcnyMUGpPVikLr0pFCVIfZK0JVAQuK8FTWpUAORphL8Y9C4ue3C45uuP8P+6jp+27aBjkA/p+xV
qPT8LMMmp4oPW7pEnF9Tf9Qq/c75xMtaPk5Z5aKEGqD4vNHdFc5p9W7scf1Iju29i/1Zs/66WM7+Wsyj
brtGA/a5oXfEzUBiwGZX7LKkMUgckmWDGDJjhQWvUBxMtRXnda3hVqSKyre9lIK9Ascoo4pLMAHe0yZN
JEajD4/pf+oek9rC1b3r+mMikyrF92xm+zb0dr9CV9jyZSXqFvwPOrQDT4KOAw4+7j/4o4Pf+38DAAD/
/1EB1DYSDAAA
`,
	},

	"/": {
		isDir: true,
		local: "",
	},

	"/definitions": {
		isDir: true,
		local: "definitions",
	},
}
