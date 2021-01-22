package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/coyim/gotk3extra"
)

func (*RealGtk) GetWidgetBuildableName(w gtki.Widget) (string, error) {
	return gotk3extra.GetWidgetBuildableName(UnwrapWidget(w))
}
