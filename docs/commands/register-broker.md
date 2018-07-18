# smctl register-broker

## Overview
`smctl register-broker` registers a broker in the Service Manager instance.

## Usage
```bash
smctl register-broker [name] [url] <description> [flags]
```

## Aliases
```bash
register-broker, rb 
```

## Flags
<details>
  <summary>help</summary>
  <p>
    <code>--help</code> (alias: <code>-h</code>)
  </p>
  <p>
    Help for <i>register-broker</i> command. 
  </p>
</details>
<details>
  <summary>basic credentials</summary>
  <p>
    <code>--basic</code> (alias: <code>-b</code>)
  </p>
  <p>
    Sets the username and password for basic authentication. Format is <i>&lt;username:passowrd&gt;</i>
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
> smctl register-broker sample-broker-1 https://demobroker.domain.com/ "Service broker providing some valuable services" -b user:pass

ID                                    Name             URL                             Description                                      Created               Updated               
------------------------------------  ---------------  ------------------------------  -----------------------------------------------  --------------------  --------------------  
a52be735-30e5-4849-af23-83d65d592464  sample-broker-1  https://demobroker.domain.com/  Service broker providing some valuable services  2018-06-22T13:04:19Z  2018-06-22T13:04:19Z
```