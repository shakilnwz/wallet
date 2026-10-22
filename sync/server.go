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

// apiPrefix is versioned because v2 sends amounts as integer paisa
// ("amount_paisa"). Older builds speak /sync/* with float amounts; giving them
// a 404 makes mixed versions fail loudly instead of merging zero amounts.
const apiPrefix = "/sync/v2"

type Server struct {
	// OnMerge, if set, is called after a peer pushes records into this device.
	// It runs on the HTTP goroutine; set it before Start.
	OnMerge func()

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
	mux.HandleFunc(apiPrefix+"/ping", s.handlePing)
	mux.HandleFunc(apiPrefix+"/package", s.handleGetPackage)
	mux.HandleFunc(apiPrefix+"/merge", s.handleMerge)

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

	if s.OnMerge != nil && ins+upd > 0 {
		s.OnMerge()
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"inserted": ins,
		"updated":  upd,
		"skipped":  skip,
	})
}

// LocalIP returns this device's LAN IPv4 address, or "" if offline. It asks
// the kernel which interface would route to a public address; no packet is
// sent. Unlike net.InterfaceAddrs this also works on Android 11+, which
// blocks the netlink calls that function needs.
func LocalIP() string {
	conn, err := net.Dial("udp4", "192.0.2.1:9") // TEST-NET-1, never routed
	if err != nil {
		return ""
	}
	defer conn.Close()
	addr, ok := conn.LocalAddr().(*net.UDPAddr)
	if !ok || addr.IP.IsLoopback() {
		return ""
	}
	return addr.IP.String()
}
