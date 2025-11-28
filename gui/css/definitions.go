package css

import (
	"bytes"
	"encoding/base64"
	"io"
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
		f.data, err = io.ReadAll(b64)
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

// _escFS returns a http.Filesystem for the embedded assets. If useLocal is true,
// the filesystem's contents are instead used.
func _escFS(useLocal bool) http.FileSystem {
	if useLocal {
		return _escLocal
	}
	return _escStatic
}

// _escDir returns a http.Filesystem for the embedded assets on a given prefix dir.
// If useLocal is true, the filesystem's contents are instead used.
func _escDir(useLocal bool, name string) http.FileSystem {
	if useLocal {
		return _escDirectory{fs: _escLocal, name: name}
	}
	return _escDirectory{fs: _escStatic, name: name}
}

// _escFSByte returns the named file from the embedded assets. If useLocal is
// true, the filesystem's contents are instead used.
func _escFSByte(useLocal bool, name string) ([]byte, error) {
	if useLocal {
		f, err := _escLocal.Open(name)
		if err != nil {
			return nil, err
		}
		b, err := io.ReadAll(f)
		f.Close()
		return b, err
	}
	f, err := _escStatic.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.data, nil
}

// _escFSMustByte is the same as _escFSByte, but panics if name is not present.
func _escFSMustByte(useLocal bool, name string) []byte {
	b, err := _escFSByte(useLocal, name)
	if err != nil {
		panic(err)
	}
	return b
}

// _escFSString is the string version of _escFSByte.
func _escFSString(useLocal bool, name string) (string, error) {
	b, err := _escFSByte(useLocal, name)
	return string(b), err
}

// _escFSMustString is the string version of _escFSMustByte.
func _escFSMustString(useLocal bool, name string) string {
	return string(_escFSMustByte(useLocal, name))
}

var _escData = map[string]*_escFile{

	"/definitions/account_notification_box.css": {
		local:   "definitions/account_notification_box.css",
		size:    27,
		modtime: 1489449600,
		compressed: `
CmJveCB7CiAgICBib3JkZXI6IG5vbmU7Cn0K
`,
	},

	"/definitions/bold_header_style.css": {
		local:   "definitions/bold_header_style.css",
		size:    60,
		modtime: 1489449600,
		compressed: `
LmJvbGQtaGVhZGVyLXN0eWxlIHsKCWZvbnQtc2l6ZTogMjAwJTsKCWZvbnQtd2VpZ2h0OiA4MDA7Cn0K
`,
	},

	"/definitions/conversation_pane_scroll_window.css": {
		local:   "definitions/conversation_pane_scroll_window.css",
		size:    96,
		modtime: 1489449600,
		compressed: `
CgpzY3JvbGxlZHdpbmRvdyB7CiAgICBib3JkZXItdG9wOiAycHggc29saWQgQGNveWltLWNvbnZlcnNh
dGlvbi1wYW5lLXNjcm9sbGVkLXdpbmRvdy1ib3JkZXI7Cn0K
`,
	},

	"/definitions/conversation_security_warning_unprotected.css": {
		local:   "definitions/conversation_security_warning_unprotected.css",
		size:    297,
		modtime: 1489449600,
		compressed: `
QGRlZmluZS1jb2xvciBjb3lpbS1jb252ZXJzYXRpb24tc2VjdXJpdHktd2FybmluZy1iYWNrZ3JvdW5k
ICNmZmYzZjM7CkBkZWZpbmUtY29sb3IgY295aW0tY29udmVyc2F0aW9uLXNlY3VyaXR5LXdhcm5pbmct
Zm9yZWdyb3VuZCAjMDAwMDAwOwoKYm94IHsKICAgIGJhY2tncm91bmQtY29sb3I6IEBjb3lpbS1jb252
ZXJzYXRpb24tc2VjdXJpdHktd2FybmluZy1iYWNrZ3JvdW5kOwogICAgY29sb3I6IEBjb3lpbS1jb252
ZXJzYXRpb24tc2VjdXJpdHktd2FybmluZy1mb3JlZ3JvdW5kOwogICAgYm9yZGVyOiAzcHg7Cn0K
`,
	},

	"/definitions/coyim.css": {
		local:   "definitions/coyim.css",
		size:    274,
		modtime: 1489449600,
		compressed: `
Ci8qIEdsb2JhbCBDU1MgZGVmaW5pdGlvbnMgZm9yIHRoZSBDb3lJTSBhcHBsaWNhdGlvbiAqLwoKLnJl
Z2lzdHJhdGlvbi5mb3JtIC5maWVsZFJlcXVpcmVkSW5zdHJ1Y3Rpb24gewogIGZvbnQtc3R5bGU6IGl0
YWxpYzsKfQoKLnJlZ2lzdHJhdGlvbi5mb3JtIC50aXRsZSB7CiAgcGFkZGluZzogMC41ZW07CiAgZm9u
dC1zaXplOiAyZW07CiAgZm9udC13ZWlnaHQ6IGJvbGQ7Cn0KCi5yZWdpc3RyYXRpb24uZm9ybSAuaW5z
dHJ1Y3Rpb25zIHsKICBmb250LXNpemU6IDEuMmVtOwp9Cg==
`,
	},

	"/definitions/dark/colors.css": {
		local:   "definitions/dark/colors.css",
		size:    291,
		modtime: 1489449600,
		compressed: `
CkBkZWZpbmUtY29sb3IgY295aW0tY29udmVyc2F0aW9uLXBhbmUtc2Nyb2xsZWQtd2luZG93LWJvcmRl
ciAjZDNkM2QzOwoKQGRlZmluZS1jb2xvciBjb3lpbS10b3Itbm90aWZpY2F0aW9uLWJhY2tncm91bmQg
I2YxZjFmMTsKQGRlZmluZS1jb2xvciBjb3lpbS10b3Itbm90aWZpY2F0aW9uLWZvcmVncm91bmQgIzAw
MDAwMDsKQGRlZmluZS1jb2xvciBjb3lpbS10b3Itbm90aWZpY2F0aW9uLWJvcmRlciAjZDNkM2QzOwoK
QGRlZmluZS1jb2xvciBjb3lpbS1zZWFyY2gtYmFyLWJhY2tncm91bmQgI2U4ZThlNzsK
`,
	},

	"/definitions/file_transfer_notification_canceled.css": {
		local:   "definitions/file_transfer_notification_canceled.css",
		size:    31,
		modtime: 1489449600,
		compressed: `
CmxhYmVsIHsKICAgIGNvbG9yOiAjY2MzNjM2Owp9Cg==
`,
	},

	"/definitions/file_transfer_notification_close_button.css": {
		local:   "definitions/file_transfer_notification_close_button.css",
		size:    55,
		modtime: 1489449600,
		compressed: `
bGFiZWwgewogICAgbWFyZ2luLXJpZ2h0OiAzcHg7CiAgICBtYXJnaW4tbGVmdDogM3B4Owp9Cg==
`,
	},

	"/definitions/file_transfer_notification_new.css": {
		local:   "definitions/file_transfer_notification_new.css",
		size:    76,
		modtime: 1489449600,
		compressed: `
Ym94IHsKICAgIGJhY2tncm91bmQtY29sb3I6ICNmZmY5ZjM7CiAgICBjb2xvcjogIzAwMDAwMDsKICAg
IGJvcmRlcjogM3B4Owp9Cg==
`,
	},

	"/definitions/file_transfer_notification_success.css": {
		local:   "definitions/file_transfer_notification_success.css",
		size:    31,
		modtime: 1489449600,
		compressed: `
CmxhYmVsIHsKICAgIGNvbG9yOiAjODlBRjhGOwp9Cg==
`,
	},

	"/definitions/light/colors.css": {
		local:   "definitions/light/colors.css",
		size:    291,
		modtime: 1489449600,
		compressed: `
CkBkZWZpbmUtY29sb3IgY295aW0tY29udmVyc2F0aW9uLXBhbmUtc2Nyb2xsZWQtd2luZG93LWJvcmRl
ciAjZDNkM2QzOwoKQGRlZmluZS1jb2xvciBjb3lpbS10b3Itbm90aWZpY2F0aW9uLWJhY2tncm91bmQg
I2YxZjFmMTsKQGRlZmluZS1jb2xvciBjb3lpbS10b3Itbm90aWZpY2F0aW9uLWZvcmVncm91bmQgIzAw
MDAwMDsKQGRlZmluZS1jb2xvciBjb3lpbS10b3Itbm90aWZpY2F0aW9uLWJvcmRlciAjZDNkM2QzOwoK
QGRlZmluZS1jb2xvciBjb3lpbS1zZWFyY2gtYmFyLWJhY2tncm91bmQgI2U4ZThlNzsK
`,
	},

	"/definitions/search_bar.css": {
		local:   "definitions/search_bar.css",
		size:    67,
		modtime: 1489449600,
		compressed: `
CnNlYXJjaGJhciB7CiAgICBiYWNrZ3JvdW5kLWNvbG9yOiBAY295aW0tc2VhcmNoLWJhci1iYWNrZ3Jv
dW5kOwp9Cg==
`,
	},

	"/definitions/search_bar_box.css": {
		local:   "definitions/search_bar_box.css",
		size:    26,
		modtime: 1489449600,
		compressed: `
Ym94IHsKICAgIGJvcmRlcjogbm9uZTsKfQo=
`,
	},

	"/definitions/search_bar_entry.css": {
		local:   "definitions/search_bar_entry.css",
		size:    32,
		modtime: 1489449600,
		compressed: `
ZW50cnkgewogICAgbWluLXdpZHRoOiAzMDBweDsKfQo=
`,
	},

	"/definitions/smp_infobar_box.css": {
		local:   "definitions/smp_infobar_box.css",
		size:    226,
		modtime: 1489449600,
		compressed: `
CkBkZWZpbmUtY29sb3IgY295aW0tc21wLWluZm9iYXItYmFja2dyb3VuZCAjZmZmM2YzOwpAZGVmaW5l
LWNvbG9yIGNveWltLXNtcC1pbmZvYmFyLWZvcmVncm91bmQgIzAwMDAwMDsKCmJveCB7CiAgICBiYWNr
Z3JvdW5kLWNvbG9yOiBAY295aW0tc21wLWluZm9iYXItYmFja2dyb3VuZDsKICAgIGNvbG9yOiBAY295
aW0tc21wLWluZm9iYXItZm9yZWdyb3VuZDsKICAgIGJvcmRlcjogMnB4Owp9Cg==
`,
	},

	"/definitions/smp_infobar_close.css": {
		local:   "definitions/smp_infobar_close.css",
		size:    105,
		modtime: 1489449600,
		compressed: `
CkBkZWZpbmUtY29sb3IgY295aW0tc21wLWluZm9iYXItY2xvc2UgI2U1ZDdkNjsKCmJveCB7CiAgICBi
YWNrZ3JvdW5kLWNvbG9yOiBAY295aW0tc21wLWluZm9iYXItY2xvc2U7Cn0K
`,
	},

	"/definitions/tor_notification_box.css": {
		local:   "definitions/tor_notification_box.css",
		size:    192,
		modtime: 1489449600,
		compressed: `
CmJveCB7CiAgICBiYWNrZ3JvdW5kLWNvbG9yOiBAY295aW0tdG9yLW5vdGlmaWNhdGlvbi1iYWNrZ3Jv
dW5kOwogICAgY29sb3I6IEBjb3lpbS10b3Itbm90aWZpY2F0aW9uLWZvcmVncm91bmQ7CiAgICBib3Jk
ZXI6IDFweCBzb2xpZCBAY295aW0tdG9yLW5vdGlmaWNhdGlvbi1ib3JkZXI7CiAgICBib3JkZXItcmFk
aXVzOiAycHg7Cn0K
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

	"/definitions/dark": {
		isDir: true,
		local: "definitions/dark",
	},

	"/definitions/light": {
		isDir: true,
		local: "definitions/light",
	},
}
