module Speech

// Keep it on 1.20, so that it can be compiled for Windows 7 too if it's compiled with Go 1.20 (it's the last version
// supporting it).
go 1.20

require (
	github.com/Edw590/sapi-go v0.0.0-20240608194156-5f813a9f8707
	github.com/go-ole/go-ole v1.3.0
)

require (
	github.com/moutend/go-wca v0.2.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
)
