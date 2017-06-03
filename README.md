# autohosts

Use the credentials in `~/.aws/credentials` to pull the hostnames from all
instances in AWS and soon, GCP.

## Installation

You can download a binary for your platform from the [Releases](https://github.com/damselem/autohosts/releases) section.

## Usage

Prints hostnames to `STDOUT`:

```
autohosts
```

Updates `/etc/hosts` with hostnames:
```
autohosts -o /etc/hosts
```
