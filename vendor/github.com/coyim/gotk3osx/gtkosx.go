// +build darwin

package gotk3osx

// #cgo pkg-config: gtk-mac-integration-gtk3 gio-2.0 glib-2.0 gobject-2.0 gtk+-3.0 gdk-3.0
// #include <gtk/gtk.h>
// #include <gdk/gdk.h>
// #include <glib-object.h>
// #include <gtkosxapplication.h>
// #include "gtkosx.go.h"
import "C"
import "github.com/gotk3/gotk3/glib"
import "github.com/gotk3/gotk3/gdk"
import "github.com/gotk3/gotk3/gtk"
import "unsafe"
import "errors"
import "github.com/coyim/gotk3osx/access"
import "github.com/coyim/gotk3adapter/gdki"
import "github.com/coyim/gotk3adapter/gdka"
import "github.com/coyim/gotk3adapter/gtki"
import "github.com/coyim/gotk3adapter/gtka"

func init() {
	tm := []glib.TypeMarshaler{
		{glib.Type(C.gtkosx_application_get_type()), marshalGtkosxApplication},
	}
	glib.RegisterGValueMarshalers(tm)
}

func gbool(b bool) C.gboolean {
	if b {
		return C.gboolean(1)
	}
	return C.gboolean(0)
}

func gobool(b C.gboolean) bool {
	return b != C.FALSE
}

func stringReturn(c *C.gchar) string {
	if c == nil {
		return ""
	}
	return C.GoString((*C.char)(c))
}

var nilPtrErr = errors.New("cgo returned unexpected nil pointer")

// GtkosxApplication represents a GTK OSX application
type GtkosxApplication struct {
	glib.InitiallyUnowned
}

// Signals available to connect to:
// - "NSApplicationBlockTermination"
// - "NSApplicationDidBecomeActive"
// - "NSApplicationOpenFile"
// - "NSApplicationWillResignActive"
// - "NSApplicationWillTerminate"

func (v *GtkosxApplication) native() *C.GtkosxApplication {
	if v == nil || v.GObject == nil {
		return nil
	}
	p := unsafe.Pointer(v.GObject)
	return C.toGtkosxApplication(p)
}

func marshalGtkosxApplication(p uintptr) (interface{}, error) {
	c := C.g_value_get_object((*C.GValue)(unsafe.Pointer(p)))
	obj := glib.Take(unsafe.Pointer(c))
	return wrapGtkosxApplication(obj), nil
}

func wrapGtkosxApplication(obj *glib.Object) *GtkosxApplication {
	return &GtkosxApplication{glib.InitiallyUnowned{obj}}
}

// GetGtkosxApplication returns the singleton application object
func GetGtkosxApplication() (*GtkosxApplication, error) {
	c := C.gtkosx_application_get()
	if c == nil {
		return nil, nilPtrErr
	}

	return wrapGtkosxApplication(glib.Take(unsafe.Pointer(c))), nil
}

type real struct{}

// Real implements the main access
var Real = &real{}

// GetApplication implements access interfaces
func (*real) GetApplication() (access.Application, error) {
	return GetGtkosxApplication()
}

// Ready informs Cocoa that application initialization is complete.
func (v *GtkosxApplication) Ready() {
	C.gtkosx_application_ready(v.native())
}

// SetUseQuartzAccelerators sets quartz accelerator handling; TRUE (default) uses quartz; FALSE uses Gtk+. Quartz accelerator handling is required for normal OS X accelerators (e.g., command-q to quit) to work.
func (v *GtkosxApplication) SetUseQuartzAccelerators(v2 bool) {
	C.gtkosx_application_set_use_quartz_accelerators(v.native(), gbool(v2))
}

// UseQuartzAccelerators returns whether we are using Quartz or Gtk+ accelerator handling?
func (v *GtkosxApplication) UseQuartzAccelerators() bool {
	return gobool(C.gtkosx_application_use_quartz_accelerators(v.native()))
}

// SetMenuBarReal sets a window's menubar as the application menu bar. Call this once for each window as you create them. It works best if the menubar is reasonably fully populated before you call it. Once set, it will stay syncronized through signals as long as you don't disconnect or block them.
func (v *GtkosxApplication) SetMenuBarReal(menuShell *gtk.MenuShell) {
	C.gtkosx_application_set_menu_bar(v.native(), nativeMenuShell(menuShell))
}

// SyncMenuBar will syncronize the active window's GtkMenuBar with the OSX menu bar. You should only need this if you have programmatically changed the menus with signals blocked or disconnected.
func (v *GtkosxApplication) SyncMenuBar() {
	C.gtkosx_application_sync_menubar(v.native())
}

// InsertAppMenuItemReal will insert a menu item in the a app menu
func (v *GtkosxApplication) InsertAppMenuItemReal(menuItem *gtk.Widget, index int) {
	C.gtkosx_application_insert_app_menu_item(v.native(), nativeWidget(menuItem), (C.gint)(index))
}

// SetWindowMenuReal sets a designated menu item already on the menu bar as the Window menu. This is the menu which contains a list of open windows for the application; by default it also provides menu items to minimize and zoom the current window and to bring all windows to the front. Call this after gtk_osx_application_set_menu_bar(). It operates on the currently active menubar. If nenu_item is NULL, it will create a new menu for you, which will not be gettext translatable.
func (v *GtkosxApplication) SetWindowMenuReal(menuItem *gtk.MenuItem) {
	C.gtkosx_application_set_window_menu(v.native(), nativeMenuItem(menuItem))
}

// SetHelpMenuReal sets a designated menu item already on the menu bar as the Help menu. Call this after gtk_osx_application_set_menu_bar(), but before gtk_osx_application_window_menu(), especially if you're letting GtkosxApplication create a Window menu for you (it helps position the Window menu correctly). It operates on the currently active menubar. If nenu_item is NULL, it will create a new menu for you, which will not be gettext translatable.
func (v *GtkosxApplication) SetHelpMenuReal(menuItem *gtk.MenuItem) {
	C.gtkosx_application_set_help_menu(v.native(), nativeMenuItem(menuItem))
}

// SetDockMenu Set a GtkMenu as the dock menu. This menu does not have a "sync" function, so changes made while signals are disconnected will not update the menu which appears in the dock, and may produce strange results or crashes if a GtkMenuItem or GtkAction associated with a dock menu item is deallocated.
func (v *GtkosxApplication) SetDockMenu(menuItem *gtk.MenuShell) {
	C.gtkosx_application_set_dock_menu(v.native(), nativeMenuShell(menuItem))
}

// SetDockIconPixbufReal sets the dock icon from a GdkPixbuf
func (v *GtkosxApplication) SetDockIconPixbufReal(pixbuf *gdk.Pixbuf) {
	C.gtkosx_application_set_dock_icon_pixbuf(v.native(), nativePixbuf(pixbuf))
}

// SetDockIconResource sets the dock icon from an image file in the bundle/
func (v *GtkosxApplication) SetDockIconResource(name, tp, subdir string) {
	cstrname := C.CString(name)
	defer C.free(unsafe.Pointer(cstrname))

	cstrtype := C.CString(tp)
	defer C.free(unsafe.Pointer(cstrtype))

	cstrsubdir := C.CString(subdir)
	defer C.free(unsafe.Pointer(cstrsubdir))

	C.gtkosx_application_set_dock_icon_resource(v.native(), cstrname, cstrtype, cstrsubdir)
}

type GtkosxApplicationAttentionType int

const (
	ATTENTION_TYPE_CRITICAL_REQUEST GtkosxApplicationAttentionType = C.CRITICAL_REQUEST
	ATTENTION_TYPE_INFO_REQUEST     GtkosxApplicationAttentionType = C.INFO_REQUEST
)

// AttentionRequest creates an attention request. If type is CRITICAL_REQUEST, the dock icon will bounce until cancelled the application receives focus; otherwise it will bounce for 1 second -- but the attention request will remain asserted until cancelled or the application receives focus. This function has no effect if the application has focus.
func (v *GtkosxApplication) AttentionRequest(tp GtkosxApplicationAttentionType) int {
	res := C.gtkosx_application_attention_request(v.native(), C.GtkosxApplicationAttentionType(tp))
	return int(res)
}

// CancelAttentionRequest cancels an attention request created with gtkosx_application_attention_request.
func (v *GtkosxApplication) CancelAttentionRequest(id int) {
	C.gtkosx_application_cancel_attention_request(v.native(), C.gint(id))
}

// SetDockIconPixbuf implements the access interface
func (v *GtkosxApplication) SetDockIconPixbuf(icon gdki.Pixbuf) {
	v.SetDockIconPixbufReal(gdka.UnwrapPixbuf(icon))
}

// InsertAppMenuItem implements the access interface
func (v *GtkosxApplication) InsertAppMenuItem(mi gtki.Widget, ix int) {
	v.InsertAppMenuItemReal(gtka.UnwrapWidget(mi), ix)
}

// SetMenuBar implements the access interface
func (v *GtkosxApplication) SetMenuBar(m gtki.MenuShell) {
	v.SetMenuBarReal(gtka.UnwrapMenuShell(m))
}

// SetHelpMenu implements the access interface
func (v *GtkosxApplication) SetHelpMenu(m gtki.MenuItem) {
	v.SetHelpMenuReal(gtka.UnwrapMenuItem(m))
}

// SetWindowMenu implements the access interface
func (v *GtkosxApplication) SetWindowMenu(m gtki.MenuItem) {
	v.SetWindowMenuReal(gtka.UnwrapMenuItem(m))
}

// GetBundlePath returns the root path of the bundle or the directory containing the executable if it isn't actually a bundle.
func GetBundlePath() string {
	return stringReturn(C.gtkosx_application_get_bundle_path())
}

// GetResourcePath returns the Resource path for the bundle or the directory containing the executable if it isn't actually a bundle. Use gtkosx_application_get_bundle_id() to check (it will return NULL if it's not a bundle).
func GetResourcePath() string {
	return stringReturn(C.gtkosx_application_get_resource_path())
}

// GetExecutablePath returns the executable path, including file name
func GetExecutablePath() string {
	return stringReturn(C.gtkosx_application_get_executable_path())
}

// GetBundleID returns the value of the CFBundleIdentifier key from the bundle's Info.plist. This will return NULL if it's not really a bundle, there's no Info.plist, or if Info.plist doesn't have a CFBundleIdentifier key (So if you need to detect being in a bundle, make sure that your bundle has that key!)
func GetBundleID() string {
	return stringReturn(C.gtkosx_application_get_bundle_id())
}

// GetBundleInfo queries the bundle's Info.plist with key. If the returned object is a string, returns that; otherwise returns NULL.
func GetBundleInfo(key string) string {
	cstr := C.CString(key)
	defer C.free(unsafe.Pointer(cstr))
	return stringReturn(C.gtkosx_application_get_bundle_info(cstr))
}
