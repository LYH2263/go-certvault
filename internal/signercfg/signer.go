package signercfg

import (
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"errors"
	"time"

	"example.com/certvault/internal/csr"
	"example.com/certvault/internal/validate"
)

// Signer 签发接口。
type Signer interface {
	Sign(req *x509.CertificateRequest, days int, now time.Time) ([]byte, error)
}

// SelfSigner 自签叶子（用内嵌 CA 或自签自身）。
type SelfSigner struct {
	CACert *x509.Certificate
	CAKey  crypto.Signer
}

func (s *SelfSigner) Sign(req *x509.CertificateRequest, days int, now time.Time) ([]byte, error) {
	if s == nil || s.CAKey == nil {
		return nil, errors.New("signercfg: nil signer")
	}
	days = validate.ClampDays(days)
	tmpl := csr.LeafTemplate(req, now, days)
	parent := s.CACert
	if parent == nil {
		parent = tmpl
	}
	return x509.CreateCertificate(rand.Reader, tmpl, parent, req.PublicKey, s.CAKey)
}

// NewEphemeralCA 生成临时自签 CA。
func NewEphemeralCA(cn string, now time.Time) (*SelfSigner, error) {
	key, err := generateEC()
	if err != nil {
		return nil, err
	}
	tmpl := csr.CATemplate(cn, now, 3650)
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		return nil, err
	}
	cert, err := x509.ParseCertificate(der)
	if err != nil {
		return nil, err
	}
	return &SelfSigner{CACert: cert, CAKey: key}, nil
}
