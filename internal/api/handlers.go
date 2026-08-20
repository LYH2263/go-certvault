package api

import (
	"net/http"
	"time"

	"example.com/certvault"
)

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "closed": s.vault.Closed()})
}

func (s *Server) handleStats(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, s.vault.Stats())
}

func (s *Server) handleList(w http.ResponseWriter, _ *http.Request) {
	views := s.vault.ListPublic()
	type pub struct {
		ID          string    `json:"id"`
		Name        string    `json:"name"`
		Subject     string    `json:"subject"`
		NotAfter    time.Time `json:"not_after"`
		Active      bool      `json:"active"`
		Revoked     bool      `json:"revoked"`
		Fingerprint string    `json:"fingerprint"`
		KeyLen      int       `json:"key_len"`
	}
	out := make([]pub, 0, len(views))
	for _, v := range views {
		out = append(out, pub{
			ID: v.ID, Name: v.Name, Subject: v.Subject, NotAfter: v.NotAfter,
			Active: v.Active, Revoked: v.Revoked, Fingerprint: v.Fingerprint, KeyLen: v.KeyLen,
		})
	}
	writeJSON(w, http.StatusOK, map[string]any{"certs": out})
}

type importBody struct {
	Name        string   `json:"name"`
	CertPEM     string   `json:"cert_pem"`
	KeyPEM      string   `json:"key_pem"`
	Tags        []string `json:"tags"`
	Description string   `json:"description"`
	Active      bool     `json:"active"`
}

func (s *Server) handleImport(w http.ResponseWriter, r *http.Request) {
	var in importBody
	if err := readJSON(r, &in); err != nil {
		writeErr(w, err)
		return
	}
	id, err := s.vault.ImportPEM([]byte(in.CertPEM), []byte(in.KeyPEM), certvault.Meta{
		Name: in.Name, Tags: in.Tags, Description: in.Description, Active: in.Active,
	})
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"id": id})
}

func (s *Server) handleGet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	ev, err := s.vault.Get(id)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"id": ev.ID, "name": ev.Name, "subject": ev.Subject, "not_after": ev.NotAfter,
		"active": ev.Active, "revoked": ev.Revoked, "fingerprint": ev.Fingerprint,
		"cert_pem": string(ev.CertPEM),
	})
}

func (s *Server) handleRevoke(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var in struct {
		Reason string `json:"reason"`
	}
	_ = readJSON(r, &in)
	if err := s.vault.Revoke(id, in.Reason); err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "revoked"})
}

type scanBody struct {
	WithinDays int `json:"within_days"`
}

func (s *Server) handleScan(w http.ResponseWriter, r *http.Request) {
	var in scanBody
	if err := readJSON(r, &in); err != nil {
		in.WithinDays = 30
	}
	if in.WithinDays <= 0 {
		in.WithinDays = 30
	}
	hits, err := s.vault.ScanExpiring(time.Duration(in.WithinDays) * 24 * time.Hour)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"hits": hits})
}

type previewBody struct {
	OldID   string `json:"old_id"`
	CertPEM string `json:"cert_pem"`
	KeyPEM  string `json:"key_pem"`
}

func (s *Server) handleRotatePreview(w http.ResponseWriter, r *http.Request) {
	var in previewBody
	if err := readJSON(r, &in); err != nil {
		writeErr(w, err)
		return
	}
	p, err := s.vault.PreviewRotate(in.OldID, []byte(in.CertPEM), []byte(in.KeyPEM))
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (s *Server) handleBundle(w http.ResponseWriter, _ *http.Request) {
	b, err := s.vault.ExportTrustBundle()
	if err != nil {
		writeErr(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/x-pem-file")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
}
