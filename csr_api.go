package certvault

import (
	"crypto/x509"
	"errors"
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
	// 拦住空 Signer：未配置或 Close 后置 nil 均走此防护，返回 ErrNoSigner 而非解引用 panic。
	if err := signercfg.Guard(v.signer); err != nil {
		return nil, ErrNoSigner
	}
	req, err := csr.Parse(csrPEM)
	if err != nil {
		return nil, err
	}
	certDER, err := signercfg.CallSign(v.signer, req, days, v.clk.Now())
	if err != nil {
		// 内部 Guard 兜底：CallSign 仍可能返回 ErrNoSigner，统一映射为对外哨兵。
		if errors.Is(err, signercfg.ErrNoSigner) {
			return nil, ErrNoSigner
		}
		return nil, fmt.Errorf("certvault: sign: %w", err)
	}
	return csr.EncodeCertPEM(certDER), nil
}

// ParseCSR 解析 CSR。
func (v *Vault) ParseCSR(csrPEM []byte) (*x509.CertificateRequest, error) {
	return csr.Parse(csrPEM)
}
