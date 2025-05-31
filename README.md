# AIRA WebRTC STT

## WebRTC Speech-to-text Implementation

This repo contains WebRTC implementation that connect to speech-to-text service.
This is a part of AIRA, WebRTC AI Voice Agent project that combines WebRTC, speech-to-text, AI, and text-to-speech. We devide this project into small part as a learning material.

## Architecture
TODO

## External Services
This projects requires some external services to provide Speech-to-text capability. Currently it supports
- Google STT V1
- Deepgram STT
- OpenAI STT

Feel free to add more STT services.

## Requirements
- Golang 1.23 above
- Service credentials (Google, Deepgram, or OpenAI)

## Running Locally
```
$ git clone git@github.com:obrel/aira-webrtc-stt.git
$ cd aira-webrtc-stt
$ go mod tidy
$ go run main.go --google-creds <credentials.json>
```

Open your browser, and go to http://localhost:4000. Click start, start talking, and you'll be receive the transcription text.

## TO DO
Currently we only support Google Speech-to-text V1. But we want to provide another STT service to make it flexible.

## LICENSE
MIT
