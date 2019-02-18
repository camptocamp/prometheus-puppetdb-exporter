Prometheus PuppetDB exporter
============================

## Usage

```
Usage:
  prometheus-puppetdb-exporter [OPTIONS]

Application Options:
      --version          Show version.
  -u, --puppetdb-url=    PuppetDB base URL. [$PUPPETDB_URL]
      --cert-file=       A PEM encoded certificate file. [$PUPPETDB_CERT_FILE]
      --key-file=        A PEM encoded private key file. [$PUPPETDB_KEY_FILE]
      --ca-file=         A PEM encoded CA's certificate. [$PUPPETDB_CA_FILE]
      --ssl-skip-verify  Skip SSL verification. [$PUPPETDB_SSL_SKIP_VERIFY]
      --scrape-interval= Duration between two scrapes. (default: 5s) [$PUPPETDB_SCRAPE_INTERVAL]
      --listen-address=  Address to listen on for web interface and telemetry. (default: 0.0.0.0:9121) [$PUPPETDB_LISTEN_ADDRESS]
      --metric-path=     Path under which to expose metrics. (default: /metrics) [$PUPPETDB_METRIC_PATH]
      --verbose          Enable debug mode [$PUPPETDB_VERBOSE]

Help Options:
  -h, --help             Show this help message
```
