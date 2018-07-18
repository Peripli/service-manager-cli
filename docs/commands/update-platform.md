# smctl update-platform

## Overview
`smctl update-platform` updates a platform with the provided name in the Service Manager instance.

## Usage
```bash
smctl update-platform [name] <json_platform> [flags]
```

## Aliases
```bash
update-platform, up 
```

## Flags
<details>
  <summary>help</summary>
  <p>
    <code>--help</code> (alias: <code>-h</code>)
  </p>
  <p>
    Help for <i>update-platform</i> command. 
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
> smctl update-platform sample-platform '{"description": "Sample platform instance"}' 

ID                                    Name             Type    Description               Created               Updated               
------------------------------------  ---------------  ------  ------------------------  --------------------  --------------------  
6352fca0-c252-43ab-9cb3-d23613749b59  sample-platform  sample  Sample platform instance  2018-07-18T07:06:40Z  2018-07-18T07:09:48Z  
```
