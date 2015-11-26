// +build libidn

package xmpp

//#cgo pkg-config: libidn
//#include <stdlib.h>
//#include <string.h>
//#include <stringprep.h>
// static inline int saslprep(char* in, char** out) {
//		return stringprep_profile(in, out, "SASLprep", 0);
// }
import "C"
import "unsafe"

func normalizePassword(password string) (string, error) {
	cpass := C.CString(password)
	defer C.free(unsafe.Pointer(cpass))

	ret := C.CString("")
	ok := C.saslprep(cpass, &ret)
	if ok != C.STRINGPREP_OK {
		return "", errPasswordContainsInvalid
	}

	return C.GoString(ret), nil
}
