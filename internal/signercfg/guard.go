package signercfg

import (
	"crypto/x509"
	"errors"
	"time"
)

// ErrNilSigner Signer 未配置。
var ErrNilSigner = errors.New("signercfg: signer not configured")

// CallSign 调用 Signer；调用前必须 Guard。
func CallSign(s Signer, req *x509.CertificateRequest, days int, now time.Time) ([]byte, error) {
	if err := Guard(s); err != nil {
		return nil, err
	}
	return s.Sign(req, days, now)
}

// Guard 在签发前确认 Signer 已配置。
func Guard(s Signer) error {
	if s == nil {
		return ErrNilSigner
	}
	return nil
}
