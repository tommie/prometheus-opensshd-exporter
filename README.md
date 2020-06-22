# Prometheus Exporter for OpenSSHd

[![Build Status](https://travis-ci.org/tommie/prometheus-opensshd-exporter.svg?branch=master)](https://travis-ci.org/tommie/prometheus-opensshd-exporter)

Tails Systemd logs of the `ssh.service` to find interesting lines.

## Metrics

* `opensshd_auth_results_total{method="password",result="failed",valid_user="1"}`
  The number of authentication attempts.
