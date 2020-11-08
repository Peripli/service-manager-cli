# smctl update-broker

## Overview
`smctl update-broker` updates a service broker with the provided name in the Service Manager instance.

## Usage
```bash
smctl update-broker [name] <json_broker> [flags]
```

## Example
```bash
smctl update-broker broker '{"name": "new-name", "description": "new-description", "broker_url": "http://broker.com", "credentials": { "basic": { "username": "admin", "password": "admin" } }}'
```

## Aliases
```bash
update-broker, ub 
```

## Flags
<details>
  <summary>help</summary>
  <p>
    <code>--help</code> (alias: <code>-h</code>)
  </p>
  <p>
    Help for <i>update-broker</i> command. 
  </p>
</details>
<details>
  <summary>output format</summary>
  <p>
    <code>--output</code> (alias: <code>-o</code>)
  </p>
  <p>
    Output format of the command. Possible opitons: <i>json, yaml, text</i>
  </p>
</details>

## Global Flags
<details>
  <summary>config</summary>
  <p>
    <code>--config</code> 
  </p>
  <p>
    Set the path for the <b>smctl</b> <i>config.json</i> file (default is <i>$HOME/.sm/config.json</i>)
  </p>
</details>
<details>
  <summary>verbose</summary>
  <p>
    <code>--verbose</code> (alias: <code>-v</code>)
  </p>
  <p>
    Use verbose mode.
  </p>
</details>

## Example
```bash
> smctl update-broker sample-broker-1 '{"description": "Updated sample-broker description"}' 

ID                                    Name             URL                             Description                        Created               Updated               
------------------------------------  ---------------  ------------------------------  ---------------------------------  --------------------  --------------------  
a52be735-30e5-4849-af23-83d65d592464  sample-broker-1  https://demobroker.domain.com/  Updated sample-broker description  2018-06-22T13:04:19Z  2018-06-22T13:04:19Z
```
