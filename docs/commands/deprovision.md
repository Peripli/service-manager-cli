# deprovision

## Overview

`smctl deprovision`

Deletes a service instance.

## Usage

`smctl deprovision [name] [flags]`

## Parameters

|Optional|Global Flag|
|--------|-----------|
| -h, --help Help for deprovision command.| No |
| -f, --force Use this parameter to delete a resource without raising a confirmation message. | No |
| --force-delete Delete the service instance and all of its associated resources from the database, including all service bindings. Use this parameter if the service instance cannot be properly deleted. This parameter can only be used by operators with technical access. | No |
| --id ID of the service instance. Required when name is ambiguous. | No |
| --mode How calls to Service Manager are performed sync or async (default "async"). | No |
| --config Set the path for the smctl config.json file (default is $HOME/.sm/config.json). | Yes |
| -v, --verbose Use the Verbose mode. | Yes |

## Example

> Hint: See smctl status below

Async execution:
```
▶ smctl deprovision sample-instance
Do you really want to delete instance with name [sample-instance] (Y/n): yes
Service Instance sample-instance successfully scheduled for deletion. To see status of the operation use:
smctl status /v1/service_instances/0c170e73-28bd-47ea-b3f4-f1ad1dbf3e0a/operations/40a748c1-c0f8-4acf-84a0-64e20914531d
```
```
▶ smctl status /v1/service_instances/0c170e73-28bd-47ea-b3f4-f1ad1dbf3e0a/operations/40a748c1-c0f8-4acf-84a0-64e20914531d

| ID     | 40a748c1-c0f8-4acf-84a0-64e20914531d  |
| Type   | delete                                |
| State  | succeeded                             |
```

Sync execution:
```
▶ smctl deprovision sample-instance --mode sync
Do you really want to delete instance with name [sample-instance] (Y/n): yes
Service Instance successfully deleted.
```
