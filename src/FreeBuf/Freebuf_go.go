package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"sync"
)

// use channel
func main() {

	url := "https://www.freebuf.com/"
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatal("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	//channel
	titles := make(chan string, 100)
	wg := sync.WaitGroup{}

	doc.Find(".title-view a").Each(func(i int, s *goquery.Selection) {
		wg.Add(1)
		go func(s *goquery.Selection) {
			//fmt.Printf("%d: %s\n", i, s.Text())
			//defer wg.Done()
			defer wg.Done()
			titleUrl, err := s.Attr("href")
			if !err {
				log.Printf("no href attribute for selection: %v", s)
				return
			}
			titleRes, err3 := http.Get(titleUrl)
			if err3 != nil {
				log.Printf("http.Get error : %v", err3)
				return
			}
			defer titleRes.Body.Close()

			if titleRes.StatusCode != 200 {
				log.Fatal("status code error： %d %s", titleRes.StatusCode, titleRes.Status)
			}

			titileDoc, err1 := goquery.NewDocumentFromReader(titleRes.Body)
			if err1 != nil {
				log.Printf("goquery error : %v", err1)
				return
			}

			titleSelection := titileDoc.Find("title")
			if titleSelection.Length() == 0 {
				log.Printf("title selection not found")
				return
			}

			title := titleSelection.Text()
			//fmt.Printf("%d: %s\n", i, title)
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
