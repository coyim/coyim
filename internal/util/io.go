package util

import "io"

// CloseAndIgnore closes the Closer and ignores the resulting error
func CloseAndIgnore(c io.Closer) {
	LogIgnoredError(c.Close(), nil, "close and ignore")
}
