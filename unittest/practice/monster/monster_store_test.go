package monster

import "testing"

func TestStore(t *testing.T) {
	m := Monster{
		Name:  "lxb",
		Skill: "各种球",
	}
	f := "C:/Users/Administrator/Desktop/test1.txt"

	res := (&m).Store(f)

	if res != nil {
		t.Fatalf("json序列化写入文件失败:%v", res)
	}

	t.Logf("json序列化写入文件成功")
}

func TestReStore(t *testing.T) {
	var m Monster
	f := "C:/Users/Administrator/Desktop/test1.txt"

	res := (&m).ReStore(f)

	if res != nil {
		t.Fatalf("json反序列化失败:%v", res)
	}

	if m.Name != "lxb" {
		t.Fatalf("json反序列化失败:%v", res)
	}

	t.Logf("json反序列化成功")
}
