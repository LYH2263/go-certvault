package csr

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"

	"example.com/certvault/internal/pemutil"
	"example.com/certvault/internal/validate"
)

// Build 生成 ECDSA P-256 CSR 与私钥 PEM。
func Build(commonName string, dns []string) (csrPEM, keyPEM []byte, err error) {
	if commonName == "" {
		return nil, nil, ErrBadCSR
	}
	var names []string
	for _, d := range dns {
		if validate.DNSName(d) {
			names = append(names, d)
		}
	}
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	tmpl := &x509.CertificateRequest{
		Subject:  pkix.Name{CommonName: commonName},
		DNSNames: names,
	}
	der, err := x509.CreateCertificateRequest(rand.Reader, tmpl, key)
	if err != nil {
		return nil, nil, err
	}
	csrPEM = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: der})
	keyPEM, err = pemutil.EncodePKCS8Key(key)
	return csrPEM, keyPEM, err
}
