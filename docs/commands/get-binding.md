# get-binding

## Overview

`smctl get-binding`

Get detailed information about the service binding with provided name.

## Usage
// TODO update
`smctl get-binding [name] [flags]`

## Aliases

get-binding, gsb

## Parameters

|Optional|Global Flag|
|--------|-----------|
| -h, --help  Help for get-binding command.| No |
| -o, --output Output format of the command. Possible opitons: json, yaml, text.| No|
| --config Set the path for the smctl config.json file (default is $HOME/.sm/config.json).|Yes|
| -v, --verbose Use verbose mode.|Yes|

## Example

```
▶ smctl get-binding sample-binding
One service binding.
| ID             | 5937785d-6740-4f56-bdd9-8d24544bddac                |
| Name           | sample-binding                                      |
| Instance Name  | sample-instance                                     |
| Credentials    | {"password":"pass","username":"usr"}                |
| Created        | 2020-04-09T10:57:50.452161Z                         |
| Updated        | 2020-04-09T10:57:51.505822Z                         |
| Ready          | true                                                |
| Labels         | tenant=tenant-id                                    |
| Last Op        | create succeeded                                    |
```