package extract

import (
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Ungz struct {
	FileName string
}

func (u *Ungz) UngzFile() error {
	fileBase := filepath.Base(u.FileName)
	fileSplit := strings.Split(fileBase, ".gz")
	write := fileSplit[0]
	wfile := filepath.Join(filepath.Dir(u.FileName), write)

	of, err := os.Create(wfile)
	if err != nil {
		return err
	}

	gzipfile, err := os.Open(u.FileName)
	if err != nil {
		return err
	}

	reader, err := gzip.NewReader(gzipfile)
	if err != nil {
		return err
	}

	defer reader.Close()

	if _, err = io.Copy(of, reader); err != nil {
		return err
	}

	return nil
}

func NewUngz(file string) *Ungz {
	return &Ungz{
		FileName: file,
	}
}
