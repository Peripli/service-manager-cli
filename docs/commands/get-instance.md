# get-instance

## Overview

`smctl get-instance`

Get detailed information about the service instance with provided name.

## Usage

`smctl get-instance [name] [flags]`

## Aliases

get-instance, gi

## Parameters

|Optional|Global Flag|
|--------|-----------|
| -h, --help  Help for get-instance command.| No |
| -o, --output Output format of the command. Possible opitons: json, yaml, text.| No|
| --config Set the path for the smctl config.json file (default is $HOME/.sm/config.json).|Yes|
| -v, --verbose Use verbose mode.|Yes|

## Example

```
▶ smctl get-instance sample-instance
One service instance.
| ID               | 0c170e73-28bd-47ea-b3f4-f1ad1dbf3e0a  |
| Name             | sample-instance                       |
| Service Plan ID  | 25304783-2fc9-4f50-8dcb-0cbfe017ad15  |
| Platform ID      | service-manager                       |
| Created          | 2020-04-09T10:42:12.175051Z           |
| Updated          | 2020-04-09T10:42:13.22521Z            |
| Ready            | true                                  |
| Usable           | true                                  |
| Labels           | tenant=tenant-id                      |
| Last Op          | create succeeded                      |
```