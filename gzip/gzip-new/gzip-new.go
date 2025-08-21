package gzip_new

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type Gzip struct {
	Src     string `json:"src"`
	Dst     string `json:"dst"`
	ZipName string `json:"-"`
}

func (gz *Gzip) Check() (err error) {
	_, err = os.Stat(gz.Src)
	if err != nil {
		return
	}

	_, err = os.Stat(gz.Dst)
	if err != nil {
		return
	}

	return

}

func (gz *Gzip) GzipFile() (err error) {
	gz.ZipName = filepath.Join(gz.Dst, filepath.Base(gz.Src)+".tar.gz")
	target, err := os.Create(gz.ZipName)
	if err != nil {
		return err
	}
	defer target.Close()

	// 创建 gzip.Writer
	gzWriter := gzip.NewWriter(target)
	defer gzWriter.Close()

	// 创建 tar.Writer
	tarWriter := tar.NewWriter(gzWriter)
	defer tarWriter.Close()

	// 遍历源目录并压缩其中的文件和子目录
	err = filepath.Walk(gz.Src, func(filePath string, info os.FileInfo, err error) error {
		// 创建 tar 文件头
		header, err := tar.FileInfoHeader(info, info.Name())
		if err != nil {
			return err
		}

		// 修改文件头中的名称，以相对路径存储
		relPath, _ := filepath.Rel(gz.Src, filePath)
		header.Name = relPath

		// 写入 tar 文件头
		if err = tarWriter.WriteHeader(header); err != nil {
			return err
		}

		// 如果是文件，复制文件内容到 tar.Writer
		if !info.IsDir() {
			file, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(tarWriter, file)
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

func (gz *Gzip) UnGzipFile() (err error) {
	fPtr, err := os.Open(gz.Src)
	if err != nil {
		return err
	}

	defer fPtr.Close()

	gzReader, _ := gzip.NewReader(fPtr)

	defer gzReader.Close()

	// 创建 tar 读取器
	tarReader := tar.NewReader(gzReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break // 已经读取完所有文件
		}

		if err != nil {
			return err
		}

		// 确定解压的文件路径
		targetFilePath := fmt.Sprintf("%s/%s", gz.Dst, header.Name)

		// 创建目录
		if header.Typeflag == tar.TypeDir {
			err = os.MkdirAll(targetFilePath, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			continue
		}

		// 创建并打开目标文件
		targetFile, err := os.Create(targetFilePath)
		if err != nil {
			return err
		}

		// 将文件内容复制到目标文件
		_, err = io.Copy(targetFile, tarReader)
		if err != nil {
			return err
		}

		// 关闭目标文件
		targetFile.Close()

	}

	return
}

func NewGzip(src, dst string) *Gzip {
	gz := &Gzip{Src: src, Dst: dst}

	err := gz.Check()
	if err != nil {
		log.Fatalln(err)
	}

	return gz
}
