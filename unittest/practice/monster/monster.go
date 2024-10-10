package monster

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Monster struct {
	Name  string `json:"name"`
	Skill string `json:"skill"`
}

func (m *Monster) Store(c string) error {
	b, err0 := json.Marshal(m)
	if err0 != nil {
		return err0
	}

	f1, err1 := os.OpenFile(c, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777)
	if err1 != nil {
		return err1
	}

	defer f1.Close()

	f2 := bufio.NewWriter(f1)
	_, err2 := f2.WriteString(string(b) + "\n")
	if err2 != nil {
		return err2
	}
	f2.Flush()

	return nil
}

func (m *Monster) ReStore(c string) error {
	f1, err1 := os.OpenFile(c, os.O_RDONLY, 0777)
	if err1 != nil {
		return err1
	}
	defer f1.Close()

	f2 := bufio.NewReader(f1)
	for {
		f3, err3 := f2.ReadString('\n')
		if err3 == io.EOF {
			break
		}
		err4 := json.Unmarshal([]byte(f3), m)
		if err4 != nil {
			return err4
		}
		fmt.Println("monster反序列化=", *m)
	}
	return nil
}
