Prometheus PuppetDB exporter
============================

## Usage

```
Usage:
  prometheus-puppetdb-exporter [OPTIONS]

Application Options:
      --version          Show version.
  -u, --puppetdb-url=    PuppetDB base URL. (default: https://puppetdb:8081/pdb/query) [$PUPPETDB_URL]
      --cert-file=       A PEM encoded certificate file. [$PUPPETDB_CERT_FILE]
      --key-file=        A PEM encoded private key file. [$PUPPETDB_KEY_FILE]
      --ca-file=         A PEM encoded CA's certificate. [$PUPPETDB_CA_FILE]
      --ssl-skip-verify  Skip SSL verification. [$PUPPETDB_SSL_SKIP_VERIFY]
      --scrape-interval= Duration between two scrapes. Set to 0 for synchronous scrapes. (default: 5s) [$PUPPETDB_SCRAPE_INTERVAL]
      --listen-address=  Address to listen on for web interface and telemetry. (default: 0.0.0.0:9635) [$PUPPETDB_LISTEN_ADDRESS]
      --metric-path=     Path under which to expose metrics. (default: /metrics) [$PUPPETDB_METRIC_PATH]
      --verbose          Enable debug mode [$PUPPETDB_VERBOSE]
      --unreported-node= Tag nodes as unreported if the latest report is older than the defined duration. (default: 2h) [$PUPPETDB_UNREPORTED_NODE]
      --categories=      Report metrics categories to scrape. (default: resources,time,changes,events) [$REPORT_METRICS_CATEGORIES]

Help Options:
  -h, --help             Show this help message
```

## Metrics

```
# HELP puppetdb_exporter_build_info puppetdb exporter build informations
# TYPE puppetdb_exporter_build_info gauge
puppetdb_exporter_build_info{build_date="2019-02-18",commit_sha="XXXXXXXXXX",golang_version="go1.11.4",version="1.0.0"} 1
# HELP puppetdb_node_report_status_count Total count of reports status by type
# TYPE puppetdb_node_report_status_count gauge
puppetdb_node_report_status_count{status="changed"} 1
puppetdb_node_report_status_count{status="failed"} 1
puppetdb_node_report_status_count{status="unchanged"} 1
```

## Development

To locally run the Puppet environment with `docker-compose`, run the following commands:

```shell
docker-compose up -d puppetdb && docker-compose logs -f puppetdb

docker-compose run agent

docker-compose up -d exporter && docker-compose logs -f exporter
curl "$(docker-compose port exporter 9635)/metrics"
```

If you want to use go run, use the following command:

```shell
PUPPETDB_URL="http://$(docker-compose port puppetdb 8080)/pdb/query" go run .
```
