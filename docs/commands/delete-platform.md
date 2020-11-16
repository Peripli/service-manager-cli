# smctl delete-platform

## Overview
`smctl delete-platform` deletes a set of platforms registered in the Service Manager instance.

## Usage
```bash
smctl delete-platform [name] <name2 <name3> ... <nameN>> [flags]
```

## Aliases
```bash
delete-platform, dp 
```

## Flags
<details>
  <summary>help</summary>
  <p>
    <code>--help</code> (alias: <code>-h</code>)
  </p>
  <p>
    Help for <i>delete-platform</i> command. 
  </p>
</details>
<details>
  <summary>cascade-delete</summary>
  <p>
    <code>--cascade</code> 
  </p>
  <p>
    Cascade delete for <i>delete-platform</i> command. 
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
> smctl delete-platform sample-platform

Platform with name: sample-platform successfully deleted
```
```bash
> smctl delete-platform sample-platform --cascade

Cascade delete successfully scheduled for platform: sample-platform. To see status of the operation use:
smctl status /v1/platforms/baea022b-64c0-43d4-a9b0-e1ae64af51cd/operations/f8ca64af-e889-4a45-ad41-f1baa2e427c2
```