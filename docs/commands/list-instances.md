# list-instances

## Overview

`smctl list-instances`

Lists all service instances.

## Usage

`smctl list-instances [flags]`

## Aliases

list-instances, li

## Parameters

|Optional|Global Flag|
|--------|-----------|
| -h, --help  Help for list-plans command.| No |
| -o, --output Output format of the command. Possible opitons: json, yaml, text.| No|
| --config Set the path for the smctl config.json file (default is $HOME/.sm/config.json).|Yes|
| -v, --verbose Use verbose mode.|Yes|

## Example

```
▶ smctl list-instances
One service instance.
ID                                    Name             Service Plan ID                       Platform ID      Created                      Updated                     Ready  Usable  Labels
------------------------------------  ---------------  ------------------------------------  ---------------  ---------------------------  --------------------------  -----  ------  ----------------
0c170e73-28bd-47ea-b3f4-f1ad1dbf3e0a  sample-instance  25304783-2fc9-4f50-8dcb-0cbfe017ad15  service-manager  2020-04-09T10:42:12.175051Z  2020-04-09T10:42:13.22521Z  true   true    tenant=tenant-id
```