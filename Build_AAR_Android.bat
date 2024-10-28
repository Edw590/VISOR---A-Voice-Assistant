:: Execute with the current directory as the project root directory

:: This needs JDK 8 to be on the PATH and a maximum Android SDK installed in Android Studio that is compiled with JDK 8.

:: Only put here below the packages that need the exported names visible. The others will be added automatically (if
:: needed?).
:: Keep compressdwarf here. It's default true, but they could change it to default false, so this way it's true for
:: sure.
mkdir bin
gomobile bind^
 -target=android^
 -x^
 -v^
 -ldflags="-v -s -w -compressdwarf=true"^
 -o="bin/MainLibraries.aar"^
 "Utils/ModsFileInfo"^
 "Utils/UtilsSWA"^
 "ACD/ACD"^
 "OICComm/OICComm"^
 "GPTComm/GPTComm"^
 "SettingsSync/SettingsSync"^
 "SpeechQueue/SpeechQueue"^
 "TEHelper/TEHelper"^
 "ULHelper/ULHelper"

echo Error code: %ERRORLEVEL%
