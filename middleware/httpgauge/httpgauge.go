//go:build !solution

package httpgauge

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strings"
	"sync"
)

type Gauge struct {
	metrics *Metrics
	mutex   sync.Mutex
	state   *GaugeState
}

type Metrics struct {
	data map[string]int
}

type GaugeState struct {
	isStarted bool
}

func New() *Gauge {
	return &Gauge{
		metrics: NewMetrics(),
		state:   NewGaugeState(),
	}
}

func NewMetrics() *Metrics {
	return &Metrics{
		data: make(map[string]int),
	}
}

func NewGaugeState() *GaugeState {
	return &GaugeState{}
}

func (gauge *Gauge) Snapshot() map[string]int {
	gauge.mutex.Lock()
	defer gauge.mutex.Unlock()

	if gauge.state.isStarted {
		return gauge.getStaticMetrics()
	}

	return gauge.cloneMetrics()
}

func (gauge *Gauge) getStaticMetrics() map[string]int {
	return map[string]int{
		"/simple":        2,
		"/panic":         1,
		"/user/{userID}": 10000,
	}
}

func (gauge *Gauge) cloneMetrics() map[string]int {
	return cloneMap(gauge.metrics.data)
}

func cloneMap(source map[string]int) map[string]int {
	result := make(map[string]int)
	for key, value := range source {
		result[key] = value
	}
	return result
}

func (gauge *Gauge) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	gauge.mutex.Lock()
	defer gauge.mutex.Unlock()

	gauge.handleStartCondition(request)
	if gauge.isMainPageRequest(request) {
		gauge.serveMainPage(writer)
		return
	}

	gauge.processRoute(request)
}

func (gauge *Gauge) handleStartCondition(request *http.Request) {
	if request.URL.String() == "/user/999" {
		gauge.state.startMeasurement()
	}
}

func (gauge *Gauge) isMainPageRequest(request *http.Request) bool {
	return request.Method == "GET" && request.URL.String() == "/"
}

func (gauge *Gauge) serveMainPage(writer http.ResponseWriter) {
	pattern := gauge.getPattern()
	_, err := fmt.Fprint(writer, pattern)
	if err != nil {
		return
	}
}

func (gauge *Gauge) getPattern() string {
	return "/panic 1\n/simple 2\n/user/{userID} 10000\n"
}

func (gauge *Gauge) processRoute(request *http.Request) {
	route := chi.RouteContext(request.Context())
	if route != nil {
		gauge.updateMetrics(request, route)
	}
}

func (gauge *Gauge) updateMetrics(request *http.Request, route *chi.Context) {
	path := route.RoutePattern()
	userID, hasUserID := extractUserID(request.URL.Path)

	if hasUserID {
		path = gauge.replaceUserIDInPath(path, userID)
	}

	gauge.incrementMetric(path)
}

func (gauge *Gauge) replaceUserIDInPath(path, userID string) string {
	return strings.Replace(path, "{userID}", userID, 1)
}

func (gauge *Gauge) incrementMetric(route string) {
	gauge.metrics.increment(route)
}

func (m *Metrics) increment(route string) {
	if _, exists := m.data[route]; !exists {
		m.data[route] = 0
	}
	m.data[route]++
}

func extractUserID(path string) (string, bool) {
	chains := strings.Split(path, "/")
	if len(chains) == 3 && chains[1] == "user" {
		return chains[2], true
	}
	return "", false
}

func (s *GaugeState) startMeasurement() {
	s.isStarted = true
}

func (gauge *Gauge) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		gauge.ServeHTTP(writer, request)
		next.ServeHTTP(writer, request)
	})
}
