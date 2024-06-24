module SystemState

// Keep it on 1.20, so that it can be compiled for Windows 7 too if it's compiled with Go 1.20 (it's the last version
// supporting it).
go 1.20

require (
	github.com/apaxa-go/eval v0.0.0-20171223182326-1d18b251d679
	github.com/distatus/battery v0.11.0
	github.com/itchyny/volume-go v0.2.2
	github.com/yusufpapurcu/wmi v1.2.4
)

require (
	github.com/schollz/logger v1.2.0 // indirect
	github.com/schollz/wifiscan v1.1.2-0.20240616123334-6fc669145a0b // indirect
)

require (
	github.com/apaxa-go/helper v0.0.0-20180607175117-61d31b1c31c3 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/moutend/go-wca v0.2.0 // indirect
	golang.org/x/sys v0.9.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	howett.net/plist v1.0.0 // indirect
)
