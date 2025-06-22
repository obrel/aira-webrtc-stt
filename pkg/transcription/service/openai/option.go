package openai

import "github.com/obrel/aira-websocket-stt/pkg/transcription"

func ApiKey(s string) transcription.Option {
	return func(o *OpenAI) {
		o.apiKey = s
	}
}

func Model(s string) transcription.Option {
	return func(o *OpenAI) {
		o.model = s
	}
}

func Language(s string) transcription.Option {
	return func(o *OpenAI) {
		o.language = s
	}
}

func Encoding(s string) transcription.Option {
	return func(o *OpenAI) {
		o.encoding = s
	}
}

func SampleRate(i int) transcription.Option {
	return func(o *OpenAI) {
		o.sampleRate = i
	}
}

func Prompt(s string) transcription.Option {
	return func(o *OpenAI) {
		o.prompt = s
	}
}

func DetectionType(s string) transcription.Option {
	return func(o *OpenAI) {
		o.detectionType = s
	}
}

func DetectionThreshold(f float32) transcription.Option {
	return func(o *OpenAI) {
		o.detectionTheshold = f
	}
}

func DetectionPrefixPadding(i int) transcription.Option {
	return func(o *OpenAI) {
		o.detectionPrefixPadding = i
	}
}

func DetectionSilenceDuration(i int) transcription.Option {
	return func(o *OpenAI) {
		o.detectionPrefixPadding = i
	}
}

func NoiseReductionType(s string) transcription.Option {
	return func(o *OpenAI) {
		o.noiseReductionType = s
	}
}

func (o *OpenAI) SetDefaultOptions() {
	o.detectionType = "server_vad"
	o.detectionTheshold = 0.5
	o.detectionPrefixPadding = 300
	o.detectionSilenceDuration = 500
	o.noiseReductionType = "near_field"
}

func (o *OpenAI) GetSampleRate() int {
	return o.sampleRate
}
