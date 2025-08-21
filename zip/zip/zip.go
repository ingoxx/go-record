package main

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

func main() {
	sourceDir := "C:\\Users\\Administrator\\Desktop\\update"
	zipName := filepath.Join("C:\\Users\\Administrator\\Desktop\\", filepath.Base(sourceDir)+".zip")
	//targetDir := "C:\\Users\\Administrator\\Desktop\\test"

	zipFile(sourceDir, zipName)
	//unZipFile(zipName, targetDir)
}

func zipFile(src, dst string) {
	_, err := os.Stat(src)
	if err != nil {
		log.Fatalln(err)
	}

	fc, err := os.Create(dst)
	if err != nil {
		log.Fatalln(err)
	}

	defer fc.Close()

	zipWriter := zip.NewWriter(fc)

	defer zipWriter.Close()

	err = filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {

		head, err := zip.FileInfoHeader(info)
		if err != nil {
			log.Fatalln(err)
		}

		file, err := filepath.Rel(src, path)
		if err != nil {
			log.Fatalln(err)
		}

		head.Name = file

		ch, err := zipWriter.CreateHeader(head)
		if err != nil {
			log.Fatalln(err)
		}

		if !info.IsDir() {
			fo, err := os.Open(path)
			if err != nil {
				log.Fatalln(err)
			}
			io.Copy(ch, fo)
		}

		return nil
	})

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("zip %s ok\n", filepath.Base(src))
}

func unZipFile(src, dst string) {
	_, err := os.Stat(src)
	if err != nil {
		log.Fatalln(err)
	}

	zipReader, err := zip.OpenReader(src)
	if err != nil {
		log.Fatalln(err)
	}

	for _, file := range zipReader.File {
		if file.FileInfo().IsDir() {
			full := filepath.Join(dst, file.Name)
			err = os.MkdirAll(full, 0775)
			if err != nil {
				log.Fatalln(err)
			}
		} else {
			fullFile := filepath.Join(dst, file.Name)
			fc, err := os.Create(fullFile)
			if err != nil {
				log.Fatalln(err)
			}
			defer fc.Close()

			rc, err := file.Open()
			if err != nil {
				log.Fatalln(err)
			}
			defer rc.Close()

			io.Copy(fc, rc)
		}

	}

	fmt.Printf("unzip %s ok\n", filepath.Base(src))
}
