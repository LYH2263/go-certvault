package certvault

import (
	"bytes"
	"testing"
	"time"
)

// TestBug01_ImportPEMBufferAlias：Import 后调用方改写 PEM 缓冲不应污染库存。
func TestBug01_ImportPEMBufferAlias(t *testing.T) {
	now := time.Now()
	cert, key := mustCertKey(t, "alias1", now.Add(-time.Hour), now.Add(24*time.Hour))
	v := New()
	defer v.Close()
	id, err := v.ImportPEM(cert, key, Meta{Name: "api", Active: true})
	if err != nil {
		t.Fatal(err)
	}
	cert[0] ^= 0xff
	got, err := v.Get(id)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Equal(got.CertPEM, cert) {
		t.Fatal("inventory cert PEM aliased caller buffer")
	}
	if got.CertPEM[0] == cert[0] {
		t.Fatal("caller mutation visible in vault")
	}
}
