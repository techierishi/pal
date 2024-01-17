# pal : Simple cli app which makes your cli interaction easier.

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
[![Snippet Demo](https://img.youtube.com/vi/_1TxNmnTArY/maxresdefault.jpg)](https://youtu.be/_1TxNmnTArY)

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
[![Shell history demo](https://img.youtube.com/vi/wlgeyqcLLdc/maxresdefault.jpg)](https://youtu.be/wlgeyqcLLdc)

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
[![Clipboard demo](https://img.youtube.com/vi/oLKWGm4od7c/maxresdefault.jpg)](https://youtu.be/oLKWGm4od7c)

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
[![Credential manager demo](https://img.youtube.com/vi/Aq59JX6Er1E/maxresdefault.jpg)](https://youtu.be/Aq59JX6Er1E)

## Installation

Currently you can clone this repo and compile using following steps
```
git clone https://github.com/techierishi/pal.git
cd pal
go mod tidy
go build
```

We will add binaries for multiple platforms in future

