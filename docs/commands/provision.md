# provision

## Overview

`smctl provision`

Create a service instance in Service Manager.

## Usage

`smctl provision [name] [offering] [plan] [flags]`

## Parameters

|Optional|Global Flag|
|--------|-----------|
| -h, --help  Help for provision command.| No |
| -b, --broker-name Name of the broker which provides the service offering. Required when offering name is ambiguous| No|
| --mode How calls to Service Manager are performed sync or async (default "async") | No |
| -c, --parameters Valid JSON object containing instance parameters | No |
| -o, --output Output format of the command. Possible opitons: json, yaml, text.| No|
| --config Set the path for the smctl config.json file (default is $HOME/.sm/config.json).|Yes|
| -v, --verbose Use verbose mode.|Yes|

## Example

> Hint: See smctl status below

Async execution:
```
▶ smctl provision sample-instance overview-service simple
Service Instance sample-instance successfully scheduled for provisioning. To see status of the operation use:
smctl status /v1/service_instances/a6b0dfe6-1bd1-453f-a646-babd425b6b05/operations/32bbbee7-a9d0-48e4-a434-bf47bc471a48
```

```
▶ smctl status /v1/service_instances/a6b0dfe6-1bd1-453f-a646-babd425b6b05/operations/32bbbee7-a9d0-48e4-a434-bf47bc471a48

| ID     | 32bbbee7-a9d0-48e4-a434-bf47bc471a48  |
| Type   | create                                |
| State  | succeeded                             |
```

Sync execution:
```
▶ smctl provision sample-instance overview-service simple --mode sync

| ID               | 0c170e73-28bd-47ea-b3f4-f1ad1dbf3e0a  |
| Name             | sample-instance                       |
| Service Plan ID  | 25304783-2fc9-4f50-8dcb-0cbfe017ad15  |
| Platform ID      | service-manager                       |
| Created          | 2020-04-09T10:42:12.175051Z           |
| Updated          | 2020-04-09T10:42:13.2252101Z          |
| Ready            | true                                  |
| Usable           | true                                  |
| Labels           | tenant=tenant-id                      |
| Last Op          | create succeeded                      |
```