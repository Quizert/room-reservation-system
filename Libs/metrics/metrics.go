package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	// HttpRequestsTotal Счётчик запросов
	HttpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	// HttpRequestDuration Гистограмма времени обработки запросов
	HttpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	// DbQueriesTotal Счётчик запросов к базе данных
	DbQueriesTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "db_queries_total",
			Help: "Total number of database queries executed",
		},
		[]string{"operation", "status"},
	)

	// DbQueryDuration Гистограмма для времени выполнения запросов
	DbQueryDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_query_duration_seconds",
			Help:    "Duration of database queries in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"}, // Лейбл: тип операции (SELECT, INSERT и т.д.)
	)
)

// RecordHttpMetrics Функция для записи метрик HTTP-запросов
func RecordHttpMetrics(method, endpoint, status string, duration float64) {
	HttpRequestsTotal.WithLabelValues(method, endpoint, status).Inc()
	HttpRequestDuration.WithLabelValues(method, endpoint).Observe(duration)
}

func RecordDataBaseMetrics(operation string, status string, duration float64) {
	DbQueriesTotal.WithLabelValues(operation, status).Inc()
	DbQueryDuration.WithLabelValues(operation).Observe(duration)
}
