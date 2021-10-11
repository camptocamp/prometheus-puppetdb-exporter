package exporter

import (
	"fmt"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"

	"github.com/camptocamp/prometheus-puppetdb-exporter/internal/puppetdb"
)

// Exporter implements the prometheus.Exporter interface, and exports PuppetDB metrics
type Exporter struct {
	client         *puppetdb.PuppetDB
	namespace      string
	unreportedNode time.Duration
	categories     map[string]struct{}
	metrics        map[string]*prometheus.GaugeVec
	interval       time.Duration
}

// NewPuppetDBExporter returns a new exporter of PuppetDB metrics.
func NewPuppetDBExporter(url, certPath, caPath, keyPath string, sslSkipVerify bool, unreportedNode time.Duration, categories map[string]struct{}, interval time.Duration) (e *Exporter, err error) {
	e = &Exporter{
		namespace:      "puppetdb",
		unreportedNode: unreportedNode,
		categories:     categories,
	}

	opts := &puppetdb.Options{
		URL:        url,
		CertPath:   certPath,
		CACertPath: caPath,
		KeyPath:    keyPath,
		SSLVerify:  sslSkipVerify,
	}

	e.client, err = puppetdb.NewClient(opts)
	if err != nil {
		log.Fatalf("failed to create new client: %s", err)
		return
	}

	if interval > 0 {
		go func() {
			for {
				err := e.scrape()
				if err != nil {
					log.Errorf("failed to scrape metrics: %s", err)
				} else {
					log.Info("scraped metrics")
				}

				time.Sleep(interval)
			}
		}()
	}

	return
}

// Describe outputs PuppetDB metric descriptions
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(e, ch)
}

// Collect fetches new metrics from the PuppetDB and updates the appropriate metrics
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	if e.interval == 0 {
		err := e.scrape()
		if err != nil {
			log.Errorf("failed to scrape metrics: %s", err)
			return
		}

		log.Info("scraped metrics")
	}

	for _, m := range e.metrics {
		m.Collect(ch)
	}
}

// scrape scrapes PuppetDB and update metrics
func (e *Exporter) scrape() (err error) {
	e.metrics = map[string]*prometheus.GaugeVec{}

	e.metrics["node_report_status_count"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: e.namespace,
		Name:      "node_report_status_count",
		Help:      "Total count of reports status by type",
	}, []string{"status"})

	for category := range e.categories {
		metricName := fmt.Sprintf("report_%s", category)
		e.metrics[metricName] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "puppet",
			Name:      metricName,
			Help:      fmt.Sprintf("Total count of %s per status", category),
		}, []string{"name", "environment", "host"})
	}

	e.metrics["report"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "puppet",
		Name:      "report",
		Help:      "Timestamp of latest report",
	}, []string{"environment", "host", "deactivated"})

	statuses := make(map[string]int)

	nodes, err := e.client.Nodes()
	if err != nil {
		log.Errorf("failed to get nodes: %s", err)
	}

	for _, node := range nodes {
		var deactivated string
		if node.Deactivated == "" {
			deactivated = "false"
		} else {
			deactivated = "true"
		}

		if node.ReportTimestamp == "" {
			if deactivated == "false" {
				statuses["unreported"]++
			}
			continue
		}
		latestReport, err := time.Parse("2006-01-02T15:04:05Z", node.ReportTimestamp)
		if err != nil {
			if deactivated == "false" {
				statuses["unreported"]++
			}
			log.Errorf("failed to parse report timestamp: %s", err)
			continue
		}
		e.metrics["report"].With(prometheus.Labels{"environment": node.ReportEnvironment, "host": node.Certname, "deactivated": deactivated}).Set(float64(latestReport.Unix()))

		if deactivated == "false" {
			if latestReport.Add(e.unreportedNode).Before(time.Now()) {
				statuses["unreported"]++
			} else if node.LatestReportStatus == "" {
				statuses["unreported"]++
			} else {
				statuses[node.LatestReportStatus]++
			}
		}

		if node.LatestReportHash != "" {
			reportMetrics, _ := e.client.ReportMetrics(node.LatestReportHash)
			for _, reportMetric := range reportMetrics {
				_, ok := e.categories[reportMetric.Category]
				if ok {
					category := fmt.Sprintf("report_%s", reportMetric.Category)
					e.metrics[category].With(prometheus.Labels{"name": strings.ReplaceAll(strings.Title(reportMetric.Name), "_", " "), "environment": node.ReportEnvironment, "host": node.Certname}).Set(reportMetric.Value)
				}
			}
		}
	}

	for statusName, statusValue := range statuses {
		e.metrics["node_report_status_count"].With(prometheus.Labels{"status": statusName}).Set(float64(statusValue))
	}

	return nil
}
