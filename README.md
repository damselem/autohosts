# autohosts

Use the credentials in `~/.aws/credentials` to pull the hostnames from all
instances in AWS and soon, GCP.

## Usage

Prints hostnames to `STDOUT`:

```
autohosts
```

Updates `/etc/hosts` with hostnames:
```
autohosts -o /etc/hosts
```
