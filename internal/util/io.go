package util

import "io"

// CloseAndIgnore closes the Closer and ignores the resulting error
func CloseAndIgnore(c io.Closer) {
    _ = c.Close()
}
