package transcribe

import (
	"io"
)

type Result struct {
	Text       string  `json:"text"`
	Confidence float32 `json:"confidence"`
	Final      bool    `json:"final"`
}

type Service interface {
	CreateStream() (Stream, error)
}

type Stream interface {
	io.Writer
	io.Closer
	Recv(chan Result, chan bool) error
	Results() <-chan Result
}
