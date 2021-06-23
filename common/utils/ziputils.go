package utils

import (
	"archive/zip"
	"github.com/pkg/errors"
	"io"
	"os"
	"path"
)

func ZipCompress(sourceFile, archiveFile string) error {
	archive, err := os.Create(archiveFile)
	if err != nil {
		return errors.Wrap(err, "create archive file error")
	}
	defer archive.Close()

	zipWriter := zip.NewWriter(archive)
	defer zipWriter.Close()

	source, err := os.Open(sourceFile)
	if err != nil {
		return errors.Wrap(err, "open sourceFile file error")
	}
	defer source.Close()

	innerName := path.Base(sourceFile)
	writer, err := zipWriter.Create(innerName)
	if err != nil {
		return errors.Wrap(err, "open sourceFile file error")
	}
	if _, err := io.Copy(writer, source); err != nil {
		return errors.Wrap(err, "archive file error")
	}

	return nil
}
