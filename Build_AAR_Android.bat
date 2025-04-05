:: Execute with the current directory as the project root directory

:: This needs JDK 8 to be on the PATH and a maximum Android SDK installed in Android Studio that is compiled with JDK 8.

:: Run "go install golang.org/x/mobile/cmd/gomobile@v0.0.0-20240707233753-b765e5d5218f" and
:: "go install golang.org/x/mobile/cmd/gobind@v0.0.0-20240707233753-b765e5d5218f" to install the gomobile tools.
:: Must be the version I wrote - it's the last one supporting Go 1.20 (that I found on WayBackMachine at least).

:: Only put here below the packages that need the exported names visible. The others will be added automatically (if
:: needed?).
:: Keep compressdwarf here. It's default true, but they could change it to default false, so this way it's true for
:: sure.
mkdir bin
gomobile bind^
 -tags=client^
 -target=android^
 -x^
 -v^
 -ldflags="-v -s -w -compressdwarf=true"^
 -o="bin/MainLibraries.aar"^
 "ACD/ACD"^
 "GMan/GMan"^
 "GPTComm/GPTComm"^
 "OICComm/OICComm"^
 "SCLink/SCLink"^
 "SettingsSync/SettingsSync"^
 "SpeechQueue/SpeechQueue"^
 "TEHelper/TEHelper"^
 "ULHelper/ULHelper"^
 "Utils/ModsFileInfo"^
 "Utils/UtilsSWA"

echo Error code: %ERRORLEVEL%
