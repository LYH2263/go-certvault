package certvault

import (
	"context"
	"errors"
	"testing"
	"time"
)

// TestBug07_ScanExpiringHonorsCancel：已取消的 ctx 不得扫完全部条目。
func TestBug07_ScanExpiringHonorsCancel(t *testing.T) {
	now := time.Now()
	v := New(WithScanStep(30 * time.Millisecond))
	defer v.Close()
	for i := 0; i < 8; i++ {
		cert, key := mustCertKey(t, "s", now.Add(-time.Hour), now.Add(12*time.Hour))
		name := "n" + string(rune('a'+i))
		if _, err := v.ImportPEM(cert, key, Meta{Name: name, Active: true}); err != nil {
			t.Fatal(err)
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := v.ScanExpiringContext(ctx, 24*time.Hour)
	if !errors.Is(err, ErrCanceled) && !errors.Is(err, context.Canceled) {
		t.Fatalf("want canceled, got %v", err)
	}
}
