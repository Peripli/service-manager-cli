# Service Manager CLI Documentation

The Service Manager CLI is the official tool to communicate with a Service Manager instance directly.

> **NOTE**: Throughout this document we will refer to the Service Manager CLI as *SM CLI* for simplicity.

## Getting Started
In order to start using the SM CLI you need to download and install it. You can get the lastest SM CLI release from [HERE][1].

## Commands
The SM CLI provides commands for creating, listing, updating and deleting service brokers and platforms in a Service Manager instance. Here's a full list of the available commands:

#### Login
* [login][2]

#### Brokers
* [register-broker][3]
* [update-broker][4]
* [list-brokers][5]
* [delete-broker][6]

#### Platforms
* [register-platform][7]
* [update-platform][8]
* [list-platforms][9]
* [delete-platform][10]

#### Marketplace
* [list-offerings][11]
* [list-plans][12]
* [marketplace][13]

#### Instances
* [provision][14]
* [get-instance][15]
* [list-instances][16]
* [deprovision][17]

#### Bindings
* [bind][18]
* [get-binding][19]
* [list-bindings][20]
* [unbind][21]

#### Status
* [status][22]

#### Misc
* [info][23]
* [version][24]
* [help][25]

[1]: https://github.com/Peripli/service-manager-cli/releases

[2]: commands/login.md

[3]: commands/register-broker.md
[4]: commands/update-broker.md
[5]: commands/list-brokers.md
[6]: commands/delete-broker.md

[7]: commands/register-platform.md
[8]: commands/update-platform.md
[9]: commands/list-platforms.md
[10]: commands/delete-platform.md

[11]: commands/list-offerings.md
[12]: commands/list-plans.md
[13]: commands/marketplace.md

[14]: commands/provision.md
[15]: commands/get-instance.md
[16]: commands/list-instances.md
[17]: commands/deprovision.md

[18]: commands/bind.md
[19]: commands/get-binding.md
[20]: commands/list-bindings.md
[21]: commands/unbind.md

[22]: commands/status.md

[23]: commands/info.md
[24]: commands/version.md
[25]: commands/help.md