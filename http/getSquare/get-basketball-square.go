package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	//url := "https://apis.map.qq.com/ws/place/v1/suggestion?region=深圳&keyword=篮球场&key=YSRBZ-GSVY3-3P23L-RNWCE-OQB3V-T6BXG&page_size=20&page_index=1"
	count := 1

	for {
		if count <= 23 {
			break
		}

		url := fmt.Sprintf("https://restapi.amap.com/v3/place/text?key=9eda726e1f541a717c96e3a0f14346ed&keywords=篮球场&types=&city=深圳市&children=1&offset=20&page=%d&extensions=all", count)
		resp, err := http.Get(url)
		if err != nil {
			log.Println(err)
			return
		}

		all, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Println(string(all))
		time.Sleep(time.Second * time.Duration(6))
		count++

	}
}
