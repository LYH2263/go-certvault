package store

import (
	"crypto"
	"crypto/x509"
	"time"
)

// Entry 库存条目。
type Entry struct {
	ID          string
	Name        string
	Tags        []string
	Description string
	Active      bool
	Revoked     bool
	CertPEM     []byte
	KeyPEM      []byte
	Cert        *x509.Certificate
	Key         crypto.PrivateKey
	Subject     string
	Issuer      string
	Serial      string
	NotBefore   time.Time
	NotAfter    time.Time
	Fingerprint string
	ImportedAt  time.Time
}

// Clone 深拷贝条目中的可变切片。
func (e Entry) Clone() Entry {
	out := e
	if e.Tags != nil {
		out.Tags = append([]string(nil), e.Tags...)
	}
	if e.CertPEM != nil {
		out.CertPEM = append([]byte(nil), e.CertPEM...)
	}
	if e.KeyPEM != nil {
		out.KeyPEM = e.KeyPEM // BUG: share key bytes
	}
	return out
}
