# smctl login

## Overview
`smctl login` authenticates user against a Service Manager instance.

## Usage
```bash
smctl login [flags]
```

## Aliases
```bash
login, l
```

## Flags
<details>
  <summary>help</summary>
  <p>
    <code>--help</code> (alias: <code>-h</code>)
  </p>
  <p>
    Help for <i>login</i> command. 
  </p>
</details>
<details>
  <summary>password</summary>
  <p>
    <code>--passowrd</code> (alias: <code>-p</code>)
  </p>
  <p>
    User password.
  </p>
</details>
<details>
  <summary>skip ssl validation</summary>
  <p>
    <code>--skip-ssl-validation</code>
  </p>
  <p>
    Skip verification of the OAuth endpoint <b>Not recommended!</b>
  </p>
</details>
<details>
  <summary>url</summary>
  <p>
    <code>--url</code> (alias: <code>-a</code>)
  </p>
  <p>
    Base URL of the Service Manager.
  </p>
</details>
<details>
  <summary>user</summary>
  <p>
    <code>--user</code> (alias: <code>-u</code>)
  </p>
  <p>
    User ID.
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