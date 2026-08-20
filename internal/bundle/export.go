package bundle

import (
	"bytes"

	"example.com/certvault/internal/store"
)

// Export 拼接活动非吊销证书 PEM。
func Export(ents []store.Entry) []byte {
	var buf bytes.Buffer
	for _, e := range ents {
		if e.Revoked || !e.Active {
			continue
		}
		if len(e.CertPEM) == 0 {
			continue
		}
		buf.Write(e.CertPEM)
		if e.CertPEM[len(e.CertPEM)-1] != '\n' {
			buf.WriteByte('\n')
		}
	}
	return buf.Bytes()
}

// Count 可导出条目数。
func Count(ents []store.Entry) int {
	n := 0
	for _, e := range ents {
		if !e.Revoked && e.Active && len(e.CertPEM) > 0 {
			n++
		}
	}
	return n
}
