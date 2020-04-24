# list-plans

## Overview

`smctl list-plans`

Lists all service plans that are available in Service Manager.

## Usage

`smctl list-plans [flags]`

## Parameters

|Optional|Global Flag|
|--------|-----------|
| -h, --help  Help for list-plans command.| No |
| -o, --output Output format of the command. Possible opitons: json, yaml, text.| No|
| --config Set the path for the smctl config.json file (default is $HOME/.sm/config.json).|Yes|
| -v, --verbose Use verbose mode.|Yes|

## Example

```
▶ smctl list-plans
2 service plans.
ID                                    Name     Description               Offering ID                           Ready  Labels
------------------------------------  -------  ------------------------  ------------------------------------  -----  ------
aec1cdac-9faa-4aa4-aeb7-6dbcc275208d  simple   A very simple plan.       a56dc9b4-70f9-45e3-a8a1-3b3a06289aa5  true
eb2ce3e0-d64c-4977-b51e-2f682ec2d835  complex  A more complicated plan.  a56dc9b4-70f9-45e3-a8a1-3b3a06289aa5  true
```