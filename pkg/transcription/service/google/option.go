package google

import "github.com/obrel/aira-websocket-stt/pkg/transcription"

func Credential(s string) transcription.Option {
	return func(g *Google) {
		g.credential = s
	}
}

func Model(s string) transcription.Option {
	return func(g *Google) {
		g.model = s
	}
}

func Language(s string) transcription.Option {
	return func(g *Google) {
		g.language = s
	}
}

func Encoding(s string) transcription.Option {
	return func(g *Google) {
		g.encoding = s
	}
}

func SampleRate(i int) transcription.Option {
	return func(g *Google) {
		g.sampleRate = i
	}
}

func AudioChannel(i int) transcription.Option {
	return func(g *Google) {
		g.audioChannel = i
	}
}

func (g *Google) GetSampleRate() int {
	return g.sampleRate
}
