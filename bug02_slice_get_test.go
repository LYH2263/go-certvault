package certvault

import (
	"bytes"
	"testing"
	"time"
)

// TestBug02_GetListKeyPEMAlias：Get/List 返回的私钥不得与库存共享底层数组。
func TestBug02_GetListKeyPEMAlias(t *testing.T) {
	now := time.Now()
	cert, key := mustCertKey(t, "alias2", now.Add(-time.Hour), now.Add(24*time.Hour))
	v := New()
	defer v.Close()
	id, err := v.ImportPEM(cert, key, Meta{Name: "web", Active: true})
	if err != nil {
		t.Fatal(err)
	}
	got, err := v.Get(id)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.KeyPEM) == 0 {
		t.Fatal("empty key")
	}
	orig0 := got.KeyPEM[0]
	got.KeyPEM[0] ^= 0x55
	again, err := v.Get(id)
	if err != nil {
		t.Fatal(err)
	}
	if again.KeyPEM[0] != orig0 {
		t.Fatal("Get key PEM shared with inventory")
	}
	list := v.List()
	if len(list) != 1 {
		t.Fatalf("list len %d", len(list))
	}
	list[0].KeyPEM[0] ^= 0x22
	third, _ := v.Get(id)
	if !bytes.Equal(third.KeyPEM, again.KeyPEM) {
		t.Fatal("List key PEM mutation leaked into vault")
	}
}
