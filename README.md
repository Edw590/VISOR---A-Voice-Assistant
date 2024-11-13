# V.I.S.O.R. - A Voice Assistant
V.I.S.O.R., my in-development assistant, written in Go

## Notice
This project is a part of a bigger project, consisting of the following:
-   [V.I.S.O.R. - Android Version Assistant](https://github.com/DADi590/VISOR---A-better-Android-assistant)
-   [V.I.S.O.R. - A Voice Assistant](https://github.com/Edw590/VISOR---A-Voice-Assistant)

## Video demo
[![VISOR demo video](https://img.youtube.com/vi/GNribewvOi4/0.jpg)](https://www.youtube.com/watch?v=GNribewvOi4)

## Introduction
These are the desktop client and server versions of my in-development assistant, V.I.S.O.R. (the Android version is just above on the link).

The server one runs 24/7 on my Raspberry Pi. This is supposed to be VISOR's "operations center", where all 24/7 things run. As an example, the RSS Feed Notifier and Email Sender modules. The notifier checks for news on the feeds and the sender sends the emails that the notifier queues on it. All always running.

## Questions
Feel free to create an Issue or a Discussion with any questions you have about this. I'm trying to make it as generic as possible for anyone to be able to use it, but I may forget to document things or something, so if you need anything, feel free to ask about it.

## Usage
- Download the latest Release from [here](https://github.com/Edw590/VISOR---A-Voice-Assistant/releases).
- To use the server program (very much recommended, to get all the features working), open port 3234 on your router so that the client-server communication can be made. Also, VISOR needs an email to send emails (I used a Gmail accounted created specifically for him).
- In case you're running the server, you'll also need to generate an SSL certificate (can be self-signed). To generate a self-signed one, execute this command on either Linux or Windows: `openssl req -x509 -newkey rsa:4096 -sha256 -keyout certificate.key -out certificate.crt -subj "/CN=Common Name" -days 600 -nodes` (write the number of days you want the certificate valid for. I've left there 600 as I saw where I copied this command from), and input the path to those 2 files on the "WebsiteBackend" section inside UserSettings_EOG.dat (open with Notepad).
- Check below the requirements for each module to work. You must have a few programs installed on your computer for VISOR to work completely.
- Download a GGUF LLM model, like [this one](https://huggingface.co/QuantFactory/Llama-3.2-3B-Instruct-abliterated-GGUF/tree/main) to use on the Communicator screen settings to be VISOR's brain.
- Start the client and the server executables and that's it. 

### - Requirements for some modules to work
Check on the modules list if they work for your operating system first!

Note: if you don't know how to put something on the PATH, just copy the program to C:\Windows (Windows) or /usr/local/bin (Linux).

| Module                     | Requirement(s)
|----------------------------|-|
| Email Sender               | The `curl` program must be on the PATH.
| GPT Communicator           | The `llama-cli` ([this version is the test one](https://github.com/ggerganov/llama.cpp/releases/tag/b3880)) program must be on the PATH.
| Online Information Checker | The `chromedriver` program must be on the PATH.
| S.M.A.R.T. Checker         | The `smartctl` program must be on the PATH.
| Speech                     | For Linux, the `festival` program must be on the PATH. For Windows, it's recommended that the [`nircmdc`](https://www.nirsoft.net/utils/nircmd.html) program be on the PATH.
| System Checker             | For Linux, the `amixer` program must be on the PATH.

### Supported OSes
- Unix-like systems except MacOS (I don't test the client very much on it though, only the server)
- Windows (but I don't test the server on it, only the client)

**NOTE:** if you want full VISOR functionality on Windows with the server, run it in WSL (Win10+ only) and not natively on Windows. It will work just fine. If you do this, copy the GGUF file into WSL or it will load VERY slowly.

## Modules
| Number | Name                           | Client/Server         | Description |
|-------|--------------------------------|-----------------------|-|
| N/A   | Utils                          | Both                  | This is not a module, it's just a separate global package (all non-library modules are `main` packages; this one is the `Utils` package), but I'm writing about it because it can (and does) have utilities to communicate between modules. For example, it has utilities to queue emails to the Email Sender, so one just needs to call the function there to queue an email.
| 0     | **V.I.S.O.R.**                 | Both                  | These are the main VISOR programs (includes the server app and the client app in the same number), and are in this repository. Through here all communication are made with him.
| 1     | **Modules Manager**            | Both                  | Manages all of VISOR's modules. It's responsible for keeping them running all the time and restarting them in case they stop for any reason.
| 2     | **S.M.A.R.T. Checker**         | Server                | Runs S.M.A.R.T. tests on the given disks and checks the S.M.A.R.T. information after the tests are done.
| 3     | **Speech**                     | Client                | This is VISOR's speech module. It speaks or queues a notification about anything that needs to be spoken.
| 4     | **RSS Feed Notifier**          | Server                | Checks RSS feeds and queues an email about any news. Currently it's tested on YouTube channels *and playlists* (something YouTube didn't do nor does), and on StackExchange feeds. May work in others, but I didn't test (haven't needed so far).
| 5     | **Email Sender**               | Server                | Sends the emails that are queued for it to send. It works with cURL (the `curl` command), so it must be installed on the system and be on the PATH. It also works by sending an EML file containing the raw email information.
| 6     | **Online Information Checker** | Server                | Checks the Internet for information like weather and news and updates a file with the information it got. This file can then be read by apps to get the information back, already ready for usage.
| 7     | **GPT Communicator**           | Server (Linux only)   | Sends and receives text to and from a local LLM (Large Language Model), like Llama3.2 through the llama.cpp project. Use with any Llama3.1 or 3.2 model.
| 8     | **Website Backend**            | Server                | It's the backend of VISOR's website. It is responsible for handling the requests from the frontend.
| 9     | **Tasks Executor**             | Client                | Checks tasks and warns/executes when one is triggered. The tasks are fetched from the server.
| 10    | **System Checker**             | Client                | Collects information about the system, like Wi-Fi networks and Bluetooth devices in range. Or the state of the Wi-Fi adapter and the Bluetooth adapter. Or the screen brightness. Or others. The client can then use this information to determine where the device is and if it's being used or not.
| 11    | **Speech Recognition**         | Client (Windows only) | Currently only checks if the phrase "Hey VISOR" is spoken and shows the UI, but later should be used to detect normal speech to interact with VISOR.
| 12    | **User Locator**               | Client                | Locates the user based on everything the client knows about the user (the user must configure some things first) and on system information. For example, if the phone is communicating and the user is always with the phone (the "AlwaysWith" device), then the user is near the phone whether it's being used or not. With the computer, it must be being used because the user may leave the computer and go have lunch but not the phone.
| 13    | **Commands Executor**          | Client                | Executes commands from a sentence given to it. The commands are processed by the ACD library.
| 14    | **Google Manager**             | Server                | Manages the Google account of the user. It can currently only *check* the calendar and tasks (won't modify anything).

## Developer notes
This began as a Python project (in 2020), but even using an IDE got the project confusing. So I translated it to Go because using Go solves the issues I was having (forces me to organize the code well enough to not have cyclic imports, since it won't compile if they exist, and forces the types on variables - ah, and I can mess with pointers, and that's nice).

### - To use the project
- Download this main project and the Advanced Commands Detection module (`git clone --recursive [repo link here] VISOR`). For the server, just go to the `ServerCode` folder and run the command `go build -tags=server .`. Finally move the file to the `bin` folder. Then for the client:
  - [Install Fyne](https://docs.fyne.io/started/)
  - About portaudio:
  - - On Windows in MSYS2, install the `portaudio` package (for Windows x64, the package is [this one](https://packages.msys2.org/packages/mingw-w64-x86_64-portaudio))
  - - On Debian, install the `portaudio19-dev` package
  - - In others I didn't check, but it should be similar to the Debian one I guess
  - Go to the `ClientCode` folder and run the command `go build -tags=client .` (it will take some time the first time)
  - Move the file to the `bin` folder
- Check the Usage section above to know what to do next (with the User Settings DAT file, rename it to end with _EOG.dat ("Exclude Only from Git")).

To change VISOR to run on Windows or Linux, just compile to the OS you want and put the binaries in the `bin` folder. Nothing else needs to be done to change things from running on either OS.

## About
### - License
This project is licensed under Apache 2.0 License - http://www.apache.org/licenses/LICENSE-2.0.
