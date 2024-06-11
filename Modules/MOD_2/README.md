# V.I.S.O.R. - S.M.A.R.T. Checker
S.M.A.R.T. Checker module of V.I.S.O.R.

This repository is a submodule on the [V.I.S.O.R. - Server Version Assistant](https://github.com/Edw590/VISOR---Server-Version-Assistant) project (the main project).

## What it does
This module runs S.M.A.R.T. tests on the given disks and checks the S.M.A.R.T. information after the tests are done.

Check the `mod_user_info.json` file in the example folder. Edit it an put it in the module-specific folder inside the data folder that the module creates upon startup, together with the mod_gen_info.json file. This file configures the module information.

**PS:** no problem in using comments in the JSON files. They're all filtered.

## Command line arguments
If --notest is passed, the module will not run the tests, only send the S.M.A.R.T. information report.

## About
### - License
This project is licensed under Apache 2.0 License - http://www.apache.org/licenses/LICENSE-2.0.
