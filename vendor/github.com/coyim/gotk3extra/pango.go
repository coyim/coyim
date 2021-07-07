package gotk3extra

// #include <gtk/gtk.h>
import "C"
import (
	"unsafe"

	"github.com/gotk3/gotk3/glib"
)

// PangoAttrList description
type PangoAttrList struct {
	pangoAttrList *C.PangoAttrList
}

// Native returns a pointer to the underlying PangoAttrList.
func (v *PangoAttrList) Native() uintptr {
	return uintptr(unsafe.Pointer(v.native()))
}

func (v *PangoAttrList) native() *C.PangoAttrList {
	return (*C.PangoAttrList)(unsafe.Pointer(v.pangoAttrList))
}

// Insert description
func (v *PangoAttrList) Insert(attr *PangoAttribute) {
	C.pango_attr_list_insert(v.pangoAttrList, attr.native())
}

// PangoAttrListNew description
func PangoAttrListNew() *PangoAttrList {
	c := C.pango_attr_list_new()
	attrList := new(PangoAttrList)
	attrList.pangoAttrList = c
	return attrList
}

// PangoAttribute description
type PangoAttribute struct {
	pangoAttribute *C.PangoAttribute
}

// Native returns a pointer to the underlying PangoAttribute.
func (v *PangoAttribute) Native() uintptr {
	return uintptr(unsafe.Pointer(v.native()))
}

func (v *PangoAttribute) native() *C.PangoAttribute {
	return (*C.PangoAttribute)(unsafe.Pointer(v.pangoAttribute))
}

// StartIndex description
func (v *PangoAttribute) StartIndex(v2 uint) {
	v.pangoAttribute.start_index = C.guint(v2)
}

//EndIndex description
func (v *PangoAttribute) EndIndex(v2 uint) {
	v.pangoAttribute.end_index = C.guint(v2)
}

// PangoGetAttributesFromList description
func PangoGetAttributesFromList(attrList *PangoAttrList) []*PangoAttribute {
	orig := C.pango_attr_list_get_iterator(attrList.native())
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

// PangoSetAttributeIndexes description
func PangoSetAttributeIndexes(attrList *PangoAttrList, startIndex, endIndex int) *PangoAttrList {
	newAttrList := PangoAttrListNew()

	for _, attr := range PangoGetAttributesFromList(attrList) {
		attr.StartIndex(uint(startIndex))
		attr.EndIndex(uint(endIndex))
		newAttrList.Insert(attr)
	}

	return newAttrList
}
