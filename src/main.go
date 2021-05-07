package main

import (
	"log"
	"os"

	"github.com/VagueCoder/Stackoverflow-Questions-Scraper/src/scraper"
)

func main() {
	// url := "https://stackoverflow.com/questions/tagged/python?sort=Newest&filters=NoAnswers,NoAcceptedAnswer&edited=true"
	// url := "https://stackoverflow.com/questions/tagged/go?sort=Newest&filters=NoAnswers,NoAcceptedAnswer&edited=true"

	if len(os.Args) < 2 {
		log.Fatal("Argument error: Send page URL as argument 1.")
	}
	url := os.Args[1]

	logger := log.New(os.Stderr, "[Stackoverflow-Questions-Scraper] ", log.Lshortfile|log.LstdFlags)

	scraper.Scrape(logger, url)
}
