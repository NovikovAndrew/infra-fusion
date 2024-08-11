package healthchecker

import (
	"encoding/json"
	"net/http"
	"sync"
)

type Check func() error

type CheckerHandler struct {
	checksMutex     sync.RWMutex
	livenessChecks  map[string]Check
	readinessChecks map[string]Check
}

func (s *CheckerHandler) LiveEndpoint(w http.ResponseWriter, r *http.Request) {
	s.handle(w, r, s.livenessChecks)
}

func (s *CheckerHandler) ReadyEndpoint(w http.ResponseWriter, r *http.Request) {
	s.handle(w, r, s.readinessChecks, s.livenessChecks)
}

func (s *CheckerHandler) AddLivenessCheck(name string, check Check) {
	s.checksMutex.Lock()
	s.livenessChecks[name] = check
	s.checksMutex.Unlock()
}

func (s *CheckerHandler) AddReadinessCheck(name string, check Check) {
	s.checksMutex.Lock()
	s.readinessChecks[name] = check
	s.checksMutex.Unlock()
}

func (s *CheckerHandler) collectChecks(checks map[string]Check, resultsOut map[string]string, statusOut *int) {
	s.checksMutex.RLock()
	defer s.checksMutex.RUnlock()
	for name, check := range checks {
		if err := check(); err != nil {
			*statusOut = http.StatusServiceUnavailable
			resultsOut[name] = err.Error()
		} else {
			resultsOut[name] = "OK"
		}
	}
}

func (s *CheckerHandler) handle(w http.ResponseWriter, r *http.Request, checks ...map[string]Check) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	checkResults := make(map[string]string)
	status := http.StatusOK
	for _, checks := range checks {
		s.collectChecks(checks, checkResults, &status)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if r.URL.Query().Get("full") != "1" {
		w.Write([]byte("{}\n"))
		return
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "    ")
	encoder.Encode(checkResults)
}
