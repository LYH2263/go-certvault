package certvault

import (
	"testing"
	"time"

	"example.com/certvault/internal/signercfg"
)

func TestSmokeImportListScan(t *testing.T) {
	now := time.Now()
	signer, err := signercfg.NewEphemeralCA("test-ca", now)
	if err != nil {
		t.Fatal(err)
	}
	v := New(WithSigner(signer))
	defer v.Close()
	cert, key := mustCertKey(t, "smoke", now.Add(-time.Hour), now.Add(48*time.Hour))
	id, err := v.ImportPEM(cert, key, Meta{Name: "smoke", Active: true})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := v.Get(id); err != nil {
		t.Fatal(err)
	}
	if len(v.List()) != 1 {
		t.Fatal("list")
	}
	hits, err := v.ScanExpiring(72 * time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	if len(hits) != 1 {
		t.Fatalf("hits %d", len(hits))
	}
	csrPEM, _, err := v.BuildCSR("x.example", nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := v.SignCSR(csrPEM, 10); err != nil {
		t.Fatal(err)
	}
}
