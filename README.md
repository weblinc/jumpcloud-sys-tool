# jcsystool

Go binary for performing [JumpCloud SystemContext API][1] calls.

## Install

Download and install the binary on a system where the JumpCloud agent is installed.

```bash
sudo su -
wget https://github.com/weblinc/jumpcloud-sys-tool/releases/download/1.0.0/jcsystool-linux-amd64-1.0.0.tar.gz
tar -xvf jcsystool-linux-amd64-1.0.0.tar.gz
mv jcsystool-linux-amd64-1.0.0 /usr/bin/jcsystool && chmod 700 /usr/bin/jcsystool
```

## Usage

```
$ jcsystool --help
jcsystool - JumpCloud System Tool 0.0.1
Usage:
  jcsystool [OPTIONS]

Application Options:
  -X, --action= HTTP method to use e.g. GET/PUT/DELETE
  -J, --json=   JSON string to use for PUT actions to system API. Alternatively, use STDIN.

Help Options:
  -h, --help    Show this help message
```

`jcsystool` uses the system's configuration file and client key to authenticate with the SystemContext API.

It expects these files to be located at `/opt/jc/jcagent.conf` and `/opt/jc/client.key`.

Alternatively, you can use `JC_CLIENT_KEY_PATH` and `JC_CONFIG_PATH` environment variables to specify where
to find these files.

The SystemContext API can only be used to read or modify information about the system its called from.
There is only one endpoint `/systems/<random-system-guid>` which accepts `GET`, `PUT`, and `DELETE`.

## Examples

```bash
# Get system Information
$ jcsystool -X GET

# Delete system from JumpCloud
$ jcsystool -X DELETE

# Update system configuration via -J flag
$ jcsystool -X PUT -J '{ "displayName": "mySystem" }'

# Update system configuration via STDIN
$ echo '{ "tags": ["staging-server"] }' | jcsystool -X PUT
```

[1]: https://github.com/TheJumpCloud/SystemContextAPI
