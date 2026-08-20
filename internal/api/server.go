package api

import (
	"net/http"

	"example.com/certvault"
)

// Options HTTP 服务选项。
type Options struct {
	WebDir    string
	AllowCORS bool
}

// Server 管理 HTTP。
type Server struct {
	vault *certvault.Vault
	opts  Options
	mux   *http.ServeMux
}

// New 构造 Server。
func New(v *certvault.Vault, opts Options) *Server {
	s := &Server{vault: v, opts: opts, mux: http.NewServeMux()}
	s.routes()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.opts.AllowCORS {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,DELETE,OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	s.mux.ServeHTTP(w, r)
}
