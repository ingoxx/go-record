package unittest02

import "testing"

//测试用例函数
func TestSub(t *testing.T) {
	res := sub(10, 6)
	if res != 4 {

		//输出错误信息，并停止程序
		t.Fatalf("sub 执行错误=%v\n", res)
	}

	//正确就输出日志
	t.Logf("sub 执行正确")
}

//测试用例函数
func TestMod(t *testing.T) {
	res := mod(10, 6)
	if res != 4 {

		//输出错误信息，并停止程序
		t.Fatalf("mod 执行错误=%v\n", res)
	}

	//正确就输出日志
	t.Logf("mod 执行正确")
}
