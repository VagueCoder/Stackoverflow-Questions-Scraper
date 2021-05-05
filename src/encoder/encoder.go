package encoder

import (
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	json "github.com/json-iterator/go"
)

type FormattedTime string

func (f *FormattedTime) MarshalJSON() ([]byte, error) {
	t, err := time.Parse("2006-01-02 15:04:05Z", fmt.Sprint(*f))
	if err != nil {
		return []byte(""), fmt.Errorf("Error at FormattedTime Marshal: %v", err)
	}
	timeString := fmt.Sprintf("%q", t.Format("02-Jan-2006 15:04:05"))

	return []byte(timeString), nil
}

var err error

type Encoder struct {
	Logger *log.Logger
	WG     *sync.WaitGroup
	mu     *sync.Mutex

	*json.Encoder
}

func NewJSONEncoder(wr io.Writer, l *log.Logger) *Encoder {

	en := json.NewEncoder(wr)

	return &Encoder{
		Encoder: en,
		Logger:  l,
		WG:      &sync.WaitGroup{},
		mu:      &sync.Mutex{},
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
