package main

import (
	"fmt"
	"syscall"
)

func main() {
	// 创建共享内存
	shmName := "/shared_memory"
	data := "Hello from Go!"

	// 打开共享内存
	fd, err := syscall.Open(shmName, syscall.O_CREAT|syscall.O_RDWR, 0666)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer syscall.Close(fd)

	// 设置共享内存大小
	syscall.Ftruncate(fd, int64(len(data)))

	// 映射到内存
	memory, err := syscall.Mmap(fd, 0, len(data), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer syscall.Munmap(memory)

	// 写入数据
	copy(memory, data)
	fmt.Println("Data written to shared memory:", data)
}
