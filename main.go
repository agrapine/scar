package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	resp, _ := http.Get("https://www.cuvva.com/")
	bytes, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("HTML:\n\n", string(bytes))
}
