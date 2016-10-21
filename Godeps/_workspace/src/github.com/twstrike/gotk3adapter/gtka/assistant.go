package gtka

import (
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/gotk3/gotk3/gtk"
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/gtki"
)

type assistant struct {
	*window
	internal *gtk.Assistant
}

func wrapAssistantSimple(v *gtk.Assistant) *assistant {
	if v == nil {
		return nil
	}
	return &assistant{wrapWindowSimple(&v.Window), v}
}

func wrapAssistant(v *gtk.Assistant, e error) (*assistant, error) {
	return wrapAssistantSimple(v), e
}

func unwrapAssistant(v gtki.Assistant) *gtk.Assistant {
	if v == nil {
		return nil
	}
	return v.(*assistant).internal
}
