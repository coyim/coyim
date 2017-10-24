package gliba

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/coyim/gotk3adapter/glibi"
)

func init() {
	glibi.APPLICATION_FLAGS_NONE = glibi.ApplicationFlags(glib.APPLICATION_FLAGS_NONE)
	glibi.APPLICATION_IS_SERVICE = glibi.ApplicationFlags(glib.APPLICATION_IS_SERVICE)
	glibi.APPLICATION_HANDLES_OPEN = glibi.ApplicationFlags(glib.APPLICATION_HANDLES_OPEN)
	glibi.APPLICATION_HANDLES_COMMAND_LINE = glibi.ApplicationFlags(glib.APPLICATION_HANDLES_COMMAND_LINE)
	glibi.APPLICATION_SEND_ENVIRONMENT = glibi.ApplicationFlags(glib.APPLICATION_SEND_ENVIRONMENT)
	glibi.APPLICATION_NON_UNIQUE = glibi.ApplicationFlags(glib.APPLICATION_NON_UNIQUE)

	glibi.TYPE_INVALID = glibi.Type(glib.TYPE_INVALID)
	glibi.TYPE_NONE = glibi.Type(glib.TYPE_NONE)
	glibi.TYPE_INTERFACE = glibi.Type(glib.TYPE_INTERFACE)
	glibi.TYPE_CHAR = glibi.Type(glib.TYPE_CHAR)
	glibi.TYPE_UCHAR = glibi.Type(glib.TYPE_UCHAR)
	glibi.TYPE_BOOLEAN = glibi.Type(glib.TYPE_BOOLEAN)
	glibi.TYPE_INT = glibi.Type(glib.TYPE_INT)
	glibi.TYPE_UINT = glibi.Type(glib.TYPE_UINT)
	glibi.TYPE_LONG = glibi.Type(glib.TYPE_LONG)
	glibi.TYPE_ULONG = glibi.Type(glib.TYPE_ULONG)
	glibi.TYPE_INT64 = glibi.Type(glib.TYPE_INT64)
	glibi.TYPE_UINT64 = glibi.Type(glib.TYPE_UINT64)
	glibi.TYPE_ENUM = glibi.Type(glib.TYPE_ENUM)
	glibi.TYPE_FLAGS = glibi.Type(glib.TYPE_FLAGS)
	glibi.TYPE_FLOAT = glibi.Type(glib.TYPE_FLOAT)
	glibi.TYPE_DOUBLE = glibi.Type(glib.TYPE_DOUBLE)
	glibi.TYPE_STRING = glibi.Type(glib.TYPE_STRING)
	glibi.TYPE_POINTER = glibi.Type(glib.TYPE_POINTER)
	glibi.TYPE_BOXED = glibi.Type(glib.TYPE_BOXED)
	glibi.TYPE_PARAM = glibi.Type(glib.TYPE_PARAM)
	glibi.TYPE_OBJECT = glibi.Type(glib.TYPE_OBJECT)
	glibi.TYPE_VARIANT = glibi.Type(glib.TYPE_VARIANT)
}
