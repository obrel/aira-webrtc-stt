# AIRA WebRTC STT

## WebRTC Speech-to-text Implementation

This repo contains WebRTC implementation that connect to speech-to-text service.
This is a part of AIRA, WebRTC AI Voice Agent project that combines WebRTC, speech-to-text, AI, and text-to-speech. We devide this project into small part as a learning material.

## Architecture
TODO

## Requirements
- Golang 1.23 above
- Google Credentials (with Speech-to-text API enabled)

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
