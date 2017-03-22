package gui

import (
	"bytes"
	"compress/zlib"
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
		b64 := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(f.compressed))
		gr, err := zlib.NewReader(b64)
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
`,
	},

	"/definitions/AccountRegistration.xml": {
		local:   "definitions/AccountRegistration.xml",
		size:    2956,
		modtime: 1489449600,
		compressed: `
eJzEVs1u2zwQvH9PQfCuz00PPckCkrR13AZxkagF2ouxptYSE5oUyJV/3r6w6MSWRSlO0qI3SdgZ7s4M
RcZSE9o5CEz+Yyw2s3sUxIQC54Z8RA8fJSiTcyazIc/887aQsbi0pkRLG6ZhgUNOkhRyRha0U0AwUzjk
G3Q8uSyMcciAObRLtIwMs5hLR2jZxlSWgRCm0hQPHil3KziZa1A7fqGMQ84K0JlCO+RGTwVogWrqyzgb
hGAZKiSMcImankN7uCikylgtiwYV1a9DvpyZ9W7ykEwXZu01+nFY11ZpATaXmidn746nDVW7EoTUeWd5
3dv+PdTYNcxQ+dYW6BzkyA8B7TWVBwR8/KYQHDLh7aQCHw1dFWhxayVbgaaGv6Af3f0/NEGogZWFkidk
Kwwj4oEfsvGtBPEgdd7PjOsSdMaTOSjXQd4GzaVS/e201o4HO2Ne5NTIyswblW+f+rvyOYpmhsgsuvIR
Alqzivaxen8qTBhVLfQe+eEQ2ES2Rn02mD5GkU/eMfIlCb2ribpmCnHdV47kfMOTUfp1+uX7XTr+/HN6
Ox5dpd0soQCyzhAGJ8A5RUAEouBJp3UhJJnyFGAglOwwmK8x7NIsZubCrFNc06FvJzhWgItQk93wJO3c
SH9B2rN/Km17mK4fQ+u0AUHS6ClYBN7/17ioiIx+On5mzeMnNKexEjXBdgEf+8nt+NNNep6OJzfTq8nt
+NfkJj2/7vwznBgX39iuK9KRP27ftL0va4peT4+uDFI8YHbKpaHPtQ57Xy8FlKXavEmJ8y3DS8ItQEcZ
zqFS1HeYnayigyX+AQ17tkiooLl/Yr9LopXMciT3BGl8ZhZdabTzGtQJ3KcxHjRqn2XwxrGdkEO+FdLz
gTckQHf0sW5zP1U8OLh//w4AAP//pCBesQ==
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
eJzkWN9z4jYQfu9fodG7h6Yz7RN4hssRjjaFDPG107541vJidAiJkdYQ/vsOGIoTyz4TJ/0x9xbH/j6t
vt1P2qUvNaFdgMDwO8b6JvmCgphQ4NyAj2n1UYIyGWcyHfChW90Z+wDO7YxN+QHAWH9jzQYt7ZmGNQ74
TurU7IKNcZKk0TwcR7/Ev0+m8cPsMb4dTaPRvN87Y/wUJEkhZ2RBOwUEicIB36Pj4egQLNub3LLNKYwK
mZOZBnWiSlEhYYBb1MTZEnSq0A640bEALVDFZ5q4gHHWO9GIpVQpO6qjQQXHxwHfJubptHGfWh/MUyHV
b+Xvqjtcg82k5uHN9y/Dv6x9efatdA8JqmKtNToHGfIyoLqkKgAeUX0h+BiWoGR2SujwfjKexo/RcB61
hX/JHcnFvsD//Pkxmtz9Ed+P7loT7CxseEg2x7aINTwFO5nSMhBLsI6HP3r19kEdKhRHlXgYzT+P2i95
yGyQGCKzrkswY/1ekdFS0nsvst6mCsZWpkURZIe/3iS0KtCaXeA2IKTOeHjzQ1uYMCpf6wvyp1pgZaf+
3ZZqHoQwuaZfT6XPXoKvqf9hwVUXnY/MW8vzyfhTbTH7WKqGGk0/NhIAkZVJTuiqL8uvz45BmS2Jsy2o
HAc8MSrlPQ9tr563WqjnvYBYSZ21EB4XFAARiCUPa2vOhySzaQPs97yxVAzFOpXaFNat6uzKU/LNFb75
/yl8uIPfxMkPNT1BE9u/YuX/iqsaiuU9cj7SZPeXnJe6yGZtXbAwInfNl78PupVOJlJJ2vNwAco1g591
jiBIboHwedeIh34Q09q28d1S/Gpbd0hxV1872OLZksW/u5j7EbZY2/U3UX7TDq/t1t7H4bdLFKsPOZHR
1Rqo9h7/tBk6qPGa5r0yPh4OFaNjsAi8ubMvNPx7nkyOj/HzqdK3W2MlaoLL6D2bT0bTaBhNZtP402w+
+XM2jYb3XVvxcopPsRUDdWuLh7fH79sfyEJJscL0qim+KWfevHXUwKw6HXG3RmsUV00iAnSQ4gJyRS2u
x68L2vWC8yj6Vbs8/6D08vKi3yv9TvVXAAAA//+3NlCH
`,
	},

	"/definitions/AskToEncrypt.xml": {
		local:   "definitions/AskToEncrypt.xml",
		size:    650,
		modtime: 1489449600,
		compressed: `
eJyEktFq20AQRd/zFcM+t/EPyCptakIosUutEvIkRquRPfV6RuyMqurvi6IaGjD4Vdx77kE7BYtT7jBS
eQdQaPOLokNMaLYOj356JjM80FfGpIcA3K7DZztVupGYp97D3AIo+qw9ZZ9A8EzrMLK0On7s1dhZJZSP
1bf65Wlbf9/t64fNttr8KFaXznXEWVtMofQ80K2osycK4BnFEjo2idZhIgvlP0uIKh0fhoyzDXScbjPp
j19FvuiQWph0gMQnAlcw/E3zh3xlBlgABWjxoBY6zWf0T1Ad2SCiQENgfBDuOKJ4muCsmcAoDpk+QDP4
29bIKYGoz/FZZt7FGMns2jJ3i6AagR8JejQbNbf38HphqaTpDWanxWrhXJIwHkngQaenZzDH7HZ/6481
g7uKLW/95WdV7bb7+nWzr7e799VitVxZeVes/ru+vwEAAP//gDnlBg==
`,
	},

	"/definitions/AuthorizeSubscription.xml": {
		local:   "definitions/AuthorizeSubscription.xml",
		size:    399,
		modtime: 1489449600,
		compressed: `
eJyEkU1qwzAQhfc+hdA++AKyoD/ChBK7rRxKV0K2p0GtIqmaMcG3L0EUEghkN4v3fcPMEy4Q5C87gawY
E3H8honY5C1iw1v62QGiPcCzsz4eOHNzw+cyn/OMiZRjgkwrC/YIDSdHHjijbAN6S3b00PAVkEu9jDhl
l8jFwDL8LoAk6n9cVjd1JxfmeNqkiO7McdkOL+Zj25nXXpsn1Q3q/dJxS3GMs/VcUl7gbrQcu6E1QVm1
U1o/tMq87ZUetn13zzAuRDFggR/3w9B32nwqbbr+GhV1ebWsRH1RwV8AAAD//xVgh6Q=
`,
	},

	"/definitions/BadUsernameNotification.xml": {
		local:   "definitions/BadUsernameNotification.xml",
		size:    538,
		modtime: 1489449600,
		compressed: `
eJx00s9qs0AQAPC7T7HsPfgCq2BA/ORLFYxQ6EVGMybb6o7sTmjy9qWaphtij8P8+82yShtG20OHcSCE
ovYdOxbdAM5FMuOP3PS0BSuFPkRSm55asPK7VAg1WZrQ8lUYGDGSIzoHR9zwdUIZZ/X/5iXd75MsbV6T
qsiLTIU/HXGwTOhOejiImWBg2MxhJDsyjIYbsAi3XWu0LV0WVkuXe9kz60QjHdEgnZ2Mexgc+o6/ushq
NAysySy3lFWeFnVS52XR/Cur/K0s6mS3Omo+4zdes++gxWHR355N+g3Pnk8Lk4zZnlf1Qqhw2eExwgfH
Y8E9GfgZFXq/4SsAAP//wA2kIg==
`,
	},

	"/definitions/CaptureInitialMasterPassword.xml": {
		local:   "definitions/CaptureInitialMasterPassword.xml",
		size:    2728,
		modtime: 1489449600,
		compressed: `
eJzsVl9v2jAQf9+nsPw8Rtu9hkgdQx1qCxVDnbSX6HAu4VZjR/YB5dtPSdj4kxBSJk2atEcn97Pu98c+
B2QYXQIKw3dCBHb2AxULpcH7nrzjl88E2qZSUNyTfch46XBoiAn0I3hG9wTer62LZQ4XIsiczdDxRhhY
YE+uycR23cmsJyZrZHg3vY++DUfR0/hr1B+MpoNJ0P2Fqd+CiTVKwQ6M18Aw09iTG/Qy7FuTULp0KBZF
LyLbNnNuyxgTWGrurCnmuQw/Xl1VEJ5SA3pbr7T1KMUcTKzR9aQ1kQKjUEdlmRTdLUzNScei0NSA7hTL
nlzN7OtWoDqNP9nXUuDn/bpq2wtwKRkZXlfarav2GSgyadvyuV3YFA3apZdhAtpjG5R1hIZh5+14MhyM
prfT4XgUPQ8m02H/9qF2o0Ka3bpOlweYoS6VyZ19RO8hRSn2UdWWdImqycuTRvAoMHdH8LwSG5HY8rva
BqvgJRLS+F7w2gqmBfoPdXRqs+8gkyG7Za2SQgTdkvHBtwzUC5m0eWd8zcDEDTbVgRLSuqmdKmJ3amsT
lDOotBt0j4xtY/TAsNvsjN67T06mFXwnsSrP6lsYrcjTjDTx5px4B+cfFNMK+OgK8LDC4wvgn3P2+q87
e3PG2v8eHXt0c6lHhxwPfp6YVLmK1kTgEJoG1pLZmt9ja1Yso+bh1Tgpvownw+/j0fTyWVG2VPZTzuYz
MTs9JvoFvF3qlCb1gnGbt8Gp2F10rvbp5hm/mOz4/g+IvuV0tY3m7kfQ3Xub/gwAAP//9Iw7vQ==
`,
	},

	"/definitions/CertificateDialog.xml": {
		local:   "definitions/CertificateDialog.xml",
		size:    21609,
		modtime: 1489449600,
		compressed: `
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
PLXqUnPRfwIAAP//PVwErQ==
`,
	},

	"/definitions/ChooseKeyToImport.xml": {
		local:   "definitions/ChooseKeyToImport.xml",
		size:    2051,
		modtime: 1489449600,
		compressed: `
eJysVcFu2zAMve8rCN3TbHc7QNsNXbCiAYpgh10CWmYczbJoSEyT/P0Qy22aRHFbrDdLIJ8e39OTM+OE
/BI1Tb4AZFz8JS2gLYaQqzupvxu0XCkwZa7K+L0vBMhazy152YHDhnIlRiwpEI8uWBQsLOVqR0FNblfM
gQChph0Ig2la9pKNnwF6vGAqh7ZH05YDKVihKy35XLFbaHSa7CKWKRin2kqyJDSiJ3LyVnds1ytjS+hE
cGhH3TJXTwVv+zlTotzwNiry+3XduSYN+so4Nfn29XTaVHVoURtXXSzvuB3WKWL3WJCN1BoKAStSrxvO
z7SxIeHafEWwNJZgx2vQ0ULNTtC4AA17AlmhA3YErTdPKLT39wp6t5eeG5AVgTVBoCDLm265vwN7xA2v
bQnW1HS4ElepsVOsNx5bNRG/pnRHNo7KHO21qGvjqmFk2rboSjVZog0XwM+blsbaYTpnZ2fj3s0P2XvL
TcE3vJ3TVqLLNe2CGh790lFnlx61GHYL9IRqmMfNWoTdSwqK4xSkJGJvyAnuD1CTu/mvxexx+uNhfj2f
zh4WP2eP0z+zh/n1/UXFz8QZItazEjeKqVenjR/JwW0HcYkYJF4uo2sq3/N2Dbl25NxnSBEj9l9STJMP
9xCWRjcqaYlrK0MBebeOcYhP0HEgJqmC4wxlMSmjjSkrkvDScrQNnkLLLkQVult4uJHZ+Kj2TQSuFfQ6
5mqvYwR7/pMmwE42O5KHmbLxq9/+vwAAAP//qbBnvQ==
`,
	},

	"/definitions/ConfirmAccountRemoval.xml": {
		local:   "definitions/ConfirmAccountRemoval.xml",
		size:    520,
		modtime: 1489449600,
		compressed: `
eJyEkUFq8zAQRvc5xaB9yAVs/+RPTQgldhs7lK7MRJ6kamWNkcZ1ffsS1EALAe+00HsP5kuME/Jn1JQt
ABI+vZMW0BZDSNVWPvYUAl7owaDliwLTpupAHX/SWmsenKgrBpD0nnvyMoHDjlIlRiwpEI8uWBQ8WUrV
REFlG3Zn4zvAyIO/2tAmq5vhvjCQZtein5ZCX3LXvPYEEw8Qhp/HiE5AOCYI5M2EW/bfXG80ruVx2XMw
YtipbFs/Ni+7onkqq2aTF3V+mFN03KJVmfiBZr/GKy9l6imm9nlVrbd583zMq3pXFnOG0yDCLkT4/7Gu
y6JqXvOqKcq/aLKKG2eLZPVr++8AAAD//+yos2o=
`,
	},

	"/definitions/ConnectingAccountInfo.xml": {
		local:   "definitions/ConnectingAccountInfo.xml",
		size:    873,
		modtime: 1489449600,
		compressed: `
eJykk0FuszAQhff/KSzvIy4ASEQ/pVYpRIFVNmggQ+LWsS3baZOevmqgDU2sqFWXI7+Z+d4bCLl0aHro
MP5HSKjaJ+wc6QRYG9HMPTPZqzkYSvg6olz2qgVDP6SEhNoojcYdiYQdRnSH1sIGZ+6okcZZ/dA8plWV
ZGnDirsyDD7lY3e35WJNTusliNmpjGinpEPpGjAI4x4f1lwdBqRWHb5k10hbtVMblKj2lsY9CIuXGL4u
ZThKB44rOfgolywt6qRmZdHcl0u2Kos6yb2jTjbOtY89hxbFQD9GRqcN1zwoBNeWvyGNF0mRlU2a52xR
sVXapMV/H4ZvyqsBTWNn9t4MCAmDgXRiJrhw8xN3leZS4vjB2LG4TQad4y94i81zWRB8M54nyVlW/CaK
LR40yPWf0vgumDyeH8Jg8nO9BwAA//92Nwbt
`,
	},

	"/definitions/ConnectionFailureNotification.xml": {
		local:   "definitions/ConnectionFailureNotification.xml",
		size:    3111,
		modtime: 1489449600,
		compressed: `
eJzUV09z2j4Qvf8+hUZ3D7/0bDxDpi71lEKG0EsuGllezDay5JGUEvrpO0KEEJDB6WTS9uY/+7S7T++t
5RSVA7PkArL/CEl1+R2EI0Jya4d07O4LtdTX3FCC1ZCiWuqSG+pDCUlbo1swbkMUb2BIG7CW15C4TQs0
Gy++sK/57e1onLN8Pp/N08FTfBxuV3qdCKktJOWDc1rRbMmlhUNcAIoVyopsC1dcJtvbIRVaOVCOcQN8
V2GsoWv9GJop9eM+LFJNywWqmmb/H9cdi17pRtegQD/YSNVdKG0QlOMOfa+er9m8yKeL0aKYTdnn2by4
m00Xo0mvAuCx5aqi2cI89MrdcFOjSqzjxvVscgcBnyYO2O7E832M/gkvQYYN2OmFHgIiGwEShOOlBJot
5t/yWOIYDqTE1uJPoNnNaDqesXwyKW5ui7uc5dOPfVe5yGsMtDa8fR1iR63BeuVodhVll5B0EOh88azl
4h5VfYGM3+ii1RaDNDvLOcmdDo4k0EcS18HtwZTba9ZoA8xPmwvi2G/Ppy7PRVFcYq1oBqrqCzEgEZY0
U1r1TvNjl0aAn1R9UaU2FRi2xsqtOqknJLVuI+HlM0+3Z/WpZo4WKjo4G7QduMcx6SCyemqxVlzucSju
oaJkxVUlwQzp4baxEEvJ8cIncohLomj8YAhTQhtI/JpJ6RQ9hp7yh0KrxF/SzKO6+Iu7KaLgt7Ddq9T5
7Lurd/ddUMP5+gRXbKmF/9J1FHhJ1B/+Ate9v33+1OjqS8hFd297fUNnh/WwOT0DRNtAi+EQcOYbFgP6
gcDCQFijqvaHTLtpSi1R/MsTosNJlyfEy4oPXj6/SAcHvwa/AgAA//+CKZDv
`,
	},

	"/definitions/ConnectionInformation.xml": {
		local:   "definitions/ConnectionInformation.xml",
		size:    6434,
		modtime: 1489449600,
		compressed: `
eJzsmV1v2j4Uxu//n8LyPcq/266mgES7lrExmGhaabuJjHMIbo0d2QfWfvuJhC6BvPCSdG017ort57F9
/NNDdXCFQjBTxqHzHyGuntwBR8Ils7ZNe3j/STCpQ0pE0KYXWingKLTqq6k2c7b6k650hLiR0REYfCSK
zaFNUaAEStAwZSVDNpHQpo9gaSd1IRkb13kyWPtZESom125caguUzJgKJJg21cqPh/xkFSXOWsVnQgYk
vpNishV/bNPlRD+sz1l0x3P9kFzwNrsuf6c5M6FQtHP2f/a06fJ4t1RetFXPiCDZi28UMx7PSst2b6GO
aOfddrl2iCYaUc+3T76H0CIzeIQOVHCIyuhfLRsxLlR4wOW4lou5SoVn5crc0xQ/z4BNQCbvY8EswSQD
28L8UWQiLMD9OvYh6+eGgKD+WHbMIucZkyJUtNPzvvrdQb839K+97tg7xOJuYVFMHxOPLzfXXv/qhz+4
vDrIxIIEHl+MdrzxzWW51nWSuubGI8bvhQr3qCZMscUQGZ/RTilERUrU0T5C1yk8i+sUcHI0O7dMLqAW
OydOdlQuy8nZC3BSFxSU9haMFVo1ETR8YYCk3ytkmVif4mZXTRuJmwr+njluUopOkfM2IqcGKw3A0pWh
fo68YTLURuBsfkqcXVVtJHFK/939G4mzgugUN28jbmqA0gApV0KFYCIjFD5H6kxT+4NyZ1kfp3+eyKOj
6/1LRlcGyFOCVfNCiMsQjZgsEGx+Mju99p5qha0ALKdkuSpum37TStuIcaBOgb1T7v9a4rMGrEfRer5A
1CrBNRKqxcFgLUa/C0VWJmIqOEOorMNW31Xwewg2O6+RUE9919yDvpYn+9BEvuQvs7Voc8HGZFU7msVf
XD4zwKq60jEFf3rTk/ijX92h1kaAwnVffpUQo3H/cuh1vf5o6H8ejfs/R0OvOygqzj7d6yyYSUu+ukFb
TuTFSl3ard3N4Eb332nm0dIJ18n8KvI7AAD//5DzJeg=
`,
	},

	"/definitions/ContactPopupMenu.xml": {
		local:   "definitions/ContactPopupMenu.xml",
		size:    1776,
		modtime: 1489449600,
		compressed: `
eJy8lctqwzAQRff9ikEfEEO7tQOlL7LopunejKVxo1qWjGackr8vfqRJC4E2drozludw78GDUuuFYoma
llcAaSjeSQtoh8yZepLqmXyrwJpM6eAFtfQvum8BUr2xzgzPp2ZXQvUwT8bK3YHRH+xnAdImhoai7MBj
TZlyWJBTIBE9OxQsHGVqR6yWD8bKYrFIk/3EEYTtm0c3IlCL3aKQgg164yhmKvi8i5GPXRQkX+mTIf7Y
LDmq9reakeqwpTmKvvSkc3sOOeZtuqYGI0qI3yszNdcqmWwOuRq1vYY10VpQWp5k8JYrGAWABGAikA3Z
CNyzz1WLXO295hJyJsoH4CX+J3QufMztpWP+NFPvpmrpqP8npgyxsGZmM489dHY1Q9ZLujm9mjczrKZp
62blyzDJ7X1bN2B9Gc612KXIO8BvlB0O0uTojvsMAAD//yQCN7o=
`,
	},

	"/definitions/Conversation.xml": {
		local:   "definitions/Conversation.xml",
		size:    664,
		modtime: 1489449600,
		compressed: `
eJyM0sFqwkAQBuBz+xTL3kM8tLckYItIaEmKDfUY1uxopo0zYXc0+vbFRmtpcvC2sP+3w89shCTg1qaC
5F6piFefUImqGuN9rOfytUSy3GmFNtYV0x6cN4JM+hRXKmodt+DkqMhsIdbdTzxo2WOfmhcv5TLNyrf8
vczybBaFFzH+gIW12TUS1ICbWnTyOJncSjq0Uuvk4RbhxfEx6FDqoDUOSHQibgcD6HFDpjkzILNqQKva
kG3AxZqprJgIKtEqHBMW/ZBY9ANV1djY/jy2hSc+9CtY8UFfYsNee/R4GjfaZQzUvOUNEPDO62RtGj9Q
d/8JOwQSc91uvkhnWTEt0jwrP2aLIn2evg5nR2Ff6tw3/C18vYjCP3/xOwAA//9didl2
`,
	},

	"/definitions/ConversationPane.xml": {
		local:   "definitions/ConversationPane.xml",
		size:    7721,
		modtime: 1489449600,
		compressed: `
eJzsWUuP4jgQvu+vsHyP2If2Bkhphu1B0w0jyHZr9hI5TgW8bcqR7fD496sk9HQDDiS8ZjQ7N1CqylVf
Vb76krQFWtAJ49D9hZC2iv4FbgmXzJgOvbcvd2pFiYg7NFIrmpsQ0k61SkHbNUE2hw5dCCMiCbRrdQbt
1utVt/FMzdUUEFRmaDdh0hz1UFoAWmaFQtq9Dz6Fo/GgPwz8YDAahk/9cTDo+Q97QfhMyLj87SrrETC7
Y7osbQ6YRUzTV/OGJe6fV33mwMK8PJQrXIA2RV35Ffre+YQMXE6SRSApsZqhkcyySEKHrsHQbthHrtdp
fnjdYJkBL8MYtBQItBscyqNAg9h1Ch1qsmi+X18VQG8d2XM4ERR3e+q2yVim7SgYO1p0dlbNWjbJMyFQ
9g1iwmfMHo1uxBSZ3MRm3IoFs0DJjGEsQXeowrCoMFRWh6UxJS0nUK0SKQe4rQp0T0cdMP4uMO9jfAXE
AePvDO8FaJGs/xI4Bf3tUX8qsiFJkU6qBV4E9bLGMEkvjLvbwWG8b7hjtGvQThl/ETit3kuwShnGFYvU
5ZAIKQ9ssR3r/HwvZ/Fy8X72e5/CSeCPg1rOyohyaf+6b95ubdW2BcSx3f1VkqCyIhG8WKEe08DO2uCX
kB376btLGGCivsoPAzzTwq69JdMocFpXCVQ23eU1B2PYFN5187E/mfj3/fDZHw8Hw/sju7yQicikV/wt
5IsFtOEO7Ed7tlus96YraxV9dNM31puHvA8OwcfRePDPaBi4x6B6HKpReigpMcfpFZ5N325ByEvN0ktT
+ReVEaaBWCbzu52oBWjCkGSYamWBb9YqIsjDh16PlitHPF8jChtPeJRZqzC88lzfeDLviqK2Ciz36TfW
xDmbQD1hthOUM/RiSFgmbc2stjQGl4K/QHwbKX2GxHCpCBccR5WEy6lUE41cGkkKZ4BDsoLsSwsHUjfT
WbXNLyS0frus0JpwraSE+FlgrJblzT8Txiq9Li9dUm8tTBEyYtpLlRR8vYFi9DDofQn9v4PRox8MenVC
zc4IVUO8BbCyTwK2AbnG65t8IXtzFQPt/l7XB2JREGTD25Jn2ijtnaYtU7ECabwIpFp6UmDOzH/WfkcF
ifXmTE8FNvDSYjo76nbFB666U33+81Z/+KEZCThG5RAJnMwEKWAscPqTCWAbkJ9M8L9ighu+emlMBQ30
wMk0sHlA/WFpYCeUmbFYLTedEc6PJ82Zw/2QfxHmmDHjJYpn5mTCyZF6Hvufw+fR+EPY++iPmzNQo49W
72ngj9NooMKtAQ38UDzQ6AXsGwbt1ruPw/8FAAD///hY98U=
`,
	},

	"/definitions/EditProxy.xml": {
		local:   "definitions/EditProxy.xml",
		size:    6725,
		modtime: 1489449600,
		compressed: `
eJzsWUtv4jAQvvdXRL6jbPd1WIVIfbCU3QpWbXrYvUSDYxK3xo7sgcK/XyWhJQEHZcujWsStJjPj8ff5
y0ymHpfI9Ago888cx1PDR0bRoQKMaZMuPl1zEComDo/apBNx/KXVbE4yW8fxUq1SpnHuSBizNkGOghEH
NUgjAGEoWJvMmSF+5unkrp774mSPMVYRCOIHesLWTA2PJYiFIRXKMOIkICPBdJsoGVKQlImwMCOOu3Cj
CReRk59TgmjlyzaZDtVscQ7buS/VrDh02cySLuiYS+Kff1jNdrn1cm3bqKt5VOwUZ3+Vjde3m3LDh4IR
Hy34bM6xNVSIalyXqs1Rq+eWSYFyGRP//GNTN6rEZCyXnl9rHdfgsUN0C0MmCoxwnrJiueq2noco3CzX
Mb+JThbrW11utoCPE4N8NCd+N/gZ/ni4D3rff4d3ve5N8C9REhA8lkWQi9tetx92+tf1ATy3wGPt9xTo
E5dxAxzYCFuACDQhfi35Nk9UaRNHz7Xm4rkWfptyfqXGQ3WpZgGbYUF9qhUqqkQr481GP0c2Npbf3boH
1fdJAjJmUfWN8rJnmO0ZvpqYZ0hTFrWJVMTdM2Hn/wlhJZFODNPbi/TBMH2SZ1O2N1yTfbDdkajnS7Zt
RFfUBRT5FHClYBuYstVyvTdo3yykA0NbElIKxuyg2oExz0pHJzE1Zby21dm7mNIFVw0Iz3tBLjjOiT8C
YWr7Qeeo1HhgbkpqNExPd1HY7vM4JzU2ZfzTu6mxYPyoi9uBwS0XN6VxB8VNaTxJqSnbn9+vsCmNRy2k
A0Nb6RIx2UWXiMlJSE3Z/vKOHSImRy2knUC7fphXo7Mq+GuD4Qw7JUPQDCo4W2bEE0Qll5PifBlWB8a2
0yrNmUTItiku/uCu1+kHF0Fv0A9vBne9P4N+cHG77di0SK/IrRiNb/WCuMpDNP/IoILTp9WBWnVE3/AO
biWdMgqZALZr3WG6+TNrdRQOshWxEUwEbp7ZN4SvrOC3g7dBHDaD0sPlA88t/ePobwAAAP//RRhBDA==
`,
	},

	"/definitions/Feedback.xml": {
		local:   "definitions/Feedback.xml",
		size:    707,
		modtime: 1489449600,
		compressed: `
eJyMUk1r20AQvedXDCr0VFv3VjakrmpMaztUKjmK0erZ2nq9I3ZHVvTvixCFFlKc2xze5/Ay6xXhxAbr
B6JM6l8wSsZxjKtkq5c9YuQzvlh2ck7INqukme8JT5R1QToEHcnzFatErTokpIF9dKxcO6ySETFZP4MG
6V1Dzl5AKhRgYG+gUfpAJ6Cp2Vyy9I/g6/oRRnzDYVwoXvRVoycHjvhACueoj9TKQFbJRjqL9Wc6SZg8
pwh9BG1k3O2X79+9PH4qWxsnnLYg8W6kgUcaQIY9mQBWEFMNVQRSEUfastJgnaML0M1NjPgbQmS14iN1
wd5YsbzX679tWtUufkxTI+PSXtO3/mmwvpFh0Um0U5BkvS2/Vc+7Q/V0LKpNfijzH/ckrtKwS9YaetyF
zitZ6NhhttrnRfG4zavd4evxHrvuVcXHmfj5Z1keD0W1+X4s8n+ZWTrPc/2QpX/N9ncAAAD//0Jv9Ik=
`,
	},

	"/definitions/FeedbackInfo.xml": {
		local:   "definitions/FeedbackInfo.xml",
		size:    1128,
		modtime: 1489449600,
		compressed: `
eJy0lE2OnDAQhfc5heU94gJARCuEQekBiWGVDSpMAc64bWQbZbh9xE+mG+FEUZTssKue9V59JQIuLeoO
GEYfCAlU8w2ZJUyAMSFN7WsmO3UBTQlvQ9ohtg2w9ZIu/YQEo1YjajsTCTcM6Q2NgR49O49Io7T6Uj8n
Ly9xmtRF9ZSUgf+z3y03g/ruMaEMes1krZI0snrCk8zwXoLYRRrNqKRBSgaQrUAd0u2jfC/4u44NXLRk
zSxBeOsxpExJi9LWoBH2XK5ZXNTbNodGvb23nTMM6qZ6lKgmQ6MOhDnZd6mU5igtWL5kXuZWlFmSV3GV
FXn9VJTZ1yKv4qvzqTXG/ezyfoUGxZHi84aKPgrPvsQmtBqkEWChERjSGQ2NYo1kVhMR/JXLnnD70WWO
kMDf3DwY9g+Ojw2H4i+AAVvm9Oe81lWqf4/tPwO4bOt8ILDf/S2Az/sz7rE7FhME7/dw8TVL8zrJP/1r
ZvdC4D/8W34EAAD//zS7XW8=
`,
	},

	"/definitions/FirstAccountDialog.xml": {
		local:   "definitions/FirstAccountDialog.xml",
		size:    2368,
		modtime: 1489449600,
		compressed: `
eJzElt9v2kAMx9/5K6x7bmF7D5HaDnVoFUwUbdpeInMx4dbDF905Bf77CQKF/ABVolLfuMRf39cfbCuR
YSE/R01xByBys3+kBbTFEPrqUV6+GbQuU2DSvkrL39tAgCj3LicvG2BcUl+tDKdudZu7YMQ4VvHj9Efy
ezhKfo6fk4fBaDqYRL2Dpj2FGLGkQDxysCg4s9RXGwoqfiYpcti4wsPc+CCAWruCpZExmIzR7vOlZEno
ll6JRcECObXk+8pxopE12aSMVtDbq/XC2BR2RBjt7e7YV68zt94X3Ubo3q1LPL9O45rVLdFnhlX89Uvd
9fHu47ntpieckS3vWlIImJE6FTSvtKWgDahhTVugkDpgJ7DAVwIE7XhussJTekB8AysCJkpBHASSnajI
YWVkAY6p2+n8cQVoZNAL5wJt48wyd14AGWhtghjOtqE34CkzQcgDAtNq+wychyVygdZujrdXlHsj3U61
1haIrX3pMVex+ILaFVGv5Fx5lqN+MZxdzkzrHDlV8RxtOJO8KZobay/ZaSqOE9XaODsSdbtRr9ZOZ1ob
9TZzgp5QXe69+0LE8Vuvz3bHpNrxbe6dN8SCx5UwngwHo+nddDgeJd/Hk+Hf8Wh693QWRmMsLtmreCtn
XNW156YkftjFnzMCje2irdEvlL5nsVzutpa/68rSy/F7d+ltC2K4n+D97HW7V4IpLX02mMP+uQrN5G2J
fRScg63PxnNYuVfhGdT39tV4DrY+AE8zsBZUDTh5eXwR9U6+mf4HAAD//1JNv2w=
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
eJzUlsFu2zAMhu97CkF3I9jOtoE2C7KgRTxkRgfsEjAy42hWRUOSm+TtB9tZXduK16Cn3WKCpL7/Jwkk
lNqh2YPA+BNjIe1+o3BMKLA24ktXfJWgKOdMZhHP2t91ImNhaahE485MwzNG/Ch1RsegJCudJM3jZfqw
/blab78nP7bzxTpdbMLZ3xp/CyedQs6cAW0VONgpjPgZLY/vsowtDVXlqIU4SJWxRoQGFTSfEX/Z0enC
6RN1T6dW0dPbvDHQgZ4pR41UWR7vQVkcvu+rIiNRO+hcSDarxTq9S1fJevu02KSr+d2jt1FD33370B9h
h6qFz2tDgvrNQLXRt6VjrkuSx97G2ibLh8VYOGs5erESRCF1Pv0mnkrQ2YR5vqK9VIrHzlTvrui27vM1
BSPccDaw+z32L7Qz56H9fBruADbYk6h36BZNIJx8AYc2yHAPlXI8Tq+W/0cT+nL7hAYjGp17bRXpLRgE
Pj2/+8o50q/3v+vfvw988pq/JZvVr2Sd+u/Zv1RTYBcqpwMBWqDiw8JbrnretLgG5t8Zz1l8RAUVH1KQ
PFynH/cRoLtTmdrNG6SPEwdJ/YT+xobtXgZHmeXo7GtJL8wM2pK0bQU0M+/mH856uf/sQAVnFwsiXlvQ
NqPC22gQtA12Jyicdf8N/gQAAP//NgNsOA==
`,
	},

	"/definitions/Importer.xml": {
		local:   "definitions/Importer.xml",
		size:    4786,
		modtime: 1489449600,
		compressed: `
eJzcWN9z4jYQfr+/QqN3H0mn7fTB+Ibj0hzTHGQIvUz74hH2YlSE1pXWCfz3Hdv8sLGMuZJ2On1DaHf1
fatPq5X9D5u1Yi9grETd57fvbzgDHWEsddLnGS28n/iH4J0vNYFZiAiCd4z5OP8DImKREtb2+T2tHqSl
J0IDnMm4z+U6RUODKMJMky0nckfG/AhVtta2HB3GjLYp9HkSLYURxogt7/1tizmiAqH3836vsqTfK6G7
WXySQmFSpQDmWeoYX/foU4MpGNoyLdbQ56/FpJeilSRR8+B+9kv4PBqHj5OncHg3nt1N/d7exx2CJCng
jIzQVgkScwV9vgXLg1GBgO2z2BXoVca0DA38mYElHnz/w02XxxJksqSjy483nS5zNDGYsFiLB7dN+2gp
VcwKtWihvGLY5y9z3PDDdjXy/vUjbsqsV+0cgHGNCWjAzPJgIZSF0/WPGI5jp17FHFS5pCp+Vs2bC5c2
rk0aLhEtMFoCE7uNYlvM2CtmKmZKroARslJMTOiYWRKGWGalTtjC4JoNcTv68t5FxLnLRqQ8IJM5qbd5
eGuMgQePg/H9JHyeDh7D58n0Uzj8PGjocxekck6OgUW0kjo5vxxsUqHjM9vjclpIpb6N1fHINTS4Z9CA
6/dOlHGJUp4ig0pBvKsDtdpQznVo58UWVnNhvBSVjLZllXicPIyGv4WDX2eTL4PZaHgp8eWV4RqU3bRn
BuCrhBPCh39PAzRx5opTPHDcBG3Q2uGdhzgsynsVqEdLab3defTK8u+A7ILdWY7z0OcYuIJaNBSWMEIZ
8+C7zgCtaXCnYghKTUHHYMDMMEkUtKfD7OxaEuLCLyKSL2Uyzp/SSghHAanNCyIj5xmBPYPjYFQFAkUC
D1PtCLqWaNSDS+C3Ol0r3fwy8ESaKhmJvLK9lW5/zi+ZStxrxdtSby9JhDsZNfHChlrScdBt7201RfmS
Oav/naL2Jz7n+VZq2lXywu5aJd3+O0qqpeGfVdHtf05FbgeH8fUd37e0b9c0fC2yubzhazxP8msFdSgM
CH6+G/yYEaE+PleKYVh/tbjQo5GgSRwfiZPp6G48G8xGk3H4eTId/T4ZzwYPlyZjLUwitUeYuh5hdbbd
vV5JqsYoEjpqvIpcSM68jIoQ7Uf8YmVeRQRXV5Eom76zHetJrEhoL4aFyBR19UpXnM0To7pBbdIvxe29
yjgBOn53qf3NDNgUtS3xFztfE4Lfq5l3BsEVZ7sk9HmehEM8XDljnfx58q3G71U+P/0VAAD//2D0Xys=
`,
	},

	"/definitions/Main.xml": {
		local:   "definitions/Main.xml",
		size:    9171,
		modtime: 1489449600,
		compressed: `
eJzsWk2P4jgQve+vsHzvZjQr7QkiNWh7Zw69s9KsZo6WY1eCd4wd2Q7p/PuVSfgIJCF8NG1Gc4Pgqnp+
ryiXKcZCOTAJZRD9htBYx/8Bc4hJau0E/+V+PGWZFIw6odV3obguMBJ8ghdUrN97O4TGmdEZGFciRRcw
wYwqkmiWWxw9U2lhPFovaF/vhJOAo5kuP78cW8shobl0pBDczXH08cOHoRZzEOnc4eiPFhMrUkXlxsA6
o0uM5lRxCWaCtSJMagukWG2bVMsxsgXNMuATrDQe1a7YXEhevW7jdKpfKxY/xfoVr9cdol4KK2LPyr8m
PyDwHNIPwR0B+K0J8NyYbXbaCFBulVg4WoJxglHZaXwAuh34C6h8Sk2do6DymBq8b3b+HrqhdMP57GBR
4Zlp5Shz1j9tAXUpsDZ7SWOQGDlDlZXU0VjCBJdgcUTWaE51mVt4yBUHI4XqSc1DypArM5hgm8eLbgK6
WNwq2ml4OXuoV95+gFuZKec9Cl8Ta5ufbsWfOH98fBzst1EOKXNiSR006yHlnLAqjQZUxY44o4rOHkVG
RyS5XDIFxUyrZZCy/Q0F8uDA2FWxvLKGCgqv4cb9OyjY76DHuNuw0+j86v3EmM5VKNV7jebNq/dNOf4m
oAiDX4/k3k7GjwEejbM5sP02yD/yb1/ApBBcvV2hQkO/Xxv/jVLrdJpK4M1KWz8kzO+feALIKtbdnJl9
Wn6d6+JLklSZH5iiHhuqwaGhbe8mzFnC+pCkDvnT6PudCidUGqa+Nbhb6luH/Dn01cZ9ddR5mUKTVxuH
piWq4L2xrto4Mi1JFetXL9zRp33J/FXhxFb44Y1atRrMPXVruoJMjvkIqm/7UzFTZm6mVSLS3Kxui89C
QnNpaNWjRo3YLmyUCDk81vBKAlUw0ghGfLC7OSSammcGEjCgGNhgFSb/bEGe8uPI5fVg42rA7yw7TP46
Vjqy7RPILIjrvwdy06v77wEeAU1tEgAe06rQB1cCnmtw1/zCrjdMuKBS30+XvzeMiHXugtTsySO76jTC
O3w/tW5XYrsMxhllP4RKj8874TWjig8Zdu4ZJkLK42fjvlWmrajmvAdT991dtcJv5WDoOHgzx1baiaT+
G8MDNUCvPRS+ZLgdgKInm3lgxJ99OALF3zEZ2qhrp62Lsr5sPj3/T95uy1b3ttnc4s6H2w/Go52/8fwf
AAD///NVoJI=
`,
	},

	"/definitions/MasterPassword.xml": {
		local:   "definitions/MasterPassword.xml",
		size:    2357,
		modtime: 1489449600,
		compressed: `
eJzUVk9v2k4Qvf8+xWjvhES/q7GUUpSiJBBRlKq9WGN7bKYsu9buGMK3r2wnAYL5E3rqce15ozdv3j47
YCPkMkwo/A8gsPFvSgQSjd731J3MvzJqmyvgtKce0Qu5J/R+ZV2qKgBAUDhbkJM1GFxQT63YpHbVKaxn
YWtUeDe9j34MR9HT+HvUH4ymg0nQfcO0txAWTQrEofEaBWNNPbUmr8JBRRYWNQ8oXomcapdShqWWzopT
manw/+vrPYTn3KB+r9ck1KElGVEwQ5Nqcj1lTZSgSUhHTbWC7is6mbFOoRbSoO7Ux55axvblVaM2Yb/Y
l0bV5+26ffYLdDkbFd7ssW6r9gUmbPJzy2d2YXMyZEuvwgy1p3NQ1jEZwc16x5PhYDS9nQ7Ho+h5MJkO
+7cPrY1qaTbnNl0eMCbdKFMt+JG8x5wUbKP2KekG1WKZJ03oCah2jszoo3sgs83zxJqM89LVc0HGmq7g
py1hxVqDsQIxAfo5vSHYb3pgjmygNMIa1rYER17QCfTtevh41SZF69VxWKhQXNm6BYCg26i186zAZM4m
P96ZXgo06ZEVt4Ey1voYnX3E5tK3uq+aYI9u0P1ginNMMjDi1huTbMXRQaej72Q2qXz+mYmW7DlmzbI+
Jd5OhGAivESh3fjwuKSP4fHPbfbm0s3uzrjz8kCCVipaE6EjPBakpYg173Ea18foeKgeTbBv48nw13g0
vTzDGkoNn+abccKdh+OrX8PPc12iOZlTes4365DtLrqN2+NWHr942PH9Xwz6mdt1rjU3L4Lu1o/SnwAA
AP//WXbTAg==
`,
	},

	"/definitions/NewCustomConversation.xml": {
		local:   "definitions/NewCustomConversation.xml",
		size:    4733,
		modtime: 1489449600,
		compressed: `
eJzsWE1z2zYQvedXoLijijOTHjoUp46iKGo9Ukdmk2kvHBBckbBBgAMsLau/vkNRX5QgiUrsaQ+9mca+
xe6+3QdAgdQIds4FhG8ICUzyAAKJUNy5Ph3h4510eI/GAiUy7VMuhKk0OlaYFBStMYQEwqiq0K75IiT4
gTGytiSaF0AY2641tgSXJfRpJnJuubV8SXtetEyvwAa9vUCCXpNL+Mab10fJlcmapCawGFQOTTEw+gms
4yiN3uRWWlOCxeUqkT5dSJ2aBSuNk43VKPot/jqexL9P7+PBcBINZ0Fvg/G7SIxNwcYLmWJOw58umaNE
BZSg5dopjjxR0KdLcDS8R26RaFgQsRf4JX8WnPy79kLDyFZwyTyFOa8UshxkliMN37192xWyzvB9F4RD
a5ZsITFnJbegkYboi07kUqVk1bSaK7b67NOnxDzTbZsckf3BPDdMf9m3Ow4jN4XJQIOpHA3nXLmj/X0o
YyVo5Lt2mM7Gw0l0G42nk/jLcBaNB7d3XRy5kgups4Oe2Jmvkt3BfZmOrEybVLP6r33j4/0KbjOpGZqS
hjfvfRGeQSUG0RQ0vDki9wLQ1V37DTjQ6TUoaxZsW9Cbd11hjYKwC1T46fBTcscTUG31XP/rEHscjWoM
PbP/yZri51NJ+Vw9VA7lfNl06K9/3EfjT3/Gs/Hoc3SNl5wrma3b/PZuPJrEw8nH0w52InzolotHqbMO
FYA5Mo7IRU7Dk+T7kGjKLsCg540l6Hno7Ur5wBSJ2YrOhnV6Oeb1qdo+Za/jx7G5EbV++eTzLBSeS17P
mO9UuFyFE5UApWagU7BgI3jGVkVYvTGzzXK6uwG0XHJEK5MKwfmW9w02vNfb1IxvF7x+e+cce9l/8Xa+
+Rfa+XslrASwLyBfXyXm/8tXV77PNMpryNdQo12uJzVNLbgu0lUqLiA3KgXLmhF0pgCj4ZcHniRgfzQ2
ex05cjLTXK1hXKB84giU5FynCmyfGh2v7hxxY0jJkdD8V+b6RXj2JeNP5DCUTcXPnRyHmLlU6jrE7ul0
Il1PqsfydeIhUNNvdMwtcHr+pvyhQjR6e0gnq8+4/T7whX/2tv95Ohv/NZ1E/vs+uWIGm/BasQmuxfbN
fS7ERnbDeLACXNOMlQNW1cexkvrE+7AFb42eUFI8QtqePKGMg2sn77vUy1M589i5aiff2a9ZxsM3CNds
/YbucI27zME3qV8nZTkwahvsLe4Wgt7eD07/BAAA//9aple6
`,
	},

	"/definitions/PeerDetails.xml": {
		local:   "definitions/PeerDetails.xml",
		size:    9977,
		modtime: 1489449600,
		compressed: `
eJzsWk1z2zYTvudX4OUd0et+5NCRNePIqqqpI83IbDLTiwcCVhRiEmAB0LL76zsk9EUJoPghp0nbmyhi
F7v77LMAluhzYUAtCYXBG4T6cvEZqEE0JlpfB2PzeMe1uTdSQYA4uw5ophQIgyMls1QHuQxCfSrjLBHa
Pu2ekXlJ4TqI6IooohR5CXqb8b0DgX7Pzjl445z/A4hsYiCx0xPGxvnM24lTJVNQ5gUJksB18MQ1X8QQ
DEKVQb+3feseHJMFxAEyiggdE0MWMVwHL6CDwRTWqJjm7du3J1o0jwSJNzoINfyJGAjQiggWg7oOpMCE
MSxgbYMUoF5NN62LNrI4yf9o7OWZWW45iWVk52H2t3sKw00MztiMGDeISmEINdWxobHUR4GhRNA85ts8
oCseM1RkoCAxLh6vg6eFfA52uXTixHv5bD34eDju1IeEqIiLYHD1/0M798OL2fbirqnGirMtKpwFh4N9
82Ej02Bw9eNxbM5ILaQxMjk2toagNkSZFnIgWBMpJddYp4RyEQWDq+/qilme7yXfOaFww+GG5M6ytigG
lMpMGGyJfCzahPA3VpPPK5eyFYl5JILBOPz14eZuMp4+jKa3TRRoiIEWRgSDcP7byC+7J/WxRkIfuYhq
eA5Lg4kxhK6CgRdzl6SRaR3Bfs9pS7/nQLUD0rlNNYA+xeY+vJmH3wQ6V38DOl3h+czZBUg4dC8qzYD+
x5OwIj9emYSfT1bAepD8C7jXAZSuqAhOH3ODLkDA6UbVKzLwa2GRdwPzOiwaCaNeynjVQGp7uNCYwZJk
sXGfa14tuq3p0CG6Xemg4I+MK8AgqHpJDZfiAsSYW6VorxStuVkhs+IapQDqvzWrMh++/7JsG66APr7P
jJHClxTbTkTnmLTmSIeYdCZJqXlT6aqXEkV7RHfL+29nY9A6839oIrjmzFxwbfofxiic3c5+QmFep/RK
ZjFDC0BcIII0VTKOgaE1F0yubUFLuOgl5Blp/icghHE7/oUK4COHdamf9ZT/cT4GiWR5sS63GZsEsbIF
WCVIicBLSTPdXHQFhIHKfdzM/TOJdSMNeiXXGJ5TInJFLRQokIqBIvVcd/fdLL3y6ngq4sP5fidTgF2t
wg23LQn3o7vRMJzMpg/3k+n4roLbG03lXuOKiAhYudtYJA/emYR3g/SapCnkmzF5uhCgivLhoRryksMf
uJwgw6I7dkCTosGBbdOsZgRtp3awFz8bOK+hbmOHEMdzyPMSVAjPxmMuVnYMc0e00E2MUXyRGdC+IYeD
tg7mc+YVePfCq793bgIvei0gL7XaSy++8BrjbTW7BLuvMa2WBLsb23XuN2vCovi3zkZEKg7CEFtZ8nIx
m09G0/CmKBi/zOaT32fT8ObufNGrR9APILLDDSRhDC9MXUrulgFzpg67hFOZZumGz/YrkLdl7lPBuNoW
4WKzNZ/PPj3czj5NL10aJgmJoBQhzJPIEyWXpZxKsWnqxlwbTBg7Z2MVU9GFKd4saconjkQ+QVXOHH0q
4/TxePna6LAfET1FtQNiexM7gmYVfU24fTWl+d0ltvKn7p/WYs+OjhRl4IEoICWAzxZoW5kfyh9YXd52
LMzN1o+NbWb3IbnLCXZYqKhE9jxHN3acdhOcKdupq3EcAk2e6jQR/QG4J9WsdR2Qdo3Ic0tbjdgVDrSP
XA1aXPQIfIxTxcWEY9ocR3JRnNXwZi/m+wh/csqUiYxAgMyPqEv/AbERRT+O5uFk6COoKzcrO0tLLiJQ
qeLC6IlYSpUQ52mwqrNTlVqu7HDlRi2797c7Ds0en970aHvXo8Ntj/b3Pdre+Gh956PhrY/GcFYQPX/Y
sr087GDI/kW/d3DZ7a8AAAD//7iO0q4=
`,
	},

	"/definitions/RegistrationForm.xml": {
		local:   "definitions/RegistrationForm.xml",
		size:    2189,
		modtime: 1489449600,
		compressed: `
eJzEVkGP2jwQvX+/wvI9yrc99ORE2m0riroCCaEeekETZwjuGjuyh2X59xVxlhBiuktaqbdkmOd57/kN
IJQhdGuQmP/HmLDFT5TEpAbvMz6hp88KtK04U2XGy/B8bGRMeFUZ0MzAFjPu0NfWeORsA6bU6LpS0lY4
S2NIqW0P1rz3MAEkN0qXrGFrQCfNa8afC/vSEoqxf7Avgfr38z7GRO1sjY4OLYln5VWhkefkdijS10+v
A7bgKmWSwhLZLc/v/n8/6LK7a280dfCYoIlTZVBUHZ/Om0eoGq0sBnR2n/gapDIVz+8+vBcmrd5tTYf8
eBU4sCdu0SMUqINHynhyO0nKGs8vkSP9YkykYeagXoN8UqZ6e5DGNSVABHLD86sGx5Bk63HAvSppw/Or
13KUFeUv0mEsBwacmi7CPNhXaC5jBQ6B/z7pDzsia04LXPQXePT1XYKsU2gIjqx4Pll+W80X0y+z5f1y
Op+tvs4X0x/z2fL+8U9DGdS0UsgkEoxE/fcyGclYWANyYLwGgkJjxg/oef6pmX1rviM5GC/fYaU8ofs3
Biza6becJsEkJa5hp2nc98NtexRr6C+ZCKuU7FVZIfkTpFdmrz/CjYAmcV36RNrrffMEqGt94Kx1IeNH
F8J57mRo5MSLYsO0EybSsz8gvwIAAP//h1KDAQ==
`,
	},

	"/definitions/Roster.xml": {
		local:   "definitions/Roster.xml",
		size:    2558,
		modtime: 1489449600,
		compressed: `
eJy8Vk2P2jAQvfdXuL5bUbX9uIRIK4oQ6nZZLelWPSHHHsC7xo7sCSz/vkoC2XyBNkjtMR6/N+N5bwZC
ZRDciguIPhAS2uQZBBKhufcjOsWX2AEs0DqgRMkRddYjOLa1EjTNEYSEwupsa3z5RUj4kTHyrCRhrDoq
rxA8pDCia7HhjjvHDzRogKTyqeYHYvgWhqO5EDYzSK7J7JFj5omw2rrh6ISLl7WzmZHXMuxBrTdI7Iqs
rMHzeGWwhURrNap0eEolrDmLmsqXB/WaZKsWyNl9cWNAujCo+SMMSoP1e20hnNUa5G9lpN3XDXeyWups
Cg5Li4zoxheIhDuWWq3EgUbT+MfyYX43G/9Z3k+eJo9hcML0U+wuUtz+iuc/b+PZuEMjNkrLqge9Y/Ok
oPEItssPTphuKcehqo9YO21vE4BLcJ7tlFeJBhqtuPbwHqTf2D2D15SbnGAAUMMONFNGgkGOyhoa3fTi
cs/E8+9z4mBrd0C2fK0EMdk2gfqUdDMcbc1K79Doay+9V2vD9RHi7J5xgWrHESQlG26kBjei1ixPx8sk
k/JASVDjKIQkxRY0XLPic0Q9aBDly6qrZ9djdbcQux/aL3dptsXkbjKOZ/P75WJ2P72b9L21OTunk4YL
26687Mxx2dmi5HwXVK1uFt1h7Gcdg9aPkPsI3HFztJmZK+OSBh0+juhUkiH4dqgePDYuLfmjb2FQhd4B
E6A1e1vVNPp8AV8LtUrqNP2fKpOXfkaZ9sQozIc//+j3z3VixvCKnVKulvLSPvCNhXBWRswLij4N035l
HZxkvxkGfa9jeqDlrzqNvvw3ozUv1IJvgTCo/eP7GwAA///UIu1L
`,
	},

	"/definitions/SimpleNotification.xml": {
		local:   "definitions/SimpleNotification.xml",
		size:    312,
		modtime: 1489449600,
		compressed: `
eJyE0M9qAyEQBvB7nkK8h30BV+gfu4QQLdXSo5h1Gmw3juiEkLcvRQotFPY2h+83M3wiZYL6HmaQG8YE
Hj9gJjYvobWRT/R5gNbCCR5TWPDEWYojj33+zjMmSsUClW4shzOM/JpyxOu2YEuUMHM5ub1/22n/bKx/
UNqpFzH8mP9XnDGGhUuqF1iN9u+2dCvQTx2UtXeT8jv9ZNb08UKEuXV4/+qc0dab/V8mht6J3IjhV1df
AQAA//98bmjD
`,
	},

	"/definitions/Test.xml": {
		local:   "definitions/Test.xml",
		size:    261,
		modtime: 1489449600,
		compressed: `
eJx8jz1uwzAMRnefguBeyEs7yRq69AbNLEt0xMQRDZmxndsHgfM3GNkI8D3gfZXlrFQ6H8hVAFbaAwWF
0PtxbPBPjzvOUWYEjg0GyROV0StLxhsOYIciAxW9QPYnajBS58+9fiXifVJ0P3VtzYP5rMwcNaH73jBC
4j6u91bj/68sa+HUyoLOmhW52+apvx7WvO2+BgAA///uaFOk
`,
	},

	"/definitions/TorNotRunningNotification.xml": {
		local:   "definitions/TorNotRunningNotification.xml",
		size:    715,
		modtime: 1489449600,
		compressed: `
eJxsksGOozAMhu/7FFHuiBcAJKpFLFoGKsqpF2SoKZlJY5SkajtPP2phWipyi+X/t/7PTiCURd1Dh9Ef
xgJqP7GzrJNgTMhT+5WpnjagOROHkAvVUwua36WMBaOmEbW9MQUnDPkJjYEjevY2Io/S+n/zkex2cZo0
SVWVVeD/6t12M9DF6yQZ9NqztaR4ZPUZVzYjjgrkbNJoRlIGORtAHSTqkE+P6tnwZ183CHlgD1wF0nuU
Ie9IWVS2AY0wc7nWsKHrtIKWrk/ZmmGgEx1RIZ0Nj3qQZhXf5SItUFmw4s5831tZZUlRx3VWFs2/ssr2
ZVHHuXPUA+NVu7Ln0KKc0s8n4kvDOg9KKUYjvpFH27hIyybJ82y7y/ZJkxR/XTFcUy4aRvcJZ70/JV3A
+G8074JF89UI/MX3/QkAAP//mUzfSQ==
`,
	},

	"/definitions/UnifiedLayout.xml": {
		local:   "definitions/UnifiedLayout.xml",
		size:    8750,
		modtime: 1489449600,
		compressed: `
eJzsWk2P2zYTvr+/Qi/vgrtJi/Zg67ALdBsgLYruoldhRI0lxjQpkJQ/8usLipYs25RkWU42AXJZLCXO
aD6eIechPWfCoFoCxeh/QTCXySekJqActF6QZ7P6yLR5MVIhCVi6IJxpo6uhnR4Ecyp5uRbajYJg/v8w
DNyzUMAagwIyjJlIcReEYTPLzQjMvsAFyZgwZNapofrTKUtzUKAU7Hs0fGLpNAVUcqmmqUiArjIlSzHR
lC2yLDe3xtJIyQ0rppmgDZhSd+p4Tld/s11SLntUlCJFxZnoSWzLkfmshbL5zGHUD9c/EFJUj6AcXPNq
mICq4VooWaAy+wpVC7JhmiUcSfSqSpzP6rf+yYYZO/VJ7j/8OTRX53IbUy41xklpjBS+T/S78ih3zgll
M34PByiIeClpqUn0O3A9OF/IuPIDOO/4AM0ZT5sEXrjwggUoMPKQDV0PH0gtM9KjW7zyyUjFUBgwzCZm
g8owCvxSsJ0gpwboiomsWzHuChDpCEuWrCu4vtmF1MwZ/ZPP2hPz5rNWdoYy1YAtkbs3z05eR/Haj1yZ
zvMwDITCLR5xInekLXFDUG4NjN9kv9kfIUF+YjivnpxL3mj+FBcqWTBGsaQ0qC9ftl8fdLt9jgQb4CVa
WPL0uJ+05Gbdes+r9+jGeRV3OTiEwy65/qLukuot7qNTXutPir15eCV0Ht0OVWHndM/6KtgZK6qQItug
jlNcQsnNLRo4wyWJhBT9cpplAnhtMmd0hSkJchApR7UgUsTNQ72FosDU7ppeoHpz4c/HhzVkh36b2X8f
PImYlIyBhAwWs0+eUSli+y+JtkykchtWUAr1fp1Izmi/vq5a7UD2PYt7eOm6T3UDXcW2tyURivS2deHh
HuuCL3L+qN0UsfHRGr0Oeny98POanf4vaTCRcuVKTdSjL7vZjxHTVEnOYfTnqq7dQDLYW/jkEqlSVCMl
16AyJmJtQBkS/TJSDC2quoVOVmG9ZYbmYVEtke2V2L2I3YvO1Xg6+sckYgr4O4p9GPzfMGvx+DSFtXTx
y3dvzWC+f3757svxy/ffSXau20RavKsiXOfN2tcli2dyjgIGRoHQHIzdSRZkj5pET1JsUOkqCLpTXQ+5
uoGwdZG1Hx1JDaZXhfgvw63Dk1GIGzv6djqS3LUkcWFb+j2JBJhS+YvHJ76WKXISNbcIV3+2OtHQcePq
qISjsMCPNYKi+UhZoIZtwGBsOwwmMo6O8g2E7aRlUXIb1nrO6GP9uJc/OuQE1U2NAB5WQ7vdcaQOloMU
3wLrpZl/gq6w0fNwWbATThVqLD9VZ/cHOkuvOlVwJ+2+Zcsq6OVOI2j2E3L+D4oUFarDjUVjZHi4rVDu
vffsafBQy7NKUuQ8PN4Fkejn1qJ4pY7C2Rr9OiDadzY2klffHQg1+FyYfQccF9ykAMpERqLfxnDnNRPx
lqUmJ9HDu55jtTEYdDaP0XX1OeIhonD+Ygi+r7gzLrAlS8+wq+6F3YmwXUqFtfj78eLGumgJxFjBQ2Ni
+eVY0eaaskLdGxbbD9b6llzo/Rgu1H+p+8I+47OSZXFgrOxzVRNFx/Wu7Ze6Kcp8y9IMzfHHF2580jNV
11cz/4wDX2lu2Vvqjl7MZ61fiPwXAAD//0H+ERQ=
`,
	},

	"/definitions/VerifyFingerprint.xml": {
		local:   "definitions/VerifyFingerprint.xml",
		size:    1378,
		modtime: 1489449600,
		compressed: `
eJysVMGO2jAQvfcrrLmjtPck0i5F26irpKLRVuolmiQDuDUeZJtl+fsKhxYczJaqe8s4896898ZJKrUj
s8CO8ndCpNz+oM6JTqG1GTy4nx8lKl6CkH0G/fB8aBQi3RjekHF7oXFNGeyk7nk32bCVTrKG/KH+3Hwr
yuZL9bWZzsp6Nk+T35gjRbeSqhdegUY18WUGzy2/HIfEFN3zyyDnvO1ST8umJzPZyd6tIP/wfjw8hlnx
mpekibcW8gUqS7eg2EjSDk+uq3kxK+u7uqjK5mk2r4vp3WOUyBs+1TG3j9iSGvyuyVpcEiRnDMmI4kqk
2B3UNWgI4fV591vnWJ8y9mUTRv3PIXyq5sX3qqzjMcSjeE1eoK1D3ZGCMfZSYod60tMCt8pB7sw2utxr
YDXswRnUVqHDVlEGe7KQT/3861xpMrgYGR5v7j9TeCYjF/sbUrhu5MlTvIGRy8ZRU9gQvEyHq3r4bpfk
7B9EcCwM2Q1rSxloBnFcagaHpUIe3Is0CZB/5fNJBJlGGUaHXufJVZqc/VV/BQAA///rVJ76
`,
	},

	"/definitions/VerifyFingerprintUnknown.xml": {
		local:   "definitions/VerifyFingerprintUnknown.xml",
		size:    1102,
		modtime: 1489449600,
		compressed: `
eJyUVE1vm0AQvfdXrOZu0d4XpCS1UpQIKhe1Ui9ogLG99XrH2l3Xyb+vME7M8mE1N4+Z9zHvIaQynuwa
a0o+CSG5+kO1F7VG52J49LuvCjVvQKgmhqb73S4KIQ+WD2T9qzC4pxhOyjR8WhzYKa/YQPJYPJW/0qz8
nv8oH5ZZsVzJ6A1zoai3Sjfi7MCgXpzHGP5W/HIRmXJ0zy+dnf7a2E/FtiG7OKnGbyH58nkoPoXZ8p43
ZIiPDpI1akf/g2KryHi8Xp2v0mVW3BVpnpU/l6sifbh7niQ6H3ydp659xop0d++enMMNgYh6FNGAYyZT
rFt7JVpCuC14f/SezTXk81iGWX84hW/5Kv2dZ0WYQ8g3CuOWv8Ac72CIG/ur0SwaWuNRe0i8PU5WOwfW
XQveonEaPVaaYnglB0n+NM8jo879ILhhZf3FuWJDquCh7Lpt3/QNefeOCP4WltyBjWuL2oG4BBFDGwQk
7znKKEC9qY0Vrn5k1PuC/AsAAP//GrNMEw==
`,
	},

	"/definitions/VerifyIdentityNotification.xml": {
		local:   "definitions/VerifyIdentityNotification.xml",
		size:    965,
		modtime: 1489449600,
		compressed: `
eJy0k99qs0AQxe+/p1j2PvgCKhgIftJUIZEWeiOjGZNtNzuyO7bx7Us0bWJj/9zkcpk5M785h/WVYbQ1
VBj+E8Kn8hkrFpUG5wIZ80tiapqDlUJtAqlMTSVYeWwVwm8sNWi5Ewb2GMg9OgdbnHHXoAzj/K64X6zX
UbwoHqNVmqSx730oTgOqndIb0RMY0LP+GciKDKPhAizCadUU2ZwOA1VJh8+2a6od7WmLBql1MqxBO/yK
MaUiq9AwsCIznJKtkkWaR3mSpcX/bJU8ZWkeLSdH9Wec31PsSyhRD/Qn1+Sl4JrnzUIjQ7btJL0Qvjfs
uMDwRhzjhlHxmxigOl7/9xRaZjLFz2Hc2NZ5zzDieUWr6u4Xd/UQB1swTgNDqTGQHToZPvTyac+v51Rg
ZhusodV8m7DOBd+7+LjvAQAA///JzCYS
`,
	},

	"/definitions/XMLConsole.xml": {
		local:   "definitions/XMLConsole.xml",
		size:    3090,
		modtime: 1489449600,
		compressed: `
eJy0Vl1zqzYQfc+v0OidOu30djod4I7jeu719MbOODTJ9EUjpAXUyBIjCWP313cAJ/4AbJP2PgJ7do/2
7B7hf96sJFqDsUKrAP/4wy1GoJjmQqUBLlzi/Yo/hze+UA5MQhmENwj5Ov4bmENMUmsD/MW9RrBxd0WS
gMFI8AAzrayWMNHKgXK4AiHk50bnYNwWKbqCADvYOIycocpK6mgsIcBbsDj8ClJqVGojuT96A9WFR03l
8KaTxe+CSp02DF7uv00aEj3VS6G4Lr1cW+GEVjj8Ev1Bnmdz8rB4JJPpPJouj4u3U8TacDCkFNxlOPzl
UrgTTkLngV/uHx7Qju5vaDyZLP6cR2Q+vp9eymnAin+qTDiMTAGXwjkktJDOy0CkmcPhz7e310J2p/x0
DcI6o7deKVzm5dTUE+C62LFMSI7q2VJUevVjgNex3uxE61L5Tm8aiZ8O49o0Mr3SKSjQhcVhQqVt1e9C
aSNAOboficVyNp1H42i2mJOn6TKaTcbfrklkc8qESk/mYh9eH3YP7zrpIzNaSuDP9ag2h16BtTQFYnff
mjHGh4naXNbCit4Z6QMxqkiiWdW9ITClic10SaiUw4AXtqm3zxnluiRumwMOwbEMuCfUteiVUIQ1LkXe
luJTL7ilWbdulRs+CThRrLK7dfX2NMEZvbq25hywNDT3VpoDDn8aggMuXOMivYvSOyaFsdp475QHJ8jF
BqT1YpC69KRQlSH2StCVQELivBU1qVADkaYS/GPQuLntwuObrj/D/uo6ftu2gY5AP6fsVaj0/CzDJqeK
D1u6RJxfU3/UKv3O+cTLWj5OWeWihBqg+LzR3RXOafVu7HH9SI7tvYv9WbP+uljO/lrMo267RgP2uaF3
xM1AYsBmV+yypDFIHJJlgxgyY4UFr1AcTLUV53Wt4Vakisq3vZSCvQLHKKOKSzAB3tMmTSRGow+P6X/q
HpPawtW96/pjIpMqxfdsZvs29Ha/QlfY8mUl6hb8Dzq0A0+CjgMOPu4/+KOD3/t/AwAA//8oH6wZ
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
