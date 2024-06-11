# V.I.S.O.R. - RSS Feed Notifier
RSS Feed Notifier module of V.I.S.O.R.

This repository is a submodule on the [V.I.S.O.R. - Server Version Assistant](https://github.com/Edw590/VISOR---Server-Version-Assistant) project (the main project).

## What it does
This module checks RSS feeds and queues an email about any news (for the Email Sender module to send).

Currently it's tested on YouTube videos and playlists, and on StackExchange feeds. May work in others, but I didn't test
(haven't needed so far).

Check the `mod_user_info.json` file in the example folder. Edit it an put it in the module-specific folder inside the data folder that the module creates upon startup, together with the mod_gen_info.json file. This file configures the module information.

**PS:** no problem in using comments in the JSON files. They're all filtered.

## About
### - License
This project is licensed under Apache 2.0 License - http://www.apache.org/licenses/LICENSE-2.0.
