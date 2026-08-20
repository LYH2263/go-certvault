package certvault

import (
	"crypto/x509"
	"fmt"

	"example.com/certvault/internal/csr"
	"example.com/certvault/internal/signercfg"
)

// BuildCSR 生成 CSR PEM（不要求 Signer）。
func (v *Vault) BuildCSR(commonName string, dns []string) (csrPEM, keyPEM []byte, err error) {
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.closed {
		return nil, nil, ErrClosed
	}
	return csr.Build(commonName, dns)
}

// SignCSR 使用已配置 Signer 签发；未配置返回 ErrNoSigner。
func (v *Vault) SignCSR(csrPEM []byte, days int) ([]byte, error) {
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.closed {
		return nil, ErrClosed
	}
	// BUG: 未检查 signer==nil，且未走 signercfg.Guard
	req, err := csr.Parse(csrPEM)
	if err != nil {
		return nil, err
	}
	certDER, err := signercfg.CallSign(v.signer, req, days, v.clk.Now())
	if err != nil {
		return nil, fmt.Errorf("certvault: sign: %w", err)
	}
	return csr.EncodeCertPEM(certDER), nil
}

// ParseCSR 解析 CSR。
func (v *Vault) ParseCSR(csrPEM []byte) (*x509.CertificateRequest, error) {
	return csr.Parse(csrPEM)
}
