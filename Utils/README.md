# V.I.S.O.R. - Utils
Utilities and development platform package of V.I.S.O.R.

This repository is a submodule on the [V.I.S.O.R. - Server Version Assistant](https://github.com/Edw590/VISOR---Server-Version-Assistant) project (the main project).

## What it does
These utilities are used on the entire V.I.S.O.R. Server Version project. All functions here are supposed to (and do) work on all main operating systems (Unix-like and Windows) - they are universal. This can also be seen asthe development platform, since all modules use functions from here and there are UtilsModules, for example the mandatory ModStartup() function.

These utilities can't also be *completely* used as a standalone package since some of the functions here use V.I.S.O.R. modules - various don't though. The Email utilities do use the Email Sender module, for example. I wanted them to be completely independent from the modules, but at least the way I tried, it created cyclic imports and Go doesn't allow that. So all utilities and exported functions are here and the modules have only private elements and use these exported ones only.

## About
### - License
This project is licensed under Apache 2.0 License - http://www.apache.org/licenses/LICENSE-2.0.
