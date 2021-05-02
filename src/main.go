package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type FormattedTime string

type Question struct {
	Qsn          string        `json:"question"`
	NoOfAns      string        `json:"no_of_answers"`
	URL          string        `json:"url"`
	Time         FormattedTime `json:"time"`
	RelativeTime string        `json:"relative_time"`
}

func (f *FormattedTime) MarshalJSON() ([]byte, error) {
	t, err := time.Parse("2006-01-02 15:04:05Z", fmt.Sprint(*f))
	if err != nil {
		return []byte(""), fmt.Errorf("Error at FormattedTime Marshal: %v", err)
	}
	timeString := fmt.Sprintf("%q", t.Format("02-Jan-2006 15:04:05"))

	return []byte(timeString), nil
}

var (
	ok  bool
	err error
	res *http.Response
)

func main() {
	url := "https://stackoverflow.com/questions/tagged/python?sort=Newest&filters=NoAnswers,NoAcceptedAnswer&edited=true"
	// url := "https://stackoverflow.com/questions/tagged/go?sort=Newest&filters=NoAnswers,NoAcceptedAnswer&edited=true"

	res, err = http.Get(url)
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

	var questions []Question

	doc.Find("div#questions div.mln24").Each(func(i int, div *goquery.Selection) {
		var q Question
		q.NoOfAns = strings.TrimSpace(div.Find("div.status strong").Text())
		qsndiv := div.Find("a.question-hyperlink")
		q.Qsn = qsndiv.Text()
		q.URL, ok = qsndiv.Attr("href")
		if !ok {
			q.URL = "NA"
		}
		q.URL = "https://stackoverflow.com" + q.URL

		timeTag := div.Find("span.relativetime")
		timeString, ok := timeTag.Attr("title")
		if !ok {
			q.Time = FormattedTime("")
		}

		q.Time = FormattedTime(timeString)

		q.RelativeTime = timeTag.Text()

		questions = append(questions, q)

	})

	for _, q := range questions {
		if err = json.NewEncoder(os.Stdout).Encode(q); err != nil {
			log.Fatalf("Failed to encode JSON: %v", err)
		}
	}
}
