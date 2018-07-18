# smctl info

## Overview
`smctl info` displays information about the Service Manager instance the *smctl* is connected to.

## Usage
```bash
smctl info [flags]
```

## Aliases
```bash
info, i 
```

## Flags
<details>
  <summary>help</summary>
  <p>
    <code>--help</code> (alias: <code>-h</code>)
  </p>
  <p>
    Help for <i>info</i> command. 
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
> smctl info

Service Manager URL: https://service-manager-url.com/
Logged user: testUser 
```