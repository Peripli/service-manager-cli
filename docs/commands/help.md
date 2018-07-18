# smctl help

## Overview
`smctl help` displays the whole list of available commands. When a specific command is provided, help information for the specified command is displayed.

## Usage
```bash
smctl help [command] [flags]
```

## Aliases
```bash
help, -h 
```

## Flags
None.

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
> smctl help

smctl controls a Service Manager instance.

Usage:
  smctl [command]

Available Commands:
  delete-broker     Deletes brokers
  delete-platform   Deletes platforms
  help              Help about any command
  info              Prints information for logged user
  list-brokers      List brokers
  list-platforms    List platforms
  login             Logs user in
  register-broker   Registers a broker
  register-platform Registers a platform
  update-broker     Updates broker
  update-platform   Updates platform
  version           Prints smctl version

Flags:
      --config string   config file (default is $HOME/.sm/config.json)
  -h, --help            help for smctl
  -v, --verbose         verbose

Use "smctl [command] --help" for more information about a command.
```