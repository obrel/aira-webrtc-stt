package transcription

import (
	"sync"

	"github.com/obrel/go-lib/pkg/log"
)

var (
	transcriptions = map[string]Factory{}
	lock           sync.RWMutex
)

type Option interface{}

type Factory func(...Option) (Transcription, error)

type Transcription interface {
	Connect() error
	Write([]byte) (int, error)
	Receive(chan Result, chan bool) error
	Close() error
}

type Result struct {
	Text       string  `json:"text"`
	Confidence float32 `json:"confidence"`
	Final      bool    `json:"final"`
}

func NewTranscription(s string, opts ...Option) (Transcription, error) {
	lock.RLock()
	defer lock.RUnlock()

	service, ok := transcriptions[s]
	if !ok {
		log.For("transcription", "new").Fatal("Transcription not found.")
	}

	return service(opts...)
}

func Register(s string, b Factory) {
	lock.Lock()
	defer lock.Unlock()

	if b == nil {
		log.For("transcription", "register").Fatal("Invalid transcription.")
	}

	if _, ok := transcriptions[s]; ok {
		log.For("transcription", "register").Fatalf("Transcription %s already registered.", s)
	}

	transcriptions[s] = b
}
