package csr

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
)

var ErrBadCSR = errors.New("csr: bad request")

// Parse 解析 CSR PEM。
func Parse(csrPEM []byte) (*x509.CertificateRequest, error) {
	block, _ := pem.Decode(csrPEM)
	if block == nil {
		return nil, fmt.Errorf("%w: no pem", ErrBadCSR)
	}
	if block.Type != "CERTIFICATE REQUEST" && block.Type != "NEW CERTIFICATE REQUEST" {
		return nil, fmt.Errorf("%w: type %q", ErrBadCSR, block.Type)
	}
	req, err := x509.ParseCertificateRequest(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrBadCSR, err)
	}
	if err := req.CheckSignature(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrBadCSR, err)
	}
	return req, nil
}

// EncodeCertPEM 编码证书 DER。
func EncodeCertPEM(der []byte) []byte {
	return pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
}
