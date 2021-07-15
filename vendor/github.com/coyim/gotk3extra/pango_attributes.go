package gotk3extra

// #include <gtk/gtk.h>
import "C"
import (
	"unsafe"

	"github.com/gotk3/gotk3/glib"
)

// PangoAttrList is our pango attribute list implementation that supports some methods
// that are missing from the pango implementation in "gotk3".
type PangoAttrList struct {
	internal *C.PangoAttrList
}

func (v *PangoAttrList) Native() uintptr {
	return uintptr(unsafe.Pointer(v.native()))
}

func (v *PangoAttrList) native() *C.PangoAttrList {
	return (*C.PangoAttrList)(unsafe.Pointer(v.internal))
}

// Insert inserts the pango attribute specified by attr into the attribute list.
func (v *PangoAttrList) Insert(attr *PangoAttribute) {
	C.pango_attr_list_insert(v.internal, attr.native())
}

// GetAttributes returns a list of attributes.
func (v *PangoAttrList) GetAttributes() []*PangoAttribute {
	orig := C.pango_attr_list_get_iterator(v.internal)
	iter := C.pango_attr_iterator_copy(orig)
	l := (*C.struct__GSList)(C.pango_attr_iterator_get_attrs(iter))
	attrs := glib.WrapSList(uintptr(unsafe.Pointer(l)))
	if attrs == nil {
		return nil
	}

	defer attrs.Free()

	ret := make([]*PangoAttribute, 0, attrs.Length())
	attrs.Foreach(func(ptr unsafe.Pointer) {
		ret = append(ret, &PangoAttribute{(*C.PangoAttribute)(ptr)})
	})

	return ret
}

// PangoAttrListNew returns a new instance of our pango attribute list implementation.
func PangoAttrListNew() *PangoAttrList {
	c := C.pango_attr_list_new()
	attrList := new(PangoAttrList)
	attrList.internal = c
	return attrList
}
