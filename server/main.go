package server

import (
	"fmt"
	"io"
	"net/http"

	"github.com/akerl/metrics/metrics"

	"github.com/akerl/timber/v2/log"
)

var logger = log.NewLogger("metrics.server")

// Cache shares a MetricSet between a writer and a reader
type Cache struct {
	MetricSet metrics.MetricSet
}

// Server defines a Prometheus-compatible metrics engine
type Server struct {
	Port  int
	Cache *Cache
}

// NewServer creates a new Server object
func NewServer(port int, cache *Cache) *Server {
	return &Server{
		Port:  port,
		Cache: cache,
	}
}

// Run starts the Server object in the foreground
func (s *Server) Run() error {
	bindStr := fmt.Sprintf(":%d", s.Port)
	logger.DebugMsgf("binding metrics server to %s", bindStr)
	mux := http.NewServeMux()
	mux.HandleFunc("/metrics", s.handleMetrics)
	return http.ListenAndServe(bindStr, mux)
}

func (s *Server) handleMetrics(w http.ResponseWriter, _ *http.Request) {
	if !s.Cache.MetricSet.Validate() {
		logger.DebugMsg("invalid metrics file requested")
		http.Error(w, "invalid metrics file", http.StatusInternalServerError)
	} else {
		logger.DebugMsg("successful metrics request")
		io.WriteString(w, s.Cache.MetricSet.String())
	}
}
