# V.I.S.O.R. - Server Version Assistant
The server version of my in-development assistant

## Notice
This project is a part of a bigger project, consisting of the following:
-   [V.I.S.O.R. - A better Android assistant](https://github.com/DADi590/VISOR---A-better-Android-assistant)
-   [V.I.S.O.R. - Server Version Assistant](https://github.com/Edw590/VISOR---Server-Version-Assistant)
-   [Advanced Commands Detection](https://github.com/DADi590/Advanced-Commands-Detection)

## Introduction
This is the server version of my in-development assistant, V.I.S.O.R.. When it's ready, it'll be running 24/7 on my Raspberry Pi. For now, the modules it has are ran automatically by the OS on the RPi (OpenMediaVault) every some minutes and close after finishing the tasks (the infinite loop is disabled on them).

This is supposed to be VISOR's "operations center", where all 24/7 things run. As an example, the RSS Feed Notifier and Email Sender modules. The notifier checks for news on the feeds and the sender sends the emails that the notifier queues on it. All always running.

## Modules
- **[Utils](https://github.com/Edw590/VISOR-Utils)** – This is not a module, it's just a separate global package (all modules are `main` packages, this one is the `Utils` package), but I'm writing about it because it can (and does) have utilities to communicate between modules. For example, it has utilities to queue emails to the Email Sender, so one just needs to call the function there to queue an email.
- **[RSS Feed Notifier](https://github.com/Edw590/VISOR-RssFeedNotifier)** – Checks RSS feeds and queues an email about any news. Currently it's tested on YouTube videos *and playlists*, and on StackExchange feeds. May work in others, but I didn't test (haven't needed so far).
- **[Email Sender](https://github.com/Edw590/VISOR-EmailSender)** – Sends the emails that are queued for it to send. It works with cURL (the `curl` command), so it must be installed on the system and be on the PATH. It also works by sending an EML file containing the raw email information.

## Developer notes
This began as a Python project (in 2021 I think), but even using an IDE got the project confusing. So I'm translating it to Go, because using Go solves the issues I was having (forces me to organize the code well enough to not have cyclic imports, since it won't compile if they exist, and forces the types on variables - ah, and I can mess with pointers, and that's nice). So this is now a GoLand project.

### - To compile the project
Download this main project and the module(s) you want to compile, and in each module folder compile the whole directory (all modules are `main` packages) with `go build .`.

All module are separated from each other *from the compiler's point of view* (so one can download each one and compile it without having any other module present) - but they can require files from the other modules when they run (for example, the RSS Feed Notifier requires the Email Sender to be present because it gets files from the latter's folders to generate the emails to be queued).

#### Supported OSes
The entire project is supposed to be able to be ran on Unix-like and Windows OSes (multi-platform project). If by chance any module is not supported on any operating system, it will refuse to run on the unsupported OS(es) - even though it can probably still be compiled for them (just not ran). In case there is a module like this, it will be warned on the Modules list above.

There is just one thing that needs to be hard-coded so far and that changes for every envionment: the path to project's directory, through the VISOR_DIR constant. Aside from that, the project can be ran anywhere, supposedly (I just changed the constant from the RPi path to a Windows path and the modules worked right there without any further modifications).

## About
### - License
This project is licensed under Apache 2.0 License - http://www.apache.org/licenses/LICENSE-2.0.
