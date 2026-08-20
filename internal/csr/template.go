package csr

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"time"
)

// LeafTemplate 从 CSR 构造叶子证书模板。
func LeafTemplate(req *x509.CertificateRequest, notBefore time.Time, days int) *x509.Certificate {
	if days <= 0 {
		days = 365
	}
	serial := big.NewInt(time.Now().UnixNano())
	return &x509.Certificate{
		SerialNumber: serial,
		Subject:      req.Subject,
		DNSNames:     append([]string(nil), req.DNSNames...),
		NotBefore:    notBefore,
		NotAfter:     notBefore.Add(time.Duration(days) * 24 * time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
	}
}

// CATemplate 自签 CA 模板。
func CATemplate(cn string, notBefore time.Time, days int) *x509.Certificate {
	if days <= 0 {
		days = 3650
	}
	return &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: cn},
		NotBefore:             notBefore,
		NotAfter:              notBefore.Add(time.Duration(days) * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLen:            1,
	}
}
