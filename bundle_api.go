package certvault

import "example.com/certvault/internal/bundle"

// ExportTrustBundle 导出活动非吊销证书的 trust bundle。
func (v *Vault) ExportTrustBundle() ([]byte, error) {
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.closed || v.st == nil {
		return nil, ErrClosed
	}
	return bundle.Export(v.st.List()), nil
}
