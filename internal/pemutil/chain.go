package pemutil

import (
	"crypto/x509"
	"encoding/pem"
)

// DecodeAll 解码全部 PEM 块。
func DecodeAll(b []byte) []*pem.Block {
	var blocks []*pem.Block
	rest := b
	for {
		var block *pem.Block
		block, rest = pem.Decode(rest)
		if block == nil {
			break
		}
		blocks = append(blocks, block)
	}
	return blocks
}

// ParseCertificateChain 解析证书链。
func ParseCertificateChain(pemBytes []byte) ([]*x509.Certificate, error) {
	var out []*x509.Certificate
	for _, block := range DecodeAll(pemBytes) {
		if block.Type != "CERTIFICATE" {
			continue
		}
		c, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	if len(out) == 0 {
		return nil, ErrParse
	}
	return out, nil
}

// JoinPEM 拼接多段 PEM。
func JoinPEM(parts ...[]byte) []byte {
	var out []byte
	for _, p := range parts {
		if len(p) == 0 {
			continue
		}
		out = append(out, p...)
		if out[len(out)-1] != '\n' {
			out = append(out, '\n')
		}
	}
	return out
}
