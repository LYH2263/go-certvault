package certvault

import (
	"fmt"

	"example.com/certvault/internal/pemutil"
	"example.com/certvault/internal/store"
	"example.com/certvault/internal/validate"
)

// ImportPEM 解析并入库证书与私钥 PEM。库存侧持有独立拷贝。
func (v *Vault) ImportPEM(certPEM, keyPEM []byte, meta Meta) (string, error) {
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.closed {
		return "", ErrClosed
	}
	if v.st == nil {
		return "", ErrClosed
	}
	if err := validate.MetaName(meta.Name); err != nil {
		return "", fmt.Errorf("%w: %v", ErrBadRequest, err)
	}
	parsed, err := pemutil.ParsePair(certPEM, keyPEM)
	if err != nil {
		return "", fmt.Errorf("import pem: %w", err)
	}
	ent := store.Entry{
		Name:        meta.Name,
		Tags:        append([]string(nil), meta.Tags...),
		Description: meta.Description,
		Active:      meta.Active || meta.Name != "",
		CertPEM:     pemutil.Clone(certPEM),
		KeyPEM:      pemutil.Clone(keyPEM),
		Cert:        parsed.Cert,
		Key:         parsed.Key,
		Subject:     parsed.Subject,
		Issuer:      parsed.Issuer,
		Serial:      parsed.Serial,
		NotBefore:   parsed.NotBefore,
		NotAfter:    parsed.NotAfter,
		Fingerprint: parsed.Fingerprint,
		ImportedAt:  v.clk.Now(),
	}
	id, err := v.st.Put(ent)
	if err != nil {
		return "", err
	}
	v.imports++
	if err := v.persistLocked(); err != nil {
		_ = v.st.Delete(id)
		v.imports--
		return "", err
	}
	return id, nil
}
