# status

## Overview

`smctl status`

Get asynchronous operation's status

## Usage

`smctl status [operation URL path] [flags]`

## Parameters

|Optional|Global Flag|
|--------|-----------|
| -h, --help  Help for status command.| No |
| -o, --output Output format of the command. Possible opitons: json, yaml, text.| No|
| --config Set the path for the smctl config.json file (default is $HOME/.sm/config.json).|Yes|
| -v, --verbose Use verbose mode.|Yes|

## Example

```
▶ smctl status /v1/service_bindings/5937785d-6740-4f56-bdd9-8d24544bddac/operations/6066bd46-79d4-4f8e-be50-9ad2e5ca035a

| ID     | 6066bd46-79d4-4f8e-be50-9ad2e5ca035a  |
| Type   | delete                                |
| State  | succeeded                             |
```