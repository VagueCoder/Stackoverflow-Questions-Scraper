package encoder

import (
	"io"
	"log"
	"sync"

	json "github.com/json-iterator/go"
)

type Encoder struct {
	Logger *log.Logger
	Syncer syncer

	*json.Encoder
}

type syncer struct {
	wg *sync.WaitGroup
	mu *sync.Mutex
}

var err error

func NewJSONEncoder(wr io.Writer, l *log.Logger, w *sync.WaitGroup, m *sync.Mutex) *Encoder {

	en := json.NewEncoder(wr)
	return &Encoder{
		Encoder: en,
		Logger:  l,
		Syncer:  syncer{w, m},
	}
}

func (en *Encoder) Encode(q interface{}) {
	en.Syncer.mu.Lock()
	err = en.Encoder.Encode(q)
	en.Syncer.mu.Unlock()

	if err != nil {
		en.Logger.Printf("Encoding Error: %v\n", err)
	}

	en.Syncer.wg.Done()
}
