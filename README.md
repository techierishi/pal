# pal : Simple cli app which makes your cli interaction easier.
 
| <img src="doc/logo.png" width="200">  |   A CLI app with features like snippet manager, shell history manager, clipboard manager, alias manager, config sync etc.   |
|---|---|

---

<p align="center">
<img src="doc/pal.gif" >
</p>


# Usage

Run `pal init` right after installation. This step ensures `ctr+r` will open interactive shell history search. It also add some alias shortcuts for commond pal commands

```
pal -  Simple cli app which makes your cli interaction easier. Its based on another awesome project (github.com/knqyf263/pet)

Usage:
  pal [command]

Available Commands:
  alias       Simple alias manager.
  backup      Backup all configs
  clip        Simple clipboard manager.
  completion  Generate the autocompletion script for the specified shell
  configure   Edit config file
  cred        Simple credential manager.
  gen         Generate aliases, key mappings
  help        Help about any command
  hist        Simple shell history manager.
  init        Initial pal
  restore     Backup all configs
  snip        Simple command-line snippet manager.
  svc         Pal backgroud service.
  sync        Sync configs
  version     Print the version number

Flags:
      --config string   config file (default is C:\Users\techi\AppData\Roaming\pal )
      --debug           debug mode
  -h, --help            help for pal

Use "pal [command] --help" for more information about a command.
```

## Snippet usage

```
snip - Simple command-line snippet manager.

Usage:
  pal snip [command]

Available Commands:
  copy        Copy the selected commands
  edit        Edit snippet file
  exec        Run the selected commands
  list        Show all snippets
  new         Create a new snippet
  search      Search snippets

Flags:
  -h, --help   help for snip

Global Flags:
      --config string   config file (default is C:\Users\techi\AppData\Roaming\pal )
      --debug           debug mode

Use "pal snip [command] --help" for more information about a command.

```
## Snippet demo

https://github.com/techierishi/pal/assets/7880021/ba677979-f399-474c-a1c2-7e19bc73eaf1



## Shell history usage

```
hist - Simple shell history manager.

Usage:
  pal hist [flags]

Flags:
  -c, --command   Show the command with the plain text before executing
  -p, --copy      Just copy command with the plain text.
  -h, --help      help for hist

Global Flags:
      --config string   config file (default is C:\Users\techi\AppData\Roaming\pal )
      --debug           debug mode
```

## Shell history demo

https://github.com/techierishi/pal/assets/7880021/c79aee98-4160-4105-a0fc-5463d5ee0b70



## Clipboard usage

```
clip - Simple clipboard manager.

Usage:
  pal clip [command]

Available Commands:
  list        Clipboard history

Flags:
  -h, --help   help for clip

Global Flags:
      --config string   config file (default is C:\Users\techi\AppData\Roaming\pal )
      --debug           debug mode

Use "pal clip [command] --help" for more information about a command.
```

## Clipboard demo

https://github.com/techierishi/pal/assets/7880021/8878ef0c-daa4-47d9-aae0-c9436cac9d61

## Credential usage

> This feature is only to support snippet manager if and when a credential is needed for some snippet
> Please note that the password is only saved os keychain not in plain yaml file
```
cred - Simple credential manager.

Usage:
  pal cred [command]

Available Commands:
  edit        Edit credential file
  list        Password list to console
  new         Create new credential
  search      Password search

Flags:
  -h, --help   help for cred

Global Flags:
      --config string   config file (default is C:\Users\techi\AppData\Roaming\pal )
      --debug           debug mode

Use "pal cred [command] --help" for more information about a command.
```
## Credential manager demo

https://github.com/techierishi/pal/assets/7880021/8ff7d8aa-bf41-4fb9-a9a7-44d241913721


## Installation

Currently you can clone this repo and compile using following steps
```
git clone https://github.com/techierishi/pal.git
cd pal
go mod tidy
go build
```

We will add binaries for multiple platforms in future

