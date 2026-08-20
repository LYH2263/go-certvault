package certvault

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestBug06_RotatePersistFailureRollsBack：持久化失败时活动证书不得切换。
func TestBug06_RotatePersistFailureRollsBack(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")
	now := time.Now()
	cert, key := mustCertKey(t, "rot1", now.Add(-time.Hour), now.Add(24*time.Hour))
	v := New(WithPersistPath(path))
	id, err := v.ImportPEM(cert, key, Meta{Name: "svc", Active: true})
	if err != nil {
		t.Fatal(err)
	}
	// 父路径是普通文件，CreateTemp/MkdirAll 必失败（跨平台）
	blockFile := filepath.Join(dir, "notadir")
	if err := os.WriteFile(blockFile, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	badPath := filepath.Join(blockFile, "snap.json")
	v.mu.Lock()
	v.persistPath = badPath
	v.mu.Unlock()

	ncert, nkey := mustCertKey(t, "rot1", now.Add(-time.Hour), now.Add(48*time.Hour))
	_, err = v.Rotate(id, ncert, nkey)
	if err == nil {
		t.Fatal("expected persist error")
	}
	active, err := v.GetActive("svc")
	if err != nil {
		t.Fatal(err)
	}
	if active.ID != id {
		t.Fatalf("active switched despite persist fail: old=%s new=%s", id, active.ID)
	}
}
