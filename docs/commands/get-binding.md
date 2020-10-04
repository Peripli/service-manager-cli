# get-binding

## Overview

`smctl get-binding`

Get detailed information about the service binding with provided name.

## Usage

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
| --binding-params  Show service binding configuration parameters.| No |

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


```
▶ smctl get-binding sample-binding --binding-params
Showing parameters for service binding id:  0c170e73-28bd-47ea-b3f4-f1ad1dbf3e0a
The parameters are:
{
   "param1":"value1",
   "param2":"value2"
}
```