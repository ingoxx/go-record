package main

import "net/http"

func main() {
	url := "https://www.google.com/"
	http.Get(url)
}
