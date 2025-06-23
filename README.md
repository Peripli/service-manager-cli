

<h1 align="center">⚠️ DEPRECATION NOTICE ⚠️</h1>

<p align="center">
  <strong>This project is no longer actively maintained and will be archived on <span style="color:red">30/09/2025</span>.</strong><br>
</p>


# Service Manager CLI

[![Build Status](https://github.com/Peripli/service-manager-cli/workflows/Go/badge.svg)](https://github.com/Peripli/service-manager-cli/actions)
[![Coverage Status](https://coveralls.io/repos/github/Peripli/service-manager-cli/badge.svg)](https://coveralls.io/github/Peripli/service-manager-cli)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/Peripli/service-manager-cli/blob/master/LICENSE)


***Service Manage CLI*** is the official command line client for [Service Manager][1]. 

## Getting started

To use the Service Manager CLI you need to download and install it first:

### Approach 1: Manual installation

### Enable Go Modules
`` export GO111MODULE=on ``

#### Download CLI
`` go get github.com/Peripli/service-manager-cli``

#### Install CLI

``go install github.com/Peripli/service-manager-cli``

#### Rename the CLI binary

``mv $GOPATH/bin/service-manager-cli $GOPATH/bin/smctl``

#### Use CLI

You're done! Now you can use the **smctl** command along with some other subcommand (*register-broker*, *list-platforms*, etc...) to interact with a Service Manager instance.

### Approach 2: Get the latest Service Manager CLI release
You can get started with the CLI by simply downloading the latest release from [HERE][2].

## Example usage of CLI:

```sh
# We need to connect and authenticate with a running Service Manager instance before doing anythign else  
smctl login -a http://service-manager-url.com -u {user} -p {pass}

# List all brokers
smctl list-brokers
ID                                    Name  URL                             Description                                      Created               Updated               
------------------------------------  ----  ------------------------------  -----------------------------------------------  --------------------  --------------------

  
# Registering a broker
smctl register-broker sample-broker-1 https://demobroker.domain.com/ "Service broker providing some valuable services" -b {user}:{pass}
ID                                    Name             URL                             Description                                      Created               Updated               
------------------------------------  ---------------  ------------------------------  -----------------------------------------------  --------------------  --------------------  
a52be735-30e5-4849-af23-83d65d592464  sample-broker-1  https://demobroker.domain.com/  Service broker providing some valuable services  2018-06-22T13:04:19Z  2018-06-22T13:04:19Z


# Registering another broker
smctl register-broker sample-broker-2 https://demobroker.domain.com/ "Another broker providing valuable services" -b {user}:{pass}
ID                                    Name             URL                             Description                                      Created               Updated               
------------------------------------  ---------------  ------------------------------  -----------------------------------------------  --------------------  -------------------- 
a52be735-30e5-4849-af23-83d65d592464  sample-broker-1  https://demobroker.domain.com/   Service broker providing some valuable services  2018-06-22T13:04:19Z  2018-06-22T13:04:19Z  
b419b538-b938-4293-86e0-7c92b0200d8e  sample-broker-2  https://demobroker.domain.com/   Another broker providing valuable services       2018-06-22T13:05:41Z  2018-06-22T13:05:41Z 

```



For a list of all available commands run: ``smctl help``

## Documentation
Documentation of the Service Manager CLI and all of it's commands can be found [HERE][3].


[1]: https://github.com/Peripli/service-manager
[2]: https://github.com/Peripli/service-manager-cli/releases
[3]: docs/README.md
