# V.I.S.O.R. - Server Version Assistant
The server version of my in-development assistant

## Notice
This project is a part of a bigger project, consisting of the following:
-   [V.I.S.O.R. - Android Version Assistant](https://github.com/DADi590/VISOR---A-better-Android-assistant)
-   [V.I.S.O.R. - Server Version Assistant](https://github.com/Edw590/VISOR---Server-Version-Assistant)

## Introduction
This is the server version of my in-development assistant, V.I.S.O.R.. When it's ready, it'll be running 24/7 on my Raspberry Pi. For now, the modules it has are ran automatically by the OS on the RPi (OpenMediaVault) every some minutes and close after finishing the tasks (the infinite loop is disabled on them).

This is supposed to be VISOR's "operations center", where all 24/7 things run. As an example, the RSS Feed Notifier and Email Sender modules. The notifier checks for news on the feeds and the sender sends the emails that the notifier queues on it. All always running.

## Modules
- **[Utils](https://github.com/Edw590/VISOR-Utils)** – This is not a module, it's just a separate global package (all non-library modules are `main` packages; this one is the `Utils` package), but I'm writing about it because it can (and does) have utilities to communicate between modules. For example, it has utilities to queue emails to the Email Sender, so one just needs to call the function there to queue an email.
- **[Advanced Commands Detection](https://github.com/DADi590/Advanced-Commands-Detection)** (library module) --> Detects commands in a sentence of words (a link because it's in another repository). It's the module that understands the user communication (voice or text - as long as it uses words). It can detect no so simple sentences of multiple commands, and understands the meaning of "don't", "it", and "and". Example of a complex sentence it can successfully understand (without the punctuation - it must not be present): `"turn it on. turn on the wifi, and and the airplane mode, get it it on. no, don't turn it on. turn off airplane mode and also the wifi, please."` (ignores/warns about the meaningless "it", turns on the Wi-Fi, and turns off the airplane mode and the Wi-Fi).
- **[Email Sender](https://github.com/Edw590/VISOR-EmailSender)** – Sends the emails that are queued for it to send. It works with cURL (the `curl` command), so it must be installed on the system and be on the PATH. It also works by sending an EML file containing the raw email information.
- **[RSS Feed Notifier](https://github.com/Edw590/VISOR-RssFeedNotifier)** – Checks RSS feeds and queues an email about any news. Currently it's tested on YouTube channels *and playlists* (something it didn't do nor does), and on StackExchange feeds. May work in others, but I didn't test (haven't needed so far).

## Developer notes
This began as a Python project (in 2020), but even using an IDE got the project confusing. So I'm translating it to Go, because using Go solves the issues I was having (forces me to organize the code well enough to not have cyclic imports, since it won't compile if they exist, and forces the types on variables - ah, and I can mess with pointers, and that's nice).

### - To compile the project
Download this main project and the module(s) you want to compile, and in each module folder compile the whole directory (all modules are `main` packages) with `go build .`.

All module are separated from each other *from the compiler's point of view*. So one can download each one and compile it without having any other module present - but they still require the `data` folder, especially the `ProgramData` one, with everything of all modules inside it. This is because each module can require files from other modules when they run (for example, the RSS Feed Notifier requires the Email Sender to be present because it gets files from the latter's folders to generate the emails to be queued).

#### Supported OSes
The entire project is supposed to be able to be ran on Unix-like and Windows OSes (multi-platform project). If by chance any module is not supported on any operating system, it will refuse to run on the unsupported OS(es) - even though it can probably still be compiled for them (just not ran). In case there is a module like this, it will be warned on the Modules list above.

To change it to run on Windows or Linux (aside from compiling with the Linux or Windows flag), go on the `bin` folder, put all the executable files there and check the JSON file. Constants on it must be modified (including on the first run, to set them since they are empty).

## About
### - License
This project is licensed under Apache 2.0 License - http://www.apache.org/licenses/LICENSE-2.0.
