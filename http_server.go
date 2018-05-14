package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/ccding/go-config-reader/config"
)

func _check(err error) {
	if err != nil {
		log.Println(err.Error())
	}
}

var imgRegexp = regexp.MustCompile("http(.*?)(jpg|jpeg|png)")
var priceRegexp = regexp.MustCompile("[0-9.]+")

type meta struct {
	Title string  `json:"title"`
	Img   string  `json:"image"`
	Price float64 `json:"price"`
}

type answer struct {
	URL  string `json:"url"`
	Meta meta   `json:"meta"`
}

var cntWorkers int

func getConfig(section string, key string) string {
	c := config.NewConfig("server.conf")
	err := c.Read()
	_check(err)
	return c.Get(section, key)
}

func handler(rw http.ResponseWriter, req *http.Request) {
	log.Println("Start handler")
	var requestURL []string
	result := make(chan answer)
	worklist := make(chan []string)
	unseen := make(chan string)

	type tstruct struct {
		Test []string
	}

	countURL := 0
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&requestURL)

	if err != nil {
		panic(err)
	}
	countURL += len(requestURL)
	go func() {
		worklist <- requestURL[:]
		close(worklist)
	}()

	for i := 0; i < cntWorkers; i++ {
		go func(i int) {
			for link := range unseen {
				parseResult := parseURL(link)
				go func() {
					result <- parseResult
				}()
				log.Println(link, i)
			}
		}(i)
	}

	seen := make(map[string]bool)
	for list := range worklist {
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				unseen <- link
			} else {
				countURL--
			}
		}
	}

	receive := 0
	response := make([]answer, 0)
	for receive < countURL {
		receive++
		r := <-result
		response = append(response, r)
	}
	close(result)
	log.Println(response)
	jsonData, err := json.MarshalIndent(response, "", " ")
	_check(err)
	fmt.Fprintln(rw, string(jsonData))
	log.Println("Stop handler")
}

func main() {
	log.Println("Start server...")
	port := getConfig("SERVER", "port")
	cntWorkers, _ = strconv.Atoi(getConfig("SERVER", "workers"))
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func parseURL(url string) answer {
	answer := answer{URL: url}

	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Println(err.Error())
		return answer
	}
	log.Printf("doc=%v\n", doc)
	doc.Find("span").Each(func(i int, s *goquery.Selection) {
		if s.HasClass("offer-price") {
			answer.Meta.Price, _ = strconv.ParseFloat(priceRegexp.FindString(s.Text()), 64)
		}
	})
	log.Printf("answer=%v\n", answer)

	rawTitle := doc.Find("span#productTitle").Text()
	answer.Meta.Title = strings.TrimSpace(rawTitle)
	img, _ := doc.Find("img#imgBlkFront").Attr("data-a-dynamic-image")
	answer.Meta.Img = imgRegexp.FindString(img)
	return answer
}
