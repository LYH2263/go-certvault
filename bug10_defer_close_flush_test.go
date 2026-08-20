package certvault

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestBug10_CloseFlushesBeforeDropStore：Close 须先 Flush 再标记关闭，快照应含已导入条目。
func TestBug10_CloseFlushesBeforeDropStore(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")
	now := time.Now()
	cert, key := mustCertKey(t, "flush", now.Add(-time.Hour), now.Add(24*time.Hour))
	v := New(WithPersistPath(path))
	if _, err := v.ImportPEM(cert, key, Meta{Name: "flush", Active: true}); err != nil {
		t.Fatal(err)
	}
	if err := v.Close(); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(b) < 10 {
		t.Fatalf("snapshot too small: %q", b)
	}
	v2 := New(WithPersistPath(path))
	if err := v2.LoadPersist(); err != nil {
		t.Fatal(err)
	}
	if v2.StoreLen() != 1 {
		t.Fatalf("loaded entries=%d", v2.StoreLen())
	}
}
