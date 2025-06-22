# AIRA WebRTC STT

## WebRTC Speech-to-text Implementation

This repo contains WebRTC implementation that connect to speech-to-text service.
This is a part of AIRA, WebRTC AI Voice Agent project that combines WebRTC, speech-to-text, AI, and text-to-speech. We devide this project into small part as a learning material.

## Architecture
TODO

## External Services
This projects requires some external services to provide Speech-to-text capability. Currently it supports
- [Cartesia](example/cartesia/)
- [Deepgram](example/deepgram/)
- Elevenlabs (coming soon)
- [Google](example/google/)
- [OpenAI](example/openai/)

Feel free to add more STT services.

## Requirements
- Golang 1.23 above
- Service credentials (Google, Deepgram, OpenAI, or Cartesia)

## Example
Check [example](example/) directory to find out how to use this module.

## LICENSE
MIT
