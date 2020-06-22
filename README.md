# Prometheus Exporter for OpenSSHd

Tails Systemd logs of the `ssh.service` to find interesting lines.

## Metrics

* `opensshd_auth_results_total{method="password",result="failed",valid_user="1"}`
  The number of authentication attempts.
