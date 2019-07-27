package xmpp

import (
	"io"
	"net/http"
	"os"
)

var (
	fileName string
)

func downloadFile(filepath, url string) (err error) {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	defer out.Close()

	return nil
}
