package embedDir

import (
	"embed"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
)

//go:embed dir

var data embed.FS

func DirFuc() {
	dir, _ := data.ReadDir("dir")
	for _, v := range dir {

		fmt.Println(v.Name(), filepath.Join("dir", v.Name()))
		file, err := os.ReadFile(filepath.Join("dir", v.Name()))
		fmt.Println(string(file), err)

	}
}
