package deepgram

import "github.com/obrel/aira-websocket-stt/pkg/transcription"

func ApiKey(s string) transcription.Option {
	return func(d *Deepgram) {
		d.apiKey = s
	}
}

func Model(s string) transcription.Option {
	return func(d *Deepgram) {
		d.model = s
	}
}

func Language(s string) transcription.Option {
	return func(d *Deepgram) {
		d.language = s
	}
}

func Encoding(s string) transcription.Option {
	return func(d *Deepgram) {
		d.encoding = s
	}
}

func SampleRate(i int) transcription.Option {
	return func(d *Deepgram) {
		d.sampleRate = i
	}
}

func (d *Deepgram) GetSampleRate() int {
	return d.sampleRate
}
