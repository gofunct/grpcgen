package runtime

import (
	"github.com/heptiolabs/healthcheck"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/spf13/viper"
	"net/http"
	"net/http/pprof"
	"time"
)

func NewMux() *http.ServeMux {
	mux := http.NewServeMux()
	check := healthcheck.NewMetricsHandler(prometheus.DefaultRegisterer, "grpc")
	check.AddLivenessCheck("goroutine_threshold", healthcheck.GoroutineCountCheck(viper.GetInt("gprc.routine_threshold")))
	mux.HandleFunc("/live", check.LiveEndpoint)
	log.Info("liveness handler registered-->", "/live")
	check.AddReadinessCheck("db_health_check", healthcheck.TCPDialCheck(viper.GetString("gprc.db_port"), 1*time.Second))
	mux.HandleFunc("/ready", check.ReadyEndpoint)
	log.Info("readiness handler registered-->", "/ready")

	mux.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	mux.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	mux.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	mux.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	mux.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
	log.Info("pprof handler registered-->", "/debug/pprof")

	mux.Handle("/metrics", promhttp.Handler())
	log.Info("metrics handler registered-->", "/metrics")

	return mux
}
