package util

import "io"

func CloseAndIgnore(c io.Closer) {
    _ = c.Close()
}
