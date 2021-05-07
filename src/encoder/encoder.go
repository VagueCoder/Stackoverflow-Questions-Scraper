package encoder

import (
	"fmt"
	"io"
	"log"
	"regexp"
	"sync"
	"time"

	json "github.com/json-iterator/go"
)

type FormattedTime string

func (f *FormattedTime) MarshalJSON() ([]byte, error) {
	pattern := regexp.MustCompile(`\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}`)
	matchedBytes := pattern.Find([]byte(fmt.Sprint(*f)))
	if len(matchedBytes) == 0 {
		return []byte(""), fmt.Errorf("Error at Regex Match: Couldn't find time string in pattern XXXX-XX-XX XX:XX:XX")
	}
	t, err := time.Parse("2006-01-02 15:04:05", string(matchedBytes))
	if err != nil {
		return []byte(""), fmt.Errorf("Error at FormattedTime Marshal: %v", err)
	}
	timeString := fmt.Sprintf("%q", t.Format("02-Jan-2006 15:04:05"))

	return []byte(timeString), nil
}

var err error

type Encoder struct {
	*json.Encoder

	Logger *log.Logger
	WG     *sync.WaitGroup
	mu     *sync.Mutex
}

func NewJSONEncoder(wr io.Writer, l *log.Logger) *Encoder {
	return &Encoder{
		Encoder: json.NewEncoder(wr),
		Logger:  l,
		WG:      &sync.WaitGroup{},

		// Unexported
		mu: &sync.Mutex{},
	}
}

func (en *Encoder) Encode(q interface{}) {
	en.mu.Lock()
	err = en.Encoder.Encode(q)
	en.mu.Unlock()

	if err != nil {
		en.Logger.Printf("Encoding Error: %v\n", err)
	}

	en.WG.Done()
}
