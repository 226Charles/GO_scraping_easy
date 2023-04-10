package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"sync"
	"time"
)

func main() {

	tr := &http.Transport{
		MaxIdleConns:    10,
		IdleConnTimeout: 30 * time.Second,
	}

	client := &http.Client{
		Transport: tr,
	}

	url := "https://www.freebuf.com"
	res, err := client.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatal("status code error :%d %s", res.StatusCode, res.Status)
	}

	doc, err1 := goquery.NewDocumentFromReader(res.Body)
	if err1 != nil {
		log.Fatal(err1)
	}

	titles := make(chan string, 100)
	wg := sync.WaitGroup{}

	doc.Find(".title-view a").Each(func(i int, s *goquery.Selection) {
		wg.Add(1)

		go func(s *goquery.Selection) {
			defer wg.Done()
			titleUrl, err := s.Attr("href")
			if !err {
				log.Printf("no herf attribute fo selection: %v", s)
				return
			}

			titleRes, err1 := client.Get(titleUrl)
			if err1 != nil {
				log.Printf("http.Get error :%v", err1)
				return
			}
			defer titleRes.Body.Close()

			if titleRes.StatusCode != 200 {
				log.Printf("status code error: %d %s", titleRes.StatusCode, titleRes.Status)
				return
			}

			doc, err2 := goquery.NewDocumentFromReader(titleRes.Body)
			if err2 != nil {
				log.Printf("goquery error %v", err)
				return
			}

			titleSelection := doc.Find("title")
			if titleSelection.Length() == 0 {
				log.Printf("title selection not found")
				return
			}

			title := titleSelection.Text()
			titles <- title

		}(s)
	})

	go func() {
		wg.Wait()
		close(titles)
	}()

	for title := range titles {
		fmt.Println(title)
	}

}
