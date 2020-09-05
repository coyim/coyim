// +build gtk_3_6 gtk_3_8 gtk_3_10

package gtka

import (
	"fmt"

	"github.com/coyim/gotk3adapter/gtki"
)

func (v *box) SetCenterWidget(gtki.Widget) {
	fmt.Println("WARNING - gotk3adapter - Box.SetCenterWidget() is not supported on your platform")
}

func (v *box) GetCenterWidget() gtki.Widget {
	fmt.Println("WARNING - gotk3adapter - Box.GetCenterWidget() is not supported on your platform")
	return nil
}
