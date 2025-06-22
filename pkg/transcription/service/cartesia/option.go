package cartesia

import "github.com/obrel/aira-websocket-stt/pkg/transcription"

func ApiKey(s string) transcription.Option {
	return func(c *Cartesia) {
		c.apiKey = s
	}
}

func Model(s string) transcription.Option {
	return func(c *Cartesia) {
		c.model = s
	}
}

func Language(s string) transcription.Option {
	return func(c *Cartesia) {
		c.language = s
	}
}

func Encoding(s string) transcription.Option {
	return func(c *Cartesia) {
		c.encoding = s
	}
}

func SampleRate(i int) transcription.Option {
	return func(c *Cartesia) {
		c.sampleRate = i
	}
}

func (c *Cartesia) GetSampleRate() int {
	return c.sampleRate
}
