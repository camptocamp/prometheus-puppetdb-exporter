FROM golang:1.11 as builder
WORKDIR /go/src/github.com/camptocamp/prometheus-puppetdb-exporter
COPY . .
RUN make prometheus-puppetdb-exporter

FROM scratch
COPY --from=builder /go/src/github.com/camptocamp/prometheus-puppetdb-exporter/prometheus-puppetdb-exporter /
ENTRYPOINT ["/prometheus-puppetdb-exporter"]
CMD [""]
