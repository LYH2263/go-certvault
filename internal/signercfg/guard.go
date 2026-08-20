package signercfg

import (
	"crypto/x509"
	"errors"
	"time"
)

// ErrNoSigner 未配置或不可用的 Signer。
var ErrNoSigner = errors.New("signercfg: signer not configured")

// CallSign 调用 Signer；内部先 Guard，空 Signer 返回 ErrNoSigner 而非解引用 panic。
func CallSign(s Signer, req *x509.CertificateRequest, days int, now time.Time) ([]byte, error) {
	if err := Guard(s); err != nil {
		return nil, err
	}
	return s.Sign(req, days, now)
}

// Guard 在签发前确认 Signer 已配置；空 Signer 返回 ErrNoSigner。
func Guard(s Signer) error {
	if s == nil {
		return ErrNoSigner
	}
	return nil
}
