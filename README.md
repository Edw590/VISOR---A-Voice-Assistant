# V.I.S.O.R. - Server Version Assistant
The server version of my in-development assistant

## Notice
This project is a part of a bigger project, consisting of the following:
-   [V.I.S.O.R. - Android Version Assistant](https://github.com/DADi590/VISOR---A-better-Android-assistant)
-   [V.I.S.O.R. - Server Version Assistant](https://github.com/Edw590/VISOR---Server-Version-Assistant)

## Introduction
This is the server version of my in-development assistant, V.I.S.O.R.. It runs 24/7 on my Raspberry Pi (but it can run on Windows too. It's supported on both). This is supposed to be VISOR's "operations center", where all 24/7 things run. As an example, the RSS Feed Notifier and Email Sender modules. The notifier checks for news on the feeds and the sender sends the emails that the notifier queues on it. All always running.

## Modules
- **[Utils](https://github.com/Edw590/VISOR-Utils)** – This is not a module, it's just a separate global package (all non-library modules are `main` packages; this one is the `Utils` package), but I'm writing about it because it can (and does) have utilities to communicate between modules. For example, it has utilities to queue emails to the Email Sender, so one just needs to call the function there to queue an email.
- **[Email Sender](https://github.com/Edw590/VISOR-EmailSender)** – Sends the emails that are queued for it to send. It works with cURL (the `curl` command), so it must be installed on the system and be on the PATH. It also works by sending an EML file containing the raw email information.
- **[GPT Communicator](https://github.com/Edw590/VISOR-GPTCommunicator)** - Sends and receives text to and from a local LLM (Large Language Model), like Llama 3 through the llama.cpp project.
- **[Modules Manager](https://github.com/Edw590/VISOR-ModulesManager)** - Manages all of VISOR's modules. It's responsible for keeping them running all the time and restarting them in case they stop for any reason.
- **[Online Information Checker](https://github.com/Edw590/VISOR-OnlineInformationChecker)** – Checks the Internet for information like weather and news and updates a file with the information it got. This file can then be read by apps to get the information back, already ready for usage.
- **[RSS Feed Notifier](https://github.com/Edw590/VISOR-RssFeedNotifier)** – Checks RSS feeds and queues an email about any news. Currently it's tested on YouTube channels *and playlists* (something YouTube didn't do nor does), and on StackExchange feeds. May work in others, but I didn't test (haven't needed so far).
- **[S.M.A.R.T. Checker](https://github.com/Edw590/VISOR-SMARTChecker)** - Runs S.M.A.R.T. tests on the given disks and checks the S.M.A.R.T. information after the tests are done.
- **[Website Backend](https://github.com/Edw590/VISOR-WebsiteBackend)** - It's the backend of VISOR's website. It is responsible for handling the requests from the frontend.

## Libraries
- **[Advanced Commands Detection](https://github.com/DADi590/Advanced-Commands-Detection)** --> Detects commands in a sentence of words (a link because it's in another repository). It's the module that understands the user communication (voice or text - as long as it uses words). It can detect no so simple sentences of multiple commands, and understands the meaning of "don't", "it", and "and". Example of a complex sentence it can successfully understand (without the punctuation - it must not be present): `"turn it on. turn on the wifi, and and the airplane mode, get it it on. no, don't turn it on. turn off airplane mode and also the wifi, please."` (ignores/warns about the meaningless "it", turns on the Wi-Fi, and turns off the airplane mode and the Wi-Fi).
- **[Online Information Getter](https://github.com/Edw590/VISOR-OnlineInformationGetter)** - Gets the information the Online Internet Checker module put on the website files.
- **[GPT](https://github.com/Edw590/VISOR-GPT)** - Sends and gets text to and from the GPT Communicator module through the Website Backend module.

## Developer notes
This began as a Python project (in 2020), but even using an IDE got the project confusing. So I'm translating it to Go, because using Go solves the issues I was having (forces me to organize the code well enough to not have cyclic imports, since it won't compile if they exist, and forces the types on variables - ah, and I can mess with pointers, and that's nice).

### - To use the project
- Download this main project and the module(s) you want to compile (to download all, `git clone --recursive [repo link here]`), and in each module folder compile the whole directory (all modules are `main` packages) with `go build .` (install Go first). After that, rename the generated executable file to the name of the folder + a suffix: MOD_1.exe for Windows or MOD_1_linux for Linux. Finally move the file to the `bin` folder.
- Next go on the `bin` folder and edit the JSON file with your values and rename the file to PersonalConsts_EOG.json. VISOR needs an email of its own btw. Also needs a website. I'll try to remove that requirement soon. But for full functionality (like communication between the app and the server) the website must exist.
- Go on each module folder and copy the JSON file to `data/UserData/MOD_[module number here]` (create the folders if they don't exist) and configure it (in case the module needs one).
- To be sure each module is supported, start each module individually (except the Manager) and see if no errors pop up.
- If no errors appear, start the Modules Manager (MOD_1), which will start and keep running all the other modules.

#### Supported OSes
The entire project is supposed to be able to be ran on Unix-like and Windows OSes (multi-platform project). If by chance any module is not supported on any operating system, it will refuse to run on the unsupported OS(es) - even though it can probably still be compiled for them (just not ran). In case there is a module like this, it will be warned on the Modules list above.

To change it to run on Windows or Linux (aside from compiling with the Linux or Windows flag), go on the `bin` folder, put all the executable files there and check the JSON file. Constants on it must be modified (including on the first run, to set them since they are empty).

## About
### - License
This project is licensed under Apache 2.0 License - http://www.apache.org/licenses/LICENSE-2.0.
