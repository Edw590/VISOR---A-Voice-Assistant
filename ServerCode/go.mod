module VISOR_Server

// Keep it on 1.20, so that it can be compiled for Windows 7 too if it's compiled with Go 1.20 (it's the last version
// supporting it).
go 1.20

require (
	github.com/spf13/pflag v1.0.7
)

require github.com/mattn/go-runewidth v0.0.9 // indirect
