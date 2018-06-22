# Service Manager CLI

[![Build Status](https://travis-ci.org/Peripli/service-manager-cli.svg?branch=master)](https://travis-ci.org/Peripli/service-manager-cli)
[![Coverage Status](https://coveralls.io/repos/github/Peripli/service-manager-cli/badge.svg)](https://coveralls.io/github/Peripli/service-manager-cli)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/Peripli/service-manager-cli/blob/master/LICENSE)


***Service Manage CLI*** is the official command line client for [Service Manager][1]. 

## Getting started

To use the Service Manager CLI you need to download and install it first.

#### Download CLI
`` go get github.com/Peripli/service-manager-cli``

#### Install CLI

``go install github.com/Peripli/service-manager-cli/smctl``

#### Use CLI

You're done! Now you can use the **smctl** command along with some other subcommand (*register-broker*, *list-platforms*, etc...) to interact with a Service Manager instance.

##### Example:

```sh
# We need to connect and authenticate with a running Service Manager instance before doing anythign else  
smctl login -a http://service-manager-url.com -u admin -p admin

# List all brokers
smctl list-brokers
ID                                    Name  URL                             Description                                      Created               Updated               
------------------------------------  ----  ------------------------------  -----------------------------------------------  --------------------  --------------------

  
# Registering a broker
smctl register-broker some-service-broker https://demobroker.domain.com/ "Service broker providing some valuable services" -b admin:admin
ID                                    Name  URL                             Description                                      Created               Updated               
------------------------------------  ----  ------------------------------  -----------------------------------------------  --------------------  --------------------  
a52be735-30e5-4849-af23-83d65d592464  abc   https://demobroker.domain.com/  Service broker providing some valuable services  2018-06-22T13:04:19Z  2018-06-22T13:04:19Z


# Registering another broker
smctl register-broker def https://demobroker.domain.com/ "Another broker" -b admin:admin
ID                                    Name                 URL                             Description                                      Created               Updated               
------------------------------------  -------------------  ------------------------------  -----------------------------------------------  --------------------  -------------------- 
a52be735-30e5-4849-af23-83d65d592464  abc                  https://demobroker.domain.com/   Service broker providing some valuable services  2018-06-22T13:04:19Z  2018-06-22T13:04:19Z  
b419b538-b938-4293-86e0-7c92b0200d8e  def                  https://demobroker.domain.com/   Another broker                                   2018-06-22T13:05:41Z  2018-06-22T13:05:41Z 

```



For a list of all available commands run: ``smctl help``


[1]: https://github.com/Peripli/service-manager