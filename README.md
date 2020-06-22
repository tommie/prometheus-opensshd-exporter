# Prometheus Exporter for OpenSSHd

[![Build Status](https://travis-ci.org/tommie/prometheus-opensshd-exporter.svg?branch=master)](https://travis-ci.org/tommie/prometheus-opensshd-exporter)

Tails Systemd logs of the `ssh.service` to find interesting lines.

## Metrics

* `opensshd_auth_results_total{method="password",result="failed",valid_user="1"}`
  The number of authentication attempts.

## Usage

```
go install ./cmd/...
$GOPATH/bin/openssd_exporter
```

See
[`docker-compose.yml`](https://github.com/tommie/prometheus-opensshd-exporter/blob/master/docker-compose.yml)
for an example of how this can be used in Docker.
