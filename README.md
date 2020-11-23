# nssh

SSH client for [SORACOM Napter](https://developers.soracom.io/en/docs/napter/). You can easily open an SSH connection with your subscriber's name.

Napter is an on-demand networking service for devices using Soracom Air for Cellular SIM cards, which enables you to quickly and securely access your devices remotely. Napter allows you to perform remote maintenance, troubleshooting, or other typical remote access tasks, without setting up any relay servers or installing agent software on the device.

## Install

1. Download the archive for your platform and architecture from [Releases](https://github.com/0x6b/nssh/releases) section.
2. Unpack the archive.
3. Move the executable to one of your `PATH` directories.

Or you can build executable from the source:

```console
$ git clone https://github.com/0x6b/nssh
$ cd nssh
$ make # and you'll get `nssh` under the root directory
```

## Usage

### One-time Setup

1. Create a SAM user with following permission:
   ```json
   {
     "statements": [
       {
         "api": [
           "Auth:auth",
           "Subscriber:listSubscribers",
           "PortMapping:listPortMappingsForSubscriber",
           "PortMapping:createPortMapping"
         ],
         "effect": "allow"
       }
     ]
   }
   ```
2. Generate authentication key for the user.
3. Save the authentication information at `$HOME/.soracom/ssh.json`.
   ```json5
   {
     "coverageType": "jp", // or "g" for global
     "authKeyId": "keyId-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
     "authKey": "secret-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
   }
   ```
4. Name your desired SIM at SORACOM User Console.

### Connect

```console
$ nssh connect pi@your-sim-name
```

You can specify coverage type, profile name, port number, connection duration, or identity file for SSH public key authentication. See `nssh connect --help`.

- Override coverage type, `jp` or `global`:
  ```console
  $ nssh --coverage-type global connect pi@your-sim-name
  ```
- Use another profile under `$HOME/.soracom/` directory, without extension `.json`:
  ```console
  $ nssh --profile-name default connect pi@your-sim-name
  ```
- Use public key authentication:
  ```console
  $ nssh connect pi@your-sim-name -i ~/.ssh/id_rsa
  ```
- Specify another port number and connection duration:
  ```console
  $ nssh connect pi@your-sim-name --port 2222 --duration 120
  ```

### Details

Global help:

```console
$ nssh --help
nssh -- SSH client for SORACOM Napter

Usage:
  nssh [command]

Available Commands:
  connect     Connect to specified subscriber via SSH.
  help        Help about any command
  list        List port mappings for specified subscriber
  version     Show version

Flags:
      --coverage-type string   Specify coverage type, "g" for Global, "jp" for Japan
  -h, --help                   help for nssh
      --profile-name string    Specify SORACOM CLI profile name (default "ssh")
```

Help for `connect` sub-command:

```console
$ nssh connect --help
Create port mappings for specified subscriber and connect via SSH. If <user>@ is not specified, "pi" will be used as default. Quote with " if name contains spaces or special characters.

Usage:
  nssh connect [<user>@]<subscriber name> [flags]

Aliases:
  connect, c

Flags:
  -d, --duration int      Specify session duration in minutes (default 60)
  -h, --help              help for connect
  -i, --identity string   Specify a path to file from which the identity for public key authentication is read
  -p, --port int          Specify port number to connect (default 22)

Global Flags:
      --coverage-type string   Specify coverage type, "g" for Global, "jp" for Japan
      --profile-name string    Specify SORACOM CLI profile name (default "ssh")
```

Help for `list` sub-command:

```console
$ nssh list --help
List port mappings for specified subscriber

Usage:
  nssh list <subscriber name> [flags]

Aliases:
  list, l

Flags:
  -h, --help   help for list

Global Flags:
      --coverage-type string   Specify coverage type, "g" for Global, "jp" for Japan
      --profile-name string    Specify SORACOM CLI profile name (default "ssh")
```

## References

- Japanese
  - [IoT プラットフォーム 株式会社ソラコム](https://soracom.jp/)
  - [SORACOM Napter とは | ユーザーガイド | SORACOM Developers](https://dev.soracom.io/jp/napter/what-is-napter/)
- English
  - [Soracom | Cellular IoT Cloud Connectivity](https://www.soracom.io/)
  - [Soracom Napter Overview | SORACOM Developers](https://developers.soracom.io/en/docs/napter/)

## License

MIT. See [LICENSE](LICENSE) for details.

## Privacy

This program will send requests to following services:

- https://checkip.amazonaws.com/, to determine your global IP address.
- https://g.api.soracom.io (Global coverage) or https://api.soracom.io (Japan coverage), to use SORACOM services

Other than that, the program does not send user action/data to any server. Please consult each provider's privacy notices.

- [AWS Privacy Notice](https://aws.amazon.com/privacy/)
- [Privacy Policy | Soracom](https://www.soracom.io/privacy-policy/)

---

**Enjoy remote connection with SORACOM Napter!**
