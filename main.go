package main

import (
	"flag"
	//"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/camptocamp/prometheus-puppetdb-exporter/internal/exporter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var (
	puppetDBURL    = flag.String("puppetdb.url", getEnv("PUPPETDB_URL", "http://puppetdb:8080"), "PuppetDB base URL. (default: http://puppetdb:8080)")
	certFile       = flag.String("puppetdb.cert-file", getEnv("CERT_FILE", "certs/client.pem"), "A PEM encoded certificate file. (default: certs/client.pem)")
	keyFile        = flag.String("puppetdb.key-file", getEnv("KEY_FILE", "certs/client.key"), "A PEM encoded private key file. (default: certs/client.key)")
	caFile         = flag.String("puppetdb.ca-file", getEnv("CA_FILE", "certs/cacert.pem"), "A PEM encoded CA's certificate file. (default: certs/cacert.pem)")
	sslVerify      = flag.Bool("puppetdb.ssl-verify", false, "Skip SSL verification")
	scrapeInterval = flag.String("scrape-interval", getEnv("PUPPETDB_SCRAPE_INTERVAL", "5s"), "Interval between two scrapes on the PuppetDB. (default: 5s)")
	listenAddress  = flag.String("web.listen-address", ":9121", "Address to listen on for web interface and telemetry.")
	metricPath     = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	isDebug        = flag.Bool("debug", false, "Output verbose debug information")
	showVersion    = flag.Bool("version", false, "Show version information and exit")
	logFormat      = flag.String("log-format", "txt", "Log format, valid options are txt and json")

	// VERSION, BUILD_DATE, GIT_COMMIT are filled in by the build script
	version    = "<<< filled in by build >>>"
	buildDate  = "<<< filled in by build >>>"
	commitSha1 = "<<< filled in by build >>>"
)

func getEnv(key string, defaultVal string) string {
	if envVal, ok := os.LookupEnv(key); ok {
		return envVal
	}
	return defaultVal
}

func main() {
	flag.Parse()

	switch *logFormat {
	case "json":
		log.SetFormatter(&log.JSONFormatter{})
	default:
		log.SetFormatter(&log.TextFormatter{})
	}
	log.Printf("PuppetDB Metrics Exporter %s    build date: %s    sha1: %s    Go: %s",
		version, buildDate, commitSha1,
		runtime.Version(),
	)
	if *isDebug {
		log.SetLevel(log.DebugLevel)
		log.Debugln("Enabling debug output")
	} else {
		log.SetLevel(log.InfoLevel)
	}

	if *showVersion {
		return
	}

	interval, err := time.ParseDuration(*scrapeInterval)
	if err != nil {
		log.Fatalf("failed to parse scrape interval duration: %s", err)
	}

	exp, err := exporter.NewPuppetDBExporter(*puppetDBURL, *certFile, *caFile, *keyFile, *sslVerify)
	if err != nil {
		log.Fatalf("failed to initialize exporter: %s", err)
	}

	go exp.Scrape(interval)

	buildInfo := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "puppetdb_exporter_build_info",
		Help: "puppetdb exporter build_info",
	}, []string{"version", "commit_sha", "build_date", "golang_version"})
	buildInfo.WithLabelValues(version, commitSha1, buildDate, runtime.Version()).Set(1)
	prometheus.MustRegister(buildInfo)

	http.Handle(*metricPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
<html>
<head><title>Prometheus PuppetDB Exporter v` + version + `</title></head>
<body>
<h1>Prometheus PuppetDB Exporter ` + version + `</h1>
<p><a href='` + *metricPath + `'>Metrics</a></p>
</body>
</html>
						`))
	})

	log.Infof("Providing metrics at %s%s", *listenAddress, *metricPath)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
