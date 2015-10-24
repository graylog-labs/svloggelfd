# Runit GELF integration
Sends service log messages as UDP/GELF to a Graylog server

## Installation

Copy `svloggelfd` binary in the same directory `svlogd` is installed.

```
  $ sudo cp svloggelfd /usr/bin/
```

Edit `log/run` for the service you want to enable GELF log forwarding

```
  $ vi /etc/sv/<service name>/log/run
```

Replace `svlogd` with `svloggelfd` for message forwarding

```shell
#!/bin/sh
exec svloggelfd -H 127.0.0.1:12201 -s elasticsearch
```

Combine `svlogd` and `svloggelfd` for disk storage _and_ GELF forwarding

```shell
#!/bin/sh
exec svloggelfd -H 127.0.0.1:12201 -s elasticsearch -e | svlogd -tt /var/log/graylog/elasticsearch
```

## Options

| Parameter | Argument       | Description                                      |
|-----------|----------------|--------------------------------------------------|
| -H        | GELF Host      | Load GELF logging module                         |
| -s        | Source field   | Overwrite source field of GELF message           |
| -f        | Facility field | Overwrite facility field of GELF message         |
| -t        | Tag field      | Add an additional tag to every message           |
| -e        | -              | Echo message back to STDOUT for command chaining |

## Build

```
  $ go get github.com/codegangsta/cli
  $ go get github.com/Graylog2/go-gelf/gelf
  $ go build -o svloggelfd main.go
```
