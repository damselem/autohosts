# autohosts

Uses the credentials in `~/.aws/credentials` and `~/.config/gcloud` to pull the
hostnames from all instances in AWS and GCP.

## Installation

You can download a binary for your platform from the [Releases](https://github.com/damselem/autohosts/releases) section.

## Usage

Prints hostnames to `STDOUT`:

```
autohosts all
```

Updates `/etc/hosts` with hostnames:
```
autohosts aws -o /etc/hosts
```
