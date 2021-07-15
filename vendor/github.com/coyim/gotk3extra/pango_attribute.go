package gotk3extra

// #include <gtk/gtk.h>
import "C"
import (
	"unsafe"

	"github.com/gotk3/gotk3/pango"
)

// PangoAttribute contains information an attribute that applies to a section of text.
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

// SetStartIndex sets the index of the start of the attribute application in the text.
func (v *PangoAttribute) SetStartIndex(v2 uint) {
	v.pangoAttribute.start_index = C.guint(v2)
}

// SetEndIndex the index of the end of the attribute application in the text.
func (v *PangoAttribute) SetEndIndex(v2 uint) {
	v.pangoAttribute.end_index = C.guint(v2)
}

// PangoAttributeFromReal takes a "real" pango attribute and returns our representation of that attribute.
func PangoAttributeFromReal(v2 *pango.Attribute) *PangoAttribute {
	return &PangoAttribute{(*C.PangoAttribute)(unsafe.Pointer(v2.Native()))}
}
