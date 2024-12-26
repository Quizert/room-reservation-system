package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func SetupMetricsRoute() *http.ServeMux {
	prometheus.MustRegister(HttpRequestsTotal, HttpRequestDuration, DbQueriesTotal, DbQueryDuration)
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	return mux
}
