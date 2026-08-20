package api

func (s *Server) routes() {
	s.mux.HandleFunc("GET /api/health", s.handleHealth)
	s.mux.HandleFunc("GET /api/stats", s.handleStats)
	s.mux.HandleFunc("GET /api/certs", s.handleList)
	s.mux.HandleFunc("POST /api/certs/import", s.handleImport)
	s.mux.HandleFunc("GET /api/certs/{id}", s.handleGet)
	s.mux.HandleFunc("POST /api/certs/{id}/revoke", s.handleRevoke)
	s.mux.HandleFunc("POST /api/scan", s.handleScan)
	s.mux.HandleFunc("POST /api/rotate/preview", s.handleRotatePreview)
	s.mux.HandleFunc("GET /api/bundle", s.handleBundle)
	if s.opts.WebDir != "" {
		s.mux.Handle("/", staticHandler(s.opts.WebDir))
	}
}
