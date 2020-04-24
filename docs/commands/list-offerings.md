# list-offerings

## Overview

`smctl list-offerings`

Lists all service offerings that are available in Service Manager.

## Usage

`smctl list-offerings [flags]`

## Aliases

list-offerings, lo

## Parameters

|Optional|Global Flag|
|--------|-----------|
| -h, --help  Help for list-offerings command.| No |
| -o, --output Output format of the command. Possible opitons: json, yaml, text.| No|
| --config Set the path for the smctl config.json file (default is $HOME/.sm/config.json).|Yes|
| -v, --verbose Use verbose mode.|Yes|

## Example

```
▶ smctl list-offerings
One service offering.
ID                                    Name              Description                                                                                       Broker ID                             Ready  Labels
------------------------------------  ----------------  ------------------------------------------------------------------------------------------------  ------------------------------------  -----  ------
54944d91-75b9-442c-aecd-f98821490740  overview-service  Provides an overview of any service instances and bindings that have been created by a platform.  46343c8e-957f-4fde-8176-ca3510d489e0  true
```