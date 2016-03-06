package gdka

import (
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/gotk3/gotk3/gdk"
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/gdki"
)

func init() {
	gdki.GDK_SHIFT_MASK = gdki.ModifierType(gdk.GDK_SHIFT_MASK)
	gdki.GDK_LOCK_MASK = gdki.ModifierType(gdk.GDK_LOCK_MASK)
	gdki.GDK_CONTROL_MASK = gdki.ModifierType(gdk.GDK_CONTROL_MASK)
	gdki.GDK_MOD1_MASK = gdki.ModifierType(gdk.GDK_MOD1_MASK)
	gdki.GDK_MOD2_MASK = gdki.ModifierType(gdk.GDK_MOD2_MASK)
	gdki.GDK_MOD3_MASK = gdki.ModifierType(gdk.GDK_MOD3_MASK)
	gdki.GDK_MOD4_MASK = gdki.ModifierType(gdk.GDK_MOD4_MASK)
	gdki.GDK_MOD5_MASK = gdki.ModifierType(gdk.GDK_MOD5_MASK)
	gdki.GDK_BUTTON1_MASK = gdki.ModifierType(gdk.GDK_BUTTON1_MASK)
	gdki.GDK_BUTTON2_MASK = gdki.ModifierType(gdk.GDK_BUTTON2_MASK)
	gdki.GDK_BUTTON3_MASK = gdki.ModifierType(gdk.GDK_BUTTON3_MASK)
	gdki.GDK_BUTTON4_MASK = gdki.ModifierType(gdk.GDK_BUTTON4_MASK)
	gdki.GDK_BUTTON5_MASK = gdki.ModifierType(gdk.GDK_BUTTON5_MASK)
	gdki.GDK_SUPER_MASK = gdki.ModifierType(gdk.GDK_SUPER_MASK)
	gdki.GDK_HYPER_MASK = gdki.ModifierType(gdk.GDK_HYPER_MASK)
	gdki.GDK_META_MASK = gdki.ModifierType(gdk.GDK_META_MASK)
	gdki.GDK_RELEASE_MASK = gdki.ModifierType(gdk.GDK_RELEASE_MASK)
	gdki.GDK_MODIFIER_MASK = gdki.ModifierType(gdk.GDK_MODIFIER_MASK)
}
