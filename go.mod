module VISOR_Server

// Keep it on 1.20, so that it can be compiled for Windows 7 too if it's compiled with Go 1.20 (it's the last version
// supporting it).
go 1.20

require github.com/cention-sany/utf7 v0.0.0-20170124080048-26cad61bd60a

require golang.org/x/text v0.14.0 // indirect
