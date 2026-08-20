package certvault

import (
	"errors"
	"testing"

	"example.com/certvault/internal/pemutil"
)

// TestBug05_ParseErrorIsSentinel：非法 PEM 应能 errors.Is(..., pemutil.ErrParse)。
func TestBug05_ParseErrorIsSentinel(t *testing.T) {
	v := New()
	defer v.Close()
	_, err := v.ImportPEM([]byte("not-a-pem"), []byte("also-bad"), Meta{Name: "bad", Active: true})
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, pemutil.ErrParse) {
		t.Fatalf("errors.Is ErrParse failed: %v", err)
	}
}
