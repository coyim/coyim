// +build !gtk_3_6
// +build !gtk_3_8
// +build !gtk_3_10
// +build !gtk_3_12
// +build !gtk_3_14
// +build !gtk_3_16
// +build !gtk_3_18
// +build !gtk_3_20
// +build !gtk_3_22

package gtka

import (
	"errors"

	"github.com/coyim/gotk3adapter/gtki"
)

func (*RealGtk) CssProviderGetDefault() (gtki.CssProvider, error) {
	return nil, errors.New("css_provider_get_default is not provided anymore")
}
