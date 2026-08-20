package certvault

import (
	"context"
	"errors"
	"testing"
	"time"
)

// TestBug08_RotateContextHonorsCancel：轮换等待 I/O 时必须响应 ctx。
func TestBug08_RotateContextHonorsCancel(t *testing.T) {
	now := time.Now()
	cert, key := mustCertKey(t, "io", now.Add(-time.Hour), now.Add(24*time.Hour))
	v := New(WithRotateIODelay(400 * time.Millisecond))
	defer v.Close()
	id, err := v.ImportPEM(cert, key, Meta{Name: "io", Active: true})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Millisecond)
	defer cancel()
	ncert, nkey := mustCertKey(t, "io", now.Add(-time.Hour), now.Add(48*time.Hour))
	start := time.Now()
	_, err = v.RotateContext(ctx, id, ncert, nkey)
	elapsed := time.Since(start)
	if !errors.Is(err, ErrCanceled) && !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("want canceled/deadline, got %v", err)
	}
	if elapsed > 250*time.Millisecond {
		t.Fatalf("wait ignored cancel, elapsed=%s", elapsed)
	}
}
