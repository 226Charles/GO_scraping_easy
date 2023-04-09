package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
)

func main() {

	url := "https://www.freebuf.com/"
	res, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatal("Status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".title-view").Each(func(i int, s *goquery.Selection) {
		fmt.Printf("%d: %s\n", i, s.Text())
	})

}
