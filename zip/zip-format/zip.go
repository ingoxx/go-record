package zip_format

import (
	"archive/zip"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

type Zip struct {
	Src     string `json:"src"`
	Dst     string `json:"dst"`
	ZipName string `json:"-"`
}

func (z *Zip) Check() (err error) {
	_, err = os.Stat(z.Src)
	if err != nil {
		return
	}

	_, err = os.Stat(z.Dst)
	if err != nil {
		return
	}

	return
}

func (z *Zip) ZipFile() (err error) {
	z.ZipName = filepath.Join(z.Dst, filepath.Base(z.Src)+".zip")
	fc, err := os.Create(z.ZipName)
	if err != nil {
		return err
	}

	defer fc.Close()

	zipWriter := zip.NewWriter(fc)

	defer zipWriter.Close()

	err = filepath.Walk(z.Src, func(path string, info fs.FileInfo, err error) error {
		head, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		file, err := filepath.Rel(z.Src, path)
		if err != nil {
			return err
		}

		head.Name = file

		ch, err := zipWriter.CreateHeader(head)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			fo, err := os.Open(path)
			if err != nil {
				return err
			}
			_, err = io.Copy(ch, fo)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return
	}

	return
}

func (z *Zip) UnZipFile() (err error) {
	zipReader, err := zip.OpenReader(z.Src)
	if err != nil {
		return err
	}

	for _, file := range zipReader.File {
		if file.FileInfo().IsDir() {
			full := filepath.Join(z.Dst, file.Name)
			err = os.MkdirAll(full, 0775)
			if err != nil {
				return
			}
		} else {
			fullFile := filepath.Join(z.Dst, file.Name)
			fc, err := os.Create(fullFile)
			if err != nil {
				return err
			}
			defer fc.Close()

			rc, err := file.Open()
			if err != nil {
				return err
			}
			defer rc.Close()

			_, err = io.Copy(fc, rc)
			if err != nil {
				return err
			}
		}

	}

	return
}

func NewZip(src, dst string) *Zip {
	z := &Zip{Src: src, Dst: dst}

	err := z.Check()
	if err != nil {
		log.Fatalln(err)
	}

	return z
}
