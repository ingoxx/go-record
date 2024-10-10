package unittest02 //如果都在同一个包里也就是在同一个目录下的文件都不需要import，这里就不需要再引用也就是import,unittest02.go文件中的Add函数

import (
	"fmt"
	"testing"
)

//查看编写的add函数是否正确
//这里不需要main也能执行的原因是：testing框架会隐藏main(),然后在import xxx_test.go,再放到main()运行
//单元测试命名必须是Test开头

//如果只想测试单个文件，一定要带上被测试的文件名，默认情况下，testing框架会扫描整个目录
//如：go test -v add_test.go unittest02.go
//测试单个方法：go test -v -test.run TestAdd
//测试整个包：go test -v

func TestAdd(t *testing.T) {
	res := add(10)
	if res != 55 {

		//输出错误信息，并停止程序
		t.Fatalf("Add 执行错误=%v\n", res)
	}

	//正确就输出日志
	t.Logf("Add 执行正确")
}

func TestHello(t *testing.T) {
	fmt.Println("hello")
}
