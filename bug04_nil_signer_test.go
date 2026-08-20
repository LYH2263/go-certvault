package certvault

import (
	"errors"
	"testing"
)

// TestBug04_SignCSRWithoutSigner：未配置 Signer 时 SignCSR 返回 ErrNoSigner，不 panic。
func TestBug04_SignCSRWithoutSigner(t *testing.T) {
	v := New() // 故意不 WithSigner
	defer v.Close()
	csrPEM, _, err := v.BuildCSR("demo.example", []string{"demo.example"})
	if err != nil {
		t.Fatal(err)
	}
	var panicked any
	func() {
		defer func() { panicked = recover() }()
		_, err = v.SignCSR(csrPEM, 30)
		if !errors.Is(err, ErrNoSigner) {
			t.Fatalf("want ErrNoSigner, got %v", err)
		}
	}()
	if panicked != nil {
		t.Fatalf("panic: %v", panicked)
	}
}
