package scraper

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/VagueCoder/Stackoverflow-Questions-Scraper/src/encoder"
)

type Question struct {
	Qsn          string                `json:"question"`
	NoOfAns      string                `json:"no_of_answers"`
	URL          string                `json:"url"`
	Time         encoder.FormattedTime `json:"time"`
	RelativeTime string                `json:"relative_time"`
}

var (
	ok     bool
	err    error
	logger *log.Logger
	en     *encoder.Encoder
)

func readDetails(_ int, div *goquery.Selection) {
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
		q.Time = encoder.FormattedTime("")
	}
	q.Time = encoder.FormattedTime(timeString)

	q.RelativeTime = timeTag.Text()

	en.WG.Add(1)
	go en.Encode(&q)
}

func Scrape(logger *log.Logger, url string) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed at HTTP get: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("Invalid response. Status code %d", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalf("Failed at goquery document creation: %v", doc)
	}

	en = encoder.NewJSONEncoder(os.Stdout, logger)

	doc.Find("div#questions div.mln24").Each(readDetails)

	en.WG.Wait()
}
