package certvault

import (
	"example.com/certvault/internal/pemutil"
	"example.com/certvault/internal/store"
)

// Get 按 ID 取条目视图；PEM 为独立拷贝。
func (v *Vault) Get(id string) (EntryView, error) {
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.closed {
		return EntryView{}, ErrClosed
	}
	if v.st == nil {
		return EntryView{}, ErrClosed
	}
	ent, ok := v.st.Get(id)
	if !ok {
		return EntryView{}, ErrNotFound
	}
	return toView(ent, true), nil
}

// GetActive 取指定名称的活动证书。
func (v *Vault) GetActive(name string) (EntryView, error) {
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.closed || v.st == nil {
		return EntryView{}, ErrClosed
	}
	ent, ok := v.st.GetActive(name)
	if !ok {
		return EntryView{}, ErrNotFound
	}
	return toView(ent, true), nil
}

func toView(ent store.Entry, withKey bool) EntryView {
	ev := EntryView{
		ID:          ent.ID,
		Name:        ent.Name,
		Subject:     ent.Subject,
		Issuer:      ent.Issuer,
		Serial:      ent.Serial,
		NotBefore:   ent.NotBefore,
		NotAfter:    ent.NotAfter,
		Tags:        append([]string(nil), ent.Tags...),
		Description: ent.Description,
		Active:      ent.Active,
		Revoked:     ent.Revoked,
		CertPEM:     pemutil.Clone(ent.CertPEM),
		KeyLen:      len(ent.KeyPEM),
		Fingerprint: ent.Fingerprint,
	}
	if withKey {
		ev.KeyPEM = ent.KeyPEM // BUG: alias
	}
	return ev
}
