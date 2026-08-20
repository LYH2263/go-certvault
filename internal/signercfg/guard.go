package signercfg

import (
	"crypto/x509"
	"time"
)

// CallSign 调用 Signer；调用前必须 Guard。
func CallSign(s Signer, req *x509.CertificateRequest, days int, now time.Time) ([]byte, error) {
	// BUG: 内部调用前不防 nil，直接解引用
	return s.Sign(req, days, now)
}

// Guard 在签发前确认 Signer 已配置。
func Guard(s Signer) error {
	// BUG: 空实现，不拦截 nil
	_ = s
	return nil
}
