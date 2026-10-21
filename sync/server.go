package sync

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/shakilnwz/wallet/core/models"
	"github.com/shakilnwz/wallet/db"
)

type Server struct {
	database *db.DB
	listener net.Listener
	server   *http.Server
	port     int
	mu       sync.Mutex
	running  bool
}

func NewServer(database *db.DB) *Server {
	return &Server{database: database}
}

func (s *Server) Start(port int) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return s.port, nil
	}

	addr := fmt.Sprintf(":%d", port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		// Fallback to random free port
		ln, err = net.Listen("tcp", ":0")
		if err != nil {
			return 0, err
		}
	}

	s.listener = ln
	s.port = ln.Addr().(*net.TCPAddr).Port

	mux := http.NewServeMux()
	mux.HandleFunc("/sync/ping", s.handlePing)
	mux.HandleFunc("/sync/package", s.handleGetPackage)
	mux.HandleFunc("/sync/merge", s.handleMerge)

	s.server = &http.Server{
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	s.running = true
	go func() {
		_ = s.server.Serve(ln)
	}()

	return s.port, nil
}

func (s *Server) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return nil
	}
	s.running = false
	return s.server.Close()
}

func (s *Server) Port() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.port
}

func (s *Server) IsRunning() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.running
}

func (s *Server) handlePing(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"device_id": s.database.DeviceID(),
		"status":    "ready",
		"timestamp": time.Now().UnixMilli(),
	})
}

func (s *Server) handleGetPackage(w http.ResponseWriter, r *http.Request) {
	pkg, err := s.database.ExportSyncPackage()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(pkg)
}

func (s *Server) handleMerge(w http.ResponseWriter, r *http.Request) {
	var pkg models.SyncPackage
	if err := json.NewDecoder(r.Body).Decode(&pkg); err != nil {
		http.Error(w, "invalid payload: "+err.Error(), http.StatusBadRequest)
		return
	}

	ins, upd, skip, err := s.database.ImportAndMergeSyncPackage(&pkg)
	if err != nil {
		http.Error(w, "merge error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"inserted": ins,
		"updated":  upd,
		"skipped":  skip,
	})
}
