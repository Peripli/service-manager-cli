# unbind

## Overview

`smctl unbind`

Deletes a service binding.

## Usage

`smctl unbind [instance-name] [binding-name] [flags]`

## Parameters

|Optional|Global Flag|
|--------|-----------|
| -h, --help  Help for unbind command. | No |
| -f, --force Use this parameter to delete a resource without raising a confirmation message. | No |
| --force-delete Delete the service binding and all of its associated resources from the database. Use this parameter if the service binding cannot be properly deleted. This parameter can only be used by operators with technical access. | No |
| --id ID of the service binding. Required when name is ambiguous. | No |
| --mode  Whether to use synchronous or asynchronous calls to Service Management. The default value is 'async'. | No |
| -o, --output The output format of the command. Options: json, yaml, text. | No |
| --config Set the path for the smctl config.json file (default is $HOME/.sm/config.json). | Yes |
| -v, --verbose Use the Verbose mode. | Yes |

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
