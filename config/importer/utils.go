package importer

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func composeProxyString(tp, user, pass, host, port string) string {
	authPrefix := ""

	if user != "" {
		authPrefix = user
		if pass != "" {
			authPrefix = authPrefix + ":" + pass
		}
		authPrefix = authPrefix + "@"
	}

	if tp == "tor" {
		tp = "socks5"
	}

	main := tp + "://" + authPrefix + host
	if port != "" {
		main = main + ":" + port
	}

	return main
}

// ifExists checks if the file in question exists
// and if it does, adds it to the argument
func ifExists(fs []string, f string) []string {
	if fi, err := os.Stat(f); err == nil && !fi.IsDir() {
		return append(fs, f)
	}
	return fs
}

// ifExistsDir will see if argument `d` is a directory
// if it is, it will add all files inside the directory
// to the result and return that
func ifExistsDir(fs []string, d string) []string {
	if fi, err := os.Stat(d); err == nil && fi.IsDir() {
		entries, err := ioutil.ReadDir(d)
		if err == nil {
			for _, e := range entries {
				if !e.IsDir() {
					fs = append(fs, filepath.Join(d, e.Name()))
				}
			}
		}
	}
	return fs
}
