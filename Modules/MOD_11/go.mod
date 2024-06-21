module SpeechRecognition

// Keep it on 1.20, so that it can be compiled for Windows 7 too if it's compiled with Go 1.20 (it's the last version
// supporting it).
go 1.20

require (
	github.com/Picovoice/porcupine/binding/go/v3 v3.0.2
	github.com/gordonklaus/portaudio v0.0.0-20230709114228-aafa478834f5
)
