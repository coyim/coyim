package gui

import "github.com/twstrike/gotk3adapter/gdki"

// shift, control, super, hyper, meta,

func hasState(evk gdki.EventKey, s gdki.ModifierType) bool {
	return evk.State()&uint(s) != 0
}

func hasShift(evk gdki.EventKey) bool {
	return hasState(evk, gdki.GDK_SHIFT_MASK)
}

// Used to indicate scroll lock or caps lock
func hasLock(evk gdki.EventKey) bool {
	return hasState(evk, gdki.GDK_LOCK_MASK)
}

func hasControl(evk gdki.EventKey) bool {
	return hasState(evk, gdki.GDK_CONTROL_MASK)
}

func hasSuper(evk gdki.EventKey) bool {
	return hasState(evk, gdki.GDK_SUPER_MASK)
}

func hasHyper(evk gdki.EventKey) bool {
	return hasState(evk, gdki.GDK_HYPER_MASK)
}

func hasMeta(evk gdki.EventKey) bool {
	return hasState(evk, gdki.GDK_META_MASK)
}

func hasControlingModifier(evk gdki.EventKey) bool {
	return hasShift(evk) ||
		hasControl(evk) ||
		hasSuper(evk) ||
		hasHyper(evk) ||
		hasMeta(evk)
}

func hasEnter(evk gdki.EventKey) bool {
	return evk.KeyVal() == gdki.KEY_Return ||
		evk.KeyVal() == gdki.KEY_KP_Enter
}
