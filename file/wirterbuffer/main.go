package main

import (
	"fmt"
	"os"
	"time"
)

//模拟缓冲写, 通过减少io操作来提高写的效率，跟普通写的方式比提升10倍以上的效率

var content = "aasd93j as0dasd-i 03jkasd-adkasd09I34KJ-o34jjadjaosd02j3kda0sik03jaldklsd921\n"

type WriteBuffer struct {
	fileHandle *os.File
	cache      []byte
	index      int
}

func NewWriteBuffer(fh *os.File, size int) *WriteBuffer {
	return &WriteBuffer{
		fileHandle: fh,
		cache:      make([]byte, size),
		index:      0,
	}
}

func (w *WriteBuffer) writeByte(data []byte) {
	if len(data) >= len(w.cache) { // 如果写入的数据大于设定的缓存容量就直接落盘
		w.flush()
		w.fileHandle.Write(data)
	} else {
		// 当接近设定的缓存容量就要落盘, 并不是一定要准确到底设定的缓存值才落盘, 要考虑边界问题
		if len(w.cache[:w.index])+w.index > len(w.cache) {
			w.flush()
		}
		copy(w.cache[w.index:], data) // 如果没接近设定值就继续往后追加数据,
		w.index += len(data)          // 记录追加后的cache的容量
	}
}

func (w *WriteBuffer) writeString(data string) {
	w.writeByte([]byte(data))
}

// 每次落盘后都要清空缓存
func (w *WriteBuffer) flush() {
	w.fileHandle.Write(w.cache[:w.index])
	w.index = 0
}

func commonWrite(fs *os.File) {
	for i := 0; i < 10000; i++ {
		fs.WriteString(content)
	}
}

func main() {
	start := time.Now()
	fs, _ := os.Create("C:/Users/Administrator/Desktop/a.txt")

	nw := NewWriteBuffer(fs, 4096)
	defer nw.flush() // 把内存里残留的写进去

	for i := 0; i < 10000; i++ {
		nw.writeString(content)
	}

	//commonWrite(fs)

	fmt.Println("cost time >>>", time.Since(start))
}
