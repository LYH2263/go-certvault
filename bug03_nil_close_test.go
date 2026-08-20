package certvault

import (
	"errors"
	"testing"
	"time"
)

// TestBug03_ImportAfterCloseNoPanic：Close 后 Import 应返回 ErrClosed，不得解引用 nil store。
func TestBug03_ImportAfterCloseNoPanic(t *testing.T) {
	now := time.Now()
	cert, key := mustCertKey(t, "closed", now.Add(-time.Hour), now.Add(24*time.Hour))
	v := New()
	if err := v.Close(); err != nil {
		t.Fatal(err)
	}
	var panicked any
	func() {
		defer func() { panicked = recover() }()
		_, err := v.ImportPEM(cert, key, Meta{Name: "x", Active: true})
		if !errors.Is(err, ErrClosed) {
			t.Fatalf("want ErrClosed, got %v", err)
		}
	}()
	if panicked != nil {
		t.Fatalf("panic after close: %v", panicked)
	}
}
