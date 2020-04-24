# bind

## Overview

`smctl bind`

Creates a binding in Service Manager to an instance with the provided name.

## Usage

`smctl bind [instance-name] [binding-name] [flags]`

## Parameters

|Optional|Global Flag|
|--------|-----------|
| -h, --help  Help for bind command.| No |
| -o, --output Output format of the command. Possible opitons: json, yaml, text.| No|
| --mode How calls to Service Manager are performed sync or async (default "async") | No |
| -c, --parameters Valid JSON object containing binding parameters | No |
| --id ID of the service instance. Required when name is ambiguous | No |
| --config Set the path for the smctl config.json file (default is $HOME/.sm/config.json).|Yes|
| -v, --verbose Use verbose mode.|Yes|

## Example

> Hint: See smctl status below

Async execution:

```
▶ smctl bind sample-instance sample-binding
Service Binding sample-binding successfully scheduled. To see status of the operation use:
smctl status /v1/service_bindings/6372815d-29b8-4561-9898-016d6671b34b/operations/ea67e94f-ad2f-4544-8b87-6924fe494327
```

```
▶ smctl status /v1/service_bindings/6372815d-29b8-4561-9898-016d6671b34b/operations/ea67e94f-ad2f-4544-8b87-6924fe494327

| ID     | ea67e94f-ad2f-4544-8b87-6924fe494327  |
| Type   | create                                |
| State  | succeeded                             |
```

Sync execution:

```
▶ smctl bind sample-instance sample-binding --mode sync

| ID                     | 5937785d-6740-4f56-bdd9-8d24544bddac                |
| Name                   | sample-binding                                      |
| Service Instance Name  | sample-instance                                     |
| Service Instance ID    | 742b0c67-37f6-4c63-83d9-e3c5d2cb69f0                |
| Credentials            | {"password":"pass","username":"usr"}                |
| Created                | 2020-04-09T10:57:50.452161Z                         |
| Updated                | 2020-04-09T10:57:51.5058215Z                        |
| Ready                  | true                                                |
| Labels                 | tenant=tenant-id                                    |
| Last Op                | create succeeded                                    |
```