package gui

import "github.com/coyim/coyim/internal/util"

func ignore[T any](v T) {}

func orErr[T any](v T, err error) T {
	return util.OrErr(v, err, nil, "GTK setup")
}

func orErrOs[T any](v T, err error) T {
	return util.OrErr(v, err, nil, "operating system setup")
}

func ignErrGtk(err error) {
	util.LogIgnoredError(err, nil, "GTK setup")
}

func ignErrOs(err error) {
	util.LogIgnoredError(err, nil, "operating system setup")
}
