package filetransfer

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func pack(dir string, zf *os.File) error {
	a := zip.NewWriter(zf)
	defer closeAndIgnore(a)

	baseDir := filepath.Base(dir)

	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, dir))

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := a.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		/* #nosec G304 */
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer closeAndIgnore(file)
		_, err = io.Copy(writer, file)
		return err
	})
}

func unpack(file string, intoDir string) error {
	reader, err := zip.OpenReader(file)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(intoDir, 0750); err != nil {
		return err
	}

	for _, file := range reader.File {
		path := filepath.Join(intoDir, file.Name)
		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(path, file.Mode()); err != nil {
				return err
			}
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer closeAndIgnore(fileReader)

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer closeAndIgnore(targetFile)

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}

	return nil
}
