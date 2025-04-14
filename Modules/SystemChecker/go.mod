module SystemChecker

// Keep it on 1.20, so that it can be compiled for Windows 7 too if it's compiled with Go 1.20 (it's the last version
// supporting it).
go 1.20

require (
	github.com/distatus/battery v0.11.0
	github.com/go-vgo/robotgo v0.110.1
	github.com/itchyny/volume-go v0.2.2
	github.com/yusufpapurcu/wmi v1.2.4
)

require (
	github.com/schollz/logger v1.2.0 // indirect
	github.com/schollz/wifiscan v1.1.2-0.20240616123334-6fc669145a0b
)

require (
	github.com/gen2brain/shm v0.0.0-20230802011745-f2460f5984f7 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/jezek/xgb v1.1.0 // indirect
	github.com/kbinani/screenshot v0.0.0-20230812210009-b87d31814237 // indirect
	github.com/lufia/plan9stats v0.0.0-20230326075908-cb1d2100619a // indirect
	github.com/lxn/win v0.0.0-20210218163916-a377121e959e // indirect
	github.com/moutend/go-wca v0.2.0 // indirect
	github.com/otiai10/gosseract v2.2.1+incompatible // indirect
	github.com/power-devops/perfstat v0.0.0-20221212215047-62379fc7944b // indirect
	github.com/robotn/xgb v0.0.0-20190912153532-2cb92d044934 // indirect
	github.com/robotn/xgbutil v0.0.0-20190912154524-c861d6f87770 // indirect
	github.com/shirou/gopsutil/v3 v3.23.8 // indirect
	github.com/shoenig/go-m1cpu v0.1.6 // indirect
	github.com/tklauser/go-sysconf v0.3.12 // indirect
	github.com/tklauser/numcpus v0.6.1 // indirect
	github.com/vcaesar/gops v0.30.2 // indirect
	github.com/vcaesar/imgo v0.40.0 // indirect
	github.com/vcaesar/keycode v0.10.1 // indirect
	github.com/vcaesar/tt v0.20.0 // indirect
	golang.org/x/image v0.24.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	howett.net/plist v1.0.0 // indirect
)
