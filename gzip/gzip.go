package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	sourceDir := "C:\\Users\\Administrator\\Desktop\\update" // 源目录路径
	gzName := filepath.Base(sourceDir) + ".tar.gz"
	targetFile := filepath.Join("C:\\Users\\Administrator\\Desktop", gzName) // 压缩后的目标文件名
	//extractDir := "C:\\Users\\Administrator\\Desktop\\test"
	err := gzipFile(sourceDir, targetFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	//unGzipFile(targetFile, extractDir)

}

func gzipFile(sourceDir, targetFile string) (err error) {
	// 创建目标文件
	target, err := os.Create(targetFile)
	if err != nil {
		fmt.Println("无法创建目标文件:", err)
		return
	}
	defer target.Close()

	// 创建 gzip.Writer
	gzWriter := gzip.NewWriter(target)
	defer gzWriter.Close()

	// 创建 tar.Writer
	tarWriter := tar.NewWriter(gzWriter)
	defer tarWriter.Close()

	// 遍历源目录并压缩其中的文件和子目录
	err = filepath.Walk(sourceDir, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 创建 tar 文件头
		header, err := tar.FileInfoHeader(info, info.Name())
		if err != nil {
			return err
		}

		// 修改文件头中的名称，以相对路径存储
		relPath, _ := filepath.Rel(sourceDir, filePath)
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
		fmt.Println("压缩目录时出错:", err)
		return
	}

	fmt.Println("目录已成功压缩到", targetFile)
	return
}

func unGzipFile(file, extractDir string) (err error) {
	fPtr, _ := os.Open(file)

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
			fmt.Println("读取 tar 文件头时出错：", err)
			os.Exit(1)
		}

		fmt.Println(header.Name)

		// 确定解压的文件路径
		targetFilePath := fmt.Sprintf("%s/%s", extractDir, header.Name)

		// 创建目录
		if header.Typeflag == tar.TypeDir {
			err := os.MkdirAll(targetFilePath, os.FileMode(header.Mode))
			if err != nil {
				fmt.Println("创建目录时出错：", err)
				os.Exit(1)
			}
			continue
		}

		// 创建并打开目标文件
		targetFile, err := os.Create(targetFilePath)
		if err != nil {
			fmt.Println("创建目标文件时出错：", err)
			os.Exit(1)
		}

		// 将文件内容复制到目标文件
		_, err = io.Copy(targetFile, tarReader)
		if err != nil {
			fmt.Println("解压文件内容时出错：", err)
			os.Exit(1)
		}

		// 关闭目标文件
		targetFile.Close()

	}

	fmt.Printf("%s, 解压完毕", file)

	return
}
