# list-bindings

## Overview

`smctl list-bindings`

Lists all service bindings created in Service Manager.

## Usage

`smctl list-bindings [flags]`

## Aliases

list-bindings, lsb

## Parameters

|Optional|Global Flag|
|--------|-----------|
| -h, --help  Help for list-bindings command.| No |
| -o, --output Output format of the command. Possible opitons: json, yaml, text.| No|
| --config Set the path for the smctl config.json file (default is $HOME/.sm/config.json).|Yes|
| -v, --verbose Use verbose mode.|Yes|

## Example

```
▶ smctl list-bindings
One service binding.
ID                                    Name            Instance Name    Credentials                           Created                      Updated                      Ready  Labels
------------------------------------  --------------  ---------------  ------------------------------------  ---------------------------  ---------------------------  -----  ----------------
5937785d-6740-4f56-bdd9-8d24544bddac  sample-binding  sample-instance  {"password":"pass","username":"usr"}  2020-04-09T10:57:50.452161Z  2020-04-09T10:57:51.505822Z  true   tenant=tenant-id
```