# unbind

## Overview

`smctl unbind`

Deletes a service binding.

## Usage

`smctl unbind [instance-name] [binding-name] [flags]`

## Parameters

|Optional|Global Flag|
|--------|-----------|
| -h, --help  Help for unbind command.| No |
| -f, --force Force delete without confirmation | No |
| --id ID of the service binding. Required when name is ambiguous| No |
| --mode  How calls to Service Manager are performed sync or async (default "async")| No |
| -o, --output Output format of the command. Possible opitons: json, yaml, text.| No|
| --config Set the path for the smctl config.json file (default is $HOME/.sm/config.json).|Yes|
| -v, --verbose Use verbose mode.|Yes|

## Example

> Hint: See smctl status below

Async execution:

```
▶ smctl unbind sample-instance sample-binding
Do you really want to delete binding with name [sample-binding] for instance with name sample-instance (Y/n): yes
Service Binding sample-binding successfully scheduled for deletion. To see status of the operation use:
smctl status /v1/service_bindings/5937785d-6740-4f56-bdd9-8d24544bddac/operations/6066bd46-79d4-4f8e-be50-9ad2e5ca035a
```

```
▶ smctl status /v1/service_bindings/5937785d-6740-4f56-bdd9-8d24544bddac/operations/6066bd46-79d4-4f8e-be50-9ad2e5ca035a

| ID     | 6066bd46-79d4-4f8e-be50-9ad2e5ca035a  |
| Type   | delete                                |
| State  | succeeded                             |
```

Sync execution:

```
▶ smctl unbind sample-instance sample-binding --mode sync
Do you really want to delete binding with name [sample-binding] for instance with name sample-instance (Y/n): yes
Service Binding successfully deleted.
```