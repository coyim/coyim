package gotk3extra

// #include <gtk/gtk.h>
import "C"
import (
	"unsafe"

	"github.com/gotk3/gotk3/gtk"
)

// LabelSetAttributes sets a PangoAttrList; the attributes in the list are applied to the label text.
func LabelSetAttributes(label *gtk.Label, attributes *PangoAttrList) {
	C.gtk_label_set_attributes((*C.GtkLabel)(unsafe.Pointer(label.Native())), attributes.native())
}

// LabelGetAttributes gets the attribute list that was set on the label, if any.
func LabelGetAttributes(label *gtk.Label) *PangoAttrList {
	pangoAttrList := C.gtk_label_get_attributes((*C.GtkLabel)(unsafe.Pointer(label.Native())))
	return &PangoAttrList{(*C.PangoAttrList)(unsafe.Pointer(pangoAttrList))}
}
