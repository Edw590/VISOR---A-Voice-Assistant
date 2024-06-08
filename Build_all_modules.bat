cd bin

:: Client-only modules, for Windows
set GOOS=windows
set GOARCH=amd64
go build -ldflags -H=windowsgui -o VISOR.exe .\..\
go build -o MOD_1.exe .\..\Modules\MOD_1\
go build -o MOD_3.exe .\..\Modules\MOD_3\

:: Server-only modules, for Linux
set GOOS=linux
set GOARCH=arm64
go build -o MOD_1_linux .\..\Modules\MOD_1\
go build -o MOD_2_linux .\..\Modules\MOD_2\
go build -o MOD_4_linux .\..\Modules\MOD_4\
go build -o MOD_5_linux .\..\Modules\MOD_5\
go build -o MOD_6_linux .\..\Modules\MOD_6\
go build -o MOD_7_linux .\..\Modules\MOD_7\
go build -o MOD_8_linux .\..\Modules\MOD_8\
