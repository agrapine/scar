package main

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

func scar(url string, in chan string, out chan bool) {
	res, e := http.Get(url)
	defer func() {
		out <- true
	}()
	if e != nil {
		return
	}
	body := res.Body
	defer body.Close() //i wish i had this in c#... using will do

	tokens := html.NewTokenizer(body)
	for {
		tokenType := tokens.Next()

		switch {
		case tokenType == html.ErrorToken:
			return
		case tokenType == html.StartTagToken:
			token := tokens.Token()
			if token.Data != "a" {
				continue
			}
			href, ok := getHref(token)
			if !ok {
				continue
			}
			if strings.HasPrefix(href, "http") {
				in <- href
			}
		}
	}
}

func getHref(token html.Token) (href string, ok bool) {
	for _, a := range token.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
	}
	return
}

func main() {
	links := make(map[string]bool)
	seed := make([]string, 0)
	seed = append(seed, "https://cuvva.com/")

	in := make(chan string)
	out := make(chan bool)
	for _, url := range seed {
		go scar(url, in, out)
	}

	for c := 0; c < len(seed); {
		select { //this is pretty cool
		case url := <-in:
			links[url] = true
		case <-out:
			c++
		}
	}

	for link := range links {
		fmt.Println(" - " + link)
	}

	close(in)
}
