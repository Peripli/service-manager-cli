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
    <code>--password</code> (alias: <code>-p</code>)
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
<details>
  <summary>auth flow</summary>
  <p>
    <code>--auth-flow</code>
  </p>
  <p>
    Options: <code>password</code> / <code>client-credentials</code> (default is <code>password</code> flow)
  </p>
</details>
<details>
  <summary>client id</summary>
  <p>
    <code>--client-id</code>
  </p>
  <p>
    The technical client ID that was generated upon the creation of binding and that is used for the <code>client-credentials</code> authorization flow.
  </p>
</details>
<details>
  <summary>client secret</summary>
  <p>
    <code>--client-secret</code>
  </p>
  <p>
    The technical client secret that was generated upon the creation of binding and that is used for the <code>client-credentials</code> authorization flow.
  </p>
</details>
<details>
  <summary>certificate</summary>
  <p>
    <code>--cert</code>
  </p>
  <p>
    A path to the file that contains the public key <code>certificate</code> that was generated upon the creation of binding and that is used for the <code>client-credentials</code> authorization flow.
  </p>
</details>
<details>
  <summary>private key</summary>
  <p>
    <code>--key</code>
  </p>
  <p>
    A path to the file that contains the private <code>key</code> that was generated upon the creation of binding and that is used for the <code>client-credentials</code> authorization flow.
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

## Example 1 - password flow
```bash
> smctl login -a https://service-manager-url.com

User: user                # entering username
Password:                 # entering password (password visibility is disabled)
Logged in successfully.
```

## Example 2 - password flow
```bash
> smctl login -a https://service-manager-url.com -u user -p pass

Logged in successfully.
```


## Example 3 - client id & secret
Requires: client-id, client-secret
```bash
> smctl login -a https://service-manager-url.com --auth-flow=client-credentials --client-id=id --client-secret=secret

Logged in successfully.
```

## Example 3 - mTLS (with a certificate)
Requires: client-id, cert, key
```bash
> smctl login -a https://service-manager-url.com --auth-flow=client-credentials --client-id=id --cert=cert.pem --key=key.pem

Logged in successfully.
```

