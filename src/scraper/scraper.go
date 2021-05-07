package scraper

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/VagueCoder/Stackoverflow-Questions-Scraper/src/encoder"
)

type Question struct {
	Qsn          string                `json:"question"`
	NoOfAns      string                `json:"no_of_answers"`
	URL          string                `json:"url"`
	PostedTime   encoder.FormattedTime `json:"question_posted_time"`
	RelativeTime string                `json:"relative_posted_time"`
	ScrapeTime   encoder.FormattedTime `json:"scraped_at"`
}

var (
	ok     bool
	err    error
	logger *log.Logger
	en     *encoder.Encoder
)

func Scrape(logger *log.Logger, url string) {
	res, err := http.Get(url)
	if err != nil {
		logger.Fatalf("Failed at HTTP get: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		logger.Fatalf("Invalid response. Status code %d", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		logger.Fatalf("Failed at goquery document creation: %v", doc)
	}

	en = encoder.NewJSONEncoder(os.Stdout, logger)

	doc.Find("div#questions div.mln24").Each(func(i int, div *goquery.Selection) {
		go scrapeDetails(i, div)
		en.WG.Add(1)
	})

	en.WG.Wait()
}

func scrapeDetails(_ int, div *goquery.Selection) {
	q := &Question{}

	q.NoOfAns = strings.TrimSpace(div.Find("div.status strong").Text())

	qsndiv := div.Find("a.question-hyperlink")

	q.Qsn = qsndiv.Text()

	q.URL, ok = qsndiv.Attr("href")
	if !ok {
		q.URL = "NA"
	} else {
		q.URL = "https://stackoverflow.com" + q.URL
	}

	timeTag := div.Find("span.relativetime")
	timeString, ok := timeTag.Attr("title")
	if !ok {
		q.PostedTime = encoder.FormattedTime("")
	}
	q.PostedTime = encoder.FormattedTime(timeString)

	q.RelativeTime = timeTag.Text()

	q.ScrapeTime = encoder.FormattedTime(fmt.Sprint(time.Now()))
	go en.Encode(&q)
}
