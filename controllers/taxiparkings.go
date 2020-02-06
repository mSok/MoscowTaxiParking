package controllers

import (
	"app/models"
	"app/parser"
	"app/utils"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

var (
	requestsCounter         prometheus.Counter
	requestProcessingTimeMs prometheus.Summary
	responseStatusCounter   *prometheus.CounterVec
)

func init() {
	requestsCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "requests_counters",
		})
	requestProcessingTimeMs = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name:       "request_processing_time_summary_ms",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		})
	responseStatusCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "response_status",
		},
		[]string{"status"},
	)
	prometheus.MustRegister(requestsCounter)
	prometheus.MustRegister(requestProcessingTimeMs)
	prometheus.MustRegister(responseStatusCounter)
}

// PrometheusMiddleware middleware for register metrics
var PrometheusMiddleware = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// inc all API requests
		var strtTime time.Time
		if strings.HasPrefix(strings.ToLower(r.URL.Path), "/api/") {
			requestsCounter.Inc()
			strtTime = time.Now()
		}
		lrw := logResponseWriter(w)
		next.ServeHTTP(lrw, r)
		if strings.HasPrefix(strings.ToLower(r.URL.Path), "/api/") && !strtTime.IsZero() {
			responseStatusCounter.With(prometheus.Labels{"status": strconv.Itoa(lrw.statusCode)}).Inc()
			log.Printf("Response with status: %d", lrw.statusCode)
			obs := float64(time.Since(strtTime).Milliseconds())
			log.Printf("processing_time: %f", obs)
			requestProcessingTimeMs.Observe(obs)
		}
	})
}

// logResponseWriter writer with custom WriteHeader
func logResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	// WriteHeader(int) is not called if our response implicitly returns 200 OK, so
	// we default to that status code.
	return &loggingResponseWriter{w, http.StatusOK}
}

// WriteHeader writer with save response http code
func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// MetricsGetHandler handler for get metrics
var MetricsGetHandler = func(w http.ResponseWriter, r *http.Request) {
	promhttp.Handler().ServeHTTP(w, r)
}

// GetTaxiParkingsByID handler for get data by id
var GetTaxiParkingsByID = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	log.Printf("[GET]: %s", params["id"])
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	client := models.GetDB()
	data, err := client.GetTaxiParking(id)
	if data == "" || (err != nil && err.Error() == "redis: nil") {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write([]byte(data))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// LoadTaxiParkings handler for load open data
var LoadTaxiParkings = func(w http.ResponseWriter, r *http.Request) {

	var conf = utils.GetConf()
	go func() {
		cnt, err := parser.LoadFromSource(conf.Source)
		if err != nil {
			log.Printf("[Error] parsed data %s\n", err)
		} else {
			log.Printf("[LoadFromSource] loaded %d recs\n", cnt)
		}

	}()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
