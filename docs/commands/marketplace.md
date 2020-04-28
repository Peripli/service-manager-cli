# marketplace

## Overview

`smctl marketplace`

Lists all service offerings with their service plans that are available in Service Manager.

## Usage

`smctl marketplace [flags]`

## Aliases

marketplace, m

## Parameters

|Optional|Global Flag|
|--------|-----------|
| -h, --help  Help for marketplace command.| No |
| -s, --service Detailed information about the plans of a specific service offering| No|
| -o, --output Output format of the command. Possible opitons: json, yaml, text.| No|
| --config Set the path for the smctl config.json file (default is $HOME/.sm/config.json).|Yes|
| -v, --verbose Use verbose mode.|Yes|

## Example

```
▶ smctl marketplace
One service offering.
Name              Plans            Description                                                                                       Broker ID
----------------  ---------------  ------------------------------------------------------------------------------------------------  ------------------------------------
overview-service  simple, complex  Provides an overview of any service instances and bindings that have been created by a platform.  46343c8e-957f-4fde-8176-ca3510d489e0
```

```
▶ smctl marketplace -s overview-service
2 service plans for this service offering.
Plan     Description               ID
-------  ------------------------  ------------------------------------
simple   A very simple plan.       25304783-2fc9-4f50-8dcb-0cbfe017ad15
complex  A more complicated plan.  52207e8e-1456-4f2e-b3df-3b97fe8d3d6f
```
