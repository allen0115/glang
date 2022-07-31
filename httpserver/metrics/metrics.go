package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"time"
)

const (
	MetricsNamespace = "ben_httpserver"
)

var (
	functionLatency = CreateExecutionTimeMetric(MetricsNamespace, "Time spent")
)

type ExecutionTimer struct {
	histo *prometheus.HistogramVec
	start time.Time
	last  time.Time
}

func Register() {
	log.Print("register function latency prometheus metric, %v", functionLatency)
	err := prometheus.Register(functionLatency)
	if err != nil {
		log.Printf("error register: %v", err)
	} else {
		log.Print(err)
	}

}

func NewTimer() *ExecutionTimer {
	return NewExecutionTimer(functionLatency)
}

func NewExecutionTimer(histo *prometheus.HistogramVec) *ExecutionTimer {
	now := time.Now()
	return &ExecutionTimer{
		histo: histo,
		start: now,
		last:  now,
	}
}

func (t *ExecutionTimer) ObserveTotal() {
	log.Print("report execution time after function done")
	(*t.histo).WithLabelValues("total").Observe(time.Now().Sub(t.start).Seconds())
}
func CreateExecutionTimeMetric(nameSpace string, help string) *prometheus.HistogramVec {
	log.Printf("CreateExecutionTimeMetric, %s, %s", nameSpace, help)
	return prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: nameSpace,
			Name:      "execution_latency_seconds",
			Help:      help,
			Buckets:   prometheus.ExponentialBuckets(0.001, 2, 15),
		}, []string{"step"},
	)
}
